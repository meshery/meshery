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

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

type componentSearchFlag struct {
	Page     int
	PageSize int
}

var cmdComponentSearchFlag componentSearchFlag

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
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(cmdComponentSearchFlag)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("%v\n\n%v", errNoArg, searchUsageMessage))
		}

		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%v\n\n%v", errMultiArg, searchUsageMessage))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		searchValue := url.Values{}
		searchValue.Add("search", args[0])

		modelData := display.DisplayDataAsync{
			UrlPath:  fmt.Sprintf("%s?%s", componentApiPath, searchValue.Encode()),
			DataType: "component",
			Header:   []string{"Name", "Model", "Version"},
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
