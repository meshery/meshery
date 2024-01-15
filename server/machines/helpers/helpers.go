package helpers

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/grafana"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/machines/prometheus"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/logger"
)

func StatusToEvent(status connections.ConnectionStatus) machines.EventType {
	switch status {
	case connections.DISCOVERED:
		return machines.Discovery
	case connections.REGISTERED:
		return machines.Register
	case connections.CONNECTED:
		return machines.Connect
	case connections.DISCONNECTED:
		return machines.Disconnect
	case connections.IGNORED:
		return machines.Ignore
	case connections.DELETED:
		return machines.Delete
	case connections.NOTFOUND:
		return machines.NotFound
	}
	return machines.EventType(machines.DefaultState)
}
func getMachine(initialState machines.StateType, mtype, id string, userID uuid.UUID, log logger.Handler) (*machines.StateMachine, error) {
	switch mtype {
	case "kubernetes":
		return kubernetes.New(id, userID, log)
	case "grafana":
		mch, err := machines.New(initialState, id, userID, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		mch.States[machines.REGISTERED] = *register.RegisterAction(&grafana.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		mch.States[machines.CONNECTED] = *connect.RegisterAction(&machines.DefaultConnectAction{})
		return mch, nil
	case "prometheus":
		mch, err := machines.New(initialState, id, userID, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		mch.States[machines.REGISTERED] = *register.RegisterAction(&prometheus.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		mch.States[machines.CONNECTED] = *connect.RegisterAction(&machines.DefaultConnectAction{})

		return mch, nil
	}
	return nil, machines.ErrInvalidType(fmt.Errorf("invlaid type requested"))
}

func InitializeMachineWithContext(
	machineCtx interface{},
	ctx context.Context,
	ID uuid.UUID,
	userID uuid.UUID,
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker,
	log logger.Handler,
	provider models.Provider,
	initialState machines.StateType,
	mtype string,
	initFunc connections.InitFunc,
) (*machines.StateMachine, error) {
	inst, ok := smInstanceTracker.Get(ID)
	if ok {
		return inst, nil
	}

	inst, err := getMachine(initialState, mtype, ID.String(), userID, log)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	inst.Provider = provider
	_, err = inst.Start(ctx, machineCtx, log, initFunc)
	smInstanceTracker.Add(ID, inst)
	if err != nil {
		return nil, err
	}

	return inst, nil
}
