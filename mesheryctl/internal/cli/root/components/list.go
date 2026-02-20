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

type componentListFlag struct {
	Count    bool
	Page     int
	PageSize int
}

var cmdComponentListFlag componentListFlag

// represents the mesheryctl component list command
var listComponentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered components",
	Long: `List all components registered in Meshery Server
Find more information at: https://docs.meshery.io/reference/mesheryctl/component/list`,
	Example: `
// View list of components
mesheryctl component list

// View list of components with specified page number (10 components per page)
mesheryctl component list --page [page-number]

// View list of components with specified page number with specified number of components per page
mesheryctl component list --page [page-number] --pagesize [page-size]

// Display the number of components present in Meshery
mesheryctl component list --count
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modelData := display.DisplayDataAsync{
			UrlPath:          componentApiPath,
			DataType:         "component",
			Header:           []string{"Name", "Model", "Category", "Version"},
			Page:             cmdComponentListFlag.Page,
			PageSize:         cmdComponentListFlag.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmdComponentListFlag.Count,
		}

		err := display.ListAsyncPagination(modelData, generateComponentDataToDisplay)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	// Add the new components commands to the ComponentsCmd
	listComponentCmd.Flags().IntVarP(&cmdComponentListFlag.Page, "page", "p", 1, "(optional) List next set of components with --page (default = 1)")
	listComponentCmd.Flags().IntVarP(&cmdComponentListFlag.PageSize, "pagesize", "s", 10, "(optional) List next set of components with --pagesize (default = 10)")
	listComponentCmd.Flags().BoolVarP(&cmdComponentListFlag.Count, "count", "c", false, "(optional) Display count only")
}
