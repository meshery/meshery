package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	_events "github.com/meshery/meshkit/utils/events"
	"github.com/meshery/schemas/models/core"
)

const testTimeout = time.Second

type testFlusher struct {
	flushes atomic.Int32
	flushCh chan struct{}
}

func (f *testFlusher) Flush() {
	f.flushes.Add(1)
	if f.flushCh != nil {
		select {
		case f.flushCh <- struct{}{}:
		default:
		}
	}
}

type safeBuffer struct {
	mu      sync.Mutex
	buf     bytes.Buffer
	writeCh chan struct{}
}

func (b *safeBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	n, err := b.buf.Write(p)
	if b.writeCh != nil {
		select {
		case b.writeCh <- struct{}{}:
		default:
		}
	}

	return n, err
}

func (b *safeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.buf.String()
}

func waitForSignal(t *testing.T, ch <-chan struct{}, name string) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(testTimeout):
		t.Fatalf("timed out waiting for %s", name)
	}
}

func newTestLogger(t *testing.T) logger.Handler {
	t.Helper()

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	return log
}

func waitForPayload(t *testing.T, ch <-chan []byte, name string) []byte {
	t.Helper()

	select {
	case payload := <-ch:
		return payload
	case <-time.After(testTimeout):
		t.Fatalf("timed out waiting for %s", name)
	}

	panic("unreachable")
}

func assertNoPayload(t *testing.T, ch <-chan []byte, name string) {
	t.Helper()

	select {
	case payload := <-ch:
		t.Fatalf("unexpected payload for %s: %s", name, payload)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestSendStreamEvent(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() (context.Context, context.CancelFunc)
		setupChan   func() chan []byte
		expectSent  bool
		expectValue []byte
	}{
		{
			name: "sends data when receiver is available",
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			setupChan: func() chan []byte {
				return make(chan []byte, 1)
			},
			expectSent:  true,
			expectValue: []byte("payload"),
		},
		{
			name: "returns false when context is cancelled",
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, func() {}
			},
			setupChan: func() chan []byte {
				return make(chan []byte)
			},
			expectSent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.setupCtx()
			defer cancel()

			respChan := tt.setupChan()
			sent := sendStreamEvent(ctx, respChan, []byte("payload"))
			if sent != tt.expectSent {
				t.Fatalf("expected sent=%v, got %v", tt.expectSent, sent)
			}

			if !tt.expectSent {
				return
			}

			select {
			case got := <-respChan:
				if string(got) != string(tt.expectValue) {
					t.Fatalf("expected %q, got %q", tt.expectValue, got)
				}
			case <-time.After(testTimeout):
				t.Fatal("timed out waiting for streamed payload")
			}
		})
	}
}

func TestWriteEventStream_StopsOnContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := newTestLogger(t)
	respChan := make(chan []byte, 1)
	flusher := &testFlusher{flushCh: make(chan struct{}, 1)}
	body := &safeBuffer{writeCh: make(chan struct{}, 1)}
	done := make(chan struct{})

	go func() {
		writeEventStream(ctx, body, respChan, log, flusher)
		close(done)
	}()

	respChan <- []byte(`{"status":"ok"}`)

	waitForSignal(t, body.writeCh, "stream write")
	waitForSignal(t, flusher.flushCh, "flush")

	if body.String() != "data: {\"status\":\"ok\"}\n\n" {
		t.Fatalf("unexpected stream output: %q", body.String())
	}

	if flusher.flushes.Load() != 1 {
		t.Fatalf("expected flusher to be called once, got %d", flusher.flushes.Load())
	}

	cancel()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("writeEventStream did not stop after context cancellation")
	}
}

func TestListenForCoreEvents_StopsBlockedSendOnCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	respChan := make(chan []byte, 1)
	respChan <- []byte("pre-filled to force a blocked send")
	done := make(chan struct{})
	subscribed := make(chan struct{}, 1)
	subscribe := func(eb *_events.EventStreamer, ch chan interface{}) {
		defaultSubscribeToEventStream(eb, ch)
		subscribed <- struct{}{}
	}

	go func() {
		listenForCoreEvents(ctx, eb, respChan, log, nil, subscribe)
		close(done)
	}()

	waitForSignal(t, subscribed, "event streamer subscription")
	eb.Publish(&meshes.EventsResponse{Summary: "stream event"})
	cancel()

	select {
	case <-done:
	case <-time.After(2 * testTimeout):
		t.Fatal("listenForCoreEvents remained blocked after cancellation")
	}
}

