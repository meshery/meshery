package relationships

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	meshkitRegistry "github.com/meshery/meshkit/registry"
	meshkit "github.com/meshery/meshkit/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
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
const minRelationshipCSVColumns = 15

var relationshipGenerateFlag cmdRelationshipGenerateFlag

var fetchSheetValues = func(id, cred string) (*sheets.ValueRange, error) {
	srv, err := meshkit.NewSheetSRV(cred)
	if err != nil {
		return nil, err
	}
	return srv.Spreadsheets.Values.Get(id, "Relationships").Do()
}

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

		resp, err := fetchSheetValues(relationshipGenerateFlag.SpreadsheetID, relationshipGenerateFlag.SpreadsheetCred)
		if err != nil {
			return err
		}

		if len(resp.Values) <= 2 {
			return ErrEmptySheetData(fmt.Errorf("no relationship data found in sheet"))
		}

		return processSheetData(resp, outputPath)
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

func processSheetData(resp *sheets.ValueRange, jsonFilePath string) error {
	var customResp []CustomValueRange

	for _, row := range resp.Values[2:] {
		if len(row) >= minRelationshipCSVColumns && row[0] != "" {
			customResp = append(customResp, CustomValueRange{
				Model:       row[0].(string),
				Version:     row[1].(string),
				KIND:        row[2].(string),
				Type:        row[3].(string),
				SubType:     row[4].(string),
				Description: row[5].(string),
				Styles:      row[7].(string),
				EvalPolicy:  row[8].(string),
			})
		}
	}

	return saveRelationshipsJSON(customResp, jsonFilePath)
}

func saveRelationshipsJSON(data []CustomValueRange, jsonFilePath string) error {
	jsonFormatter := display.NewJSONOutputFormatter(data).(*display.JSONOutputFormatter[[]CustomValueRange])
	saver := display.NewJSONOutputFormatterSaver(*jsonFormatter)
	return saver.
		WithFilePath(jsonFilePath).
		Save()
}
