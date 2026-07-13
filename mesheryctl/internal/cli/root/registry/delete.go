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
	"fmt"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [connection-id]",
	Short: "Delete all models of a registrant connection",
	Long:  `Delete all model definitions, components, policies, and relationships associated with a registrant connection ID.`,
	Example: `
// Delete all models of a registrant connection by ID
mesheryctl registry delete [connection-id]
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(errors.New("Please provide a connection ID. Use 'mesheryctl registry delete --help' to display usage guide.\n"))
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
			return ErrDeleteRegistry(err, connectionID)
		}

		if resp.StatusCode == http.StatusNotFound {
			utils.Log.Warnf("No registrant connection with ID %q found or no models to delete", connectionID)
			return nil
		}

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code from Meshery Server: %d", resp.StatusCode)
		}

		utils.Log.Infof("Successfully deleted all models associated with registrant connection ID %q", connectionID)
		return nil
	},
}
