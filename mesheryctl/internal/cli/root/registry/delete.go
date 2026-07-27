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

package registry

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [connection-id]",
	Short: "Delete all models associated with a connection",
	Long:  `Delete all cached model definitions, components, policies, and relationships associated with a registrant connection ID. Note that this command does not delete the connection itself.`,
	Example: `
// Delete all models associated with a connection by ID
mesheryctl registry delete [connection-id]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("Please provide a connection ID. Use 'mesheryctl registry delete --help' to display usage guide.\n"))
		}
		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", args[0]))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionID := args[0]

		url := fmt.Sprintf("api/meshmodels/registrants/%s/models", connectionID)
		resp, err := api.Delete(url)
		if err != nil {
			if errors.GetCode(err) == utils.ErrNotFoundCode {
				errStr := err.Error()
				var bodyMap map[string]string
				if json.Unmarshal([]byte(errStr), &bodyMap) == nil {
					if errMsg, ok := bodyMap["error"]; ok {
						utils.Log.Warnf("%s", errMsg)
						return nil
					}
				}
				utils.Log.Warnf("No registrant connection with ID %q found or no models to delete", connectionID)
				return nil
			}
			utils.Log.Error(err)
			return ErrDeleteRegistry(err, connectionID)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			var bodyMap map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&bodyMap); err == nil {
				if errMsg, ok := bodyMap["error"]; ok {
					utils.Log.Warnf("%s", errMsg)
					return nil
				}
			}
			utils.Log.Warnf("No registrant connection with ID %q found or no models to delete", connectionID)
			return nil
		}

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			err := fmt.Errorf("unexpected status code from Meshery Server: %d", resp.StatusCode)
			utils.Log.Error(err)
			return ErrDeleteRegistry(err, connectionID)
		}

		type DeleteModelsResponse struct {
			Message        string `json:"message"`
			Count          int    `json:"count"`
			ConnectionName string `json:"connectionName"`
		}

		if resp.StatusCode == http.StatusOK {
			var successResp DeleteModelsResponse
			if err := json.NewDecoder(resp.Body).Decode(&successResp); err == nil {
				utils.Log.Infof("Deleted %d models for registrant %s.", successResp.Count, successResp.ConnectionName)
				return nil
			}
		}

		utils.Log.Infof("Successfully deleted all models associated with registrant connection ID %q", connectionID)
		return nil
	},
}
