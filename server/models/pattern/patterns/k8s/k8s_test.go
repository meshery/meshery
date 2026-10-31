package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/meshery/schemas/models/v1beta2/component"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest/fake"
)

// newTestRESTMapper constructs a mock meta.RESTMapper configured with irregular plural Kubernetes resources for unit testing.
func newTestRESTMapper() meta.RESTMapper {
	groupVersions := []schema.GroupVersion{
		{Group: "", Version: "v1"},
		{Group: "apps", Version: "v1"},
		{Group: "networking.k8s.io", Version: "v1"},
		{Group: "storage.k8s.io", Version: "v1"},
	}
	mapper := meta.NewDefaultRESTMapper(groupVersions)

	// Add irregular plural kinds
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingresses"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "ingress"},
		meta.RESTScopeNamespace,
	)
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "NetworkPolicy"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicies"},
		schema.GroupVersionResource{Group: "networking.k8s.io", Version: "v1", Resource: "networkpolicy"},
		meta.RESTScopeNamespace,
	)
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "storage.k8s.io", Version: "v1", Kind: "StorageClass"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclasses"},
		schema.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "storageclass"},
		meta.RESTScopeRoot,
	)
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Endpoints"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "endpoints"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "endpoint"},
		meta.RESTScopeNamespace,
	)
	mapper.AddSpecific(
		schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"},
		schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespace"},
		meta.RESTScopeRoot,
	)

	return mapper
}

// newFakeRESTClient constructs a fake Kubernetes REST client using the provided HTTP round-trip handler function.
func newFakeRESTClient(fn func(req *http.Request) (*http.Response, error)) *fake.RESTClient {
	scheme := runtime.NewScheme()
	codecs := serializer.NewCodecFactory(scheme)
	return &fake.RESTClient{
		NegotiatedSerializer: codecs.WithoutConversion(),
		GroupVersion:         schema.GroupVersion{Group: "", Version: "v1"},
		Client:               fake.CreateHTTPClient(fn),
	}
}

// jsonResponse serializes an object into JSON and wraps it in an HTTP response with the specified status code.
func jsonResponse(statusCode int, obj interface{}) (*http.Response, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	resp := &http.Response{
		StatusCode: statusCode,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(data)),
	}
	return resp, nil
}

// TestResolveResourcePath_IrregularPlurals tests dynamic resource path resolution for irregular plural Kubernetes resources.
func TestResolveResourcePath_IrregularPlurals(t *testing.T) {
	mapper := newTestRESTMapper()

	tests := []struct {
		name             string
		k8sResource      map[string]interface{}
		namespace        string
		useMapper        bool
		wantBasePath     string
		wantResourcePath string
	}{
		{
			name: "Ingress with RESTMapper discovers ingresses",
			k8sResource: map[string]interface{}{
				"apiVersion": "networking.k8s.io/v1",
				"kind":       "Ingress",
			},
			namespace:        "test-ns",
			useMapper:        true,
			wantBasePath:     "/apis/networking.k8s.io/v1/namespaces/test-ns",
			wantResourcePath: "ingresses",
		},
		{
			name: "NetworkPolicy with RESTMapper discovers networkpolicies",
			k8sResource: map[string]interface{}{
				"apiVersion": "networking.k8s.io/v1",
				"kind":       "NetworkPolicy",
			},
			namespace:        "prod",
			useMapper:        true,
			wantBasePath:     "/apis/networking.k8s.io/v1/namespaces/prod",
			wantResourcePath: "networkpolicies",
		},
		{
			name: "StorageClass with RESTMapper discovers storageclasses and root scope",
			k8sResource: map[string]interface{}{
				"apiVersion": "storage.k8s.io/v1",
				"kind":       "StorageClass",
			},
			namespace:        "ignored-ns",
			useMapper:        true,
			wantBasePath:     "/apis/storage.k8s.io/v1",
			wantResourcePath: "storageclasses",
		},
		{
			name: "Endpoints with RESTMapper discovers endpoints",
			k8sResource: map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Endpoints",
			},
			namespace:        "default",
			useMapper:        true,
			wantBasePath:     "/api/v1/namespaces/default",
			wantResourcePath: "endpoints",
		},
		{
			name: "Fallback when mapper is nil uses strings.ToLower(kind) + s",
			k8sResource: map[string]interface{}{
				"apiVersion": "networking.k8s.io/v1",
				"kind":       "Ingress",
			},
			namespace:        "test-ns",
			useMapper:        false,
			wantBasePath:     "/apis/networking.k8s.io/v1/namespaces/test-ns",
			wantResourcePath: "ingresss", // naive fallback
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var m meta.RESTMapper
			if tc.useMapper {
				m = mapper
			}
			basePath, resourcePath, err := resolveResourcePath(tc.k8sResource, tc.namespace, m)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if basePath != tc.wantBasePath {
				t.Errorf("basePath = %q, want %q", basePath, tc.wantBasePath)
			}
			if resourcePath != tc.wantResourcePath {
				t.Errorf("resourcePath = %q, want %q", resourcePath, tc.wantResourcePath)
			}
		})
	}
}

