// Copyright 2023 Layer5, Inc.
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
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	apiCorev1 "k8s.io/api/core/v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// IsPodRequired checks if a given pod is specified in the required pods
func IsPodRequired(requiredPods []string, pod string) bool {
	for _, rp := range requiredPods {
		if rp == pod {
			return true
		}
	}
	return false
}

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs",
	Long: `Print history of Meshery's logs and begin tailing them.

It also shows the logs of a specific component.`,
	Args: cobra.ArbitraryArgs,
	Example: `
// Starts tailing Meshery server debug logs (works with components also)
mesheryctl system logs --verbose
mesheryctl system logs meshery-istio
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

		// switch statement for multiple platform
		switch currPlatform {
		case "docker":
			ok, err := utils.AreMesheryComponentsRunning(currPlatform)
			if err != nil {
				return err
			}
			if !ok {
				log.Error("No logs to show. Meshery is not running.")
				return nil
			}
			log.Info("Starting Meshery logging...")

			if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
				log.Errorf("%s does not exists", utils.DockerComposeFile)
				log.Info("run \"mesheryctl system start\" again to download and generate docker-compose based on your context")
				return nil
			}

			cmdlog := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "logs", "-f")

			cmdReader, err := cmdlog.StdoutPipe()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to create stdout pipe"))
			}
			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}()
			if err := cmdlog.Start(); err != nil {
				return errors.Wrap(err, utils.SystemError("failed start logger"))
			}
			if err := cmdlog.Wait(); err != nil {
				return errors.Wrap(err, utils.SystemError("failed to wait for exec process"))
			}
		case "kubernetes":
			// if the platform is kubernetes, use kubernetes go-client to
			// display pod status in the MesheryNamespace

			ok, err := utils.AreMesheryComponentsRunning(currPlatform)
			if err != nil {
				return err
			}
			if !ok {
				log.Error("No logs to show. Meshery is not running.")
				return nil
			}

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return err
			}

			// List the pods in the MesheryNamespace
			podList, err := utils.GetPodList(client, utils.MesheryNamespace)
			availablePods := podList.Items

			if err != nil {
				return err
			}

			var data []string
			var requiredPods []string

			// If the user specified logs from any particular pods, then show only that
			if len(args) > 0 {
				// Get the actual pod names even when the user specifes incomplete pod names
				requiredPods, err = utils.GetRequiredPods(args, availablePods)

				// error when the specified pod is invalid
				if err != nil {
					return err
				}
			}

			log.Info("Starting Meshery logging...")
			// List all the pods similar to kubectl get pods -n MesheryNamespace
			for _, pod := range podList.Items {

				// Get the values from the pod status
				name := pod.GetName()

				// Only print the logs from the required pods
				if len(requiredPods) > 0 {
					if !IsPodRequired(requiredPods, name) {
						continue
					}
				}

				// If a pod has multiple containers, get the logs from all the containers
				for container := range pod.Spec.Containers {
					containerName := pod.Spec.Containers[container].Name

					// Get the logs from a container within the pod
					podLogOpts := apiCorev1.PodLogOptions{
						Container: containerName,
					}

					req := client.KubeClient.CoreV1().Pods(utils.MesheryNamespace).GetLogs(name, &podLogOpts)

					logs, err := req.Stream(context.TODO())
					if err != nil {
						return err
					}
					defer logs.Close()

					buf := new(bytes.Buffer)
					_, err = io.Copy(buf, logs)
					if err != nil {
						return fmt.Errorf("error in copy information from logs to buf")
					}

					// Append this to data to be printed
					for _, str := range strings.Split(buf.String(), "\n") {
						data = append(data, fmt.Sprintf("%s\t|\t%s", name, str))
					}
					data = append(data, "\n")
				}
			}

			// Print the data
			for _, str := range data {
				log.Print(str)
			}

		}

		return nil
	},
}
