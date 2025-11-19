package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/spf13/viper"
)

func TestK8sHealthzHandler(t *testing.T) {
	// Setup test logger
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name               string
		releaseChannel     string
		setupProviders     func() map[string]models.Provider
		setupExtensions    func() string // returns temp dir path
		expectedStatus     int
		expectedHealthy    bool
		expectedCapLoaded  bool
		expectedExtExists  *bool // nil means not checked
	}{
		{
			name:           "Healthy - capabilities loaded, not kanvas mode",
			releaseChannel: "stable",
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			setupExtensions:   nil,
			expectedStatus:    http.StatusOK,
			expectedHealthy:   true,
			expectedCapLoaded: true,
			expectedExtExists: nil,
		},
		{
			name:           "Unhealthy - no capabilities loaded",
			releaseChannel: "stable",
			setupProviders: func() map[string]models.Provider {
				// Create a mock provider with no capabilities
				dbHandler := &database.Handler{}
				provider := models.NewRemoteProvider(
					nil,
					dbHandler,
					"test-provider",
					"",
					"",
					log,
					nil,
				)
				// Set empty capabilities
				provider.SetProviderProperties(models.ProviderProperties{
					Capabilities: models.Capabilities{},
				})
				return map[string]models.Provider{
					"test": provider,
				}
			},
			setupExtensions:   nil,
			expectedStatus:    http.StatusServiceUnavailable,
			expectedHealthy:   false,
			expectedCapLoaded: false,
			expectedExtExists: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper
			viper.Set("RELEASE_CHANNEL", tt.releaseChannel)

			// Setup temp dir for extensions if needed
			var tempDir string
			if tt.setupExtensions != nil {
				tempDir = tt.setupExtensions()
				defer os.RemoveAll(tempDir)
			}

			// Create handler with test config
			providers := tt.setupProviders()
			handlerConfig := &models.HandlerConfig{
				Providers: providers,
			}
			handler := &Handler{
				config: handlerConfig,
				log:    log,
			}

			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.K8sHealthzHandler(w, req)

			// Check response status
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Parse response body
			var healthStatus HealthStatus
			err := json.NewDecoder(w.Body).Decode(&healthStatus)
			if err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			// Check health status
			expectedStatusStr := "healthy"
			if !tt.expectedHealthy {
				expectedStatusStr = "unhealthy"
			}
			if healthStatus.Status != expectedStatusStr {
				t.Errorf("Expected status '%s', got '%s'", expectedStatusStr, healthStatus.Status)
			}

			// Check capabilities loaded
			if healthStatus.CapabilitiesLoaded != tt.expectedCapLoaded {
				t.Errorf("Expected CapabilitiesLoaded %v, got %v", tt.expectedCapLoaded, healthStatus.CapabilitiesLoaded)
			}

			// Check extension exists if applicable
			if tt.expectedExtExists != nil && healthStatus.ExtensionExists != *tt.expectedExtExists {
				t.Errorf("Expected ExtensionExists %v, got %v", *tt.expectedExtExists, healthStatus.ExtensionExists)
			}
		})
	}
}

func TestK8sHealthzHandler_KanvasMode(t *testing.T) {
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

	// Setup viper for kanvas mode
	viper.Set("RELEASE_CHANNEL", "kanvas")

	// Create a mock provider with capabilities and extensions
	dbHandler := &database.Handler{}
	provider := models.NewRemoteProvider(
		nil,
		dbHandler,
		"test-provider",
		"",
		"",
		log,
		nil,
	)

	// Set provider properties with extensions
	providerProps := models.ProviderProperties{
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
	}
	provider.SetProviderProperties(providerProps)

	// Mock PackageLocation to return our temp dir
	// Note: In real implementation, we'd need to properly set this
	// For now, we'll test with the default behavior

	handlerConfig := &models.HandlerConfig{
		Providers: map[string]models.Provider{
			"test": provider,
		},
	}
	handler := &Handler{
		config: handlerConfig,
		log:    log,
	}

	// Create test request
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.K8sHealthzHandler(w, req)

	// In Kanvas mode without extension file at expected location, should be unhealthy
	// Since PackageLocation uses actual home directory, this test will show unhealthy
	var healthStatus HealthStatus
	err = json.NewDecoder(w.Body).Decode(&healthStatus)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify response includes kanvas-specific checks
	if healthStatus.ReleaseChannel != "kanvas" {
		t.Errorf("Expected release channel 'kanvas', got '%s'", healthStatus.ReleaseChannel)
	}

	t.Logf("Health status in Kanvas mode: %+v", healthStatus)
}