// TestDryRun_UndeployNotFoundSuppression verifies that 404/NotFound errors during delete/undeploy dry-runs are suppressed.
func TestDryRun_UndeployNotFoundSuppression(t *testing.T) {
	t.Run("suppresses 404 NotFound on delete dry-run", func(t *testing.T) {
		rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Fatalf("expected DELETE request, got %s", req.Method)
			}
			if req.URL.Query().Get("dryRun") != "All" {
				t.Errorf("expected dryRun=All query param, got %s", req.URL.Query().Get("dryRun"))
			}
			status := metav1.Status{
				TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
				Status:   metav1.StatusFailure,
				Message:  "deployments.apps \"my-deploy\" not found",
				Reason:   metav1.StatusReasonNotFound,
				Code:     http.StatusNotFound,
				Details: &metav1.StatusDetails{
					Name: "my-deploy",
					Kind: "deployments",
				},
			}
			return jsonResponse(http.StatusNotFound, status)
		})

		resource := map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "my-deploy",
			},
		}

		st, success, err := dryRun(rClient, resource, "default", true, DryRunOptions{})
		if err != nil {
			t.Fatalf("expected nil error on undeploy dry-run 404, got: %v", err)
		}
		if !success {
			t.Fatalf("expected success=true on undeploy dry-run 404, got false")
		}
		if st == nil || st["status"] != "Success" {
			t.Errorf("expected status=Success, got: %v", st)
		}
	})

	t.Run("does not suppress non-404 error on delete dry-run", func(t *testing.T) {
		rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
			status := metav1.Status{
				TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
				Status:   metav1.StatusFailure,
				Message:  "forbidden",
				Reason:   metav1.StatusReasonForbidden,
				Code:     http.StatusForbidden,
			}
			return jsonResponse(http.StatusForbidden, status)
		})

		resource := map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "my-deploy",
			},
		}

		st, success, _ := dryRun(rClient, resource, "default", true, DryRunOptions{})
		if success {
			t.Fatalf("expected success=false on delete forbidden, got true")
		}
		if st == nil || st["status"] != "Failure" {
			t.Errorf("expected status=Failure, got: %v", st)
		}
	})
}

