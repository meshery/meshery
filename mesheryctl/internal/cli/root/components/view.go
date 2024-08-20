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
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// represents the mesheryctl components view [component-name] subcommand.
var viewComponentCmd = &cobra.Command{
	Use:   "view",
	Short: "view registered components",
	Long:  "view a component registered in Meshery Server",
	Example: `
// View details of a specific component
mesheryctl components view [component-name]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp component view [component-name]\nRun 'mesheryctl exp component view --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("component name isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		componentDefinition := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/components?search=%s&pagesize=all", baseUrl, componentDefinition)
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

		componentResponse := &models.MeshmodelComponentsAPIResponse{}
		err = json.Unmarshal(data, componentResponse)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var selectedComponent component.ComponentDefinition

		if componentResponse.Count == 0 {
			fmt.Println("No component(s) found for the given name: ", componentDefinition)
			return nil
		} else if componentResponse.Count == 1 {
			selectedComponent = componentResponse.Components[0] // Update the type of selectedModel
		} else {
			selectedComponent = selectComponentPrompt(componentResponse.Components)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormatFlag = strings.ToLower(outFormatFlag)

		if outFormatFlag != "json" && outFormatFlag != "yaml" {
			return errors.New("output-format choice is invalid or not provided, use [json|yaml]")
		}
		// Get the home directory of the user to save the output file
		homeDir, _ := os.UserHomeDir()
		componentString := strings.ReplaceAll(fmt.Sprintf("%v", selectedComponent.DisplayName), " ", "_")

		if outFormatFlag == "yaml" {
			if output, err = yaml.Marshal(selectedComponent); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			if saveFlag {
				fmt.Println("Saving output as YAML file")
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".yaml", output, 0666)
				if err != nil {
					return errors.Wrap(err, "failed to save output as YAML file")
				}
				fmt.Println("Output saved as YAML file in ~/.meshery/component_" + componentString + ".yaml")
			} else {
				fmt.Print(string(output))
			}
		} else if outFormatFlag == "json" {
			if saveFlag {
				fmt.Println("Saving output as JSON file")
				output, err = json.MarshalIndent(selectedComponent, "", "  ")
				if err != nil {
					return errors.Wrap(err, "failed to format output in JSON")
				}
				err = os.WriteFile(homeDir+"/.meshery/component_"+componentString+".json", output, 0666)
				if err != nil {
					return errors.Wrap(err, "failed to save output as JSON file")
				}
				fmt.Println("Output saved as JSON file in ~/.meshery/component_" + componentString + ".json")
				return nil
			}
			return OutputJson(selectedComponent)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	viewComponentCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewComponentCmd.Flags().BoolVarP(&saveFlag, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
