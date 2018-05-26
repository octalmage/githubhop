package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/gosuri/uiprogress"
)

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func getUrl(date time.Time) string {
	return fmt.Sprintf("http://data.gharchive.org/%02d-%02d-%02d-%d.json.gz", date.Year(), int(date.Month()), date.Day(), date.Hour())
}

type incr func()

func decodeFromUrl(url string, channel chan *gabs.Container, wg *sync.WaitGroup, done incr) {
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
			fmt.Println(name)
		}
	}
}

func main() {
	now := time.Now()
	aYearAgo := now.AddDate(-1, -1, 0)

	channel := make(chan *gabs.Container)
	wg := sync.WaitGroup{}

	hours := 24
	bar := uiprogress.AddBar(hours).AppendCompleted().PrependElapsed()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return fmt.Sprintf("Fetching hours (%d/%d)", b.Current(), hours)
	})

	uiprogress.Start()
	bar.Incr()
	for hour := 0; hour <= hours-1; hour++ {
		dateForUrl := time.Date(aYearAgo.Year(), aYearAgo.Month(), aYearAgo.Day(), hour, 0, 0, 0, aYearAgo.Location())
		ghUrl := getUrl(dateForUrl)
		wg.Add(1)
		go decodeFromUrl(ghUrl, channel, &wg, func() {
			bar.Incr()
		})
	}

	go func() {
		wg.Wait()
		uiprogress.Stop()
		close(channel)
	}()

	for item := range channel {
		fmt.Println(item)
	}

	// io.Copy(os.Stdout, uncompressed_resp)
	// out, _ := os.Create("output.txt")
	// defer out.Close()
	// io.Copy(out, uncompressed_resp)
}
