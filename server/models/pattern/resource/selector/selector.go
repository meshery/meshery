package selector

import (
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
)

const (
	CoreResource = "pattern.meshery.io/core"
	MeshResource = "pattern.meshery.io/mesh/workload"
	K8sResource  = "pattern.meshery.io/k8s"
)

type Selector struct {
	registry *meshmodel.RegistryManager
}

func New(reg *meshmodel.RegistryManager) *Selector {
	return &Selector{
		registry: reg,
	}
}
