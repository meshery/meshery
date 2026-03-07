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
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

type relationshipListFlag struct {
	Count    bool
	Page     int
	PageSize int
}

var relationshipListFlags relationshipListFlag

// represents the mesheryctl exp relationships list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered relationships",
	Long:  "List all relationships registered in Meshery Server",
	Example: `
// List of relationships
mesheryctl exp relationship list

// List of relationships for a specified page
mesheryctl relationship list --page [page-number]

// Display number of available relationships in Meshery
mesheryctl relationship list --count
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(relationshipListFlags)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		dataToDisplay := display.DisplayDataAsync{
			UrlPath:          relationshipApiPath,
			DataType:         "relationship",
			Header:           []string{"kind", "API Version", "Model name", "Sub Type", "Evaluation Policy"},
			Page:             relationshipListFlags.Page,
			PageSize:         relationshipListFlags.PageSize,
			DisplayCountOnly: relationshipListFlags.Count,
			IsPage:           cmd.Flags().Changed("page"),
		}
		return display.ListAsyncPagination(dataToDisplay, generateRelationshipDataToDisplay)

	},
}

func init() {
	// Add the new exp relationship commands to the listRelationshipsCmd
	listCmd.Flags().IntVarP(&relationshipListFlags.Page, "page", "p", 1, "(optional) List next set of relationships with --page (default = 1)")
listCmd.Flags().IntVar(&relationshipListFlags.PageSize, "pagesize", 10, "(optional) List next set of relationships with --pagesize (default = 10)")
	listCmd.Flags().BoolVarP(&relationshipListFlags.Count, "count", "c", false, "(optional) Display count only")
}
