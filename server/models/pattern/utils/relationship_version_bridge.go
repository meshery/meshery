package utils

import (
	capabilityv1alpha1 "github.com/meshery/schemas/models/v1alpha1/capability"
	relationshipv1alpha3 "github.com/meshery/schemas/models/v1alpha3/relationship"
	capabilityv1beta1 "github.com/meshery/schemas/models/v1beta1/capability"
	relationshipv1beta2 "github.com/meshery/schemas/models/v1beta2/relationship"
)

// RelationshipV1alpha3ToV1beta2 converts the registry's v1alpha3
// RelationshipDefinition into the v1beta2 shape used by the Go policy engine.
// The conversion is a shallow typed copy: shared maps, selector/model pointers,
// and same-typed scalar pointers remain aliased, while fields whose pointed-to
// named types differ across versions are copied explicitly.
func RelationshipV1alpha3ToV1beta2(src *relationshipv1alpha3.RelationshipDefinition) *relationshipv1beta2.RelationshipDefinition {
	if src == nil {
		return nil
	}
	dst := &relationshipv1beta2.RelationshipDefinition{
		ID:                   src.ID,
		Capabilities:         relationshipCapabilitiesV1alpha3ToV1beta2(src.Capabilities),
		EvaluationQuery:      src.EvaluationQuery,
		Kind:                 relationshipv1beta2.RelationshipDefinitionKind(src.Kind),
		RelationshipMetadata: relationshipMetadataV1alpha3ToV1beta2(src.Metadata),
		Model:                src.Model,
		ModelId:              src.ModelId,
		SchemaVersion:        src.SchemaVersion,
		Selectors:            selectorSetV1alpha3ToV1beta2(src.Selectors),
		SubType:              src.SubType,
		Status:               convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionStatus, relationshipv1beta2.RelationshipDefinitionStatus](src.Status),
		RelationshipType:     src.RelationshipType,
		Version:              src.Version,
	}
	return dst
}

// RelationshipV1beta2ToV1alpha3 is the inverse bridge used by tests and future
// callers that need to round-trip a relationship definition back into the
// registry-facing schema version.
func RelationshipV1beta2ToV1alpha3(src *relationshipv1beta2.RelationshipDefinition) *relationshipv1alpha3.RelationshipDefinition {
	if src == nil {
		return nil
	}
	dst := &relationshipv1alpha3.RelationshipDefinition{
		ID:               src.ID,
		Capabilities:     relationshipCapabilitiesV1beta2ToV1alpha3(src.Capabilities),
		EvaluationQuery:  src.EvaluationQuery,
		Kind:             relationshipv1alpha3.RelationshipDefinitionKind(src.Kind),
		Metadata:         relationshipMetadataV1beta2ToV1alpha3(src.RelationshipMetadata),
		Model:            src.Model,
		ModelId:          src.ModelId,
		SchemaVersion:    src.SchemaVersion,
		Selectors:        selectorSetV1beta2ToV1alpha3(src.Selectors),
		SubType:          src.SubType,
		Status:           convertEnumPtr[relationshipv1beta2.RelationshipDefinitionStatus, relationshipv1alpha3.RelationshipDefinitionStatus](src.Status),
		RelationshipType: src.RelationshipType,
		Version:          src.Version,
	}
	return dst
}

func relationshipCapabilitiesV1alpha3ToV1beta2(src *[]capabilityv1alpha1.Capability) *[]capabilityv1beta1.Capability {
	if src == nil {
		return nil
	}
	dst := make([]capabilityv1beta1.Capability, len(*src))
	for i, capability := range *src {
		dst[i] = capabilityV1alpha1ToV1beta1(capability)
	}
	return &dst
}

func relationshipCapabilitiesV1beta2ToV1alpha3(src *[]capabilityv1beta1.Capability) *[]capabilityv1alpha1.Capability {
	if src == nil {
		return nil
	}
	dst := make([]capabilityv1alpha1.Capability, len(*src))
	for i, capability := range *src {
		dst[i] = capabilityV1beta1ToV1alpha1(capability)
	}
	return &dst
}

