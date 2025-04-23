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

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environments",
	Long: `Create a new environments by providing the name and description of the environment
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/create`,
	Example: `
// Create a new environment
mesheryctl environment create --orgID [orgID] --name [name] --description [description]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "[ Organization ID | Name | Description ] aren't specified\n    Usage: mesheryctl environment create --orgID [orgID] --name [name] --description [description]\n    Run 'mesheryctl environment create --help' to see detailed help message"

		// Check if all three flags are set
		orgIDFlag, _ := cmd.Flags().GetString("orgID")
		nameFlag, _ := cmd.Flags().GetString("name")
		descriptionFlag, _ := cmd.Flags().GetString("description")

		if orgIDFlag == "" || nameFlag == "" || descriptionFlag == "" {
			return errors.New(display.GetErrorLog(utils.ErrInvalidArgument(errors.New(errMsg))))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/environments", baseUrl)

		orgID, _ := cmd.Flags().GetString("orgID")
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		payload := &environments.EnvironmentPayload{
			Name:        name,
			Description: description,
			OrgID:       orgID,
		}

		err = api.Add(url, payload)
		if err != nil {

			return errors.New(display.GetErrorLog(ErrEnvironmentCreate(err)))
		}

		utils.Log.Info(fmt.Sprintf("Environment named %s created in organization id %s", payload.Name, payload.OrgID))
		return nil
	},
}

func init() {
	createEnvironmentCmd.Flags().StringP("orgID", "o", "", "Organization ID")
	createEnvironmentCmd.Flags().StringP("name", "n", "", "Name of the environment")
	createEnvironmentCmd.Flags().StringP("description", "d", "", "Description of the environment")
}
