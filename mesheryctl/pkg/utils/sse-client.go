package utils

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

// Event represents a Server-Sent Event
type Event struct {
	Name string
	ID   string
	Data EventData
}

type EventData struct {
	Details     string `json:"details"`
	OperationID string `json:"operation_id"`
	Summary     string `json:"summary"`
}

// ConvertRespToSSE converts a connection to a stream of server sent events
func ConvertRespToSSE(ctx context.Context, resp *http.Response) (<-chan Event, error) {
	events := make(chan Event)
	reader := bufio.NewReader(resp.Body)

	go loop(ctx, resp, reader, events)

	return events, nil
}

func loop(ctx context.Context, resp *http.Response, reader *bufio.Reader, events chan<- Event) {
	var once sync.Once
	closeEvents := func() {
		once.Do(func() {
			close(events)
		})
	}
	defer closeEvents()
	defer func() {
		_ = resp.Body.Close()
	}()

	stopCtxMonitor := make(chan struct{})
	defer close(stopCtxMonitor)

	go func() {
		select {
		case <-ctx.Done():
			_ = resp.Body.Close()
		case <-stopCtxMonitor:
		}
	}()

	ev := Event{}

	var buf bytes.Buffer

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err != nil {
			if ctx.Err() == nil {
				fmt.Fprintf(os.Stderr, "error during resp.Body read:%s\n", err)
			}
			return
		}

		switch {
		case hasPrefix(line, ":"):
			// Comment, do nothing
		case hasPrefix(line, "retry:"):
			// Retry, do nothing for now

		// id of event
		case hasPrefix(line, "id: "):
			ev.ID = string(line[4:])
		case hasPrefix(line, "id:"):
			ev.ID = string(line[3:])

		// name of event
		case hasPrefix(line, "event: "):
			ev.Name = string(line[7 : len(line)-1])
		case hasPrefix(line, "event:"):
			ev.Name = string(line[6 : len(line)-1])

		// event data
		case hasPrefix(line, "data: "):
			buf.Write(line[6:])
		case hasPrefix(line, "data:"):
			buf.Write(line[5:])

		// end of event
		case bytes.Equal(line, []byte("\n")):
			b := buf.Bytes()

			if hasPrefix(b, "{") {
				var data EventData

				err := json.Unmarshal(b, &data)

				if err == nil {
					ev.Data = data
					buf.Reset()
					select {
					case events <- ev:
					case <-ctx.Done():
						return
					}
					ev = Event{}
				}
			}

		default:
			return
		}
	}
}

func hasPrefix(s []byte, prefix string) bool {
	return bytes.HasPrefix(s, []byte(prefix))
}
