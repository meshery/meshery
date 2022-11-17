package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestConvertRespToSSE(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{
			name: "Convert Response to SSE with sample body",
			body: "Test data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := http.Response{
				Body: io.NopCloser(bytes.NewBufferString(tt.body)),
			}

			_, err := ConvertRespToSSE(&resp)

			if err != nil {
				t.Fatalf("ConvertRespToSSE() error = %v", err)
				return
			}
		})
	}
}

func Test_loop(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		fixture string
		golden  string
	}{
		{
			name:    "Successful SMI test event",
			body:    "data: {\"summary\":\"Smi conformance test completed successfully\",\"details\":\"Result-Id: ce336fe4-6e01-4779-949b-6526e44e3d55\",\"operation_id\":\"b33e2b2d-4826-4ea7-8c41-f8f76e8e4bf7\"}\nid: test1\nevent: testevent1\n:\"Comment to be ignored\"\n\n",
			fixture: "loop.success.golden",
			golden:  "loop.expect.success.golden",
		},
		{
			name:    "Failed SMI test event",
			body:    "data:{\"event_type\":2,\"summary\":\"Error while running SMI Conformance test\",\"details\":\"Get \\\"https://192.168.49.2:8443/api?timeout=32s\\\": dial tcp 192.168.49.2:8443: i/o timeout\",\"operation_id\":\"95c27959-a4d7-4b35-95d2-50faa04ab641\"}\nid:test2\nevent:testevent2\nretry:\"test\"\n\n",
			fixture: "loop.error.golden",
			golden:  "loop.expect.error.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Read fixture
			bodyBytes, err := os.ReadFile("fixtures/sse_client/" + tt.fixture)
			if err != nil {
				t.Fatalf("Error while reading fixture: %v", err)
			}

			body := string(bodyBytes)

			resp := http.Response{
				Body: io.NopCloser(bytes.NewBufferString(body)),
			}

			reader := bufio.NewReader(resp.Body)
			events := make(chan Event)

			go loop(reader, events)

			gotEvents := []Event{}
			i := 0

			for event := range events {
				gotEvents = append(gotEvents, event)
				i++
			}

			// Read golden file
			golden := NewGoldenFile(t, tt.golden, "testdata/sse_client/")
			if *update {
				data, _ := json.Marshal(gotEvents)
				golden.Write(string(data))
			}

			wantEvents := []Event{}
			err = json.Unmarshal([]byte(golden.Load()), &wantEvents)

			if err != nil {
				t.Fatalf("Error in unmarshalling golden file: %v", err)
			}

			if i > len(wantEvents) {
				t.Fatalf("Error: Got more events [%v] than expected [%v]", i, wantEvents)
			}

			if !reflect.DeepEqual(gotEvents, wantEvents) {
				t.Fatalf("Error: Wanted events \n%v\nGot events\n%v\n", wantEvents, gotEvents)
			}
		})
	}
}

func Test_hasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		prefix string
		want   bool
	}{
		{
			name:   "Required prefix present",
			s:      "Test String",
			prefix: "Test",
			want:   true,
		},
		{
			name:   "Required prefix not present",
			s:      "Test String",
			prefix: "1234",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasPrefix([]byte(tt.s), tt.prefix); got != tt.want {
				t.Errorf("hasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
