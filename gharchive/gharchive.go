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

type Cache interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, bool)
}

var gharchiveUrl = "https://data.gharchive.org"

// TODO: Add some disk caching.

// DownloadEventsForDay Download GitHub events for a day.
func DownloadEventsForDay(date time.Time, username string, cache Cache, progress chan bool) []*gabs.Container {
	// Create buffered channel with a size of 24, since we know we'll have 24 workers.
	channel := make(chan []*gabs.Container, 24)
	// Decided to use sizedwaitgroup since making 24 HTTP requests at once ended up slowing us down.
	wg := sizedwaitgroup.New(12)

	hours := 24

	// Start kicking off HTTP requests.
	for hour := 0; hour <= hours-1; hour++ {
		dateForUrl := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		ghUrl := buildUrl(dateForUrl)

		wg.Add()
		go decodeFromUrl(ghUrl, username, channel, func() {
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
	fmt.Println(flattenedEvents)
	stringEvents := gabs.New()
	for _, event := range flattenedEvents {
		// stringEvents.Index(i).Set(event)
		stringEvents.ArrayAppend(event.Data())
	}

	stringEvents.ArrayRemove(0)

	fmt.Println(stringEvents.String())

	cache.Set(date.Format("20060102"), []byte(stringEvents.String()))

	// Since events can come back in any order, sort them!
	sort.Slice(flattenedEvents, func(i, j int) bool {
		format := "2006-01-02T15:04:05Z"
		created_at1 := flattenedEvents[i].Path("created_at").Data().(string)
		created_at2 := flattenedEvents[j].Path("created_at").Data().(string)
		t1, _ := time.Parse(format, created_at1)
		t2, _ := time.Parse(format, created_at2)
		return t1.Before(t2)
	})

	return flattenedEvents
}

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func buildUrl(date time.Time) string {
	return fmt.Sprintf("%s/%02d-%02d-%02d-%d.json.gz", gharchiveUrl, date.Year(), int(date.Month()), date.Day(), date.Hour())
}

// Callback function for when decodeFromUrl is done.
type done func()

func decodeFromUrl(url string, username string, channel chan []*gabs.Container, done done) {
	defer done()
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	uncompressed_resp, _ := gzip.NewReader(resp.Body)

	dec := json.NewDecoder(uncompressed_resp)
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
