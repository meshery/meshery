package policies

import (
	"fmt"
	"testing"

	"github.com/meshery/meshkit/logger"
)

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
	comp := map[string]interface{}{
		"component": map[string]interface{}{"kind": "Namespace"},
		"model":     map[string]interface{}{"name": "kubernetes"},
	}

	tests := []struct {
		name     string
		selector map[string]interface{}
		expected bool
	}{
		{
			"wildcard",
			map[string]interface{}{"kind": "*", "model": map[string]interface{}{"name": "*"}},
			true,
		},
		{
			"exact match",
			map[string]interface{}{"kind": "Namespace", "model": map[string]interface{}{"name": "kubernetes"}},
			true,
		},
		{
			"kind mismatch",
			map[string]interface{}{"kind": "Service", "model": map[string]interface{}{"name": "kubernetes"}},
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
	fromDecl := map[string]interface{}{
		"id":        "from-namespace",
		"component": map[string]interface{}{"kind": "Namespace"},
		"model": map[string]interface{}{
			"name":       "kubernetes",
			"registrant": map[string]interface{}{"kind": "github"},
		},
	}
	toDecl := map[string]interface{}{
		"id":        "to-namespace",
		"component": map[string]interface{}{"kind": "Namespace"},
		"model": map[string]interface{}{
			"name":       "kubernetes",
			"registrant": map[string]interface{}{"kind": "github"},
		},
	}

	denySelectors := map[string]interface{}{
		"from": []interface{}{
			map[string]interface{}{
				"kind": "Namespace",
				"model": map[string]interface{}{
					"name":       "kubernetes",
					"registrant": "*",
				},
			},
		},
		"to": []interface{}{
			map[string]interface{}{
				"kind": "Namespace",
				"model": map[string]interface{}{
					"name":       "kubernetes",
					"registrant": "*",
				},
			},
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
		rel      map[string]interface{}
		expected bool
	}{
		{
			"alias",
			map[string]interface{}{"kind": "hierarchical", "type": "parent", "subType": "alias"},
			true,
		},
		{
			"alias uppercase",
			map[string]interface{}{"kind": "Hierarchical", "type": "Parent", "subType": "Alias"},
			true,
		},
		{
			"inventory",
			map[string]interface{}{"kind": "hierarchical", "type": "parent", "subType": "inventory"},
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
		rel            map[string]interface{}
		expectedPolicy string
	}{
		{
			"hierarchical inventory",
			map[string]interface{}{"kind": "hierarchical", "type": "parent", "subType": "inventory"},
			"hierarchical_parent_child",
		},
		{
			"edge non-binding",
			map[string]interface{}{"kind": "edge", "type": "non-binding"},
			"edge-non-binding",
		},
		{
			"alias",
			map[string]interface{}{"kind": "hierarchical", "type": "parent", "subType": "alias"},
			"alias_relationships_policy",
		},
		{
			"sibling",
			map[string]interface{}{"type": "sibling"},
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
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{"id": "comp-a"},
			map[string]interface{}{"id": "comp-b"},
		},
	}

	rel := map[string]interface{}{
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{map[string]interface{}{"id": "comp-a"}},
					"to":   []interface{}{map[string]interface{}{"id": "comp-b"}},
				},
			},
		},
	}

	if !fromAndToComponentsExist(rel, design) {
		t.Error("Expected both components to exist")
	}

	relMissing := map[string]interface{}{
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{map[string]interface{}{"id": "comp-a"}},
					"to":   []interface{}{map[string]interface{}{"id": "comp-missing"}},
				},
			},
		},
	}

	if fromAndToComponentsExist(relMissing, design) {
		t.Error("Expected missing component to be detected")
	}
}

func TestValidateRelationshipsInDesign(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{"id": "comp-a"},
		},
		"relationships": []interface{}{
			map[string]interface{}{
				"id":      "rel-1",
				"kind":    "edge",
				"type":    "non-binding",
				"status":  "approved",
				"subType": "",
				"selectors": []interface{}{
					map[string]interface{}{
						"allow": map[string]interface{}{
							"from": []interface{}{map[string]interface{}{"id": "comp-a"}},
							"to":   []interface{}{map[string]interface{}{"id": "comp-deleted"}},
						},
					},
				},
			},
		},
	}

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

	binding := map[string]interface{}{
		"kind":    "edge",
		"type":    "binding",
		"subType": "",
	}
	if !p.IsImplicatedBy(binding) {
		t.Error("Expected edge binding to be implicated")
	}

	nonBinding := map[string]interface{}{
		"kind":    "edge",
		"type":    "non-binding",
		"subType": "",
	}
	if p.IsImplicatedBy(nonBinding) {
		t.Error("Expected edge non-binding to not be implicated")
	}
}

