// # Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gocarina/gocsv"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/spf13/cobra"
)

const sampleModelJSONFixture = `{
  "id": "11111111-1111-1111-1111-111111111111",
  "schemaVersion": "models.meshery.io/v1beta2",
  "version": "v1.0.0",
  "name": "sample-model",
  "displayName": "Sample Model Display",
  "description": "A sample model for testing",
  "status": "enabled",
  "registrant": {
    "name": "Github",
    "kind": "github"
  },
  "category": {
    "name": "Cloud Native Storage"
  },
  "subCategory": "Storage",
  "metadata": {
    "primaryColor": "#00D3A9",
    "secondaryColor": "#7AA116",
    "shape": "circle",
    "styleOverrides": "custom-override",
    "shapePolygonPoints": "10,20",
    "logoURL": "https://example.com/logo.png",
    "svgColor": "<svg id=\"color\"></svg>",
    "svgWhite": "<svg id=\"white\"></svg>",
    "svgComplete": "<svg id=\"complete\"></svg>",
    "source_uri": "git://github.com/meshery/sample/helm",
    "styles": {
      "stroke": "black"
    },
    "defaultData": {
      "replica": 1
    },
    "capabilities": [
      { "kind": "action", "displayName": "Scale" }
    ],
    "isAnnotation": false,
    "published": true
  },
  "model": {
    "version": "v2.5.0"
  }
}`

const sampleComponentJSONFixture = `{
  "id": "22222222-2222-2222-2222-222222222222",
  "schemaVersion": "components.meshery.io/v1beta2",
  "version": "v1.0.0",
  "displayName": "Sample Component",
  "description": "A sample component for testing",
  "status": "enabled",
  "model": {
    "name": "sample-model",
    "displayName": "Sample Model Display",
    "version": "v1.0.0",
    "registrant": {
      "name": "Github",
      "kind": "github"
    },
    "model": {
      "version": "v2.5.0"
    }
  },
  "styles": {
    "primaryColor": "#00D3A9",
    "secondaryColor": "#7AA116",
    "shape": "rectangle",
    "svgColor": "<svg id=\"comp-color\"></svg>",
    "svgWhite": "<svg id=\"comp-white\"></svg>",
    "svgComplete": "<svg id=\"comp-complete\"></svg>",
    "styleOverrides": "comp-override",
    "shapePolygonPoints": "5,15"
  },
  "capabilities": [
    { "kind": "mutate", "displayName": "Configure" }
  ],
  "metadata": {
    "genealogy": "parent-node",
    "logoURL": "https://example.com/comp.png",
    "isAnnotation": true,
    "defaultData": {
      "port": 8080
    }
  },
  "component": {
    "kind": "SampleKind",
    "version": "apps/v1",
    "schema": "{\"properties\":{\"replicas\":{\"type\":\"integer\"}}}"
  }
}`

