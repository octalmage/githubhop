package gharchive

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/remeh/sizedwaitgroup"
)

// DownloadEventsForDay Download GitHub events for a day.
func DownloadEventsForDay(date time.Time, username string, progress chan bool) []*gabs.Container {
	channel := make(chan *gabs.Container)
	// Decided to use sizedwaitgroup since making 24 HTTP requests at once ended up slowing us down.
	wg := sizedwaitgroup.New(12)

	hours := 24

	// Listen to events and add them to our array.
	var events = []*gabs.Container{}
	go func() {
		for event := range channel {
			events = append(events, event)
		}
	}()

	// Start kicking off HTTP requests.
	for hour := 0; hour <= hours; hour++ {
		dateForUrl := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())
		ghUrl := buildUrl(dateForUrl)
		wg.Add()
		go decodeFromUrl(ghUrl, username, channel, &wg, func() {
			progress <- true
		})
	}

	wg.Wait()
	close(channel)
	close(progress)

	return events
}

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func buildUrl(date time.Time) string {
	return fmt.Sprintf("http://data.gharchive.org/%02d-%02d-%02d-%d.json.gz", date.Year(), int(date.Month()), date.Day(), date.Hour())
}

// Callback function for when decodeFromUrl is done.
type done func()

func decodeFromUrl(url string, username string, channel chan *gabs.Container, wg *sizedwaitgroup.SizedWaitGroup, done done) {
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

		if name == username {
			channel <- parsed
		}
	}
}
