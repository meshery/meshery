package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestKroModelSeedingAndComponents validates that the Kro model is correctly
// structured, discovered by GetModelDirectoryPaths, and bundles valid
// ResourceGraphDefinition, Graph, and GraphRevision components as well as
// required relationships.
func TestKroModelSeedingAndComponents(t *testing.T) {
	seededDirs, err := GetModelDirectoryPaths(ModelsPath)
	if err != nil {
		t.Fatalf("GetModelDirectoryPaths(%q): %v", ModelsPath, err)
	}

	var kroDefDir string
	for _, dir := range seededDirs {
		if strings.Contains(filepath.ToSlash(dir), "/kro/") {
			kroDefDir = dir
			break
		}
	}
	if kroDefDir == "" {
		t.Fatalf("kro model directory not found in seeded model directories")
	}

	// 1. Verify model.json
	modelJSONPath := filepath.Join(kroDefDir, "model.json")
	modelBytes, err := os.ReadFile(modelJSONPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", modelJSONPath, err)
	}

	var modelData map[string]interface{}
	if err := json.Unmarshal(modelBytes, &modelData); err != nil {
		t.Fatalf("model.json is not valid JSON: %v", err)
	}

	if name, ok := modelData["name"].(string); !ok || name != "kro" {
		t.Errorf("expected model name 'kro', got %v", modelData["name"])
	}
	if disp, ok := modelData["displayName"].(string); !ok || disp != "Kro" {
		t.Errorf("expected displayName 'Kro', got %v", modelData["displayName"])
	}
	cat, ok := modelData["category"].(map[string]interface{})
	if !ok || cat["name"] != "Orchestration & Management" {
		t.Errorf("expected category 'Orchestration & Management', got %v", cat)
	}

	// 2. Verify Components
	compsDir := filepath.Join(kroDefDir, "components")
	expectedComps := map[string]struct {
		kind         string
		version      string
		isNamespaced bool
	}{
		"ResourceGraphDefinition.json": {
			kind:         "ResourceGraphDefinition",
			version:      "kro.run/v1alpha1",
			isNamespaced: false,
		},
		"Graph.json": {
			kind:         "Graph",
			version:      "kro.run/v1alpha1",
			isNamespaced: true,
		},
		"GraphRevision.json": {
			kind:         "GraphRevision",
			version:      "internal.kro.run/v1alpha1",
			isNamespaced: false,
		},
	}

	for fileName, exp := range expectedComps {
		compPath := filepath.Join(compsDir, fileName)
		compBytes, err := os.ReadFile(compPath)
		if err != nil {
			t.Errorf("missing component %s: %v", fileName, err)
			continue
		}

		var compData struct {
			SchemaVersion string `json:"schemaVersion"`
			DisplayName   string `json:"displayName"`
			Metadata      struct {
				IsNamespaced bool `json:"isNamespaced"`
			} `json:"metadata"`
			Component struct {
				Kind    string `json:"kind"`
				Version string `json:"version"`
				Schema  string `json:"schema"`
			} `json:"component"`
		}

		if err := json.Unmarshal(compBytes, &compData); err != nil {
			t.Errorf("component %s failed JSON unmarshal: %v", fileName, err)
			continue
		}

		if compData.Component.Kind != exp.kind {
			t.Errorf("component %s: expected kind %s, got %s", fileName, exp.kind, compData.Component.Kind)
		}
		if compData.Component.Version != exp.version {
			t.Errorf("component %s: expected version %s, got %s", fileName, exp.version, compData.Component.Version)
		}
		if compData.Metadata.IsNamespaced != exp.isNamespaced {
			t.Errorf("component %s: expected isNamespaced=%v, got %v", fileName, exp.isNamespaced, compData.Metadata.IsNamespaced)
		}
		if len(compData.Component.Schema) == 0 {
			t.Errorf("component %s: empty schema", fileName)
		}

		// Verify schema itself is valid JSON
		var schemaCheck map[string]interface{}
		if err := json.Unmarshal([]byte(compData.Component.Schema), &schemaCheck); err != nil {
			t.Errorf("component %s: embedded schema is not valid JSON: %v", fileName, err)
		}
	}

	// 3. Verify Relationships
	relsDir := filepath.Join(kroDefDir, "relationships")
	relEntries, err := os.ReadDir(relsDir)
	if err != nil {
		t.Fatalf("failed to read relationships dir: %v", err)
	}
	if len(relEntries) < 2 {
		t.Errorf("expected at least 2 relationships for kro, got %d", len(relEntries))
	}

	for _, entry := range relEntries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		relBytes, err := os.ReadFile(filepath.Join(relsDir, entry.Name()))
		if err != nil {
			t.Errorf("failed to read relationship %s: %v", entry.Name(), err)
			continue
		}
		var relData map[string]interface{}
		if err := json.Unmarshal(relBytes, &relData); err != nil {
			t.Errorf("relationship %s failed JSON unmarshal: %v", entry.Name(), err)
		}
	}
}