func TestIsValidBinding(t *testing.T) {
	role := map[string]interface{}{
		"id":        "role-1",
		"component": map[string]interface{}{"kind": "Role"},
		"name":      "my-role",
	}
	roleBinding := map[string]interface{}{
		"id":        "rb-1",
		"component": map[string]interface{}{"kind": "RoleBinding"},
		"roleRef":   map[string]interface{}{"name": "my-role"},
	}
	selector := map[string]interface{}{
		"kind": "Role",
		"match": map[string]interface{}{
			"from": []interface{}{map[string]interface{}{
				"kind":       "self",
				"mutatorRef": []interface{}{[]interface{}{"name"}},
			}},
			"to": []interface{}{map[string]interface{}{
				"kind":       "RoleBinding",
				"mutatedRef": []interface{}{[]interface{}{"roleRef", "name"}},
			}},
		},
	}

	if !isValidBinding(role, roleBinding, selector) {
		t.Error("Expected valid binding")
	}

	// Mismatched value
	badBinding := map[string]interface{}{
		"id":        "rb-2",
		"component": map[string]interface{}{"kind": "RoleBinding"},
		"roleRef":   map[string]interface{}{"name": "other-role"},
	}
	if isValidBinding(role, badBinding, selector) {
		t.Error("Expected invalid binding")
	}
}

func TestHierarchicalIdentifyRelationship(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "default"},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"namespace": "default",
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def-1",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "Namespace",
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "name"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "Deployment",
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{
									[]interface{}{"configuration", "namespace"},
								},
							},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("Expected 1 identified relationship, got %d", len(identified))
	}
	if getMapString(identified[0], "status") != "approved" {
		t.Error("Expected status to be 'approved'")
	}
}

func TestHierarchicalSideEffects(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"configuration": map[string]interface{}{"name": "production"},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"configuration": map[string]interface{}{
					"namespace": "default",
				},
			},
		},
	}

	rel := map[string]interface{}{
		"id":      "rel-1",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"status":  "approved",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"id": "ns-1",
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "name"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"id": "deploy-1",
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{
									[]interface{}{"configuration", "namespace"},
								},
							},
						},
					},
				},
			},
		},
	}

	actions := p.SideEffects(rel, design)
	if len(actions) == 0 {
		t.Fatal("Expected side effect actions")
	}
	if actions[0].Value["value"] != "production" {
		t.Errorf("Expected namespace to be updated to 'production', got %v", actions[0].Value["value"])
	}
}

func TestIdentifyAdditions(t *testing.T) {
	// A deployment references namespace "default" but no Namespace component exists.
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"namespace": "default",
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def-1",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "name"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Deployment",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{
									[]interface{}{"configuration", "namespace"},
								},
							},
						},
					},
				},
			},
		},
	}

	actions := identifyAdditions(relDef, design)
	if len(actions) == 0 {
		t.Fatal("Expected addition actions for missing Namespace")
	}
	if actions[0].Op != AddComponentOp {
		t.Errorf("Expected add_component op, got %s", actions[0].Op)
	}
	item := getMapMap(actions[0].Value, "item")
	if item == nil {
		t.Fatal("Expected item in action value")
	}
	comp := getMapMap(item, "component")
	if getMapString(comp, "kind") != "Namespace" {
		t.Errorf("Expected Namespace kind, got %s", getMapString(comp, "kind"))
	}
}

func TestIdentifyAdditionsNoAction(t *testing.T) {
	// Both namespace and deployment exist, so no additions needed.
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "default"},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"namespace": "default",
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def-1",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "name"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Deployment",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{
									[]interface{}{"configuration", "namespace"},
								},
							},
						},
					},
				},
			},
		},
	}

	actions := identifyAdditions(relDef, design)
	if len(actions) != 0 {
		t.Errorf("Expected no additions, got %d", len(actions))
	}
}

