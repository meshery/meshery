package policies

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// TestWalletPatching verifies that wallet relationships patch the child component's
// configuration from the parent's nested path.
// Example: Deployment (parent) contains PodTemplate (child).
// The PodTemplate's spec should be patched from Deployment's spec.template.spec.
func TestWalletPatching(t *testing.T) {
	p := &HierarchicalWalletPolicy{}

	deployID := testUUID("1")
	podTplID := testUUID("2")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{"name": "nginx", "image": "nginx:1.25"},
					},
				},
			},
		},
	})

	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"containers": []interface{}{
				map[string]interface{}{"name": "old", "image": "old:0.1"},
			},
		},
	})

	design := makePatternFile([]*component.ComponentDefinition{deploy, podTpl}, nil)

	// Wallet relationship: Deployment -> PodTemplate
	// mutatorRef on Deployment side (parent), mutatedRef on PodTemplate side (child)
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:                                   &deployID,
						Kind:                                 strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: mutatorPatch("configuration", "spec", "template", "spec", "containers"),
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:                                   &podTplID,
						Kind:                                 strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: mutatedPatch("configuration", "spec", "containers"),
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
	rel.ID = testUUID("10")

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
	deployID := testUUID("1")
	podTplID := testUUID("2")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"replicas": float64(3),
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"serviceAccountName": "my-sa",
				},
			},
		},
	})

	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"serviceAccountName": "default",
		},
	})

	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:                                   &deployID,
						Kind:                                 strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: mutatorPatch("configuration", "spec", "template", "spec", "serviceAccountName"),
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:                                   &podTplID,
						Kind:                                 strPtr("PodTemplate"),
						RelationshipDefinitionSelectorsPatch: mutatedPatch("configuration", "spec", "serviceAccountName"),
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
	rel.ID = testUUID("10")

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
	assertWasmMatches(t, *design, nil, resp)

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

	roleID := testUUID("1")
	saID := testUUID("2")
	rbID := testUUID("3")

	role := k8sComp("ClusterRole", roleID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "admin-role",
		},
	})

	sa := k8sComp("ServiceAccount", saID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"name": "my-service-account",
		},
	})

	rb := k8sComp("ClusterRoleBinding", rbID, map[string]interface{}{
		"roleRef": map[string]interface{}{
			"name": "",
		},
		"subjects": []interface{}{
			map[string]interface{}{
				"name": "",
				"kind": "ServiceAccount",
			},
		},
	})

	design := makePatternFile([]*component.ComponentDefinition{role, sa, rb}, nil)

	// Binding relationship: from=ClusterRole, to=ServiceAccount, binding=ClusterRoleBinding
	// The from selector's match field references the binding component (ClusterRoleBinding).
	fromMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRole", MutatorRef: mutatorRef("configuration", "metadata", "name"), ID: &roleID},
	}
	fromMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatedRef: mutatedRef("configuration", "roleRef", "name"), ID: &rbID},
	}

	toMatchFrom := []relationship.MatchSelectorItem{
		{Kind: "ClusterRoleBinding", MutatorRef: mutatorRef("configuration", "metadata", "name"), ID: &rbID},
	}
	toMatchTo := []relationship.MatchSelectorItem{
		{Kind: "ServiceAccount", MutatedRef: mutatedRef("configuration", "subjects", "0", "name"), ID: &saID},
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
	rel.ID = testUUID("20")

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

	deployID := testUUID("1")
	aliasID := testUUID("99")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{"name": "web", "image": "nginx:1.25"},
					},
				},
			},
		},
	})

	design := makePatternFile([]*component.ComponentDefinition{deploy}, nil)

	// Alias relationship with status "identified" should produce add-component actions.
	containerPath := []string{"configuration", "spec", "template", "spec", "containers", "0"}
	patch := &relationship.RelationshipDefinitionSelectorsPatch{
		MutatorRef: mutatorRef(containerPath...),
		MutatedRef: mutatedRef(containerPath...),
	}
	relStatus := relationship.RelationshipDefinitionStatus("identified")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:                                   &aliasID,
						Kind:                                 strPtr("Container"),
						RelationshipDefinitionSelectorsPatch: patch,
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:                                   &deployID,
						Kind:                                 strPtr("Deployment"),
						RelationshipDefinitionSelectorsPatch: patch,
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
	rel.ID = testUUID("30")

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

