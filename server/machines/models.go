package machines

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"

)

// Represents an event in the system/machine
type EventType string

// Represents an state in the system/machine
type StateType string

// Action to be executed in a given state.
type Action interface{
	Action(ctx interface{}) (State, EventType, error)
}

// Represents the mapping between event and the next state in the event's response
type Events map[EventType]StateType

type State struct {
	Events Events
	Action Action

	// Used as guards/prerequisites checks and actions to be performed when the machine enters a given state.
	EntryAction Action

	// Used for cleanup actions to perform when the machine exits a given state
	ExitAction Action
}

// Represents mapping between state name and the state
type States map[StateType]State

var DefaultState StateType = ""

type StateMachine struct {
	// ID to trace the events originated from the machine, also used in logs
	ID uuid.UUID

	// Given name for the machine, used in logs to track issues
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

	log logger.Handler
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
				sm.log.Info("next state: ", nextState)
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
		sm.log.Error(err)
		return err
	}

	state, ok := sm.States[nextState]
	if !ok || state.Action == nil {
		return ErrInvalidTransition
	}

	exitAction := sm.States[sm.CurrentState].Action
	if exitAction != nil {
		exitAction.Action(sm.Context)
	}


	prerequisiteAction := state.EntryAction
	if prerequisiteAction != nil {
		_, _, err := prerequisiteAction.Action(sm.Context)
		if err != nil {
			sm.log.Error(err)
			return err
		}
	}

	_, nextEvent, err := state.Action.Action(sm.Context)
	if err != nil {
		sm.log.Error(err)
		return err
	}

	sm.log.Info(nextEvent)
	sm.PreviousState = sm.CurrentState
	sm.CurrentState = nextState
	sm.log.Info("previous state for ", sm.Name, " ", sm.States[sm.PreviousState])
	sm.log.Info("next state for ", sm.Name, " ", sm.States[sm.CurrentState])
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
