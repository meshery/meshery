package planner

import "sync"

// Queue is a thread-safe FIFO queue.
type Queue struct {
	internalQueue []string

	mu sync.RWMutex
}

// Enqueue adds an item to the end of the queue. Any string, including "", is valid.
func (q *Queue) Enqueue(str string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.internalQueue = append(q.internalQueue, str)
}

// Dequeue removes and returns the first item from the queue.
// The bool is false if the queue was empty.
func (q *Queue) Dequeue() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.internalQueue) == 0 {
		return "", false
	}

	val := q.internalQueue[0]
	q.internalQueue[0] = "" // release the reference for garbage collection
	q.internalQueue = q.internalQueue[1:]
	if len(q.internalQueue) == 0 {
		q.internalQueue = nil
	}

	return val, true
}

// Length returns the number of items in the queue.
func (q *Queue) Length() int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return len(q.internalQueue)
}

// Peek returns the first item without removing it.
// The bool is false if the queue is empty.
func (q *Queue) Peek() (string, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if len(q.internalQueue) == 0 {
		return "", false
	}

	return q.internalQueue[0], true
}
