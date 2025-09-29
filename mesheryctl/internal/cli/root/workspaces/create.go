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
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/workspace"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var workspacePayload workspace.WorkspacePayload

var createWorkspaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace under an organization",
	Long: `Create a new workspace by providing the name, description, and organization ID
Documentation for models can be found at https://docs.meshery.io/reference/mesheryctl/exp/workspace/create`,
	Example: `
// Create a new workspace in an organization
mesheryctl exp workspace create --orgId [orgId] --name [name] --description [description]
`,

	Args: func(cmd *cobra.Command, args []string) error {
		const errorMsg = "[ Organization ID | Workspace name | Workspace description ] aren't specified\n\nUsage: \nmesheryctl exp workspace create --orgId [Organization ID] --name [name] --description [description]\nmesheryctl exp workspace create --help' to see detailed help message"

		if workspacePayload.OrganizationID == "" || workspacePayload.Name == "" || workspacePayload.Description == "" {

			return utils.ErrInvalidArgument(errors.New(errorMsg))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		payloadBytes, err := json.Marshal(workspacePayload)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		_, err = api.Add(workspacesApiPath, bytes.NewBuffer(payloadBytes))
		if err != nil {
			if meshkitErr, ok := err.(*mErrors.Error); ok {
				if meshkitErr.Code == utils.ErrFailRequestCode {
					errMsg := utils.WorkspaceSubError("Request failed.\nEnsure meshery server is available.", "create")
					return utils.ErrFailRequest(errors.New(errMsg))
				}
				if meshkitErr.Code == utils.ErrFailReqStatusCode {
					errMsg := utils.WorkspaceSubError(fmt.Sprintf("Failed to create \"%s\" workspace in \"%s\" organization.\nEnsure you provide a valid organization ID.\n", workspacePayload.Name, workspacePayload.OrganizationID), "create")
					return utils.ErrNotFound(fmt.Errorf("%s", errMsg))
				}
			}
			return err
		}

		utils.Log.Info("Workspace ", workspacePayload.Name, " created in organization ", workspacePayload.OrganizationID)
		return nil
	},
}

func init() {
	createWorkspaceCmd.Flags().StringVarP(&workspacePayload.OrganizationID, "orgId", "o", "", "Organization ID")
	createWorkspaceCmd.Flags().StringVarP(&workspacePayload.Name, "name", "n", "", "Name of the workspace")
	createWorkspaceCmd.Flags().StringVarP(&workspacePayload.Description, "description", "d", "", "Description of the workspace")
}
