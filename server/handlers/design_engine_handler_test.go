package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshkit/logger"
	meshkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	kubeerror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest/fake"
)

func mutateTestComponent(displayName, kind string) *component.ComponentDefinition {
	return &component.ComponentDefinition{
		ID:          uuid.Must(uuid.NewV4()),
		DisplayName: displayName,
		Component:   component.Component{Kind: kind},
	}
}

func newMutateActionProvider(t *testing.T) *serviceActionProvider {
	t.Helper()

	log, err := logger.New("mutate-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	return &serviceActionProvider{log: log}
}

// The dependency Mutate writes has to be one the execution plan can resolve: a
// dependency names another component of the same design by its name, and no
// component depends on itself.
func TestMutateOrdersComponentsBehindACustomResourceDefinition(t *testing.T) {
	crd := mutateTestComponent("widgets.example.com", "CustomResourceDefinition")
	widget := mutateTestComponent("a widget", "Widget")

	design := pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       "mutate-test",
		Components: []*component.ComponentDefinition{crd, widget},
	}

	newMutateActionProvider(t).Mutate(&design)

	dependsOn, err := planner.DeclaredDependencies(design.Name, widget)
	if err != nil {
		t.Fatalf("the dependency Mutate wrote cannot be read back: %v", err)
	}

	if len(dependsOn) != 1 || dependsOn[0] != crd.DisplayName {
		t.Errorf("%q depends on %v, want [%q]", widget.DisplayName, dependsOn, crd.DisplayName)
	}

	for _, dep := range dependsOn {
		if dep == widget.DisplayName || dep == widget.ID.String() {
			t.Errorf("%q was made to depend on itself through %q", widget.DisplayName, dep)
		}
	}

	if crdDependsOn, err := planner.DeclaredDependencies(design.Name, crd); err != nil || len(crdDependsOn) != 0 {
		t.Errorf("the custom resource definition depends on %v (err %v), want nothing", crdDependsOn, err)
	}

	// The plan is what consumes the dependency, so building one is what proves
	// Mutate wrote a value that is actually usable.
	plan, err := planner.CreatePlan(design, nil, false)
	if err != nil {
		t.Fatalf("no plan could be built from the design Mutate produced: %v", err)
	}

	log, err := logger.New("mutate-test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build a logger: %v", err)
	}

	order := []string{}
	if err := plan.Execute(func(id string, comp component.ComponentDefinition) bool {
		order = append(order, comp.DisplayName)
		return true
	}, func(string, component.ComponentDefinition, string) {}, log); err != nil {
		t.Fatalf("the plan could not be executed: %v", err)
	}

	if len(order) != 2 || order[0] != crd.DisplayName {
		t.Errorf("deployed in the order %v, want the custom resource definition first", order)
	}
}

// A design that already declares the dependency Mutate would add is left alone,
// rather than accumulating the same name twice.
func TestMutateDoesNotRepeatADependencyTheDesignAlreadyDeclares(t *testing.T) {
	crd := mutateTestComponent("widgets.example.com", "CustomResourceDefinition")
	widget := mutateTestComponent("a widget", "Widget")
	widget.Metadata.AdditionalProperties = map[string]interface{}{
		"dependsOn": []interface{}{crd.DisplayName},
	}

	design := pattern.PatternFile{
		ID:         uuid.Must(uuid.NewV4()),
		Name:       "mutate-test",
		Components: []*component.ComponentDefinition{crd, widget},
	}

	newMutateActionProvider(t).Mutate(&design)

	dependsOn, err := planner.DeclaredDependencies(design.Name, widget)
	if err != nil {
		t.Fatalf("the dependency Mutate wrote cannot be read back: %v", err)
	}

	if len(dependsOn) != 1 || dependsOn[0] != crd.DisplayName {
		t.Errorf("%q depends on %v, want [%q]", widget.DisplayName, dependsOn, crd.DisplayName)
	}
}

