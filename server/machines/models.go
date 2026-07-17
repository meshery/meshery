package machines

import (
	"context"
	"fmt"
	"sync"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
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
	Exit       EventType = "exit"

	DISCOVERED   StateType = "discovered"
	REGISTERED   StateType = "registered"
	CONNECTED    StateType = "connected"
	IGNORED      StateType = "ignored"
	MAINTENANCE  StateType = "maintenance"
	DISCONNECTED StateType = "disconnected"
	DELETED      StateType = "deleted"
	NOTFOUND     StateType = "not found"

	Init EventType = "initialize"
)

var (
	DefaultState StateType = ""
	InitialState StateType = "initialized"
)

type Payload struct {
	Connection connections.Connection
	Credential models.Credential
}

// Represents an event in the system/machine
type EventType string

// Represents the mapping between event and the next state in the event's response
type Events map[EventType]StateType

type StateMachine struct {
	// ID to trace the events originated from the machine, also used in logs
	ID core.Uuid

	UserID core.Uuid

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

func (sm *StateMachine) AssignProvider(provider models.Provider) *StateMachine {
	sm.Provider = provider
	return sm
}

func (sm *StateMachine) Start(ctx context.Context, machinectx interface{}, log logger.Handler, init connections.InitFunc) (*events.Event, error) {
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
	sm.Log.Debug("inside getNextState: ", event, ok)
	if ok {
		events := state.Events
		if events != nil {
			nextState, ok := events[event]
			if ok {
				sm.Log.Debug("next state: ", nextState)
				return nextState, nil
			}
		}
	}
	return DefaultState, ErrInvalidTransitionEvent(sm.CurrentState, event)
}

// Returns events.Event and error. The func invoking the SendEvent should handle the error and publish the event.
// wherever possible use the userID and systemID from context as the events can be created from other comps or actors and not only user actors.
// In cases when the event is received as part of some other event and not explicitly created by an actor, use the useID and systemID of the actor who initially invoked the machine.
func (sm *StateMachine) SendEvent(ctx context.Context, eventType EventType, payload interface{}) (*events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*core.Uuid)
	userUUID := user.ID
	ctx = context.WithValue(ctx, models.ProviderCtxKey, sm.Provider)
	defaultEvent := events.NewEvent().WithDescription(fmt.Sprintf("Invalid status change requested to %s for connection type %s.", eventType, sm.Name)).ActedUpon(sm.ID).FromOwner(userUUID).FromSystem(*sysID).WithSeverity(events.Error)
	sm.mx.Lock()
	defer sm.mx.Unlock()
	var event *events.Event

	// eventType is reassigned inside the transition loop (to the recovery event,
	// then ultimately to NoOp), so capture the caller's original event up front.
	// It is needed after the loop to (a) skip the status update on an Exit event
	// and (b) restrict failure de-duplication to the background Discovery path.
	originalEventType := eventType

	// actionErr and actionErrEvent capture the FIRST genuine action failure: the
	// non-nil error returned by a state's entry/execute action, together with the
	// Error event that action built. A failing action may return a non-NoOp
	// "recovery" event (e.g. NotFound) instead of NoOp, so the loop does not
	// early-return; it transitions onward and the recovery state's
	// exit/entry/execute actions then reassign "event" (frequently to nil). That
	// previously masked the original failure and made SendEvent return an
	// Informational "success" with a nil error, so every caller that gates on
	// "if err != nil" persisted/published nothing for real connection failures.
	// Capturing here — and returning it after the status update — makes the
	// failure reach callers.
	//
	// NOTE: the loop-scoped "err" below is declared with ":=" inside the
	// for-block and shadows any function-scope err, so it MUST NOT be relied upon
	// for the terminal return. The captured values are the source of truth.
	var actionErr error
	var actionErrEvent *events.Event

	for eventType != NoOp {
		nextState, err := sm.getNextState(eventType)
		if err != nil {
			sm.Log.Error(err)
			event = defaultEvent.WithMetadata(map[string]interface{}{"error": err}).Build()
			sm.Log.Debug(defaultEvent.WithMetadata(map[string]interface{}{"error": err}).Build())
			break
		}

		sm.Log.Debugf("%s: transitioning to next state: %s", sm.Name, nextState)

		// next state to transition
		state, ok := sm.States[nextState]
		if !ok || state.Action == nil {
			sm.Log.Error(err)
			event = defaultEvent.WithMetadata(map[string]interface{}{"error": ErrInvalidTransition(sm.CurrentState, nextState)}).Build()
			sm.Log.Debug(event)
			break
		}

		// Execute exit actions before entering new state.
		action := sm.States[sm.CurrentState].Action
		if action != nil {
			_, event, err = action.ExecuteOnExit(ctx, sm.Context, nil)
			if err != nil {
				sm.Log.Error(err)
				return event, err
			}
		}

		if state.Action != nil {
			// Execute entry actions for the state entered.
			eventType, event, err = state.Action.ExecuteOnEntry(ctx, sm.Context, nil)
			sm.Log.Debugf("%s: entry action executed, event emitted %v", sm.Name, eventType)

			if err != nil {
				sm.Log.Error(err)
				sm.Log.Debug(event)
				if actionErr == nil {
					actionErr = err
					actionErrEvent = event
				}
				if eventType == NoOp {
					return event, err
				}
			} else {
				eventType, event, err = state.Action.Execute(ctx, sm.Context, payload)

				sm.Log.Debugf("%s: inside action executed, event emitted %v", sm.Name, eventType)
				if err != nil {
					sm.Log.Error(err)
					sm.Log.Debug(event)
					if actionErr == nil {
						actionErr = err
						actionErrEvent = event
					}
					if eventType == NoOp {
						return event, err
					}

				}
			}
		}

		sm.PreviousState = sm.CurrentState
		sm.CurrentState = nextState
	}

	// statusChanged tracks whether this send actually moves the persisted
	// connection status. K8sFSMMiddleware re-runs discovery on EVERY request for
	// each connected context, so once failures propagate to callers a cluster
	// that stays down would emit a persisted + broadcast Error event on every
	// request (notification spam). Emitting the failure only on an actual status
	// transition de-duplicates that noise while still surfacing the first
	// failure. Defaults to true so paths without a provider (or the Exit event)
	// still report failures.
	statusChanged := true

	if sm.Provider != nil && originalEventType != Exit {
		token, _ := ctx.Value(models.TokenCtxKey).(string)
		connection, _, err := sm.Provider.GetConnectionByID(token, sm.ID)

		if err != nil {

			return events.NewEvent().WithDescription(fmt.Sprintf("Failed to retrieve the connection with id %s to update status.", sm.ID)).WithMetadata(map[string]interface{}{"error": err}).FromSystem(*sysID).FromOwner(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").Build(), err
		}

		// Defensive guard: a provider that reports no error but hands back a nil
		// connection would panic on the field access below. Treat it as a
		// retrieval failure so the caller sees an error instead of a crash.
		if connection == nil {
			err = ErrConnectionNotFound(sm.ID.String())
			return events.NewEvent().WithDescription(fmt.Sprintf("Failed to retrieve the connection with id %s to update status.", sm.ID)).WithMetadata(map[string]interface{}{"error": err}).FromSystem(*sysID).FromOwner(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").Build(), err
		}

		// Compare the currently persisted status against the state the machine
		// settled in. Only an actual change should produce a user-facing failure
		// event (see statusChanged above).
		statusChanged = connection.Status != connections.ConnectionStatus(sm.CurrentState)

		connectionPayload := &connections.ConnectionPayload{
			ID:       sm.ID,
			Kind:     connection.Kind,
			MetaData: connection.Metadata,
			Status:   connections.ConnectionStatus(sm.CurrentState)}

		if connectionPayload.MetaData == nil {
			connectionPayload.MetaData = map[string]interface{}{}
		}

		connection, err = sm.Provider.UpdateConnectionById(token, connectionPayload, sm.ID.String())

		if err != nil {
			// In this case should the current state be again set to previous state i.e. should we rollback. But not only state should be rollback but other actions as well, rn we don't rollback state.
			return events.NewEvent().WithDescription(fmt.Sprintf("Operation succeeded but failed to update the status of the connection to %s.", sm.CurrentState)).WithMetadata(map[string]interface{}{"error": err}).FromSystem(*sysID).FromOwner(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").Build(), err
		}

		sm.Log.Debugf("%s: updated \"status\" for connection with id: %s to \"%s\"", sm.Name, connection.ID, sm.CurrentState)
	}

	// A genuine action failure occurred earlier in the transition loop. Return
	// its ORIGINAL Error event together with the non-nil error so callers persist
	// + broadcast it — this is the core fix for failures previously masked by a
	// later recovery transition. The connection-status update above has already
	// run, so the connection has been moved to its recovery state (e.g. NOTFOUND)
	// regardless of what we return here. Suppress the emit when the persisted
	// status did not change, so a repeatedly re-discovered down connection does
	// not spam a notification on every request.
	if actionErr != nil {
		// De-duplicate ONLY the background discovery path. K8sFSMMiddleware
		// re-runs SendEvent(Discovery) on every request, so a cluster that stays
		// down must not re-emit an Error event once its status is already the
		// recovery state. User-initiated events (Register, Connect, ...) always
		// propagate their failure — a manual action that fails must be reported
		// even when the persisted status does not change (e.g. it was already
		// NOTFOUND/DISCONNECTED), otherwise the request silently "succeeds".
		if !statusChanged && originalEventType == Discovery {
			return nil, nil
		}
		// Guarantee a non-nil event on the error return: callers such as the K8s
		// middleware and addK8SConfig dereference the event (*event) whenever the
		// error is non-nil. Real actions always build an Error event on failure,
		// but fall back to a generic Error event so a misbehaving action can never
		// turn a surfaced failure into a nil-pointer panic.
		if actionErrEvent == nil {
			actionErrEvent = defaultEvent.WithMetadata(map[string]interface{}{"error": actionErr}).Build()
		}
		return actionErrEvent, actionErr
	}

	// The action func only emits event when an error occurs.
	// If "event" is nil, it indicates actions were execeuted successfully, hence send an confirmation that request was processed successsfully.
	if event == nil {
		event = events.NewEvent().WithDescription(fmt.Sprintf("%s connection changed to %s", sm.Name, sm.CurrentState)).FromSystem(*sysID).FromOwner(userUUID).ActedUpon(sm.ID).WithCategory("connection").WithAction("update").WithMetadata(map[string]interface{}{
			"previousStatus": sm.PreviousState,
			"currentStatus":  sm.CurrentState,
		}).WithSeverity(events.Informational).Build()
	}

	return event, nil
}
