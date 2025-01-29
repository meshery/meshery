package relationships

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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

func TestGenerate(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	//initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		Token            string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Generate registered relationships",
			Args:             []string{"generate", "$CRED", "", "-s", "1"},
			URL:              testContext.BaseURL + "/api/meshmodels/relationships",
			Fixture:          "generate.relationship.api.response.golden",
			ExpectedResponse: "generate.relationship.output.golden",
			Token:            filepath.Join(fixturesDir, "token.golden"),
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			utils.TokenFlag = tt.Token

			httpmock.RegisterResponder("GET", tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(currDir, "testdata")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			RelationshipCmd.SetArgs(tt.Args)
			RelationshipCmd.SetOutput(rescueStdout)
			err := RelationshipCmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// write it in file
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					utils.Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			actualResponse := string(out)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Log("Generate experimental relationship test passed")
	}

	utils.StopMockery(t)
}

func TestGenerateArguments(t *testing.T) {

	const errMsg = "Usage: mesheryctl exp relationship generate $CRED [google-sheets-credential] --sheetId [sheet-id]\nRun 'mesheryctl exp relationship generate --help' to see detailed help message"

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		ExpectedResponse string
		ExpectError      bool
	}{
		{
			Name:             "Missing Credentials",
			Args:             []string{"generate"},
			ExpectedResponse: "Google Sheet Credentials is required\n" + errMsg,
			ExpectError:      false,
		},
		{
			Name:             "Missing Sheet ID",
			Args:             []string{"generate", "$CRED", "--sheetId", ""},
			ExpectedResponse: "Sheet ID is required\n" + errMsg,
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()

			utils.Equals(t, errors.New(utils.RelationshipsError(tt.ExpectedResponse, "generate")), err)
		})
		t.Log("Generate relationships test passed")
	}
}
