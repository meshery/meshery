package models

import (
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

// TestBroadcastDeliversToSubscriber is the happy path: a published value reaches
// a live subscriber.
func TestBroadcastDeliversToSubscriber(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	b.Publish(id, 42)

	select {
	case got := <-ch:
		if got != 42 {
			t.Fatalf("got %v, want 42", got)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for published value")
	}
}

// TestBroadcastBurstWithinBufferNoLoss guards the data-loss regression: a burst
// of events published before the subscriber drains must all be delivered, in
// order. The old Publish used a destructive IsClosed check that consumed a
// buffered event and then skipped the send, losing both.
func TestBroadcastBurstWithinBufferNoLoss(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	const n = 10 // well within listenerBufferSize
	for i := 0; i < n; i++ {
		b.Publish(id, i)
	}

	for i := 0; i < n; i++ {
		select {
		case got := <-ch:
			if got != i {
				t.Fatalf("event %d: got %v, want %d (order/loss)", i, got, i)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out after %d/%d events; earlier events were lost", i, n)
		}
	}
}

// TestBroadcastNonBlockingWhenBufferFull guards the slow-consumer regression:
// Publish must never block, even when a subscriber never drains. Overflow beyond
// the buffer is dropped rather than blocking the publisher.
func TestBroadcastNonBlockingWhenBufferFull(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	defer unsubscribe()

	done := make(chan struct{})
	go func() {
		// Publish more than the buffer can hold, without anyone draining.
		for i := 0; i < listenerBufferSize+64; i++ {
			b.Publish(id, i)
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on a full subscriber buffer")
	}

	// The buffer retained exactly listenerBufferSize items; the rest were dropped.
	count := 0
	for {
		select {
		case <-ch:
			count++
		default:
			if count != listenerBufferSize {
				t.Fatalf("buffered %d items, want %d", count, listenerBufferSize)
			}
			return
		}
	}
}

// TestBroadcastConcurrentPublishUnsubscribeNoPanic guards the send-on-closed
// panic: publishers and churning subscribers race for the same id. Run with
// -race. A regression (unlocked send racing a close) panics and fails the test.
func TestBroadcastConcurrentPublishUnsubscribeNoPanic(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	var wg sync.WaitGroup
	stop := make(chan struct{})

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					b.Publish(id, 1)
				}
			}
		}()
	}

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ch, unsubscribe := b.Subscribe(id)
					select {
					case <-ch:
					default:
					}
					unsubscribe()
				}
			}
		}()
	}

	time.Sleep(200 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// TestBroadcastUnsubscribeIsIdempotent ensures a second unsubscribe is a no-op
// rather than a double-close panic.
func TestBroadcastUnsubscribeIsIdempotent(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	_, unsubscribe := b.Subscribe(id)
	unsubscribe()
	unsubscribe() // must not panic
}

// TestBroadcastUnsubscribeStopsDelivery ensures that after unsubscribe the
// listener channel is closed and receives no further events.
func TestBroadcastUnsubscribeStopsDelivery(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsubscribe := b.Subscribe(id)
	unsubscribe()

	// Publishing to an id whose only listener has unsubscribed must not panic
	// and must not deliver.
	b.Publish(id, 7)

	select {
	case v, open := <-ch:
		if open {
			t.Fatalf("received %v after unsubscribe; channel should be drained/closed", v)
		}
	case <-time.After(time.Second):
		t.Fatal("closed channel did not yield")
	}
}

// TestBroadcastUnsubscribePrunesEmptyClientsMap ensures that when all subscribers
// for a user unsubscribe, the user entry is deleted from the clients map to prevent
// an unbounded memory leak.
func TestBroadcastUnsubscribePrunesEmptyClientsMap(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	_, unsubscribe1 := b.Subscribe(id)
	_, unsubscribe2 := b.Subscribe(id)

	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("expected client entry to be present in map while subscribed")
	}

	unsubscribe1()
	// Still has one subscriber, entry must remain in map.
	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("client entry pruned prematurely while one subscriber remains")
	}

	unsubscribe2()
	// All subscribers gone, entry must be pruned from map.
	if _, ok := b.clients.Load(id); ok {
		t.Fatal("expected client entry to be pruned from map after all subscribers unsubscribed")
	}
}

// TestBroadcastResubscribeAfterUnsubscribe ensures that resubscribing after
// full eviction creates a fresh active entry and receives published data normally.
func TestBroadcastResubscribeAfterUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch1, unsub1 := b.Subscribe(id)
	unsub1()

	if _, ok := b.clients.Load(id); ok {
		t.Fatal("entry should be pruned after unsub")
	}

	ch2, unsub2 := b.Subscribe(id)
	defer unsub2()

	if _, ok := b.clients.Load(id); !ok {
		t.Fatal("entry should be present after resubscribe")
	}

	b.Publish(id, "resubscribed-event")

	select {
	case msg := <-ch2:
		if msg != "resubscribed-event" {
			t.Fatalf("got %v, want resubscribed-event", msg)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on resubscribed listener")
	}

	select {
	case _, open := <-ch1:
		if open {
			t.Fatal("old channel should be closed")
		}
	default:
	}
}

// TestBroadcastConcurrentPruneAndSubscribe hammers Subscribe, Publish, and Unsubscribe
// on the same ID to verify that atomic eviction does not race with new subscribers
// or panic with -race.
func TestBroadcastConcurrentPruneAndSubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Publishers
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					b.Publish(id, val)
				}
			}
		}(i)
	}

	// Churning subscribers triggering eviction repeatedly
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ch, unsubscribe := b.Subscribe(id)
					select {
					case <-ch:
					default:
					}
					unsubscribe()
				}
			}
		}()
	}

	time.Sleep(300 * time.Millisecond)
	close(stop)
	wg.Wait()

	// After all churn finishes, resubscribing must work cleanly
	ch, unsub := b.Subscribe(id)
	defer unsub()

	b.Publish(id, 999)
	select {
	case v := <-ch:
		if v != 999 {
			t.Fatalf("got %v, want 999", v)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for event on final subscriber")
	}
}
