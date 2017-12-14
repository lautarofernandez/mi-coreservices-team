package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mercadolibre/coreservices-team/migrator/process"
	"github.com/mercadolibre/coreservices-team/migrator/test/task"
)

const (
	// WorkersCount sets the number of concurrent files to process
	WorkersCount = 2

	// Throughput is currently not used
	Throughput = 100

	// MigrationDir is the directory where the migration CSV are found
	MigrationDir = "./exports"
)

func main() {
	if _, err := os.Stat(MigrationDir); os.IsNotExist(err) {
		log.Fatalf("Path %v does not exists", MigrationDir)
	}

	f, err := os.OpenFile("./progress.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	files, err := ioutil.ReadDir(MigrationDir)
	if err != nil {
		log.Fatalf("Error reading files from dir: %v", err)
	}

	wg := sync.WaitGroup{}

	// Parse pending files and send them to a channel
	pending := make(chan string, 1000)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}

		pending <- file.Name()
		wg.Add(1)
	}

	// Throughput is not used in migrator
	throughputPerWorker := int64(Throughput) / WorkersCount
	t := task.NewMovementMigrator(throughputPerWorker)

	os.Chdir(MigrationDir)

	// Start the corresponding number of workers
	for i := 0; i < WorkersCount; i++ {
		go func() {
			for file := range pending {
				proc := process.NewProcess(t, 2, 0.1)
				if err := proc.Run(file); err != nil {
					log.Printf("migration for file %v ended with error: %v", file, err)
				}

				wg.Done()
			}
		}()
	}

	wg.Wait()
}
