package files

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

var inputFile *os.File
var trackerFile *os.File
var outputFile *os.File
var logFile *os.File

var readerInput *bufio.Reader
var readerTracker *bufio.Reader
var actualFileNumber int
var actualLineNumber int

var existsTrackFile bool
var trackFiles []trackFile

type trackFile struct {
	name     string
	nroLines int
}

//OpenFiles open the files of read,log and compare line file
func OpenFiles(fileName string) error {
	var err error

	// For read access.
	inputFile, err = os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Error open file %v", fileName)
	}
	readerInput = bufio.NewReader(inputFile)
	Log("Opening input file %s", fileName)

	i := len(trackFiles)
	if i > 0 {
		existsTrackFile = true
		actualFileNumber = i - 1
		nro := strconv.Itoa(i - 1)
		trackerfileName := fileName + "." + nro
		trackerFile, err = os.Open(trackerfileName)
		if err != nil {
			return fmt.Errorf("Error open file %v", trackerfileName)
		}
		readerTracker = bufio.NewReader(trackerFile)
		Log("Opening track file %s", trackerfileName)
	}

	nro := strconv.Itoa(i)
	outputfileName := fileName + "." + nro
	outputFile, err = os.OpenFile(outputfileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Error open file %v", outputfileName)
	}
	Log("Opening output file %s", outputfileName)
	return nil
}

//GetNextlineToProcess returns the next line no processed
func GetNextlineToProcess() (string, error) {
	var line string
	var err error
	var ok bool

	line, err = readerInput.ReadString('\n')
	if err != io.EOF && err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if !existsTrackFile {
		actualLineNumber++
	} else {
		//search the next line without processing
		for true {
			ok, err = verifyPreviousExecution(line)
			if err != nil {
				return "", err
			}
			if ok {
				SetOk(line)
			} else {
				break
			}
			actualLineNumber++
			line, err = readerInput.ReadString('\n')
			line = strings.TrimSpace(line)
		}
	}
	return line, nil
}

//verifyPreviousExecution verified if current line was processed
func verifyPreviousExecution(line string) (bool, error) {
	var ok bool
	var err error
	previousline, err := readerTracker.ReadString('\n')
	if err != io.EOF && err != nil {
		return false, err
	}
	if previousline == "" {
		previousline, ok, err = getLineOfPreviousFile()
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

func getLineOfPreviousFile() (string, bool, error) {
	var openFile = false
	var line string
	var err error
	var ok bool

	if actualFileNumber > 0 {
		trackerFile.Close()
		for actualFileNumber > 0 {
			actualFileNumber--
			Log("Actual lines %v, analyze file %v with %v lines", actualLineNumber, trackFiles[actualFileNumber].name, trackFiles[actualFileNumber].nroLines)
			if trackFiles[actualFileNumber].nroLines > actualLineNumber {
				openFile = true
				trackerFile, err := os.Open(trackFiles[actualFileNumber].name)
				if err != nil {
					fmt.Printf("Error open file %v", trackFiles[actualFileNumber].name)
					return "", false, err
				}
				readerTracker = bufio.NewReader(trackerFile)
				for i := 0; i <= actualLineNumber && err == nil; i++ {
					line, err = readerTracker.ReadString('\n')
				}
				if err != io.EOF && err != nil {
					return "", false, err
				}
				ok = true
				break
			}
		}
		if actualFileNumber == 0 && openFile == false {
			existsTrackFile = false
			return "", false, nil
		}
	} else {
		existsTrackFile = false
		return "", false, nil
	}
	return line, ok, err
}

//SetOk set in log file that current line was processed ok
func SetOk(line string) error {
	_, err := outputFile.WriteString(line + ",OK\n")
	return err
}

//SetNok set in log file that current line was processed wrong
func SetNok(line string) error {
	_, err := outputFile.WriteString(line + ",ER\n")
	return err
}

//OpenLogFile open a log process file
func OpenLogFile() (err error) {
	logFile, err = os.OpenFile("migrator-progress.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("Error open log file")
	}
	return err
}

//Log to logfile
func Log(line string, parameters ...interface{}) error {
	logline := fmt.Sprintf(line, parameters...)
	logDate := fmt.Sprintf("%v - %s\n", time.Now().Format(time.RFC3339), logline)
	_, err := logFile.WriteString(logDate)
	logFile.Sync()
	return err
}

//CloseFiles closeFiles
func CloseFiles() {
	inputFile.Close()
	trackerFile.Close()
	outputFile.Close()
}

//LoadTrackFilesData iterate in logs files and load data about theses
func LoadTrackFilesData(fileName string) error {
	var i = 0
	var err error
	origFilename := fileName
	nro := strconv.Itoa(i)
	fileName = origFilename + "." + nro
	trackFiles = make([]trackFile, 0)

	for true {
		trackerFile, err = os.Open(fileName)
		if err != nil {
			break
		}
		readerTracker = bufio.NewReader(trackerFile)
		l := 0
		for true {
			line, err := readerTracker.ReadString('\n')
			if err != io.EOF && err != nil {
				return err
			}
			if line == "" {
				break
			}
			l++
		}
		trackerFile.Close()
		tmpLogFile := trackFile{
			name:     fileName,
			nroLines: l,
		}
		trackFiles = append(trackFiles, tmpLogFile)
		Log("Append %s trackfile with %d lines", fileName, l)
		i++
		nro = strconv.Itoa(i)
		fileName = origFilename + "." + nro
	}
	return nil
}
