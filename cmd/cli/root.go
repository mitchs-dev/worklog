/*
Copyright © 2025 @mitchs-dev <github@mitch.dev>
*/
package cli

import (
	"os"

	"github.com/mitchs-dev/worklog/internal/configuration"
	"github.com/spf13/cobra"
)

// rootCli represents the base command when called without any subcommands
var rootCli = &cobra.Command{
	Use:   "worklog",
	Short: "Worklog is a CLI tool to help you track your work",
	Long:  `Worklog is a CLI tool to help you track your work.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		configuration.ConfigInit()
	},
}

func Execute() {

	err := rootCli.Execute()
	if err != nil {
		os.Exit(1)
	}

}

func init() {
	// Add the flags to the root command
	rootCli.PersistentFlags().BoolVar(&configuration.EnableDebugMode, "debug", false, "Enable debug mode")
	rootCli.PersistentFlags().StringVarP(&configuration.ConfigurationPath, "config", "c", "", "Path to the configuration file")

}
