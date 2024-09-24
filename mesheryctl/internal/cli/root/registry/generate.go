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
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/models/meshmodel/entity"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/generators"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/store"
	"github.com/meshery/schemas/models/v1beta1/component"
	v1beta1Model "github.com/meshery/schemas/models/v1beta1/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/sheets/v4"
)

var (
	componentSpredsheetGID         int64
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
var (
	artifactHubCount        = 0
	artifactHubRateLimit    = 100
	artifactHubRateLimitDur = 5 * time.Minute
	artifactHubMutex        sync.Mutex
)
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long:  "Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under `/server/models/<model-name>`.",
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred
// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]
// Generate a specific Model from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet).
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred --model "[model-name]"
// Generate Meshery Models and Component from csv files in a local directory.
mesheryctl registry generate -directory <DIRECTORY_PATH>
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Prerequisite check is needed - https://github.com/meshery/meshery/issues/10369
		// TODO: Include a prerequisite check to confirm that this command IS being the executED from within a fork of the Meshery repo, and is being executed at the root of that fork.
		//

		err := os.MkdirAll(logDirPath, 0755)
		if err != nil {
			return ErrUpdateRegistry(err, modelLocation)
		}
		utils.Log.SetLevel(logrus.DebugLevel)
		logFilePath := filepath.Join(logDirPath, "model-generation.log")
		logFile, err = os.Create(logFilePath)
		if err != nil {
			return err
		}

		utils.LogError.SetLevel(logrus.ErrorLevel)
		logErrorFilePath := filepath.Join(logDirPath, "registry-errors.log")
		errorLogFile, err = os.Create(logErrorFilePath)
		if err != nil {
			return err
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

		// isCsvPathPresent := modelCSVFilePath != "" && componentCSVFilePath != ""
		if csvDirectory == "" {
			srv, err = mutils.NewSheetSRV(spreadsheeetCred)
			if err != nil {
				utils.LogError.Error(ErrUpdateRegistry(err, modelLocation))
				return err
			}

			resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
			if err != nil || resp.HTTPStatusCode != 200 {
				utils.LogError.Error(ErrUpdateRegistry(err, outputLocation))
				return err
			}

			// Collect list of Models by name from spreadsheet
			sheetGID = GetSheetIDFromTitle(resp, "Models")
			// Collect list of corresponding Components by name from spreadsheet
			componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")
		} else {
			// Get all files in the directory
			files, err := os.ReadDir(csvDirectory)
			if err != nil {
				return fmt.Errorf("error reading the directory: %v", err)
			}
			for _, file := range files {
				filePath := filepath.Join(csvDirectory, file.Name())
				if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
					headers, secondRow, err := getCSVHeader(filePath)
					if utils.Contains("modelDisplayName", headers) != -1 || utils.Contains("modelDisplayName", secondRow) != -1 {
						modelCSVFilePath = filePath
					} else if utils.Contains("component", headers) != -1 || utils.Contains("component", secondRow) != -1 { // Check if the file matches the ComponentCSV structure
						componentCSVFilePath = filePath
					}
					if err != nil {
						return fmt.Errorf("error checking file %s: %v", file.Name(), err)
					}

				}
			}

			if modelCSVFilePath == "" || componentCSVFilePath == "" {
				return fmt.Errorf("both ModelCSV and ComponentCSV files must be present in the directory")
			}
		}

		err = InvokeGenerationFromSheet(&wg)
		if err != nil {
			// meshkit
			utils.LogError.Error(err)
			return err
		}

		return err
	},
}

type compGenerateTracker struct {
	totalComps int
	version    string
}

var modelToCompGenerateTracker = store.NewGenericThreadSafeStore[compGenerateTracker]()