func TestListenForCoreEvents_SingleUserReceivesAndStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	respChan := make(chan []byte, 1)
	done := make(chan struct{})
	subscribed := make(chan struct{}, 1)
	subscribe := func(eb *_events.EventStreamer, ch chan interface{}) {
		defaultSubscribeToEventStream(eb, ch)
		subscribed <- struct{}{}
	}

	go func() {
		listenForCoreEvents(ctx, eb, respChan, log, nil, subscribe)
		close(done)
	}()

	waitForSignal(t, subscribed, "single-user subscription")
	eb.Publish(&meshes.EventsResponse{Summary: "single-user event"})

	payload := waitForPayload(t, respChan, "single-user event")
	if !bytes.Contains(payload, []byte("single-user event")) {
		t.Fatalf("expected streamed payload to contain event summary, got %s", payload)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("single-user listener did not stop after cancellation")
	}
}

func TestListenForCoreEvents_MultiUserBroadcast(t *testing.T) {
	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	respChans := []chan []byte{make(chan []byte, 1), make(chan []byte, 1)}
	done := []chan struct{}{make(chan struct{}), make(chan struct{})}
	subscribed := make(chan struct{}, len(respChans))
	cancels := make([]context.CancelFunc, 0, len(respChans))

	for i := range respChans {
		ctx, cancel := context.WithCancel(context.Background())
		cancels = append(cancels, cancel)

		go func(ctx context.Context, respChan chan []byte, done chan struct{}) {
			listenForCoreEvents(ctx, eb, respChan, log, nil, func(eb *_events.EventStreamer, ch chan interface{}) {
				defaultSubscribeToEventStream(eb, ch)
				subscribed <- struct{}{}
			})
			close(done)
		}(ctx, respChans[i], done[i])
	}

	for i := 0; i < len(respChans); i++ {
		waitForSignal(t, subscribed, "multi-user subscription")
	}

	eb.Publish(&meshes.EventsResponse{Summary: "broadcast event"})

	for i, respChan := range respChans {
		payload := waitForPayload(t, respChan, "multi-user event")
		if !bytes.Contains(payload, []byte("broadcast event")) {
			t.Fatalf("listener %d did not receive broadcast payload: %s", i, payload)
		}
	}

	for _, cancel := range cancels {
		cancel()
	}

	for i, doneChan := range done {
		select {
		case <-doneChan:
		case <-time.After(testTimeout):
			t.Fatalf("listener %d did not stop after cancellation", i)
		}
	}
}

func TestListenForCoreEvents_MultiUserBlockedSendsExitOnCancellation(t *testing.T) {
	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	const listeners = 3
	subscribed := make(chan struct{}, listeners)
	done := make([]chan struct{}, 0, listeners)
	cancels := make([]context.CancelFunc, 0, listeners)

	for range listeners {
		ctx, cancel := context.WithCancel(context.Background())
		respChan := make(chan []byte, 1)
		respChan <- []byte("pre-filled")
		doneChan := make(chan struct{})
		done = append(done, doneChan)
		cancels = append(cancels, cancel)

		go func(ctx context.Context, respChan chan []byte, done chan struct{}) {
			listenForCoreEvents(ctx, eb, respChan, log, nil, func(eb *_events.EventStreamer, ch chan interface{}) {
				defaultSubscribeToEventStream(eb, ch)
				subscribed <- struct{}{}
			})
			close(done)
		}(ctx, respChan, doneChan)
	}

	for range listeners {
		waitForSignal(t, subscribed, "blocked-send listener subscription")
	}

	eb.Publish(&meshes.EventsResponse{Summary: "blocked event"})

	for _, cancel := range cancels {
		cancel()
	}

	for i, doneChan := range done {
		select {
		case <-doneChan:
		case <-time.After(2 * testTimeout):
			t.Fatalf("blocked listener %d remained stuck after cancellation", i)
		}
	}
}

