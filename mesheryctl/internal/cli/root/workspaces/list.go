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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

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
		// Check if the orgID is provided
		orgIdFlag, _ := cmd.Flags().GetString("orgId")
		if orgIdFlag == "" {
			const errorMsg = "[ Organization ID ] isn't specified\n\nUsage: \nmesheryctl  exp workspace list --orgId [Organization ID]\nmesheryctl  exp workspace list --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errorMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, _ := cmd.Flags().GetString("orgId")

		workspaceResponse, err := api.Fetch[models.WorkspacePage](fmt.Sprintf("%s?orgID=%s", workspacesApiPath, orgID))
		if err != nil {
			return err
		}

		header := []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"}
		rows := [][]string{}
		for _, workspace := range workspaceResponse.Workspaces {
			rows = append(rows, []string{workspace.ID.String(), workspace.Name, workspace.OrganizationID.String(), workspace.Description, workspace.CreatedAt.String(), workspace.UpdatedAt.String()})
		}

		count, _ := cmd.Flags().GetBool("count")
		dataToDisplay := display.DisplayedData{
			DataType:         "workspaces",
			Header:           header,
			Rows:             rows,
			Count:            int64(workspaceResponse.TotalCount),
			DisplayCountOnly: count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		err = display.List(dataToDisplay)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	listWorkspaceCmd.Flags().BoolP("count", "", false, "total number of registered workspaces")
	listWorkspaceCmd.Flags().StringP("orgId", "o", "", "Organization ID")
}
