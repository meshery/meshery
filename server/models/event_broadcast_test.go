package models

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
)

// newBroadcastTestID returns a fresh random id for broadcaster tests. It
// returns core.Uuid (converting explicitly) so callers match the
// Subscribe/Publish signature regardless of whether core.Uuid stays a type
// alias for uuid.UUID.
func newBroadcastTestID(t *testing.T) core.Uuid {
	t.Helper()
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate uuid: %v", err)
	}
	return core.Uuid(id)
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

	// Every Subscribe was paired with an unsubscribe, so the map entry must be
	// gone rather than left behind as an empty (leaked) entry.
	if _, ok := b.clients.Load(id); ok {
		t.Fatalf("expected the map entry to be removed after all unsubscribes")
	}

	// A final publish after all the churn must still be safe.
	b.Publish(id, "final")
}

// TestBroadcast_MapEntryRemovedAfterLastUnsubscribe verifies the map entry for
// an id is reclaimed once its last listener unsubscribes, preventing unbounded
// growth as unique ids come and go.
func TestBroadcast_MapEntryRemovedAfterLastUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	_, unsubscribe := b.Subscribe(id)
	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("expected a map entry to exist while subscribed")
	}

	unsubscribe()
	if _, ok := b.clients.Load(id); ok {
		t.Fatal("expected the map entry to be removed after the last unsubscribe")
	}

	// Subscribing again after cleanup must still work and deliver events.
	ch, unsubscribe2 := b.Subscribe(id)
	defer unsubscribe2()
	b.Publish(id, "again")
	if got := <-ch; got != "again" {
		t.Fatalf("got %v, want %q", got, "again")
	}
}

// TestBroadcast_MapEntryRetainedWhileSubscribersRemain verifies the entry is
// only removed once the *last* listener leaves, not on the first unsubscribe.
func TestBroadcast_MapEntryRetainedWhileSubscribersRemain(t *testing.T) {
	b := NewBroadcaster("test")
	id := newBroadcastTestID(t)

	_, unsubscribe1 := b.Subscribe(id)
	ch2, unsubscribe2 := b.Subscribe(id)

	unsubscribe1()
	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("expected the map entry to remain while a listener is still subscribed")
	}

	// The surviving listener must still receive events.
	b.Publish(id, "still-here")
	if got := <-ch2; got != "still-here" {
		t.Fatalf("got %v, want %q", got, "still-here")
	}

	unsubscribe2()
	if _, ok := b.clients.Load(id); ok {
		t.Fatal("expected the map entry to be removed after the last unsubscribe")
	}
}
