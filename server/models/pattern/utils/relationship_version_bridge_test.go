package utils

import (
	"encoding/json"
	"testing"

	capabilityv1alpha1 "github.com/meshery/schemas/models/v1alpha1/capability"
	relationshipv1alpha3 "github.com/meshery/schemas/models/v1alpha3/relationship"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	relationshipv1beta2 "github.com/meshery/schemas/models/v1beta2/relationship"
)

func TestRelationshipV1alpha3ToV1beta2_PreservesAliasedFields(t *testing.T) {
	src := newRelationshipV1alpha3Fixture()

	dst := RelationshipV1alpha3ToV1beta2(src)
	if dst == nil {
		t.Fatal("RelationshipV1alpha3ToV1beta2() returned nil")
	}
	if got, want := dst.RelationshipType, src.RelationshipType; got != want {
		t.Fatalf("RelationshipType = %q, want %q", got, want)
	}
	if dst.RelationshipMetadata == nil || dst.RelationshipMetadata.Styles == nil {
		t.Fatal("typed bridge lost relationship metadata/styles")
	}
	if got := dst.RelationshipMetadata.AdditionalProperties["owner"]; got != "meshery" {
		t.Fatalf("metadata additionalProperties not preserved, got %#v", got)
	}
	if got := *dst.RelationshipMetadata.Styles.SourceLabel; got != "mounted-by" {
		t.Fatalf("SourceLabel = %q, want %q", got, "mounted-by")
	}
	if got := string(*dst.RelationshipMetadata.Styles.CurveStyle); got != "bezier" {
		t.Fatalf("CurveStyle = %q, want %q", got, "bezier")
	}
	if got := (*dst.Selectors)[0].Allow.From[0].RelationshipDefinitionSelectorsPatch; got == nil {
		t.Fatal("typed bridge lost selector patch")
	}

	dst.RelationshipMetadata.AdditionalProperties["audited"] = true
	if got := src.Metadata.AdditionalProperties["audited"]; got != true {
		t.Fatalf("metadata alias lost, got %#v", got)
	}

	*dst.RelationshipMetadata.Styles.SourceLabel = "depends-on"
	if got := *src.Metadata.Styles.SourceLabel; got != "depends-on" {
		t.Fatalf("style pointer alias lost, got %q", got)
	}

	if (*dst.Capabilities)[0].Metadata == nil {
		t.Fatal("typed bridge lost capability metadata")
	}
	(*(*dst.Capabilities)[0].Metadata)["scope"] = "namespace"
	if got := (*(*src.Capabilities)[0].Metadata)["scope"]; got != "namespace" {
		t.Fatalf("capability metadata alias lost, got %#v", got)
	}

	(*dst.Selectors)[0].Allow.From[0].Model.Name = "meshery"
	if got := (*src.Selectors)[0].Allow.From[0].Model.Name; got != "meshery" {
		t.Fatalf("selector model pointer alias lost, got %q", got)
	}

	(*(*dst.Selectors)[0].Allow.From[0].MatchStrategyMatrix)[0][0] = "from.metadata.name"
	if got := (*(*src.Selectors)[0].Allow.From[0].MatchStrategyMatrix)[0][0]; got != "from.metadata.name" {
		t.Fatalf("selector matchStrategyMatrix alias lost, got %q", got)
	}

	(*(*dst.Selectors)[0].Allow.From[0].RelationshipDefinitionSelectorsPatch.MutatedRef)[0][0] = "status"
	if got := (*(*src.Selectors)[0].Allow.From[0].Patch.MutatedRef)[0][0]; got != "status" {
		t.Fatalf("selector mutatedRef alias lost, got %q", got)
	}

	roundtripped := RelationshipV1beta2ToV1alpha3(dst)
	if roundtripped == nil {
		t.Fatal("RelationshipV1beta2ToV1alpha3() returned nil")
	}
	if got := roundtripped.Metadata.AdditionalProperties["audited"]; got != true {
		t.Fatalf("round-tripped metadata missing aliased mutation, got %#v", got)
	}
	if got := *roundtripped.Metadata.Styles.SourceLabel; got != "depends-on" {
		t.Fatalf("round-tripped SourceLabel = %q, want %q", got, "depends-on")
	}
	if got := (*(*roundtripped.Selectors)[0].Allow.From[0].Patch.MutatedRef)[0][0]; got != "status" {
		t.Fatalf("round-tripped MutatedRef = %q, want %q", got, "status")
	}
}

