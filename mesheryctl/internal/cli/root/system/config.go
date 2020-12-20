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
	"path"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

// TODO: https://github.com/layer5io/meshery/issues/1022

// GETCONTEXTS endpoint points to the URL return the contexts available
const GETCONTEXTS = "http://localhost:9081/api/k8sconfig/contexts"

// SETCONTEXT endpoint points to set context
const SETCONTEXT = "http://localhost:9081/api/k8sconfig"

const paramName = "k8sfile"
const contextName = "contextName"

var tokenPath string

type k8sContext struct {
	contextName    string
	clusterName    string
	currentContext bool
}

func getContexts(configFile, tokenPath string) ([]string, error) {
	client := &http.Client{}

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

		switch args[0] {
		case "minikube":
			if err := utils.GenerateConfigMinikube(); err != nil {
				log.Fatal("Error generating config:", err)
				return
			}
		case "gke":
			SAName := "sa-meshery-" + utils.StringWithCharset(8)
			if err := utils.GenerateConfigGKE(SAName, "default"); err != nil {
				log.Fatal("Error generating config:", err)
				return
			}
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

			// Write AKS compatible config to the filesystem
			if err := utils.GenerateConfigAKS(resourceGroup, aksName); err != nil {
				log.Fatal("Error generating kubeconfig: ", err)
				return
			}
		case "eks":
			log.Info("Configuring Meshery to access EKS...")
			var regionName, clusterName string

			// Prompt user for AWS region name
			log.Info("Please enter the AWS region name:")
			_, err := fmt.Scanf("%s", &regionName)
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

			// Write EKS compatible config to the filesystem
			if err := utils.GenerateConfigEKS(regionName, clusterName); err != nil {
				log.Fatal("Error generating kubeconfig: ", err)
				return
			}
		default:
			log.Fatal("The argument has to be one of gke | minikube | aks | eks")
		}

		// TODO: Assumes Mac or Linux. Make arch-specific
		// Issue: https://github.com/layer5io/meshery/issues/1894
		configPath := ""
		usr, err := user.Current()
		if err != nil {
			configPath = ".meshery/kubeconfig.yaml"
		} else {
			configPath = path.Join(usr.HomeDir, ".meshery/kubeconfig.yaml")
		}

		log.Info(tokenPath)
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
