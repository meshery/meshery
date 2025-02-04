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

package orgs

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	maxRowsPerPage int = 10
)

var ListOrgCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered orgs",
	Long: `Print all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
	
	Documentation for organizations can be found at 
	https://docs.layer5.io/cloud/identity/organizations/
	`,
	Example: `
	// list all organizations
	mesheryctl exp org list
	// list organizations (using flags)
	mesheryctl exp org --page [page_no] --size [size]
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		orgs, err := getAllOrgs()
		if err != nil {
			utils.Log.Error(err)
		}

		var orgsData [][]string
		columnNames := []string{"NAME", "ID", "CREATED-AT"}

		for _, org := range orgs.Organizations {
			orgsData = append(orgsData, []string{org.Name, org.ID.String(), fmt.Sprintf("%v/%v/%v", org.CreatedAt.Year(), org.CreatedAt.Month(), org.CreatedAt.Day())})
		}
		if len(orgsData) == 0 {
			utils.Log.Info("No organizations found")
			return nil
		}
		err = utils.HandlePagination(maxRowsPerPage, "organizations", orgsData, columnNames)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		return nil
	},
}
