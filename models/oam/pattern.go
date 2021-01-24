package oam

import (
	"fmt"

	"github.com/layer5io/meshery/models/oam/core/v1alpha1"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Pattern is the golang representation of the Pattern
// config file model
type Pattern struct {
	Name     string              `yaml:"name,omitempty"`
	Services map[string]*Service `yaml:"services,omitempty"`
}

// Service represents the services defined within the appfile
type Service struct {
	Type      string   `yaml:"type,omitempty"`
	Namespace string   `yaml:"namespace,omitempty"`
	DependsOn []string `yaml:"dependsOn,omitempty"`

	Settings map[string]interface{} `yaml:"settings,omitempty"`
	Traits   map[string]interface{} `yaml:"traits,omitempty"`
}

// NewPatternFile takes in raw yaml and encodes it into a construct
func NewPatternFile(yml []byte) (af Pattern, err error) {
	err = yaml.Unmarshal(yml, &af)

	for _, svc := range af.Services {
		svc.Settings = RecursiveCastMapStringInterfaceToMapStringInterface(svc.Settings)
		svc.Traits = RecursiveCastMapStringInterfaceToMapStringInterface(svc.Traits)
	}
	return
}

// GetApplicationComponent generates OAM Application Components from the
// the given Pattern file
func (af *Pattern) GetApplicationComponent(name string) (v1alpha1.Component, error) {
	svc, ok := af.Services[name]
	if !ok {
		return v1alpha1.Component{}, fmt.Errorf("invalid service name")
	}

	comp := v1alpha1.Component{
		TypeMeta:   v1.TypeMeta{Kind: "Component", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{Name: name, Namespace: svc.Namespace},
		Spec: v1alpha1.ComponentSpec{
			Type:     svc.Type,
			Settings: svc.Settings,
		},
	}

	return comp, nil
}

// GenerateApplicationConfiguration generates OAM Application Configuration from the
// the given Pattern file for a particular deploymnet
func (af *Pattern) GenerateApplicationConfiguration() (v1alpha1.Configuration, error) {
	config := v1alpha1.Configuration{
		TypeMeta:   v1.TypeMeta{Kind: "ApplicationConfiguration", APIVersion: "core.oam.dev/v1alpha2"},
		ObjectMeta: v1.ObjectMeta{Name: af.Name},
	}

	// Create configs for each component
	for _, v := range af.Services {
		// Indicates that map for properies is not empty
		if len(v.Traits) > 0 {
			specComp := v1alpha1.ConfigurationSpecComponent{
				ComponentName: v.Type,
			}

			for k2, v2 := range v.Traits {
				castToMap, ok := v2.(map[string]interface{})

				trait := v1alpha1.ConfigurationSpecComponentTrait{
					Name: k2,
				}

				if !ok {
					castToMap = map[string]interface{}{}
				}

				trait.Properties = castToMap

				specComp.Traits = append(specComp.Traits, trait)
			}

			config.Spec.Components = append(config.Spec.Components, specComp)
		}
	}

	return config, nil
}

// RecursiveCastMapStringInterfaceToMapStringInterface will convert a
// map[string]interface{} recursively => map[string]interface{}
func RecursiveCastMapStringInterfaceToMapStringInterface(in map[string]interface{}) map[string]interface{} {
	interfaceMap := make(map[interface{}]interface{})

	for k, v := range in {
		interfaceMap[k] = v
	}

	return CastMapInterfaceInterfaceToMapStringInterface(interfaceMap)
}

// CastMapInterfaceInterfaceToMapStringInterface tries to convert map[interface{}]interface{} => map[string]interface{}
func CastMapInterfaceInterfaceToMapStringInterface(in map[interface{}]interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	for k, v := range in {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			out[fmt.Sprint(k)] = CastMapInterfaceInterfaceToMapStringInterface(v2)
		default:
			out[fmt.Sprint(k)] = v
		}
	}

	return out
}
