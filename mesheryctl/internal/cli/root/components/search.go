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

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var usageErrorMessage = "Usage: mesheryctl exp component search [query-text]\nRun 'mesheryctl component search --help' to see detailed help message"

type componentSearchFlags struct {
	PageSize int `json:"pageSize" validate:"gte=1"`
	Page     int `json:"page" validate:"gte=1"`
}

var componentSearchCmdFlags componentSearchFlags

// represents the mesheryctl component search [query-text] subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered components",
	Long: `Search components registered in Meshery Server based on kind
Find more information at: https://docs.meshery.io/reference/mesheryctl/component/search`,
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
		searchValue.Set("search", componentName)

		// rows, componentsCount := generateComponentDataToDisplay(componentsResponse)
		searchUrlPath := fmt.Sprintf("%s?%s", componentApiPath, searchValue.Encode())

		componentData := display.DisplayDataAsync{
			UrlPath:          searchUrlPath,
			DataType:         "component",
			Header:           []string{"ID", "Name", "Model", "Version"},
			Page:             componentSearchCmdFlags.Page,
			PageSize:         componentSearchCmdFlags.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: false,
		}

		return display.ListAsyncPagination(componentData, generateComponentDataToDisplay)
	},
}

func init() {
	searchComponentsCmd.Flags().IntVarP(&componentSearchCmdFlags.PageSize, "pagesize", "", 10, "(optional) List next set of components with --pagesize (default = 10)")
	searchComponentsCmd.Flags().IntVarP(&componentSearchCmdFlags.Page, "page", "", 1, "(optional) List next set of components with --page (default = 1)")
}
