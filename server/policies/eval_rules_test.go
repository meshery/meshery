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

func TestValidateRelationshipsInDesign_SeveredConfigurationPathWhenComponentsExist(t *testing.T) {
	podComp := newTestComponent("pod-1", "Pod")
	// Port initially exists in containers configuration
	podComp.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{
					"name": "nginx",
					"ports": []interface{}{
						map[string]interface{}{
							"containerPort": 80,
						},
					},
				},
			},
		},
	}

	aliasComp := newTestComponent("port-alias", "Container")

	path := []string{"configuration", "spec", "containers", "0", "ports", "0"}
	mutatorRef := relationship.MutatorRef{path}
	mutatedRef := relationship.MutatedRef{path}
	patch := relationship.RelationshipDefinitionSelectorsPatch{
		MutatorRef: &mutatorRef,
		MutatedRef: &mutatedRef,
	}

	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "alias",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{ID: &aliasComp.ID, RelationshipDefinitionSelectorsPatch: &patch}},
					To:   []relationship.SelectorItem{{ID: &podComp.ID, RelationshipDefinitionSelectorsPatch: &patch}},
				},
			},
		},
	}
	setRelStatus(rel, StatusApproved)

	// Both components exist in the design
	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{podComp, aliasComp},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	policy := &AliasPolicy{}

	// Step 1: Intact port is valid and produces no deletion action
	if policy.IsInvalid(rel, design) {
		t.Fatalf("expected policy.IsInvalid to return false when referenced configuration path exists")
	}
	initialActions := validateRelationshipsInDesign(design, policy)
	if len(initialActions) != 0 {
		t.Fatalf("expected 0 deletion actions for intact alias relationship, got %d", len(initialActions))
	}

	// Step 2: Remove the port so configuration path is severed
	podComp.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{
					"name": "nginx",
					// ports array removed
				},
			},
		},
	}

	// Step 3: policy.IsInvalid returns true (the relationship is semantically invalid)
	if !policy.IsInvalid(rel, design) {
		t.Fatalf("expected policy.IsInvalid to return true for severed alias configuration path")
	}

	// Step 4: Both components exist (fromOrToComponentsDontExist returns false)
	if fromOrToComponentsDontExist(rel, design) {
		t.Fatalf("expected fromOrToComponentsDontExist to return false because both components are declared in design")
	}

	// Step 5: Run validateRelationshipsInDesign.
	actions := validateRelationshipsInDesign(design, policy)
	if len(actions) != 1 {
		t.Fatalf("expected 1 validation action marking severed alias relationship as deleted, got %d actions", len(actions))
	}
	if actions[0].StringValue != StatusDeleted {
		t.Fatalf("expected status %s, got %v", StatusDeleted, actions[0].StringValue)
	}
}

