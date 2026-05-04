package policies

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/storage/inmem"
)

// benchDesign returns a realistic design with components and relationships as a map
// (for OPA benchmarking and alignment tests).
func benchDesign() map[string]interface{} {
	return map[string]interface{}{
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
					"replicas":  float64(3),
				},
			},
			map[string]interface{}{
				"id":        "svc-1",
				"component": map[string]interface{}{"kind": "Service"},
				"model":     map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{
					"namespace": "default",
					"port":      float64(80),
				},
			},
		},
		"relationships": []interface{}{
			map[string]interface{}{
				"id":      "rel-1",
				"kind":    "hierarchical",
				"type":    "parent",
				"subType": "inventory",
				"status":  "approved",
				"model":   map[string]interface{}{"name": "kubernetes"},
				"selectors": []interface{}{
					map[string]interface{}{
						"allow": map[string]interface{}{
							"from": []interface{}{map[string]interface{}{
								"id":   "ns-1",
								"kind": "Namespace",
								"patch": map[string]interface{}{
									"mutatorRef": []interface{}{
										[]interface{}{"configuration", "name"},
									},
								},
							}},
							"to": []interface{}{map[string]interface{}{
								"id":   "deploy-1",
								"kind": "Deployment",
								"patch": map[string]interface{}{
									"mutatedRef": []interface{}{
										[]interface{}{"configuration", "namespace"},
									},
								},
							}},
						},
					},
				},
			},
		},
		"preferences": map[string]interface{}{},
	}
}

// benchRelDefs returns relationship definitions for benchmarking as maps
// (for OPA benchmarking).
func benchRelDefs() []map[string]interface{} {
	return []map[string]interface{}{
		{
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
		},
	}
}

// benchTypedDesign returns the bench design as typed structs for Go engine functions.
func benchTypedDesign() *pattern.PatternFile {
	ns, deploy, svc := benchTypedComponents()
	rel := benchTypedRelationship(ns, deploy)
	return &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{ns, deploy, svc},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}
}

