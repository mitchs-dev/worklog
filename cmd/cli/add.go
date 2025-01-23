/*
Copyright © 2025 @mitchs-dev <github@mitchs.dev>
*/
package cli

import (
	"strings"

	"github.com/mitchs-dev/worklog/internal/logManager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// addCli represents the add command
var addCli = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "Add a new entry to your worklog",
	Long:    `This command will add a new entry to your worklog and then display the ID associated with the entry.`,
	Run: func(Cli *cobra.Command, args []string) {

		log.Debug("Running the add command")

		if args == nil {
			log.Fatal("No arguments provided")
		}

		logEntryArgs := args

		log.Debug("logEntry: ", logEntryArgs)

		logEntry := strings.Join(logEntryArgs, " ")

		log.Debug("Calling logManager.Action(\"add\")")
		addedEntry, logIds := logManager.Action("add", logEntry, "", "")

		if len(addedEntry.Entries) == 1 && len(logIds) == 1 {

			logEntry := addedEntry.Entries[logIds[0]]

			if logEntry.Status == logManager.EntryStatusAdded {

				log.Debug("Entry added successfully")
				log.Info("Entry ID: " + logIds[0])

			} else {
				log.Fatal("Unexpected status: ", logEntry.Status)
			}
		} else {
			log.Error("Failed to add")
		}

	},
}

func init() {
	rootCli.AddCommand(addCli)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCli.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCli.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