// TestGetNamespaceNameFromComponent tests extracting namespace names from component metadata, configuration, or display name fallback.
func TestGetNamespaceNameFromComponent(t *testing.T) {
	t.Run("returns DisplayName even when metadata.name differs", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "staging",
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": "production",
				},
			},
		}
		if name := getNamespaceNameFromComponent(cmp); name != "staging" {
			t.Errorf("got %q, want %q", name, "staging")
		}
	})

	t.Run("returns DisplayName without configuration", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "fallback-display-name",
		}
		if name := getNamespaceNameFromComponent(cmp); name != "fallback-display-name" {
			t.Errorf("got %q, want %q", name, "fallback-display-name")
		}
	})

	t.Run("returns empty string for nil component", func(t *testing.T) {
		if name := getNamespaceNameFromComponent(nil); name != "" {
			t.Errorf("got %q, want empty string", name)
		}
	})
}

// TestCollectDeclaredNamespaces pins that a Namespace is registered under the name deploy creates it with,
// its DisplayName, and never under a differing metadata.name.
func TestCollectDeclaredNamespaces(t *testing.T) {
	ns := mutateTestComponent("staging", "Namespace")
	ns.Configuration = map[string]interface{}{
		"metadata": map[string]interface{}{"name": "production"},
	}
	workload := mutateTestComponent("web", "Deployment")

	got := collectDeclaredNamespaces([]*component.ComponentDefinition{ns, workload, nil})

	if len(got) != 1 || !got["staging"] {
		t.Fatalf("got %v, want only %q declared", got, "staging")
	}
	if got["production"] {
		t.Errorf("metadata.name %q must not be registered; deploy creates the Namespace as %q", "production", "staging")
	}
}

// fakeCluster is a stand-in API server for DryRun tests. It serves discovery for core/v1 and
// networking.k8s.io/v1 unless discovery is broken, records the path of every dry-run request, and
// answers dry-runs with the response chosen by respond.
type fakeCluster struct {
	server *httptest.Server
	mu     sync.Mutex
	paths  []string
}

func newFakeCluster(t *testing.T, brokenDiscovery bool, respond func(path string, body []byte) (int, interface{})) *fakeCluster {
	t.Helper()

	discovery := map[string]interface{}{
		"/api": metav1.APIVersions{
			TypeMeta: metav1.TypeMeta{Kind: "APIVersions"},
			Versions: []string{"v1"},
		},
		"/apis": metav1.APIGroupList{
			TypeMeta: metav1.TypeMeta{Kind: "APIGroupList", APIVersion: "v1"},
			Groups: []metav1.APIGroup{{
				Name:             "networking.k8s.io",
				Versions:         []metav1.GroupVersionForDiscovery{{GroupVersion: "networking.k8s.io/v1", Version: "v1"}},
				PreferredVersion: metav1.GroupVersionForDiscovery{GroupVersion: "networking.k8s.io/v1", Version: "v1"},
			}},
		},
		"/api/v1": metav1.APIResourceList{
			TypeMeta:     metav1.TypeMeta{Kind: "APIResourceList", APIVersion: "v1"},
			GroupVersion: "v1",
			APIResources: []metav1.APIResource{{Name: "namespaces", Kind: "Namespace", Verbs: []string{"create", "delete", "get"}}},
		},
		"/apis/networking.k8s.io/v1": metav1.APIResourceList{
			TypeMeta:     metav1.TypeMeta{Kind: "APIResourceList", APIVersion: "v1"},
			GroupVersion: "networking.k8s.io/v1",
			APIResources: []metav1.APIResource{{Name: "ingresses", Namespaced: true, Kind: "Ingress", Verbs: []string{"create", "delete", "get"}}},
		},
	}

	fc := &fakeCluster{}
	fc.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("dryRun") != "All" {
			obj, ok := discovery[r.URL.Path]
			if brokenDiscovery || !ok {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(obj)
			return
		}

		body, _ := io.ReadAll(r.Body)
		fc.mu.Lock()
		fc.paths = append(fc.paths, r.URL.Path)
		fc.mu.Unlock()

		code, obj := respond(r.URL.Path, body)
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(obj)
	}))
	t.Cleanup(fc.server.Close)
	return fc
}

