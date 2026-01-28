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

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

type environmentListFlags struct {
	count    bool
	orgID    string
	page     int
	pagesize int
}

var environmentListFlagsProvided environmentListFlags

var listEnvironmentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered environments",
	Long: `List detailed information of all registered environments
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/list`,
	Example: `
// List all registered environment
mesheryctl environment list --orgID [orgID]

// List count of all registered environment
mesheryctl environment list --orgID [orgID] --count

// List all registered environment at a specific page
mesheryctl environment list --orgID [orgID] --page [page]

// List all registered environment with a specific page size
mesheryctl environment list --orgID [orgID] --pagesize [pagesize]
`,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		if environmentListFlagsProvided.orgID == "" {
			const errMsg = "[ orgID ] isn't specified\n\nUsage: mesheryctl environment list --orgID [orgID]\nRun 'mesheryctl environment list --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(environmentListFlagsProvided.orgID) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid orgID: %s", environmentListFlagsProvided.orgID))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		data := display.DisplayDataAsync{
			UrlPath:          fmt.Sprintf("%s?orgID=%s", environmentApiPath, environmentListFlagsProvided.orgID),
			DataType:         "environments",
			Header:           []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"},
			Page:             environmentListFlagsProvided.page,
			PageSize:         environmentListFlagsProvided.pagesize,
			DisplayCountOnly: environmentListFlagsProvided.count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(data, processEnvironmentData)
	},
}

func processEnvironmentData(environmentResponse *environment.EnvironmentPage) ([][]string, int64) {
	rows := [][]string{}
	for _, environment := range environmentResponse.Environments {
		row := []string{environment.ID.String(), environment.Name, environment.OrganizationID.String(), environment.Description, environment.CreatedAt.String(), environment.UpdatedAt.String()}
		rows = append(rows, row)
	}
	return rows, int64(environmentResponse.TotalCount)
}

func init() {
	listEnvironmentCmd.Flags().BoolVarP(&environmentListFlagsProvided.count, "count", "c", false, "(optional) Display count only")
	listEnvironmentCmd.Flags().StringVarP(&environmentListFlagsProvided.orgID, "orgID", "", "", "Organization ID")
	listEnvironmentCmd.Flags().IntVarP(&environmentListFlagsProvided.page, "page", "", 1, "(optional) Page number of paginated results")
	listEnvironmentCmd.Flags().IntVarP(&environmentListFlagsProvided.pagesize, "pagesize", "", 10, "(optional) Number of results per page")
}
