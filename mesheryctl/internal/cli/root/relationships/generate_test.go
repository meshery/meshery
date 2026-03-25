package relationships

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/sheets/v4"
)

func resetGenerateFlags() {
	relationshipGenerateFlag.SpreadsheetCred = ""
	relationshipGenerateFlag.SpreadsheetID = ""
	relationshipGenerateFlag.File = ""
	relationshipGenerateFlag.Output = ""
	// Reset Cobra's internal flag state to prevent bleed between tests
	generateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
		f.Value.Set(f.DefValue)
	})
}

func TestGenerateErrorOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// test scenarios for fetching data
	tests := []struct {
		Name           string
		Args           []string
		ExpectError    bool
		ErrorSubstring string
	}{
		{
			Name:           "Generate without any input source",
			Args:           []string{"generate"},
			ExpectError:    true,
			ErrorSubstring: "at least one of the flags",
		},
		{
			Name:           "Generate with spreadsheet-id but without spreadsheet-cred",
			Args:           []string{"generate", "--spreadsheet-id", "1"},
			ExpectError:    true,
			ErrorSubstring: "spreadsheet-cred",
		},
		{
			Name:           "Generate with both file and spreadsheet-id",
			Args:           []string{"generate", "--file", "test.csv", "--spreadsheet-id", "1", "--spreadsheet-cred", "cred"},
			ExpectError:    true,
			ErrorSubstring: "if any flags in the group",
		},
		{
			Name:           "Generate with nonexistent CSV file",
			Args:           []string{"generate", "--file", "/nonexistent/file.csv"},
			ExpectError:    true,
			ErrorSubstring: "File read error",
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer resetGenerateFlags()

			mesheryctlflags.InitValidators(RelationshipCmd)
			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()

			if err != nil {
				if tt.ExpectError {
					assert.True(t, strings.Contains(err.Error(), tt.ErrorSubstring),
						"expected error containing %q, got: %v", tt.ErrorSubstring, err)
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}
		})
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

			defer resetGenerateFlags()

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

func TestGenerateCSVDataOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	tests := []struct {
		Name        string
		CSVFixture  string
		ExpectError bool
	}{
		{
			Name:        "Generate relationships from CSV file",
			CSVFixture:  "generate.relationship.csv.data.golden",
			ExpectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer resetGenerateFlags()

			fixturesDir := filepath.Join(currDir, "fixtures")
			csvPath := filepath.Join(fixturesDir, tt.CSVFixture)

			outputPath := filepath.Join(t.TempDir(), "output.json")

			originalStdout := os.Stdout
			_ = utils.SetupMeshkitLoggerTesting(t, false)
			defer func() {
				os.Stdout = originalStdout
			}()

			mesheryctlflags.InitValidators(RelationshipCmd)
			RelationshipCmd.SetArgs([]string{"generate", "--file", csvPath, "--output", outputPath})
			RelationshipCmd.SetOut(originalStdout)
			err := RelationshipCmd.Execute()

			if err != nil {
				if tt.ExpectError {
					return
				}
				t.Fatal(err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			// Validate generated json file matches expected output
			relationshipData, err := os.ReadFile(outputPath)
			if err != nil {
				t.Fatal("Error reading generated JSON file:", err)
			}

			jsonFilePath := "./testdata/relationships-data-csv-test.json"
			expectedRelationshipData, err := os.ReadFile(jsonFilePath)
			if err != nil {
				t.Fatal("Error reading expected JSON file:", err)
			}
			assert.JSONEqf(t, string(expectedRelationshipData), string(relationshipData), "Generated JSON data from CSV does not match expected data.\n Difference: %s", cmp.Diff(relationshipData, expectedRelationshipData))
		})
		t.Log("Generate CSV relationship test passed")
	}
}