// TestDryRun_InDesignNamespaces tests dry-run validation when resources reference namespaces declared within the same design pattern.
func TestDryRun_InDesignNamespaces(t *testing.T) {
	t.Run("substitutes default namespace when target namespace is in declaredNamespaces and missing", func(t *testing.T) {
		mapper := newTestRESTMapper()
		var requestedPaths []string

		rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
			requestedPaths = append(requestedPaths, req.URL.Path)

			if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/my-custom-ns/ingresses" {
				// Initial call targeting undeployed in-design namespace fails with 404
				status := metav1.Status{
					TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
					Status:   metav1.StatusFailure,
					Message:  "namespaces \"my-custom-ns\" not found",
					Reason:   metav1.StatusReasonNotFound,
					Code:     http.StatusNotFound,
					Details: &metav1.StatusDetails{
						Name: "my-custom-ns",
						Kind: "namespaces",
					},
				}
				return jsonResponse(http.StatusNotFound, status)
			}

			if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/default/ingresses" {
				// Re-run substituting default succeeds
				status := metav1.Status{
					TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
					Status:   metav1.StatusSuccess,
				}
				return jsonResponse(http.StatusOK, status)
			}

			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		})

		resource := map[string]interface{}{
			"apiVersion": "networking.k8s.io/v1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name":      "my-ingress",
				"namespace": "my-custom-ns",
			},
		}

		opts := DryRunOptions{
			RESTMapper: mapper,
			DeclaredNamespaces: map[string]bool{
				"my-custom-ns": true,
			},
		}

		st, success, err := dryRun(rClient, resource, "my-custom-ns", false, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !success {
			t.Fatalf("expected dry-run to succeed via default namespace fallback")
		}
		if st == nil || st["status"] != "Success" {
			t.Errorf("expected status=Success, got: %v", st)
		}

		// Verify that both paths were requested: first custom-ns, then default
		if len(requestedPaths) != 2 {
			t.Fatalf("expected 2 requests, got %d: %v", len(requestedPaths), requestedPaths)
		}
		if requestedPaths[0] != "/apis/networking.k8s.io/v1/namespaces/my-custom-ns/ingresses" {
			t.Errorf("request 0 = %s, want custom-ns path", requestedPaths[0])
		}
		if requestedPaths[1] != "/apis/networking.k8s.io/v1/namespaces/default/ingresses" {
			t.Errorf("request 1 = %s, want default path", requestedPaths[1])
		}

		// Verify original resource object was not mutated
		metaMap := resource["metadata"].(map[string]interface{})
		if metaMap["namespace"] != "my-custom-ns" {
			t.Errorf("expected original resource namespace to remain my-custom-ns, got %v", metaMap["namespace"])
		}
	})

	t.Run("surfaces workload field errors when validation against default fails", func(t *testing.T) {
		mapper := newTestRESTMapper()

		rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
			if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/my-custom-ns/ingresses" {
				status := metav1.Status{
					TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
					Status:   metav1.StatusFailure,
					Message:  "namespaces \"my-custom-ns\" not found",
					Reason:   metav1.StatusReasonNotFound,
					Code:     http.StatusNotFound,
					Details: &metav1.StatusDetails{
						Name: "my-custom-ns",
						Kind: "namespaces",
					},
				}
				return jsonResponse(http.StatusNotFound, status)
			}

			if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/default/ingresses" {
				// Workload itself has invalid fields
				status := metav1.Status{
					TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
					Status:   metav1.StatusFailure,
					Message:  "Ingress.networking.k8s.io \"my-ingress\" is invalid: spec.rules: Required value",
					Reason:   metav1.StatusReasonInvalid,
					Code:     http.StatusUnprocessableEntity,
					Details: &metav1.StatusDetails{
						Causes: []metav1.StatusCause{
							{
								Type:    metav1.CauseTypeFieldValueRequired,
								Message: "Required value",
								Field:   "spec.rules",
							},
						},
					},
				}
				return jsonResponse(http.StatusUnprocessableEntity, status)
			}

			return nil, fmt.Errorf("unexpected path: %s", req.URL.Path)
		})

		resource := map[string]interface{}{
			"apiVersion": "networking.k8s.io/v1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name":      "my-ingress",
				"namespace": "my-custom-ns",
			},
		}

		opts := DryRunOptions{
			RESTMapper: mapper,
			DeclaredNamespaces: map[string]bool{
				"my-custom-ns": true,
			},
		}

		st, success, _ := dryRun(rClient, resource, "my-custom-ns", false, opts)
		if success {
			t.Fatalf("expected validation against default to fail with workload field errors, got success")
		}
		if st == nil || st["status"] != "Failure" {
			t.Errorf("expected status=Failure from workload validation, got: %v", st)
		}
	})

	t.Run("does not substitute when target namespace was not declared in design", func(t *testing.T) {
		mapper := newTestRESTMapper()
		callCount := 0

		rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
			callCount++
			status := metav1.Status{
				TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
				Status:   metav1.StatusFailure,
				Message:  "namespaces \"undeclared-ns\" not found",
				Reason:   metav1.StatusReasonNotFound,
				Code:     http.StatusNotFound,
				Details: &metav1.StatusDetails{
					Name: "undeclared-ns",
					Kind: "namespaces",
				},
			}
			return jsonResponse(http.StatusNotFound, status)
		})

		resource := map[string]interface{}{
			"apiVersion": "networking.k8s.io/v1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name":      "my-ingress",
				"namespace": "undeclared-ns",
			},
		}

		opts := DryRunOptions{
			RESTMapper: mapper,
			DeclaredNamespaces: map[string]bool{
				"other-ns": true,
			},
		}

		_, success, _ := dryRun(rClient, resource, "undeclared-ns", false, opts)
		if success {
			t.Fatalf("expected failure when namespace was not declared, got success")
		}
		if callCount != 1 {
			t.Errorf("expected exactly 1 call (no fallback), got %d", callCount)
		}
	})
}

