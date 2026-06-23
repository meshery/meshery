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

	"github.com/go-playground/validator/v10"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
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
)

type cmdRegistryGenerateFlags struct {
	SpreadsheetID     string        `json:"spreadsheet-id" validate:"omitempty"`
	SpreadsheetCred   string        `json:"spreadsheet-cred" validate:"omitempty"`
	RegistrantDef     string        `json:"registrant-def" validate:"omitempty"`
	RegistrantCred    string        `json:"registrant-cred" validate:"omitempty"`
	Directory         string        `json:"directory" validate:"omitempty"`
	ModelCSV          string        `json:"model-csv" validate:"omitempty,filepath,registry-generate-model-csv-requires-component"`
	ComponentCSV      string        `json:"component-csv" validate:"omitempty,filepath,registry-generate-component-csv-requires-model"`
	RelationshipCSV   string        `json:"relationship-csv" validate:"omitempty,filepath,registry-generate-relationship-csv-requires-model-and-component"`
	Timeout           time.Duration `json:"timeout" validate:"omitempty"`
	LatestVersionOnly bool          `json:"latest-only" validate:"boolean"`
}

var registryGenerateFlags cmdRegistryGenerateFlags

func registerRegistryGenerateCustomValidations() error {
	flagValidator := mesheryctlflags.GetFlagValidator()

	err := flagValidator.Validator.RegisterValidation("registry-generate-model-csv-requires-component", func(fl validator.FieldLevel) bool {
		flags, ok := fl.Parent().Interface().(cmdRegistryGenerateFlags)
		if !ok || flags.ModelCSV == "" {
			return true
		}

		return flags.ComponentCSV != ""
	})
	if err != nil {
		return err
	}

	err = flagValidator.Validator.RegisterValidation("registry-generate-component-csv-requires-model", func(fl validator.FieldLevel) bool {
		flags, ok := fl.Parent().Interface().(cmdRegistryGenerateFlags)
		if !ok || flags.ComponentCSV == "" {
			return true
		}

		return flags.ModelCSV != ""
	})
	if err != nil {
		return err
	}

	err = flagValidator.Validator.RegisterValidation("registry-generate-relationship-csv-requires-model-and-component", func(fl validator.FieldLevel) bool {
		flags, ok := fl.Parent().Interface().(cmdRegistryGenerateFlags)
		if !ok || flags.RelationshipCSV == "" {
			return true
		}

		return flags.ModelCSV != "" && flags.ComponentCSV != ""
	})
	if err != nil {
		return err
	}

	flagValidator.CustomErrors["registry-generate-model-csv-requires-component"] = "--component-csv is required when --model-csv is provided"
	flagValidator.CustomErrors["registry-generate-component-csv-requires-model"] = "--model-csv is required when --component-csv is provided"
	flagValidator.CustomErrors["registry-generate-relationship-csv-requires-model-and-component"] = "--relationship-csv can only be used with --model-csv and --component-csv"

	return nil
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long: `Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under "/server/models/<model-name>".
Find more information at: https://docs.meshery.io/reference/mesheryctl/registry/generate`,
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

		if err := registerRegistryGenerateCustomValidations(); err != nil {
			return err
		}

		if err := mesheryctlflags.ValidateCmdFlags(cmd, &registryGenerateFlags); err != nil {
			return err
		}

		spreadsheeetID = registryGenerateFlags.SpreadsheetID
		spreadsheeetCred = registryGenerateFlags.SpreadsheetCred
		pathToRegistrantConnDefinition = registryGenerateFlags.RegistrantDef
		pathToRegistrantCredDefinition = registryGenerateFlags.RegistrantCred
		csvDirectory = registryGenerateFlags.Directory

		// Check if individual CSV flags are provided
		hasIndividualCSVs := registryGenerateFlags.ModelCSV != "" && registryGenerateFlags.ComponentCSV != ""

		if registryGenerateFlags.SpreadsheetID == "" && registryGenerateFlags.RegistrantDef == "" && registryGenerateFlags.Directory == "" && !hasIndividualCSVs {
			return utils.ErrInvalidArgument(errors.New(utils.RegistryError(errorMsg, "generate")))
		}

		if registryGenerateFlags.SpreadsheetID != "" && registryGenerateFlags.SpreadsheetCred == "" {
			return utils.ErrInvalidArgument(errors.New(utils.RegistryError("Spreadsheet Credentials is required\n\nUsage: \nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED\nmesheryctl registry generate --spreadsheet-id [Spreadsheet ID] --spreadsheet-cred $CRED --model \"[model-name]\"\nRun 'mesheryctl registry generate --help'", "generate")))
		}

		if registryGenerateFlags.RegistrantDef != "" && registryGenerateFlags.RegistrantCred == "" {
			return utils.ErrInvalidArgument(errors.New(utils.RegistryError("Registrant Credentials is required\n\nUsage: mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]\nRun 'mesheryctl registry generate --help'", "generate")))
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
		fmt.Printf("   Per-model timeout: %v\n", registryGenerateFlags.Timeout)
		if registryGenerateFlags.LatestVersionOnly {
			fmt.Printf("   Mode: Latest version only\n")
		}
		if modelName != "" {
			fmt.Printf("   Target model: %s\n", modelName)
		}
		fmt.Println()

		// Configure generation options
		genOpts := meshkitRegistryUtils.GenerationOptions{
			ModelTimeout:      registryGenerateFlags.Timeout,
			LatestVersionOnly: registryGenerateFlags.LatestVersionOnly,
			ProgressCallback: func(modelName string, current, total int) {
				remaining := total - current
				fmt.Printf("📦 [%d/%d] Processing: %s (remaining: %d)\n", current, total, modelName, remaining)
			},
		}

		// Handle individual CSV files (new feature)
		if registryGenerateFlags.ModelCSV != "" && registryGenerateFlags.ComponentCSV != "" {
			modelCSVFilePath = registryGenerateFlags.ModelCSV
			componentCSVFilePath = registryGenerateFlags.ComponentCSV
			relationshipCSVFilePath = registryGenerateFlags.RelationshipCSV
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
				utils.Log.Error(ErrUpdateRegistry(err, outputLocation))
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
		fmt.Println("   (Each model has a timeout of", registryGenerateFlags.Timeout, ")")
		fmt.Println()

		err = meshkitRegistryUtils.InvokeGenerationFromSheetWithOptions(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv, genOpts)

		elapsed := time.Since(startTime).Round(time.Second)
		fmt.Println()
		if err != nil {
			fmt.Printf("❌ Model generation completed with errors after %v\n", elapsed)
			utils.Log.Error(err)
			fmt.Printf("   Check error logs at: %s\n", filepath.Join(logDirPath, "registry-errors.log"))
		} else {
			fmt.Printf("✅ Model generation completed successfully in %v\n", elapsed)
		}

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
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.SpreadsheetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.SpreadsheetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")

	generateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")

	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.RegistrantDef, "registrant-def", "", "path pointing to the registrant connection definition")
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.RegistrantCred, "registrant-cred", "", "path pointing to the registrant credential definition")

	generateCmd.MarkFlagsRequiredTogether("registrant-def", "registrant-cred")

	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "registrant-def")
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-cred", "registrant-cred")
	generateCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "", "specific model name to be generated")
	generateCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")

	generateCmd.PersistentFlags().StringVarP(&registryGenerateFlags.Directory, "directory", "d", "", "Directory containing the Model and Component CSV files")

	// New flags for individual CSV files
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.ModelCSV, "model-csv", "", "path to the model CSV file")
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.ComponentCSV, "component-csv", "", "path to the component CSV file")
	generateCmd.PersistentFlags().StringVar(&registryGenerateFlags.RelationshipCSV, "relationship-csv", "", "path to the relationship CSV file (optional)")

	// Mark individual CSV flags as required together (model and component are required, relationship is optional)
	generateCmd.MarkFlagsRequiredTogether("model-csv", "component-csv")

	// Mark mutual exclusivity between different input methods
	generateCmd.MarkFlagsMutuallyExclusive("spreadsheet-id", "directory", "model-csv")

	// New flags for per-model timeout and latest version only
	generateCmd.PersistentFlags().DurationVar(&registryGenerateFlags.Timeout, "timeout", meshkitRegistryUtils.DefaultModelTimeout, "timeout duration for generating each model (e.g., 5m, 10m, 1h), default: 5m")
	generateCmd.PersistentFlags().BoolVar(&registryGenerateFlags.LatestVersionOnly, "latest-only", false, "generate only the latest version of each model")
}