// benchTypedComponents returns the three typed components used in benchmarks.
func benchTypedComponents() (ns, deploy, svc *component.ComponentDefinition) {
	ns = &component.ComponentDefinition{
		Component:      component.Component{Kind: "Namespace"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"name": "default"},
	}
	// Use deterministic UUIDs that correspond to the map-based IDs above.
	ns.ID = staticUUID("ns-1")

	deploy = &component.ComponentDefinition{
		Component:      component.Component{Kind: "Deployment"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "default", "replicas": float64(3)},
	}
	deploy.ID = staticUUID("deploy-1")

	svc = &component.ComponentDefinition{
		Component:      component.Component{Kind: "Service"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration:  map[string]interface{}{"namespace": "default", "port": float64(80)},
	}
	svc.ID = staticUUID("svc-1")

	return ns, deploy, svc
}


// benchTypedRelationship builds the approved hierarchical relationship between ns and deploy.
func benchTypedRelationship(ns, deploy *component.ComponentDefinition) *relationship.RelationshipDefinition {
	mutatorRef := relationship.MutatorRef{[]string{"configuration", "name"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "namespace"}}
	relStatus := relationship.RelationshipDefinitionStatus("approved")
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						ID:   &ns.ID,
						Kind: strPtr("Namespace"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						ID:   &deploy.ID,
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
		Kind:             relationship.RelationshipDefinitionKind("hierarchical"),
		RelationshipType: "parent",
		SubType:          "inventory",
		Status:           &relStatus,
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	rel.ID, _ = uuidFromString("rel-1")
	return rel
}

// benchTypedRelDefs returns the relationship definitions for benchmarking as typed structs.
func benchTypedRelDefs() []*relationship.RelationshipDefinition {
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
	relDef.ID, _ = uuidFromString("rel-def-1")
	return []*relationship.RelationshipDefinition{relDef}
}

// loadOPAModules loads all rego policy modules (excluding tests).
func loadOPAModules(tb testing.TB) []func(*rego.Rego) {
	policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"
	policyFiles, err := collectRegoFiles(policiesDir)
	if err != nil {
		tb.Fatalf("Failed to collect policy files: %v", err)
	}

	var modules []func(*rego.Rego)
	for _, file := range policyFiles {
		if strings.Contains(file, "/tests/") {
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			tb.Fatalf("Failed to read file %s: %v", file, err)
		}
		modules = append(modules, rego.Module(file, string(content)))
	}
	return modules
}

// prepareOPAQuery compiles and prepares an OPA query with the given design and data.
func prepareOPAQuery(tb testing.TB, query string, design map[string]interface{}, data map[string]interface{}) rego.PreparedEvalQuery {
	modules := loadOPAModules(tb)
	ctx := context.Background()

	store := inmem.NewFromObject(data)

	opts := append(modules,
		rego.Query(query),
		rego.Input(design),
		rego.Store(store),
	)

	r := rego.New(opts...)
	pq, err := r.PrepareForEval(ctx)
	if err != nil {
		tb.Fatalf("OPA prepare failed for query %q: %v", query, err)
	}
	return pq
}

// TestBehaviorAlignment verifies that Go and OPA produce equivalent results.
func TestBehaviorAlignment(t *testing.T) {
	design := benchDesign()
	ctx := context.Background()

	t.Run("match_values", func(t *testing.T) {
		// Go
		goResult := matchValues("default", "default", "equal")
		if !goResult {
			t.Fatal("Go: matchValues should return true for equal values")
		}

		// OPA
		pq := prepareOPAQuery(t,
			`data.eval_rules.match_values("default", "default", "equal")`,
			design, map[string]interface{}{})
		rs, err := pq.Eval(ctx)
		if err != nil {
			t.Fatalf("OPA eval error: %v", err)
		}
		if len(rs) == 0 || len(rs[0].Expressions) == 0 {
			t.Fatal("OPA: match_values returned no result")
		}
		if rs[0].Expressions[0].Value != true {
			t.Fatalf("OPA: match_values returned %v, want true", rs[0].Expressions[0].Value)
		}
	})

	t.Run("from_and_to_components_exist", func(t *testing.T) {
		typedDesign := benchTypedDesign()
		rel := typedDesign.Relationships[0]

		// Go
		goResult := fromAndToComponentsExist(rel, typedDesign)
		if !goResult {
			t.Fatal("Go: fromAndToComponentsExist should return true")
		}

		// OPA: use the map-based design for OPA input
		pq := prepareOPAQuery(t,
			`data.eval_rules.from_and_to_components_exist(input.relationships[0], input)`,
			design, map[string]interface{}{})
		rs, err := pq.Eval(ctx)
		if err != nil {
			t.Fatalf("OPA eval error: %v", err)
		}
		if len(rs) == 0 || len(rs[0].Expressions) == 0 {
			t.Fatal("OPA: from_and_to_components_exist returned no result")
		}
		if rs[0].Expressions[0].Value != true {
			t.Fatalf("OPA returned %v, want true", rs[0].Expressions[0].Value)
		}
	})

	t.Run("from_and_to_missing_component", func(t *testing.T) {
		missingDesign := deepCopyMap(design)
		missingDesign["components"] = []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "default"},
			},
		}

		// Build typed design with only ns component but keeping the relationship
		// The relationship still references deploy (which is absent from the design).
		ns, deploy, _ := benchTypedComponents()
		rel := benchTypedRelationship(ns, deploy)
		typedMissingDesign := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{ns},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}

		// Go
		goResult := fromAndToComponentsExist(rel, typedMissingDesign)
		if goResult {
			t.Fatal("Go: fromAndToComponentsExist should return false when deploy-1 is missing")
		}

		// OPA: use map-based design for OPA input
		pq := prepareOPAQuery(t,
			`data.eval_rules.from_or_to_components_dont_exist(input.relationships[0], input)`,
			missingDesign, map[string]interface{}{})
		rs, err := pq.Eval(ctx)
		if err != nil {
			t.Fatalf("OPA eval error: %v", err)
		}
		if len(rs) == 0 || len(rs[0].Expressions) == 0 {
			t.Fatal("OPA: from_or_to_components_dont_exist returned no result")
		}
		if rs[0].Expressions[0].Value != true {
			t.Fatalf("OPA returned %v, want true", rs[0].Expressions[0].Value)
		}
	})

	t.Run("validation_invalidates_orphaned_relationship", func(t *testing.T) {
		orphanDesign := deepCopyMap(design)
		orphanDesign["components"] = []interface{}{
			map[string]interface{}{
				"id":            "ns-1",
				"component":     map[string]interface{}{"kind": "Namespace"},
				"model":         map[string]interface{}{"name": "kubernetes"},
				"configuration": map[string]interface{}{"name": "default"},
			},
		}

		// Build typed orphan design: only ns, but relationship references deploy (missing).
		ns, deploy, _ := benchTypedComponents()
		rel := benchTypedRelationship(ns, deploy)
		typedOrphanDesign := &pattern.PatternFile{
			Components:    []*component.ComponentDefinition{ns},
			Relationships: []*relationship.RelationshipDefinition{rel},
		}

		// Go: validation should produce a delete action for the orphaned relationship
		policy := &HierarchicalParentChildPolicy{}
		goActions := validateRelationshipsInDesign(typedOrphanDesign, policy)
		if len(goActions) != 1 {
			t.Fatalf("Go: expected 1 validation action, got %d", len(goActions))
		}
		if goActions[0].Op != UpdateRelationshipOp {
			t.Fatalf("Go: expected %s, got %s", UpdateRelationshipOp, goActions[0].Op)
		}
		if goActions[0].StringValue != "deleted" {
			t.Fatalf("Go: expected status deleted, got %v", goActions[0].StringValue)
		}

		// OPA: validate_relationships_in_design should produce equivalent action
		pq := prepareOPAQuery(t,
			`data.eval.validate_relationships_in_design(input, "hierarchical_parent_child")`,
			orphanDesign, map[string]interface{}{})
		rs, err := pq.Eval(ctx)
		if err != nil {
			t.Fatalf("OPA eval error: %v", err)
		}
		if len(rs) == 0 || len(rs[0].Expressions) == 0 {
			t.Fatal("OPA: validate returned no result")
		}

		opaResult, ok := rs[0].Expressions[0].Value.([]interface{})
		if !ok {
			t.Fatalf("OPA: expected set/array result, got %T", rs[0].Expressions[0].Value)
		}
		if len(opaResult) != 1 {
			t.Fatalf("OPA: expected 1 validation action, got %d", len(opaResult))
		}
		action, ok := opaResult[0].(map[string]interface{})
		if !ok {
			t.Fatalf("OPA: expected map action, got %T", opaResult[0])
		}
		if action["op"] != UpdateRelationshipOp {
			t.Fatalf("OPA: expected op %s, got %v", UpdateRelationshipOp, action["op"])
		}
		val, ok := action["value"].(map[string]interface{})
		if !ok || val["value"] != "deleted" {
			t.Fatalf("OPA: expected status deleted, got %v", val)
		}
	})

	t.Run("identify_relationships", func(t *testing.T) {
		typedDesign := benchTypedDesign()
		typedRelDefs := benchTypedRelDefs()
		relDefs := benchRelDefs()

		// Go
		policy := &HierarchicalParentChildPolicy{}
		goActions := identifyRelationshipsInDesign(typedDesign, typedRelDefs, policy)

		// OPA: identify_relationships_in_design needs relationships as data
		opaData := map[string]interface{}{
			"relationships": []interface{}{relDefs[0]},
		}
		pq := prepareOPAQuery(t,
			`data.eval.identify_relationships_in_design(input, data.relationships, "hierarchical_parent_child")`,
			design, opaData)
		rs, err := pq.Eval(ctx)
		if err != nil {
			t.Fatalf("OPA eval error: %v", err)
		}

		opaCount := 0
		if len(rs) > 0 && len(rs[0].Expressions) > 0 {
			if arr, ok := rs[0].Expressions[0].Value.([]interface{}); ok {
				opaCount = len(arr)
			}
		}

		// Both should find 0 new relationships (rel-1 already exists in design)
		t.Logf("Go identified %d new relationships, OPA identified %d", len(goActions), opaCount)
		if len(goActions) != opaCount {
			t.Errorf("Mismatch: Go=%d, OPA=%d", len(goActions), opaCount)
		}
	})
}

