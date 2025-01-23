/*
Copyright © 2025 @mitchs-dev <github@mitchs.dev>
*/
package cli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mitchs-dev/library-go/generator"
	"github.com/mitchs-dev/worklog/internal/configuration"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// syncCli represents the sync command
var syncCli = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"sy"},
	Short:   "sync your worklog to Git",
	Long:    `This command will sync your worklog to Git.`,
	Run: func(Cli *cobra.Command, args []string) {

		log.Debug("Running the sync command")

		forceFlag, err := Cli.Flags().GetBool("force")
		if err != nil {
			log.Fatal("Failed to get force flag")
		}

		// First, ensure that Git is enabled in the configuration file
		if !configuration.GitSync {
			log.Fatal("Uh oh! Git is not enabled in the configuration file. Please enable Git and configure it and then try again.")
		}

		// Change to the logs path
		err = os.Chdir(configuration.LogsPath)
		if err != nil {
			log.Fatal("Failed to change to logs path: ", err)
		}

		snapshotTimestamp := generator.StringTimestamp(configuration.ScheduleWorkdayTimezone)
		commitMessage := "SNAPSHOT: " + snapshotTimestamp
		var commitHash string

		// Check that the logs path is a Git repository
		_, err = os.Stat(".git")
		if err != nil {
			// Ask the user if they want to initialize a new Git repository
			log.Warn("It looks like the logs path is not a Git repository. Would you like to initialize it as a new Git repository? (y/n)")
			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				log.Fatal("Failed to read response: ", err)
			}

			if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
				// Initialize the Git repository
				log.Debug("Initializing Git repository")
				cmd := exec.Command("git", "init")
				err := cmd.Run()
				if err != nil {
					log.Fatal("Failed to initialize Git repository: ", err)
				}
			} else {
				log.Fatal("No problem. You won't be able to sync your worklog to Git until you initialize the logs path as a Git repository.")
			}

			log.Info("Git repository initialized")
			log.Warn("There is one extra step to complete. You will need to create the remote repository on your Git server (Github, Gitlab, etc) and then we can sync your worklog to Git.")
			log.Info("Once you have created the remote repository, I will go ahead and add it as the remote origin for you.")
			usernameAndRepository := strings.Split(configuration.GitUri, ":")[1]
			usernameAndRepository = strings.TrimSuffix(usernameAndRepository, ".git")
			log.Warn("Make sure that you use the same exact account and repository (" + usernameAndRepository + ") on your Git server.")
			log.Info("Please press any key to continue once you have completed this step.")
			fmt.Scanln()

			// Add the remote origin
			log.Debug("Adding remote origin")
			cmd := exec.Command("git", "remote", "add", "origin", configuration.GitUri)
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to add remote origin: ", err)
			}

			log.Info("Remote origin added")

			log.Debug("Add the branch: ", configuration.GitBranch)
			cmd = exec.Command("git", "checkout", "-b", configuration.GitBranch)
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to add branch: ", err)
			}

			log.Debug("Running git add .")
			cmd = exec.Command("git", "add", ".")
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to add files: ", err)
			}

			log.Debug("Running git commit -m \"Initial commit\"")
			cmd = exec.Command("git", "commit", "-m", commitMessage)
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to commit files: ", err)
			}

			// Get the commit hash
			log.Debug("Getting the commit hash")
			cmd = exec.Command("git", "rev-parse", "HEAD")
			output, err := cmd.Output()
			if err != nil {
				log.Fatal("Failed to get commit hash: ", err)
			}

			commitHash = strings.TrimSpace(string(output))

			log.Debug("Running git push origin " + configuration.GitBranch)
			cmd = exec.Command("git", "push", "origin", configuration.GitBranch)
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to push changes: ", err)
			}
		} else {

			log.Debug("Fetching the remote origin")
			cmd := exec.Command("git", "fetch", "origin")
			err := cmd.Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to fetch remote origin: %v\n", err)
				log.Fatal("Failed to fetch remote origin: ", err)
			}

			// Check if we're up to date
			cmd = exec.Command("git", "rev-list", "HEAD...origin/"+configuration.GitBranch, "--count")
			output, err := cmd.Output()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to check remote status: %v\n", err)
				log.Fatal("Failed to check remote status: ", err)
			}

			if string(bytes.TrimSpace(output)) == "0" {
				fmt.Println("You're up to date!")
				os.Exit(0)
			}

			// Stash any changes
			log.Debug("Stashing any changes")
			cmd = exec.Command("git", "stash")
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to stash changes: ", err)
			}

			log.Debug("Checking for any upstream changes")
			cmd = exec.Command("git", "pull", "origin", configuration.GitBranch)
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			err = cmd.Run()
			if err != nil {
				log.Fatalf("Failed to pull changes: %s\n%s\n", err, stderr.String())
			}

			// Pop the stash
			log.Debug("Popping the stash")
			cmd = exec.Command("git", "stash", "pop")
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to pop the stash: ", err)
			}

			log.Debug("Running git add .")
			cmd = exec.Command("git", "add", ".")
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to add files: ", err)
			}

			log.Debug("Running git commit -m \"Snapshot\"")
			cmd = exec.Command("git", "commit", "-m", commitMessage)
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to commit files: ", err)
			}

			// Get the commit hash
			log.Debug("Getting the commit hash")
			cmd = exec.Command("git", "rev-parse", "HEAD")
			output, err = cmd.Output()
			if err != nil {
				log.Fatal("Failed to get commit hash: ", err)
			}

			commitHash = strings.TrimSpace(string(output))

			log.Debug("Running git push origin " + configuration.GitBranch)
			if forceFlag {
				cmd = exec.Command("git", "push", "origin", configuration.GitBranch, "--force")
			} else {
				cmd = exec.Command("git", "push", "origin", configuration.GitBranch)
			}
			err = cmd.Run()
			if err != nil {
				log.Fatal("Failed to push changes: ", err)
			}
		}
		log.Info("Worklog synced to Git (Commit: " + commitHash + ")")

	},
}

func init() {
	rootCli.AddCommand(syncCli)

	syncCli.Flags().BoolP("force", "", false, "Force your worklog to sync to Git")
}
