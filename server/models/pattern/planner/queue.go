package planner

import "sync"

// Queue is the thread safe queue
type Queue struct {
	internalQueue []string

	sync.RWMutex
}

// Enqueue data to the queue
func (q *Queue) Enqueue(str string) {
	q.Lock()
	defer q.Unlock()

	q.internalQueue = append(q.internalQueue, str)
}

// Dequeue removes data from the queue
func (q *Queue) Dequeue() string {
	q.Lock()
	defer q.Unlock()

	val := q.internalQueue[0]
	q.internalQueue = q.internalQueue[1:]

	return val
}

// Length returns the length of the queue
func (q *Queue) Length() int {
	q.Lock()
	defer q.Unlock()

	return len(q.internalQueue)
}

// Peek returns the first item in the queue
func (q *Queue) Peek() string {
	q.Lock()
	defer q.Unlock()

	if len(q.internalQueue) > 0 {
		return q.internalQueue[0]
	}

	return ""
}
