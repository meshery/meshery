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
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long:  `Stop all Meshery containers, remove their instances and prune their connected volumes.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if tempContext != "" {
			return utils.PreReqCheck(cmd.Use, tempContext)
		}
		return utils.PreReqCheck(cmd.Use, "")
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := stop(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop Meshery"))
		}
		return nil
	},
}

func stop() error {
	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// if a temp context is set using the -c flag, use it as the current context
	currCtx, err := mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	// Get the current platform and the specified adapters in the config.yaml
	RequestedAdapters := currCtx.Adapters

	switch currCtx.Platform {
	case "docker":
		// if the platform is docker, then stop all the running containers
		ok, err := utils.IsMesheryRunning(currCtx.Platform)
		if err != nil {
			return err
		}
		if !ok {
			log.Info("Meshery is not running. Nothing to stop.")
			return nil
		}
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to mkdir %s", utils.MesheryFolder)))
			}
		}

		log.Info("Stopping Meshery...")

		// Stop all Docker containers
		stop := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "stop")
		stop.Stdout = os.Stdout
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop meshery - could not stop some containers."))
		}

		// Remove all Docker containers
		stop = exec.Command("docker-compose", "-f", utils.DockerComposeFile, "rm", "-f")
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop meshery"))
		}

		// Mesheryctl uses a docker volume for persistence. This volume should only be cleared when user wants
		// to start from scratch with a fresh install.
		// if err := exec.Command("docker", "volume", "prune", "-f").Run(); err != nil {
		// 	log.Fatal("[ERROR] Please install docker-compose. The error message: \n", err)
		// }

	case "kubernetes":
		// if the platform is kubernetes, stop the deployment by deleting the manifest files
		ok, err := utils.IsMesheryRunning(currCtx.Platform)
		if err != nil {
			return err
		}
		if !ok {
			log.Info("Meshery is not running. Nothing to stop.")
			return nil
		}

		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
		}

		if !userResponse {
			log.Info("Stop aborted.")
			return nil
		}

		// create an kubernetes client
		client, err := meshkitkube.New([]byte(""))

		if err != nil {
			return err
		}

		// check if the manifest folder exists on the machine
		if _, err := os.Stat(filepath.Join(utils.MesheryFolder, utils.ManifestsFolder)); os.IsNotExist(err) {
			log.Errorf("%s folder does not exist.", utils.ManifestsFolder)
			return err
		}

		version := currCtx.Version
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
		// get correct manfestsURL based on version
		manifestsURL, err := utils.GetManifestTreeURL(version)
		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}
		// pick all the manifest files stored in minfestsURL
		manifests, err := utils.ListManifests(manifestsURL)

		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}

		log.Info("Stopping Meshery...")

		// delete the Meshery deployment using the manifest files to stop Meshery
		err = utils.ApplyManifestFiles(manifests, RequestedAdapters, client, false, true)

		if err != nil {
			return err
		}
	}

	log.Info("Meshery is stopped.")

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to reset meshery config"))
		}
	}
	return nil
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
