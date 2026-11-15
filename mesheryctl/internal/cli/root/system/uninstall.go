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
	"path"
	"time"

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
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
		// A Docker installation that is not there has nothing to stop, so a
		// missing daemon must not turn the no-op into a preflight failure.
		if skipPreflightForAbsentDockerInstallation() {
			return nil
		}
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

// skipPreflightForAbsentDockerInstallation reports whether the preflight can be
// skipped because the current context is a Docker installation that has no
// compose file to remove.
//
// The Docker preflight fails, or even tries to start the daemon, whenever Docker
// is not reachable, which would make `system uninstall` refuse to report that
// there is nothing to uninstall. Only the Docker platform is affected, the
// Kubernetes preflight is never skipped, and a context that cannot be resolved
// falls through to the normal preflight so its own error is reported.
func skipPreflightForAbsentDockerInstallation() bool {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return false
	}
	if err := mctlCfg.SetCurrentContext(tempContext); err != nil {
		return false
	}
	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		return false
	}

	platform := currCtx.GetPlatform()
	if utils.PlatformFlag != "" {
		if utils.PlatformFlag != platformDocker && utils.PlatformFlag != platformKubernetes {
			// Let the preflight report the unsupported platform.
			return false
		}
		platform = utils.PlatformFlag
	}

	return skipPreflightForPlatform(platform, dockerInstallationPresent())
}

// skipPreflightForPlatform reports whether the preflight can be skipped for a
// resolved platform: only a Docker installation that is absent from disk has
// nothing to remove, and a Kubernetes preflight is never skipped.
func skipPreflightForPlatform(platform string, dockerInstallationPresent bool) bool {
	return platform == platformDocker && !dockerInstallationPresent
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
		// The Kubernetes configuration is resolved here, once, and every teardown
		// stage works from this client: the custom resources through its dynamic
		// client, the CRDs through its REST configuration and the namespace through
		// its clientset. Nothing reloads the kubeconfig afterwards, so a context
		// switch mid-uninstall cannot point a later stage at another cluster.
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
			// The stages say what failed; classify the failure as an uninstall
			// failure so its error code and recovery instructions reach the user.
			return ErrUninstallMeshery(err)
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
//
// The order is a dependency chain, not a preference. The custom resources are
// deleted, and their deletion confirmed, before the Helm release, because that
// release is what removes the operator able to clear their finalizers; the
// namespace goes last so its deletion covers whatever is left. A normal
// uninstall stops at the first stage that fails, which is what keeps the cluster
// recoverable when a stage could not finish. --force deliberately breaks the
// chain: it exists for the case where a stage cannot complete at all, so it runs
// every stage and aggregates the failures instead of stopping.
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
		// Delete meshery components by uninstalling the helm releases
		stages = append(stages, uninstallStage{
			description: "cannot uninstall the Meshery Helm release",
			run:         func() error { return uninstallHelmReleases(client.KubeClient, client, currCtx) },
		})
	}

	stages = append(stages,
		uninstallStage{
			description: "cannot delete Meshery custom resource definitions",
			// The configuration resolved when the command started is reused, so
			// every stage targets the cluster the teardown began on.
			run: func() error { return deleteMesheryCRDs(&client.RestConfig) },
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

// The Meshery Helm releases, and the labels Helm's Secret storage driver puts on
// every release revision it keeps in the cluster.
const (
	mesheryReleaseName         = "meshery"
	mesheryOperatorReleaseName = "meshery-operator"

	helmReleaseOwnerLabel      = "owner"
	helmReleaseOwnerLabelValue = "helm"
	helmReleaseNameLabel       = "name"
)

// applyHelmChart uninstalls a single Helm release. It is a variable so tests can
// observe which releases are uninstalled, and with what, without a cluster.
var applyHelmChart = func(client *meshkitkube.Client, cfg meshkitkube.ApplyHelmChartConfig) error {
	return client.ApplyHelmChart(cfg)
}

// uninstallHelmReleases uninstalls the Meshery and Meshery Operator Helm
// releases, treating a release Helm no longer stores as successfully
// uninstalled.
//
// Each release is uninstalled on its own so that one already-absent release
// cannot mask the outcome of the other, and the genuine failures are reported
// together the way applyHelmCharts reports them for install. Uninstall has to be
// re-runnable: an earlier attempt can remove the releases and then fail while
// deleting the custom resource definitions, and the retry must reach the
// remaining stages instead of stopping here.
func uninstallHelmReleases(kubeClient kubernetes.Interface, chartClient *meshkitkube.Client, currCtx *config.Context) error {
	var failures []error

	for _, release := range []struct {
		releaseName string
		chartName   string
	}{
		{mesheryReleaseName, utils.HelmChartName},
		{mesheryOperatorReleaseName, utils.HelmChartOperatorName},
	} {
		if err := uninstallHelmRelease(kubeClient, chartClient, currCtx, release.releaseName, release.chartName); err != nil {
			failures = append(failures, errors.Wrapf(err, "cannot uninstall the %s Helm release", release.releaseName))
		}
	}

	return utilerrors.NewAggregate(failures)
}

// uninstallHelmRelease uninstalls one Meshery Helm release, skipping the
// uninstall when the release is not installed.
//
// MeshKit's uninstall action leaves Helm's IgnoreNotFound unset, so uninstalling
// an absent release fails with Helm's "release: not found", and MeshKit renders
// that cause into the message of a generic apply-chart error instead of wrapping
// it, leaving nothing to match with errors.Is. The release is therefore looked up
// before it is uninstalled instead of its failure being recognised afterwards.
func uninstallHelmRelease(kubeClient kubernetes.Interface, chartClient *meshkitkube.Client, currCtx *config.Context, releaseName, chartName string) error {
	installed, err := isHelmReleaseInstalled(kubeClient, utils.MesheryNamespace, releaseName)
	if err != nil {
		return errors.Wrapf(err, "cannot tell whether the %s release is installed", releaseName)
	}
	if !installed {
		return nil
	}

	// The install path passes no chart version for "latest", which makes MeshKit
	// resolve the latest published chart; reuse whatever the context pins.
	chartVersion := ""
	if version := currCtx.GetVersion(); version != "latest" {
		chartVersion = version
	}

	return applyHelmChart(chartClient, meshkitkube.ApplyHelmChartConfig{
		Namespace:       utils.MesheryNamespace,
		ReleaseName:     releaseName,
		CreateNamespace: true,
		ChartLocation: meshkitkube.HelmChartLocation{
			Repository: utils.HelmChartURL,
			Chart:      chartName,
			Version:    chartVersion,
		},
		Action: meshkitkube.UNINSTALL,
		// Keep the download location the install path uses, so a chart an
		// earlier start already downloaded is reused rather than fetched again.
		DownloadLocation: path.Join(utils.MesheryFolder, utils.ManifestsFolder),
	})
}

// isHelmReleaseInstalled is the lookup uninstallHelmRelease performs. It is a
// variable so the Helm stage can be exercised without a cluster.
var isHelmReleaseInstalled = helmReleaseExists

// helmReleaseExists reports whether Helm still stores a release of the given
// name in the namespace.
//
// Helm records every release revision as a Secret in the release namespace - the
// Secret driver is MeshKit's default - labelled owner=helm, name=<release>. That
// label pair is the same one Helm's storage query uses to decide that a release
// does not exist, so asking the cluster for it answers the question uninstall
// needs without matching error strings.
func helmReleaseExists(client kubernetes.Interface, namespace, releaseName string) (bool, error) {
	secrets, err := client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s,%s=%s",
			helmReleaseOwnerLabel, helmReleaseOwnerLabelValue,
			helmReleaseNameLabel, releaseName),
	})
	if err != nil {
		return false, err
	}

	return len(secrets.Items) > 0, nil
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

