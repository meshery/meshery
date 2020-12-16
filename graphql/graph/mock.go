package graph

import "github.com/layer5io/meshery/graphql/graph/model"

type pod struct {
	Name         string
	ID           string
	Nodeid       string
	Deploymentid string
}

var (

	istio = []*model.Istio{
		&model.Istio{
			Name: "istio-1.7",
			Type: "Service Mesh",
			Version: "1.7",
		},
		&model.Istio{
			Name: "istio-1.5",
			Type: "Service Mesh",
			Version: "1.5",
		},
	}

	nodes = []*model.Node{
		&model.Node{
			Name: "kube-master-7f4a",
			ID:   "kube-master-7f4a",
			Istio: istio[0],
		},
		&model.Node{
			Name: "kube-node-45cd",
			ID:   "kube-node-45cd",
			Istio: istio[1],
		},
	}

	namespaces = []*model.Namespace{
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

	deployments = []*model.Deployment{
		&model.Deployment{
			Name:     "coredns",
			ID:       "coredns",
			Parentid: "kube-system",
		},
		&model.Deployment{
			Name:     "ratings-v1",
			ID:       "ratings-v1",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "reviews-v2",
			ID:       "reviews-v2",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "details-v1",
			ID:       "details-v1",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "reviews-v1",
			ID:       "reviews-v1",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "reviews-v3",
			ID:       "reviews-v3",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "istio-ingressgateway",
			ID:       "istio-ingressgateway",
			Parentid: "istio-system",
		},
		&model.Deployment{
			Name:     "istio-egressgateway",
			ID:       "istio-egressgateway",
			Parentid: "istio-system",
		},
		&model.Deployment{
			Name:     "istiod",
			ID:       "istiod",
			Parentid: "istio-system",
		},
		&model.Deployment{
			Name:     "productpage-v1",
			ID:       "productpage-v1",
			Parentid: "default",
		},
		&model.Deployment{
			Name:     "local-path-provisioner",
			ID:       "local-path-provisioner",
			Parentid: "kube-system",
		},
	}

	pods = []*pod{
		&pod{
			Name:         "coredns-8655855d6-f4bj2",
			ID:           "coredns-8655855d6-f4bj2",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "coredns",
		},
		&pod{
			Name:         "svclb-istio-ingressgateway-dk2st",
			ID:           "svclb-istio-ingressgateway-dk2st",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "istio-ingressgateway",
		},
		&pod{
			Name:         "ratings-v1-7dc98c7588-bk9ts",
			ID:           "ratings-v1-7dc98c7588-bk9ts2",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "ratings-v1",
		},
		&pod{
			Name:         "reviews-v2-7d79d5bd5d-74tsv",
			ID:           "reviews-v2-7d79d5bd5d-74tsv",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "reviews-v2",
		},
		&pod{
			Name:         "svclb-istio-ingressgateway-7sf7h",
			ID:           "svclb-istio-ingressgateway-7sf7h",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "istio-ingressgateway",
		},
		&pod{
			Name:         "details-v1-558b8b4b76-zf8wm",
			ID:           "details-v1-558b8b4b76-zf8wm",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "details-v1",
		},
		&pod{
			Name:         "reviews-v1-7f99cc4496-djssb",
			ID:           "reviews-v1-7f99cc4496-djssb",
			Nodeid:       "kube-master-7f4a",
			Deploymentid: "reviews-v1",
		},
		&pod{
			Name:         "reviews-v3-7dbcdcbc56-94bkl",
			ID:           "reviews-v3-7dbcdcbc56-94bkl",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "reviews-v3",
		},
		&pod{
			Name:         "istio-ingressgateway-669cfc876b-rqcsv",
			ID:           "istio-ingressgateway-669cfc876b-rqcsv",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istio-ingressgateway",
		},
		&pod{
			Name:         "istio-egressgateway-85c567665-bkgbt",
			ID:           "istio-egressgateway-85c567665-bkgbt",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istio-egressgateway",
		},
		&pod{
			Name:         "istiod-5b9488f8bc-fs27g",
			ID:           "istiod-5b9488f8bc-fs27g",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "istiod",
		},
		&pod{
			Name:         "productpage-v1-6987489c74-jxgzl",
			ID:           "productpage-v1-6987489c74-jxgzl",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "productpage-v1",
		},
		&pod{
			Name:         "local-path-provisioner-6d59f47c7-fwq8w",
			ID:           "local-path-provisioner-6d59f47c7-fwq8w",
			Nodeid:       "kube-node-45cd",
			Deploymentid: "local-path-provisioner",
		},
	}
)
