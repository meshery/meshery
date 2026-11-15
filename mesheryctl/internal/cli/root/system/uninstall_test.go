package system

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	mesheryoperator "github.com/meshery/meshery-operator/api/v1alpha1"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mesheryerrors "github.com/meshery/meshkit/errors"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apiextensionfake "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/fake"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/kubernetes"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	k8stesting "k8s.io/client-go/testing"
)

// TestUninstallCmdRegistration verifies that the uninstall command is wired
// into `mesheryctl system` and exposes the flags that drive the teardown
// behavior.
func TestUninstallCmdRegistration(t *testing.T) {
	assert.Contains(t, availableSubcommands, uninstallCmd)
	assert.Equal(t, "uninstall", uninstallCmd.Use)

	// --force drives the ForceCleanupCluster path and the best-effort teardown.
	assert.NotNil(t, uninstallCmd.Flags().Lookup("force"))

	// --keep-namespace skips the Meshery namespace deletion.
	assert.NotNil(t, uninstallCmd.Flags().Lookup("keep-namespace"))
}

// TestUninstallCmdHelp verifies the help output of the uninstall command. Help
// short-circuits before the preflight health checks, so this test runs without
// a Docker daemon or a Kubernetes cluster.
func TestUninstallCmdHelp(t *testing.T) {
	buf := setupSystemOutCmdTest(t)
	defer buf.Reset()

	SystemCmd.SetArgs([]string{"uninstall", "--help"})
	err = SystemCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Uninstall Meshery")
	assert.Contains(t, output, "--keep-namespace")
	assert.Contains(t, output, "--force")

	// Docker cleanup already forces container removal, so --force only changes
	// Kubernetes teardown. The help output must say so.
	assert.Contains(t, output, "has no effect for Docker installations")
}

// TestUninstallDockerErrorRemediation verifies that a Docker Compose removal
// failure points at Docker remediation and never at the Kubernetes-only
// --force flag.
func TestUninstallDockerErrorRemediation(t *testing.T) {
	err := utils.ErrDockerComposeRemove(errors.New("docker compose remove failed"))

	var meshkitErr *mesheryerrors.Error
	require.ErrorAs(t, err, &meshkitErr)
	assert.Equal(t, utils.ErrDockerComposeRemoveCode, meshkitErr.Code)
	assert.Equal(t, "Failed to delete Meshery containers", meshkitErr.ShortDescription[0])

	remediation := strings.Join(meshkitErr.SuggestedRemediation, " ")
	assert.Contains(t, remediation, "Docker")
	assert.NotContains(t, remediation, "--force")
}

// TestUninstallKubernetesErrorRemediation verifies that Kubernetes teardown
// failures keep the --force remediation, because --force is the flag that
// changes the Kubernetes uninstall path.
func TestUninstallKubernetesErrorRemediation(t *testing.T) {
	err := ErrUninstallMeshery(errors.New("helm uninstall failed"))

	var meshkitErr *mesheryerrors.Error
	require.ErrorAs(t, err, &meshkitErr)
	assert.Equal(t, ErrUninstallMesheryCode, meshkitErr.Code)

	remediation := strings.Join(meshkitErr.SuggestedRemediation, " ")
	assert.Contains(t, remediation, "mesheryctl system uninstall --force")
	assert.Contains(t, remediation, "cluster")
	assert.NotContains(t, remediation, "Docker")
}

// TestSystemUninstallErrorReference verifies the "uninstall" error formatting
// points at the uninstall command reference page.
func TestSystemUninstallErrorReference(t *testing.T) {
	msg := utils.SystemLifeCycleError("this command takes no arguments.\n", "uninstall")

	assert.Contains(t, msg, "this command takes no arguments")
	assert.Contains(t, msg, "https://docs.meshery.io/reference/mesheryctl/system/uninstall")
}

// TestUninstallHealthCheckOptions pins the teardown preflight: version checks
// are skipped so an outdated cluster or kubectl cannot block an uninstall, but
// the PreRunE mode (and therefore the API connectivity checks) is preserved.
func TestUninstallHealthCheckOptions(t *testing.T) {
	options := uninstallHealthCheckOptions(uninstallCmd)

	assert.True(t, options.IsPreRunE, "uninstall must keep running preflight checks")
	assert.False(t, options.PrintLogs, "uninstall must not switch to report-style checks")
	assert.Equal(t, "uninstall", options.Subcommand)
	assert.True(t, options.SkipKubernetesVersionChecks, "version checks must be skipped for teardown")
}