func TestValidateRelationshipsInDesign_PendingAliasPreservedWhenComponentMissing(t *testing.T) {
	podComp := newTestComponent("pod-1", "Pod")
	podComp.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{
					"name": "nginx",
					"ports": []interface{}{
						map[string]interface{}{
							"containerPort": 80,
						},
					},
				},
			},
		},
	}

	aliasCompID := uuid.Must(uuid.NewV4())
	path := []string{"configuration", "spec", "containers", "0", "ports", "0"}
	mutatorRef := relationship.MutatorRef{path}
	mutatedRef := relationship.MutatedRef{path}
	patch := relationship.RelationshipDefinitionSelectorsPatch{
		MutatorRef: &mutatorRef,
		MutatedRef: &mutatedRef,
	}

	rel := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "alias",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{ID: &aliasCompID, RelationshipDefinitionSelectorsPatch: &patch}},
					To:   []relationship.SelectorItem{{ID: &podComp.ID, RelationshipDefinitionSelectorsPatch: &patch}},
				},
			},
		},
	}
	setRelStatus(rel, StatusPending)

	// Design only contains podComp; aliasComp has not been created yet by Phase 3 side effects
	design := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{podComp},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	policy := &AliasPolicy{}

	// Endpoint absence evaluates to true because aliasComp is not yet in the design
	if !fromOrToComponentsDontExist(rel, design) {
		t.Fatalf("expected fromOrToComponentsDontExist to return true for missing alias component")
	}

	// Pending alias relationships must not be treated as invalid
	if policy.IsInvalid(rel, design) {
		t.Fatalf("expected policy.IsInvalid to return false for pending alias relationship")
	}

	// validateRelationshipsInDesign must NOT mark the pending relationship as deleted
	actions := validateRelationshipsInDesign(design, policy)
	if len(actions) != 0 {
		t.Fatalf("expected 0 actions for pending alias with missing component, got %d", len(actions))
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
				{Kind: "Service"},                          // index 0 has no mutatorRef
				{Kind: "Service", MutatorRef: &mutatorRef}, // index 1 has mutatorRef
			},
			To: &[]relationship.MatchSelectorItem{
				{Kind: "Pod"},                          // index 0 has no mutatedRef
				{Kind: "Pod", MutatedRef: &mutatedRef}, // index 1 has mutatedRef
			},
		},
	}

	// 1. Verify patchBindingSelectorIDsTyped stamps all items without mutating the source template
	patchedSel := patchBindingSelectorIDsTyped(&sel, fromID, toID)
	for i, item := range *patchedSel.Match.From {
		if item.ID == nil || item.ID.String() != fromID {
			t.Fatalf("expected Match.From[%d].ID to be %s, got %v", i, fromID, item.ID)
		}
	}
	for i, item := range *patchedSel.Match.To {
		if item.ID == nil || item.ID.String() != toID {
			t.Fatalf("expected Match.To[%d].ID to be %s, got %v", i, toID, item.ID)
		}
	}
	sel = patchedSel

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

