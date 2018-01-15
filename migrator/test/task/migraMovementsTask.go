package task

import (
	"fmt"

	"os"
	"strconv"

	"github.com/mercadolibre/coreservices-team/migrator/tasks"
	"github.com/mercadolibre/go-meli-toolkit/gobigqueue"
)

var publisher gobigqueue.Publisher

//MigraTask is the task that migrate the activities
type MigraMovementsTask struct {
	tasks.Task
}

func init() {
	os.Setenv("GO_ENVIRONMENT", "production")
	publisher = gobigqueue.NewPublisher("default", []string{"mango-mpcs-movements-v1-migration"})
}

//NewMigraTask returns a tasks for migrate activities
func NewMigraMovementsTask() *MigraMovementsTask {
	return &MigraMovementsTask{}
}

//Do is the especific task that migrate the activities
//sending a msg to bigQ
func (migraTask *MigraMovementsTask) Do(data interface{}) error {

	line, ok := data.(string)
	if !ok {
		return fmt.Errorf("Error in cast data to string")
	}
	id, err := strconv.ParseUint(line, 10, 64)
	if err != nil {
		return fmt.Errorf("Error parsing id: " + line)
	}
	msg := map[string]interface{}{
		"id": id,
	}
	fmt.Println("Movement ID: " + line)
	return publisher.Send(&gobigqueue.Payload{msg, nil, nil})
}