// TestUninstallCmdInSubcommandList guards the subcommand registration against
// accidental removal that would silently change `mesheryctl system --help`.
func TestUninstallCmdInSubcommandList(t *testing.T) {
	names := make([]string, 0, len(availableSubcommands))
	for _, cmd := range availableSubcommands {
		names = append(names, cmd.Use)
	}

	assert.True(t, slices.Contains(names, "uninstall"))
}

// TestUninstallForceFlag verifies that --force is wired to the global that
// switches the Kubernetes teardown to best-effort.
func TestUninstallForceFlag(t *testing.T) {
	previous := forceDelete
	t.Cleanup(func() { forceDelete = previous })

	require.NoError(t, uninstallCmd.Flags().Set("force", "true"))
	t.Cleanup(func() { _ = uninstallCmd.Flags().Set("force", "false") })

	assert.True(t, forceDelete)
}

// TestUninstallKeepNamespaceFlag verifies that --keep-namespace is wired to the
// global that skips namespace deletion.
func TestUninstallKeepNamespaceFlag(t *testing.T) {
	previous := utils.KeepNamespace
	t.Cleanup(func() { utils.KeepNamespace = previous })

	require.NoError(t, uninstallCmd.Flags().Set("keep-namespace", "true"))
	t.Cleanup(func() { _ = uninstallCmd.Flags().Set("keep-namespace", "false") })

	assert.True(t, utils.KeepNamespace)
}

// newTestStage builds a teardown stage that records that it ran.
func newTestStage(ran *[]string, description string, failure error) uninstallStage {
	return uninstallStage{
		description: description,
		run: func() error {
			*ran = append(*ran, description)
			return failure
		},
	}
}

// TestRunUninstallStages pins the difference between a normal uninstall and a
// forced one: without force the first failure ends the teardown, with force
// every stage still runs and the genuine failures come back together.
func TestRunUninstallStages(t *testing.T) {
	build := func() ([]uninstallStage, *[]string) {
		ran := &[]string{}
		return []uninstallStage{
			newTestStage(ran, "first", nil),
			newTestStage(ran, "second", errors.New("boom")),
			newTestStage(ran, "third", errors.New("bang")),
		}, ran
	}

	t.Run("normal mode stops at the first failure", func(t *testing.T) {
		stages, ran := build()

		err := runUninstallStages(stages, false)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "second")
		assert.NotContains(t, err.Error(), "third")
		assert.Equal(t, []string{"first", "second"}, *ran, "stages after the failure must not run")
	})

	t.Run("force mode runs every stage and aggregates the failures", func(t *testing.T) {
		stages, ran := build()

		err := runUninstallStages(stages, true)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "second")
		assert.Contains(t, err.Error(), "third")
		assert.Equal(t, []string{"first", "second", "third"}, *ran)

		aggregate, ok := err.(interface{ Errors() []error })
		require.True(t, ok, "expected an aggregated error, got %T", err)
		assert.Len(t, aggregate.Errors(), 2)
	})

	t.Run("reports success when every stage succeeds", func(t *testing.T) {
		ran := &[]string{}
		stages := []uninstallStage{
			newTestStage(ran, "first", nil),
			newTestStage(ran, "second", nil),
		}

		assert.NoError(t, runUninstallStages(stages, true))
		assert.Equal(t, []string{"first", "second"}, *ran)
	})
}

// TestKubernetesUninstallStageOrder pins the teardown order. The custom
// resources have to go before the Helm release, and the namespace last, so that
// a failing early stage cannot strand the later ones.
func TestKubernetesUninstallStageOrder(t *testing.T) {
	previousForceDelete := forceDelete
	t.Cleanup(func() { forceDelete = previousForceDelete })

	descriptions := func(stages []uninstallStage) []string {
		names := make([]string, 0, len(stages))
		for _, stage := range stages {
			names = append(names, stage.description)
		}
		return names
	}

	forceDelete = false
	assert.Equal(t, []string{
		"cannot delete Meshery custom resources",
		"cannot uninstall the Meshery Helm release",
		"cannot delete Meshery custom resource definitions",
		"cannot delete the Meshery namespace",
	}, descriptions(kubernetesUninstallStages(nil, nil)))

	forceDelete = true
	assert.Equal(t, []string{
		"cannot delete Meshery custom resources",
		"cannot clean up Meshery resources",
		"cannot delete Meshery custom resource definitions",
		"cannot delete the Meshery namespace",
	}, descriptions(kubernetesUninstallStages(nil, nil)))
}

