package registry

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/store"
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

	// Tracks the indexed mapping between component spreadsheet columns.
	// Used when generating component definition from spreadsheet itself, for eg: Compnent of Meshery core model.
	// The GoogleSpreadsheetAPI doesn't return column names, hence when invoking generation columns names are retrived by dumoing the sheet in CSV format then extrcting the columns (ComponentCSVHelper)
	componentSpreadsheetCols []string

	// current working directory location
	cwd string

	registryLocation    string
	totalAggregateModel int
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import Models",
	Long:  "Import models from spreadsheet, GitHub or ArtifactHub repositories",
	Example: `
    // Import models from Meshery Integration Spreadsheet
    mesheryctl registry import --spreadsheet_url <url> --spreadsheet_cred <base64 encoded spreadsheet credential>
    
    // Directly import models from one of the supported registrants by using Registrant Connection Definition and (optional) Registrant Credential Definition
    mesheryctl registry import --registrant_def <path to connection definition> --registrant_cred <path to credential definition>
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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

		cwd, _ = os.Getwd()
		registryLocation = filepath.Join(cwd, outputLocation)

		if pathToRegistrantConnDefinition != "" {
			utils.Log.Info("Model Generation from registrant definitions not yet supproted.")
			return nil
		}
		var err error

		srv, err = mutils.NewSheetSRV(spreadsheeetCred)
		if err != nil {
			fmt.Println(err, utils.Log.GetLevel(), ErrUpdateRegistry(err, modelLocation), ErrUpdateRegistry(err, modelLocation).Error())
			utils.Log.Error(ErrUpdateRegistry(err, modelLocation))
			return err
		}

		resp, err := srv.Spreadsheets.Get(spreadsheeetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(ErrUpdateRegistry(err, outputLocation))
			return err
		}

		sheetGID = GetSheetIDFromTitle(resp, "Models")
		componentSpredsheetGID = GetSheetIDFromTitle(resp, "Components")

		err = InvokeGenerationFromSheet()
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

func InvokeGenerationFromSheet() error {
	utils.Log.UpdateLogOutput(logFile)

	totalAvailableModels := 0
	defer func() {
		_ = logFile.Close()
		utils.Log.UpdateLogOutput(os.Stdout)
		utils.Log.Info(fmt.Sprintf("Generated %d models and %d components", totalAggregateModel, totalAggregateComponents))

		utils.Log.Info("refer ", logDirPath, " for detailed registry generate logs")

		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()

	url := GoogleSpreadSheetURL + spreadsheeetID

	modelCSVHelper, err := utils.NewModelCSVHelper(url, "Models", sheetGID)
	if err != nil {
		return err
	}

	err = modelCSVHelper.ParseModelsSheet(false)
	if err != nil {
		return ErrGenerateModel(err, "all")
	}

	componentCSVHelper, err := utils.NewComponentCSVHelper(url, "Components", componentSpredsheetGID)
	if err != nil {
		return err
	}

	componentSpreadsheetCols, _ = componentCSVHelper.GetColumns()

	weightedSem := semaphore.NewWeighted(20)

	var wg sync.WaitGroup

	// Iterate models from the spreadsheet
	for _, model := range modelCSVHelper.Models {
		totalAvailableModels++

		ctx := context.Background()

		err := weightedSem.Acquire(ctx, 1)
		if err != nil {
			break
		}
		utils.Log.Info("Current model: ", model.Model)
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
				time.Sleep(1 * time.Second)
			}
			pkg, err := generator.GetPackage()
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			version := pkg.GetVersion()

			modelDefPath, _, err := writeModelDefToFileSystem(&model, version)
			if err != nil {
				utils.Log.Error(err)
				return
			}

			comps, err := pkg.GenerateComponents()
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			utils.Log.Info("Extracted ", len(comps), "for model: %s", model.ModelDisplayName)

			compDirName, err := createVersionDirectoryForModel(modelDefPath, version)
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}

			for _, comp := range comps {
				err := comp.WriteComponentDefinition(compDirName)
				if err != nil {
					utils.Log.Info(err)
				}
			}
			modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
				totalComps: len(comps),
				version:    version,
			})
		}(model)

	}
	wg.Wait()
	logModelGenerationSummary(modelToCompGenerateTracker)

	return nil
}

// For registrants eg: meshery, whose components needs to be directly created by referencing the sheet.
// the sourceURL contains the "rangeString" : the value that describes the sheet in which the components are listed, and their starting row to ending row.
func GenerateDefsForCoreRegistrant(model utils.ModelCSV) error {
	totalComps := 0
	version := "1.0.0"
	defer func() {
		modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
			totalComps: totalComps,
			version:    version,
		})
	}()

	modelPath, modelDef, err := writeModelDefToFileSystem(&model, version) // how to infer this? @Beginner86 any idea? new column?
	if err != nil {
		return ErrGenerateModel(err, model.Model)
	}
	isModelPublished, _ := modelDef.Metadata["published"].(bool)
	components, err := srv.Spreadsheets.Values.Get(spreadsheeetID, model.SourceURL).Do()
	if err != nil {
		return ErrGenerateModel(err, model.Model)
	}
	compDirName, err := createVersionDirectoryForModel(modelPath, modelDef.Version)
	if err != nil {
		err = ErrGenerateModel(err, modelDef.Name)
		return err
	}

	if len(componentSpreadsheetCols) > 0 {
		for _, comp := range components.Values {
			totalComps++
			component := make(map[string]interface{}, len(comp))
			for i, compValue := range comp {
				component[componentSpreadsheetCols[i]] = compValue
			}
			compName, _ := component["component"].(string)
			compCSV, err := mutils.MarshalAndUnmarshal[map[string]interface{}, utils.ComponentCSV](component)
			if err != nil {
				err = ErrGenerateComponent(err, model.Model, compName)
				utils.Log.Error(err)
				continue
			}

			componentDef, err := compCSV.CreateComponentDefinition(isModelPublished)
			if err != nil {
				err = ErrGenerateComponent(err, model.Model, compName)
				utils.Log.Error(err)
				continue
			}
			componentDef.Model = *modelDef // remove this, left for backward compatibility

			err = componentDef.WriteComponentDefinition(compDirName)

			if err != nil {
				err = ErrGenerateComponent(err, model.Model, compName)
				utils.Log.Error(err)
				continue
			}
		}
	}

	return nil
}

func createVersionDirectoryForModel(modelDefPath, version string) (string, error) {
	versionDirPath := filepath.Join(modelDefPath, version)
	err := mutils.CreateDirectory(versionDirPath)
	return versionDirPath, err
}

func writeModelDefToFileSystem(model *utils.ModelCSV, version string) (string, *v1alpha1.Model, error) {
	modelDef := model.CreateModelDefinition(version)
	modelDefPath := filepath.Join(registryLocation, modelDef.Name)

	err := modelDef.WriteModelDefinition(modelDefPath)
	if err != nil {
		return "", nil, err
	}

	return modelDefPath, &modelDef, nil
}

func logModelGenerationSummary(modelToCompGenerateTracker *store.GenerticThreadSafeStore[compGenerateTracker]) {
	for key, val := range modelToCompGenerateTracker.GetAllPairs() {
		utils.Log.Info(fmt.Sprintf("For model %s-%s, generated %d components.", key, val.version, val.totalComps))
		totalAggregateComponents += val.totalComps
		totalAggregateModel++
	}

	utils.Log.Info(fmt.Sprintf("Generated %d models and %d components", totalAggregateModel, totalAggregateComponents))
}

func init() {
	importCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet_id", "", "spreadsheet it for the integration spreadsheet")
	importCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet_cred", "", "base64 encoded credential to download the spreadsheet")

	importCmd.MarkFlagsRequiredTogether("spreadsheet_id", "spreadsheet_cred")

	importCmd.PersistentFlags().StringVar(&pathToRegistrantConnDefinition, "registrant_def", "", "path pointing to the registrant connection definition")
	importCmd.PersistentFlags().StringVar(&pathToRegistrantCredDefinition, "registrant_cred", "", "path pointing to the registrant credetial definition")

	importCmd.MarkFlagsRequiredTogether("registrant_def", "registrant_cred")

	importCmd.MarkFlagsMutuallyExclusive("spreadsheet_id", "registrant_def")
	importCmd.MarkFlagsMutuallyExclusive("spreadsheet_cred", "registrant_cred")

	importCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")
}
