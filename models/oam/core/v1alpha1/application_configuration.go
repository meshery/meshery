package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Configuration is the structure for OAM Application Configuration
type Configuration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ConfigurationSpec `json:"spec,omitempty"`
}

// ConfigurationSpec is the structure for the OAM Application
// Configuration Spec
type ConfigurationSpec struct {
	Components []ConfigurationSpecComponent
}

// ConfigurationSpecComponent is the struct for OAM Appl
type ConfigurationSpecComponent struct {
	ComponentName string
	Traits        []ConfigurationSpecComponentTrait
	Scopes        []ConfigurationSpecComponentScope
}

// ConfigurationSpecComponentTrait is the struct
type ConfigurationSpecComponentTrait struct {
	Name       string
	Properties map[string]interface{}
}

type ConfigurationSpecComponentScope struct {
	ScopeRef ConfigurationSpecComponentScopeRef
}

type ConfigurationSpecComponentScopeRef struct {
	metav1.TypeMeta `json:",inline"`
	Name            string
}
