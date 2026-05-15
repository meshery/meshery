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

package workspaces

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deleteWorkspaceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a workspace",
	Long: `Delete a workspace by providing the workspace ID
Find more information at: https://docs.meshery.io/reference/mesheryctl/workspace/delete`,
	Example: `
// Delete a workspace
mesheryctl workspace delete [workspaceId]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			const errMsg = "[ Workspace ID ] isn't specified\n\nUsage: mesheryctl workspace delete [workspaceId]\nRun 'mesheryctl workspace delete --help' to see detailed help message"
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %s", args[0]))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := api.Delete(fmt.Sprintf("%s/%s", workspacesApiPath, args[0]))
		if err != nil {
			return err
		}

		utils.Log.Infof("Workspace with ID %s has been deleted", args[0])
		return nil
	},
}
