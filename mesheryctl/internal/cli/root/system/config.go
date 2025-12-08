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
	"net/http"
	"os"
	"os/exec"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
)

	var (
		okeRegionFlag string
		okeClusterIDFlag string
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

	res, err := api.Add("api/system/kubernetes/contexts", req.Body)
	if err != nil {
		return nil, err
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

var okeConfigCmd = &cobra.Command{
	Use:   "oke",
	Short: "Configure Meshery to use OKE cluster",
	Long:  `Configure Meshery to connect to OKE cluster`,
	Example: `
// Interactive mode - prompts for region and cluster ID
mesheryctl system config oke

// Specify cluster ID and region (non-interactive)
mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx --region us-phoenix-1

// Using short flag for region
mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx -r us-phoenix-1

// Configure with authentication token (for CI/CD or automation)
mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx -r us-phoenix-1 --token auth.json
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) >= 1 {
			return errors.New("more than one config name provided")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		okeCheck := exec.Command("oci", "--version")
		okeCheck.Stdout = os.Stdout
		okeCheck.Stderr = os.Stderr
		err := okeCheck.Run()
		if err != nil {
			log.Fatalf("OCI CLI not found. Please install OCI CLI and try again. \nSee https://docs.oracle.com/en-us/iaas/Content/API/SDKDocs/cliinstall.htm ")
		}
		log.Info("Configuring Meshery to access OKE...")
		var regionName, clusterID string

		// Respect flags if provided
		regionName = okeRegionFlag
		clusterID = okeClusterIDFlag

		// Prompt user for OCI region name if not provided via flag
		if regionName == "" {
			fmt.Print("Please enter the OCI region name: ")
			_, err = fmt.Scanf("%s", &regionName)
			if err != nil {
				log.Warnf("Error reading OCI region name: %s", err.Error())
				fmt.Print("Let's try again. Please enter the OCI region name: ")
				_, err = fmt.Scanf("%s", &regionName)
				if err != nil {
					log.Fatalf("Error reading OCI region name: %s", err.Error())
				}
			}
		}

		// Prompt user for OKE cluster ID if not provided via flag
		if clusterID == "" {
			fmt.Print("Please enter the OKE cluster ID: ")
			_, err = fmt.Scanf("%s", &clusterID)
			if err != nil {
				log.Warnf("Error reading OKE cluster ID: %s", err.Error())
				fmt.Print("Let's try again. Please enter the OKE cluster ID: ")
				_, err = fmt.Scanf("%s", &clusterID)
				if err != nil {
					log.Fatalf("Error reading OKE cluster ID: %s", err.Error())
				}
			}
		}

		// Build the OCI CLI syntax to fetch cluster config in kubeconfig.yaml file
		okeCmd := exec.Command("oci", "ce", "cluster", "create-kubeconfig", "--cluster-id", clusterID, "--region", regionName, "--file", utils.ConfigPath, "--overwrite")
		okeCmd.Stdout = os.Stdout
		okeCmd.Stderr = os.Stderr
		// Write OKE compatible config to the filesystem
		err = okeCmd.Run()
		if err != nil {
			log.Fatalf("Error generating kubeconfig: %s", err.Error())
			return err
		}
		log.Debugf("OKE configuration is written to: %s", utils.ConfigPath)

		// set the token in the chosen context
		setToken()
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Meshery",
	Long: `Configure the Kubernetes cluster used by Meshery.

This command provides subcommands for configuring Meshery to work with various
Kubernetes cluster types. Each subcommand automates the process of generating
kubeconfig and registering the cluster context with Meshery.

Supported Kubernetes Platforms:
  - aks:      Azure Kubernetes Service (Microsoft Azure)
  - eks:      Elastic Kubernetes Service (Amazon Web Services)  
  - gke:      Google Kubernetes Engine (Google Cloud Platform)
  - minikube: Minikube (Local Kubernetes)
  - oke:      Oracle Kubernetes Engine (Oracle Cloud Infrastructure)

Each platform has specific prerequisites and configuration requirements:

Azure Kubernetes Service (AKS):
  Prerequisites: Azure CLI installed and configured
  Authentication: Uses Azure CLI credentials
  
Amazon Elastic Kubernetes Service (EKS):
  Prerequisites: AWS CLI installed and configured
  Authentication: Uses AWS CLI credentials and IAM roles
  
Google Kubernetes Engine (GKE):
  Prerequisites: Google Cloud SDK installed and configured
  Authentication: Uses gcloud credentials and service accounts
  
Minikube:
  Prerequisites: Minikube installed and running
  Authentication: Uses local kubeconfig from minikube
  
Oracle Kubernetes Engine (OKE):
  Prerequisites: OCI CLI installed and configured
  Authentication: Uses OCI CLI credentials and IAM policies
  Required IAM policies:
    - Allow group <group> to manage cluster-family in compartment <compartment>
    - Allow group <group> to use virtual-network-family in compartment <compartment>

Use 'mesheryctl system config <platform> --help' for platform-specific details and examples.`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = `Usage: mesheryctl system config [aks|eks|gke|minikube|oke]
Example: mesheryctl system config oke
Description: Configure the Kubernetes cluster used by Meshery.

Supported platforms:
  aks      - Azure Kubernetes Service (Microsoft Azure)
  eks      - Amazon Elastic Kubernetes Service (AWS)
  gke      - Google Kubernetes Engine (Google Cloud)
  minikube - Minikube (local Kubernetes development)
  oke      - Oracle Kubernetes Engine (Oracle Cloud Infrastructure)

Prerequisites by platform:
  aks      - Azure CLI (az) installed and authenticated
  eks      - AWS CLI (aws) installed and configured
  gke      - Google Cloud SDK (gcloud) installed and authenticated
  minikube - Minikube installed and cluster running
  oke      - OCI CLI (oci) installed and configured with proper IAM policies

For platform-specific help and examples:
  mesheryctl system config <platform> --help`

		if len(args) == 0 {
			return fmt.Errorf("name of kubernetes cluster to configure Meshery not provided\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("expected one argument received multiple arguments")
		}
		return nil
	},
	Example: `
// Configure Meshery for Azure Kubernetes Service
mesheryctl system config aks

// Configure Meshery for Amazon EKS  
mesheryctl system config eks

// Configure Meshery for Google Kubernetes Engine
mesheryctl system config gke

// Configure Meshery for Minikube
mesheryctl system config minikube

// Configure Meshery for Oracle Kubernetes Engine (interactive)
mesheryctl system config oke

// Configure OKE with authentication token
mesheryctl system config oke --token auth.json

// Configure OKE with specific cluster and region (non-interactive)
mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx --region us-phoenix-1


// Configure OKE with specific cluster and region (non-interactive)
mesheryctl system config oke --cluster-id ocid1.cluster.oc1.phx.xxx --region us-phoenix-1`,

}

func init() {
	availableSubcommands = []*cobra.Command{
		aksConfigCmd,
		eksConfigCmd,
		gkeConfigCmd,
		minikubeConfigCmd,
		okeConfigCmd,
	}

	aksConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	eksConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	gkeConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	minikubeConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	okeConfigCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
	// OKE specific flags
	okeConfigCmd.Flags().StringVarP(&okeRegionFlag, "region", "r", "", "OCI region for OKE cluster")
	okeConfigCmd.Flags().StringVar(&okeClusterIDFlag, "cluster-id", "", "OKE cluster OCID")

	configCmd.AddCommand(availableSubcommands...)
}

// Given the token path, get the context and set the token in the chosen context
func setToken() {
	log.Debugf("Token path: %s", utils.TokenFlag)
	contexts, err := getContexts(utils.ConfigPath)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}
	if len(contexts) < 1 {
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