package components

import (
	"github.com/mercadolibre/go-meli-toolkit/gobigqueue"
)

//Bigqueue is a bigq wrapper
type Bigqueue struct{}

//NewBigqueue returns a bigq wrapper
func NewBigqueue() *Bigqueue {
	return &Bigqueue{}
}

//SendNotification send notification by bigq
func (sender *Bigqueue) SendNotification(topicName string, id string, resourceName string) error {

	publisher := gobigqueue.NewPublisher(
		"default",
		[]string{
			topicName,
		})

	msg := map[string]interface{}{
		"uid":      id,
		"resource": resourceName,
	}
	return publisher.Send(&gobigqueue.Payload{msg, nil, nil})
}
