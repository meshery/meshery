package utils

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertRespToSSE(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Event
		hasError bool
	}{
		{
			name:  "Single event",
			input: "id: 1\nevent: test\ndata: {\"details\":\"detail1\",\"operation_id\":\"op1\",\"summary\":\"summary1\"}\n\n",
			expected: []Event{
				{
					ID:   "1",
					Name: "test",
					Data: EventData{
						Details:     "detail1",
						OperationID: "op1",
						Summary:     "summary1",
					},
				},
			},
			hasError: false,
		},
		{
			name: "Multiple events",
			input: "id: 1\nevent: test1\ndata: {\"details\":\"detail1\",\"operation_id\":\"op1\",\"summary\":\"summary1\"}\n\n" +
				"id: 2\nevent: test2\ndata: {\"details\":\"detail2\",\"operation_id\":\"op2\",\"summary\":\"summary2\"}\n\n",
			expected: []Event{
				{
					ID:   "1",
					Name: "test1",
					Data: EventData{
						Details:     "detail1",
						OperationID: "op1",
						Summary:     "summary1",
					},
				},
				{
					ID:   "2",
					Name: "test2",
					Data: EventData{
						Details:     "detail2",
						OperationID: "op2",
						Summary:     "summary2",
					},
				},
			},
			hasError: false,
		},
		{
			name:  "Comment and retry",
			input: ": this is a comment\nretry: 1000\nid: 1\nevent: test\ndata: {\"details\":\"detail1\",\"operation_id\":\"op1\",\"summary\":\"summary1\"}\n\n",
			expected: []Event{
				{
					ID:   "1",
					Name: "test",
					Data: EventData{
						Details:     "detail1",
						OperationID: "op1",
						Summary:     "summary1",
					},
				},
			},
			hasError: false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: []Event{},
			hasError: false,
		},

	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP response
			resp := httptest.NewRecorder()
			resp.Body = bytes.NewBufferString(tt.input)
			resp.Header().Set("Content-Type", "text/event-stream")

			events, err := ConvertRespToSSE(resp.Result())
			if tt.hasError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			var received []Event
			for event := range events {
				received = append(received, event)
			}

			assert.Equal(t, len(tt.expected), len(received), "expected %d events, got %d", len(tt.expected), len(received))

			for i, event := range received {
				assert.Equal(t, tt.expected[i].ID, strings.TrimSpace(event.ID), "event %d: expected ID %s, got %s", i, tt.expected[i].ID, event.ID)
				assert.Equal(t, tt.expected[i].Name, event.Name, "event %d: expected Name %s, got %s", i, tt.expected[i].Name, event.Name)
				assert.Equal(t, tt.expected[i].Data, event.Data, "event %d: expected Data %+v, got %+v", i, tt.expected[i].Data, event.Data)
			}
		})
	}
}
