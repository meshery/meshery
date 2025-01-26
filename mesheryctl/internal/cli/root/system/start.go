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
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
	"net/http"
	"io"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	pkgconstants "github.com/layer5io/meshery/mesheryctl/pkg/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	dockerCmd "github.com/docker/cli/cli/command"
	cliconfig "github.com/docker/cli/cli/config"
	dockerconfig "github.com/docker/cli/cli/config"
	cliflags "github.com/docker/cli/cli/flags"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"

	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipUpdateFlag  bool
	skipBrowserFlag bool
	configureCluster string
)

func getContexts(configFile string) ([]string, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
		return nil, nil
	}

	// GETCONTEXTS endpoint points to the URL return the contexts available
	GETCONTEXTS := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes/contexts"

	req, err := utils.UploadFileWithParams(GETCONTEXTS, nil, utils.ParamName, configFile)
	if err != nil {
		return nil, ErrUploadFileParams(err)
	}

	res, err := utils.MakeRequest(req)
	if err != nil {
		return nil, utils.ErrRequestResponse(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, utils.ErrReadResponseBody(err)
	}

	log.Debugf("Get context API response: %s", string(body))
	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, utils.ErrUnmarshal(err)
	}

	if results == nil {
		return nil, ErrK8sContext(err, configFile)
	}

	var contextNames []string
	for _, ctx := range results {
		ctxname, ok := ctx["name"].(string)
		if !ok {
			errstr := "Invalid context name: context name should be a string"
			return nil, errors.New(errstr)
		}
		contextNames = append(contextNames, ctxname)
	}
	log.Debugf("Available contexts: %s", contextNames)
	return contextNames, nil
}

