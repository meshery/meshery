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

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/spf13/cobra"
)

type componentViewFlags struct {
	OutputFormat string
	Save         bool
}

var cmdComponentViewFlags componentViewFlags

// represents the mesheryctl component view [component-name] subcommand.
var viewComponentCmd = &cobra.Command{
	Use:   "view",
	Short: "View registered components",
	Long: `View a component registered in Meshery Server
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/view`,
	Example: `
// View details of a specific component
mesheryctl component view [component-name]

// View details of a specific component in specifed format 
mesheryctl component view [component-name] -o [json|yaml]

// View details of a specific component in specified format and save it as a file
mesheryctl component view [component-name] -o [json|yaml] --save
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return display.ValidateOutputFormat(cmdComponentViewFlags.OutputFormat)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("[component name] is required but not specified\n\n%s", errViewCmdMsg))
		} else if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", errViewCmdMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		componentDefinition := args[0]

		viewUrlValue := url.Values{}
		viewUrlValue.Add("search", componentDefinition)
		viewUrlValue.Add("pagesize", "all")

		urlPath := fmt.Sprintf("%s?%s", componentApiPath, viewUrlValue.Encode())

		componentResponse, err := api.Fetch[models.MeshmodelComponentsAPIResponse](urlPath)
		if err != nil {
			return err
		}

		if componentResponse.Count == 0 {
			utils.Log.Info("No component(s) found for the given name: ", componentDefinition)
			return nil
		}

		var selectedComponent component.ComponentDefinition

		if componentResponse.Count == 1 {
			selectedComponent = componentResponse.Components[0] // Update the type of selectedModel
		} else {
			selectedComponent = selectComponentPrompt(componentResponse.Components)
		}

		outputFormatterFactory := display.OutputFormatterFactory[component.ComponentDefinition]{}
		outputFormatter, err := outputFormatterFactory.New(cmdComponentViewFlags.OutputFormat, selectedComponent)
		// outputFormatter.WithOutput(cmd.OutOrStdout())

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

// selectComponentPrompt lets user to select a model if models are more than one
func selectComponentPrompt(components []component.ComponentDefinition) component.ComponentDefinition {
	componentNames := make([]string, len(components))

	for i, component := range components {
		componentNames[i] = fmt.Sprintf("%s, version: %s", component.DisplayName, component.Component.Version)
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

		return components[i]
	}
}