func InvokeGenerationFromSheet(wg *sync.WaitGroup) error {
	weightedSem := semaphore.NewWeighted(20)
	url := GoogleSpreadSheetURL + spreadsheeetID
	totalAvailableModels := 0
	spreadsheeetChan := make(chan utils.SpreadsheetData)

	defer func() {
		logModelGenerationSummary(modelToCompGenerateTracker)

		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)
		utils.Log.Info(fmt.Sprintf("Summary: %d models, %d components generated.", totalAggregateModel, totalAggregateComponents))

		utils.Log.Info("See ", logDirPath, " for detailed logs.")

		_ = logFile.Close()
		_ = errorLogFile.Close()
		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()

	modelCSVHelper, err := parseModelSheet(url)
	if err != nil {
		return err
	}

	componentCSVHelper, err := parseComponentSheet(url)
	if err != nil {
		return err
	}
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	multiErrorWriter := io.MultiWriter(os.Stdout, errorLogFile)

	utils.Log.UpdateLogOutput(multiWriter)
	utils.LogError.UpdateLogOutput(multiErrorWriter)
	var wgForSpreadsheetUpdate sync.WaitGroup
	wgForSpreadsheetUpdate.Add(1)

	go func() {
		utils.ProcessModelToComponentsMap(componentCSVHelper.Components)
		utils.VerifyandUpdateSpreadsheet(spreadsheeetCred, &wgForSpreadsheetUpdate, srv, spreadsheeetChan, spreadsheeetID, modelCSVFilePath, componentCSVFilePath)
	}()
	// Iterate models from the spreadsheet
	for _, model := range modelCSVHelper.Models {
		if modelName != "" && modelName != model.Model {
			continue
		}
		totalAvailableModels++
		ctx := context.Background()

		err := weightedSem.Acquire(ctx, 1)
		if err != nil {
			break
		}
		wg.Add(1)
		go func(model utils.ModelCSV) {
			defer func() {
				wg.Done()
				weightedSem.Release(1)
			}()
			if mutils.ReplaceSpacesAndConvertToLowercase(model.Registrant) == "meshery" {
				err = GenerateDefsForCoreRegistrant(model, componentCSVHelper)
				if err != nil {
					utils.LogError.Error(err)
				}
				return
			}

			generator, err := generators.NewGenerator(model.Registrant, model.SourceURL, model.Model)
			if err != nil {
				utils.LogError.Error(ErrGenerateModel(err, model.Model))
				return
			}

			if mutils.ReplaceSpacesAndConvertToLowercase(model.Registrant) == "artifacthub" {
				rateLimitArtifactHub()

			}
			pkg, err := generator.GetPackage()
			if err != nil {
				utils.LogError.Error(ErrGenerateModel(err, model.Model))
				return
			}

			version := pkg.GetVersion()
			modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model)
			if err != nil {
				utils.LogError.Error(ErrGenerateModel(err, model.Model))
				return
			}
			modelDef, alreadyExsit, err := writeModelDefToFileSystem(&model, version, modelDirPath)
			if err != nil {
				utils.LogError.Error(err)
				return
			}
			if alreadyExsit {
				totalAvailableModels--
			}
			comps, err := pkg.GenerateComponents()
			if err != nil {
				utils.LogError.Error(ErrGenerateModel(err, model.Model))
				return
			}
			lengthOfComps := len(comps)

			for _, comp := range comps {
				comp.Version = defVersion
				// Assign the component status corresponding to model status.
				// i.e. If model is enabled comps are also "enabled". Ultimately all individual comps itself will have ability to control their status.
				// The status "enabled" indicates that the component will be registered inside the registry.
				if modelDef.Metadata == nil {
					modelDef.Metadata = &v1beta1Model.ModelDefinition_Metadata{}
				}
				if modelDef.Metadata.AdditionalProperties == nil {
					modelDef.Metadata.AdditionalProperties = make(map[string]interface{})
				}

				if comp.Model.Metadata.AdditionalProperties != nil {
					modelDef.Metadata.AdditionalProperties["source_uri"] = comp.Model.Metadata.AdditionalProperties["source_uri"]
				}
				comp.Model = *modelDef

				utils.AssignDefaultsForCompDefs(&comp, modelDef)
				compAlreadyExist, err := comp.WriteComponentDefinition(compDirPath, "json")
				if compAlreadyExist {
					lengthOfComps--
				}
				if err != nil {
					utils.Log.Info(err)
				}
			}
			if !alreadyExsit {
				if len(comps) == 0 {
					utils.LogError.Error(ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model))
				} else {
					utils.Log.Info("Current model: ", model.Model)
					utils.Log.Info(" extracted ", lengthOfComps, " components for ", model.ModelDisplayName, " (", model.Model, ")")
				}
			} else {
				if len(comps) > 0 {
					utils.Log.Info("Model already exists: ", model.Model)
				} else {
					utils.LogError.Error(ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model))
				}
			}
			spreadsheeetChan <- utils.SpreadsheetData{
				Model:      &model,
				Components: comps,
			}

			modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
				totalComps: lengthOfComps,
				version:    version,
			})
		}(model)

	}
	wg.Wait()
	close(spreadsheeetChan)
	wgForSpreadsheetUpdate.Wait()
	return nil
}

