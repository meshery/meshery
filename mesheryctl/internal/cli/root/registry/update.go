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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/store"
	comp "github.com/meshery/schemas/models/v1beta1/component"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	csvDir                   string
	modelLocation            string
	logFile                  *os.File
	errorLogFile             *os.File
	sheetGID                 int64
	totalAggregateComponents int
	logDirPath               = filepath.Join(mutils.GetHome(), ".meshery", "logs", "registry")
	componentCSVHelper       *utils.ComponentCSVHelper
)

// This command is used for retreving the information of components based on the sheet. It updates the components with the actual values of the fetched for sheet.
// Look the utils.ComponentCSV to see the values fetched.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the registry with latest data.",
	Long:  "Updates the component metadata (SVGs, shapes, styles and other) by referring from a Google Spreadsheet or CSV directory.",
	Example: `
// Update models from Meshery Integration Spreadsheet
mesheryctl registry update --spreadsheet-id [id] --spreadsheet-cred [base64 encoded spreadsheet credential] -i [path to the directory containing models].

// Updating models in the meshery/meshery repo
mesheryctl registry update --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --spreadsheet-cred $CRED
// Updating models in the meshery/meshery repo based on flag
mesheryctl registry update --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --spreadsheet-cred $CRED --model "[model-name]"

// Updating models from a local CSV directory
mesheryctl registry update --csv-dir /path/to/csv-directory

// Example to update with both Google Sheets and CSV, CSV takes precedence
mesheryctl registry update --spreadsheet-id [id] --spreadsheet-cred [base64 encoded spreadsheet credential] --csv-dir /path/to/csv-directory
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		err := os.MkdirAll(logDirPath, 0755)
		if err != nil {
			return ErrUpdateRegistry(err, modelLocation)
		}
		utils.Log.SetLevel(logrus.DebugLevel)
		logFilePath := filepath.Join(logDirPath, "registry-update")
		logFile, err = os.Create(logFilePath)
		if err != nil {
			return ErrUpdateRegistry(err, modelLocation)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if csvDir != "" {
			utils.Log.Info("Using local CSV directory: ", csvDir)
			err := InvokeCompUpdateFromCSV()
			if err != nil {
				utils.Log.Error(err)
				return err
			}
			return nil
		}

		if spreadsheeetID != "" && spreadsheeetCred != "" {
			utils.Log.Info("Using Google Sheet with ID: ", spreadsheeetID)
			err := InvokeCompUpdateFromSheets()
			if err != nil {
				utils.Log.Error(err)
				return err
			}
			return nil
		}
		return errors.New("please provide a Google Sheet or a local csv directory")
	},
}
var (
	ExcludeDirs = []string{"relationships", "policies"}
)

type compUpdateTracker struct {
	totalComps        int
	totalCompsUpdated int
	version           string
}

func InvokeCompUpdateFromSheets() error {
	utils.Log.UpdateLogOutput(logFile)

	defer func() {
		_ = logFile.Close()
		utils.Log.UpdateLogOutput(os.Stdout)

		// Additionally log the summary to the terminal
		utils.Log.Info(fmt.Sprintf("Updated %d models and %d components", totalAggregateModel, totalAggregateComponents))
		utils.Log.Info("refer ", logDirPath, " for detailed registry update logs")

		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()
	modelToCompUpdateTracker := store.NewGenericThreadSafeStore[[]compUpdateTracker]()

	srv, err := mutils.NewSheetSRV(spreadsheeetCred)
	if err != nil {
		return err
	}
	resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		return fmt.Errorf("failed to get google sheet: %w", err)
	}

	sheetGID = GetSheetIDFromTitle(resp, "Components")

	url := GoogleSpreadSheetURL + spreadsheeetID
	componentCSVHelper, err = utils.NewComponentCSVHelper(url, "Components", sheetGID, componentCSVFilePath)
	if err != nil {
		return err
	}

	err = componentCSVHelper.ParseComponentsSheet(modelName)
	if err != nil {
		err = ErrUpdateRegistry(err, modelLocation)
		utils.Log.Error(err)
		return nil
	}

	utils.Log.Info("Total Registrants: ", len(componentCSVHelper.Components))

	err = updateRegistryComponents(componentCSVHelper.Components, modelToCompUpdateTracker)
	if err != nil {
		return err
	}

	logModelUpdateSummary(modelToCompUpdateTracker)
	return nil
}

func InvokeCompUpdateFromCSV() error {
	utils.Log.UpdateLogOutput(logFile)
	defer func() {
		_ = logFile.Close()
		utils.Log.UpdateLogOutput(os.Stdout)
		utils.Log.Info(fmt.Sprintf("Updated %d models and %d components", totalAggregateModel, totalAggregateComponents))
		utils.Log.Info("refer ", logDirPath, " for detailed registry update logs")
		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()

	modelToCompUpdateTracker := store.NewGenericThreadSafeStore[[]compUpdateTracker]()

	// Parse local CSV files
	localComponents, err := parseLocalCSVDirectory(csvDir)
	if err != nil {
		return err
	}

	utils.Log.Info("Total Registrants (from CSV): ", len(localComponents))

	// Reuse the exact same iteration logic
	err = updateRegistryComponents(localComponents, modelToCompUpdateTracker)
	if err != nil {
		return err
	}

	logModelUpdateSummary(modelToCompUpdateTracker)
	return nil
}

func updateRegistryComponents(
	components map[string]map[string][]utils.ComponentCSV,
	modelToCompUpdateTracker *store.GenerticThreadSafeStore[[]compUpdateTracker],
) error {
	pwd, _ := os.Getwd()

	// var wg sync.WaitGroup
	for registrant, model := range components {
		if registrant == "" {
			continue
		}

		// Iterate all models
		for modelName, components := range model {
			availableComponentsPerModelPerVersion := 0
			modelPath := filepath.Join(pwd, modelLocation, modelName)
			utils.Log.Info("Starting to update components of model ", modelName)

			modelContents, err := os.ReadDir(modelPath)
			if err != nil {
				err = ErrUpdateModel(err, modelName)
				utils.Log.Error(err)
				continue
			}

			// Iterate over all content inside model
			// Comps, relationships, policies
			compUpdateArray := []compUpdateTracker{}
			for _, content := range modelContents {
				totalCompsUpdatedPerModelPerVersion := 0

				if content.IsDir() {
					if utils.Contains(content.Name(), ExcludeDirs) != -1 {
						continue
					}

					// A model can have components with multiple versions
					versionPath := filepath.Join(modelPath, content.Name(), "v1.0.0") 
					entries, _ := os.ReadDir(versionPath)
					availableComponentsPerModelPerVersion += len(entries)

					utils.Log.Info("Updating component of model ", modelName, " with version: ", content.Name())

					for _, component := range components {
						compPath := filepath.Join(versionPath, "components", fmt.Sprintf("%s.json", component.Component))
						componentByte, err := os.ReadFile(compPath)
						if err != nil {
							utils.Log.Error(ErrUpdateComponent(err, modelName, component.Component))
							continue
						}
						componentDef := comp.ComponentDefinition{}
						err = json.Unmarshal(componentByte, &componentDef)
						if err != nil {
							utils.Log.Error(ErrUpdateComponent(err, modelName, component.Component))
							continue
						}

						err = component.UpdateCompDefinition(&componentDef)
						if err != nil {
							utils.Log.Error(ErrUpdateComponent(err, modelName, component.Component))
							continue
						}
						tmpFilePath := filepath.Join(versionPath, "components", "tmp_model.json")

						// Ensure the temporary file is removed regardless of what happens
						defer func() {
							_ = os.Remove(tmpFilePath)
						}()

						if _, err := os.Stat(compPath); err == nil {
							existingData, err := os.ReadFile(compPath)
							if err != nil {
								utils.Log.Error(err)
								continue
							}

							err = mutils.WriteJSONToFile[comp.ComponentDefinition](tmpFilePath, componentDef)
							if err != nil {
								utils.Log.Error(err)
								continue
							}

							newData, err := os.ReadFile(tmpFilePath)
							if err != nil {
								utils.Log.Error(err)
								continue
							}

							if bytes.Equal(existingData, newData) {
								utils.Log.Info("No changes detected for ", componentDef.Component.Kind)
								continue
							} else {
								err = mutils.WriteJSONToFile[comp.ComponentDefinition](compPath, componentDef)
								if err != nil {
									utils.Log.Error(err)
									continue
								}
								totalCompsUpdatedPerModelPerVersion++

							}
						}
					}

					compUpdateArray = append(compUpdateArray, compUpdateTracker{
						totalComps:        availableComponentsPerModelPerVersion,
						totalCompsUpdated: totalCompsUpdatedPerModelPerVersion,
						version:           content.Name(),
					})
				}
			}
			modelToCompUpdateTracker.Set(modelName, compUpdateArray)
			utils.Log.Info("\n")
		}

	}
	return nil
}

func logModelUpdateSummary(modelToCompUpdateTracker *store.GenerticThreadSafeStore[[]compUpdateTracker]) {
	values := modelToCompUpdateTracker.GetAllPairs()
	for key, val := range values {
		for _, value := range val {
			utils.Log.Info(fmt.Sprintf("For model %s-%s, updated %d out of %d components.", key, value.version, value.totalCompsUpdated, value.totalComps))
			totalAggregateComponents += value.totalCompsUpdated
		}
	}

	utils.Log.Info(fmt.Sprintf("For %d models updated %d components", len(values), totalAggregateComponents))
}

func parseLocalCSVDirectory(dirPath string) (map[string]map[string][]utils.ComponentCSV, error) {
	localComps := make(map[string]map[string][]utils.ComponentCSV)
	var allErrors []error

	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("could not access csv directory '%s': %w", dirPath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("provided path is not a directory: %s", dirPath)
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// optionally recurse or skip
			continue
		}
		if filepath.Ext(entry.Name()) == ".csv" {
			csvPath := filepath.Join(dirPath, entry.Name())
			fileComps, err := parseOneCSVFile(csvPath)
			if err != nil {
				allErrors = append(allErrors, err)
				continue
			}
			// Merge fileComps into localComps
			for registrant, modelMap := range fileComps {
				if _, ok := localComps[registrant]; !ok {
					localComps[registrant] = make(map[string][]utils.ComponentCSV)
				}
				for modelName, comps := range modelMap {
					localComps[registrant][modelName] = append(localComps[registrant][modelName], comps...)
				}
			}
		}
	}
	if len(allErrors) > 0 {
		return nil, errors.Join(allErrors...)
	}
	return localComps, nil
}

// parseOneCSVFile reads a single CSV file, enforces the required columns, and returns
// a nested map structure: map[registrant]map[modelName][]utils.ComponentInfo
func parseOneCSVFile(csvPath string) (map[string]map[string][]utils.ComponentCSV, error) {
	result := make(map[string]map[string][]utils.ComponentCSV)

	f, err := os.Open(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %w", csvPath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV file '%s': %w", csvPath, err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV file '%s' must have at least one header row and one data row", csvPath)
	}

	header := records[0]
	colIndex := map[string]int{}
	for i, col := range header {
		colIndex[col] = i
	}

	requiredCols := []string{"registrant", "model", "component"}
	for _, rc := range requiredCols {
		if _, ok := colIndex[rc]; !ok {
			return nil, fmt.Errorf("missing required column '%s' in CSV file '%s'", rc, csvPath)
		}
	}

	for i, row := range records[1:] {
		registrant := row[colIndex["registrant"]]
		modelVal := row[colIndex["model"]]
		compVal := row[colIndex["component"]]
		if registrant == "" || modelVal == "" || compVal == "" {
			return nil, fmt.Errorf("missing data in CSV file '%s', row number %d", csvPath, i+2)
		}

		if _, ok := result[registrant]; !ok {
			result[registrant] = make(map[string][]utils.ComponentCSV)
		}
		result[registrant][modelVal] = append(
			result[registrant][modelVal],
			utils.ComponentCSV{
				Component: compVal,
			},
		)
	}

	return result, nil
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&modelLocation, "input", "i", "../server/meshmodel", "relative or absolute input path to the models directory")
	_ = updateCmd.MarkPersistentFlagRequired("path")

	updateCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet-id", "", "spreadsheet it for the integration spreadsheet")
	updateCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")
	updateCmd.PersistentFlags().StringVarP(&modelName, "model", "m", "", "specific model name to be generated")
	updateCmd.PersistentFlags().StringVar(&csvDir, "csv-dir", "", "Path to directory containing local CSV files for model and component updates")

	updateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")

}
