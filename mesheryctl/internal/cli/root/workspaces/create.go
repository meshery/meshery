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

	googleUUID "github.com/google/uuid"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/workspace"

	"github.com/spf13/cobra"
)

type cmdWorkspaceCreateFlags struct {
	OrganizationID string `json:"orgId" validate:"required,uuid"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"omitempty"`
}

var workspaceCreateFlags cmdWorkspaceCreateFlags

var createWorkspaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace under an organization",
	Long: `Create a new workspace by providing the name, description, and organization ID
Find more information at: https://docs.meshery.io/reference/mesheryctl/workspace/create`,
	Example: `
// Create a new workspace in an organization
mesheryctl workspace create --orgId [orgId] --name [name] --description [description]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &workspaceCreateFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		organizationID, err := googleUUID.Parse(workspaceCreateFlags.OrganizationID)
		if err != nil {
			return err
		}

		workspaceCreatePayload := workspace.WorkspacePayload{
			OrganizationID: organizationID,
			Name:           workspaceCreateFlags.Name,
			Description:    workspaceCreateFlags.Description,
		}
		payloadBytes, err := json.Marshal(workspaceCreatePayload)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		_, err = api.Add(workspacesApiPath, bytes.NewBuffer(payloadBytes), nil)
		if err != nil {
			if mErrors.GetCode(err) == utils.ErrNotFoundCode {
				return returnFailedCreateWorkspaceError(workspaceCreateFlags.Name, workspaceCreateFlags.OrganizationID)
			}
			return err
		}

		utils.Log.Infof("Workspace %s created in organization %s", workspaceCreateFlags.Name, workspaceCreateFlags.OrganizationID)
		return nil
	},
}

func init() {
	createWorkspaceCmd.Flags().StringVarP(&workspaceCreateFlags.OrganizationID, "orgId", "o", "", "Organization ID")
	createWorkspaceCmd.Flags().StringVarP(&workspaceCreateFlags.Name, "name", "n", "", "Name of the workspace")
	createWorkspaceCmd.Flags().StringVarP(&workspaceCreateFlags.Description, "description", "d", "", "(Optional) Description of the workspace")
}

func returnFailedCreateWorkspaceError(name, organizationID string) error {
	errMsg := utils.WorkspaceSubError(fmt.Sprintf("Failed to create \"%s\" workspace in \"%s\" organization.\nEnsure you provide a valid organization ID.\n", name, organizationID), "create")
	return utils.ErrNotFound(fmt.Errorf("%s", errMsg))
}