func (fc *fakeCluster) kubeconfig() string {
	return fmt.Sprintf(`apiVersion: v1
kind: Config
clusters:
- name: fake
  cluster:
    server: %s
contexts:
- name: fake
  context:
    cluster: fake
    user: fake
current-context: fake
users:
- name: fake
  user:
    token: fake
`, fc.server.URL)
}

func (fc *fakeCluster) requestedPaths() []string {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return append([]string(nil), fc.paths...)
}

// echoDryRun answers a dry-run create the way the API server does on success: with the object itself.
func echoDryRun(_ string, body []byte) (int, interface{}) {
	return http.StatusCreated, json.RawMessage(body)
}

// namespaceNotFoundStatus is the Status the API server returns when creating into a missing namespace.
func namespaceNotFoundStatus(ns string) metav1.Status {
	status := kubeerror.NewNotFound(schema.GroupResource{Resource: "namespaces"}, ns).ErrStatus
	status.TypeMeta = metav1.TypeMeta{Kind: "Status", APIVersion: "v1"}
	return status
}

func dryRunTestComponent(displayName, kind, version string, configuration map[string]interface{}) *component.ComponentDefinition {
	cmp := mutateTestComponent(displayName, kind)
	cmp.Component.Version = version
	cmp.Configuration = configuration
	return cmp
}

// TestDryRunCreateFallsBackOnlyForDeclaredNamespaces runs DryRun end to end on the create path: a workload
// in a Namespace the design declares is validated against default, while a workload in a namespace the
// design does not create fails with the original NotFound.
func TestDryRunCreateFallsBackOnlyForDeclaredNamespaces(t *testing.T) {
	const (
		stagingPath    = "/apis/networking.k8s.io/v1/namespaces/staging/ingresses"
		productionPath = "/apis/networking.k8s.io/v1/namespaces/production/ingresses"
		defaultPath    = "/apis/networking.k8s.io/v1/namespaces/default/ingresses"
	)
	cluster := newFakeCluster(t, false, func(path string, body []byte) (int, interface{}) {
		switch path {
		case stagingPath:
			return http.StatusNotFound, namespaceNotFoundStatus("staging")
		case productionPath:
			return http.StatusNotFound, namespaceNotFoundStatus("production")
		}
		return echoDryRun(path, body)
	})

	// Deploy creates this Namespace as "staging", its DisplayName, not as its metadata.name.
	namespace := dryRunTestComponent("staging", "Namespace", "v1", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "production"},
	})
	inDeclared := dryRunTestComponent("web", "Ingress", "networking.k8s.io/v1", map[string]interface{}{
		"metadata": map[string]interface{}{"namespace": "staging"},
	})
	inUndeclared := dryRunTestComponent("api", "Ingress", "networking.k8s.io/v1", map[string]interface{}{
		"metadata": map[string]interface{}{"namespace": "production"},
	})

	sap := newMutateActionProvider(t)
	sap.ctxTokubeconfig = map[string]string{"ctx-a": cluster.kubeconfig()}

	resp, err := sap.DryRun([]*component.ComponentDefinition{namespace, inDeclared, inUndeclared})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !resp["staging"]["ctx-a"].Success {
		t.Errorf("expected the Namespace to pass, got %+v", resp["staging"]["ctx-a"].Error)
	}
	if !resp["web"]["ctx-a"].Success {
		t.Errorf("expected the workload in the declared namespace to pass against default, got %+v", resp["web"]["ctx-a"].Error)
	}
	if resp["api"]["ctx-a"].Success {
		t.Errorf("expected the workload in the undeclared namespace to fail")
	}

	defaultRetries := 0
	for _, path := range cluster.requestedPaths() {
		if path == defaultPath {
			defaultRetries++
		}
	}
	if defaultRetries != 1 {
		t.Errorf("expected exactly one retry against default (for web only), got %d in %v", defaultRetries, cluster.requestedPaths())
	}
}