func TestListenForCoreEvents_UnsubscribesOnExit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	respChan := make(chan []byte, 1)
	done := make(chan struct{})
	subscribed := make(chan struct{}, 1)
	var subscribedChannel chan interface{}
	subscribe := func(eb *_events.EventStreamer, ch chan interface{}) {
		subscribedChannel = ch
		defaultSubscribeToEventStream(eb, ch)
		subscribed <- struct{}{}
	}

	go func() {
		listenForCoreEvents(ctx, eb, respChan, log, nil, subscribe)
		close(done)
	}()

	waitForSignal(t, subscribed, "unsubscribe test subscription")
	cancel()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("listener did not stop after cancellation")
	}

	eb.Publish(&meshes.EventsResponse{Summary: "post-cancel event"})

	select {
	case payload := <-subscribedChannel:
		t.Fatalf("unsubscribed channel unexpectedly received payload: %#v", payload)
	case <-time.After(100 * time.Millisecond):
	}
}

func TestListenForCoreEvents_IgnoresUnexpectedPayloads(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := newTestLogger(t)
	eb := _events.NewEventStreamer()
	respChan := make(chan []byte, 1)
	done := make(chan struct{})
	subscribed := make(chan struct{}, 1)

	go func() {
		listenForCoreEvents(ctx, eb, respChan, log, nil, func(eb *_events.EventStreamer, ch chan interface{}) {
			defaultSubscribeToEventStream(eb, ch)
			subscribed <- struct{}{}
		})
		close(done)
	}()

	waitForSignal(t, subscribed, "edge-case subscription")
	eb.Publish("not a mesh event")
	assertNoPayload(t, respChan, "unexpected payload")

	eb.Publish(&meshes.EventsResponse{Summary: "valid event"})
	payload := waitForPayload(t, respChan, "valid edge-case event")
	if !bytes.Contains(payload, []byte("valid event")) {
		t.Fatalf("expected valid event payload after ignored message, got %s", payload)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(testTimeout):
		t.Fatal("edge-case listener did not stop after cancellation")
	}
}

// TestEventStreamHandler_FlushesHeadersBeforeFirstEvent verifies that the stream
// opens right away. A client doesn't see the response until the headers are
// sent, so without the flush mesheryctl adapter would wait with no timer running
// until an adapter sent its first event.
func TestEventStreamHandler_FlushesHeadersBeforeFirstEvent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	h := &Handler{log: newTestLogger(t), EventsBuffer: _events.NewEventStreamer()}
	req := httptest.NewRequest(http.MethodGet, "/api/system/adapter/events", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.EventStreamHandler(rec, req, &models.Preference{}, nil, nil)

	if !rec.Flushed {
		t.Error("expected the response headers to be flushed before any event")
	}
	if got := rec.Header().Get("Content-Type"); got != "text/event-stream" {
		t.Errorf("expected content type text/event-stream, got %q", got)
	}
}

// eventsSpyProvider mocks models.Provider to inspect event handling invocations and record arguments.
type eventsSpyProvider struct {
	*models.DefaultLocalProvider
	getEventsCalled          atomic.Bool
	getEventTypesCalled      atomic.Bool
	updateEventStatusCalled  atomic.Bool
	bulkUpdateCalled         atomic.Bool
	bulkDeleteCalled         atomic.Bool
	deleteEventCalled        atomic.Bool
	persistClientEventCalled atomic.Bool

	lastUserID core.Uuid
	lastSysID  core.Uuid
}

// newEventsSpyProvider instantiates a new eventsSpyProvider with an embedded DefaultLocalProvider.
func newEventsSpyProvider() *eventsSpyProvider {
	base := &models.DefaultLocalProvider{}
	return &eventsSpyProvider{DefaultLocalProvider: base}
}

// GetEvents mocks provider.GetEvents, capturing the user and system ID arguments.
func (s *eventsSpyProvider) GetEvents(token string, filter *events.EventsFilter, page int, userID core.Uuid, sysID core.Uuid) (*models.EventsResponse, error) {
	s.getEventsCalled.Store(true)
	s.lastUserID = userID
	s.lastSysID = sysID
	return &models.EventsResponse{Events: []*events.Event{}}, nil
}

// GetEventTypes mocks provider.GetEventTypes, returning an EventTypesResponse with category1.
func (s *eventsSpyProvider) GetEventTypes(token string, userID core.Uuid, sysID core.Uuid) (models.EventTypesResponse, error) {
	s.getEventTypesCalled.Store(true)
	s.lastUserID = userID
	s.lastSysID = sysID
	return models.EventTypesResponse{Category: []string{"category1"}}, nil
}

// UpdateEventStatus mocks provider.UpdateEventStatus, recording the call.
func (s *eventsSpyProvider) UpdateEventStatus(token string, eventID core.Uuid, status string) error {
	s.updateEventStatusCalled.Store(true)
	return nil
}

