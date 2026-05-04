package policies

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/relationship"
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

func TestIsRelationshipFeasible(t *testing.T) {
	comp := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
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

	design := makePatternFile([]*component.ComponentDefinition{compA}, []*relationship.RelationshipDefinition{rel})

	policy := &EdgeNonBindingPolicy{}
	actions := validateRelationshipsInDesign(design, policy)

	if len(actions) != 1 {
		t.Fatalf("Expected 1 validation action, got %d", len(actions))
	}
	if actions[0].Op != UpdateRelationshipOp {
		t.Error("Expected update_relationship action")
	}
	if actions[0].StringValue != "deleted" {
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

	role := &component.ComponentDefinition{
		Component:     component.Component{Kind: "Role"},
		Configuration: map[string]interface{}{"name": "my-role"},
	}
	role.ID = roleID

	roleBinding := &component.ComponentDefinition{
		Component:     component.Component{Kind: "RoleBinding"},
		Configuration: map[string]interface{}{"roleRef": map[string]interface{}{"name": "my-role"}},
	}
	roleBinding.ID = rbID

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
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("Deployment"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
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
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID: &deployID,
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
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
	if actions[0].UpdateValue != "production" {
		t.Errorf("Expected namespace to be updated to 'production', got %v", actions[0].UpdateValue)
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
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
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
		expected bool
	}{
		{"pending", false},
		{"approved", false},
		{"deleted", true},
		{"malformed", true},
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
