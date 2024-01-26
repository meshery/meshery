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
	"fmt"
	"path/filepath"
	"strings"

	// "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	// "github.com/utils/errors"

	// log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
	system                        string
	googleSheetCredential         string
	sheetID                       string
	modelsOutputPath              string
	imgsOutputPath                string
	GoogleSpreadSheetURL          = "https://docs.google.com/spreadsheets/d/"
	models                        = []utils.ModelCSV{}
	components                    = map[string]map[string][]utils.ComponentCSV{}
	pathToIntegrationsLayer5      string
	pathToIntegrationsMeshery     string
	pathToIntegrationsMesheryDocs string
)

// publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery",
	Long:  `Publishes metadata about Meshery Models to Websites, Remote Provider and Meshery by reading from a Google Spreadsheet.`,
	Example: `
// Publish To System
mesheryctl exp registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]

// Publish To Meshery
mesheryctl exp registry publish meshery GoogleCredential GoogleSheetID <repo>/server/meshmodel

// Publish To Remote Provider
mesheryctl exp registry publish remote-provider GoogleCredential GoogleSheetID <repo>/meshmodels/models <repo>/ui/public/img/meshmodels

// Publish To Website
mesheryctl exp registry publish website GoogleCredential GoogleSheetID <repo>/integrations <repo>/ui/public/img/meshmodels
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 5 {
			return cmd.Help()
		}

		system = args[0]
		googleSheetCredential = args[1]
		sheetID = args[2]
		// modelsOutputPath = args[3]
		// imgsOutputPath = args[4]

		pathToIntegrationsLayer5 = args[3]
		pathToIntegrationsMeshery = args[4]
		pathToIntegrationsMesheryDocs = args[5]

		srv, err := utils.NewSheetSRV(googleSheetCredential)
		if err != nil {
			return err
		}
		resp, err := srv.Spreadsheets.Get(sheetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			return err
		}

		modelCSVHelper := &utils.ModelCSVHelper{}
		componentCSVHelper := &utils.ComponentCSVHelper{}
		GoogleSpreadSheetURL += sheetID

		for _, v := range resp.Sheets {
			switch v.Properties.Title {
			case "Models":
				modelCSVHelper, err = utils.NewModelCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
				if err != nil {
					return err
				}
				modelCSVHelper.ParseModelsSheet()
			case "Components":
				componentCSVHelper, err = utils.NewComponentCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
				if err != nil {
					return err
				}
				componentCSVHelper.ParseComponentsSheet()
			}
		}

		models = modelCSVHelper.Models
		components = componentCSVHelper.Components

		switch system {
		case "meshery":
			return mesherySystem()
		case "remote-provider":
			return remoteProviderSystem()
		case "website":
			return websiteSystem()
		default:
			fmt.Errorf("invalid system: %s", system) // update to meshkit
		}

		err = modelCSVHelper.Cleanup()
		if err != nil {
			return err
		}

		err = componentCSVHelper.Cleanup()
		if err != nil {
			return err
		}

		return nil
	},
}

func mesherySystem() error {
	return nil
}

func remoteProviderSystem() error {
	return nil
}

func websiteSystem() error {
	mesheryioDocsJSON := "const data = ["
	for _, model := range models {
		pathForLayer5ioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsLayer5))
		pathForMesheryioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsMeshery))
		pathForMesheryDocsIntegrations, _ := filepath.Abs(filepath.Join("../../", pathToIntegrationsMesheryDocs))

		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			fmt.Println("no components found for ", model.Model)
			comps = []utils.ComponentCSV{}
		}

		err := utils.GenerateLayer5Docs(model, comps, pathForLayer5ioIntegrations)
		if err != nil {
			fmt.Printf("Error generating layer5 docs for model %s: %v\n", model.Model, err.Error())
		}

		mesheryioDocsJSON, err = utils.GenerateMesheryioDocs(model, pathForMesheryioIntegrations, mesheryioDocsJSON)
		if err != nil {
			fmt.Printf("Error generating mesheryio docs for model %s: %v\n", model.Model, err.Error())
		}

		err = utils.GenerateMesheryDocs(model, comps, pathForMesheryDocsIntegrations)
		if err != nil {
			fmt.Printf("Error generating meshery docs for model %s: %v\n", model.Model, err.Error())
		}

	}

	mesheryioDocsJSON = strings.TrimSuffix(mesheryioDocsJSON, ",")
	mesheryioDocsJSON += "]; export default data"
	if err := utils.WriteToFile(filepath.Join("../../../", pathToIntegrationsMeshery, "data.js"), mesheryioDocsJSON); err != nil {
		fmt.Printf("Error writing to file: %v\n", err.Error())
		return err
	}

	return nil
}
