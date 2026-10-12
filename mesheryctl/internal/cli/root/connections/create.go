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
	"path/filepath"
	"slices"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	supportedConnectionTypes = []string{"aks", "eks", "gke", "kubernetes", "minikube"}
	connectionType           string
	kubeconfigPath           string
	contextFlag              string
)

type userPrompt struct {
	request                 string
	errorReadingResourceMsg string
}

var createConnectionCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new connection",
	Long: `Create a new connection to a Kubernetes cluster or other supported platform.
	Find more information at: https://docs.meshery.io/reference/references/mesheryctl/connection/create`,
	Example: `
// Create a new Kubernetes connection from a kubeconfig file
mesheryctl connection create --file ~/.kube/config
mesheryctl connection create --file ~/.kube/config --context my-cluster

// Create a new Kubernetes connection using a specific type
mesheryctl connection create --type aks
mesheryctl connection create --type eks
mesheryctl connection create --type gke
mesheryctl connection create --type minikube

// Create a connection with a token
mesheryctl connection create --type gke --token auth.json
	`,
	Args: func(_ *cobra.Command, args []string) error {
		if connectionType == "" && kubeconfigPath == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("either --file or --type is required. Specify a kubeconfig with --file or use --type (%s)", strings.Join(supportedConnectionTypes, "|")))
		}
		if connectionType != "" && !slices.Contains(supportedConnectionTypes, connectionType) {
			return errInvalidConnectionType(connectionType)
		}
		if kubeconfigPath != "" && connectionType != "" && connectionType != "kubernetes" {
			return utils.ErrInvalidArgument(fmt.Errorf("--file cannot be used with --type %s", connectionType))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if kubeconfigPath != "" {
			return createKubeconfigConnection(kubeconfigPath)
		}
		switch connectionType {
		case "kubernetes":
			return createKubeconfigConnection(utils.KubeConfig)
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
	err = setToken(contextFlag)
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
	err = setToken(contextFlag)
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
	err := setToken(contextFlag)
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
	if err := writeKubeconfigSafely(kubeConfig, utils.ConfigPath); err != nil {
		return err
	}
	utils.Log.Debugf("Minikube configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err = setToken(contextFlag)
	if err != nil {
		return err
	}

	utils.Log.Info("Minikube connection created.")
	return nil
}

// writeKubeconfigSafely writes the flattened kubeconfig to a temporary file with
// restricted permissions (0600) in the destination directory and atomically replaces
// destPath to avoid leaving partially written or corrupted configuration on failure.
func writeKubeconfigSafely(kubeConfig *clientcmdapi.Config, destPath string) error {
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return errWriteKubeConfig(err)
	}

	tmpFile, err := os.CreateTemp(destDir, "kubeconfig-*.tmp")
	if err != nil {
		return errWriteKubeConfig(err)
	}
	tmpPath := tmpFile.Name()

	// Ensure the temporary file descriptor is closed before writing or renaming,
	// and clean up the temporary file if any subsequent step fails.
	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return errWriteKubeConfig(err)
	}
	_ = os.Chmod(tmpPath, 0600)

	if err := clientcmd.WriteToFile(*kubeConfig, tmpPath); err != nil {
		_ = os.Remove(tmpPath)
		return errWriteKubeConfig(err)
	}

	if err := os.Rename(tmpPath, destPath); err != nil {
		_ = os.Remove(tmpPath)
		return errWriteKubeConfig(err)
	}

	return nil
}

// createKubeconfigConnection loads a generic kubeconfig file, validates that it
// contains at least one context and that any explicitly requested --context exists,
// securely writes the flattened configuration to utils.ConfigPath, and registers
// the connection with Meshery.
func createKubeconfigConnection(filePath string) error {
	utils.Log.Info("Configuring Meshery to access Kubernetes cluster...")

	if strings.HasPrefix(filePath, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			filePath = filepath.Join(home, filePath[2:])
		}
	} else if filePath == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			filePath = home
		}
	}

	if _, err := os.Stat(filePath); err != nil {
		return errReadKubeConfig(err)
	}
	kubeConfig, err := clientcmd.LoadFromFile(filePath)
	if kubeConfig == nil || err != nil {
		return errReadKubeConfig(err)
	}
	if len(kubeConfig.Contexts) == 0 {
		return utils.ErrGetKubernetesContexts(fmt.Errorf("no contexts found in %s", filePath))
	}

	// Validate explicitly requested context exists in the kubeconfig before writing
	if contextFlag != "" {
		if _, exists := kubeConfig.Contexts[contextFlag]; !exists {
			var availableContexts []string
			for name := range kubeConfig.Contexts {
				availableContexts = append(availableContexts, name)
			}
			slices.Sort(availableContexts)
			return utils.ErrInvalidArgument(fmt.Errorf("context %q not found in kubeconfig (available contexts: %s)", contextFlag, strings.Join(availableContexts, ", ")))
		}
	}

	// Flatten the config file
	err = clientcmdapi.FlattenConfig(kubeConfig)
	if err != nil {
		return errReadKubeConfig(err)
	}

	// Securely write flattened config to destination via private temp file
	if err := writeKubeconfigSafely(kubeConfig, utils.ConfigPath); err != nil {
		return err
	}
	utils.Log.Debugf("Kubernetes configuration is written to: %s", utils.ConfigPath)

	// set the token in the chosen context
	err = setToken(contextFlag)
	if err != nil {
		return err
	}

	utils.Log.Info("Kubernetes connection created.")
	return nil
}

