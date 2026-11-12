package policies

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// kubernetesInventoryRelationship returns the canonical k8s
// hierarchical-parent-inventory relationship: any namespaced kind in any
// model points to a Namespace via metadata.namespace; the new Namespace's
// displayName is set from that value.
func kubernetesInventoryRelationship() *relationship.RelationshipDefinition {
	mutatedRef := relationship.MutatedRef{{"configuration", "metadata", "namespace"}}
	mutatorRef := relationship.MutatorRef{{"displayName"}}
	wildcard := "*"
	namespace := "Namespace"
	return &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{
						Kind: &wildcard,
						Model: &modelv1beta1.ModelReference{
							Name: "*",
						},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					}},
					To: []relationship.SelectorItem{{
						Kind: &namespace,
						Model: &modelv1beta1.ModelReference{
							Name: "kubernetes",
						},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					}},
				},
			},
		},
	}
}

// crossModelInventoryRelationship returns a relationship modeled on the
// real azure-event-grid EventSubscription → azure-storage StorageAccount
// dependency: the parent (StorageAccount) is in a different model than
// the child (EventSubscription). Used to verify the auto-added parent
// follows the to-side model, not the child's model.
func crossModelInventoryRelationship() *relationship.RelationshipDefinition {
	mutatedRef := relationship.MutatedRef{{"configuration", "spec", "destination", "storageAccount"}}
	mutatorRef := relationship.MutatorRef{{"displayName"}}
	fromKind := "EventSubscription"
	toKind := "StorageAccount"
	return &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{
						Kind:  &fromKind,
						Model: &modelv1beta1.ModelReference{Name: "azure-event-grid"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatedRef: &mutatedRef,
						},
					}},
					To: []relationship.SelectorItem{{
						Kind:  &toKind,
						Model: &modelv1beta1.ModelReference{Name: "azure-storage"},
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					}},
				},
			},
		},
	}
}

// makePodWithNamespace builds a single-Pod design with metadata.namespace
// set to the given value — the canonical bug-fixture: a Kubernetes resource
// referencing a Namespace that does not yet exist.
func makePodWithNamespace(t *testing.T, ns string) (*pattern.PatternFile, *component.ComponentDefinition) {
	t.Helper()
	id, err := uuid.NewV4()
	require.NoError(t, err)
	pod := &component.ComponentDefinition{
		ID:        id,
		Component: component.Component{Kind: "Pod", Version: "v1"},
		Model: &modelv1beta1.ModelDefinition{
			Name:    "kubernetes",
			Version: "v1.25.0",
			Model:   modelv1beta1.Model{Version: "v1.25.0"},
		},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": ns,
			},
		},
	}
	return makePatternFile([]*component.ComponentDefinition{pod}, nil), pod
}

// Bug regression — Pod with metadata.namespace and no Namespace component:
// the Go engine must auto-emit an add_component action for a Namespace
// in the kubernetes model, displayName matching the referenced value.
// Without identifyInventoryAdditions wired into engine.evaluate, this
// returns nothing — the parity bug behind issue #19090.
func TestIdentifyInventoryAdditions_PodImpliesNamespace(t *testing.T) {
	t.Parallel()
	design, _ := makePodWithNamespace(t, "default")
	rel := kubernetesInventoryRelationship()

	actions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{rel})

	require.Len(t, actions, 1, "exactly one Namespace must be auto-added for one Pod with one missing parent")
	a := actions[0]
	assert.Equal(t, AddComponentOp, a.Op)
	require.NotNil(t, a.Component)
	assert.Equal(t, "Namespace", a.Component.Component.Kind, "auto-added kind must be Namespace, not the from-side wildcard")
	require.NotNil(t, a.Component.Model, "Model must be populated so processEvaluationResponse can hydrate styles")
	assert.Equal(t, "kubernetes", a.Component.Model.Name, "Model.Name must follow the to-side selector, not the wildcard")
	assert.Equal(t, "kubernetes", a.Component.ModelReference.Name)
	assert.Equal(t, "default", string(a.Component.DisplayName), "DisplayName must come from the mutator-path patch (mutatorRef=displayName ↔ mutatedValue='default')")
}

