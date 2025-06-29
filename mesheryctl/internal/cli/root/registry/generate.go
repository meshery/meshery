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
	"io"
	"strings"
	"encoding/csv"
	"encoding/json"
	"net/http"

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

	// CVS file paths for models, components, and relationships
	componentCSVFilePath    string
	modelCSVFilePath        string
	relationshipCSVFilePath string
	spreadsheeetID          string
	spreadsheeetCred        string
	csvDirectory            string
	sheetGID                int64
	modelName               string
	logFile                 *os.File
	errorLogFile            *os.File

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

			err = applySchemaDefaults(srv, spreadsheeetID, componentSpredsheetGID)
			if err != nil {
				utils.LogError.Warn("Failed to apply schema defaults:", err)
			}
		} else {
			modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, err = meshkitRegistryUtils.GetCsv(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			if modelCSVFilePath == "" || componentCSVFilePath == "" || relationshipCSVFilePath == "" {
				return fmt.Errorf("ModelCSV, ComponentCSV and RelationshipCSV files must be present in the directory")
			}

			err = applySchemaDefaultsToCSV(componentCSVFilePath)
			if err != nil {
				utils.LogError.Warn("Failed to apply schema defaults to CSV:", err)
			}
		}

		err = meshkitRegistryUtils.InvokeGenerationFromSheet(&wg, registryLocation, sheetGID, componentSpredsheetGID, spreadsheeetID, modelName, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath, relationshipSpredsheetGID, srv)
		if err != nil {
			// meshkit
			utils.LogError.Error(err)
			return nil
		}
		_ = logFile.Close()
		_ = errorLogFile.Close()

		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)
		return err
	},
}

// applySchemaDefaults applies schema defaults to the components in the Google Sheet
func applySchemaDefaults(srv *sheets.Service, spreadsheetID string, componentSheetGID int64) error {
    schemaDefaults, err := getCapabilitiesFromSchema()
    if err != nil {
        return fmt.Errorf("failed to read schema defaults: %v", err)
    }
    
    if schemaDefaults.Capabilities == "" {
        utils.Log.Info("No default capabilities defined in schema")
        return nil
    }
    
    readRange := "Components!A:R"
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
    if err != nil {
        return fmt.Errorf("unable to read components: %v", err)
    }
    
    tempFile, err := createTempFileWithSchemaDefaults(resp.Values, schemaDefaults)
    if err != nil {
        return fmt.Errorf("failed to create temp file with defaults: %v", err)
    }
    
    componentCSVFilePath = tempFile
    
    utils.Log.Infof("Applied schema-driven default capabilities: '%s'", schemaDefaults.Capabilities)
    return nil
}

func applySchemaDefaultsToCSV(componentCSVFilePath string) error {
    schemaDefaults, err := getCapabilitiesFromSchema()
    if err != nil {
        return fmt.Errorf("failed to read schema defaults: %v", err)
    }
    
    if schemaDefaults.Capabilities == "" {
        utils.Log.Info("No default capabilities defined in schema")
        return nil
    }
    
    file, err := os.Open(componentCSVFilePath)
    if err != nil {
        return fmt.Errorf("failed to open CSV file: %v", err)
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return fmt.Errorf("failed to read CSV: %v", err)
    }
    
    if len(records) == 0 {
        return nil
    }
    
    capabilitiesCol := -1
    for i, header := range records[0] {
        if strings.ToLower(strings.TrimSpace(header)) == "capabilities" {
            capabilitiesCol = i
            break
        }
    }
    
    if capabilitiesCol == -1 {
        utils.Log.Warn("Capabilities column not found in CSV, skipping schema defaults")
    	return nil
    }
    
    updatedCount := 0
    for i := 1; i < len(records); i++ {
        if capabilitiesCol < len(records[i]) {
            if strings.TrimSpace(records[i][capabilitiesCol]) == "" {
                records[i][capabilitiesCol] = schemaDefaults.Capabilities
                updatedCount++
            }
        }
    }
    
    outputFile, err := os.Create(componentCSVFilePath)
    if err != nil {
        return fmt.Errorf("failed to create updated CSV: %v", err)
    }
    defer outputFile.Close()
    
    writer := csv.NewWriter(outputFile)
    defer writer.Flush()
    
    for _, record := range records {
        if err := writer.Write(record); err != nil {
            return fmt.Errorf("failed to write CSV record: %v", err)
        }
    }
    
    utils.Log.Infof("Applied schema defaults to %d components in CSV", updatedCount)
    return nil
}

type ComponentSchemaDefaults struct {
    Capabilities string `json:"capabilities"`
	// Can add more fields as needed based on the schema
}

// getCapabilitiesFromSchema acquires default capabilities from the schema
func getCapabilitiesFromSchema() (*ComponentSchemaDefaults, error) {
    schemaPath := filepath.Join(cwd, "schemas", "constructs", "v1beta1", "component", "component.json")
    if _, err := os.Stat(schemaPath); err == nil {
        utils.Log.Info("Reading schema from local file:", schemaPath)
        return readSchemaFromFile(schemaPath)
    }
    
    utils.Log.Info("Local schema not found, fetching from GitHub")
    return readSchemaFromGitHub()
}