// TestDockerInstallationPresent verifies the check that decides whether there is
// a Docker installation to remove.
func TestDockerInstallationPresent(t *testing.T) {
	originalComposeFile := utils.DockerComposeFile
	t.Cleanup(func() { utils.DockerComposeFile = originalComposeFile })

	utils.DockerComposeFile = filepath.Join(t.TempDir(), "meshery.yaml")
	assert.False(t, dockerInstallationPresent(), "a missing compose file means no installation")

	require.NoError(t, os.WriteFile(utils.DockerComposeFile, []byte("services: {}\n"), 0o600))
	assert.True(t, dockerInstallationPresent())
}

// TestSkipPreflightForPlatform pins when the uninstall preflight is skipped: a
// Docker installation that is not there has nothing to stop, so a missing
// daemon must not block the no-op. A Kubernetes teardown always keeps its
// preflight, because it is the only thing that verifies API connectivity.
func TestSkipPreflightForPlatform(t *testing.T) {
	tests := []struct {
		name                     string
		platform                 string
		dockerInstallationOnDisk bool
		want                     bool
	}{
		{"docker platform without a compose file", platformDocker, false, true},
		{"docker platform with a compose file", platformDocker, true, false},
		{"kubernetes platform without a compose file", platformKubernetes, false, false},
		{"kubernetes platform with a compose file", platformKubernetes, true, false},
		{"unknown platform", "unknown", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, skipPreflightForPlatform(tt.platform, tt.dockerInstallationOnDisk))
		})
	}
}

// TestUninstallDockerInstallationWithoutComposeFile verifies that uninstalling a
// Docker installation that is not there is a no-op. In particular it must not
// create ~/.meshery just to look for the compose file.
func TestUninstallDockerInstallationWithoutComposeFile(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	originalFolder, originalComposeFile := utils.MesheryFolder, utils.DockerComposeFile
	t.Cleanup(func() {
		utils.MesheryFolder, utils.DockerComposeFile = originalFolder, originalComposeFile
	})

	utils.MesheryFolder = filepath.Join(t.TempDir(), "meshery")
	utils.DockerComposeFile = filepath.Join(utils.MesheryFolder, "meshery.yaml")

	require.NoError(t, uninstallDockerInstallation())

	_, err := os.Stat(utils.MesheryFolder)
	assert.True(t, os.IsNotExist(err), "uninstall must leave %s untouched", utils.MesheryFolder)
}

// TestDeleteNamespace pins the namespace deletion contract: a namespace that is
// already gone is success, a rejected deletion is a failure.
func TestDeleteNamespace(t *testing.T) {
	t.Run("returns an error when deletion is rejected", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		client.PrependReactor("delete", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: "namespaces"},
				action.(k8stesting.DeleteAction).GetName(), errors.New("namespace deletion is forbidden"))
		})

		err := deleteNamespace(utils.MesheryNamespace, client)
		require.Error(t, err)
		assert.True(t, apierrors.IsForbidden(err))
	})

	t.Run("treats an already deleted namespace as success", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset()

		assert.NoError(t, deleteNamespace(utils.MesheryNamespace, client))
	})

	t.Run("deletes an existing namespace", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})

		require.NoError(t, deleteNamespace(utils.MesheryNamespace, client))

		_, err := client.CoreV1().Namespaces().Get(t.Context(), utils.MesheryNamespace, metav1.GetOptions{})
		assert.True(t, apierrors.IsNotFound(err))
	})
}

