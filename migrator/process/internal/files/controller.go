package files

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type File struct {
	inputFile        *os.File
	trackerFile      *os.File
	outputFile       *os.File
	logFile          *os.File
	readerInput      *bufio.Reader
	readerTracker    *bufio.Reader
	actualFileNumber int
	actualLineNumber int
	existsTrackFile  bool
	trackFiles       []trackFile
}

type trackFile struct {
	name     string
	nroLines int
}

//OpenFiles open the files of read,log and compare line file
func (f *File) OpenFiles(fileName string) error {
	var err error

	// For read access.
	f.inputFile, err = os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error opening file %v: %v", fileName, err)
	}
	f.readerInput = bufio.NewReader(f.inputFile)
	f.Log("Opening input file %s", fileName)

	i := len(f.trackFiles)
	if i > 0 {
		f.existsTrackFile = true
		f.actualFileNumber = i - 1
		nro := strconv.Itoa(i - 1)
		trackerfileName := fileName + "." + nro
		f.trackerFile, err = os.Open(trackerfileName)
		if err != nil {
			return fmt.Errorf("Error opening file %v: %v", fileName, err)
		}
		f.readerTracker = bufio.NewReader(f.trackerFile)
		f.Log("Opening track file %s", trackerfileName)
	}

	nro := strconv.Itoa(i)
	outputfileName := fileName + "." + nro
	f.outputFile, err = os.OpenFile(outputfileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Error open file %v", outputfileName)
	}
	f.Log("Opening output file %s", outputfileName)
	return nil
}

//GetNextlineToProcess returns the next line no processed
func (f *File) GetNextlineToProcess() (string, error) {
	var line string
	var err error
	var ok bool

	line, err = f.readerInput.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if !f.existsTrackFile {
		f.actualLineNumber++
	} else {
		//search the next line without processing
		for true {
			ok, err = f.verifyPreviousExecution(line)
			if err != nil {
				return "", err
			}
			if ok {
				f.SetOk(line)
			} else {
				break
			}
			f.actualLineNumber++
			line, err = f.readerInput.ReadString('\n')
			line = strings.TrimSpace(line)
		}
	}
	return line, nil
}

//verifyPreviousExecution verified if current line was processed
func (f *File) verifyPreviousExecution(line string) (bool, error) {
	var ok bool
	var err error
	previousline, err := f.readerTracker.ReadString('\n')
	if err != io.EOF && err != nil {
		return false, err
	}
	if previousline == "" {
		previousline, ok, err = f.getLineOfPreviousFile()
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	previousline = strings.TrimSpace(previousline)
	l := len(previousline)
	textline := previousline[0 : l-3]
	if textline != line {
		fmt.Printf("lines are deferents (%v)(%v)\n", textline, line)
		return false, nil
	}
	statusline := previousline[l-2 : l]
	if statusline == "OK" {
		return true, nil
	}
	return false, nil
}

func (f *File) getLineOfPreviousFile() (string, bool, error) {
	var openFile = false
	var line string
	var err error
	var ok bool

	if f.actualFileNumber > 0 {
		f.trackerFile.Close()
		for f.actualFileNumber > 0 {
			f.actualFileNumber--
			f.Log("Actual lines %v, analyze file %v with %v lines", f.actualLineNumber, f.trackFiles[f.actualFileNumber].name, f.trackFiles[f.actualFileNumber].nroLines)
			if f.trackFiles[f.actualFileNumber].nroLines > f.actualLineNumber {
				openFile = true
				f.trackerFile, err = os.Open(f.trackFiles[f.actualFileNumber].name)
				if err != nil {
					fmt.Printf("Error open file %v", f.trackFiles[f.actualFileNumber].name)
					return "", false, err
				}
				f.readerTracker = bufio.NewReader(f.trackerFile)
				for i := 0; i <= f.actualLineNumber && err == nil; i++ {
					line, err = f.readerTracker.ReadString('\n')
				}
				if err != io.EOF && err != nil {
					return "", false, err
				}
				ok = true
				break
			}
		}
		if f.actualFileNumber == 0 && openFile == false {
			f.existsTrackFile = false
			return "", false, nil
		}
	} else {
		f.existsTrackFile = false
		return "", false, nil
	}
	return line, ok, err
}

//SetOk set in log file that current line was processed ok
func (f *File) SetOk(line string) error {
	_, err := f.outputFile.WriteString(line + ",OK\n")
	return err
}

//SetNok set in log file that current line was processed wrong
func (f *File) SetNok(line string) error {
	_, err := f.outputFile.WriteString(line + ",ER\n")
	return err
}

//Log to logfile
func (f *File) Log(line string, parameters ...interface{}) error {
	logline := fmt.Sprintf(line, parameters...)

	log.Println(logline)

	return nil
}

//CloseFiles closeFiles
func (f *File) CloseFiles() {
	f.inputFile.Close()
	f.trackerFile.Close()
	f.outputFile.Close()
}

//LoadTrackFilesData iterate in logs files and load data about theses
func (f *File) LoadTrackFilesData(fileName string) error {
	var i = 0
	var err error
	origFilename := fileName
	nro := strconv.Itoa(i)
	fileName = origFilename + "." + nro
	f.trackFiles = make([]trackFile, 0)

	for true {
		f.trackerFile, err = os.Open(fileName)
		if err != nil {
			break
		}
		f.readerTracker = bufio.NewReader(f.trackerFile)
		l := 0
		for true {
			line, err := f.readerTracker.ReadString('\n')
			if err != io.EOF && err != nil {
				return err
			}
			if line == "" {
				break
			}
			l++
		}
		f.trackerFile.Close()
		tmpLogFile := trackFile{
			name:     fileName,
			nroLines: l,
		}
		f.trackFiles = append(f.trackFiles, tmpLogFile)
		f.Log("Append %s trackfile with %d lines", fileName, l)
		i++
		nro = strconv.Itoa(i)
		fileName = origFilename + "." + nro
	}
	return nil
}
