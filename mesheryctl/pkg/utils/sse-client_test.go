package utils

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestSSEConversion(t *testing.T) {
	chunks := []string{
		"id:id1\n",
		"event:event1\n",
		`data: {"details":"details1","operation_id":"op_id1","summary":"summary1"}` + "\n",
		"\n",
		"id: id2\n",
		"event: event2\n",
		`data: {"details":"details2","operation_id":"op_id2","summary":"summary2"}` + "\n",
		"\n",
	}

	stream := io.NopCloser(strings.NewReader(strings.Join(chunks, "")))
	defer func() {
		_ = stream.Close()
	}()

	resp := &http.Response{
		Body: stream,
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	ch, err := ConvertRespToSSE(resp)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		id          string
		name        string
		details     string
		operationID string
		summary     string
	}{
		{"id1", "event1", "details1", "op_id1", "summary1"},
		{"id2", "event2", "details2", "op_id2", "summary2"},
	}

	for i, expected := range tests {
		select {
		case ev := <-ch:
			if ev.ID != expected.id {
				t.Errorf("event %d: expected id %s but got %s", i+1, expected.id, ev.ID)
			}
			if ev.Name != expected.name {
				t.Errorf("event %d: expected name %s but got %s", i+1, expected.name, ev.Name)
			}
			if ev.Data.Details != expected.details {
				t.Errorf("event %d: expected details %s but got %s", i+1, expected.details, ev.Data.Details)
			}
			if ev.Data.OperationID != expected.operationID {
				t.Errorf("event %d: expected operation_id %s but got %s", i+1, expected.operationID, ev.Data.OperationID)
			}
			if ev.Data.Summary != expected.summary {
				t.Errorf("event %d: expected summary %s but got %s", i+1, expected.summary, ev.Data.Summary)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("parser timed out waiting for event %d", i+1)
		}
	}

	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel closure")
	}
}
