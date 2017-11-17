package process

import (
	"fmt"

	"github.com/mercadolibre/coreservices-migrator/src/files"
	"github.com/mercadolibre/coreservices-migrator/src/tasks"
)

//Process is the interface for the Process
type Process struct {
	task         tasks.Task
	rowsToInform int
	rateToStop   float32
}

//NewProcess returns a NewProcess instance
func NewProcess(task tasks.Task, rowsToInform int, rateToStop float32) *Process {
	return &Process{
		task:         task,
		rowsToInform: rowsToInform,
		rateToStop:   rateToStop,
	}
}

//Run runs the migrations process
func (process Process) Run(fileName string) error {

	defer files.CloseFiles()

	var count int
	var countOk float32
	var countNok float32

	err := files.OpenLogFile()
	if err != nil {
		return fmt.Errorf("Error opening log file migrator-progress.txt")
	}
	files.Log(" The migration process runs with %v file", fileName)

	//Load tracks files data
	err = files.LoadTrackFilesData(fileName)
	if err != nil {
		files.Log("Error in LoadTrackFilesData - %v", err)
		return err
	}
	//Opens the input and trackfile
	err = files.OpenFiles(fileName)
	if err != nil {
		files.Log("Error in OpenFiles - %v", err)
		return err
	}

	//for each line of the open file
	for err == nil {
		//Get the next line without processing
		//it's possible that the first lines have
		// been executed in previous executions
		line, err := files.GetNextlineToProcess()
		if err != nil {
			files.Log("Error in GetNextlineToProcess - %v", err)
			return fmt.Errorf("Error reading lines %v", err)
		}
		//ends of the file
		if line == "" {
			files.Log("The migration process ends %v file", fileName)
			break
		}
		//do the work
		err = process.task.Do(line)
		count++
		if err == nil {
			countOk++
			err = files.SetOk(line)
		} else {
			countNok++
			err = files.SetNok(line)
		}
		//error in trak file
		if err != nil {
			files.Log("Error writing track file - %v", err)
			return fmt.Errorf("Error writing track file - %v", err)
		}

		//informs if it corresponds
		if process.rowsToInform != 0 && count%process.rowsToInform == 0 {
			files.Log("Processed %v rows", count)
			countNok = 0
			countOk = 0
		}
		//stops if it corresponds
		if countOk+countNok > float32(process.rowsToInform)*0.1 && (countOk == 0 || countNok/countOk > process.rateToStop) {
			files.Log("The error rate exceeded %v percent, the process will stops", process.rateToStop)
			return fmt.Errorf("the error rate exceeded %v percent, the process stops", process.rateToStop)
		}
	}
	return err
}