func TestRelationshipV1alpha3ToV1beta2_PreservesFieldsDroppedByJSON(t *testing.T) {
	src := newRelationshipV1alpha3Fixture()

	typed := RelationshipV1alpha3ToV1beta2(src)

	raw, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("json.Marshal() unexpected error: %v", err)
	}
	var legacy relationshipv1beta2.RelationshipDefinition
	if err := json.Unmarshal(raw, &legacy); err != nil {
		t.Fatalf("json.Unmarshal() unexpected error: %v", err)
	}

	if legacy.RelationshipMetadata != nil && legacy.RelationshipMetadata.Styles != nil && legacy.RelationshipMetadata.Styles.SourceLabel != nil {
		t.Fatalf("legacy JSON path unexpectedly preserved SourceLabel: %q", *legacy.RelationshipMetadata.Styles.SourceLabel)
	}
	if legacy.Selectors != nil && (*legacy.Selectors)[0].Allow.From[0].MatchStrategyMatrix != nil {
		t.Fatalf("legacy JSON path unexpectedly preserved MatchStrategyMatrix: %#v", *(*legacy.Selectors)[0].Allow.From[0].MatchStrategyMatrix)
	}

	if typed.RelationshipMetadata == nil || typed.RelationshipMetadata.Styles == nil || typed.RelationshipMetadata.Styles.SourceLabel == nil {
		t.Fatal("typed bridge lost SourceLabel")
	}
	if got := *typed.RelationshipMetadata.Styles.SourceLabel; got != "mounted-by" {
		t.Fatalf("typed bridge SourceLabel = %q, want %q", got, "mounted-by")
	}
	if typed.Selectors == nil || (*typed.Selectors)[0].Allow.From[0].MatchStrategyMatrix == nil {
		t.Fatal("typed bridge lost MatchStrategyMatrix")
	}
}

func newRelationshipV1alpha3Fixture() *relationshipv1alpha3.RelationshipDefinition {
	description := "relationship bridge"
	isAnnotation := true
	sourceLabel := "mounted-by"
	fontFamily := "Inter"
	targetArrowColor := "#00B39F"
	curveStyle := relationshipv1alpha3.RelationshipDefinitionMetadataStylesCurveStyle("bezier")
	status := relationshipv1alpha3.RelationshipDefinitionStatus("approved")
	patchStrategy := relationshipv1alpha3.RelationshipDefinitionSelectorsPatchStrategy("copy")
	kind := "Deployment"
	matchStrategyMatrix := [][]string{{"from.kind", "to.kind"}}
	capabilityMetadata := map[string]interface{}{"scope": "cluster"}
	modelRef := &modelv1beta1.ModelReference{Name: "kubernetes"}
	mutatedRef := relationshipv1alpha3.MutatedRef{{"spec", "template"}}
	mutatorRef := relationshipv1alpha3.MutatorRef{{"metadata", "labels", "app"}}

	return &relationshipv1alpha3.RelationshipDefinition{
		Capabilities: &[]capabilityv1alpha1.Capability{
			{
				DisplayName: "Patch",
				Key:         "patch",
				Kind:        "action",
				Type:        "mutation",
				SubType:     "copy",
				EntityState: []string{"enabled"},
				Metadata:    &capabilityMetadata,
				Status:      capabilityv1alpha1.Enabled,
				Version:     "0.0.1",
			},
		},
		EvaluationQuery: stringPtr("rego.rule"),
		Kind:            relationshipv1alpha3.RelationshipDefinitionKind("edge"),
		Metadata: &relationshipv1alpha3.RelationshipMetadata{
			Description:  &description,
			IsAnnotation: &isAnnotation,
			Styles: &relationshipv1alpha3.RelationshipDefinitionMetadataStyles{
				CurveStyle:       &curveStyle,
				FontFamily:       &fontFamily,
				PrimaryColor:     "#00B39F",
				SourceLabel:      &sourceLabel,
				SvgColor:         "#111111",
				SvgWhite:         "#FFFFFF",
				TargetArrowColor: &targetArrowColor,
			},
			AdditionalProperties: map[string]interface{}{
				"owner": "meshery",
			},
		},
		Model:         modelv1beta1.ModelReference{Name: "kubernetes"},
		SchemaVersion: "relationships.meshery.io/v1alpha3",
		Selectors: &relationshipv1alpha3.SelectorSet{
			{
				Allow: relationshipv1alpha3.Selector{
					From: []relationshipv1alpha3.SelectorItem{
						{
							Kind:                &kind,
							MatchStrategyMatrix: &matchStrategyMatrix,
							Model:               modelRef,
							Patch: &relationshipv1alpha3.RelationshipDefinitionSelectorsPatch{
								MutatedRef:    &mutatedRef,
								MutatorRef:    &mutatorRef,
								PatchStrategy: &patchStrategy,
							},
						},
					},
					To: []relationshipv1alpha3.SelectorItem{
						{
							Kind:  &kind,
							Model: modelRef,
						},
					},
				},
			},
		},
		SubType:          "mount",
		Status:           &status,
		RelationshipType: "non-binding",
		Version:          "0.0.1",
	}
}

func stringPtr(v string) *string {
	return &v
}
