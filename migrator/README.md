MercadoPago CoreServices Migrator
=

Skeleton library for migrations


Instalation
---

```
$ go get github.com/mercadolibre/coreservices-team/migrator
```

Usage
---

The library is responsible for iterating a text file and for each line a task is executed. The library supports re run saving the result of the execution of a line in a track file.  

Uso

To use this library it is necessary only define the task implementing the following interface:

```go
type Task interface {
	// Do is the function to do of the task
	 Do(data interface{}) error
}
```

Example
===

First it's necesary make the task.

```go
package mock

import (
	"fmt"
	"github.com/mercadolibre/coreservices-migrator/src/tasks"
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
	
```

And create the process intance with the folowing parameters: 

```go

	proc := process.NewProcess(mock.NewMockTask(), rowsToInform, rateToStop)
	err := proc.Run("migra.csv")
	if err == nil {
		fmt.Printf("The process ends without errors")
	}else {
		fmt.Printf("The process ends with error: %v\n", err)
		os.Exit(1)
	}

```

mock.NewMockTask(): is the task to execute.

rowsToInform: is the number of lines that will pass for the status to be reported.

rateToStop: is the rate of task executes with errors that stops the process.

If you have a process that reports the results every 1000 rows and a rate of 0.1, the process will stop after 100 failed executions. 



Changelog
---

0.0.1 - 2017-11-17 

- Initial commit. 

