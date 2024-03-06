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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/store"
	"github.com/sirupsen/logrus"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/spf13/cobra"
)

var (
	modelLocation            string
	logFile                  *os.File
	sheetGID                 int64
	totalAggregateComponents int
	logDirPath               = filepath.Join(mutils.GetHome(), ".meshery", "logs")
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the registry with latest data.",
	Long:  "`Updates the component metadata (SVGs, shapes, styles and other) by referring from a Google Spreadsheet.`",
	Example: `
// Update models from Meshery Integration Spreadsheet
mesheryctl registry update --spreadsheet-id [id] --spreadsheet-cred [base64 encoded spreadsheet credential] -i [path to the directory containing models].

// Updating models in the meshery/meshery repo
mesheryctl registry update --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --spreadsheet-cred $CRED
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

		srv, err := mutils.NewSheetSRV(spreadsheeetCred)
		if err != nil {
			utils.Log.Error(ErrUpdateRegistry(err, modelLocation))
			return err
		}
		resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(ErrUpdateRegistry(err, outputLocation))
			return err
		}

		sheetGID = GetSheetIDFromTitle(resp, "Components")

		err = InvokeCompUpdate()
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		return nil
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

func InvokeCompUpdate() error {
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

	url := GoogleSpreadSheetURL + spreadsheeetID
	componentCSVHelper, err := utils.NewComponentCSVHelper(url, "Components", sheetGID)
	if err != nil {
		return err
	}

	err = componentCSVHelper.ParseComponentsSheet()
	if err != nil {
		err = ErrUpdateRegistry(err, modelLocation)
		utils.Log.Error(err)
		return err
	}

	utils.Log.Info("Total Registrants: ", len(componentCSVHelper.Components))

	// Since component update doesn't take long skip doing it concurrently
	// weightedSem := semaphore.NewWeighted(20)
	pwd, _ := os.Getwd()

	// var wg sync.WaitGroup

	for registrant, model := range componentCSVHelper.Components {
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
					versionPath := filepath.Join(modelPath, content.Name())
					entries, _ := os.ReadDir(versionPath)
					availableComponentsPerModelPerVersion += len(entries)

					utils.Log.Info("Updating component of model ", modelName, " with version: ", content.Name())

					for _, component := range components {
						utils.Log.Info("Updating ", component.Component)
						compPath := fmt.Sprintf("%s/%s.json", versionPath, component.Component)
						componentByte, err := os.ReadFile(compPath)
						if err != nil {
							utils.Log.Error(ErrUpdateComponent(err, modelName, component.Component))
							continue
						}
						componentDef := v1alpha1.ComponentDefinition{}
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
						err = mutils.WriteJSONToFile[v1alpha1.ComponentDefinition](compPath, componentDef)
						if err != nil {
							utils.Log.Error(err)
							continue
						}
						totalCompsUpdatedPerModelPerVersion++
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
	logModelUpdateSummary(modelToCompUpdateTracker)
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

	utils.Log.Info(fmt.Sprintf("Updated %d models and %d components", len(values), totalAggregateComponents))
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&modelLocation, "input", "i", "../server/meshmodel", "relative or absolute input path to the models directory")
	_ = updateCmd.MarkPersistentFlagRequired("path")

	updateCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet-id", "", "spreadsheet it for the integration spreadsheet")
	updateCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")

	updateCmd.MarkFlagsRequiredTogether("spreadsheet-id", "spreadsheet-cred")

}