// TestDryRunHelper_NilClient verifies that passing a nil Kubernetes client to DryRunHelper returns an error.
func TestDryRunHelper_NilClient(t *testing.T) {
	comp := component.ComponentDefinition{
		Component: component.Component{
			Kind: "Pod",
		},
	}
	_, success, err := DryRunHelper(nil, comp, false, DryRunOptions{})
	if success || err == nil {
		t.Fatalf("expected failure with nil client, got success=%v, err=%v", success, err)
	}
}

// TestIsMissingNamespaceError_StatusError tests detection of missing namespace errors from Kubernetes StatusError objects.
func TestIsMissingNamespaceError_StatusError(t *testing.T) {
	statusErr := &kubeerror.StatusError{
		ErrStatus: metav1.Status{
			Status: metav1.StatusFailure,
			Reason: metav1.StatusReasonNotFound,
			Code:   http.StatusNotFound,
			Details: &metav1.StatusDetails{
				Name: "test-ns",
				Kind: "namespaces",
			},
		},
	}

	if !isMissingNamespaceError(statusErr, nil, "test-ns") {
		t.Errorf("expected isMissingNamespaceError to be true for statusErr matching target namespace")
	}

	if isMissingNamespaceError(statusErr, nil, "other-ns") {
		t.Errorf("expected isMissingNamespaceError to be false for statusErr with different namespace name")
	}
}

// TestResolveResourcePath_MissingFields tests that resolveResourcePath validates that kind and apiVersion are non-empty.
func TestResolveResourcePath_MissingFields(t *testing.T) {
	_, _, err := resolveResourcePath(map[string]interface{}{"kind": "Pod"}, "default", nil)
	if err == nil {
		t.Errorf("expected error when apiVersion is empty")
	}

	_, _, err = resolveResourcePath(map[string]interface{}{"apiVersion": "v1"}, "default", nil)
	if err == nil {
		t.Errorf("expected error when kind is empty")
	}
}

// TestDryRun_InDesignNamespaces_DefaultAlsoNotFound tests fallback behavior when the fallback default namespace also returns NotFound.
func TestDryRun_InDesignNamespaces_DefaultAlsoNotFound(t *testing.T) {
	mapper := newTestRESTMapper()

	// Both custom-ns and fallback default-ns return 404 missing namespace
	rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
		ns := "my-custom-ns"
		if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/default/ingresses" {
			ns = "default"
		}
		return jsonResponse(http.StatusNotFound, namespaceNotFoundStatus(ns))
	})

	resource := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1",
		"kind":       "Ingress",
		"metadata": map[string]interface{}{
			"name":      "my-ingress",
			"namespace": "my-custom-ns",
		},
	}

	opts := DryRunOptions{
		RESTMapper: mapper,
		DeclaredNamespaces: map[string]bool{
			"my-custom-ns": true,
		},
	}

	st, success, err := dryRun(rClient, resource, "my-custom-ns", false, opts)
	if err != nil || !success {
		t.Fatalf("expected a non-fatal pass when default is missing too, got success=%v err=%v st=%v", success, err, st)
	}
	// The component must keep its original configuration, not a synthetic {"status": "Success"}.
	if st["kind"] != "Ingress" || st["status"] != nil {
		t.Errorf("expected the workload itself to be returned, got: %v", st)
	}
	metadata, _ := st["metadata"].(map[string]interface{})
	if metadata["namespace"] != "my-custom-ns" {
		t.Errorf("expected the original namespace to be kept, got: %v", metadata["namespace"])
	}
}

