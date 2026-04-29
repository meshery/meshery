package selector

import (
	"fmt"
	"strings"

	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	componentv1beta2 "github.com/meshery/schemas/models/v1beta2/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
)

// GetDefinition looks up a component definition from the registry and
// returns it in the v1beta2 representation that PatternFile.Components
// uses. The registry itself stores v1beta3/component.ComponentDefinition
// (the canonical-casing version that implements entity.Entity); the
// v1beta3 → v1beta2 bridge happens in FindCompDefinitionWithVersion.
func (s *Selector) GetDefinition(name string, version string, modelName string, apiVersion string, allowEmptyAPIVersion bool) (componentv1beta2.ComponentDefinition, error) {
	name = strings.Split(name, ".")[0]
	fmt.Println(name, modelName, version, apiVersion)
	if modelName == "" {
		return componentv1beta2.ComponentDefinition{}, fmt.Errorf("model name is required")
	}

	if apiVersion == "" && !allowEmptyAPIVersion {
		return componentv1beta2.ComponentDefinition{}, fmt.Errorf("apiVersion is required")
	}

	entities, _, _, _ := s.registry.GetEntities(&regv1beta1.ComponentFilter{
		Name:       name,
		ModelName:  modelName,
		APIVersion: apiVersion,
	})

	comp, found := FindCompDefinitionWithVersion(entities, version)
	if !found || comp == nil {
		return componentv1beta2.ComponentDefinition{}, fmt.Errorf("could not find component with name: %s, model: %s, apiVersion: %s", name, modelName, apiVersion)
	}
	return *comp, nil
}

// FindCompDefinitionWithVersion walks a list of registry entities, casts
// the v1beta3 component-definitions they contain, prefers an exact
// Model.Model.Version match, and returns the resulting component
// converted back to the v1beta2 representation so pattern stages can
// consume it directly.
func FindCompDefinitionWithVersion(entities []entity.Entity, version string) (*componentv1beta2.ComponentDefinition, bool) {
	var match *componentv1beta3.ComponentDefinition
	found := false
	for _, en := range entities {
		if en == nil {
			continue
		}
		comp, ok := en.(*componentv1beta3.ComponentDefinition)
		if !ok {
			continue
		}
		found = true
		match = comp
		if comp.Model != nil && comp.Model.Model.Version == version {
			// prefer the exact version match when available
			break
		}
	}
	if !found || match == nil {
		return &componentv1beta2.ComponentDefinition{}, found
	}
	return patternutils.ComponentV1beta3ToV1beta2(match), true
}
