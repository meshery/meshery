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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pageNo   int
	pageSize int
	name     string
)

type orgsStruct struct {
	Orgs []models.Organization `json:"organizations"`
}

var OrgCmd = &cobra.Command{
	Use:   "org",
	Short: "List registered orgs",
	Long: `Print all registered organizations with their id, name and date of creation. Organization will be returned based on provider you logged in.
	
	Documentation for organizations can be found at 
	https://docs.layer5.io/cloud/identity/organizations/
	`,
	Example: `
	// list all organizations
	mesheryctl exp org
	// list organizations (using flags)
	mesheryctl exp org --page [page_no] --size [size]
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(utils.ErrLoadConfig(err))
			return nil

		}
		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/identity/orgs?page=%d&pagesize=%d", baseUrl, pageNo, pageSize)

		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer resp.Body.Close()
		JsonData, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var orgs orgsStruct
		err = json.Unmarshal(JsonData, &orgs)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		var orgsData [][]string
		columnNames := []string{"NAME", "ID", "CREATED-AT"}

		for _, org := range orgs.Orgs {
			orgsData = append(orgsData, []string{org.Name, org.ID.String(), org.CreatedAt.String()})
		}

		utils.PrintToTable(columnNames, orgsData)

		return nil
	},
}

func init() {
	OrgCmd.Flags().IntVarP(&pageNo, "page", "p", 0, "page number to fetch")
	OrgCmd.Flags().IntVarP(&pageSize, "page_size", "s", 10, "page size")
	OrgCmd.Flags().StringVarP(&name, "grab", "g", "", "name of organization to search")
}
