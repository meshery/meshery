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

var verboseStatus bool

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Meshery status",
	Args:  cobra.NoArgs,
	Long:  `Check status of Meshery and Meshery components.`,
	Example: `
// Check status of Meshery, Meshery adapters, Meshery Operator and its controllers.
mesheryctl system status 
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
			return errors.Wrapf(err, "failed to initialize healthchecker")
		}
		// execute healthchecks
		err = hc.RunPreflightHealthChecks()
		if err != nil {
			cmd.SilenceUsage = true
		}

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
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

		currPlatform := currCtx.GetPlatform()

		ok, err := utils.AreMesheryComponentsRunning(currPlatform)
		if err != nil {
			return err
		}
		if !ok {
			log.Error("Meshery is not running. Run `mesheryctl system start` to start Meshery.")
			return nil
		}

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
			}

			hcOptions := &HealthCheckOptions{
				PrintLogs:           false,
				IsPreRunE:           false,
				Subcommand:          "status",
				RunKubernetesChecks: true,
			}
			hc, err := NewHealthChecker(hcOptions)
			if err != nil {
				return errors.Wrapf(err, "failed to initialize healthchecker")

			}
			// If k8s is available print the status of pods in the MesheryNamespace
			if err = hc.Run(); err != nil {
				return nil
			}

			fallthrough
		case "kubernetes":
			// if the platform is kubernetes, use kubernetes go-client to
			// display pod status in the MesheryNamespace

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return err
			}

			// List the pods in the MesheryNamespace
			podList, err := utils.GetPodList(client, utils.MesheryNamespace)

			if err != nil {
				return err
			}

			var data [][]string
			columnNames := []string{"Name", "Ready", "Status", "Restarts", "Age"}
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

				if len(pod.Spec.Containers) > 0 && len(podStatus.ContainerStatuses) > 0 {
					// If a pod has multiple containers, get the status from all
					for container := range pod.Spec.Containers {
						containerRestarts += podStatus.ContainerStatuses[container].RestartCount
						if podStatus.ContainerStatuses[container].Ready {
							containerReady++
						}
						totalContainers++
					}
				}

				// Get the values from the pod status
				name := utils.GetCleanPodName(pod.GetName())
				ready := fmt.Sprintf("%v/%v", containerReady, containerReady)
				status := fmt.Sprintf("%v", podStatus.Phase)
				restarts := fmt.Sprintf("%v", containerRestarts)
				ageS := age.String()
				row := []string{name, ready, status, restarts, ageS}

				// Append this to data to be printed in a table
				if verboseStatus {
					row = append(row, pod.Name)
					row = append(row, podStatus.PodIP)
				}
				data = append(data, row)
			}
			if verboseStatus {
				columnNames = append(columnNames, "Pod-Names")
				columnNames = append(columnNames, "Pod-IP")
			}
			// Print the data to a table for readability
			utils.PrintToTable(columnNames, data)

			log.Info("\nMeshery endpoint is " + currCtx.GetEndpoint())
		}
		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&verboseStatus, "verbose", "v", false, "(optional) Extra data in status table")
}
