// Package utils contains helpers that bridge the v1beta1/component,
// v1beta2/component, and v1beta3/component representations of a
// component definition. The pairings currently provided are
// v1beta2 <-> v1beta3 (used by the pattern pipeline, which holds
// PatternFile.Components as v1beta2) and v1beta3 -> v1beta1 (used by
// the relationship-policy evaluator carve-out in
// server/handlers/policy_relationship_handler.go, which is pinned to
// v1beta1/pattern.EvaluationRequest/EvaluationResponse).
//
// # Why these bridges exist
//
// meshery/schemas pins v1beta3/design.PatternFile.Components to
// *v1beta2/component.ComponentDefinition (see schemas' v1beta3/const.go
// commentary). At the same time the meshkit registry, orchestration
// helpers (EnrichComponentWithMesheryMetadata) and dependency handlers
// (DependencyHandler.HandleDependents) operate on
// *v1beta3/component.ComponentDefinition, because v1beta3/component is
// the only version whose Entity interface is implemented and which
// AutoMigrate binds to the component_definition_dbs table.
//
// The relationship-policy evaluator stays on v1beta1/pattern by design
// (the EvaluationRequest / EvaluationResponse types and the in-tree
// policy engine all consume v1beta1), so registry lookups inside that
// handler also need to bridge the v1beta3 result back to v1beta1.
//
// The struct types are nominally distinct in Go, so any consumer that
// wants to cross a meshkit boundary must convert. The shared inner
// types (*v1beta1/model.ModelDefinition, *core.ComponentStyles,
// capabilities slice, configuration map, Metadata AdditionalProperties
// map) are pointer-identical across all three representations, so each
// conversion is a shallow field copy rather than a deep clone.
// Mutations through the aliased pointer targets on the returned value
// therefore remain visible on the source.
package utils

import (
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	componentv1beta2 "github.com/meshery/schemas/models/v1beta2/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
)

// ComponentV1beta1ToV1beta2 returns a v1beta2/component.ComponentDefinition
// whose inner reference-typed fields (Model pointer, Capabilities slice,
// Configuration map, Styles pointer, Metadata.AdditionalProperties map,
// ModelId pointer) alias the source. Mutations through those aliased maps /
// pointer targets on the returned value remain visible on the source; scalar
// field updates on the returned value do NOT.
func ComponentV1beta1ToV1beta2(src *componentv1beta1.ComponentDefinition) *componentv1beta2.ComponentDefinition {
	if src == nil {
		return nil
	}
	dst := &componentv1beta2.ComponentDefinition{
		ID:             src.ID,
		SchemaVersion:  src.SchemaVersion,
		Version:        src.Version,
		DisplayName:    src.DisplayName,
		Description:    src.Description,
		Format:         componentv1beta2.ComponentDefinitionFormat(src.Format),
		Model:          src.Model,
		ModelReference: src.ModelReference,
		Styles:         src.Styles,
		Capabilities:   src.Capabilities,
		Metadata: componentv1beta2.ComponentDefinition_Metadata{
			Genealogy:             src.Metadata.Genealogy,
			IsAnnotation:          src.Metadata.IsAnnotation,
			IsNamespaced:          src.Metadata.IsNamespaced,
			Published:             src.Metadata.Published,
			InstanceDetails:       src.Metadata.InstanceDetails,
			ConfigurationUISchema: src.Metadata.ConfigurationUISchema,
			AdditionalProperties:  src.Metadata.AdditionalProperties,
		},
		Configuration: src.Configuration,
		Component: componentv1beta2.Component{
			Kind:    src.Component.Kind,
			Version: src.Component.Version,
			Schema:  src.Component.Schema,
		},
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		ModelId:   src.ModelId,
	}
	if src.Status != nil {
		st := componentv1beta2.ComponentDefinitionStatus(*src.Status)
		dst.Status = &st
	}
	return dst
}

