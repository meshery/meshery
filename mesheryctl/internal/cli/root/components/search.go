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

var (
	usageErrorMessage = "Usage: mesheryctl exp component search [query-text]\nRun 'mesheryctl exp component search --help' to see detailed help message"
)

type componentSearchFlag struct {
	Count    bool
	Page     int
	PageSize int
}

var cmdComponentSearchFlag componentSearchFlag

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

		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("multi-word filter names must be enclosed in double quotes\n\n%v", usageErrorMessage))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		componentName := strings.Join(args, "_")
		searchValue := url.Values{}
		searchValue.Add("search", componentName)

		modelData := display.DisplayDataAsync{
			UrlPath:  fmt.Sprintf("%s?%s", componentApiPath, searchValue.Encode()),
			DataType: "component",
			Header:   []string{"Name", "Model", "Category", "Version"},
			Page:     cmdComponentSearchFlag.Page,
			PageSize: cmdComponentSearchFlag.PageSize,
			IsPage:   cmd.Flags().Changed("page"),
		}

		err := display.ListAsyncPagination(modelData, generateComponentDataToDisplay)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	searchComponentsCmd.Flags().IntVarP(&cmdComponentSearchFlag.Page, "page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	searchComponentsCmd.Flags().IntVarP(&cmdComponentSearchFlag.PageSize, "pagesize", "s", 10, "(optional) List next set of components with --pagesize (default = 10)")
}
