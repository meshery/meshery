// # Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	mutils "github.com/layer5io/meshkit/utils"

	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
)

var (
	componentSpredsheetGID         int64
	relationshipSpredsheetGID      int64
	outputLocation                 string
	pathToRegistrantConnDefinition string
	pathToRegistrantCredDefinition string
	GoogleSpreadSheetURL           = "https://docs.google.com/spreadsheets/d/"
	srv                            *sheets.Service

	// current working directory location
	cwd string

	registryLocation    string
	totalAggregateModel int
	defVersion          = "v1.0.0"
)
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long:  "Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under `/server/models/<model-name>`.",
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
./mesheryctl registry generate --spreadsheet-id "1t0OqpPI_TRNwNB4Rm8yBRIG_LwaDtAiVmozYlg9FE-k" --spreadsheet-cred $CRED
// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]
// Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred --model "[model-name]"
// Generate Meshery Models and Component from csv files in a local directory.
mesheryctl registry generate --directory <DIRECTORY_PATH>
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prerequisite check is needed - https://github.com/meshery/meshery/issues/10369
		// TODO: Include a prerequisite check to confirm that this command IS being the executED from within a fork of the Meshery repo, and is being executed at the root of that fork.

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		var wg sync.WaitGroup
		cwd, _ = os.Getwd()
		registryLocation = filepath.Join(cwd, outputLocation)

		if pathToRegistrantConnDefinition != "" {
			utils.Log.Info("Model generation from Registrant definitions not yet supported.")
			return nil
		}
		var err error

		if csvDirectory == "" {
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				utils.LogError.Error(ErrUpdateRegistry(err, modelLocation))
				return nil
			}

			resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
			if err != nil || resp.HTTPStatusCode != 200 {
				utils.LogError.Error(ErrUpdateRegistry(err, outputLocation))
				return nil
			}

			sheetGID = GetSheetIDFromTitle(resp, "Models")
			componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")
			// Collect list of corresponding relationship by name from spreadsheet
			relationshipSpredsheetGID = GetSheetIDFromTitle(resp, "Relationships")
		} else {
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = utils.GetCsv(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			if modelCSVFilePath == "" || componentCSVFilePath == "" || relationshipCSVFilePath == "" {
				return fmt.Errorf("both ModelCSV and ComponentCSV files must be present in the directory")
			}
		}

		err = utils.InvokeGenerationFromSheet(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID)
		if err != nil {
			// meshkit
			utils.LogError.Error(err)
			return nil
		}
		return err
	},
}

func init() {
	generateCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	generateCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")

	generateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")

	generateCmd.PersistentFlags().StringVar(&pathToRegistrantConnDefinition, "registrant-def", "", "path pointing to the registrant connection definition")
	generateCmd.PersistentFlags().StringVar(&pathToRegistrantCredDefinition, "registrant-cred", "", "path pointing to the registrant credential definition")

	generateCmd.MarkFlagsRequiredTogether("registrant-def", "registrant-cred")

	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "registrant-def")
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-cred", "registrant-cred")
	generateCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "", "specific model name to be generated")
	generateCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")

	generateCmd.PersistentFlags().StringVarP(&csvDirectory, "directory", "d", "", "Directory containing the Model and Component CSV files")

}
