package utils

import (
	"encoding/json"
	"testing"

	"github.com/gofrs/uuid"
	core "github.com/meshery/schemas/models/core"
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	patternv1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	relationshipv1beta2 "github.com/meshery/schemas/models/v1beta2/relationship"
	designv1beta3 "github.com/meshery/schemas/models/v1beta3/design"
)

func TestPatternV1beta1ToV1beta3_PreservesAliasedFields(t *testing.T) {
	src := newPatternV1beta1Fixture()

	dst, err := PatternV1beta1ToV1beta3(src)
	if err != nil {
		t.Fatalf("PatternV1beta1ToV1beta3() unexpected error: %v", err)
	}
	if dst == nil {
		t.Fatal("PatternV1beta1ToV1beta3() returned nil")
	}
	if got, want := dst.Name, src.Name; got != want {
		t.Fatalf("Name = %q, want %q", got, want)
	}
	if len(dst.Components) != 1 {
		t.Fatalf("expected 1 component, got %d", len(dst.Components))
	}
	if dst.Relationships[0] != src.Relationships[0] {
		t.Fatal("relationship pointers should stay aliased across the pattern bridge")
	}
	if got := dst.Metadata.AdditionalProperties["team"]; got != "meshery" {
		t.Fatalf("metadata additionalProperties not preserved, got %#v", got)
	}
	if got := (*dst.Metadata.ResolvedAliases)["service"].ImmediateParentId.String(); got != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("resolved alias metadata not preserved, got %q", got)
	}
	if got := dst.Components[0].Metadata.AdditionalProperties["componentKey"]; got != "componentValue" {
		t.Fatalf("component metadata additionalProperties not preserved, got %#v", got)
	}

	dst.Preferences.Layers["annotations"] = true
	if got := src.Preferences.Layers["annotations"]; got != true {
		t.Fatalf("preference layers mutation did not reflect on source, got %#v", got)
	}

	dst.Metadata.AdditionalProperties["owner"] = "catalog"
	if got := src.Metadata.AdditionalProperties["owner"]; got != "catalog" {
		t.Fatalf("metadata alias lost, got %#v", got)
	}
	resolvedAlias := (*dst.Metadata.ResolvedAliases)["service"]
	resolvedAlias.ImmediateRefFieldPath[0] = "spec.template.metadata.labels.app"
	if got := (*src.Metadata.ResolvedAliases)["service"].ImmediateRefFieldPath[0]; got != "spec.template.metadata.labels.app" {
		t.Fatalf("resolved alias slice alias lost, got %#v", got)
	}

	dst.Components[0].Configuration["replicas"] = 3
	if got := src.Components[0].Configuration["replicas"]; got != 3 {
		t.Fatalf("component configuration alias lost, got %#v", got)
	}

	dst.Components[0].Metadata.AdditionalProperties["ui"] = "dark"
	if got := src.Components[0].Metadata.AdditionalProperties["ui"]; got != "dark" {
		t.Fatalf("component metadata alias lost, got %#v", got)
	}

	roundtripped, err := PatternV1beta3ToV1beta1(dst)
	if err != nil {
		t.Fatalf("PatternV1beta3ToV1beta1() unexpected error: %v", err)
	}
	if roundtripped == nil {
		t.Fatal("PatternV1beta3ToV1beta1() returned nil")
	}
	if got := roundtripped.Components[0].Configuration["replicas"]; got != 3 {
		t.Fatalf("round-tripped component configuration missing mutation, got %#v", got)
	}
	if got := roundtripped.Metadata.AdditionalProperties["owner"]; got != "catalog" {
		t.Fatalf("round-tripped metadata missing mutation, got %#v", got)
	}
}

func TestPatternV1beta1ToV1beta3_PreservesAliasContractLostByJSON(t *testing.T) {
	src := newPatternV1beta1Fixture()

	typed, err := PatternV1beta1ToV1beta3(src)
	if err != nil {
		t.Fatalf("PatternV1beta1ToV1beta3() unexpected error: %v", err)
	}

	raw, err := json.Marshal(src)
	if err != nil {
		t.Fatalf("json.Marshal() unexpected error: %v", err)
	}
	var legacy designv1beta3.PatternFile
	if err := json.Unmarshal(raw, &legacy); err != nil {
		t.Fatalf("json.Unmarshal() unexpected error: %v", err)
	}

	typed.Preferences.Layers["typed-layer"] = "present"
	if got := src.Preferences.Layers["typed-layer"]; got != "present" {
		t.Fatalf("typed bridge lost preference aliasing, got %#v", got)
	}
	legacy.Preferences.Layers["json-layer"] = "detached"
	if _, ok := src.Preferences.Layers["json-layer"]; ok {
		t.Fatal("legacy JSON path unexpectedly preserved preference aliasing")
	}

	typed.Components[0].Configuration["typed-config"] = "present"
	if got := src.Components[0].Configuration["typed-config"]; got != "present" {
		t.Fatalf("typed bridge lost component configuration aliasing, got %#v", got)
	}
	legacy.Components[0].Configuration["json-config"] = "detached"
	if _, ok := src.Components[0].Configuration["json-config"]; ok {
		t.Fatal("legacy JSON path unexpectedly preserved component configuration aliasing")
	}
}

func newPatternV1beta1Fixture() *patternv1beta1.PatternFile {
	return &patternv1beta1.PatternFile{
		Name:          "typed-bridge",
		SchemaVersion: "designs.meshery.io/v1beta1",
		Version:       "0.0.1",
		Metadata: &patternv1beta1.PatternFile_Metadata{
			ResolvedAliases: &map[string]core.ResolvedAlias{
				"service": {
					AliasComponentId:      uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000001")),
					ImmediateParentId:     uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111")),
					ImmediateRefFieldPath: []string{"spec.template.spec.containers.0.image"},
					RelationshipId:        uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222")),
					ResolvedParentId:      uuid.Must(uuid.FromString("33333333-3333-3333-3333-333333333333")),
					ResolvedRefFieldPath:  []string{"spec.template.metadata.labels"},
				},
			},
			AdditionalProperties: map[string]interface{}{
				"team": "meshery",
			},
		},
		Components: []*componentv1beta1.ComponentDefinition{
			{
				DisplayName: "nginx",
				ModelReference: modelv1beta1.ModelReference{
					Name: "kubernetes",
				},
				Metadata: componentv1beta1.ComponentDefinition_Metadata{
					AdditionalProperties: map[string]interface{}{
						"componentKey": "componentValue",
					},
				},
				Configuration: map[string]interface{}{
					"replicas": 1,
				},
				Component: componentv1beta1.Component{
					Kind:    "Deployment",
					Version: "apps/v1",
				},
			},
		},
		Preferences: &patternv1beta1.DesignPreferences{
			Layers: map[string]interface{}{
				"relationships": false,
			},
		},
		Relationships: []*relationshipv1beta2.RelationshipDefinition{
			{
				Kind:             relationshipv1beta2.RelationshipDefinitionKind("hierarchical"),
				RelationshipType: "parent",
				SubType:          "inventory",
				Version:          "0.0.1",
			},
		},
	}
}
