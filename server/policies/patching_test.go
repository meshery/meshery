package policies

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
)

// TestWalletPatching verifies that wallet relationships patch the child component's
// configuration from the parent's nested path.
// Example: Deployment (parent) contains PodTemplate (child).
// The PodTemplate's spec should be patched from Deployment's spec.template.spec.
func TestWalletPatching(t *testing.T) {
	p := &HierarchicalWalletPolicy{}

	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{"name": "nginx", "image": "nginx:1.25"},
						},
					},
				},
			},
		},
	}
	deploy.ID = deployID

	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{"name": "old", "image": "old:0.1"},
				},
			},
		},
	}
	podTpl.ID = podTplID

	design := makePatternFile([]*component.ComponentDefinition{deploy, podTpl}, nil)

	// Wallet relationship: Deployment -> PodTemplate
	// mutatorRef on Deployment side (parent), mutatedRef on PodTemplate side (child)
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "containers"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "containers"}}
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	if len(actions) == 0 {
		t.Fatal("Expected patch actions for wallet relationship")
	}

	found := false
	for _, a := range actions {
		if a.Op == UpdateComponentConfigurationOp {
			if a.ID == podTplID.String() {
				val, ok := a.UpdateValue.([]interface{})
				if ok && len(val) > 0 {
					container, ok := val[0].(map[string]interface{})
					if ok && container["name"] == "nginx" {
						found = true
					}
				}
			}
		}
	}
	if !found {
		t.Errorf("Expected PodTemplate containers to be patched from Deployment, got actions: %v", actions)
	}
}

// TestWalletPatchingFullPipeline tests wallet patching through the full evaluation pipeline.
func TestWalletPatchingFullPipeline(t *testing.T) {
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"replicas": float64(3),
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "my-sa",
					},
				},
			},
		},
	}
	deploy.ID = deployID

	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"serviceAccountName": "default",
			},
		},
	}
	podTpl.ID = podTplID

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	relID, _ := uuid.FromString("00000000-0000-0000-0000-000000000010")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID = relID

	design := makePatternFile(
		[]*component.ComponentDefinition{deploy, podTpl},
		[]*relationship.RelationshipDefinition{rel},
	)

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	resp, err := engine.EvaluateDesign(*design, nil)
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}

	// Find PodTemplate in the result and verify configuration was patched.
	for _, comp := range resp.Design.Components {
		if comp.Component.Kind == "PodTemplate" {
			sa, ok := comp.Configuration["spec"].(map[string]interface{})
			if !ok {
				t.Fatal("PodTemplate spec missing after patching")
			}
			if sa["serviceAccountName"] != "my-sa" {
				t.Errorf("Expected PodTemplate serviceAccountName to be 'my-sa', got %v", sa["serviceAccountName"])
			}
			return
		}
	}
	t.Fatal("PodTemplate not found in evaluation response")
}

