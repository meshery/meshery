package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/layer5io/meshery/graphql/graph/generated"
	"github.com/layer5io/meshery/graphql/graph/model"
)

func (r *queryResolver) Cluster(ctx context.Context, id *string) (*model.Cluster, error) {
	// id is for multi clusters - no use right now

	// nodes first
	nodes := nodes

	for _, node := range nodes {
		nodeid := node.ID
		var retPods []*model.Pod

		for _, pod := range pods {
			if pod.Nodeid == nodeid {
				p := &model.Pod{
					Name:     pod.Name,
					ID:       pod.ID,
					Parentid: pod.Nodeid,
				}

				retPods = append(retPods, p)
			}
		}
		node.Pods = retPods
	}

	// now namespaces
	namespaces := namespaces
	for _, namespace := range namespaces {
		namespaceid := namespace.ID
		// input deployments
		for _, deployment := range deployments {
			if deployment.Parentid == namespaceid {
				// adding pods to the deployment
				var retPods []*model.Pod

				for _, pod := range pods {
					if pod.Deploymentid == deployment.ID {
						p := &model.Pod{
							Name:     pod.Name,
							ID:       pod.ID,
							Parentid: deployment.ID,
						}

						retPods = append(retPods, p)
					}
				}
				deployment.Pods = retPods

				namespace.Deployments = append(namespace.Deployments, deployment)
			}
		}
	}

	return &model.Cluster{
		Node:      nodes,
		Namespace: namespaces,
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
