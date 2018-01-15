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
	bigqCluster = "default"
	bigqTopic   = "mango-mpcs-movements-v1-migration"
)

type MovementMigrator struct {
	Publisher   gobigqueue.Publisher
	TimePerPush time.Duration
}

func NewMovementMigrator(rpm int64) *MovementMigrator {
	os.Setenv("GO_ENVIRONMENT", "production")

	timePerPush := time.Duration((60 * time.Second).Nanoseconds() / rpm)

	return &MovementMigrator{
		Publisher:   gobigqueue.NewPublisher(bigqCluster, []string{bigqTopic}),
		TimePerPush: timePerPush,
	}
}

func (m *MovementMigrator) Do(data interface{}) error {
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
		Msg: map[string]interface{}{"id": id},
	}

	err = m.Publisher.Send(payload)
	if err != nil {
		log.Printf("error sending id to bigqueue: %v", err)
	}
	return err
}