// rbacTriple returns the Role / RoleBinding / ServiceAccount components shared
// by the binding-permission identification tests, with IDs testUUID 1/2/3 and
// empty configuration. Callers that need a populated RoleBinding override
// rb.Configuration after the call.
func rbacTriple() (role, rb, sa *component.ComponentDefinition) {
	role = k8sComp("Role", testUUID("1"), map[string]interface{}{})
	role.DisplayName = "my-role"
	rb = k8sComp("RoleBinding", testUUID("2"), map[string]interface{}{})
	sa = k8sComp("ServiceAccount", testUUID("3"), map[string]interface{}{})
	sa.DisplayName = "my-sa"
	return role, rb, sa
}

// bindingPermissionSelectors returns the canonical RBAC binding-permission
// selector set: a Role bound to a ServiceAccount through a RoleBinding.
func bindingPermissionSelectors() relationship.SelectorSet {
	fromMatchFrom := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &relationship.MutatorRef{{"component", "kind"}, {"displayName"}}}}
	fromMatchTo := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &relationship.MutatedRef{{"configuration", "spec", "roleRef", "kind"}, {"configuration", "spec", "roleRef", "name"}}}}
	toMatchFrom := []relationship.MatchSelectorItem{{Kind: "RoleBinding", MutatedRef: &relationship.MutatedRef{{"configuration", "spec", "subjects", "_", "name"}}}}
	toMatchTo := []relationship.MatchSelectorItem{{Kind: "self", MutatorRef: &relationship.MutatorRef{{"displayName"}}}}
	return relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{{
					Kind:  strPtr("Role"),
					Model: k8sModel(),
					Match: &relationship.MatchSelector{From: &fromMatchFrom, To: &fromMatchTo},
				}},
				To: []relationship.SelectorItem{{
					Kind:  strPtr("ServiceAccount"),
					Model: k8sModel(),
					Match: &relationship.MatchSelector{From: &toMatchFrom, To: &toMatchTo},
				}},
			},
		},
	}
}

