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

package connections

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	supportedConnectionTypes = []string{"aks", "eks", "gke", "minikube"}
	connectionType           string
)

type userPrompt struct {
	request                 string
	errorReadingResourceMsg string
}

var createConnectionCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new connection",
	Long:  `Create a new connection to a Kubernetes cluster or other supported platform`,
	Example: `
// Create a new Kubernetes connection using a specific type
mesheryctl connection create --type aks
mesheryctl connection create --type eks
mesheryctl connection create --type gke
mesheryctl connection create --type minikube

// Create a connection with a token
mesheryctl connection create --type gke --token auth.json
	`,
	Args: func(_ *cobra.Command, args []string) error {
		if connectionType == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("connection type is required. Use --type flag to specify the type (%s)", strings.Join(supportedConnectionTypes, "|")))
		}
		if !slices.Contains(supportedConnectionTypes, connectionType) {
			return errInvalidConnectionType(connectionType)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch connectionType {
		case "aks":
			return createAKSConnection()
		case "eks":
			return createEKSConnection()
		case "gke":
			return createGKEConnection()
		case "minikube":
			return createMinikubeConnection()
		default:
			return fmt.Errorf("unsupported connection type: %s", connectionType)
		}
	},
}

func getUserPrompt(userPrompt userPrompt) (string, error) {
	var prompt string
	utils.Log.Info(userPrompt.request)
	_, err := fmt.Scanf("%s", &prompt)
	if err != nil {
		utils.Log.Warnf("Error reading %s: %s", userPrompt.errorReadingResourceMsg, err.Error())
		utils.Log.Info(fmt.Sprintf("Let's try again. %s", userPrompt.request))
		_, err = fmt.Scanf("%s", &prompt)
		if err != nil {
			return "", utils.ErrReadInput(err)
		}
	}
	return prompt, nil
}

