package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/octalmage/githubhop/github"
	"github.com/spf13/cobra"
)

var Username string
var Date string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "githubhop",
	Short: "Timehop for GitHub",
	Long: `Timehop for GitHub

Running githubhop on it's own will call the run command.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initFlags)
	RootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "GitHub username")
	RootCmd.PersistentFlags().StringVarP(&Date, "date", "d", "", "Date to pull, default is a year ago. Format: YYYY-MM-DD")
}

func initFlags() {
	if Username == "" {
		Username = getGithubUsername()
	}
	if Date == "" {
		now := time.Now()
		aYearAgo := now.AddDate(-1, 0, 0)
		Date = aYearAgo.Format("2006-01-02")
	}
}

func getGithubUsername() string {
	cmd := exec.Command("git", "config", "user.email")
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return ""
	}

	email := string(stdout)
	username, _ := github.GetUsernameForEmail(email)

	return username
}
