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
	Short: "Create a new workspaces",
	Long:  `Create a new workspaces by providing the name and description of the workspace`,
	Example: `
// Create a new workspace
mesheryctl exp workspace create --orgId [orgId] --name [workspace-name] --description [workspace-description]

// Documentation for workspace can be found at:
https://docs.layer5.io/cloud/spaces/workspaces/
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return system.ErrGetCurrentContext(err)
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},

	Args: func(cmd *cobra.Command, args []string) error {

		// Check if all three flags are set
		orgIdFlag, _ := cmd.Flags().GetString("orgId")
		nameFlag, _ := cmd.Flags().GetString("name")
		descriptionFlag, _ := cmd.Flags().GetString("description")

		if orgIdFlag == "" || nameFlag == "" || descriptionFlag == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New("all three flags are required")
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

		payload := &models.WorkspacePayload{
			Name:           name,
			Description:    description,
			OrganizationID: orgID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		req, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		_, err = utils.MakeRequest(req)
		if strings.Contains(err.Error(), "201") {
			utils.Log.Info("Workspace created successfully")
			return nil
		}

		utils.Log.Error(errors.New("Unable to create workspace. " + err.Error()))
		return nil
	},
}
