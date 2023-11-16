package machines

import (
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"

)

const (
	Discovery EventType = "discovery"
	Register EventType = "register"
	Connect EventType = "connect"
	Disconnect EventType = "disconnect"
	Ignore EventType = "ignore"
	NotFound EventType = "not found"
	Delete EventType = "delete"

	DISCOVERED   StateType = "discovered"
	REGISTERED   StateType = "registered"
	CONNECTED    StateType = "connected"
	IGNORED      StateType = "ignored"
	MAINTENANCE  StateType = "maintenance"
	DISCONNECTED StateType = "disconnected"
	DELETED      StateType = "deleted"
	NOTFOUND     StateType = "not found"
	
)

// Represents an event in the system/machine
type EventType string

// Represents an state in the system/machine
type StateType string

// Action to be executed in a given state.
type Action interface {

	// Used as guards/prerequisites checks and actions to be performed when the machine enters a given state.
	ExecuteOnEntry(ctx interface{}) (EventType, error)

	Execute(ctx interface{}) (State, EventType, error)

	// Used for cleanup actions to perform when the machine exits a given state
	ExecuteOnExit(ctx interface{}) (EventType, error)
}

// Represents the mapping between event and the next state in the event's response
type Events map[EventType]StateType

type State struct {
	Events Events
	Action Action
}

// Represents mapping between state name and the state
type States map[StateType]State

var DefaultState StateType = ""

type StateMachine struct {
	// ID to trace the events originated from the machine, also used in Logs
	ID uuid.UUID

	// Given name for the machine, used in Logs to track issues
	Name string

	// Configuration of states managed by the machine
	States States
	
	// Represent the previous state of the machine 
	PreviousState StateType

	// The current state of the machine
	CurrentState StateType

	// The initial state at which machine gets initialized
	InitialState StateType

	// Machine specific parameters/context.
	// Provided at initialization of the machine. 
	Context interface{} 
	
	mx sync.RWMutex

	Log logger.Handler
}

func (sm *StateMachine) getNextState(event EventType) (StateType, error) {
	sm.mx.RLock()
	defer sm.mx.RUnlock()

	state, ok := sm.States[sm.CurrentState]
	if ok {
		events := state.Events
		if events != nil {
			nextState, ok := events[event]
			if ok {
				sm.Log.Info("next state: ", nextState)
				return nextState, nil
			}
		}
	}
	return DefaultState, ErrInvalidTransition	
}

func (sm *StateMachine) SendEvent(event EventType) (error) {
	sm.mx.Lock()
	defer sm.mx.Unlock()

	nextState, err := sm.getNextState(event)
	if err != nil {
		sm.Log.Error(err)
		return err
	}

	// next state to transition
	state, ok := sm.States[nextState]
	if !ok || state.Action == nil {
		return ErrInvalidTransition
	}

	// Execute exit actions before entreing new state.
	action := sm.States[sm.CurrentState].Action
	if action != nil {
		action.ExecuteOnExit(sm.Context)
	} 

	
	var nextEvent EventType
	if state.Action != nil {

		// Execute entry actions for the state entered.
		_, err :=  state.Action.ExecuteOnEntry(sm.Context)
		if err != nil {
			sm.Log.Error(err)
			return err
		}
	
		_, nextEvent, err = state.Action.Execute(sm.Context)
		if err != nil {
			sm.Log.Error(err)
			return err
		}
	}

	sm.Log.Info(nextEvent)
	sm.PreviousState = sm.CurrentState
	sm.CurrentState = nextState
	sm.Log.Info("previous state for ", sm.Name, " ", sm.States[sm.PreviousState])
	sm.Log.Info("next state for ", sm.Name, " ", sm.States[sm.CurrentState])
	return nil
}

// var tm = Machine{
// 	States: States{
// 		test: State{
// 			Events: Events{
// 				event1:test2,
// 				event2:test3,
// 			},
// 		},
// 	},
// }

// var event1 EventType = "1"
// var event2 EventType = "2"
// var event3 EventType = "2"
// var test StateType = "test"
// var test2 StateType = "test"
// var test3 StateType = "test"
