package machines

import (
	"context"
	"sync"

	"github.com/meshery/schemas/models/core"
)

type inFlightState struct {
	ch   chan struct{}
	inst *StateMachine
	err  error
}

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[core.Uuid]*StateMachine
	inFlight             map[core.Uuid]*inFlightState
	mx                   sync.RWMutex
}

func (smt *ConnectionToStateMachineInstanceTracker) GetOrInitialize(id core.Uuid, initFn func() (*StateMachine, error)) (*StateMachine, error) {
	smt.mx.Lock()
	if inst, ok := smt.ConnectToInstanceMap[id]; ok {
		smt.mx.Unlock()
		return inst, nil
	}

	if smt.inFlight == nil {
		smt.inFlight = make(map[core.Uuid]*inFlightState)
	}

	state, inFlight := smt.inFlight[id]
	if !inFlight {
		state = &inFlightState{ch: make(chan struct{})}
		smt.inFlight[id] = state
	}
	smt.mx.Unlock()

	if inFlight {
		<-state.ch
		// The other goroutine finished initialization. Return its exact result.
		if state.err != nil {
			return nil, state.err
		}
		// If there is no error, it must be in the map (or state.inst is populated).
		// We can return the cached instance safely.
		return state.inst, nil
	}

	// We are the one initializing
	defer func() {
		smt.mx.Lock()
		delete(smt.inFlight, id)
		close(state.ch) // wake up waiters
		smt.mx.Unlock()
	}()

	inst, err := initFn()

	smt.mx.Lock()
	defer smt.mx.Unlock()

	state.inst = inst
	state.err = err

	if err == nil && inst != nil {
		if smt.ConnectToInstanceMap == nil {
			smt.ConnectToInstanceMap = make(map[core.Uuid]*StateMachine)
		}
		smt.ConnectToInstanceMap[id] = inst
	}

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

func (smt *ConnectionToStateMachineInstanceTracker) RemoveIfMatch(id core.Uuid, expectedInst *StateMachine) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	if inst, ok := smt.ConnectToInstanceMap[id]; ok && inst == expectedInst {
		delete(smt.ConnectToInstanceMap, id)
	}
}

func (smt *ConnectionToStateMachineInstanceTracker) RemoveIfMatchAndState(id core.Uuid, expectedInst *StateMachine, allowedStates ...StateType) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	if inst, ok := smt.ConnectToInstanceMap[id]; ok && inst == expectedInst {
		currentState := inst.GetCurrentState()
		for _, s := range allowedStates {
			if currentState == s {
				delete(smt.ConnectToInstanceMap, id)
				return
			}
		}
	}
}

func (smt *ConnectionToStateMachineInstanceTracker) RemoveIfMatchAndGeneration(id core.Uuid, expectedInst *StateMachine, expectedGeneration context.Context) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	if inst, ok := smt.ConnectToInstanceMap[id]; ok && inst == expectedInst {
		if inst.GetLifecycleCtx() == expectedGeneration {
			delete(smt.ConnectToInstanceMap, id)
		}
	}
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
