package task

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mercadolibre/go-meli-toolkit/gobigqueue"
)

const (
	paymentBigqCluster = "default2"
	paymentBigqTopic   = "reindex.mpcs-payments"
)

type PaymentMigrator struct {
	Publisher   gobigqueue.Publisher
	TimePerPush time.Duration
}

func NewPaymentMigrator(rpm int64) *PaymentMigrator {
	os.Setenv("GO_ENVIRONMENT", "production")

	timePerPush := time.Duration((60 * time.Second).Nanoseconds() / rpm)

	log.Printf("Posting messages to %s %s", paymentBigqCluster, paymentBigqTopic)

	return &PaymentMigrator{
		Publisher:   gobigqueue.NewPublisher(paymentBigqCluster, []string{paymentBigqTopic}),
		TimePerPush: timePerPush,
	}
}

func (m *PaymentMigrator) Do(data interface{}) error {
	line, ok := data.(string)
	if !ok {
		log.Printf("Error casting data to string: %v", data)
		return fmt.Errorf("Error casting data to string: %v", data)
	}

	id, err := strconv.ParseUint(line, 10, 64)
	if err != nil {
		log.Printf("Error casting data to uint: %v", err)
		return fmt.Errorf("Error parsing string into uint: %v", line)
	}

	payload := &gobigqueue.Payload{
		Msg: map[string]interface{}{"payment_id": line},
	}

	err = m.Publisher.Send(payload)
	if err != nil {
		log.Printf("error sending id %d to bigqueue: %v", id, err)
	}
	return err
}