// TestIdentifyBindingRelationships_MatchSelectorIDIsolation verifies that when multiple
// binding relationships are identified in a design (e.g. multiple ClusterRoleBindings or Services),
// ID stamping deep-copies the MatchSelector so that subsequent iterations do not mutate the
// MatchSelector IDs of previously identified relationships in-place.
func TestIdentifyBindingRelationships_MatchSelectorIDIsolation(t *testing.T) {
	relDefID := uuid.Must(uuid.NewV4())
	mutatorRef := [][]string{{"configuration", "spec", "roleName"}}
	mutatedRef := [][]string{{"configuration", "spec", "boundRole"}}

	fromKind := "ClusterRole"
	bindingKind := "ClusterRoleBinding"
	toKind := "ServiceAccount"

	relDef := &relationship.RelationshipDefinition{
		ID:               relDefID,
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "edge",
		SubType:          "binding",
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							Kind: &fromKind,
							Match: &relationship.MatchSelector{
								From: &[]relationship.MatchSelectorItem{
									{Kind: "self", MutatorRef: &mutatorRef},
								},
								To: &[]relationship.MatchSelectorItem{
									{Kind: bindingKind, MutatedRef: &mutatedRef},
								},
							},
						},
					},
					To: []relationship.SelectorItem{
						{
							Kind: &toKind,
							Match: &relationship.MatchSelector{
								From: &[]relationship.MatchSelectorItem{
									{Kind: bindingKind, MutatorRef: &mutatorRef},
								},
								To: &[]relationship.MatchSelectorItem{
									{Kind: "self", MutatedRef: &mutatedRef},
								},
							},
						},
					},
				},
			},
		},
	}

	// Pair 1
	cr1 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: fromKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"roleName": "viewer"},
		},
	}
	crb1 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: bindingKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"boundRole": "viewer", "roleName": "viewer"},
		},
	}
	sa1 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: toKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"boundRole": "viewer"},
		},
	}

	// Pair 2
	cr2 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: fromKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"roleName": "admin"},
		},
	}
	crb2 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: bindingKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"boundRole": "admin", "roleName": "admin"},
		},
	}
	sa2 := &component.ComponentDefinition{
		ID:        uuid.Must(uuid.NewV4()),
		Component: component.Component{Kind: toKind},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"boundRole": "admin"},
		},
	}

	design := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{cr1, crb1, sa1, cr2, crb2, sa2},
	}

	identified := identifyBindingRelationships(relDef, design)
	if len(identified) != 2 {
		t.Fatalf("expected 2 identified binding relationships, got %d", len(identified))
	}

	// Find the relationship corresponding to Pair 1
	var rel1, rel2 *relationship.RelationshipDefinition
	for _, rel := range identified {
		fromID := (*rel.Selectors)[0].Allow.From[0].ID
		if fromID != nil && *fromID == cr1.ID {
			rel1 = rel
		} else if fromID != nil && *fromID == cr2.ID {
			rel2 = rel
		}
	}

	if rel1 == nil || rel2 == nil {
		t.Fatalf("could not isolate both distinct identified relationships: rel1=%v, rel2=%v", rel1, rel2)
	}

	// Assert rel1's Match.From and Match.To retain cr1 and crb1 IDs
	rel1FromMatch := (*rel1.Selectors)[0].Allow.From[0].Match
	if rel1FromMatch == nil || rel1FromMatch.From == nil || (*rel1FromMatch.From)[0].ID == nil {
		t.Fatalf("rel1 Match.From is nil or unstamped")
	}
	if *(*rel1FromMatch.From)[0].ID != cr1.ID {
		t.Errorf("rel1 Match.From[0].ID clobbered: want %s (cr1), got %s", cr1.ID, *(*rel1FromMatch.From)[0].ID)
	}
	if *(*rel1FromMatch.To)[0].ID != crb1.ID {
		t.Errorf("rel1 Match.To[0].ID clobbered: want %s (crb1), got %s", crb1.ID, *(*rel1FromMatch.To)[0].ID)
	}

	// Assert rel2's Match.From and Match.To retain cr2 and crb2 IDs
	rel2FromMatch := (*rel2.Selectors)[0].Allow.From[0].Match
	if *(*rel2FromMatch.From)[0].ID != cr2.ID {
		t.Errorf("rel2 Match.From[0].ID invalid: want %s (cr2), got %s", cr2.ID, *(*rel2FromMatch.From)[0].ID)
	}
	if *(*rel2FromMatch.To)[0].ID != crb2.ID {
		t.Errorf("rel2 Match.To[0].ID invalid: want %s (crb2), got %s", crb2.ID, *(*rel2FromMatch.To)[0].ID)
	}
}

func TestCopyMatchSelector_PreservesRefs(t *testing.T) {
	refs := [][]string{
		{"configuration", "spec", "port"},
		{"configuration", "spec", "targetPort"},
	}
	fromID := uuid.Must(uuid.NewV4())
	toID := uuid.Must(uuid.NewV4())
	orig := &relationship.MatchSelector{
		From: &[]relationship.MatchSelectorItem{
			{Kind: "Service", ID: &fromID},
		},
		To: &[]relationship.MatchSelectorItem{
			{Kind: "Pod", ID: &toID},
		},
		Refs: &refs,
	}

	copied := copyMatchSelector(orig)
	if copied == nil {
		t.Fatal("expected non-nil copy")
	}
	if copied.Refs == nil {
		t.Fatal("expected copied Refs to be non-nil")
	}
	if len(*copied.Refs) != len(refs) {
		t.Fatalf("expected %d refs, got %d", len(refs), len(*copied.Refs))
	}
	// Compare each copied reference path with original before mutation
	for i, refPath := range refs {
		copiedPath := (*copied.Refs)[i]
		if len(copiedPath) != len(refPath) {
			t.Fatalf("expected ref[%d] length %d, got %d", i, len(refPath), len(copiedPath))
		}
		for j, elem := range refPath {
			if copiedPath[j] != elem {
				t.Errorf("expected copied.Refs[%d][%d] to be %q, got %q", i, j, elem, copiedPath[j])
			}
		}
	}

	// Verify deep copy isolation
	(*copied.Refs)[0][0] = "mutated"
	if refs[0][0] == "mutated" {
		t.Error("modifying copied Refs affected original refs slice")
	}
}

