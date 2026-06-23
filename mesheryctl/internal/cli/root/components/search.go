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
	"github.com/spf13/pflag"
)

type cmdComponentSearchFlags struct {
	Page         int    `json:"page"          validate:"omitempty,gte=1"`
	PageSize     int    `json:"page-size"     validate:"omitempty,gte=1"`
	Model        string `json:"model"         validate:"omitempty"`
	OutputFormat string `json:"output-format" validate:"omitempty,oneof=table json yaml"`
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

// Search for components within a specific model
mesheryctl component search --model [model-name]

// Combine search query with model filter
mesheryctl component search [query-text] --model [model-name]

// Search list of components of specified page [int]
mesheryctl component search [query-text] [--page 1]

// Output results in JSON format
mesheryctl component search [query-text] -o json

// Output results in YAML format
mesheryctl component search [query-text] --model [model-name] -o yaml
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &componentSearchFlags)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("at most one argument (search query) can be provided\n\n%s", searchUsageMsg))
		}
		if len(args) == 0 && !cmd.Flags().Changed("model") {
			return utils.ErrInvalidArgument(fmt.Errorf("provide a search query or use --model flag\n\n%s", searchUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		urlPath := buildComponentSearchUrl(args)

		// Handle structured output formats (json, yaml)
		if componentSearchFlags.OutputFormat != "table" {
			return handleStructuredOutput(urlPath)
		}

		// Default table output
		modelData := display.DisplayDataAsync{
			UrlPath:  urlPath,
			DataType: "component",
			Header:   []string{"ID", "Name", "Model", "Version"},
			Page:     componentSearchFlags.Page,
			PageSize: componentSearchFlags.PageSize,
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
	searchComponentsCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	searchComponentsCmd.Flags().IntVarP(&componentSearchFlags.Page, "page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	searchComponentsCmd.Flags().IntVarP(&componentSearchFlags.PageSize, "pagesize", "s", 10, "(optional) List next set of components with --pagesize (default = 10)")
	searchComponentsCmd.Flags().StringVarP(&componentSearchFlags.Model, "model", "m", "", "(optional) Filter components by model name")
	searchComponentsCmd.Flags().StringVarP(&componentSearchFlags.OutputFormat, "output-format", "o", "table", "(optional) Output format [table|json|yaml] (default = table)")
}

// buildComponentSearchUrl constructs the API URL for the component search endpoint.
func buildComponentSearchUrl(args []string) string {
	base := componentApiPath

	if componentSearchFlags.Model != "" {
		base = fmt.Sprintf("api/meshmodels/models/%s/components",
			url.PathEscape(componentSearchFlags.Model))
	}

	params := url.Values{}
	if len(args) > 0 && args[0] != "" {
		params.Set("search", args[0])
	}

	if len(params) > 0 {
		return base + "?" + params.Encode()
	}
	return base
}

// handleStructuredOutput fetches all components and renders them as JSON or YAML.
func handleStructuredOutput(urlPath string) error {
	if err := display.ValidateOutputFormat(componentSearchFlags.OutputFormat); err != nil {
		return err
	}

	// Fetch all results by requesting pagesize=all
	fetchUrl := urlPath
	if strings.Contains(fetchUrl, "?") {
		fetchUrl += "&pagesize=all"
	} else {
		fetchUrl += "?pagesize=all"
	}

	data, err := api.Fetch[models.MeshmodelComponentsAPIResponse](fetchUrl)
	if err != nil {
		return err
	}

	outputFormatterFactory := display.OutputFormatterFactory[models.MeshmodelComponentsAPIResponse]{}
	outputFormatter, err := outputFormatterFactory.New(componentSearchFlags.OutputFormat, *data)
	if err != nil {
		return err
	}

	return outputFormatter.Display()
}
