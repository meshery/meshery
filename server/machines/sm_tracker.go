package machines

import (
	"sync"

	"github.com/meshery/schemas/models/core"
)

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[core.Uuid]*StateMachine
	mx                   sync.RWMutex
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
