package sse

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestStream_HeartbeatEmitted shrinks HeartbeatInterval to ~20 ms so the
// heartbeat path can be observed without waiting the production 25 s. The
// var is restored via t.Cleanup. Tests run sequentially in this file so
// the global mutation is safe; t.Parallel is not used.

func TestStream_HeadersAndInitialFlush(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan any)

	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	// Wait for the initial flush so we know headers are committed.
	waitForFlush(t, rec, 1, 500*time.Millisecond)

	got := rec.Header()
	if got.Get("Content-Type") != "text/event-stream" {
		t.Errorf("Content-Type = %q, want text/event-stream", got.Get("Content-Type"))
	}
	if got.Get("Cache-Control") != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache", got.Get("Cache-Control"))
	}
	if got.Get("Connection") != "keep-alive" {
		t.Errorf("Connection = %q, want keep-alive", got.Get("Connection"))
	}
	if got.Get("X-Accel-Buffering") != "no" {
		t.Errorf("X-Accel-Buffering = %q, want no", got.Get("X-Accel-Buffering"))
	}

	cancel()
	close(events)
	if err := <-done; err != nil {
		t.Errorf("Stream returned error on clean shutdown: %v", err)
	}
}

func TestStream_WritesJSONEvent(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 1)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	events <- map[string]string{"hello": "world"}

	waitForBody(t, rec, `data: {"hello":"world"}`, time.Second)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error: %v", err)
	}
}

func TestStream_NamedEvent(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 1)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	events <- Event{Name: "k8s_contexts", Data: map[string]int{"totalCount": 2}}

	waitForBody(t, rec, "event: k8s_contexts\ndata: {\"totalCount\":2}", time.Second)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error: %v", err)
	}
}

func TestStream_NamedEventViaPointer(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 1)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	// Publishing &Event{...} is idiomatic Go; the type switch must handle
	// the pointer form or it would silently degrade to an unnamed message
	// containing the serialized envelope (i.e. data: {"Name":"...","Data":...}).
	events <- &Event{Name: "controllers_status", Data: []int{1, 2, 3}}

	waitForBody(t, rec, "event: controllers_status\ndata: [1,2,3]", time.Second)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error: %v", err)
	}
}

func TestStream_CompactsMultilinePayload(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 1)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	// json.RawMessage with embedded newlines is the common way a multi-line
	// payload sneaks into the stream. SSE treats bare LFs as field
	// terminators, so without compaction the client would see a truncated
	// or split frame. Verify the published value is collapsed onto one line.
	events <- json.RawMessage("{\n  \"a\": 1,\n  \"b\": 2\n}")

	waitForBody(t, rec, `data: {"a":1,"b":2}`, time.Second)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error: %v", err)
	}
}

func TestStream_MarshalFailureSkipsEventButKeepsStreamOpen(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 2)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	// A channel value isn't JSON-marshalable. Sending it should be logged
	// and skipped — the stream must stay open so the next (valid) event
	// still gets delivered. This is the contract called out in the marshal-
	// failure branch comment in Stream.
	events <- make(chan int)
	events <- map[string]string{"ok": "after-bad"}

	waitForBody(t, rec, `data: {"ok":"after-bad"}`, time.Second)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error: %v", err)
	}
}

func TestStream_ContextCancellationStopsLoop(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	events := make(chan any)

	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	// Wait for the loop to be running (initial flush committed).
	waitForFlush(t, rec, 1, 500*time.Millisecond)

	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Stream returned error on ctx cancel: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Stream did not exit after ctx cancel")
	}
}

func TestStream_EventsChannelCloseStopsLoop(t *testing.T) {
	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	waitForFlush(t, rec, 1, 500*time.Millisecond)

	close(events)
	select {
	case err := <-done:
		if err != nil {
			t.Errorf("Stream returned error on events close: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Stream did not exit after events channel closed")
	}
}

func TestStream_WriteErrorStopsLoop(t *testing.T) {
	rec := &errorWriter{
		ResponseWriter: httptest.NewRecorder(),
		failAfter:      0, // fail every Write call; initial flusher.Flush() doesn't go through Write
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any, 1)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	events <- map[string]string{"a": "b"}

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Stream returned nil on write error, want non-nil")
		}
		if !errors.Is(err, errInjected) && !strings.Contains(err.Error(), "injected") {
			t.Errorf("Stream returned %v, want injected write error", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Stream did not exit after write error")
	}
}

func TestStream_HeartbeatEmitted(t *testing.T) {
	orig := HeartbeatInterval
	HeartbeatInterval = 20 * time.Millisecond
	t.Cleanup(func() { HeartbeatInterval = orig })

	rec := newFlushRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events := make(chan any)
	done := make(chan error, 1)
	go func() {
		done <- Stream(ctx, rec, events, nil)
	}()

	waitForBody(t, rec, ": keepalive", 500*time.Millisecond)

	cancel()
	if err := <-done; err != nil {
		t.Errorf("Stream returned error on shutdown: %v", err)
	}
}

func TestStream_FlusherUnavailable(t *testing.T) {
	// nonFlushingWriter intentionally does NOT implement http.Flusher.
	w := &nonFlushingWriter{header: http.Header{}}
	events := make(chan any)
	err := Stream(context.Background(), w, events, nil)
	if !errors.Is(err, ErrFlusherUnavailable) {
		t.Errorf("Stream() = %v, want ErrFlusherUnavailable", err)
	}
}

// --- helpers ---

// flushRecorder is an http.ResponseWriter that records the body, headers,
// and number of Flush() calls. It is concurrency-safe so the test goroutine
// and the Stream goroutine can both touch it.
type flushRecorder struct {
	mu      sync.Mutex
	header  http.Header
	body    strings.Builder
	flushes int
}

func newFlushRecorder() *flushRecorder {
	return &flushRecorder{header: http.Header{}}
}

func (r *flushRecorder) Header() http.Header {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.header
}

func (r *flushRecorder) Write(b []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.body.Write(b)
}

func (r *flushRecorder) WriteHeader(int) {}

func (r *flushRecorder) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.flushes++
}

func (r *flushRecorder) flushCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.flushes
}

func (r *flushRecorder) bodyString() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.body.String()
}

func waitForFlush(t *testing.T, r *flushRecorder, want int, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if r.flushCount() >= want {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for flush count %d (got %d)", want, r.flushCount())
}

func waitForBody(t *testing.T, r *flushRecorder, contains string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if strings.Contains(r.bodyString(), contains) {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for body to contain %q (got %q)", contains, r.bodyString())
}

// errorWriter wraps an httptest.ResponseRecorder but starts returning
// errInjected on Write after failAfter successful writes. The first write
// from Stream is just the header flush via flusher.Flush() which doesn't
// touch Write — so callers need to size failAfter against actual data
// writes.
var errInjected = errors.New("injected write error")

type errorWriter struct {
	http.ResponseWriter
	mu        sync.Mutex
	writes    int
	failAfter int
}

func (e *errorWriter) Write(b []byte) (int, error) {
	e.mu.Lock()
	e.writes++
	n := e.writes
	e.mu.Unlock()
	if n > e.failAfter {
		return 0, errInjected
	}
	return e.ResponseWriter.Write(b)
}

func (e *errorWriter) Flush() {
	if f, ok := e.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// nonFlushingWriter is the minimum http.ResponseWriter — no Flush method.
type nonFlushingWriter struct {
	header http.Header
}

func (w *nonFlushingWriter) Header() http.Header         { return w.header }
func (w *nonFlushingWriter) Write(b []byte) (int, error) { return len(b), nil }
func (w *nonFlushingWriter) WriteHeader(int)             {}
