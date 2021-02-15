package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery/internal/graphql/generated"
	"github.com/layer5io/meshery/internal/graphql/model"
)

func (r *mutationResolver) ChangeAddonStatus(ctx context.Context, selector *model.MeshType, targetStatus *model.Status) (*model.StatusResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ChangeOperatorStatus(ctx context.Context, targetStatus *model.Status) (*model.StatusResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAvailableAddons(ctx context.Context, selector *model.MeshType) ([]*model.AddonList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetControlPlanes(ctx context.Context, filter *model.ControlPlaneFilter) ([]*model.ControlPlane, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) ListenToAddonEvents(ctx context.Context, selector *model.MeshType) (<-chan []*model.AddonList, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) ListenToControlPlaneEvents(ctx context.Context, filter *model.ControlPlaneFilter) (<-chan []*model.ControlPlane, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *subscriptionResolver) ListenToOperatorEvents(ctx context.Context) (<-chan []*model.OperatorStatus, error) {
	panic(fmt.Errorf("not implemented"))
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
func (r *subscriptionResolver) GetControlPlanes(ctx context.Context, filter *model.ControlPlaneFilter) (<-chan []*model.ControlPlane, error) {
	panic(fmt.Errorf("not implemented"))
}
