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
	"net/url"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/spf13/cobra"
)

type componentViewFlags struct {
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
}

var cmdComponentViewFlags componentViewFlags

// represents the mesheryctl component view [component-name] subcommand.
var viewComponentCmd = &cobra.Command{
	Use:   "view",
	Short: "View registered components",
	Long: `View a component registered in Meshery Server
Find more information at: https://docs.meshery.io/reference/mesheryctl/component/view`,
	Example: `
// View details of a specific component
mesheryctl component view [component-name | component-id]

// View details of a specific component in specified format
mesheryctl component view [component-name | component-id] -o [json|yaml]

// View details of a specific component in specified format and save it as a file
mesheryctl component view [component-name | component-id] -o [json|yaml] --save
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &cmdComponentViewFlags)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errInvalidArg, viewUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		componentDefinition := args[0]
		urlPath := componentApiPath
		searchTerm := ""

		// build url for uuid
		if utils.IsUUID(componentDefinition) {
			viewUrlValue := url.Values{}
			viewUrlValue.Add("id", componentDefinition)

			urlPath = fmt.Sprintf("%s?%s", urlPath, viewUrlValue.Encode())
		} else {
			searchTerm = componentDefinition
		}

		selectedComponent := new(component.ComponentDefinition)

		err := display.PromptAsyncPagination(
			display.DisplayDataAsync{
				UrlPath:        urlPath,
				ErrNotFoundMsg: fmt.Sprintf("%s%s", errNoComponentFound, componentDefinition),
				SearchTerm:     searchTerm,
			},
			formatLabel,
			func(data *models.MeshmodelComponentsAPIResponse) ([]component.ComponentDefinition, int64) {
				return data.Components, data.TotalCount
			},
			selectedComponent,
		)
		if err != nil {
			return err
		}

		outputFormatterFactory := display.OutputFormatterFactory[component.ComponentDefinition]{}
		outputFormatter, err := outputFormatterFactory.New(cmdComponentViewFlags.OutputFormat, *selectedComponent)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if cmdComponentViewFlags.Save {
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[component.ComponentDefinition]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(cmdComponentViewFlags.OutputFormat, outputFormatter)
			if err != nil {
				return err
			}

			componentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedComponent.DisplayName), " ", "_")
			fileName := filepath.Join(utils.MesheryFolder, fmt.Sprintf("component_%s.%s", componentString, cmdComponentViewFlags.OutputFormat))

			outputFormatterSaver = outputFormatterSaver.WithFilePath(fileName)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	viewComponentCmd.Flags().StringVarP(&cmdComponentViewFlags.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewComponentCmd.Flags().BoolVarP(&cmdComponentViewFlags.Save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}

func formatLabel(components []component.ComponentDefinition) []string {
	labels := []string{}

	for _, component := range components {
		name := fmt.Sprintf("%s, version: %s", component.DisplayName, component.Component.Version)
		labels = append(labels, name)
	}
	return labels
}