func capabilityV1alpha1ToV1beta1(src capabilityv1alpha1.Capability) capabilityv1beta1.Capability {
	return capabilityv1beta1.Capability{
		Description:   src.Description,
		DisplayName:   src.DisplayName,
		EntityState:   src.EntityState,
		Key:           src.Key,
		Kind:          src.Kind,
		Metadata:      src.Metadata,
		SchemaVersion: src.SchemaVersion,
		Status:        capabilityv1beta1.CapabilityStatus(src.Status),
		SubType:       src.SubType,
		Type:          src.Type,
		Version:       src.Version,
	}
}

func capabilityV1beta1ToV1alpha1(src capabilityv1beta1.Capability) capabilityv1alpha1.Capability {
	return capabilityv1alpha1.Capability{
		Description:   src.Description,
		DisplayName:   src.DisplayName,
		EntityState:   src.EntityState,
		Key:           src.Key,
		Kind:          src.Kind,
		Metadata:      src.Metadata,
		SchemaVersion: src.SchemaVersion,
		Status:        capabilityv1alpha1.CapabilityStatus(src.Status),
		SubType:       src.SubType,
		Type:          src.Type,
		Version:       src.Version,
	}
}

func relationshipMetadataV1alpha3ToV1beta2(src *relationshipv1alpha3.RelationshipMetadata) *relationshipv1beta2.RelationshipMetadata {
	if src == nil {
		return nil
	}
	// Keep the explicit field copy here instead of a direct struct conversion:
	// the metadata envelope is structurally close across versions, but its
	// Styles pointer targets version-specific structs whose enum-typed pointer
	// fields differ by package and therefore cannot be converted wholesale.
	return &relationshipv1beta2.RelationshipMetadata{
		Description:          src.Description,
		IsAnnotation:         src.IsAnnotation,
		Styles:               relationshipStylesV1alpha3ToV1beta2(src.Styles),
		AdditionalProperties: src.AdditionalProperties,
	}
}

func relationshipMetadataV1beta2ToV1alpha3(src *relationshipv1beta2.RelationshipMetadata) *relationshipv1alpha3.RelationshipMetadata {
	if src == nil {
		return nil
	}
	// Same rationale as relationshipMetadataV1alpha3ToV1beta2: the nested Styles
	// payload carries version-specific enum pointer types, so the bridge must
	// copy explicitly rather than rely on a direct struct conversion.
	return &relationshipv1alpha3.RelationshipMetadata{
		Description:          src.Description,
		IsAnnotation:         src.IsAnnotation,
		Styles:               relationshipStylesV1beta2ToV1alpha3(src.Styles),
		AdditionalProperties: src.AdditionalProperties,
	}
}

func relationshipStylesV1alpha3ToV1beta2(src *relationshipv1alpha3.RelationshipDefinitionMetadataStyles) *relationshipv1beta2.RelationshipDefinitionMetadataStyles {
	if src == nil {
		return nil
	}
	return &relationshipv1beta2.RelationshipDefinitionMetadataStyles{
		ArrowScale:          src.ArrowScale,
		Color:               src.Color,
		CurveStyle:          convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesCurveStyle, relationshipv1beta2.RelationshipDefinitionMetadataStylesCurveStyle](src.CurveStyle),
		EdgeAnimation:       src.EdgeAnimation,
		FontFamily:          src.FontFamily,
		FontSize:            src.FontSize,
		FontStyle:           src.FontStyle,
		FontWeight:          src.FontWeight,
		Label:               src.Label,
		LineCap:             convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesLineCap, relationshipv1beta2.RelationshipDefinitionMetadataStylesLineCap](src.LineCap),
		LineColor:           src.LineColor,
		LineOpacity:         src.LineOpacity,
		LineStyle:           convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesLineStyle, relationshipv1beta2.RelationshipDefinitionMetadataStylesLineStyle](src.LineStyle),
		MidTargetArrowColor: src.MidTargetArrowColor,
		MidTargetArrowFill:  convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesMidTargetArrowFill, relationshipv1beta2.RelationshipDefinitionMetadataStylesMidTargetArrowFill](src.MidTargetArrowFill),
		MidTargetArrowShape: convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesMidTargetArrowShape, relationshipv1beta2.RelationshipDefinitionMetadataStylesMidTargetArrowShape](src.MidTargetArrowShape),
		Opacity:             src.Opacity,
		PrimaryColor:        src.PrimaryColor,
		SecondaryColor:      src.SecondaryColor,
		SourceLabel:         src.SourceLabel,
		SvgColor:            src.SvgColor,
		SvgComplete:         src.SvgComplete,
		SvgWhite:            src.SvgWhite,
		TargetArrowColor:    src.TargetArrowColor,
		TargetArrowFill:     convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesTargetArrowFill, relationshipv1beta2.RelationshipDefinitionMetadataStylesTargetArrowFill](src.TargetArrowFill),
		TargetArrowShape:    convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesTargetArrowShape, relationshipv1beta2.RelationshipDefinitionMetadataStylesTargetArrowShape](src.TargetArrowShape),
		TargetLabel:         src.TargetLabel,
		TextOpacity:         src.TextOpacity,
		TextTransform:       convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionMetadataStylesTextTransform, relationshipv1beta2.RelationshipDefinitionMetadataStylesTextTransform](src.TextTransform),
		ZIndex:              src.ZIndex,
	}
}