func newTestComponentWithModel(name, kind, modelName string) *component.ComponentDefinition {
	id := uuid.Must(uuid.NewV4())
	return &component.ComponentDefinition{
		ID:          id,
		DisplayName: name,
		Component: component.Component{
			Kind: kind,
		},
		ModelReference: modelv1beta1.ModelReference{
			Name: modelName,
		},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": name,
			},
		},
	}
}

// TestEdgeNonBindingPolicy_CallerPath_SelectorModelIsolation tests that mutator-based
// relationship identification respects selector model via the exported policy caller path.
func TestEdgeNonBindingPolicy_CallerPath_SelectorModelIsolation(t *testing.T) {
	policy := &EdgeNonBindingPolicy{}

	compRB := newTestComponentWithModel("app-rb", "RoleBinding", "kubernetes")
	compRB.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"roleRef": map[string]interface{}{
				"name": "app-role",
			},
		},
	}

	compAwsRole := newTestComponentWithModel("app-role", "Role", "aws-iam-controller")
	compK8sRole := newTestComponentWithModel("app-role", "Role", "kubernetes")

	design := makePatternFile([]*component.ComponentDefinition{compRB, compAwsRole, compK8sRole}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "roleRef", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "metadata", "name"}}

	fromSel := relationship.SelectorItem{
		Kind:  strPtr("RoleBinding"),
		Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatorRef: &mutatorRef,
		},
	}
	toSel := relationship.SelectorItem{
		Kind:  strPtr("Role"),
		Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatedRef: &mutatedRef,
		},
	}

	relDef := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.Edge,
		RelationshipType: "non-binding",
		SubType:          "reference",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors: &relationship.SelectorSet{
			{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{fromSel},
					To:   []relationship.SelectorItem{toSel},
				},
			},
		},
	}

	identified := policy.IdentifyRelationship(relDef, design)
	if len(identified) != 1 {
		t.Fatalf("expected exactly 1 identified relationship (positive control), got %d", len(identified))
	}

	identifiedToID := (*identified[0].Selectors)[0].Allow.To[0].ID
	if identifiedToID == nil || *identifiedToID != compK8sRole.ID {
		t.Errorf("expected relationship to target kubernetes Role %s, got %v", compK8sRole.ID, identifiedToID)
	}
	if *identifiedToID == compAwsRole.ID {
		t.Errorf("relationship wrongly identified with aws-iam-controller Role %s", compAwsRole.ID)
	}
}

// TestMatchComponentPairs_HelperPath_SelectorModelEnforcement tests that matchComponentPairs
// directly enforces the selector model and ignores cross-model same-kind components.
func TestMatchComponentPairs_HelperPath_SelectorModelEnforcement(t *testing.T) {
	compRB := newTestComponentWithModel("app-rb", "RoleBinding", "kubernetes")
	compRB.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"roleRef": map[string]interface{}{
				"name": "app-role",
			},
		},
	}

	compAwsRole := newTestComponentWithModel("app-role", "Role", "aws-iam-controller")
	compK8sRole := newTestComponentWithModel("app-role", "Role", "kubernetes")

	design := makePatternFile([]*component.ComponentDefinition{compRB, compAwsRole, compK8sRole}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "roleRef", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "metadata", "name"}}

	fromSel := relationship.SelectorItem{
		Kind:  strPtr("RoleBinding"),
		Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatorRef: &mutatorRef,
		},
	}
	toSel := relationship.SelectorItem{
		Kind:  strPtr("Role"),
		Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatedRef: &mutatedRef,
		},
	}

	relDef := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.Edge,
		RelationshipType: "non-binding",
		SubType:          "reference",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
	}

	matches := matchComponentPairs(fromSel, toSel, design.Components, relDef, design)
	if len(matches) != 1 {
		t.Fatalf("expected 1 match (kubernetes Role), got %d", len(matches))
	}

	matchedToID := (*matches[0].Selectors)[0].Allow.To[0].ID
	if matchedToID == nil || *matchedToID != compK8sRole.ID {
		t.Errorf("expected match target %s, got %v", compK8sRole.ID, matchedToID)
	}
}