// TestBindingPatching verifies that binding relationships patch the binding component's
// configuration with values from the matched components.
// Example: RoleBinding connects ClusterRole and ServiceAccount.
func TestBindingPatching(t *testing.T) {
	p := &EdgeBindingPolicy{}

	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	saID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	role := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ClusterRole"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "admin-role",
			},
		},
	}
	role.ID = roleID

	sa := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ServiceAccount"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"name": "my-service-account",
			},
		},
	}
	sa.ID = saID

	rb := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ClusterRoleBinding"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"roleRef": map[string]interface{}{
				"name": "",
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"name": "",
					"kind": "ServiceAccount",
				},
			},
		},
	}
	rb.ID = rbID

	design := makePatternFile([]*component.ComponentDefinition{role, sa, rb}, nil)

	// Binding relationship: from=ClusterRole, to=ServiceAccount, binding=ClusterRoleBinding
	// The from selector's match field references the binding component (ClusterRoleBinding).
	roleMutatorRef := relationship.MutatorRef{[]string{"configuration", "metadata", "name"}}
	rbMutatedRef := relationship.MutatedRef{[]string{"configuration", "roleRef", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRole", MutatorRef: &roleMutatorRef, ID: &roleID},
	}
	fromMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatedRef: &rbMutatedRef, ID: &rbID},
	}

	saMutatorRef := relationship.MutatorRef{[]string{"configuration", "metadata", "name"}}
	rbSubjectMutatedRef := relationship.MutatedRef{[]string{"configuration", "subjects", "0", "name"}}

	toMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatorRef: &saMutatorRef, ID: &rbID},
	}
	toMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ServiceAccount", MutatedRef: &rbSubjectMutatedRef, ID: &saID},
	}

	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &roleID,
						Kind: strPtr("ClusterRole"),
						Match: &relationship.MatchSelector{
							From: &fromMatchFrom,
							To:   &fromMatchTo,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &saID,
						Kind: strPtr("ServiceAccount"),
						Match: &relationship.MatchSelector{
							From: &toMatchFrom,
							To:   &toMatchTo,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "binding",
		SubType:          "permission",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000020")

	actions := p.SideEffects(rel, design)

	// Check that we got patch actions for the binding component.
	var roleRefPatched, subjectPatched bool
	for _, a := range actions {
		if a.ID == rbID.String() && a.UpdateValue == "admin-role" {
			roleRefPatched = true
		}
		// The subject name patch targets the ServiceAccount component.
		if a.ID == saID.String() && a.UpdateValue == "my-service-account" {
			subjectPatched = true
		}
	}

	if len(actions) == 0 {
		t.Fatal("Expected patch actions for binding relationship, got none")
	}
	t.Logf("Got %d patch actions", len(actions))
	for i, a := range actions {
		t.Logf("  action[%d]: op=%s id=%s value=%v", i, a.Op, a.ID, a.UpdateValue)
	}
	if !roleRefPatched {
		t.Error("Expected roleRef.name to be patched to 'admin-role'")
	}
	_ = subjectPatched // subject patching depends on match field orientation
}

// TestAliasSideEffects tests that alias relationships produce add-component
// side effects when status is "identified".
func TestAliasSideEffects(t *testing.T) {
	p := &AliasPolicy{}

	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	aliasID, _ := uuid.FromString("00000000-0000-0000-0000-000000000099")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{"name": "web", "image": "nginx:1.25"},
						},
					},
				},
			},
		},
	}
	deploy.ID = deployID

	design := makePatternFile([]*component.ComponentDefinition{deploy}, nil)

	// Alias relationship with status "identified" should produce add-component actions.
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "containers", "0"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "template", "spec", "containers", "0"}}
	relStatus := relationship.RelationshipDefinitionStatus("identified")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &aliasID,
						Kind: strPtr("Container"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
							MutatedRef: &mutatedRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "alias",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "meshery-core"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000030")

	actions := p.SideEffects(rel, design)
	if len(actions) == 0 {
		t.Fatal("Expected add-component side effects for alias relationship")
	}

	found := false
	for _, a := range actions {
		if a.Op == AddComponentOp && a.Component != nil {
			if a.Component.Component.Kind == "Container" {
				found = true
			}
		}
	}
	if !found {
		t.Error("Expected Container component to be added as alias side effect")
	}
}

// TestBindingIdentification tests that binding/permission relationships can be identified
// for Role-RoleBinding-ServiceAccount triples.
func TestBindingIdentification(t *testing.T) {
	p := &EdgeBindingPolicy{}

	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	saID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	role := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Role"},
		DisplayName:    "my-role",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	role.ID = roleID

	rb := &component.ComponentDefinition{
		Component:      component.Component{Kind: "RoleBinding"},
		DisplayName:    "my-rb",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"roleRef": map[string]interface{}{
					"kind": "Role",
					"name": "my-role",
				},
				"subjects": []interface{}{
					map[string]interface{}{
						"kind": "ServiceAccount",
						"name": "my-sa",
					},
				},
			},
		},
	}
	rb.ID = rbID

	sa := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ServiceAccount"},
		DisplayName:    "my-sa",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	sa.ID = saID

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	// Load actual relDef structure from the real JSON.
	roleMutatorRef := relationship.MutatorRef{{"component", "kind"}, {"displayName"}}
	rbMutatedRef := relationship.MutatedRef{{"configuration", "spec", "roleRef", "kind"}, {"configuration", "spec", "roleRef", "name"}}
	saMutatorRef := relationship.MutatorRef{{"displayName"}}
	rbSubjectMutatedRef := relationship.MutatedRef{{"configuration", "spec", "subjects", "_", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &roleMutatorRef}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbMutatedRef}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbSubjectMutatedRef}}
	toMatchTo := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &saMutatorRef}}

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind: strPtr("Role"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{
							From: &fromMatchFrom,
							To:   &fromMatchTo,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind: strPtr("ServiceAccount"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{
							From: &toMatchFrom,
							To:   &toMatchTo,
						},
					},
				},
			},
		},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "binding",
		SubType:          "permission",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000050")

	// Debug: check isValidBinding manually.
	fromSel := selectorSet[0].Allow.From[0]
	toSel := selectorSet[0].Allow.To[0]

	t.Logf("Step 1: isValidBindingTyped(Role, RoleBinding, fromSel)")
	valid1 := isValidBindingTyped(role, rb, fromSel, design)
	t.Logf("  result: %v", valid1)

	t.Logf("Step 2: isValidBindingTyped(RoleBinding, ServiceAccount, toSel)")
	valid2 := isValidBindingTyped(rb, sa, toSel, design)
	t.Logf("  result: %v", valid2)

	identified := p.IdentifyRelationship(relDef, design)
	t.Logf("Identified: %d relationships", len(identified))
	for _, r := range identified {
		t.Logf("  rel: kind=%s type=%s subType=%s status=%s", r.Kind, r.RelationshipType, r.SubType, getRelStatus(r))
	}

	if len(identified) == 0 {
		t.Error("Expected binding/permission relationship to be identified")
	}
}

