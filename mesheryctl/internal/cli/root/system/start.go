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
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/spf13/viper"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"k8s.io/client-go/kubernetes"
)

var (
	skipUpdateFlag bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery",
	Long:  `Run 'docker-compose' to start Meshery and each of its service mesh adapters.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		return utils.PreReqCheck(cmd.Use)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := start(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to start Meshery"))
		}
		return nil
	},
}

// ValidateComposeFileForRecreation validates the docker-compose.yaml file
func ValidateComposeFileForRecreation(CurrentServices map[string]utils.Service, RequestedServices []string) error {
	valid := true
	for _, v := range RequestedServices {
		_, ok := CurrentServices[v]
		if !ok {
			valid = false
			break
		}
	}
	if !valid {
		if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file from %s", utils.DockerComposeFile, fileURL)))
		}
		err := utils.ViperCompose.ReadInConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

// applyManifest is a wrapper function for client.ApplyManifest
func applyManifest(manifest []byte, client *meshkitkube.Client) error {
	// ApplyManifest applies the given manifest file to the Kubernetes cluster
	err := client.ApplyManifest(manifest, meshkitkube.ApplyOptions{
		Namespace: "meshery",
		Update:    true,
		Delete:    false,
	})

	if err != nil {
		return errors.Wrap(err, "failed to apply manifests")
	}
	return nil
}

// ApplyManifestFiles applies all the required manifests into the Kubernetes cluster
func ApplyManifestFiles(requestedAdapters []string, client *meshkitkube.Client) error {
	// path to the manifest files ~/.meshery/manifests
	manifestFiles := filepath.Join(utils.MesheryFolder, utils.ManifestsFolder)

	// read the manifest files as strings
	// other than the adapters, meshery-deployment.yaml, meshery-service.yaml and service-account.yaml should be applied
	MesheryDeploymentManifest, err := ioutil.ReadFile(filepath.Join(manifestFiles, utils.MesheryDeployment))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	mesheryServiceManifest, err := ioutil.ReadFile(filepath.Join(manifestFiles, utils.MesheryService))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	serviceAccountManifest, err := ioutil.ReadFile(filepath.Join(manifestFiles, utils.ServiceAccount))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}

	// apply the manifest files
	if err = applyManifest(MesheryDeploymentManifest, client); err != nil {
		return errors.Wrap(err, "failed to apply manifest files")
	}
	if err = applyManifest(mesheryServiceManifest, client); err != nil {
		return errors.Wrap(err, "failed to apply manifest files")
	}
	if err = applyManifest(serviceAccountManifest, client); err != nil {
		return errors.Wrap(err, "failed to apply manifest files")
	}

	// loop through the required adapters as specified in the config.yaml file and apply each
	for _, adapter := range requestedAdapters {
		// for each adapter, there is a meshery-adapterName-deployment.yaml and meshery-adapterName-service.yaml
		// manifest file. See- https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		adapterFile := filepath.Join(manifestFiles, adapter)
		adapterDeployment := adapterFile + "-deployment.yaml"
		adapterService := adapterFile + "-service.yaml"

		// read manifest files as strings and apply
		manifestDepl, err := ioutil.ReadFile(adapterDeployment)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}
		manifestService, err := ioutil.ReadFile(adapterService)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}

		if err = applyManifest(manifestDepl, client); err != nil {
			return errors.Wrap(err, "failed to apply manifest files")
		}
		if err = applyManifest(manifestService, client); err != nil {
			return errors.Wrap(err, "failed to apply manifest files")
		}
	}
	log.Debug("applied manifests to the Kubernetes cluster.")

	return nil
}

func start() error {
	if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
		if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to make %s directory", utils.MesheryFolder)))
		}
	}

	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// TODO: Centralize docker-compose.yaml file pull for this and reset
	// get the platform, channel and the version of the current context
	currPlatform := mctlCfg.GetContextContent().Platform
	RequestedAdapters := mctlCfg.GetContextContent().Adapters // Requested Adapters / Services
	// TODO: Context support for system start
	// currChannel := mctlCfg.Contexts[currentContext].Channel
	// currVersion := mctlCfg.Contexts[currentContext].Version
	// fileURL := ""

	// Deploy to platform specified in the config.yaml
	switch currPlatform {
	case "docker":

		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := utils.DownloadFile(utils.DockerComposeFile, fileURL); err != nil {
				return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file from %s", utils.DockerComposeFile, fileURL)))
			}
		}

		// Viper instance used for docker compose
		utils.ViperCompose.SetConfigFile(utils.DockerComposeFile)
		err = utils.ViperCompose.ReadInConfig()
		if err != nil {
			return err
		}

		compose := &utils.DockerCompose{}
		err = utils.ViperCompose.Unmarshal(&compose)
		if err != nil {
			return err
		}

		services := compose.Services                                        // Current Services
		err = ValidateComposeFileForRecreation(services, RequestedAdapters) // Validates docker-compose file and recreates and sync's if required
		if err != nil {
			return err
		}

		err = utils.ViperCompose.Unmarshal(&compose)
		if err != nil {
			return err
		}
		services = compose.Services // Current Services
		RequiredService := []string{"meshery", "watchtower"}

		AllowedServices := map[string]utils.Service{}
		for _, v := range RequestedAdapters {
			if services[v].Image == "" {
				log.Fatalf("Invalid adapter specified %s", v)
			}
			AllowedServices[v] = services[v]
		}
		for _, v := range RequiredService {
			AllowedServices[v] = services[v]
		}

		utils.ViperCompose.Set("services", AllowedServices)
		err = utils.ViperCompose.WriteConfig()
		if err != nil {
			return err
		}

		//fmt.Println("Services", services);
		//fmt.Println("RequiredAdapters", RequestedAdapters);
		//fmt.Println("AllowedServices", AllowedServices);
		//fmt.Println("version", utils.ViperCompose.GetString("version")) // Works here

		//////// FLAGS
		// Control whether to pull for new Meshery container images
		if skipUpdateFlag {
			log.Info("Skipping Meshery update...")
		} else {
			err := utils.UpdateMesheryContainers()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to update meshery containers"))
			}
		}

		// Reset Meshery config file to default settings
		if utils.ResetFlag {
			err := resetMesheryConfig()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to reset meshery config"))
			}
		}

		log.Info("Starting Meshery...")
		start := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "up", "-d")
		start.Stdout = os.Stdout
		start.Stderr = os.Stderr

		if err := start.Run(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to run meshery server"))
		}
		checkFlag := 0 //flag to check

		//connection to docker-client
		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to create new env client"))
		}

		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
		if err != nil {
			return errors.Wrap(err, utils.SystemError("failed to fetch the list of containers"))
		}

		//check for container meshery_meshery_1 running status
		for _, container := range containers {
			if container.Names[0] == "/meshery_meshery_1" {
				log.Info("Opening Meshery in your browser. If Meshery does not open, please point your browser to http://localhost:9081 to access Meshery.")

				//check for os of host machine
				if runtime.GOOS == "windows" {
					// Meshery running on Windows host
					err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
					if err != nil {
						return errors.Wrap(err, utils.SystemError("failed to exec command"))
					}
				} else if runtime.GOOS == "linux" {
					// Meshery running on Linux host
					_, err = exec.LookPath("xdg-open")
					if err != nil {
						break
					}
					err = exec.Command("xdg-open", url).Start()
					if err != nil {
						return errors.Wrap(err, utils.SystemError("failed to exec command"))
					}
				} else {
					// Assume Meshery running on MacOS host
					err = exec.Command("open", url).Start()
					if err != nil {
						return errors.Wrap(err, utils.SystemError("failed to exec command"))
					}
				}

				//check flag to check successful deployment
				checkFlag = 0
				break
			} else {
				checkFlag = 1
			}
		}

		//if meshery_meshery_1 failed to start showing logs
		//code for logs
		if checkFlag == 1 {
			log.Info("Starting Meshery logging . . .")
			cmdlog := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "logs", "-f")
			cmdReader, err := cmdlog.StdoutPipe()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to create stdout pipe"))
			}
			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					log.Println(scanner.Text())
				}
			}()
			if err := cmdlog.Start(); err != nil {
				return errors.Wrap(err, utils.SystemError("failed to start logging"))
			}
			if err := cmdlog.Wait(); err != nil {
				return errors.Wrap(err, utils.SystemError("failed to wait for command to execute"))
			}
		}

	case "kubernetes":

		log.Debug("detecting kubeconfig file...")

		// Detect user's kubeconfig file for the Kubernetes cluster
		config, err := meshkitkube.DetectKubeConfig()

		if err != nil {
			return errors.Wrap(err, "failed to detect Kube Config file")
		}

		log.Debug("creating new Clientset...")
		// Create a new Clientset for given config
		clientSet, err := kubernetes.NewForConfig(config)

		if err != nil {
			return errors.Wrap(err, "error setting clientset")
		}

		// Create a new client
		client, err := meshkitkube.New(clientSet, *config)

		if err != nil {
			return errors.Wrap(err, "failed to create new client")
		}

		log.Debug("fetching required Kubernetes manifest files...")
		// pick all the manifest files stored in https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		manifests, err := utils.ListManifests(manifestsURL)

		if err != nil {
			return errors.Wrap(err, "failed to make GET request")
		}

		log.Debug("creating ~/.meshery/manifests folder...")
		// create a manifests folder under ~/.meshery to store the manifest files
		if _, err := os.Stat(utils.ManifestsFolder); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Join(utils.MesheryFolder, utils.ManifestsFolder), os.ModePerm); err != nil {
				return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to make %s directory", utils.ManifestsFolder)))
			}
			log.Debug("created manifests folder...")
		}

		log.Debugf("downloading manifest files from %s", gitHubFolder)

		// download all the manifest files to the ~/.meshery/manifests folder
		err = utils.DownloadManifests(manifests, rawManifestsURL)

		if err != nil {
			return errors.Wrap(err, "failed to download manifests")
		}

		// downloaded required files successfully now apply the manifest files
		log.Info("Starting Meshery...")

		log.Debug("applying the manifests to Kubernetes cluster...")

		// apply the adapters mentioned in the config.yaml file to the Kubernetes cluster
		err = ApplyManifestFiles(RequestedAdapters, client)

		if err != nil {
			return errors.Wrap(err, "failed to apply manifests")
		}
	}

	return nil
}

func init() {
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
	startCmd.Flags().BoolVarP(&utils.SilentFlag, "silent", "", false, "(optional) silently create Meshery's configuration file with default settings.")
}
