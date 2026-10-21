package handlers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/server/models"
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
