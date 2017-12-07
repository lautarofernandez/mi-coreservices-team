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
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: "Error retrieving id",
		})
		return
	}
	//retrieve a Lock
	lock, err := orchestrator.lockClient.Lock(id)
	if err != nil {
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: "error in lock",
		})
		return
	}

	//always at the end, release the lock
	defer orchestrator.lockClient.Unlock(lock)
	log.Info(
		"get_kvs_item",
		logger.Attrs{
			"id": id,
		})

	//retireve the kvsItem
	exportItem, err := orchestrator.getKvsItem(id)
	if err != nil {
		errors.ReturnError(c, &errors.Error{
			Cause:   err.Error(),
			Code:    errors.InternalServerApiError,
			Message: "Error retriving id in kvs",
		})
		return
	}
	//Verified the status
	if exportItem.ExportStatus == models.ExportPending {
		exportItem, err := orchestrator.process.Process(c, exportItem, lock)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: "Error in Process Export",
			})
			return
		}
		err = orchestrator.saveKvsItem(exportItem)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: "Error saving in kvs",
			})
			return
		}
	}
	if exportItem.ExportStatus == models.ExportDone && !exportItem.SendNotification && exportItem.NotifyCompletion.BqTopicName != "" {
		//Send Notification
		err := orchestrator.sender.SendNotification(exportItem.NotifyCompletion.BqTopicName, exportItem.ID, exportItem.ResourceName)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: "Error sending the notification",
			})
			return
		}
		//Save send notification in kvs
		exportItem.SendNotification = true
		err = orchestrator.saveKvsItem(exportItem)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Cause:   err.Error(),
				Code:    errors.InternalServerApiError,
				Message: "Error saving in kvs",
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
