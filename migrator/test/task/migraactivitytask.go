package tasks

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mercadolibre/coreservices-team/migrator/tasks"
	"github.com/mercadolibre/go-meli-toolkit/gobigqueue"
)

var publisher gobigqueue.Publisher

//MigraTask is the task that migrate the activities
type MigraTask struct {
	tasks.Task
}

func init() {
	topicName := os.Getenv("TOPIC_NAME")
	cluster := os.Getenv("CLUSTER_NAME")

	if cluster == "" {
		cluster = "default"
	}

	publisher = gobigqueue.NewPublisher(
		cluster,
		[]string{
			topicName,
		})
}

//NewMigraTask returns a tasks for migrate activities
func NewMigraTask() *MigraTask {
	return &MigraTask{}
}

//Do is the especific task that migrate the activities
//sending a msg to bigQ
func (migraTask *MigraTask) Do(data interface{}) error {
	var msg map[string]interface{}

	line, ok := data.(string)
	if !ok {
		return fmt.Errorf("Error in cast data to string")
	}
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		return fmt.Errorf("Error in split line, detects %d fields and got 4", len(fields))
	}

	userId, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return err
	}
	//Para las activities de mango, se hace un tratamiento especial
	if fields[1] == "payment_v1_gateway" && fields[2] == "payment_v1" {
		msg = map[string]interface{}{
			"user_id": userId,
			"resource": map[string]interface{}{
				"type": fields[2],
			},
			"type": "gateway",
			"id":   fields[3],
		}
	} else {
		msg = map[string]interface{}{
			"user_id": userId,
			"resource": map[string]interface{}{
				"type": fields[1],
			},
			"type": fields[2],
			"id":   fields[3],
		}
	}
	return publisher.Send(&gobigqueue.Payload{msg, nil, nil})
}
