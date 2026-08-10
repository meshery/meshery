package relationships

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	registrycmd "github.com/meshery/meshery/mesheryctl/internal/cli/root/registry"

	meshkitRegistry "github.com/meshery/meshkit/registry"
	meshkit "github.com/meshery/meshkit/utils"

	"github.com/spf13/cobra"
)

type cmdRelationshipGenerateFlag struct {
	SpreadsheetID   string `json:"spreadsheet-id" validate:"omitempty"`
	SpreadsheetCred string `json:"spreadsheet-cred" validate:"omitempty"`
	File            string `json:"file" validate:"omitempty,filepath"`
	Output          string `json:"output" validate:"omitempty"`
}

// minRelationshipCSVColumns defines the minimum number of columns required in a CSV row
// to be considered a valid relationship entry.
// Refer to the Meshery relationship CSV template:
// https://github.com/meshery/meshery/blob/master/mesheryctl/templates/template-csvs/Relationships.csv

const googleSpreadsheetURL = "https://docs.google.com/spreadsheets/d/"

var relationshipGenerateFlag cmdRelationshipGenerateFlag

var getRelationshipSheetID = func(cred, spreadsheetID string) (int64, error) {
	srv, err := meshkit.NewSheetSRV(cred)
	if err != nil {
		return 0, err
	}

	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Fields().Do()
	if err != nil {
		return 0, err
	}

	return registrycmd.GetSheetIDFromTitle(spreadsheet, "Relationships"), nil
}

var newRelationshipCSVHelper = meshkitRegistry.NewRelationshipCSVHelper

var relationshipsOutputPath = "../docs/data/RelationshipsData.json"

type CustomValueRange = meshkitRegistry.RelationshipCSV

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate relationships documents",
	Long:  "Generate relationships documents from a CSV file or Google Spreadsheet",
	Example: `
// Generate relationships documents from a CSV file
mesheryctl relationship generate --file <path-to-relationships.csv>

// Generate relationships documents with a custom output path
mesheryctl relationship generate --file <path-to-relationships.csv> --output <path-to-output.json>

// Generate relationships documents from a Google Spreadsheet
mesheryctl relationship generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &relationshipGenerateFlag)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		outputPath := relationshipGenerateFlag.Output
		if outputPath == "" {
			outputPath = relationshipsOutputPath
		}

		if relationshipGenerateFlag.File != "" {
			data, err := generateRelationshipsFromCSV(relationshipGenerateFlag.File)
			if err != nil {
				return err
			}
			return saveRelationshipsJSON(data, outputPath)
		}

		relationshipSheetID, err := getRelationshipSheetID(
			relationshipGenerateFlag.SpreadsheetCred,
			relationshipGenerateFlag.SpreadsheetID,
		)
		if err != nil {
			return err
		}

		if relationshipSheetID == -1 {
			return ErrEmptySheetData(fmt.Errorf("relationships sheet not found"))
		}

		helper, err := newRelationshipCSVHelper(
			googleSpreadsheetURL+relationshipGenerateFlag.SpreadsheetID,
			"Relationships",
			relationshipSheetID,
			"",
		)
		if err != nil {
			return err
		}

		if err := helper.ParseRelationshipsSheet(""); err != nil {
			return err
		}

		if len(helper.Relationships) == 0 {
			return ErrEmptySheetData(fmt.Errorf("no relationship data found in sheet"))
		}

		return saveRelationshipsJSON(helper.Relationships, outputPath)
	},
}

func init() {
	generateCmd.Flags().StringVarP(&relationshipGenerateFlag.File, "file", "f", "", "path to the relationships CSV file")
	generateCmd.Flags().StringVar(&relationshipGenerateFlag.SpreadsheetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	generateCmd.Flags().StringVar(&relationshipGenerateFlag.SpreadsheetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")
	generateCmd.Flags().StringVarP(&relationshipGenerateFlag.Output, "output", "o", "", "path to the output JSON file")

	generateCmd.MarkFlagsOneRequired("spreadsheet-id", "file")
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "file")
	generateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")
}

func generateRelationshipsFromCSV(filePath string) ([]CustomValueRange, error) {
	helper, err := meshkitRegistry.NewRelationshipCSVHelper("", "", 0, filePath)
	if err != nil {
		return nil, err
	}

	if err := helper.ParseRelationshipsSheet(""); err != nil {
		return nil, err
	}

	if len(helper.Relationships) == 0 {
		return nil, ErrEmptyCSVData(
			fmt.Errorf("no valid relationship rows found in CSV file: %s", filePath),
		)
	}

	return helper.Relationships, nil
}

func saveRelationshipsJSON(data []CustomValueRange, jsonFilePath string) error {
	jsonFormatter := display.NewJSONOutputFormatter(data).(*display.JSONOutputFormatter[[]CustomValueRange])
	saver := display.NewJSONOutputFormatterSaver(*jsonFormatter)
	return saver.
		WithFilePath(jsonFilePath).
		Save()
}
