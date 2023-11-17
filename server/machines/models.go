package machines

import (
	"context"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
)

const (
	Discovery  EventType = "discovery"
	Register   EventType = "register"
	Connect    EventType = "connect"
	Disconnect EventType = "disconnect"
	Ignore     EventType = "ignore"
	NotFound   EventType = "not found"
	Delete     EventType = "delete"
	NoOp       EventType = "noop"

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
	ExecuteOnEntry(context context.Context, machinectx interface{}) (EventType, *events.Event, error)

	Execute(context context.Context, machinectx interface{}) (EventType, *events.Event, error)

	// Used for cleanup actions to perform when the machine exits a given state
	ExecuteOnExit(context context.Context, machinectx interface{}) (EventType, *events.Event, error)
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

var InitialState StateType = "initialized"

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

	Log logger.Handler
}

func (sm *StateMachine) ResetState() {
	sm.mx.Lock()
	defer sm.mx.Unlock()

	sm.CurrentState = DISCOVERED
}

func (sm *StateMachine) getNextState(event EventType) (StateType, error) {
	state, ok := sm.States[sm.CurrentState]
	sm.Log.Info("inside getNextState: ", event, ok)
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
	return DefaultState, ErrInvalidTransitionEvent(sm.CurrentState, event)
}

// This should handle error and event publishing . THis should return the events.Event and error and the func in which SendEvent is called should publish the event.
// The statusCode for eg: in updateSconnectionStatusById is lost so the func that calls send event will not know the original statuscode, we will be sending 500 only, if possible capture the statuscode inside the event itself.

func (sm *StateMachine) SendEvent(ctx context.Context, eventType EventType, payload interface{}) error {
	sm.mx.Lock()
	defer sm.mx.Unlock()
	sm.Log.Info("inside send event entering loop")
	for {
		var nextEventTransition EventType
		if nextEventTransition == NoOp {
			break
		}
		
		nextState, err := sm.getNextState(eventType)
		if err != nil {
			sm.Log.Error(err)
			return err
		}
		sm.Log.Info("inside models send event line 130: next state", nextState)

		// next state to transition
		state, ok := sm.States[nextState]
		if !ok || state.Action == nil {
			return ErrInvalidTransition(sm.CurrentState, nextState)
		}

		// Execute exit actions before entreing new state.
		action := sm.States[sm.CurrentState].Action
		if action != nil {
			_, _, err := action.ExecuteOnExit(ctx, sm.Context)
			if err != nil {
				sm.Log.Error(err)
				break
			}
		}

		
		if state.Action != nil {

			// Execute entry actions for the state entered.
			et, event, err := state.Action.ExecuteOnEntry(ctx, sm.Context)
			sm.Log.Info("inside models.go entry action executed event emiited", et, "events.Event", event)

			if err != nil {
				sm.Log.Error(err)
				// event publishing
				sm.Log.Info(event)
				nextEventTransition = et
				
			} else {	
				nextEventTransition, event, err = state.Action.Execute(ctx, sm.Context)
				sm.Log.Info("inside models.go action executed, nextEventTransition: ", nextEventTransition)
				if err != nil {
					sm.Log.Error(err)
					// event publishing
					sm.Log.Info(event)
				}
			}
			
			if nextEventTransition != NoOp {
				eventType = nextEventTransition
			}
		}

		sm.Log.Info(nextEventTransition)
		sm.PreviousState = sm.CurrentState
		sm.CurrentState = nextState
		sm.Log.Info("previous state for ", sm.Name, " ", sm.States[sm.PreviousState])
		sm.Log.Info("next state for ", sm.Name, " ", sm.States[sm.CurrentState])
	}
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
