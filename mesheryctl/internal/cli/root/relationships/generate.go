package relationships

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkit "github.com/meshery/meshkit/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

var (
	spreadsheeetID   string
	spreadsheeetCred string
)

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
	SelectorsAllowFrom   string `json:"selectorsAllwowFrom"`
	SelectorsAllowTo     string `json:"selectorsAllowTo"`
	CompleteDefinition   string `json:"CompleteDefinition"`
	VisualizationExample string `json:"VisualizationExample"`
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate relationships documents",
	Long:  "Generate relationships documents from the google spreadsheets",
	Example: `
// Generate relationships documentss
mesheryctl exp relationship generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED
`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "[ Spreadsheet ID | Spreadsheet Credentials ] aren't specified\n\nUsage: mesheryctl exp relationship generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nRun 'mesheryctl exp relationship generate --help' to see detailed help message"

		// Check if flag is set
		if spreadsheeetID == "" || spreadsheeetCred == "" {
			return errors.New(utils.RelationshipsError(errMsg, "generate"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, err := meshkit.NewSheetSRV(spreadsheeetCred)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		resp, err := srv.Spreadsheets.Values.Get(spreadsheeetID, "Relationships").Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(err)
			return nil
		}
		if len(resp.Values) == 0 {
			utils.Log.Info("No data(relationships) found in the sheet")
		}

		// If no error, fetch the data from the sheet
		err = createJsonFile(resp, "../docs/_data/RelationshipsData.json")
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		utils.Log.Info("Relationships data generated in docs/_data/RelationshipsData.json")
		return nil
	},
}

func init() {
	generateCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	generateCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")
}

func createJsonFile(resp *sheets.ValueRange, jsonFilePath string) error {

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

	jsonData, err := json.MarshalIndent(customResp, "", "    ")
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}
	defer jsonFile.Close()
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}
	return nil
}
