package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/octalmage/githubhop/gharchive"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Pull GitHub events from this day last year.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		getEvents(Username, Date)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
}

// Convert event_type to english.
func convertEventType(eventType string) string {
	var types = map[string]string{
		"CommitCommentEvent":            "commented a commit on the repo",
		"CreateEvent":                   "created the",
		"DeleteEvent":                   "deleted the",
		"ForkEvent":                     "forked the repo",
		"IssueCommentEvent":             "commented on an issue on",
		"IssuesEvent":                   "created an issue on",
		"MemberEvent":                   "were added to the repo",
		"PullRequestEvent":              "created a PR on",
		"PullRequestReviewCommentEvent": "commented on a PR on the repo",
		"PushEvent":                     "pushed commits to",
		"ReleaseEvent":                  "published a new release of",
		"WatchEvent":                    "watched the repo",
	}

	return types[eventType]
}

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func getEvents(username string, date string) {
	aYearAgo, _ := time.Parse("2006-01-02", date)

	hours := 24

	bar := uiprogress.AddBar(hours).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Fetching hours (%d/%d)", b.Current(), hours)
	})

	uiprogress.Start()

	progress := make(chan bool)
	go func() {
		for range progress {
			bar.Incr()
		}
	}()

	events := gharchive.DownloadEventsForDay(aYearAgo, username, progress)

	for _, event := range events {
		eventType := event.Path("type").Data().(string)
		action := convertEventType(eventType)

		if eventType == "CreateEvent" {
			refType := event.Path("payload.ref_type").Data().(string)
			action += fmt.Sprintf(" %s", refType)
			if refType != "repository" {
				ref := event.Path("payload.ref").Data().(string)
				action += fmt.Sprintf(" %s on", ref)
			}
		}

		if eventType == "DeleteEvent" {
			refType := event.Path("payload.ref_type").Data().(string)
			ref := event.Path("payload.ref").Data().(string)
			action += fmt.Sprintf(" %s %s on", refType, ref)
		}

		target := event.Path("repo.name").Data().(string)

		// Parse and format the created_at date.
		created_at := event.Path("created_at").Data().(string)
		t, _ := time.Parse("2006-01-02T15:04:05Z", created_at)
		date := t.Format("2006-01-02 3:04pm")

		// Print a formatted message to the screen.
		fmt.Printf(
			"At %s, you %s %s\n",
			date,
			action,
			target,
		)
	}
}