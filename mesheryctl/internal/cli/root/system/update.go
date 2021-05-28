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
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Pull new Meshery images/manifest files and operator manifests.",
	Long:  `Pull new Meshery container images/manifests and operator manifests from Github/Docker Hub and pull if new version(s) are available.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		return RunPreflightHealthChecks(true, cmd.Use)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		currCtx, err := mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return err
		}
		if currCtx.Version != "latest" {
			// ask confirmation if user has pinned the version in config
			log.Infof("You have pinned version: %s in your current context", currCtx.Version)
			userResponse := false
			if utils.SilentFlag {
				userResponse = true
			} else {
				userResponse = utils.AskForConfirmation("Updating Meshery container images/manifest files and Meshery operator manifests will supersede the version to latest. Are you sure you want to continue")
			}

			if !userResponse {
				log.Info("Update aborted.")
				return nil
			}
			currCtx.Version = "latest"
		}

		switch currCtx.Platform {
		case "docker":
			if !utils.SkipResetFlag {

				log.Printf("Fetching latest docker-compose file for channel: %s...\n", currCtx.Channel)
				err = utils.DownloadDockerComposeFile(currCtx, true)
				if err != nil {
					return errors.Wrap(err, "failed to fetch docker-compose file")
				}

				log.Info("Updating Meshery...")

				err = utils.UpdateMesheryContainers()
				if err != nil {
					return errors.Wrap(err, utils.SystemError("failed to update Meshery containers"))
				}

				// applying operator manifest
				kubeClient, err := meshkitkube.New([]byte(""))
				if err != nil {
					return err
				}

				running, err := utils.IsMesheryRunning(currCtx.Platform)
				if err != nil {
					return err
				}

				if !running {
					err = utils.CreateManifestsFolder()

					if err != nil {
						return err
					}
				}

				err = utils.DownloadOperatorManifest()

				if err != nil {
					return err
				}

				err = utils.ApplyOperatorManifest(kubeClient, true, false)

				if err != nil {
					return err
				}
			}

			skipUpdateFlag = true
			err = restart()

			if err != nil {
				return err
			}

			utils.ViperDocker.SetConfigFile(utils.DefaultConfigPath)
			err := utils.ViperDocker.ReadInConfig()
			if err != nil {
				return err
			}

			dockerCompose := &config.MesheryCtlConfig{}
			err = utils.ViperDocker.Unmarshal(&dockerCompose)
			if err != nil {
				return err
			}

			utils.ViperDocker.Set("contexts."+mctlCfg.CurrentContext, currCtx)

			err = utils.ViperDocker.WriteConfig()
			if err != nil {
				return err
			}

		case "kubernetes":
			// If the user skips reset, then just restart the pods else fetch updated manifest files and apply them
			if !utils.SkipResetFlag {

				log.Debug("creating new Clientset...")
				// Create a new client
				client, err := meshkitkube.New([]byte(""))

				if err != nil {
					return errors.Wrap(err, "failed to create new client")
				}

				version := currCtx.Version
				RequestedAdapters := currCtx.Adapters

				if version == "latest" {
					if currCtx.Channel == "edge" {
						version = "master"
					} else {
						version, err = utils.GetLatestStableReleaseTag()
						if err != nil {
							return err
						}
					}
				}

				// fetch the manifest files corresponding to the version specified
				manifests, err := utils.FetchManifests(version)

				if err != nil {
					return err
				}

				// downloaded required files successfully now apply the manifest files
				log.Info("Updating Meshery...")

				log.Info("applying the manifests to Kubernetes cluster...")

				// apply the adapters mentioned in the config.yaml file to the Kubernetes cluster
				err = utils.ApplyManifestFiles(manifests, RequestedAdapters, client, true, false)

				if err != nil {
					return err
				}

				err = utils.ApplyOperatorManifest(client, true, false)

				if err != nil {
					return err
				}
			}

			skipUpdateFlag = true

			// restart the pods in meshery namespace
			err = restart()

			if err != nil {
				return err
			}

			utils.ViperK8s.SetConfigFile(utils.DefaultConfigPath)
			err := utils.ViperK8s.ReadInConfig()
			if err != nil {
				return err
			}

			kubeCompose := &config.MesheryCtlConfig{}
			err = utils.ViperK8s.Unmarshal(&kubeCompose)
			if err != nil {
				return err
			}

			utils.ViperK8s.Set("contexts."+mctlCfg.CurrentContext, currCtx)

			err = utils.ViperK8s.WriteConfig()
			if err != nil {
				return err
			}

			log.Info("... updated Meshery in the Kubernetes Cluster.")
		}

		log.Info("Meshery is now up-to-date")
		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&utils.SkipResetFlag, "skip-reset", "", false, "(optional) skip checking for new Meshery manifest files.")
}
