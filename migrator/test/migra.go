package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/mercadolibre/coreservices-team/migrator/process"
	"github.com/mercadolibre/coreservices-team/migrator/test/task"
)

const (
	// DefaultWorkersCount sets the number of concurrent files to process
	DefaultWorkersCount = 10

	// Throughput is currently not used
	Throughput = 200000

	// DefaultMigrationDir is the directory where the migration CSV are found
	DefaultMigrationDir = "./exports"
)

func main() {
	if len(os.Args) != 3 {
		log.Printf("Usage: %s [concurrency] [csv directory]", os.Args[0])
		os.Exit(1)
	}

	workers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("could not parse concurrency count: %v", err)
	}

	directory := os.Args[2]
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		log.Fatalf("Path %v does not exists", directory)
	}

	log.Printf("Process started with %d workers and reading files from %s", workers, directory)
	log.Printf("Redirecting output to progress.log")

	f, err := os.OpenFile("./progress.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatalf("Error reading files from dir: %v", err)
	}

	wg := sync.WaitGroup{}

	// Parse pending files and send them to a channel
	pending := make(chan string, 5000)
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".csv") {
			continue
		}

		path, err := filepath.Abs(filepath.Join(directory, file.Name()))
		if err != nil {
			log.Fatalf("could not get absolute path for file %v: %v", file.Name(), err)
		}

		pending <- path
		wg.Add(1)
	}

	log.Printf("Read %d files from directory", len(pending))

	// Throughput is not used in migrator
	throughputPerWorker := int64(Throughput) / int64(workers)
	t := task.NewMovementMigrator(throughputPerWorker)

	os.Chdir(directory)

	// Start the corresponding number of workers
	log.Printf("Spawning workers...")
	for i := 0; i < workers; i++ {
		go func() {
			for file := range pending {
				log.Printf(`[START] Worker for file "%s" started.`, path.Base(file))

				proc := process.NewProcess(t, 10000, 0.1)
				if err := proc.Run(file); err != nil {
					log.Printf("migration for file %v ended with error: %v", path.Base(file), err)
				}

				log.Printf(`[FINISHED] Worker for file "%s" finished.`, path.Base(file))
				wg.Done()
			}
		}()
	}

	wg.Wait()
	log.Printf("All workers done!")
}
