package prometheus

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
)

const prometheus = "prometheus"

type PromConn struct {
	URL  string
	Name string
}

type PromCred struct {
	Name string
	// Default strategy
	APIKey string
	// formatted as username:password
	BasicAuth string
}

type PrometheusPayload struct {
	PromConn
	PromCred
}

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
		RegisterEvent(machines.Connect, machines.REGISTERED).
		RegisterEvent(machines.Ignore, machines.IGNORED).
		RegisterAction(&RegisterAction{})
}

func Connected() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Disconnect, machines.DISCONNECTED).
		RegisterAction(&machines.DefaultConnectAction{})
}

func Initial() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Discovery, machines.DISCOVERED).
		RegisterEvent(machines.Register, machines.REGISTERED).
		RegisterEvent(machines.Connect, machines.CONNECTED)
}

type MachineCtx struct {
	PromClient *models.PrometheusClient
	PromConn   PromConn
	PromCred   PromCred
	provider   models.Provider
}

func New(initialState machines.StateType, ID string, log logger.Handler) (*machines.StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	log.Info("initialising prometheus machine for connection Id", connectionID)
	if err != nil {
		return nil, machines.ErrInititalizeK8sMachine(err)
	}

	return &machines.StateMachine{
		ID:            connectionID,
		Name:          "prometheus",
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
