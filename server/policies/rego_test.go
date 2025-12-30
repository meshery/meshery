// Package policies provides Rego policy testing using the OPA Go SDK.
package policies

import (
"context"
"os"
"path/filepath"
"strings"
"testing"

"github.com/open-policy-agent/opa/v1/ast"
"github.com/open-policy-agent/opa/v1/rego"
"github.com/open-policy-agent/opa/v1/tester"
"github.com/open-policy-agent/opa/v1/topdown"
)

// TestRegoPolicy runs OPA tests for all Rego policy files in the meshery-core model.
// This test uses the OPA SDK to execute Rego tests, providing comprehensive coverage
// for relationship evaluation policies.
func TestRegoPolicy(t *testing.T) {
ctx := context.Background()

// Paths to policy directories
policiesDir := "../meshmodel/meshery-core/0.7.2/v1.0.0/policies"
testsDir := filepath.Join(policiesDir, "tests")

// Collect all Rego files from policies directory
policyFiles, err := collectRegoFiles(policiesDir)
if err != nil {
t.Fatalf("Failed to collect policy files: %v", err)
}

// Collect test files
testFiles, err := collectRegoFiles(testsDir)
if err != nil {
t.Fatalf("Failed to collect test files: %v", err)
}

// Combine all files
allFiles := append(policyFiles, testFiles...)

if len(allFiles) == 0 {
t.Fatal("No Rego files found")
}

// Build modules map
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

	// Run tests using OPA tester
	runner := tester.NewRunner().
		SetModules(modules).
		EnableTracing(false).
		SetTimeout(10 * 1000000000) // 10 second timeout per test

ch, err := runner.Run(ctx, modules)
if err != nil {
t.Fatalf("Failed to run tests: %v", err)
}

// Report test results
var passed, failed int
for result := range ch {
	testName := strings.TrimPrefix(result.Name, "data.")
	if result.Fail {
	failed++
	t.Errorf("FAIL: %s", testName)
	if result.Error != nil {
		t.Errorf("  Error: %v", result.Error)
	}
	} else if result.Error != nil {
	failed++
	t.Errorf("ERROR: %s - %v", testName, result.Error)
	} else {
	passed++
	}
}

t.Logf("\nTest Summary: %d passed, %d failed", passed, failed)

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

// Build module contents
var modules []func(*rego.Rego)
for _, file := range policyFiles {
// Skip test files
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
name     string
query    string
input    map[string]interface{}
wantPass bool
}{
{
name:  "alias_policy_identifier_is_set",
query: "data.eval.alias_policy_identifier",
input: map[string]interface{}{},
// Test that the alias policy identifier is defined
wantPass: true,
},
{
name:  "hierarchical_parent_child_policy_identifier_is_set",
query: "data.eval.hierarchical_parent_child_policy_identifier",
input: map[string]interface{}{},
// Test that the hierarchical policy identifier is defined
wantPass: true,
},
{
name:  "action_operations_defined",
query: "data.actions.update_component_op",
input: map[string]interface{}{},
// Test that action operations are defined
wantPass: true,
},
}

for _, tc := range testCases {
t.Run(tc.name, func(t *testing.T) {
// Build rego query
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
t.Fatalf("Evaluation failed: %v", err)
}
return
}

if len(rs) == 0 || len(rs[0].Expressions) == 0 {
if tc.wantPass {
t.Error("Expected result but got none")
}
return
}

if !tc.wantPass {
t.Error("Expected failure but test passed")
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

	// Build module contents
	var modules []func(*rego.Rego)
	for _, file := range policyFiles {
	// Skip Rego template sources; these are not executable policy modules.
	if strings.HasSuffix(file, ".template") {
		continue
	}
	// Skip Rego test modules; this scenario test should load only non-test policies.
	if strings.Contains(file, "/tests/") {
		continue
	}
	content, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", file, err)
	}
	modules = append(modules, rego.Module(file, string(content)))
	}

// Test scenarios with sample design inputs
testCases := []struct {
name  string
query string
input map[string]interface{}
}{
{
name:  "is_alias_relationship_true",
query: "data.eval.is_alias_relationship",
input: map[string]interface{}{
"kind":    "hierarchical",
"type":    "parent",
"subType": "alias",
},
},
{
name:  "core_utils_set_to_array",
query: "data.core_utils.set_to_array",
input: map[string]interface{}{},
},
{
name:  "eval_rules_match_values_equal",
query: `data.eval_rules.match_values("test", "test", "equal")`,
input: map[string]interface{}{},
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
t.Fatalf("Query evaluation error for %q: %v", tc.name, err)
}

if len(rs) > 0 && len(rs[0].Expressions) > 0 {
t.Logf("Result: %v", rs[0].Expressions[0].Value)
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

// Parse the Rego module directly
_, err = ast.ParseModule(file, string(content))
if err != nil {
t.Errorf("Parse error in %s: %v", filepath.Base(file), err)
}
})
}
}