func createAKSConnection() error {
	aksCheck := exec.Command("az", "version")
	aksCheck.Stdout = os.Stdout
	aksCheck.Stderr = os.Stderr
	err := aksCheck.Run()
	if err != nil {
		return errAzureCliNotFound(err)
	}

	utils.Log.Info("Configuring Meshery to access AKS...")
	var resourceGroup, aksName string

	resourceGroup, err = getUserPrompt(userPrompt{request: "Please enter the Azure resource group name:", errorReadingResourceMsg: "Azure resource group name"})
	if err != nil {
		return err
	}

	aksName, err = getUserPrompt(userPrompt{request: "Please enter the AKS cluster name:", errorReadingResourceMsg: "AKS cluster name"})
	if err != nil {
		return err
	}

	// Build the Azure CLI syntax to fetch cluster config in kubeconfig.yaml file
	aksCmd := exec.Command("az", "aks", "get-credentials", "--resource-group", resourceGroup, "--name", aksName, "--file", utils.ConfigPath)
	aksCmd.Stdout = os.Stdout
	aksCmd.Stderr = os.Stderr
	// Write AKS compatible config to the filesystem
	err = aksCmd.Run()
	if err != nil {
		return errAzureAksGetCredentials(err)
	}
	utils.Log.Debugf("AKS configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err = setToken()
	if err != nil {
		return err
	}

	utils.Log.Infof("AKS connection on cluster %s created.", aksName)
	return nil
}

func createEKSConnection() error {
	eksCheck := exec.Command("aws", "--version")
	eksCheck.Stdout = os.Stdout
	eksCheck.Stderr = os.Stderr
	err := eksCheck.Run()
	if err != nil {
		return errAwsCliNotFound(err)

	}

	utils.Log.Info("Configuring Meshery to access EKS...")
	var regionName, clusterName string

	regionName, err = getUserPrompt(userPrompt{request: "Please enter the AWS region name:", errorReadingResourceMsg: "AWS region name"})
	if err != nil {
		return err
	}

	clusterName, err = getUserPrompt(userPrompt{request: "Please enter the EKS cluster name:", errorReadingResourceMsg: "EKS cluster name"})
	if err != nil {
		return err
	}

	// Build the aws CLI syntax to fetch cluster config in kubeconfig.yaml file
	eksCmd := exec.Command("aws", "eks", "--region", regionName, "update-kubeconfig", "--name", clusterName, "--kubeconfig", utils.ConfigPath)
	eksCmd.Stdout = os.Stdout
	eksCmd.Stderr = os.Stderr
	// Write EKS compatible config to the filesystem
	err = eksCmd.Run()
	if err != nil {
		return errAwsEksGetCredentials(err)
	}
	utils.Log.Debugf("EKS configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err = setToken()
	if err != nil {
		return err
	}

	utils.Log.Infof("EKS connection on cluster %s created.", clusterName)
	return nil
}

func createGKEConnection() error {
	// TODO: move the GenerateConfigGKE logic to meshkit/client-go
	utils.Log.Info("Configuring Meshery to access GKE...")
	SAName := "sa-meshery-" + utils.StringWithCharset(8)
	if err := utils.GenerateConfigGKE(utils.ConfigPath, SAName, "default"); err != nil {
		return errGcpGKEGetCredentials(err)
	}
	utils.Log.Debugf("GKE configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err := setToken()
	if err != nil {
		return err
	}

	utils.Log.Info("GKE connection created.")
	return nil
}

func createMinikubeConnection() error {
	utils.Log.Info("Configuring Meshery to access Minikube...")
	// Get the config from the default config path
	if _, err := os.Stat(utils.KubeConfig); err != nil {
		return errReadKubeConfig(err)
	}
	kubeConfig, err := clientcmd.LoadFromFile(utils.KubeConfig)
	if kubeConfig == nil || err != nil {
		return errReadKubeConfig(err)
	}
	// Flatten the config file
	err = clientcmdapi.FlattenConfig(kubeConfig)
	if err != nil {
		return errReadKubeConfig(err)
	}
	// write the flattened config to kubeconfig.yaml file
	err = clientcmd.WriteToFile(*kubeConfig, utils.ConfigPath)
	if err != nil {
		return errWriteKubeConfig(err)
	}
	utils.Log.Debugf("Minikube configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err = setToken()
	if err != nil {
		return err
	}

	utils.Log.Info("Minikube connection created.")
	return nil
}

func getContexts(configFile string) ([]string, error) {
	client := &http.Client{}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, err
	}

	// getContextsURL endpoint points to the URL returning the available contexts
	getContextsURL := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes/contexts"

	req, err := utils.UploadFileWithParams(getContextsURL, nil, utils.ParamName, configFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to upload file with parameters")
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, utils.ErrRequestResponse(err)
	}
	if res.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(res.Body)
		_ = res.Body.Close()
		return nil, fmt.Errorf("failed to get contexts: received status code %d with body %s", res.StatusCode, string(body))
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
	client := &http.Client{}
	contextParams := map[string]string{
		"contextName": cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}

	// setContextURL endpoint points to set context
	setContextURL := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes"
	req, err := utils.UploadFileWithParams(setContextURL, contextParams, utils.ParamName, configFile)

	if err != nil {
		return utils.ErrUploadFileWithParams(err, configFile)
	}
	res, err := client.Do(req)
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

// Given the token path, get the context and set the token in the chosen context
func setToken() error {
	utils.Log.Debugf("Token path: %s", utils.TokenFlag)
	contexts, err := getContexts(utils.ConfigPath)
	if err != nil {
		return utils.ErrGetKubernetesContexts(err)
	}

	utils.Log.Debugf("Available contexts: %s", contexts)
	if len(contexts) < 1 {
		return utils.ErrGetKubernetesContexts(fmt.Errorf("no contexts found"))
	}

	chosenCtx := contexts[0]
	if len(contexts) > 1 {
		fmt.Println("List of available contexts: ")

		prompt := promptui.Select{
			Label: "Select context for the connection",
			Items: contexts,
		}

		for {
			i, _, err := prompt.Run()
			if err != nil {
				continue
			}

			chosenCtx = contexts[i]
			break
		}

	}
	utils.Log.Debugf("Chosen context : %s out of the %d available contexts", chosenCtx, len(contexts))

	err = setContext(utils.ConfigPath, chosenCtx)
	if err != nil {
		return utils.ErrSetKubernetesContext(err)
	}

	utils.Log.Infof("Token set in context %s", chosenCtx)
	return nil
}

func init() {
	createConnectionCmd.Flags().StringVarP(&connectionType, "type", "t", "", "Type of connection to create (aks|eks|gke|minikube)")
	createConnectionCmd.Flags().StringVar(&utils.TokenFlag, "token", "", "Path to token for authenticating to Meshery API")
}
