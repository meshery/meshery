package oam

import (
	"fmt"

	"github.com/layer5io/meshery/helpers/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pattern is the golang representation of the Pattern
// config file model
type Pattern struct {
	Name     string             `yaml:"name,omitempty"`
	Services map[string]Service `yaml:"services,omitempty"`
}

// Service represents the services defined within the appfile
type Service struct {
	Type       string                 `yaml:"type,omitempty"`
	Variant    string                 `yaml:"variant,omitempty"`
	Namespace  string                 `yaml:"namespace,omitempty"`
	Enable     bool                   `yaml:"enable,omitempty"`
	DependsOn  []string               `yaml:"dependsOn,omitempty"`
	Properties map[string]interface{} `yaml:"properties,omitempty"`
}

// NewPatternFile takes in raw yaml and encodes it into a construct
func NewPatternFile(yml []byte) (af Pattern, err error) {
	err = yaml.Unmarshal(yml, &af)
	return
}

// ToApplicationComponents generates OAM Application Components from the
// the given AppFile
func (af *Pattern) ToApplicationComponents() ([]v1alpha1.Component, error) {
	var comps []v1alpha1.Component

	for k, v := range af.Services {
		comp := v1alpha1.Component{
			TypeMeta:   v1.TypeMeta{Kind: "Component", APIVersion: "core.oam.dev/v1alpha2"},
			ObjectMeta: v1.ObjectMeta{Name: k, Namespace: v.Namespace},
			Spec: v1alpha1.ComponentSpec{
				Type: v.Type,
				Settings: map[string]interface{}{
					"variant":   v.Variant,
					"enable":    v.Enable,
					"dependsOn": v.DependsOn,
				},
			},
		}

		comps = append(comps, comp)

		fmt.Printf("Component:\n%+v\n\n", comp)
	}

	return comps, nil
}

// ToApplicationConfiguration generates OAM Application Configuration from the
// the given AppFile for a particular application
func (af *Pattern) ToApplicationConfiguration() (v1alpha1.Configuration, error) {
	config := v1alpha1.Configuration{
		TypeMeta:   v1.TypeMeta{Kind: "ApplicationConfiguration", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{Name: af.Name},
	}

	// Create configs for each component
	for _, v := range af.Services {
		// Indicates that map for properies is not empty
		if len(v.Properties) > 0 {
			specComp := v1alpha1.ConfigurationSpecComponent{
				ComponentName: v.Variant,
			}

			for k2, v2 := range v.Properties {
				castToMap, ok := v2.(map[interface{}]interface{})
				trait := v1alpha1.ConfigurationSpecComponentTrait{
					Name: k2,
				}
				if ok {
					cMap, _ := CastMapInterfaceInterfaceToMapStringInterface(castToMap)
					trait.Properties = cMap
				}

				specComp.Traits = append(specComp.Traits, trait)
			}

			config.Spec.Components = append(config.Spec.Components, specComp)
		}
	}

	return config, nil
}

// CastMapInterfaceInterfaceToMapStringInterface tries to convert map[interface{}]interface{} => map[string]interface{}
func CastMapInterfaceInterfaceToMapStringInterface(in map[interface{}]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	for k, v := range in {
		cStr, ok := k.(string)
		if ok {
			out[cStr] = v
			continue
		}

		return out, fmt.Errorf("can't be converted to map[string]interface{}")
	}

	return out, nil
}
