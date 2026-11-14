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

package connections

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/cobra"
)

var deleteConnectionCmd = &cobra.Command{
	Use:   "delete [connection-id]",
	Short: "Delete a connection",
	Long: `Delete a connection.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/connection/delete`,

	Example: `
// Delete a connection
mesheryctl connection delete [connection_id]
`,

	Args: func(_ *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errNoArgMsg, deleteUsageMsg))
		}
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s\n%s", errMultiArgMsg, deleteUsageMsg))
		}

		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid connection ID: %q", args[0]))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.Delete(fmt.Sprintf("%s/%s", connectionApiPath, args[0]))
		if err != nil {
			if errors.GetCode(err) == utils.ErrNotFoundCode {
				utils.Log.Warnf("No connection with ID \"%s\" found", args[0])
				return nil
			}

			return err
		}
		defer utils.SafeClose(resp.Body)

		utils.Log.Infof("Connection with ID: %q is deleted", args[0])
		return nil
	},
}
