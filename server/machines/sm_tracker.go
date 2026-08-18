package machines

import (
	"sync"

	"github.com/meshery/schemas/models/core"
)

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[core.Uuid]*StateMachine
	inFlight             map[core.Uuid]chan struct{}
	mx                   sync.RWMutex
}

func (smt *ConnectionToStateMachineInstanceTracker) GetOrInitialize(id core.Uuid, initFn func() (*StateMachine, error)) (*StateMachine, error) {
	smt.mx.Lock()
	if inst, ok := smt.ConnectToInstanceMap[id]; ok {
		smt.mx.Unlock()
		return inst, nil
	}

	if smt.inFlight == nil {
		smt.inFlight = make(map[core.Uuid]chan struct{})
	}

	ch, inFlight := smt.inFlight[id]
	if !inFlight {
		ch = make(chan struct{})
		smt.inFlight[id] = ch
	}
	smt.mx.Unlock()

	if inFlight {
		<-ch
		// The other goroutine finished initialization. We can now just get it from the map.
		smt.mx.Lock()
		inst, ok := smt.ConnectToInstanceMap[id]
		smt.mx.Unlock()
		if ok {
			return inst, nil
		}
		// If it's not in the map, the initialization failed, so we should try again.
		return smt.GetOrInitialize(id, initFn)
	}

	// We are the one initializing
	inst, err := initFn()

	smt.mx.Lock()
	defer smt.mx.Unlock()

	if err == nil && inst != nil {
		if smt.ConnectToInstanceMap == nil {
			smt.ConnectToInstanceMap = make(map[core.Uuid]*StateMachine)
		}
		smt.ConnectToInstanceMap[id] = inst
	}

	delete(smt.inFlight, id)
	close(ch) // wake up waiters

	return inst, err
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
