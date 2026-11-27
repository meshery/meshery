package core

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/registry"
)

func TestDesignNameFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		expectedName string
	}{
		{
			name:         "given regular yaml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yaml",
			expectedName: "my-deployment",
		},
		{
			name:         "given regular yml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yml",
			expectedName: "my-deployment",
		},
		{
			name:         "given tar.gz file when DesignNameFromFileName then strip compound extension",
			fileName:     "my-chart.tar.gz",
			expectedName: "my-chart",
		},
		{
			name:         "given json file when DesignNameFromFileName then return name",
			fileName:     "config.json",
			expectedName: "config",
		},
		{
			name:         "given filename with multiple dots when DesignNameFromFileName then strip only last extension",
			fileName:     "my.k8s.deployment.yaml",
			expectedName: "my.k8s.deployment",
		},
		{
			name:         "given empty filename when DesignNameFromFileName then return empty",
			fileName:     "",
			expectedName: "",
		},
		{
			name:         "given filename without extension when DesignNameFromFileName then return as is",
			fileName:     "mydesign",
			expectedName: "mydesign",
		},
		{
			name:         "given tgz file when DesignNameFromFileName then strip extension",
			fileName:     "helm-chart.tgz",
			expectedName: "helm-chart",
		},
		{
			name:         "given extension only when DesignNameFromFileName then return empty",
			fileName:     ".yaml",
			expectedName: "",
		},
		{
			name:         "given tar.gz extension only when DesignNameFromFileName then return empty",
			fileName:     ".tar.gz",
			expectedName: "",
		},
		{
			name:         "given unsupported xls extension when DesignNameFromFileName then still strip extension",
			fileName:     "spreadsheet.xls",
			expectedName: "spreadsheet",
		},
		{
			name:         "given unsupported zip extension when DesignNameFromFileName then still strip extension",
			fileName:     "archive.zip",
			expectedName: "archive",
		},
		{
			name:         "given unsupported exe extension when DesignNameFromFileName then still strip extension",
			fileName:     "installer.exe",
			expectedName: "installer",
		},
		{
			name:         "given unsupported tar extension when DesignNameFromFileName then strip last suffix only",
			fileName:     "archive.tar",
			expectedName: "archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DesignNameFromFileName(tt.fileName)
			if result != tt.expectedName {
				t.Errorf("expected %q, got %q", tt.expectedName, result)
			}
		})
	}
}

// newEmptyTestRegistry builds a file-backed registry holding zero registered
// components, reproducing the startup window in which a fresh server has not
// finished model registration yet (meshmodel schema exists, but no rows).
func newEmptyTestRegistry(t *testing.T) *registry.RegistryManager {
	t.Helper()

	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: filepath.Join(t.TempDir(), "registry.db")})
	if err != nil {
		t.Fatalf("failed to open a database: %v", err)
	}

	reg, err := registry.NewRegistryManager(&db)
	if err != nil {
		t.Fatalf("failed to build a registry manager: %v", err)
	}

	return reg
}

func TestNewPatternFileFromK8sManifest_AllResourcesUnmatched(t *testing.T) {
	// Arrange: a valid two-resource manifest and an empty registry
	reg := newEmptyTestRegistry(t)
	manifest := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: web
---
apiVersion: v1
kind: Service
metadata:
  name: web
`

	// Act
	design, err := NewPatternFileFromK8sManifest(manifest, "web.yaml", true, reg)

	// Assert: the import must fail loudly and name the unmatched resources
	// instead of silently returning a design with zero components.
	if err == nil {
		t.Fatalf("expected an error when no manifest resources match a registered component, got a design with %d components and no error", len(design.Components))
	}
	for _, want := range []string{"apps/v1 Deployment", "v1 Service"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("expected error to name %q, got: %v", want, err)
		}
	}
}

func TestNewPatternFileFromK8sManifest_EmptyManifest(t *testing.T) {
	// Arrange
	reg := newEmptyTestRegistry(t)

	// Act
	_, err := NewPatternFileFromK8sManifest("", "empty.yaml", true, reg)

	// Assert: the pre-existing empty-manifest error must be preserved
	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Errorf("expected an \"empty manifest\" error, got: %v", err)
	}
}
