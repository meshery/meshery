// Copyright 2024 Layer5, Inc.
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
	"net/http"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateWorkspaceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace",
	Long:  `Create a new workspace`,
	Example: `
// Create a new workspace
mesheryctl exp workspace create orgId --name [workspace-name] --description [workspace-description]`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			utils.Log.Error(err)
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(system.ErrGetCurrentContext(err))
			return nil
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return nil
	},

	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp workspace list \nRun 'mesheryctl exp workspace list --help' to see detailed help message"
		if len(args) >= 3 {
			return utils.ErrInvalidArgument(fmt.Errorf("%s: expected 3 arguments (name, orgId, description)", errMsg))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/workspaces", baseUrl)

		if name == "" {
			return utils.ErrInvalidArgument(errors.New("name is required"))
		}

		if description == "" {
			return utils.ErrInvalidArgument(errors.New("description is required"))
		}
		payload := &models.WorkspacePayload{
			Name:           name,
			Description:    description,
			OrganizationID: orgID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		req, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		_, err = utils.MakeRequest(req)
		// http code 201 will be returned if the workspace is created successfully

		if strings.Contains(err.Error(), "201") {
			fmt.Println("Workspace created successfully")
			return nil
		} else {
			fmt.Println("Error creating workspace")
			return err
		}

		return nil
	},
}
