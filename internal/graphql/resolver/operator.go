package resolver

import (
	"context"

	operatorClient "github.com/layer5io/meshery-operator/pkg/client"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/errors"
	"github.com/layer5io/meshkit/utils/broadcast"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (r *Resolver) changeOperatorStatus(ctx context.Context, provider models.Provider, status model.Status, ctxID string) (model.Status, error) {
	delete := true

	// Tell operator status subscription that operation is starting
	r.Broadcast.Submit(broadcast.BroadcastMessage{
		Source: broadcast.OperatorSyncChannel,
		Data:   true,
		Type:   "health",
	})

	if status == model.StatusEnabled {
		r.Log.Info("Installing Operator")
		delete = false
	}

	var kubeclient *mesherykube.Client
	var k8scontext models.K8sContext
	var err error
	if ctxID != "" {
		k8scontext, err = provider.GetK8sContext("give token", ctxID)
		if err != nil {
			return model.StatusUnknown, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return model.StatusUnknown, ErrMesheryClient(err)
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return model.StatusUnknown, ErrMesheryClient(nil)
		}
		k8scontext = k8scontexts[0]
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return model.StatusUnknown, ErrMesheryClient(err)
		}
	}
	if kubeclient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data:   ErrNilClient,
			Type:   "error",
		})
		return model.StatusUnknown, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		err := model.Initialize(kubeclient, del, r.Config.AdapterTracker)
		if err != nil {
			r.Log.Error(err)
			r.Broadcast.Submit(broadcast.BroadcastMessage{
				Source: broadcast.OperatorSyncChannel,
				Data:   err,
				Type:   "error",
			})
			return
		}
		r.Log.Info("Operator operation executed")

		if !del {
			_, err := r.resyncCluster(context.TODO(), provider, &model.ReSyncActions{
				ReSync:  "false",
				ClearDb: "true",
			})
			if err != nil {
				r.Log.Error(err)
				r.Broadcast.Submit(broadcast.BroadcastMessage{
					Source: broadcast.OperatorSyncChannel,
					Data:   false,
					Type:   "health",
				})
				return
			}

			endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn, connectionTrackerSingleton)
			r.Log.Debug("Endpoint: ", endpoint)
			if err != nil {
				r.Log.Error(err)
				r.Broadcast.Submit(broadcast.BroadcastMessage{
					Source: broadcast.OperatorSyncChannel,
					Data:   false,
					Type:   "health",
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

		r.Broadcast.Submit(broadcast.BroadcastMessage{
			Source: broadcast.OperatorSyncChannel,
			Data:   false,
			Type:   "health",
		})
	}(delete, kubeclient)

	return model.StatusProcessing, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context, provider models.Provider, ctxID string) (*model.OperatorStatus, error) {
	status := model.StatusUnknown
	version := string(model.StatusUnknown)

	var kubeclient *mesherykube.Client
	var err error
	if ctxID != "" {
		k8scontext, err := provider.GetK8sContext("give token", ctxID)
		if err != nil {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontexts[0].GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
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

	controllers, err := model.GetControllersInfo(kubeclient, r.BrokerConn, r.meshsyncLivenessChannel)
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

func (r *Resolver) getMeshsyncStatus(ctx context.Context, provider models.Provider, selector *model.K8sContext) (*model.OperatorControllerStatus, error) {
	var kubeclient *mesherykube.Client
	var err error
	ctxID := *selector.ID
	if ctxID != "" {
		k8scontext, err := provider.GetK8sContext("give token", ctxID)
		if err != nil {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontexts[0].GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
	}
	mesheryclient, err := operatorClient.New(&kubeclient.RestConfig)
	if err != nil {
		return nil, err
	}

	status, err := model.GetMeshSyncInfo(mesheryclient, kubeclient, r.meshsyncLivenessChannel)
	if err != nil {
		return &status, err
	}
	return &status, nil
}

func (r *Resolver) getNatsStatus(ctx context.Context, provider models.Provider, selector *model.K8sContext) (*model.OperatorControllerStatus, error) {
	var kubeclient *mesherykube.Client
	var err error
	ctxID := *selector.ID
	if ctxID != "" {
		k8scontext, err := provider.GetK8sContext("give token", ctxID)
		if err != nil {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontext.GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
	} else {
		k8scontexts, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			return nil, ErrMesheryClient(nil)
		}
		kubeclient, err = k8scontexts[0].GenerateKubeHandler()
		if err != nil {
			return nil, ErrMesheryClient(err)
		}
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

func (r *Resolver) listenToOperatorState(ctx context.Context, provider models.Provider, selector *model.K8sContext) (<-chan *model.OperatorStatus, error) {
	operatorChannel := make(chan *model.OperatorStatus)

	if r.operatorSyncChannel == nil {
		r.operatorSyncChannel = make(chan bool)
	}
	if r.meshsyncLivenessChannel == nil {
		r.meshsyncLivenessChannel = make(chan struct{})
	}

	operatorSyncChannel := make(chan broadcast.BroadcastMessage)
	r.Broadcast.Register(operatorSyncChannel)

	go func() {
		r.Log.Info("Operator subscription started")
		err := r.connectToBroker(ctx, provider, *selector.ID)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			// The subscription should remain live to send future messages and only die when context is done
			// return
		}

		// Enforce enable operator
		status, err := r.getOperatorStatus(ctx, provider, *selector.ID)
		if err != nil {
			r.Log.Error(ErrOperatorSubscription(err))
			return
		}
		if status.Status != model.StatusEnabled {
			_, err = r.changeOperatorStatus(ctx, provider, model.StatusEnabled, *selector.ID)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				// return
			}
		}
		for {
			select {
			case processing := <-operatorSyncChannel:
				r.Log.Info("Operator sync channel called")
				status, err := r.getOperatorStatus(ctx, provider, *selector.ID)
				if err != nil {
					r.Log.Error(ErrOperatorSubscription(err))
					r.Log.Info("Operator subscription flushed")
					close(operatorChannel)
					// return
					continue
				}

				if processing.Source == broadcast.OperatorSyncChannel {
					switch processing.Data.(type) {
					case bool:
						if processing.Data.(bool) {
							status.Status = model.StatusProcessing
						}
					case *errors.Error:
						status.Error = &model.Error{
							Code:        processing.Data.(*errors.Error).Code,
							Description: processing.Data.(*errors.Error).Error(),
						}
					case error:
						status.Error = &model.Error{
							Code:        "",
							Description: processing.Data.(error).Error(),
						}
					}
				}
				operatorChannel <- status
			case <-ctx.Done():
				r.Log.Info("Operator subscription flushed")
				close(operatorChannel)
				r.Broadcast.Unregister(operatorSyncChannel)
				close(operatorSyncChannel)
				return
			}
		}
	}()

	return operatorChannel, nil
}
