package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/layer5io/meshery/graphql/graph/generated"
	"github.com/layer5io/meshery/graphql/graph/model"
)

func (r *clusterResolver) ClusterNodes(ctx context.Context, obj *model.Cluster, nodeid *string) ([]*model.ClusterNode, error) {
	if nodeid != nil && *nodeid != "" {
		for _, node := range nodes {
			if node.ID == *nodeid {
				return []*model.ClusterNode{node}, nil
			}
		}
	} else {
		return nodes, nil
	}

	return []*model.ClusterNode{}, nil
}

func (r *clusterResolver) Namespaces(ctx context.Context, obj *model.Cluster, namespaceid *string) ([]*model.Namespace, error) {
	if namespaceid != nil && *namespaceid != "" {
		for _, namespace := range namespaces {
			if namespace.ID == *namespaceid {
				return []*model.Namespace{namespace}, nil
			}
		}
	} else {
		return namespaces, nil
	}

	return []*model.Namespace{}, nil
}

func (r *clusterNodeResolver) Pods(ctx context.Context, obj *model.ClusterNode) ([]*model.Pod, error) {
	var resPods []*model.Pod
	for _, pod := range pods {
		if pod.Nodeid == obj.ID {
			resPods = append(resPods, &model.Pod{
				Name:     pod.Name,
				ID:       pod.ID,
				Parentid: obj.ID,
			})
		}
	}
	return resPods, nil
}

func (r *deploymentResolver) Pods(ctx context.Context, obj *model.Deployment) ([]*model.Pod, error) {
	var resPods []*model.Pod
	for _, pod := range pods {
		if pod.Deploymentid == obj.ID {
			resPods = append(resPods, &model.Pod{
				Name:     pod.Name,
				ID:       pod.ID,
				Parentid: obj.ID,
			})
		}
	}
	return resPods, nil
}

func (r *namespaceResolver) Deployments(ctx context.Context, obj *model.Namespace, deploymentid *string) ([]*model.Deployment, error) {
	namespaceid := obj.ID
	var resDeployments []*model.Deployment
	for _, deployment := range deployments {
		if deployment.Parentid == namespaceid {
			if deploymentid == nil || *deploymentid == "" {
				resDeployments = append(resDeployments, deployment)
			} else {
				if deployment.ID == *deploymentid {
					return []*model.Deployment{deployment}, nil
				}
			}
		}
	}
	return resDeployments, nil
}

func (r *namespaceResolver) Services(ctx context.Context, obj *model.Namespace, serviceid *string) ([]*model.Service, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Cluster(ctx context.Context, clusterid *string) ([]*model.Cluster, error) {
	// id is for multi clusters - no use right now

	// collectedFields := graphql.CollectFieldsCtx(ctx, nil)
	// for _, f := range collectedFields {
	// 	fmt.Println(f.Name)
	// 	for _, arg := range f.Arguments {
	// 		fmt.Println(arg.Name, arg.Value)
	// 	}
	// }

	// rezCtx := graphql.GetResolverContext(ctx)
	// fmt.Println(rezCtx.Args)
	// fmt.Println(rezCtx.Path())
	// for k, v := range rezCtx.Parent.Parent.Parent.Args {
	// 	fmt.Println(k, "=", v)
	// }

	// nodes first
	// nodes := nodes

	// for _, node := range nodes {
	// 	nodeid := node.ID
	// 	var retPods []*model.Pod

	// 	for _, pod := range pods {
	// 		if pod.Nodeid == nodeid {
	// 			p := &model.Pod{
	// 				Name:     pod.Name,
	// 				ID:       pod.ID,
	// 				Parentid: pod.Nodeid,
	// 			}

	// 			retPods = append(retPods, p)
	// 		}
	// 	}
	// 	node.Pods = retPods
	// }

	// now namespaces
	// namespaces := namespaces
	// for _, namespace := range namespaces {
	// 	namespaceid := namespace.ID
	// 	// input deployments
	// 	for _, deployment := range deployments {
	// 		if deployment.Parentid == namespaceid {
	// 			// adding pods to the deployment
	// 			var retPods []*model.Pod

	// 			for _, pod := range pods {
	// 				if pod.Deploymentid == deployment.ID {
	// 					p := &model.Pod{
	// 						Name:     pod.Name,
	// 						ID:       pod.ID,
	// 						Parentid: deployment.ID,
	// 					}

	// 					retPods = append(retPods, p)
	// 				}
	// 			}
	// 			deployment.Pods = retPods

	// 			namespace.Deployments = append(namespace.Deployments, deployment)
	// 		}
	// 	}
	// }

	return []*model.Cluster{&model.Cluster{
		Name: "cluster-1",
		ID:   "cluster-1",
	}}, nil
}

// Cluster returns generated.ClusterResolver implementation.
func (r *Resolver) Cluster() generated.ClusterResolver { return &clusterResolver{r} }

// ClusterNode returns generated.ClusterNodeResolver implementation.
func (r *Resolver) ClusterNode() generated.ClusterNodeResolver { return &clusterNodeResolver{r} }

// Deployment returns generated.DeploymentResolver implementation.
func (r *Resolver) Deployment() generated.DeploymentResolver { return &deploymentResolver{r} }

// Namespace returns generated.NamespaceResolver implementation.
func (r *Resolver) Namespace() generated.NamespaceResolver { return &namespaceResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type clusterResolver struct{ *Resolver }
type clusterNodeResolver struct{ *Resolver }
type deploymentResolver struct{ *Resolver }
type namespaceResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
