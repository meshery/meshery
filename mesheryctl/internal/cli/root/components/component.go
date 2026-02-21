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

	"github.com/meshery/meshery/server/models"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type componentFlags struct {
	count bool
}

var (
	availableSubcommands = []*cobra.Command{listComponentCmd, viewComponentCmd, searchComponentsCmd}

	componentApiPath       = "api/meshmodels/components"
	componentFlagsProvided = &componentFlags{}
)

// ComponentCmd represents the mesheryctl component command
var ComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Manage Meshery components",
	Long: `List, search and view component(s) and detailed informations
Find more information at: https://docs.meshery.io/reference/mesheryctl/component`,
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
		argsIsEmpty := len(args) == 0 || (len(args) == 1 && args[0] == "")
		if argsIsEmpty && !componentFlagsProvided.count {
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

		err := cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	ComponentCmd.AddCommand(availableSubcommands...)
	ComponentCmd.Flags().BoolVarP(&componentFlagsProvided.count, "count", "", false, "(optional) Get the number of components in total")
}

func generateComponentDataToDisplay(componentsResponse *models.MeshmodelComponentsAPIResponse) ([][]string, int64) {
	rows := [][]string{}
	for _, component := range componentsResponse.Components {
		componentName := component.DisplayName
		if componentName == "" {
			componentName = "N/A"
		}
		modelName := component.Model.Name
		if modelName == "" {
			modelName = "N/A"
		}
		componentVersion := component.Component.Version
		if componentVersion == "" {
			componentVersion = "N/A"
		}
		componenttKind := component.Component.Kind
		if componenttKind == "" {
			componenttKind = "N/A"
		}
		rows = append(rows, []string{componentName, modelName, componenttKind, componentVersion})
	}

	return rows, int64(componentsResponse.Count)
}

func listComponents(cmd *cobra.Command, apiPath string) error {
	page, _ := cmd.Flags().GetInt("page")

	modelData := display.DisplayDataAsync{
		UrlPath:          componentApiPath,
		DataType:         "component",
		Header:           []string{"Name", "Model", "Category", "Version"},
		Page:             page,
		IsPage:           cmd.Flags().Changed("page"),
		DisplayCountOnly: cmd.Flags().Changed("count"),
	}

	return display.ListAsyncPagination(modelData, generateComponentDataToDisplay)
}
