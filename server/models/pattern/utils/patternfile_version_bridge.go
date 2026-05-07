// This file bridges the v1beta1/pattern.PatternFile and
// v1beta3/design.PatternFile representations with typed shallow field
// copies. It exists because the evaluation-engine (rego policies,
// EvaluationRequest / EvaluationResponse) is a documented carve-out that
// still consumes v1beta1/pattern, while meshkit's HydratePattern /
// DehydratePattern helpers operate exclusively on v1beta3/design.
//
// The nested component representation differs between the two versions
// (v1beta1/component vs v1beta2/component), while relationships stay on
// the same v1beta2 type in both versions. The conversion therefore
// reuses the component shallow-copy bridges and preserves aliasing for
// shared maps, slices, and pointer targets instead of silently dropping
// non-JSON fields such as AdditionalProperties.
package utils

import (
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	patternv1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	componentv1beta2 "github.com/meshery/schemas/models/v1beta2/component"
	designv1beta3 "github.com/meshery/schemas/models/v1beta3/design"
)

// PatternV1beta3ToV1beta1 converts a v1beta3/design.PatternFile into the
// v1beta1/pattern.PatternFile shape that the evaluation engine still
// consumes. The conversion is a shallow typed copy, so the error return
// is preserved only for caller stability.
func PatternV1beta3ToV1beta1(src *designv1beta3.PatternFile) (*patternv1beta1.PatternFile, error) {
	if src == nil {
		return nil, nil
	}
	dst := &patternv1beta1.PatternFile{
		ID:            src.ID,
		Name:          src.Name,
		SchemaVersion: src.SchemaVersion,
		Version:       src.Version,
		Metadata:      patternMetadataV1beta3ToV1beta1(src.Metadata),
		Preferences:   designPreferencesV1beta3ToV1beta1(src.Preferences),
		Relationships: src.Relationships,
	}
	if src.Components != nil {
		dst.Components = make([]*componentv1beta1.ComponentDefinition, len(src.Components))
		for i, component := range src.Components {
			dst.Components[i] = ComponentV1beta2ToV1beta1(component)
		}
	}
	return dst, nil
}

// PatternV1beta1ToV1beta3 converts a v1beta1/pattern.PatternFile back to
// the canonical-casing v1beta3/design.PatternFile so meshkit helpers
// (Hydrate/Dehydrate) can operate on it. Inverse of PatternV1beta3ToV1beta1.
func PatternV1beta1ToV1beta3(src *patternv1beta1.PatternFile) (*designv1beta3.PatternFile, error) {
	if src == nil {
		return nil, nil
	}
	dst := &designv1beta3.PatternFile{
		ID:            src.ID,
		Name:          src.Name,
		SchemaVersion: src.SchemaVersion,
		Version:       src.Version,
		Metadata:      patternMetadataV1beta1ToV1beta3(src.Metadata),
		Preferences:   designPreferencesV1beta1ToV1beta3(src.Preferences),
		Relationships: src.Relationships,
	}
	if src.Components != nil {
		dst.Components = make([]*componentv1beta2.ComponentDefinition, len(src.Components))
		for i, component := range src.Components {
			dst.Components[i] = ComponentV1beta1ToV1beta2(component)
		}
	}
	return dst, nil
}

func patternMetadataV1beta1ToV1beta3(src *patternv1beta1.PatternFile_Metadata) *designv1beta3.PatternFile_Metadata {
	if src == nil {
		return nil
	}
	return &designv1beta3.PatternFile_Metadata{
		ResolvedAliases:      src.ResolvedAliases,
		AdditionalProperties: src.AdditionalProperties,
	}
}

func patternMetadataV1beta3ToV1beta1(src *designv1beta3.PatternFile_Metadata) *patternv1beta1.PatternFile_Metadata {
	if src == nil {
		return nil
	}
	return &patternv1beta1.PatternFile_Metadata{
		ResolvedAliases:      src.ResolvedAliases,
		AdditionalProperties: src.AdditionalProperties,
	}
}

func designPreferencesV1beta1ToV1beta3(src *patternv1beta1.DesignPreferences) *designv1beta3.DesignPreferences {
	if src == nil {
		return nil
	}
	return &designv1beta3.DesignPreferences{
		Layers: src.Layers,
	}
}

func designPreferencesV1beta3ToV1beta1(src *designv1beta3.DesignPreferences) *patternv1beta1.DesignPreferences {
	if src == nil {
		return nil
	}
	return &patternv1beta1.DesignPreferences{
		Layers: src.Layers,
	}
}
