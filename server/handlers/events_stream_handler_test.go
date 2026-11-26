package handlers

import (
	"context"
	"fmt"
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

func newTestEventWithID(userID uuid.UUID, desc string) *events.Event {
	ev := events.NewEvent().
		FromOwner(userID).
		WithCategory("test").
		WithAction("streamed").
		WithSeverity(events.Informational).
		WithDescription(desc).
		Build()
	ev.ID = uuid.Must(uuid.NewV4())
	return ev
}

func waitForBodyContains(t *testing.T, rec *sseRecorder, writeCh chan struct{}, substr string) {
	t.Helper()
	deadline := time.After(testTimeout)
	for !strings.Contains(rec.body.String(), substr) {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for %q in SSE body; got: %q", substr, rec.body.String())
		case <-writeCh:
		case <-time.After(20 * time.Millisecond):
		}
	}
}

func TestSubscribeEventsHandler_FramesEventID(t *testing.T) {
	h := newEventsTestHandler(t)
	userID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: userID}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx)

	writeCh := make(chan struct{}, 16)
	rec := newSSERecorder(writeCh)

	done := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec, req, nil, user, nil)
		close(done)
	}()

	eventWithID := newTestEventWithID(userID, "with-id-event")
	deadline := time.After(testTimeout)
	for !strings.Contains(rec.body.String(), "with-id-event") {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for event; got: %q", rec.body.String())
		default:
		}
		h.config.EventBroadcaster.Publish(userID, eventWithID)
		select {
		case <-writeCh:
		case <-time.After(20 * time.Millisecond):
		}
	}

	wantIDLine := fmt.Sprintf("id: %s:", eventWithID.ID.String())
	if !strings.Contains(rec.body.String(), wantIDLine) {
		t.Fatalf("expected SSE frame to contain %q, got: %q", wantIDLine, rec.body.String())
	}

	// An event with uuid.Nil ID must omit the `id:` line.
	nilIDEvent := events.NewEvent().
		FromOwner(userID).
		WithCategory("test").
		WithAction("streamed").
		WithSeverity(events.Informational).
		WithDescription("nil-id-event").
		Build()
	nilIDEvent.ID = uuid.Nil
	h.config.EventBroadcaster.Publish(userID, nilIDEvent)
	waitForBodyContains(t, rec, writeCh, "nil-id-event")

	for _, frame := range strings.Split(rec.body.String(), "\n\n") {
		if strings.Contains(frame, "nil-id-event") {
			if strings.Contains(frame, "id: ") {
				t.Fatalf("expected nil-ID event to omit id: line, got frame %q", frame)
			}
			break
		}
	}

	cancel()
	<-done
}

func TestSubscribeEventsHandler_ReplaysMissedEventsOnLastEventID(t *testing.T) {
	h := newEventsTestHandler(t)
	userID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: userID}

	e1 := newTestEventWithID(userID, "event-one")
	e2 := newTestEventWithID(userID, "event-two")
	e3 := newTestEventWithID(userID, "event-three")

	// Published before the client reconnects with Last-Event-ID = e1.ID.
	h.config.EventBroadcaster.Publish(userID, e1)
	h.config.EventBroadcaster.Publish(userID, e2)
	h.config.EventBroadcaster.Publish(userID, e3)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx)
	req.Header.Set("Last-Event-ID", e1.ID.String())

	writeCh := make(chan struct{}, 16)
	rec := newSSERecorder(writeCh)

	done := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec, req, nil, user, nil)
		close(done)
	}()

	waitForBodyContains(t, rec, writeCh, "event-three")

	body := rec.body.String()
	if strings.Contains(body, "event-one") {
		t.Fatalf("did not expect Last-Event-ID anchor event-one to be replayed; got: %q", body)
	}
	idx2 := strings.Index(body, "event-two")
	idx3 := strings.Index(body, "event-three")
	if idx2 == -1 || idx3 == -1 || idx2 >= idx3 {
		t.Fatalf("expected event-two before event-three in replay; got: %q", body)
	}
	if !strings.Contains(body, fmt.Sprintf("id: %s:", e2.ID.String())) ||
		!strings.Contains(body, fmt.Sprintf("id: %s:", e3.ID.String())) {
		t.Fatalf("replayed events missing id: framing; got: %q", body)
	}

	// Subsequent live event still streams after replay.
	e4 := newTestEventWithID(userID, "event-four")
	h.config.EventBroadcaster.Publish(userID, e4)
	waitForBodyContains(t, rec, writeCh, "event-four")

	body = rec.body.String()
	if strings.Index(body, "event-three") >= strings.Index(body, "event-four") {
		t.Fatalf("expected live event-four after replayed event-three; got: %q", body)
	}

	cancel()
	<-done
}

