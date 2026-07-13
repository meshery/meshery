package machines

import (
	"sync"

	"github.com/meshery/schemas/models/core"
	"golang.org/x/sync/singleflight"
)

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[core.Uuid]*StateMachine
	mx                   sync.RWMutex
	// initGroup deduplicates concurrent get-or-create calls for the same
	// connection ID so that exactly one StateMachine is ever built per
	// connection, even when multiple goroutines miss the cache simultaneously.
	initGroup singleflight.Group
}

func (smt *ConnectionToStateMachineInstanceTracker) Get(id core.Uuid) (*StateMachine, bool) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	inst, ok := smt.ConnectToInstanceMap[id]
	return inst, ok
}

func (smt *ConnectionToStateMachineInstanceTracker) Remove(id core.Uuid) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	delete(smt.ConnectToInstanceMap, id)
}

func (smt *ConnectionToStateMachineInstanceTracker) Add(id core.Uuid, inst *StateMachine) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	smt.ConnectToInstanceMap[id] = inst
}

// GetOrInit returns the StateMachine tracked for id, building and storing
// exactly one instance if none exists yet. The check-build-store sequence is
// made atomic per connection ID: concurrent callers for the same id share a
// single execution of build and all receive the same instance, upholding the
// "one FSM per connection" invariant even when several goroutines miss the
// cache at once. The instance is cached only when build succeeds, so a failed
// initialization is never observed as valid by later callers.
func (smt *ConnectionToStateMachineInstanceTracker) GetOrInit(id core.Uuid, build func() (*StateMachine, error)) (*StateMachine, error) {
	if inst, ok := smt.Get(id); ok {
		return inst, nil
	}

	v, err, _ := smt.initGroup.Do(id.String(), func() (interface{}, error) {
		// Re-check under the singleflight: a prior in-flight initialization
		// for this id may have completed and cached an instance between the
		// fast-path miss above and entering this call.
		if inst, ok := smt.Get(id); ok {
			return inst, nil
		}

		inst, err := build()
		if err != nil {
			return nil, err
		}
		smt.Add(id, inst)
		return inst, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*StateMachine), nil
}

// Range calls fn for a snapshot of the tracked connection/state-machine
// pairs. Iteration happens over a copy taken under the read lock, so fn may
// safely interact with the tracker or perform slow work.
func (smt *ConnectionToStateMachineInstanceTracker) Range(fn func(id core.Uuid, inst *StateMachine) bool) {
	smt.mx.RLock()
	snapshot := make(map[core.Uuid]*StateMachine, len(smt.ConnectToInstanceMap))
	for id, inst := range smt.ConnectToInstanceMap {
		snapshot[id] = inst
	}
	smt.mx.RUnlock()

	for id, inst := range snapshot {
		if !fn(id, inst) {
			return
		}
	}
}