// TestUninstallNamespace pins the behavior around --keep-namespace and the wait
// for the namespace to actually disappear.
func TestUninstallNamespace(t *testing.T) {
	// uninstallNamespace logs through utils.Log, which tests must initialize.
	utils.SetupMeshkitLoggerTesting(t, false)

	// stubNamespaceWait replaces the real 300s wait and restores it afterwards.
	stubNamespaceWait := func(t *testing.T, fn func() (bool, error)) {
		t.Helper()
		previous := waitForMesheryNamespaceDeletion
		waitForMesheryNamespaceDeletion = fn
		t.Cleanup(func() { waitForMesheryNamespaceDeletion = previous })
	}

	assertUninstallError := func(t *testing.T, err error) {
		t.Helper()
		require.Error(t, err)
		var meshkitErr *mesheryerrors.Error
		require.ErrorAs(t, err, &meshkitErr)
		assert.Equal(t, ErrUninstallMesheryCode, meshkitErr.Code)
	}

	t.Run("keeps the namespace when requested", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		waited := false
		stubNamespaceWait(t, func() (bool, error) {
			waited = true
			return true, nil
		})

		require.NoError(t, uninstallNamespace(true, utils.MesheryNamespace, client))
		assert.False(t, waited, "--keep-namespace must not wait for namespace deletion")

		_, err := client.CoreV1().Namespaces().Get(t.Context(), utils.MesheryNamespace, metav1.GetOptions{})
		assert.NoError(t, err, "--keep-namespace must leave the namespace in place")
	})

	t.Run("fails when namespace deletion is rejected", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		client.PrependReactor("delete", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: "namespaces"},
				action.(k8stesting.DeleteAction).GetName(), errors.New("namespace deletion is forbidden"))
		})
		waited := false
		stubNamespaceWait(t, func() (bool, error) {
			waited = true
			return true, nil
		})

		assertUninstallError(t, uninstallNamespace(false, utils.MesheryNamespace, client))
		assert.False(t, waited, "a rejected deletion must not wait on the namespace")
	})

	t.Run("fails when the namespace does not disappear within the timeout", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		stubNamespaceWait(t, func() (bool, error) {
			return false, errors.New("timed out waiting for the condition")
		})

		assertUninstallError(t, uninstallNamespace(false, utils.MesheryNamespace, client))
	})

	t.Run("fails when the namespace is reported still present", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		stubNamespaceWait(t, func() (bool, error) {
			return false, nil
		})

		assertUninstallError(t, uninstallNamespace(false, utils.MesheryNamespace, client))
	})

	t.Run("succeeds when the namespace disappears", func(t *testing.T) {
		client := k8sfake.NewSimpleClientset(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: utils.MesheryNamespace}})
		stubNamespaceWait(t, func() (bool, error) {
			return true, nil
		})

		require.NoError(t, uninstallNamespace(false, utils.MesheryNamespace, client))
	})
}

func newDynamicClient(objects ...runtime.Object) *dynamicfake.FakeDynamicClient {
	return dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), objects...)
}

func newCRInstanceOfKind(name, kind string) *unstructured.Unstructured {
	return &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": mesheryoperator.GroupVersion.String(),
		"kind":       kind,
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": utils.MesheryNamespace,
		},
	}}
}

// forbidCRDelete makes every deletion of resource fail with Forbidden and counts
// the attempts.
func forbidCRDelete(client *dynamicfake.FakeDynamicClient, resource string, attempts *int) {
	client.PrependReactor("delete", resource, func(action k8stesting.Action) (bool, runtime.Object, error) {
		*attempts++
		return true, nil, apierrors.NewForbidden(
			schema.GroupResource{Resource: resource},
			action.(k8stesting.DeleteAction).GetName(),
			errors.New(resource+" deletion is forbidden"))
	})
}

// crExists reports whether a custom resource instance is still present.
func crExists(t *testing.T, client *meshkitkube.Client, resource, name string) bool {
	t.Helper()

	_, err := client.DynamicKubeClient.Resource(schema.GroupVersionResource{
		Group:    mesheryoperator.GroupVersion.Group,
		Version:  mesheryoperator.GroupVersion.Version,
		Resource: resource,
	}).Namespace(utils.MesheryNamespace).Get(t.Context(), name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		return false
	}
	require.NoError(t, err)

	return true
}

// TestDeleteCRInstance pins the contract of a single custom resource deletion: a
// resource that is already gone is success, a rejected deletion is a failure.
func TestDeleteCRInstance(t *testing.T) {
	const (
		brokerResource = "brokers"
		brokerInstance = "meshery-broker"
	)

	t.Run("treats an already deleted custom resource as success", func(t *testing.T) {
		client := &meshkitkube.Client{DynamicKubeClient: newDynamicClient()}

		assert.NoError(t, deleteCRInstance(brokerResource, brokerInstance, client))
	})

	t.Run("returns genuine deletion failures", func(t *testing.T) {
		dynamicClient := newDynamicClient(newCRInstanceOfKind(brokerInstance, "Broker"))
		attempts := 0
		forbidCRDelete(dynamicClient, brokerResource, &attempts)
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		err := deleteCRInstance(brokerResource, brokerInstance, client)
		require.Error(t, err)
		assert.True(t, apierrors.IsForbidden(err), "Forbidden must not be swallowed: %v", err)
	})
}

// withShortCRDeletionWait shortens the bounded wait for a custom resource
// instance to disappear so the timeout path is testable, and restores it after.
func withShortCRDeletionWait(t *testing.T, interval, timeout time.Duration) {
	t.Helper()

	previousInterval, previousTimeout := crDeletionPollInterval, crDeletionTimeout
	crDeletionPollInterval, crDeletionTimeout = interval, timeout
	t.Cleanup(func() {
		crDeletionPollInterval, crDeletionTimeout = previousInterval, previousTimeout
	})
}

