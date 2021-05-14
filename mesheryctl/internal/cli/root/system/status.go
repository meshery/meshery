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
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"

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

				log.Info("Meshery endpoint is " + mctlCfg.Contexts[mctlCfg.CurrentContext].Endpoint)

			} else {
				log.Info("Meshery is not running. Run `mesheryctl system start` to start Meshery.")
			}

		case "kubernetes":
			// if the platform is kubernetes, use kubernetes go-client to
			// display pod status in the MesheryNamespace

			ok, err := utils.IsMesheryRunning(currPlatform)
			if err != nil {
				return err
			}
			if !ok {
				log.Error("Meshery is not running. Run `mesheryctl system start` to start Meshery.")
				return nil
			}

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return err
			}

			// List the pods in the MesheryNamespace
			podList, err := utils.GetPods(client, utils.MesheryNamespace)

			if err != nil {
				return err
			}

			var data [][]string

			// List all the pods similar to kubectl get pods -n MesheryNamespace
			for _, pod := range podList.Items {
				// Calculate the age of the pod
				podCreationTime := pod.GetCreationTimestamp()
				age := time.Since(podCreationTime.Time).Round(time.Second)

				// Get the status of each of the pods
				podStatus := pod.Status

				var containerRestarts int32
				var containerReady int
				var totalContainers int

				// If a pod has multiple containers, get the status from all
				for container := range pod.Spec.Containers {
					containerRestarts += podStatus.ContainerStatuses[container].RestartCount
					if podStatus.ContainerStatuses[container].Ready {
						containerReady++
					}
					totalContainers++
				}

				// Get the values from the pod status
				name := pod.GetName()
				ready := fmt.Sprintf("%v/%v", containerReady, containerReady)
				status := fmt.Sprintf("%v", podStatus.Phase)
				restarts := fmt.Sprintf("%v", containerRestarts)
				ageS := age.String()

				// Append this to data to be printed in a table
				data = append(data, []string{name, ready, status, restarts, ageS})
			}

			// Print the data to a table for readability
			utils.PrintToTable([]string{"Name", "Ready", "Status", "Restarts", "Age"}, data)

			log.Info("\nMeshery endpoint is " + mctlCfg.Contexts[mctlCfg.CurrentContext].Endpoint)

		}
		return nil
	},
}
