package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery/internal/graphql/generated"
	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

func (r *mutationResolver) ChangeAddonStatus(ctx context.Context, input *model.AddonStatusInput) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if input.Selector != nil {
		return r.changeAddonStatus(ctx, provider)
	}

	return model.StatusUnknown, ErrInvalidRequest
}

func (r *mutationResolver) ChangeOperatorStatus(ctx context.Context, input *model.OperatorStatusInput) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.changeOperatorStatus(ctx, provider, input.TargetStatus)
}

func (r *queryResolver) GetAvailableAddons(ctx context.Context, selector *model.MeshType) ([]*model.AddonList, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if selector != nil {
		return r.getAvailableAddons(ctx, provider, selector)
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

func (r *queryResolver) GetOperatorStatus(ctx context.Context) (*model.OperatorStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getOperatorStatus(ctx, provider)
}

func (r *queryResolver) ResyncCluster(ctx context.Context, selector *model.ReSyncActions) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.resyncCluster(ctx, provider, selector)
}

func (r *queryResolver) GetMeshsyncStatus(ctx context.Context) (*model.OperatorControllerStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getMeshsyncStatus(ctx, provider)
}

func (r *queryResolver) DeployMeshsync(ctx context.Context) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.deployMeshsync(ctx, provider)
}

func (r *queryResolver) GetNatsStatus(ctx context.Context) (*model.OperatorControllerStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getNatsStatus(ctx, provider)
}

func (r *queryResolver) ConnectToNats(ctx context.Context) (model.Status, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.connectToNats(ctx, provider)
}

func (r *queryResolver) GetAvailableNamespaces(ctx context.Context) ([]*model.NameSpace, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.getAvailableNamespaces(ctx, provider)
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

func (r *subscriptionResolver) ListenToAddonState(ctx context.Context, selector *model.MeshType) (<-chan []*model.AddonList, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if selector != nil {
		return r.listenToAddonState(ctx, provider, selector)
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

func (r *subscriptionResolver) ListenToOperatorState(ctx context.Context) (<-chan *model.OperatorStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.listenToOperatorState(ctx, provider)
}

func (r *subscriptionResolver) ListenToMeshSyncEvents(ctx context.Context) (<-chan *model.OperatorControllerStatus, error) {
	provider := ctx.Value(models.ProviderCtxKey).(models.Provider)
	return r.listenToMeshSyncEvents(ctx, provider)
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
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *subscriptionResolver) SubscribePerfResult(ctx context.Context, id string) (<-chan *model.MesheryResult, error) {
	panic(fmt.Errorf("not implemented"))
}