// BulkUpdateEventStatus mocks provider.BulkUpdateEventStatus, recording the call.
func (s *eventsSpyProvider) BulkUpdateEventStatus(token string, eventIDs []*core.Uuid, status string) error {
	s.bulkUpdateCalled.Store(true)
	return nil
}

// BulkDeleteEvent mocks provider.BulkDeleteEvent, recording the call.
func (s *eventsSpyProvider) BulkDeleteEvent(token string, eventIDs []*core.Uuid) error {
	s.bulkDeleteCalled.Store(true)
	return nil
}

// DeleteEvent mocks provider.DeleteEvent, recording the call.
func (s *eventsSpyProvider) DeleteEvent(token string, eventID core.Uuid) error {
	s.deleteEventCalled.Store(true)
	return nil
}

// GetProviderToken returns a mock provider token.
func (s *eventsSpyProvider) GetProviderToken(req *http.Request) (string, error) {
	return "test-token", nil
}

// PersistClientEvent mocks provider.PersistClientEvent, recording the call and returning a persisted event.
func (s *eventsSpyProvider) PersistClientEvent(event events.Event, token string) (models.PersistedEvent, error) {
	s.persistClientEventCalled.Store(true)
	return models.PersistedEvent{
		Event:       &event,
		Disposition: "created",
	}, nil
}

// TestGetAllEvents_NilUser verifies that GetAllEvents returns 401 Unauthorized when user is nil.
func TestGetAllEvents_NilUser(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()
	req := httptest.NewRequest(http.MethodGet, "/api/system/events", nil)
	rec := httptest.NewRecorder()

	h.GetAllEvents(rec, req, nil, nil, spy)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	if spy.getEventsCalled.Load() {
		t.Fatal("provider.GetEvents should not be called when user is nil")
	}
}

// TestGetAllEvents_NilProvider verifies that GetAllEvents returns 500 when provider is nil.
func TestGetAllEvents_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodGet, "/api/system/events", nil)
	rec := httptest.NewRecorder()

	h.GetAllEvents(rec, req, nil, user, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
}

// TestGetAllEvents_NilSystemID verifies that GetAllEvents safely defaults to uuid.Nil when h.SystemID is nil.
func TestGetAllEvents_NilSystemID(t *testing.T) {
	h := &Handler{log: newTestLogger(t), SystemID: nil}
	spy := newEventsSpyProvider()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodGet, "/api/system/events", nil)
	rec := httptest.NewRecorder()

	h.GetAllEvents(rec, req, nil, user, spy)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if !spy.getEventsCalled.Load() {
		t.Fatal("provider.GetEvents should be called")
	}
	if spy.lastSysID != uuid.Nil {
		t.Fatalf("expected sysID to be uuid.Nil when h.SystemID is nil, got %v", spy.lastSysID)
	}
	if spy.lastUserID != user.ID {
		t.Fatalf("expected userID %v, got %v", user.ID, spy.lastUserID)
	}
}

// TestGetEventTypes_NilUser verifies that GetEventTypes returns 401 Unauthorized when user is nil.
func TestGetEventTypes_NilUser(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/types", nil)
	rec := httptest.NewRecorder()

	h.GetEventTypes(rec, req, nil, nil, spy)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	if spy.getEventTypesCalled.Load() {
		t.Fatal("provider.GetEventTypes should not be called when user is nil")
	}
}

// TestGetEventTypes_NilProvider verifies that GetEventTypes returns 500 when provider is nil.
func TestGetEventTypes_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/types", nil)
	rec := httptest.NewRecorder()

	h.GetEventTypes(rec, req, nil, user, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
}

// TestGetEventTypes_NilSystemID verifies that GetEventTypes safely defaults to uuid.Nil when h.SystemID is nil.
func TestGetEventTypes_NilSystemID(t *testing.T) {
	h := &Handler{log: newTestLogger(t), SystemID: nil}
	spy := newEventsSpyProvider()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodGet, "/api/system/events/types", nil)
	rec := httptest.NewRecorder()

	h.GetEventTypes(rec, req, nil, user, spy)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if !spy.getEventTypesCalled.Load() {
		t.Fatal("provider.GetEventTypes should be called")
	}
	if spy.lastSysID != uuid.Nil {
		t.Fatalf("expected sysID to be uuid.Nil when h.SystemID is nil, got %v", spy.lastSysID)
	}
}

