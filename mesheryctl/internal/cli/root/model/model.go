// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by a, filepath.Dir(${1:}modelDefPathpplicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/layer5io/meshery/server/models"
	"github.com/manifoldco/promptui"
	"github.com/meshery/schemas/models/v1beta1/model"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// flag used to specify the page number in list command
	pageNumberFlag int
	// flag used to specify format of output of view {model-name} command
	outFormatFlag string

	// flag used to specify output location of export {model-name} command
	outLocationFlag string
	// flag used to specify format of output of export {model-name} command
	outTypeFlag string
	// flag used to specify whether to discard components in the model
	discardComponentsFlag bool
	// flag used to specify whether to discard relationships in the model
	discardRelationshipsFlag bool
	// flag used to specify the version of the model
	versionFlag string

	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25

	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)

	availableSubcommands = []*cobra.Command{listModelCmd, viewModelCmd, searchModelCmd, importModelCmd, exportModal}

	countFlag bool
)

// represents the mesheryctl model view [model-name] subcommand.

// represents the mesheryctl model search [query-text] subcommand.

// ModelCmd represents the mesheryctl model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "View list of models and detail of models",
	Long:  "View list of models and detailed information of a specific model",
	Example: `
// To view total of available models
mesheryctl model --count

// To view list of models
mesheryctl model list

// To view a specific model
mesheryctl model view [model-name]

// To search for a specific model
mesheryctl model search [model-name]
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !countFlag {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if countFlag {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				log.Fatalln(err, "error processing config")
			}

			baseUrl := mctlCfg.GetBaseMesheryURL()
			url := fmt.Sprintf("%s/api/meshmodels/models?page=1", baseUrl)
			return listModel(cmd, url, countFlag)
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl model --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	listModelCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	viewModelCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")

	exportModal.Flags().StringVarP(&outFormatFlag, "output-format", "t", "yaml", "(optional) format to display in [json|yaml] (default = yaml)")
	exportModal.Flags().StringVarP(&outLocationFlag, "output-location", "l", "./", "(optional) output location (default = current directory)")
	exportModal.Flags().StringVarP(&outTypeFlag, "output-type", "o", "oci", "(optional) format to display in [oci|tar] (default = oci)")
	exportModal.Flags().BoolVarP(&discardComponentsFlag, "discard-components", "c", false, "(optional) whether to discard components in the exported model definition (default = false)")
	exportModal.Flags().BoolVarP(&discardRelationshipsFlag, "discard-relationships", "r", false, "(optional) whether to discard relationships in the exported model definition (default = false)")
	exportModal.Flags().StringVarP(&versionFlag, "version", "", "", "(optional) model version to export (default = \"\")")

	ModelCmd.AddCommand(availableSubcommands...)
	ModelCmd.Flags().BoolVarP(&countFlag, "count", "", false, "(optional) Get the number of models in total")
}

// selectModelPrompt lets user to select a model if models are more than one
func selectModelPrompt(models []model.ModelDefinition) model.ModelDefinition {
	modelArray := []model.ModelDefinition{}
	modelNames := []string{}

	modelArray = append(modelArray, models...)

	for _, model := range modelArray {
		modelName := fmt.Sprintf("%s, version: %s", model.DisplayName, model.Version)
		modelNames = append(modelNames, modelName)
	}

	prompt := promptui.Select{
		Label: "Select a model",
		Items: modelNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return modelArray[i]
	}
}

func outputJson(model model.ModelDefinition) error {
	if err := prettifyJson(model); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(model, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}

// prettifyJson takes a model.ModelDefinition struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyJson(model model.ModelDefinition) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(model)
}

func listModel(cmd *cobra.Command, url string, displayCountOnly bool) error {
	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	// defers the closing of the response body after its use, ensuring that the resources are properly released.
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	modelsResponse := &models.MeshmodelsAPIResponse{}
	err = json.Unmarshal(data, modelsResponse)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	header := []string{"Model", "Category", "Version"}
	rows := [][]string{}

	for _, model := range modelsResponse.Models {
		if len(model.DisplayName) > 0 {
			rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
		}
	}

	if len(rows) == 0 {
		// if no model is found
		// fmt.Println("No model(s) found")
		whiteBoardPrinter.Println("No model(s) found")
		return nil
	}

	utils.DisplayCount("models", modelsResponse.Count)

	if displayCountOnly {
		return nil
	}

	if cmd.Flags().Changed("page") {
		utils.PrintToTable(header, rows)
	} else {
		err := utils.HandlePagination(maxRowsPerPage, "models", rows, header)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
	}

	return nil
}

func exportModel(modelName string, _ *cobra.Command, url string, _ bool) error {
	// Find the entity with the model name
	req, err := utils.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	// ensure proper cleaning of resources
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	var exportedModelPath string
	if outTypeFlag != "oci" {
		exportedModelPath = filepath.Join(outLocationFlag, modelName+"."+"tar.gz")
	} else {
		exportedModelPath = outLocationFlag + modelName + ".tar"
	}
	err = os.WriteFile(exportedModelPath, data, 0644)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	utils.Log.Infof("Exported model to %s", exportedModelPath)
	return nil
}
