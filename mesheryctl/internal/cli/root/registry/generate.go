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
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	mutils "github.com/meshery/meshkit/utils"
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

	componentTracker = &ComponentTracker{
		modelComponents: make(map[string]map[string]bool),
		mutex:           &sync.Mutex{},
	}
)

type ComponentTracker struct {
	modelComponents map[string]map[string]bool 
	mutex           *sync.Mutex
}

func (ct *ComponentTracker) RegisterComponents(modelName string, components []string) {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if _, exists := ct.modelComponents[modelName]; !exists {
		ct.modelComponents[modelName] = make(map[string]bool)
	}

	for _, component := range components {
		ct.modelComponents[modelName][component] = true
	}
}

func (ct *ComponentTracker) GetComponents(modelName string) []string {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	if componentSet, exists := ct.modelComponents[modelName]; exists {
		components := make([]string, 0, len(componentSet))
		for component := range componentSet {
			components = append(components, component)
		}
		return components
	}

	return []string{}
}

func (ct *ComponentTracker) GetAllComponents() map[string][]string {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	result := make(map[string][]string)
	for model, componentSet := range ct.modelComponents {
		components := make([]string, 0, len(componentSet))
		for component := range componentSet {
			components = append(components, component)
		}
		result[model] = components
	}

	return result
}

func setupDetailedLogging() {
	logrus.AddHook(&detailedLogHook{})
}

type detailedLogHook struct{}

func (h *detailedLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (h *detailedLogHook) Fire(entry *logrus.Entry) error {
	msg := entry.Message

	componentCountRegex := regexp.MustCompile(` extracted (\d+) components for (.+?) \((.+?)\)`)
	if matches := componentCountRegex.FindStringSubmatch(msg); len(matches) >= 4 {
		count := matches[1]
		modelDisplayName := matches[2]
		modelName := matches[3]

		if modelName != "" {
			extractComponentsFromDirectory(modelName)
		}

		components := componentTracker.GetComponents(modelName)

		if len(components) > 0 {
			componentsList := strings.Join(components, ", ")
			entry.Message = fmt.Sprintf("Extracted %s components for model %s (%s): %s",
				count, modelDisplayName, modelName, componentsList)
		} else {
			entry.Message = fmt.Sprintf("Extracted %s components for model %s (%s)",
				count, modelDisplayName, modelName)
		}
	}

	if matches := regexp.MustCompile(`Generated (\d+) components for model \[(.+?)\] (.+)`).FindStringSubmatch(msg); len(matches) >= 4 {
		count, _ := strconv.Atoi(matches[1])
		if count > 0 {
			modelName := matches[2]
			version := matches[3]

			if modelName != "" {
				extractComponentsFromDirectory(modelName)
			}

			components := componentTracker.GetComponents(modelName)

			if len(components) > 0 {
				componentsList := strings.Join(components, ", ")
				entry.Message = fmt.Sprintf("Generated %d components for model [%s] version %s: %s",
					count, modelName, version, componentsList)
			} else {
				entry.Message = fmt.Sprintf("Generated %d components for model [%s] version %s",
					count, modelName, version)
			}
		}
	}

	appendingRegex := regexp.MustCompile(`Appending (\d+) in the components sheet`)
	if matches := appendingRegex.FindStringSubmatch(msg); len(matches) >= 2 {
		countStr := matches[1]
		count, err := strconv.Atoi(countStr)
		if err == nil && count > 0 {
			allComponents := componentTracker.GetAllComponents()

			if len(allComponents) > 0 {
				var detailedMsg strings.Builder
				detailedMsg.WriteString(fmt.Sprintf("Appending %d in the components sheet\n", count))
				detailedMsg.WriteString("Generated components for models:\n")

				for model, components := range allComponents {
					if len(components) > 0 {
						detailedMsg.WriteString(fmt.Sprintf("  - Model [%s]: %s\n",
							model, strings.Join(components, ", ")))
					}
				}

				entry.Message = detailedMsg.String()
			} else {
				entry.Message = fmt.Sprintf("Appending %d in the components sheet\nNo component details available.", count)
			}
		}
	}

	return nil
}

func extractComponentsFromDirectory(modelName string) {
	modelDir := filepath.Join(registryLocation, modelName)
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		return
	}

	findComponentsInDir(modelDir, modelName)

	versionDirs, err := os.ReadDir(modelDir)
	if err == nil {
		for _, versionDir := range versionDirs {
			if versionDir.IsDir() && strings.HasPrefix(versionDir.Name(), "v") {
				versionPath := filepath.Join(modelDir, versionDir.Name())

				componentsPath := filepath.Join(versionPath, "components")
				if _, err := os.Stat(componentsPath); err == nil {
					findComponentsInDir(componentsPath, modelName)
				} else {
					findComponentsInDir(versionPath, modelName)
				}
			}
		}
	}
}

func findComponentsInDir(dirPath string, modelName string) {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".json") && !strings.Contains(info.Name(), "model") {
			componentName := strings.TrimSuffix(info.Name(), ".json")
			if componentName != "" {
				componentTracker.RegisterComponents(modelName, []string{componentName})
			}
		}
		return nil
	})

	if err != nil {
		utils.Log.Debug("Error extracting components from directory:", err)
	}
}

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

		setupDetailedLogging()

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

		utils.Log.Debug("Starting model and component generation process...")

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
