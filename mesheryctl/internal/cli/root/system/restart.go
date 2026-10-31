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
	"fmt"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type cmdSystemRestartFlags struct {
	SkipUpdate bool   `json:"skip-update" validate:"boolean"`
	Provider   string `json:"provider" validate:"omitempty"`
}

var systemRestartFlags cmdSystemRestartFlags

var (
	silentFlagSet bool
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Stop, then start Meshery",
	Long: `Restart all Meshery containers / pods.
Find more information at: https://docs.meshery.io/reference/mesheryctl/system/restart`,
	Example: `
// Restart all Meshery containers, their instances and their connected volumes
mesheryctl system restart

// (optional) skip checking for new updates available in Meshery.
mesheryctl system restart --skip-update
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return utils.ErrInvalidArgument(
				fmt.Errorf("restart does not take any arguments, see '%s --help' for more information", cmd.CommandPath()),
			)
		}
		return nil
	},

	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags
		err := mesheryctlflags.ValidateCmdFlags(cmd, &systemRestartFlags)
		if err != nil {
			return err
		}

		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return ErrHealthCheckFailed(err)
		}
		// execute healthchecks
		err = hc.RunPreflightHealthChecks()
		if err != nil {
			cmd.SilenceUsage = true
		}

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		skipUpdateFlag = systemRestartFlags.SkipUpdate
		providerFlag = systemRestartFlags.Provider
		return restart()
	},
}

func restart() error {
	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}
	// get the platform, channel and the version of the current context
	// if a temp context is set using the -c flag, use it as the current context
	if tempContext != "" {
		err = mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return ErrSettingTemporaryContext(err)
		}
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return ErrRetrievingCurrentContext(err)
	}

	currPlatform := currCtx.GetPlatform()

	running, err := utils.AreMesheryComponentsRunning(currPlatform)
	if err != nil {
		return ErrRestartMeshery(err)
	}
	if !running { // Meshery is not running
		if err := start(); err != nil {
			return ErrRestartMeshery(err)
		}
	} else {
		silentFlagSet = utils.SilentFlag
		if currPlatform == platformKubernetes {
			userResponse := false
			if utils.SilentFlag {
				userResponse = true
			} else {
				// ask user for confirmation
				userResponse = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
			}
			if !userResponse {
				utils.Log.Info("Restart aborted.")
				return nil
			}
			// skips asking for confirmation in the stop() function
			utils.SilentFlag = true
		}

		silentFlagSet = utils.SilentFlag
		utils.Log.Info("Restarting Meshery...")

		if err := stop(); err != nil {
			return ErrRestartMeshery(err)
		}

		// reset the silent flag to avoid overriding the flag for start command
		utils.SilentFlag = silentFlagSet

		if err := start(); err != nil {
			return ErrRestartMeshery(err)
		}
	}
	return nil
}

func init() {
	restartCmd.Flags().BoolVarP(&systemRestartFlags.SkipUpdate, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	restartCmd.Flags().StringVar(&systemRestartFlags.Provider, "provider", "", "Provider to use with the Meshery server")
}
