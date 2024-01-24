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
	"os"
	"path/filepath"
	"strings"

	// "github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	// "github.com/layer5io/meshery/mesheryctl/pkg/utils"
	// "github.com/pkg/errors"

	// log "github.com/sirupsen/logrus"


	"github.com/spf13/cobra"
	// "github.com/spf13/viper"
)

var (
	system string
	googleSheetCredential string
	sheetID string
	modelsOutputPath string
	imgsOutputPath string
)

// publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery",
	Long:  `Publishes metadata about Meshery Models to Websites, Remote Provider and Meshery by reading from a Google Spreadsheet.`,
	Example: `
// Publish To System
mesheryctl publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]

// Publish To Meshery
mesheryctl registry publish meshery GoogleCredential GoogleSheetID <repo>/server/meshmodel

// Publish To Remote Provider
mesheryctl registry publish remote-provider GoogleCredential GoogleSheetID <repo>/meshmodels/models <repo>/ui/public/img/meshmodels

// Publish To Website
mesheryctl registry publish website GoogleCredential $GoogleSheetID <repo>/integrations <repo>/ui/public/img/meshmodels
	`,
	// PreRunE: func(cmd *cobra.Command, args []string) error {
	// 	//Check prerequisite
	// },
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 5 {
			return cmd.Help()
		}

		system = args[0]
		
		switch system {
			case "meshery":
				return mesherySystem()
			case "remote-provider":
				return remoteProviderSystem()
			case "website":
				return websiteSystem()
			default:
			return fmt.Errorf("invalid system: %s", system) // update to meshkit
		}
	},
}

func mesherySystem() error {
	return nil
}

func remoteProviderSystem() error {
	return nil
}

func websiteSystem() error {
	pathToIntegrationsLayer5 := os.Args[4]
	pathToIntegrationsMeshery := os.Args[5]
	pathToIntegrationsMesheryDocs := os.Args[6]
	mesheryioDocsJSON := "const data = ["
	for _, model := range models {
		pathForLayer5ioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsLayer5))
		pathForMesheryioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsMeshery))
		pathForMesheryDocsIntegrations, _ := filepath.Abs(filepath.Join("../../", pathToIntegrationsMesheryDocs))

		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			fmt.Println("no components found for ", model.Model)
			comps = []pkg.ComponentCSV{}
		}

		err := pkg.GenerateLayer5Docs(model, comps, pathForLayer5ioIntegrations)
		if err != nil {
			fmt.Printf("Error generating layer5 docs for model %s: %v\n", model.Model, err.Error())
		}

		mesheryioDocsJSON, err = pkg.GenerateMesheryioDocs(model, pathForMesheryioIntegrations, mesheryioDocsJSON)
		if err != nil {
			fmt.Printf("Error generating mesheryio docs for model %s: %v\n", model.Model, err.Error())
		}

		err = pkg.GenerateMesheryDocs(model, comps, pathForMesheryDocsIntegrations)
		if err != nil {
			fmt.Printf("Error generating meshery docs for model %s: %v\n", model.Model, err.Error())
		}

	}

	mesheryioDocsJSON = strings.TrimSuffix(mesheryioDocsJSON, ",")
	mesheryioDocsJSON += "]; export default data"
	if err := pkg.WriteToFile(filepath.Join("../../../", pathToIntegrationsMeshery, "data.js"), mesheryioDocsJSON); err != nil {
		log.Fatal(err)
	}
}