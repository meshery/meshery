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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// viewCredentialCmd represents the view command by Credential ID
var viewCredentialCmd = &cobra.Command{
	Use:   "view",
	Short: "View a credential",
	Long:  `View a credential details from list of Credentials by providing credential ID`,
	Example: `
// View a credential:
mesheryctl exp credential view [credential_ID]
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
		if len(args) != 1 {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(errors.New("credential ID"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		baseUrl := mctlCfg.GetBaseMesheryURL()

		credID := args[0]

		url := fmt.Sprintf("%s/api/integrations/credentials/%s", baseUrl, credID)

		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		credentialResponse := models.Credential{}
		err = json.Unmarshal(body, &credentialResponse)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		if credentialResponse.ID.String() == "" {
			utils.Log.Info("No credentials found for the given ID: ", credID)
			return nil
		}

		homedir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		// To save the credential in a file as json format
		credFile, err := os.Create(homedir + "/.meshery/credential.json")
		if err != nil {
			return err
		}
		defer credFile.Close()

		credBytes, err := json.MarshalIndent(credentialResponse, "", "  ")
		if err != nil {
			return err
		}

		_, err = credFile.Write(credBytes)
		if err != nil {
			return err
		}

		utils.Log.Info("Credential saved to file: $/homedir/.meshery/credential.json")

		return nil
	},
}
