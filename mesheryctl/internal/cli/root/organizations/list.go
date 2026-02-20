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
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	Long: `List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
Find more information at: https://docs.meshery.io/reference/mesheryctl/exp/organizations/list
	`,
	Example: `
// list all organizations
mesheryctl exp organization list

// list organizations for a specified page
mesheryctl exp organization list --page [page-number]

// list organizations for a specified page
mesheryctl exp organization list --count
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		page, _ := cmd.Flags().GetInt("page")
		pagesize, _ := cmd.Flags().GetInt("pagesize")
		count, _ := cmd.Flags().GetBool("count")
		data := display.DisplayDataAsync{ 
			UrlPath: organizationsApiPath,
			Page: page,
			PageSize: pagesize,
			DataType: "organizations",
			Header: []string{"ID", "NAME", "CREATED-AT"},
			DisplayCountOnly: count,
			IsPage: cmd.Flags().Changed("page"),
		}
		return display.ListAsyncPagination(data, processOrgData)
	},
}



func processOrgData(orgs *models.OrganizationsPage) ([][]string,int64) {
	var rows [][]string
	for _,org := range orgs.Organizations {
		rows = append(rows, []string{
			org.Id.String(),
			org.Name,
			fmt.Sprintf("%v/%v/%v", org.CreatedAt.Year(), org.CreatedAt.Month(), org.CreatedAt.Day()),
		})

	}
	return rows, int64(orgs.TotalCount)
}


func init() {
	listOrgCmd.Flags().IntP("page", "p", 1, "(optional) Page number of paginated results")
	listOrgCmd.Flags().IntP("pagesize", "s", 10, "(optional) Number of organizations per page")
	listOrgCmd.Flags().BoolP("count", "", false, "total number of registered orgs")



}