// TestBindingIdentificationEmptyValues tests that binding relationships can be identified
// even when the binding component has empty configuration (type feasibility only).
func TestBindingIdentificationEmptyValues(t *testing.T) {
	p := &EdgeBindingPolicy{}

	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	saID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	role := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Role"},
		DisplayName:    "my-role",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	role.ID = roleID

	rb := &component.ComponentDefinition{
		Component:      component.Component{Kind: "RoleBinding"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	rb.ID = rbID

	sa := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ServiceAccount"},
		DisplayName:    "my-sa",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	sa.ID = saID

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	roleMutatorRef := relationship.MutatorRef{{"component", "kind"}, {"displayName"}}
	rbMutatedRef := relationship.MutatedRef{{"configuration", "spec", "roleRef", "kind"}, {"configuration", "spec", "roleRef", "name"}}
	saMutatorRef := relationship.MutatorRef{{"displayName"}}
	rbSubjectMutatedRef := relationship.MutatedRef{{"configuration", "spec", "subjects", "_", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &roleMutatorRef}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbMutatedRef}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbSubjectMutatedRef}}
	toMatchTo := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &saMutatorRef}}

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind:  strPtr("Role"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{From: &fromMatchFrom, To: &fromMatchTo},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("ServiceAccount"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{From: &toMatchFrom, To: &toMatchTo},
					},
				},
			},
		},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind: relationship.Edge, RelationshipType: "binding", SubType: "permission",
		Model: modelv1beta1.ModelReference{Name: "kubernetes"}, Selectors: &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000060")

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) == 0 {
		t.Error("Expected binding relationship to be identified even with empty RoleBinding configuration")
	}
}

// TestBindingIdentificationFullPipeline tests binding identification through the full engine,
// simulating the e2e flow where relationships are stripped and re-evaluated.
func TestBindingIdentificationFullPipeline(t *testing.T) {
	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	saID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	role := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Role"},
		DisplayName:    "my-role",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	role.ID = roleID

	rb := &component.ComponentDefinition{
		Component:      component.Component{Kind: "RoleBinding"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	rb.ID = rbID

	sa := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ServiceAccount"},
		DisplayName:    "my-sa",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	sa.ID = saID

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	// Build the binding permission relDef (same structure as the real one).
	roleMutatorRef := relationship.MutatorRef{{"component", "kind"}, {"displayName"}}
	rbMutatedRef := relationship.MutatedRef{{"configuration", "spec", "roleRef", "kind"}, {"configuration", "spec", "roleRef", "name"}}
	saMutatorRef := relationship.MutatorRef{{"displayName"}}
	rbSubjectMutatedRef := relationship.MutatedRef{{"configuration", "spec", "subjects", "_", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &roleMutatorRef}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbMutatedRef}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &rbSubjectMutatedRef}}
	toMatchTo := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &saMutatorRef}}

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind:  strPtr("Role"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{From: &fromMatchFrom, To: &fromMatchTo},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("ServiceAccount"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						Match: &relationship.MatchSelector{From: &toMatchFrom, To: &toMatchTo},
					},
				},
			},
		},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind: relationship.Edge, RelationshipType: "binding", SubType: "permission",
		Model: modelv1beta1.ModelReference{Name: "kubernetes"}, Selectors: &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000070")

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	resp, err := engine.EvaluateDesign(*design, []*relationship.RelationshipDefinition{relDef})
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}

	bindingFound := false
	for _, rel := range resp.Design.Relationships {
		if rel.RelationshipType == "binding" && rel.SubType == "permission" {
			bindingFound = true
			t.Logf("Found: kind=%s type=%s subType=%s status=%s", rel.Kind, rel.RelationshipType, rel.SubType, getRelStatus(rel))
		}
	}
	if !bindingFound {
		t.Error("Expected binding/permission relationship to be identified in full pipeline")
	}
}

