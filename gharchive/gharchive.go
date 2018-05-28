package gharchive

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/remeh/sizedwaitgroup"
)

var gharchiveUrl = "https://data.gharchive.org"

// TODO: Add some disk caching.

type SafeArray struct {
	array [][]*gabs.Container
	m     sync.Mutex
}

// DownloadEventsForDay Download GitHub events for a day.
func DownloadEventsForDay(date time.Time, username string, progress chan bool) []*gabs.Container {
	channel := make(chan []*gabs.Container, 1)
	// Decided to use sizedwaitgroup since making 24 HTTP requests at once ended up slowing us down.
	wg := sizedwaitgroup.New(12)

	hours := 24

	var events SafeArray

	// Start kicking off HTTP requests.
	for hour := 0; hour <= hours-1; hour++ {
		dateForUrl := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		ghUrl := buildUrl(dateForUrl)

		wg.Add()
		go decodeFromUrl(ghUrl, username, channel, func() {
			progress <- true

			// Read if we can.
			select {
			case event := <-channel:
				// Lock array for mutating
				events.m.Lock()
				defer events.m.Unlock()
				events.array = append(events.array, event)
			default:
			}

			wg.Done()
		})
	}

	wg.Wait()
	close(channel)
	close(progress)

	// Our channel returns [][]*gabs.Container, we want []*gabs.Container.
	var flattenedEvents []*gabs.Container
	for _, eventArray := range events.array {
		for _, event := range eventArray {
			flattenedEvents = append(flattenedEvents, event)
		}
	}

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
