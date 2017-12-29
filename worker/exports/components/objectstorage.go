package components

import (
	"fmt"
	"net/http"

	"github.com/mercadolibre/go-meli-toolkit/restful/rest"
)

//ObjectStorage is a Objecto Storage wrapper
type ObjectStorage struct {
	writeClient rest.RequestBuilder
}

//NewObjectStorage returns a Objecto Storage wrapper
func NewObjectStorage(writeClient rest.RequestBuilder) *ObjectStorage {
	return &ObjectStorage{
		writeClient: writeClient,
	}
}

//PutFile put a file in fury Objecto Storage
func (obj *ObjectStorage) PutFile(path string, data []byte) error {

	response := obj.writeClient.Put(path, data)
	if response.Err != nil {
		return fmt.Errorf("error saving in object storage")
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status returned while perfoming PUT into storage: %d", response.StatusCode)
	}
	return nil
}
