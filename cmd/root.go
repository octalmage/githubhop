package cmd

import (
	"fmt"
	"os"
	"time"

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
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "GitHub Username")
	RootCmd.PersistentFlags().StringVarP(&Date, "date", "d", "", "Date to pull")
}

func initFlags() {
	if Username == "" {
		// TODO: Pull GitHub username using "git config user.email" and curl https://api.github.com/search/users?q=jason@stallin.gs+in:email
		Username = "octalmage"
	}
	if Date == "" {
		now := time.Now()
		aYearAgo := now.AddDate(-1, 0, -2)
		Date = aYearAgo.Format("2006-01-02")
	}
}
