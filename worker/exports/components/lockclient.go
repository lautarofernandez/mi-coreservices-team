package components

import (
	"fmt"

	"github.com/mercadolibre/go-meli-toolkit/golockclient"
)

//LockClient is a fury lock wrapper
type LockClient struct {
	ttl          int
	golockclient golockclient.Client
}

//NewLock returns a new fury lock wrapper
func NewLock(golockclient golockclient.Client, ttl int) *LockClient {
	return &LockClient{
		golockclient: golockclient,
		ttl:          ttl,
	}
}

//Lock lock a resource
func (lockClient *LockClient) Lock(resource string) (interface{}, error) {
	loc, err := lockClient.golockclient.Lock(resource, lockClient.ttl)
	if err != nil {
		return nil, err
	}
	return &loc, nil
}

//KeepAlive maintains a resource lock
func (client *LockClient) KeepAlive(lockInterface interface{}) (interface{}, error) {
	lock, ok := lockInterface.(*golockclient.Lock)
	if !ok {
		return nil, fmt.Errorf("Error, can't cast interfaz to loc")
	}
	return client.golockclient.KeepAlive(*lock)
}

//Unlock unlock a resource
func (client *LockClient) Unlock(lockInterface interface{}) error {
	lock, ok := lockInterface.(*golockclient.Lock)
	if !ok {
		return fmt.Errorf("Error, can't cast interfaz to loc")
	}

	return client.golockclient.Unlock(*lock)
}
