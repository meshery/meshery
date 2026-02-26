package relationships

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func TestGenerateCreateJsonFile(t *testing.T) {

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
	jsonFilePath := "./testdata/relationships-data-test.json"

	err = createJsonFile(resp, jsonFilePath)
	if err != nil {
		t.Fatal("Error in createJsonFile function: ", err)
	}

	// comparing the generated JSON file with the expected JSON file
	goldenFilePath := "./fixtures/generate.relationship.json.output.golden"

	relationshipData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		t.Fatal("Error in reading file 'relationships-data-test.json': ", err)
	}

	expectedRelationshipData, err := os.ReadFile(goldenFilePath)
	if err != nil {
		t.Fatal("Error in reading file 'generate.relationship.json.output.golden': ", err)
	}

	var relationshipDataJson, expectedRelationshipDataJson []map[string]interface{}
	err = json.Unmarshal([]byte(relationshipData), &relationshipDataJson)
	if err != nil {
		t.Fatal("Error in parsing file 'relationships-data-test.json': ", err)
	}

	err = json.Unmarshal([]byte(expectedRelationshipData), &expectedRelationshipDataJson)
	if err != nil {
		t.Fatal("Error in parsing file 'relationships-data-test.json': ", err)
	}

	assert.JSONEqf(t, string(expectedRelationshipData), string(relationshipData), "Generated JSON data does not match expected data.\n Difference: %s", cmp.Diff(relationshipData, expectedRelationshipDataJson))

	t.Log("Create JSON file test passed")
}

func TestGenerateErrorOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		Fixture          string
		ExpectedResponse string
		ExpectError      bool
		IsOutputGolden   bool  `default:"true"`
		ExpectedError    error `default:"nil"`
	}{

		{
			Name:           "Generate registered relationships without spreadsheet creadentials",
			Args:           []string{"generate", "--spreadsheet-id", "1"},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrSpreadsheetFlagMissing(fmt.Errorf("%s", errMsg)),
		},
		{
			Name:           "Generate registered relationships without spreadsheet id",
			Args:           []string{"generate", "--spreadsheet-cred", "$CRED"},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrSpreadsheetFlagMissing(fmt.Errorf("%s", errMsg)),
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			defer func() {
				spreadsheetCred = ""
				spreadsheetID = ""
			}()

			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()

			// to validate the expected errors
			if err != nil {
				if tt.ExpectError {

					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}

				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

		})
		t.Log("Generate experimental relationship test for argument validation has passed")
	}
}

func TestGenerateDataOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// test scenarios for fetching data
	tests := []struct {
		Name             string
		Args             []string
		Fixture          string
		ExpectedResponse string
		ExpectError      bool
		IsOutputGolden   bool  `default:"true"`
		ExpectedError    error `default:"nil"`
	}{
		{
			Name:             "Generate registered relationships",
			Args:             []string{"generate", "--spreadsheet-cred", "$CRED", "--spreadsheet-id", "1"},
			Fixture:          "generate.relationship.sheet.data.golden",
			ExpectedResponse: "generate.relationship.output.golden",
			IsOutputGolden:   true,
			ExpectError:      false,
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			defer func() {
				spreadsheetCred = ""
				spreadsheetID = ""
			}()

			fixturesDir := filepath.Join(currDir, "fixtures")

			originalStdout := os.Stdout
			b := utils.SetupMeshkitLoggerTesting(t, false)
			defer func() {
				os.Stdout = originalStdout
			}()

			originalPath := relationshipsOutputPath
			defer func() { relationshipsOutputPath = originalPath }()
			relationshipsOutputPath = "./fixtures/generate.relationship.json.output.golden"

			originalFetch := fetchSheetValues
			defer func() { fetchSheetValues = originalFetch }()

			golden := utils.NewGoldenFile(t, tt.Fixture, fixturesDir)
			sheetData := golden.Load()

			var sheetDataParsed map[string]interface{}
			if err := json.Unmarshal([]byte(sheetData), &sheetDataParsed); err != nil {
				t.Fatal("Error parsing golden sheet data:", err)
			}

			fetchSheetValues = func(id, cred string) (*sheets.ValueRange, error) {
				return &sheets.ValueRange{
					MajorDimension: "ROWS",
					Range:          "Relationships!A1:O1000",
					Values: [][]interface{}{
						{},
						{},
						sheetDataParsed["ROW3"].([]interface{}),
					},
				}, nil
			}

			RelationshipCmd.SetArgs(tt.Args)
			RelationshipCmd.SetOut(originalStdout)
			err := RelationshipCmd.Execute()

			// to validate the expected errors
			if err != nil {
				if tt.ExpectError {

					utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}

				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			actualResponse := b.String()

			testdataDir := filepath.Join(currDir, "testdata")
			golden = utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := utils.CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := utils.CleanStringFromHandlePagination(expectedResponse)

			utils.Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Log("Generate experimental relationship test for sheetdata validation has passed")
	}
}
