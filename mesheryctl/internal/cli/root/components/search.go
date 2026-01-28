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
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

var usageErrorMessage = "Usage: mesheryctl exp component search [query-text]\nRun 'mesheryctl exp component search --help' to see detailed help message"

// represents the mesheryctl component search [query-text] subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered components",
	Long: `Search components registered in Meshery Server based on kind
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/search`,
	Example: `
// Search for components using a query
mesheryctl component search [query-text]
	`,
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("[search term] isn't specified. Please enter component name to search\n\n%v", usageErrorMessage))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		componentName := strings.Join(args, " ")
		searchValue := url.Values{}
		searchValue.Add("search", componentName)
		searchValue.Add("pagesize", "all")

		componentsResponse, err := api.Fetch[models.MeshmodelComponentsAPIResponse](fmt.Sprintf("%s?%s", componentApiPath, searchValue.Encode()))

		if err != nil {
			return err
		}

		header := []string{"Name", "Model", "kind", "Version"}

		rows, componentsCount := generateComponentDataToDisplay(componentsResponse)

		dataToDisplay := display.DisplayedData{
			DataType:         "components",
			Header:           header,
			Rows:             rows,
			Count:            componentsCount,
			DisplayCountOnly: false,
			IsPage:           false,
		}

		err = display.List(dataToDisplay)
		if err != nil {
			return err
		}

		return nil
	},
}
