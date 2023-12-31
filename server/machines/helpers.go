package machines

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines/grafana"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/machines/prometheus"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/machines"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshsync/pkg/model"
)

func StatusToEvent(status string) string {
	return strings.TrimSuffix(strings.ToLower(status), "ed")
}

func getMachine(initialState machines.StateType, mtype, id string, log logger.Handler) (*machines.StateMachine, error) {
	switch mtype {
	case "kubernetes":
		return kubernetes.New(id, log)
	case "grafana":
		mch, err := New(initialState, id, log, mtype)
		if err != nil {
			return mch, err
		}
		register := mch.States[machines.REGISTERED]
		mch.States[machines.REGISTERED] = *register.RegisterAction(&grafana.RegisterAction{})

		connect := mch.States[machines.CONNECTED]
		mch.States[machines.CONNECTED] = *connect.RegisterAction(&machines.DefaultConnectAction{})
		return mch, nil
	case "prometheus":
		mch, err := New(initialState, id, log, mtype)
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


// func processRegistration() {
// 	arh := helpers.GetAutoRegistrationHelperSingleton()
// 	if arh == nil {
// 		return
// 	}

// 	for obj := range arh.Queue.Ch {
// 		go func(obj *model.KubernetesResource) {
// 			// Ideally iterate all Connection defs, extract fingerprint composite key and try to match with the given obj,
// 			// For all connections that match the fingerprint and autoRegsiter is set to true, try to do auto registration.

// 			// connectionDefs, _, _ := v1alpha1.GetMeshModelComponents(&arh.dbHandler, *connectionCompFilter)
// 			// for _, connectionDef := range connectionDefs {
// 			// 	capabilities, err := utils.Cast[map[string]interface{}](connectionDef.Metadata["capabilities"])
// 			// 	if err != nil {
// 			// 		arh.log.Error(err)
// 			// 		continue
// 			// 	}
// 			// 	autoRegister, ok := capabilities["autoRegister"].(bool)
// 			// 	if ok && autoRegister {
// 			// 		fmt.Println("TEST:: inside for loop extracted capabilities.autoRegister")
// 			// 		// ch
// 			// 	}
// 			// }

// 			// For now, the auto-registration for Prometheus/Grafana is hard-coded.
// 			connType := getTypeOfConnection(obj)
// 			if connType != "" {
// 				id, _ := uuid.NewV4() // id should be hash of somehting.
// 				machineInst, err := InitializeMachineWithContext(nil, context.TODO(), id, arh.SMInstanceTracker, arh.Log, nil, machines.DISCOVERED, connType, nil)
// 				if err != nil {
// 					arh.Log.Error(ErrAutoRegister(err, connType))
// 				}
// 				machineInst.Provider = nil // set provider somehow
// 			}
// 		}(&obj)
// 	}
// }


func getTypeOfConnection(obj *model.KubernetesResource) string {
	if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
		return "grafana"
	} else if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") {
		return "prometheus"
	}
	return ""
}