// TestGetNamespaceForComponent tests extraction and defaulting of namespaces from component definition configuration and metadata.
func TestGetNamespaceForComponent(t *testing.T) {
	t.Run("extracts namespace from metadata.namespace", func(t *testing.T) {
		comp := &component.ComponentDefinition{
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{
					"namespace": "meta-ns",
				},
			},
		}
		ns := getNamespaceForComponent(comp)
		if ns != "meta-ns" {
			t.Errorf("got %q, want %q", ns, "meta-ns")
		}
	})

	t.Run("ignores a top-level Configuration namespace", func(t *testing.T) {
		// The converter drops a top-level namespace key from the manifest, so it must not steer apply or dry-run.
		comp := &component.ComponentDefinition{
			Configuration: map[string]interface{}{
				"namespace": "top-ns",
			},
		}
		comp.Metadata.AdditionalProperties = map[string]interface{}{
			"isNamespaced": true,
		}
		ns := getNamespaceForComponent(comp)
		if ns != "default" {
			t.Errorf("got %q, want %q", ns, "default")
		}
	})

	t.Run("defaults to default when isNamespaced is true and no namespace specified", func(t *testing.T) {
		comp := &component.ComponentDefinition{}
		comp.Metadata.AdditionalProperties = map[string]interface{}{
			"isNamespaced": true,
		}
		ns := getNamespaceForComponent(comp)
		if ns != "default" {
			t.Errorf("got %q, want %q", ns, "default")
		}
	})
}

// TestIsMissingNamespaceError_SingularKind tests detection of missing namespace errors when error details specify singular kind "namespace".
func TestIsMissingNamespaceError_SingularKind(t *testing.T) {
	st := map[string]interface{}{
		"status": "Failure",
		"reason": "NotFound",
		"details": map[string]interface{}{
			"kind": "namespace", // singular
			"name": "my-ns",
		},
	}

	if !isMissingNamespaceError(nil, st, "my-ns") {
		t.Errorf("expected isMissingNamespaceError to be true for singular kind 'namespace'")
	}
}

// namespaceNotFoundStatus builds the Status the API server returns when creating into a missing namespace.
func namespaceNotFoundStatus(ns string) metav1.Status {
	return metav1.Status{
		TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
		Status:   metav1.StatusFailure,
		Message:  fmt.Sprintf("namespaces %q not found", ns),
		Reason:   metav1.StatusReasonNotFound,
		Code:     http.StatusNotFound,
		Details:  &metav1.StatusDetails{Name: ns, Kind: "namespaces"},
	}
}

// TestIsMissingNamespaceError_RejectsUnverifiedFailures pins that only a NotFound Status naming the target
// namespace counts as a missing namespace; anything else must keep its original error.
func TestIsMissingNamespaceError_RejectsUnverifiedFailures(t *testing.T) {
	cases := []struct {
		name string
		err  error
		st   map[string]interface{}
	}{
		{
			name: "NotFound for another resource whose message mentions the namespace",
			err:  kubeerror.NewNotFound(schema.GroupResource{Resource: "configmaps"}, "app-config"),
			st: map[string]interface{}{
				"status":  "Failure",
				"reason":  "NotFound",
				"message": "configmaps \"app-config\" not found in namespace \"my-ns\"",
				"details": map[string]interface{}{"kind": "configmaps", "name": "app-config"},
			},
		},
		{
			name: "non-NotFound failure whose message mentions a missing namespace",
			err:  kubeerror.NewForbidden(schema.GroupResource{Resource: "namespaces"}, "my-ns", fmt.Errorf("namespace not found in allow list")),
			st: map[string]interface{}{
				"status":  "Failure",
				"reason":  "Forbidden",
				"message": "namespace my-ns not found in allow list",
				"details": map[string]interface{}{"kind": "namespaces", "name": "my-ns"},
			},
		},
		{
			name: "message-only status without reason or details",
			st: map[string]interface{}{
				"status":  "Failure",
				"message": "namespaces \"my-ns\" not found",
			},
		},
		{
			name: "NotFound for namespaces without a name",
			st: map[string]interface{}{
				"status":  "Failure",
				"reason":  "NotFound",
				"details": map[string]interface{}{"kind": "namespaces"},
			},
		},
		{
			name: "NotFound for a different namespace",
			err:  kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, "other-ns"),
		},
		{
			name: "plain error mentioning a namespace",
			err:  fmt.Errorf("namespace my-ns not found"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if isMissingNamespaceError(tc.err, tc.st, "my-ns") {
				t.Errorf("expected isMissingNamespaceError to be false")
			}
		})
	}

	t.Run("accepts a wrapped namespace NotFound", func(t *testing.T) {
		err := fmt.Errorf("dry run: %w", kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, "my-ns"))
		if !isMissingNamespaceError(err, nil, "my-ns") {
			t.Errorf("expected isMissingNamespaceError to be true for wrapped namespace NotFound")
		}
	})

	t.Run("rejects an empty target namespace", func(t *testing.T) {
		err := kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, "")
		if isMissingNamespaceError(err, nil, "") {
			t.Errorf("expected isMissingNamespaceError to be false for empty target")
		}
	})
}

