package helpers

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/grafana"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/machines/prometheus"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
)

func StatusToEvent(status string) string {
	return strings.TrimSuffix(strings.ToLower(status), "ed")
}

func getMachine(initialState machines.StateType, mtype, id string, log logger.Handler) (*machines.StateMachine, error) {
	switch mtype {
	case "kubernetes":
		return kubernetes.New(id, log)
	case "grafana":
		mch, err := machines.New(initialState, id, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		mch.States[machines.REGISTERED] = *register.RegisterAction(&grafana.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		mch.States[machines.CONNECTED] = *connect.RegisterAction(&machines.DefaultConnectAction{})
		return mch, nil
	case "prometheus":
		mch, err := machines.New(initialState, id, log, mtype)
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
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker,
	log logger.Handler,
	provider models.Provider,
	initialState machines.StateType,
	mtype string,
	initFunc models.InitFunc,
) (*machines.StateMachine, error) {
	smInstanceTracker.Mx.Lock()
	defer smInstanceTracker.Mx.Unlock()

	inst, ok := smInstanceTracker.ConnectToInstanceMap[ID]
	if ok {
		return inst, nil
	}

	inst, err := getMachine(initialState, mtype, ID.String(), log)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	inst.Provider = provider
	_, err = inst.Start(ctx, machineCtx, log, initFunc)
	smInstanceTracker.ConnectToInstanceMap[ID] = inst
	if err != nil {
		return nil, err
	}

	return inst, nil
}