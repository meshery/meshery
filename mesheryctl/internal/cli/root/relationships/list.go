// Copyright Meshery Authors
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

package relationships

import (
	"fmt"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cmdRelationshipListFlags struct {
	Count    bool `json:"count"    validate:"boolean"`
	Page     int  `json:"page"     validate:"omitempty,gte=1"`
	PageSize int  `json:"pagesize" validate:"omitempty,gte=1"`
}

var relationshipListFlags cmdRelationshipListFlags

// represents the mesheryctl exp relationships list command.
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered relationships",
	Long:  "List all relationships registered in Meshery Server",
	Example: `
// List all relationships
mesheryctl relationship list

// List relationships for a specified page
mesheryctl relationship list --page [page-number]

// List relationships with a custom page size
mesheryctl relationship list --pagesize [page-size]

// Display the total number of available relationships in Meshery
mesheryctl relationship list --count
`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &relationshipListFlags)
	},

	Args: func(_ *cobra.Command, args []string) error {
		if len(args) != 0 {
			errMsg := "Usage: mesheryctl exp relationship list\nRun 'mesheryctl exp relationship list --help' to see detailed help message"
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		dataToDisplay := display.DisplayDataAsync{
			UrlPath:          relationshipApiPath,
			DataType:         "relationship",
			Header:           []string{"ID", "Kind", "API Version", "Model Name", "Sub Type", "Type"},
			Page:             relationshipListFlags.Page,
			PageSize:         relationshipListFlags.PageSize,
			DisplayCountOnly: relationshipListFlags.Count,
			IsPage:           cmd.Flags().Changed("page"),
		}
		return display.ListAsyncPagination(dataToDisplay, generateRelationshipDataToDisplay)
	},
}

func init() {
	listCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	listCmd.Flags().IntVarP(&relationshipListFlags.Page, "page", "p", 1, "(optional) List next set of relationships with --page (default = 1)")
	listCmd.Flags().IntVar(&relationshipListFlags.PageSize, "pagesize", 10, "(optional) Number of results per page (default = 10)")
	listCmd.Flags().BoolVarP(&relationshipListFlags.Count, "count", "c", false, "(optional) Display the total count of relationships only")
}
