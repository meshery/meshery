package utils

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	legacycoremodel "github.com/meshery/schemas/models/core"
	componentv1beta1 "github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	patternv1beta1 "github.com/meshery/schemas/models/v1beta1/pattern"
	coremodelv1beta2 "github.com/meshery/schemas/models/v1beta2/core"
	relationshipv1beta2 "github.com/meshery/schemas/models/v1beta2/relationship"
	designv1beta3 "github.com/meshery/schemas/models/v1beta3/design"

	"github.com/gofrs/uuid"
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

func TestPatternV1beta1ToV1beta3_ConvertsResolvedAliasesToCanonicalWireShape(t *testing.T) {
	aliasID := mustUUID(t, "27423b67-76bc-4126-a199-f6d1aa37fe58")
	parentID := mustUUID(t, "b66ef7d1-7202-4ba1-bb8a-ef3f6c9a7b7b")
	relationshipID := mustUUID(t, "618e58aa-7733-4185-a13b-799da087e614")
	legacyAliases := map[string]legacycoremodel.ResolvedAlias{
		aliasID.String(): {
			AliasComponentId:      aliasID,
			ImmediateParentId:     parentID,
			ImmediateRefFieldPath: []string{"configuration", "spec", "containers", "1"},
			RelationshipId:        relationshipID,
			ResolvedParentId:      parentID,
			ResolvedRefFieldPath:  []string{"configuration", "spec", "containers", "1"},
		},
	}
	src := &patternv1beta1.PatternFile{
		Metadata: &patternv1beta1.PatternFile_Metadata{
			ResolvedAliases: &legacyAliases,
		},
	}

	dst, err := PatternV1beta1ToV1beta3(src)
	if err != nil {
		t.Fatalf("PatternV1beta1ToV1beta3() unexpected error: %v", err)
	}
	got := (*dst.Metadata.ResolvedAliases)[aliasID.String()]
	want := coremodelv1beta2.ResolvedAlias{
		AliasComponentId:      aliasID,
		ImmediateParentId:     parentID,
		ImmediateRefFieldPath: []string{"configuration", "spec", "containers", "1"},
		RelationshipId:        relationshipID,
		ResolvedParentId:      parentID,
		ResolvedRefFieldPath:  []string{"configuration", "spec", "containers", "1"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("resolved alias = %#v, want %#v", got, want)
	}

	raw, err := json.Marshal(dst.Metadata.ResolvedAliases)
	if err != nil {
		t.Fatalf("json.Marshal() unexpected error: %v", err)
	}
	body := string(raw)
	if !strings.Contains(body, "resolvedParentId") || strings.Contains(body, "resolved_parent_id") {
		t.Fatalf("resolved aliases should marshal with canonical camelCase keys, got %s", body)
	}

	roundtripped, err := PatternV1beta3ToV1beta1(dst)
	if err != nil {
		t.Fatalf("PatternV1beta3ToV1beta1() unexpected error: %v", err)
	}
	legacyGot := (*roundtripped.Metadata.ResolvedAliases)[aliasID.String()]
	if !reflect.DeepEqual(legacyGot, legacyAliases[aliasID.String()]) {
		t.Fatalf("round-tripped legacy alias = %#v, want %#v", legacyGot, legacyAliases[aliasID.String()])
	}
}

func newPatternV1beta1Fixture() *patternv1beta1.PatternFile {
	return &patternv1beta1.PatternFile{
		Name:          "typed-bridge",
		SchemaVersion: "designs.meshery.io/v1beta1",
		Version:       "0.0.1",
		Metadata: &patternv1beta1.PatternFile_Metadata{
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

func mustUUID(t *testing.T, raw string) uuid.UUID {
	t.Helper()
	id, err := uuid.FromString(raw)
	if err != nil {
		t.Fatalf("uuid.FromString(%q): %v", raw, err)
	}
	return id
}
