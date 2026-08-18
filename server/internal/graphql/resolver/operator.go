package resolver

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/handlers"
	"github.com/meshery/meshery/server/internal/graphql/model"
	mhelpers "github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/meshkit/utils/broadcast"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
)

/*
*

	Contains resolvers for,
	1. Performing the synthetic test for the Operator and its controllers and returns their latest status.

	[Deprecated, the connection states should be used to control the behaviour [Connected/Disconnected]]
	2. Invoking action on the Operator (Provisoning/Deprovisioning)

*
*/
type operatorStatusK8sContext struct {
	ctxID      string
	processing interface{}
}

// operatorControllersHelper resolves the controllers helper that owns Meshery
// Operator lifecycle for a Kubernetes context, by way of that connection's
// state machine - the same route the resync and controller-status paths take.
//
// Operator install and removal go through the helper rather than a meshkit
// handler read straight out of the request context, because the helper is what
// holds the resolved Helm chart version and refuses an install when no
// published version could be pinned. Reaching around it is how an unresolvable
// version would reach Helm from this entry point while every other one refused.
func operatorControllersHelper(ctx context.Context, k8scontext models.K8sContext) (*models.MesheryControllersHelper, error) {
	connectionDetail := "Kubernetes context " + k8scontext.ID + " (connection " + k8scontext.ConnectionID + ")"

	h, ok := ctx.Value(models.HandlerKey).(*handlers.Handler)
	if !ok || h == nil {
		return nil, ErrOperatorControllersHelper("No Meshery handler is present in the request context.")
	}
	tracker := h.ConnectionToStateMachineInstanceTracker
	if tracker == nil {
		return nil, ErrOperatorControllersHelper("No connection state machine tracker is available on the Meshery handler.")
	}
	if k8scontext.ConnectionID == "" {
		return nil, ErrOperatorControllersHelper("Kubernetes context " + k8scontext.ID + " carries no connection id.")
	}
	inst, ok := tracker.Get(uuid.FromStringOrNil(k8scontext.ConnectionID))
	if !ok || inst == nil || !mhelpers.HasMachineContext(inst) {
		return nil, ErrOperatorControllersHelper(connectionDetail + " has no ready state machine.")
	}
	machinectx, castErr := utils.Cast[*kubernetes.MachineCtx](inst.Context)
	if castErr != nil {
		return nil, ErrOperatorControllersHelper(castErr.Error())
	}
	if machinectx.MesheryCtrlsHelper == nil {
		return nil, ErrOperatorControllersHelper(connectionDetail + " has no controllers helper attached.")
	}
	return machinectx.MesheryCtrlsHelper, nil
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
		allContexts, ok := ctx.Value(models.AllKubeClusterKey).([]*models.K8sContext)
		if !ok || len(allContexts) == 0 {
			r.Log.Error(ErrNilClient)
			return model.StatusUnknown, ErrNilClient
		}
		for _, ctx := range allContexts {
			if ctx == nil {
				continue
			}
			if ctx.ID == ctxID {
				k8scontext = *ctx
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

	// Resolve the lifecycle owner before reporting the request as accepted, so a
	// connection that cannot be acted on fails the mutation outright instead of
	// reporting "processing" and then going quiet.
	ctrlHelper, err := operatorControllersHelper(ctx, k8scontext)
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
		return model.StatusUnknown, err
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		if r.Config.OperatorTracker.DisableOperator { //Do not deploy operator is explicitly in disabled mode
			r.Log.Info("skipping operator deployment (in disabled mode)")
			return
		}

		// SetOperatorDeployment carries the same guard the connect-time path
		// uses, so an unresolvable chart version is refused identically here
		// rather than being handed to Helm as an opaque chart-not-found. It also
		// re-resolves that version first, so clicking deploy again after the
		// chart repository recovers is a retry that can actually succeed.
		err := ctrlHelper.SetOperatorDeployment(k8scontext, !del)
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