// crDeletionPollInterval and crDeletionTimeout bound the wait for a deleted
// custom resource instance to disappear. They are variables so tests can
// exercise the timeout without waiting for the real one.
var (
	crDeletionPollInterval = 2 * time.Second
	crDeletionTimeout      = 2 * time.Minute
)

// deleteCRInstance deletes a single custom resource instance in the Meshery
// namespace and waits until it is really gone, treating an already deleted
// instance as success.
//
// An accepted delete request is not a completed deletion. The Broker and
// MeshSync instances carry finalizers that the Meshery Operator clears while it
// is still running, so the instance can stay behind in Terminating for as long
// as the operator needs to finish. Returning as soon as the request is accepted
// would let the next stage remove the operator, and with it the only thing able
// to clear those finalizers: the instances then stay stuck and the namespace
// never becomes deletable.
func deleteCRInstance(resourceName, instanceName string, client *meshkitkube.Client) error {
	resource := client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    v1alpha1.GroupVersion.Group,
		Version:  v1alpha1.GroupVersion.Version,
		Resource: resourceName,
	}).Namespace(utils.MesheryNamespace)

	err := resource.Delete(context.TODO(), instanceName, metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	err = wait.PollUntilContextTimeout(context.TODO(), crDeletionPollInterval, crDeletionTimeout, true,
		func(ctx context.Context) (bool, error) {
			_, err := resource.Get(ctx, instanceName, metav1.GetOptions{})
			if apierrors.IsNotFound(err) {
				return true, nil
			}
			// A successful Get means the instance is still there, whatever its
			// deletion timestamp: keep waiting. Any other error is a genuine API
			// failure and must be reported rather than retried as though the
			// instance were merely still present.
			return false, err
		})
	if err != nil {
		return errors.Wrapf(err, "cannot confirm deletion of %s within %s", instanceName, crDeletionTimeout)
	}

	return nil
}

// apiExtensionsClient builds the client used to delete Meshery's CRDs from the
// REST configuration the rest of the teardown already resolved. It is a variable
// so tests can substitute a fake clientset.
var apiExtensionsClient = func(cfg *rest.Config) (apiextension.Interface, error) {
	return apiextension.NewForConfig(cfg)
}

// deleteMesheryCRDs deletes the Broker and MeshSync CRDs, attempting both and
// reporting every genuine failure together. As with the custom resources, an
// already deleted CRD is not a failure.
func deleteMesheryCRDs(cfg *rest.Config) error {
	const (
		brokerCRDName   = "brokers.meshery.io"
		meshsyncCRDName = "meshsyncs.meshery.io"
	)

	client, err := apiExtensionsClient(cfg)
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
