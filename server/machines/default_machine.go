package machines

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
)

func Discovered() State {
	state := &State{}
	return *state.
		RegisterEvent(Register, REGISTERED).
		RegisterEvent(Ignore, IGNORED)
	// RegisterAction(DiscoverAction)
}

// Registered carries a no-op action so kinds without a verification step can
// enter REGISTERED: SendEvent rejects a transition into a state with no action.
// getMachine replaces it with a kind-specific action (grafana, prometheus).
func Registered() State {
	state := &State{
		Action: &DefaultAction{},
	}
	return *state.
		RegisterEvent(Connect, CONNECTED).
		RegisterEvent(Ignore, IGNORED)
}

func Connected() State {
	state := &State{}
	return *state.
		RegisterEvent(Disconnect, DISCONNECTED).
		RegisterEvent(Delete, DELETED)
}

func Disconnected() State {
	state := &State{}
	return *state.
		RegisterEvent(Connect, CONNECTED).
		RegisterEvent(Delete, DELETED)
}

func Initial() State {
	state := &State{}
	return *state.
		RegisterEvent(Discovery, DISCOVERED).
		RegisterEvent(Register, REGISTERED).
		RegisterEvent(Connect, CONNECTED)
}

func New(initialState StateType, ID string, userID core.Uuid, log logger.Handler, mtype string) (*StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	if err != nil {
		return nil, ErrInititalizeK8sMachine(err)
	}

	return &StateMachine{
		ID:            connectionID,
		UserID:        userID,
		Name:          mtype,
		PreviousState: DefaultState,
		InitialState:  initialState,
		CurrentState:  initialState,
		Log:           log,
		States: States{
			DISCOVERED:   Discovered(),
			REGISTERED:   Registered(),
			CONNECTED:    Connected(),
			DISCONNECTED: Disconnected(),
			InitialState: Initial(),
		},
	}, nil
}