func TestGoEngineCreation(t *testing.T) {
	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	if engine == nil {
		t.Fatal("Expected non-nil engine")
	}
	if len(engine.policies) != 5 {
		t.Errorf("Expected 5 policies, got %d", len(engine.policies))
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
	u2 := newUUID(seed)
	// UUIDs should differ because newUUID incorporates time
	if u1 == u2 {
		t.Error("newUUID should produce different UUIDs on successive calls (time-based)")
	}
}

// --- policy_alias.go tests ---

func TestAliasIsInvalidStatusMatrix(t *testing.T) {
	p := &AliasPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "from-comp",
				"component": map[string]interface{}{"kind": "Container"},
				"configuration": map[string]interface{}{
					"ports": []interface{}{map[string]interface{}{"containerPort": 80}},
				},
			},
			map[string]interface{}{
				"id":        "to-comp",
				"component": map[string]interface{}{"kind": "Deployment"},
				"configuration": map[string]interface{}{
					"ports": []interface{}{map[string]interface{}{"containerPort": 80}},
				},
			},
		},
	}

	makeRel := func(status string) map[string]interface{} {
		return map[string]interface{}{
			"kind":    "hierarchical",
			"type":    "parent",
			"subType": "alias",
			"status":  status,
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"id": "from-comp",
								"patch": map[string]interface{}{
									"mutatorRef": []interface{}{
										[]interface{}{"configuration", "ports"},
									},
								},
							},
						},
						"to": []interface{}{
							map[string]interface{}{
								"id": "to-comp",
							},
						},
					},
				},
			},
		}
	}

	tests := []struct {
		status   string
		expected bool // IsInvalid returns true when invalid
	}{
		{"pending", false},   // pending is valid
		{"approved", false},  // approved with existing components is valid
		{"deleted", true},    // deleted is invalid
		{"malformed", true},  // unknown status is invalid
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			rel := makeRel(tt.status)
			result := p.IsInvalid(rel, design)
			if result != tt.expected {
				t.Errorf("IsInvalid(status=%q) = %v, want %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestAliasRefFromRelationshipMissingSelectorParts(t *testing.T) {
	t.Run("no selectors", func(t *testing.T) {
		rel := map[string]interface{}{}
		ref := aliasRefFromRelationship(rel)
		if ref != nil {
			t.Errorf("Expected nil, got %v", ref)
		}
	})

	t.Run("empty selectors", func(t *testing.T) {
		rel := map[string]interface{}{
			"selectors": []interface{}{},
		}
		ref := aliasRefFromRelationship(rel)
		if ref != nil {
			t.Errorf("Expected nil, got %v", ref)
		}
	})

	t.Run("selector without allow", func(t *testing.T) {
		rel := map[string]interface{}{
			"selectors": []interface{}{
				map[string]interface{}{},
			},
		}
		ref := aliasRefFromRelationship(rel)
		if ref != nil {
			t.Errorf("Expected nil, got %v", ref)
		}
	})

	t.Run("allow without from", func(t *testing.T) {
		rel := map[string]interface{}{
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{},
				},
			},
		}
		ref := aliasRefFromRelationship(rel)
		if ref != nil {
			t.Errorf("Expected nil, got %v", ref)
		}
	})

	t.Run("from without patch", func(t *testing.T) {
		rel := map[string]interface{}{
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{"id": "comp-1"},
						},
					},
				},
			},
		}
		ref := aliasRefFromRelationship(rel)
		if ref != nil {
			t.Errorf("Expected nil, got %v", ref)
		}
	})

	t.Run("valid extraction", func(t *testing.T) {
		rel := map[string]interface{}{
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"id": "comp-1",
								"patch": map[string]interface{}{
									"mutatorRef": []interface{}{
										[]interface{}{"configuration", "containers", "0"},
									},
								},
							},
						},
					},
				},
			},
		}
		ref := aliasRefFromRelationship(rel)
		if len(ref) != 3 || ref[0] != "configuration" || ref[1] != "containers" || ref[2] != "0" {
			t.Errorf("Expected [configuration containers 0], got %v", ref)
		}
	})
}

func TestGetArrayAwareConfigPaths(t *testing.T) {
	design := map[string]interface{}{}

	t.Run("direct ref", func(t *testing.T) {
		comp := map[string]interface{}{
			"id":            "comp-1",
			"configuration": map[string]interface{}{"name": "test"},
		}
		ref := []string{"configuration", "name"}
		paths := getArrayAwareConfigPaths(ref, comp, design)
		if len(paths) != 1 {
			t.Fatalf("Expected 1 path, got %d", len(paths))
		}
		if len(paths[0]) != 2 || paths[0][0] != "configuration" || paths[0][1] != "name" {
			t.Errorf("Expected [configuration name], got %v", paths[0])
		}
	})

	t.Run("direct ref missing value", func(t *testing.T) {
		comp := map[string]interface{}{
			"id":            "comp-1",
			"configuration": map[string]interface{}{},
		}
		ref := []string{"configuration", "name"}
		paths := getArrayAwareConfigPaths(ref, comp, design)
		if paths != nil {
			t.Errorf("Expected nil for missing value, got %v", paths)
		}
	})

	t.Run("wildcard ref", func(t *testing.T) {
		comp := map[string]interface{}{
			"id": "comp-1",
			"configuration": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{"name": "nginx"},
					map[string]interface{}{"name": "sidecar"},
				},
			},
		}
		ref := []string{"configuration", "containers", "_"}
		paths := getArrayAwareConfigPaths(ref, comp, design)
		if len(paths) != 2 {
			t.Fatalf("Expected 2 paths, got %d", len(paths))
		}
	})

	t.Run("wildcard ref empty array", func(t *testing.T) {
		comp := map[string]interface{}{
			"id": "comp-1",
			"configuration": map[string]interface{}{
				"containers": []interface{}{},
			},
		}
		ref := []string{"configuration", "containers", "_"}
		paths := getArrayAwareConfigPaths(ref, comp, design)
		if paths != nil {
			t.Errorf("Expected nil for empty array, got %v", paths)
		}
	})
}

