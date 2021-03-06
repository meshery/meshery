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

	meshkitutils "github.com/layer5io/meshkit/utils"
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

// deleteManifest is a wrapper function for client.ApplyManifest
// to delete deployments with the specified manifest file
func deleteManifest(manifest []byte, client *meshkitkube.Client) error {
	// ApplyManifest applies/deletes the given manifest file to/from the Kubernetes cluster
	err := client.ApplyManifest(manifest, meshkitkube.ApplyOptions{
		Namespace: utils.MesheryNamespace,
		Update:    false,
		Delete:    true,
	})

	if err != nil {
		return errors.Wrap(err, "failed to delete deployments")
	}
	return nil
}

// DeleteManifestFiles deletes all the specified deployments in the meshery namespace
func DeleteManifestFiles(manifestArr []utils.Manifest, requestedAdapters []string, client *meshkitkube.Client) error {
	// TODO: @navendu-pottekkat Merge the ApplyManifestFiles and DeleteManifestFiles to a single function
	// path to the manifest files ~/.meshery/manifests
	manifestFiles := filepath.Join(utils.MesheryFolder, utils.ManifestsFolder)

	// read the manifest files as strings
	// other than the adapters, meshery-deployment.yaml, meshery-service.yaml and service-account.yaml should be deleted
	MesheryDeploymentManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, utils.MesheryDeployment))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	mesheryServiceManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, utils.MesheryService))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	serviceAccountManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, utils.ServiceAccount))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}

	// delete the manifest files
	if err = deleteManifest([]byte(MesheryDeploymentManifest), client); err != nil {
		return errors.Wrap(err, "failed to stop Meshery")
	}
	if err = deleteManifest([]byte(mesheryServiceManifest), client); err != nil {
		return errors.Wrap(err, "failed to stop Meshery")
	}
	if err = deleteManifest([]byte(serviceAccountManifest), client); err != nil {
		return errors.Wrap(err, "failed to stop Meshery")
	}

	// loop through the required adapters as specified in the config.yaml file and delete each
	for _, adapter := range requestedAdapters {
		// for each adapter, there is a meshery-adapterName-deployment.yaml and meshery-adapterName-service.yaml
		// manifest file. See- https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		adapterFile := filepath.Join(manifestFiles, adapter)
		adapterDeployment := adapterFile + "-deployment.yaml"
		adapterService := adapterFile + "-service.yaml"

		if utils.IsAdapterValid(manifestArr, adapter+"-deployment.yaml") == false {
			return fmt.Errorf("invalid adapter %s specified. Please check %s/config.yaml file", adapter, utils.MesheryFolder)
		}

		// read manifest files as strings and delete
		manifestDepl, err := meshkitutils.ReadLocalFile(adapterDeployment)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}
		manifestService, err := meshkitutils.ReadLocalFile(adapterService)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}

		if err = deleteManifest([]byte(manifestDepl), client); err != nil {
			return errors.Wrap(err, "failed to stop Meshery")
		}
		if err = deleteManifest([]byte(manifestService), client); err != nil {
			return errors.Wrap(err, "failed to stop Meshery")
		}
	}
	log.Debug("deleted the Meshery deployments.")

	return nil
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
	currPlatform := currCtx.Platform
	RequestedAdapters := currCtx.Adapters

	switch currPlatform {
	case "docker":
		// if the platform is docker, then stop all the running containers
		if !utils.IsMesheryRunning() {
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

		// ask user for confirmation
		userResponse := utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
		if !userResponse {
			log.Info("Stop aborted.")
			return nil
		}

		// create an kubernetes client
		client, err := utils.CreateKubeClient()

		if err != nil {
			return err
		}

		// check if the manifest folder exists on the machine
		if _, err := os.Stat(filepath.Join(utils.MesheryFolder, utils.ManifestsFolder)); os.IsNotExist(err) {
			log.Errorf("%s folder does not exist.", utils.ManifestsFolder)
			return err
		}

		// pick all the manifest files stored in https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		manifests, err := utils.ListManifests(manifestsURL)

		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}

		log.Info("Stopping Meshery...")

		// delete the Meshery deployment using the manifest files to stop Meshery
		err = DeleteManifestFiles(manifests, RequestedAdapters, client)

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
