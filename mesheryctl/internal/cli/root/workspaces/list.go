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
	"net/url"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/schemas/models/v1beta1/workspace"

	"github.com/spf13/cobra"
)

type cmdWorkspaceListFlags struct {
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"pagesize" validate:"min=1,max=100"`
	Count    bool   `json:"count" validate:"boolean"`
	OrgId    string `json:"orgId" validate:"required,uuid"`
}

var workspaceListFlags cmdWorkspaceListFlags

var listWorkspaceCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered workspaces",
	Long: `List name of all registered workspaces
Find more information at: https://docs.meshery.io/reference/mesheryctl/workspace/list`,
	Example: `
// List of workspace under a specific organization
mesheryctl workspace list --orgId [orgId]

// List of workspace under a specific organization for a specified page
mesheryctl workspace list --orgId [orgId] --page [page-number]

// Display number of available  workspace under a specific organization
mesheryctl workspace list --orgId [orgId] --count
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &workspaceListFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		urlQueryParams := url.Values{}
		urlQueryParams.Add("orgID", workspaceListFlags.OrgId)
		if cmd.Flags().Changed("page") {
			// Adjusting page number to be zero-based index for API compatibility,
			// while keeping the user-facing flag as one-based index for better UX
			urlQueryParams.Add("page", fmt.Sprint(workspaceListFlags.Page-1))
		}
		// API will use default page size if pagesize flag is not provided,
		// so only add pagesize to query params if the flag is explicitly set by the user
		if cmd.Flags().Changed("pagesize") {
			urlQueryParams.Add("pagesize", fmt.Sprint(workspaceListFlags.PageSize))
		}
		urlPath := fmt.Sprintf("%s?%s", workspacesApiPath, urlQueryParams.Encode())

		modelData := display.DisplayDataAsync{
			UrlPath:          urlPath,
			DataType:         "workspace",
			Header:           []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"},
			Page:             workspaceListFlags.Page,
			PageSize:         workspaceListFlags.PageSize,
			IsPage:           cmd.Flags().Changed("page"),
			DisplayCountOnly: workspaceListFlags.Count,
		}

		return display.ListAsyncPagination(modelData, processDataToDisplay)
	},
}

func init() {
	listWorkspaceCmd.Flags().BoolVarP(&workspaceListFlags.Count, "count", "", false, "total number of registered workspaces")
	listWorkspaceCmd.Flags().StringVarP(&workspaceListFlags.OrgId, "orgId", "o", "", "Organization ID")
	listWorkspaceCmd.Flags().IntVarP(&workspaceListFlags.Page, "page", "", 1, "page number for paginated results. (default: 1)")
	listWorkspaceCmd.Flags().IntVarP(&workspaceListFlags.PageSize, "pagesize", "", 10, "number of items to be displayed per page for paginated results. (default: 10, max limit: 100)")
}

func processDataToDisplay(workspaceResponse *workspace.WorkspacePage) ([][]string, int64) {
	rows := [][]string{}

	for _, workspace := range workspaceResponse.Workspaces {
		description := "N/A"
		if workspace.Description != "" {
			description = workspace.Description
		}
		rows = append(rows, []string{workspace.ID.String(), workspace.Name, workspace.OrganizationId.String(), description, workspace.CreatedAt.String(), workspace.UpdatedAt.String()})
	}

	return rows, int64(workspaceResponse.TotalCount)
}
