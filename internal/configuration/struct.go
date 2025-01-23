package configuration

// This file holds the structure of the configuration file

// Generated using github.com/mitchs-dev/build-struct@v1.2.1
type Configuration struct {
	Settings struct {
		Schedule struct {
			Days struct {
				Start string `yaml:"start"`
				End   string `yaml:"end"`
			} `yaml:"days"`
			Workday struct {
				End      string `yaml:"end,omitempty"`
				Timezone string `yaml:"timezone,omitempty"`
				Enabled  bool   `yaml:"enabled"`
				Start    string `yaml:"start,omitempty"`
			} `yaml:"workday"`
		} `yaml:"schedule"`
		Logs struct {
			Path string `yaml:"path"`
		} `yaml:"logs"`
		Git struct {
			Sync   bool   `yaml:"sync"`
			Uri    string `yaml:"uri,omitempty"`
			Branch string `yaml:"branch,omitempty"`
		} `yaml:"git"`
	} `yaml:"settings"`
}