func TestAliasIdentifyRelationshipMultiplePaths(t *testing.T) {
	p := &AliasPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{"name": "nginx"},
						map[string]interface{}{"name": "sidecar"},
					},
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "alias-def",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "alias",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Container",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "containers", "_"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Deployment",
							"model": map[string]interface{}{"name": "kubernetes"},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 2 {
		t.Fatalf("Expected 2 identified relationships (one per container), got %d", len(identified))
	}
	for _, rel := range identified {
		if getMapString(rel, "status") != "identified" {
			t.Error("Expected status 'identified'")
		}
	}
}

func TestAliasAlreadyExistsDuplicateDetection(t *testing.T) {
	patch := map[string]interface{}{
		"patchStrategy": "replace",
		"mutatorRef":    []interface{}{[]interface{}{"configuration", "containers", "0"}},
	}

	existing := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "alias",
		"status":  "approved",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{map[string]interface{}{"id": "from-1"}},
					"to": []interface{}{map[string]interface{}{
						"kind":  "Deployment",
						"id":    "to-1",
						"patch": patch,
					}},
				},
			},
		},
	}

	newRel := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "alias",
		"status":  "identified",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{map[string]interface{}{"id": "from-new"}},
					"to": []interface{}{map[string]interface{}{
						"kind":  "Deployment",
						"id":    "to-1",
						"patch": patch,
					}},
				},
			},
		},
	}

	t.Run("duplicate detected", func(t *testing.T) {
		design := map[string]interface{}{
			"relationships": []interface{}{existing},
		}
		p := &AliasPolicy{}
		if !p.AlreadyExists(newRel, design) {
			t.Error("Expected duplicate to be detected")
		}
	})

	t.Run("deleted existing skipped", func(t *testing.T) {
		deletedExisting := deepCopyMap(existing)
		deletedExisting["status"] = "deleted"
		design := map[string]interface{}{
			"relationships": []interface{}{deletedExisting},
		}
		p := &AliasPolicy{}
		if p.AlreadyExists(newRel, design) {
			t.Error("Expected deleted relationship to be skipped")
		}
	})

	t.Run("no duplicate", func(t *testing.T) {
		design := map[string]interface{}{
			"relationships": []interface{}{},
		}
		p := &AliasPolicy{}
		if p.AlreadyExists(newRel, design) {
			t.Error("Expected no duplicate in empty design")
		}
	})
}

func TestAliasSideEffectsByStatus(t *testing.T) {
	p := &AliasPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{"id": "from-comp"},
			map[string]interface{}{"id": "to-comp"},
		},
	}

	rel := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "alias",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"id":   "from-comp",
							"kind": "Container",
							"model": map[string]interface{}{
								"name": "kubernetes",
							},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{
									[]interface{}{"configuration", "containers", "0"},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{"id": "to-comp"},
					},
				},
			},
		},
	}

	t.Run("identified triggers add", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "identified"
		actions := p.SideEffects(r, design)
		if len(actions) == 0 {
			t.Fatal("Expected add component actions for identified status")
		}
		if actions[0].Op != AddComponentOp {
			t.Errorf("Expected add_component, got %s", actions[0].Op)
		}
	})

	t.Run("pending triggers add", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "pending"
		actions := p.SideEffects(r, design)
		if len(actions) == 0 {
			t.Fatal("Expected add component actions for pending status")
		}
		if actions[0].Op != AddComponentOp {
			t.Errorf("Expected add_component, got %s", actions[0].Op)
		}
	})

	t.Run("deleted triggers delete", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "deleted"
		actions := p.SideEffects(r, design)
		if len(actions) == 0 {
			t.Fatal("Expected delete component actions for deleted status")
		}
		if actions[0].Op != DeleteComponentOp {
			t.Errorf("Expected delete_component, got %s", actions[0].Op)
		}
	})

	t.Run("approved returns nil", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "approved"
		actions := p.SideEffects(r, design)
		if actions != nil {
			t.Errorf("Expected nil actions for approved status, got %d", len(actions))
		}
	})
}

