package machines

import (
	"sync"
	"testing"
	"time"
)

// Jobs must run in the order they were submitted. This is the guarantee that
// keeps a connect's deploy from being overtaken by a following disconnect's
// undeploy (the state-vs-cluster divergence in issue #20671).
func TestSerialExecutor_RunsJobsInSubmissionOrder(t *testing.T) {
	const n = 200
	e := NewSerialExecutor(n)

	var mu sync.Mutex
	order := make([]int, 0, n)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		i := i
		if !e.Submit(func() {
			mu.Lock()
			order = append(order, i)
			mu.Unlock()
			wg.Done()
		}) {
			t.Fatalf("Submit returned false for job %d on an open executor", i)
		}
	}
	wg.Wait()

	for pos, got := range order {
		if got != pos {
			t.Fatalf("job ran out of order: position %d holds job %d; full order: %v", pos, got, order)
		}
	}
}

// Jobs must never run concurrently. counter is mutated without synchronization,
// so the race detector (go test -race) flags this test if the executor ever
// runs two jobs at once - the data race on the shared MesheryCtrlsHelper that
// issue #20671 reported.
func TestSerialExecutor_ConfinesJobsToSingleGoroutine(t *testing.T) {
	const n = 1000
	e := NewSerialExecutor(n)

	counter := 0
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		e.Submit(func() {
			counter++ // safe only because jobs are guaranteed to run one at a time
			wg.Done()
		})
	}
	wg.Wait()

	if counter != n {
		t.Fatalf("counter = %d, want %d: jobs appear to have run concurrently", counter, n)
	}
}

// After Close the executor rejects new work and its worker goroutine exits once
// the already-queued jobs have run, so it does not outlive a deleted connection.
func TestSerialExecutor_SubmitAfterCloseIsRejected(t *testing.T) {
	e := NewSerialExecutor(4)

	ran := make(chan struct{}, 1)
	if !e.Submit(func() { ran <- struct{}{} }) {
		t.Fatal("Submit before Close returned false")
	}
	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("queued job did not run")
	}

	e.Close()
	e.Close() // idempotent: must not panic

	if e.Submit(func() { t.Error("a job submitted after Close must not run") }) {
		t.Fatal("Submit after Close returned true; want false")
	}

	select {
	case <-e.done:
	case <-time.After(2 * time.Second):
		t.Fatal("worker goroutine did not exit after Close")
	}
}
