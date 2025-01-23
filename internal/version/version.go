package version

import (
	"embed"
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//go:embed version.json
var versionFileEmbed embed.FS

var (
	// VersionFilePath is the path to the version file
	VersionFilePath = "version.json"
)

// Version is the version of the tool
type Version struct {
	Semantic string `json:"semantic" yaml:"semantic"`
	Build    struct {
		Commit string `json:"commit" yaml:"commit"`
		Date   string `json:"date" yaml:"date"`
	} `json:"build" yaml:"build"`
}

func GetVersion(format string) string {
	versionData, err := versionFileEmbed.ReadFile(VersionFilePath)
	if err != nil {
		log.Fatalf("Error reading version: %v", err)
	}
	if len(versionData) == 0 {
		log.Fatalf("Error version is empty")
	}

	var VersionStruct Version
	format = strings.ToLower(format)
	log.Debugf("Version format: %s", format)

	switch format {
	case "plain":
		if err = json.Unmarshal(versionData, &VersionStruct); err != nil {
			log.Fatalf("Error unmarshalling version: %v", err)
		}
		return `Version: ` + VersionStruct.Semantic + `
Commit: ` + VersionStruct.Build.Commit + `
Build Date: ` + VersionStruct.Build.Date
	case "json":
		if err = json.Unmarshal(versionData, &VersionStruct); err != nil {
			log.Fatalf("Error unmarshalling version: %v", err)
		}
		version, err := json.Marshal(VersionStruct) // Changed from 'version' to 'VersionStruct'
		if err != nil {
			log.Fatalf("Error marshalling version: %v", err)
		}
		return strings.TrimSpace(string(version))
	case "yaml":
		if err = yaml.Unmarshal(versionData, &VersionStruct); err != nil {
			log.Fatalf("Error unmarshalling version: %v", err)
		}
		version, err := yaml.Marshal(VersionStruct) // Changed from 'version' to 'VersionStruct'
		if err != nil {
			log.Fatalf("Error marshalling version: %v", err)
		}
		return strings.TrimSpace(string(version))
	default:
		log.Fatalf("Error unknown format: %s", format)
		return ""
	}
}
