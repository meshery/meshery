package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/layer5io/meshery/server/internal/graphql/generated"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/models/controllers"
)

func (r *mutationResolver) ChangeOperatorStatus(ctx context.Context, input *model.OperatorStatusInput) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.changeOperatorStatus(ctx, provider, input.TargetStatus, input.ContextID)
}

func (r *queryResolver) GetAvailableAddons(ctx context.Context, filter *model.ServiceMeshFilter) ([]*model.AddonList, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.getAvailableAddons(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *queryResolver) GetControlPlanes(ctx context.Context, filter *model.ServiceMeshFilter) ([]*model.ControlPlane, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.getControlPlanes(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *queryResolver) GetDataPlanes(ctx context.Context, filter *model.ServiceMeshFilter) ([]*model.DataPlane, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.getDataPlanes(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *queryResolver) GetOperatorStatus(ctx context.Context, k8scontextID string) (*model.OperatorStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getOperatorStatus(ctx, provider, k8scontextID)
}

func (r *queryResolver) ResyncCluster(ctx context.Context, selector *model.ReSyncActions, k8scontextID string) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.resyncCluster(ctx, provider, selector, k8scontextID)
}

func (r *queryResolver) GetMeshsyncStatus(ctx context.Context, k8scontextID string) (*model.OperatorControllerStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getMeshsyncStatus(ctx, provider, k8scontextID)
}

func (r *queryResolver) DeployMeshsync(ctx context.Context, k8scontextID string) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.deployMeshsync(ctx, provider, k8scontextID)
}

func (r *queryResolver) GetNatsStatus(ctx context.Context, k8scontextID string) (*model.OperatorControllerStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getNatsStatus(ctx, provider, k8scontextID)
}

func (r *queryResolver) ConnectToNats(ctx context.Context, k8scontextID string) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.connectToNats(ctx, provider, k8scontextID)
}

func (r *queryResolver) GetAvailableNamespaces(ctx context.Context, k8sClusterIDs []string) ([]*model.NameSpace, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getAvailableNamespaces(ctx, provider, k8sClusterIDs)
}

func (r *queryResolver) GetPerfResult(ctx context.Context, id string) (*model.MesheryResult, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getPerfResult(ctx, provider, id)
	// panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FetchResults(ctx context.Context, selector model.PageFilter, profileID string) (*model.PerfPageResult, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.fetchResults(ctx, provider, selector, profileID)
}

func (r *queryResolver) GetPerformanceProfiles(ctx context.Context, selector model.PageFilter) (*model.PerfPageProfiles, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getPerformanceProfiles(ctx, provider, selector)
}

func (r *queryResolver) FetchAllResults(ctx context.Context, selector model.PageFilter) (*model.PerfPageResult, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.fetchAllResults(ctx, provider, selector)
}

func (r *queryResolver) FetchPatterns(ctx context.Context, selector model.PageFilter) (*model.PatternPageResult, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.fetchPatterns(ctx, provider, selector)
}

func (r *queryResolver) GetWorkloads(ctx context.Context, name *string, id *string, trim *bool) ([]*model.OAMCapability, error) {
	return r.getWorkloads(ctx, name, id, trim)
}

func (r *queryResolver) GetTraits(ctx context.Context, name *string, id *string, trim *bool) ([]*model.OAMCapability, error) {
	return r.getTraits(ctx, name, id, trim)
}

func (r *queryResolver) GetScopes(ctx context.Context, name *string, id *string, trim *bool) ([]*model.OAMCapability, error) {
	return r.getScopes(ctx, name, id, trim)
}

func (r *queryResolver) GetKubectlDescribe(ctx context.Context, name string, kind string, namespace string) (*model.KctlDescribeDetails, error) {
	return r.getKubectlDescribe(ctx, name, kind, namespace)
}