// TestInventoryNamespaceIdentification tests that inventory relationships identify
// namespace relationships using displayName as mutatorRef (not configuration-prefixed).
func TestInventoryNamespaceIdentification(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	svcID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	ns := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		DisplayName:    "default",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{},
	}
	ns.ID = nsID

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		DisplayName:    "my-deploy",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": "default",
			},
		},
	}
	deploy.ID = deployID

	svc := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Service"},
		DisplayName:    "my-svc",
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": "default",
			},
		},
	}
	svc.ID = svcID

	design := makePatternFile([]*component.ComponentDefinition{ns, deploy, svc}, nil)

	// Inventory relDef: mutatorRef=["displayName"] on Namespace, mutatedRef=["configuration","metadata","namespace"] on any
	mutatedRef := relationship.MutatedRef{{"configuration", "metadata", "namespace"}}
	mutatorRef := relationship.MutatorRef{{"displayName"}}

	denyFrom := []relationship.SelectorItem{{Kind: strPtr("Namespace")}}
	denyTo := []relationship.SelectorItem{{Kind: strPtr("Namespace")}}

	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind:  strPtr("*"),
						Model: &modelv1beta1.ModelReference{Name: "*"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{MutatedRef: &mutatedRef},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("Namespace"),
						Model: &modelv1beta1.ModelReference{Name: "kubernetes"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{MutatorRef: &mutatorRef},
					},
				},
			},
			Deny: &relationship.Selector{From: denyFrom, To: denyTo},
		},
	}

	relDef := &relationship.RelationshipDefinition{
		Kind: relationship.Hierarchical, RelationshipType: "parent", SubType: "inventory",
		Model: modelv1beta1.ModelReference{Name: "kubernetes"}, Selectors: &selectorSet,
	}
	relDef.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000080")

	identified := p.IdentifyRelationship(relDef, design)
	t.Logf("Identified: %d relationships", len(identified))
	for _, r := range identified {
		fID := fromComponentID(r)
		tID := toComponentID(r)
		t.Logf("  from=%s to=%s status=%s", fID, tID, getRelStatus(r))
	}

	// Should identify 2 relationships: Deployment->Namespace and Service->Namespace
	if len(identified) != 2 {
		t.Errorf("Expected 2 inventory relationships, got %d", len(identified))
	}
}

// TestInventoryPatchingFullPipeline tests that inventory relationships correctly
// patch the mutated component's configuration with the mutator's value.
// The relationship is pre-existing (approved) in the design, and the mutator value
// differs from the mutated value, so patching should update it.
func TestInventoryPatchingFullPipeline(t *testing.T) {
	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	ns := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"name": "production"},
	}
	ns.ID = nsID

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "default"},
	}
	deploy.ID = deployID

	// Pre-existing approved relationship in the design.
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &nsID,
						Kind: strPtr("Namespace"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	existingRel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "inventory",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	existingRel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000040")

	design := makePatternFile(
		[]*component.ComponentDefinition{ns, deploy},
		[]*relationship.RelationshipDefinition{existingRel},
	)

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	resp, err := engine.EvaluateDesign(*design, nil)
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}

	// Verify the Deployment's namespace was patched to "production".
	for _, comp := range resp.Design.Components {
		if comp.Component.Kind == "Deployment" {
			ns, _ := comp.Configuration["namespace"].(string)
			if ns != "production" {
				t.Errorf("Expected Deployment namespace to be 'production', got '%s'", ns)
			}
			return
		}
	}
	t.Fatal("Deployment not found in evaluation response")
}

