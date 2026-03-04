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
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var verboseStatus bool

var linkDocStatus = map[string]string{
	"link":    "![status-usage](/assets/img/mesheryctl/status.png)",
	"caption": "Usage of mesheryctl system status",
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check Meshery status",
	Long:  `Check status of Meshery and Meshery components.`,
	Example: `
mesheryctl system status
mesheryctl system status --verbose
`,
	Annotations: linkDocStatus,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}

		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			utils.Log.Error(ErrHealthCheckFailed(err))
			return nil
		}

		if err := hc.RunPreflightHealthChecks(); err != nil {
			cmd.SilenceUsage = true
			return err
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(
				utils.SystemLifeCycleError(
					fmt.Sprintf(
						"this command takes no arguments. See '%s --help' for more information.",
						cmd.CommandPath(),
					),
					"status",
				),
			)
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		if tempContext != "" {
			if err := mctlCfg.SetCurrentContext(tempContext); err != nil {
				utils.Log.Error(ErrSetCurrentContext(err))
				return nil
			}
		}

		currCtx, err := mctlCfg.CheckIfCurrentContextIsValid()

		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(err))
			return nil
		}

		currPlatform := currCtx.GetPlatform()

		ok, err := utils.AreMesheryComponentsRunning(currPlatform)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		if !ok {
			utils.Log.Error(utils.ErrMesheryServerNotRunning(currPlatform))
			return nil
		}

		switch currPlatform {

		case "docker":
			composeClient, err := utils.NewComposeClient()
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to create compose client"))
			}

			output, err := composeClient.GetPsOutput(context.Background(), utils.DockerComposeFile)
			if err != nil {
				return errors.Wrap(err, utils.SystemError("failed to get Meshery status"))
			}

			if strings.Contains(output, "meshery") {
				log.Info(output)
			}

			hcOptions := &HealthCheckOptions{
				PrintLogs:           false,
				IsPreRunE:           false,
				Subcommand:          "status",
				RunKubernetesChecks: true,
			}

			hc, err := NewHealthChecker(hcOptions)
			if err == nil {
				_ = hc.Run()
			}

			fallthrough

		case "kubernetes":
			client, err := meshkitkube.New([]byte(""))
			if err != nil {
				return err
			}

			podList, err := utils.GetPodList(client, utils.MesheryNamespace)
			if err != nil {
				return err
			}

			var data [][]string
			columns := []string{"Name", "Ready", "Status", "Restarts", "Age"}

			for _, pod := range podList.Items {
				age := time.Since(pod.CreationTimestamp.Time).Round(time.Second)

				var ready, restarts int
				for _, cs := range pod.Status.ContainerStatuses {
					restarts += int(cs.RestartCount)
					if cs.Ready {
						ready++
					}
				}

				row := []string{
					utils.GetCleanPodName(pod.Name),
					fmt.Sprintf("%d/%d", ready, len(pod.Spec.Containers)),
					string(pod.Status.Phase),
					fmt.Sprintf("%d", restarts),
					age.String(),
				}

				if verboseStatus {
					row = append(row, pod.Name, pod.Status.PodIP)
				}

				data = append(data, row)
			}

			if verboseStatus {
				columns = append(columns, "Pod-Name", "Pod-IP")
			}

			utils.PrintToTable(columns, data, nil)
			log.Info("\nMeshery endpoint is " + currCtx.GetEndpoint())
		}

		return nil
	},
}

func init() {
	statusCmd.Flags().BoolVarP(&verboseStatus, "verbose", "v", false, "(optional) Extra data in status table")
}
