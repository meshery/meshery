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
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
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

func printLogs(logs string, podName string) {
	for _, logMsg := range strings.Split(logs, "\n") {
		logStr := fmt.Sprintf("%s\t|\t%s", podName, logMsg)

		utils.Log.Info(logStr)
	}
}

type cmdSystemLogsFlags struct {
	Follow bool `json:"follow" validate:"boolean"`
}

var systemLogsFlags cmdSystemLogsFlags

const BYTE_SIZE = 2000

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs",
	Long: `Print history of Meshery's logs and begin tailing them.

It also shows the logs of a specific component.`,
	Args: cobra.ArbitraryArgs,
	Example: `
// Show logs (without tailing)
mesheryctl system logs

// Starts tailing Meshery server debug logs (works with components also)
mesheryctl system logs --follow
mesheryctl system logs meshery-istio
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		if err := mesheryctlflags.ValidateCmdFlags(cmd, &systemLogsFlags); err != nil {
			return err
		}
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return ErrHealthCheckFailed(err)
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
			return err
		}
		// get the platform, channel and the version of the current context
		// if a temp context is set using the -c flag, use it as the current context
		if tempContext != "" {
			err = mctlCfg.SetCurrentContext(tempContext)
			if err != nil {
				return ErrSetCurrentContext(err)
			}
		}

		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return ErrGetCurrentContext(err)
		}

		currPlatform := currCtx.GetPlatform()

		// switch statement for multiple platform
		switch currPlatform {
		case platformDocker:
			ok, err := utils.AreMesheryComponentsRunning(currPlatform)
			if err != nil {
				return err
			}
			if !ok {
				return utils.ErrMesheryServerNotRunning(currPlatform)
			}
			utils.Log.Info("Starting Meshery logging...")

			if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
				return utils.ErrDockerComposeFileMissing(err)
			}

			// Use compose library instead of exec.Command
			composeClient, err := utils.NewComposeClient()
			if err != nil {
				return utils.ErrDockerComposeClient(err)
			}

			if err := composeClient.Logs(context.Background(), utils.DockerComposeFile, true, os.Stdout); err != nil {
				return errors.Wrap(err, utils.SystemError("failed to get logs"))
			}
		case platformKubernetes:
			// if the platform is kubernetes, use kubernetes go-client to
			// display pod status in the MesheryNamespace

			ok, err := utils.AreMesheryComponentsRunning(currPlatform)
			if err != nil {
				return err
			}
			if !ok {
				return utils.ErrMesheryServerNotRunning(currPlatform)
			}

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return err
			}

			// Get and display current context
			currentContext, err := utils.GetCurrentK8sContext(client)
			if err != nil {
				utils.Log.Warn(err)
			} else {
				utils.Log.Info("Using Kubernetes context: ", currentContext)
			}

			// List the pods in the MesheryNamespace
			podList, err := utils.GetPodList(client, utils.MesheryNamespace)
			if err != nil {
				return err
			}

			availablePods := podList.Items

			requiredPodsMap := make(map[string]string)

			// If the user specified logs from any particular pods, then show only that
			if len(args) > 0 {
				// Get the actual pod names even when the user specifes incomplete pod names
				requiredPodsMap, err = utils.GetRequiredPods(args, availablePods)
				// error when the specified pod is invalid
				if err != nil {
					return err
				}
			}

			utils.Log.Info("Starting Meshery logging...")
			wg := &sync.WaitGroup{}

			// List all the pods similar to kubectl get pods -n MesheryNamespace
			for _, pod := range podList.Items {

				// Get the values from the pod status
				name := pod.GetName()

				// Only print the logs from the required pods
				if len(requiredPodsMap) > 0 {
					_, ok := requiredPodsMap[name]
					if !ok {
						continue
					}
				}

				// If a pod has multiple containers, get the logs from all the containers
				for container := range pod.Spec.Containers {
					containerName := pod.Spec.Containers[container].Name

					// Get the logs from a container within the pod
					podLogOpts := apiCorev1.PodLogOptions{
						Container: containerName,
						Follow:    systemLogsFlags.Follow,
					}

					req := client.KubeClient.CoreV1().Pods(utils.MesheryNamespace).GetLogs(name, &podLogOpts)

					logs, err := req.Stream(context.TODO())
					if err != nil {
						return err
					}
					defer func() { _ = logs.Close() }()
					var logBuf []byte
					if !systemLogsFlags.Follow {
						logBuf, err = io.ReadAll(logs)
						if err != nil {
							return utils.ErrReadResponseBody(err)
						}
						printLogs(string(logBuf), name)
					} else {
						wg.Add(1)
						go func() {
							defer wg.Done()
							for {
								buf := make([]byte, BYTE_SIZE)
								numBytes, err := logs.Read(buf)
								if numBytes == 0 {
									continue
								}
								if err == io.EOF {
									break
								}
								if err != nil {
									utils.Log.Errorf("error occurred while processing logs %v", err)
									break
								}
								logBuf = buf[0:numBytes]
								printLogs(string(logBuf), name)
							}
						}()
					}
				}
			}
			wg.Wait()
		}
		return nil
	},
}

func init() {
	logsCmd.Flags().BoolVarP(&systemLogsFlags.Follow, "follow", "f", false, "(Optional) Follow the stream of the Meshery's logs. Defaults to false.")
}