// TestDeleteCRInstanceWaitsForDeletion pins the contract the teardown order
// depends on: a deletion is only complete once the instance is gone, so the
// operator that clears the instance's finalizers is still running when
// deleteCRInstance returns. An accepted delete request is not enough.
func TestDeleteCRInstanceWaitsForDeletion(t *testing.T) {
	const (
		brokerResource = "brokers"
		brokerInstance = "meshery-broker"
	)

	// keepInstance answers every Get with the instance, as the API server does
	// while the instance is still terminating.
	keepInstance := func(client *dynamicfake.FakeDynamicClient, gets *int) {
		client.PrependReactor("get", brokerResource, func(action k8stesting.Action) (bool, runtime.Object, error) {
			*gets++
			return true, newCRInstanceOfKind(brokerInstance, "Broker"), nil
		})
	}

	t.Run("polls until the instance is gone", func(t *testing.T) {
		dynamicClient := newDynamicClient(newCRInstanceOfKind(brokerInstance, "Broker"))
		withShortCRDeletionWait(t, time.Millisecond, time.Second)

		// The instance outlives the delete request by two polls.
		gets := 0
		dynamicClient.PrependReactor("get", brokerResource, func(action k8stesting.Action) (bool, runtime.Object, error) {
			gets++
			if gets <= 2 {
				return true, newCRInstanceOfKind(brokerInstance, "Broker"), nil
			}
			return false, nil, nil
		})
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		require.NoError(t, deleteCRInstance(brokerResource, brokerInstance, client))
		assert.Equal(t, 3, gets, "deletion must not be reported as complete while the instance is still present")
	})

	t.Run("fails with a bounded timeout when the instance never disappears", func(t *testing.T) {
		dynamicClient := newDynamicClient(newCRInstanceOfKind(brokerInstance, "Broker"))
		gets := 0
		keepInstance(dynamicClient, &gets)
		withShortCRDeletionWait(t, time.Millisecond, 20*time.Millisecond)
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		err := deleteCRInstance(brokerResource, brokerInstance, client)
		require.Error(t, err, "an instance that never disappears must not be reported as deleted")
		assert.Contains(t, err.Error(), "cannot confirm deletion of "+brokerInstance)
		// A timeout is allowed to land after the single immediate poll: only one
		// invocation is guaranteed. Repeated polling is pinned by the
		// disappearance test above.
		assert.GreaterOrEqual(t, gets, 1, "the wait must poll at least once")
	})

	t.Run("reports an API error instead of waiting for the timeout", func(t *testing.T) {
		dynamicClient := newDynamicClient(newCRInstanceOfKind(brokerInstance, "Broker"))
		gets := 0
		dynamicClient.PrependReactor("get", brokerResource, func(action k8stesting.Action) (bool, runtime.Object, error) {
			gets++
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: brokerResource}, brokerInstance, errors.New("get is forbidden"))
		})
		// A long timeout proves the failure comes from the API error, not the wait.
		withShortCRDeletionWait(t, time.Millisecond, time.Minute)
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		err := deleteCRInstance(brokerResource, brokerInstance, client)
		require.Error(t, err)
		assert.Equal(t, 1, gets, "a genuine API error must not be retried as if the instance were present")

		var statusErr *apierrors.StatusError
		require.ErrorAs(t, err, &statusErr)
		assert.True(t, apierrors.IsForbidden(statusErr), "the original API error must survive the wait")
	})

	t.Run("does not wait for an instance that is already gone", func(t *testing.T) {
		withShortCRDeletionWait(t, time.Millisecond, 20*time.Millisecond)
		client := &meshkitkube.Client{DynamicKubeClient: newDynamicClient()}

		assert.NoError(t, deleteCRInstance(brokerResource, brokerInstance, client))
	})
}

