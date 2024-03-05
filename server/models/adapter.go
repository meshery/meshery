package models

import (
	"context"

	"github.com/layer5io/meshery/server/meshes"
)

// Available Meshery adapters
var (
	Istio     = Adapter{Name: "meshery-istio", Location: "10000"}
	Linkerd   = Adapter{Name: "meshery-linkerd", Location: "10001"}
	Consul    = Adapter{Name: "meshery-consul", Location: "10002"}
	NSM       = Adapter{Name: "meshery-nsm", Location: "10004"}
	AppMesh   = Adapter{Name: "meshery-app-mesh", Location: "10005"}
	Traefik   = Adapter{Name: "meshery-traefik-mesh", Location: "10006"}
	Kuma      = Adapter{Name: "meshery-kuma", Location: "10007"}
	Nginx     = Adapter{Name: "meshery-nginx-sm", Location: "10010"}
	Cilium    = Adapter{Name: "meshery-cilium", Location: "10012"}
	Nighthawk = Adapter{Name: "meshery-nighthawk", Location: "10013"}
)

var ListAvailableAdapters = []Adapter{Istio, Linkerd, Consul, Kuma, NSM, Nginx, AppMesh, Cilium, Traefik, Nighthawk}

// Adapter represents an adapter in Meshery
type Adapter struct {
	Location     string                       `json:"adapter_location"`
	Name         string                       `json:"name"`
	Version      string                       `json:"version"`
	GitCommitSHA string                       `json:"git_commit_sha"`
	Ops          []*meshes.SupportedOperation `json:"ops"`
}

// AdaptersTrackerInterface defines the methods a type should implement to be an adapter tracker
type AdaptersTrackerInterface interface {
	AddAdapter(context.Context, Adapter)
	RemoveAdapter(context.Context, Adapter)
	GetAdapters(context.Context) []Adapter
	DeployAdapter(context.Context, Adapter) error
	UndeployAdapter(context.Context, Adapter) error
}