func setContext(configFile, cname string) error {
	client := &http.Client{}
	extraParams1 := map[string]string{
		"contextName": cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
		return nil
	}

	// SETCONTEXT endpoint points to set context
	SETCONTEXT := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes"
	req, err := utils.UploadFileWithParams(SETCONTEXT, extraParams1, utils.ParamName, configFile)
	if err != nil {
		return ErrUploadFileParams(err)
	}
	res, err := client.Do(req)
	if err != nil {
		return utils.ErrRequestResponse(err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrReadResponseBody(err)
	}
	// TODO: Pretty print the output
	log.Debugf("Set context API response: %s", string(body))
	return nil
}

var aksConfigCmd = &cobra.Command{
	Use:   "aks",
	Short: "Configure Meshery to use AKS cluster",
	Long:  `Configure Meshery to connect to AKS cluster`,
	Example: `
// Configure Meshery to connect to AKS cluster using auth token
mesheryctl system config aks --token auth.json

// Configure Meshery to connect to AKS cluster (if session is logged in using login subcommand)
mesheryctl system config aks
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return errors.New("more than one config name provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		aksCheck := exec.Command("az", "version")
		aksCheck.Stdout = os.Stdout
		aksCheck.Stderr = os.Stderr
		err := aksCheck.Run()
		if err != nil {
			log.Fatalf("Azure CLI not found. Please install Azure CLI and try again. \nSee https://docs.microsoft.com/en-us/cli/azure/install-azure-cli ")
		}
		log.Info("Configuring Meshery to access AKS...")
		var resourceGroup, aksName string

		// Prompt user for Azure resource name
		log.Info("Please enter the Azure resource group name:")
		_, err = fmt.Scanf("%s", &resourceGroup)
		if err != nil {
			log.Warnf("Error reading Azure resource group name: %s", err.Error())
			log.Info("Let's try again. Please enter the Azure resource group name:")
			_, err = fmt.Scanf("%s", &resourceGroup)
			if err != nil {
				log.Fatalf("Error reading Azure resource group name: %s", err.Error())
			}
		}

		// Prompt user for AKS cluster name
		log.Info("Please enter the AKS cluster name:")
		_, err = fmt.Scanf("%s", &aksName)
		if err != nil {
			log.Warnf("Error reading AKS cluster name: %s", err.Error())
			log.Info("Let's try again. Please enter the AKS cluster name:")
			_, err = fmt.Scanf("%s", &aksName)
			if err != nil {
				log.Fatalf("Error reading AKS cluster name: %s", err.Error())
			}
		}

		// Build the Azure CLI syntax to fetch cluster config in kubeconfig.yaml file
		aksCmd := exec.Command("az", "aks", "get-credentials", "--resource-group", resourceGroup, "--name", aksName, "--file", utils.ConfigPath)
		aksCmd.Stdout = os.Stdout
		aksCmd.Stderr = os.Stderr
		// Write AKS compatible config to the filesystem
		err = aksCmd.Run()
		if err != nil {
			log.Fatalf("Error generating kubeconfig: %s", err.Error())
			return err
		}
		log.Debugf("AKS configuration is written to: %s", utils.ConfigPath)

		// set the token in the chosen context
		setToken()
		return nil
	},
}

var eksConfigCmd = &cobra.Command{
	Use:   "eks",
	Short: "Configure Meshery to use EKS cluster",
	Long:  `Configure Meshery to connect to EKS cluster`,
	Example: `
// Configure Meshery to connect to EKS cluster using auth token
mesheryctl system config eks --token auth.json

// Configure Meshery to connect to EKS cluster (if session is logged in using login subcommand)
mesheryctl system config eks
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return errors.New("more than one config name provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		eksCheck := exec.Command("aws", "--version")
		eksCheck.Stdout = os.Stdout
		eksCheck.Stderr = os.Stderr
		err := eksCheck.Run()
		if err != nil {
			log.Fatalf("AWS CLI not found. Please install AWS CLI and try again. \nSee https://docs.aws.amazon.com/cli/latest/reference/ ")
		}
		log.Info("Configuring Meshery to access EKS...")
		var regionName, clusterName string

		// Prompt user for AWS region name
		log.Info("Please enter the AWS region name:")
		_, err = fmt.Scanf("%s", &regionName)
		if err != nil {
			log.Warnf("Error reading AWS region name: %s", err.Error())
			log.Info("Let's try again. Please enter the AWS region name:")
			_, err = fmt.Scanf("%s", &regionName)
			if err != nil {
				log.Fatalf("Error reading AWS region name: %s", err.Error())
			}
		}

		// Prompt user for AWS cluster name
		log.Info("Please enter the AWS cluster name:")
		_, err = fmt.Scanf("%s", &clusterName)
		if err != nil {
			log.Warnf("Error reading AWS cluster name: %s", err.Error())
			log.Info("Let's try again. Please enter the AWS cluster name:")
			_, err = fmt.Scanf("%s", &clusterName)
			if err != nil {
				log.Fatalf("Error reading AWS cluster name: %s", err.Error())
			}
		}

		// Build the aws CLI syntax to fetch cluster config in kubeconfig.yaml file
		eksCmd := exec.Command("aws", "eks", "--region", regionName, "update-kubeconfig", "--name", clusterName, "--kubeconfig", utils.ConfigPath)
		eksCmd.Stdout = os.Stdout
		eksCmd.Stderr = os.Stderr
		// Write EKS compatible config to the filesystem
		err = eksCmd.Run()
		if err != nil {
			log.Fatalf("Error generating kubeconfig: %s", err.Error())
			return err
		}
		log.Debugf("EKS configuration is written to: %s", utils.ConfigPath)

		// set the token in the chosen context
		setToken()
		return nil
	},
}

var gkeConfigCmd = &cobra.Command{
	Use:   "gke",
	Short: "Configure Meshery to use GKE cluster",
	Long:  `Configure Meshery to connect to GKE cluster`,
	Example: `
// Configure Meshery to connect to GKE cluster using auth token
mesheryctl system config gke --token auth.json

// Configure Meshery to connect to GKE cluster (if session is logged in using login subcommand)
mesheryctl system config gke
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return errors.New("more than one config name provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: move the GenerateConfigGKE logic to meshkit/client-go
		log.Info("Configuring Meshery to access GKE...")
		SAName := "sa-meshery-" + utils.StringWithCharset(8)
		if err := utils.GenerateConfigGKE(utils.ConfigPath, SAName, "default"); err != nil {
			log.Fatal("Error generating config:", err)
			return err
		}
		log.Debugf("GKE configuration is written to: %s", utils.ConfigPath)

		// set the token in the chosen context
		setToken()
		return nil
	},
}

var minikubeConfigCmd = &cobra.Command{
	Use:   "minikube",
	Short: "Configure Meshery to use minikube cluster",
	Long:  `Configure Meshery to connect to minikube cluster`,
	Example: `
// Configure Meshery to connect to minikube cluster using auth token
mesheryctl system config minikube --token auth.json

// Configure Meshery to connect to minikube cluster (if session is logged in using login subcommand)
mesheryctl system config minikube
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return errors.New("more than one config name provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Configuring Meshery to access Minikube...")
		// Get the config from the default config path
		if _, err = os.Stat(utils.KubeConfig); err != nil {
			log.Fatal("Could not find the default kube config:", err)
			return err
		}

		// Minifies and flattens kubeconfig and writes it to kubeconfig.yaml
		_, _, err := meshkitkube.ProcessConfig(utils.KubeConfig, utils.ConfigPath)
		if err != nil {
			log.Fatal("Error writing config to file:", err)
			return err
		}

		log.Debugf("Minikube configuration is written to: %s", utils.ConfigPath)

		// set the token in the chosen context
		setToken()
		return nil
	},
}

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
		if configureCluster != "" && utils.TokenFlag == "" {
			// authenticate the user if not already
			log.Printf("Please login first to configure %s", configureCluster)
			loginCmd.PersistentFlags().StringVarP(&providerFlag, "provider", "p", "", "login Meshery with specified provider")
			log.Printf("authenticated")
		}

		return nil
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		utils.CheckMesheryctlClientVersion(constants.GetMesheryctlVersion())
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

	if providerFlag != "" {
		currCtx.SetProvider(providerFlag)
	}

	// update the context to config
	err = config.UpdateContextInConfig(currCtx, mctlCfg.GetCurrentContextName())
	if err != nil {
		return err
	}

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}

	callbackURL := viper.GetString(pkgconstants.CallbackURLENV)
	providerURL := viper.GetString(pkgconstants.ProviderURLsENV)
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
		//extracting the custom user port from config.yaml
		userPort := strings.Split(currCtx.GetEndpoint(), ":")
		//extracting container port from the docker-compose
		containerPort := strings.Split(utils.Services["meshery"].Ports[0], ":")
		userPortMapping := userPort[len(userPort)-1] + ":" + containerPort[len(containerPort)-1]
		utils.Services["meshery"].Ports[0] = userPortMapping

		RequiredService := []string{"meshery", "watchtower"}

		AllowedServices := map[string]utils.Service{}
		for _, v := range currCtx.GetComponents() {
			if utils.Services[v].Image == "" {
				log.Fatalf("Invalid component specified %s", v)
			}

			temp, ok := utils.Services[v]
			if !ok {
				return errors.New(fmt.Sprintf("No Docker Compose service exists for Meshery component `%s`.", v))
			}

			spliter := strings.Split(temp.Image, ":")
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
			utils.Services[v] = temp
			AllowedServices[v] = utils.Services[v]
			utils.ViperCompose.Set(fmt.Sprintf("services.%s", v), utils.Services[v])
			err = utils.ViperCompose.WriteConfig()
			if err != nil {
				// failure while adding a service to docker compose file is not a fatal error
				// mesheryctl will continue deploying with required services (meshery, watchtower)
				log.Infof("Encountered an error while adding `%s` service to Docker Compose file. Verify permission to write to `.meshery/meshery.yaml` file.", v)
			}
		}

		for _, v := range RequiredService {
			if v == "watchtower" {
				AllowedServices[v] = utils.Services[v]
				continue
			}

			temp, ok := utils.Services[v]
			if !ok {
				return errors.New("unable to extract meshery version")
			}

			spliter := strings.Split(temp.Image, ":")
			temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), "latest")
			if v == "meshery" {
				callbackEnvVaridx, ok := utils.FindInSlice(pkgconstants.CallbackURLENV, temp.Environment)
				if !ok {
					temp.Environment = append(temp.Environment, fmt.Sprintf("%s=%s", pkgconstants.CallbackURLENV, callbackURL))
				} else if callbackURL != "" {
					if ok {
						temp.Environment[callbackEnvVaridx] = fmt.Sprintf("%s=%s", pkgconstants.CallbackURLENV, callbackURL)
					}
				}

				providerEnvVar := currCtx.GetProvider()
				// If user has specified provider using --provider flag use that.
				if providerFlag != "" {
					providerEnvVar = providerFlag
				}
				proivderEnvVaridx, ok := utils.FindInSlice(pkgconstants.ProviderENV, temp.Environment)

				if !ok {
					temp.Environment = append(temp.Environment, fmt.Sprintf("%s=%s", pkgconstants.ProviderENV, providerEnvVar))
				} else if providerEnvVar != "" {
					temp.Environment[proivderEnvVaridx] = fmt.Sprintf("%s=%s", pkgconstants.ProviderENV, providerEnvVar)
				}

				temp.Image = fmt.Sprintf("%s:%s-%s", spliter[0], currCtx.GetChannel(), mesheryImageVersion)
			}
			utils.Services[v] = temp
			AllowedServices[v] = utils.Services[v]
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

		// group, err := user.LookupGroup("docker")
		// if err != nil {
		// 	return errors.Wrap(err, utils.SystemError("unable to get GID of docker group"))
		// }

		// // Create the group_add option and add GID of docker group to meshery container
		// groupAdd := viper.GetStringSlice("services.meshery.group_add")
		// groupAdd = append(groupAdd, group.Gid)
		// utils.ViperCompose.Set("services.meshery.group_add", groupAdd)

		// // Write the modified configuration back to the Docker Compose file
		// if err := utils.ViperCompose.WriteConfig(); err != nil {
		// 	return errors.Wrap(err, utils.SystemError("unable to add group_add option. Meshery Server cannot perform this privileged action"))
		// }

		log.Info("Starting Meshery...")
		start := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "up", "-d")
		start.Stdout = os.Stdout
		start.Stderr = os.Stderr

		if err := start.Run(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to run Meshery Server"))
		}

		checkFlag := 0 //flag to check

		// Get the Docker configuration
		dockerCfg, err := cliconfig.Load(dockerconfig.Dir())
		if err != nil {
			return ErrCreatingDockerClient(err)
		}

		//connection to docker-client
		cli, err := dockerCmd.NewAPIClientFromFlags(cliflags.NewClientOptions(), dockerCfg)
		if err != nil {
			return ErrCreatingDockerClient(err)
		}
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{
			Filters: filters.NewArgs(),
		})
		//fetch the list of containers
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

		// log.Info("Starting Meshery...")

		spinner := utils.CreateDefaultSpinner("Deploying Meshery on Kubernetes", "\nMeshery deployed on Kubernetes.")
		spinner.Start()

		if err := utils.CreateManifestsFolder(); err != nil {
			utils.Log.Error(ErrCreateManifestsFolder(err))
			return err
		}

		// Applying Meshery Helm charts for installing Meshery
		if err = applyHelmCharts(kubeClient, currCtx, mesheryImageVersion, false, meshkitkube.INSTALL, callbackURL, providerURL); err != nil {
			return err
		}

		// checking if Meshery is ready
		time.Sleep(20 * time.Second) // sleeping 10 seconds to countermeasure time to apply helm charts
		ready, err := mesheryReadinessHealthCheck()
		if err != nil {
			log.Info(err)
		}

		spinner.Stop()

		if !ready {
			log.Info("\nTimeout. Meshery pod(s) is not running, yet.\nCheck status of Meshery pod(s) by executing “mesheryctl system status`. Expose Meshery UI with `mesheryctl system dashboard` as needed.")
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
	availableSubcommands = []*cobra.Command{
		aksConfigCmd,
		eksConfigCmd,
		gkeConfigCmd,
		minikubeConfigCmd,
	}

	aksConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	eksConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	gkeConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	minikubeConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")

	startCmd.PersistentFlags().StringVarP(&utils.PlatformFlag, "platform", "p", "", "platform to deploy Meshery to.")
	startCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
	startCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
	startCmd.Flags().BoolVarP(&skipBrowserFlag, "skip-browser", "", false, "(optional) skip opening of MesheryUI in browser.")
	startCmd.Flags().StringVarP(&configureCluster, "config-cluster", "", "", "(optional) configure the Kubernetes cluster used by Meshery.")
	startCmd.PersistentFlags().StringVar(&providerFlag, "provider", "", "(optional) Defaults to the provider specified in the current context")

	startCmd.AddCommand(availableSubcommands...)
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

// Given the token path, get the context and set the token in the chosen context
func setToken() {
	log.Debugf("Token path: %s", utils.TokenFlag)
	contexts, err := getContexts(utils.ConfigPath)
	if err != nil {
		utils.Log.Error(err)
	}
	if contexts == nil || len(contexts) < 1 {
		log.Fatalf("Error getting context: %s", fmt.Errorf("no contexts found"))
	}
	choosenCtx := contexts[0]
	if len(contexts) > 1 {
		fmt.Println("List of available contexts: ")
		for i, ctx := range contexts {
			fmt.Printf("(%d) %s \n", i+1, ctx)
		}
		var choice int
		fmt.Print("Enter choice (number): ")
		_, err = fmt.Scanf("%d", &choice)
		if err != nil {
			log.Fatalf("Error reading input:  %s", err.Error())
		}
		choosenCtx = contexts[choice-1]
	}

	log.Debugf("Chosen context : %s out of the %d available contexts", choosenCtx, len(contexts))
	err = setContext(utils.ConfigPath, choosenCtx)
	if err != nil {
		log.Fatalf("Error setting context: %s", err.Error())
	}
}
