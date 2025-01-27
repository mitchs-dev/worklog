package logManager

import (
	"encoding/json"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/mitchs-dev/library-go/processor"
	"github.com/mitchs-dev/worklog/internal/calendarManager"
)

// CreateYearIfNotExist creates a year and week if it does not exist
func CreateWeekIfNotExist(logFilePath string) {

	log.Debug("Using log file path: ", logFilePath)

	// Get the year from the end of the filename.Dir
	basePath := filepath.Dir(logFilePath)

	log.Debug("basePath: ", basePath)

	// Get the week from the end of the filename
	weekFile := filepath.Base(logFilePath)

	log.Debug("weekFile: ", weekFile)

	// Check if the year exists
	if !processor.DirectoryOrFileExists(basePath) {

		// Create the year
		if !processor.CreateDirectory(basePath) {
			log.Fatalf("Failed to create basePath %v", basePath)
		}
	} else {
		log.Debugf("basePath %v already exists", basePath)
	}

	// Check if the week exists
	if processor.DirectoryOrFileExists(logFilePath) {
		log.Debugf("logFilePath %v already exists", logFilePath)
		return
	} else {
		// Create the empty week
		emptyWeek := calendarManager.WeekTree{Weeks: make(map[int]calendarManager.MonthDayTree)}

		// Marshal the empty week
		emptyWeekData, err := json.Marshal(emptyWeek)
		if err != nil {
			log.Fatal("Error creating week (", logFilePath, "): ", err)
		}

		// Create the week
		if !processor.CreateFileAsByte(logFilePath, emptyWeekData) {
			log.Fatal("Error creating week (", logFilePath, ")")
		} else {
			log.Debugf("Created week %v", logFilePath)
		}
	}
}
