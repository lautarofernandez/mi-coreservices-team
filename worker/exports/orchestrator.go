package exports

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"github.com/mercadolibre/coreservices-team/libs/go/logger"
	"github.com/mercadolibre/coreservices-team/worker/exports/config"
	"github.com/mercadolibre/coreservices-team/worker/exports/dependencies"
	"github.com/mercadolibre/coreservices-team/worker/exports/models"
	"github.com/mercadolibre/coreservices-team/worker/exports/utils"
	"github.com/mercadolibre/go-meli-toolkit/godog"
	"github.com/mercadolibre/go-meli-toolkit/gokvsclient"
)

//Orchestrator contains the dependencies of the orchestrators
type Orchestrator struct {
	idFinder   dependencies.IDFinderer
	process    dependencies.ProcessExport
	lockClient dependencies.LockClient
	kvsClient  dependencies.KvsClient
	sender     dependencies.SenderNotification
}

//NewOrchestrator returns a new orchestrator
func NewOrchestrator(idFinder dependencies.IDFinderer, process dependencies.ProcessExport, kvsClient dependencies.KvsClient, lockClient dependencies.LockClient, sender dependencies.SenderNotification) *Orchestrator {

	return &Orchestrator{
		idFinder:   idFinder,
		lockClient: lockClient,
		kvsClient:  kvsClient,
		sender:     sender,
		process:    process,
	}
}

//Export main flow of export orchestrator
func (orchestrator *Orchestrator) Export(c *gin.Context) {
	defer utils.StartSegment(c, config.OrchestratorSeg).End()
	log := logger.LoggerWithName(c, "export orchestrator")
	//retrieve the id
	id, err := orchestrator.idFinder.GetID(c)
	if err != nil {
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_get_id").ToArray()...)
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("Error retrieving id %v", id),
		})
		return
	}
	log.Info(
		"Begin export and lock id",
		logger.Attrs{
			"id": id,
		})
	//retrieve a Lock
	lock, err := orchestrator.lockClient.Lock(id)
	if err != nil {
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_lock").ToArray()...)
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("error in lock id %v", id),
		})
		return
	}

	//always at the end, release the lock
	defer orchestrator.lockClient.Unlock(lock)
	log.Info(
		"Get kvs",
		logger.Attrs{
			"id": id,
		})

	//retireve the kvsItem
	exportItem, err := orchestrator.getKvsItem(id)
	if err != nil {
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_get_kvs").ToArray()...)
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: fmt.Sprintf("Error retriving id %v in kvs", id),
		})
		return
	}
	//Verified the status
	if exportItem.ExportStatus == models.ExportPending {
		log.Info(
			"Call to export",
			logger.Attrs{
				"id": id,
			})
		exportItem, err := orchestrator.process.Process(c, exportItem, lock)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("Error in Process Export id %v", id),
			})
			return
		}
		log.Info(
			"Save kvs status",
			logger.Attrs{
				"id": id,
			})
		err = orchestrator.saveKvsItem(exportItem)
		if err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_save_kvs").ToArray()...)
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("Error saving in kvs id", id),
			})
			return
		}
		godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "ok").ToArray()...)
	}
	if exportItem.ExportStatus == models.ExportDone && !exportItem.SendNotification && exportItem.NotifyCompletion.BqTopicName != "" {
		log.Info(
			"Send Notification",
			logger.Attrs{
				"id": id,
			})
		//Send Notification
		err := orchestrator.sender.SendNotification(exportItem.NotifyCompletion.BqTopicName, exportItem.ID, exportItem.ResourceName)
		if err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_notigication").ToArray()...)
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("Error sending the notification id %v", id),
			})
			return
		}
		//Save send notification in kvs
		log.Info(
			"Save send Notification",
			logger.Attrs{
				"id": id,
			})
		exportItem.SendNotification = true
		err = orchestrator.saveKvsItem(exportItem)
		if err != nil {
			godog.RecordSimpleMetric(config.GodogExecuteMetric, 1, new(godog.Tags).Add("results", "error_save_kvs").ToArray()...)
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: fmt.Sprintf("Error saving in kvs id %v", id),
			})
			return
		}
	}
	return
}

func (orchestrator *Orchestrator) getKvsItem(id string) (*models.ExportItem, error) {
	var exportItem models.ExportItem

	kvsItem, errKvs := orchestrator.kvsClient.Get(id)
	if errKvs != nil {
		return nil, fmt.Errorf("error in retrieve the export item")
	}
	if kvsItem == nil {
		return nil, fmt.Errorf("error export doesn't exists")
	}
	errKvs = kvsItem.GetValue(&exportItem)
	if errKvs != nil {
		return nil, fmt.Errorf("error retrieving item from kvsItem")
	}

	return &exportItem, nil
}

func (orchestrator *Orchestrator) saveKvsItem(exportItem *models.ExportItem) error {
	keysItem := gokvsclient.MakeItem(exportItem.ID, exportItem)
	errKvs := orchestrator.kvsClient.Save(keysItem)
	if errKvs != nil {
		return fmt.Errorf("Error saving in kvs")
	}
	return nil
}
