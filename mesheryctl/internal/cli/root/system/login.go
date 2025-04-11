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

package system

import (
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	providerFlag string
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to a Meshery Server",
	Long: `
Authenticate to the Local or a Remote Provider of a Meshery Server

The authentication mode is web-based browser flow`,
	Args: cobra.MinimumNArgs(0),
	Example: `
// Login with the Meshery Provider of your choice: the Local Provider or a Remote Provider.
mesheryctl system login

// Login with the Meshery Provider by specifying it via -p or --provider flag.
mesheryctl system login -p Meshery
    `,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrProcessingLoginConfig(err)
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrRetrievingCurrentContext(err)
		}

		isRunning, err := utils.IsMesheryRunning(currCtx.GetPlatform())
		if err != nil {
			return ErrCheckServerStatus(err)
		}

		if !isRunning {
			return ErrMesheryServerNotRunning()
		}

		var tokenData []byte
		if providerFlag != "" {
			var provider = providerFlag
			tokenData, err = utils.InitiateLogin(mctlCfg, provider)
		} else {
			tokenData, err = utils.InitiateLogin(mctlCfg, "")
		}

		if err != nil {
			return ErrAuthenticationFailed(mctlCfg.GetBaseMesheryURL())
		}

		log.Println("authenticated")

		token, err := mctlCfg.GetTokenForContext(mctlCfg.GetCurrentContextName())
		if err != nil {
			token.Location = utils.AuthConfigFile
			if err := config.AddTokenToConfig(token, utils.DefaultConfigPath); err != nil {
				return ErrAddTokenToConfig(err)
			}
		}

		if err := os.WriteFile(token.GetLocation(), tokenData, 0666); err != nil {
			return ErrWriteTokenToFile(err)
		}

		return nil
	},
}

func init() {
	loginCmd.PersistentFlags().StringVarP(&providerFlag, "provider", "p", "", "login Meshery with specified provider")
	loginCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		// Don't show immediate help, just return the structured error
		return ErrInvalidFlag(err)
	})
}
