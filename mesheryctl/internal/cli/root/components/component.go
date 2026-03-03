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
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type componentFlags struct {
	Count bool
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(componentFlagsProvided)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	Args: func(cmd *cobra.Command, args []string) error {
		argsIsEmpty := len(args) == 0 || (len(args) == 1 && args[0] == "")
		if argsIsEmpty && !componentFlagsProvided.Count {
			return utils.ErrInvalidArgument(errors.New("please provide a subcommand"))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if componentFlagsProvided.Count {
			componentData := display.DisplayDataAsync{
				UrlPath:          componentApiPath,
				DataType:         "component",
				Header:           []string{},
				Page:             cmdComponentListFlag.Page,
				PageSize:         10,
				IsPage:           false,
				DisplayCountOnly: componentFlagsProvided.Count,
			}
			return display.ListAsyncPagination(componentData, generateComponentDataToDisplay)
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
	ComponentCmd.Flags().BoolVarP(&componentFlagsProvided.Count, "count", "", false, "(optional) Get the number of components in total")
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
		rows = append(rows, []string{componentName, modelName, componentVersion})
	}

	return rows, int64(componentsResponse.Count)
}
