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
	stderrors "errors"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitErrors "github.com/meshery/meshkit/errors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cmdSystemLoginFlags struct {
	Provider string `json:"provider" validate:"omitempty"`
}

var systemLoginFlags cmdSystemLoginFlags

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate to a Meshery Server",
	Long: `
Authenticate to the Local or a Remote Provider of a Meshery Server

The authentication mode is web-based browser flow`,
	Args: cobra.NoArgs,
	Example: `
// Login with the Meshery Provider of your choice: the Local Provider or a Remote Provider.
mesheryctl system login

// Login with the Meshery Provider by specifying it via -p or --provider flag.
mesheryctl system login -p Meshery
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &systemLoginFlags)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		// if a temp context is set using the -c flag, use it as the current context
		if err := mctlCfg.SetCurrentContext(tempContext); err != nil {
			return ErrGetCurrentContext(err)
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}

		isRunning, err := utils.IsMesheryRunning(currCtx.GetPlatform())
		if err != nil {
			return utils.ErrMesheryCheckRunningStatus(err)
		}

		if !isRunning {
			return utils.ErrMesheryServerNotRunning(currCtx.GetPlatform())
		}

		var tokenData []byte
		if systemLoginFlags.Provider != "" {
			var provider = systemLoginFlags.Provider
			tokenData, err = utils.InitiateLogin(mctlCfg, provider)
		} else {
			tokenData, err = utils.InitiateLogin(mctlCfg, "")
		}

		if err != nil {
			// A structured MeshKit error (for example an unknown --provider)
			// already names its cause and remediation; wrapping it as a
			// connectivity failure would mislabel it.
			var structured *meshkitErrors.Error
			if stderrors.As(err, &structured) {
				utils.Log.Error(err)
				return err
			}
			return errors.Wrap(err, "authentication failed: Unable to reach Meshery Server. Verify system readiness with `mesheryctl system check`.")
		}

		utils.Log.Info("authenticated")

		token, err := utils.EnsureContextToken(mctlCfg, mctlCfg.GetCurrentContextName())
		if err != nil {
			return utils.ErrAttachAuthToken(err)
		}

		if err := utils.WriteTokenFile(token.GetLocation(), tokenData); err != nil {
			return utils.ErrCreateFile(token.GetLocation(), err)
		}

		return nil
	},
}

func init() {
	loginCmd.PersistentFlags().StringVarP(&systemLoginFlags.Provider, "provider", "p", "", "login Meshery with specified provider (run 'mesheryctl system provider list' to see options)")
}
