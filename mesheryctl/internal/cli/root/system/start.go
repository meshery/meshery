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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
)

var (
	skipUpdateFlag bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery",
	Long:  `Start Meshery and each of its service mesh adapters.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		err := RunPreflightHealthChecks(true, cmd.Use)
		if err != nil {
			cmd.SilenceUsage = true

		}

		return err

	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := start(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to start Meshery"))

		}
		return nil

	},
}

func start() error {
	if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
		if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to make %s directory", utils.MesheryFolder)))
		}
	}

	kubeClient, err := meshkitkube.New([]byte(""))
	if err != nil {
		return err
	}

	err = utils.CreateManifestsFolder()

	if err != nil {
		return err
	}

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

	if utils.PlatformFlag != "" {
		currCtx.Platform = utils.PlatformFlag
		err := utils.ChangePlatform(mctlCfg.CurrentContext, currCtx)

		if err != nil {
			return err
		}
	}

	currPlatform := currCtx.Platform
	RequestedAdapters := currCtx.Adapters // Requested Adapters / Services

	// Deploy to platform specified in the config.yaml
	switch currPlatform {
	case "docker":

		// download the docker-compose.yaml file corresponding to the current version
		if err := utils.DownloadDockerComposeFile(currCtx, true); err != nil {
			return errors.Wrapf(err, utils.SystemError(fmt.Sprintf("failed to download %s file", utils.DockerComposeFile)))
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

		//changing the port mapping in docker compose
		services := compose.Services // Current Services
		//extracting the custom user port from config.yaml
		userPort := strings.Split(currCtx.Endpoint, ":")
		//extracting container port from the docker-compose
		containerPort := strings.Split(services["meshery"].Ports[0], ":")
		userPortMapping := userPort[len(userPort)-1] + ":" + containerPort[len(containerPort)-1]
		services["meshery"].Ports[0] = userPortMapping

		RequiredService := []string{"meshery", "watchtower"}

		AllowedServices := map[string]utils.Service{}
		for _, v := range RequestedAdapters {
			if services[v].Image == "" {
				log.Fatalf("Invalid adapter specified %s", v)
			}

			temp, ok := services[v]
			if !ok {
				return errors.New("unable to extract adapter version")
			}

			spliter := strings.Split(temp.Image, ":")
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.Channel, "latest")
			services[v] = temp
			AllowedServices[v] = services[v]
		}

		for _, v := range RequiredService {
			if v == "watchtower" {
				AllowedServices[v] = services[v]
				continue
			}

			temp, ok := services[v]
			if !ok {
				return errors.New("unable to extract meshery version")
			}

			spliter := strings.Split(temp.Image, ":")
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.Channel, "latest")
			if v == "meshery" {
				temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.Channel, currCtx.Version)
			}
			services[v] = temp
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
				return errors.Wrap(err, utils.SystemError("failed to update Meshery containers"))
			}
		}

		// Reset Meshery config file to default settings
		if utils.ResetFlag {
			err := resetMesheryConfig()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to reset meshery config"))
			}
		}

		var endpoint meshkitutils.HostPort

		userResponse := false

		//skip asking confirmation if -y flag used or host in meshconfig is already localhost
		if utils.SilentFlag || strings.HasSuffix(userPort[1], "localhost") {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("The endpoint address will be changed to localhost. Are you sure you want to continue?")
		}

		if userResponse {
			endpoint.Address = utils.EndpointProtocol + "://localhost"
			currCtx.Endpoint = endpoint.Address + ":" + userPort[len(userPort)-1]

			err = utils.ChangeConfigEndpoint(mctlCfg.CurrentContext, currCtx)
			if err != nil {
				return err
			}
		} else {
			endpoint.Address = userPort[0]
		}

		tempPort, err := strconv.Atoi(userPort[len(userPort)-1])
		if err != nil {
			return err
		}
		endpoint.Port = int32(tempPort)

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

		var mockEndpoint *meshkitutils.MockOptions
		mockEndpoint = nil

		res := meshkitutils.TcpCheck(&endpoint, mockEndpoint)
		if res {
			return errors.New("the endpoint is not accessible")
		}

		//check for container meshery_meshery_1 running status
		for _, container := range containers {
			if container.Names[0] == "/meshery_meshery_1" {
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

		version := currCtx.Version
		channel := currCtx.Channel
		if version == "latest" {
			if channel == "edge" {
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
		// path to the manifest files ~/.meshery/manifests
		manifestFiles := filepath.Join(utils.MesheryFolder, utils.ManifestsFolder)

		// change version in meshery-deployment manifest
		err = utils.ChangeManifestVersion(utils.MesheryDeployment, channel, version, filepath.Join(manifestFiles, utils.MesheryDeployment))
		if err != nil {
			return err
		}

		// downloaded required files successfully now apply the manifest files
		log.Info("Starting Meshery...")

		// apply the adapters mentioned in the config.yaml file to the Kubernetes cluster
		err = utils.ApplyManifestFiles(manifests, RequestedAdapters, kubeClient, false, false)

		if err != nil {
			return err
		}

		log.Info("...Meshery deployed on Kubernetes.")

		clientset := kubeClient.KubeClient

		var opts meshkitkube.ServiceOptions
		opts.Name = "meshery"
		opts.Namespace = utils.MesheryNamespace
		opts.APIServerURL = kubeClient.RestConfig.Host

		var endpoint *meshkitutils.Endpoint

		deadline := time.Now().Add(3 * time.Second)

		//polling for endpoint to be available within the three second deadline
		for !(time.Now().After(deadline)) {
			endpoint, err = meshkitkube.GetServiceEndpoint(context.TODO(), clientset, &opts)
			if err == nil {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}

		if err != nil {
			log.Info("service endpoint not discovered. Meshery deployed on default endpoint")
		} else {
			currCtx.Endpoint = utils.EndpointProtocol + "://" + endpoint.External.Address + ":" + strconv.Itoa(int(endpoint.External.Port))

			err = utils.ChangeConfigEndpoint(mctlCfg.CurrentContext, currCtx)
			if err != nil {
				return err
			}
		}

		// switch to default case if the platform specified is not supported
	default:
		return errors.New(fmt.Sprintf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes\nPlease check %s/config.yaml file.", currPlatform, utils.MesheryFolder))
	}

	err = utils.DownloadOperatorManifest()

	if err != nil {
		return err
	}

	if !skipUpdateFlag {
		err = utils.ApplyOperatorManifest(kubeClient, true, false)

		if err != nil {
			return err
		}
	} else {
		// skip applying update on operators when the flag is used
		err = utils.ApplyOperatorManifest(kubeClient, false, false)

		if err != nil {
			return err
		}
	}

	log.Info("Opening Meshery in your browser. If Meshery does not open, please point your browser to " + currCtx.Endpoint + " to access Meshery.")

	err = utils.NavigateToBrowser(currCtx.Endpoint)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	startCmd.PersistentFlags().StringVarP(&utils.PlatformFlag, "platform", "p", "", "platform to deploy Meshery to.")
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
