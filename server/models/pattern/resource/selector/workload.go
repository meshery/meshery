package selector

import (
	"fmt"
	"strings"

	"github.com/meshery/schemas/models/v1beta1/component"

	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
)

func (s *Selector) GetDefinition(name string, version string, modelName string, apiVersion string) (component.ComponentDefinition, error) {
	var comp *component.ComponentDefinition
	name = strings.Split(name, ".")[0]
	fmt.Println(name, modelName, version, apiVersion)
	if modelName == "" {
		return component.ComponentDefinition{}, fmt.Errorf("model name is required")
	}

	if apiVersion == "" {
		return component.ComponentDefinition{}, fmt.Errorf("apiVersion is required")
	}

	entities, _, _, _ := s.registry.GetEntities(&regv1beta1.ComponentFilter{
		Name:       name,
		ModelName:  modelName,
		APIVersion: apiVersion,
	})

	found := false
	for _, en := range entities {
		if en != nil {
			var ok bool
			comp, ok = en.(*component.ComponentDefinition)
			if ok {
				found = true
			}
			if comp.Model.Version == version { //prefer to use the correct version, if available
				break
			}
		}
	}
	if !found || comp == nil {
		component := component.ComponentDefinition{}
		return component, fmt.Errorf("could not find component with name: %s, model: %s, apiVersion: %s", name, modelName, apiVersion)
	}
	return *comp, nil
}
