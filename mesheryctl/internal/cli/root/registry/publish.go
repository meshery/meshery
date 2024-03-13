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

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitUtils "github.com/layer5io/meshkit/utils"
)

var (
	system                string
	googleSheetCredential string
	sheetID               string
	modelsOutputPath      string
	imgsOutputPath        string
	models                = []utils.ModelCSV{}
	components            = map[string]map[string][]utils.ComponentCSV{}
	outputFormat          string
)

// Example publishing to meshery docs
// cd docs;
// mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docs/pages/integrations docs/assets/img/integrations -o md

// Example publishing to mesheryio docs
// mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.io/integrations meshery.io/assets/images/integration -o js

// Example publishing to layer5 docs
// mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw layer5/src/collections/integrations layer5/src/collections/integrations -o mdx

// publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery Server",
	Long:  `Publishes metadata about Meshery Models to Websites, Remote Provider, or Meshery Server, including model and component icons by reading from a Google Spreadsheet and outputing to markdown or json format.`,
	Example: `
// Publish To System
mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]

// Publish To Meshery
mesheryctl registry publish meshery GoogleCredential GoogleSheetID [repo]/server/meshmodel

// Publish To Remote Provider
mesheryctl registry publish remote-provider GoogleCredential GoogleSheetID [repo]/meshmodels/models [repo]/ui/public/img/meshmodels

// Publish To Website
mesheryctl registry publish website GoogleCredential GoogleSheetID [repo]/integrations [repo]/ui/public/img/meshmodels

// Publishing to meshery docs
cd docs;
mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docs/pages/integrations docs/assets/img/integrations -o md

// Publishing to mesheryio site
mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.io/integrations meshery.io/assets/images/integration -o js

// Publishing to layer5 site
mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw layer5/src/collections/integrations layer5/src/collections/integrations -o mdx

// Publishing to any website
mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw path/to/models path/to/icons -o mdx
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 5 {
			return cmd.Help()
		}

		system = args[0]
		googleSheetCredential = args[1]
		sheetID = args[2]
		modelsOutputPath = args[3]
		imgsOutputPath = args[4]

		if outputFormat != "md" && outputFormat != "mdx" && outputFormat != "js" {
			return errors.New(utils.RegistryError(fmt.Sprintf("invalid output format: %s", outputFormat), "publish"))
		}

		// move to meshkit
		srv, err := meshkitUtils.NewSheetSRV(googleSheetCredential)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		resp, err := srv.Spreadsheets.Get(sheetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			utils.Log.Error(err)
			return nil
		}

		modelCSVHelper := &utils.ModelCSVHelper{}
		componentCSVHelper := &utils.ComponentCSVHelper{}
		GoogleSpreadSheetURL += sheetID

		for _, v := range resp.Sheets {
			switch v.Properties.Title {
			case "Models":
				modelCSVHelper, err = utils.NewModelCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				err := modelCSVHelper.ParseModelsSheet(true)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			case "Components":
				componentCSVHelper, err = utils.NewComponentCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				err := componentCSVHelper.ParseComponentsSheet()
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			}
		}

		models = modelCSVHelper.Models
		components = componentCSVHelper.Components

		switch system {
		case "meshery":
			err = mesherySystem()
		case "remote-provider":
			err = remoteProviderSystem()
		case "website":
			err = websiteSystem()
		default:
			err = fmt.Errorf("invalid system: %s", system) // update to meshkit
		}

		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		err = modelCSVHelper.Cleanup()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		err = componentCSVHelper.Cleanup()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		return nil
	},
}

// TODO
func mesherySystem() error {
	return nil
}

// TODO
func remoteProviderSystem() error {
	return nil
}

func websiteSystem() error {
	var err error

	docsJSON := "const data = ["
	for _, model := range models {
		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			utils.Log.Debug("no components found for ", model.Model)
			comps = []utils.ComponentCSV{}
		}

		switch outputFormat {
		case "mdx":
			err := utils.GenerateMDXStyleDocs(model, comps, modelsOutputPath, imgsOutputPath) // creates mdx file
			if err != nil {
				log.Fatalln(fmt.Printf("Error generating layer5 docs for model %s: %v\n", model.Model, err.Error()))
			}
		case "md":
			err := utils.GenerateMDStyleDocs(model, comps, modelsOutputPath, imgsOutputPath) // creates md file
			if err != nil {
				log.Fatalln(fmt.Printf("Error generating meshery docs for model %s: %v\n", model.Model, err.Error()))
			}
		case "js":
			docsJSON, err = utils.GenerateJSStyleDocs(model, docsJSON, imgsOutputPath) // json file
			if err != nil {
				log.Fatalln(fmt.Printf("Error generating mesheryio docs for model %s: %v\n", model.Model, err.Error()))
			}
		}

	}

	if outputFormat == "js" {
		docsJSON = strings.TrimSuffix(docsJSON, ",")
		docsJSON += "]; export default data"
		mOut, _ := filepath.Abs(filepath.Join("../", modelsOutputPath, "data.js"))
		if err := meshkitUtils.WriteToFile(mOut, docsJSON); err != nil {
			utils.Log.Error(err)
			return nil
		}
	}

	return nil
}

func init() {
	// these flags are making the command too long. So currently using args instead of flags @theBeginner86

	// publishCmd.Flags().StringVarP(&system, "system", "s", "", "system to publish to")
	// publishCmd.Flags().StringVarP(&googleSheetCredential, "google-sheet-credential", "g", "", "google sheet credential")
	// publishCmd.Flags().StringVarP(&sheetID, "sheet-id", "i", "", "sheet id")
	// publishCmd.Flags().StringVarP(&modelsOutputPath, "models-output-path", "m", "", "models output path")
	// publishCmd.Flags().StringVarP(&imgsOutputPath, "imgs-output-path", "p", "", "images output path")

	publishCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "", "output format [md | mdx | js]")
	err := publishCmd.MarkFlagRequired("output-format")
	if err != nil {
		utils.Log.Error(err)
	}

	// publishCmd.MarkFlagRequired("system")
	// publishCmd.MarkFlagRequired("google-sheet-credential")
	// publishCmd.MarkFlagRequired("sheet-id")
	// publishCmd.MarkFlagRequired("models-output-path")
	// publishCmd.MarkFlagRequired("imgs-output-path")
}
