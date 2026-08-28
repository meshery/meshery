package machines

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
)

func TestGetOrInitializeThunderingHerd(t *testing.T) {
	tracker := &ConnectionToStateMachineInstanceTracker{}
	connID, _ := uuid.NewV4()

	initAttempts := 0
	expectedErr := fmt.Errorf("mock initialization failure")

	// initFn will block until we close the barrier, then fail
	barrier := make(chan struct{})
	initFn := func() (*StateMachine, error) {
		initAttempts++
		<-barrier
		return nil, expectedErr
	}

	var wg sync.WaitGroup
	numWaiters := 5

	results := make([]error, numWaiters)

	// Launch multiple concurrent callers
	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := tracker.GetOrInitialize(core.Uuid(connID), initFn)
			results[idx] = err
		}(i)
	}

	// Wait a moment to ensure all goroutines are blocked waiting on the same in-flight state
	time.Sleep(100 * time.Millisecond)

	// Unblock the initialization
	close(barrier)

	wg.Wait()

	// Verify that initialization was only attempted ONCE, despite 5 concurrent callers
	if initAttempts != 1 {
		t.Fatalf("expected initialization to be attempted exactly 1 time, got %d", initAttempts)
	}

	// Verify that ALL callers received the exact same error
	for i, err := range results {
		if err != expectedErr {
			t.Errorf("caller %d expected error %v, got %v", i, expectedErr, err)
		}
	}

	// Verify that a SUBSEQUENT call triggers a NEW initialization attempt
	initAttempts = 0
	barrier2 := make(chan struct{})
	close(barrier2) // unblock immediately
	initFnSuccess := func() (*StateMachine, error) {
		initAttempts++
		<-barrier2
		return &StateMachine{}, nil
	}

	inst, err := tracker.GetOrInitialize(core.Uuid(connID), initFnSuccess)
	if err != nil {
		t.Fatalf("expected subsequent initialization to succeed, got error: %v", err)
	}
	if inst == nil {
		t.Fatalf("expected non-nil StateMachine instance")
	}
	if initAttempts != 1 {
		t.Fatalf("expected subsequent initialization to be attempted exactly 1 time, got %d", initAttempts)
	}
}

// TestGetOrInitializePanicYieldsErrorToAllWaiters verifies that if an
// initialization function panics, the panic is safely recovered and converted
// into an error, ensuring that no waiter receives a raw (nil, nil) result and
// that the tracker cleans up in-flight state so subsequent attempts can succeed.
func TestGetOrInitializePanicYieldsErrorToAllWaiters(t *testing.T) {
	tracker := &ConnectionToStateMachineInstanceTracker{}
	connID, _ := uuid.NewV4()

	initAttempts := 0

	barrier := make(chan struct{})
	initFn := func() (*StateMachine, error) {
		initAttempts++
		<-barrier
		panic("simulated fatal panic during initialization")
	}

	var wg sync.WaitGroup
	numWaiters := 5
	results := make([]error, numWaiters)

	for i := 0; i < numWaiters; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			inst, err := tracker.GetOrInitialize(core.Uuid(connID), initFn)
			if inst == nil && err == nil {
				t.Errorf("waiter %d received (nil, nil) after initialization panic", idx)
			}
			results[idx] = err
		}(i)
	}

	time.Sleep(100 * time.Millisecond)
	close(barrier)
	wg.Wait()

	if initAttempts != 1 {
		t.Fatalf("expected initialization to be attempted exactly 1 time, got %d", initAttempts)
	}

	for i, err := range results {
		if err == nil {
			t.Errorf("caller %d expected an error due to panic, got nil", i)
		}
	}

	// Verify that a SUBSEQUENT call triggers a NEW initialization attempt
	initAttempts = 0
	barrier2 := make(chan struct{})
	close(barrier2) // unblock immediately
	initFnSuccess := func() (*StateMachine, error) {
		initAttempts++
		<-barrier2
		return &StateMachine{}, nil
	}

	inst, err := tracker.GetOrInitialize(core.Uuid(connID), initFnSuccess)
	if err != nil {
		t.Fatalf("expected subsequent initialization to succeed, got error: %v", err)
	}
	if inst == nil {
		t.Fatalf("expected non-nil StateMachine instance")
	}
	if initAttempts != 1 {
		t.Fatalf("expected subsequent initialization to be attempted exactly 1 time, got %d", initAttempts)
	}
}
