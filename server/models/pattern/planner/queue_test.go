package planner

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestQueue_EmptyQueueBehavior(t *testing.T) {
	q := &Queue{}

	if q.Length() != 0 {
		t.Errorf("expected length 0, got %d", q.Length())
	}
	if got, ok := q.Peek(); ok || got != "" {
		t.Errorf("expected (\"\", false) from Peek on empty queue, got (%q, %v)", got, ok)
	}
	// This should NOT panic on an empty queue
	if got, ok := q.Dequeue(); ok || got != "" {
		t.Errorf("expected (\"\", false) from Dequeue on empty queue, got (%q, %v)", got, ok)
	}
}

func TestQueue_EnqueueEmptyString(t *testing.T) {
	q := &Queue{}

	q.Enqueue("")
	if q.Length() != 1 {
		t.Errorf("expected length 1 after Enqueue empty string, got %d", q.Length())
	}
	if got, ok := q.Peek(); !ok || got != "" {
		t.Errorf("expected (\"\", true) from Peek, got (%q, %v)", got, ok)
	}
	if got, ok := q.Dequeue(); !ok || got != "" {
		t.Errorf("expected (\"\", true) from Dequeue, got (%q, %v)", got, ok)
	}
	if q.Length() != 0 {
		t.Errorf("expected length 0 after Dequeue, got %d", q.Length())
	}
	if got, ok := q.Dequeue(); ok || got != "" {
		t.Errorf("expected (\"\", false) from Dequeue on drained queue, got (%q, %v)", got, ok)
	}
}

func TestQueue_FIFOOrder(t *testing.T) {
	q := &Queue{}

	q.Enqueue("a")
	q.Enqueue("b")
	q.Enqueue("c")

	if q.Length() != 3 {
		t.Errorf("expected length 3, got %d", q.Length())
	}
	if got, ok := q.Peek(); !ok || got != "a" {
		t.Errorf("expected Peek to return (\"a\", true), got (%q, %v)", got, ok)
	}

	if got, ok := q.Dequeue(); !ok || got != "a" {
		t.Errorf("expected Dequeue to return (\"a\", true), got (%q, %v)", got, ok)
	}
	if got, ok := q.Dequeue(); !ok || got != "b" {
		t.Errorf("expected Dequeue to return (\"b\", true), got (%q, %v)", got, ok)
	}
	if got, ok := q.Dequeue(); !ok || got != "c" {
		t.Errorf("expected Dequeue to return (\"c\", true), got (%q, %v)", got, ok)
	}

	if q.Length() != 0 {
		t.Errorf("expected length 0 after draining, got %d", q.Length())
	}
	// Dequeue on now-empty queue must not panic and must return ok=false
	if got, ok := q.Dequeue(); ok || got != "" {
		t.Errorf("expected (\"\", false) from Dequeue on empty queue after drain, got (%q, %v)", got, ok)
	}
}

func TestQueue_PeekDoesNotRemove(t *testing.T) {
	q := &Queue{}
	q.Enqueue("x")

	if got, ok := q.Peek(); !ok || got != "x" {
		t.Errorf("expected Peek to return (\"x\", true), got (%q, %v)", got, ok)
	}
	// Peek must not consume the item
	if q.Length() != 1 {
		t.Errorf("expected length 1 after Peek, got %d", q.Length())
	}
	if got, ok := q.Dequeue(); !ok || got != "x" {
		t.Errorf("expected Dequeue to return (\"x\", true), got (%q, %v)", got, ok)
	}
}

func TestQueue_ConcurrentAccess(t *testing.T) {
	q := &Queue{}
	const goroutines = 20
	const itemsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				q.Enqueue("item")
				_ = q.Length()
				_, _ = q.Peek()
			}
		}()
	}

	wg.Wait()
	if q.Length() != goroutines*itemsPerGoroutine {
		t.Errorf("expected %d items after concurrent enqueue, got %d", goroutines*itemsPerGoroutine, q.Length())
	}
}

func TestQueue_ConcurrentConsumers(t *testing.T) {
	const total = 1000
	const consumers = 10

	q := &Queue{}
	for i := 0; i < total; i++ {
		q.Enqueue(fmt.Sprintf("item-%d", i))
	}

	// Each consumer claims an item before dequeuing, so the number of Dequeue
	// calls is exact and ok confirms each item was dequeued successfully.
	var remaining atomic.Int64
	remaining.Store(total)

	received := make(chan string, total)

	var wg sync.WaitGroup
	wg.Add(consumers)

	for i := 0; i < consumers; i++ {
		go func() {
			defer wg.Done()
			for remaining.Add(-1) >= 0 {
				v, ok := q.Dequeue()
				if !ok {
					t.Errorf("Dequeue returned ok=false while items remained")
					return
				}
				received <- v
			}
		}()
	}

	wg.Wait()
	close(received)

	seen := make(map[string]int)
	for v := range received {
		seen[v]++
	}

	if got := len(seen); got != total {
		t.Errorf("expected %d distinct items, got %d", total, got)
	}
	for item, count := range seen {
		if count != 1 {
			t.Errorf("item %q received %d times, want exactly 1", item, count)
		}
	}
	if q.Length() != 0 {
		t.Errorf("expected empty queue, got length %d", q.Length())
	}
	if _, ok := q.Dequeue(); ok {
		t.Errorf("expected ok=false from Dequeue on drained queue")
	}
}

func TestQueue_ConcurrentProducersAndConsumers(t *testing.T) {
	const (
		numProducers     = 10
		itemsPerProducer = 100
		totalItems       = numProducers * itemsPerProducer
		numConsumers     = 10
	)

	q := &Queue{}
	received := make(chan string, totalItems)

	var producersWg sync.WaitGroup
	producersWg.Add(numProducers)

	for p := 0; p < numProducers; p++ {
		go func(producerID int) {
			defer producersWg.Done()
			for i := 0; i < itemsPerProducer; i++ {
				q.Enqueue(fmt.Sprintf("p%d-item%d", producerID, i))
			}
		}(p)
	}

	done := make(chan struct{})
	var readersWg sync.WaitGroup
	readersWg.Add(4)
	for r := 0; r < 4; r++ {
		go func() {
			defer readersWg.Done()
			for {
				select {
				case <-done:
					return
				default:
					_ = q.Length()
					_, _ = q.Peek()
				}
			}
		}()
	}

	var remaining atomic.Int64
	remaining.Store(totalItems)

	var consumersWg sync.WaitGroup
	consumersWg.Add(numConsumers)

	for c := 0; c < numConsumers; c++ {
		go func() {
			defer consumersWg.Done()
			for remaining.Add(-1) >= 0 {
				for {
					v, ok := q.Dequeue()
					if ok {
						received <- v
						break
					}
					time.Sleep(10 * time.Microsecond)
				}
			}
		}()
	}

	producersWg.Wait()
	consumersWg.Wait()
	close(done)
	readersWg.Wait()
	close(received)

	seen := make(map[string]int)
	for v := range received {
		seen[v]++
	}

	if len(seen) != totalItems {
		t.Errorf("expected %d distinct items, got %d", totalItems, len(seen))
	}
	for item, count := range seen {
		if count != 1 {
			t.Errorf("item %q consumed %d times, want 1", item, count)
		}
	}
	if q.Length() != 0 {
		t.Errorf("expected empty queue, got %d", q.Length())
	}
	if _, ok := q.Dequeue(); ok {
		t.Errorf("expected ok=false from Dequeue on drained queue")
	}
}
