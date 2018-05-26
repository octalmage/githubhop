package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Actor struct {
	Id          int    `json:"id"`
	Login       string `json:"login"`
	Gravatar_id string `json:"gravatar_id"`
	Url         string `json:"url"`
	Avatar_url  string `json:"avatar_url"`
}
type Payload struct {
	Ref           string `json:"ref"`
	Ref_type      string `json:"ref_type"`
	Master_branch string `json:"master_branch"`
	Description   string `json:"description"`
	Pusher_type   string `json:"pusher_type"`
}

type Event struct {
	Id        string                 `json:"id"`
	Type      string                 `json:"Type"`
	CreatedAt string                 `json:"created_at"`
	Actor     Actor                  `json:"actor"`
	Payload   map[string]interface{} `json:"payload"`
	Public    bool                   `json:"public"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getUrl(date time.Time) string {
	return fmt.Sprintf("http://data.gharchive.org/%02d-%02d-%02d-%d.json.gz", date.Year(), int(date.Month()), date.Day(), date.Hour())
}

func decodeFromUrl(url string, channel chan Event, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	uncompressed_resp, _ := gzip.NewReader(resp.Body)

	dec := json.NewDecoder(uncompressed_resp)

	// Decode event.
	for dec.More() {
		var e Event
		err := dec.Decode(&e)
		if err != nil {
			log.Fatal(err)
		}
		if e.Actor.Login == "octalmage" {
			channel <- e
			fmt.Println(e)
		}
	}
	wg.Done()
}

func main() {
	now := time.Now()
	aYearAgo := now.AddDate(-1, 0, 0)

	channel := make(chan Event)
	wg := sync.WaitGroup{}
	for hour := 0; hour <= 23; hour++ {
		dateForUrl := time.Date(aYearAgo.Year(), aYearAgo.Month(), aYearAgo.Day(), hour, 0, 0, 0, aYearAgo.Location())
		ghUrl := getUrl(dateForUrl)
		wg.Add(1)
		go decodeFromUrl(ghUrl, channel, &wg)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for item := range channel {
		fmt.Println(item.Type)
	}

	// io.Copy(os.Stdout, uncompressed_resp)
	// out, _ := os.Create("output.txt")
	// defer out.Close()
	// io.Copy(out, uncompressed_resp)
}
