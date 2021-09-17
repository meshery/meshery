// Copyright 2020 Layer5, Inc.
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
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Stop, then start Meshery",
	Long:  `Restart all Meshery containers / pods.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
		return restart()
	},
}

func restart() error {
	log.Info("Restarting Meshery...")

	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return ErrProcessingMctlConfig(err)
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

	running, err := utils.IsMesheryRunning(currPlatform)
	if err != nil {
		return err
	}
	if !running { // Meshery is not running
		if err := start(); err != nil {
			return ErrRestartMeshery(err)
		}
	} else {
		if err := stop(); err != nil {
			return ErrRestartMeshery(err)
		}

		if err := start(); err != nil {
			return ErrRestartMeshery(err)
		}
	}
	return nil
}

func init() {
	restartCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
}
