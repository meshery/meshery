package kubernetes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	mkkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/schemas/models/core"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type mockProviderForDiscoverTest struct {
	models.Provider
	savedContext models.K8sContext
}

func (m *mockProviderForDiscoverTest) SaveK8sContext(token string, k8sContext models.K8sContext, metadata map[string]any) (connections.Connection, error) {
	m.savedContext = k8sContext
	return connections.Connection{ID: core.Uuid(uuid.Must(uuid.NewV4()))}, nil
}

func TestDiscoverAction_PropagatesServerIDAndVersionToMachineCtx(t *testing.T) {
	clusterUID := "30e099e0-ef31-477a-9f2d-754ee9b1d41a"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/namespaces/kube-system" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"metadata": map[string]string{"uid": clusterUID},
			})
			return
		}
		if r.URL.Path == "/version" {
			_ = json.NewEncoder(w).Encode(map[string]string{
				"major": "1",
				"minor": "30",
				"gitVersion": "v1.30.0",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	clientset, err := k8sclient.NewForConfig(&rest.Config{Host: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatal(err)
	}

	mockProv := &mockProviderForDiscoverTest{}
	machineCtx := &MachineCtx{
		K8sContext: models.K8sContext{
			Name:   "test-ctx",
			Server: server.URL,
		},
		clientset: &mkkube.Client{KubeClient: clientset},
		log:       log,
	}

	ctx := context.WithValue(context.Background(), models.ProviderCtxKey, mockProv)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")

	da := &DiscoverAction{}
	nextEvent, _, err := da.Execute(ctx, machineCtx, nil)
	if err != nil {
		t.Fatalf("unexpected DiscoverAction error: %v", err)
	}
	if nextEvent != machines.Register {
		t.Fatalf("expected next event %v, got %v", machines.Register, nextEvent)
	}

	// Verify that machineCtx.K8sContext has been updated in memory
	if machineCtx.K8sContext.KubernetesServerID == nil {
		t.Fatalf("expected machineCtx.K8sContext.KubernetesServerID to be populated, but it was nil")
	}
	if machineCtx.K8sContext.KubernetesServerID.String() != clusterUID {
		t.Fatalf("expected KubernetesServerID %s, got %s", clusterUID, machineCtx.K8sContext.KubernetesServerID.String())
	}
	if machineCtx.K8sContext.Version != "v1.30.0" {
		t.Fatalf("expected Version v1.30.0, got %s", machineCtx.K8sContext.Version)
	}

	// Verify K8sContextSnapshot also reflects the updated server ID
	snap := machineCtx.K8sContextSnapshot()
	if snap.KubernetesServerID == nil || snap.KubernetesServerID.String() != clusterUID {
		t.Fatalf("expected K8sContextSnapshot().KubernetesServerID to be %s, got %v", clusterUID, snap.KubernetesServerID)
	}
}

func TestRefreshK8sContext_SynchronizesMissingServerID(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatal(err)
	}

	existingID := core.Uuid(uuid.Must(uuid.NewV4()))
	m := &MachineCtx{
		K8sContext: models.K8sContext{
			Name:   "cluster-1",
			Server: "https://k8s.example.com:6443",
		},
		clientset: &mkkube.Client{},
		log:       log,
	}

	newCtx := models.K8sContext{
		Name:               "cluster-1",
		Server:             "https://k8s.example.com:6443",
		KubernetesServerID: &existingID,
	}

	changed, err := m.RefreshK8sContext(newCtx, nil, log)
	if err != nil {
		t.Fatalf("unexpected RefreshK8sContext error: %v", err)
	}
	if !changed {
		t.Fatalf("expected changed to be true when populating missing server ID")
	}
	if m.K8sContext.KubernetesServerID == nil || *m.K8sContext.KubernetesServerID != existingID {
		t.Fatalf("expected KubernetesServerID to be %s, got %v", existingID, m.K8sContext.KubernetesServerID)
	}
}

func TestRefreshK8sContext_ReturnsErrorOnInvalidCredentials(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatal(err)
	}

	m := &MachineCtx{
		K8sContext: models.K8sContext{
			Name:   "cluster-1",
			Server: "https://k8s.example.com:6443",
		},
		log: log,
	}

	badCtx := models.K8sContext{
		Name:   "cluster-1",
		Server: "https://k8s-new.example.com:6443",
	}

	changed, err := m.RefreshK8sContext(badCtx, events.NewEvent(), log)
	if err == nil {
		t.Fatalf("expected error on invalid credentials, got nil")
	}
	if changed {
		t.Fatalf("expected changed=false on error, got true")
	}
}