// If a Namespace named "default" already exists, no add_component action
// must be emitted — the rego rule's `every` guard avoids creating a
// duplicate parent.
func TestIdentifyInventoryAdditions_SkipsWhenParentAlreadyPresent(t *testing.T) {
	t.Parallel()
	design, pod := makePodWithNamespace(t, "default")
	nsID, err := uuid.NewV4()
	require.NoError(t, err)
	ns := &component.ComponentDefinition{
		ID:             nsID,
		Component:      component.Component{Kind: "Namespace", Version: "v1"},
		DisplayName:    "default",
		Model:          &modelv1beta1.ModelDefinition{Name: "kubernetes"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
	}
	design.Components = append(design.Components, ns)
	_ = pod

	actions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{kubernetesInventoryRelationship()})
	assert.Empty(t, actions, "must not emit add_component when a matching parent already exists in the design")
}

// makeDeploymentInNamespace builds a Deployment with metadata.namespace set
// to the given value, for designs where several children share one parent.
func makeDeploymentInNamespace(t *testing.T, ns string) *component.ComponentDefinition {
	t.Helper()
	id, err := uuid.NewV4()
	require.NoError(t, err)
	return &component.ComponentDefinition{
		ID:             id,
		Component:      component.Component{Kind: "Deployment", Version: "apps/v1"},
		Model:          &modelv1beta1.ModelDefinition{Name: "kubernetes"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": ns,
			},
		},
	}
}

// Issue #22778 regression: several Deployments in one namespace, evaluated
// against several Namespace inventory relationships (the registry ships one
// per kubernetes model version), must add exactly one Namespace. Each
// candidate gets a fresh ID, so deduplicating by ID added one per pair.
func TestIdentifyInventoryAdditions_SharedNamespaceAddedOnce(t *testing.T) {
	t.Parallel()
	design := makePatternFile([]*component.ComponentDefinition{
		makeDeploymentInNamespace(t, "web"),
		makeDeploymentInNamespace(t, "web"),
		makeDeploymentInNamespace(t, "web"),
	}, nil)
	rels := []*relationship.RelationshipDefinition{
		kubernetesInventoryRelationship(),
		kubernetesInventoryRelationship(),
	}

	actions := identifyInventoryAdditions(design, rels)

	require.Len(t, actions, 1, "one shared namespace must produce exactly one auto-added Namespace")
	require.NotNil(t, actions[0].Component)
	assert.Equal(t, "Namespace", actions[0].Component.Component.Kind)
	assert.Equal(t, "web", string(actions[0].Component.DisplayName))
}

// Deduplication must not collapse distinct parents: Deployments in two
// namespaces still get one Namespace each.
func TestIdentifyInventoryAdditions_DistinctNamespacesAddedSeparately(t *testing.T) {
	t.Parallel()
	design := makePatternFile([]*component.ComponentDefinition{
		makeDeploymentInNamespace(t, "web"),
		makeDeploymentInNamespace(t, "web"),
		makeDeploymentInNamespace(t, "db"),
	}, nil)

	actions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{kubernetesInventoryRelationship()})

	require.Len(t, actions, 2, "each distinct namespace must produce its own Namespace")
	names := []string{}
	for _, a := range actions {
		require.NotNil(t, a.Component)
		assert.Equal(t, "Namespace", a.Component.Component.Kind)
		names = append(names, string(a.Component.DisplayName))
	}
	assert.ElementsMatch(t, []string{"web", "db"}, names)
}

