package machines

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"
)

func Discovered() State {
	state := &State{}
	return *state.
		RegisterEvent(Register, REGISTERED).
		RegisterEvent(Ignore, IGNORED)
	// RegisterAction(DiscoverAction)
}

func Registered() State {
	state := &State{}
	return *state.
		RegisterEvent(Connect, CONNECTED).
		RegisterEvent(Ignore, IGNORED)
}

func Connected() State {
	state := &State{}
	return *state.
		RegisterEvent(Disconnect, DISCONNECTED)
}

func Initial() State {
	state := &State{}
	return *state.
		RegisterEvent(Discovery, DISCOVERED).
		RegisterEvent(Register, REGISTERED).
		RegisterEvent(Connect, CONNECTED)
}

func New(initialState StateType, ID string, log logger.Handler, mtype string) (*StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	if err != nil {
		return nil, ErrInititalizeK8sMachine(err)
	}

	return &StateMachine{
		ID:            connectionID,
		Name:          mtype,
		PreviousState: DefaultState,
		InitialState:  initialState,
		CurrentState:  initialState,
		Log:           log,
		States: States{
			DISCOVERED:   Discovered(),
			REGISTERED:   Registered(),
			CONNECTED:    Connected(),
			InitialState: Initial(),
		},
	}, nil
}
