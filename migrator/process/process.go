package process

import (
	"context"
	"fmt"

	"golang.org/x/time/rate"

	"github.com/mercadolibre/coreservices-team/migrator/process/internal/files"
	"github.com/mercadolibre/coreservices-team/migrator/tasks"
)

type Process interface {
	Run(fileName string) error
}

//Process is the interface for the Process
type process struct {
	task         tasks.Task
	rowsToInform int
	rateToStop   float32
	limiter      *rate.Limiter
}

//NewProcess returns a NewProcess instance
func NewProcess(task tasks.Task, rowsToInform int, rateToStop float32) Process {
	return &process{
		task:         task,
		rowsToInform: rowsToInform,
		rateToStop:   rateToStop,
	}
}

func NewProcessWithLimit(task tasks.Task, rowsToInform int, rateToStop float32, limiter *rate.Limiter) Process {
	return &process{
		task:         task,
		rowsToInform: rowsToInform,
		rateToStop:   rateToStop,
		limiter:      limiter,
	}
}

//Run runs the migrations process
func (p *process) Run(fileName string) error {
	f := files.File{}

	defer f.CloseFiles()

	var count int
	var countOk float32
	var countNok float32
	var countTotalNok float32

	f.Log("%v: starting migrator process", fileName)

	//Load tracks files data
	err := f.LoadTrackFilesData(fileName)
	if err != nil {
		f.Log("Error in LoadTrackFilesData - %v", err)
		return err
	}
	//Opens the input and trackfile
	err = f.OpenFiles(fileName)
	if err != nil {
		f.Log("Error in OpenFiles - %v", err)
		return err
	}

	//for each line of the open file
	for err == nil {
		//Get the next line without processing
		//it's possible that the first lines have
		// been executed in previous executions
		line, err := f.GetNextlineToProcess()
		if err != nil {
			f.Log("Error in GetNextlineToProcess - %v", err)
			return fmt.Errorf("Error reading lines %v", err)
		}
		//ends of the file
		if line == "" {
			f.Log("The migration process ends %v file with %v lines processed and %v lines with error", fileName, count, countTotalNok)
			break
		}

		//do the work
		if p.limiter != nil {
			if err := p.limiter.Wait(context.Background()); err != nil {
				f.Log("Error in Rate limiter Wait - %v", err)
				return fmt.Errorf("error waiting for rate limit token: %v", err)
			}
		}

		err = p.task.Do(line)
		count++
		if err == nil {
			countOk++
			err = f.SetOk(line)
		} else {
			countNok++
			countTotalNok++
			err = f.SetNok(line)
		}
		//error in trak file
		if err != nil {
			f.Log("Error writing track file - %v", err)
			return fmt.Errorf("Error writing track file - %v", err)
		}

		//informs if it corresponds
		if p.rowsToInform != 0 && count%p.rowsToInform == 0 {
			f.Log("%v: processed %v rows", fileName, count)
			countNok = 0
			countOk = 0
		}
		//stops if it corresponds
		if countOk+countNok > float32(p.rowsToInform)*0.1 && (countOk == 0 || countNok/countOk > p.rateToStop) {
			f.Log("%v: The error rate exceeded %v percent, the process will stops", fileName, p.rateToStop)
			return fmt.Errorf("%v: the error rate exceeded %v percent, the process stops", fileName, p.rateToStop)
		}
	}
	return err
}