func TestParseModelJSONFile_FullyPopulated(t *testing.T) {
	tempDir := t.TempDir()
	modelFile := filepath.Join(tempDir, "model.json")
	if err := os.WriteFile(modelFile, []byte(sampleModelJSONFixture), 0644); err != nil {
		t.Fatalf("failed to write test fixture: %v", err)
	}

	modelCSV, err := ParseModelJSONFile(modelFile)
	if err != nil {
		t.Fatalf("ParseModelJSONFile failed: %v", err)
	}

	if modelCSV.Model != "sample-model" {
		t.Errorf("expected Model sample-model, got %s", modelCSV.Model)
	}
	if modelCSV.ModelDisplayName != "Sample Model Display" {
		t.Errorf("expected ModelDisplayName 'Sample Model Display', got %s", modelCSV.ModelDisplayName)
	}
	if modelCSV.Registrant != "Github" {
		t.Errorf("expected Registrant 'Github', got %s", modelCSV.Registrant)
	}
	if string(modelCSV.Category) != "Cloud Native Storage" {
		t.Errorf("expected Category 'Cloud Native Storage', got %s", modelCSV.Category)
	}
	if string(modelCSV.SubCategory) != "Storage" {
		t.Errorf("expected SubCategory 'Storage', got %s", modelCSV.SubCategory)
	}
	if modelCSV.Description != "A sample model for testing" {
		t.Errorf("expected Description, got %s", modelCSV.Description)
	}
	if modelCSV.Shape != "circle" {
		t.Errorf("expected Shape 'circle', got %s", modelCSV.Shape)
	}
	if modelCSV.PrimaryColor != "#00D3A9" {
		t.Errorf("expected PrimaryColor '#00D3A9', got %s", modelCSV.PrimaryColor)
	}
	if modelCSV.SecondaryColor != "#7AA116" {
		t.Errorf("expected SecondaryColor '#7AA116', got %s", modelCSV.SecondaryColor)
	}
	if modelCSV.SourceURL != "git://github.com/meshery/sample/helm" {
		t.Errorf("expected SourceURL 'git://github.com/meshery/sample/helm', got %s", modelCSV.SourceURL)
	}
	if modelCSV.SVGColor != "<svg id=\"color\"></svg>" {
		t.Errorf("expected SVGColor, got %s", modelCSV.SVGColor)
	}
	if modelCSV.IsAnnotation != "FALSE" {
		t.Errorf("expected IsAnnotation 'FALSE', got %s", modelCSV.IsAnnotation)
	}
	if modelCSV.PublishToRegistry != "TRUE" {
		t.Errorf("expected PublishToRegistry 'TRUE', got %s", modelCSV.PublishToRegistry)
	}
	if !strings.Contains(modelCSV.Styles, "stroke") {
		t.Errorf("expected Styles to contain stroke, got %s", modelCSV.Styles)
	}
	if !strings.Contains(modelCSV.DefaultData, "replica") {
		t.Errorf("expected DefaultData to contain replica, got %s", modelCSV.DefaultData)
	}
	if !strings.Contains(modelCSV.Capabilities, "Scale") {
		t.Errorf("expected Capabilities to contain Scale, got %s", modelCSV.Capabilities)
	}
}

func TestParseModelJSONFile_MissingOptionalFields(t *testing.T) {
	tempDir := t.TempDir()
	minimalJSON := `{"name": "minimal-model"}`
	modelFile := filepath.Join(tempDir, "model.json")
	if err := os.WriteFile(modelFile, []byte(minimalJSON), 0644); err != nil {
		t.Fatalf("failed to write test fixture: %v", err)
	}

	modelCSV, err := ParseModelJSONFile(modelFile)
	if err != nil {
		t.Fatalf("ParseModelJSONFile failed on minimal model: %v", err)
	}

	if modelCSV.Model != "minimal-model" {
		t.Errorf("expected Model 'minimal-model', got %s", modelCSV.Model)
	}
	if modelCSV.ModelDisplayName != "minimal-model" {
		t.Errorf("expected fallback ModelDisplayName 'minimal-model', got %s", modelCSV.ModelDisplayName)
	}
	if modelCSV.Registrant != "Github" {
		t.Errorf("expected default Registrant 'Github', got %s", modelCSV.Registrant)
	}
	if modelCSV.PublishToRegistry != "FALSE" {
		t.Errorf("expected default PublishToRegistry 'FALSE', got %s", modelCSV.PublishToRegistry)
	}
	if modelCSV.IsAnnotation != "FALSE" {
		t.Errorf("expected default IsAnnotation 'FALSE', got %s", modelCSV.IsAnnotation)
	}
	if modelCSV.Category != "" || modelCSV.SubCategory != "" || modelCSV.Shape != "" {
		t.Errorf("expected optional fields to be empty string")
	}
}

