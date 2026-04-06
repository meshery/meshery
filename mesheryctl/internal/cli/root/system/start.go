// Copyright Meshery Authors
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
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/constants"
	pkgconstants "github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	meshkitutils "github.com/meshery/meshkit/utils"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
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

// Specify Provider to use.
mesheryctl system start --provider Meshery
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
			utils.Log.Error(err)
			return nil
		}
		ctx, err := cfg.GetCurrentContext()
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}
		err = ctx.ValidateVersion()
		if err != nil {
			utils.Log.Error(err)
			return nil
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
		utils.CheckMesheryctlClientVersion(constants.GetMesheryctlVersion())
	},
}

func start() error {
	if err := ensureMesheryFolder(); err != nil {
		return err
	}
	utils.Log.Debug("Meshery folder exists. Preparing start context...")
	mctlCfg, currCtx, mesheryImageVersion, callbackURL, providerURL, err := prepareStartContext()
	if err != nil {
		return err
	}

	utils.Log.Debugf("Starting Meshery using platform: %s", currCtx.GetPlatform())
	switch currCtx.GetPlatform() {
	case platformDocker:
		if err := startDockerDeployment(mctlCfg, currCtx, mesheryImageVersion, callbackURL); err != nil {
			return err
		}
	case platformKubernetes:
		if err := startKubernetesDeployment(currCtx, mesheryImageVersion, callbackURL, providerURL); err != nil {
			return err
		}
	default:
		return unsupportedPlatformError(currCtx.GetPlatform())
	}

	if skipBrowserFlag {
		utils.Log.Info("Meshery deployed. Use `mesheryctl system dashboard` to access the UI.")
		return nil
	}
	if err := dashboardCmd.RunE(nil, nil); err != nil {
		utils.Log.Warnf("Meshery deployed successfully, but could not open dashboard: %v\nRun `mesheryctl system dashboard` to access the UI.", err)
	}
	return nil
}

func ensureMesheryFolder() error {
	if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
		if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
			return ErrCreateDir(err, utils.MesheryFolder)
		}
	}

	return nil
}

func prepareStartContext() (*config.MesheryCtlConfig, *config.Context, string, string, string, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, nil, "", "", "", errors.Wrap(err, "error processing config")
	}

	if tempContext != "" {
		if err := mctlCfg.SetCurrentContext(tempContext); err != nil {
			return nil, nil, "", "", "", errors.Wrap(err, "failed to set temporary context")
		}
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return nil, nil, "", "", "", err
	}

	if utils.PlatformFlag != "" {
		switch utils.PlatformFlag {
		case platformDocker, platformKubernetes:
			currCtx.SetPlatform(utils.PlatformFlag)
		default:
			return nil, nil, "", "", "", ErrUnsupportedPlatform(utils.PlatformFlag, utils.CfgFile)
		}
	}

	if providerFlag != "" {
		currCtx.SetProvider(providerFlag)
	}

	if err := config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName()); err != nil {
		return nil, nil, "", "", "", err
	}

	if utils.ResetFlag {
		if err := resetMesheryConfig(); err != nil {
			return nil, nil, "", "", "", ErrResetMeshconfig(err)
		}
	}

	callbackURL := viper.GetString(pkgconstants.CallbackURLENV)
	providerURL := viper.GetString(pkgconstants.ProviderURLsENV)

	return mctlCfg, currCtx, resolveMesheryImageVersion(currCtx), callbackURL, providerURL, nil
}

func resolveMesheryImageVersion(currCtx *config.Context) string {
	if currCtx.GetChannel() == "stable" && currCtx.GetVersion() == "latest" {
		return "latest"
	}

	return currCtx.GetVersion()
}

func startDockerDeployment(mctlCfg *config.MesheryCtlConfig, currCtx *config.Context, mesheryImageVersion, callbackURL string) error {
	if err := utils.DownloadDockerComposeFile(currCtx, true); err != nil {
		return utils.ErrDownloadFile(err, utils.DockerComposeFile)
	}

	utils.ViperCompose.SetConfigFile(utils.DockerComposeFile)
	if err := utils.ViperCompose.ReadInConfig(); err != nil {
		return err
	}

	compose := &utils.DockerCompose{}
	if err := utils.ViperCompose.Unmarshal(&compose); err != nil {
		return ErrUnmarshalDockerCompose(err, utils.DockerComposeFile)
	}

	userPort, err := configureDockerServices(currCtx, mesheryImageVersion, callbackURL)
	if err != nil {
		return err
	}

	if skipUpdateFlag {
		log.Info("Skipping Meshery update...")
	} else if err := utils.UpdateMesheryContainers(); err != nil {
		return errors.Wrap(err, utils.SystemError("failed to update Meshery containers"))
	}

	endpoint, err := resolveDockerEndpoint(mctlCfg, currCtx, userPort)
	if err != nil {
		return err
	}

	composeClient, err := utils.NewComposeClient()
	if err != nil {
		return errors.Wrap(err, utils.SystemError("failed to create compose client"))
	}

	if err := startDockerCompose(composeClient); err != nil {
		return err
	}

	if err := waitForMesheryContainer(composeClient); err != nil {
		return err
	}

	waitForMesheryEndpoint(endpoint)
	return nil
}