func readSchemaFromFile(schemaPath string) (*ComponentSchemaDefaults, error) {
    file, err := os.ReadFile(schemaPath)
    if err != nil {
        return nil, fmt.Errorf("failed to read schema file: %v", err)
    }
    
    var schema map[string]interface{}
    if err := json.Unmarshal(file, &schema); err != nil {
        return nil, fmt.Errorf("failed to parse schema JSON: %v", err)
    }
    
    defaults := &ComponentSchemaDefaults{
        Capabilities: "",
    }
    
    if properties, ok := schema["properties"].(map[string]interface{}); ok {
        if capsProp, ok := properties["capabilities"].(map[string]interface{}); ok {
            if defaultVal, ok := capsProp["default"].(string); ok {
                defaults.Capabilities = defaultVal
            }
        }
    }
    
    return defaults, nil
}

func readSchemaFromGitHub() (*ComponentSchemaDefaults, error) {
    url := "https://raw.githubusercontent.com/meshery/schemas/master/schemas/constructs/v1beta1/component/component.json"
    
    resp, err := http.Get(url)
    if err != nil {
        utils.Log.Warn("Failed to fetch schema from GitHub:", err)
        return &ComponentSchemaDefaults{Capabilities: ""}, nil
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        utils.Log.Warn(fmt.Sprintf("GitHub returned status %d, using empty default", resp.StatusCode))
        return &ComponentSchemaDefaults{Capabilities: ""}, nil
    }
    
	// Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        utils.Log.Warn("Failed to read response body:", err)
        return &ComponentSchemaDefaults{Capabilities: ""}, nil
    }
    
    // Analyze the JSON schema
    var schema map[string]interface{}
    if err := json.Unmarshal(body, &schema); err != nil {
        utils.Log.Warn("Failed to parse schema JSON:", err)
        return &ComponentSchemaDefaults{Capabilities: ""}, nil
    }
    
    defaults := &ComponentSchemaDefaults{
        Capabilities: "",
    }
    
    // Exract the default capabilities from the schema
    if properties, ok := schema["properties"].(map[string]interface{}); ok {
        if capsProp, ok := properties["capabilities"].(map[string]interface{}); ok {
            if defaultVal, ok := capsProp["default"].(string); ok {
                defaults.Capabilities = defaultVal
                utils.Log.Info("Found default capabilities in GitHub schema:", defaultVal)
            } else if defaultArray, ok := capsProp["default"].([]interface{}); ok {
                caps := make([]string, len(defaultArray))
                for i, cap := range defaultArray {
                    caps[i] = fmt.Sprintf("%v", cap)
                }
                defaults.Capabilities = strings.Join(caps, ",")
                utils.Log.Info("Found default capabilities array in GitHub schema:", defaults.Capabilities)
            }
        }
    }
    
    if defaults.Capabilities == "" {
        utils.Log.Info("No default capabilities found in GitHub schema")
    }
    
    utils.Log.Info("Successfully read schema from GitHub")
    return defaults, nil
}

func createTempFileWithSchemaDefaults(values [][]interface{}, defaults *ComponentSchemaDefaults) (string, error) {
    tempFile, err := os.CreateTemp("", "components_with_defaults_*.csv")
    if err != nil {
        return "", err
    }
    defer tempFile.Close()
    
    writer := csv.NewWriter(tempFile)
    defer writer.Flush()
    
    if len(values) == 0 {
        return tempFile.Name(), nil
    }
    
    header := make([]string, len(values[0]))
    capabilitiesCol := -1
    for i, col := range values[0] {
        header[i] = fmt.Sprintf("%v", col)
        if strings.ToLower(header[i]) == "capabilities" {
            capabilitiesCol = i
        }
    }
    writer.Write(header)
    
    updatedCount := 0
    for _, row := range values[1:] {
        record := make([]string, len(row))
        for i, cell := range row {
            record[i] = fmt.Sprintf("%v", cell)
        }
        
        if capabilitiesCol >= 0 && capabilitiesCol < len(record) {
            if strings.TrimSpace(record[capabilitiesCol]) == "" {
                record[capabilitiesCol] = defaults.Capabilities
                updatedCount++
            }
        }
        
        writer.Write(record)
    }
    
    utils.Log.Infof("Applied schema defaults to %d components", updatedCount)
    return tempFile.Name(), nil
}

func ErrUpdateRegistry(err error, location string) error {
    return fmt.Errorf("failed to update registry at location %s: %w", location, err)
}

func GetSheetIDFromTitle(spreadsheet *sheets.Spreadsheet, title string) int64 {
    for _, sheet := range spreadsheet.Sheets {
        if sheet.Properties != nil && sheet.Properties.Title == title {
            return sheet.Properties.SheetId
        }
    }
    utils.Log.Warning(fmt.Sprintf("Sheet with title '%s' not found", title))
    return 0
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
