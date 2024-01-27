package registry

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/registry/helpers"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

// gid: 1502185467
// sid: 17_ZLXFkBUdTQu39t49sZTSU4sVMumBRhXOxIP8Jw-nw
//https://docs.google.com/spreadsheets/d/e/2PACX-1vRPea3kydCnCYHBe2nmj8fV-XJeVi1nTQMkMWX0MYwyKeORpooM7at2v5RNpYEu-iRQoCi3xS-JP4gO/pub?output=csv

var (
	spreadsheeetID   string
	spreadsheeetCred string

	outputLocation string

	pathToRegistrantConnDefinition string
	pathToRegistrantCredDefinition string
	GoogleSpreadSheetURL           = "https://docs.google.com/spreadsheets/d/17_ZLXFkBUdTQu39t49sZTSU4sVMumBRhXOxIP8Jw-nw"
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

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("test ", spreadsheeetCred, spreadsheeetID, pathToRegistrantConnDefinition, pathToRegistrantCredDefinition)
		var err error
		if spreadsheeetID != "" {
			// err := os.Setenv("CRED", spreadsheeetCred)
			// fmt.Println("SET ENV ERR : ", err)
			err = InvokeGenerationFromSheet()
			if err != nil {
				// meshkit
				utils.Log.Error(err)
				return err
			}
		} else {
			fmt.Println("test.....", utils.Log.GetLevel())
			utils.Log.Info("Model Generation from registrant definitions not yet supproted.")
			return nil
		}

		return err
	},
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println(args)
	// },
}

func InvokeGenerationFromSheet() error {
	// _spreadsheetID, _ := strconv.Atoi(spreadsheeetID)
	modelCSVHelper, err := helpers.NewModelCSVHelper(GoogleSpreadSheetURL, "Models", 1502185467)
	if err != nil {
		return err // add meshkit err
	}
	if err != nil {
		return err // add meshkit err
	}
	modelCSVHelper.ParseModelsSheet()
	fmt.Println("length of models: ", len(modelCSVHelper.Models))
	weightedSem := semaphore.NewWeighted(20)
	var wg sync.WaitGroup
	for _, model := range modelCSVHelper.Models {
		ctx := context.Background()

		err := weightedSem.Acquire(ctx, 1)
		if err != nil {
			break // or sometheing else?
		}
		utils.Log.Info("Current model: ", model.Model)
		wg.Add(1)
		go func(model helpers.ModelCSV) {
			defer func() {
				defer wg.Done()
				weightedSem.Release(1)
			}()

			generator, err := generators.NewGenerator(model.Registrant, model.SourceURL, model.Model)
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			if model.Registrant == "artifacthub" {
				fmt.Println("QWERTY", model.Model)
				time.Sleep(1 * time.Second)
			}
			pkg, err := generator.GetPackage()
			fmt.Println("AFTER GET PACKAGE : ERR", err)
			if err != nil {
				utils.Log.Error(ErrGenerateModel(err, model.Model))
				return
			}
			comps, err := pkg.GenerateComponents()
			if err != nil {
				fmt.Println("AFTER GENERATE COMPS : ERR", err.Error())
				utils.Log.Error(ErrGenerateComponent(err, model.Model))
				return
			}
			fmt.Println("AFTER GENERATE COMP NO ERR")
			utils.Log.Info("Extracted", len(comps), "for model %s", model.ModelDisplayName)

			dirName := filepath.Join(outputLocation, model.Model)
			_, err = os.Stat(dirName)

			if errors.Is(err, os.ErrNotExist) {
				pwd, _ := os.Getwd()
				err = os.MkdirAll(filepath.Join(pwd, dirName), 0755)
				if err != nil {
					utils.Log.Error(ErrGenerateComponent(err, model.Model))
					return
				}
			}

			for _, comp := range comps {
				location := fmt.Sprintf("%s%s", filepath.Join(dirName, comp.Kind), ".json")
				err := writeCompsToFileSystem[v1alpha1.ComponentDefinition](location, comp)
				if err != nil {
					fmt.Println("INSIDE COMPS : ERR", err)
					utils.Log.Info(err)
				}
			}

		}(model)
		// Create an new logger instance and write logs to file instead?

	}
	wg.Wait()
	return nil
}

func writeCompsToFileSystem[K any](outputPath string, data K) error {
	byt, err := json.MarshalIndent(data, "", "")
	if err != nil {
		return utils.ErrMarshal(err)
	}
	fmt.Println("LINE 153 : ")
	file, err := os.Create(outputPath)
	if err != nil {
		return mutils.ErrCreateFile(err, outputPath)
	}
	fmt.Println("LINE 158 : ")
	_, err = file.Write(byt)
	if err != nil {
		return mutils.ErrWriteFile(err, outputPath)
	}
	fmt.Println("LINE 163 : ")
	return nil

}

func init() {
	importCmd.PersistentFlags().StringVar(&spreadsheeetID, "spreadsheet_id", "", "spreadsheet it for the integration spreadsheet")
	// importCmd.PersistentFlags().StringVar(&spreadsheeetCred, "spreadsheet_cred", "", "base64 encoded credential to download the spreadsheet")

	// importCmd.MarkFlagsRequiredTogether("spreadsheet_id", "spreadsheet_cred")

	importCmd.PersistentFlags().StringVar(&pathToRegistrantConnDefinition, "registrant_def", "", "path pointing to the registrant connection definition")
	importCmd.PersistentFlags().StringVar(&pathToRegistrantCredDefinition, "registrant_cred", "", "path pointing to the registrant credetial definition")

	importCmd.MarkFlagsRequiredTogether("registrant_def", "registrant_cred")

	importCmd.MarkFlagsMutuallyExclusive("spreadsheet_id", "registrant_def")
	// importCmd.MarkFlagsMutuallyExclusive("spreadsheet_cred", "registrant_cred")
	importCmd.PersistentFlags().StringVarP(&outputLocation, "output", "o", "../server/meshmodels", "location to output generated models, defaults to ../server/meshmodels")
	// importCmd.MarkFlagsMutuallyExclusive()

}
