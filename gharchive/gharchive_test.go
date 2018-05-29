package gharchive

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDownloadEventsForDay(t *testing.T) {
	fakeEvent := `{"type":"WatchEvent","actor":{"login":"octalmage"},"repo":{"name":"octalmage/robotjs"},"payload":{"action":"started"},"created_at":"2015-01-01T15:01:57Z"}`

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(gzipString(fakeEvent))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	// Point gharchive to mock server.
	gharchiveURL = ts.URL

	// Create a buffered channel to prevent blocking, we don't care about checking progress.
	progress := make(chan bool, 50)
	events := DownloadEventsForDay(time.Now(), "octalmage", progress)

	if len(events) != 24 {
		t.Errorf("Did not get 24 events back, got %d", len(events))
	}
}

func gzipString(toGzip string) []byte {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(toGzip)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}

	return b.Bytes()
}