func (r *queryResolver) FetchPatternCatalogContent(ctx context.Context, selector *model.CatalogSelector) ([]*model.CatalogPattern, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.fetchCatalogPattern(ctx, provider, selector)
}

func (r *queryResolver) FetchFilterCatalogContent(ctx context.Context, selector *model.CatalogSelector) ([]*model.CatalogFilter, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.fetchCatalogFilter(ctx, provider, selector)
}

func (r *queryResolver) GetClusterResources(ctx context.Context, k8scontextIDs []string, namespace string) (*model.ClusterResources, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getClusterResources(ctx, provider, k8scontextIDs, namespace)
}

func (r *subscriptionResolver) ListenToAddonState(ctx context.Context, filter *model.ServiceMeshFilter) (<-chan []*model.AddonList, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.listenToAddonState(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *subscriptionResolver) ListenToControlPlaneState(ctx context.Context, filter *model.ServiceMeshFilter) (<-chan []*model.ControlPlane, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.listenToControlPlaneState(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *subscriptionResolver) ListenToDataPlaneState(ctx context.Context, filter *model.ServiceMeshFilter) (<-chan []*model.DataPlane, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if filter != nil {
		return r.listenToDataPlaneState(ctx, provider, filter)
	}

	return nil, ErrInvalidRequest
}

func (r *subscriptionResolver) ListenToOperatorState(ctx context.Context, k8scontextIDs []string) (<-chan *model.OperatorStatusPerK8sContext, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.listenToOperatorsState(ctx, provider, k8scontextIDs)
}

func (r *subscriptionResolver) ListenToMeshSyncEvents(ctx context.Context, k8scontextIDs []string) (<-chan *model.OperatorControllerStatusPerK8sContext, error) {
	// provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	// return r.listenToMeshSyncEvents(ctx, provider)
	return nil, nil
}

func (r *subscriptionResolver) SubscribePerfProfiles(ctx context.Context, selector model.PageFilter) (<-chan *model.PerfPageProfiles, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.subscribePerfProfiles(ctx, provider, selector)
}

func (r *subscriptionResolver) SubscribePerfResults(ctx context.Context, selector model.PageFilter, profileID string) (<-chan *model.PerfPageResult, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.subscribePerfResults(ctx, provider, selector, profileID)
}

func (r *subscriptionResolver) SubscribeBrokerConnection(ctx context.Context) (<-chan bool, error) {
	return r.subscribeBrokerConnection(ctx)
}

func (r *subscriptionResolver) SubscribeMesheryControllersStatus(ctx context.Context, k8scontextIDs []string) (<-chan []*model.MesheryControllersStatusListItem, error) {
	resChan := make(chan []*model.MesheryControllersStatusListItem)
	controllerHandlersPerContext, ok := ctx.Value(models.MesheryControllerHandlersKey).(map[string]map[models.MesheryController]controllers.IMesheryController)
	if !ok || len(controllerHandlersPerContext) == 0 || controllerHandlersPerContext == nil {
		er := model.ErrMesheryControllersStatusSubscription(fmt.Errorf("controller handlers are not configured for any of the contexts"))
		r.Log.Error(er)
		return nil, er
	}
	statusMapPerCtx := make(map[string]map[models.MesheryController]controllers.MesheryControllerStatus)
	// initialize the map
	for ctxID, ctrlHandlers := range controllerHandlersPerContext {
		for controller, handler := range ctrlHandlers {
			if _, ok := statusMapPerCtx[ctxID]; !ok {
				statusMapPerCtx[ctxID] = make(map[models.MesheryController]controllers.MesheryControllerStatus)
			}
			statusMapPerCtx[ctxID][controller] = handler.GetStatus()
		}
	}
	go func() {
		ctrlsStatusList := make([]*model.MesheryControllersStatusListItem, 0)
		// first send the initial status of the controllers
		for ctxID, statusMap := range statusMapPerCtx {
			for controller, status := range statusMap {
				ctrlsStatusList = append(ctrlsStatusList, &model.MesheryControllersStatusListItem{
					ContextID:  ctxID,
					Controller: model.GetInternalController(controller),
					Status:     model.GetInternalControllerStatus(status),
				})
			}
		}
		resChan <- ctrlsStatusList
		ctrlsStatusList = make([]*model.MesheryControllersStatusListItem, 0)
		// do this every 5 seconds
		for {
			for ctxID, ctrlHandlers := range controllerHandlersPerContext {
				for controller, handler := range ctrlHandlers {
					newStatus := handler.GetStatus()
					// if the status has changed, send that to the subscription
					if newStatus != statusMapPerCtx[ctxID][controller] {
						ctrlsStatusList = append(ctrlsStatusList, &model.MesheryControllersStatusListItem{
							ContextID:  ctxID,
							Controller: model.GetInternalController(controller),
							Status:     model.GetInternalControllerStatus(newStatus),
						})
						resChan <- ctrlsStatusList
					}
					// update the status list with newStatus
					statusMapPerCtx[ctxID][controller] = newStatus
					ctrlsStatusList = make([]*model.MesheryControllersStatusListItem, 0)
				}
			}
			// establish a watch connection to get updates, ideally in meshery-operator
			time.Sleep(time.Second * 5)
		}
	}()
	return resChan, nil
}

func (r *subscriptionResolver) SubscribeMeshSyncEvents(ctx context.Context, k8scontextIDs []string) (<-chan *model.MeshSyncEvent, error) {
	resChan := make(chan *model.MeshSyncEvent)
	// get handlers
	meshSyncDataHandlers, ok := ctx.Value(models.MeshSyncDataHandlersKey).(map[string]models.MeshsyncDataHandler)
	if !ok || len(meshSyncDataHandlers) == 0 || meshSyncDataHandlers == nil {
		er := model.ErrMeshSyncEventsSubscription(fmt.Errorf("meshsync data handlers are not configured for any of the contexts"))
		r.Log.Error(er)
		return nil, er
	}
	for ctxID, dataHandler := range meshSyncDataHandlers {
		brokerEventsChan := make(chan *broker.Message)
		err := dataHandler.ListenToMeshSyncEvents(brokerEventsChan)
		if err != nil {
			r.Log.Warn(err)
			r.Log.Info("skipping meshsync events subscription for contexId: %s", ctxID)
			continue
		}
		go func(ctxID string, brokerEventsChan chan *broker.Message) {
			for event := range brokerEventsChan {
				if event.EventType == broker.ErrorEvent {
					// TODO: Handle errors accordingly
					continue
				}
				// handle the events
				res := &model.MeshSyncEvent{
					ContextID: ctxID,
					Type:      string(event.EventType),
					Object:    event.Object,
				}
				resChan <- res
				go r.Config.DashboardK8sResourcesChan.PublishDashboardK8sResources()
			}
		}(ctxID, brokerEventsChan)
	}
	return resChan, nil
}

func (r *subscriptionResolver) SubscribeConfiguration(ctx context.Context, applicationSelector model.PageFilter, patternSelector model.PageFilter, filterSelector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.subscribeConfiguration(ctx, provider, applicationSelector, patternSelector, filterSelector)
}

func (r *subscriptionResolver) SubscribeClusterResources(ctx context.Context, k8scontextIDs []string, namespace string) (<-chan *model.ClusterResources, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.subscribeClusterResources(ctx, provider, k8scontextIDs, namespace)
}

func (r *subscriptionResolver) SubscribeK8sContext(ctx context.Context, selector model.PageFilter) (<-chan *model.K8sContextsPage, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.subscribeK8sContexts(ctx, provider, selector)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *subscriptionResolver) SubscribePerfResult(ctx context.Context, id string) (<-chan *model.MesheryResult, error) {
	panic(fmt.Errorf("not implemented"))
}
