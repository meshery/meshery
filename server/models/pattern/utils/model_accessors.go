package utils

import (
	componentv1beta2 "github.com/meshery/schemas/models/v1beta2/component"
)

// ModelName returns the name of the model a component belongs to.
//
// Model is only set by hydration against the registry, which skips rather than
// fails on a component it cannot find, so it is nil on every component of a
// design that names an unregistered one. ModelReference carries the same name
// from the design file itself, so it is the fallback.
func ModelName(comp *componentv1beta2.ComponentDefinition) string {
	if comp == nil {
		return ""
	}
	if comp.Model != nil {
		return comp.Model.Name
	}
	return comp.ModelReference.Name
}

// ModelVersion returns the version of the model a component belongs to, with
// the same fallback as ModelName.
func ModelVersion(comp *componentv1beta2.ComponentDefinition) string {
	if comp == nil {
		return ""
	}
	if comp.Model != nil {
		return comp.Model.Model.Version
	}
	return comp.ModelReference.Model.Version
}
