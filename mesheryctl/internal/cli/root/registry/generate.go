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
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	mutils "github.com/meshery/meshkit/utils"

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

	// Individual CSV file paths (new flags)
	modelCSVFlag        string
	componentCSVFlag    string
	relationshipCSVFlag string

	// Per-model timeout duration (default 5 minutes)
	modelTimeout time.Duration
	// Whether to generate only the latest version of each model
	latestVersionOnly bool
)
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long: `Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under "/server/models/<model-name>".
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/registry/generate`,
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED"

// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]

// Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred --model "[model-name]"

// Generate Meshery Models and Component from csv files in a local directory.
mesheryctl registry generate --directory [DIRECTORY_PATH]

// Generate Meshery Models from individual CSV files.
mesheryctl registry generate --model-csv [path/to/models.csv] --component-csv [path/to/components.csv] --relationship-csv [path/to/relationships.csv]

// Generate models with a custom per-model timeout (e.g., 10 minutes per model).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED" --timeout 10m

// Generate only the latest version of each model.
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred "$CRED" --latest-only
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prerequisite check is needed - https://github.com/meshery/meshery/issues/10369
		// TODO: Include a prerequisite check to confirm that this command IS being the executED from within a fork of the Meshery repo, and is being executed at the root of that fork.
		const errorMsg = "[ Spreadsheet ID | Registrant Connection Definition Path | Local Directory | Individual CSV files ] isn't specified\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nmesheryctl registry generate --model-csv [path] --component-csv [path] --relationship-csv [path]\nRun 'mesheryctl registry generate --help' to see detailed help message"

		spreadsheetIdFlag, _ := cmd.Flags().GetString("spreadsheet-id")
		registrantDefFlag, _ := cmd.Flags().GetString("registrant-def")
		directory, _ := cmd.Flags().GetString("directory")
		modelCSV, _ := cmd.Flags().GetString("model-csv")
		componentCSV, _ := cmd.Flags().GetString("component-csv")

		// Check if individual CSV flags are provided
		hasIndividualCSVs := modelCSV != "" && componentCSV != ""

		if spreadsheetIdFlag == "" && registrantDefFlag == "" && directory == "" && !hasIndividualCSVs {
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

		// Validate individual CSV files if provided
		if hasIndividualCSVs {
			if _, err := os.Stat(modelCSV); os.IsNotExist(err) {
				return errors.New(utils.RegistryError(fmt.Sprintf("Model CSV file not found: %s", modelCSV), "generate"))
			}
			if _, err := os.Stat(componentCSV); os.IsNotExist(err) {
				return errors.New(utils.RegistryError(fmt.Sprintf("Component CSV file not found: %s", componentCSV), "generate"))
			}
			relationshipCSV, _ := cmd.Flags().GetString("relationship-csv")
			if relationshipCSV != "" {
				if _, err := os.Stat(relationshipCSV); os.IsNotExist(err) {
					return errors.New(utils.RegistryError(fmt.Sprintf("Relationship CSV file not found: %s", relationshipCSV), "generate"))
				}
			}
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

		// Print start message with timestamp
		startTime := time.Now()
		fmt.Printf("\n⏱️ Starting model generation at %s\n", startTime.Format("15:04:05"))
		fmt.Printf("   Output directory: %s\n", registryLocation)
		fmt.Printf("   Logs directory: %s\n", logDirPath)
		fmt.Printf("   Per-model timeout: %v\n", modelTimeout)
		if latestVersionOnly {
			fmt.Printf("   Mode: Latest version only\n")
		}
		if modelName != "" {
			fmt.Printf("   Target model: %s\n", modelName)
		}
		fmt.Println()

		// Configure generation options
		genOpts := meshkitRegistryUtils.GenerationOptions{
			ModelTimeout:      modelTimeout,
			LatestVersionOnly: latestVersionOnly,
			ProgressCallback: func(modelName string, current, total int) {
				remaining := total - current
				fmt.Printf("📦 [%d/%d] Processing: %s (remaining: %d)\n", current, total, modelName, remaining)
			},
		}

		// Handle individual CSV files (new feature)
		if modelCSVFlag != "" && componentCSVFlag != "" {
			modelCSVFilePath = modelCSVFlag
			componentCSVFilePath = componentCSVFlag
			relationshipCSVFilePath = relationshipCSVFlag
			fmt.Println("📁 Using individual CSV files for generation:")
			fmt.Printf("   Model CSV: %s\n", modelCSVFilePath)
			fmt.Printf("   Component CSV: %s\n", componentCSVFilePath)
			if relationshipCSVFilePath != "" {
				fmt.Printf("   Relationship CSV: %s\n", relationshipCSVFilePath)
			}
			fmt.Println()
		} else if csvDirectory == "" {
			// Using Google Spreadsheet
			fmt.Println("📊 Connecting to Google Spreadsheet...")
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				return errors.New(utils.RegistryError("Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token.", "generate"))
			}

			resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
			if err != nil || resp.HTTPStatusCode != 200 {
				utils.LogError.Error(ErrUpdateRegistry(err, outputLocation))
				return err
			}
			fmt.Println("✅ Connected to spreadsheet successfully")

			sheetGID = GetSheetIDFromTitle(resp, "Models")
			componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")
			relationshipSpredsheetGID = GetSheetIDFromTitle(resp, "Relationships")
			fmt.Println()
		} else {
			// Using directory with CSV files
			fmt.Printf("📁 Reading CSV files from directory: %s\n", csvDirectory)
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = meshkitRegistryUtils.GetCsv(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			if modelCSVFilePath == "" || componentCSVFilePath == "" {
				return fmt.Errorf("ModelCSV and ComponentCSV files must be present in the directory")
			}
			fmt.Println("✅ CSV files found successfully")
			fmt.Println()
		}

		fmt.Println("🔄 Generating models and components...")
		fmt.Println("   (Each model has a timeout of", modelTimeout, ")")
		fmt.Println()

		err = meshkitRegistryUtils.InvokeGenerationFromSheetWithOptions(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv, genOpts)

		elapsed := time.Since(startTime).Round(time.Second)
		fmt.Println()
		if err != nil {
			fmt.Printf("❌ Model generation completed with errors after %v\n", elapsed)
			utils.LogError.Error(err)
			fmt.Printf("   Check error logs at: %s\n", filepath.Join(logDirPath, "registry-errors.log"))
		} else {
			fmt.Printf("✅ Model generation completed successfully in %v\n", elapsed)
		}

		if logFile != nil {
			_ = logFile.Close()
		}
		if errorLogFile != nil {
			_ = errorLogFile.Close()
		}

		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)

		// Print final summary
		fmt.Printf("\n📋 Generation Summary:\n")
		fmt.Printf("   Output: %s\n", registryLocation)
		fmt.Printf("   Logs: %s\n", logDirPath)
		fmt.Printf("   Total time: %v\n", elapsed)
		fmt.Println()

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

	// New flags for individual CSV files
	generateCmd.PersistentFlags().StringVar(&modelCSVFlag, "model-csv", "", "path to the model CSV file")
	generateCmd.PersistentFlags().StringVar(&componentCSVFlag, "component-csv", "", "path to the component CSV file")
	generateCmd.PersistentFlags().StringVar(&relationshipCSVFlag, "relationship-csv", "", "path to the relationship CSV file (optional)")

	// Mark individual CSV flags as required together (model and component are required, relationship is optional)
	generateCmd.MarkFlagsRequiredTogether("model-csv", "component-csv")

	// Mark mutual exclusivity between different input methods
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "directory", "model-csv")

	// New flags for per-model timeout and latest version only
	generateCmd.PersistentFlags().DurationVar(&modelTimeout, "timeout", meshkitRegistryUtils.DefaultModelTimeout, "timeout duration for generating each model (e.g., 5m, 10m, 1h)")
	generateCmd.PersistentFlags().BoolVar(&latestVersionOnly, "latest-only", false, "generate only the latest version of each model")
}
