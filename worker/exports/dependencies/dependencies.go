package dependencies

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/worker/exports/models"
	"github.com/mercadolibre/go-meli-toolkit/godsclient"
	"github.com/mercadolibre/go-meli-toolkit/gokvsclient"
	"github.com/mercadolibre/go-meli-toolkit/golockclient"
	"github.com/mercadolibre/go-meli-toolkit/goosclient"
)

// DsClient is used so that we can accept the minimum required functionality
// from the injected DS client.
type DsClient interface {
	SearchBuilder() godsclient.SearchBuilder
	ScrollBuilder() godsclient.ScrollBuilder
	CountBuilder() godsclient.CountBuilder
}

//LockClient is used to accept the minimum required lock functionality
type LockClient interface {
	Lock(resource string, ttl int) (golockclient.Lock, error)
	KeepAlive(lock golockclient.Lock) (golockclient.Lock, error)
	Unlock(lock golockclient.Lock) error
}

//IDFinderer is the interface used with diferents idfinders
type IDFinderer interface {
	GetID(c *gin.Context) (string, error)
}

//KvsClient is used to accept the minimum required to kvsclient functionality
type KvsClient interface {
	Save(gokvsclient.Item) error
	Update(gokvsclient.Item) error
	Get(key string) (gokvsclient.Item, error)
}

// StorageClient is used to accept the minimum required to object storage functionality
type StorageClient interface {
	Multipart(goosclient.MultipartInput) error
}

//SenderNotification is the interface used with diferents sender
type SenderNotification interface {
	SendNotification(exportItem *models.ExportItem) error
}

//ProcessExport is the interface used with diferents export process
type ProcessExport interface {
	Process(c *gin.Context, exportItem *models.ExportItem, lock golockclient.Lock) (*models.ExportItem, error)
}
