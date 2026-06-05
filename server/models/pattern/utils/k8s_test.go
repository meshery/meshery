package utils

import (
	"fmt"
	"testing"

	"github.com/meshery/meshkit/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	dynamicfake "k8s.io/client-go/dynamic/fake"
	k8stesting "k8s.io/client-go/testing"
)

type fakePatternLogger struct {
	logger.Handler
}

func (f *fakePatternLogger) Info(_ ...interface{})                 {}
func (f *fakePatternLogger) Infof(_ string, _ ...interface{})      {}
func (f *fakePatternLogger) Debug(_ ...interface{})                {}
func (f *fakePatternLogger) Debugf(_ string, _ ...interface{})     {}
func (f *fakePatternLogger) Warn(_ error)                          {}
func (f *fakePatternLogger) Error(_ error)                         {}
func (f *fakePatternLogger) KV(_ string, _ interface{}) interface{} { return nil }

func TestCreateK8sResource_ClusterScoped_DoesNotAttemptNamespaceApply(t *testing.T) {
	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme())
	log := &fakePatternLogger{}

	client.PrependReactor("patch", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction, ok := action.(k8stesting.PatchAction)
		if !ok {
			return true, nil, fmt.Errorf("unexpected action type for namespace patch: %T", action)
		}

		if patchAction.GetName() == "" {
			return true, nil, fmt.Errorf("empty namespace patch should never be attempted")
		}

		return false, nil, nil
	})

	clusterRole := map[string]interface{}{
		"apiVersion": "rbac.authorization.k8s.io/v1",
		"kind":       "ClusterRole",
		"metadata": map[string]interface{}{
			"name": "test-cluster-role",
		},
		"rules": []interface{}{},
	}

	err := CreateK8sResource(
		client,
		"rbac.authorization.k8s.io",
		"v1",
		"clusterroles",
		clusterRole,
		false,
		log,
	)
	if err != nil {
		t.Fatalf("CreateK8sResource() returned error for cluster-scoped resource: %v", err)
	}
}

func TestCreateK8sResource_Namespaced_AttemptsNamespaceApply(t *testing.T) {
	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme())
	log := &fakePatternLogger{}

	namespacePatchCount := 0
	client.PrependReactor("patch", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
		patchAction, ok := action.(k8stesting.PatchAction)
		if !ok {
			return true, nil, fmt.Errorf("unexpected action type for namespace patch: %T", action)
		}

		if patchAction.GetName() == "my-namespace" {
			namespacePatchCount++
		}

		return true, nil, nil
	})

	configMap := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "test-config",
			"namespace": "my-namespace",
		},
		"data": map[string]interface{}{
			"key": "value",
		},
	}

	err := CreateK8sResource(
		client,
		"",
		"v1",
		"configmaps",
		configMap,
		false,
		log,
	)
	if err != nil {
		t.Fatalf("CreateK8sResource() returned error for namespaced resource: %v", err)
	}

	if namespacePatchCount == 0 {
		t.Fatalf("expected namespace apply attempt for namespaced resource, got none")
	}
}

func TestCreateK8sResource_ForceRecreate_ReturnsWithoutPatch(t *testing.T) {
	existing := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      "force-test-config",
				"namespace": "default",
				"labels": map[string]interface{}{
					"controller": "not-meshery",
				},
			},
			"data": map[string]interface{}{
				"old": "value",
			},
		},
	}

	client := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), existing)
	log := &fakePatternLogger{}

	client.PrependReactor("patch", "namespaces", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, nil
	})

	client.PrependReactor("patch", "configmaps", func(action k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, fmt.Errorf("patch should not be called after successful force recreate")
	})

	configMap := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata": map[string]interface{}{
			"name":      "force-test-config",
			"namespace": "default",
		},
		"data": map[string]interface{}{
			"new": "value",
		},
	}

	err := CreateK8sResource(
		client,
		"",
		"v1",
		"configmaps",
		configMap,
		true,
		log,
	)
	if err != nil {
		t.Fatalf("CreateK8sResource() returned error for force recreate flow: %v", err)
	}
}

