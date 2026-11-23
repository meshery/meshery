package planner

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestQueue_BasicOperations(t *testing.T) {
	var q Queue

	if length := q.Length(); length != 0 {
		t.Fatalf("expected initial length 0, got %d", length)
	}

	if peek := q.Peek(); peek != "" {
		t.Fatalf("expected empty peek on empty queue, got %q", peek)
	}

	// Dequeue on empty queue must return empty string safely without panicking
	if dequeued := q.Dequeue(); dequeued != "" {
		t.Fatalf("expected empty string on dequeuing empty queue, got %q", dequeued)
	}

	q.Enqueue("node-1")
	q.Enqueue("node-2")
	q.Enqueue("node-3")

	if length := q.Length(); length != 3 {
		t.Fatalf("expected length 3, got %d", length)
	}

	if peek := q.Peek(); peek != "node-1" {
		t.Fatalf("expected peek 'node-1', got %q", peek)
	}

	// Verify FIFO order
	expected := []string{"node-1", "node-2", "node-3"}
	for _, exp := range expected {
		val := q.Dequeue()
		if val != exp {
			t.Fatalf("expected dequeue %q, got %q", exp, val)
		}
	}

	if length := q.Length(); length != 0 {
		t.Fatalf("expected length 0 after dequeuing all, got %d", length)
	}

	// Dequeue once again on now-empty queue must not panic
	if dequeued := q.Dequeue(); dequeued != "" {
		t.Fatalf("expected empty string on dequeuing depleted queue, got %q", dequeued)
	}
}

func TestQueue_ConcurrentAccess(t *testing.T) {
	var q Queue
	const numGoroutines = 50
	const opsPerGoroutine = 100
	const totalEnqueued = numGoroutines * opsPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	var dequeued int64

	// Producers: concurrently enqueue items
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				q.Enqueue(fmt.Sprintf("item-%d-%d", id, j))
			}
		}(i)
	}

	// Consumers: concurrently dequeue items, counting non-empty returns
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerGoroutine; j++ {
				if val := q.Dequeue(); val != "" {
					atomic.AddInt64(&dequeued, 1)
				}
			}
		}()
	}

	wg.Wait()

	// Assert item conservation: successfully dequeued non-empty items + items
	// remaining in the queue must exactly equal the total number of items enqueued.
	remaining := q.Length()
	if totalAccounted := remaining + int(dequeued); totalAccounted != totalEnqueued {
		t.Fatalf("item conservation violated: enqueued %d, got %d dequeued + %d remaining = %d accounted",
			totalEnqueued, dequeued, remaining, totalAccounted)
	}
}
