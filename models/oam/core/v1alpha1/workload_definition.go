package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type WorkloadDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec WorkloadDefinitionSpec `json:"spec,omitempty"`
}

type WorkloadDefinitionSpec struct {
	DefinitionRef DefinitionRef `json:"definitionRef,omitempty"`
}
