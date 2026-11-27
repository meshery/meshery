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
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var listUsageMsg = "Usage: mesheryctl environment list --orgId [orgId]\nRun 'mesheryctl environment list --help' to see detailed help message"

type environmentListFlags struct {
	Count    bool   `json:"count" validate:"boolean"`
	OrgID    string `json:"orgId"`
	Page     int    `json:"page" validate:"gte=1"`
	PageSize int    `json:"pagesize" validate:"gte=1"`
}

var environmentListFlagsProvided environmentListFlags

var listEnvironmentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered environments",
	Long: `List detailed information of all registered environments
Find more information at: https://docs.meshery.io/reference/mesheryctl/environment/list`,
	Example: `
// List all registered environments
mesheryctl environment list --orgId [orgId]

// List count of all registered environments
mesheryctl environment list --orgId [orgId] --count

// List all registered environments at a specific page
mesheryctl environment list --orgId [orgId] --page [page]

// List all registered environments with a specific page size
mesheryctl environment list --orgId [orgId] --pagesize [pagesize]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments specified\n\n%s", listUsageMsg))
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if environmentListFlagsProvided.OrgID == "" {
			errMsg := "[ orgId ] isn't specified\n\n" + listUsageMsg
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(environmentListFlagsProvided.OrgID) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid orgId: %s", environmentListFlagsProvided.OrgID))
		}

		return mesheryctlflags.ValidateCmdFlags(cmd, &environmentListFlagsProvided)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		data := display.DisplayDataAsync{
			UrlPath:          fmt.Sprintf("%s?orgId=%s", environmentApiPath, environmentListFlagsProvided.OrgID),
			DataType:         "environments",
			Header:           []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"},
			Page:             environmentListFlagsProvided.Page,
			PageSize:         environmentListFlagsProvided.PageSize,
			DisplayCountOnly: environmentListFlagsProvided.Count,
			IsPage:           cmd.Flags().Changed("page"),
		}

		return display.ListAsyncPagination(data, processEnvironmentData)
	},
}

func processEnvironmentData(environmentResponse *environment.EnvironmentPage) ([][]string, int64) {
	if environmentResponse == nil {
		return [][]string{}, 0
	}

	rows := [][]string{}
	for _, environment := range environmentResponse.Environments {
		description := "N/A"
		if strings.TrimSpace(environment.Description) != "" {
			description = environment.Description
		}
		row := []string{environment.ID.String(), environment.Name, environment.OrganizationID.String(), description, environment.CreatedAt.String(), environment.UpdatedAt.String()}
		rows = append(rows, row)
	}
	return rows, int64(environmentResponse.TotalCount)
}

func init() {
	listEnvironmentCmd.Flags().BoolVarP(&environmentListFlagsProvided.Count, "count", "c", false, "(optional) Display count only")
	listEnvironmentCmd.Flags().StringVarP(&environmentListFlagsProvided.OrgID, "orgId", "", "", "Organization ID")
	listEnvironmentCmd.Flags().IntVarP(&environmentListFlagsProvided.Page, "page", "", 1, "(optional) Page number of paginated results")
	listEnvironmentCmd.Flags().IntVarP(&environmentListFlagsProvided.PageSize, "pagesize", "", 10, "(optional) Number of results per page")
}
