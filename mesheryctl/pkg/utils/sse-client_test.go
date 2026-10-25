package utils

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func getMockSSE(chunks []string) *http.Response {
	r, w := io.Pipe()
	go func() {
		for _, chunk := range chunks {
			time.Sleep(200 * time.Millisecond)
			_, _ = fmt.Fprint(w, chunk)
		}
		_ = w.Close()
	}()

	return &http.Response{
		Body: r,
	}
}

func TestSSEConversion(t *testing.T) {
	chunks := []string{
		"id: id1\n",
		"event: event1\n",
		`data: {"details":"details1","operation_id":"op_id1","summary":"summary1"}` + "\n",
		"\n",
	}
	resp := getMockSSE(chunks)
	ch, err := ConvertRespToSSE(resp)
	if err != nil {
		t.Fatal(err)
	}
	ev := <-ch
	if ev.ID != "id1" {
		t.Errorf("expected id id1 but got %s", ev.ID)
	}
	if ev.Name != "event1" {
		t.Errorf("expected event event1 but got %s", ev.Name)
	}
	if ev.Data.Details != "details1" {
		t.Errorf("expected details details1 but got %s", ev.Data.Details)
	}
	if ev.Data.OperationID != "op_id1" {
		t.Errorf("expected operations id op_id1 but got %s", ev.Data.OperationID)
	}
	if ev.Data.Summary != "summary1" {
		t.Errorf("expected summary summary1 but got %s", ev.Data.Summary)
	}
}
