package exports

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"github.com/mercadolibre/coreservices-team/libs/go/logger"
	"github.com/mercadolibre/coreservices-team/worker/exports/config"
	"github.com/mercadolibre/coreservices-team/worker/exports/dependencies"
	"github.com/mercadolibre/coreservices-team/worker/exports/models"
	"github.com/mercadolibre/coreservices-team/worker/exports/utils"
	"github.com/mercadolibre/go-meli-toolkit/godog"
	"github.com/mercadolibre/go-meli-toolkit/gokvsclient"
	"github.com/mercadolibre/go-meli-toolkit/golockclient"
)

const (
	// ExportExpiryDelta is the delta used for expiring a message. The delta is between
	// now and the time the export was requested.
	ExportExpiryDelta = 12 * time.Hour
)

//Orchestrator contains the dependencies of the orchestrators
type Orchestrator struct {
	idFinder   dependencies.IDFinderer
	process    dependencies.ProcessExport
	kvsClient  dependencies.KvsClient
	sender     dependencies.SenderNotification
	lockClient dependencies.LockClient
	lockTTL    time.Duration
}

//NewOrchestrator returns a new orchestrator
func NewOrchestrator(idFinder dependencies.IDFinderer, process dependencies.ProcessExport, kvsClient dependencies.KvsClient, lockClient dependencies.LockClient, lockTTL time.Duration, sender dependencies.SenderNotification) *Orchestrator {
	return &Orchestrator{
		idFinder:   idFinder,
		lockClient: lockClient,
		lockTTL:    lockTTL,
		kvsClient:  kvsClient,
		sender:     sender,
		process:    process,
	}
}

//Export main flow of export orchestrator
func (o *Orchestrator) Export(c *gin.Context) {
	defer utils.StartSegment(c, config.OrchestratorSeg).End()
	log := logger.LoggerWithName(c, "export orchestrator")

	// Using the bigq message as body, use the id finder to get the notified export ID.
	exportID, err := o.idFinder.GetID(c)
	if err != nil {
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_get_id")
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("error retrieving id %v", exportID),
		})
		return
	}

	log.Debug("begin_export", logger.Attrs{"export_id": exportID})

	// With the export ID we are going to try and lock the export process. If the lock fails to acquire
	// this would mean that there's another process already working this export. In this case we
	// return 422 Unprocessable Entity as the request was correct but we chose not to process.
	lock, err := o.lockClient.Lock(exportID, int(o.lockTTL.Seconds()))
	if err != nil {
		if err == golockclient.ErrLocked {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.UnprocessableEntityApiError,
				Message: fmt.Sprintf("export id %v is running", exportID),
			})
			return
		}

		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_lock")
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("error in lock id %v", exportID),
		})
		return
	}
	defer o.lockClient.Unlock(lock)

	log.Debug("kvs_get_export", logger.Attrs{"export_id": exportID})

	exportItem, err := o.getKvsItem(exportID)
	if err != nil {
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_get_kvs")
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("error getting export %s from KVS: %v", exportID, err),
		})
		return
	}

	// If we have a requested time inside the export item request, then we calculate a delta from
	// now. If the delta is higher than the defined expiry one, then we flag the export as
	// error. We continue this process so this export is notified as well.
	if t := exportItem.RequestedAt; t != nil && time.Since(*t) > ExportExpiryDelta {
		exportItem.ExportStatus = models.ExportDone
		exportItem.ErrorDescription = "export request expired due to long completion wait"

		log.Info("export_expired", logger.Attrs{"export_id": exportItem.ID, "requested_at": t.Format(time.RFC3339), "delta": time.Since(*t)})
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_export_expired")

		if err := o.saveKvsItem(exportItem); err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_save_kvs")
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("error saving id %s to kvs", exportID),
			})
			return
		}
	}

	// If export status is pending, and we where able to lock the export ID, then this export is fine
	// to start processing. We'll do just that, and delegate the export process to the exporter.
	if exportItem.ExportStatus == models.ExportPending {
		log.Debug("start_export_process", logger.Attrs{"export_id": exportID})

		exportItem, err := o.process.Process(c, exportItem, lock)
		if err != nil {
			log.Error("export_process_error", logger.Attrs{"export_id": exportID, "error": err.Error()})
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("error processing export %s", exportID),
			})
			return
		}

		if err = o.saveKvsItem(exportItem); err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_save_kvs")
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("error saving id %s to kvs", exportID),
			})
			return
		}

		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:ok")
	}

	// The export was executed in the step before, or the current notified export was done prior to this
	// call. Either way, we check if we need to notify it's result to the caller. If the notification
	// was already sent, or there's no notification data then we bail out as we have nothing to do.
	if exportItem.ExportStatus == models.ExportDone || exportItem.ExportStatus == models.ExportError {
		if exportItem.NotificationSent || exportItem.NotifyCompletion.BqTopicName == "" {
			return
		}

		log.Debug("export_notify_result", logger.Attrs{"export_id": exportID, "result": string(exportItem.ExportStatus)})

		if err := o.sender.SendNotification(exportItem); err != nil {
			// We failed notifying the export because the given topic does not exist. In this case we are going
			// to answer 200OK to bigq, because there's nothing else we can do.
			if strings.Contains(err.Error(), "topic not found") {
				log.Warning("export_topic_not_found", logger.Attrs{"export_id": exportID, "error": err.Error()})
				godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_topic_not_found")

				c.Status(http.StatusOK)
				return
			}

			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_notification")
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("error sending the notification id %v", exportID),
			})
			return
		}

		// Once the notification was successfully sent, we need to
		// flag it as so in the export data stored in KVS.
		exportItem.NotificationSent = true

		if err := o.saveKvsItem(exportItem); err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, "results:error_save_kvs")
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("error saving in kvs id %v", exportID),
			})
			return
		}
	}
}

func (o *Orchestrator) getKvsItem(id string) (*models.ExportItem, error) {
	kvsItem, err := o.kvsClient.Get(id)
	if err != nil {
		return nil, fmt.Errorf("error retrieving export item: %v", err)
	}

	if kvsItem == nil {
		return nil, fmt.Errorf("error retrieving export: %s does not exist", id)
	}

	var exportItem models.ExportItem
	if err := kvsItem.GetValue(&exportItem); err != nil {
		return nil, fmt.Errorf("error unmarshalling export %s from kvs: %v", id, err)
	}

	return &exportItem, nil
}

func (o *Orchestrator) saveKvsItem(exportItem *models.ExportItem) error {
	item := gokvsclient.MakeItem(exportItem.ID, exportItem)
	if err := o.kvsClient.Save(item); err != nil {
		return fmt.Errorf("error saving export to kvs: %v", err)
	}

	return nil
}