func TestParseComponentJSONFile_FullyPopulated(t *testing.T) {
	tempDir := t.TempDir()
	compFile := filepath.Join(tempDir, "SampleKind.json")
	if err := os.WriteFile(compFile, []byte(sampleComponentJSONFixture), 0644); err != nil {
		t.Fatalf("failed to write test fixture: %v", err)
	}

	compCSV, err := ParseComponentJSONFile(compFile)
	if err != nil {
		t.Fatalf("ParseComponentJSONFile failed: %v", err)
	}

	if compCSV.Model != "sample-model" {
		t.Errorf("expected Model 'sample-model', got %s", compCSV.Model)
	}
	if compCSV.Component != "SampleKind" {
		t.Errorf("expected Component 'SampleKind', got %s", compCSV.Component)
	}
	if compCSV.Registrant != "Github" {
		t.Errorf("expected Registrant 'Github', got %s", compCSV.Registrant)
	}
	if compCSV.Description != "A sample component for testing" {
		t.Errorf("expected Description, got %s", compCSV.Description)
	}
	if compCSV.Shape != "rectangle" {
		t.Errorf("expected Shape 'rectangle', got %s", compCSV.Shape)
	}
	if compCSV.PrimaryColor != "#00D3A9" {
		t.Errorf("expected PrimaryColor '#00D3A9', got %s", compCSV.PrimaryColor)
	}
	if compCSV.SecondaryColor != "#7AA116" {
		t.Errorf("expected SecondaryColor '#7AA116', got %s", compCSV.SecondaryColor)
	}
	if compCSV.SVGColor != "<svg id=\"comp-color\"></svg>" {
		t.Errorf("expected SVGColor, got %s", compCSV.SVGColor)
	}
	if compCSV.StyleOverrides != "comp-override" {
		t.Errorf("expected StyleOverrides 'comp-override', got %s", compCSV.StyleOverrides)
	}
	if compCSV.ShapePolygonPoints != "5,15" {
		t.Errorf("expected ShapePolygonPoints '5,15', got %s", compCSV.ShapePolygonPoints)
	}
	if compCSV.Genealogy != "parent-node" {
		t.Errorf("expected Genealogy 'parent-node', got %s", compCSV.Genealogy)
	}
	if compCSV.IsAnnotation != "TRUE" {
		t.Errorf("expected IsAnnotation 'TRUE', got %s", compCSV.IsAnnotation)
	}
	if compCSV.Version != "apps/v1" {
		t.Errorf("expected Version 'apps/v1', got %s", compCSV.Version)
	}
	if compCSV.Status != "enabled" {
		t.Errorf("expected Status 'enabled', got %s", compCSV.Status)
	}
	if !strings.Contains(compCSV.Schema, "replicas") {
		t.Errorf("expected Schema to contain replicas, got %s", compCSV.Schema)
	}
	if !strings.Contains(compCSV.Capabilities, "Configure") {
		t.Errorf("expected Capabilities to contain Configure, got %s", compCSV.Capabilities)
	}
}

func TestParseComponentJSONFile_MissingOptionalFields(t *testing.T) {
	tempDir := t.TempDir()
	minimalJSON := `{
		"displayName": "FallbackKind",
		"model": {
			"name": "fallback-model"
		}
	}`
	compFile := filepath.Join(tempDir, "FallbackKind.json")
	if err := os.WriteFile(compFile, []byte(minimalJSON), 0644); err != nil {
		t.Fatalf("failed to write test fixture: %v", err)
	}

	compCSV, err := ParseComponentJSONFile(compFile)
	if err != nil {
		t.Fatalf("ParseComponentJSONFile failed on minimal component: %v", err)
	}

	if compCSV.Model != "fallback-model" {
		t.Errorf("expected Model 'fallback-model', got %s", compCSV.Model)
	}
	if compCSV.Component != "FallbackKind" {
		t.Errorf("expected Component 'FallbackKind', got %s", compCSV.Component)
	}
	if compCSV.Registrant != "Github" {
		t.Errorf("expected Registrant 'Github', got %s", compCSV.Registrant)
	}
	if compCSV.Status != "enabled" {
		t.Errorf("expected default Status 'enabled', got %s", compCSV.Status)
	}
	if compCSV.IsAnnotation != "FALSE" {
		t.Errorf("expected default IsAnnotation 'FALSE', got %s", compCSV.IsAnnotation)
	}
}

