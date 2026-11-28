package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

// TestLoadExtensionFromPackageOpensPluginOnce checks that a failed plugin load
// is remembered rather than retried: a second plugin.Open of a package whose
// first load failed is a fatal runtime error that no handler can recover from.
func TestLoadExtensionFromPackageOpensPluginOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	props := models.ProviderProperties{
		ProviderName:   "TestProvider",
		PackageVersion: "v0.0.0-test",
	}
	props.Extensions.GraphQL = models.GraphQLExtensions{{Path: "graphql/plugin.so"}}
	provider := newMockProvider(props)

	pluginPath := filepath.Join(provider.PackageLocation(), "graphql", "plugin.so")
	if err := os.MkdirAll(filepath.Dir(pluginPath), 0o755); err != nil {
		t.Fatal(err)
	}
	// Not a shared object, so plugin.Open fails cleanly instead of loading.
	if err := os.WriteFile(pluginPath, []byte("not a plugin"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		mx.Lock()
		delete(loadedExtensionPackages, pluginPath)
		mx.Unlock()
	})

	h := &Handler{}
	first := h.LoadExtensionFromPackage(nil, nil, provider)
	if first == nil {
		t.Fatal("expected loading an invalid plugin to fail")
	}
	second := h.LoadExtensionFromPackage(nil, nil, provider)
	if second != first {
		t.Fatalf("expected the first load's error to be returned again without reopening the plugin; got %v, want %v", second, first)
	}

	mx.Lock()
	recorded, ok := loadedExtensionPackages[pluginPath]
	mx.Unlock()
	if !ok || recorded != first {
		t.Fatalf("expected the load outcome to be recorded for %s", pluginPath)
	}
}

// mockExtensionProxyProvider is a test provider mock for ExtensionProxy scenarios.
type mockExtensionProxyProvider struct {
	*models.DefaultLocalProvider
	err error
}

func (m *mockExtensionProxyProvider) ExtensionProxy(req *http.Request) (*models.ExtensionProxyResponse, error) {
	return nil, m.err
}

// TestExtensionsHandler_InvalidPathReturnsBadRequest verifies that client path validation
// errors return HTTP 400 Bad Request with the structured MeshKit error.
func TestExtensionsHandler_InvalidPathReturnsBadRequest(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatal(err)
	}
	h := &Handler{
		log: log,
	}

	mockP := &mockExtensionProxyProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		err:                  models.ErrInvalidExtensionProxyPath("/invalid/path"),
	}

	req := httptest.NewRequest(http.MethodGet, "/invalid/path", nil)
	rec := httptest.NewRecorder()

	h.ExtensionsHandler(rec, req, nil, nil, mockP)

	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), models.ErrInvalidExtensionProxyPathCode) {
		t.Fatalf("expected HTTP 400 with error code %q, got %d (body: %s)", models.ErrInvalidExtensionProxyPathCode, rec.Code, rec.Body.String())
	}
}
