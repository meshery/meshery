package oam

import "sync"

// Queue is the thread safe queue
type Queue struct {
	InternalQueue []string

	sync.RWMutex
}

// Enqueue data to the queue
func (q *Queue) Enqueue(str string) {
	q.Lock()
	defer q.Unlock()

	q.InternalQueue = append(q.InternalQueue, str)
}

// Dequeue removes data from the queue
func (q *Queue) Dequeue() string {
	q.Lock()
	defer q.Unlock()

	val := q.InternalQueue[0]
	q.InternalQueue = q.InternalQueue[1:]

	return val
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	q.Lock()
	defer q.Unlock()

	return len(q.InternalQueue)
}

// Peek returns the first item in the queue
func (q *Queue) Peek() string {
	q.Lock()
	defer q.Unlock()

	if len(q.InternalQueue) > 0 {
		return q.InternalQueue[0]
	}

	return ""
}