// TestDryRun_InDesignNamespaces_KeepsUnrelatedNotFound pins that a NotFound for something other than the
// declared namespace is surfaced as-is instead of being retried against default.
func TestDryRun_InDesignNamespaces_KeepsUnrelatedNotFound(t *testing.T) {
	callCount := 0
	rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
		callCount++
		status := metav1.Status{
			TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
			Status:   metav1.StatusFailure,
			Message:  "serviceaccounts \"builder\" not found in namespace \"my-custom-ns\"",
			Reason:   metav1.StatusReasonNotFound,
			Code:     http.StatusNotFound,
			Details:  &metav1.StatusDetails{Name: "builder", Kind: "serviceaccounts"},
		}
		return jsonResponse(http.StatusNotFound, status)
	})

	resource := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1",
		"kind":       "Ingress",
		"metadata":   map[string]interface{}{"name": "my-ingress", "namespace": "my-custom-ns"},
	}
	opts := DryRunOptions{
		RESTMapper:         newTestRESTMapper(),
		DeclaredNamespaces: map[string]bool{"my-custom-ns": true},
	}

	st, success, _ := dryRun(rClient, resource, "my-custom-ns", false, opts)
	if success {
		t.Fatalf("expected the unrelated NotFound to fail the dry-run, got success")
	}
	if callCount != 1 {
		t.Errorf("expected exactly 1 call (no default fallback), got %d", callCount)
	}
	if st == nil || st["message"] != "serviceaccounts \"builder\" not found in namespace \"my-custom-ns\"" {
		t.Errorf("expected the original Status to be returned, got: %v", st)
	}
}

// TestIsDeleteTargetNotFound pins that a delete dry-run 404 counts as "not deployed" only when the
// NotFound Status names the object being deleted, or the namespace holding it.
func TestIsDeleteTargetNotFound(t *testing.T) {
	deployment := map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]interface{}{"name": "web", "namespace": "team"},
	}

	t.Run("object being deleted is not found", func(t *testing.T) {
		err := kubeerror.NewNotFound(schema.GroupResource{Group: "apps", Resource: "deployments"}, "web")
		if !isDeleteTargetNotFound(err, nil, deployment, "team", newTestRESTMapper()) {
			t.Errorf("expected true for a NotFound naming the deleted object")
		}
	})

	t.Run("namespace holding the object is not found", func(t *testing.T) {
		err := kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, "team")
		if !isDeleteTargetNotFound(err, nil, deployment, "team", newTestRESTMapper()) {
			t.Errorf("expected true for a NotFound naming the object's namespace")
		}
	})

	notMatching := []struct {
		name string
		err  error
		st   map[string]interface{}
	}{
		{
			name: "missing CRD or unknown resource path",
			st: map[string]interface{}{
				"status":  "Failure",
				"reason":  "NotFound",
				"code":    float64(http.StatusNotFound),
				"message": "the server could not find the requested resource",
				"details": map[string]interface{}{},
			},
		},
		{
			name: "bare 404 without reason or details",
			st:   map[string]interface{}{"status": "Failure", "code": float64(http.StatusNotFound)},
		},
		{
			name: "NotFound for another object of the same resource",
			err:  kubeerror.NewNotFound(schema.GroupResource{Group: "apps", Resource: "deployments"}, "api"),
		},
		{
			name: "NotFound for a different resource",
			err:  kubeerror.NewNotFound(schema.GroupResource{Resource: "configmaps"}, "web"),
		},
		{
			name: "NotFound for a different namespace",
			err:  kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, "other"),
		},
		{
			name: "non-NotFound failure",
			err:  kubeerror.NewForbidden(schema.GroupResource{Group: "apps", Resource: "deployments"}, "web", nil),
		},
	}
	for _, tc := range notMatching {
		t.Run(tc.name, func(t *testing.T) {
			if isDeleteTargetNotFound(tc.err, tc.st, deployment, "team", newTestRESTMapper()) {
				t.Errorf("expected false")
			}
		})
	}
}

