package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
func ConvertRespToSSE(resp *http.Response) (chan Event, error) {
	events := make(chan Event)
	reader := bufio.NewReader(resp.Body)

	go loop(reader, events)

	return events, nil
}

func loop(reader *bufio.Reader, events chan Event) {
	ev := Event{}

	var buf bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "error during resp.Body read:%s\n", err)

			close(events)
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
					events <- ev
					ev = Event{}
				}
			}

		default:
			close(events)
			return
		}
	}
}

func hasPrefix(s []byte, prefix string) bool {
	return bytes.HasPrefix(s, []byte(prefix))
}
