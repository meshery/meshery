package machines

import (
	"sync"

	"github.com/gofrs/uuid"
)

type ConnectionToStateMachineInstanceTracker struct {
	ConnectToInstanceMap map[uuid.UUID]*StateMachine
	Mx                   sync.RWMutex
}
