package handlers

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEvaluationTracker_SingleLeader(t *testing.T) {
	tr := newEvaluationTracker()

	leader, wait := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}
	if wait != nil {
		t.Fatal("leader should not have a wait channel")
	}

	// publish with no followers should not panic or block
	tr.publish("d1", evalResult{})

	// after publish, the next caller should again be a leader
	leader2, _ := tr.acquire("d1")
	if !leader2 {
		t.Fatal("after publish, next caller should be a new leader")
	}
}

func TestEvaluationTracker_CoalescesConcurrent(t *testing.T) {
	tr := newEvaluationTracker()

	// Leader acquires first and does not publish yet.
	leader, _ := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}

	const followers = 50
	waits := make([]<-chan evalResult, 0, followers)
	for range followers {
		isLeader, w := tr.acquire("d1")
		if isLeader {
			t.Fatal("subsequent callers should be followers")
		}
		waits = append(waits, w)
	}

	// Leader finishes and publishes once. All followers must receive the same result.
	sentinelErr := errors.New("boom")
	tr.publish("d1", evalResult{err: sentinelErr})

	var wg sync.WaitGroup
	var received int32
	for _, w := range waits {
		wg.Add(1)
		go func(ch <-chan evalResult) {
			defer wg.Done()
			select {
			case r := <-ch:
				if !errors.Is(r.err, sentinelErr) {
					t.Errorf("follower got wrong err: %v", r.err)
				}
				atomic.AddInt32(&received, 1)
			case <-time.After(2 * time.Second):
				t.Error("follower timed out waiting for result")
			}
		}(w)
	}
	wg.Wait()

	if got := atomic.LoadInt32(&received); got != followers {
		t.Fatalf("expected %d followers to receive result, got %d", followers, got)
	}
}

func TestEvaluationTracker_PublishIsIdempotent(t *testing.T) {
	tr := newEvaluationTracker()
	_, _ = tr.acquire("d1")

	tr.publish("d1", evalResult{})
	// second publish must be a no-op (in particular, no panic).
	tr.publish("d1", evalResult{})
}

func TestEvaluationTracker_DistinctDesignsAreIndependent(t *testing.T) {
	tr := newEvaluationTracker()

	leader1, _ := tr.acquire("d1")
	leader2, _ := tr.acquire("d2")
	if !leader1 || !leader2 {
		t.Fatal("different designs should each get their own leader")
	}
}
