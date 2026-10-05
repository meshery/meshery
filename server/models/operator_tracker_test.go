package models

import (
	"fmt"
	"sync"
	"testing"
)

// TestOperatorTracker_ConcurrentAccess exercises Undeployed and IsUndeployed
// concurrently across distinct context IDs. Without locking this fails under
// -race or with a fatal concurrent map access.
func TestOperatorTracker_ConcurrentAccess(t *testing.T) {
	ot := NewOperatorTracker(false)

	const contexts = 8
	const ops = 100

	var wg sync.WaitGroup
	for i := 0; i < contexts; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctxID := fmt.Sprintf("ctx-%d", i)
			for j := 0; j < ops; j++ {
				ot.Undeployed(ctxID, j%2 == 0)
				_ = ot.IsUndeployed(ctxID)
			}
		}(i)
	}
	wg.Wait()

	// Final state must remain readable after the burst.
	for i := 0; i < contexts; i++ {
		_ = ot.IsUndeployed(fmt.Sprintf("ctx-%d", i))
	}
}

// TestOperatorTracker_ConcurrentLazyInit covers the nil-map branch: a zero
// value tracker lazily creates the map on first use, which must also hold
// the lock or two goroutines each allocate and assign a fresh map.
func TestOperatorTracker_ConcurrentLazyInit(t *testing.T) {
	ot := &OperatorTracker{}

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ctxID := fmt.Sprintf("lazy-ctx-%d", i)
			ot.Undeployed(ctxID, true)
			if !ot.IsUndeployed(ctxID) {
				t.Errorf("IsUndeployed(%q) = false, want true", ctxID)
			}
		}(i)
	}
	wg.Wait()
}

// TestOperatorTracker_DisabledStaysNoOp keeps the existing DisableOperator
// semantics pinned: writes are dropped and reads always report undeployed.
func TestOperatorTracker_DisabledStaysNoOp(t *testing.T) {
	ot := NewOperatorTracker(true)

	ot.Undeployed("any", false)
	if !ot.IsUndeployed("any") {
		t.Fatal("IsUndeployed with DisableOperator = false, want true")
	}
}
