package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
)

// sseRecorder is an http.ResponseWriter that also satisfies http.Flusher, so
// SubscribeEventsHandler takes its streaming path. Writes are funnelled through
// a safeBuffer because the handler writes from its own goroutine while the test
// reads the captured body.
type sseRecorder struct {
	header http.Header
	body   *safeBuffer
	*testFlusher
}

func newSSERecorder(writeCh chan struct{}) *sseRecorder {
	return &sseRecorder{
		header:      make(http.Header),
		body:        &safeBuffer{writeCh: writeCh},
		testFlusher: &testFlusher{},
	}
}

func (r *sseRecorder) Header() http.Header         { return r.header }
func (r *sseRecorder) Write(p []byte) (int, error) { return r.body.Write(p) }
func (r *sseRecorder) WriteHeader(int)             {}

func newEventsTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("Events")},
		log:    newTestLogger(t),
	}
}

func TestSubscribeEventsHandler_StreamsPublishedEvent(t *testing.T) {
	h := newEventsTestHandler(t)
	userID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: userID}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx)

	writeCh := make(chan struct{}, 8)
	rec := newSSERecorder(writeCh)

	done := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec, req, nil, user, nil)
		close(done)
	}()

	event := events.NewEvent().
		FromOwner(userID).
		WithCategory("test").
		WithAction("streamed").
		WithSeverity(events.Informational).
		WithDescription("a thing happened").
		Build()

	// The handler subscribes after writing the initial headers, so a single
	// Publish can race the Subscribe. Retry until the framed event lands in the
	// body (or the test times out).
	deadline := time.After(testTimeout)
	for !strings.Contains(rec.body.String(), "a thing happened") {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for event frame; got: %q", rec.body.String())
		default:
		}
		h.config.EventBroadcaster.Publish(userID, event)
		select {
		case <-writeCh:
		case <-time.After(20 * time.Millisecond):
		}
	}

	body := rec.body.String()
	if !strings.Contains(body, "data: ") || !strings.Contains(body, "\n\n") {
		t.Fatalf("event not framed as SSE data message: %q", body)
	}
	if got := rec.header.Get("Content-Type"); got != "text/event-stream" {
		t.Fatalf("Content-Type = %q, want text/event-stream", got)
	}

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("handler did not return after context cancellation")
	}
}

func TestSubscribeEventsHandler_StopsOnContextCancellation(t *testing.T) {
	h := newEventsTestHandler(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx)
	rec := newSSERecorder(nil)

	done := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec, req, nil, user, nil)
		close(done)
	}()

	cancel()
	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("handler did not return after context cancellation")
	}
}
