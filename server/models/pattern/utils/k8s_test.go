package utils

import (
	"context"
	"strings"
	"testing"

	"github.com/meshery/meshkit/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
)

var (
	configMapGVR  = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "configmaps"}
	namespaceGVR  = schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}
	testGVRToList = map[schema.GroupVersionResource]string{
		configMapGVR: "ConfigMapList",
		namespaceGVR: "NamespaceList",
	}
)

func testLogger(t *testing.T) logger.Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{Format: logger.JsonLogFormat})
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	return log
}

func seedNamespace(name string) *unstructured.Unstructured {
	ns := &unstructured.Unstructured{}
	ns.SetAPIVersion("v1")
	ns.SetKind("Namespace")
	ns.SetName(name)
	return ns
}

func seedConfigMap(name, ns string, labels map[string]string) *unstructured.Unstructured {
	cm := &unstructured.Unstructured{}
	cm.SetAPIVersion("v1")
	cm.SetKind("ConfigMap")
	cm.SetName(name)
	cm.SetNamespace(ns)
	if len(labels) > 0 {
		cm.SetLabels(labels)
	}
	return cm
}

// configMapInput builds the object passed to CreateK8sResource. Use nested map[string]interface{}
// for metadata (not unstructured.Object) so JSON round-trips match what the fake dynamic client expects.
func configMapInput(name, ns string, labels map[string]string) map[string]interface{} {
	meta := map[string]interface{}{
		"name":      name,
		"namespace": ns,
	}
	if len(labels) > 0 {
		meta["labels"] = labels
	}
	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   meta,
	}
}

func newFakeDynamicClient(objs ...runtime.Object) *fake.FakeDynamicClient {
	scheme := runtime.NewScheme()
	_ = corev1.AddToScheme(scheme)
	client := fake.NewSimpleDynamicClientWithCustomListKinds(scheme, testGVRToList, objs...)
	// CreateNamespace uses ApplyPatchType against namespaces; the fake ObjectTracker's Apply path
	// uses strategic merge against Unstructured and can error ("unable to find api field ...
	// metadata"). Short-circuit namespace server-side apply like the live API would succeed.
	client.PrependReactor("patch", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patch, ok := action.(k8stesting.PatchActionImpl)
		if !ok || patch.GetPatchType() != types.ApplyPatchType {
			return false, nil, nil
		}
		u := &unstructured.Unstructured{}
		u.SetAPIVersion("v1")
		u.SetKind("Namespace")
		u.SetName(patch.GetName())
		return true, u, nil
	})
	return client
}

func TestCreateK8sResource(t *testing.T) {
	t.Parallel()

	t.Run("force_true_recreates_unowned_resource_returns_nil", func(t *testing.T) {
		t.Parallel()
		ns := seedNamespace("default")
		existing := seedConfigMap("my-cm", "default", map[string]string{"controller": "someone-else"})
		client := newFakeDynamicClient(ns, existing)
		log := testLogger(t)

		err := CreateK8sResource(client, "", "v1", "configmaps",
			configMapInput("my-cm", "default", nil), true, log)
		if err != nil {
			t.Fatalf("expected nil after successful force recreate, got %v", err)
		}

		got, err := client.Resource(configMapGVR).Namespace("default").
			Get(context.Background(), "my-cm", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Get after recreate: %v", err)
		}
		labels := got.GetLabels()
		if labels["controller"] != "meshery" || labels["source"] != "pattern" {
			t.Fatalf("expected meshery labels after recreate, got %#v", labels)
		}
	})

	t.Run("force_false_unowned_resource_returns_error", func(t *testing.T) {
		t.Parallel()
		ns := seedNamespace("default")
		existing := seedConfigMap("my-cm", "default", map[string]string{"controller": "someone-else"})
		client := newFakeDynamicClient(ns, existing)
		log := testLogger(t)

		err := CreateK8sResource(client, "", "v1", "configmaps",
			configMapInput("my-cm", "default", nil), false, log)
		if err == nil {
			t.Fatal("expected error when force=false and resource not owned by meshery")
		}
		if !strings.Contains(err.Error(), "already exists") {
			t.Fatalf("expected 'already exists' in error, got: %v", err)
		}
	})

	t.Run("meshery_owned_resource_is_patched_returns_nil", func(t *testing.T) {
		t.Parallel()
		ns := seedNamespace("default")
		existing := seedConfigMap("owned-cm", "default", map[string]string{
			"controller": "meshery",
			"source":     "pattern",
		})
		client := newFakeDynamicClient(ns, existing)
		log := testLogger(t)

		err := CreateK8sResource(client, "", "v1", "configmaps",
			configMapInput("owned-cm", "default", nil), false, log)
		if err != nil {
			t.Fatalf("expected nil on patch path for meshery-owned resource, got %v", err)
		}

		_, err = client.Resource(configMapGVR).Namespace("default").
			Get(context.Background(), "owned-cm", metav1.GetOptions{})
		if err != nil {
			t.Fatalf("Get after patch: %v", err)
		}
	})

	t.Run("new_resource_creates_returns_nil", func(t *testing.T) {
		t.Parallel()
		ns := seedNamespace("default")
		client := newFakeDynamicClient(ns)
		log := testLogger(t)

		err := CreateK8sResource(client, "", "v1", "configmaps",
			configMapInput("new-cm", "default", nil), false, log)
		if err != nil {
			t.Fatalf("expected nil when creating new resource, got %v", err)
		}
	})
}
