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

package organizations

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

type orgListFlags struct {
	Count    bool `json:"count" validate:"boolean"`
	Page     int  `json:"page" validate:"gte=1"`
	PageSize int  `json:"pagesize" validate:"gte=1"`
}

var orgListFlagsProvided orgListFlags

var listUsageMsg = "Usage: mesheryctl organization list\nRun 'mesheryctl organization list --help' to see detailed help message"

var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	Long: `List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
Find more information at: https://docs.meshery.io/reference/mesheryctl/organization/list`,
	Example: `
// list all organizations
mesheryctl organization list

// list organizations for a specified page
mesheryctl organization list --page [page-number]

// Display number of available organizations
mesheryctl organization list --count
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &orgListFlagsProvided)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		data := display.DisplayDataAsync{
			UrlPath:          organizationsApiPath,
			Page:             orgListFlagsProvided.Page,
			PageSize:         orgListFlagsProvided.PageSize,
			DataType:         "organizations",
			Header:           []string{"ID", "NAME", "CREATED-AT"},
			DisplayCountOnly: orgListFlagsProvided.Count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(data, processOrgData)
	},
}

func processOrgData(orgs *models.OrganizationsPage) ([][]string, int64) {
	var rows [][]string
	for _, org := range orgs.Organizations {
		rows = append(rows, []string{
			org.ID.String(),
			org.Name,
			org.CreatedAt.Format("01-02-2006"),
		})

	}
	return rows, int64(orgs.TotalCount)
}

func init() {
	listOrgCmd.Flags().IntVarP(&orgListFlagsProvided.Page, "page", "p", 1, "(optional) Page number of paginated results")
	listOrgCmd.Flags().IntVarP(&orgListFlagsProvided.PageSize, "pagesize", "s", 10, "(optional) Number of organizations per page")
	listOrgCmd.Flags().BoolVarP(&orgListFlagsProvided.Count, "count", "", false, "total number of registered orgs")
}
