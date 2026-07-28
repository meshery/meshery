package machines

import "sync"

// SerialExecutor runs submitted jobs one at a time, in submission order, on a
// single background goroutine.
//
// A Kubernetes connection's lifecycle actions (connect/disconnect/delete) run
// their cluster side effects - deploying/undeploying the operator and MeshSync
// and wiring data handlers - against a MesheryControllersHelper that is shared
// across all of that connection's actions. Previously each action ran its side
// effects in its own detached goroutine, which let them (a) finish out of order
// so a connect's deploy could land after a following disconnect's undeploy and
// leave the persisted connection status diverged from the cluster, and (b)
// mutate the shared helper concurrently, a data race. Funnelling the side
// effects through one per-connection SerialExecutor keeps them ordered and
// confines the helper to a single goroutine.
//
// Submissions for a given connection are already serialized by the state
// machine's lock (SendEvent holds it across the action), so jobs are enqueued in
// event order and the worker runs them in that same order.
type SerialExecutor struct {
	mu     sync.Mutex
	jobs   chan func()
	closed bool
	// done is closed once the worker goroutine has drained every queued job and
	// returned, i.e. after Close.
	done chan struct{}
}

// NewSerialExecutor starts the executor's worker goroutine. buffer bounds how
// many jobs may be queued before Submit blocks; since submissions for a single
// connection are serialized upstream, blocking here just applies back-pressure
// to that one connection rather than dropping work.
func NewSerialExecutor(buffer int) *SerialExecutor {
	e := &SerialExecutor{
		jobs: make(chan func(), buffer),
		done: make(chan struct{}),
	}
	go e.run()
	return e
}

func (e *SerialExecutor) run() {
	defer close(e.done)
	for job := range e.jobs {
		job()
	}
}

// Submit enqueues job to run after every previously-submitted job has run. It
// returns false, dropping job, if the executor has been closed - which happens
// once the connection is deleted and no further side effects should occur.
func (e *SerialExecutor) Submit(job func()) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return false
	}
	e.jobs <- job
	return true
}

// Close stops the executor after all already-submitted jobs have run. It is
// idempotent and safe to call while the connection is being torn down.
func (e *SerialExecutor) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return
	}
	e.closed = true
	close(e.jobs)
}
