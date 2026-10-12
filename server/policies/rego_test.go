// Package policies provides Rego policy testing using the OPA Go SDK.
package policies

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/tester"
	"github.com/open-policy-agent/opa/v1/topdown"

	"github.com/meshery/schemas/models/v1beta1/component"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// designFileName extracts a human-friendly design file name from a test location.
// If no location is available, a placeholder is returned for logging clarity.
func designFileName(loc *ast.Location) string {
	if loc == nil || loc.File == "" {
		return "design file: not specified"
	}

	return filepath.Base(loc.File)
}

// TestRegoPolicy runs OPA tests for all Rego policy files in the meshery-core model.
// This test uses the OPA SDK to execute Rego tests, providing comprehensive coverage
// for relationship evaluation policies.
func TestRegoPolicy(t *testing.T) {
	ctx := context.Background()

	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"
	testsDir := filepath.Join(policiesDir, "tests")

	policyFiles, err := collectRegoFiles(policiesDir)
	if err != nil {
		t.Fatalf("Failed to collect policy files: %v", err)
	}

	testFiles, err := collectRegoFiles(testsDir)
	if err != nil {
		t.Fatalf("Failed to collect test files: %v", err)
	}

	allFiles := append(policyFiles, testFiles...)
	if len(allFiles) == 0 {
		t.Fatal("No Rego files found")
	}

	modules := make(map[string]*ast.Module)
	for _, file := range allFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", file, err)
		}

		module, err := ast.ParseModule(file, string(content))
		if err != nil {
			t.Fatalf("Failed to parse module %s: %v", file, err)
		}

		modules[file] = module
	}

	runner := tester.NewRunner().
		SetModules(modules).
		EnableTracing(false).
		SetTimeout(10 * time.Second)

	runner.SetModules(modules)
	ch, err := runner.RunTests(ctx, nil) // Pass nil for transaction to match previous behavior (Run() used nil internally)
	if err != nil {
		t.Fatalf("Failed to run tests: %v", err)
	}

	var passed, failed int
	for result := range ch {
		testName := strings.TrimPrefix(result.Name, "data.")
		moduleFile := "unknown"
		if result.Location != nil && result.Location.File != "" {
			moduleFile = result.Location.File
		}

		designFile := designFileName(result.Location)

		switch {
		case result.Fail:
			failed++
			t.Errorf("FAIL: %s (design: %s, module: %s)", testName, designFile, moduleFile)
			if result.Error != nil {
				t.Errorf("  Error: %v", result.Error)
			}
		case result.Error != nil:
			failed++
			t.Errorf("ERROR: %s (design: %s, module: %s) - %v", testName, designFile, moduleFile, result.Error)
		default:
			passed++
			t.Logf("PASS: %s (design: %s, module: %s)", testName, designFile, moduleFile)
		}
	}

	t.Logf("Test Summary: %d passed, %d failed", passed, failed)

	if failed > 0 {
		t.Fail()
	}
}

// TestRegoPolicyRules tests specific policy rules with sample inputs.
// These are sanity checks that verify policy rule definitions and basic evaluation behavior.
func TestRegoPolicyRules(t *testing.T) {
	ctx := context.Background()

	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"

	policyFiles, err := collectRegoFiles(policiesDir)
	if err != nil {
		t.Fatalf("Failed to collect policy files: %v", err)
	}

	var modules []func(*rego.Rego)
	for _, file := range policyFiles {
		if strings.Contains(file, "/tests/") {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", file, err)
		}

		modules = append(modules, rego.Module(file, string(content)))
	}

	testCases := []struct {
		name          string
		query         string
		input         map[string]interface{}
		wantPass      bool
		designFileRef string
	}{
		{
			name:          "alias_policy_identifier_is_set",
			query:         "data.eval.alias_policy_identifier",
			input:         map[string]interface{}{},
			wantPass:      true,
			designFileRef: "inline",
		},
		{
			name:          "hierarchical_parent_child_policy_identifier_is_set",
			query:         "data.eval.hierarchical_parent_child_policy_identifier",
			input:         map[string]interface{}{},
			wantPass:      true,
			designFileRef: "inline",
		},
		{
			name:          "action_operations_defined",
			query:         "data.actions.update_component_op",
			input:         map[string]interface{}{},
			wantPass:      true,
			designFileRef: "inline",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := append(modules,
				rego.Query(tc.query),
				rego.Input(tc.input),
				rego.EnablePrintStatements(true),
				rego.PrintHook(topdown.NewPrintHook(os.Stderr)),
			)

			r := rego.New(opts...)
			rs, err := r.Eval(ctx)
			if err != nil {
				if tc.wantPass {
					t.Fatalf("Evaluation failed (%s): %v", tc.designFileRef, err)
				}
				return
			}

			if len(rs) == 0 || len(rs[0].Expressions) == 0 {
				if tc.wantPass {
					t.Errorf("Expected result but got none (design: %s)", tc.designFileRef)
				}
				return
			}

			if !tc.wantPass {
				t.Errorf("Expected failure but test passed (design: %s)", tc.designFileRef)
			}
		})
	}
}