// TestUpdateEventStatus_NilProvider verifies that UpdateEventStatus returns 500 when provider is nil.
func TestUpdateEventStatus_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/123", nil)
	rec := httptest.NewRecorder()

	h.UpdateEventStatus(rec, req, nil, nil, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for nil provider, got %d", rec.Code)
	}
}

// TestUpdateEventStatus_MalformedJSON verifies that UpdateEventStatus returns 400 Bad Request on invalid JSON.
func TestUpdateEventStatus_MalformedJSON(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/123", bytes.NewBufferString("not valid json"))
	rec := httptest.NewRecorder()

	h.UpdateEventStatus(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed json, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if spy.updateEventStatusCalled.Load() {
		t.Fatal("provider.UpdateEventStatus should not be called on malformed json")
	}
}

// TestUpdateEventStatus_MissingStatus verifies that UpdateEventStatus returns 400 Bad Request on empty status.
func TestUpdateEventStatus_MissingStatus(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/123", bytes.NewBufferString(`{"status":""}`))
	rec := httptest.NewRecorder()

	h.UpdateEventStatus(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty status, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if spy.updateEventStatusCalled.Load() {
		t.Fatal("provider.UpdateEventStatus should not be called on empty status")
	}
}

// TestBulkUpdateEventStatus_NilProvider verifies that BulkUpdateEventStatus returns 500 when provider is nil.
func TestBulkUpdateEventStatus_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/bulk", nil)
	rec := httptest.NewRecorder()

	h.BulkUpdateEventStatus(rec, req, nil, nil, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for nil provider, got %d", rec.Code)
	}
}

// TestBulkUpdateEventStatus_EmptyOrMalformed verifies that BulkUpdateEventStatus returns 400 on malformed or empty payloads.
func TestBulkUpdateEventStatus_EmptyOrMalformed(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()

	// Malformed JSON
	req := httptest.NewRequest(http.MethodPut, "/api/system/events/status/bulk", bytes.NewBufferString("not json"))
	rec := httptest.NewRecorder()
	h.BulkUpdateEventStatus(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed JSON, got %d", rec.Code)
	}

	// Empty IDs
	req = httptest.NewRequest(http.MethodPut, "/api/system/events/status/bulk", bytes.NewBufferString(`{"status":"read","ids":[]}`))
	rec = httptest.NewRecorder()
	h.BulkUpdateEventStatus(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty IDs, got %d", rec.Code)
	}

	// Empty Status with non-empty IDs
	validID := uuid.Must(uuid.NewV4()).String()
	req = httptest.NewRequest(http.MethodPut, "/api/system/events/status/bulk", bytes.NewBufferString(fmt.Sprintf(`{"status":"","ids":["%s"]}`, validID)))
	rec = httptest.NewRecorder()
	h.BulkUpdateEventStatus(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty status, got %d", rec.Code)
	}
	if spy.bulkUpdateCalled.Load() {
		t.Fatal("provider.BulkUpdateEventStatus should not be called on invalid input")
	}
}

// TestBulkDeleteEvent_NilProvider verifies that BulkDeleteEvent returns 500 when provider is nil.
func TestBulkDeleteEvent_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodDelete, "/api/system/events/bulk", nil)
	rec := httptest.NewRecorder()

	h.BulkDeleteEvent(rec, req, nil, nil, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for nil provider, got %d", rec.Code)
	}
}

// TestBulkDeleteEvent_EmptyOrMalformed verifies that BulkDeleteEvent returns 400 on malformed or empty payloads.
func TestBulkDeleteEvent_EmptyOrMalformed(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()

	// Malformed JSON
	req := httptest.NewRequest(http.MethodDelete, "/api/system/events/bulk", bytes.NewBufferString("{bad"))
	rec := httptest.NewRecorder()
	h.BulkDeleteEvent(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for malformed JSON, got %d", rec.Code)
	}

	// Empty IDs
	req = httptest.NewRequest(http.MethodDelete, "/api/system/events/bulk", bytes.NewBufferString(`{"ids":[]}`))
	rec = httptest.NewRecorder()
	h.BulkDeleteEvent(rec, req, nil, nil, spy)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty IDs, got %d", rec.Code)
	}
	if spy.bulkDeleteCalled.Load() {
		t.Fatal("provider.BulkDeleteEvent should not be called on invalid input")
	}
}

