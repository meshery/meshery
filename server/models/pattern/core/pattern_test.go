package core

import (
	"testing"

	"github.com/meshery/meshery/server/models/pattern/planner"
)

func TestDesignNameFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		expectedName string
	}{
		{
			name:         "given regular yaml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yaml",
			expectedName: "my-deployment",
		},
		{
			name:         "given regular yml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yml",
			expectedName: "my-deployment",
		},
		{
			name:         "given tar.gz file when DesignNameFromFileName then strip compound extension",
			fileName:     "my-chart.tar.gz",
			expectedName: "my-chart",
		},
		{
			name:         "given json file when DesignNameFromFileName then return name",
			fileName:     "config.json",
			expectedName: "config",
		},
		{
			name:         "given filename with multiple dots when DesignNameFromFileName then strip only last extension",
			fileName:     "my.k8s.deployment.yaml",
			expectedName: "my.k8s.deployment",
		},
		{
			name:         "given empty filename when DesignNameFromFileName then return empty",
			fileName:     "",
			expectedName: "",
		},
		{
			name:         "given filename without extension when DesignNameFromFileName then return as is",
			fileName:     "mydesign",
			expectedName: "mydesign",
		},
		{
			name:         "given tgz file when DesignNameFromFileName then strip extension",
			fileName:     "helm-chart.tgz",
			expectedName: "helm-chart",
		},
		{
			name:         "given extension only when DesignNameFromFileName then return empty",
			fileName:     ".yaml",
			expectedName: "",
		},
		{
			name:         "given tar.gz extension only when DesignNameFromFileName then return empty",
			fileName:     ".tar.gz",
			expectedName: "",
		},
		{
			name:         "given unsupported xls extension when DesignNameFromFileName then still strip extension",
			fileName:     "spreadsheet.xls",
			expectedName: "spreadsheet",
		},
		{
			name:         "given unsupported zip extension when DesignNameFromFileName then still strip extension",
			fileName:     "archive.zip",
			expectedName: "archive",
		},
		{
			name:         "given unsupported exe extension when DesignNameFromFileName then still strip extension",
			fileName:     "installer.exe",
			expectedName: "installer",
		},
		{
			name:         "given unsupported tar extension when DesignNameFromFileName then strip last suffix only",
			fileName:     "archive.tar",
			expectedName: "archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DesignNameFromFileName(tt.fileName)
			if result != tt.expectedName {
				t.Errorf("expected %q, got %q", tt.expectedName, result)
			}
		})
	}
}

// The UI serializes designs to YAML for deploy/undeploy/validate, and the
// schema types only capture unknown metadata keys (such as a component's
// "dependsOn") through their custom JSON unmarshalers. Without normalizing
// YAML input to JSON, those keys were silently dropped and the provision
// planner saw no ordering edges.
func TestNewPatternFileKeepsDependsOnFromYAML(t *testing.T) {
	yml := []byte(`
name: dependson-design
components:
  - id: 22222222-2222-2222-2222-222222222222
    displayName: guestbook
    metadata:
      isNamespaced: true
  - id: 11111111-1111-1111-1111-111111111111
    displayName: frontend
    metadata:
      isNamespaced: true
      dependsOn:
        - guestbook
`)

	pf, err := NewPatternFile(yml)
	if err != nil {
		t.Fatalf("NewPatternFile returned an error: %v", err)
	}
	if len(pf.Components) != 2 {
		t.Fatalf("expected 2 components, got %d", len(pf.Components))
	}

	var frontendDeps interface{}
	var found bool
	for _, comp := range pf.Components {
		if comp.DisplayName == "frontend" {
			frontendDeps, found = comp.Metadata.AdditionalProperties["dependsOn"]
			if !comp.Metadata.IsNamespaced {
				t.Fatal("typed metadata (isNamespaced) was lost while parsing the YAML design")
			}
		}
	}
	if !found {
		t.Fatalf("dependsOn was dropped while parsing the YAML design: %v", frontendDeps)
	}

	// The planner must be able to consume the wire-decoded shape
	// ([]interface{}) and produce an ordering edge from it.
	plan, err := planner.CreatePlan(pf, false)
	if err != nil {
		t.Fatalf("CreatePlan rejected the parsed design: %v", err)
	}
	if len(plan.Edges) == 0 {
		t.Fatal("expected the design's dependsOn to produce at least one ordering edge")
	}

	// Traversal dereferences edge endpoints; make sure a YAML-parsed
	// dependsOn does not panic when the plan is evaluated.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Plan.IsFeasible panicked on YAML-parsed dependsOn: %v", r)
		}
	}()
	_ = plan.IsFeasible()
}
