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
	"os"
	"strconv"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	controllerConfig "sigs.k8s.io/controller-runtime/pkg/client/config"

	"github.com/meshery/meshery-operator/api/v1alpha1"
)

const desiredReplicasAnnotation = "meshery.io/desired-replicas"

var (
	// forceDelete used to clean-up meshery resources forcefully (shared with system uninstall)
	forceDelete bool
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Meshery",
	Long: `Stop Meshery without uninstalling it.

Docker: stops Meshery containers (does not remove them).
Kubernetes: scales Meshery deployments to 0 and keeps the Helm release,
custom resources, CRDs, and namespace.

Use "mesheryctl system uninstall" to fully remove Meshery resources.
	Find more information at: https://docs.meshery.io/reference/references/mesheryctl/system/stop`,
	Example: `
// Stop Meshery
mesheryctl system stop

// Reset Meshery's configuration file to default settings.
mesheryctl system stop --reset
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
			return ErrHealthCheckFailed(err)
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "stop"))
		}
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

	ok, err := utils.AreMesheryComponentsRunning(currCtx.GetPlatform())
	if err != nil {
		return err
	}
	if !ok {
		utils.Log.Info("Meshery resources are not running. Nothing to stop.")
		return nil
	}

	switch currCtx.GetPlatform() {
	case platformDocker:
		if _, err := os.Stat(utils.MesheryFolder); os.IsNotExist(err) {
			if err := os.Mkdir(utils.MesheryFolder, 0777); err != nil {
				return ErrCreateDir(err, utils.MesheryFolder)
			}
		}

		utils.Log.Info("Stopping Meshery...")

		composeClient, err := utils.NewComposeClient()
		if err != nil {
			return utils.ErrDockerComposeClient(err)
		}

		// Stop only — do not remove containers (uninstall owns removal)
		if err := composeClient.Stop(context.Background(), utils.DockerComposeFile); err != nil {
			return utils.ErrDockerComposeStop(err)
		}
		utils.Log.Info("Meshery stopped. Containers retained. Use `mesheryctl system uninstall` to remove Meshery.")
	case platformKubernetes:
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			userResponse = utils.AskForConfirmation("Meshery deployments will be scaled down (stopped) but not uninstalled. Are you sure you want to continue")
		}

		if !userResponse {
			utils.Log.Info("Stop aborted.")
			return nil
		}

		utils.Log.Info("Stopping Meshery (scaling deployments to 0)...")

		if err := scaleMesheryDeployments(client, 0, true); err != nil {
			return ErrStopMeshery(err)
		}

		utils.Log.Info("Meshery stopped. Resources retained (Helm release, CRs, CRDs, namespace). Use `mesheryctl system uninstall` to remove Meshery.")
	default:
		return ErrUnsupportedPlatform(currCtx.GetPlatform(), utils.CfgFile)
	}

	// Reset Meshery config file to default settings
	if utils.ResetFlag {
		err := resetMesheryConfig()
		if err != nil {
			return ErrResetMeshconfig(err)
		}
	}
	return nil
}

// scaleMesheryDeployments scales Meshery-related deployments in the Meshery namespace.
// When saveDesired is true and targetReplicas is 0, the current non-zero replica count
// is stored in the meshery.io/desired-replicas annotation so start can resume later.
func scaleMesheryDeployments(client *meshkitkube.Client, targetReplicas int32, saveDesired bool) error {
	deploymentInterface := client.KubeClient.AppsV1().Deployments(utils.MesheryNamespace)
	deploymentList, err := deploymentInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	scaled := 0
	for _, deployment := range deploymentList.Items {
		if !strings.Contains(deployment.GetName(), "meshery") {
			continue
		}

		dep, err := deploymentInterface.Get(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if saveDesired && targetReplicas == 0 {
			desired := int32(1)
			if dep.Spec.Replicas != nil && *dep.Spec.Replicas > 0 {
				desired = *dep.Spec.Replicas
			}
			if dep.Annotations == nil {
				dep.Annotations = map[string]string{}
			}
			dep.Annotations[desiredReplicasAnnotation] = strconv.Itoa(int(desired))
			dep, err = deploymentInterface.Update(context.TODO(), dep, metav1.UpdateOptions{})
			if err != nil {
				return err
			}
		}

		if targetReplicas > 0 {
			if dep.Annotations != nil {
				if v, ok := dep.Annotations[desiredReplicasAnnotation]; ok {
					if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
						targetReplicas = int32(parsed)
					}
				}
			}
		}

		replicas := targetReplicas
		dep.Spec.Replicas = &replicas
		if _, err := deploymentInterface.Update(context.TODO(), dep, metav1.UpdateOptions{}); err != nil {
			return err
		}
		scaled++
		utils.Log.Debug(fmt.Sprintf("Scaled deployment %s to %d replicas", dep.Name, targetReplicas))
	}

	if scaled == 0 {
		return errors.New("no Meshery deployments found to scale in namespace " + utils.MesheryNamespace)
	}
	return nil
}

// invokeDeleteCRs is a wrapper of deleteCR to delete CR instances (brokers and meshsyncs)
// Kept here for use by system uninstall (see PR #22230 / related work).
func invokeDeleteCRs(client *meshkitkube.Client) error {
	const (
		brokerResourceName   = "brokers"
		brokerInstanceName   = "meshery-broker"
		meshsyncResourceName = "meshsyncs"
		meshsyncInstanceName = "meshery-meshsync"
	)

	if err := deleteCR(brokerResourceName, brokerInstanceName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CR "+brokerInstanceName))
		if !forceDelete {
			return err
		}

		utils.Log.Debug(err)
	}

	if err := deleteCR(meshsyncResourceName, meshsyncInstanceName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CR "+meshsyncInstanceName))
		if !forceDelete {
			return err
		}

		utils.Log.Debug(err)
	}

	return nil
}

// deleteCR delete the specified CR instance in the clusters
func deleteCR(resourceName, instanceName string, client *meshkitkube.Client) error {
	return client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace(utils.MesheryNamespace).Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
}

// invokeDeleteCRDs is a wrapper of deleteCRD to delete CRDs (brokers and meshsyncs)
func invokeDeleteCRDs() error {
	const (
		brokerCRDName   = "brokers.meshery.io"
		meshsyncCRDName = "meshsyncs.meshery.io"
	)

	cfg := controllerConfig.GetConfigOrDie()
	client, err := apiextension.NewForConfig(cfg)
	if err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot invoke delete CRDs"))
		if !forceDelete {
			return err
		}

		utils.Log.Debug(err)
	}

	if err = deleteCRD(brokerCRDName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CRD "+brokerCRDName))
		if !forceDelete {
			return err
		}

		utils.Log.Debug(err)
	}

	if err = deleteCRD(meshsyncCRDName, client); err != nil {
		err = ErrStopMeshery(errors.Wrap(err, "cannot delete CRD "+meshsyncCRDName))
		if !forceDelete {
			return err
		}

		utils.Log.Debug(err)
	}

	return nil
}

// deleteCRD delete the specified CRD in the clusters
func deleteCRD(name string, client *apiextension.Clientset) error {
	return client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func deleteNs(ns string, client *kubernetes.Clientset) error {
	return client.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})
}

func init() {
	stopCmd.Flags().BoolVarP(&utils.ResetFlag, "reset", "", false, "(optional) reset Meshery's configuration file to default settings.")
}
