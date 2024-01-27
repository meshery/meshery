package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	// "sync"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/spf13/cobra"
)

var (
	modelLocation string
)

var updateCmd = &cobra.Command{
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
		fmt.Println("model location: ", modelLocation)

		if spreadsheeetID != "" {
			// err := os.Setenv("CRED", spreadsheeetCred)
			// fmt.Println("SET ENV ERR : ", err)
			err := InvokeGenerationFromSheet()
			if err != nil {
				// meshkit
				utils.Log.Error(err)
				return err
			}
		}

		return nil
	},
}
var (
	ExcludeDirs = []string{"relationships", "policies"}
)

func InvokeCompUpdate() error {
	componentCSVHelper, err := utils.NewComponentCSVHelper(GoogleSpreadSheetURL, "Components", 1502185467)
	if err != nil {
		return err // add meshkit err
	}

	componentCSVHelper.ParseComponentsSheet()
	utils.Log.Info("total components: ", len(componentCSVHelper.Components))

	// Since component update doesn;t take long skip doing it concurrently
	// weightedSem := semaphore.NewWeighted(20)
	pwd, _ := os.Getwd()

	// var wg sync.WaitGroup

	for _, model := range componentCSVHelper.Components {
		for modelName, components := range model {
			modelPath := filepath.Join(pwd, modelLocation, modelName)
			utils.Log.Info("Updating component for model %s", modelName)

			modelVersions, err := os.ReadDir(modelPath)
			if err != nil {
				err = ErrUpdateModel(err, modelName)
				utils.Log.Error(err)
				return err
			}
			for _, version := range modelVersions {
				entries, _ := os.ReadDir(version.Name())
				for _, entry := range entries {
					if entry.IsDir() {
						if utils.Contains(entry.Name(), ExcludeDirs) != -1 {
							continue
						}
						componentsDirPath := filepath.Join(modelPath, version.Name(), entry.Name())
						for _, component := range components {
							utils.Log.Debug("Updating", component.Component)
							compPath := fmt.Sprintf("%s/%s", componentsDirPath, component.Component)
							componentByte, err := os.ReadFile(compPath)
							if err != nil {
								utils.Log.Error(ErrUpdateComponent(err, component.Component))
								continue
							}
							componentDef := v1alpha1.ComponentDefinition{}
							err = json.Unmarshal(componentByte, &componentDef)
							if err != nil {
								utils.Log.Error(ErrUpdateComponent(err, component.Component))
								continue
							}

							err = component.UpdateCompDefinition(&componentDef)
							if err != nil {
								utils.Log.Error(ErrUpdateComponent(err, component.Component))
								continue
							}
						}
					}
				}
			}

		}

	}
	return nil
}

func init() {
	updateCmd.PersistentFlags().StringVarP(&modelLocation, "path", "p", "", "relative or absolute path to the models directory")
}
