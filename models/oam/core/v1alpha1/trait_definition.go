package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TraitDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TraitDefinitionSpec `json:"spec,omitempty"`
}

type TraitDefinitionSpec struct {
	AppliesToWorkloads []string      `json:"appliesToWorkloads,omitempty"`
	DefinitionRef      DefinitionRef `json:"definitionRef,omitempty"`
	RevisionEnabled    bool          `json:"revisionEnabled,omitempty"`
}
