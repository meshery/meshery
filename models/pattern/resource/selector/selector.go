package selector

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/models/pattern/core"
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
	helpers Helpers
}

func New(helpers Helpers) *Selector {
	return &Selector{
		helpers: helpers,
	}
}

func GetAnnotationsForWorkload(w core.WorkloadCapability) map[string]string {
	res := map[string]string{}

	metadata := w.OAMDefinition.Spec.Metadata
	typ, ok := metadata["@type"]
	if ok {
		for k, v := range metadata {
			if k == "@type" {
				continue
			}

			res[fmt.Sprintf("%s.%s", strings.ReplaceAll(typ, "/", "."), k)] = v
		}
	}

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
