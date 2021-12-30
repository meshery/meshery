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
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
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
			return err
		}
		cfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		ctx, err := cfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := start(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to start Meshery"))
		}
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

func start() error {
	if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
		if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
			return ErrCreateDir(err, utils.MesheryFolder)
		}
	}

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
	mesheryImageVersion := currCtx.GetVersion()
	if currCtx.GetChannel() == "stable" && currCtx.GetVersion() == "latest" {
		mesheryImageVersion = "latest"
	}

	if utils.PlatformFlag != "" {
		if utils.PlatformFlag == "docker" || utils.PlatformFlag == "kubernetes" {
			currCtx.SetPlatform(utils.PlatformFlag)
		} else {
			return ErrUnsupportedPlatform(utils.PlatformFlag, utils.CfgFile)
		}
	}

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}

	// deploy to platform specified in the config.yaml
	switch currCtx.GetPlatform() {
	case "docker":
		// download the docker-compose.yaml file corresponding to the current version
		if err := utils.DownloadDockerComposeFile(currCtx, true); err != nil {
			return ErrDownloadFile(err, utils.DockerComposeFile)
		}

		// viper instance used for docker compose
		utils.ViperCompose.SetConfigFile(utils.DockerComposeFile)
		err = utils.ViperCompose.ReadInConfig()
		if err != nil {
			return err
		}

		compose := &utils.DockerCompose{}
		err = utils.ViperCompose.Unmarshal(&compose)
		if err != nil {
			return ErrUnmarshal(err, utils.DockerComposeFile)
		}

		//changing the port mapping in docker compose
		services := compose.Services // Current Services
		//extracting the custom user port from config.yaml
		userPort := strings.Split(currCtx.GetEndpoint(), ":")
		//extracting container port from the docker-compose
		containerPort := strings.Split(services["meshery"].Ports[0], ":")
		userPortMapping := userPort[len(userPort)-1] + ":" + containerPort[len(containerPort)-1]
		services["meshery"].Ports[0] = userPortMapping

		RequiredService := []string{"meshery", "watchtower"}

		AllowedServices := map[string]utils.Service{}
		for _, v := range currCtx.GetAdapters() {
			if services[v].Image == "" {
				log.Fatalf("Invalid adapter specified %s", v)
			}

			temp, ok := services[v]
			if !ok {
				return errors.New("unable to extract adapter version")
			}

			spliter := strings.Split(temp.Image, ":")
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
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
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
			if v == "meshery" {
				if !utils.ContainsStringPrefix(temp.Environment, "MESHERY_SERVER_CALLBACK_URL") {
					temp.Environment = append(temp.Environment, fmt.Sprintf("%s=%s", "MESHERY_SERVER_CALLBACK_URL", viper.GetString("MESHERY_SERVER_CALLBACK_URL")))
				}

				temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), mesheryImageVersion)
			}
			services[v] = temp
			AllowedServices[v] = services[v]
		}

		utils.ViperCompose.Set("services", AllowedServices)
		err = utils.ViperCompose.WriteConfig()
		if err != nil {
			return err
		}

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
			currCtx.SetEndpoint(endpoint.Address + ":" + userPort[len(userPort)-1])

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
		kubeClient, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		log.Info("Starting Meshery...")
		spinner := utils.CreateDefaultSpinner("Deploying Meshery on Kubernetes", "\nMeshery deployed on Kubernetes.")
		spinner.Start()

		if err := utils.CreateManifestsFolder(); err != nil {
			return err
		}

		// get value overrides to install the helm chart
		overrideValues := utils.SetOverrideValues(currCtx, mesheryImageVersion)

		// install the helm charts with specified override values
		var chartVersion string
		if mesheryImageVersion != "latest" {
			chartVersion = mesheryImageVersion
		}
		if err = kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
			Namespace:       utils.MesheryNamespace,
			CreateNamespace: true,
			ChartLocation: meshkitkube.HelmChartLocation{
				Repository: utils.HelmChartURL,
				Chart:      utils.HelmChartName,
				Version:    chartVersion,
			},
			OverrideValues: overrideValues,
			Action:         meshkitkube.INSTALL,
			// the helm chart will be downloaded to ~/.meshery/manifests if it doesn't exist
			DownloadLocation: path.Join(utils.MesheryFolder, utils.ManifestsFolder),
		}); err != nil {
			return errors.Wrap(err, "cannot start Meshery")
		}

		deadline := time.Now().Add(20 * time.Second)

		// check if all the pods are running
		for !(time.Now().After(deadline)) {
			podsStatus, err := utils.AreAllPodsRunning()
			if err != nil {
				return err
			}

			if podsStatus {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}

		spinner.Stop()

		podsStatus, err := utils.AreAllPodsRunning()
		if !podsStatus {
			log.Info("\nSome Meshery pods have not come up yet.\nPlease check the status of the pods by executing “mesheryctl system status” before using meshery.")
		} else {
			log.Info("Meshery is starting...")
		}

		clientset := kubeClient.KubeClient

		var opts meshkitkube.ServiceOptions
		opts.Name = "meshery"
		opts.Namespace = utils.MesheryNamespace
		opts.APIServerURL = kubeClient.RestConfig.Host

		var endpoint *meshkitutils.Endpoint

		deadline = time.Now().Add(3 * time.Second)

		//polling for endpoint to be available within the three second deadline
		for !(time.Now().After(deadline)) {
			endpoint, err = meshkitkube.GetServiceEndpoint(context.TODO(), clientset, &opts)
			if err == nil {
				break
			} else {
				time.Sleep(1 * time.Second)
			}
		}

		currCtx.SetEndpoint(fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, endpoint.Internal.Address, endpoint.Internal.Port))
		if !meshkitutils.TcpCheck(&meshkitutils.HostPort{
			Address: endpoint.Internal.Address,
			Port:    endpoint.Internal.Port,
		}, nil) {
			currCtx.SetEndpoint(fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, endpoint.External.Address, endpoint.External.Port))
			if !meshkitutils.TcpCheck(&meshkitutils.HostPort{
				Address: endpoint.External.Address,
				Port:    endpoint.External.Port,
			}, nil) {
				u, _ := url.Parse(opts.APIServerURL)
				if meshkitutils.TcpCheck(&meshkitutils.HostPort{
					Address: u.Hostname(),
					Port:    endpoint.External.Port,
				}, nil) {
					currCtx.SetEndpoint(fmt.Sprintf("%s://%s:%d", utils.EndpointProtocol, u.Hostname(), endpoint.External.Port))
				}
			}
		}

		if err == nil {
			err = utils.ChangeConfigEndpoint(mctlCfg.CurrentContext, currCtx)
			if err != nil {
				return err
			}
		}

		// switch to default case if the platform specified is not supported
	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes\nPlease check %s/config.yaml file", currCtx.GetPlatform(), utils.MesheryFolder)
	}

	// Check for Meshery status before opening it in browser
	client := &http.Client{}
	url := currCtx.GetEndpoint()
	// Can not wait for ever, so setting a limit of 10 seconds
	waittime := time.Now().Add(10 * time.Second)

	for !(time.Now().After(waittime)) {
		// Request to check whether endpoint is up or not
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Info("To open Meshery in browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
			return nil
		}

		resp, err := client.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil || resp.StatusCode != 200 {
			// wait and try accessing the endpoint after one second
			time.Sleep(1 * time.Second)
		} else {
			// meshery server is up and responding, break out of loop and open Meshery in browser
			break
		}
	}

	log.Info("Opening Meshery (" + currCtx.GetEndpoint() + ") in browser.")

	err = utils.NavigateToBrowser(currCtx.GetEndpoint())
	if err != nil {
		log.Warn("Failed to open Meshery in browser, please point your browser to " + currCtx.GetEndpoint() + " to access Meshery.")
	}

	return nil
}

func init() {
	startCmd.PersistentFlags().StringVarP(&utils.PlatformFlag, "platform", "p", "", "platform to deploy Meshery to.")
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
