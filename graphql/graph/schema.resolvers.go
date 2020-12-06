package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/layer5io/meshery/graphql/graph/generated"
	"github.com/layer5io/meshery/graphql/graph/model"
)

func (r *queryResolver) Node(ctx context.Context) ([]*model.Node, error) {
	nodes := []*model.Node{
		&model.Node{
			Name: "kube-master-7f4a",
			ID:   "48F1733E-6163-4FE9-8E43-C35AEA88E816",
		},
		&model.Node{
			Name: "kube-node-45cd",
			ID:   "4C42AE40-2678-4A3F-AED9-A0ABB799CAA4",
		},
	}

	return nodes, nil
}

func (r *queryResolver) Namespace(ctx context.Context) ([]*model.Namespace, error) {
	namespaces := []*model.Namespace{
		&model.Namespace{
			Name: "kube-system",
			ID:   "kube-system",
		},
		&model.Namespace{
			Name: "kube-public",
			ID:   "kube-public",
		},
		&model.Namespace{
			Name: "kube-node-lease",
			ID:   "kube-node-lease",
		},
		&model.Namespace{
			Name: "istio-system",
			ID:   "istio-system",
		},
		&model.Namespace{
			Name: "default",
			ID:   "default",
		},
	}

	return namespaces, nil
}

func (r *queryResolver) Deployment(ctx context.Context, layout *model.DeploymentFilter) ([]*model.Deployment, error) {
	deployments := []*model.Deployment{
		&model.Deployment{
			Name:        "coredns",
			ID:          "coredns",
			Namespaceid: "kube-system",
		},
		&model.Deployment{
			Name:        "ratings-v1",
			ID:          "ratings-v1",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "reviews-v2",
			ID:          "reviews-v2",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "details-v1",
			ID:          "details-v1",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "reviews-v1",
			ID:          "reviews-v1",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "reviews-v3",
			ID:          "reviews-v3",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "istio-ingressgateway",
			ID:          "istio-ingressgateway",
			Namespaceid: "istio-system",
		},
		&model.Deployment{
			Name:        "istio-egressgateway",
			ID:          "istio-egressgateway",
			Namespaceid: "istio-system",
		},
		&model.Deployment{
			Name:        "istiod",
			ID:          "istiod",
			Namespaceid: "istio-system",
		},
		&model.Deployment{
			Name:        "productpage-v1",
			ID:          "productpage-v1",
			Namespaceid: "default",
		},
		&model.Deployment{
			Name:        "local-path-provisioner",
			ID:          "local-path-provisioner",
			Namespaceid: "kube-system",
		},
	}

	if layout == nil || layout.Namespaceid == nil || *layout.Namespaceid == "" {
		return deployments, nil
	}

	var result []*model.Deployment
	for _, deployment := range deployments {
		if deployment.Namespaceid == *layout.Namespaceid {
			result = append(result, deployment)
		}
	}

	return result, nil
}

func (r *queryResolver) Pod(ctx context.Context, layout *model.PodFilter) ([]*model.Pod, error) {
	pods := []*model.Pod{
		&model.Pod{
			Name:         "coredns-8655855d6-f4bj2",
			ID:           "coredns-8655855d6-f4bj2",
			Namespaceid:  "kube-system",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "coredns",
		},
		&model.Pod{
			Name:         "svclb-istio-ingressgateway-dk2st",
			ID:           "svclb-istio-ingressgateway-dk2st",
			Namespaceid:  "istio-system",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "istio-ingressgateway",
		},
		&model.Pod{
			Name:         "ratings-v1-7dc98c7588-bk9ts",
			ID:           "ratings-v1-7dc98c7588-bk9ts2",
			Namespaceid:  "default",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "ratings-v1",
		},
		&model.Pod{
			Name:         "reviews-v2-7d79d5bd5d-74tsv",
			ID:           "reviews-v2-7d79d5bd5d-74tsv",
			Namespaceid:  "default",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "reviews-v2",
		},
		&model.Pod{
			Name:         "svclb-istio-ingressgateway-7sf7h",
			ID:           "svclb-istio-ingressgateway-7sf7h",
			Namespaceid:  "istio-system",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "istio-ingressgateway",
		},
		&model.Pod{
			Name:         "details-v1-558b8b4b76-zf8wm",
			ID:           "details-v1-558b8b4b76-zf8wm",
			Namespaceid:  "default",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "details-v1",
		},
		&model.Pod{
			Name:         "reviews-v1-7f99cc4496-djssb",
			ID:           "reviews-v1-7f99cc4496-djssb",
			Namespaceid:  "default",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "reviews-v1",
		},
		&model.Pod{
			Name:         "reviews-v3-7dbcdcbc56-94bkl",
			ID:           "reviews-v3-7dbcdcbc56-94bkl",
			Namespaceid:  "default",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "reviews-v3",
		},
		&model.Pod{
			Name:         "istio-ingressgateway-669cfc876b-rqcsv",
			ID:           "istio-ingressgateway-669cfc876b-rqcsv",
			Namespaceid:  "istio-system",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istio-ingressgateway",
		},
		&model.Pod{
			Name:         "istio-egressgateway-85c567665-bkgbt",
			ID:           "istio-egressgateway-85c567665-bkgbt",
			Namespaceid:  "istio-system",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istio-egressgateway",
		},
		&model.Pod{
			Name:         "istiod-5b9488f8bc-fs27g",
			ID:           "istiod-5b9488f8bc-fs27g",
			Namespaceid:  "istio-system",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istiod",
		},
		&model.Pod{
			Name:         "productpage-v1-6987489c74-jxgzl",
			ID:           "productpage-v1-6987489c74-jxgzl",
			Namespaceid:  "default",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "productpage-v1",
		},
		&model.Pod{
			Name:         "local-path-provisioner-6d59f47c7-fwq8w",
			ID:           "local-path-provisioner-6d59f47c7-fwq8w",
			Namespaceid:  "kube-system",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "local-path-provisioner",
		},
	}

	if layout == nil {
		return pods, nil
	}

	if layout.Namespaceid != nil && *layout.Namespaceid != "" {
		var temp []*model.Pod
		for _, pod := range pods {
			if pod.Namespaceid == *layout.Namespaceid {
				temp = append(temp, pod)
			}
		}
		pods = temp
	}

	if layout.Nodeid != nil && *layout.Nodeid != "" {
		var temp []*model.Pod
		for _, pod := range pods {
			if pod.Nodeid == *layout.Nodeid {
				temp = append(temp, pod)
			}
		}
		pods = temp
	}

	if layout.Deploymentid != nil && *layout.Deploymentid != "" {
		var temp []*model.Pod
		for _, pod := range pods {
			if pod.Deploymentid == *layout.Deploymentid {
				temp = append(temp, pod)
			}
		}
		pods = temp
	}

	return pods, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