func configureDockerServices(currCtx *config.Context, mesheryImageVersion, callbackURL string) ([]string, error) {
	userPort := strings.Split(currCtx.GetEndpoint(), ":")
	containerPort := strings.Split(utils.Services["meshery"].Ports[0], ":")
	userPortMapping := userPort[len(userPort)-1] + ":" + containerPort[len(containerPort)-1]
	utils.Services["meshery"].Ports[0] = userPortMapping

	allowedServices, err := buildAllowedDockerServices(currCtx, mesheryImageVersion, callbackURL)
	if err != nil {
		return nil, err
	}

	for name, service := range allowedServices {
		utils.ViperCompose.Set(fmt.Sprintf("services.%s", name), service)
		if err := utils.ViperCompose.WriteConfig(); err != nil {
			log.Errorf("Encountered an error while adding `%s` service to Docker Compose file. Verify permission to write to `.meshery/meshery.yaml` file.", name)
		}
	}

	return userPort, nil
}

func buildAllowedDockerServices(currCtx *config.Context, mesheryImageVersion, callbackURL string) (map[string]utils.Service, error) {
	requiredServices := []string{"meshery", "watchtower"}
	allowedServices := map[string]utils.Service{}

	for _, component := range currCtx.GetComponents() {
		if utils.Services[component].Image == "" {
			log.Fatalf("Invalid component specified %s", component)
		}

		service, ok := utils.Services[component]
		if !ok {
			return nil, errors.Errorf("No Docker Compose service exists for Meshery component `%s`.", component)
		}

		spliter := strings.Split(service.Image, ":")
		service.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
		utils.Services[component] = service
		allowedServices[component] = utils.Services[component]
	}

	for _, serviceName := range requiredServices {
		if serviceName == "watchtower" {
			allowedServices[serviceName] = utils.Services[serviceName]
			continue
		}

		service, ok := utils.Services[serviceName]
		if !ok {
			return nil, errors.New("unable to extract meshery version")
		}

		spliter := strings.Split(service.Image, ":")
		service.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
		if serviceName == "meshery" {
			service.Environment = upsertServiceEnvVar(service.Environment, pkgconstants.CallbackURLENV, callbackURL)

			providerEnvVar := currCtx.GetProvider()
			if providerFlag != "" {
				providerEnvVar = providerFlag
			}
			service.Environment = upsertServiceEnvVar(service.Environment, pkgconstants.ProviderENV, providerEnvVar)
			service.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), mesheryImageVersion)

			for key, value := range currCtx.GetEnvs() {
				service.Environment = append(service.Environment, fmt.Sprintf("%s=%v", strings.ToUpper(key), value))
			}
		}

		utils.Services[serviceName] = service
		allowedServices[serviceName] = utils.Services[serviceName]
	}

	return allowedServices, nil
}

func upsertServiceEnvVar(environment []string, key, value string) []string {
	idx, found := utils.FindInSlice(key, environment)
	if !found {
		return append(environment, fmt.Sprintf("%s=%s", key, value))
	}
	if value != "" {
		environment[idx] = fmt.Sprintf("%s=%s", key, value)
	}
	return environment
}

func resolveDockerEndpoint(mctlCfg *config.MesheryCtlConfig, currCtx *config.Context, userPort []string) (meshkitutils.HostPort, error) {
	var endpoint meshkitutils.HostPort

	userResponse := false
	if utils.SilentFlag || strings.HasSuffix(userPort[1], "localhost") {
		userResponse = true
	} else {
		userResponse = utils.AskForConfirmation("The endpoint address will be changed to localhost. Are you sure you want to continue?")
	}

	if userResponse {
		endpoint.Address = utils.EndpointProtocol + "://localhost"
		currCtx.SetEndpoint(endpoint.Address + ":" + userPort[len(userPort)-1])
		if err := config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName()); err != nil {
			return endpoint, err
		}
	} else {
		endpoint.Address = userPort[0]
	}

	tempPort, err := strconv.Atoi(userPort[len(userPort)-1])
	if err != nil {
		return endpoint, err
	}
	endpoint.Port = int32(tempPort)

	return endpoint, nil
}

func startDockerCompose(composeClient *utils.ComposeClient) error {
	spinner := utils.CreateDefaultSpinner("Starting Meshery...", "\nMeshery containers started.")
	spinner.Start()

	if err := composeClient.Up(context.Background(), utils.DockerComposeFile); err != nil {
		spinner.Stop()
		return errors.Wrap(err, utils.SystemError("failed to run Meshery Server"))
	}

	spinner.Stop()
	return nil
}

