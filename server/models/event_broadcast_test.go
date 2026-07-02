package models

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid"
)

// newBroadcastTestID returns a fresh random id for broadcaster tests.
// core.Uuid is an alias for uuid.UUID, so this value is usable directly as the
// id argument to Subscribe/Publish.
func newBroadcastTestID(t *testing.T) uuid.UUID {
	t.Helper()
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}
	return id
}

// TestBroadcast_PublishDeliversToSubscriber verifies the basic contract: an
// event published for an id is delivered to that id's subscriber.
func TestBroadcast_PublishDeliversToSubscriber(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	b.Publish(id, "hello")

	select {
	case got := <-ch:
		if got != "hello" {
			t.Fatalf("got %v, want %q", got, "hello")
		}
	default:
		t.Fatal("expected the published event to be buffered, channel was empty")
	}
}

// TestBroadcast_BufferedEventNotDrained guards against the regression where the
// closed-channel check consumed the buffered event: the first event must remain
// queued for the subscriber. Because the buffer has capacity 1, the second
// event is intentionally dropped (back-pressure) rather than replacing the
// first or being silently drained out-of-band.
func TestBroadcast_BufferedEventNotDrained(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	b.Publish(id, "first")  // buffered (capacity 1)
	b.Publish(id, "second") // buffer full -> dropped, must not drain "first"

	got := <-ch
	if got != "first" {
		t.Fatalf("buffered event was lost or replaced: got %v, want %q", got, "first")
	}
}

// TestBroadcast_UnsubscribeClosesChannelAndPublishIsSafe verifies unsubscribe
// closes the listener channel and that publishing afterwards (with no
// listeners) does not panic.
func TestBroadcast_UnsubscribeClosesChannelAndPublishIsSafe(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	ch, unsubscribe := b.Subscribe(id)
	unsubscribe()

	if _, open := <-ch; open {
		t.Fatal("expected channel to be closed after unsubscribe")
	}

	b.Publish(id, "after-unsub") // must be a no-op, not a panic
}

// TestBroadcast_DoubleUnsubscribeDoesNotPanic verifies unsubscribe is
// idempotent and never double-closes the channel.
func TestBroadcast_DoubleUnsubscribeDoesNotPanic(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	_, unsubscribe := b.Subscribe(id)
	unsubscribe()
	unsubscribe() // must be a no-op
}

// TestBroadcast_ConcurrentPublishAndUnsubscribe drives publishers and
// subscribing/unsubscribing goroutines against the same id simultaneously.
// Before the fix this panicked with "send on closed channel" and tripped the
// race detector while iterating the listeners slice. It must now run clean
// under `go test -race`.
func TestBroadcast_ConcurrentPublishAndUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	const (
		publishers  = 8
		subscribers = 8
		iterations  = 500
	)

	var wg sync.WaitGroup

	for p := 0; p < publishers; p++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				b.Publish(id, i) // must never panic sending to a closed channel
			}
		}()
	}

	for s := 0; s < subscribers; s++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				ch, unsubscribe := b.Subscribe(id)
				// Drain whatever happens to be buffered without blocking, then
				// unsubscribe (closing the channel) concurrently with publishers.
				select {
				case <-ch:
				default:
				}
				unsubscribe()
			}
		}()
	}

	wg.Wait()

	// Every Subscribe was paired with an unsubscribe, so no listeners should
	// remain leaked in the slice.
	if actual, ok := b.clients.Load(id); ok {
		cl := actual.(*clients)
		cl.mu.Lock()
		n := len(cl.listeners)
		cl.mu.Unlock()
		if n != 0 {
			t.Fatalf("expected no listeners after all unsubscribes, got %d", n)
		}
	}

	// A final publish after all the churn must still be safe.
	b.Publish(id, "final")
}
