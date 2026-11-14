package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

// mockAdaptersTracker is an in-memory implementation of models.AdaptersTrackerInterface for testing.
type mockAdaptersTracker struct {
	mu       sync.Mutex
	adapters map[string]models.Adapter
}

func newMockAdaptersTracker(initial ...models.Adapter) *mockAdaptersTracker {
	m := &mockAdaptersTracker{
		adapters: make(map[string]models.Adapter),
	}
	for _, a := range initial {
		m.adapters[a.Location] = a
	}
	return m
}

func (m *mockAdaptersTracker) AddAdapter(_ context.Context, adapter models.Adapter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.adapters[adapter.Location] = adapter
}

func (m *mockAdaptersTracker) RemoveAdapter(_ context.Context, adapter models.Adapter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.adapters, adapter.Location)
}

func (m *mockAdaptersTracker) GetAdapters(_ context.Context) []models.Adapter {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]models.Adapter, 0, len(m.adapters))
	for _, a := range m.adapters {
		res = append(res, a)
	}
	return res
}

func (m *mockAdaptersTracker) DeployAdapter(_ context.Context, _ models.Adapter) error {
	return nil
}

func (m *mockAdaptersTracker) UndeployAdapter(_ context.Context, _ models.Adapter) error {
	return nil
}

// TestDeleteAdapterRemovesFromAdapterTracker verifies that deleteAdapter removes the adapter
// from both the preference slice and the runtime AdapterTracker to prevent resurrection on session sync.
func TestDeleteAdapterRemovesFromAdapterTracker(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	adapterA := &models.Adapter{
		Location: "meshery-istio:10000",
		Name:     "meshery-istio",
	}
	adapterB := &models.Adapter{
		Location: "meshery-linkerd:10001",
		Name:     "meshery-linkerd",
	}

	tracker := newMockAdaptersTracker(*adapterA, *adapterB)

	h := &Handler{
		config: &models.HandlerConfig{
			AdapterTracker: tracker,
		},
		log: log,
	}

	meshAdapters := []*models.Adapter{adapterA, adapterB}

	// Delete adapterA via DELETE request
	req := httptest.NewRequest(http.MethodDelete, "/api/system/adapter/manage?adapter=meshery-istio:10000", nil)
	rec := httptest.NewRecorder()

	updatedAdapters, err := h.deleteAdapter(meshAdapters, rec, req)
	if err != nil {
		t.Fatalf("deleteAdapter failed unexpectedly: %v", err)
	}

	// Verify slice removal
	if len(updatedAdapters) != 1 {
		t.Fatalf("expected 1 remaining adapter in slice, got %d", len(updatedAdapters))
	}
	if updatedAdapters[0].Location != "meshery-linkerd:10001" {
		t.Errorf("expected remaining adapter to be meshery-linkerd:10001, got %q", updatedAdapters[0].Location)
	}

	// Verify removal from AdapterTracker
	remainingInTracker := tracker.GetAdapters(context.Background())
	if len(remainingInTracker) != 1 {
		t.Fatalf("expected 1 remaining adapter in tracker, got %d", len(remainingInTracker))
	}
	if remainingInTracker[0].Location != "meshery-linkerd:10001" {
		t.Errorf("expected adapter in tracker to be meshery-linkerd:10001, got %q", remainingInTracker[0].Location)
	}

	// Verify non-existent adapter returns error and does not mutate tracker
	reqNonExistent := httptest.NewRequest(http.MethodDelete, "/api/system/adapter/manage?adapter=non-existent:9999", nil)
	recNonExistent := httptest.NewRecorder()

	_, err = h.deleteAdapter(updatedAdapters, recNonExistent, reqNonExistent)
	if err != ErrValidAdapter {
		t.Errorf("expected ErrValidAdapter for non-existent adapter, got %v", err)
	}
	if recNonExistent.Code != http.StatusBadRequest {
		t.Errorf("expected status %d for non-existent adapter, got %d", http.StatusBadRequest, recNonExistent.Code)
	}
}

// TestMeshOpsHandlerFailsGracefullyOnInvalidKubeConfig verifies that MeshOpsHandler writes an
// error response when generating a kubeconfig fails, instead of returning an empty unwritten HTTP response.
func TestMeshOpsHandlerFailsGracefullyOnInvalidKubeConfig(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	adapter := &models.Adapter{
		Location: "meshery-istio:10000",
		Name:     "meshery-istio",
	}

	h := &Handler{
		config: &models.HandlerConfig{
			AdapterTracker: newMockAdaptersTracker(*adapter),
		},
		log: log,
	}

	prefObj := &models.Preference{
		MeshAdapters: []*models.Adapter{adapter},
	}
	user := &models.User{
		FirstName: "test",
	}

	form := url.Values{}
	form.Set("adapter", "meshery-istio:10000")
	form.Set("query", "test-op")

	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/operation", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Inject an unmarshalable K8sContext so GenerateKubeConfig returns ErrMarshal
	invalidCluster := models.K8sContext{
		Name: "broken-cluster",
		Cluster: map[string]interface{}{
			"unmarshalable": make(chan int),
		},
	}
	ctx := context.WithValue(req.Context(), models.KubeClustersKey, []models.K8sContext{invalidCluster})
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	h.MeshOpsHandler(rec, req, prefObj, user, &models.DefaultLocalProvider{})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected HTTP %d on invalid kubeconfig, got %d", http.StatusBadRequest, rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "broken-cluster") && !strings.Contains(body, "kube config") {
		t.Errorf("expected error body to reference context or kube config error, got %q", body)
	}
}
