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
	"bytes"
	"encoding/json"

	googleUUID "github.com/google/uuid"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/spf13/cobra"
)

type cmdEnvironmentCreateFlags struct {
	OrganizationID string `json:"orgId" validate:"required,uuid"`
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
}

var createEnvironmentFlags cmdEnvironmentCreateFlags

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long: `Create a new environment by providing the name and description of the environment
Find more information at: https://docs.meshery.io/reference/mesheryctl/environment/create`,
	Example: `
// Create a new environment
mesheryctl environment create --orgId [orgId] --name [name] --description [description]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &createEnvironmentFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		organizationID, err := googleUUID.Parse(createEnvironmentFlags.OrganizationID)
		if err != nil {
			return err
		}

		createEnvironmentPayload := environment.EnvironmentPayload{
			OrgId:       organizationID,
			Name:        createEnvironmentFlags.Name,
			Description: createEnvironmentFlags.Description,
		}
		payloadBytes, err := json.Marshal(&createEnvironmentPayload)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		_, err = api.Add(environmentApiPath, bytes.NewBuffer(payloadBytes), nil)
		if err != nil {
			if mErrors.GetCode(err) == utils.ErrFailReqStatusCode {
				return errCreateEnvironment(createEnvironmentFlags.Name, createEnvironmentFlags.OrganizationID)
			}
			return err
		}

		utils.Log.Infof("Environment named %s created in organization id %s", createEnvironmentFlags.Name, createEnvironmentFlags.OrganizationID)
		return nil
	},
}

func init() {
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlags.OrganizationID, "orgId", "o", "", "Organization ID")
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlags.Name, "name", "n", "", "Name of the environment")
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlags.Description, "description", "d", "", "Description of the environment")
}
