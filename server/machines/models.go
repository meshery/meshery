package machines

import (
	"context"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
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

	Provider models.Provider
}

type initFunc func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error)

func (sm *StateMachine) Start(ctx context.Context, machinectx interface{}, log logger.Handler, init initFunc) (*events.Event, error) {
	var mCtx interface{}
	var event *events.Event
	var err error
	if init != nil {
		mCtx, event, err = init(ctx, machinectx, log)
		if err != nil {
			return event, err
		}
	}
	sm.Context = mCtx
	return nil, nil
}

func (sm *StateMachine) ResetState() {
	sm.mx.Lock()
	defer sm.mx.Unlock()

	sm.CurrentState = InitialState
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

// This should handle error and event publishing . This should return the events.Event and error. The func invoking the SendEvent should publish the event.

func (sm *StateMachine) SendEvent(ctx context.Context, eventType EventType, payload interface{}) (*events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	defaultEvent := events.NewEvent().WithDescription(fmt.Sprintf("Invalid status change requested to %s for connection type %s.", eventType, sm.Name)).ActedUpon(sm.ID).FromUser(userUUID).FromSystem(*sysID).WithSeverity(events.Error)
	sm.mx.Lock()
	defer sm.mx.Unlock()
	var event *events.Event
	var err error
	for {
		if eventType == NoOp {
			break
		}

		nextState, err := sm.getNextState(eventType)
		if err != nil {
			sm.Log.Error(err)
			sm.Log.Info(defaultEvent.WithMetadata(map[string]interface{}{"error": err}).Build())
			eventType = NoOp
			break
		}

		sm.Log.Info("transitioning to next state: ", nextState)

		// next state to transition
		state, ok := sm.States[nextState]
		if !ok || state.Action == nil {
			sm.Log.Error(err)
			sm.Log.Info(defaultEvent.WithMetadata(map[string]interface{}{"error": ErrInvalidTransition(sm.CurrentState, nextState)}).Build())
			eventType = NoOp
			break
		}

		// Execute exit actions before entering new state.
		action := sm.States[sm.CurrentState].Action
		if action != nil {
			_, event, err = action.ExecuteOnExit(ctx, sm.Context)
			if err != nil {
				sm.Log.Error(err)
				return event, err
			}
		}

		if state.Action != nil {
			// Execute entry actions for the state entered.
			eventType, event, err = state.Action.ExecuteOnEntry(ctx, sm.Context)
			sm.Log.Info("entry action executed, event emitted ", eventType)

			if err != nil {
				sm.Log.Error(err)
				sm.Log.Info(event)
			} else {
				eventType, event, err = state.Action.Execute(ctx, sm.Context)
				sm.Log.Info("inside action executed, event emitted ", eventType)
				if err != nil {
					sm.Log.Error(err)
					sm.Log.Info(event)
				}
			}
		}

		sm.PreviousState = sm.CurrentState
		sm.CurrentState = nextState
	}

	token, _ := ctx.Value(models.TokenCtxKey).(string)
	connection, statusCode, err := sm.Provider.UpdateConnectionStatusByID(token, sm.ID, connections.ConnectionStatus(sm.CurrentState))

	if err != nil {
		// In this case should the current state be again set to previous state i.e. should we rollback. But not only state should be rollback but other actions as well, rn we don't rollback state.
		return events.NewEvent().WithDescription(fmt.Sprintf("Operation succeeded but failed to update the status of the connection to %s.", sm.CurrentState)).WithMetadata(map[string]interface{}{"error": err}).FromSystem(*sysID).FromUser(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").Build(), err
	}

	sm.Log.Debug("HTTP status:", statusCode, "updated status for connection", connection.ID)

	// The action func only emits event when an error occurs.
	// If "event" is nil, it indicates actions were execeuted successfully, hence send an confirmation that request was processed successsfully.
	if event == nil {
		event = events.NewEvent().WithDescription(fmt.Sprintf("%s connection changed to %s", sm.Name, sm.CurrentState)).FromSystem(*sysID).FromUser(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").WithMetadata(map[string]interface{}{
			"previous_status": sm.PreviousState,
			"current_status":  sm.CurrentState,
		}).WithSeverity(events.Informational).Build()
	}

	return event, err
}