// TestMatchComponentPairs_SyntheticTwoModelKindCollision verifies that when two different models
// share the same component kind and property values, the selector's model restricts the matches.
func TestMatchComponentPairs_SyntheticTwoModelKindCollision(t *testing.T) {
	compFrom := newTestComponentWithModel("db-source", "Database", "provider-a")
	compFrom.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"target": "db-prod",
		},
	}

	compWrongModel := newTestComponentWithModel("db-target-wrong", "DatabaseInstance", "provider-b")
	compWrongModel.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "db-prod",
		},
	}

	compRightModel := newTestComponentWithModel("db-target-right", "DatabaseInstance", "provider-a")
	compRightModel.Configuration = map[string]interface{}{
		"spec": map[string]interface{}{
			"name": "db-prod",
		},
	}

	design := makePatternFile([]*component.ComponentDefinition{compFrom, compWrongModel, compRightModel}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "target"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "name"}}

	fromSel := relationship.SelectorItem{
		Kind:  strPtr("Database"),
		Model: &modelv1beta1.ModelReference{Name: "provider-a"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatorRef: &mutatorRef,
		},
	}
	toSel := relationship.SelectorItem{
		Kind:  strPtr("DatabaseInstance"),
		Model: &modelv1beta1.ModelReference{Name: "provider-a"},
		RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
			MutatedRef: &mutatedRef,
		},
	}

	relDef := &relationship.RelationshipDefinition{
		ID:               uuid.Must(uuid.NewV4()),
		Kind:             relationship.Edge,
		RelationshipType: "non-binding",
		SubType:          "reference",
	}

	matches := matchComponentPairs(fromSel, toSel, design.Components, relDef, design)
	if len(matches) != 1 {
		t.Fatalf("expected exactly 1 match, got %d", len(matches))
	}
	toID := (*matches[0].Selectors)[0].Allow.To[0].ID
	if toID == nil || *toID != compRightModel.ID {
		t.Errorf("expected target to be compRightModel (%s), got %v", compRightModel.ID, toID)
	}
}

// TestIsRelationshipFeasible_ModelSemanticsPin verifies the contract for selector model matching:
// nil, empty string, and wildcard model references match permissively.
func TestIsRelationshipFeasible_ModelSemanticsPin(t *testing.T) {
	comp := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Role"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
	}

	tests := []struct {
		name     string
		selector relationship.SelectorItem
		want     bool
	}{
		{
			name:     "nil Model is permissive",
			selector: relationship.SelectorItem{Kind: strPtr("Role"), Model: nil},
			want:     true,
		},
		{
			name:     "empty Model.Name is permissive",
			selector: relationship.SelectorItem{Kind: strPtr("Role"), Model: &modelv1beta1.ModelReference{Name: ""}},
			want:     true,
		},
		{
			name:     "wildcard Model.Name matches any model",
			selector: relationship.SelectorItem{Kind: strPtr("Role"), Model: &modelv1beta1.ModelReference{Name: "*"}},
			want:     true,
		},
		{
			name:     "matching Model.Name",
			selector: relationship.SelectorItem{Kind: strPtr("Role"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
			want:     true,
		},
		{
			name:     "mismatched Model.Name",
			selector: relationship.SelectorItem{Kind: strPtr("Role"), Model: &modelv1beta1.ModelReference{Name: "aws-iam-controller"}},
			want:     false,
		},
		{
			name:     "kind mismatch even if model matches",
			selector: relationship.SelectorItem{Kind: strPtr("Pod"), Model: &modelv1beta1.ModelReference{Name: "kubernetes"}},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isRelationshipFeasible(tt.selector, comp)
			if got != tt.want {
				t.Errorf("isRelationshipFeasible() = %v, want %v", got, tt.want)
			}
		})
	}
}