// makeComponentInNamespace builds a namespaced component of the given kind
// with metadata.namespace set, for designs where children of different
// kinds share one parent.
func makeComponentInNamespace(t *testing.T, kind, apiVersion, ns string) *component.ComponentDefinition {
	t.Helper()
	id, err := uuid.NewV4()
	require.NoError(t, err)
	return &component.ComponentDefinition{
		ID:             id,
		Component:      component.Component{Kind: kind, Version: apiVersion},
		Model:          &modelv1beta1.ModelDefinition{Name: "kubernetes"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"metadata": map[string]interface{}{
				"namespace": ns,
			},
		},
	}
}

// kubernetesModelVersions mirrors the registry, which ships one copy of
// hierarchical-parent-inventory-kmjea.json per kubernetes model version.
var kubernetesModelVersions = []string{
	"v1.25.0", "v1.26.0", "v1.27.0", "v1.28.0", "v1.29.0", "v1.30.0", "v1.31.0",
	"v1.32.0", "v1.33.0", "v1.34.0", "v1.35.0-alpha.1", "v1.35.0-alpha.2",
}

// kubernetesInventoryRelationshipsAcrossModelVersions returns one copy of the
// Namespace inventory relationship per kubernetes model version, stamped the
// way the shipped JSONs are: the relationship's own model carries the version
// while the to-side (Namespace) selector's model version stays empty.
func kubernetesInventoryRelationshipsAcrossModelVersions(versions []string) []*relationship.RelationshipDefinition {
	rels := make([]*relationship.RelationshipDefinition, 0, len(versions))
	for _, v := range versions {
		rel := kubernetesInventoryRelationship()
		rel.Model.Name = "kubernetes"
		rel.Model.Model.Version = core.SemverString(v)
		rels = append(rels, rel)
	}
	return rels
}

// kubernetesInventoryRelationshipsWithVersionedParents is the hypothetical
// shape where each copy also pins its model version on the to-side Namespace
// selector, so the version flows through buildInventoryParentCandidate into
// the candidate's ModelReference.
func kubernetesInventoryRelationshipsWithVersionedParents(versions []string) []*relationship.RelationshipDefinition {
	rels := kubernetesInventoryRelationshipsAcrossModelVersions(versions)
	for i, rel := range rels {
		for si := range *rel.Selectors {
			for ti := range (*rel.Selectors)[si].Allow.To {
				to := &(*rel.Selectors)[si].Allow.To[ti]
				if to.Model == nil {
					continue
				}
				to.Model.Version = core.SemverString(versions[i])
				to.Model.Model.Version = core.SemverString(versions[i])
			}
		}
	}
	return rels
}

// Issue #23360 regression: a Deployment and a Service in `default`, evaluated
// against the registry's per-kubernetes-version copies of the Namespace
// inventory relationship, must add exactly one Namespace. Each copy differs
// only in its own model version; the to-side selector's model version is
// empty, so every copy builds the same candidate and the key collapses them.
func TestIdentifyInventoryAdditions_OneNamespaceAcrossKubernetesModelVersions(t *testing.T) {
	t.Parallel()
	design := makePatternFile([]*component.ComponentDefinition{
		makeComponentInNamespace(t, "Deployment", "apps/v1", "default"),
		makeComponentInNamespace(t, "Service", "v1", "default"),
	}, nil)
	rels := kubernetesInventoryRelationshipsAcrossModelVersions(kubernetesModelVersions)
	require.Len(t, rels, len(kubernetesModelVersions))

	actions := identifyInventoryAdditions(design, rels)

	require.Len(t, actions, 1, "one missing namespace must produce exactly one add-Namespace action across all model-version copies")
	require.NotNil(t, actions[0].Component)
	assert.Equal(t, "Namespace", actions[0].Component.Component.Kind)
	assert.Equal(t, "default", string(actions[0].Component.DisplayName))
	assert.Equal(t, "kubernetes", actions[0].Component.ModelReference.Name)
}

