package policies

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

func newTestComponent(name, kind string) *component.ComponentDefinition {
	id := uuid.Must(uuid.NewV4())
	return &component.ComponentDefinition{
		ID:          id,
		DisplayName: name,
		Component: component.Component{
			Kind: kind,
		},
		ModelReference: modelv1beta1.ModelReference{
			Name: "kubernetes",
		},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

// TestFromAndToComponentsExist_MultiComponentRelationship verifies that
// fromAndToComponentsExist checks all referenced components in a multi-component
// relationship (such as sibling matchLabels with 3+ components), not just index 0.
func TestFromAndToComponentsExist_MultiComponentRelationship(t *testing.T) {
	compA := newTestComponent("pod-a", "Pod")
	compB := newTestComponent("pod-b", "Pod")
	compC := newTestComponent("pod-c", "Pod")

	kind := "Pod"
	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
					To: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
				},
			},
		},
	}

	t.Run("all components exist", func(t *testing.T) {
		design := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{compA, compB, compC},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}
		if !fromAndToComponentsExist(rel, design) {
			t.Fatal("expected fromAndToComponentsExist to be true when all components exist")
		}
		if fromOrToComponentsDontExist(rel, design) {
			t.Fatal("expected fromOrToComponentsDontExist to be false when all components exist")
		}
	})

	t.Run("secondary component compC deleted", func(t *testing.T) {
		// compC is removed from design
		design := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{compA, compB},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}
		if fromAndToComponentsExist(rel, design) {
			t.Fatal("expected fromAndToComponentsExist to be false when secondary component compC is missing")
		}
		if !fromOrToComponentsDontExist(rel, design) {
			t.Fatal("expected fromOrToComponentsDontExist to be true when secondary component compC is missing")
		}
	})

	t.Run("secondary component compB deleted", func(t *testing.T) {
		// compB is removed from design
		design := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{compA, compC},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}
		if fromAndToComponentsExist(rel, design) {
			t.Fatal("expected fromAndToComponentsExist to be false when secondary component compB is missing")
		}
		if !fromOrToComponentsDontExist(rel, design) {
			t.Fatal("expected fromOrToComponentsDontExist to be true when secondary component compB is missing")
		}
	})
}

// TestValidateRelationshipsInDesign_OrphanedMultiComponentRelationship verifies
// that validateRelationshipsInDesign marks a multi-component relationship as deleted
// when any of its constituent components are removed from the design.
func TestValidateRelationshipsInDesign_OrphanedMultiComponentRelationship(t *testing.T) {
	compA := newTestComponent("deploy-a", "Deployment")
	compB := newTestComponent("deploy-b", "Deployment")
	compC := newTestComponent("deploy-c", "Deployment")

	kind := "Deployment"
	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
					To: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
				},
			},
		},
	}

	// compC was deleted from the canvas/design
	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{compA, compB},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	policy := &MatchLabelsPolicy{}
	actions := validateRelationshipsInDesign(design, policy)

	if len(actions) != 1 {
		t.Fatalf("expected 1 validation action marking orphaned relationship as deleted, got %d", len(actions))
	}
	if actions[0].Op != UpdateRelationshipOp {
		t.Fatalf("expected op %s, got %s", UpdateRelationshipOp, actions[0].Op)
	}
	if actions[0].StringValue != StatusDeleted {
		t.Fatalf("expected status %s, got %v", StatusDeleted, actions[0].StringValue)
	}
}

// TestFromAndToComponentsExist_MultiSelectorSets verifies that fromAndToComponentsExist
// inspects all selector sets in *rel.Selectors, not only the first selector set (*rel.Selectors)[0].
func TestFromAndToComponentsExist_MultiSelectorSets(t *testing.T) {
	compA := newTestComponent("comp-a", "Deployment")
	compB := newTestComponent("comp-b", "Service")
	compX := newTestComponent("comp-x", "ConfigMap")
	compY := newTestComponent("comp-y", "Secret")

	kindDep := "Deployment"
	kindSvc := "Service"
	kindCM := "ConfigMap"
	kindSec := "Secret"

	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "binding",
		SubType:          "permission",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{ID: &compA.ID, Kind: &kindDep}},
					To:   []relationship.SelectorItem{{ID: &compB.ID, Kind: &kindSvc}},
				},
			},
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{ID: &compX.ID, Kind: &kindCM}},
					To:   []relationship.SelectorItem{{ID: &compY.ID, Kind: &kindSec}},
				},
			},
		},
	}

	t.Run("all selector sets exist", func(t *testing.T) {
		design := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{compA, compB, compX, compY},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}
		if !fromAndToComponentsExist(rel, design) {
			t.Fatal("expected fromAndToComponentsExist to be true when all selector set components exist")
		}
	})

	t.Run("second selector set component missing", func(t *testing.T) {
		// compY is removed from design
		design := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{compA, compB, compX},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}
		if fromAndToComponentsExist(rel, design) {
			t.Fatal("expected fromAndToComponentsExist to be false when second selector set component compY is missing")
		}
	})
}

