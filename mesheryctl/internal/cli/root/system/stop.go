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
	"os"
	"os/exec"
	"time"

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
	Long:  `Stop all Meshery containers / remove all Meshery pods.`,
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
		return hc.RunPreflightHealthChecks()
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
	err = mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return err
	}

	ok, err := utils.IsMesheryRunning(currCtx.GetPlatform())
	if err != nil {
		return err
	}
	if !ok {
		log.Info("Meshery is not running. Nothing to stop.")
		return nil
	}

	// Get the current platform and the specified adapters in the config.yaml
	RequestedAdapters := currCtx.GetAdapters()

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

	switch currCtx.GetPlatform() {
	case "docker":
		// if the platform is docker, then stop all the running containers
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

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
			return ErrStopMeshery(err)
		}
	}

	// If k8s is available in case of platform docker then we remove check the meshery namespace and remove all pods
	hcOptions := &HealthCheckOptions{
		PrintLogs:           false,
		IsPreRunE:           false,
		Subcommand:          "",
		RunKubernetesChecks: true,
	}
	hc, err := NewHealthChecker(hcOptions)
	if err != nil {
		return ErrHealthCheckFailed(err)
	}
	// stopping meshery pods if k8s is running
	if err = hc.Run(); err == nil {
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		_, version, err := utils.GetChannelAndVersion(currCtx)
		if err != nil {
			log.Println("failed to get version")
		}
		// get correct manifestsURL based on version
		manifestsURL, err := utils.GetManifestTreeURL(version)
		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}
		// pick all the manifest files stored in manifestsURL
		manifests, err := utils.ListManifests(manifestsURL)

		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}

		log.Info("Stopping Meshery...")

		// delete the Meshery deployment using the manifest files to stop Meshery
		err = utils.ApplyManifestFiles(manifests, RequestedAdapters, client, false, true)
		if err != nil {
			return ErrApplyManifest(err, false, true)
		}

		err = utils.ApplyOperatorManifest(client, false, true)
		if err != nil {
			return ErrApplyOperatorManifest(err, false, true)
		}

		s := utils.CreateDefaultSpinner("Terminating Meshery pods", "\nPods terminated.")
		s.Start()

		deadline := time.Now().Add(20 * time.Second)

		for !(time.Now().After(deadline)) {
			ok, err := utils.IsMesheryRunning("kubernetes")

			if err != nil {
				return err
			}

			if !ok {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}
		s.Stop()
	}

	log.Info("Meshery is stopped.")

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}
	return nil
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
