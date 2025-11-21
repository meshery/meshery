package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

// mockProvider embeds DefaultLocalProvider and overrides specific methods for testing
type mockProvider struct {
	*models.DefaultLocalProvider
	props models.ProviderProperties
}

func newMockProvider(props models.ProviderProperties) *mockProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &mockProvider{
		DefaultLocalProvider: base,
		props:                props,
	}
}

func (m *mockProvider) GetProviderProperties() models.ProviderProperties {
	return m.props
}

func (m *mockProvider) PackageLocation() string {
	return m.props.PackageLocation()
}

func (m *mockProvider) SetProviderProperties(props models.ProviderProperties) {
	m.props = props
}

func TestK8sHealthzHandler(t *testing.T) {
	// Setup test logger
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name             string
		verbose          bool
		setupProviders   func() map[string]models.Provider
		expectedStatus   int
		expectedBody     string
		expectedContains []string // strings that should be in the response
	}{
		{
			name:    "Healthy - capabilities loaded",
			verbose: false,
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:    "Healthy - capabilities loaded, verbose mode shows extension info",
			verbose: true,
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			expectedStatus:   http.StatusOK,
			expectedContains: []string{"[+]capabilities ok", "[i]extension", "healthz check passed"},
		},
		{
			name:    "Unhealthy - no capabilities loaded (extension status doesn't affect health)",
			verbose: false,
			setupProviders: func() map[string]models.Provider {
				// Create a mock provider with no capabilities
				provider := newMockProvider(models.ProviderProperties{
					ProviderName: "test-provider",
					Capabilities: models.Capabilities{},
				})
				return map[string]models.Provider{
					"test": provider,
				}
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedContains: []string{"healthz check failed", "capabilities", "not loaded"},
		},
		{
			name:    "Unhealthy - no capabilities loaded, verbose mode shows extension info",
			verbose: true,
			setupProviders: func() map[string]models.Provider {
				provider := newMockProvider(models.ProviderProperties{
					ProviderName: "test-provider",
					Capabilities: models.Capabilities{},
				})
				return map[string]models.Provider{
					"test": provider,
				}
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedContains: []string{"[-]capabilities failed", "not loaded", "[i]extension"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with test config
			providers := tt.setupProviders()
			handlerConfig := &models.HandlerConfig{
				Providers: providers,
			}
			handler := &Handler{
				config: handlerConfig,
				log:    log,
			}

			// Create test request with optional verbose parameter
			url := "/healthz/ready"
			if tt.verbose {
				url += "?verbose=1"
			}
			req := httptest.NewRequest(http.MethodGet, url, nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.K8sHealthzHandler(w, req)

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, "text/plain") {
				t.Errorf("Expected Content-Type to start with 'text/plain', got '%s'", contentType)
			}

			// Check response body
			body := w.Body.String()
			if tt.expectedBody != "" && body != tt.expectedBody {
				t.Errorf("Expected body '%s', got '%s'", tt.expectedBody, body)
			}

			// Check expected content strings
			for _, expected := range tt.expectedContains {
				if !strings.Contains(body, expected) {
					t.Errorf("Expected body to contain '%s', got '%s'", expected, body)
				}
			}

			t.Logf("Response body: %s", body)
		})
	}
}

func TestK8sHealthzHandler_ExtensionInfo(t *testing.T) {
	// Setup test logger
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create a temp directory structure for testing
	tempDir, err := os.MkdirTemp("", "meshery-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create extension file structure
	packagePath := path.Join(tempDir, "provider", "test-provider", "1.0.0")
	err = os.MkdirAll(packagePath, 0755)
	if err != nil {
		t.Fatalf("Failed to create package path: %v", err)
	}

	// Create a dummy extension file
	extensionFile := path.Join(packagePath, "extension.so")
	err = os.WriteFile(extensionFile, []byte("dummy"), 0644)
	if err != nil {
		t.Fatalf("Failed to create extension file: %v", err)
	}

	// Create a mock provider with capabilities and extensions
	provider := newMockProvider(models.ProviderProperties{
		ProviderName:   "test-provider",
		PackageVersion: "1.0.0",
		Capabilities: models.Capabilities{
			{
				Feature:  models.SyncPrefs,
				Endpoint: "/user/preferences",
			},
		},
		Extensions: models.Extensions{
			Navigator: models.NavigatorExtensions{
				{
					Title: "Test Extension",
					Href: models.Href{
						URI: "/extension/test",
					},
				},
			},
			GraphQL: models.GraphQLExtensions{
				{
					Component: "graphql",
					Path:      "extension.so",
					Type:      "navigator",
				},
			},
		},
	})

	handlerConfig := &models.HandlerConfig{
		Providers: map[string]models.Provider{
			"test": provider,
		},
	}
	handler := &Handler{
		config: handlerConfig,
		log:    log,
	}

	// Test 1: Without verbose - should be healthy even if extension not found (informational only)
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	w := httptest.NewRecorder()
	handler.K8sHealthzHandler(w, req)

	// Should be healthy because extension status doesn't affect health
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d (extension status should not affect health)", http.StatusOK, w.Code)
	}

	body := w.Body.String()
	// Non-verbose mode shouldn't mention extensions in error messages
	if body != "ok" {
		t.Errorf("Expected 'ok', got: %s", body)
	}

	t.Logf("Extension check response (non-verbose): %s", body)

	// Test 2: With verbose mode - should show extension info
	req = httptest.NewRequest(http.MethodGet, "/healthz/ready?verbose=1", nil)
	w = httptest.NewRecorder()
	handler.K8sHealthzHandler(w, req)

	body = w.Body.String()
	if !strings.Contains(body, "[+]capabilities ok") {
		t.Errorf("Expected verbose output to show capabilities check, got: %s", body)
	}
	// Should show informational extension status with [i] prefix
	if !strings.Contains(body, "[i]extension") {
		t.Errorf("Expected verbose output to show extension info with [i] prefix, got: %s", body)
	}

	t.Logf("Extension check response (verbose): %s", body)
}
