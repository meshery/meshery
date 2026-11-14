package system

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	mesheryoperator "github.com/meshery/meshery-operator/api/v1alpha1"
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
	k8sfake "k8s.io/client-go/kubernetes/fake"
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
		apiExtensionsClient = func() (apiextension.Interface, error) { return client, nil }
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

		require.NoError(t, deleteMesheryCRDs())

		for _, name := range []string{"brokers.meshery.io", "meshsyncs.meshery.io"} {
			_, err := client.ApiextensionsV1().CustomResourceDefinitions().Get(t.Context(), name, metav1.GetOptions{})
			assert.True(t, apierrors.IsNotFound(err), "%s must be deleted", name)
		}
	})

	t.Run("succeeds when both CRDs are already gone", func(t *testing.T) {
		withClient(t, newCRDClient())

		assert.NoError(t, deleteMesheryCRDs())
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

		err := deleteMesheryCRDs()
		require.Error(t, err)
		assert.Equal(t, 2, attempts, "both CRDs must be attempted")
		assert.Contains(t, err.Error(), "brokers.meshery.io")
		assert.Contains(t, err.Error(), "meshsyncs.meshery.io")
	})

	t.Run("reports a client that cannot be built", func(t *testing.T) {
		apiExtensionsClient = func() (apiextension.Interface, error) {
			return nil, errors.New("no kubeconfig")
		}
		t.Cleanup(func() { apiExtensionsClient = previousClient })

		err := deleteMesheryCRDs()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no kubeconfig")
	})
}
