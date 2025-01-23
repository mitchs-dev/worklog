/*
Copyright © 2025 @mitchs-dev <github@mitchs.dev>
*/
package cli

import (
	"fmt"

	"github.com/mitchs-dev/worklog/internal/logManager"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// listCli represents the list command
var listCli = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list entries in your worklog",
	Long: `This command will list entries in your worklog.

By default, it will list entries only for the day. However, you can specify a period using the --period,-p flag

Available Periods:

  Single Day:
    • today       - Entries from current day (default)
    • yesterday   - Entries from previous day
  
  Multi Day:
    • 3day       - Last 3 days including today
    • week       - Last 7 days from today
    • cweek      - Current week (Start and end date defined in the configuration file)
  
  Extended Periods:
    • month      - Last 30 days
    • quarter    - Last 90 days
    • year       - Last 365 days`,
	Run: func(Cli *cobra.Command, args []string) {

		period, err := Cli.Flags().GetString("period")
		if err != nil {
			log.Fatal("Failed to get period flag")
		}

		log.Debug("Period: ", period)

		log.Debug("Running the list command")
		listEntries, logIds := logManager.Action("list", "", "", period)

		if len(listEntries.Entries) > 0 {
			fmt.Println("Period:", period) // period should be defined earlier in the code
			fmt.Println("Worklog:")
			for _, logId := range logIds {
				logEntry := listEntries.Entries[logId]
				fmt.Printf("- [%s] %s\n", logId, logEntry.Message)
			}
		} else {
			log.Info("No entries found")
		}

	},
}

func init() {
	rootCli.AddCommand(listCli)

	listCli.Flags().StringP("period", "p", "today", "The period to list entries for")
}
