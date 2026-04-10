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
