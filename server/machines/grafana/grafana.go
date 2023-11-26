package grafana

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
)

const grafana = "grafana"

type GrafanaConn struct {
	URL  string
	Name string
}

type GrafanaCred struct {
	Name string `json:"name,omitempty"`
	// If Basic then it should be formatted as username:password
	APIKeyOrBasicAuth string `json:"credential,omitempty"`
}

type GrafanaPayload struct {
	GrafanaConn
	GrafanaCred
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
		RegisterAction(&RegisterAction{})
}

func Initial() machines.State {
	state := &machines.State{}
	return *state.
		RegisterEvent(machines.Discovery, machines.DISCOVERED).
		RegisterEvent(machines.Register, machines.REGISTERED).
		RegisterEvent(machines.Connect, machines.CONNECTED)
}

type MachineCtx struct {
	GrafanaClient *models.GrafanaClient
	GrafanaConn   GrafanaConn
	GrafanaCred   GrafanaCred
	provider      models.Provider
}

func New(initialState machines.StateType, ID string, log logger.Handler) (*machines.StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	if err != nil {
		return nil, machines.ErrInititalizeK8sMachine(err)
	}
	log.Info("initialising grafana machine for connetion Id", connectionID)

	return &machines.StateMachine{
		ID:            connectionID,
		Name:          grafana,
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
