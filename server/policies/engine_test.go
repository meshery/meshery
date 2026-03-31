package policies

import (
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
)

// strPtr returns a pointer to the given string.
func strPtr(s string) *string { return &s }

// makePatternFile builds a PatternFile from typed components and relationships.
func makePatternFile(comps []*component.ComponentDefinition, rels []*relationship.RelationshipDefinition) *pattern.PatternFile {
	return &pattern.PatternFile{
		Components:    comps,
		Relationships: rels,
	}
}

func TestObjectGetNested(t *testing.T) {
	obj := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": "value",
			},
		},
		"arr": []interface{}{"x", "y", "z"},
	}

	tests := []struct {
		name     string
		path     []string
		expected interface{}
	}{
		{"simple", []string{"a", "b", "c"}, "value"},
		{"missing", []string{"a", "x"}, "default"},
		{"empty path", nil, obj},
		{"array index", []string{"arr", "1"}, "y"},
		{"array out of bounds", []string{"arr", "5"}, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := objectGetNested(obj, tt.path, "default")
			if result == nil && tt.expected == nil {
				return
			}
			// For map comparisons, just check non-nil
			if _, ok := result.(map[string]interface{}); ok {
				return
			}
			if result != tt.expected {
				t.Errorf("objectGetNested(%v) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestMatchName(t *testing.T) {
	tests := []struct {
		name, pattern string
		expected      bool
	}{
		{"Namespace", "*", true},
		{"Namespace", "Namespace", true},
		{"Namespace", "Service", false},
		{"Deployment", "Deploy.*", true},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_"+tt.pattern, func(t *testing.T) {
			result := matchName(tt.name, tt.pattern)
			if result != tt.expected {
				t.Errorf("matchName(%q, %q) = %v, want %v", tt.name, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestMatchValues(t *testing.T) {
	tests := []struct {
		name     string
		from, to interface{}
		strategy string
		expected bool
	}{
		{"equal_match", "test", "test", "equal", true},
		{"equal_mismatch", "test", "other", "equal", false},
		{"equal_as_strings", 42, "42", "equal_as_strings", true},
		{"not_null_both", "a", "b", "not_null", true},
		{"not_null_nil", nil, "b", "not_null", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchValues(tt.from, tt.to, tt.strategy)
			if result != tt.expected {
				t.Errorf("matchValues(%v, %v, %q) = %v, want %v", tt.from, tt.to, tt.strategy, result, tt.expected)
			}
		})
	}
}

func TestIsRelationshipFeasible(t *testing.T) {
	comp := &component.ComponentDefinition{
		Component:       component.Component{Kind: "Namespace"},
		ModelReference:  modelv1beta1.ModelReference{Name: "kubernetes"},
	}

	tests := []struct {
		name     string
		selector relationship.SelectorItem
		expected bool
	}{
		{
			"wildcard",
			relationship.SelectorItem{Kind: strPtr("*"), Model: &modelv1beta1.ModelReference{Name: "*"}},
			true,
		},
		{
			"exact match",
			relationship.SelectorItem{Kind: strPtr("Namespace"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
			true,
		},
		{
			"kind mismatch",
			relationship.SelectorItem{Kind: strPtr("Service"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isRelationshipFeasible(tt.selector, comp)
			if result != tt.expected {
				t.Errorf("isRelationshipFeasible() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsRelationshipDenied(t *testing.T) {
	fromDecl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
	}
	toDecl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
	}

	denySelectors := &relationship.Selector{
		From: []relationship.SelectorItem{
			{Kind: strPtr("Namespace"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
		},
		To: []relationship.SelectorItem{
			{Kind: strPtr("Namespace"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
		},
	}

	result := isRelationshipDenied(fromDecl, toDecl, denySelectors)
	if !result {
		t.Error("Expected Namespace->Namespace to be denied, but it was not")
	}
}

func TestIsAliasRelationship(t *testing.T) {
	tests := []struct {
		name     string
		rel      *relationship.RelationshipDefinition
		expected bool
	}{
		{
			"alias",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
				RelationshipType: "parent",
				SubType:          "alias",
			},
			true,
		},
		{
			"alias uppercase",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("Hierarchical"),
				RelationshipType: "Parent",
				SubType:          "Alias",
			},
			true,
		},
		{
			"inventory",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
				RelationshipType: "parent",
				SubType:          "inventory",
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAliasRelationship(tt.rel)
			if result != tt.expected {
				t.Errorf("isAliasRelationship() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPolicyImplication(t *testing.T) {
	policies := []RelationshipPolicy{
		&HierarchicalParentChildPolicy{},
		&EdgeNonBindingPolicy{},
		&AliasPolicy{},
		&MatchLabelsPolicy{},
	}

	tests := []struct {
		name           string
		rel            *relationship.RelationshipDefinition
		expectedPolicy string
	}{
		{
			"hierarchical inventory",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
				RelationshipType: "parent",
				SubType:          "inventory",
			},
			"hierarchical_parent_child",
		},
		{
			"edge non-binding",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("edge"),
				RelationshipType: "non-binding",
			},
			"edge-non-binding",
		},
		{
			"alias",
			&relationship.RelationshipDefinition{
				Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
				RelationshipType: "parent",
				SubType:          "alias",
			},
			"alias_relationships_policy",
		},
		{
			"sibling",
			&relationship.RelationshipDefinition{
				RelationshipType: "sibling",
			},
			"sibling_match_labels_policy",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, p := range policies {
				if p.IsImplicatedBy(tt.rel) {
					if p.Identifier() != tt.expectedPolicy {
						t.Errorf("Expected policy %q, got %q", tt.expectedPolicy, p.Identifier())
					}
					found = true
					break
				}
			}
			if !found {
				t.Errorf("No policy matched for rel %v", tt.rel)
			}
		})
	}
}

func TestApplyAllActionsToDesign(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "comp-1",
				"component": map[string]interface{}{"kind": "Namespace"},
			},
			map[string]interface{}{
				"id":        "comp-2",
				"component": map[string]interface{}{"kind": "Service"},
			},
		},
		"relationships": []interface{}{
			map[string]interface{}{
				"id":     "rel-1",
				"status": "approved",
			},
		},
	}

	actions := []PolicyAction{
		{
			Op: DeleteComponentOp,
			Value: map[string]interface{}{
				"id": "comp-2",
			},
		},
		{
			Op: AddRelationshipOp,
			Value: map[string]interface{}{
				"item": map[string]interface{}{
					"id":     "rel-2",
					"status": "identified",
				},
			},
		},
		{
			Op: UpdateRelationshipOp,
			Value: map[string]interface{}{
				"id":    "rel-1",
				"path":  "/status",
				"value": "deleted",
			},
		},
	}

	result := applyAllActionsToDesign(design, actions)

	comps := extractMapSlice(result, "components")
	if len(comps) != 1 {
		t.Fatalf("Expected 1 component after delete, got %d", len(comps))
	}
	if getMapString(comps[0], "id") != "comp-1" {
		t.Error("Wrong component remained after delete")
	}

	rels := extractMapSlice(result, "relationships")
	if len(rels) != 2 {
		t.Fatalf("Expected 2 relationships, got %d", len(rels))
	}

	// Check rel-1 was updated
	for _, rel := range rels {
		if getMapString(rel, "id") == "rel-1" {
			if getMapString(rel, "status") != "deleted" {
				t.Error("Expected rel-1 status to be 'deleted'")
			}
		}
	}
}

func TestFromAndToComponentsExist(t *testing.T) {
	compA := &component.ComponentDefinition{}
	compA.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000001")

	compB := &component.ComponentDefinition{}
	compB.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000002")

	design := makePatternFile([]*component.ComponentDefinition{compA, compB}, nil)

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{{ID: &compA.ID}},
				To:   []relationship.SelectorItem{{ID: &compB.ID}},
			},
		},
	}
	rel := &relationship.RelationshipDefinition{Selectors: &selectorSet}

	if !fromAndToComponentsExist(rel, design) {
		t.Error("Expected both components to exist")
	}

	missingID, _ := uuid.FromString("00000000-0000-0000-0000-000000000099")
	selectorSetMissing := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{{ID: &compA.ID}},
				To:   []relationship.SelectorItem{{ID: &missingID}},
			},
		},
	}
	relMissing := &relationship.RelationshipDefinition{Selectors: &selectorSetMissing}

	if fromAndToComponentsExist(relMissing, design) {
		t.Error("Expected missing component to be detected")
	}
}

func TestValidateRelationshipsInDesign(t *testing.T) {
	compA := &component.ComponentDefinition{}
	compA.ID, _ = uuid.FromString("00000000-0000-0000-0000-0000000000aa")

	compDeleted := &component.ComponentDefinition{}
	compDeleted.ID, _ = uuid.FromString("00000000-0000-0000-0000-00000000dead")

	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{{ID: &compA.ID}},
				To:   []relationship.SelectorItem{{ID: &compDeleted.ID}},
			},
		},
	}
	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "non-binding",
		Status:           &relStatus,
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000001")

	// design only has compA, not compDeleted
	design := makePatternFile([]*component.ComponentDefinition{compA}, []*relationship.RelationshipDefinition{rel})

	policy := &EdgeNonBindingPolicy{}
	actions := validateRelationshipsInDesign(design, policy)

	if len(actions) != 1 {
		t.Fatalf("Expected 1 validation action, got %d", len(actions))
	}
	if actions[0].Op != UpdateRelationshipOp {
		t.Error("Expected update_relationship action")
	}
	if actions[0].Value["value"] != "deleted" {
		t.Error("Expected status to be set to 'deleted'")
	}
}

func TestEdgeBindingPolicyImplication(t *testing.T) {
	p := &EdgeBindingPolicy{}

	binding := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "binding",
	}
	if !p.IsImplicatedBy(binding) {
		t.Error("Expected edge binding to be implicated")
	}

	nonBinding := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "non-binding",
	}
	if p.IsImplicatedBy(nonBinding) {
		t.Error("Expected edge non-binding to not be implicated")
	}
}

func TestIsValidBinding(t *testing.T) {
	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	// role has "name" in configuration
	role := &component.ComponentDefinition{
		Component:     component.Component{Kind: "Role"},
		Configuration: map[string]interface{}{"name": "my-role"},
	}
	role.ID = roleID

	// roleBinding has "roleRef.name" in configuration
	roleBinding := &component.ComponentDefinition{
		Component:     component.Component{Kind: "RoleBinding"},
		Configuration: map[string]interface{}{"roleRef": map[string]interface{}{"name": "my-role"}},
	}
	roleBinding.ID = rbID

	// mutatorRef/mutatedRef paths are relative to the component map root (not configuration prefix)
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "roleRef", "name"}}
	selector := relationship.SelectorItem{
		Kind: strPtr("Role"),
		Match: &relationship.MatchSelector{
			From: &[]relationship.MatchSelectorItem{
				{Kind: "self", MutatorRef: &mutatorRef},
			},
			To: &[]relationship.MatchSelectorItem{
				{Kind: "RoleBinding", MutatedRef: &mutatedRef},
			},
		},
	}

	design := makePatternFile([]*component.ComponentDefinition{role, roleBinding}, nil)

	if !isValidBindingTyped(role, roleBinding, selector, design) {
		t.Error("Expected valid binding")
	}

	// Mismatched value
	badRbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")
	badBinding := &component.ComponentDefinition{
		Component:     component.Component{Kind: "RoleBinding"},
		Configuration: map[string]interface{}{"roleRef": map[string]interface{}{"name": "other-role"}},
	}
	badBinding.ID = badRbID

	design2 := makePatternFile([]*component.ComponentDefinition{role, badBinding}, nil)
	if isValidBindingTyped(role, badBinding, selector, design2) {
		t.Error("Expected invalid binding")
	}
}

