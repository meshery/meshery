package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
)

func (r *mutationResolver) changeAddonStatus(ctx context.Context) (*model.Status, error) {
	return nil, nil
}

func (r *queryResolver) getAvailableAddons(ctx context.Context) ([]*model.AddonList, error) {
	return nil, nil
}

func (r *subscriptionResolver) listenToAddonEvents(ctx context.Context) (<-chan []*model.AddonList, error) {
	return nil, nil
}
