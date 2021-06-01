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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Stop, then start Meshery",
	Long:  `Restart all Meshery containers/ pods.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return restart()
	},
}

func restart() error {
	log.Info("Restarting Meshery...")

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
		if err := stop(); err != nil {
			return errors.Wrap(err, utils.SystemError("Failed to restart Meshery"))
		}

		if err := start(); err != nil {
			return errors.Wrap(err, utils.SystemError("Failed to restart Meshery"))
		}

	case "kubernetes":
		// create an kubernetes client
		client, err := meshkitkube.New([]byte(""))

		if err != nil {
			return err
		}

		// Create a pod interface for the MesheryNamespace
		podInterface := client.KubeClient.CoreV1().Pods(utils.MesheryNamespace)

		// List the pods in the MesheryNamespace
		podList, err := podInterface.List(context.TODO(), v1.ListOptions{})
		if err != nil {
			return err
		}

		// List all the pods similar to kubectl get pods -n MesheryNamespace
		for _, pod := range podList.Items {
			// Get the values from the pod status
			name := pod.GetName()
			log.Info("Deleting pod ", name)
			err := client.KubeClient.CoreV1().Pods(utils.MesheryNamespace).Delete(context.TODO(), name, v1.DeleteOptions{})
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Restarting pod ", name)
		}
	}
	return nil
}

func init() {
	restartCmd.Flags().BoolVarP(&skipUpdateFlag, "skip-update", "", false, "(optional) skip checking for new Meshery's container images.")
}
