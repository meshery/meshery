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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: https://github.com/layer5io/me shery/issues/1022

const paramName = "k8sfile"
const kubeConfigYaml = "kubeconfig.yaml"

var tokenPath string

func getContexts(configFile, tokenPath string) ([]string, error) {
	client := &http.Client{}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, errors.Wrap(err, "error processing config")
	}

	// GETCONTEXTS endpoint points to the URL return the contexts available
	GETCONTEXTS := mctlCfg.GetBaseMesheryURL() + "/api/k8sconfig/contexts"

	req, err := utils.UploadFileWithParams(GETCONTEXTS, nil, paramName, configFile)
	if err != nil {
		return nil, err
	}

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}

	var contexts []string
	for _, item := range results {
		contexts = append(contexts, item["contextName"].(string))
	}
	return contexts, nil
}

func setContext(configFile, cname, tokenPath string) error {
	client := &http.Client{}
	extraParams1 := map[string]string{
		"contextName": cname,
	}
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// SETCONTEXT endpoint points to set context
	SETCONTEXT := mctlCfg.GetBaseMesheryURL() + "/api/k8sconfig"
	req, err := utils.UploadFileWithParams(SETCONTEXT, extraParams1, paramName, configFile)
	if err != nil {
		return err
	}
	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	// TODO: Pretty print the output
	fmt.Printf("%v\n", string(body))
	return nil
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure Meshery",
	Long:  `Configure the Kubernetes cluster used by Meshery.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if len(tokenPath) < 0 {
			log.Fatal("fetch me a token path invalid")

		}
		if tokenPath == "" {
			log.Fatal("Token path invalid")
		}
		// Define the path where the kubeconfig.yaml will be written to
		configPath := ""
		kubeConfig := ""
		usr, err := user.Current()
		if err != nil {
			configPath = filepath.Join(".meshery", kubeConfigYaml)
			kubeConfig = filepath.Join(".kube", "config")
		} else {
			configPath = filepath.Join(usr.HomeDir, ".meshery", kubeConfigYaml)
			kubeConfig = filepath.Join(usr.HomeDir, ".kube", "config")
		}

		// create the .meshery folder where the kubeconfig.yaml will be written to
		configDir := filepath.Dir(configPath)
		if _, err = os.Stat(configDir); err != nil {
			err = os.Mkdir(configDir, os.ModeDir)
			if err != nil {
				log.Fatal("Error while creating .meshery folder for config:", err)
				return
			}
		}

		switch args[0] {
		case "minikube":
			log.Info("Configuring Meshery to access Minikube...")
			// Get the config from the default config path
			if _, err = os.Stat(kubeConfig); err != nil {
				log.Fatal("Could not find the default kube config:", err)
				return
			}
			config, _ := clientcmd.LoadFromFile(kubeConfig)
			if config == nil {
				log.Fatal("Error reading the default kube config:", err)
				return
			}
			// Flatten the config file
			err = clientcmdapi.FlattenConfig(config)
			if err != nil {
				log.Fatal("Error flattening config:", err)
				return
			}
			// write the flattened config to kubeconfig.yaml file
			err = clientcmd.WriteToFile(*config, configPath)
			if err != nil {
				log.Fatal("Error writing config to file:", err)
				return
			}
			log.Debugf("Minikube configuration is written to: %s", configPath)

		case "gke":
			// TODO: move the GenerateConfigGKE logic to meshkit/client-go
			log.Info("Configuring Meshery to access GKE...")
			SAName := "sa-meshery-" + utils.StringWithCharset(8)
			if err := utils.GenerateConfigGKE(configPath, SAName, "default"); err != nil {
				log.Fatal("Error generating config:", err)
				return
			}
			log.Debugf("GKE configuration is written to: %s", configPath)
		case "aks":
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
			aksCmd := exec.Command("az", "aks", "get-credentials", "--resource-group", resourceGroup, "--name", aksName, "--file", configPath)
			aksCmd.Stdout = os.Stdout
			aksCmd.Stderr = os.Stderr
			// Write AKS compatible config to the filesystem
			err = aksCmd.Run()
			if err != nil {
				log.Fatalf("Error generating kubeconfig: %s", err.Error())
				return
			}
			log.Debugf("AKS configuration is written to: %s", configPath)
		case "eks":
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
			eksCmd := exec.Command("aws", "eks", "--region", regionName, "update-kubeconfig", "--name", clusterName, "--kubeconfig", configPath)
			eksCmd.Stdout = os.Stdout
			eksCmd.Stderr = os.Stderr
			// Write EKS compatible config to the filesystem
			err = eksCmd.Run()
			if err != nil {
				log.Fatalf("Error generating kubeconfig: %s", err.Error())
				return
			}
			log.Debugf("EKS configuration is written to: %s", configPath)
		default:
			log.Fatal("The argument has to be one of gke | minikube | aks | eks")
		}

		log.Debugf("Token path: %s", tokenPath)
		contexts, err := getContexts(configPath, tokenPath)
		if err != nil || contexts == nil || len(contexts) < 1 {
			log.Fatalf("Error getting context: %s", err.Error())
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

		log.Debugf("Chosen context : %s", choosenCtx)
		err = setContext(configPath, choosenCtx, tokenPath)
		if err != nil {
			log.Fatalf("Error setting context: %s", err.Error())
		}
	},
}

func init() {
	configCmd.Flags().StringVarP(&tokenPath, "token", "t", utils.AuthConfigFile, "Path to token for authenticating to Meshery API")
	_ = configCmd.MarkFlagRequired("tokenPath")
}
