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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
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

// ComponentsCmd represents the mesheryctl exp components command
var ComponentsCmd = &cobra.Command{
	Use:   "components",
	Short: "View list of components and detail of components",
	Long:  "View list of components and detailed information of a specific component",
	Example: `
// To view the number of components present in Meshery
mesheryctl exp components --count
// To view list of components
mesheryctl exp components list
// To view a specific component
mesheryctl exp components view [component-name]
// To search for a specific component
mesheryctl exp components search [component-name]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			_ = cmd.Help()
		}

		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		err = cmd.Usage()
		if err != nil {
			return err
		}

		// Run count functionality only if count flag is set
		if countFlag {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return err
			}

			baseUrl := mctlCfg.GetBaseMesheryURL()

			// Since we are not searching for particular component, we don't need to pass any search parameter
			url := fmt.Sprintf("%s/api/meshmodels/components?pagesize=all", baseUrl)

			req, err := utils.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				utils.Log.Error(err)
				return err
			}

			resp, err := utils.MakeRequest(req)
			if err != nil {
				utils.Log.Error(err)
				return err
			}

			// defers the closing of the response body after its use, ensuring that the resources are properly released.
			defer resp.Body.Close()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				utils.Log.Error(err)
				return err
			}

			componentsResponse := &models.MeshmodelComponentsAPIResponse{}
			err = json.Unmarshal(data, componentsResponse)
			if err != nil {
				utils.Log.Error(err)
				return err
			}

			fmt.Println("Total components present: ", componentsResponse.Count)
			return nil
		}

		return nil
	},
}

// selectComponentPrompt lets user to select a model if models are more than one
func selectComponentPrompt(components []v1alpha1.ComponentDefinition) v1alpha1.ComponentDefinition {
	componentNames := []string{}
	componentArray := []v1alpha1.ComponentDefinition{}

	componentArray = append(componentArray, components...)

	for _, component := range componentArray {
		componentName := fmt.Sprintf("%s, version: %s", component.DisplayName, component.TypeMeta.APIVersion)
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

func outputComponentJson(component v1alpha1.ComponentDefinition) error {
	if err := prettifyComponentJson(component); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(component, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}

// prettifyJson takes a v1alpha1.Model struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyComponentJson(component v1alpha1.ComponentDefinition) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(component)
}

// Common function to get the min of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	ComponentsCmd.AddCommand(availableSubcommands...)
	ComponentsCmd.Flags().BoolVarP(&countFlag, "count", "c", false, "(optional) Get the number of components in total")
}
