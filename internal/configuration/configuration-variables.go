package configuration

// This file is used to load the configuration file variables inside the application to make them easier to reference

// Logs variables
var (
	LogsPath string
)

// Git variables
var (
	GitSync   bool
	GitUri    string
	GitBranch string
)

// Schedule variables
var (
	ScheduleDaysStart       string
	ScheduleDaysEnd         string
	ScheduleWorkdayEnabled  bool
	ScheduleWorkdayStart    string
	ScheduleWorkdayEnd      string
	ScheduleWorkdayTimezone string
)

// Misc variables
var (
	AllowedOutputFormats = []string{"json", "yaml", "text"}
)
