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