// --- Benchmarks: same operation, both engines ---

// BenchmarkValidation/Go vs BenchmarkValidation/OPA
func BenchmarkValidation(b *testing.B) {
	design := benchDesign()

	// Design with orphaned relationship (deploy-1 removed)
	orphanDesign := deepCopyMap(design)
	orphanDesign["components"] = []interface{}{
		map[string]interface{}{
			"id":            "ns-1",
			"component":     map[string]interface{}{"kind": "Namespace"},
			"model":         map[string]interface{}{"name": "kubernetes"},
			"configuration": map[string]interface{}{"name": "default"},
		},
	}

	// Typed orphan design for Go engine: relationship still references deploy (absent).
	ns, deploy, _ := benchTypedComponents()
	rel := benchTypedRelationship(ns, deploy)
	typedOrphanDesign := &pattern.PatternFile{
		Components:    []*component.ComponentDefinition{ns},
		Relationships: []*relationship.RelationshipDefinition{rel},
	}

	b.Run("Go", func(b *testing.B) {
		policy := &HierarchicalParentChildPolicy{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			validateRelationshipsInDesign(typedOrphanDesign, policy)
		}
	})

	b.Run("OPA", func(b *testing.B) {
		ctx := context.Background()
		pq := prepareOPAQuery(b,
			`data.eval.validate_relationships_in_design(input, "hierarchical_parent_child")`,
			orphanDesign, map[string]interface{}{})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = pq.Eval(ctx)
		}
	})
}

