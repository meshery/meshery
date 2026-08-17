package resources

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// mockMesheryClient simulates Meshery backend responses for testing
type mockMesheryClient struct {
	shouldFail bool
}

func (m *mockMesheryClient) GetConnections(ctx context.Context) ([]ConnectionResource, error) {
	if m.shouldFail {
		return nil, errors.New("failed to connect to meshery backend")
	}
	return []ConnectionResource{
		{
			ID:       "c1",
			Name:     "minikube",
			Type:     "kubernetes",
			Status:   "connected",
			Endpoint: "https://127.0.0.1:8443",
		},
	}, nil
}

func (m *mockMesheryClient) GetProviders(ctx context.Context) ([]ProviderResource, error) {
	if m.shouldFail {
		return nil, errors.New("failed to fetch providers")
	}
	return []ProviderResource{
		{Name: "Local", Status: "active", URL: "http://localhost:9081"},
	}, nil
}

func (m *mockMesheryClient) GetAdapters(ctx context.Context) ([]AdapterResource, error) {
	if m.shouldFail {
		return nil, errors.New("failed to fetch adapters")
	}
	return []AdapterResource{
		{Name: "meshery-istio", Version: "v1.10.0", Status: "running"},
	}, nil
}

func (m *mockMesheryClient) GetHealth(ctx context.Context) (HealthResponse, error) {
	if m.shouldFail {
		return HealthResponse{}, errors.New("health check failed")
	}
	return HealthResponse{
		Status:     "healthy",
		Components: map[string]string{"database": "ok"},
	}, nil
}

func (m *mockMesheryClient) GetEnvironments(ctx context.Context) ([]EnvironmentResource, error) {
	if m.shouldFail {
		return nil, errors.New("failed to fetch environments")
	}
	return []EnvironmentResource{
		{ID: "e1", Name: "dev", ConnectionIDs: []string{"c1"}},
	}, nil
}

func TestListResources(t *testing.T) {
	res := ListResources()
	if len(res) != 5 {
		t.Fatalf("expected 5 resources, got %d", len(res))
	}
}

func TestReadResource_Success(t *testing.T) {
	client := &mockMesheryClient{shouldFail: false}
	ctx := context.Background()

	tests := []struct {
		uri          string
		expectedText string
	}{
		{"meshery://connections", "minikube"},
		{"meshery://providers", "Local"},
		{"meshery://adapters", "meshery-istio"},
		{"meshery://health", "healthy"},
		{"meshery://environments", "dev"},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			data, err := ReadResource(ctx, tt.uri, client)
			if err != nil {
				t.Fatalf("unexpected error reading %s: %v", tt.uri, err)
			}

			output := string(data)
			if !strings.Contains(output, tt.expectedText) {
				t.Errorf("expected output to contain %q, got: %s", tt.expectedText, output)
			}
			if !strings.Contains(output, "timestamp") {
				t.Errorf("expected output to contain timestamp, got: %s", output)
			}
		})
	}
}

func TestReadResource_DisconnectedGracefulFallback(t *testing.T) {
	client := &mockMesheryClient{shouldFail: true}
	ctx := context.Background()

	// Health check when disconnected
	data, err := ReadResource(ctx, "meshery://health", client)
	if err != nil {
		t.Fatalf("expected graceful response, got error: %v", err)
	}
	if !strings.Contains(string(data), "unavailable") {
		t.Errorf("expected health status to be unavailable, got: %s", string(data))
	}

	// Connections check when disconnected
	data, err = ReadResource(ctx, "meshery://connections", client)
	if err != nil {
		t.Fatalf("expected graceful response for connections, got error: %v", err)
	}
	if !strings.Contains(string(data), `"connections":[]`) {
		t.Errorf("expected empty connections list, got: %s", string(data))
	}
}

func TestReadResource_InvalidURI(t *testing.T) {
	client := &mockMesheryClient{}
	ctx := context.Background()

	_, err := ReadResource(ctx, "meshery://invalid-resource", client)
	if err == nil {
		t.Fatal("expected error for invalid URI, got nil")
	}
}

func TestSanitizeJSON(t *testing.T) {
	rawWithSecrets := []byte(`{
		"id": "c1",
		"kubeconfig": "apiVersion: v1...",
		"token": "secret-jwt-token",
		"status": "connected"
	}`)

	sanitized, err := SanitizeJSON(rawWithSecrets)
	if err != nil {
		t.Fatalf("sanitizer failed: %v", err)
	}

	str := string(sanitized)
	if strings.Contains(str, "kubeconfig") || strings.Contains(str, "secret-jwt-token") {
		t.Errorf("sanitizer failed to strip sensitive fields, got: %s", str)
	}
	if !strings.Contains(str, "c1") || !strings.Contains(str, "connected") {
		t.Errorf("sanitizer improperly removed non-sensitive fields, got: %s", str)
	}
}