// TestDeleteMesheryCRs verifies that both custom resources are attempted and
// that every genuine failure is reported, so --force can still make progress on
// a cluster where one of them is stuck.
func TestDeleteMesheryCRs(t *testing.T) {
	t.Run("deletes both custom resources", func(t *testing.T) {
		client := &meshkitkube.Client{DynamicKubeClient: newDynamicClient(
			newCRInstanceOfKind("meshery-broker", "Broker"),
			newCRInstanceOfKind("meshery-meshsync", "MeshSync"),
		)}

		require.NoError(t, deleteMesheryCRs(client))
		assert.False(t, crExists(t, client, "brokers", "meshery-broker"))
		assert.False(t, crExists(t, client, "meshsyncs", "meshery-meshsync"))
	})

	t.Run("succeeds when both custom resources are already gone", func(t *testing.T) {
		client := &meshkitkube.Client{DynamicKubeClient: newDynamicClient()}

		assert.NoError(t, deleteMesheryCRs(client))
	})

	t.Run("attempts both custom resources and reports both failures", func(t *testing.T) {
		dynamicClient := newDynamicClient(
			newCRInstanceOfKind("meshery-broker", "Broker"),
			newCRInstanceOfKind("meshery-meshsync", "MeshSync"),
		)
		attempts := 0
		forbidCRDelete(dynamicClient, "brokers", &attempts)
		forbidCRDelete(dynamicClient, "meshsyncs", &attempts)
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		err := deleteMesheryCRs(client)
		require.Error(t, err)
		assert.Equal(t, 2, attempts, "both custom resources must be attempted")
		assert.Contains(t, err.Error(), "meshery-broker")
		assert.Contains(t, err.Error(), "meshery-meshsync")

		aggregate, ok := err.(interface{ Errors() []error })
		require.True(t, ok, "expected an aggregated error, got %T", err)
		assert.Len(t, aggregate.Errors(), 2)
	})

	t.Run("reports a failing custom resource while the other is deleted", func(t *testing.T) {
		dynamicClient := newDynamicClient(
			newCRInstanceOfKind("meshery-broker", "Broker"),
			newCRInstanceOfKind("meshery-meshsync", "MeshSync"),
		)
		attempts := 0
		forbidCRDelete(dynamicClient, "brokers", &attempts)
		client := &meshkitkube.Client{DynamicKubeClient: dynamicClient}

		err := deleteMesheryCRs(client)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "meshery-broker")
		assert.NotContains(t, err.Error(), "meshery-meshsync")
		assert.False(t, crExists(t, client, "meshsyncs", "meshery-meshsync"), "the other custom resource must still be deleted")
		assert.True(t, crExists(t, client, "brokers", "meshery-broker"))
	})
}

// TestDeleteMesheryCRDs verifies CRD deletion: both CRDs are attempted, an
// already deleted CRD is success, and a client that cannot be built is reported
// instead of panicking.
func TestDeleteMesheryCRDs(t *testing.T) {
	previousClient := apiExtensionsClient
	t.Cleanup(func() { apiExtensionsClient = previousClient })

	withClient := func(t *testing.T, client apiextension.Interface) {
		t.Helper()
		apiExtensionsClient = func(*rest.Config) (apiextension.Interface, error) { return client, nil }
		t.Cleanup(func() { apiExtensionsClient = previousClient })
	}

	newCRDClient := func(names ...string) *apiextensionfake.Clientset {
		crds := make([]runtime.Object, 0, len(names))
		for _, name := range names {
			crds = append(crds, &apiextensionsv1.CustomResourceDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: name},
			})
		}

		return apiextensionfake.NewClientset(crds...)
	}

	t.Run("deletes both CRDs", func(t *testing.T) {
		client := newCRDClient("brokers.meshery.io", "meshsyncs.meshery.io")
		withClient(t, client)

		require.NoError(t, deleteMesheryCRDs(&rest.Config{}))

		for _, name := range []string{"brokers.meshery.io", "meshsyncs.meshery.io"} {
			_, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(t.Context(), name, metav1.GetOptions{})
			assert.True(t, apierrors.IsNotFound(err), "%s must be deleted", name)
		}
	})

	t.Run("succeeds when both CRDs are already gone", func(t *testing.T) {
		withClient(t, newCRDClient())

		assert.NoError(t, deleteMesheryCRDs(&rest.Config{}))
	})

	t.Run("reports every genuine failure", func(t *testing.T) {
		client := newCRDClient("brokers.meshery.io", "meshsyncs.meshery.io")
		attempts := 0
		client.PrependReactor("delete", "customresourcedefinitions", func(action k8stesting.Action) (bool, runtime.Object, error) {
			attempts++
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: "customresourcedefinitions"},
				action.(k8stesting.DeleteAction).GetName(),
				errors.New("CRD deletion is forbidden"))
		})
		withClient(t, client)

		err := deleteMesheryCRDs(&rest.Config{})
		require.Error(t, err)
		assert.Equal(t, 2, attempts, "both CRDs must be attempted")
		assert.Contains(t, err.Error(), "brokers.meshery.io")
		assert.Contains(t, err.Error(), "meshsyncs.meshery.io")
	})

	t.Run("deletes through the configuration it is given", func(t *testing.T) {
		cfg := &rest.Config{Host: "https://cluster-a.example:6443"}
		var got *rest.Config
		apiExtensionsClient = func(cfg *rest.Config) (apiextension.Interface, error) {
			got = cfg
			return newCRDClient(), nil
		}
		t.Cleanup(func() { apiExtensionsClient = previousClient })

		require.NoError(t, deleteMesheryCRDs(cfg))
		assert.Same(t, cfg, got, "the CRD client must be built from the supplied configuration")
	})

	t.Run("reports a client that cannot be built", func(t *testing.T) {
		apiExtensionsClient = func(*rest.Config) (apiextension.Interface, error) {
			return nil, errors.New("no kubeconfig")
		}
		t.Cleanup(func() { apiExtensionsClient = previousClient })

		err := deleteMesheryCRDs(&rest.Config{})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no kubeconfig")
	})
}

