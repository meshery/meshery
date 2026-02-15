// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package workspaces

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/workspace"

	"github.com/spf13/cobra"
)

type workspaceListFlags struct {
	OrganizationID string
	Count          bool
	Page           int
	PageSize       int
}

var workspaceListFlagsProvided workspaceListFlags
var listUsageErrorMessage = "Usage: mesheryctl exp workspace list --orgId [Organization ID]\nRun 'mesheryctl exp workspace list --help' to see detailed help message"

var listWorkspaceCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered workspaces",
	Long: `List name of all registered workspaces
Documentation for models can be found at https://docs.meshery.io/reference/mesheryctl/exp/workspace/list`,
	Example: `
// List of workspace under a specific organization
mesheryctl exp workspace list --orgId [orgId]

// List of workspace under a specific organization for a specified page
mesheryctl exp workspace list --orgId [orgId] --page [page-number]

// Display number of available  workspace under a specific organization
mesheryctl exp workspace list --orgId [orgId] --count
`,

	Args: func(cmd *cobra.Command, args []string) error {
		if workspaceListFlagsProvided.OrganizationID == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("[ Organization ID ] isn't specified\n\n%s", listUsageErrorMessage))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID := workspaceListFlagsProvided.OrganizationID

		modelData := display.DisplayDataAsync{
			UrlPath:          workspacesApiPath + "?orgID=" + orgID,
			DataType:         "workspace",
			Header:           []string{"ID", "Name", "Description", "Created At", "Updated At"},
			Page:             workspaceListFlagsProvided.Page,
			PageSize:         workspaceListFlagsProvided.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: cmd.Flags().Changed("count"),
		}

		return display.ListAsyncPagination(modelData, generateModelDataToDisplay)

	},
}

func generateModelDataToDisplay(workspacePage *workspace.WorkspacePage) ([][]string, int64) {
	rows := [][]string{}

	for _, workspace := range workspacePage.Workspaces {
		rows = append(rows, []string{workspace.ID.String(), workspace.Name, workspace.Description, workspace.CreatedAt.String(), workspace.UpdatedAt.String()})
	}

	return rows, int64(workspacePage.TotalCount)
}

func init() {
	listWorkspaceCmd.Flags().BoolVarP(&workspaceListFlagsProvided.Count, "count", "c", false, "[optional] count total of registered workspaces")
	listWorkspaceCmd.Flags().StringVarP(&workspaceListFlagsProvided.OrganizationID, "orgId", "o", "", "Organization ID")
	listWorkspaceCmd.Flags().IntVarP(&workspaceListFlagsProvided.Page, "page", "p", 0, "[optional] page number for paginated results")
	listWorkspaceCmd.Flags().IntVarP(&workspaceListFlagsProvided.PageSize, "pagesize", "s", 0, "[optional] number of workspaces per page for paginated results")
}
