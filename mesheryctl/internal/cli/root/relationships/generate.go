package relationships

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
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

var relationshipGenerateFlag cmdRelationshipGenerateFlag

var fetchSheetValues = func(id, cred string) (*sheets.ValueRange, error) {
	srv, err := meshkit.NewSheetSRV(cred)
	if err != nil {
		return nil, err
	}
	return srv.Spreadsheets.Values.Get(id, "Relationships").Do()
}

var relationshipsOutputPath = "docs/data/RelationshipsData.json"

type CustomValueRange struct {
	Model                string `json:"Model"`
	Version              string `json:"Version"`
	Kind                 string `json:"kind"`
	Type                 string `json:"type"`
	SubType              string `json:"subType"`
	MetadataDescription  string `json:"metadataDescription"`
	Docs                 string `json:"docs"`
	MetadataStyles       string `json:"metadataStyles"`
	EvalPolicy           string `json:"evalPolicy"`
	SelectorsDenyFrom    string `json:"selectorsDenyFrom"`
	SelectorsDenyTo      string `json:"selectorsDenyTo"`
	SelectorsAllowFrom   string `json:"selectorsAllowFrom"`
	SelectorsAllowTo     string `json:"selectorsAllowTo"`
	CompleteDefinition   string `json:"CompleteDefinition"`
	VisualizationExample string `json:"VisualizationExample"`
}

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
			return generateFromCSV(relationshipGenerateFlag.File, outputPath)
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
	generateCmd.Flags().StringVarP(&relationshipGenerateFlag.Output, "output", "o", "", "path to the output JSON file (default: docs/data/RelationshipsData.json)")

	generateCmd.MarkFlagsOneRequired("spreadsheet-id", "file")
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "file")
	generateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")
}

func generateFromCSV(filePath, outputPath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return utils.ErrFileRead(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return utils.ErrFileRead(err)
	}

	// First two rows are headers, data starts at row 3
	if len(records) <= 2 {
		return ErrEmptyCSVData(fmt.Errorf("no relationship data found in CSV file: %s", filePath))
	}

	var customResp []CustomValueRange
	for _, row := range records[2:] {
		if len(row) >= 15 && row[0] != "" {
			customResp = append(customResp, CustomValueRange{
				Model:                row[0],
				Version:              row[1],
				Kind:                 row[2],
				Type:                 row[3],
				SubType:              row[4],
				MetadataDescription:  row[5],
				Docs:                 row[6],
				MetadataStyles:       row[7],
				EvalPolicy:           row[8],
				SelectorsDenyFrom:    row[9],
				SelectorsDenyTo:      row[10],
				SelectorsAllowFrom:   row[11],
				SelectorsAllowTo:     row[12],
				CompleteDefinition:   row[13],
				VisualizationExample: row[14],
			})
		}
	}

	if len(customResp) == 0 {
		return ErrEmptyCSVData(fmt.Errorf("no valid relationship rows found in CSV file: %s", filePath))
	}

	return saveRelationshipsJSON(customResp, outputPath)
}

func processSheetData(resp *sheets.ValueRange, jsonFilePath string) error {
	var customResp []CustomValueRange

	for _, row := range resp.Values[2:] {
		if len(row) >= 15 && row[0] != "" {
			customResp = append(customResp, CustomValueRange{
				Model:                row[0].(string),
				Version:              row[1].(string),
				Kind:                 row[2].(string),
				Type:                 row[3].(string),
				SubType:              row[4].(string),
				MetadataDescription:  row[5].(string),
				Docs:                 row[6].(string),
				MetadataStyles:       row[7].(string),
				EvalPolicy:           row[8].(string),
				SelectorsDenyFrom:    row[9].(string),
				SelectorsDenyTo:      row[10].(string),
				SelectorsAllowFrom:   row[11].(string),
				SelectorsAllowTo:     row[12].(string),
				CompleteDefinition:   row[13].(string),
				VisualizationExample: row[14].(string),
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