// --- policy_hierarchical.go tests ---

func TestHierarchicalDenySelectorRejection(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes", "registrant": map[string]interface{}{"kind": "github"}},
				"configuration": map[string]interface{}{"name": "default"},
			},
			map[string]interface{}{
				"id":            "ns-2",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes", "registrant": map[string]interface{}{"kind": "github"}},
				"configuration": map[string]interface{}{"namespace": "default"},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "name"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "namespace"}},
							},
						},
					},
				},
				"deny": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes", "registrant": "*"},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes", "registrant": "*"},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 0 {
		t.Errorf("Expected denied relationship to be rejected, got %d identified", len(identified))
	}
}

func TestHierarchicalFeasibilityRejection(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "default"},
			},
			map[string]interface{}{
				"id":            "svc-1",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "istio"}, // different model
				"configuration": map[string]interface{}{"namespace": "default"},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "name"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Deployment",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "namespace"}},
							},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 0 {
		t.Errorf("Expected infeasible relationship to be rejected, got %d identified", len(identified))
	}
}

func TestHierarchicalApprovedStatusOnValidPair(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "prod"},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"namespace": "prod",
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "rel-def",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Namespace",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "name"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Deployment",
							"model": map[string]interface{}{"name": "kubernetes"},
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "namespace"}},
							},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("Expected 1 identified relationship, got %d", len(identified))
	}
	if getMapString(identified[0], "status") != "approved" {
		t.Errorf("Expected status 'approved', got %q", getMapString(identified[0], "status"))
	}
}

func TestHierarchicalSideEffectsSkippedWhenDeleted(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"configuration": map[string]interface{}{"name": "production"},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"configuration": map[string]interface{}{
					"namespace": "default",
				},
			},
		},
	}

	rel := map[string]interface{}{
		"id":      "rel-1",
		"kind":    "hierarchical",
		"type":    "parent",
		"subType": "inventory",
		"status":  "deleted",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"id": "ns-1",
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "name"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"id": "deploy-1",
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "namespace"}},
							},
						},
					},
				},
			},
		},
	}

	actions := p.SideEffects(rel, design)
	if len(actions) != 0 {
		t.Errorf("Expected no side effects for deleted status, got %d actions", len(actions))
	}
}

// --- policy_binding.go tests ---

func TestBindingTripleMatching(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "role-1",
				"component": map[string]interface{}{"kind": "ClusterRole"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"name":      "admin-role",
			},
			map[string]interface{}{
				"id":        "rb-1",
				"component": map[string]interface{}{"kind": "ClusterRoleBinding"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"roleRef":   map[string]interface{}{"name": "admin-role"},
				"subjects":  []interface{}{map[string]interface{}{"name": "my-sa"}},
			},
			map[string]interface{}{
				"id":        "sa-1",
				"component": map[string]interface{}{"kind": "ServiceAccount"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"name":      "my-sa",
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "binding-def",
		"kind":    "edge",
		"type":    "binding",
		"subType": "permission",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "ClusterRole",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatorRef": []interface{}{[]interface{}{"name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatedRef": []interface{}{[]interface{}{"roleRef", "name"}},
									},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "ServiceAccount",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatorRef": []interface{}{[]interface{}{"subjects", "0", "name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatedRef": []interface{}{[]interface{}{"name"}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	p := &EdgeBindingPolicy{}
	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("Expected 1 binding relationship, got %d", len(identified))
	}
	if getMapString(identified[0], "status") != "approved" {
		t.Error("Expected status 'approved'")
	}
}

func TestBindingMismatchRejection(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "role-1",
				"component": map[string]interface{}{"kind": "ClusterRole"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"name":      "admin-role",
			},
			map[string]interface{}{
				"id":        "rb-1",
				"component": map[string]interface{}{"kind": "ClusterRoleBinding"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"roleRef":   map[string]interface{}{"name": "different-role"}, // mismatch
				"subjects":  []interface{}{map[string]interface{}{"name": "my-sa"}},
			},
			map[string]interface{}{
				"id":        "sa-1",
				"component": map[string]interface{}{"kind": "ServiceAccount"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"name":      "my-sa",
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "binding-def",
		"kind":    "edge",
		"type":    "binding",
		"subType": "permission",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "ClusterRole",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatorRef": []interface{}{[]interface{}{"name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatedRef": []interface{}{[]interface{}{"roleRef", "name"}},
									},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "ServiceAccount",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatorRef": []interface{}{[]interface{}{"subjects", "0", "name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatedRef": []interface{}{[]interface{}{"name"}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	p := &EdgeBindingPolicy{}
	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 0 {
		t.Errorf("Expected 0 binding relationships for mismatched role, got %d", len(identified))
	}
}

func TestBindingDenySelectorEnforcement(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "role-1",
				"component": map[string]interface{}{"kind": "ClusterRole"},
				"model":     map[string]interface{}{"name": "kubernetes", "registrant": map[string]interface{}{"kind": "github"}},
				"name":      "admin-role",
			},
			map[string]interface{}{
				"id":        "rb-1",
				"component": map[string]interface{}{"kind": "ClusterRoleBinding"},
				"model":     map[string]interface{}{"name": "kubernetes", "registrant": map[string]interface{}{"kind": "github"}},
				"roleRef":   map[string]interface{}{"name": "admin-role"},
				"subjects":  []interface{}{map[string]interface{}{"name": "my-sa"}},
			},
			map[string]interface{}{
				"id":        "sa-1",
				"component": map[string]interface{}{"kind": "ServiceAccount"},
				"model":     map[string]interface{}{"name": "kubernetes", "registrant": map[string]interface{}{"kind": "github"}},
				"name":      "my-sa",
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "binding-def",
		"kind":    "edge",
		"type":    "binding",
		"subType": "permission",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "ClusterRole",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatorRef": []interface{}{[]interface{}{"name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatedRef": []interface{}{[]interface{}{"roleRef", "name"}},
									},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "ServiceAccount",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRoleBinding",
										"mutatorRef": []interface{}{[]interface{}{"subjects", "0", "name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatedRef": []interface{}{[]interface{}{"name"}},
									},
								},
							},
						},
					},
				},
				"deny": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "ClusterRole",
							"model": map[string]interface{}{"name": "kubernetes", "registrant": "*"},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "ServiceAccount",
							"model": map[string]interface{}{"name": "kubernetes", "registrant": "*"},
						},
					},
				},
			},
		},
	}

	p := &EdgeBindingPolicy{}
	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 0 {
		t.Errorf("Expected denied binding to produce 0 results, got %d", len(identified))
	}
}

