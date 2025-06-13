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
	"io"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
		// Check if all three flags are set
		orgIdFlag, _ := cmd.Flags().GetString("orgId")
		nameFlag, _ := cmd.Flags().GetString("name")
		descriptionFlag, _ := cmd.Flags().GetString("description")

		const errorMsg = "[ Organization ID | Workspace name | Workspace description ] aren't specified\n\nUsage: \nmesheryctl  exp workspace --orgId [Organization ID] --name [name] --description [description]\nmesheryctl  exp workspace --help' to see detailed help message"

		if orgIdFlag == "" || nameFlag == "" || descriptionFlag == "" {

			return utils.ErrInvalidArgument(errors.New(errorMsg))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		orgIdFlag, _ := cmd.Flags().GetString("orgId")

		// verify orgId
		orgsUrl := fmt.Sprintf("%s/api/identity/orgs", baseUrl)
		utils.Log.Info("🔍 sending request to verify organization ID: ", orgsUrl)

		reqOrg, err := utils.NewRequest(http.MethodGet, orgsUrl, nil)
		if err != nil {
			utils.Log.Info("failed to create request to fetch orgs")
			return err
		}

		respOrg, err := utils.MakeRequest(reqOrg)
		if err != nil {
			utils.Log.Info("failed to fetch organizations")
			return err
		}
		defer respOrg.Body.Close()

		utils.Log.Info("received response status code: ", respOrg.StatusCode)

		bodyBytes, _ := io.ReadAll(respOrg.Body)

		type organization struct {
			ID string `json:"id"`
		}

		type orgListResponse struct {
			Organizations []organization `json:"organizations"`
		}

		var orgResp orgListResponse
		if err := json.Unmarshal(bodyBytes, &orgResp); err != nil {
			utils.Log.Info("failed to parse organization list")
			return err
		}

		found := false
		for _, org := range orgResp.Organizations {
			if org.ID == orgIdFlag {
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid organization ID: not found in your account's organization list, please make sure your Org is created and existed ;)")
		}
		// Organization ID verified
		utils.Log.Info("✅ Organization ID verified")
		url := fmt.Sprintf("%s/%s", baseUrl, workspacesApiPath)

		nameFlag, _ := cmd.Flags().GetString("name")
		descriptionFlag, _ := cmd.Flags().GetString("description")


		payload := &models.WorkspacePayload{
			Name:           nameFlag,
			Description:    descriptionFlag,
			OrganizationID: orgIdFlag,
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			utils.Log.Info("Failed to marshal payload")
			return err
		}

		req, err := utils.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.Log.Info("Failed to create request")
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		if resp.StatusCode == http.StatusCreated {
			utils.Log.Info("Workspace ", nameFlag, " created")
			return nil
		}

		utils.Log.Info("Failed to create ", nameFlag, " workspace")
		return nil
	},
}

func init() {
	createWorkspaceCmd.Flags().StringP("orgId", "o", "", "Organization ID")
	createWorkspaceCmd.Flags().StringP("name", "n", "", "Name of the workspace")
	createWorkspaceCmd.Flags().StringP("description", "d", "", "Description of the workspace")
}
