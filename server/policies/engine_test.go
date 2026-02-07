package policies

import (
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
