package core

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/pattern/pattern"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta2/component"
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

func TestToCytoscapeJSDuplicateStyles(t *testing.T) {
	log, _ := logger.New("test", logger.Options{Format: logger.SyslogLogFormat})

	// Setup component definition with overlapping styles
	cmpId, _ := uuid.NewV4()
	bgColor := "blue"
	textOpacity := float32(1.0)
	cmp := component.ComponentDefinition{
		ID: cmpId,
		Styles: &core.ComponentStyles{
			BackgroundColor: &bgColor,
			TextOpacity:     &textOpacity,
			AdditionalProperties: map[string]interface{}{
				"backgroundColor": "red",
				"textOpacity":     0.5,
				"borderWidth":     "2px",
				"border-width":    "2px",
				"transparent":     "true", // ordinary CSS value
				"inherit":         "true",
			},
		},
	}

	pf := &pattern.PatternFile{
		Components: []*component.ComponentDefinition{&cmp},
	}

	cy, err := ToCytoscapeJS(pf, log)
	if err != nil {
		t.Fatalf("ToCytoscapeJS failed: %v", err)
	}

	if len(cy.Elements) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(cy.Elements))
	}

	elem := cy.Elements[0]
	scratch, ok := elem.Scratch.(map[string]interface{})
	if !ok {
		t.Fatalf("Scratch is not map[string]interface{}")
	}

	data, ok := scratch["_data"].(map[string]interface{})
	if !ok {
		t.Fatalf("_data is not map[string]interface{}")
	}

	styles, ok := data["styles"].(map[string]interface{})
	if !ok {
		t.Fatalf("styles is not map[string]interface{}")
	}

	// Verify duplicate camelCase/kebab-case resolution
	if _, exists := styles["backgroundColor"]; exists {
		t.Errorf("Expected camelCase 'backgroundColor' to be deleted")
	}
	if _, exists := styles["textOpacity"]; exists {
		t.Errorf("Expected camelCase 'textOpacity' to be deleted")
	}
	if _, exists := styles["borderWidth"]; exists {
		t.Errorf("Expected camelCase 'borderWidth' to be deleted")
	}

	// Verify legitimate kebab-case properties exist
	if _, exists := styles["background-color"]; !exists {
		t.Errorf("Expected kebab-case 'background-color' to be preserved")
	}
	if _, exists := styles["text-opacity"]; !exists {
		t.Errorf("Expected kebab-case 'text-opacity' to be preserved")
	}
	if _, exists := styles["border-width"]; !exists {
		t.Errorf("Expected kebab-case 'border-width' to be preserved")
	}

	// Verify other unrelated values
	if _, exists := styles["transparent"]; !exists {
		t.Errorf("Expected 'transparent' to be preserved")
	}
	if _, exists := styles["inherit"]; !exists {
		t.Errorf("Expected 'inherit' to be preserved")
	}
}
