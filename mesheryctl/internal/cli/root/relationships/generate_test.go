package relationships

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func expectedViewFlagError(spreadsheetId string, spreadsheetCred string) error {
	fv := mesheryctlflags.NewFlagValidator()
	return fv.Validate(&cmdRelationshipGenerateFlag{SpreadsheetID: spreadsheetId, SpreadsheetCred: spreadsheetCred})
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
		IsOutputGolden   bool
		ExpectedError    error
	}{

		{
			Name:           "Generate registered relationships without spreadsheet credentials",
			Args:           []string{"generate", "--spreadsheet-id", "1"},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  expectedViewFlagError("1", ""),
		},
		{
			Name:           "Generate registered relationships without spreadsheet id",
			Args:           []string{"generate", "--spreadsheet-cred", "$CRED"},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  expectedViewFlagError("", "$CRED"),
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			defer func() {
				relationshipGenerateFlag.SpreadsheetCred = ""
				relationshipGenerateFlag.SpreadsheetID = ""
			}()

			mesheryctlflags.InitValidators(RelationshipCmd)
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
		IsOutputGolden   bool
		ExpectedError    error
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
				relationshipGenerateFlag.SpreadsheetCred = ""
				relationshipGenerateFlag.SpreadsheetID = ""
			}()

			fixturesDir := filepath.Join(currDir, "fixtures")

			originalStdout := os.Stdout
			b := utils.SetupMeshkitLoggerTesting(t, false)
			defer func() {
				os.Stdout = originalStdout
			}()

			originalPath := relationshipsOutputPath
			defer func() { relationshipsOutputPath = originalPath }()
			relationshipsOutputPath = "./testdata/generate.relationship.json.output.golden"

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

			mesheryctlflags.InitValidators(RelationshipCmd)
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

			// Validate generated json file
			relationshipData, err := os.ReadFile(relationshipsOutputPath)
			if err != nil {
				t.Fatal("Error in reading file 'relationships-data-test.json': ", err)
			}

			jsonFilePath := "./testdata/relationships-data-test.json"
			expectedRelationshipData, err := os.ReadFile(jsonFilePath)
			if err != nil {
				t.Fatal("Error in reading file 'generate.relationship.json.output.golden': ", err)
			}
			assert.JSONEqf(t, string(expectedRelationshipData), string(relationshipData), "Generated JSON data does not match expected data.\n Difference: %s", cmp.Diff(relationshipData, expectedRelationshipData))

		})
		t.Log("Generate experimental relationship test for sheetdata validation has passed")
	}
}
