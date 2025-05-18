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
	"fmt"

	"github.com/layer5io/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/component"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
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
	saveFlag       bool

	componentApiPath = "api/meshmodels/components"
)

// ComponentCmd represents the mesheryctl component command
var ComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Manage components",
	Long: `List, search and view component(s) and detailed informations
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component`,
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
		count, _ := cmd.Flags().GetBool("count")
		if len(args) == 0 && !count {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetBool("count")
		if count {
			return listComponents(cmd, componentApiPath)
		}

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return utils.ErrInvalidArgument(fmt.Errorf("'%s' is an invalid subcommand. Please provide required options from [list, search, view]. Use 'mesheryctl component --help' to display usage guide", args[0]))
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
	ComponentCmd.Flags().BoolP("count", "", false, "(optional) Get the number of components in total")
}

func listComponents(cmd *cobra.Command, url string) error {
	componentsResponse, err := api.Fetch[models.MeshmodelComponentsAPIResponse](url)

	if err != nil {
		return err
	}

	header := []string{"Model", "kind", "Version"}
	rows := [][]string{}

	for _, component := range componentsResponse.Components {
		if len(component.DisplayName) > 0 {
			rows = append(rows, []string{component.Model.Name, component.Component.Kind, component.Component.Version})
		}
	}

	count, _ := cmd.Flags().GetBool("count")

	dataToDisplay := display.DisplayedData{
		DataType:         "components",
		Header:           header,
		Rows:             rows,
		Count:            int64(componentsResponse.Count),
		DisplayCountOnly: count,
		IsPage:           cmd.Flags().Changed("page"),
	}

	err = display.List(dataToDisplay)
	if err != nil {
		return err
	}

	return nil
}
