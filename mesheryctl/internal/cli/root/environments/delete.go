// Copyright Meshery Authorsayer5, Inc.
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
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteEnvironmentCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an environment",
	Long: `Delete an environment by providing the environment ID
Documentation for environment can be found at Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/delete`,
	Example: `
// delete a new environment
mesheryctl environment delete [environmentId]
`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			const errMsg = "[ Environment ID ] isn't specified\n\nUsage: mesheryctl environment delete [environmentId]\nRun 'mesheryctl environment delete --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/environments/%s", baseUrl, args[0])
		req, err := utils.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		// Check if the response status code is 200
		if resp.StatusCode == http.StatusOK {
			utils.Log.Info(fmt.Sprintf("Environment with ID %s has been deleted", args[0]))
			return nil
		}

		return utils.ErrBadRequest(errors.New(fmt.Sprintf("failed to delete environment with id %s", args[0])))
	},
}