func waitForMesheryContainer(composeClient *utils.ComposeClient) error {
	spinner := utils.CreateDefaultSpinner("Waiting for Meshery containers to be ready...", "\nMeshery containers are ready.")
	spinner.Start()

	mesheryRunning := false
	const maxRetries = 60
	for i := 0; i < maxRetries; i++ {
		time.Sleep(5 * time.Second)

		containers, err := composeClient.Ps(context.Background(), utils.DockerComposeFile)
		if err != nil {
			spinner.Stop()
			return errors.Wrap(err, utils.SystemError("failed to fetch the list of containers"))
		}

		for _, container := range containers {
			if container.Service == "meshery" && container.State == "running" {
				mesheryRunning = true
				break
			}
		}

		if mesheryRunning {
			spinner.Stop()
			return nil
		}

		if i > 0 && i%6 == 0 {
			spinner.Stop()
			log.Infof("Still waiting for Meshery to start... (%d seconds elapsed)", i*5)
			spinner = utils.CreateDefaultSpinner("Waiting for Meshery containers to be ready...", "\nMeshery containers are ready.")
			spinner.Start()
		}
	}

	spinner.Stop()
	log.Info("Timeout waiting for Meshery container to start. Checking container status...")
	containers, err := composeClient.Ps(context.Background(), utils.DockerComposeFile)
	if err != nil {
		return errors.Wrap(err, utils.SystemError("failed to fetch the list of containers"))
	}

	log.Info("Container status:")
	for _, container := range containers {
		log.Infof("  %s: %s", container.Service, container.State)
	}

	log.Info("\nShowing Meshery logs:")
	if err := composeClient.Logs(context.Background(), utils.DockerComposeFile, false, os.Stdout); err != nil {
		return errors.Wrap(err, utils.SystemError("failed to get logs"))
	}

	return errors.New("timeout: Meshery container did not start within the expected time. Please check the logs above for more details")
}

func waitForMesheryEndpoint(endpoint meshkitutils.HostPort) {
	spinner := utils.CreateDefaultSpinner("Verifying Meshery endpoint accessibility...", "\nMeshery endpoint is accessible.")
	spinner.Start()
	defer spinner.Stop()

	const maxEndpointRetries = 30
	for i := 0; i < maxEndpointRetries; i++ {
		time.Sleep(2 * time.Second)
		if meshkitutils.TcpCheck(&endpoint, nil) {
			log.Info("Meshery is now running!")
			return
		}
	}

	log.Warn("Warning: Meshery endpoint is not yet accessible. The server may still be initializing.")
	log.Info("You can check the status later with: mesheryctl system status")
}

func startKubernetesDeployment(currCtx *config.Context, mesheryImageVersion, callbackURL, providerURL string) error {
	kubeClient, err := meshkitkube.New([]byte(""))
	if err != nil {
		return err
	}

	spinner := utils.CreateDefaultSpinner("Deploying Meshery on Kubernetes", "\nMeshery deployed on Kubernetes.")
	spinner.Start()
	defer spinner.Stop()

	if err := utils.CreateManifestsFolder(); err != nil {
		utils.Log.Error(utils.ErrCreateManifestsFolder(err))
		return err
	}

	if err = applyHelmCharts(kubeClient, currCtx, mesheryImageVersion, false, meshkitkube.INSTALL, callbackURL, providerURL); err != nil {
		return err
	}

	time.Sleep(20 * time.Second)
	ready, err := mesheryReadinessHealthCheck()
	if err != nil {
		log.Info(err)
	}

	if !ready {
		log.Info("\nTimeout. Meshery pod(s) is not running, yet.\nCheck status of Meshery pod(s) by executing “mesheryctl system status`. Expose Meshery UI with `mesheryctl system dashboard` as needed.")
		return nil
	}
	log.Info("Meshery is starting...")

	return nil
}

func init() {
	startCmd.PersistentFlags().StringVarP(&utils.PlatformFlag, "platform", "p", "", "platform to deploy Meshery to.")
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
	startCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
	startCmd.PersistentFlags().StringVar(&providerFlag, "provider", "", "(optional) Defaults to the provider specified in the current context")
}

// Apply Meshery helm charts
func applyHelmCharts(kubeClient *meshkitkube.Client, currCtx *config.Context, mesheryImageVersion string, dryRun bool, act meshkitkube.HelmChartAction, callbackURL, providerURL string) error {
	// get value overrides to install the helm chart
	overrideValues := utils.SetOverrideValues(currCtx, mesheryImageVersion, callbackURL, providerURL)

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
		return fmt.Errorf("could not %s Meshery Server: %s\ncould not %s meshery-operator: %s", action, errServer.Error(), action, errOperator.Error())
	}
	if errServer != nil {
		return fmt.Errorf("%s success for Meshery Operator, but failed for Meshery Server: %s", action, errServer.Error())
	}
	if errOperator != nil {
		return fmt.Errorf("%s success for Meshery Server, but failed for Meshery Operator: %s", action, errOperator.Error())
	}
	return nil
}