// TestDryRunMultiClusterWithFailingDiscovery pins that a context whose discovery fails falls back to the
// default resource mapping without affecting another context, and that both contexts get a result.
func TestDryRunMultiClusterWithFailingDiscovery(t *testing.T) {
	healthy := newFakeCluster(t, false, echoDryRun)
	broken := newFakeCluster(t, true, echoDryRun)

	web := dryRunTestComponent("web", "Ingress", "networking.k8s.io/v1", map[string]interface{}{
		"metadata": map[string]interface{}{"namespace": "team"},
	})

	sap := newMutateActionProvider(t)
	sap.ctxTokubeconfig = map[string]string{
		"healthy": healthy.kubeconfig(),
		"broken":  broken.kubeconfig(),
	}

	resp, err := sap.DryRun([]*component.ComponentDefinition{web})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp["web"]) != 2 {
		t.Fatalf("expected a result for both contexts, got %v", resp["web"])
	}
	for ctxID, r := range resp["web"] {
		if !r.Success {
			t.Errorf("context %s: expected success, got %+v", ctxID, r.Error)
		}
	}

	if got := healthy.requestedPaths(); len(got) != 1 || got[0] != "/apis/networking.k8s.io/v1/namespaces/team/ingresses" {
		t.Errorf("healthy context should use the discovered path, got %v", got)
	}
	if got := broken.requestedPaths(); len(got) != 1 || got[0] != "/apis/networking.k8s.io/v1/namespaces/team/ingresss" {
		t.Errorf("broken context should use the default mapping, got %v", got)
	}
}

// TestDryRunComponent_PreservesComponentOnDelete verifies that dry-run delete operations do not strip component configuration and preserve original spec.
func TestDryRunComponent_PreservesComponentOnDelete(t *testing.T) {
	originalConfig := map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      "my-workload",
			"namespace": "default",
		},
		"spec": map[string]interface{}{
			"replicas": 3,
		},
	}
	cmp := &component.ComponentDefinition{
		DisplayName: "my-workload",
		Component: component.Component{
			Kind:    "Deployment",
			Version: "apps/v1",
		},
		Configuration: originalConfig,
	}

	rClient := &fake.RESTClient{
		NegotiatedSerializer: serializer.NewCodecFactory(runtime.NewScheme()).WithoutConversion(),
		GroupVersion:         schema.GroupVersion{Group: "apps", Version: "v1"},
		Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete {
				t.Fatalf("expected DELETE request, got %s", req.Method)
			}
			status := metav1.Status{
				TypeMeta: metav1.TypeMeta{Kind: "Status", APIVersion: "v1"},
				Status:   metav1.StatusSuccess,
			}
			data, err := json.Marshal(status)
			if err != nil {
				return nil, err
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/json"}},
				Body:       io.NopCloser(bytes.NewReader(data)),
			}, nil
		}),
	}
	fakeClient := &meshkube.Client{
		KubeClient: kubernetes.New(rClient),
	}

	dResp, err := dryRunComponent(fakeClient, cmp, true, nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dResp.Success {
		t.Fatalf("expected dResp.Success to be true, got false, Error: %+v", dResp.Error)
	}
	if dResp.Component == nil || dResp.Component.Configuration["spec"] == nil {
		t.Errorf("expected dResp.Component.Configuration to retain spec, got: %v", dResp.Component)
	}
	// Original component configuration must remain untouched
	if cmp.Configuration["spec"] == nil {
		t.Errorf("original component configuration was mutated: %v", cmp.Configuration)
	}
}
