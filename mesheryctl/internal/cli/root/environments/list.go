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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listEnvironmentCmd = &cobra.Command{
	Use:   "list",
	Short: "List registered environments",
	Long: `List name of all registered environments
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/list`,
	Example: `
// List all registered environment
mesheryctl environment list --orgID [orgID]
`,

	Args: func(cmd *cobra.Command, args []string) error {
		// Check if all flag is set
		orgIDFlag, _ := cmd.Flags().GetString("orgID")

		if orgIDFlag == "" {
			const errMsg = "[ orgID ] isn't specified\n\nUsage: mesheryctl environment list --orgID [orgID]\nRun 'mesheryctl environment list --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID, _ := cmd.Flags().GetString("orgID")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()

		url := fmt.Sprintf("%s/api/environments?orgID=%s", baseUrl, orgID)

		environmentResponse, err := api.Fetch[environments.EnvironmentPage](url)

		if err != nil {
			return err
		}

		header := []string{"ID", "Name", "Organization ID", "Description", "Created At", "Updated At"}
		rows := [][]string{}
		for _, environment := range environmentResponse.Environments {
			rows = append(rows, []string{environment.ID.String(), environment.Name, environment.OrganizationID.String(), environment.Description, environment.CreatedAt.String(), environment.UpdatedAt.String()})
		}

		dataToDisplay := display.DisplayedData{
			DataType:         "environments",
			Header:           header,
			Rows:             rows,
			Count:            int64(environmentResponse.TotalCount),
			DisplayCountOnly: false,
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
	listEnvironmentCmd.Flags().StringP("orgID", "", "", "Organization ID")
}
