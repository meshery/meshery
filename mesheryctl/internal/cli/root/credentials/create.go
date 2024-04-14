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

package credentials

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/manifoldco/promptui"
)

var createCredentialCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new credential",
	Long:  `Create a new credential by providing the name, user ID, type, and secret of the credential`,
	Example: `
// Create a new credential
mesheryctl exp credential create --user-id [user-id]
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
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
		user_id, _ = cmd.Flags().GetString("user-id")

		if user_id == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return errors.New("user-id is required")
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		// Prompt for input
		prompt := promptui.Prompt{
			Label: "Name",
		}
		name, err := prompt.Run()
		if err != nil {
			return err
		}

		prompt = promptui.Prompt{
			Label: "Type",
		}
		credentialType, err := prompt.Run()
		if err != nil {
			return err
		}

		prompt = promptui.Prompt{
			Label: "Secret",
		}
		secret, err := prompt.Run()
		if err != nil {
			return err
		}

		if name == "" || credentialType == "" || secret == "" {
			return utils.ErrInvalidArgument(errors.New("name, type, and secret are required"))
		}

		baseURL := mctlCfg.GetBaseMesheryURL()
		url := fmt.Sprintf("%s/api/integrations/credentials", baseURL)

		// Generate a unique identifier for the credential
		id, err := uuid.NewV4()
		if err != nil {
			return err
		}

		// Parse the user_id as UUID
		parsedUserID, err := uuid.FromString(user_id)
		if err != nil {
			return utils.ErrInvalidArgument(errors.New("invalid user_id format"))
		}

		secretMap := make(map[string]interface{})
		secretMap["key"] = secret

		// Construct the payload according to the schema
		payload := &models.Credential{
			ID:        id,
			UserID:    &parsedUserID,
			Name:      name,
			Type:      credentialType,
			Secret:    secretMap,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: sql.NullTime{}, // Set to zero value
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		req, err := utils.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		if resp.StatusCode == http.StatusOK {
			utils.Log.Info("Credential created successfully")
			return nil
		}
		utils.Log.Info("Error creating credential")
		return nil
	},
}