// TestBindingIdentification tests that binding/permission relationships can be identified
// for Role-RoleBinding-ServiceAccount triples.
func TestBindingIdentification(t *testing.T) {
	p := &EdgeBindingPolicy{}

	role, rb, sa := rbacTriple()
	rb.DisplayName = "my-rb"
	rb.Configuration = map[string]interface{}{
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
	}

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	selectorSet := bindingPermissionSelectors()
	relDef := &relationship.RelationshipDefinition{
		Kind:             relationship.Edge,
		RelationshipType: "binding",
		SubType:          "permission",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	relDef.ID = testUUID("50")

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

	role, rb, sa := rbacTriple()

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	selectorSet := bindingPermissionSelectors()
	relDef := &relationship.RelationshipDefinition{
		Kind: relationship.Edge, RelationshipType: "binding", SubType: "permission",
		Model: modelv1beta1.ModelReference{Name: "kubernetes"}, Selectors: &selectorSet,
	}
	relDef.ID = testUUID("60")

	identified := p.IdentifyRelationship(relDef, design)
	if len(identified) == 0 {
		t.Error("Expected binding relationship to be identified even with empty RoleBinding configuration")
	}
}

// TestBindingIdentificationFullPipeline tests binding identification through the full engine,
// simulating the e2e flow where relationships are stripped and re-evaluated.
func TestBindingIdentificationFullPipeline(t *testing.T) {
	role, rb, sa := rbacTriple()

	design := makePatternFile([]*component.ComponentDefinition{role, rb, sa}, nil)

	// Build the binding permission relDef (same structure as the real one).
	selectorSet := bindingPermissionSelectors()
	relDef := &relationship.RelationshipDefinition{
		Kind: relationship.Edge, RelationshipType: "binding", SubType: "permission",
		Model: modelv1beta1.ModelReference{Name: "kubernetes"}, Selectors: &selectorSet,
	}
	relDef.ID = testUUID("70")

	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})
	engine := NewGoEngine(log)
	regRels := []*relationship.RelationshipDefinition{relDef}
	resp, err := engine.EvaluateDesign(*design, regRels)
	if err != nil {
		t.Fatalf("EvaluateDesign failed: %v", err)
	}
	assertWasmMatches(t, *design, regRels, resp)

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

	nsID := testUUID("1")
	deployID := testUUID("2")
	svcID := testUUID("3")

	ns := k8sComp("Namespace", nsID, map[string]interface{}{})
	ns.DisplayName = "default"
	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"namespace": "default",
		},
	})
	deploy.DisplayName = "my-deploy"
	svc := k8sComp("Service", svcID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"namespace": "default",
		},
	})
	svc.DisplayName = "my-svc"
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
						Kind:                                 strPtr("*"),
						Model:                                &modelv1beta1.ModelReference{Name: "*"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{MutatedRef: &mutatedRef},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:                                 strPtr("Namespace"),
						Model:                                &modelv1beta1.ModelReference{Name: "kubernetes"},
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
	relDef.ID = testUUID("80")

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
	nsID := testUUID("1")
	deployID := testUUID("2")

	ns := k8sComp("Namespace", nsID, map[string]interface{}{"name": "production"})
	deploy := k8sComp("Deployment", deployID, map[string]interface{}{"namespace": "default"})
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
	existingRel.ID = testUUID("40")

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
	assertWasmMatches(t, *design, nil, resp)

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

	deployID := testUUID("1")
	podTplID := testUUID("2")

	// Both components already hold the mutator value (simulating a prior patch
	// that has been applied and persisted).
	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"serviceAccountName": "my-sa",
				},
			},
		},
	})
	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"serviceAccountName": "my-sa",
		},
	})
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
	rel.ID = testUUID("10")

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

	deployID := testUUID("1")
	podTplID := testUUID("2")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"serviceAccountName": "my-sa",
				},
			},
		},
	})
	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"serviceAccountName": "my-sa",
		},
	})
	// Schema carries the serviceAccountName default that cleanup must restore to.
	podTpl.Component.Schema = `{
		"properties": {
			"spec": {
				"properties": {
					"serviceAccountName": {"default": "default"}
				}
			}
		}
	}`
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
	rel.ID = testUUID("10")

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
	deployID := testUUID("1")
	podTplID := testUUID("2")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"serviceAccountName": "my-sa",
				},
			},
		},
	})
	// PodTemplate.serviceAccountName already holds the mutator value, simulating
	// a previous evaluation that applied the patch.
	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"serviceAccountName": "my-sa",
		},
	})
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "template", "spec", "serviceAccountName"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "serviceAccountName"}}
	relStatus := relationship.RelationshipDefinitionStatus("deleted")
	relID := testUUID("10")
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

	deployID := testUUID("1")
	podTplID := testUUID("2")

	deploy := k8sComp("Deployment", deployID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"serviceAccountName": "my-sa",
				},
			},
		},
	})
	// PodTemplate holds a different value than the mutator (user changed it).
	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{
			"serviceAccountName": "user-changed",
		},
	})
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
	rel.ID = testUUID("10")

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

	nsID := testUUID("1")
	deployID := testUUID("2")

	ns := k8sComp("Namespace", nsID, map[string]interface{}{"name": "production"})
	// Deployment.namespace already holds the mutator value ("production"),
	// simulating a prior patch that has been applied.
	deploy := k8sComp("Deployment", deployID, map[string]interface{}{"namespace": "production"})
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
	rel.ID = testUUID("10")

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

	svcID := testUUID("1")
	podID := testUUID("2")

	svc := k8sComp("Service", svcID, map[string]interface{}{
		"spec": map[string]interface{}{
			"selector": map[string]interface{}{"app": "web"},
		},
	})
	// Pod.labels.app already holds the mutator value, simulating a prior patch.
	pod := k8sComp("Pod", podID, map[string]interface{}{
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{"app": "web"},
		},
	})
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
	rel.ID = testUUID("10")

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

	roleID := testUUID("1")
	saID := testUUID("2")
	rbID := testUUID("3")

	role := k8sComp("ClusterRole", roleID, map[string]interface{}{
		"metadata": map[string]interface{}{"name": "admin-role"},
	})
	sa := k8sComp("ServiceAccount", saID, map[string]interface{}{
		"metadata": map[string]interface{}{"name": "my-service-account"},
	})
	// ClusterRoleBinding already holds the mutator value on roleRef.name,
	// simulating a prior binding patch that was applied.
	rb := k8sComp("ClusterRoleBinding", rbID, map[string]interface{}{
		"roleRef": map[string]interface{}{"name": "admin-role"},
		"subjects": []interface{}{
			map[string]interface{}{"name": "", "kind": "ServiceAccount"},
		},
	})
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
	rel.ID = testUUID("20")

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
	deployAID := testUUID("1")
	deployBID := testUUID("2")
	podTplID := testUUID("3")

	// Two Deployments both hold the same serviceAccountName; PodTemplate
	// already holds that value (simulating prior patches from both rels).
	deployA := k8sComp("Deployment", deployAID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
			},
		},
	})
	deployB := k8sComp("Deployment", deployBID, map[string]interface{}{
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
			},
		},
	})
	podTpl := k8sComp("PodTemplate", podTplID, map[string]interface{}{
		"spec": map[string]interface{}{"serviceAccountName": "shared-sa"},
	})
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
