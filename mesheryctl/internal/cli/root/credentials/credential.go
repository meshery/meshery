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
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/system"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	user_id              string
	name                 string
	credential_type      string
	secret               string
	availableSubcommands []*cobra.Command
)

var CredentialCmd = &cobra.Command{
	Use:   "credential",
	Short: "Manage credentials",
	Long: `View and manage list of credentials for Meshery Environment.
Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
	Example: `
// Base command for credentials:
mesheryctl exp credential [subcommands]
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
			return utils.ErrInvalidArgument(errors.New("Please provide a subcommand with the command"))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.CredentialsError(fmt.Sprintf("invalid subcommand: %s. Please provide options from [view]. Use 'mesheryctl exp credential --help' to display usage guide.\n", args[0]), "credential"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		if err := cmd.Usage(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	CredentialCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")
	createCredentialCmd.Flags().StringVarP(&user_id, "user-id", "u", "", "User ID")
	createCredentialCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the credential")
	createCredentialCmd.Flags().StringVarP(&credential_type, "type", "y", "", "Type of the credential")
	createCredentialCmd.Flags().StringVarP(&secret, "secret", "s", "", "Secret of the credential")
	availableSubcommands = []*cobra.Command{listCredentialCmd, createCredentialCmd, DeleteCredenetialCmd, viewCredentialCmd}
	CredentialCmd.AddCommand(availableSubcommands...)
}