// For registrants eg: meshery, whose components needs to be directly created by referencing meshery/schemas repo.
// the sourceURL contains the path of models component definitions
func GenerateDefsForCoreRegistrant(model utils.ModelCSV, ComponentCSVHelper *utils.ComponentCSVHelper) error {
	var version string
	parts := strings.Split(model.SourceURL, "/")
	// Assuming the URL is always of the format "protocol://github.com/owner/repo/tree/definitions/{model-name}/version/components"
	// We know the version is the 7th element (0-indexed) in the split URL
	if len(parts) >= 8 {
		version = parts[8] // Fetch the version from the expected position
	} else {
		return fmt.Errorf("invalid SourceURL format: %s", model.SourceURL)
	}
	isModelPublished, _ := strconv.ParseBool(model.PublishToRegistry)
	var compDefComps []component.ComponentDefinition
	alreadyExist := false
	actualCompCount := 0
	defer func() {
		modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
			totalComps: len(compDefComps) - actualCompCount,
			version:    version,
		})
	}()
	status := entity.Ignored
	if isModelPublished {
		status = entity.Enabled
	}
	_status := component.ComponentDefinitionStatus(status)
	modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model)
	if err != nil {
		err = ErrGenerateModel(err, model.Model)
		return err
	}
	modelDef, alreadyExists, err := writeModelDefToFileSystem(&model, version, modelDirPath)
	if err != nil {
		return ErrGenerateModel(err, model.Model)
	}
	isModelPublishToSite, _ := strconv.ParseBool(model.PublishToSites)
	alreadyExist = alreadyExists
	for registrant, models := range ComponentCSVHelper.Components {
		if registrant != "meshery" {
			continue
		}
		for _, comps := range models {
			for _, comp := range comps {
				if comp.Model != model.Model {
					continue
				}
				var componentDef component.ComponentDefinition
				componentDef, err = comp.CreateComponentDefinition(isModelPublishToSite, "v1.0.0")
				if err != nil {
					utils.Log.Error(ErrUpdateComponent(err, modelName, comp.Component))
					continue
				}
				componentDef.Status = &_status
				componentDef.Model = *modelDef
				alreadyExists, err = componentDef.WriteComponentDefinition(compDirPath, "json")
				if err != nil {
					err = ErrGenerateComponent(err, comp.Model, componentDef.DisplayName)
					utils.LogError.Error(err)
				}
				if alreadyExists {
					actualCompCount++
				}
				compDefComps = append(compDefComps, componentDef)
			}
		}
	}

	if !alreadyExist {
		if len(compDefComps) == 0 {
			utils.LogError.Error(ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model))
		} else if len(compDefComps)-actualCompCount == 0 {
			utils.Log.Info("Current model: ", model.Model)
			utils.Log.Info(" no change in components for ", model.ModelDisplayName, " (", model.Model, ")")
		} else {
			utils.Log.Info("Current model: ", model.Model)
			utils.Log.Info(" extracted ", len(compDefComps)-actualCompCount, " components for ", model.ModelDisplayName, " (", model.Model, ")")
		}
	} else {
		if len(compDefComps) > 0 {
			if len(compDefComps)-actualCompCount == 0 {
				utils.Log.Info("Model already exists: ", model.Model)
			} else {
				utils.Log.Info("Current model: ", model.Model)
				utils.Log.Info(" extracted ", len(compDefComps)-actualCompCount, " components for ", model.ModelDisplayName, " (", model.Model, ")")
			}
		} else {
			utils.LogError.Error(ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model))
		}
	}

	return nil
}

