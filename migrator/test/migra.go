package main

import (
	"fmt"
	"os"

	"github.com/mercadolibre/coreservices-team/migrator/process"
	"github.com/mercadolibre/coreservices-team/migrator/test/mock"
)

func main() {

	fmt.Printf("args %v\n", os.Args)

	proc := process.NewProcess(mock.NewMockTask(), 100, 0.1)

	os.Chdir("..")

	err := proc.Run("migra.csv")

	if err != nil {
		fmt.Printf("The process ends with error: %v", err)
		os.Exit(1)
	}
}
