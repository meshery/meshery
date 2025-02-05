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
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/sirupsen/logrus"

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
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred $CRED
// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]
// Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred $CRED --model "[model-name]"
// Generate Meshery Models and Component from csv files in a local directory.
mesheryctl registry generate --directory <DIRECTORY_PATH>
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prerequisite check is needed - https://github.com/meshery/meshery/issues/10369
		// TODO: Include a prerequisite check to confirm that this command IS being the executED from within a fork of the Meshery repo, and is being executed at the root of that fork.
		const errorMsg = "[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory ] isn't specified\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nRun 'mesheryctl registry generate --help' to see detailed help message"

		spreadsheetIdFlag, _ := cmd.Flags().GetString("spreadsheet-id")
		registrantDefFlag, _ := cmd.Flags().GetString("registrant-def")
		directory, _ := cmd.Flags().GetString("directory")

		if spreadsheetIdFlag == "" && registrantDefFlag == "" && directory == "" {
			return errors.New(utils.RegistryError(errorMsg, "generate"))
		}

		spreadsheetCredFlag, _ := cmd.Flags().GetString("spreadsheet-cred")
		registrantCredFlag, _ := cmd.Flags().GetString("registrant-cred")

		if spreadsheetIdFlag != "" && spreadsheetCredFlag == "" {
			return errors.New(utils.RegistryError("Spreadsheet Credentials is required\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nRun 'mesheryctl registry generate --help'", "generate"))
		}

		if registrantDefFlag != "" && registrantCredFlag == "" {
			return errors.New(utils.RegistryError("Registrant Credentials is required\n\nUsage: mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]\nRun 'mesheryctl registry generate --help'", "generate"))
		}

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

		utils.Log.Info("Starting model generation...")
		utils.Log.Infof("Model output location: %s", registryLocation)

		if err := setupLogging(); err != nil {
			utils.LogError.Error(err)
		}

		defer cleanup()

		if csvDirectory == "" {
			utils.Log.Info("Initializing Google Sheets service...")
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				utils.LogError.Error(errors.New(utils.RegistryError("Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token.", "generate")))
				return nil
			}

			utils.Log.Info("Fetching data from spreadsheet ID: ", spreadsheeetID)
			resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
			if err != nil || resp.HTTPStatusCode != 200 {
				utils.LogError.Error(errors.New(utils.RegistryError(fmt.Sprintf("%s\nPlease verify the spreadsheet ID and permissions.", err), "generate")))
				return nil
			}

			utils.Log.Debug("Extracting sheet IDs from spreadsheet")
			sheetGID = GetSheetIDFromTitle(resp, "Models")
			componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")
			relationshipSpredsheetGID = GetSheetIDFromTitle(resp, "Relationships")

			utils.Log.Debugf("Sheet IDs: \n -Models: %d, \n -Components: %d, \n -Relationships: %d",
				sheetGID, componentSpredsheetGID, relationshipSpredsheetGID)
		} else {
			utils.Log.Info("Reading CSV files from directory: ", csvDirectory)
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = utils.GetCsv(csvDirectory)
			if err != nil {
				utils.LogError.Error(errors.New(utils.RegistryError(fmt.Sprintf("Error reading directory %s: %v", csvDirectory, err), "generate")))
				return nil
			}

			if modelCSVFilePath == "" || componentCSVFilePath == "" || relationshipCSVFilePath == "" {
				errMsg := "Required CSV files missing. Directory must contain:\n" +
					"- Model definition CSV file\n" +
					"- Component definition CSV file\n" +
					"- Relationship definition CSV file"
				utils.LogError.Error(errors.New(utils.RegistryError(errMsg, "generate")))
				return nil
			}

			utils.Log.Info("Successfully located all required CSV files")
			utils.Log.Debugf("Found CSV files: \n -Models: %s \n -Components: %s \n -Relationships: %s",
				filepath.Base(modelCSVFilePath),
				filepath.Base(componentCSVFilePath),
				filepath.Base(relationshipCSVFilePath))
		}

		err = utils.InvokeGenerationFromSheet(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv)
		if err != nil {
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

func setupLogging() error {
	if err := os.MkdirAll(logDirPath, 0755); err != nil {
		return fmt.Errorf("Failed to create log directory: %v", err)
	}

	utils.Log.SetLevel(logrus.DebugLevel)
	logFile, err := os.Create(filepath.Join(logDirPath, "model-generation.log"))
	if err != nil {
		return fmt.Errorf("Failed to create log file: %v", err)
	}

	utils.LogError.SetLevel(logrus.ErrorLevel)
	errorLogFile, err := os.Create(filepath.Join(logDirPath, "registry-errors.log"))
	if err != nil {
		return fmt.Errorf("Failed to create error log file: %v", err)
	}

	utils.Log.UpdateLogOutput(io.MultiWriter(os.Stdout, logFile))
	utils.LogError.UpdateLogOutput(io.MultiWriter(os.Stdout, errorLogFile))

	return nil
}

func cleanup() {
	logFile.Close()
	errorLogFile.Close()

	utils.Log.UpdateLogOutput(os.Stdout)
	utils.LogError.UpdateLogOutput(os.Stdout)
}
