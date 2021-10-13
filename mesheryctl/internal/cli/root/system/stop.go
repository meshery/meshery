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
	"os"
	"os/exec"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"

	"github.com/layer5io/meshery-operator/api/v1alpha1"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long:  `Stop all Meshery containers / remove all Meshery pods.`,
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
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
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := stop(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop Meshery"))
		}
		return nil
	},
}

func stop() error {
	// Get viper instance used for context
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// if a temp context is set using the -c flag, use it as the current context
	err = mctlCfg.SetCurrentContext(tempContext)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve current-context")
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return err
	}

	ok, err := utils.IsMesheryRunning(currCtx.GetPlatform())
	if err != nil {
		return err
	}
	if !ok {
		log.Info("Meshery is not running. Nothing to stop.")
		return nil
	}

	switch currCtx.GetPlatform() {
	case "docker":
		// if the platform is docker, then stop all the running containers
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

		log.Info("Stopping Meshery...")

		// Stop all Docker containers
		stop := exec.Command("docker-compose", "-f", utils.DockerComposeFile, "stop")
		stop.Stdout = os.Stdout
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to stop meshery - could not stop some containers."))
		}

		// Remove all Docker containers
		stop = exec.Command("docker-compose", "-f", utils.DockerComposeFile, "rm", "-f")
		stop.Stderr = os.Stderr

		if err := stop.Run(); err != nil {
			return ErrStopMeshery(err)
		}
	case "kubernetes":
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}
		// if the platform is kubernetes, stop the deployment by deleting the manifest files
		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			// ask user for confirmation
			userResponse = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
		}

		if !userResponse {
			log.Info("Stop aborted.")
			return nil
		}

		log.Info("Stopping Meshery...")

		// Delete the CR instances for brokers and meshsyncs
		if err = invokeDeleteCRs(client); err != nil {
			return err
		}

		// TODO: Delete the CRDs for brokers and meshsyncs

		// Delete the helm chart installation
		// Note: this doesn't delete the CRDs (broker and meshsync)
		if err = client.ApplyHelmChart(meshkitkube.ApplyHelmChartConfig{
			Namespace: utils.MesheryNamespace,
			ChartLocation: meshkitkube.HelmChartLocation{
				Repository: utils.HelmChartURL,
				Chart:      utils.HelmChartName,
			},
			Action: meshkitkube.UNINSTALL,
		}); err != nil {
			return errors.Wrap(err, "cannot stop Meshery")
		}
	}

	// If k8s is available in case of platform docker than we remove operator
	hcOptions := &HealthCheckOptions{
		PrintLogs:           false,
		IsPreRunE:           false,
		Subcommand:          "",
		RunKubernetesChecks: true,
	}
	hc, err := NewHealthChecker(hcOptions)
	if err != nil {
		return ErrHealthCheckFailed(err)
	}

	if err = hc.Run(); err != nil {
		return ErrHealthCheckFailed(err)
	}

	log.Info("Meshery is stopped.")

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}
	return nil
}

// invokeDeleteCRs is a wrapper of deleteCR to delete CR instances (brokers and meshsyncs)
func invokeDeleteCRs(client *meshkitkube.Client) error {
	const (
		brokerResourceName   = "brokers"
		brokerInstanceName   = "meshery-broker"
		meshsyncResourceName = "meshsyncs"
		meshsyncInstanceName = "meshery-meshsync"
	)

	if err := deleteCR(brokerResourceName, brokerInstanceName, client); err != nil {
		return errors.Wrap(err, "cannot delete CR "+brokerInstanceName)
	}

	if err := deleteCR(meshsyncResourceName, meshsyncInstanceName, client); err != nil {
		return errors.Wrap(err, "cannot delete CR "+meshsyncInstanceName)
	}

	return nil
}
// deleteCRs delete the specified CR instance in the clusters
func deleteCR(resourceName, instanceName string, client *meshkitkube.Client) error {
	return client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace(utils.MesheryNamespace).Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
