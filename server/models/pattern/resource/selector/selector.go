package selector

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

const (
	CoreResource = "pattern.meshery.io/core"
	MeshResource = "pattern.meshery.io/mesh/workload"
	K8sResource  = "pattern.meshery.io/k8s"
)

type Helpers interface {
	GetServiceMesh() (name string, version string)
	GetAPIVersionForKind(kind string) string
}

type Selector struct {
	registry *meshmodel.RegistryManager
	helpers  Helpers
}

func New(reg *meshmodel.RegistryManager, helpers Helpers) *Selector {
	return &Selector{
		registry: reg,
		helpers:  helpers,
	}
}

const annotationsPrefix = "design.meshmodel.io"

func GetAnnotationsForWorkload(w v1alpha1.ComponentDefinition) map[string]string {
	res := map[string]string{}

	for key, val := range w.Metadata {
		if v, ok := val.(string); ok {
			res[strings.ReplaceAll(fmt.Sprintf("%s.%s", annotationsPrefix, key), " ", "")] = v
		}
	}
	res[fmt.Sprintf("%s.model.name", annotationsPrefix)] = w.Model.Name
	res[fmt.Sprintf("%s.model.version", annotationsPrefix)] = w.Model.Version
	res[fmt.Sprintf("%s.model.category", annotationsPrefix)] = w.Model.Category
	return res
}

func generateWorkloadKey(name string) string {
	return fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		"core.oam.dev/v1alpha1",
		"WorkloadDefinition",
		name,
	)
}

func generateTraitKey(name string) string {
	return fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		"core.oam.dev/v1alpha1",
		"TraitDefinition",
		name,
	)
}
