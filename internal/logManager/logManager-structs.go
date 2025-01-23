package logManager

// This file holds the structures of the log files

// LogFile holds the structure of the log file
type LogFile struct {
	Log  map[string]map[int]string    `json:"Log"`
	Time map[string]map[int]TimeEntry `json:"time,omitempty"`
}

// TimeEntry holds the time entries of the logs
type TimeEntry struct {
	Start  int64 `json:"s,omitempty"`
	Pause  int64 `json:"p,omitempty"`
	Resume int64 `json:"r,omitempty"`
	End    int64 `json:"e,omitempty"`
	Total  int64 `json:"t,omitempty"`
}

// LogEntry represents a single entry in the log
type LogEntry struct {
	Status  string `yaml:"Status"`
	Time    string `yaml:"Time"`
	Message string `yaml:"Message"`
}

// LogFileEntries holds the structure of the log file entries
type LogFileEntries struct {
	Entries map[string]LogEntry `yaml:"Entries"`
}