// Collapsing across model-version copies must not collapse across
// namespaces: children in two namespaces still get one Namespace each.
func TestIdentifyInventoryAdditions_DistinctNamespacesAcrossKubernetesModelVersions(t *testing.T) {
	t.Parallel()
	design := makePatternFile([]*component.ComponentDefinition{
		makeComponentInNamespace(t, "Deployment", "apps/v1", "web"),
		makeComponentInNamespace(t, "Service", "v1", "web"),
		makeComponentInNamespace(t, "Deployment", "apps/v1", "db"),
	}, nil)
	rels := kubernetesInventoryRelationshipsAcrossModelVersions(kubernetesModelVersions)

	actions := identifyInventoryAdditions(design, rels)

	require.Len(t, actions, 2, "each distinct namespace must produce its own Namespace, once")
	names := []string{}
	for _, a := range actions {
		require.NotNil(t, a.Component)
		assert.Equal(t, "Namespace", a.Component.Component.Kind)
		names = append(names, string(a.Component.DisplayName))
	}
	assert.ElementsMatch(t, []string{"web", "db"}, names)
}

// The key carries the candidate's full ModelReference, not just its name:
// if a to-side selector pins a distinct model version, that version reaches
// the candidate and the candidates must stay distinct rather than being
// merged across versions.
func TestIdentifyInventoryAdditions_VersionedParentSelectorsStayDistinct(t *testing.T) {
	t.Parallel()
	versions := []string{"v1.33.0", "v1.34.0", "v1.35.0-alpha.2"}
	design := makePatternFile([]*component.ComponentDefinition{
		makeComponentInNamespace(t, "Deployment", "apps/v1", "default"),
		makeComponentInNamespace(t, "Service", "v1", "default"),
	}, nil)
	rels := kubernetesInventoryRelationshipsWithVersionedParents(versions)

	actions := identifyInventoryAdditions(design, rels)

	require.Len(t, actions, len(versions), "a to-side selector pinning a distinct model version must yield a distinct Namespace candidate")
	gotVersions := []string{}
	for _, a := range actions {
		require.NotNil(t, a.Component)
		assert.Equal(t, "Namespace", a.Component.Component.Kind)
		assert.Equal(t, "default", string(a.Component.DisplayName))
		gotVersions = append(gotVersions, string(a.Component.ModelReference.Version))
	}
	assert.ElementsMatch(t, versions, gotVersions)
}

// inventoryCandidateKey must ignore the random component ID and key on the
// candidate's kind, full ModelReference, displayName and configuration.
func TestInventoryCandidateKey(t *testing.T) {
	t.Parallel()
	build := func(t *testing.T, modelVersion, displayName string) *component.ComponentDefinition {
		t.Helper()
		id, err := uuid.NewV4()
		require.NoError(t, err)
		return &component.ComponentDefinition{
			ID:          id,
			Component:   component.Component{Kind: "Namespace"},
			DisplayName: core.InputString(displayName),
			ModelReference: modelv1beta1.ModelReference{
				Name:    "kubernetes",
				Version: core.SemverString(modelVersion),
			},
		}
	}

	a := build(t, "", "default")
	b := build(t, "", "default")
	require.NotEqual(t, a.ID, b.ID)
	assert.Equal(t, inventoryCandidateKey(a), inventoryCandidateKey(b), "the random ID must not be part of the key")

	assert.NotEqual(t, inventoryCandidateKey(a), inventoryCandidateKey(build(t, "", "web")), "displayName must be part of the key")
	assert.NotEqual(t, inventoryCandidateKey(a), inventoryCandidateKey(build(t, "v1.34.0", "default")), "the ModelReference version must be part of the key")

	withConfig := build(t, "", "default")
	withConfig.Configuration = map[string]interface{}{"metadata": map[string]interface{}{"labels": map[string]interface{}{"team": "web"}}}
	assert.NotEqual(t, inventoryCandidateKey(a), inventoryCandidateKey(withConfig), "configuration must be part of the key")

	// A configuration JSON cannot encode must still produce a deterministic
	// key that distinguishes candidates rather than colliding on an error.
	nan := build(t, "", "default")
	nan.Configuration = map[string]interface{}{"weight": math.NaN()}
	nan2 := build(t, "", "default")
	nan2.Configuration = map[string]interface{}{"weight": math.NaN()}
	assert.Equal(t, inventoryCandidateKey(nan), inventoryCandidateKey(nan2), "fallback key must be deterministic")
	assert.NotEqual(t, inventoryCandidateKey(nan), inventoryCandidateKey(a), "fallback key must not collide with the JSON key of a different candidate")
}