// TestRelationshipsAreSame_MultiComponentDistinction verifies that RelationshipsAreSame
// does not falsely treat two distinct relationships as identical when they share index 0
// but differ on subsequent selector items.
func TestRelationshipsAreSame_MultiComponentDistinction(t *testing.T) {
	compA := newTestComponent("comp-a", "Pod")
	compB := newTestComponent("comp-b", "Pod")
	compC := newTestComponent("comp-c", "Pod")

	kind := "Pod"

	// rel1 connects [compA, compB] -> [compA, compB]
	rel1 := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
					},
					To: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compB.ID, Kind: &kind},
					},
				},
			},
		},
	}

	// rel2 connects [compA, compC] -> [compA, compC].
	// Both rel1 and rel2 have equal slice lengths (2 From items and 2 To items) and share
	// the primary component compA at index 0. On master, RelationshipsAreSame only checked
	// index 0 and falsely evaluated them as duplicates.
	// Equal lengths ensure this test exercises multiset content comparison (compB vs compC),
	// rather than passing merely on slice length differences.
	rel2 := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
					To: []relationship.SelectorItem{
						{ID: &compA.ID, Kind: &kind},
						{ID: &compC.ID, Kind: &kind},
					},
				},
			},
		},
	}

	if RelationshipsAreSame(rel1, rel2) {
		t.Fatal("RelationshipsAreSame falsely reported rel1 and rel2 as identical because they share From[0] and To[0]")
	}

	// rel3 connects [compB, compA] -> [compB, compA] (same components as rel1, permuted order).
	// Verifies multiset order independence: same items in different order are recognized as identical.
	rel3 := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &compB.ID, Kind: &kind},
						{ID: &compA.ID, Kind: &kind},
					},
					To: []relationship.SelectorItem{
						{ID: &compB.ID, Kind: &kind},
						{ID: &compA.ID, Kind: &kind},
					},
				},
			},
		},
	}

	if !RelationshipsAreSame(rel1, rel3) {
		t.Fatal("RelationshipsAreSame failed to recognize rel1 and rel3 as identical under permuted selector item order")
	}

	// Selector-less relationships must return false to retain master's deduplication behavior.
	relNil1 := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
	}
	relNil2 := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "sibling",
		SubType:          "matchlabels",
	}
	if RelationshipsAreSame(relNil1, relNil2) {
		t.Fatal("RelationshipsAreSame should return false for selector-less relationships to preserve master dedup semantics")
	}
}

