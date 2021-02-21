package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
)

func (r *Resolver) getControlPlanes(ctx context.Context) ([]*model.ControlPlane, error) {
	return nil, nil
}

func (r *Resolver) listenToControlPlaneEvents(ctx context.Context) (<-chan []*model.ControlPlane, error) {
	return nil, nil
}
