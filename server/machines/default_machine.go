package machines

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/logger"
)

func Discovered() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Register, machines.REGISTERED).
		RegisterEvent(machines.Ignore, machines.IGNORED)
	// RegisterAction(machines.DiscoverAction)
}

func Registered() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Connect, machines.CONNECTED).
		RegisterEvent(machines.Ignore, machines.IGNORED)
}

func Connected() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Disconnect, machines.DISCONNECTED)
}

func Initial() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Discovery, machines.DISCOVERED).
		RegisterEvent(machines.Register, machines.REGISTERED).
		RegisterEvent(machines.Connect, machines.CONNECTED)
}

func New(initialState machines.StateType, ID string, log logger.Handler, mtype string) (*machines.StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	if err != nil {
		return nil, machines.ErrInititalizeK8sMachine(err)
	}

	return &machines.StateMachine{
		ID:            connectionID,
		Name:          mtype,
		PreviousState: machines.DefaultState,
		InitialState:  initialState,
		CurrentState:  initialState,
		Log:           log,
		States: machines.States{
			machines.DISCOVERED:   Discovered(),
			machines.REGISTERED:   Registered(),
			machines.CONNECTED:    Connected(),
			machines.InitialState: Initial(),
		},
	}, nil
}
