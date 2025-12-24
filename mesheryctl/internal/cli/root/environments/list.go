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

package environments

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models/environments"
)

var listEnvironmentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered environments",
	Long: `List the names of all registered environments.

Documentation:
https://docs.meshery.io/reference/mesheryctl/environment/list`,
	Example: `
# List all registered environments for an organization
mesheryctl environment list --orgID <org-id>
`,
	Args: func(cmd *cobra.Command, args []string) error {
        orgID, err := cmd.Flags().GetString("orgID")
        if err != nil {
            return err
        }
		if orgID == "" {
			return utils.ErrInvalidArgument(
				errors.New(
					"[ orgID ] is required\n\nUsage: mesheryctl environment list --orgID <org-id>",
				),
			)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, _ := cmd.Flags().GetString("orgID")

		// Validate orgID is a valid UUID
		if _, err := uuid.Parse(orgID); err != nil {
			return utils.ErrInvalidOrgID(err)
		}

		// Fetch environments from API
		environmentResponse, err := api.Fetch[environments.EnvironmentPage](
			fmt.Sprintf("api/environments?orgID=%s", orgID),
		)
		if err != nil {
			return utils.ErrFetchEnvironments(err)
		}

		// Prepare table output
		header := []string{
			"ID",
			"Name",
			"Organization ID",
			"Description",
			"Created At",
			"Updated At",
		}

		rows := [][]string{}
		for _, env := range environmentResponse.Environments {
			rows = append(rows, []string{
				env.ID.String(),
				env.Name,
				env.OrganizationID.String(),
				env.Description,
				env.CreatedAt.String(),
				env.UpdatedAt.String(),
			})
		}

		data := display.DisplayedData{
			DataType:         "environments",
			Header:           header,
			Rows:             rows,
			Count:            int64(environmentResponse.TotalCount),
			DisplayCountOnly: false,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.List(data)
	},
}

func init() {
	listEnvironmentCmd.Flags().String(
		"orgID",
		"",
		"Organization ID (required)",
	)
}
