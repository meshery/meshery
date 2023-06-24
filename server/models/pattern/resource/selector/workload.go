package selector

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

func (s *Selector) Workload(name string, version string, model string, apiVersion string) (v1alpha1.ComponentDefinition, error) {
	var comp v1alpha1.ComponentDefinition
	name = strings.Split(name, ".")[0]
	fmt.Println(name, model, version)
	if model == "" && name == "Application" { //If model is not passed, default to core
		model = "core"
	}
	if apiVersion == "core.oam.dev/v1alpha1" { //For backwards compatibility with older designs which were created using OAM
		apiVersion = ""
	}
	entities, _ := s.registry.GetEntities(&v1alpha1.ComponentFilter{
		Name:       name,
		ModelName:  model,
		APIVersion: apiVersion,
	})
	found := false
	if len(entities) != 0 {
		for _, en := range entities {
			if en != nil {
				var ok bool
				comp, ok = en.(v1alpha1.ComponentDefinition)
				if ok {
					found = true
				}
				if comp.Model.Version == version { //prefer to use the correct version, if available
					break
				}
			}
		}
	}
	if !found {
		return comp, fmt.Errorf(fmt.Sprintf("could not find component with name: %s, model: %s, apiVersion: %s", name, model, apiVersion))
	}
	return comp, nil
}

func getResourceType(metadata map[string]string) string {
	typ, ok := metadata["@type"]
	if !ok {
		// Legacy resource => For now mark it as core
		return CoreResource
	}

	return typ
}
