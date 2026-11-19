package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func newTestDB(t *testing.T) *database.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create in-memory test database: %v", err)
	}
	return &database.Handler{DB: db, Mutex: &sync.Mutex{}}
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
		setupDB          func(t *testing.T) *database.Handler
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
			expectedContains: []string{"[+]database ok", "[+]capabilities ok", "[i]extension", "healthz check passed"},
		},
		{
			name:    "Unhealthy - no capabilities loaded (extension status doesn't affect health)",
			verbose: false,
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
			expectedContains: []string{"[+]database ok", "[-]capabilities failed", "not loaded", "[i]extension"},
		},
		{
			name:    "Unhealthy - database connection pool closed",
			verbose: false,
			setupDB: func(t *testing.T) *database.Handler {
				db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to open in-memory sqlite: %v", err)
				}
				sqlDB, err := db.DB()
				if err != nil {
					t.Fatalf("failed to get sql.DB: %v", err)
				}
				_ = sqlDB.Close()
				return &database.Handler{DB: db, Mutex: &sync.Mutex{}}
			},
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedContains: []string{"healthz check failed", "database: database ping failed"},
		},
		{
			name:    "Unhealthy - database connection pool closed, verbose mode",
			verbose: true,
			setupDB: func(t *testing.T) *database.Handler {
				db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
				if err != nil {
					t.Fatalf("failed to open in-memory sqlite: %v", err)
				}
				sqlDB, err := db.DB()
				if err != nil {
					t.Fatalf("failed to get sql.DB: %v", err)
				}
				_ = sqlDB.Close()
				return &database.Handler{DB: db, Mutex: &sync.Mutex{}}
			},
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedContains: []string{"[-]database failed: database ping failed", "[+]capabilities ok"},
		},
		{
			name:    "Unhealthy - database handler is nil",
			verbose: false,
			setupDB: func(t *testing.T) *database.Handler {
				return nil
			},
			setupProviders: func() map[string]models.Provider {
				provider := &models.DefaultLocalProvider{}
				provider.Initialize()
				return map[string]models.Provider{
					"Local": provider,
				}
			},
			expectedStatus:   http.StatusServiceUnavailable,
			expectedContains: []string{"healthz check failed", "database: database handler not initialized"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			providers := tt.setupProviders()
			handlerConfig := &models.HandlerConfig{
				Providers: providers,
			}
			var dbHandler *database.Handler
			if tt.setupDB != nil {
				dbHandler = tt.setupDB(t)
			} else {
				dbHandler = newTestDB(t)
			}

			handler := &Handler{
				config:    handlerConfig,
				log:       log,
				dbHandler: dbHandler,
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
			if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
				t.Errorf("Expected X-Content-Type-Options 'nosniff', got '%s'", nosniff)
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
	t.Cleanup(func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("Failed to clean temp dir: %v", err)
		}
	})

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
		config:    handlerConfig,
		log:       log,
		dbHandler: newTestDB(t),
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

func TestK8sLivenessHandler(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Liveness must return 200 OK even when provider capabilities are empty
	// and dbHandler is nil.
	provider := newMockProvider(models.ProviderProperties{
		ProviderName: "uninitialized-provider",
		Capabilities: models.Capabilities{},
	})
	handler := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{"uninit": provider},
		},
		log:       log,
		dbHandler: nil, // DB not initialized
	}

	t.Run("Liveness returns 200 without verbose", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
		w := httptest.NewRecorder()

		handler.K8sLivenessHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if body := w.Body.String(); body != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", body)
		}
		if contentType := w.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "text/plain") {
			t.Errorf("Expected text/plain content type, got '%s'", contentType)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected X-Content-Type-Options 'nosniff', got '%s'", nosniff)
		}
	})

	t.Run("Liveness returns 200 with verbose", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/live?verbose=1", nil)
		w := httptest.NewRecorder()

		handler.K8sLivenessHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "[+]live ok") {
			t.Errorf("Expected verbose output to contain '[+]live ok', got '%s'", body)
		}
		if !strings.Contains(body, "healthz check passed") {
			t.Errorf("Expected verbose output to contain 'healthz check passed', got '%s'", body)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected X-Content-Type-Options 'nosniff', got '%s'", nosniff)
		}
	})
}

