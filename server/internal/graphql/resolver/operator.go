package resolver

import (
	"context"
	"sync"

	"github.com/go-errors/errors"
	operatorClient "github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/controllers"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

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

func (r *Resolver) getOperatorStatus(ctx context.Context, _ models.Provider, ctxID string) (*model.OperatorStatus, error) {
	status := model.StatusUnknown

	var kubeclient *mesherykube.Client
	var err error
	if ctxID != "" {
		k8scontexts, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return nil, model.ErrMesheryClientNil
		}
		for _, ctx := range k8scontexts {
			if ctx.ID == ctxID {
				kubeclient, err = ctx.GenerateKubeHandler()
				if err != nil {
					return nil, model.ErrMesheryClient(err)
				}
				break
			}
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return nil, model.ErrMesheryClientNil
		}
		kubeclient, err = k8scontexts[0].GenerateKubeHandler()
		if err != nil {
			return nil, model.ErrMesheryClient(err)
		}
	}
	if kubeclient == nil {
		return nil, model.ErrMesheryClientNil
	}
	name, version, err := model.GetOperator(kubeclient)
	if err != nil {
		r.Log.Error(err)
		return &model.OperatorStatus{
			Status: status,
			Error: &model.Error{
				Code:        "",
				Description: err.Error(),
			},
		}, nil
	}
	if name == "" {
		status = model.StatusDisabled
	} else {
		status = model.StatusEnabled
	}

	controllers, err := model.GetControllersInfo(kubeclient, r.BrokerConn)
	if err != nil {
		r.Log.Error(err)
		return &model.OperatorStatus{
			Status: status,
			Error: &model.Error{
				Code:        "",
				Description: err.Error(),
			},
		}, nil
	}

	return &model.OperatorStatus{
		Status:      status,
		Version:     version,
		Controllers: controllers,
	}, nil
}

func (r *Resolver) getMeshsyncStatus(ctx context.Context, provider models.Provider, connectionID string) (*model.OperatorControllerStatus, error) {
	var kubeclient *mesherykube.Client
	var err error
	
	tokenString := ctx.Value(models.TokenCtxKey).(string)

	if connectionID != "" {
		k8scontext, err := provider.GetK8sContext(tokenString, connectionID)
		if err != nil {
			return nil, model.ErrMesheryClientNil
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return nil, model.ErrMesheryClient(err)
		}
	}

	if kubeclient == nil {
		return nil, model.ErrMesheryClientNil
	}
	mesheryclient, err := operatorClient.New(&kubeclient.RestConfig)
	if err != nil {
		return nil, err
	}

	status, err := model.GetMeshSyncInfo(mesheryclient, kubeclient)
	if err != nil {
		return &status, err
	}
	return &status, nil
}

func (r *Resolver) getNatsStatus(ctx context.Context, provider models.Provider, connectionID string) (*model.OperatorControllerStatus, error) {
	var kubeclient *mesherykube.Client
	var err error
	
	tokenString := ctx.Value(models.TokenCtxKey).(string)

	if connectionID != "" {
		k8scontext, err := provider.GetK8sContext(tokenString, connectionID)
		if err != nil {
			return nil, model.ErrMesheryClientNil
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return nil, model.ErrMesheryClient(err)
		}
	}

	if kubeclient == nil {
		return nil, model.ErrMesheryClientNil
	}
	mesheryclient, err := operatorClient.New(&kubeclient.RestConfig)
	if err != nil {
		return nil, err
	}

	status, err := model.GetBrokerInfo(mesheryclient, kubeclient, r.BrokerConn)
	if err != nil {
		return &status, err
	}
	return &status, nil
}

func (r *Resolver) listenToOperatorsState(ctx context.Context, provider models.Provider, k8scontextIDs []string) (<-chan *model.OperatorStatusPerK8sContext, error) {
	operatorChannel := make(chan *model.OperatorStatusPerK8sContext)

	k8sctxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
	if !ok || len(k8sctxs) == 0 {
		return nil, ErrNilClient
	}
	var k8sContexts []models.K8sContext
	if len(k8scontextIDs) == 1 && k8scontextIDs[0] == "all" {
		k8sContexts = k8sctxs
	} else if len(k8scontextIDs) != 0 {
		var k8sContextIDsMap = make(map[string]bool)
		for _, k8sContext := range k8scontextIDs {
			k8sContextIDsMap[k8sContext] = true
		}
		for _, k8Context := range k8sctxs {
			if k8sContextIDsMap[k8Context.ID] {
				k8sContexts = append(k8sContexts, k8Context)
			}
		}
	}
	var group sync.WaitGroup
	for _, k8scontext := range k8sContexts {
		group.Add(1)
		go func(k8scontext models.K8sContext) {
			defer group.Done()
			operatorSyncChannel := make(chan broadcast.BroadcastMessage)
			r.Broadcast.Register(operatorSyncChannel)
			r.Log.Info("Operator subscription started for ", k8scontext.Name)

			// Enforce enable operator
			status, err := r.getOperatorStatus(ctx, provider, k8scontext.ID)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				return
			}
			statusWithContext := model.OperatorStatusPerK8sContext{
				ContextID:      k8scontext.ID,
				OperatorStatus: status,
			}
			operatorChannel <- &statusWithContext
			err = r.connectToBroker(ctx, provider, k8scontext.ID)
			if err != nil && err != ErrNoMeshSync {
				r.Log.Error(err)
				// The subscription should remain live to send future messages and only die when context is done
				// return
			}
			for {
				select {
				case processing := <-operatorSyncChannel:
					if processing.Source == broadcast.OperatorSyncChannel {
						r.Log.Info("Operator sync channel called for ", k8scontext.Name)
						status, err := r.getOperatorStatus(ctx, provider, k8scontext.ID)
						if err != nil {
							r.Log.Error(ErrOperatorSubscription(err))
							return
						}
						switch processing.Data.(type) {
						case operatorStatusK8sContext:
							if processing.Data.(operatorStatusK8sContext).ctxID != k8scontext.ID {
								continue
							}
							switch processing.Data.(operatorStatusK8sContext).processing.(type) {
							case bool:
								if processing.Data.(operatorStatusK8sContext).processing.(bool) {
									status.Status = model.StatusProcessing
								}
							case *errors.Error:
								status.Error = &model.Error{
									Code:        "",
									Description: processing.Data.(operatorStatusK8sContext).processing.(*errors.Error).Error(),
								}
							case error:
								status.Error = &model.Error{
									Code:        "",
									Description: processing.Data.(operatorStatusK8sContext).processing.(error).Error(),
								}
							}
						}
						statusWithContext := model.OperatorStatusPerK8sContext{
							ContextID:      k8scontext.ID,
							OperatorStatus: status,
						}
						operatorChannel <- &statusWithContext
					}
				case <-ctx.Done():
					r.Log.Info("Operator subscription flushed for ", k8scontext.Name)
					r.Broadcast.Unregister(operatorSyncChannel)
					close(operatorSyncChannel)

					return
				}
			}
		}(k8scontext)
	}
	go func() {
		group.Wait()
		close(operatorChannel)
	}()
	return operatorChannel, nil
}
