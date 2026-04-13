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
	inFlight map[string][]chan evalResult
}

func newEvaluationTracker() *evaluationTracker {
	return &evaluationTracker{
		inFlight: make(map[string][]chan evalResult),
	}
}

// acquire returns (leader=true, nil) for the first caller per designID;
// subsequent callers get (false, waitCh) and must read one value from waitCh.
func (t *evaluationTracker) acquire(designID string) (leader bool, wait <-chan evalResult) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.inFlight[designID]; exists {
		ch := make(chan evalResult, 1)
		t.inFlight[designID] = append(t.inFlight[designID], ch)
		return false, ch
	}
	t.inFlight[designID] = nil
	return true, nil
}

// publish broadcasts the result to all waiters. Idempotent: subsequent
// calls for the same designID after the entry is cleared are no-ops.
func (t *evaluationTracker) publish(designID string, result evalResult) {
	t.mu.Lock()
	waiters, ok := t.inFlight[designID]
	if ok {
		delete(t.inFlight, designID)
	}
	t.mu.Unlock()

	for _, ch := range waiters {
		ch <- result
	}
}