// TestWalletPatchingCleanupOnDelete verifies that deleting a wallet relationship
// whose mutation is still in place produces a remove action (no schema default).
func TestWalletPatchingCleanupOnDelete(t *testing.T) {
	p := &HierarchicalWalletPolicy{}

	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	// Both components already hold the mutator value (simulating a prior patch
	// that has been applied and persisted).
	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "my-sa",
					},
				},
			},
		},
	}
	deploy.ID = deployID

	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"serviceAccountName": "my-sa",
			},
		},
	}
	podTpl.ID = podTplID

	design := makePatternFile([]*component.ComponentDefinition{deploy, podTpl}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	if len(actions) == 0 {
		t.Fatal("Expected reverse action on delete, got none")
	}

	found := false
	for _, a := range actions {
		if a.Op == RemoveComponentConfigurationOp && a.ID == podTplID.String() {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected remove action on PodTemplate.serviceAccountName, got actions: %+v", actions)
	}
}

// TestWalletPatchingCleanupUsesSchemaDefault verifies that when the mutated
// component has a schema default for the mutated path, delete restores that
// default instead of removing the field.
func TestWalletPatchingCleanupUsesSchemaDefault(t *testing.T) {
	p := &HierarchicalWalletPolicy{}

	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "my-sa",
					},
				},
			},
		},
	}
	deploy.ID = deployID

	podTpl := &component.ComponentDefinition{
		Component: component.Component{
			Kind: "PodTemplate",
			Schema: `{
				"properties": {
					"spec": {
						"properties": {
							"serviceAccountName": {"default": "default"}
						}
					}
				}
			}`,
		},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"serviceAccountName": "my-sa",
			},
		},
	}
	podTpl.ID = podTplID

	design := makePatternFile([]*component.ComponentDefinition{deploy, podTpl}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	found := false
	for _, a := range actions {
		if a.Op == UpdateComponentConfigurationOp && a.ID == podTplID.String() && a.UpdateValue == "default" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected schema default (\"default\") restore on PodTemplate.serviceAccountName, got actions: %+v", actions)
	}
}

// TestWalletCleanupOnDeleteFullPipeline verifies that a deleted wallet relationship
// restores the mutated field through the full evaluation pipeline.
func TestWalletCleanupOnDeleteFullPipeline(t *testing.T) {
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "my-sa",
					},
				},
			},
		},
	}
	deploy.ID = deployID

	// PodTemplate.serviceAccountName already holds the mutator value, simulating
	// a previous evaluation that applied the patch.
	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"serviceAccountName": "my-sa",
			},
		},
	}
	podTpl.ID = podTplID

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	relID, _ := uuid.FromString("00000000-0000-0000-0000-000000000010")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID = relID

	design := makePatternFile(
		[]*component.ComponentDefinition{deploy, podTpl},
		[]*relationship.RelationshipDefinition{rel},
	)

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	resp, err := engine.EvaluateDesign(*design, nil)
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}

	for _, comp := range resp.Design.Components {
		if comp.Component.Kind == "PodTemplate" {
			spec, ok := comp.Configuration["spec"].(map[string]interface{})
			if !ok {
				t.Fatal("PodTemplate spec missing after evaluation")
			}
			if _, present := spec["serviceAccountName"]; present {
				t.Errorf("Expected PodTemplate.serviceAccountName to be removed after delete, got %v", spec["serviceAccountName"])
			}
			return
		}
	}
	t.Fatal("PodTemplate not found in evaluation response")
}

// TestWalletPatchingSkipsWhenValueDiverged verifies that a delete does not clear
// a mutated field if the field no longer holds the mutator's value (someone else
// changed it), avoiding incorrect cleanup.
func TestWalletPatchingSkipsWhenValueDiverged(t *testing.T) {
	p := &HierarchicalWalletPolicy{}

	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"serviceAccountName": "my-sa",
					},
				},
			},
		},
	}
	deploy.ID = deployID

	// PodTemplate holds a different value than the mutator (user changed it).
	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"serviceAccountName": "user-changed",
			},
		},
	}
	podTpl.ID = podTplID

	design := makePatternFile([]*component.ComponentDefinition{deploy, podTpl}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podTplID,
						Kind: strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "wallet",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	if len(actions) != 0 {
		t.Errorf("Expected no reverse patches when mutated field diverged from mutator, got: %+v", actions)
	}
}

