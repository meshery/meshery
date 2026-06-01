package relationships

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkit "github.com/meshery/meshkit/utils"
	meshkitcsv "github.com/meshery/meshkit/utils/csv"
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

var relationshipsOutputPath = "../docs/data/RelationshipsData.json"

type CustomValueRange struct {
	Model                string `json:"Model" csv:"Model"`
	Version              string `json:"Version" csv:"Version"`
	Kind                 string `json:"kind" csv:"Kind"`
	Type                 string `json:"type" csv:"Type"`
	SubType              string `json:"subType" csv:"SubType"`
	MetadataDescription  string `json:"metadataDescription" csv:"MetadataDescription"`
	Docs                 string `json:"docs" csv:"Docs"`
	MetadataStyles       string `json:"metadataStyles" csv:"MetadataStyles"`
	EvalPolicy           string `json:"evalPolicy" csv:"EvalPolicy"`
	SelectorsDenyFrom    string `json:"selectorsDenyFrom" csv:"SelectorsDenyFrom"`
	SelectorsDenyTo      string `json:"selectorsDenyTo" csv:"SelectorsDenyTo"`
	SelectorsAllowFrom   string `json:"selectorsAllowFrom" csv:"SelectorsAllowFrom"`
	SelectorsAllowTo     string `json:"selectorsAllowTo" csv:"SelectorsAllowTo"`
	CompleteDefinition   string `json:"CompleteDefinition" csv:"CompleteDefinition"`
	VisualizationExample string `json:"VisualizationExample" csv:"VisualizationExample"`
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
	ch := make(chan CustomValueRange, 1)
	errorChan := make(chan error, 1)

	csvParser, err := meshkitcsv.NewCSVParser[CustomValueRange](filePath, 1, nil, func(_ []string, row []string) bool {
		return len(row) > 0 && row[0] != ""
	})
	if err != nil {
		return nil, err
	}

	var customResp []CustomValueRange

	go func() {
		if err := csvParser.Parse(ch, errorChan); err != nil {
			errorChan <- err
		}
	}()

	for {
		select {
		case data := <-ch:
			customResp = append(customResp, data)
		case err := <-errorChan:
			utils.Log.Error(err)
		case <-csvParser.Context.Done():
			if len(customResp) == 0 {
				return nil, ErrEmptyCSVData(fmt.Errorf("no valid relationship rows found in CSV file: %s", filePath))
			}
			return customResp, nil
		}
	}
}

func processSheetData(resp *sheets.ValueRange, jsonFilePath string) error {
	var customResp []CustomValueRange

	// Row index 1 = headers
	headers := resp.Values[1]
	colIndex := map[string]int{}
	for i, h := range headers {
		colIndex[fmt.Sprintf("%v", h)] = i
	}

	getCol := func(row []interface{}, name string) string {
		idx, ok := colIndex[name]
		if !ok || idx >= len(row) {
			return ""
		}
		return fmt.Sprintf("%v", row[idx])
	}

	for _, row := range resp.Values[2:] {
		if getCol(row, "Model") == "" {
			continue
		}
		customResp = append(customResp, CustomValueRange{
			Model:                getCol(row, "Model"),
			Version:              getCol(row, "Version"),
			Kind:                 getCol(row, "Kind"),
			Type:                 getCol(row, "Type"),
			SubType:              getCol(row, "SubType"),
			MetadataDescription:  getCol(row, "MetadataDescription"),
			Docs:                 getCol(row, "Docs"),
			MetadataStyles:       getCol(row, "MetadataStyles"),
			EvalPolicy:           getCol(row, "EvalPolicy"),
			SelectorsDenyFrom:    getCol(row, "SelectorsDenyFrom"),
			SelectorsDenyTo:      getCol(row, "SelectorsDenyTo"),
			SelectorsAllowFrom:   getCol(row, "SelectorsAllowFrom"),
			SelectorsAllowTo:     getCol(row, "SelectorsAllowTo"),
			CompleteDefinition:   getCol(row, "CompleteDefinition"),
			VisualizationExample: getCol(row, "VisualizationExample"),
		})
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
