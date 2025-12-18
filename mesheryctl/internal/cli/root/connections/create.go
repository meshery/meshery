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

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	connectionType string
)

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
			return errors.New("connection type is required. Use --type flag to specify the type (aks|eks|gke|minikube)")
		}
		validTypes := []string{"aks", "eks", "gke", "minikube"}
		for _, t := range validTypes {
			if connectionType == t {
				return nil
			}
		}
		return fmt.Errorf("invalid connection type '%s'. Valid types are: aks, eks, gke, minikube", connectionType)
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

func createAKSConnection() error {
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
}

func createEKSConnection() error {
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
}

func createGKEConnection() error {
	// TODO: move the GenerateConfigGKE logic to meshkit/client-go
	log.Info("Configuring Meshery to access GKE...")
	SAName := "sa-meshery-" + utils.StringWithCharset(8)
	if err := utils.GenerateConfigGKE(utils.ConfigPath, SAName, "default"); err != nil {
		log.Errorf("Error generating config: %v", err)
		return err
	}
	log.Debugf("GKE configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	setToken()
	return nil
}

func createMinikubeConnection() error {
	log.Info("Configuring Meshery to access Minikube...")
	// Get the config from the default config path
	if _, err := os.Stat(utils.KubeConfig); err != nil {
		log.Fatal("Could not find the default kube config:", err)
		return err
	}
	kubeConfig, err := clientcmd.LoadFromFile(utils.KubeConfig)
	if kubeConfig == nil || err != nil {
		log.Fatal("Error reading the default kube config:", err)
		return err
	}
	// Flatten the config file
	err = clientcmdapi.FlattenConfig(kubeConfig)
	if err != nil {
		log.Fatal("Error flattening config:", err)
		return err
	}
	// write the flattened config to kubeconfig.yaml file
	err = clientcmd.WriteToFile(*kubeConfig, utils.ConfigPath)
	if err != nil {
		log.Fatal("Error writing config to file:", err)
		return err
	}
	log.Debugf("Minikube configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	setToken()
	return nil
}

func getContexts(configFile string) ([]string, error) {
	client := &http.Client{}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
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
		defer res.Body.Close()
		return nil, fmt.Errorf("failed to get contexts: received status code %d with body %s", res.StatusCode, string(body))
	}
	defer res.Body.Close()

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
	log.Debugf("Available contexts: %s", contextNames)
	return contextNames, nil
}

func setContext(configFile, cname string) error {
	client := &http.Client{}
	contextParams := map[string]string{
		"contextName": cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		utils.Log.Error(err)
		return err
	}

	// setContextURL endpoint points to set context
	setContextURL := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes"
	req, err := utils.UploadFileWithParams(setContextURL, contextParams, utils.ParamName, configFile)
	if err != nil {
		return errors.Wrap(err, "failed to upload file with parameters")
	}
	res, err := client.Do(req)
	if err != nil {
		return utils.ErrRequestResponse(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return utils.ErrReadResponseBody(err)
	}
	// TODO: Pretty print the output
	log.Debugf("Set context API response: %s", string(body))
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
		log.Errorf("Error getting context: %s", fmt.Errorf("no contexts found"))
	}
	chosenCtx := contexts[0]
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
		if choice < 1 || choice > len(contexts) {
			log.Fatalf("Invalid choice: %d. Please select a number between 1 and %d.", choice, len(contexts))
		}
		chosenCtx = contexts[choice-1]
	}

	log.Debugf("Chosen context : %s out of the %d available contexts", chosenCtx, len(contexts))
	err = setContext(utils.ConfigPath, chosenCtx)
	if err != nil {
		log.Fatalf("Error setting context: %s", err.Error())
	}
}

func init() {
	createConnectionCmd.Flags().StringVarP(&connectionType, "type", "t", "", "Type of connection to create (aks|eks|gke|minikube)")
	createConnectionCmd.Flags().StringVar(&utils.TokenFlag, "token", "", "Path to token for authenticating to Meshery API")
}