// getContexts uploads the specified kubeconfig file to the Meshery contexts API
// and returns the names of all discovered contexts.
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

// registeredK8sContext is the subset of the server's per-context registration
// response that the CLI needs. The `/api/system/kubernetes` response
// (server SaveK8sContextResponse) has no schemas-generated equivalent, so this
// is a parse-only DTO — not a duplicate of a schemas type.
type registeredK8sContext struct {
	Name         string `json:"name"`
	ConnectionID string `json:"connectionId"`
}

type saveK8sContextResponse struct {
	RegisteredContexts []registeredK8sContext `json:"registeredContexts"`
	ConnectedContexts  []registeredK8sContext `json:"connectedContexts"`
}

// setContext registers the chosen kubeconfig context with Meshery and returns
// the id of the connection created (or reused) for that context.
func setContext(configFile, cname string) (string, error) {
	selectedJSON, err := json.Marshal([]string{cname})
	if err != nil {
		return "", err
	}
	contextParams := map[string]string{
		"selectedContexts": string(selectedJSON),
		"contextName":      cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return "", err
	}

	// setContextURL endpoint points to set context
	setContextURL := mctlCfg.GetBaseMesheryURL() + "/api/system/kubernetes"
	req, err := utils.UploadFileWithParams(setContextURL, contextParams, utils.ParamName, configFile)

	if err != nil {
		return "", utils.ErrUploadFileWithParams(err, configFile)
	}
	res, err := utils.MakeRequest(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", utils.ErrReadResponseBody(err)
	}
	utils.Log.Debugf("Set context API response: %s", string(body))

	var response saveK8sContextResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// The context was still registered; only id resolution is best-effort.
		utils.Log.Debugf("Unable to parse set context response for connection id: %s", err.Error())
		return "", nil
	}

	return connectionIDForContext(response, cname), nil
}

// connectionIDForContext resolves the connection id for the given context name.
// It returns the id only on an exact name match; when the server does not echo
// the requested context name it returns "" rather than guessing. A
// multi-context kubeconfig registers every context, so an arbitrary fallback
// could surface a different context's id than the one the user selected and
// mis-target `connection view`/`connection delete`.
func connectionIDForContext(response saveK8sContextResponse, cname string) string {
	all := append(append([]registeredK8sContext{}, response.ConnectedContexts...), response.RegisteredContexts...)
	for _, ctx := range all {
		if ctx.Name == cname && ctx.ConnectionID != "" {
			return ctx.ConnectionID
		}
	}
	return ""
}

// Given the token path, get the context and set the token in the chosen context
func setToken(specifiedContext string) error {
	utils.Log.Debugf("Token path: %s", utils.TokenFlag)
	contexts, err := getContexts(utils.ConfigPath)
	if err != nil {
		return utils.ErrGetKubernetesContexts(err)
	}

	utils.Log.Debugf("Available contexts: %s", contexts)
	if len(contexts) < 1 {
		return utils.ErrGetKubernetesContexts(fmt.Errorf("no contexts found"))
	}

	var chosenCtx string
	if specifiedContext != "" {
		if !slices.Contains(contexts, specifiedContext) {
			return utils.ErrInvalidArgument(fmt.Errorf("context %q not found in kubeconfig (available contexts: %s)", specifiedContext, strings.Join(contexts, ", ")))
		}
		chosenCtx = specifiedContext
	} else if len(contexts) == 1 {
		chosenCtx = contexts[0]
	} else {
		i, err := utils.RunSelectPrompt("Select context for the connection", contexts)
		if err != nil {
			return err
		}
		chosenCtx = contexts[i]
	}
	utils.Log.Debugf("Chosen context : %s out of the %d available contexts", chosenCtx, len(contexts))

	connectionID, err := setContext(utils.ConfigPath, chosenCtx)
	if err != nil {
		return utils.ErrSetKubernetesContext(err)
	}

	utils.Log.Infof("Token set in context %s", chosenCtx)
	// Emit the created/reused connection id in a stable, parseable form so it
	// can be fed to `connection view`/`connection delete` (and captured by e2e
	// tests). The prefix is a documented contract — do not reformat casually.
	if connectionID != "" {
		utils.Log.Infof("connection_id: %s", connectionID)
	}
	return nil
}

func init() {
	createConnectionCmd.Flags().StringVarP(&connectionType, "type", "t", "", "Type of connection to create (aks|eks|gke|kubernetes|minikube)")
	createConnectionCmd.Flags().StringVarP(&kubeconfigPath, "file", "f", "", "Path to kubeconfig file")
	createConnectionCmd.Flags().StringVarP(&contextFlag, "context", "c", "", "Context name to select from the kubeconfig (optional)")
	createConnectionCmd.Flags().StringVar(&utils.TokenFlag, "token", "", "Path to token for authenticating to Meshery API")
}