func relationshipStylesV1beta2ToV1alpha3(src *relationshipv1beta2.RelationshipDefinitionMetadataStyles) *relationshipv1alpha3.RelationshipDefinitionMetadataStyles {
	if src == nil {
		return nil
	}
	return &relationshipv1alpha3.RelationshipDefinitionMetadataStyles{
		ArrowScale:          src.ArrowScale,
		Color:               src.Color,
		CurveStyle:          convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesCurveStyle, relationshipv1alpha3.RelationshipDefinitionMetadataStylesCurveStyle](src.CurveStyle),
		EdgeAnimation:       src.EdgeAnimation,
		FontFamily:          src.FontFamily,
		FontSize:            src.FontSize,
		FontStyle:           src.FontStyle,
		FontWeight:          src.FontWeight,
		Label:               src.Label,
		LineCap:             convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesLineCap, relationshipv1alpha3.RelationshipDefinitionMetadataStylesLineCap](src.LineCap),
		LineColor:           src.LineColor,
		LineOpacity:         src.LineOpacity,
		LineStyle:           convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesLineStyle, relationshipv1alpha3.RelationshipDefinitionMetadataStylesLineStyle](src.LineStyle),
		MidTargetArrowColor: src.MidTargetArrowColor,
		MidTargetArrowFill:  convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesMidTargetArrowFill, relationshipv1alpha3.RelationshipDefinitionMetadataStylesMidTargetArrowFill](src.MidTargetArrowFill),
		MidTargetArrowShape: convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesMidTargetArrowShape, relationshipv1alpha3.RelationshipDefinitionMetadataStylesMidTargetArrowShape](src.MidTargetArrowShape),
		Opacity:             src.Opacity,
		PrimaryColor:        src.PrimaryColor,
		SecondaryColor:      src.SecondaryColor,
		SourceLabel:         src.SourceLabel,
		SvgColor:            src.SvgColor,
		SvgComplete:         src.SvgComplete,
		SvgWhite:            src.SvgWhite,
		TargetArrowColor:    src.TargetArrowColor,
		TargetArrowFill:     convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesTargetArrowFill, relationshipv1alpha3.RelationshipDefinitionMetadataStylesTargetArrowFill](src.TargetArrowFill),
		TargetArrowShape:    convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesTargetArrowShape, relationshipv1alpha3.RelationshipDefinitionMetadataStylesTargetArrowShape](src.TargetArrowShape),
		TargetLabel:         src.TargetLabel,
		TextOpacity:         src.TextOpacity,
		TextTransform:       convertEnumPtr[relationshipv1beta2.RelationshipDefinitionMetadataStylesTextTransform, relationshipv1alpha3.RelationshipDefinitionMetadataStylesTextTransform](src.TextTransform),
		ZIndex:              src.ZIndex,
	}
}

