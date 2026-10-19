package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
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
	t.Run("returns name from metadata configuration", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "my-display-name",
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{
					"name": "custom-ns-name",
				},
			},
		}
		name := getNamespaceNameFromComponent(cmp)
		if name != "custom-ns-name" {
			t.Errorf("got %q, want %q", name, "custom-ns-name")
		}
	})

	t.Run("does not fall back to DisplayName when metadata name is empty", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "my-display-name",
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{"name": ""},
			},
		}
		if name := getNamespaceNameFromComponent(cmp); name != "" {
			t.Errorf("got %q, want empty string", name)
		}
	})

	t.Run("does not fall back to DisplayName when metadata name is not a string", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "my-display-name",
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{"name": 42},
			},
		}
		if name := getNamespaceNameFromComponent(cmp); name != "" {
			t.Errorf("got %q, want empty string", name)
		}
	})

	t.Run("does not fall back to DisplayName when metadata is malformed", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "my-display-name",
			Configuration: map[string]interface{}{
				"metadata": "not-a-map",
			},
		}
		if name := getNamespaceNameFromComponent(cmp); name != "" {
			t.Errorf("got %q, want empty string", name)
		}
	})

	t.Run("ignores a top-level Configuration name", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "my-display-name",
			Configuration: map[string]interface{}{
				"metadata": map[string]interface{}{"name": "real-ns"},
				"name":     "top-level-ns-name",
			},
		}
		if name := getNamespaceNameFromComponent(cmp); name != "real-ns" {
			t.Errorf("got %q, want %q", name, "real-ns")
		}
	})

	t.Run("falls back to DisplayName when metadata name missing", func(t *testing.T) {
		cmp := &component.ComponentDefinition{
			DisplayName: "fallback-display-name",
		}
		name := getNamespaceNameFromComponent(cmp)
		if name != "fallback-display-name" {
			t.Errorf("got %q, want %q", name, "fallback-display-name")
		}
	})

	t.Run("returns empty string for nil component", func(t *testing.T) {
		name := getNamespaceNameFromComponent(nil)
		if name != "" {
			t.Errorf("got %q, want empty string", name)
		}
	})
}

// TestCollectDeclaredNamespaces pins that a Namespace whose metadata.name is set is registered under that
// name only, so a workload targeting the component's DisplayName is not treated as in-design.
func TestCollectDeclaredNamespaces(t *testing.T) {
	ns := mutateTestComponent("display-ns", "Namespace")
	ns.Configuration = map[string]interface{}{
		"metadata": map[string]interface{}{"name": "real-ns"},
	}
	unnamed := mutateTestComponent("unnamed-ns", "Namespace")
	workload := mutateTestComponent("web", "Deployment")
	workload.Configuration = map[string]interface{}{
		"metadata": map[string]interface{}{"name": "other-name"},
	}

	got := collectDeclaredNamespaces([]*component.ComponentDefinition{ns, unnamed, workload, nil})

	want := map[string]bool{"real-ns": true, "unnamed-ns": true}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for name := range want {
		if !got[name] {
			t.Errorf("expected %q to be declared, got %v", name, got)
		}
	}
	if got["display-ns"] {
		t.Errorf("DisplayName %q must not be registered when metadata.name is set", "display-ns")
	}
}

// TestIsDryRunNotFoundError tests detection of Kubernetes 404/NotFound status errors and responses during dry run operations.
func TestIsDryRunNotFoundError(t *testing.T) {
	t.Run("detects kubeerror.StatusError NotFound", func(t *testing.T) {
		err := kubeerror.NewNotFound(schema.GroupResource{Group: "apps", Resource: "deployments"}, "web")
		if !isDryRunNotFoundError(err, nil) {
			t.Errorf("expected isDryRunNotFoundError to be true for kubeerror NotFound")
		}
	})

	t.Run("detects status map with reason NotFound", func(t *testing.T) {
		st := map[string]interface{}{
			"status": "Failure",
			"reason": "NotFound",
			"code":   float64(404),
		}
		if !isDryRunNotFoundError(nil, st) {
			t.Errorf("expected isDryRunNotFoundError to be true for status map reason NotFound")
		}
	})

	t.Run("does not treat plain message substring as not found", func(t *testing.T) {
		st := map[string]interface{}{
			"status":  "Failure",
			"message": "Resource was not found on server",
		}
		if isDryRunNotFoundError(nil, st) {
			t.Errorf("expected isDryRunNotFoundError to be false for plain message without 404 or NotFound reason")
		}
	})

	t.Run("returns false for non-404 error", func(t *testing.T) {
		err := kubeerror.NewForbidden(schema.GroupResource{Group: "apps", Resource: "deployments"}, "web", nil)
		if isDryRunNotFoundError(err, nil) {
			t.Errorf("expected isDryRunNotFoundError to be false for forbidden error")
		}
	})
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