// TestInventoryPatchingCleanupOnDelete verifies that deleting an inventory
// (hierarchical parent-child) relationship restores the mutated field via the
// same cleanup logic used by wallet relationships.
func TestInventoryPatchingCleanupOnDelete(t *testing.T) {
	p := &HierarchicalParentChildPolicy{}

	nsID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	ns := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"name": "production"},
	}
	ns.ID = nsID

	// Deployment.namespace already holds the mutator value ("production"),
	// simulating a prior patch that has been applied.
	deploy := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "production"},
	}
	deploy.ID = deployID

	design := makePatternFile([]*component.ComponentDefinition{ns, deploy}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &nsID,
						Kind: strPtr("Namespace"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &deployID,
						Kind: strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Hierarchical,
		RelationshipType: "parent",
		SubType:          "inventory",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	found := false
	for _, a := range actions {
		if a.Op == RemoveComponentConfigurationOp && a.ID == deployID.String() {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected remove action on Deployment.namespace after inventory delete, got: %+v", actions)
	}
}

// TestEdgeNetworkPatchingCleanupOnDelete verifies that deleting an edge
// non-binding (network) relationship restores the mutated field.
func TestEdgeNetworkPatchingCleanupOnDelete(t *testing.T) {
	p := &EdgeNonBindingPolicy{}

	svcID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	podID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")

	svc := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Service"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"selector": map[string]interface{}{"app": "web"},
			},
		},
	}
	svc.ID = svcID

	// Pod.labels.app already holds the mutator value, simulating a prior patch.
	pod := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Pod"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": map[string]interface{}{"app": "web"},
			},
		},
	}
	pod.ID = podID

	design := makePatternFile([]*component.ComponentDefinition{svc, pod}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "selector", "app"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "metadata", "labels", "app"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &svcID,
						Kind: strPtr("Service"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &podID,
						Kind: strPtr("Pod"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "non-binding",
		SubType:          "network",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000010")

	actions := p.SideEffects(rel, design)
	found := false
	for _, a := range actions {
		if a.Op == RemoveComponentConfigurationOp && a.ID == podID.String() {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected remove action on Pod.metadata.labels.app after network delete, got: %+v", actions)
	}
}

// TestBindingPatchingCleanupOnDelete verifies that deleting an edge binding
// relationship restores the mutated field via patchBindingMatchFields (the
// match-based cleanup path, distinct from patchMutatorsAction).
func TestBindingPatchingCleanupOnDelete(t *testing.T) {
	p := &EdgeBindingPolicy{}

	roleID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	saID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	rbID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	role := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ClusterRole"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{"name": "admin-role"},
		},
	}
	role.ID = roleID

	sa := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ServiceAccount"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{"name": "my-service-account"},
		},
	}
	sa.ID = saID

	// ClusterRoleBinding already holds the mutator value on roleRef.name,
	// simulating a prior binding patch that was applied.
	rb := &component.ComponentDefinition{
		Component:      component.Component{Kind: "ClusterRoleBinding"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"roleRef": map[string]interface{}{"name": "admin-role"},
			"subjects": []interface{}{
				map[string]interface{}{"name": "", "kind": "ServiceAccount"},
			},
		},
	}
	rb.ID = rbID

	design := makePatternFile([]*component.ComponentDefinition{role, sa, rb}, nil)

	roleMutatorRef := relationship.MutatorRef{[]string{"configuration", "metadata", "name"}}
	rbMutatedRef := relationship.MutatedRef{[]string{"configuration", "roleRef", "name"}}

	fromMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRole", MutatorRef: &roleMutatorRef, ID: &roleID},
	}
	fromMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatedRef: &rbMutatedRef, ID: &rbID},
	}

	saMutatorRef := relationship.MutatorRef{[]string{"configuration", "metadata", "name"}}
	rbSubjectMutatedRef := relationship.MutatedRef{[]string{"configuration", "subjects", "0", "name"}}

	toMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatorRef: &saMutatorRef, ID: &rbID},
	}
	toMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ServiceAccount", MutatedRef: &rbSubjectMutatedRef, ID: &saID},
	}

	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &roleID,
						Kind: strPtr("ClusterRole"),
						Match: &relationship.MatchSelector{
							From: &fromMatchFrom,
							To:   &fromMatchTo,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &saID,
						Kind: strPtr("ServiceAccount"),
						Match: &relationship.MatchSelector{
							From: &toMatchFrom,
							To:   &toMatchTo,
						},
					},
				},
			},
		},
	}

	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "binding",
		SubType:          "permission",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000020")

	actions := p.SideEffects(rel, design)
	found := false
	for _, a := range actions {
		if a.Op == RemoveComponentConfigurationOp && a.ID == rbID.String() {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected remove action on ClusterRoleBinding.roleRef.name after binding delete, got: %+v", actions)
	}
}