func TestSubscribeEventsHandler_LastEventIDEdgeCases(t *testing.T) {
	cases := []struct {
		name         string
		lastEventID  string
		expectReplay bool
	}{
		{name: "malformed", lastEventID: "not-a-uuid", expectReplay: false},
		{name: "unknown", lastEventID: uuid.Must(uuid.NewV4()).String(), expectReplay: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newEventsTestHandler(t)
			userID := uuid.Must(uuid.NewV4())
			user := &models.User{ID: userID}

			prior := newTestEventWithID(userID, "prior-event")
			h.config.EventBroadcaster.Publish(userID, prior)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			req := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx)
			req.Header.Set("Last-Event-ID", tc.lastEventID)

			writeCh := make(chan struct{}, 8)
			rec := newSSERecorder(writeCh)

			done := make(chan struct{})
			go func() {
				h.SubscribeEventsHandler(rec, req, nil, user, nil)
				close(done)
			}()

			live := newTestEventWithID(userID, "live-event")
			deadline := time.After(testTimeout)
			for !strings.Contains(rec.body.String(), "live-event") {
				select {
				case <-deadline:
					t.Fatalf("timed out waiting for live-event; got: %q", rec.body.String())
				default:
				}
				h.config.EventBroadcaster.Publish(userID, live)
				select {
				case <-writeCh:
				case <-time.After(20 * time.Millisecond):
				}
			}

			hasPrior := strings.Contains(rec.body.String(), "prior-event")
			if hasPrior != tc.expectReplay {
				t.Fatalf("prior-event replay = %v, want %v for %s Last-Event-ID; got: %q",
					hasPrior, tc.expectReplay, tc.name, rec.body.String())
			}

			cancel()
			<-done
		})
	}
}

func TestSubscribeEventsHandler_RepeatedEventIDsReplay(t *testing.T) {
	h := newEventsTestHandler(t)
	userID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: userID}

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	req1 := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx1)
	writeCh1 := make(chan struct{}, 8)
	rec1 := newSSERecorder(writeCh1)

	done1 := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec1, req1, nil, user, nil)
		close(done1)
	}()

	sharedID := uuid.Must(uuid.NewV4())
	firstA := newTestEventWithID(userID, "first-A")
	firstA.ID = sharedID

	deadline := time.After(testTimeout)
	for !strings.Contains(rec1.body.String(), "first-A") {
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for first-A; got: %q", rec1.body.String())
		default:
		}
		h.config.EventBroadcaster.Publish(userID, firstA)
		select {
		case <-writeCh1:
		case <-time.After(20 * time.Millisecond):
		}
	}

	var firstCursor string
	for _, frame := range strings.Split(rec1.body.String(), "\n\n") {
		if strings.Contains(frame, "first-A") {
			for _, line := range strings.Split(frame, "\n") {
				if strings.HasPrefix(line, "id: ") {
					firstCursor = strings.TrimPrefix(line, "id: ")
					break
				}
			}
			break
		}
	}
	if firstCursor == "" {
		t.Fatalf("failed to extract SSE id cursor from first-A frame: %q", rec1.body.String())
	}

	cancel1()
	<-done1

	// Publish B and a second A with the same UUID while the client is disconnected.
	eventB := newTestEventWithID(userID, "event-B")
	secondA := newTestEventWithID(userID, "second-A")
	secondA.ID = sharedID
	h.config.EventBroadcaster.Publish(userID, eventB)
	h.config.EventBroadcaster.Publish(userID, secondA)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	req2 := httptest.NewRequest(http.MethodGet, "/api/system/events/subscribe", nil).WithContext(ctx2)
	req2.Header.Set("Last-Event-ID", firstCursor)
	writeCh2 := make(chan struct{}, 8)
	rec2 := newSSERecorder(writeCh2)

	done2 := make(chan struct{})
	go func() {
		h.SubscribeEventsHandler(rec2, req2, nil, user, nil)
		close(done2)
	}()

	waitForBodyContains(t, rec2, writeCh2, "second-A")
	body := rec2.body.String()
	if strings.Contains(body, "first-A") {
		t.Fatalf("did not expect first-A to be replayed; got: %q", body)
	}
	idxB := strings.Index(body, "event-B")
	idxA2 := strings.Index(body, "second-A")
	if idxB == -1 || idxA2 == -1 || idxB >= idxA2 {
		t.Fatalf("expected event-B before second-A on replay; got: %q", body)
	}

	cancel2()
	<-done2
}
