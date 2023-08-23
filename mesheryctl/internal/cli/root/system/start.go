// Copyright 2023 Layer5, Inc.
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
	"path"
	"strconv"
	"strings"
	"time"

	c "github.com/layer5io/meshery/mesheryctl/pkg/constants"
	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	dockerCmd "github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types"
	dockerconfig "github.com/docker/docker/cli/config"

	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipUpdateFlag  bool
	skipBrowserFlag bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Meshery",
	Long:  `Start Meshery and each of its cloud native components.`,
	Args:  cobra.NoArgs,
	Example: `
// Start meshery
mesheryctl system start

// (optional) skip opening of MesheryUI in browser.
mesheryctl system start --skip-browser

// (optional) skip checking for new updates available in Meshery.
mesheryctl system start --skip-update

// Reset Meshery's configuration file to default settings.
mesheryctl system start --reset

// Specify Platform to deploy Meshery to.
mesheryctl system start -p docker
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
		latestVersions, err := meshkitutils.GetLatestReleaseTagsSorted(c.GetMesheryGitHubOrg(), c.GetMesheryGitHubRepo())
		version := constants.GetMesheryctlVersion()
		if err == nil {
			if len(latestVersions) == 0 {
				log.Warn("no versions found for Meshery")
				return
			}
			latest := latestVersions[len(latestVersions)-1]
			if latest != version {
				log.Printf("A new release of mesheryctl is available: %s → %s", version, latest)
				log.Printf("https://github.com/layer5io/meshery/releases/tag/%s", latest)
				log.Print("Check https://docs.meshery.io/guides/upgrade#upgrading-meshery-cli for instructions on how to update mesheryctl\n")
			}
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

			// update the context to config
			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
			if err != nil {
				return err
			}
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
			return ErrUnmarshalDockerCompose(err, utils.DockerComposeFile)
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
		for _, v := range currCtx.GetComponents() {
			if services[v].Image == "" {
				log.Fatalf("Invalid component specified %s", v)
			}

			temp, ok := services[v]
			if !ok {
				return errors.New("unable to extract component version")
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

				if currCtx.GetProvider() != "" {
					temp.Environment = append(temp.Environment, fmt.Sprintf("%s=%s", "PROVIDER", currCtx.GetProvider()))
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

			err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
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

		// Get the Docker configuration
		dockerCfg, err := cliconfig.Load(dockerconfig.Dir())
		if err != nil {
			return ErrCreatingDockerClient(err)
		}

		//connection to docker-client
		cli, err := dockerCmd.NewAPIClientFromFlags(cliflags.NewCommonOptions(), dockerCfg)
		if err != nil {
			return ErrCreatingDockerClient(err)
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
			}

			checkFlag = 1
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

		// Applying Meshery Helm charts for installing Meshery
		if err = applyHelmCharts(kubeClient, currCtx, mesheryImageVersion, false, meshkitkube.INSTALL); err != nil {
			return err
		}

		// checking if Meshery is ready
		time.Sleep(10 * time.Second) // sleeping 10 seconds to countermeasure time to apply helm charts
		ready, err := mesheryReadinessHealthCheck()
		if err != nil {
			log.Info(err)
		}

		spinner.Stop()

		if !ready {
			log.Info("\nFew Meshery pods have not come up yet.\nPlease check the status of the pods by executing “mesheryctl system status” and Meshery-UI endpoint with “mesheryctl system dashboard” before using meshery.")
			return nil
		}
		log.Info("Meshery is starting...")

		// switch to default case if the platform specified is not supported
	default:
		return fmt.Errorf("the platform %s is not supported currently. The supported platforms are:\ndocker\nkubernetes\nPlease check %s/config.yaml file", currCtx.GetPlatform(), utils.MesheryFolder)
	}

	// execute dashboard command to fetch and navigate to Meshery UI
	return dashboardCmd.RunE(nil, nil)
}

func init() {
	startCmd.PersistentFlags().StringVarP(&utils.PlatformFlag, "platform", "p", "", "platform to deploy Meshery to.")
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
	startCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
}

// Apply Meshery helm charts
func applyHelmCharts(kubeClient *meshkitkube.Client, currCtx *config.Context, mesheryImageVersion string, dryRun bool, act meshkitkube.HelmChartAction) error {
	// get value overrides to install the helm chart
	overrideValues := utils.SetOverrideValues(currCtx, mesheryImageVersion)

	// install the helm charts with specified override values
	var chartVersion string
	if mesheryImageVersion != "latest" {
		chartVersion = mesheryImageVersion
	}
	action := "install"
	if act == meshkitkube.UNINSTALL {
		action = "uninstall"
	}
	errServer := kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
		Namespace:       utils.MesheryNamespace,
		ReleaseName:     "meshery",
		CreateNamespace: true,
		ChartLocation: meshkitkube.HelmChartLocation{
			Repository: utils.HelmChartURL,
			Chart:      utils.HelmChartName,
			Version:    chartVersion,
		},
		OverrideValues: overrideValues,
		Action:         act,
		// the helm chart will be downloaded to ~/.meshery/manifests if it doesn't exist
		DownloadLocation: path.Join(utils.MesheryFolder, utils.ManifestsFolder),
		DryRun:           dryRun,
	})
	errOperator := kubeClient.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
		Namespace:       utils.MesheryNamespace,
		ReleaseName:     "meshery-operator",
		CreateNamespace: true,
		ChartLocation: meshkitkube.HelmChartLocation{
			Repository: utils.HelmChartURL,
			Chart:      utils.HelmChartOperatorName,
			Version:    chartVersion,
		},
		Action: act,
		// the helm chart will be downloaded to ~/.meshery/manifests if it doesn't exist
		DownloadLocation: path.Join(utils.MesheryFolder, utils.ManifestsFolder),
		DryRun:           dryRun,
	})
	if errServer != nil && errOperator != nil {
		return fmt.Errorf("could not %s meshery server: %s\ncould not %s meshery-operator: %s", action, errServer.Error(), action, errOperator.Error())
	}
	if errServer != nil {
		return fmt.Errorf("%s success for operator but failed for meshery server: %s", action, errServer.Error())
	}
	if errOperator != nil {
		return fmt.Errorf("%s success for meshery server but failed for meshery operator: %s", action, errOperator.Error())
	}
	return nil
}
