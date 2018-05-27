package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/gosuri/uiprogress"
	"github.com/remeh/sizedwaitgroup"
)

// TODO: Bring in cobra for command line options.

// Convert event_type to english.
func convertEventType(eventType string) string {
	var types = map[string]string{
		"CommitCommentEvent":            "commented a commit on the repo",
		"CreateEvent":                   "created the",
		"DeleteEvent":                   "deleted the",
		"ForkEvent":                     "forked the repo",
		"IssueCommentEvent":             "commented on an issue on",
		"IssuesEvent":                   "created an issue on",
		"MemberEvent":                   "was added to the repo",
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

func getUrl(date time.Time) string {
	return fmt.Sprintf("http://data.gharchive.org/%02d-%02d-%02d-%d.json.gz", date.Year(), int(date.Month()), date.Day(), date.Hour())
}

type incr func()

func decodeFromUrl(url string, channel chan *gabs.Container, wg *sizedwaitgroup.SizedWaitGroup, done incr) {
	defer done()
	defer wg.Done()
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	uncompressed_resp, _ := gzip.NewReader(resp.Body)

	dec := json.NewDecoder(uncompressed_resp)
	// Decode event.
	for dec.More() {
		parsed, err := gabs.ParseJSONDecoder(dec)
		check(err)
		name, _ := parsed.Path("actor.login").Data().(string)

		if name == "octalmage" {
			channel <- parsed
		}
	}
}

func main() {
	now := time.Now()
	aYearAgo := now.AddDate(-1, -2, 0)

	channel := make(chan *gabs.Container)
	wg := sizedwaitgroup.New(runtime.NumCPU())

	hours := 24
	bar := uiprogress.AddBar(hours).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Fetching hours (%d/%d)", b.Current(), hours)
	})

	var events = []*gabs.Container{}
	go func() {
		for event := range channel {
			events = append(events, event)
		}
	}()

	uiprogress.Start()
	for hour := 0; hour <= hours; hour++ {
		dateForUrl := time.Date(aYearAgo.Year(), aYearAgo.Month(), aYearAgo.Day(), hour, 0, 0, 0, aYearAgo.Location())
		ghUrl := getUrl(dateForUrl)
		wg.Add()
		go decodeFromUrl(ghUrl, channel, &wg, func() {
			bar.Incr()
		})
	}

	wg.Wait()
	close(channel)
	uiprogress.Stop()

	for _, event := range events {
		eventType := event.Path("type").Data().(string)
		action := convertEventType(eventType)

		if eventType == "CreateEvent" {
			refType := event.Path("payload.ref_type").Data().(string)
			ref := event.Path("payload.ref").Data().(string)
			action += fmt.Sprintf(" %s", refType)
			if refType != "repository" {
				action += fmt.Sprintf(" %s on", ref)
			}
		}

		if eventType == "DeleteEvent" {
			refType := event.Path("payload.ref_type").Data().(string)
			ref := event.Path("payload.ref").Data().(string)
			action += fmt.Sprintf(" %s %s on", refType, ref)
		}

		target := event.Path("repo.name").Data().(string)

		fmt.Printf(
			"At %s, you %s %s\n",
			event.Path("created_at").Data().(string),
			action,
			target,
		)
	}
}