// ComponentV1beta2ToV1beta1 is the inverse bridge: a shallow field copy that
// keeps pointer/reference-typed inner fields (Model, Styles, Capabilities,
// Configuration, Metadata.AdditionalProperties, ModelId) aliased across both
// versions so mutations remain visible across the bridge.
func ComponentV1beta2ToV1beta1(src *componentv1beta2.ComponentDefinition) *componentv1beta1.ComponentDefinition {
	if src == nil {
		return nil
	}
	dst := &componentv1beta1.ComponentDefinition{
		ID:             src.ID,
		SchemaVersion:  src.SchemaVersion,
		Version:        src.Version,
		DisplayName:    src.DisplayName,
		Description:    src.Description,
		Format:         componentv1beta1.ComponentDefinitionFormat(src.Format),
		Model:          src.Model,
		ModelReference: src.ModelReference,
		Styles:         src.Styles,
		Capabilities:   src.Capabilities,
		Metadata: componentv1beta1.ComponentDefinition_Metadata{
			Genealogy:             src.Metadata.Genealogy,
			IsAnnotation:          src.Metadata.IsAnnotation,
			IsNamespaced:          src.Metadata.IsNamespaced,
			Published:             src.Metadata.Published,
			InstanceDetails:       src.Metadata.InstanceDetails,
			ConfigurationUISchema: src.Metadata.ConfigurationUISchema,
			AdditionalProperties:  src.Metadata.AdditionalProperties,
		},
		Configuration: src.Configuration,
		Component: componentv1beta1.Component{
			Kind:    src.Component.Kind,
			Version: src.Component.Version,
			Schema:  src.Component.Schema,
		},
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		ModelId:   src.ModelId,
	}
	if src.Status != nil {
		st := componentv1beta1.ComponentDefinitionStatus(*src.Status)
		dst.Status = &st
	}
	return dst
}

// ComponentV1beta2ToV1beta3 returns a v1beta3/component.ComponentDefinition
// whose inner reference-typed fields (Model pointer, Capabilities slice,
// Configuration map, Styles pointer, Metadata.AdditionalProperties map,
// ModelId pointer) alias the source. Mutations to those aliased maps /
// pointer targets on the returned value are visible on the source;
// scalar field updates on the returned value are NOT. Use this at the
// meshkit-boundary only; prefer working on the native v1beta2 value
// everywhere else.
func ComponentV1beta2ToV1beta3(src *componentv1beta2.ComponentDefinition) *componentv1beta3.ComponentDefinition {
	if src == nil {
		return nil
	}
	dst := &componentv1beta3.ComponentDefinition{
		ID:             src.ID,
		SchemaVersion:  src.SchemaVersion,
		Version:        src.Version,
		DisplayName:    src.DisplayName,
		Description:    src.Description,
		Format:         componentv1beta3.ComponentDefinitionFormat(src.Format),
		Model:          src.Model,
		ModelReference: src.ModelReference,
		Styles:         src.Styles,
		Capabilities:   src.Capabilities,
		// Metadata is structurally identical across v1beta1, v1beta2,
		// and v1beta3 (the generated schemas only differ in docstring
		// content), so direct conversion copies every field — including
		// future additions — instead of risking drift via a
		// hand-rolled per-field copy. ApplyV1beta3MetadataChanges below
		// uses the same pattern.
		Metadata:      componentv1beta3.ComponentDefinition_Metadata(src.Metadata),
		Configuration: src.Configuration,
		Component: componentv1beta3.Component{
			Kind:    src.Component.Kind,
			Version: src.Component.Version,
			Schema:  src.Component.Schema,
		},
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		ModelID:   src.ModelId,
	}
	if src.Status != nil {
		st := componentv1beta3.ComponentDefinitionStatus(*src.Status)
		dst.Status = &st
	}
	return dst
}

// ComponentV1beta3ToV1beta2 is the inverse bridge: a shallow field copy
// that keeps pointer/reference-typed inner fields (Model, Styles,
// Capabilities, Configuration, Metadata.AdditionalProperties, ModelID)
// aliased across both versions so mutations remain visible across the
// bridge. Use it at read-boundaries where a meshkit helper has handed
// back a v1beta3 value but the rest of the pattern pipeline expects
// v1beta2.
func ComponentV1beta3ToV1beta2(src *componentv1beta3.ComponentDefinition) *componentv1beta2.ComponentDefinition {
	if src == nil {
		return nil
	}
	dst := &componentv1beta2.ComponentDefinition{
		ID:             src.ID,
		SchemaVersion:  src.SchemaVersion,
		Version:        src.Version,
		DisplayName:    src.DisplayName,
		Description:    src.Description,
		Format:         componentv1beta2.ComponentDefinitionFormat(src.Format),
		Model:          src.Model,
		ModelReference: src.ModelReference,
		Styles:         src.Styles,
		Capabilities:   src.Capabilities,
		// See ComponentV1beta2ToV1beta3 above for the rationale; same
		// applies in the reverse direction.
		Metadata:      componentv1beta2.ComponentDefinition_Metadata(src.Metadata),
		Configuration: src.Configuration,
		Component: componentv1beta2.Component{
			Kind:    src.Component.Kind,
			Version: src.Component.Version,
			Schema:  src.Component.Schema,
		},
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		ModelId:   src.ModelID,
	}
	if src.Status != nil {
		st := componentv1beta2.ComponentDefinitionStatus(*src.Status)
		dst.Status = &st
	}
	return dst
}