func TestBindingSelfReferenceExclusion(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "role-1",
				"component": map[string]interface{}{"kind": "ClusterRole"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"name":      "admin-role",
				"roleRef":   map[string]interface{}{"name": "admin-role"},
				"subjects":  []interface{}{map[string]interface{}{"name": "admin-role"}},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "binding-def",
		"kind":    "edge",
		"type":    "binding",
		"subType": "permission",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "ClusterRole",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatorRef": []interface{}{[]interface{}{"name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRole",
										"mutatedRef": []interface{}{[]interface{}{"roleRef", "name"}},
									},
								},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "ClusterRole",
							"match": map[string]interface{}{
								"from": []interface{}{
									map[string]interface{}{
										"kind":       "ClusterRole",
										"mutatorRef": []interface{}{[]interface{}{"subjects", "0", "name"}},
									},
								},
								"to": []interface{}{
									map[string]interface{}{
										"kind":       "self",
										"mutatedRef": []interface{}{[]interface{}{"name"}},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	p := &EdgeBindingPolicy{}
	identified := p.IdentifyRelationship(relDef, design)
	// With only one component, from==binding or binding==to should be excluded
	if len(identified) != 0 {
		t.Errorf("Expected self-reference to be excluded, got %d", len(identified))
	}
}

// --- policy_edge_network.go tests ---

func TestEdgeNonBindingImplication(t *testing.T) {
	p := &EdgeNonBindingPolicy{}

	t.Run("positive", func(t *testing.T) {
		rel := map[string]interface{}{"kind": "edge", "type": "non-binding"}
		if !p.IsImplicatedBy(rel) {
			t.Error("Expected edge non-binding to be implicated")
		}
	})

	t.Run("negative binding", func(t *testing.T) {
		rel := map[string]interface{}{"kind": "edge", "type": "binding"}
		if p.IsImplicatedBy(rel) {
			t.Error("Expected edge binding to not be implicated by non-binding policy")
		}
	})

	t.Run("negative hierarchical", func(t *testing.T) {
		rel := map[string]interface{}{"kind": "hierarchical", "type": "parent"}
		if p.IsImplicatedBy(rel) {
			t.Error("Expected hierarchical to not be implicated by non-binding policy")
		}
	})
}

func TestEdgeNonBindingIdentifyPositive(t *testing.T) {
	p := &EdgeNonBindingPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":        "svc-1",
				"component": map[string]interface{}{"kind": "Service"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"selector": map[string]interface{}{"app": "web"},
				},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"selector": map[string]interface{}{
						"matchLabels": map[string]interface{}{"app": "web"},
					},
				},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "net-def",
		"kind":    "edge",
		"type":    "non-binding",
		"subType": "network",
		"model":   map[string]interface{}{"name": "kubernetes"},
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind": "Service",
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "selector"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind": "Deployment",
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "selector", "matchLabels"}},
							},
						},
					},
				},
			},
		},
	}

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("Expected 1 identified network relationship, got %d", len(identified))
	}
}

