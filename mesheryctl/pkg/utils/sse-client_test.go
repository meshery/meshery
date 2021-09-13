package utils

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net/http"
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
				Body: ioutil.NopCloser(bytes.NewBufferString(tt.body)),
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
		name string
		body string
		want []Event
	}{
		{
			name: "Successful SMI test event",
			body: "data: {\"summary\":\"Smi conformance test completed successfully\",\"details\":\"Result-Id: ce336fe4-6e01-4779-949b-6526e44e3d55\",\"operation_id\":\"b33e2b2d-4826-4ea7-8c41-f8f76e8e4bf7\"}\nid: test1\nevent: testevent1\n:\"Comment to be ignored\"\n\n",
			want: []Event{
				{
					ID:   "test1\n",
					Name: "testevent1",
					Data: EventData{Details: "Result-Id: ce336fe4-6e01-4779-949b-6526e44e3d55", OperationID: "b33e2b2d-4826-4ea7-8c41-f8f76e8e4bf7", Summary: "Smi conformance test completed successfully"},
				},
			},
		},
		{
			name: "Failed SMI test event",
			body: "data:{\"event_type\":2,\"summary\":\"Error while running SMI Conformance test\",\"details\":\"Get \\\"https://192.168.49.2:8443/api?timeout=32s\\\": dial tcp 192.168.49.2:8443: i/o timeout\",\"operation_id\":\"95c27959-a4d7-4b35-95d2-50faa04ab641\"}\nid:test2\nevent:testevent2\nretry:\"test\"\n\n",
			want: []Event{
				{
					ID:   "test2\n",
					Name: "testevent2",
					Data: EventData{Details: "Get \"https://192.168.49.2:8443/api?timeout=32s\": dial tcp 192.168.49.2:8443: i/o timeout", OperationID: "95c27959-a4d7-4b35-95d2-50faa04ab641", Summary: "Error while running SMI Conformance test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := http.Response{
				Body: ioutil.NopCloser(bytes.NewBufferString(tt.body)),
			}

			reader := bufio.NewReader(resp.Body)
			events := make(chan Event)

			go loop(reader, events)

			i := 0
			for event := range events {
				if i == len(tt.want) {
					t.Fatalf("Error: Got more events than expected")
				}

				if !reflect.DeepEqual(event, tt.want[i]) {
					t.Fatalf("Error: Wanted event \n%v\nGot event\n%v\n", tt.want[i], event)
				}

				i++
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
