// Package utils contains helpers that bridge the v1beta2/component and
// v1beta3/component representations of a component definition.
//
// Why this bridge exists
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
// The two struct types are nominally distinct in Go, so a pattern
// iteration (v1beta2) that wants to call a meshkit helper (v1beta3) must
// convert. The shared inner types (*v1beta1/model.ModelDefinition,
// *core.ComponentStyles, capabilities slice, configuration map, Metadata
// AdditionalProperties map) are pointer-identical across both
// representations, so the conversion is a shallow field copy rather than
// a deep clone. Mutations through the aliased pointer targets on the
// returned value therefore remain visible on the source.
package utils

import (
	componentv1beta2 "github.com/meshery/schemas/models/v1beta2/component"
	componentv1beta3 "github.com/meshery/schemas/models/v1beta3/component"
)

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
		Metadata: componentv1beta3.ComponentDefinition_Metadata{
			Genealogy:             src.Metadata.Genealogy,
			IsAnnotation:          src.Metadata.IsAnnotation,
			IsNamespaced:          src.Metadata.IsNamespaced,
			Published:             src.Metadata.Published,
			InstanceDetails:       src.Metadata.InstanceDetails,
			ConfigurationUISchema: src.Metadata.ConfigurationUISchema,
			AdditionalProperties:  src.Metadata.AdditionalProperties,
		},
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