func TestEdgeNonBindingSideEffectsDeletedNoop(t *testing.T) {
	p := &EdgeNonBindingPolicy{}

	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "svc-1",
				"component":     map[string]interface{}{"kind": "Service"},
				"configuration": map[string]interface{}{"selector": map[string]interface{}{"app": "web"}},
			},
			map[string]interface{}{
				"id":        "deploy-1",
				"component": map[string]interface{}{"kind": "Deployment"},
				"configuration": map[string]interface{}{
					"selector": map[string]interface{}{"matchLabels": map[string]interface{}{"app": "old"}},
				},
			},
		},
	}

	rel := map[string]interface{}{
		"id":      "rel-1",
		"kind":    "edge",
		"type":    "non-binding",
		"subType": "network",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"id": "svc-1",
							"patch": map[string]interface{}{
								"mutatorRef": []interface{}{[]interface{}{"configuration", "selector"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"id": "deploy-1",
							"patch": map[string]interface{}{
								"mutatedRef": []interface{}{[]interface{}{"configuration", "selector", "matchLabels"}},
							},
						},
					},
				},
			},
		},
	}

	t.Run("deleted returns nil", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "deleted"
		actions := p.SideEffects(r, design)
		if actions != nil {
			t.Errorf("Expected nil side effects for deleted, got %d", len(actions))
		}
	})

	t.Run("approved produces patch", func(t *testing.T) {
		r := deepCopyMap(rel)
		r["status"] = "approved"
		actions := p.SideEffects(r, design)
		if len(actions) == 0 {
			t.Fatal("Expected patch actions for approved status with differing values")
		}
	})
}

// --- policy_matchlabels.go tests ---

func TestMatchLabelsGroupingLogic(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "svc-1",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}},
			},
			map[string]interface{}{
				"id":            "svc-2",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}},
			},
			map[string]interface{}{
				"id":            "svc-3",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "api"}},
			},
		},
	}

	relDef := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "sibling",
		"subType": "matchlabels",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
							"match": map[string]interface{}{
								"refs": []interface{}{[]interface{}{"configuration", "labels"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
						},
					},
				},
			},
		},
	}

	groups := identifyMatchlabels(design, relDef)
	// svc-1 and svc-2 share app=web, svc-3 has app=api so no group for it
	found := false
	for _, g := range groups {
		if g.Field == "app" && g.Value == "web" {
			found = true
			if len(g.Components) != 2 {
				t.Errorf("Expected 2 components in web group, got %d", len(g.Components))
			}
		}
	}
	if !found {
		t.Error("Expected to find group for app=web")
	}
}

func TestMatchLabelsDedupInGroup(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "svc-1",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web", "env": "prod"}},
			},
			map[string]interface{}{
				"id":            "svc-2",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web", "env": "prod"}},
			},
		},
	}

	relDef := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "sibling",
		"subType": "matchlabels",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
							"match": map[string]interface{}{
								"refs": []interface{}{[]interface{}{"configuration", "labels"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
						},
					},
				},
			},
		},
	}

	groups := identifyMatchlabels(design, relDef)
	for _, g := range groups {
		ids := make(map[string]bool)
		for _, comp := range g.Components {
			id := getMapString(comp, "id")
			if ids[id] {
				t.Errorf("Duplicate component %q in group %s=%v", id, g.Field, g.Value)
			}
			ids[id] = true
		}
	}
}

func TestMatchLabelsMaxCap(t *testing.T) {
	// Create more than maxMatchLabels unique label groups
	var comps []interface{}
	for i := 0; i < maxMatchLabels+5; i++ {
		label := fmt.Sprintf("label-%d", i)
		comps = append(comps, map[string]interface{}{
			"id":            fmt.Sprintf("svc-%d-a", i),
			"component":     map[string]interface{}{"kind": "Service"},
			"model":         map[string]interface{}{"name": "kubernetes"},
			"configuration": map[string]interface{}{"labels": map[string]interface{}{label: "val"}},
		})
		comps = append(comps, map[string]interface{}{
			"id":            fmt.Sprintf("svc-%d-b", i),
			"component":     map[string]interface{}{"kind": "Service"},
			"model":         map[string]interface{}{"name": "kubernetes"},
			"configuration": map[string]interface{}{"labels": map[string]interface{}{label: "val"}},
		})
	}

	design := map[string]interface{}{
		"components": comps,
	}

	relDef := map[string]interface{}{
		"kind":    "hierarchical",
		"type":    "sibling",
		"subType": "matchlabels",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
							"match": map[string]interface{}{
								"refs": []interface{}{[]interface{}{"configuration", "labels"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
						},
					},
				},
			},
		},
	}

	groups := identifyMatchlabels(design, relDef)
	if len(groups) > maxMatchLabels {
		t.Errorf("Expected at most %d groups, got %d", maxMatchLabels, len(groups))
	}
}