func selectorSetV1alpha3ToV1beta2(src *relationshipv1alpha3.SelectorSet) *relationshipv1beta2.SelectorSet {
	if src == nil {
		return nil
	}
	// Rebuild the outer selector slice because the item types differ across
	// versions; inner pointer- and slice-typed fields stay aliased via the
	// nested helpers so in-place mutations remain visible to both sides.
	dst := make(relationshipv1beta2.SelectorSet, len(*src))
	for i, item := range *src {
		dst[i] = selectorSetItemV1alpha3ToV1beta2(item)
	}
	return &dst
}

func selectorSetV1beta2ToV1alpha3(src *relationshipv1beta2.SelectorSet) *relationshipv1alpha3.SelectorSet {
	if src == nil {
		return nil
	}
	// Rebuild the outer selector slice because the item types differ across
	// versions; inner pointer- and slice-typed fields stay aliased via the
	// nested helpers so in-place mutations remain visible to both sides.
	dst := make(relationshipv1alpha3.SelectorSet, len(*src))
	for i, item := range *src {
		dst[i] = selectorSetItemV1beta2ToV1alpha3(item)
	}
	return &dst
}

func selectorSetItemV1alpha3ToV1beta2(src relationshipv1alpha3.SelectorSetItem) relationshipv1beta2.SelectorSetItem {
	return relationshipv1beta2.SelectorSetItem{
		Allow: selectorV1alpha3ToV1beta2(src.Allow),
		Deny:  selectorV1alpha3ToV1beta2Ptr(src.Deny),
	}
}

func selectorSetItemV1beta2ToV1alpha3(src relationshipv1beta2.SelectorSetItem) relationshipv1alpha3.SelectorSetItem {
	return relationshipv1alpha3.SelectorSetItem{
		Allow: selectorV1beta2ToV1alpha3(src.Allow),
		Deny:  selectorV1beta2ToV1alpha3Ptr(src.Deny),
	}
}

func selectorV1alpha3ToV1beta2(src relationshipv1alpha3.Selector) relationshipv1beta2.Selector {
	return relationshipv1beta2.Selector{
		From: selectorItemsV1alpha3ToV1beta2(src.From),
		To:   selectorItemsV1alpha3ToV1beta2(src.To),
	}
}

func selectorV1beta2ToV1alpha3(src relationshipv1beta2.Selector) relationshipv1alpha3.Selector {
	return relationshipv1alpha3.Selector{
		From: selectorItemsV1beta2ToV1alpha3(src.From),
		To:   selectorItemsV1beta2ToV1alpha3(src.To),
	}
}

func selectorV1alpha3ToV1beta2Ptr(src *relationshipv1alpha3.Selector) *relationshipv1beta2.Selector {
	if src == nil {
		return nil
	}
	dst := selectorV1alpha3ToV1beta2(*src)
	return &dst
}

func selectorV1beta2ToV1alpha3Ptr(src *relationshipv1beta2.Selector) *relationshipv1alpha3.Selector {
	if src == nil {
		return nil
	}
	dst := selectorV1beta2ToV1alpha3(*src)
	return &dst
}

func selectorItemsV1alpha3ToV1beta2(src []relationshipv1alpha3.SelectorItem) []relationshipv1beta2.SelectorItem {
	if src == nil {
		return nil
	}
	dst := make([]relationshipv1beta2.SelectorItem, len(src))
	for i, item := range src {
		dst[i] = relationshipv1beta2.SelectorItem{
			ID:                                   item.ID,
			Kind:                                 item.Kind,
			Match:                                matchSelectorV1alpha3ToV1beta2(item.Match),
			MatchStrategyMatrix:                  item.MatchStrategyMatrix,
			Model:                                item.Model,
			RelationshipDefinitionSelectorsPatch: relationshipSelectorsPatchV1alpha3ToV1beta2(item.Patch),
		}
	}
	return dst
}