func TestHierarchicalIdentifyRelationship(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	ns := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"name": "default"},
	}
	ns.ID = nsID

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "default"},
	}
	deploy.ID = deployID

	design := makePatternFile(
		[]*component.ComponentDefinition{ns, deploy},
		[]*relationship.RelationshipDefinition{},
	)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind:  strPtr("Namespace"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Patch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("Deployment"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Patch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}
	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("Expected 1 identified relationship, got %d", len(identified))
	}
	if getRelStatus(identified[0]) != "approved" {
		t.Error("Expected status to be 'approved'")
	}
}

func TestHierarchicalSideEffects(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	ns := &component.ComponentDefinition{
		Component:     component.Component{Kind: "Namespace"},
		Configuration: map[string]interface{}{"name": "production"},
	}
	ns.ID = nsID

	deploy := &component.ComponentDefinition{
		Component:     component.Component{Kind: "Deployment"},
		Configuration: map[string]interface{}{"namespace": "default"},
	}
	deploy.ID = deployID

	design := makePatternFile([]*component.ComponentDefinition{ns, deploy}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID: &nsID,
						Patch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID: &deployID,
						Patch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}
	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Status:           &relStatus,
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000001")

	actions := p.SideEffects(rel, design)
	if len(actions) == 0 {
		t.Fatal("Expected side effect actions")
	}
	if actions[0].Value["value"] != "production" {
		t.Errorf("Expected namespace to be updated to 'production', got %v", actions[0].Value["value"])
	}
}