func TestMatchLabelsDeterministicID(t *testing.T) {
	design := map[string]interface{}{
		"components": []interface{}{
			map[string]interface{}{
				"id":            "svc-1",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}},
			},
			map[string]interface{}{
				"id":            "svc-2",
				"component":     map[string]interface{}{"kind": "Service"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}},
			},
		},
		"relationships": []interface{}{},
	}

	relDef := map[string]interface{}{
		"id":      "ml-def",
		"kind":    "hierarchical",
		"type":    "sibling",
		"subType": "matchlabels",
		"selectors": []interface{}{
			map[string]interface{}{
				"allow": map[string]interface{}{
					"from": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
							"match": map[string]interface{}{
								"refs": []interface{}{[]interface{}{"configuration", "labels"}},
							},
						},
					},
					"to": []interface{}{
						map[string]interface{}{
							"kind":  "Service",
							"model": map[string]interface{}{"name": "kubernetes"},
						},
					},
				},
			},
		},
	}

	p := &MatchLabelsPolicy{}
	rels1 := p.IdentifyRelationship(relDef, design)
	rels2 := p.IdentifyRelationship(relDef, design)

	if len(rels1) != len(rels2) {
		t.Fatalf("Different number of identified rels: %d vs %d", len(rels1), len(rels2))
	}

	ids1 := make(map[string]bool)
	for _, r := range rels1 {
		ids1[getMapString(r, "id")] = true
	}
	for _, r := range rels2 {
		id := getMapString(r, "id")
		if !ids1[id] {
			t.Errorf("ID %q from second run not found in first run", id)
		}
	}
}

func TestMatchLabelsMalformedInput(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		design := map[string]interface{}{
			"components": []interface{}{
				map[string]interface{}{
					"id":        "svc-1",
					"component": map[string]interface{}{"kind": "Service"},
					"model":     map[string]interface{}{"name": "kubernetes"},
				},
				map[string]interface{}{
					"id":        "svc-2",
					"component": map[string]interface{}{"kind": "Service"},
					"model":     map[string]interface{}{"name": "kubernetes"},
				},
			},
		}

		relDef := map[string]interface{}{
			"kind":    "hierarchical",
			"type":    "sibling",
			"subType": "matchlabels",
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"kind":  "Service",
								"model": map[string]interface{}{"name": "kubernetes"},
								"match": map[string]interface{}{
									"refs": []interface{}{[]interface{}{"configuration", "labels"}},
								},
							},
						},
						"to": []interface{}{
							map[string]interface{}{
								"kind":  "Service",
								"model": map[string]interface{}{"name": "kubernetes"},
							},
						},
					},
				},
			},
		}

		groups := identifyMatchlabels(design, relDef)
		if len(groups) != 0 {
			t.Errorf("Expected 0 groups for nil config, got %d", len(groups))
		}
	})

	t.Run("empty match refs", func(t *testing.T) {
		design := map[string]interface{}{
			"components": []interface{}{
				map[string]interface{}{
					"id":            "svc-1",
					"component":     map[string]interface{}{"kind": "Service"},
					"model":         map[string]interface{}{"name": "kubernetes"},
					"configuration": map[string]interface{}{"labels": map[string]interface{}{"app": "web"}},
				},
			},
		}

		relDef := map[string]interface{}{
			"kind":    "hierarchical",
			"type":    "sibling",
			"subType": "matchlabels",
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"kind":  "Service",
								"model": map[string]interface{}{"name": "kubernetes"},
								"match": map[string]interface{}{
									"refs": []interface{}{},
								},
							},
						},
						"to": []interface{}{
							map[string]interface{}{
								"kind":  "Service",
								"model": map[string]interface{}{"name": "kubernetes"},
							},
						},
					},
				},
			},
		}

		groups := identifyMatchlabels(design, relDef)
		if len(groups) != 0 {
			t.Errorf("Expected 0 groups for empty refs, got %d", len(groups))
		}
	})
}
