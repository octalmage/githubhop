package gharchive

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/remeh/sizedwaitgroup"
)

var gharchiveURL = "https://data.gharchive.org"

// TODO: Add some disk caching.

// DownloadEventsForDay Download GitHub events for a day.
func DownloadEventsForDay(date time.Time, username string, progress chan bool) []*gabs.Container {
	// Create buffered channel with a size of 24, since we know we'll have 24 workers.
	channel := make(chan []*gabs.Container, 24)
	// Decided to use sizedwaitgroup since making 24 HTTP requests at once ended up slowing us down.
	wg := sizedwaitgroup.New(12)

	hours := 24

	// Start kicking off HTTP requests.
	for hour := 0; hour <= hours-1; hour++ {
		dateForURL := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		ghURL := buildURL(dateForURL)

		wg.Add()
		go decodeFromURL(ghURL, username, channel, func() {
			progress <- true
			wg.Done()
		})
	}

	wg.Wait()
	close(channel)
	close(progress)

	// Pull events from our channel.
	var flattenedEvents []*gabs.Container
	for events := range channel {
		for _, event := range events {
			flattenedEvents = append(flattenedEvents, event)
		}
	}

	// Since events can come back in any order, sort them!
	sort.Slice(flattenedEvents, func(i, j int) bool {
		format := "2006-01-02T15:04:05Z"
		createdAt1 := flattenedEvents[i].Path("created_at").Data().(string)
		createdAt2 := flattenedEvents[j].Path("created_at").Data().(string)
		t1, _ := time.Parse(format, createdAt1)
		t2, _ := time.Parse(format, createdAt2)
		return t1.Before(t2)
	})

	return flattenedEvents
}

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func buildURL(date time.Time) string {
	return fmt.Sprintf("%s/%02d-%02d-%02d-%d.json.gz", gharchiveURL, date.Year(), int(date.Month()), date.Day(), date.Hour())
}

// Callback function for when decodeFromUrl is done.
type done func()

func decodeFromURL(url string, username string, channel chan []*gabs.Container, done done) {
	defer done()
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	uncompressedResp, _ := gzip.NewReader(resp.Body)

	dec := json.NewDecoder(uncompressedResp)
	var events []*gabs.Container
	// Decode event.
	for dec.More() {
		parsed, err := gabs.ParseJSONDecoder(dec)
		check(err)
		name, _ := parsed.Path("actor.login").Data().(string)

		if name == username {
			events = append(events, parsed)
		}
	}
	channel <- events
}
