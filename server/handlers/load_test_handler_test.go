package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

type mockPublishProvider struct {
	models.DefaultLocalProvider
	publishFn func(req *http.Request, result *models.MesheryResult, profileID string) (string, error)
}

func (m *mockPublishProvider) PublishResults(req *http.Request, result *models.MesheryResult, profileID string) (string, error) {
	if m.publishFn != nil {
		return m.publishFn(req, result, profileID)
	}
	return "mock-result-id", nil
}

func TestPersistPerformanceTestResult_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	expectedUUID := uuid.Must(uuid.NewV4())
	mockP := &mockPublishProvider{
		publishFn: func(req *http.Request, result *models.MesheryResult, profileID string) (string, error) {
			if profileID != "test-profile-123" {
				t.Errorf("expected profileID 'test-profile-123', got %s", profileID)
			}
			return expectedUUID.String(), nil
		},
	}
	respChan := make(chan *models.LoadTestResponse, 5)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/perf/profile/test-profile-123/run", nil)
	result := &models.MesheryResult{
		Name: "test-run",
		Mesh: "istio",
	}

	h.persistPerformanceTestResult(context.Background(), req, result, "uuid-1", "test-profile-123", nil, nil, mockP, respChan)
	close(respChan)

	var responses []*models.LoadTestResponse
	for resp := range respChan {
		responses = append(responses, resp)
	}

	if len(responses) != 2 {
		t.Fatalf("expected 2 responses (info + success), got %d", len(responses))
	}
	if responses[0].Status != models.LoadTestInfo {
		t.Errorf("expected first response to be LoadTestInfo, got %v", responses[0].Status)
	}
	if responses[1].Status != models.LoadTestSuccess {
		t.Errorf("expected second response to be LoadTestSuccess, got %v", responses[1].Status)
	}
	if responses[1].Result.ID != expectedUUID {
		t.Errorf("expected result ID %s, got %s", expectedUUID, responses[1].Result.ID)
	}
}

func TestPersistPerformanceTestResult_ProviderError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	mockP := &mockPublishProvider{
		publishFn: func(req *http.Request, result *models.MesheryResult, profileID string) (string, error) {
			return "", errors.New("network timeout")
		},
	}
	respChan := make(chan *models.LoadTestResponse, 5)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "/api/perf/run", nil)
	result := &models.MesheryResult{Name: "err-run"}

	h.persistPerformanceTestResult(context.Background(), req, result, "uuid-1", "p-1", nil, nil, mockP, respChan)
	close(respChan)

	var responses []*models.LoadTestResponse
	for resp := range respChan {
		responses = append(responses, resp)
	}

	if len(responses) != 1 {
		t.Fatalf("expected 1 error response, got %d", len(responses))
	}
	if responses[0].Status != models.LoadTestError {
		t.Errorf("expected LoadTestError status, got %v", responses[0].Status)
	}
}

func TestPersistPerformanceTestResult_ContextCancelled(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	mockP := &mockPublishProvider{
		publishFn: func(req *http.Request, result *models.MesheryResult, profileID string) (string, error) {
			return "some-id", nil
		},
	}
	// Unbuffered channel that nobody is reading from
	respChan := make(chan *models.LoadTestResponse)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/api/perf/run", nil)
		h.persistPerformanceTestResult(ctx, req, &models.MesheryResult{}, "u", "p", nil, nil, mockP, respChan)
		close(done)
	}()

	select {
	case <-done:
		// Succeeded in exiting cleanly without hanging on respChan
	case <-time.After(1 * time.Second):
		t.Fatal("persistPerformanceTestResult hung on cancelled context")
	}
}

func TestLoadTest_ConcurrentContextIsolation(t *testing.T) {
	// Verify that multiple goroutines modifying clusterResultMap simultaneously
	// do NOT suffer from race conditions or cross-context map contamination.
	baseResultsMap := map[string]interface{}{
		"load-generator": "fortio",
		"qps":            100.0,
		"duration":       "10s",
	}

	const contextCount = 10
	var wg sync.WaitGroup
	results := make([]*models.MesheryResult, contextCount)

	for i := 0; i < contextCount; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			serverVersion := fmt.Sprintf("v1.2%d.0", idx)
			nodes := []*models.K8SNode{{InternalIP: fmt.Sprintf("10.0.0.%d", idx)}}
			clusterResultMap := buildClusterResultMap(baseResultsMap, serverVersion, nodes, nil)

			res := &models.MesheryResult{
				Name:   fmt.Sprintf("test-%d", idx),
				Result: clusterResultMap,
			}
			results[idx] = res
		}(i)
	}

	wg.Wait()

	// Verify isolation across all results
	for i := 0; i < contextCount; i++ {
		rMap := results[i].Result
		k8s, ok := rMap["kubernetes"].(map[string]interface{})
		if !ok {
			t.Fatalf("result %d missing kubernetes map", i)
		}
		expectedVersion := fmt.Sprintf("v1.2%d.0", i)
		if k8s["server_version"] != expectedVersion {
			t.Errorf("context %d contaminated: expected %s, got %v", i, expectedVersion, k8s["server_version"])
		}
	}
}

func TestLoadTestHelperHandler_SSEStreamingAndTeardown(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	rec := httptest.NewRecorder()

	ctx, cancel := context.WithCancel(context.Background())
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/api/perf/run", nil)

	prefObj := &models.Preference{}
	loadOpts := &models.LoadTestOptions{
		LoadGenerator: models.FortioLG,
	}
	mockP := &mockPublishProvider{}

	// Cancel context after a short duration to simulate client disconnection
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	done := make(chan struct{})
	go func() {
		h.loadTestHelperHandler(rec, req, "profile-1", "test-1", "istio", "uuid-1", prefObj, loadOpts, mockP)
		close(done)
	}()

	select {
	case <-done:
		// Clean exit without panic or hang!
	case <-time.After(2 * time.Second):
		t.Fatal("loadTestHelperHandler failed to terminate on client disconnection")
	}

	// Verify SSE headers
	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected Content-Type text/event-stream, got %s", rec.Header().Get("Content-Type"))
	}
}
