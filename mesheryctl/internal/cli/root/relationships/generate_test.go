package relationships

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistry "github.com/meshery/meshkit/registry"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestGenerateErrorOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	// test scenarios for fetching data
	tests := []struct {
		Name           string
		Args           []string
		ExpectError    bool
		ErrorSubstring string
		ExpectedError  error
	}{
		{
			Name:           "Given no input flags, when generate runs, then it errors with missing required flag",
			Args:           []string{"generate"},
			ExpectError:    true,
			ErrorSubstring: "at least one of the flags",
		},
		{
			Name:           "Given only spreadsheet-id, when generate runs, then it errors with missing spreadsheet-cred",
			Args:           []string{"generate", "--spreadsheet-id", "1"},
			ExpectError:    true,
			ErrorSubstring: "spreadsheet-cred",
		},
		{
			Name:           "Given nonexistent CSV file path, when generate runs, then it errors with file read error",
			Args:           []string{"generate", "--file", "/nonexistent/file.csv"},
			ExpectError:    true,
			ErrorSubstring: "no such file or directory",
		},
		{
			Name:           "Given both file and spreadsheet-id, when generate runs, then it errors with mutually exclusive flags",
			Args:           []string{"generate", "--file", "test.csv", "--spreadsheet-id", "1", "--spreadsheet-cred", "cred"},
			ExpectError:    true,
			ErrorSubstring: "if any flags in the group",
		},
	}

	// run tests
	for _, tt := range tests {
		tt := tt
		t.Run(tt.Name, func(t *testing.T) {
			defer utils.ResetCommandFlags(RelationshipCmd, t)
			defer utils.ResetCommandFlags(generateCmd, t)
			relationshipGenerateFlag = cmdRelationshipGenerateFlag{}

			generateCmd.Flags().VisitAll(func(f *pflag.Flag) {
				f.Changed = false
			})
			mesheryctlflags.InitValidators(RelationshipCmd)
			RelationshipCmd.SetArgs(tt.Args)
			err := RelationshipCmd.Execute()

			if err != nil {
				if tt.ExpectError {
					if tt.ErrorSubstring != "" {
						assert.Contains(t, err.Error(), tt.ErrorSubstring)
					}

					if tt.ExpectedError != nil {
						utils.AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					}

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

func TestGenerateSpreadsheetDataOutput(t *testing.T) {
	utils.SetupContextEnv(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test directory")
	}

	currDir := filepath.Dir(filename)

	originalGetSheetID := getRelationshipSheetID
	originalNewHelper := newRelationshipCSVHelper
	originalOutputPath := relationshipsOutputPath

	defer func() {
		getRelationshipSheetID = originalGetSheetID
		newRelationshipCSVHelper = originalNewHelper
		relationshipsOutputPath = originalOutputPath
	}()

	relationshipsOutputPath = filepath.Join(
		t.TempDir(),
		"relationships.json",
	)

	// Mock spreadsheet metadata lookup.
	getRelationshipSheetID = func(cred, spreadsheetID string) (int64, error) {
		assert.Equal(t, "$CRED", cred)
		assert.Equal(t, "test-spreadsheet-id", spreadsheetID)
		return 1410291737, nil
	}

	// Mock MeshKit relationship helper.
	newRelationshipCSVHelper = func(
		sheetURL string,
		sheetName string,
		sheetID int64,
		localCSVPath string,
	) (*meshkitRegistry.RelationshipCSVHelper, error) {
		assert.Equal(
			t,
			"https://docs.google.com/spreadsheets/d/test-spreadsheet-id",
			sheetURL,
		)
		assert.Equal(t, "Relationships", sheetName)
		assert.Equal(t, int64(1410291737), sheetID)
		assert.Empty(t, localCSVPath)

		fixturesDir := filepath.Join(currDir, "fixtures")

		return &meshkitRegistry.RelationshipCSVHelper{
			SpreadsheetID:  sheetID,
			SpreadsheetURL: sheetURL,
			Title:          sheetName,
			CSVPath: filepath.Join(
				fixturesDir,
				"generate.relationship.csv.data.golden",
			),
			Relationships: []meshkitRegistry.RelationshipCSV{},
		}, nil
	}

	defer utils.ResetCommandFlags(RelationshipCmd, t)

	mesheryctlflags.InitValidators(RelationshipCmd)

	RelationshipCmd.SetArgs([]string{
		"generate",
		"--spreadsheet-id", "test-spreadsheet-id",
		"--spreadsheet-cred", "$CRED",
		"--output", relationshipsOutputPath,
	})

	err := RelationshipCmd.Execute()
	if err != nil {
		t.Fatalf("generate command failed: %v", err)
	}

	actual, err := os.ReadFile(relationshipsOutputPath)
	if err != nil {
		t.Fatalf("failed to read generated JSON: %v", err)
	}

	expected := `[
  {
    "Model": "kubernetes",
    "Version": "v1.25.2",
    "kind": "Hierarchical",
    "type": "",
    "status": "",
    "subType": "Inventory",
    "PublishToRegistry": "",
    "metadata.description": "A hierarchical inventory relationship in which the configuration of (parent) component is patched with the configuration of other (child) component. Eg: The configuration of the EnvoyFilter (parent) component is patched with the configuration as received from WASMFilter (child) component.",
    "metadata.isAnnotation": "",
    "metadata.styles": "",
    "evalPolicy": "hierarchical_inventory_relationship",
    "selector": "",
    "filename": ""
  }
]`

	assert.JSONEq(t, expected, string(actual))
}

func TestGenerateCSVDataOutput(t *testing.T) {
	// setup current context
	utils.SetupContextEnv(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to determine test directory")
	}

	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	csvPath := filepath.Join(
		fixturesDir,
		"generate.relationship.csv.data.golden",
	)

	outputPath := filepath.Join(t.TempDir(), "output.json")

	// Initialize Meshery logger used by display package.
	utils.SetupMeshkitLoggerTesting(t, false)

	defer utils.ResetCommandFlags(RelationshipCmd, t)

	mesheryctlflags.InitValidators(RelationshipCmd)

	RelationshipCmd.SetArgs([]string{
		"generate",
		"--file", csvPath,
		"--output", outputPath,
	})

	err := RelationshipCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	actual, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read generated JSON: %v", err)
	}

	expectedPath := filepath.Join(
		currDir,
		"testdata",
		"relationships-data-csv-test.json",
	)

	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("failed to read expected JSON: %v", err)
	}

	assert.JSONEqf(
		t,
		string(expected),
		string(actual),
		"Generated JSON does not match expected data.\nDifference: %s",
		cmp.Diff(expected, actual),
	)
}
