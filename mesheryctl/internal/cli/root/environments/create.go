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

package environments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models/environments"
	"github.com/manifoldco/promptui"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CreateEnvironmentCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environments",
	Long:  `Create a new environments by providing the name and description of the environment`,
	Example: `
// Create a new environment
mesheryctl exp environment create [orgID] 
// Documentation for environment can be found at:
https://docs.layer5.io/cloud/spaces/environments/
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
		if len(args) == 0 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New(utils.EnvironmentSubError("Please provide a orgID", "create"))
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/environments", baseUrl)

		orgID := args[0]

		prompt := promptui.Prompt{
			Label: "Enter the name of the environment",
		}
		name, err := prompt.Run()
		if err != nil {
			return err
		}

		prompt = promptui.Prompt{
			Label: "Enter the description of the environment",
		}
		description, err := prompt.Run()
		if err != nil {
			return err
		}

		if name == "" || description == "" {
			return utils.ErrInvalidArgument(errors.New("name is required"))
		}

		payload := &environments.EnvironmentPayload{
			Name:        name,
			Description: description,
			OrgID:       orgID,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		req, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return err
		}

		resp, _ := utils.MakeRequest(req)

		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("environment created successfully")
			return nil
		}
		utils.Log.Info("Error creating environment")
		return nil
	},
}
