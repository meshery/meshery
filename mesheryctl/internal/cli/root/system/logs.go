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
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Print logs",
	Long:  `Print history of Meshery's container logs and begin tailing them.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if ok := utils.IsMesheryRunning(); !ok {
			log.Error("No logs to show. Meshery is not running.")
			return nil
		}

		log.Info("Starting Meshery logging...")

		if _, err := os.Stat(utils.DockerComposeFile); os.IsNotExist(err) {
			log.Errorf("%s does not exists", utils.DockerComposeFile)
			log.Info("run \"mesheryctl system start\" again to download and generate docker-compose based on your context")
			return nil
		}

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

		var cmdlog *exec.Cmd
		// switch statement for multiple platform
		switch currPlatform {
		case "docker":
			cmdlog = exec.Command("docker-compose", "-f", utils.DockerComposeFile, "logs", "-f")

		case "kubernetes":
			// Create a new client
			client, err := meshkitkube.New([]byte(""))
			if err != nil {
				return errors.Wrap(err, "failed to create new client")
			}

			// Create a deployment interface for the MesheryNamespace
			deploymentInterface := client.KubeClient.AppsV1().Deployments(utils.MesheryNamespace)

			// List the deployments in the MesheryNamespace
			deploymentList, err := deploymentInterface.List(context.TODO(), v1.ListOptions{})

			if err != nil {
				return err
			}

			// List logs for all the deployments similar to kubectl logs MesheryNamespace
			for _, deployment := range deploymentList.Items {

			}

			cmdlog = exec.Command("kubectl logs meshery")

		}
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
		return nil
	},
}