// Cross-model regression: if the relationship has a different model on the
// to-side than the from-side (e.g. azure-event-grid EventSubscription →
// azure-storage StorageAccount), the auto-added parent must inherit the
// to-side model. This is the exact case where my originally-suggested
// "use mutated_component.model" Rego edit would have broken behavior; this
// test pins the correct contract.
func TestIdentifyInventoryAdditions_PreservesToSideModelForCrossModelRelationships(t *testing.T) {
	t.Parallel()
	id, err := uuid.NewV4()
	require.NoError(t, err)
	sub := &component.ComponentDefinition{
		ID:        id,
		Component: component.Component{Kind: "EventSubscription", Version: "v1"},
		Model: &modelv1beta1.ModelDefinition{
			Name: "azure-event-grid",
		},
		ModelReference: modelv1beta1.ModelReference{Name: "azure-event-grid"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"destination": map[string]interface{}{
					"storageAccount": "saproduction",
				},
			},
		},
	}
	design := makePatternFile([]*component.ComponentDefinition{sub}, nil)

	actions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{crossModelInventoryRelationship()})
	require.Len(t, actions, 1)
	added := actions[0].Component
	require.NotNil(t, added)
	assert.Equal(t, "StorageAccount", added.Component.Kind)
	require.NotNil(t, added.Model)
	assert.Equal(t, "azure-storage", added.Model.Name,
		"cross-model parent must follow the to-side model (azure-storage), not the from-side (azure-event-grid)")
	assert.Equal(t, "saproduction", string(added.DisplayName))
}

// Non-inventory relationships (alias, edge, etc.) must be ignored. Without
// this guard, a stray inventoryActions side-effect could fire for any
// hierarchical relationship.
func TestIdentifyInventoryAdditions_IgnoresNonInventoryRelationships(t *testing.T) {
	t.Parallel()
	design, _ := makePodWithNamespace(t, "default")
	notInventory := kubernetesInventoryRelationship()
	notInventory.SubType = "alias"

	actions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{notInventory})
	assert.Empty(t, actions, "alias relationships must not produce inventory add_component actions")
}

