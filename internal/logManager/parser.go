package logManager

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mitchs-dev/library-go/customTime"
	"github.com/mitchs-dev/library-go/processor"
	log "github.com/sirupsen/logrus"
)

// This file is used to parse the log files

// GetLogFile opens the log file and returns the contents
func (l *LogFile) GetLogFile(logFilePath string) error {

	log.Debug("Opening log file: " + logFilePath)

	// Open the log file
	logFileData := processor.ReadFile(logFilePath)
	if len(logFileData) == 0 || logFileData == nil {
		return errors.New("log file (" + logFilePath + ") is empty")
	}

	log.Debug("Parsing log file: " + logFilePath)

	// Parse the log file
	err := json.Unmarshal(logFileData, &l)
	if err != nil {
		return errors.New("error parsing log file (" + logFilePath + "): " + err.Error())
	}

	log.Debug("Log file parsed: " + logFilePath)

	return nil
}

// ConvertTime converts epoch time to human readable time (I.e 1hr30min, 1day2hr, etc)
func ConvertTime(epochTime int64) string {
	// Format epoch time to string time
	strTime := time.Unix(epochTime, 0).String()

	// Parse the time
	parsedTime, err := customTime.ParseDuration(strTime)
	if err != nil {
		log.Fatal("Error parsing time: ", err)
	}

	// Format the time to string
	parsedTimeStr := parsedTime.String()

	return parsedTimeStr

}

// SaveLogFile saves the log file
func (l *LogFile) SaveLogFile(logFilePath string) error {

	log.Debug("Saving log file: " + logFilePath)

	// Marshal the log file
	logFileData, err := json.Marshal(l)
	if err != nil {
		return errors.New("error marshaling log file (" + logFilePath + "): " + err.Error())
	}

	// Backup the log file
	if !processor.CopyFile(logFilePath, logFilePath+".bak") {
		return errors.New("error backing up log file (" + logFilePath + ")")
	}

	// Delete the log file
	if !processor.FileDelete(logFilePath) {
		return errors.New("error deleting log file (" + logFilePath + ")")
	}

	// Save the log file
	if !processor.CreateFileAsByte(logFilePath, logFileData) {
		// Restore the log file
		if !processor.CopyFile(logFilePath+".bak", logFilePath) {
			return errors.New("error restoring log file (" + logFilePath + ")")
		}
		// Delete the backup log file
		if !processor.FileDelete(logFilePath + ".bak") {
			return errors.New("error deleting backup log file (" + logFilePath + ".bak)")
		}
		return errors.New("error saving log file (" + logFilePath + ")")
	}

	// Delete the backup log file
	if !processor.FileDelete(logFilePath + ".bak") {
		return errors.New("error deleting backup log file (" + logFilePath + ".bak)")
	}

	log.Debug("Log file saved: " + logFilePath)

	return nil
}