// ApplyV1beta3MetadataChanges copies mutations that meshkit helpers make
// to a v1beta3 component back onto the v1beta2 source. This covers the
// handful of meshkit helpers (notably orchestration.EnrichComponentWithMesheryMetadata)
// that may REPLACE Configuration rather than mutate it. Because
// ComponentV1beta2ToV1beta3 aliases the Configuration map between the
// two values, mutations made via the returned v1beta3 value are already
// visible on the source — this copy-back only matters when the helper
// replaces a map.
func ApplyV1beta3MetadataChanges(
	src *componentv1beta3.ComponentDefinition,
	dst *componentv1beta2.ComponentDefinition,
) {
	if src == nil || dst == nil {
		return
	}
	dst.Configuration = src.Configuration
	// v1beta2/componentv1beta2.ComponentDefinition.Metadata and v1beta3's
	// Metadata are structurally identical (same package-level type across
	// the two generated schemas), so a direct conversion copies every
	// field — AdditionalProperties, Published, IsAnnotation, IsNamespaced,
	// Genealogy, Shape, SvgColor, etc. — rather than leaving new fields
	// out of sync between the registry-hydrated v1beta3 and the pattern-
	// level v1beta2 components.
	dst.Metadata = componentv1beta2.ComponentDefinition_Metadata(src.Metadata)
}

// ComponentV1beta3ToV1beta1 is a shallow field copy from
// v1beta3/component.ComponentDefinition to v1beta1/component.ComponentDefinition,
// keeping the same aliasing semantics as ComponentV1beta3ToV1beta2: inner
// reference-typed fields (Model pointer, Styles, Capabilities slice,
// Configuration map, Metadata.AdditionalProperties map, ModelID pointer)
// alias the source so mutations remain visible across the bridge.
//
// Used by the relationship-policy evaluator carve-out (server/handlers/
// policy_relationship_handler.go), which stays on v1beta1/pattern even
// though the meshkit registry now returns v1beta3-typed components.
func ComponentV1beta3ToV1beta1(src *componentv1beta3.ComponentDefinition) *componentv1beta1.ComponentDefinition {
	if src == nil {
		return nil
	}
	dst := &componentv1beta1.ComponentDefinition{
		ID:             src.ID,
		SchemaVersion:  src.SchemaVersion,
		Version:        src.Version,
		DisplayName:    src.DisplayName,
		Description:    src.Description,
		Format:         componentv1beta1.ComponentDefinitionFormat(src.Format),
		Model:          src.Model,
		ModelReference: src.ModelReference,
		Styles:         src.Styles,
		Capabilities:   src.Capabilities,
		// v1beta1 and v1beta3 ComponentDefinition_Metadata are
		// structurally identical (verified against the generated
		// schemas), so direct conversion copies every field — including
		// any future additions — and avoids the maintenance hazard of a
		// hand-rolled field-by-field copy. Mirrors the pattern used by
		// ApplyV1beta3MetadataChanges below.
		Metadata:      componentv1beta1.ComponentDefinition_Metadata(src.Metadata),
		Configuration: src.Configuration,
		Component: componentv1beta1.Component{
			Kind:    src.Component.Kind,
			Version: src.Component.Version,
			Schema:  src.Component.Schema,
		},
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		ModelId:   src.ModelID,
	}
	if src.Status != nil {
		st := componentv1beta1.ComponentDefinitionStatus(*src.Status)
		dst.Status = &st
	}
	return dst
}
