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
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/generators/github"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/store"
	"github.com/layer5io/meshkit/utils/walker"
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

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Models",
	Long:  "Prerequisite: Excecute this command from the root of a meshery/meshery repo fork.\n\nGiven a Google Sheet with a list of model names and source locations, generate models and components any Registrant (e.g. GitHub, Artifact Hub) repositories.\n\nGenerated Model files are written to local filesystem under `/server/models/<model-name>`.",
	Example: `
// Generate Meshery Models from a Google Spreadsheet (i.e. "Meshery Integrations" spreadsheet). 
mesheryctl registry generate --spreadsheet-id "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw" --spreadsheet-cred
// Directly generate models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
mesheryctl registry generate --registrant-def [path to connection definition] --registrant-cred [path to credential definition]
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
		logFilePath := filepath.Join(logDirPath, "registry-generate")
		logFile, err = os.Create(logFilePath)
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

		srv, err = mutils.NewSheetSRV(spreadsheeetCred)
		if err != nil {
			utils.Log.Error(ErrUpdateRegistry(err, modelLocation))
			return err
		}

		resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(ErrUpdateRegistry(err, outputLocation))
			return err
		}

		// Collect list of Models by name from spreadsheet
		sheetGID = GetSheetIDFromTitle(resp, "Models")
		// Collect list of corresponding Components by name from spreadsheet
		componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")

		err = InvokeGenerationFromSheet(&wg)
		if err != nil {
			// meshkit
			utils.Log.Error(err)
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
		utils.Log.Info(fmt.Sprintf("Summary: %d models, %d components generated.", totalAggregateModel, totalAggregateComponents))

		utils.Log.Info("See ", logDirPath, " for detailed logs.")

		_ = logFile.Close()
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

	utils.Log.UpdateLogOutput(logFile)

	var wgForSpreadsheetUpdate sync.WaitGroup
	wgForSpreadsheetUpdate.Add(1)
	go func() {
		utils.ProcessModelToComponentsMap(componentCSVHelper.Components)
		utils.VerifyandUpdateSpreadsheet(spreadsheeetCred, &wgForSpreadsheetUpdate, srv, spreadsheeetChan, spreadsheeetID)
	}()

	// Iterate models from the spreadsheet
	for _, model := range modelCSVHelper.Models {
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
			if model.Registrant == "meshery" {
				err = GenerateDefsForCoreRegistrant(model)
				if err != nil {
					utils.Log.Error(err)
				}
				return
			}

			generator, err := generators.NewGenerator(model.Registrant, model.SourceURL, model.Model)
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}

			if model.Registrant == "artifacthub" {
				time.Sleep(10 * time.Second)
			}
			pkg, err := generator.GetPackage()
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}

			version := pkg.GetVersion()
			modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model)
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			modelDef, err := writeModelDefToFileSystem(&model, version, modelDirPath)
			if err != nil {
				utils.Log.Error(err)
				return
			}

			comps, err := pkg.GenerateComponents()
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			utils.Log.Info("Current model: ", model.Model)
			utils.Log.Info(" extracted ", len(comps), " components for ", model.ModelDisplayName, " (", model.Model, ")")
			for _, comp := range comps {
				comp.Version = defVersion
				if comp.Metadata == nil {
					comp.Metadata = make(map[string]interface{})
				}
				// Assign the component status corresponding to model status.
				// i.e. If model is enabled comps are also "enabled". Ultimately all individual comps itself will have ability to control their status.
				// The status "enabled" indicates that the component will be registered inside the registry.
				comp.Model = *modelDef
				assignDefaultsForCompDefs(&comp, modelDef)
				err := comp.WriteComponentDefinition(compDirPath)
				if err != nil {
					utils.Log.Info(err)
				}
			}

			spreadsheeetChan <- utils.SpreadsheetData{
				Model:      &model,
				Components: comps,
			}

			modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
				totalComps: len(comps),
				version:    version,
			})
		}(model)

	}
	wg.Wait()
	close(spreadsheeetChan)
	wgForSpreadsheetUpdate.Wait()
	return nil
}

func assignDefaultsForCompDefs(componentDef *v1beta1.ComponentDefinition, modelDef *v1beta1.Model) {
	componentDef.Metadata["status"] = modelDef.Status
	for k, v := range modelDef.Metadata {
		componentDef.Metadata[k] = v
	}
}

// For registrants eg: meshery, whose components needs to be directly created by referencing meshery/schemas repo.
// the sourceURL contains the path of models component definitions
func GenerateDefsForCoreRegistrant(model utils.ModelCSV) error {
	totalComps := 0
	var version string
	defer func() {
		modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
			totalComps: totalComps,
			version:    version,
		})
	}()

	path, err := url.Parse(model.SourceURL)
	if err != nil {
		err = ErrGenerateModel(err, model.Model)
		utils.Log.Error(err)
		return nil
	}
	gitRepo := github.GitRepo{
		URL:         path,
		PackageName: model.Model,
	}
	owner, repo, branch, root, err := gitRepo.ExtractRepoDetailsFromSourceURL()
	if err != nil {
		err = ErrGenerateModel(err, model.Model)
		utils.Log.Error(err)
		return nil
	}

	isModelPublished, _ := strconv.ParseBool(model.PublishToRegistry)
	//Initialize walker
	gitWalker := walker.NewGit()
	if isModelPublished {
		gw := gitWalker.
			Owner(owner).
			Repo(repo).
			Branch(branch).
			Root(root).
			RegisterFileInterceptor(func(f walker.File) error {
				// Check if the file has a JSON extension
				if filepath.Ext(f.Name) != ".json" {
					return nil
				}
				contentBytes := []byte(f.Content)
				var componentDef v1beta1.ComponentDefinition
				if err := json.Unmarshal(contentBytes, &componentDef); err != nil {
					return err
				}
				version = componentDef.Model.Model.Version
				modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model)
				if err != nil {
					err = ErrGenerateModel(err, model.Model)
					return err
				}
				_, err = writeModelDefToFileSystem(&model, version, modelDirPath) // how to infer this? @Beginner86 any idea? new column?
				if err != nil {
					return ErrGenerateModel(err, model.Model)
				}

				err = componentDef.WriteComponentDefinition(compDirPath)
				if err != nil {
					err = ErrGenerateComponent(err, model.Model, componentDef.DisplayName)
					utils.Log.Error(err)
				}
				return nil
			})
		err = gw.Walk()
		if err != nil {
			return err
		}
	}

	return nil
}

func parseModelSheet(url string) (*utils.ModelCSVHelper, error) {
	modelCSVHelper, err := utils.NewModelCSVHelper(url, "Models", sheetGID)
	if err != nil {
		return nil, err
	}

	err = modelCSVHelper.ParseModelsSheet(false)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return modelCSVHelper, nil
}

func parseComponentSheet(url string) (*utils.ComponentCSVHelper, error) {
	compCSVHelper, err := utils.NewComponentCSVHelper(url, "Components", componentSpredsheetGID)
	if err != nil {
		return nil, err
	}
	err = compCSVHelper.ParseComponentsSheet()
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

func writeModelDefToFileSystem(model *utils.ModelCSV, version, modelDefPath string) (*v1beta1.Model, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	err := modelDef.WriteModelDefinition(modelDefPath)
	if err != nil {
		return nil, err
	}

	return &modelDef, nil
}

func logModelGenerationSummary(modelToCompGenerateTracker *store.GenerticThreadSafeStore[compGenerateTracker]) {
	for key, val := range modelToCompGenerateTracker.GetAllPairs() {
		utils.Log.Info(fmt.Sprintf("Generated %d components for model [%s] %s", val.totalComps, key, val.version))
		totalAggregateComponents += val.totalComps
		totalAggregateModel++
	}

	utils.Log.Info(fmt.Sprintf("-----------------------------\n-----------------------------\nGenerated %d models and %d components", totalAggregateModel, totalAggregateComponents))
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

	generateCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")
}
