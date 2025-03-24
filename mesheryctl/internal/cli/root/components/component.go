// Copyright 2024 Meshery Authors
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

package components

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/component"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands = []*cobra.Command{listComponentCmd, viewComponentCmd, searchComponentsCmd}

	pageNumberFlag int
	outFormatFlag  string
	saveFlag       bool
	countFlag      bool
)

// ComponentCmd represents the mesheryctl component command
var ComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Manage components",
	Long:  "List, search and view component(s) and detailed informations",
	Example: `
// Display number of available components in Meshery
mesheryctl component --count

// List available component(s)
mesheryctl component list

// Search for component(s)
mesheryctl component search [component-name]

// View a specific component
mesheryctl component view [component-name]
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
			url := fmt.Sprintf("%s/api/meshmodels/components?page=1", baseUrl)
			return listComponents(cmd, url, countFlag)
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl component --help' to display usage guide.\n", args[0]), "model"))
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

// selectComponentPrompt lets user to select a model if models are more than one
func selectComponentPrompt(components []component.ComponentDefinition) component.ComponentDefinition {
	componentNames := []string{}
	componentArray := []component.ComponentDefinition{}

	componentArray = append(componentArray, components...)

	for _, component := range componentArray {
		componentName := fmt.Sprintf("%s, version: %s", component.DisplayName, component.Component.Version)
		componentNames = append(componentNames, componentName)
	}

	prompt := promptui.Select{
		Label: "Select component",
		Items: componentNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return componentArray[i]
	}
}

func init() {
	ComponentCmd.AddCommand(availableSubcommands...)
	ComponentCmd.Flags().BoolVarP(&countFlag, "count", "", false, "(optional) Get the number of components in total")
}

func listComponents(cmd *cobra.Command, url string, displayCountOnly bool) error {
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

	componentsResponse := &models.MeshmodelComponentsAPIResponse{}
	err = json.Unmarshal(data, componentsResponse)
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	header := []string{"Model", "kind", "Version"}
	rows := [][]string{}

	for _, component := range componentsResponse.Components {
		if len(component.DisplayName) > 0 {
			rows = append(rows, []string{component.Model.Name, component.Component.Kind, component.Component.Version})
		}
	}

	if len(rows) == 0 {
		// if no component is found
		fmt.Println("No components(s) found")
		return nil
	}

	utils.DisplayCount("components", componentsResponse.Count)
	if displayCountOnly {
		return nil
	}

	if cmd.Flags().Changed("page") {
		utils.PrintToTable(header, rows)
	} else {
		maxRowsPerPage := 25
		err := utils.HandlePagination(maxRowsPerPage, "components", rows, header)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
	}
	return nil
}
