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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/spf13/cobra"
)

var listOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered organizations",
	Long: `List all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
Documentation for organizations can be found at https://docs.meshery.io/reference/mesheryctl/exp/organizations/list
	`,
	Example: `
// list all organizations
mesheryctl exp organizations list

// list organizations for a specified page
mesheryctl exp organizations list --page [page-number]

// list organizations for a specified page
mesheryctl exp organizations list --count
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		orgs, err := getAllOrganizations()
		if err != nil {
			return err
		}

		var rows [][]string
		header := []string{"NAME", "ID", "CREATED-AT"}

		for _, org := range orgs.Organizations {
			rows = append(rows, []string{org.Name, org.ID.String(), fmt.Sprintf("%v/%v/%v", org.CreatedAt.Year(), org.CreatedAt.Month(), org.CreatedAt.Day())})
		}

		count, _ := cmd.Flags().GetBool("count")
		dataToDisplay := display.DisplayedData{
			DataType:         "organizations",
			Header:           header,
			Rows:             rows,
			Count:            int64(orgs.TotalCount),
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
	listOrgCmd.Flags().IntP("page", "p", 0, "page number")
	listOrgCmd.Flags().BoolP("count", "", false, "total number of registered orgs")
}
