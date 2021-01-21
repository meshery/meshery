package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Component is the structure for the core OAM Application Component
type Component struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ComponentSpec `json:"spec,omitempty"`
}

// ComponentSpec is the structure for the core OAM Application Component Spec
type ComponentSpec struct {
	Type       string                 `json:"type,omitempty"`
	Settings   map[string]interface{} `json:"settings,omitempty"`
	Parameters []ComponentParameter   `json:"parameters,omitempty"`
}

// ComponentParameter is the structure for the core OAM Application Component
// Paramater
type ComponentParameter struct {
	Name        string   `json:"name"`
	FieldPaths  []string `json:"fieldPaths"`
	Required    *bool    `json:"required,omitempty"`
	Description *string  `json:"description,omitempty"`
}

// Type returns the type of the component
func (comp *Component) Type() string {
	return comp.Spec.Type
}
