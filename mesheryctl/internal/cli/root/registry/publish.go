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
	"slices"
	"sort"
	"strings"

	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	meshkitUtils "github.com/meshery/meshkit/utils"
)

var (
	system                string
	googleSheetCredential string
	sheetID               string
	modelsOutputPath      string
	imgsOutputPath        string
	models                = []meshkitRegistryUtils.ModelCSV{}
	components            = map[string]map[string][]meshkitRegistryUtils.ComponentCSV{}
	relationships         = []meshkitRegistryUtils.RelationshipCSV{}
	outputFormat          string
)

var supportedPublishSystems = []string{"meshery", "remote-provider", "website"}

const publishUsage = `Usage:
mesheryctl registry publish --system [system] --spreadsheet-cred [credential] --spreadsheet-id [sheet-id] --models-output-path [models-output-path] --imgs-output-path [imgs-output-path]
mesheryctl registry publish --system [system] --spreadsheet-cred [credential] --spreadsheet-id [sheet-id] --models-output-path [models-output-path] --imgs-output-path [imgs-output-path] -o [output-format]
Run 'mesheryctl registry publish --help' for usage instructions`

// publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery Server",
	Long: `Publishes metadata about Meshery Models to Websites, Remote Provider, or Meshery Server, including model and component icons by reading from a Google Spreadsheet and outputing to markdown or json format.
Find more information at: https://docs.meshery.io/reference/references/mesheryctl/registry/publish`,
	Example: `
// Publish To System
mesheryctl registry publish --system [system] --spreadsheet-cred [credential] --spreadsheet-id [sheet-id] --models-output-path [models-output-path] --imgs-output-path [imgs-output-path] -o [output-format]

// Publish To Meshery
mesheryctl registry publish --system meshery --spreadsheet-cred "$CRED" --spreadsheet-id GoogleSheetID --models-output-path [repo]/models

// Publish To Remote Provider
mesheryctl registry publish --system remote-provider --spreadsheet-cred "$CRED" --spreadsheet-id GoogleSheetID --models-output-path [repo]/meshmodels/models --imgs-output-path [repo]/ui/public/img/meshmodels

// Publish To Website
mesheryctl registry publish --system website --spreadsheet-cred "$CRED" --spreadsheet-id GoogleSheetID --models-output-path [repo]/integrations --imgs-output-path [repo]/ui/public/img/meshmodels -o md

// Publishing to meshery docs
cd docs;
mesheryctl registry publish --system website --spreadsheet-cred "$CRED" --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --models-output-path docs/pages/integrations --imgs-output-path docs/assets/img/integrations -o md

// Publishing to mesheryio site
mesheryctl registry publish --system website --spreadsheet-cred "$CRED" --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --models-output-path meshery.io/integrations --imgs-output-path meshery.io/assets/images/integration -o js

// Publishing to any website
mesheryctl registry publish --system website --spreadsheet-cred "$CRED" --spreadsheet-id 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --models-output-path path/to/models --imgs-output-path path/to/icons -o mdx
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return ErrPublishInvalidArgs(fmt.Sprintf("unexpected positional arguments: %s. registry publish now takes flags instead of positional arguments\n\n%s", strings.Join(args, " "), publishUsage))
		}

		requiredFlags := []string{"system", "spreadsheet-cred", "spreadsheet-id", "models-output-path"}
		var missing []string
		for _, name := range requiredFlags {
			value, _ := cmd.Flags().GetString(name)
			if strings.TrimSpace(value) == "" {
				missing = append(missing, "--"+name)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			return ErrPublishInvalidArgs(fmt.Sprintf("missing required flag(s): %s\n\n%s", strings.Join(missing, ", "), publishUsage))
		}

		sys, _ := cmd.Flags().GetString("system")
		if !slices.Contains(supportedPublishSystems, sys) {
			return ErrPublishInvalidArgs(fmt.Sprintf("invalid system: '%s'. Supported systems are %s\n\n%s", sys, strings.Join(supportedPublishSystems, ", "), publishUsage))
		}

		// remote-provider and website always write icons, so imgs-output-path is mandatory.
		imgsPath, _ := cmd.Flags().GetString("imgs-output-path")
		if sys != "meshery" && strings.TrimSpace(imgsPath) == "" {
			return ErrPublishInvalidArgs(fmt.Sprintf("--imgs-output-path is required for the '%s' system\n\n%s", sys, publishUsage))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		srv, err := meshkitUtils.NewSheetSRV(googleSheetCredential)
		if err != nil {
			return errors.New(utils.RegistryError("Invalid JWT Token: Ensure the provided token is a base64-encoded, valid Google Spreadsheets API token.", "publish"))
		}
		resp, err := srv.Spreadsheets.Get(sheetID).Fields().Do()
		if err != nil || resp.HTTPStatusCode != 200 {
			errMsg := fmt.Sprintf("Request to Google Spreadsheet did not succeed.\n\nReturned error: %s", err.Error())
			return errors.New(utils.RegistryError(errMsg, "publish"))
		}

		modelCSVHelper := &meshkitRegistryUtils.ModelCSVHelper{}
		componentCSVHelper := &meshkitRegistryUtils.ComponentCSVHelper{}
		relationshipCSVHelper := &meshkitRegistryUtils.RelationshipCSVHelper{}
		GoogleSpreadSheetURL += sheetID

		for _, v := range resp.Sheets {
			switch v.Properties.Title {
			case "Models":
				modelCSVHelper, err = meshkitRegistryUtils.NewModelCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId, modelCSVFilePath)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				err := modelCSVHelper.ParseModelsSheet(true, modelName)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			case "Components":
				componentCSVHelper, err = meshkitRegistryUtils.NewComponentCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId, componentCSVFilePath)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				err := componentCSVHelper.ParseComponentsSheet(modelName)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			case "Relationships":
				relationshipCSVHelper, err = meshkitRegistryUtils.NewRelationshipCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId, relationshipCSVFilePath)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
				err = relationshipCSVHelper.ParseRelationshipsSheet(modelName)
				if err != nil {
					utils.Log.Error(err)
					return nil
				}
			}
		}

		models = modelCSVHelper.Models
		components = componentCSVHelper.Components
		relationships = relationshipCSVHelper.Relationships

		switch system {
		case "meshery":
			err = mesherySystem()
		case "remote-provider":
			err = remoteProviderSystem()
		case "website":
			if outputFormat != "md" && outputFormat != "mdx" && outputFormat != "js" {
				return errors.New(utils.RegistryError(fmt.Sprintf("invalid output format: %s", outputFormat), "publish"))
			}
			err = websiteSystem()
		default:
			return ErrPublish(fmt.Errorf("invalid system: %s", system), system)
		}

		if err != nil {
			return ErrPublish(err, system)
		}

		err = modelCSVHelper.Cleanup()
		if err != nil {
			return ErrPublish(err, system)
		}

		err = componentCSVHelper.Cleanup()
		if err != nil {
			return ErrPublish(err, system)
		}

		return nil
	},
}

// TODO
func mesherySystem() error {
	return nil
}

// Create models definitions to remote provider path
// and add models icons to image output path
func remoteProviderSystem() error {
	// Construct absolute path to store models
	outputPath, _ := filepath.Abs(filepath.Join("../", modelsOutputPath))
	modelDir := filepath.Join(outputPath)
	totalModelsPublished := 0
	for _, model := range models {
		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			utils.Log.Debug("no components found for ", model.Model)
			comps = []meshkitRegistryUtils.ComponentCSV{}
		}

		err := utils.GenerateIcons(model, comps, imgsOutputPath)
		if err != nil {
			utils.Log.Debug(utils.ErrGeneratingIcons(err, imgsOutputPath))
			utils.Log.Fatalf("Error generating icons for model %s: %v", model.Model, err.Error())
		}

		_, _, err = WriteModelDefToFileSystem(&model, "", modelDir)
		if err != nil {
			return ErrGenerateModel(err, model.Model)
		}
		totalModelsPublished++
	}
	utils.Log.Info("Total model published: ", totalModelsPublished)
	return nil
}

func websiteSystem() error {
	var err error

	relationshipMap := make(map[string][]meshkitRegistryUtils.RelationshipCSV)
	for _, rel := range relationships {
		relationshipMap[rel.Model] = append(relationshipMap[rel.Model], rel)
	}
	docsJSON := "const data = ["
	for _, model := range models {
		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			utils.Log.Debug("no components found for ", model.Model)
			comps = []meshkitRegistryUtils.ComponentCSV{}
		}

		relnships, ok := relationshipMap[model.Model]
		if !ok || len(relnships) == 0 {
			utils.Log.Debug("no relationships found for ", model.Model)
			relnships = []meshkitRegistryUtils.RelationshipCSV{}
		}
		switch outputFormat {
		case "mdx":
			err := utils.GenerateMDXStyleDocs(model, comps, modelsOutputPath, imgsOutputPath) // creates mdx file
			if err != nil {
				utils.Log.Fatalf("Error generating remote provider docs for model %s: %v", model.Model, err.Error())
			}
		case "md":
			err := utils.GenerateMDStyleDocs(model, comps, relnships, modelsOutputPath, imgsOutputPath) // creates md file
			if err != nil {
				utils.Log.Fatalf("Error generating meshery docs for model %s: %v\n", model.Model, err.Error())
			}
		case "js":
			docsJSON, err = utils.GenerateJSStyleDocs(model, docsJSON, comps, relnships, modelsOutputPath, imgsOutputPath) // json file
			if err != nil {
				utils.Log.Fatalf("Error generating mesheryio docs for model %s: %v\n", model.Model, err.Error())
			}
		}

	}

	if outputFormat == "js" {
		docsJSON = strings.TrimSuffix(docsJSON, ",")
		docsJSON += "]; export default data"
		mOut, _ := filepath.Abs(filepath.Join(modelsOutputPath, "data.js"))
		if err := meshkitUtils.WriteToFile(mOut, docsJSON); err != nil {
			utils.Log.Error(err)
			return nil
		}
	}

	return nil
}

func init() {
	publishCmd.Flags().StringVarP(&system, "system", "s", "", "system to publish to [meshery | remote-provider | website]")
	publishCmd.Flags().StringVar(&googleSheetCredential, "spreadsheet-cred", "", "base64 encoded credential to download the spreadsheet")
	publishCmd.Flags().StringVar(&sheetID, "spreadsheet-id", "", "spreadsheet ID for the integration spreadsheet")
	publishCmd.Flags().StringVarP(&modelsOutputPath, "models-output-path", "m", "", "path to write the published models to")
	publishCmd.Flags().StringVarP(&imgsOutputPath, "imgs-output-path", "i", "", "path to write model and component icons to (required for remote-provider and website)")
	publishCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "", "output format [md | mdx | js]")
}

func WriteModelDefToFileSystem(model *meshkitRegistryUtils.ModelCSV, version string, location string) (string, *model.ModelDefinition, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	modelDefPath := filepath.Join(location, modelDef.Name)
	err := modelDef.WriteModelDefinition(modelDefPath+"/model.json", "json")
	if err != nil {
		return "", nil, err
	}

	return modelDefPath, &modelDef, nil
}