func TestScanCommittedModels_MalformedAndDeduplication(t *testing.T) {
	tempDir := t.TempDir()

	// Create structure:
	// models/model-a/v1.0.0/v1.0.0/model.json (valid)
	// models/model-a/v1.0.0/v1.0.1/model.json (duplicate model-a:v1.0.0 - should be deduplicated)
	// models/model-b/v2.0.0/v1.0.0/model.json (valid)
	// models/model-c/v1.0.0/v1.0.0/model.json (malformed JSON: invalid syntax)

	dirA1 := filepath.Join(tempDir, "model-a", "v1.0.0", "v1.0.0")
	dirA2 := filepath.Join(tempDir, "model-a", "v1.0.0", "v1.0.1")
	dirB := filepath.Join(tempDir, "model-b", "v2.0.0", "v1.0.0")
	dirC := filepath.Join(tempDir, "model-c", "v1.0.0", "v1.0.0")

	for _, d := range []string{dirA1, dirA2, dirB, dirC} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", d, err)
		}
	}

	_ = os.WriteFile(filepath.Join(dirA1, "model.json"), []byte(`{"name": "model-a", "displayName": "Model A"}`), 0644)
	_ = os.WriteFile(filepath.Join(dirA2, "model.json"), []byte(`{"name": "model-a", "displayName": "Model A Duplicate"}`), 0644)
	_ = os.WriteFile(filepath.Join(dirB, "model.json"), []byte(`{"name": "model-b", "displayName": "Model B"}`), 0644)
	_ = os.WriteFile(filepath.Join(dirC, "model.json"), []byte(`{ malformed json content }`), 0644)

	// Scan all models
	models, err := ScanCommittedModels(tempDir, "")
	if err != nil {
		t.Fatalf("ScanCommittedModels failed: %v", err)
	}

	// Expect exactly 2 models: model-a (deduplicated) and model-b (model-c skipped due to malformed JSON)
	if len(models) != 2 {
		t.Fatalf("expected 2 models after deduplication and skipping malformed, got %d", len(models))
	}

	foundA, foundB := false, false
	for _, m := range models {
		if m.Model == "model-a" {
			foundA = true
		}
		if m.Model == "model-b" {
			foundB = true
		}
	}
	if !foundA || !foundB {
		t.Errorf("expected both model-a and model-b, got foundA=%v, foundB=%v", foundA, foundB)
	}

	// Test targeting only model-b
	targetedModels, err := ScanCommittedModels(tempDir, "model-b")
	if err != nil {
		t.Fatalf("ScanCommittedModels with targetModel failed: %v", err)
	}
	if len(targetedModels) != 1 || targetedModels[0].Model != "model-b" {
		t.Errorf("expected 1 model-b, got %v", targetedModels)
	}
}

func TestScanCommittedComponents_MalformedAndDeduplication(t *testing.T) {
	tempDir := t.TempDir()

	// Structure:
	// models/model-a/v1.0.0/v1.0.0/components/Comp1.json (valid)
	// models/model-a/v1.0.0/v1.0.1/components/Comp1.json (duplicate version - should deduplicate)
	// models/model-a/v1.0.0/v1.0.0/components/Comp2.json (valid)
	// models/model-a/v1.0.0/v1.0.0/components/BadComp.json (malformed JSON)
	// models/model-b/v1.0.0/v1.0.0/components/Comp3.json (valid)

	compDirA1 := filepath.Join(tempDir, "model-a", "v1.0.0", "v1.0.0", "components")
	compDirA2 := filepath.Join(tempDir, "model-a", "v1.0.0", "v1.0.1", "components")
	compDirB := filepath.Join(tempDir, "model-b", "v1.0.0", "v1.0.0", "components")

	for _, d := range []string{compDirA1, compDirA2, compDirB} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", d, err)
		}
	}

	_ = os.WriteFile(filepath.Join(compDirA1, "Comp1.json"), []byte(`{"model": {"name": "model-a"}, "component": {"kind": "Comp1", "version": "v1"}}`), 0644)
	_ = os.WriteFile(filepath.Join(compDirA2, "Comp1.json"), []byte(`{"model": {"name": "model-a"}, "component": {"kind": "Comp1", "version": "v1"}}`), 0644)
	_ = os.WriteFile(filepath.Join(compDirA1, "Comp2.json"), []byte(`{"model": {"name": "model-a"}, "component": {"kind": "Comp2", "version": "v1"}}`), 0644)
	_ = os.WriteFile(filepath.Join(compDirA1, "BadComp.json"), []byte(`{ invalid json syntax ]`), 0644)
	_ = os.WriteFile(filepath.Join(compDirB, "Comp3.json"), []byte(`{"model": {"name": "model-b"}, "component": {"kind": "Comp3", "version": "v1"}}`), 0644)

	// Scan all components
	comps, err := ScanCommittedComponents(tempDir, "")
	if err != nil {
		t.Fatalf("ScanCommittedComponents failed: %v", err)
	}

	// Expect 3 components: Comp1 (deduplicated), Comp2, Comp3 (BadComp skipped)
	if len(comps) != 3 {
		t.Fatalf("expected 3 components, got %d", len(comps))
	}

	// Test target model
	targetedComps, err := ScanCommittedComponents(tempDir, "model-a")
	if err != nil {
		t.Fatalf("ScanCommittedComponents with targetModel failed: %v", err)
	}
	if len(targetedComps) != 2 {
		t.Errorf("expected 2 components for model-a, got %d", len(targetedComps))
	}
}

