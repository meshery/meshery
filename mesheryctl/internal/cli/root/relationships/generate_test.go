package relationships

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"google.golang.org/api/sheets/v4"
)

func TestGenerateCreateJsonFile(t *testing.T) {
	row3 := []interface{}{
		"kubernetes",
		"v1.25.2",
		"Hierarchical",
		"",
		"Inventory",
		"A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component.\"",
		"",
		"",
		"hierarchical_inventory_relationship",
		"",
		"",
		"{\"kind\": \"PersistentVolume\",\r\n            \"model\": \"kubernetes\",\r\n            \"patch\": {\r\n              \"patchStrategy\": \"replace\",\r\n              \"mutatorRef\": [\r\n                [\"name\"]\r\n              ],\r\n              \"description\": \"PersistentVolume in Kubernetes provides durable storage for Pods, ensuring data persistence across Pod restarts and rescheduling.\"\r\n            }\r\n          }",
		"{\"kind\": \"VolumeAttachment\",\r\n            \"model\": \"kubernetes\",\r\n            \"patch\": {\r\n              \"patchStrategy\": \"replace\",\r\n              \"mutatedRef\": [\r\n                [\"settings\", \"spec\", \"source\", \"persistentVolumeName\"]\r\n              ],\r\n              \"description\": \"VolumeAttachment in Kubernetes facilitates the attachment of a PersistentVolume to a node, enabling its utilization by Pods on that node.\"\r\n            }\r\n          }",
		"{\"apiVersion\": \"core.meshery.io/v1alpha1\",\r\"kind\":\"Hierarchical\",\r\"metadata\": {\r  \"description\": \"A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component.\"\",\r  \"styles\": { \"\",}\r  },\r  \"model\": {\r    \"name\": \"kubernetes\",\r    \"version\": \"v1.25.2\",\r    \"category\": {\r      \"name\": \"Orchestration \\u0026 Management\",\r      \"metadata\": null\r    },\r    \"metadata\": {}\r},\r\"subType\": \"Inventory\",\r\"rego_query\": \"hierarchical_inventory_relationship\",\r\"selectors\": {\r  \"allow\": {\r    \"from\": [{\"kind\": \"PersistentVolume\",\r\r            \"model\": \"kubernetes\",\r\r            \"patch\": {\r\r              \"patchStrategy\": \"replace\",\r\r              \"mutatorRef\": [\r\r                [\"name\"]\r\r              ],\r\r              \"description\": \"PersistentVolume in Kubernetes provides durable storage for Pods, ensuring data persistence across Pod restarts and rescheduling.\"\r\r            }\r\r          }],\r    \"to\": [{\"kind\": \"VolumeAttachment\",\r\r            \"model\": \"kubernetes\",\r\r            \"patch\": {\r\r              \"patchStrategy\": \"replace\",\r\r              \"mutatedRef\": [\r\r                [\"settings\", \"spec\", \"source\", \"persistentVolumeName\"]\r\r              ],\r\r              \"description\": \"VolumeAttachment in Kubernetes facilitates the attachment of a PersistentVolume to a node, enabling its utilization by Pods on that node.\"\r\r            }\r\r          }]\r    },\r  \"deny\": {\r    \"from\": [],\r    \"to\": []\r    }\r  }\r}",
		"",
	}

	resp := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Relationships!A1:O1000",
		Values: [][]interface{}{
			{},
			{},
			row3,
		},
	}

	jsonFilePath := t.TempDir() + "/RelationshipDataTest.json"

	if err := createJsonFile(resp, jsonFilePath); err != nil {
		t.Fatal("Error in createJsonFile function: ", err)
	}

	generated, err := os.ReadFile(jsonFilePath)
	if err != nil {
		t.Fatalf("Error in reading generated JSON file: %v", err)
	}

	var got []CustomValueRange
	if err := json.Unmarshal(generated, &got); err != nil {
		t.Fatalf("Error in parsing generated JSON file: %v", err)
	}

	expected := []CustomValueRange{
		{
			Model:                row3[0].(string),
			Version:              row3[1].(string),
			Kind:                 row3[2].(string),
			Type:                 row3[3].(string),
			SubType:              row3[4].(string),
			MetadataDescription:  row3[5].(string),
			Docs:                 row3[6].(string),
			MetadataStyles:       row3[7].(string),
			EvalPolicy:           row3[8].(string),
			SelectorsDenyFrom:    row3[9].(string),
			SelectorsDenyTo:      row3[10].(string),
			SelectorsAllowFrom:   row3[11].(string),
			SelectorsAllowTo:     row3[12].(string),
			CompleteDefinition:   row3[13].(string),
			VisualizationExample: row3[14].(string),
		},
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("Generated JSON data does not match expected (-want +got):\n%s", diff)
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		Name                string
		Args                []string
		ExpectedErrContains []string
		ExpectError         bool
	}{
		{
			Name:             "Generate registered relationships without spreadsheet id",
			Args:             []string{"generate", "--spreadsheet-cred", "$CRED"},
			ExpectedErrContains: []string{"[ Spreadsheet ID | Spreadsheet Credentials ] aren't specified"},
			ExpectError:         true,
		},
		{
			Name:             "Generate registered relationships without spreadsheet creadentials",
			Args:             []string{"generate", "--spreadsheet-id", "1"},
			ExpectedErrContains: []string{"[ Spreadsheet ID | Spreadsheet Credentials ] aren't specified"},
			ExpectError:         true,
		},
		{
			Name:             "Generate registered relationships",
			Args:             []string{"generate", "--spreadsheet-cred", "$CRED", "--spreadsheet-id", "1"},
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// Execute command and validate errors (if any). No golden snapshots.
			b := utils.SetupMeshkitLoggerTesting(t, false)
			RelationshipCmd.SetArgs(tt.Args)
			RelationshipCmd.SetOut(b)
			err := RelationshipCmd.Execute()
			if err != nil {
				if !tt.ExpectError {
					t.Fatal(err)
				}
				for _, s := range tt.ExpectedErrContains {
					if !strings.Contains(err.Error(), s) {
						t.Fatalf("expected error to contain %q, got %q", s, err.Error())
					}
				}
				spreadsheeetCred = ""
				spreadsheeetID = ""
				return
			}
			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}
			// reset the global variables
			spreadsheeetCred = ""
			spreadsheeetID = ""
		})
	}
}
