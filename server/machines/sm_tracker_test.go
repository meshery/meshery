package machines

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
)

func newTestTracker() *ConnectionToStateMachineInstanceTracker {
	return &ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*StateMachine),
	}
}

// TestGetOrInit_ConcurrentBuildsSingleInstance guards the fix for #20627: the
// get-or-create sequence used to be a non-atomic check-then-act, so concurrent
// requests for the same connection ID each constructed and started their own
// StateMachine, defeating the per-machine lock. GetOrInit must build exactly
// one machine and hand that same instance to every concurrent caller.
func TestGetOrInit_ConcurrentBuildsSingleInstance(t *testing.T) {
	tracker := newTestTracker()
	id := core.Uuid(uuid.Must(uuid.NewV4()))

	const goroutines = 64
	var builds int32

	start := make(chan struct{})
	results := make([]*StateMachine, goroutines)
	errs := make([]error, goroutines)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			<-start // release all goroutines at once to maximize contention
			results[idx], errs[idx] = tracker.GetOrInit(id, func() (*StateMachine, error) {
				atomic.AddInt32(&builds, 1)
				return &StateMachine{ID: id}, nil
			})
		}(i)
	}
	close(start)
	wg.Wait()

	if got := atomic.LoadInt32(&builds); got != 1 {
		t.Fatalf("build func called %d times, want exactly 1", got)
	}

	tracked, ok := tracker.Get(id)
	if !ok {
		t.Fatal("expected the machine to be tracked after GetOrInit")
	}
	for i := 0; i < goroutines; i++ {
		if errs[i] != nil {
			t.Fatalf("goroutine %d got error: %v", i, errs[i])
		}
		if results[i] != tracked {
			t.Fatalf("goroutine %d received a StateMachine other than the single tracked instance", i)
		}
	}
}

// TestGetOrInit_FailedBuildNotCached guards the fix for #20627: a failed
// initialization must not be cached as valid, and a later call must be free to
// retry the build.
func TestGetOrInit_FailedBuildNotCached(t *testing.T) {
	tracker := newTestTracker()
	id := core.Uuid(uuid.Must(uuid.NewV4()))
	wantErr := errors.New("init failed")

	inst, err := tracker.GetOrInit(id, func() (*StateMachine, error) {
		return nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("got error %v, want %v", err, wantErr)
	}
	if inst != nil {
		t.Fatalf("expected nil instance on build failure, got %v", inst)
	}
	if _, ok := tracker.Get(id); ok {
		t.Fatal("failed initialization must not be cached in the tracker")
	}

	// A subsequent call must be able to retry and succeed.
	built := &StateMachine{ID: id}
	inst, err = tracker.GetOrInit(id, func() (*StateMachine, error) {
		return built, nil
	})
	if err != nil {
		t.Fatalf("retry after failure returned error: %v", err)
	}
	if inst != built {
		t.Fatal("retry did not return the newly built instance")
	}
	if tracked, ok := tracker.Get(id); !ok || tracked != built {
		t.Fatal("successful retry was not cached")
	}
}

// TestGetOrInit_ExistingInstanceSkipsBuild verifies the fast path: when a
// machine is already tracked, GetOrInit returns it without invoking build.
func TestGetOrInit_ExistingInstanceSkipsBuild(t *testing.T) {
	tracker := newTestTracker()
	id := core.Uuid(uuid.Must(uuid.NewV4()))
	existing := &StateMachine{ID: id}
	tracker.Add(id, existing)

	inst, err := tracker.GetOrInit(id, func() (*StateMachine, error) {
		t.Fatal("build must not be called when an instance already exists")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inst != existing {
		t.Fatal("GetOrInit did not return the already-tracked instance")
	}
}