func TestGoEngineCreation(t *testing.T) {
	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	if engine == nil {
		t.Fatal("Expected non-nil engine")
	}
	if len(engine.policies) != 6 {
		t.Errorf("Expected 6 policies, got %d", len(engine.policies))
	}
}

// --- utils.go tests ---

func TestPopLast(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{"multiple", []string{"a", "b", "c"}, []string{"a", "b"}},
		{"two", []string{"a", "b"}, []string{"a"}},
		{"single", []string{"a"}, nil},
		{"empty", []string{}, nil},
		{"nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := popLast(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("popLast(%v) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("popLast(%v)[%d] = %v, want %v", tt.input, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestIsDirectReference(t *testing.T) {
	tests := []struct {
		name     string
		ref      []string
		expected bool
	}{
		{"all named segments is direct", []string{"configuration", "name"}, true},
		{"wildcard segment makes ref indirect", []string{"configuration", "containers", "_"}, false},
		{"empty ref is direct", []string{}, true},
		{"single named segment is direct", []string{"name"}, true},
		{"single wildcard segment is indirect", []string{"_"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isDirectReference(tt.ref)
			if result != tt.expected {
				t.Errorf("isDirectReference(%v) = %v, want %v", tt.ref, result, tt.expected)
			}
		})
	}
}

func TestStringSliceToInterface(t *testing.T) {
	tests := []struct {
		name  string
		input []string
	}{
		{"empty slice returns empty interface slice", []string{}},
		{"single element string converts to interface", []string{"a"}},
		{"multiple elements preserve order and values", []string{"a", "b", "c"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringSliceToInterface(tt.input)
			if len(result) != len(tt.input) {
				t.Fatalf("Expected %d elements, got %d", len(tt.input), len(result))
			}
			for i, v := range result {
				s, ok := v.(string)
				if !ok {
					t.Errorf("Element %d is not a string", i)
				}
				if s != tt.input[i] {
					t.Errorf("Element %d = %q, want %q", i, s, tt.input[i])
				}
			}
		})
	}
}

func TestRelationshipPreferenceKey(t *testing.T) {
	tests := []struct {
		name     string
		rel      map[string]interface{}
		expected string
	}{
		{
			"hierarchical inventory",
			map[string]interface{}{"kind": "Hierarchical", "type": "Parent", "subType": "Inventory"},
			"hierarchical-parent-inventory",
		},
		{
			"edge non-binding",
			map[string]interface{}{"kind": "Edge", "type": "Non-Binding", "subType": ""},
			"edge-non-binding-",
		},
		{
			"empty",
			map[string]interface{}{},
			"--",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := relationshipPreferenceKey(tt.rel)
			if result != tt.expected {
				t.Errorf("relationshipPreferenceKey() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestResolvePath(t *testing.T) {
	obj := map[string]interface{}{
		"containers": []interface{}{"nginx", "sidecar", "init"},
	}

	t.Run("no wildcard", func(t *testing.T) {
		path := []string{"name"}
		result := resolvePath(path, obj)
		if len(result) != 1 || result[0] != "name" {
			t.Errorf("Expected [name], got %v", result)
		}
	})

	t.Run("wildcard present", func(t *testing.T) {
		path := []string{"containers", "_", "name"}
		result := resolvePath(path, obj)
		// Should resolve "_" to last index (2)
		if len(result) != 3 || result[1] != "2" {
			t.Errorf("Expected wildcard resolved to index 2, got %v", result)
		}
	})

	t.Run("wildcard empty array", func(t *testing.T) {
		emptyObj := map[string]interface{}{
			"containers": []interface{}{},
		}
		path := []string{"containers", "_"}
		result := resolvePath(path, emptyObj)
		// Empty array should resolve to "0"
		if len(result) != 2 || result[1] != "0" {
			t.Errorf("Expected wildcard resolved to 0 for empty array, got %v", result)
		}
	})

	t.Run("wildcard non-array target", func(t *testing.T) {
		nonArrayObj := map[string]interface{}{
			"name": "test",
		}
		path := []string{"name", "_"}
		result := resolvePath(path, nonArrayObj)
		// Non-array should resolve to "0"
		if len(result) != 2 || result[1] != "0" {
			t.Errorf("Expected wildcard resolved to 0 for non-array, got %v", result)
		}
	})
}

func TestCanonicalSeedDeterminism(t *testing.T) {
	map1 := map[string]interface{}{"a": "1", "b": "2"}
	map2 := map[string]interface{}{"a": "1", "b": "2"}

	s1 := canonicalSeed(map1)
	s2 := canonicalSeed(map2)
	if s1 != s2 {
		t.Errorf("canonicalSeed not deterministic: %q != %q", s1, s2)
	}

	// Different maps should produce different seeds
	map3 := map[string]interface{}{"a": "1", "b": "3"}
	s3 := canonicalSeed(map3)
	if s1 == s3 {
		t.Error("Expected different seeds for different maps")
	}

	// String seed
	s4 := canonicalSeed("hello")
	s5 := canonicalSeed("hello")
	if s4 != s5 {
		t.Errorf("canonicalSeed not deterministic for strings: %q != %q", s4, s5)
	}
}

func TestStaticUUIDDeterminism(t *testing.T) {
	seed := map[string]interface{}{"key": "value"}
	u1 := staticUUID(seed)
	u2 := staticUUID(seed)
	if u1 != u2 {
		t.Errorf("staticUUID not deterministic: %v != %v", u1, u2)
	}
}

func TestNewUUIDNonDeterminism(t *testing.T) {
	seed := "test-seed"
	u1 := newUUID(seed)
	time.Sleep(time.Millisecond) // ensure UnixNano changes between calls
	u2 := newUUID(seed)
	// UUIDs should differ because newUUID incorporates time
	if u1 == u2 {
		t.Error("newUUID should produce different UUIDs on successive calls (time-based)")
	}
}

// --- policy_alias.go tests ---

func TestAliasIsInvalidStatusMatrix(t *testing.T) {
	p := &AliasPolicy{}

	fromID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	toID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	fromComp := &component.ComponentDefinition{
		Component: component.Component{Kind: "Container"},
		Configuration: map[string]interface{}{
			"ports": []interface{}{map[string]interface{}{"containerPort": 80}},
		},
	}
	fromComp.ID = fromID

	toComp := &component.ComponentDefinition{
		Component: component.Component{Kind: "Deployment"},
		Configuration: map[string]interface{}{
			"ports": []interface{}{map[string]interface{}{"containerPort": 80}},
		},
	}
	toComp.ID = toID

	makeRel := func(status string) *relationship.RelationshipDefinition {
		mutatorRef := relationship.MutatorRef{[]string{"configuration", "ports"}}
		relStatus := relationship.RelationshipDefinitionStatus(status)
		selectorSet := relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							ID: &fromID,
							Patch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRef,
							},
						},
					},
					To: []relationship.SelectorItem{
						{ID: &toID},
					},
				},
			},
		}
		rel := &relationship.RelationshipDefinition{
			Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
			RelationshipType: "parent",
			SubType:          "alias",
			Status:           &relStatus,
			Selectors:        &selectorSet,
		}
		rel.ID = staticUUID(fmt.Sprintf("test-alias-%s", status))
		return rel
	}

	tests := []struct {
		status   string
		expected bool // IsInvalid returns true when invalid
	}{
		{"pending", false},  // pending is valid
		{"approved", false}, // approved with existing components is valid
		{"deleted", true},   // deleted is invalid
		{"malformed", true}, // unknown status is invalid
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			rel := makeRel(tt.status)
			design := makePatternFile([]*component.ComponentDefinition{fromComp, toComp}, []*relationship.RelationshipDefinition{rel})
			result := p.IsInvalid(rel, design)
			if result != tt.expected {
				t.Errorf("IsInvalid(status=%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}
