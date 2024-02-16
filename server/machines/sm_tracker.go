package machines

import (
	"sync"

	"github.com/gofrs/uuid"
)

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[uuid.UUID]*StateMachine
	mx                   sync.RWMutex
}

func (smt *ConnectionToStateMachineInstanceTracker) Get(id uuid.UUID) (*StateMachine, bool) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	inst, ok := smt.ConnectToInstanceMap[id]
	return inst, ok
}

func (smt *ConnectionToStateMachineInstanceTracker) Remove(id uuid.UUID) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	delete(smt.ConnectToInstanceMap, id)
}

func (smt *ConnectionToStateMachineInstanceTracker) Add(id uuid.UUID, inst *StateMachine) {
	smt.mx.Lock()
	defer smt.mx.Unlock()
	smt.ConnectToInstanceMap[id] = inst
}
