// The configuration package is responsible for loading the configuration file.
package configuration

import (
	"embed"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	mConfiguration "github.com/mitchs-dev/library-go/configuration"
	"github.com/mitchs-dev/library-go/loggingFormatter"
	"github.com/mitchs-dev/library-go/processor"
	"gopkg.in/yaml.v2"
)

var (

	//go:embed default.yaml
	defaultConfigEmbed embed.FS

	// defaultConfigPath is the path to the default configuration file
	defaultConfigEmbedPath = "default.yaml"

	// Configuration is the configuration for the application
	configurationContext Configuration

	// ConfigurationPath is the path to the configuration file
	ConfigurationPath string

	// WorkLogHomeDir is the path to the worklog home directory
	WorkLogHomeDir = userHomeDir() + "/.worklog"

	// DefaultConfigurationPath is the path to the default configuration file
	DefaultConfigurationPath = WorkLogHomeDir + "/config"

	// EnableDebugMode is a flag to enable debug mode
	EnableDebugMode bool
)

// Since we need to set the variables before the init function is called
// We will create a function to initialize the configuration
// Instead of using the init function
func ConfigInit() {

	// Set the logging output, format, level
	log.SetOutput(os.Stdout)
	log.SetFormatter(&loggingFormatter.Formatter{})
	if EnableDebugMode {
		log.SetLevel(log.DebugLevel)
		log.Debug("Debug mode enabled")
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// Since we will need to unmarshal the default config either way
	// We will go ahead and load it here

	// Load the default configuration
	defaultConfigurationData, err := defaultConfigEmbed.ReadFile(defaultConfigEmbedPath)
	if err != nil {
		log.Fatal("Error loading default configuration: ", err)
	}

	// Check if the configuration path is provided before
	// attempting to load the configuration

	// If no config path provided, use the default
	if ConfigurationPath == "" {

		log.Debug("No configuration path provided, using default configuration path")

		ConfigurationPath = DefaultConfigurationPath

		// Check if the configuration file exists
		if !processor.DirectoryOrFileExists(ConfigurationPath) {

			log.Debug("Configuration file does not exist")

			log.Info("Creating configuration file at: ", ConfigurationPath)

			if !processor.CreateDirectory(WorkLogHomeDir) {
				log.Fatal("Error creating worklog home directory")
			}

			log.Debug("Creating worklog home directory")

			// Create the configuration file
			if !processor.CreateFileAsByte(ConfigurationPath, defaultConfigurationData) {
				log.Fatal("Error creating configuration file")
			}

			log.Debug("Unmarshalling default configuration")

			// Unmarshal the default configuration
			err = yaml.Unmarshal(defaultConfigurationData, &configurationContext)
			if err != nil {
				log.Fatal("Error unmarshalling default configuration: ", err)
			}

			log.Info("Configuration file created at: ", ConfigurationPath)
			log.Info("Make sure to customize this file to your needs")

		} else {

			log.Debug("Configuration file exists")

			log.Debug("Loading configuration")

			// Load the configuration
			configurationData, err := os.ReadFile(ConfigurationPath)
			if err != nil {
				log.Fatal("Error loading configuration: ", err)
			}

			log.Debug("Unmarshalling configuration")

			// Unmarshal the configuration
			err = yaml.Unmarshal(configurationData, &configurationContext)
			if err != nil {
				log.Fatal("Error unmarshalling configuration: ", err)
			}
		}

	} else {

		log.Debug("Configuration path provided: ", ConfigurationPath)

		// Check if the configuration file exists
		if !processor.DirectoryOrFileExists(ConfigurationPath) {
			log.Fatal("Configuration file does not exist")
		} else {
			log.Debug("Configuration file exists")
		}

		// Define the default configuration model
		defaultConfiguration := map[interface{}]interface{}{}

		// Define the user configuration model
		userConfiguration := map[interface{}]interface{}{}

		log.Debug("Unmarshalling default configuration")

		// Unmarshal the default configuration
		err = yaml.Unmarshal(defaultConfigurationData, &defaultConfiguration)
		if err != nil {
			log.Fatal("Error unmarshalling default configuration: ", err)
		}

		log.Debug("Loading user configuration")

		// Load the configuration
		configurationData, err := os.ReadFile(ConfigurationPath)
		if err != nil {
			log.Fatal("Error loading configuration: ", err)
		}

		log.Debug("Unmarshalling user configuration")

		// Unmarshal the user configuration
		err = yaml.Unmarshal(configurationData, &userConfiguration)
		if err != nil {
			log.Fatal("Error unmarshalling configuration: ", err)
		}

		log.Debug("Merging configurations")

		// Merge the default and user configurations
		mergedConfiguration := mConfiguration.MergeWithDefault(defaultConfiguration, userConfiguration)

		log.Debug("Marshalling merged configuration")

		// Marshal the merged configuration
		mergedConfigurationMarshalled, err := yaml.Marshal(mergedConfiguration)
		if err != nil {
			log.Fatal("Error marshalling merged configuration: ", err)
		}

		log.Debug("Unmarshalling merged configuration")

		// Unmarshal the merged configuration
		err = yaml.Unmarshal(mergedConfigurationMarshalled, &configurationContext)
		if err != nil {
			log.Fatal("Error unmarshalling merged configuration: ", err)
		}
	}

	setConfigVariables()

	if LogsPath == "" {
		log.Fatal("Logs path not provided in configuration")
	} else {
		if strings.Contains(LogsPath, "~") || strings.Contains(LogsPath, "$HOME") {
			LogsPath = strings.Replace(LogsPath, "~", userHomeDir(), -1)
			LogsPath = strings.Replace(LogsPath, "$HOME", userHomeDir(), -1)
		}
		log.Debug("Logs path: ", LogsPath)
	}
}

func userHomeDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error retrieving user home directory: ", err)
	}
	return homeDir
}

func setConfigVariables() {

	// Set the Logs variables
	log.Debug("Setting Logs variables")
	log.Debug("Setting LogsPath")
	LogsPath = configurationContext.Settings.Logs.Path

	// Set the Git variables
	log.Debug("Setting Git variables")
	log.Debug("Setting GitSync")
	GitSync = configurationContext.Settings.Git.Sync
	log.Debug("Setting GitUri")
	GitUri = configurationContext.Settings.Git.Uri
	log.Debug("Setting GitBranch")
	GitBranch = configurationContext.Settings.Git.Branch

	// Set the Schedule variables
	log.Debug("Setting Schedule variables")
	log.Debug("Setting ScheduleDays variables")
	log.Debug("Setting ScheduleDaysStart")
	ScheduleDaysStart = configurationContext.Settings.Schedule.Days.Start
	log.Debug("Setting ScheduleDaysEnd")
	ScheduleDaysEnd = configurationContext.Settings.Schedule.Days.End
	log.Debug("Setting ScheduleWorkday variables")
	log.Debug("Setting ScheduleWorkdayEnabled")
	ScheduleWorkdayEnabled = configurationContext.Settings.Schedule.Workday.Enabled
	log.Debug("Setting ScheduleWorkdayStart")
	ScheduleWorkdayStart = configurationContext.Settings.Schedule.Workday.Start
	log.Debug("Setting ScheduleWorkdayEnd")
	ScheduleWorkdayEnd = configurationContext.Settings.Schedule.Workday.End
	log.Debug("Setting ScheduleWorkdayTimezone")
	ScheduleWorkdayTimezone = configurationContext.Settings.Schedule.Workday.Timezone

	log.Debug("Configuration variables set")
}
