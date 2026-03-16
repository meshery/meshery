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
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

type cmdComponentSearchFlags struct {
	Page         int    `json:"page" validate:"omitempty,gte=1"`
	PageSize     int    `json:"page-size" validate:"omitempty,gte=1"`
	Model        string `json:"model" validate:"omitempty"`
	OutputFormat string `json:"output-format" validate:"omitempty"`
}

var componentSearchFlags cmdComponentSearchFlags

// represents the mesheryctl component search [query-text] subcommand.
var searchComponentsCmd = &cobra.Command{
	Use:   "search",
	Short: "Search registered components",
	Long: `Search components registered in Meshery Server based on kind
Find more information at: https://docs.meshery.io/reference/mesheryctl/component/search`,
	Example: `
// Search for components using a query
mesheryctl component search [query-text]

// Search for multi-word component names (must be quoted)
mesheryctl component search "Component name"

// Search list of components of specified page [int]
mesheryctl component search [query-text] [--page 1]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &componentSearchFlags)
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 && componentSearchFlags.Model == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide a query text or use --model flag\n\n%v", searchUsageMsg))
		}
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%v", searchUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		searchValue := url.Values{}
		if len(args) > 0 {
			searchValue.Add("search", args[0])
		}

		path := componentApiPath
		if componentSearchFlags.Model != "" {
			path = fmt.Sprintf("api/meshmodels/models/%s/components", url.PathEscape(componentSearchFlags.Model))
		}

		// If output format is specified, we fetch the raw data and display it
		if componentSearchFlags.OutputFormat != "" {
			queryPath := fmt.Sprintf("%s?%s", path, searchValue.Encode())
			// Add pagination params for one-shot fetch if not interactive
			if !strings.Contains(queryPath, "page=") {
				queryPath = fmt.Sprintf("%s&page=%d&pagesize=%d", queryPath, componentSearchFlags.Page-1, componentSearchFlags.PageSize)
			}

			response, err := api.Fetch[models.MeshmodelComponentsAPIResponse](queryPath)
			if err != nil {
				return err
			}

			if componentSearchFlags.OutputFormat == "table" {
				// Fallback to table if explicitly requested, but non-interactive
				rows, count := generateComponentDataToDisplay(response)
				utils.DisplayCount("component", count)
				utils.PrintToTable([]string{"ID", "Name", "Model", "Version"}, rows, nil)
				return nil
			}

			outputFormatterFactory := display.OutputFormatterFactory[models.MeshmodelComponentsAPIResponse]{}
			outputFormatter, err := outputFormatterFactory.New(strings.ToLower(componentSearchFlags.OutputFormat), *response)
			if err != nil {
				return err
			}

			return outputFormatter.Display()
		}

		modelData := display.DisplayDataAsync{
			UrlPath:  fmt.Sprintf("%s?%s", path, searchValue.Encode()),
			DataType: "component",
			Header:   []string{"ID", "Name", "Model", "Version"},
			Page:     componentSearchFlags.Page,
			PageSize: componentSearchFlags.PageSize,
			IsPage:   cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(modelData, generateComponentDataToDisplay)
	},
}

func init() {
	searchComponentsCmd.Flags().IntVarP(&componentSearchFlags.Page, "page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	searchComponentsCmd.Flags().IntVarP(&componentSearchFlags.PageSize, "pagesize", "s", 10, "(optional) List next set of components with --pagesize (default = 10)")
	searchComponentsCmd.Flags().StringVarP(&componentSearchFlags.Model, "model", "m", "", "(optional) Search components of a particular model name")
	searchComponentsCmd.Flags().StringVarP(&componentSearchFlags.OutputFormat, "output-format", "o", "", "(optional) format to display in [json|yaml]")
}