func TestExportModelsToCSV_RoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "models.csv")

	rows := []meshkitRegistryUtils.ModelCSV{
		{
			Registrant:        "Github",
			ModelDisplayName:  "Model One",
			Model:             "model-one",
			Category:          "Storage",
			SubCategory:       "Cloud",
			Description:       "First model",
			SourceURL:         "https://github.com/test/one",
			PublishToRegistry: "TRUE",
		},
		{
			Registrant:        "Artifact Hub",
			ModelDisplayName:  "Model Two",
			Model:             "model-two",
			Category:          "Networking",
			SubCategory:       "Service Mesh",
			Description:       "Second model",
			SourceURL:         "https://github.com/test/two",
			PublishToRegistry: "FALSE",
		},
	}

	if err := ExportModelsToCSV(rows, csvPath); err != nil {
		t.Fatalf("ExportModelsToCSV failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		t.Fatalf("CSV file was not created at %s", csvPath)
	}

	// Verify CSV headers using standard encoding/csv
	f, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to open exported CSV: %v", err)
	}
	reader := csv.NewReader(f)
	header, err := reader.Read()
	f.Close()
	if err != nil {
		t.Fatalf("failed to read CSV header: %v", err)
	}

	expectedHeaders := []string{
		"registrant", "modelDisplayName", "model", "category", "subCategory",
		"description", "sourceURL", "website", "docs", "shape", "primaryColor",
		"secondaryColor", "styleOverrides", "styles", "shapePolygonPoints",
		"defaultData", "capabilities", "logoURL", "svgColor", "svgWhite",
		"svgComplete", "isAnnotation", "publishToRegistry", "group",
	}

	if len(header) != len(expectedHeaders) {
		t.Fatalf("expected %d headers, got %d: %v", len(expectedHeaders), len(header), header)
	}
	for i, h := range expectedHeaders {
		if header[i] != h {
			t.Errorf("header[%d] expected %s, got %s", i, h, header[i])
		}
	}

	// Verify round-trip unmarshaling with gocsv
	f2, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to reopen CSV: %v", err)
	}
	defer f2.Close()

	var parsedRows []meshkitRegistryUtils.ModelCSV
	if err := gocsv.UnmarshalFile(f2, &parsedRows); err != nil {
		t.Fatalf("failed to unmarshal CSV back into ModelCSV: %v", err)
	}

	if len(parsedRows) != 2 {
		t.Fatalf("expected 2 parsed rows, got %d", len(parsedRows))
	}
	if parsedRows[0].Model != "model-one" || parsedRows[1].Model != "model-two" {
		t.Errorf("unexpected parsed models: %v", parsedRows)
	}
	if parsedRows[0].PublishToRegistry != "TRUE" || parsedRows[1].PublishToRegistry != "FALSE" {
		t.Errorf("unexpected publishToRegistry values: %v", parsedRows)
	}
}

