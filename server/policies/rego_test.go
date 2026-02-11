// Package policies provides Rego policy testing using the OPA Go SDK.
package policies

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/open-policy-agent/opa/v1/ast"
	"github.com/open-policy-agent/opa/v1/rego"
	"github.com/open-policy-agent/opa/v1/tester"
	"github.com/open-policy-agent/opa/v1/topdown"
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

	policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"
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

	policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"

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

	policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"

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
	policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"

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
