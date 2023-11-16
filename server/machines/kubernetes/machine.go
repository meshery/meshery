package machines

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/kubernetes/actions"
	"github.com/layer5io/meshkit/logger"
)

// One FSM per connection
// ID of the machine corresponds to ID of the machine.
// Mapping b/w Connection and its associated machine is created if not already exist in ConnectionStateMachineMap.
// Maintenance state is not considered.

var Discovered = machines.State{
	Events: machines.Events{
		machines.Register: machines.REGISTERED,
		machines.NotFound: machines.NOTFOUND,
		machines.Delete: machines.DELETED,
	},
	Action: &actions.DiscoverAction{},
}


var Registered = machines.State{
	Events: machines.Events{
		machines.Connect: machines.CONNECTED,
		machines.Ignore: machines.IGNORED,
	},
	Action: &actions.RegisterAction{},
}

var Connected = machines.State{
	Events: machines.Events{
		machines.Disconnect: machines.DISCONNECTED,
		machines.Delete: machines.DELETED,
		machines.NotFound: machines.NOTFOUND,
	},
	Action: &actions.ConnectAction{},
}

var Ignored = machines.State{
	Events: machines.Events{
		machines.Delete: machines.DELETED,
		machines.Register: machines.REGISTERED,
	},
	Action: &actions.IgnoreAction{},
}

var Disconnected = machines.State{
	Events: machines.Events{
		machines.Connect: machines.CONNECTED,
		machines.Delete: machines.DELETED,
	},
	Action: &actions.DisconnectAction{},
}

var NotFound = machines.State{
	Events: machines.Events{
		machines.Discovery: machines.DISCOVERED,
	},
	Action: nil,
}

var Delete = machines.State{
	Events: machines.Events{},
	Action: &actions.DeleteAction{},
}

func NewK8SMachine(initialState machines.StateType, connectionID uuid.UUID, machineName string, log logger.Handler) *machines.StateMachine {
	return &machines.StateMachine{
		ID: connectionID,
		Name: machineName,
		PreviousState: machines.DefaultState,
		InitialState: initialState,
		CurrentState: initialState,
		Context: nil,
		Log: log,
		States: machines.States{
			machines.DISCOVERED: Discovered, 
			machines.REGISTERED: Registered,
			machines.CONNECTED: Connected,
			machines.DISCONNECTED: Disconnected,
			machines.IGNORED: Ignored,
			machines.DELETED: Delete,

		},
	}
}