// TestRelationshipEvaluationScenarios tests relationship evaluation with realistic design inputs.
func TestRelationshipEvaluationScenarios(t *testing.T) {
	ctx := context.Background()

	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"

	policyFiles, err := collectRegoFiles(policiesDir)
	if err != nil {
		t.Fatalf("Failed to collect policy files: %v", err)
	}

	var modules []func(*rego.Rego)
	for _, file := range policyFiles {
		if strings.HasSuffix(file, ".template") {
			continue
		}

		if strings.Contains(file, "/tests/") {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("Failed to read file %s: %v", file, err)
		}

		modules = append(modules, rego.Module(file, string(content)))
	}

	testCases := []struct {
		name        string
		query       string
		input       map[string]interface{}
		expectError bool
		designFile  string
		expected    interface{}
	}{
		{
			name:  "is_alias_relationship_true",
			query: "data.eval.is_alias_relationship",
			input: map[string]interface{}{
				"kind":    "hierarchical",
				"type":    "parent",
				"subType": "alias",
			},
			expectError: true,
			designFile:  "inline",
		},
		{
			name:        "core_utils_set_to_array",
			query:       "data.core_utils.set_to_array",
			input:       map[string]interface{}{},
			expectError: true,
			designFile:  "inline",
		},
		{
			name:        "eval_rules_match_values_equal",
			query:       "data.eval_rules.match_values(\"test\", \"test\", \"equal\")",
			input:       map[string]interface{}{},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
		{
			name:  "namespace_cannot_be_parent_of_namespace",
			query: "data.relationship_evaluation_policy.is_relationship_denied(input.from, input.to, input.deny_selectors)",
			input: map[string]interface{}{
				"from": map[string]interface{}{
					"id":        "from-namespace",
					"component": map[string]interface{}{"kind": "Namespace"},
					"model": map[string]interface{}{
						"name": "kubernetes",
						"registrant": map[string]interface{}{
							"kind": "github",
						},
					},
				},
				"to": map[string]interface{}{
					"id":        "to-namespace",
					"component": map[string]interface{}{"kind": "Namespace"},
					"model": map[string]interface{}{
						"name": "kubernetes",
						"registrant": map[string]interface{}{
							"kind": "github",
						},
					},
				},
				"deny_selectors": map[string]interface{}{
					"from": []map[string]interface{}{
						{
							"kind": "Namespace",
							"model": map[string]interface{}{
								"name":       "kubernetes",
								"registrant": "*",
							},
						},
					},
					"to": []map[string]interface{}{
						{
							"kind": "Namespace",
							"model": map[string]interface{}{
								"name":       "kubernetes",
								"registrant": "*",
							},
						},
					},
				},
			},
			expectError: false,
			designFile:  "namespace_parent_inline",
			expected:    true,
		},
		{
			// Components in current designs carry the model name under
			// modelReference, not under a full model object (which is a
			// nullable pointer in the schema and absent from e2e design
			// fixtures). Feasibility must resolve the name from either
			// place, like the Go engine does (engine.go normalizes
			// Model.Name into ModelReference.Name before matching).
			name:  "feasibility_matches_model_reference_only_component",
			query: "data.feasibility_evaluation_utils.is_relationship_feasible(input.selector, input.component)",
			input: map[string]interface{}{
				"selector": map[string]interface{}{
					"kind":  "Pod",
					"model": map[string]interface{}{"name": "kubernetes"},
				},
				"component": map[string]interface{}{
					"component":      map[string]interface{}{"kind": "Pod"},
					"model":          nil,
					"modelReference": map[string]interface{}{"name": "kubernetes"},
				},
			},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
		{
			// Legacy shape keeps working: model object only, no
			// modelReference.
			name:  "feasibility_matches_model_only_component",
			query: "data.feasibility_evaluation_utils.is_relationship_feasible(input.selector, input.component)",
			input: map[string]interface{}{
				"selector": map[string]interface{}{
					"kind":  "Pod",
					"model": map[string]interface{}{"name": "kubernetes"},
				},
				"component": map[string]interface{}{
					"component": map[string]interface{}{"kind": "Pod"},
					"model":     map[string]interface{}{"name": "kubernetes"},
				},
			},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
		{
			// A JSON-null modelReference.name also falls back to the model
			// object instead of failing the match outright.
			name:  "feasibility_falls_back_on_null_model_reference_name",
			query: "data.feasibility_evaluation_utils.is_relationship_feasible(input.selector, input.component)",
			input: map[string]interface{}{
				"selector": map[string]interface{}{
					"kind":  "Pod",
					"model": map[string]interface{}{"name": "kubernetes"},
				},
				"component": map[string]interface{}{
					"component":      map[string]interface{}{"kind": "Pod"},
					"model":          map[string]interface{}{"name": "kubernetes"},
					"modelReference": map[string]interface{}{"name": nil},
				},
			},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
		{
			// When both fields are present with different names,
			// modelReference wins: the selector matching the modelReference
			// name is feasible.
			name:  "feasibility_prefers_model_reference_over_model_object",
			query: "data.feasibility_evaluation_utils.is_relationship_feasible(input.selector, input.component)",
			input: map[string]interface{}{
				"selector": map[string]interface{}{
					"kind":  "Pod",
					"model": map[string]interface{}{"name": "kubernetes"},
				},
				"component": map[string]interface{}{
					"component":      map[string]interface{}{"kind": "Pod"},
					"model":          map[string]interface{}{"name": "legacy-model"},
					"modelReference": map[string]interface{}{"name": "kubernetes"},
				},
			},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
		{
			// An empty modelReference.name falls back to the model object,
			// mirroring the Go engine's normalization (which only copies
			// Model.Name when ModelReference.Name is empty).
			name:  "feasibility_falls_back_on_empty_model_reference_name",
			query: "data.feasibility_evaluation_utils.is_relationship_feasible(input.selector, input.component)",
			input: map[string]interface{}{
				"selector": map[string]interface{}{
					"kind":  "Pod",
					"model": map[string]interface{}{"name": "kubernetes"},
				},
				"component": map[string]interface{}{
					"component":      map[string]interface{}{"kind": "Pod"},
					"model":          map[string]interface{}{"name": "kubernetes"},
					"modelReference": map[string]interface{}{"name": ""},
				},
			},
			expectError: false,
			designFile:  "inline",
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			opts := append(modules,
				rego.Query(tc.query),
				rego.Input(tc.input),
			)

			r := rego.New(opts...)
			rs, err := r.Eval(ctx)
			if err != nil {
				if !tc.expectError {
					t.Fatalf("Query evaluation error for %q (design: %s): %v", tc.name, tc.designFile, err)
				}

				t.Logf("Expected error for %q (design: %s): %v", tc.name, tc.designFile, err)
				return
			}

			if len(rs) == 0 || len(rs[0].Expressions) == 0 {
				if tc.expected != nil {
					t.Fatalf("No result returned for %q (design: %s)", tc.name, tc.designFile)
				}

				return
			}

			value := rs[0].Expressions[0].Value
			t.Logf("Result for %q (design: %s): %v", tc.name, tc.designFile, value)

			if tc.expected != nil && !reflect.DeepEqual(value, tc.expected) {
				t.Fatalf("Unexpected result for %q (design: %s): got %v, want %v", tc.name, tc.designFile, value, tc.expected)
			}
		})
	}
}

// collectRegoFiles collects all .rego files from a directory (recursively).
func collectRegoFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".rego") && !strings.HasSuffix(path, ".template") {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// TestRegoSyntax verifies all Rego files can be parsed without errors.
func TestRegoSyntax(t *testing.T) {
	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"

	files, err := collectRegoFiles(policiesDir)
	if err != nil {
		t.Fatalf("Failed to collect Rego files: %v", err)
	}

	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}

			_, err = ast.ParseModule(file, string(content))
			if err != nil {
				t.Errorf("Parse error in %s: %v", filepath.Base(file), err)
			}
		})
	}
}