// TestPatchMutatorsAction_MultiItemPatching verifies that patchMutatorsAction
// processes mutators across all paired selector items in a selector set, not only index 0.
func TestPatchMutatorsAction_MultiItemPatching(t *testing.T) {
	compA := newTestComponent("source-a", "Deployment")
	compA.Configuration = map[string]interface{}{"port": float64(8080)}
	compB := newTestComponent("target-b", "Service")
	compB.Configuration = map[string]interface{}{"port": float64(80)}

	compC := newTestComponent("source-c", "Deployment")
	compC.Configuration = map[string]interface{}{"host": "backend.local"}
	compD := newTestComponent("target-d", "Service")
	compD.Configuration = map[string]interface{}{"host": "default.local"}

	mutatorRefA := relationship.MutatorRef{[]string{"configuration", "port"}}
	mutatedRefB := relationship.MutatedRef{[]string{"configuration", "port"}}
	mutatorRefC := relationship.MutatorRef{[]string{"configuration", "host"}}
	mutatedRefD := relationship.MutatedRef{[]string{"configuration", "host"}}

	kindDep := "Deployment"
	kindSvc := "Service"

	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "non-binding",
		SubType:          "reference",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							ID:   &compA.ID,
							Kind: &kindDep,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRefA,
							},
						},
						{
							ID:   &compC.ID,
							Kind: &kindDep,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRefC,
							},
						},
					},
					To: []relationship.SelectorItem{
						{
							ID:   &compB.ID,
							Kind: &kindSvc,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatedRef: &mutatedRefB,
							},
						},
						{
							ID:   &compD.ID,
							Kind: &kindSvc,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatedRef: &mutatedRefD,
							},
						},
					},
				},
			},
		},
	}

	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{compA, compB, compC, compD},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	actions := patchMutatorsAction(rel, design)
	if len(actions) != 2 {
		t.Fatalf("expected 2 mutator patch actions across both paired selector items, got %d", len(actions))
	}

	// 1-to-N fan-out: 1 source component mutates multiple target components
	relFanOut := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "non-binding",
		SubType:          "reference",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							ID:   &compA.ID,
							Kind: &kindDep,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRefA,
							},
						},
					},
					To: []relationship.SelectorItem{
						{
							ID:   &compB.ID,
							Kind: &kindSvc,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatedRef: &mutatedRefB,
							},
						},
						{
							ID:   &compD.ID,
							Kind: &kindSvc,
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatedRef: &mutatedRefB,
							},
						},
					},
				},
			},
		},
	}
	actionsFanOut := patchMutatorsAction(relFanOut, design)
	if len(actionsFanOut) != 2 {
		t.Fatalf("expected 2 mutator patch actions in 1-to-N fan-out, got %d", len(actionsFanOut))
	}
}

func TestBindingMatchSelectors_MultiItemMatchResolutionAndIDStamping(t *testing.T) {
	fromID := "11111111-1111-1111-1111-111111111111"
	toID := "22222222-2222-2222-2222-222222222222"

	// MatchSelector with multiple MatchSelectorItems in From and To
	mutatorRef := [][]string{{"name"}}
	mutatedRef := [][]string{{"serviceName"}}

	sel := relationship.SelectorItem{
		Match: &relationship.MatchSelector{
			From: &[]relationship.MatchSelectorItem{
				{Kind: "Service"}, // index 0 has no mutatorRef
				{Kind: "Service", MutatorRef: &mutatorRef}, // index 1 has mutatorRef
			},
			To: &[]relationship.MatchSelectorItem{
				{Kind: "Pod"}, // index 0 has no mutatedRef
				{Kind: "Pod", MutatedRef: &mutatedRef}, // index 1 has mutatedRef
			},
		},
	}

	// 1. Verify patchBindingSelectorIDsTyped stamps all items
	patchBindingSelectorIDsTyped(&sel, fromID, toID)
	for i, item := range *sel.Match.From {
		if item.ID == nil || item.ID.String() != fromID {
			t.Fatalf("expected Match.From[%d].ID to be %s, got %v", i, fromID, item.ID)
		}
	}
	for i, item := range *sel.Match.To {
		if item.ID == nil || item.ID.String() != toID {
			t.Fatalf("expected Match.To[%d].ID to be %s, got %v", i, toID, item.ID)
		}
	}

	// 2. Verify extractMutatorFromMatch finds mutator even at index 1
	comp1 := &component.ComponentDefinition{ID: uuid.FromStringOrNil(fromID)}
	comp2 := &component.ComponentDefinition{ID: uuid.FromStringOrNil(toID)}
	mComp, mPaths := extractMutatorFromMatch(sel.Match, comp1, comp2)
	if mComp != comp1 || len(mPaths) != 1 || mPaths[0][0] != "name" {
		t.Fatalf("extractMutatorFromMatch failed to resolve mutator at index 1: comp=%v, paths=%v", mComp, mPaths)
	}

	// 3. Verify extractMutatedFromMatch finds mutated even at index 1
	mComp2, mPaths2 := extractMutatedFromMatch(sel.Match, comp1, comp2)
	if mComp2 != comp2 || len(mPaths2) != 1 || mPaths2[0][0] != "serviceName" {
		t.Fatalf("extractMutatedFromMatch failed to resolve mutated at index 1: comp=%v, paths=%v", mComp2, mPaths2)
	}
}

