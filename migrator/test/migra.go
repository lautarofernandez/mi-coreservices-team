package main

import (
	"fmt"
	"os"

	"github.com/mercadolibre/coreservices-team/migrator/process"
	"github.com/mercadolibre/coreservices-team/migrator/test/task"
)

func main() {

	proc := process.NewProcess(tasks.NewMigraTask(), 100, 0.1)
	err := proc.Run(os.Args[1])
	if err != nil {
		fmt.Printf("The process ends with error: %v", err)
		os.Exit(1)
	}
}
