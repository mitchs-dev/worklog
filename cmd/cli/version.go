/*
Copyright © 2025 @mitchs-dev <github@mitch.dev>
*/
package cli

import (
	"fmt"

	"github.com/mitchs-dev/worklog/internal/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// versionCli represents the base command when called without any subcommands
var versionCli = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Worklog",
	Long:  `Print the version number of Worklog.`,
	Run: func(cmd *cobra.Command, args []string) {

		log.Debug("Running the version command")

		format, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal("Failed to get format flag")
		}
		fmt.Println(version.GetVersion(format))
	},
}

func init() {
	rootCli.AddCommand(versionCli)

	versionCli.Flags().StringP("output", "o", "plain", "Output format (plain,json, yaml)")
}
