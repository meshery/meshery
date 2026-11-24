package utils

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestConvertRespToSSE(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(
			"id: 123\n" +
				"event: update\n" +
				"data: {\"details\":\"testing\",\"operation_id\":\"op-1\",\"summary\":\"complete\"}\n" +
				"\n",
		)),
	}

	events, err := ConvertRespToSSE(resp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	event, ok := <-events

	if !ok {
		t.Fatal("expected event channel to return an event")
	}

	if event.ID != "123\n" {
		t.Errorf("expected ID 123, got %q", event.ID)
	}

	if event.Name != "update" {
		t.Errorf("expected event name update, got %q", event.Name)
	}

	if event.Data.Details != "testing" {
		t.Errorf(
			"expected details testing, got %q",
			event.Data.Details,
		)
	}

	if event.Data.OperationID != "op-1" {
		t.Errorf(
			"expected operation id op-1, got %q",
			event.Data.OperationID,
		)
	}

	if event.Data.Summary != "complete" {
		t.Errorf(
			"expected summary complete, got %q",
			event.Data.Summary,
		)
	}
}

func TestConvertRespToSSEInvalidJSON(t *testing.T) {
	resp := &http.Response{
		Body: io.NopCloser(strings.NewReader(
			"event: update\n" +
				"data: invalid-json\n" +
				"\n",
		)),
	}

	events, err := ConvertRespToSSE(resp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, ok := <-events

	if ok {
		t.Error("expected channel to close without emitting invalid event")
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		name   string
		input  []byte
		prefix string
		want   bool
	}{
		{
			name:   "matching prefix",
			input:  []byte("data: hello"),
			prefix: "data:",
			want:   true,
		},
		{
			name:   "non matching prefix",
			input:  []byte("event: hello"),
			prefix: "data:",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPrefix(tt.input, tt.prefix)

			if got != tt.want {
				t.Errorf(
					"hasPrefix() = %v, want %v",
					got,
					tt.want,
				)
			}
		})
	}
}