// TestDryRun_UndeployMissingCRD pins that deleting a custom resource whose CRD is not installed is reported
// as a failure, not as an already-undeployed object.
func TestDryRun_UndeployMissingCRD(t *testing.T) {
	var requested string
	rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
		requested = req.URL.Path
		status := metav1.Status{
			TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
			Status:   metav1.StatusFailure,
			Message:  "the server could not find the requested resource",
			Reason:   metav1.StatusReasonNotFound,
			Code:     http.StatusNotFound,
			Details:  &metav1.StatusDetails{},
		}
		return jsonResponse(http.StatusNotFound, status)
	})

	widget := map[string]interface{}{
		"apiVersion": "example.com/v1",
		"kind":       "Widget",
		"metadata":   map[string]interface{}{"name": "my-widget", "namespace": "default"},
	}

	st, success, _ := dryRun(rClient, widget, "default", true, DryRunOptions{RESTMapper: newTestRESTMapper()})
	if success {
		t.Fatalf("expected a missing CRD to fail the undeploy dry-run, got success")
	}
	if st == nil || st["status"] != "Failure" {
		t.Errorf("expected the original Failure Status, got: %v", st)
	}
	if requested != "/apis/example.com/v1/namespaces/default/widgets/my-widget" {
		t.Errorf("unexpected request path %q", requested)
	}
}

// TestResolveResourcePath_UnservedVersion pins that a kind the cluster serves only at another version fails with
// an explicit error naming the preferred version, and that no request is sent to another version's endpoint.
func TestResolveResourcePath_UnservedVersion(t *testing.T) {
	ingress := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1beta1",
		"kind":       "Ingress",
		"metadata":   map[string]interface{}{"name": "web", "namespace": "team"},
	}

	_, _, err := resolveResourcePath(ingress, "team", newTestRESTMapper())
	if err == nil {
		t.Fatalf("expected an error for an unserved apiVersion")
	}
	want := "apiVersion networking.k8s.io/v1beta1 is not served by this cluster (preferred: networking.k8s.io/v1)"
	if !strings.Contains(err.Error(), want) {
		t.Errorf("error = %q, want it to contain %q", err.Error(), want)
	}

	requests := 0
	rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
		requests++
		return jsonResponse(http.StatusOK, map[string]interface{}{"kind": "Ingress"})
	})
	_, success, _ := dryRun(rClient, ingress, "team", false, DryRunOptions{RESTMapper: newTestRESTMapper()})
	if success {
		t.Errorf("expected the dry-run to fail for an unserved apiVersion")
	}
	if requests != 0 {
		t.Errorf("expected no request to the API server, got %d", requests)
	}
}

// TestDryRun_InDesignNamespaces_FallbackNameConflict pins that an unrelated object of the same name in default
// does not fail a workload whose declared namespace is missing, and that the workload, not a Status, is returned.
func TestDryRun_InDesignNamespaces_FallbackNameConflict(t *testing.T) {
	rClient := newFakeRESTClient(func(req *http.Request) (*http.Response, error) {
		if req.URL.Path == "/apis/networking.k8s.io/v1/namespaces/default/ingresses" {
			status := kubeerror.NewAlreadyExists(schema.GroupResource{Group: "networking.k8s.io", Resource: "ingresses"}, "web").ErrStatus
			status.TypeMeta = metav1.TypeMeta{Kind: "Status", APIVersion: "v1"}
			return jsonResponse(http.StatusConflict, status)
		}
		return jsonResponse(http.StatusNotFound, namespaceNotFoundStatus("staging"))
	})

	ingress := map[string]interface{}{
		"apiVersion": "networking.k8s.io/v1",
		"kind":       "Ingress",
		"metadata":   map[string]interface{}{"name": "web", "namespace": "staging"},
		"spec":       map[string]interface{}{"ingressClassName": "nginx"},
	}
	opts := DryRunOptions{
		RESTMapper:         newTestRESTMapper(),
		DeclaredNamespaces: map[string]bool{"staging": true},
	}

	st, success, err := dryRun(rClient, ingress, "staging", false, opts)
	if err != nil || !success {
		t.Fatalf("expected success on a fallback name conflict, got success=%v err=%v st=%v", success, err, st)
	}
	if st["kind"] != "Ingress" || st["spec"] == nil {
		t.Errorf("expected the workload to be returned, got %v", st)
	}

	t.Run("AlreadyExists for another name is not a pass", func(t *testing.T) {
		if isFallbackNameConflict(nil, map[string]interface{}{
			"reason":  "AlreadyExists",
			"details": map[string]interface{}{"name": "other"},
		}, ingress) {
			t.Errorf("expected false for a conflict on another object")
		}
	})
}
