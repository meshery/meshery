package handlers

import (
	"sync"

	"github.com/meshery/schemas/models/v1beta1/pattern"
)

type evalResult struct {
	resp pattern.EvaluationResponse
	err  error
}

// evaluationTracker coalesces concurrent evaluations of the same design:
// the first caller runs the evaluation, the rest wait for its result.
type evaluationTracker struct {
	mu       sync.Mutex
	lastGen  uint64
	inFlight map[string]*evalSlot
}

// evalSlot tracks one in-flight evaluation. The generation ties publishes to
// the acquire that started them: an evaluation abandoned on timeout/cancel
// can complete late, and without the generation check its publish would
// deliver a stale result to the waiters of a newer evaluation of the same
// design and drop that evaluation's real result.
type evalSlot struct {
	gen     uint64
	waiters []chan evalResult
}

func newEvaluationTracker() *evaluationTracker {
	return &evaluationTracker{
		inFlight: make(map[string]*evalSlot),
	}
}

// acquire returns (leader=true, gen, nil) for the first caller per designID;
// subsequent callers get (false, gen, waitCh) and must read one value from
// waitCh. The leader must pass gen back to publish.
func (t *evaluationTracker) acquire(designID string) (leader bool, gen uint64, wait <-chan evalResult) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if slot, exists := t.inFlight[designID]; exists {
		ch := make(chan evalResult, 1)
		slot.waiters = append(slot.waiters, ch)
		return false, slot.gen, ch
	}

	t.lastGen++
	t.inFlight[designID] = &evalSlot{gen: t.lastGen}
	return true, t.lastGen, nil
}

// publish broadcasts the result to the waiters of the evaluation started by
// the acquire that returned gen. Publishes for a cleared or superseded
// generation are no-ops.
func (t *evaluationTracker) publish(designID string, gen uint64, result evalResult) {
	t.mu.Lock()
	slot, ok := t.inFlight[designID]
	if !ok || slot.gen != gen {
		t.mu.Unlock()
		return
	}
	delete(t.inFlight, designID)
	t.mu.Unlock()

	for _, ch := range slot.waiters {
		ch <- result
	}
}
