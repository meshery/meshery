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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	c "github.com/layer5io/meshery/mesheryctl/pkg/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkDocUpdate = map[string]string{
	"link":    "![update-usage](/assets/img/mesheryctl/update.png)",
	"caption": "Usage of mesheryctl system update",
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images/manifest files.",
	Long:  `Pull new Meshery container images and manifests from artifact repository.`,
	Example: `
// Pull new Meshery images from Docker Hub. This does not update mesheryctl. This command may be executed while Meshery is running.
mesheryctl system update

// Pull the latest manifest files alone
mesheryctl system update --skip-reset
	`,
	Annotations: linkDocUpdate,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "update"))
		}

		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		if tempContext != "" {
			err = mctlCfg.SetCurrentContext(tempContext)
			if err != nil {
				utils.Log.Error(ErrSetCurrentContext(err))
				return nil
			}
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		err = currCtx.ValidateVersion()
		if err != nil {
			return err
		}
		if currCtx.GetVersion() != "latest" {
			// ask confirmation if user has pinned the version in config
			log.Infof("You have pinned version: %s in your current context", currCtx.GetVersion())
			userResponse := false
			if utils.SilentFlag {
				userResponse = true
			} else {
				userResponse = utils.AskForConfirmation("Updating Meshery container images/manifest files will supersede the version to latest. Are you sure you want to continue")
			}

			if !userResponse {
				log.Info("Update aborted.")
				return nil
			}
			currCtx.SetVersion("latest")
		}

		log.Info("Updating Meshery...")

		switch currCtx.GetPlatform() {
		case "docker":
			log.Info("Updating Meshery containers")
			err = utils.UpdateMesheryContainers()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to update Meshery containers"))
			}

			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())

			if err != nil {
				return err
			}

		case "kubernetes":
			// create a client
			kubeClient, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}
			mesheryImageVersion := currCtx.GetVersion()
			providerURL := viper.GetString(c.ProviderURLsENV)
			// If the user skips reset, then just restart the pods else fetch updated manifest files and apply them
			if !utils.SkipResetFlag {

				// Apply the latest helm chart along with the default image tag specified in the charts "stable-latest"
				if err = applyHelmCharts(kubeClient, currCtx, mesheryImageVersion, false, meshkitkube.UPGRADE, "", providerURL); err != nil {
					return errors.Wrap(err, "cannot update Meshery")
				}
			}

			// run k8s checks to make sure if k8s cluster is running
			hcOptions := &HealthCheckOptions{
				PrintLogs:           false,
				IsPreRunE:           false,
				Subcommand:          "",
				RunKubernetesChecks: true,
			}
			hc, err := NewHealthChecker(hcOptions)
			if err != nil {
				utils.Log.Error(err)
				return nil
			}
			// If k8s is available in case of platform docker than we deploy operator
			if err = hc.Run(); err != nil {
				return ErrHealthCheckFailed(err)
			}

			running, err := utils.AreMesheryComponentsRunning(currCtx.GetPlatform())
			if err != nil {
				return err
			}
			if !running {
				// Meshery is not running, run the start command
				if err := start(); err != nil {
					return ErrRestartMeshery(err)
				}
			}

			currCtx.SetVersion("latest")
			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
			if err != nil {
				return err
			}

			log.Info("... updated Meshery in the Kubernetes Cluster.")
		}

		log.Info("Meshery is now up-to-date")
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.CheckMesheryctlClientVersion(constants.GetMesheryctlVersion())
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&utils.SkipResetFlag, "skip-reset", "", false, "(optional) skip checking for new Meshery manifest files.")
}
