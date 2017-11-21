package mock

import (
	"fmt"
	"github.com/mercadolibre/coreservices-team/migrator/tasks"
	"time"
)

//MockTask is a mock of a task
type MockTask struct {
	tasks.Task
}

//NewMockTask returns a mock
func NewMockTask() *MockTask {
	return &MockTask{}
}

//Do is a mock of a task
func (mockTask *MockTask) Do(data interface{}) error {
	fmt.Printf("Tarea %v\n", data)
	time.Sleep(100 * time.Millisecond)
	return nil
}