// BenchmarkIdentify/Go vs BenchmarkIdentify/OPA
func BenchmarkIdentify(b *testing.B) {
	design := benchDesign()
	typedDesign := benchTypedDesign()
	typedRelDefs := benchTypedRelDefs()
	relDefs := benchRelDefs()

	b.Run("Go", func(b *testing.B) {
		policy := &HierarchicalParentChildPolicy{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			identifyRelationshipsInDesign(typedDesign, typedRelDefs, policy)
		}
	})

	b.Run("OPA", func(b *testing.B) {
		ctx := context.Background()
		opaData := map[string]interface{}{
			"relationships": []interface{}{relDefs[0]},
		}
		pq := prepareOPAQuery(b,
			`data.eval.identify_relationships_in_design(input, data.relationships, "hierarchical_parent_child")`,
			design, opaData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = pq.Eval(ctx)
		}
	})
}

// BenchmarkActions/Go vs BenchmarkActions/OPA
func BenchmarkActions(b *testing.B) {
	design := benchDesign()
	typedDesign := benchTypedDesign()

	b.Run("Go", func(b *testing.B) {
		policy := &HierarchicalParentChildPolicy{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			generateActionsToApplyOnDesign(typedDesign, policy)
		}
	})

	b.Run("OPA", func(b *testing.B) {
		ctx := context.Background()
		pq := prepareOPAQuery(b,
			`data.eval.generate_actions_to_apply_on_design(input, "hierarchical_parent_child")`,
			design, map[string]interface{}{})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = pq.Eval(ctx)
		}
	})
}

// BenchmarkFullPipeline/Go benchmarks the full Go evaluation pipeline.
func BenchmarkFullPipeline(b *testing.B) {
	design := benchDesign()
	typedRelDefs := benchTypedRelDefs()
	relDefs := benchRelDefs()

	b.Run("Go", func(b *testing.B) {
		log, _ := logger.New("bench", logger.Options{Format: logger.SyslogLogFormat})
		engine := NewGoEngine(log)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			engine.evaluate(deepCopyDesign(benchTypedDesign()), typedRelDefs)
		}
	})

	b.Run("OPA_ColdCompile", func(b *testing.B) {
		modules := loadOPAModules(b)
		ctx := context.Background()
		store := inmem.NewFromObject(map[string]interface{}{
			"relationships": []interface{}{relDefs[0]},
		})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			opts := append(modules,
				rego.Query("data.relationship_evaluation_policy.evaluate"),
				rego.Input(design),
				rego.Store(store),
			)
			r := rego.New(opts...)
			_, _ = r.Eval(ctx)
		}
	})

	b.Run("OPA_Prepared", func(b *testing.B) {
		ctx := context.Background()
		opaData := map[string]interface{}{
			"relationships": []interface{}{relDefs[0]},
		}
		pq := prepareOPAQuery(b,
			`data.relationship_evaluation_policy.evaluate`,
			design, opaData)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = pq.Eval(ctx)
		}
	})
}

// BenchmarkMatchValues compares a single rule evaluation.
func BenchmarkMatchValues(b *testing.B) {
	b.Run("Go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			matchValues("default", "default", "equal")
		}
	})

	b.Run("OPA", func(b *testing.B) {
		ctx := context.Background()
		pq := prepareOPAQuery(b,
			`data.eval_rules.match_values("default", "default", "equal")`,
			map[string]interface{}{}, map[string]interface{}{})
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = pq.Eval(ctx)
		}
	})
}
