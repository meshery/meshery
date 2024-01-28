package registry

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var (
	outputLocation string

	pathToRegistrantConnDefinition string
	pathToRegistrantCredDefinition string
	GoogleSpreadSheetURL           = "https://docs.google.com/spreadsheets/d/"
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

		if pathToRegistrantConnDefinition != "" {
			utils.Log.Info("Model Generation from registrant definitions not yet supproted.")
			return nil
		}

		srv, err := mutils.NewSheetSRV(spreadsheeetCred)
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

		err = InvokeGenerationFromSheet()
		if err != nil {
			// meshkit
			utils.Log.Error(err)
			return err
		}

		return err
	},
}

func InvokeGenerationFromSheet() error {
	utils.Log.UpdateLogOutput(logFile)
	fmt.Println("TEST LOG LEVEL ", utils.Log.GetLevel())
	defer func() {
		_ = logFile.Close()
		utils.Log.UpdateLogOutput(os.Stdout)
		// Change this
		utils.Log.Info(fmt.Sprintf("Updated %d models, updated %d components", totalAggregateModel, totalAggregateComponents))

		utils.Log.Info("refer ", logDirPath, " for detailed registry generate logs")

		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()

	url := GoogleSpreadSheetURL + spreadsheeetID

	modelCSVHelper, err := utils.NewModelCSVHelper(url, "Models", sheetGID)
	if err != nil {
		return err
	}

	modelCSVHelper.ParseModelsSheet(false)
	utils.Log.Info("total models: ", len(modelCSVHelper.Models))
	weightedSem := semaphore.NewWeighted(20)
	pwd, _ := os.Getwd()
	// return nil
	var wg sync.WaitGroup
	for _, model := range modelCSVHelper.Models {
		ctx := context.Background()

		err := weightedSem.Acquire(ctx, 1)
		if err != nil {
			break // or sometheing else?
		}
		utils.Log.Info("Current model: ", model.Model)
		wg.Add(1)
		go func(model utils.ModelCSV) {
			defer func() {
				wg.Done()
				weightedSem.Release(1)
			}()

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
			utils.Log.Info("\nAFTER GET PACKAGE FOR MODEL", model.Model, " : ERR", err)
			version := pkg.GetVersion()

			modelDef := model.CreateModelDefinition(version)
			modelDefPath := filepath.Join(pwd, outputLocation, modelDef.Name)
			err = os.MkdirAll(modelDefPath, 0755)
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				utils.Log.Error(err)
				return
			}
			modelFilePath := fmt.Sprintf("%s/model.json", modelDefPath)
			err = mutils.WriteJSONToFile[v1alpha1.Model](modelFilePath, modelDef)
			if err != nil {
				utils.Log.Info("ERR GENERATE MODEL DEFINITION FOR MODEL : ", model.Model)
				utils.Log.Error(ErrGenerateModel(err, modelDefPath))
				return
			}

			utils.Log.Info("\nAFTER GET PACKAGE NO ERR", version)
			comps, err := pkg.GenerateComponents()
			if err != nil {
				utils.Log.Info("\nAFTER GENERATE COMPS FOR MODEL", model.Model, ": ERR")
				utils.Log.Error(ErrGenerateComponent(err, model.Model))
				return
			}
			utils.Log.Info("\nAFTER GENERATE COMP NO ERR")
			utils.Log.Info("Extracted", len(comps), "for model %s", model.ModelDisplayName)

			dirName := filepath.Join(outputLocation, model.Model, version)
			_, err = os.Stat(dirName)

			if errors.Is(err, os.ErrNotExist) {
				err = os.MkdirAll(filepath.Join(pwd, dirName), 0755)
				if err != nil {
					utils.Log.Error(ErrGenerateComponent(err, model.Model))
					return
				}
			}

			for _, comp := range comps {
				location := fmt.Sprintf("%s%s", filepath.Join(dirName, comp.Kind), ".json")
				err := mutils.WriteJSONToFile[v1alpha1.ComponentDefinition](location, comp)
				if err != nil {
					utils.Log.Info("INSIDE COMPS : ERR", err)
					utils.Log.Info(err)
				}
			}

		}(model)

	}
	wg.Wait()
	return nil
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
