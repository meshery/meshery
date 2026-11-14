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

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"

	v1alpha1 "github.com/meshery/meshery-operator/api/v1alpha1"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes"
	controllerConfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Meshery",
	Long: `Uninstall Meshery and completely remove the Meshery installation.

Stop and remove all Meshery containers (Docker) or delete all Meshery resources
from the cluster, including custom resource instances, custom resource
definitions, the Helm release and the Meshery namespace (Kubernetes).

For Docker installations, containers are always removed forcefully, so --force
has no effect. For Kubernetes installations, --force skips the Helm release
uninstallation and forcefully deletes the remaining Meshery resources; use it
when a normal uninstall fails. A normal uninstall stops at the first failure, so
a failed run leaves the cluster untouched from that point on. --force instead
attempts every cleanup stage and reports all the failures it hit.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/system/uninstall`,
	Example: `
// Uninstall Meshery
mesheryctl system uninstall

// (optional) keep the Meshery namespace during uninstallation
mesheryctl system uninstall --keep-namespace

// Uninstall Meshery forcefully on Kubernetes (use it when system uninstall doesn't work)
mesheryctl system uninstall --force
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hc, err := NewHealthChecker(uninstallHealthCheckOptions(cmd))
		if err != nil {
			return ErrHealthCheckFailed(err)
		}
		return hc.RunPreflightHealthChecks()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemLifeCycleError(fmt.Sprintf("this command takes no arguments. See '%s --help' for more information.\n", cmd.CommandPath()), "uninstall"))
		}
		if err := uninstallMeshery(); err != nil {
			return errors.Wrap(err, utils.SystemError("failed to uninstall Meshery"))
		}
		return nil
	},
}

// uninstallHealthCheckOptions returns the preflight configuration used by
// uninstall. Teardown must work against an unhealthy or outdated installation,
// so only the Kubernetes/kubectl minimum-version checks are skipped; client
// initialization and API connectivity checks still run.
func uninstallHealthCheckOptions(cmd *cobra.Command) *HealthCheckOptions {
	return &HealthCheckOptions{
		IsPreRunE:                   true,
		PrintLogs:                   false,
		Subcommand:                  cmd.Use,
		SkipKubernetesVersionChecks: true,
	}
}

func uninstallMeshery() error {
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

	switch currCtx.GetPlatform() {
	case platformDocker:
		return uninstallDockerInstallation()
	case platformKubernetes:
		client, err := meshkitkube.New([]byte(""))
		if err != nil {
			return err
		}

		// ask user for confirmation
		userResponse := false
		if utils.SilentFlag {
			userResponse = true
		} else {
			userResponse, err = utils.AskForConfirmation("Meshery deployments will be deleted from your cluster. Are you sure you want to continue")
			if err != nil {
				return ErrUninstallMeshery(errors.Wrap(err, "cannot read the confirmation prompt"))
			}
		}

		if !userResponse {
			utils.Log.Info("Uninstall aborted.")
			return nil
		}

		utils.Log.Info("Uninstalling Meshery resources...")

		if err := uninstallKubernetesInstallation(client, currCtx); err != nil {
			return err
		}

		if utils.KeepNamespace {
			utils.Log.Info("Meshery uninstalled. The Meshery namespace was kept.")
		} else {
			utils.Log.Info("Meshery uninstalled.")
		}
	default:
		return ErrUnsupportedPlatform(currCtx.GetPlatform(), utils.CfgFile)
	}

	return nil
}

// uninstallDockerInstallation stops and removes the Meshery containers of a
// Docker installation.
//
// Compose only removes what the compose file in ~/.meshery declares, so this is
// a container-level cleanup: named volumes, downloaded images and the contents
// of ~/.meshery are left in place. Uninstall must not create ~/.meshery either,
// because a missing compose file means there is nothing to remove.
func uninstallDockerInstallation() error {
	if !dockerInstallationPresent() {
		utils.Log.Info("Meshery is not installed. Nothing to uninstall.")
		return nil
	}

	// Ask for confirmation
	userResponse := false
	if utils.SilentFlag {
		userResponse = true
	} else {
		var err error
		userResponse, err = utils.AskForConfirmation("Meshery containers will be stopped and removed. Are you sure you want to continue")
		if err != nil {
			return ErrUninstallMeshery(errors.Wrap(err, "cannot read the confirmation prompt"))
		}
	}

	if !userResponse {
		utils.Log.Info("Uninstall aborted.")
		return nil
	}

	utils.Log.Info("Uninstalling Meshery resources...")

	composeClient, err := utils.NewComposeClient()
	if err != nil {
		return utils.ErrDockerComposeClient(err)
	}

	// Stop all Docker containers
	if err := composeClient.Stop(context.Background(), utils.DockerComposeFile); err != nil {
		return utils.ErrDockerComposeStop(err)
	}

	// Remove all Docker containers. Remove forces container removal, so the
	// Kubernetes-only --force flag does not change this path.
	if err := composeClient.Remove(context.Background(), utils.DockerComposeFile); err != nil {
		return utils.ErrDockerComposeRemove(err)
	}
	utils.Log.Info("Meshery uninstalled.")

	return nil
}

// dockerInstallationPresent reports whether a Docker installation exists.
// It only inspects the filesystem: an uninstall with nothing to remove must not
// create the compose file, or anything else, while looking for it.
func dockerInstallationPresent() bool {
	_, err := os.Stat(utils.DockerComposeFile)
	return err == nil
}

// uninstallStage is one step of a cluster teardown. Describing the step lets the
// caller say which part failed without every step building the same message.
type uninstallStage struct {
	description string
	run         func() error
}

// runUninstallStages runs the teardown stages in order.
//
// Without force the first failure is returned immediately: a normal uninstall
// reports what blocked it and changes nothing after that point. With force every
// stage is attempted and the genuine failures are returned as one aggregate, so
// --force keeps making progress on a cluster that resists cleanup.
func runUninstallStages(stages []uninstallStage, force bool) error {
	var failures []error

	for _, stage := range stages {
		err := stage.run()
		if err == nil {
			continue
		}

		err = errors.Wrap(err, stage.description)
		if !force {
			return err
		}

		failures = append(failures, err)
	}

	return utilerrors.NewAggregate(failures)
}

// uninstallKubernetesInstallation removes Meshery from the current cluster.
func uninstallKubernetesInstallation(client *meshkitkube.Client, currCtx *config.Context) error {
	return runUninstallStages(kubernetesUninstallStages(client, currCtx), forceDelete)
}

// kubernetesUninstallStages returns the teardown stages in the order they run.
// The order matters: the custom resources are deleted before the Helm release,
// otherwise they can no longer be found.
func kubernetesUninstallStages(client *meshkitkube.Client, currCtx *config.Context) []uninstallStage {
	stages := []uninstallStage{
		{
			description: "cannot delete Meshery custom resources",
			run:         func() error { return deleteMesheryCRs(client) },
		},
	}

	if forceDelete {
		stages = append(stages, uninstallStage{
			description: "cannot clean up Meshery resources",
			run:         utils.ForceCleanupCluster,
		})
	} else {
		// Delete meshery components by uninstalling the helm release
		stages = append(stages, uninstallStage{
			description: "cannot uninstall the Meshery Helm release",
			run: func() error {
				return applyHelmCharts(client, currCtx, currCtx.GetVersion(), false, meshkitkube.UNINSTALL, "", "")
			},
		})
	}

	stages = append(stages,
		uninstallStage{
			description: "cannot delete Meshery custom resource definitions",
			run:         deleteMesheryCRDs,
		},
		uninstallStage{
			description: "cannot delete the Meshery namespace",
			run: func() error {
				return uninstallNamespace(utils.KeepNamespace, utils.MesheryNamespace, client.KubeClient)
			},
		},
	)

	return stages
}

// deleteMesheryCRs deletes the Broker and MeshSync custom resource instances,
// attempting both and reporting every genuine failure together. A custom
// resource that is already gone is not a failure: uninstall has to be
// re-runnable.
func deleteMesheryCRs(client *meshkitkube.Client) error {
	const (
		brokerResourceName   = "brokers"
		brokerInstanceName   = "meshery-broker"
		meshsyncResourceName = "meshsyncs"
		meshsyncInstanceName = "meshery-meshsync"
	)

	var failures []error

	for _, resource := range []struct {
		resourceName string
		instanceName string
	}{
		{brokerResourceName, brokerInstanceName},
		{meshsyncResourceName, meshsyncInstanceName},
	} {
		if err := deleteCRInstance(resource.resourceName, resource.instanceName, client); err != nil {
			failures = append(failures, errors.Wrapf(err, "cannot delete CR %s", resource.instanceName))
		}
	}

	return utilerrors.NewAggregate(failures)
}

// deleteCRInstance deletes a single custom resource instance in the Meshery
// namespace, treating an already deleted instance as success.
func deleteCRInstance(resourceName, instanceName string, client *meshkitkube.Client) error {
	err := client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace(utils.MesheryNamespace).Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

// apiExtensionsClient builds the client used to delete Meshery's CRDs. It is a
// variable so tests can substitute a fake clientset.
var apiExtensionsClient = func() (apiextension.Interface, error) {
	cfg, err := controllerConfig.GetConfig()
	if err != nil {
		return nil, errors.Wrap(err, "cannot load the Kubernetes configuration")
	}

	return apiextension.NewForConfig(cfg)
}

// deleteMesheryCRDs deletes the Broker and MeshSync CRDs, attempting both and
// reporting every genuine failure together. As with the custom resources, an
// already deleted CRD is not a failure.
func deleteMesheryCRDs() error {
	const (
		brokerCRDName   = "brokers.meshery.io"
		meshsyncCRDName = "meshsyncs.meshery.io"
	)

	client, err := apiExtensionsClient()
	if err != nil {
		return err
	}

	var failures []error
	for _, name := range []string{brokerCRDName, meshsyncCRDName} {
		err := client.ApiextensionsV1().CustomResourceDefinitions().Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil && !apierrors.IsNotFound(err) {
			failures = append(failures, errors.Wrapf(err, "cannot delete CRD %s", name))
		}
	}

	return utilerrors.NewAggregate(failures)
}

// deleteNamespace deletes the given namespace, treating an already deleted
// namespace as success.
func deleteNamespace(ns string, client kubernetes.Interface) error {
	err := client.CoreV1().Namespaces().Delete(context.TODO(), ns, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}

	return err
}

// waitForMesheryNamespaceDeletion is a variable so tests can exercise the
// failure and timeout paths without waiting on a real cluster.
var waitForMesheryNamespaceDeletion = utils.CheckMesheryNsDelete

// uninstallNamespace deletes the Meshery namespace unless the user asked to
// keep it, and waits for the namespace to actually disappear. Uninstall cannot
// report success while the namespace (and therefore Meshery's resources) is
// still present.
func uninstallNamespace(keepNamespace bool, ns string, client kubernetes.Interface) error {
	if keepNamespace {
		return nil
	}

	utils.Log.Info("Deleting Meshery Namespace...")
	if err := deleteNamespace(ns, client); err != nil {
		return ErrUninstallMeshery(errors.Wrapf(err, "cannot delete namespace %s", ns))
	}

	deleted, err := waitForMesheryNamespaceDeletion()
	if err != nil {
		return ErrUninstallMeshery(errors.Wrapf(err, "namespace %s was not deleted", ns))
	}
	if !deleted {
		return ErrUninstallMeshery(errors.Errorf("namespace %s is still present after the deletion timeout", ns))
	}

	return nil
}

func init() {
	uninstallCmd.Flags().BoolVar(&forceDelete, "force", false, "(optional) forcefully delete Kubernetes resources; has no effect for Docker installations")
	uninstallCmd.Flags().BoolVar(&utils.KeepNamespace, "keep-namespace", false, "(optional) keep the Meshery namespace during uninstallation")
}
