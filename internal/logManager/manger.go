package logManager

import (
	"fmt"
	"path/filepath"

	"github.com/mitchs-dev/library-go/generator"
	"github.com/mitchs-dev/library-go/processor"
	"github.com/mitchs-dev/worklog/internal/calendarManager"
	"github.com/mitchs-dev/worklog/internal/configuration"
	log "github.com/sirupsen/logrus"
)

var (
	allowedActions = []string{
		// Basic actions
		"add",
		"remove",
		"list",
		// Code actions (To be implemented)
		"edit",
		// Time actions (To be implemented)
		"start",
		"pause",
		"resume",
		"end",
	}
)

// validateAction checks if the action is allowed
func validateAction(action string) bool {
	for _, allowedAction := range allowedActions {
		if allowedAction == action {
			return true
		}
	}
	return false
}

// Action is the main function for the log manager
func Action(action, logMessage, logId string, period string) (LogFileEntries, []string) {

	if !validateAction(action) {
		log.Fatal("Invalid action: ", action)
	}

	switch action {
	case "add":
		return actionAdd(logMessage)
	case "remove":
		return actionRemove(logId)
	case "list":
		return actionList(period)
	case "edit":
		// To be implemented
	case "start":
		// To be implemented
	case "pause":
		// To be implemented
	case "resume":
		// To be implemented
	case "end":
		// To be implemented
	}

	return LogFileEntries{}, nil

}

// actionAdd adds a log entry
func actionAdd(logMessage string) (LogFileEntries, []string) {

	dirs, _, _, today, err := calendarManager.PeriodFetch("today")
	if err != nil {
		log.Fatal("Error fetching period: ", err)
	}

	if len(dirs) > 1 {
		log.Fatal("Multiple weeks returned but expected only one")
	}

	logFilePath := configuration.LogsPath + "/" + dirs[0]

	log.Debug("Using log file: ", logFilePath)

	logFileDir := filepath.Dir(logFilePath)

	log.Debug("Log file directory: ", logFileDir)

	// Create the log file directory if it doesn't exist
	log.Debug("Checking if log file directory exists: ", logFileDir)
	if !processor.DirectoryOrFileExists(logFileDir) {
		if !processor.CreateDirectory(logFileDir) {
			log.Fatal("Error creating log file directory: ", logFileDir)
		}
		log.Debug("Created log file directory: ", logFileDir)
	}

	var lf LogFile
	lf.GetLogFile(logFilePath)

	// Find the highest log id for the day
	var highestLogId int
	for logId := range lf.Log[today] {
		if logId > highestLogId {
			highestLogId = logId
		}
	}

	// Now we can set the new log id
	newLogId := highestLogId + 1

	// Add the log entry
	lf.Log[today][newLogId] = logMessage

	// We also need to set the time entry
	lf.Time[today][newLogId] = TimeEntry{
		Start:  int64(generator.EpochTimestamp(configuration.ScheduleWorkdayTimezone)),
		Pause:  0,
		Resume: 0,
		End:    0,
		Total:  0,
	}

	// Save the log file
	err = lf.SaveLogFile(logFilePath)
	if err != nil {
		log.Fatal("Error saving log file: ", err)
	}

	return LogFileEntries{
		Entries: map[string]LogEntry{
			today + "-" + fmt.Sprint(newLogId): {
				Status:  EntryStatusAdded,
				Time:    ConvertTime(lf.Time[today][newLogId].Total),
				Message: logMessage,
			},
		},
	}, []string{today + "-" + fmt.Sprint(newLogId)}

}

// actionRemove removes a log entry
func actionRemove(logId string) (LogFileEntries, []string) {

	return LogFileEntries{}, nil

}

func actionList(period string) (LogFileEntries, []string) {

	_, useYearTree, start, end, err := calendarManager.PeriodFetch(period)
	if err != nil {
		log.Fatal("Error fetching period: ", err)
	}

	log.Debug("Period: ", period)
	log.Debug("Start: ", start)
	log.Debug("End: ", end)

	// Initialize return values
	entries := LogFileEntries{
		Entries: make(map[string]LogEntry),
	}
	var entryIDs []string

	for year := range useYearTree.Years {
		log.Debug("Iterating year: ", year)
		for week := range useYearTree.Years[year].Weeks {
			log.Debug("Iterating week: ", week)
			weekStr := fmt.Sprint(week)
			if len(weekStr) == 1 {
				weekStr = "0" + fmt.Sprint(week)
			}
			fileName := configuration.LogsPath + "/" + fmt.Sprintf("%s/%s", fmt.Sprint(year), weekStr)
			log.Debug("Using log file: ", fileName)
			var lf LogFile
			lf.GetLogFile(fileName)
			log.Debug("Month days: ", useYearTree.Years[year].Weeks[week].MonthDays)
			for monthDayIndex, _ := range useYearTree.Years[year].Weeks[week].MonthDays {
				monthDay := useYearTree.Years[year].Weeks[week].MonthDays[monthDayIndex]
				monthDayStr := fmt.Sprint(monthDay)
				log.Debug("Iterating month day: ", monthDay)
				totalEntries := len(lf.Log[monthDayStr])
				log.Debug("Total entries: ", totalEntries)
				for logId := range lf.Log[monthDayStr] {
					log.Debug("Iterating log id: ", logId)
					var status string
					var timeStr string
					if lf.Time[monthDayStr][logId].End != 0 {
						status = EntryStatusCompleted
						timeStr = ConvertTime(lf.Time[monthDayStr][logId].Total)
					} else if lf.Time[monthDayStr][logId].Resume != 0 {
						status = EntryStatusResumed
						timeStr = ConvertTime(lf.Time[monthDayStr][logId].Resume - lf.Time[monthDayStr][logId].Start)
					} else if lf.Time[monthDayStr][logId].Pause != 0 {
						status = EntryStatusPaused
						timeStr = ConvertTime(lf.Time[monthDayStr][logId].Pause)
					} else {
						status = EntryStatusStarted
						timeStr = ConvertTime(lf.Time[monthDayStr][logId].Start)
					}
					entries.Entries[monthDayStr+"-"+fmt.Sprint(logId)] = LogEntry{
						Status:  status,
						Time:    timeStr,
						Message: lf.Log[monthDayStr][logId],
					}
					entryIDs = append(entryIDs, monthDayStr+"-"+fmt.Sprint(logId))
				}
			}
		}

	}

	return entries, entryIDs
}