func selectorItemsV1beta2ToV1alpha3(src []relationshipv1beta2.SelectorItem) []relationshipv1alpha3.SelectorItem {
	if src == nil {
		return nil
	}
	dst := make([]relationshipv1alpha3.SelectorItem, len(src))
	for i, item := range src {
		dst[i] = relationshipv1alpha3.SelectorItem{
			ID:                  item.ID,
			Kind:                item.Kind,
			Match:               matchSelectorV1beta2ToV1alpha3(item.Match),
			MatchStrategyMatrix: item.MatchStrategyMatrix,
			Model:               item.Model,
			Patch:               relationshipSelectorsPatchV1beta2ToV1alpha3(item.RelationshipDefinitionSelectorsPatch),
		}
	}
	return dst
}

func matchSelectorV1alpha3ToV1beta2(src *relationshipv1alpha3.MatchSelector) *relationshipv1beta2.MatchSelector {
	if src == nil {
		return nil
	}
	return &relationshipv1beta2.MatchSelector{
		From: matchSelectorItemsV1alpha3ToV1beta2(src.From),
		Refs: src.Refs,
		To:   matchSelectorItemsV1alpha3ToV1beta2(src.To),
	}
}

func matchSelectorV1beta2ToV1alpha3(src *relationshipv1beta2.MatchSelector) *relationshipv1alpha3.MatchSelector {
	if src == nil {
		return nil
	}
	return &relationshipv1alpha3.MatchSelector{
		From: matchSelectorItemsV1beta2ToV1alpha3(src.From),
		Refs: src.Refs,
		To:   matchSelectorItemsV1beta2ToV1alpha3(src.To),
	}
}

func matchSelectorItemsV1alpha3ToV1beta2(src *[]relationshipv1alpha3.MatchSelectorItem) *[]relationshipv1beta2.MatchSelectorItem {
	if src == nil {
		return nil
	}
	dst := make([]relationshipv1beta2.MatchSelectorItem, len(*src))
	for i, item := range *src {
		dst[i] = relationshipv1beta2.MatchSelectorItem{
			ID:         item.ID,
			Kind:       item.Kind,
			MutatedRef: item.MutatedRef,
			MutatorRef: item.MutatorRef,
		}
	}
	return &dst
}

func matchSelectorItemsV1beta2ToV1alpha3(src *[]relationshipv1beta2.MatchSelectorItem) *[]relationshipv1alpha3.MatchSelectorItem {
	if src == nil {
		return nil
	}
	dst := make([]relationshipv1alpha3.MatchSelectorItem, len(*src))
	for i, item := range *src {
		dst[i] = relationshipv1alpha3.MatchSelectorItem{
			ID:         item.ID,
			Kind:       item.Kind,
			MutatedRef: item.MutatedRef,
			MutatorRef: item.MutatorRef,
		}
	}
	return &dst
}

func relationshipSelectorsPatchV1alpha3ToV1beta2(src *relationshipv1alpha3.RelationshipDefinitionSelectorsPatch) *relationshipv1beta2.RelationshipDefinitionSelectorsPatch {
	if src == nil {
		return nil
	}
	return &relationshipv1beta2.RelationshipDefinitionSelectorsPatch{
		MutatedRef:    src.MutatedRef,
		MutatorRef:    src.MutatorRef,
		PatchStrategy: convertEnumPtr[relationshipv1alpha3.RelationshipDefinitionSelectorsPatchStrategy, relationshipv1beta2.RelationshipDefinitionSelectorsPatchStrategy](src.PatchStrategy),
	}
}

func relationshipSelectorsPatchV1beta2ToV1alpha3(src *relationshipv1beta2.RelationshipDefinitionSelectorsPatch) *relationshipv1alpha3.RelationshipDefinitionSelectorsPatch {
	if src == nil {
		return nil
	}
	return &relationshipv1alpha3.RelationshipDefinitionSelectorsPatch{
		MutatedRef:    src.MutatedRef,
		MutatorRef:    src.MutatorRef,
		PatchStrategy: convertEnumPtr[relationshipv1beta2.RelationshipDefinitionSelectorsPatchStrategy, relationshipv1alpha3.RelationshipDefinitionSelectorsPatchStrategy](src.PatchStrategy),
	}
}

func convertEnumPtr[Src ~string, Dst ~string](src *Src) *Dst {
	if src == nil {
		return nil
	}
	dst := Dst(*src)
	return &dst
}