func TestBindingMatchFields_MultiSelectorItemPatching(t *testing.T) {
	fromUUID := uuid.Must(uuid.NewV4())
	toUUID1 := uuid.Must(uuid.NewV4())
	toUUID2 := uuid.Must(uuid.NewV4())

	fromComp := &component.ComponentDefinition{
		ID:            fromUUID,
		Component:     component.Component{Kind: "Service"},
		Configuration: map[string]interface{}{"spec": map[string]interface{}{"serviceName": "api-service"}},
	}
	toComp1 := &component.ComponentDefinition{
		ID:            toUUID1,
		Component:     component.Component{Kind: "Deployment"},
		Configuration: map[string]interface{}{"spec": map[string]interface{}{"targetService": "placeholder"}},
	}
	toComp2 := &component.ComponentDefinition{
		ID:            toUUID2,
		Component:     component.Component{Kind: "Deployment"},
		Configuration: map[string]interface{}{"spec": map[string]interface{}{"targetService": "placeholder"}},
	}

	mutatorRef := [][]string{{"configuration", "spec", "serviceName"}}
	mutatedRef := [][]string{{"configuration", "spec", "targetService"}}

	match := &relationship.MatchSelector{
		From: &[]relationship.MatchSelectorItem{
			{ID: &fromUUID, MutatorRef: &mutatorRef},
		},
		To: &[]relationship.MatchSelectorItem{
			{ID: &toUUID1, MutatedRef: &mutatedRef},
			{ID: &toUUID2, MutatedRef: &mutatedRef},
		},
	}

	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "edge",
		SubType:          "binding",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &fromUUID, Match: match},
					},
					To: []relationship.SelectorItem{
						{ID: &toUUID1},
						{ID: &toUUID2},
					},
				},
			},
		},
	}

	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{fromComp, toComp1, toComp2},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	actions := patchBindingMatchFields(rel, design)
	if len(actions) != 2 {
		t.Fatalf("expected 2 patch actions (one for each mutated deployment), got %d", len(actions))
	}

	targets := make(map[string]bool)
	for _, act := range actions {
		targets[act.ID] = true
	}
	if !targets[toUUID1.String()] || !targets[toUUID2.String()] {
		t.Fatalf("expected patch actions targeting both deployments (%s, %s), got targets: %v", toUUID1, toUUID2, targets)
	}
}

func TestAliasPolicy_MultiPathResolution(t *testing.T) {
	podUUID := uuid.Must(uuid.NewV4())
	pod := &component.ComponentDefinition{
		ID:        podUUID,
		Component: component.Component{Kind: "Pod"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{"name": "app-container"},
				},
				"initContainers": []interface{}{
					map[string]interface{}{"name": "init-container"},
				},
			},
		},
	}

	mutatorRef := [][]string{
		{"configuration", "spec", "containers", "_"},
		{"configuration", "spec", "initContainers", "_"},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "hierarchical",
		SubType:          "alias",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							Kind: &[]string{"Container"}[0],
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRef,
							},
						},
					},
					To: []relationship.SelectorItem{
						{
							Kind: &[]string{"Pod"}[0],
						},
					},
				},
			},
		},
	}

	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{pod},
		Relationships: []*relationship.RelationshipDefinition{},
	}

	identified := identifyAliasRelationships(pod, relDef, design)
	if len(identified) != 2 {
		t.Fatalf("expected 2 alias relationships identified (one for containers, one for initContainers), got %d", len(identified))
	}
}

func TestFromAndToComponentID_RangingSelectors(t *testing.T) {
	fromUUID := uuid.Must(uuid.NewV4())
	toUUID := uuid.Must(uuid.NewV4())

	// Selector set 0 has empty Allow.From and Allow.To; Selector set 1 has the IDs
	rel := &relationship.RelationshipDefinition{
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{},
					To:   []relationship.SelectorItem{},
				},
			},
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{ID: &fromUUID},
					},
					To: []relationship.SelectorItem{
						{ID: &toUUID},
					},
				},
			},
		},
	}

	if got := fromComponentID(rel); got != fromUUID.String() {
		t.Fatalf("fromComponentID failed to range selectors: want %s, got %s", fromUUID, got)
	}
	if got := toComponentID(rel); got != toUUID.String() {
		t.Fatalf("toComponentID failed to range selectors: want %s, got %s", toUUID, got)
	}
}