func TestExportComponentsToCSV_RoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	csvPath := filepath.Join(tempDir, "components.csv")

	rows := []meshkitRegistryUtils.ComponentCSV{
		{
			Registrant:   "Github",
			Model:        "model-one",
			Component:    "Service",
			Description:  "Service component",
			Shape:        "circle",
			PrimaryColor: "#00D3A9",
			Version:      "v1",
			Status:       "enabled",
			IsAnnotation: "FALSE",
		},
		{
			Registrant:   "Github",
			Model:        "model-one",
			Component:    "Deployment",
			Description:  "Deployment component",
			Shape:        "rectangle",
			PrimaryColor: "#7AA116",
			Version:      "apps/v1",
			Status:       "enabled",
			IsAnnotation: "TRUE",
		},
	}

	if err := ExportComponentsToCSV(rows, csvPath); err != nil {
		t.Fatalf("ExportComponentsToCSV failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(csvPath); os.IsNotExist(err) {
		t.Fatalf("CSV file was not created at %s", csvPath)
	}

	// Verify CSV headers using standard encoding/csv
	f, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to open exported CSV: %v", err)
	}
	reader := csv.NewReader(f)
	header, err := reader.Read()
	f.Close()
	if err != nil {
		t.Fatalf("failed to read CSV header: %v", err)
	}

	expectedHeaders := []string{
		"registrant", "model", "component", "description", "shape",
		"primaryColor", "secondaryColor", "svgColor", "svgWhite", "svgComplete",
		"schema", "docs", "styleOverrides", "styles", "shapePolygonPoints",
		"defaultData", "capabilities", "logoURL", "genealogy", "isAnnotation",
		"version", "status",
	}

	if len(header) != len(expectedHeaders) {
		t.Fatalf("expected %d headers, got %d: %v", len(expectedHeaders), len(header), header)
	}
	for i, h := range expectedHeaders {
		if header[i] != h {
			t.Errorf("header[%d] expected %s, got %s", i, h, header[i])
		}
	}

	// Verify round-trip unmarshaling with gocsv
	f2, err := os.Open(csvPath)
	if err != nil {
		t.Fatalf("failed to reopen CSV: %v", err)
	}
	defer f2.Close()

	var parsedRows []meshkitRegistryUtils.ComponentCSV
	if err := gocsv.UnmarshalFile(f2, &parsedRows); err != nil {
		t.Fatalf("failed to unmarshal CSV back into ComponentCSV: %v", err)
	}

	if len(parsedRows) != 2 {
		t.Fatalf("expected 2 parsed rows, got %d", len(parsedRows))
	}
	if parsedRows[0].Component != "Service" || parsedRows[1].Component != "Deployment" {
		t.Errorf("unexpected parsed components: %v", parsedRows)
	}
}

func TestGenerateCmd_ExportOnlyPreRunE(t *testing.T) {
	// Verify that PreRunE succeeds when only export-models-csv is passed
	cmd := &cobra.Command{
		PreRunE: generateCmd.PreRunE,
	}
	cmd.Flags().String("spreadsheet-id", "", "")
	cmd.Flags().String("spreadsheet-cred", "", "")
	cmd.Flags().String("registrant-def", "", "")
	cmd.Flags().String("registrant-cred", "", "")
	cmd.Flags().String("directory", "", "")
	cmd.Flags().String("model-csv", "", "")
	cmd.Flags().String("component-csv", "", "")
	cmd.Flags().String("relationship-csv", "", "")
	cmd.Flags().String("export-models-csv", "/tmp/models.csv", "")
	cmd.Flags().String("export-components-csv", "", "")

	if err := cmd.PreRunE(cmd, []string{}); err != nil {
		t.Fatalf("expected PreRunE to succeed with only --export-models-csv, got error: %v", err)
	}

	// Verify that PreRunE also succeeds when only export-components-csv is passed
	cmd2 := &cobra.Command{
		PreRunE: generateCmd.PreRunE,
	}
	cmd2.Flags().String("spreadsheet-id", "", "")
	cmd2.Flags().String("spreadsheet-cred", "", "")
	cmd2.Flags().String("registrant-def", "", "")
	cmd2.Flags().String("registrant-cred", "", "")
	cmd2.Flags().String("directory", "", "")
	cmd2.Flags().String("model-csv", "", "")
	cmd2.Flags().String("component-csv", "", "")
	cmd2.Flags().String("relationship-csv", "", "")
	cmd2.Flags().String("export-models-csv", "", "")
	cmd2.Flags().String("export-components-csv", "/tmp/components.csv", "")

	if err := cmd2.PreRunE(cmd2, []string{}); err != nil {
		t.Fatalf("expected PreRunE to succeed with only --export-components-csv, got error: %v", err)
	}

	// Verify that PreRunE fails when NO inputs or export flags are passed
	cmd3 := &cobra.Command{
		PreRunE: generateCmd.PreRunE,
	}
	cmd3.Flags().String("spreadsheet-id", "", "")
	cmd3.Flags().String("spreadsheet-cred", "", "")
	cmd3.Flags().String("registrant-def", "", "")
	cmd3.Flags().String("registrant-cred", "", "")
	cmd3.Flags().String("directory", "", "")
	cmd3.Flags().String("model-csv", "", "")
	cmd3.Flags().String("component-csv", "", "")
	cmd3.Flags().String("relationship-csv", "", "")
	cmd3.Flags().String("export-models-csv", "", "")
	cmd3.Flags().String("export-components-csv", "", "")

	if err := cmd3.PreRunE(cmd3, []string{}); err == nil {
		t.Fatalf("expected PreRunE to fail when no inputs or export flags are specified, but got nil")
	}
}
