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
	Short: "Pull new Meshery images from Docker Hub.",
	Long:  `Pull Docker Hub for new Meshery container images and pulls if new image version(s) are available.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		if tempContext != "" {
			return utils.PreReqCheck(cmd.Use, tempContext)
		}
		return utils.PreReqCheck(cmd.Use, "")
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
			log.Infof("You have pinned version: %s in your current conext", currCtx.Version)
			userResponse := false
			if utils.SilentFlag {
				userResponse = true
			} else {
				userResponse = utils.AskForConfirmation("Updating Meshery container images will supersede the version to latest. Are you sure you want to continue")
			}

			if !userResponse {
				log.Info("Update aborted.")
				return nil
			}
			currCtx.Version = "latest"
		}

		switch currCtx.Platform {
		case "docker":
			log.Printf("Fetching latest docker-compose file for channel: %s...\n", currCtx.Channel)
			err = utils.DownloadDockerComposeFile(currCtx, true)
			if err != nil {
				return errors.Wrap(err, "failed to fetch docker-compose file")
			}

			err = utils.UpdateMesheryContainers()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to update meshery containers"))
			}

		case "kubernetes":

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
			log.Info("... updated Meshery in the Kubernetes Cluster.")

		}

		log.Info("Meshery is now up-to-date")
		return nil
	},
}
