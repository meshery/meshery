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
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta1/environment"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type createEnvironmentFlags struct {
	orgID       string
	name        string
	description string
}

var createEnvironmentFlagsProvided createEnvironmentFlags

var createEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long: `Create a new environment by providing the name and description of the environment
Documentation for environment can be found at https://docs.meshery.io/reference/mesheryctl/environment/create`,
	Example: `
// Create a new environment
mesheryctl environment create --orgID [orgID] --name [name] --description [description]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		const errMsg = "[ Organization ID | Name | Description ] aren't specified\n\nUsage: mesheryctl environment create --orgID [orgID] --name [name] --description [description]\nRun 'mesheryctl environment create --help' to see detailed help message"

		if createEnvironmentFlagsProvided.orgID == "" || createEnvironmentFlagsProvided.name == "" || createEnvironmentFlagsProvided.description == "" {
			return utils.ErrInvalidArgument(errors.New(errMsg))
		}

		if !utils.IsUUID(createEnvironmentFlagsProvided.orgID) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid Environment ID: %s", createEnvironmentFlagsProvided.orgID))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		payload := &environment.EnvironmentPayload{
			Name:        createEnvironmentFlagsProvided.name,
			Description: createEnvironmentFlagsProvided.description,
			OrgId:       createEnvironmentFlagsProvided.orgID, // TODO update OrgId in schema to OrgID
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		_, err = api.Add(environmentApiPath, bytes.NewBuffer(payloadBytes), nil)
		if err != nil {
			if meshkitErr, ok := err.(*mErrors.Error); ok {
				if meshkitErr.Code == utils.ErrFailReqStatusCode {
					return errCreateEnvironment(createEnvironmentFlagsProvided.name, createEnvironmentFlagsProvided.orgID)
				}
			}
			return err
		}

		utils.Log.Infof("Environment named %s created in organization id %s", createEnvironmentFlagsProvided.name, createEnvironmentFlagsProvided.orgID)
		return nil
	},
}

func init() {
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlagsProvided.orgID, "orgID", "o", "", "Organization ID")
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlagsProvided.name, "name", "n", "", "Name of the environment")
	createEnvironmentCmd.Flags().StringVarP(&createEnvironmentFlagsProvided.description, "description", "d", "", "Description of the environment")
}
