package kubernetes

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils/kubernetes"
)

// One FSM per connection
// ID of the machine corresponds to ID of the machine.
// Mapping b/w Connection and its associated machine is created if not already exist in ConnectionStateMachineMap.
// Maintenance state is not considered.

func Discovered(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Register: machines.REGISTERED,
			machines.NotFound: machines.NOTFOUND,
			machines.Delete:   machines.DELETED,
		},
		Action: &DiscoverAction{},
	}
}

func Registered(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Connect: machines.CONNECTED,
			machines.Ignore:  machines.IGNORED,
		},
		Action: &RegisterAction{},
	}
}

func Connected(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Disconnect: machines.DISCONNECTED,
			machines.Delete:     machines.DELETED,
			machines.NotFound:   machines.NOTFOUND,
		},
		Action: &ConnectAction{},
	}
}

func Ignored(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Delete:   machines.DELETED,
			machines.Register: machines.REGISTERED,
		},
		Action: &IgnoreAction{},
	}
}

func Disconnected(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Connect: machines.CONNECTED,
			machines.Delete:  machines.DELETED,
		},
		Action: &DisconnectAction{},
	}
}

func NotFound(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Discovery: machines.DISCOVERED,
			machines.Delete:    machines.DELETED,
		},
		Action: &NotFoundAction{},
	}
}

func Delete(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{},
		Action: &DeleteAction{},
	}
}

func Initial(log logger.Handler) machines.State {
	return machines.State{
		Events: machines.Events{
			machines.Discovery:  machines.DISCOVERED,
			machines.Register:   machines.REGISTERED,
			machines.Connect:    machines.CONNECTED,
			machines.Disconnect: machines.DISCONNECTED,
			machines.Ignore:     machines.IGNORED,
			machines.Delete:     machines.DELETED,
			machines.NotFound:   machines.NOTFOUND,
		},
		Action: nil,
	}
}

type MachineCtx struct {
	K8sContext         models.K8sContext
	MesheryCtrlsHelper *models.MesheryControllersHelper
	K8sCompRegHelper   *models.ComponentsRegistrationHelper
	EventBroadcaster   *models.Broadcast
	clientset          *kubernetes.Client
	log                logger.Handler
	Provider           models.Provider
	OperatorTracker    *models.OperatorTracker
	K8scontextChannel  *models.K8scontextChan
	RegistryManager    *meshmodel.RegistryManager
}

const (
	machineName = "kubernetes"
)

func NewK8SMachine(ID string, log logger.Handler) (*machines.StateMachine, error) {
	connectionID, err := uuid.FromString(ID)
	log.Info("initialising K8s machine for connetion Id", connectionID)
	if err != nil {
		return nil, machines.ErrInititalizeK8sMachine(err)
	}

	return &machines.StateMachine{
		ID:            connectionID,
		Name:          machineName,
		PreviousState: machines.DefaultState,
		InitialState:  machines.InitialState,
		CurrentState:  machines.InitialState,
		Log:           log,
		States: machines.States{
			machines.DISCOVERED:   Discovered(log),
			machines.REGISTERED:   Registered(log),
			machines.CONNECTED:    Connected(log),
			machines.DISCONNECTED: Disconnected(log),
			machines.IGNORED:      Ignored(log),
			machines.DELETED:      Delete(log),
			machines.NOTFOUND:     NotFound(log),
			machines.InitialState: Initial(log),
		},
	}, nil
}

func AssignInitialCtx(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error) {
	user, _ := ctx.Value(models.UserCtxKey).(*models.User)
	sysID, _ := ctx.Value(models.SystemIDKey).(*uuid.UUID)
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("register").FromSystem(*sysID).FromUser(userUUID) // pass userID and systemID in acted upon first pass user id if we can get context then update with connection Id
	machinectx, err := GetMachineCtx(machineCtx, eventBuilder)
	if err != nil {
		return nil, eventBuilder.Build(), err
	}
	err = AssignClientSetToContext(machinectx, eventBuilder)
	if err != nil {
		return nil, eventBuilder.Build(), err
	}
	machinectx.log = log
	return machinectx, nil, nil
}
