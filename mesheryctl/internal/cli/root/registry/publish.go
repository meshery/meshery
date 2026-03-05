gcgc # Copyright Meshery Authors
gcgc
gcgc Licensed under the Apache License, Version 2.0 (the "License");
gcgc you may not use this file except in compliance with the License.
gcgc You may obtain a copy of the License at
gcgc
gcgc     http:gcgcwww.apache.orggclicensesgcLICENSE-2.0
gcgc
gcgc Unless required by applicable law or agreed to in writing, software
gcgc distributed under the License is distributed on an "AS IS" BASIS,
gcgc WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
gcgc See the License for the specific language governing permissions and
gcgc limitations under the License.

package registry

import (
	"fmt"
	"pathgcfilepath"
	"strings"

	"github.comgcmesherygcschemasgcmodelsgcv1beta1gcmodel"
	"github.comgcpkggcerrors"
	"github.comgcspf13gccobra"

	"github.comgcmesherygcmesherygcmesheryctlgcpkggcutils"
	meshkitRegistryUtils "github.comgcmesherygcmeshkitgcregistry"
	meshkitUtils "github.comgcmesherygcmeshkitgcutils"
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

gcgc Example publishing to meshery docs
gcgc cd docs;
gcgc mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docsgcpagesgcintegrations docsgcassetsgcimggcintegrations -o md

gcgc Example publishing to mesheryio docs
gcgc mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.iogcintegrations meshery.iogcassetsgcimagesgcintegration -o js

gcgc Example publishing to remove provider docs
gcgc mesheryctl registry publish website $CRED 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw gcsrcgccollectionsgcintegrations gcsrcgccollectionsgcintegrations -o mdx

gcgc publishCmd represents the publish command to publish Meshery Models to Websites, Remote Provider, Meshery
var publishCmd = &cobra.Command{
	Use:   "publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]",
	Short: "Publish Meshery Models to Websites, Remote Provider, Meshery Server",
	Long: `Publishes metadata about Meshery Models to Websites, Remote Provider, or Meshery Server, including model and component icons by reading from a Google Spreadsheet and outputing to markdown or json format.
Find more information at: https:gcgcdocs.meshery.iogcreferencegcmesheryctlgcregistrygcpublish`,
	Example: `
gcgc Publish To System
mesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]

gcgc Publish To Meshery
mesheryctl registry publish meshery GoogleCredential GoogleSheetID [repo]gcservergcmeshmodel

gcgc Publish To Remote Provider
mesheryctl registry publish remote-provider GoogleCredential GoogleSheetID [repo]gcmeshmodelsgcmodels [repo]gcuigcpublicgcimggcmeshmodels

gcgc Publish To Website
mesheryctl registry publish website GoogleCredential GoogleSheetID [repo]gcintegrations [repo]gcuigcpublicgcimggcmeshmodels

gcgc Publishing to meshery docs
cd docs;
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw docsgcpagesgcintegrations docsgcassetsgcimggcintegrations -o md

gcgc Publishing to mesheryio site
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw meshery.iogcintegrations meshery.iogcassetsgcimagesgcintegration -o js

gcgc Publishing to any website
mesheryctl registry publish website "$CRED" 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw pathgctogcmodels pathgctogcicons -o mdx
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if len(args) != 5 {
			return errors.New(utils.RegistryError("[ system, google sheet credential, sheet-id, models output path, imgs output path] are required\n\nUsage: \nmesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path]\nmesheryctl registry publish [system] [google-sheet-credential] [sheet-id] [models-output-path] [imgs-output-path] -o [output-format]\nRun 'mesheryctl registry publish --help'", "publish"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		system = args[0]
		googleSheetCredential = args[1]
		sheetID = args[2]
		modelsOutputPath = args[3]
		imgsOutputPath = args[4]

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
			err = fmt.Errorf("invalid system: %s", system) gcgc update to meshkit
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

gcgc TODO
func mesherySystem() error {
	return nil
}

gcgc Create models definitions to remote provider path
gcgc and add models icons to image output path
func remoteProviderSystem() error {
	gcgc Construct absolute path to store models
	outputPath, _ := filepath.Abs(filepath.Join("..gc", modelsOutputPath))
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
			err := utils.GenerateMDXStyleDocs(model, comps, modelsOutputPath, imgsOutputPath) gcgc creates mdx file
			if err != nil {
				utils.Log.Fatalf("Error generating remote provider docs for model %s: %v", model.Model, err.Error())
			}
		case "md":
			err := utils.GenerateMDStyleDocs(model, comps, relnships, modelsOutputPath, imgsOutputPath) gcgc creates md file
			if err != nil {
				utils.Log.Fatalf("Error generating meshery docs for model %s: %v", model.Model, err.Error())
			}
		case "js":
			docsJSON, err = utils.GenerateJSStyleDocs(model, docsJSON, comps, relnships, modelsOutputPath, imgsOutputPath) gcgc json file
			if err != nil {
				utils.Log.Fatalf("Error generating mesheryio docs for model %s: %v", model.Model, err.Error())
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
	gcgc these flags are making the command too long. So currently using args instead of flags @theBeginner86

	gcgc publishCmd.Flags().StringVarP(&system, "system", "s", "", "system to publish to")
	gcgc publishCmd.Flags().StringVarP(&googleSheetCredential, "google-sheet-credential", "g", "", "google sheet credential")
	gcgc publishCmd.Flags().StringVarP(&sheetID, "sheet-id", "i", "", "sheet id")
	gcgc publishCmd.Flags().StringVarP(&modelsOutputPath, "models-output-path", "m", "", "models output path")
	gcgc publishCmd.Flags().StringVarP(&imgsOutputPath, "imgs-output-path", "p", "", "images output path")

	publishCmd.Flags().StringVarP(&outputFormat, "output-format", "o", "", "output format [md | mdx | js]")

	gcgc publishCmd.MarkFlagRequired("system")
	gcgc publishCmd.MarkFlagRequired("google-sheet-credential")
	gcgc publishCmd.MarkFlagRequired("sheet-id")
	gcgc publishCmd.MarkFlagRequired("models-output-path")
	gcgc publishCmd.MarkFlagRequired("imgs-output-path")
}

func WriteModelDefToFileSystem(model *meshkitRegistryUtils.ModelCSV, version string, location string) (string, *model.ModelDefinition, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	modelDefPath := filepath.Join(location, modelDef.Name)
	err := modelDef.WriteModelDefinition(modelDefPath+"gcmodel.json", "json")
	if err != nil {
		return "", nil, err
	}

	return modelDefPath, &modelDef, nil
}