func TestK8sHealthzHandler_PathDispatch(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Handler with empty capabilities and nil database (unready state)
	unreadyProvider := newMockProvider(models.ProviderProperties{
		ProviderName: "test-provider",
		Capabilities: models.Capabilities{},
	})
	unreadyHandler := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{"test": unreadyProvider},
		},
		log:       log,
		dbHandler: nil,
	}

	// Handler with initialized capabilities and healthy database (ready state)
	readyProvider := &models.DefaultLocalProvider{}
	readyProvider.Initialize()
	readyHandler := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{"Local": readyProvider},
		},
		log:       log,
		dbHandler: newTestDB(t),
	}

	t.Run("/healthz/live routes to liveness and succeeds even when unready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
		w := httptest.NewRecorder()

		unreadyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK on /healthz/live, got %d (body: %s)", w.Code, w.Body.String())
		}
		if w.Body.String() != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", w.Body.String())
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("/healthz/live/ with trailing slash routes to liveness and succeeds", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/live/", nil)
		w := httptest.NewRecorder()

		unreadyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK on /healthz/live/, got %d (body: %s)", w.Code, w.Body.String())
		}
		if w.Body.String() != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", w.Body.String())
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("/healthz/ready routes to readiness and reports 503 when unready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
		w := httptest.NewRecorder()

		unreadyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected 503 on /healthz/ready when unready, got %d", w.Code)
		}
		body := w.Body.String()
		if !strings.Contains(body, "healthz check failed") {
			t.Errorf("Expected body to contain failure details, got '%s'", body)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("/healthz/ready/ with trailing slash routes to readiness and reports 503 when unready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/ready/", nil)
		w := httptest.NewRecorder()

		unreadyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected 503 on /healthz/ready/ when unready, got %d", w.Code)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("/healthz/ready succeeds when dependencies are ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
		w := httptest.NewRecorder()

		readyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK on /healthz/ready when ready, got %d", w.Code)
		}
		if w.Body.String() != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", w.Body.String())
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("/healthz/ready/ with trailing slash succeeds when dependencies are ready", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/ready/", nil)
		w := httptest.NewRecorder()

		readyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200 OK on /healthz/ready/ when ready, got %d", w.Code)
		}
		if w.Body.String() != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", w.Body.String())
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	// Negative path tests: Unrelated and noncanonical paths must NOT be dispatched to liveness or readiness
	unrelatedPaths := []string{
		"/other/live",
		"/other/live/",
		"/other/ready",
		"/other/ready/",
		"/healthz",
		"/healthz/",
		"/api/system/healthz/live",
		"/healthz/../healthz/live",
		"/healthz/../healthz/ready",
		"/healthz/live/subpath",
		"/healthz/ready/subpath",
		"/unknown",
		"/",
	}

	for _, p := range unrelatedPaths {
		t.Run("Arbitrary or noncanonical path rejected with 404: "+p, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, p, nil)
			w := httptest.NewRecorder()

			readyHandler.K8sHealthzHandler(w, req)

			if w.Code != http.StatusNotFound {
				t.Errorf("Expected 404 Not Found for arbitrary path %q, got status %d (body: %s)", p, w.Code, w.Body.String())
			}
		})
	}

	t.Run("Empty path rejected with 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
		req.URL.Path = ""
		w := httptest.NewRecorder()

		readyHandler.K8sHealthzHandler(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found for empty path, got status %d (body: %s)", w.Code, w.Body.String())
		}
	})
}

// TestK8sReadinessHandler_ContextCancellation tests context cancellation error propagation.
// Note: Testing an already-canceled context validates that PingContext detects and propagates
// context cancellation, yielding HTTP 503 with "context canceled" rather than hanging.
func TestK8sReadinessHandler_ContextCancellation(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	provider := &models.DefaultLocalProvider{}
	provider.Initialize()

	handler := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{"Local": provider},
		},
		log:       log,
		dbHandler: newTestDB(t),
	}

	// Create a request with an already-canceled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.K8sReadinessHandler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected status %d when context is canceled, got %d", http.StatusServiceUnavailable, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "database ping failed: context canceled") {
		t.Errorf("Expected body to indicate context canceled, got: %s", body)
	}
	if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("Expected nosniff header, got '%s'", nosniff)
	}
}

// TestK8sReadinessHandler_ContextDeadlineExceeded tests context deadline expiration error propagation.
// Note: Passing an already-expired context tests deterministic handling of context.DeadlineExceeded.
func TestK8sReadinessHandler_ContextDeadlineExceeded(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	provider := &models.DefaultLocalProvider{}
	provider.Initialize()

	handler := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{"Local": provider},
		},
		log:       log,
		dbHandler: newTestDB(t),
	}

	// Create a request with an already-expired deadline
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-10*time.Millisecond))
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	handler.K8sReadinessHandler(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Fatalf("Expected status %d when context deadline is exceeded, got %d", http.StatusServiceUnavailable, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "database ping failed: context deadline exceeded") {
		t.Errorf("Expected body to indicate context deadline exceeded, got: %s", body)
	}
	if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("Expected nosniff header, got '%s'", nosniff)
	}
}

func TestK8sReadinessHandler_Direct(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	t.Run("Direct readiness call returns 200 when ready", func(t *testing.T) {
		provider := &models.DefaultLocalProvider{}
		provider.Initialize()

		handler := &Handler{
			config: &models.HandlerConfig{
				Providers: map[string]models.Provider{"Local": provider},
			},
			log:       log,
			dbHandler: newTestDB(t),
		}

		req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
		w := httptest.NewRecorder()

		handler.K8sReadinessHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if body := w.Body.String(); body != "ok" {
			t.Errorf("Expected body 'ok', got '%s'", body)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})

	t.Run("Direct readiness call returns 503 when database is down", func(t *testing.T) {
		provider := &models.DefaultLocalProvider{}
		provider.Initialize()

		dbHandler := newTestDB(t)
		sqlDB, err := dbHandler.DB.DB()
		if err != nil {
			t.Fatalf("Failed to get sql.DB: %v", err)
		}
		_ = sqlDB.Close()

		handler := &Handler{
			config: &models.HandlerConfig{
				Providers: map[string]models.Provider{"Local": provider},
			},
			log:       log,
			dbHandler: dbHandler,
		}

		req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
		w := httptest.NewRecorder()

		handler.K8sReadinessHandler(w, req)

		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("Expected status %d, got %d", http.StatusServiceUnavailable, w.Code)
		}
		if body := w.Body.String(); !strings.Contains(body, "database: database ping failed") {
			t.Errorf("Expected body to contain database failure, got '%s'", body)
		}
		if nosniff := w.Header().Get("X-Content-Type-Options"); nosniff != "nosniff" {
			t.Errorf("Expected nosniff header, got '%s'", nosniff)
		}
	})
}
