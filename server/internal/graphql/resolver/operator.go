package resolver

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/controllers"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

/**
	Contains resolvers for,
	1. Performing the synthetic test for the Operator and its controllers and returns their latest status.
	
	[Deprecated, the connection states should be used to control the behaviour [Connected/Disconnected]]
	2. Invoking action on the Operator (Provisoning/Deprovisioning)
**/
type operatorStatusK8sContext struct {
	ctxID      string
	processing interface{}
}

func (r *Resolver) changeOperatorStatus(ctx context.Context, provider models.Provider, status model.Status, ctxID string) (model.Status, error) {
	deleteOperator := true

	// Tell operator status subscription that operation is starting
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data: operatorStatusK8sContext{
			processing: true,
			ctxID:      ctxID,
		},
		Type: "health",
	})

	if status == model.StatusEnabled {
		r.Log.Info("Installing Operator")
		deleteOperator = false
	} else {
		r.Log.Info("Uninstalling Operator in context ", ctxID)
	}

	var kubeclient *mesherykube.Client
	var k8scontext models.K8sContext
	var err error
	if ctxID != "" {
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(allContexts) == 0 {
			r.Log.Error(ErrNilClient)
			return model.StatusUnknown, ErrNilClient
		}
		for _, ctx := range allContexts {
			if ctx.ID == ctxID {
				k8scontext = ctx
				break
			}
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return model.StatusUnknown, model.ErrMesheryClient(err)
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return model.StatusUnknown, model.ErrMesheryClientNil
		}
		k8scontext = k8scontexts[0]
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return model.StatusUnknown, model.ErrMesheryClient(err)
		}
	}
	if kubeclient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data: operatorStatusK8sContext{
				processing: ErrNilClient,
				ctxID:      ctxID,
			},
			Type: "error",
		})
		return model.StatusUnknown, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		if r.Config.OperatorTracker.DisableOperator { //Do not deploy operator is explicitly in disabled mode
			r.Log.Info("skipping operator deployment (in disabled mode)")
			return
		}
		op, _ := ctx.Value(models.MesheryControllerHandlersKey).(map[string]map[models.MesheryController]controllers.IMesheryController)

		var err error
		if del {
			err = op[ctxID][models.MesheryOperator].Undeploy()
		} else {
			err = op[ctxID][models.MesheryOperator].Deploy(true)
		}
		if err != nil {
			r.Log.Error(err)
			r.Broadcast.Submit(broadcast.BroadcastMessage{
				Source: broadcast.OperatorSyncChannel,
				Data: operatorStatusK8sContext{
					processing: err,
					ctxID:      ctxID,
				},
				Type: "error",
			})
			return
		}

		if del {
			r.Config.OperatorTracker.Undeployed(ctxID, true)
		} else {
			r.Config.OperatorTracker.Undeployed(ctxID, false)
		}

		r.Log.Info("Operator operation executed")

		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data: operatorStatusK8sContext{
				processing: false,
				ctxID:      ctxID,
			},
			Type: "health",
		})
		if !del {
			endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn, connectionTrackerSingleton)
			r.Log.Debug("Endpoint: ", endpoint)
			if err != nil {
				r.Log.Error(err)
				r.Broadcast.Submit(broadcast.BroadcastMessage{
					Source: broadcast.OperatorSyncChannel,
					Data: operatorStatusK8sContext{
						processing: err,
						ctxID:      ctxID,
					},
					Type: "health",
				})
				return
			}
			connectionTrackerSingleton.Set(k8scontext.ID, endpoint)
			r.Log.Info("Connected to broker at:", endpoint)
			connectionTrackerSingleton.Log(r.Log)
		}

		r.Log.Info("Meshsync operation executed")

		// r.operatorChannel <- &model.OperatorStatus{
		// 	Status: status,
		// }
	}(deleteOperator, kubeclient)

	return model.StatusProcessing, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context, _ models.Provider, connectionID string) (*model.MesheryControllersStatusListItem, error) {
	unknowStatus := &model.MesheryControllersStatusListItem{
		ConnectionID: connectionID,
		Status:       model.MesheryControllerStatusUnkown,
		Controller:   model.GetInternalController(models.MesheryOperator),
	}

	handler, ok := ctx.Value(models.HandlerKey).(*handlers.Handler)
	if !ok {
		return unknowStatus, nil
	}

	inst, ok := handler.ConnectionToStateMachineInstanceTracker.Get(uuid.FromStringOrNil(connectionID))
	// If machine instance is not present or points to nil, return unknown status
	if !ok || inst == nil {
		return unknowStatus, nil
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		r.Log.Error(model.ErrMesheryControllersStatusSubscription(err))
		return unknowStatus, nil
	}

	controllerhandler := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
	if !ok {
		return unknowStatus, nil
	}
	status := controllerhandler[models.MesheryOperator].GetStatus()
	return &model.MesheryControllersStatusListItem{
		ConnectionID: connectionID,
		Status:       model.GetInternalControllerStatus(status),
		Controller:   model.GetInternalController(models.MesheryOperator),
	}, nil
}

func (r *Resolver) getMeshsyncStatus(ctx context.Context, provider models.Provider, connectionID string) (*model.OperatorControllerStatus, error) {
	unknowStatus := &model.OperatorControllerStatus{
		ConnectionID: connectionID,
		Status:       model.Status(model.MesheryControllerStatusUnkown.String()),
		Name:         model.GetInternalController(models.Meshsync).String(),
	}

	handler, ok := ctx.Value(models.HandlerKey).(*handlers.Handler)
	if !ok {
		return unknowStatus, nil
	}

	inst, ok := handler.ConnectionToStateMachineInstanceTracker.Get(uuid.FromStringOrNil(connectionID))
	// If machine instance is not present or points to nil, return unknown status
	if !ok || inst == nil {
		return unknowStatus, nil
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if err != nil {
		r.Log.Error(model.ErrMesheryControllersStatusSubscription(err))
		return unknowStatus, nil
	}

	status := model.GetMeshSyncInfo(
		machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()[models.Meshsync],
		machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()[models.MesheryBroker],
		r.Log,
	)
	status.ConnectionID = connectionID
	return &status, nil
}

func (r *Resolver) getNatsStatus(ctx context.Context, provider models.Provider, connectionID string) (*model.OperatorControllerStatus, error) {

	unknowStatus := &model.OperatorControllerStatus{
		ConnectionID: connectionID,
		Status:       model.Status(model.MesheryControllerStatusUnkown.String()),
		Name:         model.GetInternalController(models.MesheryBroker).String(),
	}

	handler, ok := ctx.Value(models.HandlerKey).(*handlers.Handler)
	if !ok {
		return unknowStatus, nil
	}

	inst, ok := handler.ConnectionToStateMachineInstanceTracker.Get(uuid.FromStringOrNil(connectionID))
	// If machine instance is not present or points to nil, return unknown status
	if !ok || inst == nil {
		return unknowStatus, nil
	}

	machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)

	if err != nil {
		r.Log.Error(model.ErrMesheryControllersStatusSubscription(err))
		return unknowStatus, nil
	}

	status := model.GetBrokerInfo(
		machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()[models.MesheryBroker],
		r.Log,
	)
	status.ConnectionID = connectionID
	return &status, nil
}