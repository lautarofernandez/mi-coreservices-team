package main

import (
	"fmt"
	"os"

	"github.com/mercadolibre/coreservices-team/migrator/process"
	"github.com/mercadolibre/coreservices-team/migrator/test/task"
	"time"
)

func main() {

	fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	fmt.Printf("args %v\n", os.Args)

	proc := process.NewProcess(task.NewMigraMovementsTask(), 100, 0.1)

	os.Chdir("..")

	err := proc.Run(os.Args[1])

	fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05-0700"))
	
	if err != nil {
		fmt.Printf("The process ends with error: %v", err)
		os.Exit(1)
	}
}
