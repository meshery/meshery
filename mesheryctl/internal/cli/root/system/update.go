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
	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images/manifest files.",
	Long:  `Pull new Meshery container images and manifests from artifact repository.`,
	Args:  cobra.NoArgs,
	Example: `
// Pull new Meshery images from Docker Hub. This does not update mesheryctl. This command may be executed while Meshery is running.
mesheryctl system update

// Pull the latest manifest files alone
mesheryctl system update --skip-reset
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return errors.Wrapf(err, "failed to initialize healthchecker")
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		if tempContext != "" {
			err = mctlCfg.SetCurrentContext(tempContext)
			if err != nil {
				return errors.Wrap(err, "failed to set temporary context")
			}
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
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
			if !utils.SkipResetFlag {
				err := fetchManifests(mctlCfg)

				if err != nil {
					return err
				}
			}

			err = utils.UpdateMesheryContainers()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to update Meshery containers"))
			}

			err = config.UpdateContextInConfig(viper.GetViper(), currCtx, mctlCfg.GetCurrentContextName())

			if err != nil {
				return err
			}

		case "kubernetes":
			// create a client
			kubeClient, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}
			// If the user skips reset, then just restart the pods else fetch updated manifest files and apply them
			if !utils.SkipResetFlag {
				// get value overrides to install the helm chart
				overrideValues := utils.SetOverrideValues(currCtx, "latest")

				// Apply the latest helm chart along with the default image tag specified in the charts "stable-latest"
				if err = kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
					Namespace:       utils.MesheryNamespace,
					CreateNamespace: true,
					ChartLocation: meshkitkube.HelmChartLocation{
						Repository: utils.HelmChartURL,
						Chart:      utils.HelmChartName,
					},
					Action:         meshkitkube.UPGRADE,
					OverrideValues: overrideValues,
				}); err != nil {
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
				return errors.Wrapf(err, "failed to initialize healthchecker")
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
			err = config.UpdateContextInConfig(viper.GetViper(), currCtx, mctlCfg.GetCurrentContextName())
			if err != nil {
				return err
			}

			log.Info("... updated Meshery in the Kubernetes Cluster.")
		}

		log.Info("Meshery is now up-to-date")
		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		latest, err := utils.GetLatestStableReleaseTag()
		version := constants.GetMesheryctlVersion()
		if err == nil && latest != version {
			log.Printf("A new release of mesheryctl is available: %s → %s", version, latest)
			log.Printf("https://github.com/layer5io/meshery/releases/tag/%s", latest)
			log.Print("Check https://docs.meshery.io/guides/upgrade#upgrading-meshery-cli for instructions on how to update mesheryctl\n")
		}
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&utils.SkipResetFlag, "skip-reset", "", false, "(optional) skip checking for new Meshery manifest files.")
}