// Cross-engine parity — runs the same Pod-references-default-Namespace
// fixture through BOTH the Rego engine (data.relationship_evaluation_policy.identify_additions)
// AND the Go engine (identifyInventoryAdditions), then asserts the emitted
// new component carries the same Kind, Model.Name, and DisplayName.
//
// Pins the behavioral contract that USE_GO_POLICY_ENGINE produces the same
// auto-added component as the OPA path. Without this test, future
// divergence (intentional or accidental) would only surface as user-visible
// bugs once again.
func TestParity_RegoAndGoEngineEmitEquivalentInventoryAdditions(t *testing.T) {
	t.Parallel()

	// --- Go engine ---
	design, _ := makePodWithNamespace(t, "default")
	rel := kubernetesInventoryRelationship()

	goActions := identifyInventoryAdditions(design, []*relationship.RelationshipDefinition{rel})
	require.Len(t, goActions, 1, "Go engine must emit one add_component for one missing parent")
	require.NotNil(t, goActions[0].Component)
	goNew := goActions[0].Component

	// --- Rego engine ---
	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"
	files, err := collectRegoFiles(policiesDir)
	require.NoError(t, err)
	var modules []func(*rego.Rego)
	for _, file := range files {
		if strings.Contains(file, "/tests/") || strings.HasSuffix(file, ".template") {
			continue
		}
		content, readErr := os.ReadFile(file)
		require.NoError(t, readErr)
		modules = append(modules, rego.Module(file, string(content)))
	}

	// Mirror the rego's expected input shape — selector JSON keys match
	// the kubernetes hierarchical-parent-inventory relationship JSON.
	regoInput := map[string]interface{}{
		"design_file": map[string]interface{}{
			"components": []interface{}{
				map[string]interface{}{
					"id":        "pod-1",
					"component": map[string]interface{}{"kind": "Pod"},
					"model": map[string]interface{}{
						"name":    "kubernetes",
						"version": "v1.25.0",
						"model":   map[string]interface{}{"version": "v1.25.0"},
					},
					"configuration": map[string]interface{}{
						"metadata": map[string]interface{}{"namespace": "default"},
					},
				},
			},
			"relationships": []interface{}{},
		},
		"relationship": map[string]interface{}{
			"kind":    "hierarchical",
			"type":    "parent",
			"subType": "inventory",
			"selectors": []interface{}{
				map[string]interface{}{
					"allow": map[string]interface{}{
						"from": []interface{}{
							map[string]interface{}{
								"kind":  "*",
								"model": map[string]interface{}{"name": "*", "version": "", "model": map[string]interface{}{"version": ""}},
								"patch": map[string]interface{}{
									"mutatedRef": []interface{}{[]interface{}{"configuration", "metadata", "namespace"}},
								},
							},
						},
						"to": []interface{}{
							map[string]interface{}{
								"kind":  "Namespace",
								"model": map[string]interface{}{"name": "kubernetes", "version": "", "model": map[string]interface{}{"version": ""}},
								"patch": map[string]interface{}{
									"mutatorRef": []interface{}{[]interface{}{"displayName"}},
								},
							},
						},
					},
					"deny": map[string]interface{}{"from": []interface{}{}, "to": []interface{}{}},
				},
			},
		},
	}

	opts := append(modules,
		rego.Query("data.relationship_evaluation_policy.identify_additions(input.design_file, input.relationship)"),
		rego.Input(regoInput),
	)

	rs, err := rego.New(opts...).Eval(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, rs)
	require.NotEmpty(t, rs[0].Expressions)

	results, ok := rs[0].Expressions[0].Value.([]interface{})
	require.True(t, ok)
	require.Len(t, results, 1, "rego engine must emit one add_component for one missing parent")

	regoNew, ok := results[0].(map[string]interface{})
	require.True(t, ok)

	// --- Compare ---
	regoKind := regoNew["component"].(map[string]interface{})["kind"]
	regoModel := regoNew["model"].(map[string]interface{})["name"]
	regoDisplay := regoNew["displayName"]

	assert.Equal(t, "Namespace", goNew.Component.Kind)
	assert.Equal(t, regoKind, goNew.Component.Kind, "Go and Rego must agree on the auto-added Kind")

	require.NotNil(t, goNew.Model)
	assert.Equal(t, "kubernetes", goNew.Model.Name)
	assert.Equal(t, regoModel, goNew.Model.Name, "Go and Rego must agree on the auto-added Model.Name")

	assert.Equal(t, "default", string(goNew.DisplayName))
	assert.Equal(t, regoDisplay, string(goNew.DisplayName), "Go and Rego must agree on the auto-added DisplayName")

	// Sanity: full JSON serialization of the relevant fields. Catches
	// future field divergences that escape the explicit assertions above.
	type comparisonShape struct {
		Kind        string `json:"kind"`
		ModelName   string `json:"modelName"`
		DisplayName string `json:"displayName"`
	}
	goShape := comparisonShape{
		Kind:        goNew.Component.Kind,
		ModelName:   goNew.Model.Name,
		DisplayName: string(goNew.DisplayName),
	}
	regoShape := comparisonShape{
		Kind:        regoKind.(string),
		ModelName:   regoModel.(string),
		DisplayName: regoDisplay.(string),
	}
	gj, _ := json.Marshal(goShape)
	rj, _ := json.Marshal(regoShape)
	assert.Equal(t, string(rj), string(gj), "Go and Rego engine outputs diverge — auto-added component shape mismatch")
}