// helmReleaseSecret builds the Secret Helm's Secret storage driver keeps a
// release revision in: the labels matter, the payload does not.
func helmReleaseSecret(releaseName string) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaseName + ".v1",
			Namespace: utils.MesheryNamespace,
			Labels: map[string]string{
				"owner": "helm",
				"name":  releaseName,
			},
		},
	}
}

// TestHelmReleaseExists checks the lookup uninstall uses to tell an installed
// release from one a previous attempt already removed.
func TestHelmReleaseExists(t *testing.T) {
	t.Run("reports a stored release", func(t *testing.T) {
		client := k8sfake.NewClientset(helmReleaseSecret(mesheryReleaseName))

		installed, err := helmReleaseExists(client, utils.MesheryNamespace, mesheryReleaseName)

		require.NoError(t, err)
		assert.True(t, installed)
	})

	t.Run("does not count another release", func(t *testing.T) {
		client := k8sfake.NewClientset(helmReleaseSecret(mesheryOperatorReleaseName))

		installed, err := helmReleaseExists(client, utils.MesheryNamespace, mesheryReleaseName)

		require.NoError(t, err)
		assert.False(t, installed, "a different release must not be mistaken for this one")
	})

	t.Run("reports an absent release", func(t *testing.T) {
		installed, err := helmReleaseExists(k8sfake.NewClientset(), utils.MesheryNamespace, mesheryReleaseName)

		require.NoError(t, err)
		assert.False(t, installed)
	})

	t.Run("reports a failed lookup", func(t *testing.T) {
		client := k8sfake.NewClientset()
		client.PrependReactor("list", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: "secrets"}, utils.MesheryNamespace, errors.New("list is forbidden"))
		})

		_, err := helmReleaseExists(client, utils.MesheryNamespace, mesheryReleaseName)

		require.Error(t, err)
	})
}

