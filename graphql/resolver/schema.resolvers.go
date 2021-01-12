package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/layer5io/meshery/graphql/generated"
	"github.com/layer5io/meshery/graphql/model"
)

func (r *queryResolver) MapEdges(ctx context.Context, id *model.MapEdgeIdentifier) ([]*model.MapEdge, error) {
	return r.Istio.MapEdges(id.View)
}

func (r *queryResolver) MapNodes(ctx context.Context, id *model.MapNodeIdentifier) ([]*model.MapNode, error) {
	return r.Istio.MapNodes(id.View)
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