// TestDeleteEvent_NilProvider verifies that DeleteEvent returns 500 when provider is nil.
func TestDeleteEvent_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	req := httptest.NewRequest(http.MethodDelete, "/api/system/events/123", nil)
	rec := httptest.NewRecorder()

	h.DeleteEvent(rec, req, nil, nil, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for nil provider, got %d", rec.Code)
	}
}

// TestClientEventHandler_NilUser verifies that ClientEventHandler returns 401 Unauthorized when user is nil.
func TestClientEventHandler_NilUser(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	spy := newEventsSpyProvider()
	req := httptest.NewRequest(http.MethodPost, "/api/system/events", nil)
	rec := httptest.NewRecorder()

	h.ClientEventHandler(rec, req, nil, nil, spy)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
}

// TestClientEventHandler_NilProvider verifies that ClientEventHandler returns 500 when provider is nil.
func TestClientEventHandler_NilProvider(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	req := httptest.NewRequest(http.MethodPost, "/api/system/events", nil)
	rec := httptest.NewRecorder()

	h.ClientEventHandler(rec, req, nil, user, nil)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
}

// TestClientEventHandler_SuccessWithNilSystemIDAndNilBroadcaster verifies that ClientEventHandler safely handles nil SystemID and nil EventBroadcaster.
func TestClientEventHandler_SuccessWithNilSystemIDAndNilBroadcaster(t *testing.T) {
	h := &Handler{
		log:      newTestLogger(t),
		SystemID: nil,
		config:   &models.HandlerConfig{EventBroadcaster: nil},
	}
	spy := newEventsSpyProvider()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	targetID := uuid.Must(uuid.NewV4())

	validEventJSON := fmt.Sprintf(`{
		"actedUpon": "%s",
		"action": "view",
		"category": "system",
		"severity": "informational",
		"description": "test event"
	}`, targetID.String())

	req := httptest.NewRequest(http.MethodPost, "/api/system/events", strings.NewReader(validEventJSON))
	rec := httptest.NewRecorder()

	h.ClientEventHandler(rec, req, nil, user, spy)
	if rec.Code != http.StatusOK && rec.Code != http.StatusCreated {
		t.Fatalf("expected 200/201, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	if !spy.persistClientEventCalled.Load() {
		t.Fatal("provider.PersistClientEvent should be called")
	}
}

type nonFlusherResponseWriter struct {
	header     http.Header
	statusCode int
	body       bytes.Buffer
}

func (w *nonFlusherResponseWriter) Header() http.Header {
	return w.header
}

func (w *nonFlusherResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *nonFlusherResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

type mockEventsProvider struct {
	models.DefaultLocalProvider
}

func TestCloseAdapterConnections_SafeAndIdempotent(t *testing.T) {
	lock := &sync.Mutex{}
	adapters := map[string]*meshes.MeshClient{
		"adapter-a:10000": {},
		"adapter-b:10001": {},
	}

	cleared := closeAdapterConnections(lock, adapters)
	if len(cleared) != 0 {
		t.Fatalf("expected cleared map to be empty, got %d", len(cleared))
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = closeAdapterConnections(lock, adapters)
		}()
	}
	wg.Wait()
}

func TestEventStreamHandler_UnsupportedFlusher(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	rw := &nonFlusherResponseWriter{header: make(http.Header)}
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/events", nil)

	prefObj := &models.Preference{}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	mockP := &mockEventsProvider{}

	h.EventStreamHandler(rw, req, prefObj, user, mockP)

	if rw.statusCode != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, rw.statusCode)
	}
}

func TestEventStreamHandler_ImmediateClientDisconnection(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	sysID := uuid.Must(uuid.NewV4())
	h.SystemID = &sysID
	h.EventsBuffer = _events.NewEventStreamer()

	rec := httptest.NewRecorder()
	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/api/events", nil)

	prefObj := &models.Preference{}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	mockP := &mockEventsProvider{}

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	done := make(chan struct{})
	go func() {
		h.EventStreamHandler(rec, req, prefObj, user, mockP)
		close(done)
	}()

	select {
	case <-done:
		// Clean exit without deadlock or panic!
	case <-time.After(2 * time.Second):
		t.Fatal("EventStreamHandler failed to exit cleanly upon client disconnection")
	}

	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", rec.Header().Get("Content-Type"))
	}
}

