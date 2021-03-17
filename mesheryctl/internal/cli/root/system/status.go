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
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Meshery status",
	Args:  cobra.NoArgs,
	Long:  `Check status of Meshery and Meshery adapters.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info("Meshery status... \n")

		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		currCtx, err := mctlCfg.SetCurrentContext(tempContext)
		if err != nil {
			return err
		}
		currPlatform := currCtx.Platform

		switch currPlatform {
		case "docker":
			// List the running Meshery containers
			start := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "ps")

			outputStd, err := start.Output()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to get Meshery status"))
			}

			outputString := string(outputStd)
			if strings.Contains(outputString, "meshery") {
				log.Info(outputString)
			} else {
				log.Info("Meshery is not running, run `mesheryctl system start` to start Meshery")
			}

		case "kubernetes":
			// if the platform is kubernetes, use kubernetes go-client to
			// display pod status in the MesheryNamespace

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return err
			}

			// Create a pod interface for the MesheryNamespace
			// podInterface := client.KubeClient.CoreV1().Pods(utils.MesheryNamespace)

			// Create a deployment interface for the MesheryNamespace
			deploymentInterface := client.KubeClient.AppsV1().Deployments(utils.MesheryNamespace)

			// List the pods in the MesheryNamespace
			// podList, err := podInterface.List(context.TODO(), v1.ListOptions{})

			// List the deployments in the MesheryNamespace
			deploymentList, err := deploymentInterface.List(context.TODO(), v1.ListOptions{})

			if err != nil {
				return err
			}

			var data [][]string

			// List all the deployments similar to kubectl get deployments -n MesheryNamespace
			for _, deployment := range deploymentList.Items {

				// Calculate the age of the deployment
				deploymentCreationTime := deployment.GetCreationTimestamp()
				age := time.Since(deploymentCreationTime.Time).Round(time.Second)

				// Get the status of each of the deployments
				deploymentStatus := deployment.Status

				// Get the values from the deployment status
				name := deployment.GetName()
				ready := fmt.Sprintf("%d/%d", deploymentStatus.ReadyReplicas, deploymentStatus.Replicas)
				updated := fmt.Sprintf("%d", deploymentStatus.UpdatedReplicas)
				available := fmt.Sprintf("%d", deploymentStatus.AvailableReplicas)
				ageS := age.String()

				// Append this to data to be printed in a table
				data = append(data, []string{name, ready, updated, available, ageS})
			}

			// Print the data to a table for readability
			utils.PrintToTable([]string{"Name", "Ready", "Up-to-date", "Available", "Age"}, data)

			// List all the pods
			// for i, pod := range podList.Items {
			// 	// Get the status from all the pods
			// 	podstatusPhase := string(pod.Status.Phase)
			// 	podCreationTime := pod.GetCreationTimestamp()
			// 	age := time.Since(podCreationTime.Time).Round(time.Second)

			// 	// Log the status
			// 	podInfo := fmt.Sprintf("[%d] Pod: %s, Phase: %s , Created: %s, Age: %s", i, pod.GetName(), podstatusPhase, podCreationTime, age.String())
			// 	fmt.Println(podInfo)
			// }
		}
		return nil
	},
}
