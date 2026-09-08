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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleRelationshipFixtureJSON = `{
  "id": "00000000-0000-0000-0000-000000000000",
  "kind": "edge",
  "type": "non-binding",
  "subType": "reference",
  "schemaVersion": "relationships.meshery.io/v1beta2",
  "model": {
    "name": "aws-ecs-controller",
    "version": "v1.0.0",
    "model": {
      "version": "v1.7.0"
    }
  },
  "metadata": {
    "description": "Sample test relationship",
    "isAnnotation": false,
    "styles": {
      "primaryColor": "#00D3A9"
    }
  },
  "selectors": [
    {
      "allow": {
        "from": [
          {
            "kind": "Service",
            "model": {
              "name": "aws-ecs-controller"
            }
          }
        ],
        "to": [
          {
            "kind": "Role",
            "model": {
              "name": "aws-ecs-controller"
            }
          }
        ]
      },
      "deny": {
        "from": [],
        "to": []
      }
    }
  ]
}`

// createTestModelsFixture creates a temporary models directory with a sample relationship JSON file.
func createTestModelsFixture(t *testing.T) (string, string) {
	t.Helper()
	tempDir := t.TempDir()
	relDir := filepath.Join(tempDir, "models", "aws-ecs-controller", "v1.7.0", "v1.0.0", "relationships")
	if err := os.MkdirAll(relDir, 0755); err != nil {
		t.Fatalf("failed to create fixture directory: %v", err)
	}
	testFilePath := filepath.Join(relDir, "edge-non-binding-reference-role.json")
	if err := os.WriteFile(testFilePath, []byte(sampleRelationshipFixtureJSON), 0644); err != nil {
		t.Fatalf("failed to write fixture file: %v", err)
	}
	return filepath.Join(tempDir, "models"), testFilePath
}

func TestParseRelationshipJSONFile(t *testing.T) {
	_, testFilePath := createTestModelsFixture(t)

	row, err := ParseRelationshipJSONFile(testFilePath)
	if err != nil {
		t.Fatalf("ParseRelationshipJSONFile failed: %v", err)
	}

	if row.Model != "aws-ecs-controller" {
		t.Errorf("expected Model aws-ecs-controller, got %s", row.Model)
	}
	if row.Kind != "edge" {
		t.Errorf("expected Kind edge, got %s", row.Kind)
	}
	if row.Type != "non-binding" {
		t.Errorf("expected Type non-binding, got %s", row.Type)
	}
	if row.SubType != "reference" {
		t.Errorf("expected SubType reference, got %s", row.SubType)
	}
	if row.FromKind != "Service" {
		t.Errorf("expected FromKind Service, got %s", row.FromKind)
	}
	if row.ToKind != "Role" {
		t.Errorf("expected ToKind Role, got %s", row.ToKind)
	}
	if row.SchemaVersion != "relationships.meshery.io/v1beta2" {
		t.Errorf("expected SchemaVersion relationships.meshery.io/v1beta2, got %s", row.SchemaVersion)
	}
	if !strings.Contains(row.CompleteDefinition, `"schemaVersion":"relationships.meshery.io/v1beta2"`) {
		t.Errorf("expected CompleteDefinition to contain schemaVersion, got %s", row.CompleteDefinition)
	}
}

func TestScanAndExportRelationshipsCSV(t *testing.T) {
	modelsDir, _ := createTestModelsFixture(t)

	rows, err := ScanCommittedRelationships(modelsDir, "aws-ecs-controller")
	if err != nil {
		t.Fatalf("ScanCommittedRelationships failed: %v", err)
	}

	if len(rows) == 0 {
		t.Fatalf("expected at least 1 relationship row, got 0")
	}

	tempCSV := filepath.Join(t.TempDir(), "test_relationships.csv")
	err = ExportRelationshipsToCSV(rows, tempCSV)
	if err != nil {
		t.Fatalf("ExportRelationshipsToCSV failed: %v", err)
	}

	csvBytes, err := os.ReadFile(tempCSV)
	if err != nil {
		t.Fatalf("failed to read exported CSV: %v", err)
	}

	csvContent := string(csvBytes)
	if !strings.Contains(csvContent, "Meshery Relationship Definitions") {
		t.Errorf("expected CSV header 'Meshery Relationship Definitions'")
	}
	if !strings.Contains(csvContent, "aws-ecs-controller") {
		t.Errorf("expected CSV to contain 'aws-ecs-controller'")
	}
	if !strings.Contains(csvContent, "Service") || !strings.Contains(csvContent, "Role") {
		t.Errorf("expected CSV to contain component kinds Service and Role")
	}
}

func TestDeduplicateRelationshipRows(t *testing.T) {
	rows := []RelationshipRow{
		{
			Model:    "aws-ecs-controller",
			Filename: "edge-non-binding-reference-role.json",
			Kind:     "edge",
			Type:     "non-binding",
			SubType:  "reference",
			FromKind: "Service",
			ToKind:   "Role",
		},
		{
			Model:    "aws-ecs-controller",
			Filename: "hierarchical-parent-inventory-cluster.json",
			Kind:     "hierarchical",
			Type:     "parent",
			SubType:  "inventory",
			FromKind: "Service",
			ToKind:   "Cluster",
		},
	}

	// Simulate existing sheet values where the first row already exists
	existingSheetValues := [][]interface{}{
		{"Meshery Relationship Definitions Header 1"},
		{"Model", "Version", "kind", "type", "subType", "", "", "", "", "", "", "", "", "", "filename"},
		{"aws-ecs-controller", "v1.0.0", "edge", "non-binding", "reference", "", "", "", "", "", "", "", "", "", "edge-non-binding-reference-role.json"},
	}

	newRows := DeduplicateRelationshipRows(rows, existingSheetValues)
	if len(newRows) != 1 {
		t.Fatalf("expected 1 new row after deduplication, got %d", len(newRows))
	}

	if len(newRows[0]) < 21 {
		t.Fatalf("expected row to have at least 21 elements, got %d", len(newRows[0]))
	}

	if newRows[0][0] != "aws-ecs-controller" {
		t.Errorf("expected Model aws-ecs-controller, got %v", newRows[0][0])
	}
	if newRows[0][14] != "hierarchical-parent-inventory-cluster.json" {
		t.Errorf("expected filename hierarchical-parent-inventory-cluster.json, got %v", newRows[0][14])
	}
}

