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
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/spf13/cobra"
)

// represents the mesheryctl component list command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long: `List all components registered in Meshery Server
Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/component/list`,
	Example: `
// View list of components
mesheryctl component list

// View list of components with specified page number (25 components per page)
mesheryctl component list --page [page-number]

// Display the number of components present in Meshery
mesheryctl component list --count
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		pageSize, _ := cmd.Flags().GetInt("pagesize")
		modelData := display.DisplayDataAsync{
			UrlPath:          componentApiPath,
			DataType:         "component",
			Header:           []string{"Model", "Category", "Version"},
			Page:             page,
			PageSize:         pageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmd.Flags().Changed("count"),
		}

		return display.ListAsyncPagination(modelData, generateComponentDataToDisplay)
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	listComponentCmd.Flags().IntP("page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	listComponentCmd.Flags().IntP("pagesize", "s", 0, "(optional) List next set of components with --pagesize (default = 0)")
	listComponentCmd.Flags().BoolP("count", "c", false, "(optional) Display count only")
}
