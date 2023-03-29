package models

import (
	"context"

	"github.com/layer5io/meshery/server/meshes"
)

// Available Meshery adapters
type MesheryAdapterName string

const (
	Istio    MesheryAdapterName = "meshery-istio"
	Linkerd  MesheryAdapterName = "meshery-linkerd"
	Consul   MesheryAdapterName = "meshery-consul"
	Kuma     MesheryAdapterName = "meshery-kuma"
	NSM      MesheryAdapterName = "meshery-nsm"
	Citrix   MesheryAdapterName = "meshery-cpx"
	Nginx    MesheryAdapterName = "meshery-nginx-sm"
	Octarine MesheryAdapterName = "meshery-octarine"
	AWS      MesheryAdapterName = "meshery-app-mesh"
	OSM      MesheryAdapterName = "meshery-osm"
	Cilium   MesheryAdapterName = "meshery-cilium"
	Traefik  MesheryAdapterName = "meshery-traefik-mesh"
)

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
