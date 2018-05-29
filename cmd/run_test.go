package cmd

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/octalmage/githubhop/gharchive"
)

func TestRunCommand(t *testing.T) {
	fakeEvent := `{"type":"WatchEvent","actor":{"login":"octalmage"},"repo":{"name":"octalmage/robotjs"},"payload":{"action":"started"},"created_at":"2015-01-01T15:01:57Z"}`

	testGetter := func(date time.Time, username string, _ gharchive.Cache, progress chan bool) []*gabs.Container {
		event, _ := gabs.ParseJSON([]byte(fakeEvent))
		progress <- true
		events := []*gabs.Container{event}
		return events
	}

	var buffer bytes.Buffer
	getEvents("octalmage", "2015-01-01", testGetter, &buffer)
	output := buffer.String()

	if !strings.Contains(output, "At 2015-01-01 3:01pm, you watched the repo octalmage/robotjs") {
		t.Errorf("Did not get expected response, got: %s", output)
	}
}
