package models

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

// TestBroadcast_MultipleSubscribersSameID asserts that two Subscribe calls
// for the same id both receive the published event. This regresses the
// value-copy clobber bug where the second Subscribe could overwrite the
// first subscriber's listener slice entry.
func TestBroadcast_MultipleSubscribersSameID(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch1, unsub1 := b.Subscribe(id)
	ch2, unsub2 := b.Subscribe(id)
	defer unsub1()
	defer unsub2()

	b.Publish(id, "hello")

	for i, ch := range []chan interface{}{ch1, ch2} {
		select {
		case v := <-ch:
			if v != "hello" {
				t.Errorf("subscriber %d got %v, want hello", i, v)
			}
		case <-time.After(200 * time.Millisecond):
			t.Errorf("subscriber %d timed out waiting for event", i)
		}
	}
}

// TestBroadcast_DoubleUnsubscribe asserts that calling the returned
// unsubscribe function twice is a no-op rather than a panic on
// close-of-closed-channel.
func TestBroadcast_DoubleUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	_, unsub := b.Subscribe(id)
	unsub()

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("second unsubscribe panicked: %v", r)
		}
	}()
	unsub()
}

// TestBroadcast_ConcurrentPublishUnsubscribe runs many goroutines that
// subscribe, publish, and unsubscribe on the same id while a background
// publisher hammers the same id. Under -race this exercises both the
// listener slice serialization and the close-after-unlock ordering in
// unsubscribe.
func TestBroadcast_ConcurrentPublishUnsubscribe(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	const subscribers = 50
	const publishesPerSubscriber = 10

	var wg sync.WaitGroup
	for i := 0; i < subscribers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch, unsub := b.Subscribe(id)
			// Drain in the background so publishes do not stall on a full buffer.
			drainDone := make(chan struct{})
			go func() {
				for range ch {
				}
				close(drainDone)
			}()
			for j := 0; j < publishesPerSubscriber; j++ {
				b.Publish(id, j)
			}
			unsub()
			<-drainDone
		}()
	}

	var bgPublishes atomic.Int64
	stop := make(chan struct{})
	bgDone := make(chan struct{})
	go func() {
		defer close(bgDone)
		for {
			select {
			case <-stop:
				return
			default:
			}
			b.Publish(id, "bg")
			bgPublishes.Add(1)
		}
	}()

	wg.Wait()
	close(stop)
	<-bgDone

	t.Logf("background publishes: %d", bgPublishes.Load())
}

// TestBroadcast_BestEffortDeliveryWhenBufferFull locks in the documented
// behaviour of Publish: each subscriber's channel has buffer 1, and a
// second publish before the subscriber drains the first is dropped
// rather than blocking or panicking.
func TestBroadcast_BestEffortDeliveryWhenBufferFull(t *testing.T) {
	b := NewBroadcaster("test")
	id := uuid.Must(uuid.NewV4())

	ch, unsub := b.Subscribe(id)
	defer unsub()

	b.Publish(id, "first")
	b.Publish(id, "dropped")

	select {
	case v := <-ch:
		if v != "first" {
			t.Errorf("got %v, want first", v)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected first event in buffer")
	}

	select {
	case v := <-ch:
		t.Errorf("expected no second event, got %v", v)
	case <-time.After(50 * time.Millisecond):
	}
}