func parseModelSheet(url string) (*utils.ModelCSVHelper, error) {
	modelCSVHelper, err := utils.NewModelCSVHelper(url, "Models", sheetGID, modelCSVFilePath)
	if err != nil {
		return nil, err
	}

	err = modelCSVHelper.ParseModelsSheet(false, modelName)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return modelCSVHelper, nil
}
func rateLimitArtifactHub() {
	artifactHubMutex.Lock()
	defer artifactHubMutex.Unlock()

	if artifactHubCount > 0 && artifactHubCount%artifactHubRateLimit == 0 {
		utils.Log.Info("Rate limit reached for Artifact Hub. Sleeping for 5 minutes...")
		time.Sleep(artifactHubRateLimitDur)
	}
	artifactHubCount++
}
func parseComponentSheet(url string) (*utils.ComponentCSVHelper, error) {
	compCSVHelper, err := utils.NewComponentCSVHelper(url, "Components", componentSpredsheetGID, componentCSVFilePath)
	if err != nil {
		return nil, err
	}

	err = compCSVHelper.ParseComponentsSheet(modelName)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return compCSVHelper, nil
}

// version corresponds to the version of the pacakge from which model was sourced and not the definition version.
// Eg: helm chart version for kube-prom-stack.
func createVersionedDirectoryForModelAndComp(version, modelName string) (string, string, error) {
	modelDirPath := filepath.Join(registryLocation, modelName, version, defVersion)
	err := mutils.CreateDirectory(modelDirPath)
	if err != nil {
		return "", "", err
	}

	compDirPath := filepath.Join(modelDirPath, "components")
	err = mutils.CreateDirectory(compDirPath)
	return modelDirPath, compDirPath, err
}

func writeModelDefToFileSystem(model *utils.ModelCSV, version, modelDefPath string) (*v1beta1Model.ModelDefinition, bool, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	filePath := filepath.Join(modelDefPath, "model.json")
	tmpFilePath := filepath.Join(modelDefPath, "tmp_model.json")

	// Ensure the temporary file is removed regardless of what happens
	defer func() {
		_ = os.Remove(tmpFilePath)
	}()

	// Check if the file exists

	if _, err := os.Stat(filePath); err == nil {
		existingData, err := os.ReadFile(filePath)
		if err != nil {
			goto NewGen
		}

		err = modelDef.WriteModelDefinition(tmpFilePath, "json")
		if err != nil {
			goto NewGen
		}

		newData, err := os.ReadFile(tmpFilePath)
		if err != nil {
			goto NewGen
		}

		// Compare the existing and new data
		if bytes.Equal(existingData, newData) {
			var oldModelDef v1beta1Model.ModelDefinition
			err = encoding.Unmarshal(existingData, &oldModelDef)
			if err != nil {
				goto NewGen
			}
			// If they are the same, return without changes
			return &oldModelDef, true, nil
		}
	}
NewGen:
	// Write the model definition to the actual file if it's new or different
	err := modelDef.WriteModelDefinition(filePath, "json")
	if err != nil {
		return nil, false, err
	}

	return &modelDef, false, nil
}

func logModelGenerationSummary(modelToCompGenerateTracker *store.GenerticThreadSafeStore[compGenerateTracker]) {
	for key, val := range modelToCompGenerateTracker.GetAllPairs() {
		utils.Log.Info(fmt.Sprintf("Generated %d components for model [%s] %s", val.totalComps, key, val.version))
		totalAggregateComponents += val.totalComps
		if val.totalComps > 0 {
			totalAggregateModel++
		}
	}

	utils.Log.Info(fmt.Sprintf("-----------------------------\n-----------------------------\nGenerated %d models and %d components", totalAggregateModel, totalAggregateComponents))
}

func getCSVHeader(filePath string) (headers, secondRow []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return headers, secondRow, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err = reader.Read() // Read the first line

	if err != nil {
		return headers, secondRow, err
	}

	secondRow, err = reader.Read()
	if err != nil {
		return headers, secondRow, err
	}
	return headers, secondRow, nil
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
