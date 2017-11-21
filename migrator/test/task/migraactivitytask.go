package task

import (
	"fmt"
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
	//os.Setenv("GO_ENVIRONMENT", "production")
	publisher = gobigqueue.NewPublisher(
		"default",
		[]string{
			"test_migracion_activities.mpcs_activities",
		})
}

//NewMigraTask returns a tasks for migrate activities
func NewMigraTask() *MigraTask {
	return &MigraTask{}
}

//Do is the especific task that migrate the activities
//sending a msg to bigQ
func (migraTask *MigraTask) Do(data interface{}) error {

	line, ok := data.(string)
	if !ok {
		return fmt.Errorf("Error in cast data to string")
	}
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		return fmt.Errorf("Error in split line, detects %d fields in the line", len(fields))
	}
	msg := map[string]interface{}{
		"user_id":       fields[0],
		"resource_type": fields[1],
		"activity_type": fields[2],
		"id":            fields[3],
	}
	return publisher.Send(&gobigqueue.Payload{msg, nil, nil})
}
