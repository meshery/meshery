package relationships

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkit "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

var (
	sheetID                string
	googleSheetCredentials string
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

var GenerateRelationshipDocsCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate relationships docs",
	Long:  "generate relationships docs from the google spreadsheets",
	Example: `
    // generate relationships docs
    mesheryctl relationships generate $CRED
`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp relationship generate $CRED [google-sheets-credential] --sheetId [sheet-id]\nRun 'mesheryctl exp relationship generate --help' to see detailed help message"

		if len(args) == 0 {
			return errors.New(utils.RelationshipsError("Google Sheet Credentials is required\n"+errMsg, "generate"))
		}

		// Check if flag is set
		sheetIdFlag, _ := cmd.Flags().GetString("sheetId")

		if sheetIdFlag == "" {
			return errors.New(utils.RelationshipsError("Sheet ID is required\n"+errMsg, "generate"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		googleSheetCredentials = args[0]
		srv, err := meshkit.NewSheetSRV(googleSheetCredentials)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		resp, err := srv.Spreadsheets.Values.Get(sheetID, "Relationships").Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(err)
			return nil
		}
		if len(resp.Values) == 0 {
			utils.Log.Info("No data(relationships) found in the sheet")
		}

		// If no error, fetch the data from the sheet
		err = createJsonFile(resp)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		utils.Log.Info("Relationships data generated successfully in docs/_data/RelationshipsData.json")
		return nil
	},
}

func init() {
	GenerateRelationshipDocsCmd.PersistentFlags().StringVarP(&sheetID, "sheetId", "s", "", "Google Sheet ID")
}
func createJsonFile(resp *sheets.ValueRange) error {

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

	jsonFile, err := os.Create("../docs/_data/RelationshipsData.json")
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
