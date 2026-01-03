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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
)

func getContexts(configFile string) ([]string, error) {

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.LogError.Error(err)
		return nil, nil
	}

	// GETCONTEXTS endpoint points to the URL return the contexts available
	GETCONTEXTS := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes/contexts"

	req, err := utils.UploadFileWithParams(GETCONTEXTS, nil, utils.ParamName, configFile)
	if err != nil {
		return nil, ErrUploadFileParams(err)
	}

	// Preserve the multipart Content-Type header when sending via the api client.
	headers := map[string]string{
		"Content-Type": req.Header.Get("Content-Type"),
	}

	res, err := api.Add("api/system/kubernetes/contexts", req.Body, headers)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, utils.ErrReadResponseBody(err)
	}

	utils.Log.Debugf("Get context API response: %s", string(body))
	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, utils.ErrUnmarshal(err)
	}

	if results == nil {
		errstr := "Error unmarshalling the context info, check " + configFile + " file"
		return nil, errors.New(errstr)
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
	utils.Log.Debugf("Available contexts: %s", contextNames)
	return contextNames, nil
}

func setContext(configFile, cname string) error {
	extraParams1 := map[string]string{
		"contextName": cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.LogError.Error(err)
		return nil
	}

	// SETCONTEXT endpoint points to set context
	SETCONTEXT := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes"
	req, err := utils.UploadFileWithParams(SETCONTEXT, extraParams1, utils.ParamName, configFile)
	if err != nil {
		return ErrUploadFileParams(err)
	}
	headers := map[string]string{
		"Content-Type": req.Header.Get("Content-Type"),
	}
	res, err := api.Add("api/system/kubernetes", req.Body, headers)
	if err != nil {
		return utils.ErrRequestResponse(err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrReadResponseBody(err)
	}
	// TODO: Pretty print the output
	utils.Log.Debugf("Set context API response: %s", string(body))
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
	Args: func(cmd *cobra.Command, args []string) error {
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
			utils.LogError.Error(errors.New("Azure CLI not found. Please install Azure CLI and try again. \nSee https://docs.microsoft.com/en-us/cli/azure/install-azure-cli "))
			os.Exit(1)
		}
		utils.Log.Info("Configuring Meshery to access AKS...")
		var resourceGroup, aksName string

		// Prompt user for Azure resource name
		utils.Log.Info("Please enter the Azure resource group name:")
		_, err = fmt.Scanf("%s", &resourceGroup)
		if err != nil {
			utils.Log.Warnf("Error reading Azure resource group name: %s", err.Error())
			utils.Log.Info("Let's try again. Please enter the Azure resource group name:")
			_, err = fmt.Scanf("%s", &resourceGroup)
			if err != nil {
				utils.LogError.Error(errors.Wrapf(err, "Error reading Azure resource group name: %s", err.Error()))
				os.Exit(1)
			}
		}

		// Prompt user for AKS cluster name
		utils.Log.Info("Please enter the AKS cluster name:")
		_, err = fmt.Scanf("%s", &aksName)
		if err != nil {
			utils.Log.Warnf("Error reading AKS cluster name: %s", err.Error())
			utils.Log.Info("Let's try again. Please enter the AKS cluster name:")
			_, err = fmt.Scanf("%s", &aksName)
			if err != nil {
				utils.LogError.Error(errors.Wrapf(err, "Error reading AKS cluster name: %s", err.Error()))
				os.Exit(1)
			}
		}

		// Build the Azure CLI syntax to fetch cluster config in kubeconfig.yaml file
		aksCmd := exec.Command("az", "aks", "get-credentials", "--resource-group", resourceGroup, "--name", aksName, "--file", utils.ConfigPath)
		aksCmd.Stdout = os.Stdout
		aksCmd.Stderr = os.Stderr
		// Write AKS compatible config to the filesystem
		err = aksCmd.Run()
		if err != nil {
			utils.LogError.Error(errors.Wrapf(err, "Error generating kubeconfig: %s", err.Error()))
			return err
		}
		utils.Log.Debugf("AKS configuration is written to: %s", utils.ConfigPath)

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
	Args: func(cmd *cobra.Command, args []string) error {
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
			utils.LogError.Error(errors.Wrap(err, "AWS CLI not found. Please install AWS CLI and try again. \nSee https://docs.aws.amazon.com/cli/latest/reference/ "))
			os.Exit(1) }
		utils.Log.Info("Configuring Meshery to access EKS...")
		var regionName, clusterName string

		// Prompt user for AWS region name
		utils.Log.Info("Please enter the AWS region name:")
		_, err = fmt.Scanf("%s", &regionName)
		if err != nil {
			utils.Log.Warnf("Error reading AWS region name: %s", err.Error())
			utils.Log.Info("Let's try again. Please enter the AWS region name:")
			_, err = fmt.Scanf("%s", &regionName)
			if err != nil {
				utils.LogError.Error(errors.Wrapf(err, "Error reading AWS region name: %s", err.Error()))
				os.Exit(1)
			}
		}

		// Prompt user for AWS cluster name
		utils.Log.Info("Please enter the AWS cluster name:")
		_, err = fmt.Scanf("%s", &clusterName)
		if err != nil {
			utils.Log.Warnf("Error reading AWS cluster name: %s", err.Error())
			utils.Log.Info("Let's try again. Please enter the AWS cluster name:")
			_, err = fmt.Scanf("%s", &clusterName)
			if err != nil {
				utils.LogError.Error(errors.Wrapf(err, "Error reading AWS cluster name: %s", err.Error()))
				os.Exit(1)
			}
		}

		// Build the aws CLI syntax to fetch cluster config in kubeconfig.yaml file
		eksCmd := exec.Command("aws", "eks", "--region", regionName, "update-kubeconfig", "--name", clusterName, "--kubeconfig", utils.ConfigPath)
		eksCmd.Stdout = os.Stdout
		eksCmd.Stderr = os.Stderr
		// Write EKS compatible config to the filesystem
		err = eksCmd.Run()
		if err != nil {
			utils.LogError.Error(errors.Wrapf(err, "Error generating kubeconfig: %s", err.Error()))
			return err
		}
		utils.Log.Debugf("EKS configuration is written to: %s", utils.ConfigPath)

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
		utils.Log.Info("Configuring Meshery to access GKE...")
		SAName := "sa-meshery-" + utils.StringWithCharset(8)
		if err := utils.GenerateConfigGKE(utils.ConfigPath, SAName, "default"); err != nil {
			utils.LogError.Error(errors.Wrap(err, "Error generating config: "))
			return err
		}
		utils.Log.Debugf("GKE configuration is written to: %s", utils.ConfigPath)

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
		utils.Log.Info("Configuring Meshery to access Minikube...")
		// Get the config from the default config path
		if _, err = os.Stat(utils.KubeConfig); err != nil {
			utils.LogError.Error(errors.Wrap(err, "Could not find the default kube config:"))
			return err
		}

		// Minifies and flattens kubeconfig and writes it to kubeconfig.yaml
		_, _, err := meshkitkube.ProcessConfig(utils.KubeConfig, utils.ConfigPath)
		if err != nil {
			utils.LogError.Error(errors.Wrap(err, "Error writing config to file: "))
			return err
		}

		utils.Log.Infof("A flattened Minikube kubeconfig file available at: %s", utils.ConfigPath)
		utils.Log.Info("A new Meshery connection has been created. Run `mesheryctl connnection list` for details.")

		// set the token in the chosen context
		setToken()
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:        "config",
	Short:      "Configure Meshery",
	Long:       `Configure the Kubernetes cluster used by Meshery.`,
	Deprecated: "Please use 'mesheryctl connection create --type <k8s-type>' instead.",
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Usage: mesheryctl system config [aks|eks|gke|minikube]
Example: mesheryctl system config eks
Description: Configure the Kubernetes cluster used by Meshery.
NOTE: This command is deprecated. Please use 'mesheryctl connection create --type <k8s-type>' instead.`
		if len(args) == 0 {
			return fmt.Errorf("name of kubernetes cluster to configure Meshery not provided\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("expected one argument received multiple arguments")
		}
		return nil
	},
	Example: `
// Set configuration according to k8s cluster (DEPRECATED: use 'mesheryctl connection create' instead)
mesheryctl system config [aks|eks|gke|minikube]

// Path to token for authenticating to Meshery API (optional, can be done alternatively using "login")
mesheryctl system config --token "~/Downloads/auth.json"

// NEW: Use the connection create command instead
mesheryctl connection create --type [aks|eks|gke|minikube]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemError(fmt.Sprintf("invalid command: \"%s\".", args[0])))
		}
		return nil
	},
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

	configCmd.AddCommand(availableSubcommands...)
}

// Given the token path, get the context and set the token in the chosen context
func setToken() {
	utils.Log.Debugf("Token path: %s", utils.TokenFlag)
	contexts, err := getContexts(utils.ConfigPath)
	if err != nil {
		utils.LogError.Error(err)
		os.Exit(1)
	}
	if len(contexts) < 1 {
		utils.LogError.Error(errors.Wrap(fmt.Errorf("no contexts found"), "Error getting context: "))
		os.Exit(1)
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
			utils.LogError.Error(errors.Wrapf(err, "Error reading input: %s", err.Error()))
			os.Exit(1)
		}
		choosenCtx = contexts[choice-1]
	}

	utils.Log.Debugf("Chosen context : %s out of the %d available contexts", choosenCtx, len(contexts))
	err = setContext(utils.ConfigPath, choosenCtx)
	if err != nil {
		utils.LogError.Error(errors.Wrapf(err, "Error setting context: %s", err.Error()))
		os.Exit(1)
	}
}
