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

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	mutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"
    "github.com/meshery/schemas/models/v1alpha1/capability"


	"github.com/spf13/cobra"
	"google.golang.org/api/sheets/v4"
	
	"encoding/json"
	"strings"
	"github.com/meshery/schemas"
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
	Long:  "Prerequisite: Execute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under `/server/models/<model-name>`.",
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred $CRED
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

		if csvDirectory == "" {
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				return errors.New(utils.RegistryError("Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token.", "generate"))
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
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = meshkitRegistryUtils.GetCsv(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			if modelCSVFilePath == "" || componentCSVFilePath == "" || relationshipCSVFilePath == "" {
				return fmt.Errorf("ModelCSV, ComponentCSV and RelationshipCSV files must be present in the directory")
			}
		}

		err = meshkitRegistryUtils.InvokeGenerationFromSheet(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv)
		if err != nil {
			// meshkit
			utils.LogError.Error(err)
			return nil
		}

		utils.Log.Info("Applying minimal UI capabilities to generated components...")
		err = enhanceGeneratedComponents(registryLocation)
		if err != nil {
			utils.LogError.Error(fmt.Errorf("Error applying minimal UI capabilities: %v", err))
		} else {
			utils.Log.Info("Successfully applied minimal UI capabilities to generated components")
		}

		_ = logFile.Close()
		_ = errorLogFile.Close()

		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)
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

func getMinimalUICapabilitiesFromSchema() ([]capability.Capability, error) {
	schema, err := schemas.Schemas.ReadFile("schemas/constructs/v1beta1/component/component.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read component schema: %v", err)
	}

	allDefaults, err := extractCapabilitiesDefaultFromSchema(schema)
	if err != nil {
		return nil, fmt.Errorf("failed to extract capabilities from schema: %v", err)
	}

	if len(allDefaults) >= 3 {
		lastThree := allDefaults[len(allDefaults)-3:]
		return convertToCapabilityStructs(lastThree)
	}

	return nil, fmt.Errorf("insufficient default capabilities in schema, found %d", len(allDefaults))
}

func extractCapabilitiesDefaultFromSchema(schema []byte) ([]interface{}, error) {
    var schemaMap map[string]interface{}
    if err := json.Unmarshal(schema, &schemaMap); err != nil {
        return nil, err
    }

    if properties, ok := schemaMap["properties"].(map[string]interface{}); ok {
        if capabilitiesSchema, ok := properties["capabilities"].(map[string]interface{}); ok {
            if defaultValue, ok := capabilitiesSchema["default"].([]interface{}); ok {
                return defaultValue, nil
            }
        }
    }

    return nil, fmt.Errorf("default capabilities not found in schema")
}

func enhanceGeneratedComponents(registryLocation string) error {
	return filepath.Walk(registryLocation, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".json") || !strings.Contains(path, "components") {
			return nil
		}

		var comp component.ComponentDefinition
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		err = json.Unmarshal(data, &comp)
		if err != nil {
			return err
		}

		if needsMinimalCapabilities(&comp) {
			err = applyCapabilitiesToComponent(&comp)
			if err != nil {
				return err
			}
			
			updatedData, err := json.MarshalIndent(comp, "", "  ")
			if err != nil {
				return err
			}

			err = os.WriteFile(path, updatedData, info.Mode())
			if err != nil {
				return err
			}

			utils.Log.Info("Applied minimal UI capabilities to: ", filepath.Base(path))
		}

		return nil
	})
}


func needsMinimalCapabilities(comp *component.ComponentDefinition) bool {
    return comp.Capabilities == nil || len(*comp.Capabilities) == 0
}


func applyCapabilitiesToComponent(comp *component.ComponentDefinition) error {
	if !needsMinimalCapabilities(comp) {
		return nil
	}

	capabilities, err := getMinimalUICapabilitiesFromSchema()
	if err != nil {
		return err
	}
	
	comp.Capabilities = &capabilities
	return nil
}

func convertToCapabilityStructs(rawCaps []interface{}) ([]capability.Capability, error) {
	var capabilities []capability.Capability
	
	for _, rawCap := range rawCaps {
		capMap, ok := rawCap.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid capability format")
		}
		
		cap := capability.Capability{
			SchemaVersion: getString(capMap, "schemaVersion"),
			Version:       getString(capMap, "version"),
			DisplayName:   getString(capMap, "displayName"),
			Description:   getString(capMap, "description"),
			Kind:          getString(capMap, "kind"),
			Type:          getString(capMap, "type"),
			SubType:       getString(capMap, "subType"),
			Key:           getString(capMap, "key"),
			Status:        capability.CapabilityStatus(getString(capMap, "status")),
			EntityState:   convertEntityState(capMap["entityState"]),
			Metadata:      convertMetadata(capMap["metadata"]),
		}
		
		capabilities = append(capabilities, cap)
	}
	
	return capabilities, nil
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func convertEntityState(raw interface{}) []capability.CapabilityEntityState {
	if arr, ok := raw.([]interface{}); ok {
		var states []capability.CapabilityEntityState
		for _, item := range arr {
			if str, ok := item.(string); ok {
				states = append(states, capability.CapabilityEntityState(str))
			}
		}
		return states
	}
	return nil
}

func convertMetadata(raw interface{}) *map[string]interface{} {
	if raw == nil {
		return nil
	}
	if metadata, ok := raw.(map[string]interface{}); ok {
		return &metadata
	}
	return nil
}
