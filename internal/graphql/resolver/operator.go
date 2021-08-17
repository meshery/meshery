package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
)

func (r *Resolver) changeOperatorStatus(ctx context.Context, provider models.Provider, status model.Status) (model.Status, error) {
	delete := true
	if status == model.StatusEnabled {
		r.Log.Info("Installing Operator")
		delete = false
	}

	if r.Config.KubeClient.KubeClient == nil {
		r.Log.Error(ErrNilClient)
		return model.StatusUnknown, ErrNilClient
	}

	go func(del bool, kubeclient *mesherykube.Client) {
		err := model.Initialize(kubeclient, del)
		if err != nil {
			r.Log.Error(err)
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Operator operation executed")

		if !del {
			status, err := r.resyncCluster(context.TODO(), provider, &model.ReSyncActions{
				ReSync:  "false",
				ClearDb: "true",
			})
			if err != nil {
				r.Log.Error(err)
				r.operatorChannel <- &model.OperatorStatus{
					Status: status,
					Error: &model.Error{
						Code:        "",
						Description: err.Error(),
					},
				}
				return
			}

			endpoint, err := model.SubscribeToBroker(provider, kubeclient, r.brokerChannel, r.BrokerConn)
			r.Log.Debug("Endpoint: ", endpoint)
			if err != nil {
				r.Log.Error(err)
				r.operatorChannel <- &model.OperatorStatus{
					Status: status,
					Error: &model.Error{
						Code:        "",
						Description: err.Error(),
					},
				}
				return
			}
			r.Log.Info("Connected to broker at:", endpoint)
		}

		// installMeshsync
		err = model.RunMeshSync(kubeclient, del)
		if err != nil {
			r.Log.Error(err)
			r.operatorChannel <- &model.OperatorStatus{
				Status: status,
				Error: &model.Error{
					Code:        "",
					Description: err.Error(),
				},
			}
			return
		}
		r.Log.Info("Meshsync operation executed")

		r.operatorChannel <- &model.OperatorStatus{
			Status: status,
		}

		r.operatorSyncChannel <- false
	}(delete, r.Config.KubeClient)

	r.operatorChannel <- &model.OperatorStatus{
		Status: model.StatusProcessing,
	}

	return model.StatusProcessing, nil
}

func (r *Resolver) getOperatorStatus(ctx context.Context, provider models.Provider) (*model.OperatorStatus, error) {
	status := model.StatusUnknown
	version := string(model.StatusUnknown)
	if r.Config.KubeClient == nil {
		return nil, ErrMesheryClient(nil)
	}

	name, version, err := model.GetOperator(r.Config.KubeClient)
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

	controllers, err := model.GetControllersInfo(r.Config.KubeClient, r.BrokerConn, r.meshsyncLivenessChannel)
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

func (r *Resolver) listenToOperatorState(ctx context.Context, provider models.Provider) (<-chan *model.OperatorStatus, error) {
	if r.operatorChannel == nil {
		r.operatorChannel = make(chan *model.OperatorStatus)
	}
	if r.operatorSyncChannel == nil {
		r.operatorSyncChannel = make(chan bool)
	}
	if r.meshsyncLivenessChannel == nil {
		r.meshsyncLivenessChannel = make(chan struct{})
	}

	go func() {
		r.Log.Info("Operator subscription started")
		err := r.connectToBroker(context.TODO(), provider)
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}

		// Enforce enable operator
		status, err := r.getOperatorStatus(ctx, provider)
		if err != nil {
			r.Log.Error(ErrOperatorSubscription(err))
			return
		}
		if status.Status != model.StatusEnabled {
			_, err = r.changeOperatorStatus(ctx, provider, model.StatusEnabled)
			if err != nil {
				r.Log.Error(ErrOperatorSubscription(err))
				return
			}
		}
		for {
			select {
			case processing := <-r.operatorSyncChannel:
				r.Log.Info("Operator sync channel called")
				status, err := r.getOperatorStatus(ctx, provider)
				if err != nil {
					r.Log.Error(ErrOperatorSubscription(err))
					return
				}

				if processing {
					status.Status = model.StatusProcessing
				}

				r.operatorChannel <- status
			case <-ctx.Done():
				r.Log.Info("Operator subscription flushed")
				return
			}
		}
	}()

	return r.operatorChannel, nil
}
