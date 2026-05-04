package handlers

import (
	"bytes"
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshkit/logger"
	_events "github.com/meshery/meshkit/utils/events"
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