// TestUninstallHelmReleases pins the re-runnability of the Helm stage: a release
// that is already gone is not a failure, so a retry after an interrupted
// uninstall reaches the CRD and namespace stages instead of stopping here.
func TestUninstallHelmReleases(t *testing.T) {
	previousApply := applyHelmChart
	t.Cleanup(func() { applyHelmChart = previousApply })

	type call struct {
		releaseName string
		chartName   string
	}

	previousLookup := isHelmReleaseInstalled
	t.Cleanup(func() { isHelmReleaseInstalled = previousLookup })

	// installedSpec substitutes the release lookup with a fixed answer per
	// release, and optionally a lookup failure.
	installedSpec := func(t *testing.T, installed ...string) {
		t.Helper()
		isHelmReleaseInstalled = func(_ kubernetes.Interface, namespace, releaseName string) (bool, error) {
			assert.Equal(t, utils.MesheryNamespace, namespace, "releases must be looked up where they were installed")
			return slices.Contains(installed, releaseName), nil
		}
		t.Cleanup(func() { isHelmReleaseInstalled = previousLookup })
	}

	// recordUninstalls substitutes a Helm uninstall, failing the given releases.
	recordUninstalls := func(t *testing.T, failures map[string]error) *[]call {
		t.Helper()
		calls := new([]call)
		applyHelmChart = func(_ *meshkitkube.Client, cfg meshkitkube.ApplyHelmChartConfig) error {
			assert.Equal(t, meshkitkube.UNINSTALL, cfg.Action, "the Helm stage must only uninstall")
			*calls = append(*calls, call{cfg.ReleaseName, cfg.ChartLocation.Chart})
			return failures[cfg.ReleaseName]
		}
		t.Cleanup(func() { applyHelmChart = previousApply })
		return calls
	}

	releases := []string{mesheryReleaseName, mesheryOperatorReleaseName}
	currCtx := &config.Context{Version: "v0.8.0"}

	t.Run("succeeds when every release is already absent", func(t *testing.T) {
		calls := recordUninstalls(t, nil)
		installedSpec(t)

		require.NoError(t, uninstallHelmReleases(nil, nil, currCtx))
		assert.Empty(t, *calls, "an absent release must not be uninstalled")
	})

	t.Run("uninstalls the release that is still installed", func(t *testing.T) {
		calls := recordUninstalls(t, nil)
		installedSpec(t, mesheryReleaseName)

		require.NoError(t, uninstallHelmReleases(nil, nil, currCtx), "an absent release must not fail the stage")
		assert.Equal(t, []call{{mesheryReleaseName, utils.HelmChartName}}, *calls)
	})

	t.Run("uninstalls both releases", func(t *testing.T) {
		calls := recordUninstalls(t, nil)
		installedSpec(t, releases...)

		require.NoError(t, uninstallHelmReleases(nil, nil, currCtx))
		assert.Equal(t, []call{
			{mesheryReleaseName, utils.HelmChartName},
			{mesheryOperatorReleaseName, utils.HelmChartOperatorName},
		}, *calls)
	})

	t.Run("reports a genuine uninstall failure", func(t *testing.T) {
		calls := recordUninstalls(t, map[string]error{
			mesheryReleaseName: errors.New("the release could not be uninstalled"),
		})
		installedSpec(t, releases...)

		err := uninstallHelmReleases(nil, nil, currCtx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "the release could not be uninstalled", "the Helm error must be preserved")
		assert.Contains(t, err.Error(), mesheryReleaseName)
		assert.Len(t, *calls, len(releases), "a failing release must not stop the other release")
	})

	t.Run("reports a release lookup failure", func(t *testing.T) {
		recordUninstalls(t, nil)
		kubeClient := k8sfake.NewClientset()
		kubeClient.PrependReactor("list", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
			return true, nil, apierrors.NewForbidden(
				schema.GroupResource{Resource: "secrets"}, utils.MesheryNamespace, errors.New("list is forbidden"))
		})
		err := uninstallHelmReleases(kubeClient, nil, currCtx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot tell whether the")
	})
}

// TestKubernetesUninstallStagesUseOneConfiguration verifies that the Helm and CRD
// stages work against the configuration the command resolved once, instead of
// re-reading the kubeconfig and possibly targeting a different cluster.
func TestKubernetesUninstallStagesUseOneConfiguration(t *testing.T) {
	previousForceDelete, previousApply := forceDelete, applyHelmChart
	previousClient := apiExtensionsClient
	t.Cleanup(func() {
		forceDelete = previousForceDelete
		applyHelmChart = previousApply
		apiExtensionsClient = previousClient
	})
	forceDelete = false

	stage := func(t *testing.T, stages []uninstallStage, description string) uninstallStage {
		t.Helper()
		for _, candidate := range stages {
			if candidate.description == description {
				return candidate
			}
		}
		t.Fatalf("the %q stage is missing", description)
		return uninstallStage{}
	}

	client := &meshkitkube.Client{RestConfig: rest.Config{Host: "https://cluster-a.example:6443"}}

	// A release that is already absent must not be uninstalled, which also keeps
	// the Helm stage from needing the chart repository at all.
	t.Run("the Helm stage skips absent releases", func(t *testing.T) {
		previousLookup := isHelmReleaseInstalled
		isHelmReleaseInstalled = func(kubernetes.Interface, string, string) (bool, error) { return false, nil }
		t.Cleanup(func() { isHelmReleaseInstalled = previousLookup })
		applyHelmChart = func(*meshkitkube.Client, meshkitkube.ApplyHelmChartConfig) error {
			t.Fatal("an absent release must not be uninstalled")
			return nil
		}
		stages := kubernetesUninstallStages(&meshkitkube.Client{}, &config.Context{})

		require.NoError(t, stage(t, stages, "cannot uninstall the Meshery Helm release").run())
	})

	t.Run("the CRD stage deletes through the resolved configuration", func(t *testing.T) {
		var got *rest.Config
		apiExtensionsClient = func(cfg *rest.Config) (apiextension.Interface, error) {
			got = cfg
			return apiextensionfake.NewClientset(), nil
		}

		stages := kubernetesUninstallStages(client, &config.Context{})
		require.NoError(t, stage(t, stages, "cannot delete Meshery custom resource definitions").run())

		require.NotNil(t, got, "the CRD stage must build its client from a configuration")
		assert.Same(t, &client.RestConfig, got,
			"CRD deletion must reuse the configuration every other stage targets")
	})
}