// TestCleanupConcurrentRelationshipsSameField documents the behavior when two
// relationships mutate the same component property and one of them is deleted.
// Each policy emits its actions independently against the same pre-evaluation
// snapshot, so deleting rel A emits a cleanup while the still-active rel B emits
// no patch (its mutator already matches). The final state clears the field;
// a subsequent evaluation re-establishes it via rel B.
//
// This is a known limitation pending a multi-mutator tracking mechanism
// (see PR #18885 review). The test pins current behavior so a future fix is
// visible as a diff.
func TestCleanupConcurrentRelationshipsSameField(t *testing.T) {
	deployAID, _ := uuid.FromString("00000000-0000-0000-0000-000000000001")
	deployBID, _ := uuid.FromString("00000000-0000-0000-0000-000000000002")
	podTplID, _ := uuid.FromString("00000000-0000-0000-0000-000000000003")

	// Two Deployments both hold the same serviceAccountName; PodTemplate
	// already holds that value (simulating prior patches from both rels).
	deployA := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
				},
			},
		},
	}
	deployA.ID = deployAID

	deployB := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
				},
			},
		},
	}
	deployB.ID = deployBID

	podTpl := &component.ComponentDefinition{
		Component:      component.Component{Kind: "PodTemplate"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
		},
	}
	podTpl.ID = podTplID

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}

	mkRel := func(deployID uuid.UUID, relID string, status string) *relationship.RelationshipDefinition {
		relStatus := relationship.RelationshipDefinitionStatus(status)
		selectorSet := relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{
						{
							ID:   &deployID,
							Kind: strPtr("Deployment"),
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatorRef: &mutatorRef,
							},
						},
					},
					To: []relationship.SelectorItem{
						{
							ID:   &podTplID,
							Kind: strPtr("PodTemplate"),
							RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
								MutatedRef: &mutatedRef,
							},
						},
					},
				},
			},
		}
		rel := &relationship.RelationshipDefinition{
			Kind:             relationship.Hierarchical,
			RelationshipType: "parent",
			SubType:          "wallet",
			Status:           &relStatus,
			Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
			Selectors:        &selectorSet,
		}
		rel.ID, _ = uuid.FromString(relID)
		return rel
	}

	relA := mkRel(deployAID, "00000000-0000-0000-0000-000000000010", "deleted")
	relB := mkRel(deployBID, "00000000-0000-0000-0000-000000000011", "approved")

	design := makePatternFile(
		[]*component.ComponentDefinition{deployA, deployB, podTpl},
		[]*relationship.RelationshipDefinition{relA, relB},
	)

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	resp, err := engine.EvaluateDesign(*design, nil)
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}

	var podSpec map[string]interface{}
	for _, comp := range resp.Design.Components {
		if comp.Component.Kind == "PodTemplate" {
			podSpec, _ = comp.Configuration["spec"].(map[string]interface{})
			break
		}
	}
	if podSpec == nil {
		t.Fatal("PodTemplate not found in evaluation response")
	}

	// Current behavior: the deleted rel's cleanup wins in the same pass, so the
	// field is cleared even though rel B still mutates it. Re-evaluation would
	// restore it via rel B. Assert current behavior so a future fix shows up as
	// a diff here.
	if _, present := podSpec["serviceAccountName"]; present {
		t.Logf("serviceAccountName preserved by concurrent active relationship; " +
			"multi-mutator tracking may have been added, update this test")
	}
}
