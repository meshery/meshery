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

	sheetData, err := os.ReadFile("./fixtures/generate.relationship.sheet.data.golden")
	if err != nil {
		t.Fatal("Error in reading file 'generate.relationship.sheet.data.golden': ", err)
	}

	var sheetDataParsed map[string]interface{}
	err = json.Unmarshal(sheetData, &sheetDataParsed)
	if err != nil {
		t.Fatal("Error in parsing file 'generate.relationship.sheet.data.golden': ", err)
	}

	resp := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Range:          "Relationships!A1:O1000",
		Values: [][]interface{}{
			{},
			{},
			sheetDataParsed["ROW3"].([]interface{}),
		},
	}

	// for testing, relative path is required in createJsonFile function
	jsonFilePath := "./testdata/RelationshipDataTest.json"

	err = createJsonFile(resp, jsonFilePath)
	if err != nil {
		t.Fatal("Error in createJsonFile function: ", err)
	}

	// comparing the generated JSON file with the expected JSON file
	goldenFilePath := "./fixtures/generate.relationship.json.output.golden"

	relationshipData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		t.Fatal("Error in reading file 'RelationshipsDataTest.json': ", err)
	}

	expectedRelationshipData, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatal("Error in reading file 'generate.relationship.json.output.golden': ", err)
	}

	var relationshipDataJson, expectedRelationshipDataJson []map[string]interface{}
	err = json.Unmarshal([]byte(relationshipData), &relationshipDataJson)
	if err != nil {
		t.Fatal("Error in parsing file 'RelationshipsDataTest.json': ", err)
	}

	err = json.Unmarshal([]byte(expectedRelationshipData), &expectedRelationshipDataJson)
	if err != nil {
		t.Fatal("Error in parsing file 'RelationshipsDataTest.json': ", err)
	}

	assert.JSONEqf(t, string(expectedRelationshipData), string(relationshipData), "Generated JSON data does not match expected data.\n Difference: %s", cmp.Diff(relationshipData, expectedRelationshipDataJson))

	t.Log("Create JSON file test passed")
}
