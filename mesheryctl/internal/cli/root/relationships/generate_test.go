package relationships

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func TestExperimentalGenerate_CreateJsonFile(t *testing.T) {

	resp := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Relationships!A1:O1000",
		Values: [][]interface{}{
			{},
			{},
			{"kubernetes", "v1.25.2", "Hierarchical", "", "Inventory", "A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component.\"", "", "", "hierarchical_inventory_relationship", "", "", "{\"kind\": \"PersistentVolume\",\r\n            \"model\": \"kubernetes\",\r\n            \"patch\": {\r\n              \"patchStrategy\": \"replace\",\r\n              \"mutatorRef\": [\r\n                [\"name\"]\r\n              ],\r\n              \"description\": \"PersistentVolume in Kubernetes provides durable storage for Pods, ensuring data persistence across Pod restarts and rescheduling.\"\r\n            }\r\n          }",
				"{\"kind\": \"VolumeAttachment\",\r\n            \"model\": \"kubernetes\",\r\n            \"patch\": {\r\n              \"patchStrategy\": \"replace\",\r\n              \"mutatedRef\": [\r\n                [\"settings\", \"spec\", \"source\", \"persistentVolumeName\"]\r\n              ],\r\n              \"description\": \"VolumeAttachment in Kubernetes facilitates the attachment of a PersistentVolume to a node, enabling its utilization by Pods on that node.\"\r\n            }\r\n          }",
				"{\"apiVersion\": \"core.meshery.io/v1alpha1\",\r\"kind\":\"Hierarchical\",\r\"metadata\": {\r  \"description\": \"A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component.\"\",\r  \"styles\": { \"\",}\r  },\r  \"model\": {\r    \"name\": \"kubernetes\",\r    \"version\": \"v1.25.2\",\r    \"category\": {\r      \"name\": \"Orchestration \\u0026 Management\",\r      \"metadata\": null\r    },\r    \"metadata\": {}\r},\r\"subType\": \"Inventory\",\r\"rego_query\": \"hierarchical_inventory_relationship\",\r\"selectors\": {\r  \"allow\": {\r    \"from\": [{\"kind\": \"PersistentVolume\",\r\r            \"model\": \"kubernetes\",\r\r            \"patch\": {\r\r              \"patchStrategy\": \"replace\",\r\r              \"mutatorRef\": [\r\r                [\"name\"]\r\r              ],\r\r              \"description\": \"PersistentVolume in Kubernetes provides durable storage for Pods, ensuring data persistence across Pod restarts and rescheduling.\"\r\r            }\r\r          }],\r    \"to\": [{\"kind\": \"VolumeAttachment\",\r\r            \"model\": \"kubernetes\",\r\r            \"patch\": {\r\r              \"patchStrategy\": \"replace\",\r\r              \"mutatedRef\": [\r\r                [\"settings\", \"spec\", \"source\", \"persistentVolumeName\"]\r\r              ],\r\r              \"description\": \"VolumeAttachment in Kubernetes facilitates the attachment of a PersistentVolume to a node, enabling its utilization by Pods on that node.\"\r\r            }\r\r          }]\r    },\r  \"deny\": {\r    \"from\": [],\r    \"to\": []\r    }\r  }\r}",
				""},
		},
	}

	// for testing, relative path is required in createJsonFile function
	jsonFilePath := "../../../../../docs/_data/RelationshipsDataTest.json"

	err := createJsonFile(resp, jsonFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// comparing the generated JSON file with the expected JSON file
	goldenFilePath := "./fixtures/generate.relationship.json.output.golden"

	relationshipData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		t.Fatal(err)
	}

	expectedRelationshipData, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatal(err)
	}

	var relationshipDataJson, expectedRelationshipDataJson []map[string]interface{}
	err = json.Unmarshal([]byte(relationshipData), &relationshipDataJson)
	if err != nil {
		t.Fatal(err)
	}

	err = json.Unmarshal([]byte(expectedRelationshipData), &expectedRelationshipDataJson)
	if err != nil {
		t.Fatal(err)
	}

	assert.JSONEqf(t, string(expectedRelationshipData), string(relationshipData), "Generated JSON data does not match expected data.\n Difference: %s", cmp.Diff(relationshipData, expectedRelationshipDataJson))

	t.Log("Create JSON file test passed")

	err = os.Remove(jsonFilePath)
	if err != nil {
		t.Fatal("Error deleting file:", err)
		return
	}
}
