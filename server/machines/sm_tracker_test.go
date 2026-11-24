package machines

import (
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionToStateMachineInstanceTracker_CompareAndDelete(t *testing.T) {
	connID := core.Uuid(uuid.Must(uuid.NewV4()))
	otherConnID := core.Uuid(uuid.Must(uuid.NewV4()))

	tracker := &ConnectionToStateMachineInstanceTracker{
		ConnectToInstanceMap: make(map[core.Uuid]*StateMachine),
	}

	inst1 := &StateMachine{ID: connID}
	inst2 := &StateMachine{ID: connID}

	t.Run("returns false when key not found", func(t *testing.T) {
		removed := tracker.CompareAndDelete(connID, inst1)
		assert.False(t, removed)
	})

	t.Run("returns false and preserves entry when tracked instance does not match", func(t *testing.T) {
		tracker.Add(connID, inst1)

		// Attempt to delete with a different instance pointer
		removed := tracker.CompareAndDelete(connID, inst2)
		assert.False(t, removed, "should not delete when instance pointer does not match")

		// Verify inst1 is still tracked
		got, ok := tracker.Get(connID)
		require.True(t, ok)
		assert.Same(t, inst1, got)
	})

	t.Run("returns true and deletes entry when tracked instance matches", func(t *testing.T) {
		tracker.Add(connID, inst1)
		removed := tracker.CompareAndDelete(connID, inst1)
		assert.True(t, removed, "should delete when instance pointer matches")

		// Verify entry is gone
		_, ok := tracker.Get(connID)
		assert.False(t, ok)
	})

	t.Run("prevents stale teardown from evicting re-imported instance", func(t *testing.T) {
		// 1. Initial machine is registered
		oldMachine := &StateMachine{ID: connID}
		tracker.Add(connID, oldMachine)

		// 2. Start asynchronous teardown of oldMachine
		addDone := make(chan struct{})
		teardownDone := make(chan struct{})
		go func() {
			// Wait until re-import has registered newMachine
			<-addDone
			tracker.CompareAndDelete(connID, oldMachine)
			close(teardownDone)
		}()

		// 3. Re-import registers newMachine before the stale teardown runs
		newMachine := &StateMachine{ID: connID}
		tracker.Add(connID, newMachine)
		close(addDone)

		// Wait for old teardown to complete
		<-teardownDone

		// 4. Verify newMachine survived and was NOT evicted by old teardown
		surviving, ok := tracker.Get(connID)
		require.True(t, ok, "newMachine must still be tracked")
		assert.Same(t, newMachine, surviving, "tracker must still hold newMachine, not nil or oldMachine")
	})

	t.Run("concurrent CompareAndDelete operations are race-free", func(t *testing.T) {
		tracker.Add(otherConnID, inst1)

		var wg sync.WaitGroup
		successCount := 0
		var countMu sync.Mutex

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if tracker.CompareAndDelete(otherConnID, inst1) {
					countMu.Lock()
					successCount++
					countMu.Unlock()
				}
			}()
		}

		wg.Wait()
		assert.Equal(t, 1, successCount, "exactly one concurrent CompareAndDelete should succeed")
		_, ok := tracker.Get(otherConnID)
		assert.False(t, ok)
	})
}