// TestParity_RegoAndGoEngineIdentifyForModelReferenceOnlyComponent runs alias
// identification for a design whose component carries its model name only
// under modelReference (the shape of current designs and the UI e2e design
// fixtures; ComponentDefinition.Model is a nullable pointer that marshals to
// null when unset) through BOTH engines and asserts they agree. The Go engine
// normalizes the two fields before matching (engine.go), so the Rego
// feasibility check must resolve modelReference too.
func TestParity_RegoAndGoEngineIdentifyForModelReferenceOnlyComponent(t *testing.T) {
	podID, err := uuid.FromString("00000000-0000-0000-0000-0000000000c1")
	if err != nil {
		t.Fatalf("parse pod id: %v", err)
	}
	pod := &component.ComponentDefinition{
		Component:      component.Component{Kind: "Pod"},
		ModelReference: modelv1beta1.ModelReference{Name: "kubernetes"},
		Configuration: map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{"name": "app"},
					map[string]interface{}{"name": "sidecar"},
				},
			},
		},
	}
	pod.ID = podID

	design := makePatternFile([]*component.ComponentDefinition{pod}, nil)

	mutatorRef := relationship.MutatorRef{[]string{"configuration", "spec", "containers", "_"}}
	mutatedRef := relationship.MutatedRef{[]string{"configuration", "spec", "containers", "_"}}
	selectorSet := relationship.SelectorSet{
		relationship.SelectorSetItem{
			Allow: relationship.Selector{
				From: []relationship.SelectorItem{
					{
						Kind: strPtr("Container"),
						RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{
							MutatorRef: &mutatorRef,
						},
					},
				},
				To: []relationship.SelectorItem{
					{
						Kind:  strPtr("Pod"),
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
		SubType:          "alias",
		Model:            modelv1beta1.ModelReference{Name: "kubernetes"},
		Selectors:        &selectorSet,
	}
	relDef.ID, err = uuid.FromString("00000000-0000-0000-0000-0000000000d1")
	if err != nil {
		t.Fatalf("parse relationship id: %v", err)
	}

	goIdentified := (&AliasPolicy{}).IdentifyRelationship(relDef, design)
	if len(goIdentified) == 0 {
		t.Fatal("Go engine: expected identified alias relationships, got none")
	}

	// Feed the byte-identical input to the Rego engine. The marshaled
	// component has "model": null and the name only under "modelReference",
	// which is exactly the shape under test.
	relJSON, err := json.Marshal(relDef)
	if err != nil {
		t.Fatalf("marshal relationship: %v", err)
	}
	designJSON, err := json.Marshal(design)
	if err != nil {
		t.Fatalf("marshal design: %v", err)
	}
	var relMap, designMap map[string]interface{}
	if err := json.Unmarshal(relJSON, &relMap); err != nil {
		t.Fatalf("unmarshal relationship: %v", err)
	}
	if err := json.Unmarshal(designJSON, &designMap); err != nil {
		t.Fatalf("unmarshal design: %v", err)
	}

	policiesDir := "../../models/meshery-core/0.7.2/v1.0.0/policies"
	policyFiles, err := collectRegoFiles(policiesDir)
	if err != nil {
		t.Fatalf("collect rego files: %v", err)
	}
	var opts []func(*rego.Rego)
	for _, file := range policyFiles {
		if strings.Contains(file, "/tests/") || strings.HasSuffix(file, ".template") {
			continue
		}
		content, readErr := os.ReadFile(file)
		if readErr != nil {
			t.Fatalf("read %s: %v", file, readErr)
		}
		opts = append(opts, rego.Module(file, string(content)))
	}
	opts = append(opts,
		rego.Query(`data.eval.identify_relationship(input.relationship, input.design_file, "alias_relationships_policy")`),
		rego.Input(map[string]interface{}{
			"relationship": relMap,
			"design_file":  designMap,
		}),
	)

	rs, err := rego.New(opts...).Eval(context.Background())
	if err != nil {
		t.Fatalf("rego eval: %v", err)
	}

	regoCount := 0
	if len(rs) > 0 && len(rs[0].Expressions) > 0 {
		identified, ok := rs[0].Expressions[0].Value.([]interface{})
		if !ok {
			t.Fatalf("unexpected rego result shape: %#v", rs[0].Expressions[0].Value)
		}
		regoCount = len(identified)
	}

	if regoCount != len(goIdentified) {
		t.Fatalf("engine divergence: Go identified %d, Rego identified %d", len(goIdentified), regoCount)
	}
}
