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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the mesheryctl components search [query-text] subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "search registered components",
	Long:  "search components registered in Meshery Server based on kind",
	Example: `
// Search for components using a query
mesheryctl components search [query-text]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp component search [query-text]\nRun 'mesheryctl exp component search --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("search term is missing. Please enter component name to search\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		queryText := args[0]
		url := fmt.Sprintf("%s/api/meshmodels/components?search=%s&pagesize=all", baseUrl, queryText)

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
			fmt.Println("No component(s) found with the search term")
			return nil
		} else {
			// Print the result in tabular format
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}
