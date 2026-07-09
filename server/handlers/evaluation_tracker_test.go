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

	leader, gen, wait := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}
	if wait != nil {
		t.Fatal("leader should not have a wait channel")
	}

	// publish with no followers should not panic or block
	tr.publish("d1", gen, evalResult{})

	// after publish, the next caller should again be a leader
	leader2, _, _ := tr.acquire("d1")
	if !leader2 {
		t.Fatal("after publish, next caller should be a new leader")
	}
}

func TestEvaluationTracker_CoalescesConcurrent(t *testing.T) {
	tr := newEvaluationTracker()

	// Leader acquires first and does not publish yet.
	leader, gen, _ := tr.acquire("d1")
	if !leader {
		t.Fatal("first caller should be the leader")
	}

	const followers = 50
	waits := make([]<-chan evalResult, 0, followers)
	for range followers {
		isLeader, _, w := tr.acquire("d1")
		if isLeader {
			t.Fatal("subsequent callers should be followers")
		}
		waits = append(waits, w)
	}

	// Leader finishes and publishes once. All followers must receive the same result.
	sentinelErr := errors.New("boom")
	tr.publish("d1", gen, evalResult{err: sentinelErr})

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
	_, gen, _ := tr.acquire("d1")

	tr.publish("d1", gen, evalResult{})
	// second publish must be a no-op (in particular, no panic).
	tr.publish("d1", gen, evalResult{})
}

func TestEvaluationTracker_DistinctDesignsAreIndependent(t *testing.T) {
	tr := newEvaluationTracker()

	leader1, _, _ := tr.acquire("d1")
	leader2, _, _ := tr.acquire("d2")
	if !leader1 || !leader2 {
		t.Fatal("different designs should each get their own leader")
	}
}

// A publish carrying the generation of an abandoned evaluation must not
// deliver into a newer evaluation of the same design. Sequence mirrors the
// handler's timeout path: leader A times out and clears the slot, leader B
// starts a fresh evaluation with follower C, then A's goroutine finishes
// late.
func TestEvaluationTracker_StaleGenerationPublishIsIgnored(t *testing.T) {
	tr := newEvaluationTracker()

	// A leads, then its HTTP side times out and clears the slot.
	leaderA, genA, _ := tr.acquire("d1")
	if !leaderA {
		t.Fatal("A should be the leader")
	}
	tr.publish("d1", genA, evalResult{err: errors.New("timeout")})

	// B starts a new evaluation of the same design; C coalesces onto it.
	leaderB, genB, _ := tr.acquire("d1")
	if !leaderB {
		t.Fatal("B should lead the new evaluation")
	}
	_, _, waitC := tr.acquire("d1")

	// A's abandoned goroutine finishes late. C must not see its result.
	staleErr := errors.New("stale result from A")
	tr.publish("d1", genA, evalResult{err: staleErr})

	select {
	case r := <-waitC:
		t.Fatalf("follower received a result from the abandoned evaluation: %v", r.err)
	default:
	}

	// B's real result must still reach C.
	realErr := errors.New("real result from B")
	tr.publish("d1", genB, evalResult{err: realErr})

	select {
	case r := <-waitC:
		if !errors.Is(r.err, realErr) {
			t.Fatalf("follower got %v, want the current leader's result", r.err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("follower timed out waiting for the current leader's result")
	}
}
