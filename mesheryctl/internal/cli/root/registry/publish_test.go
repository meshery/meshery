package registry

import (
	"os"
	"path/filepath"
	"testing"

	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
)

func TestResolveModelOCIVersion(t *testing.T) {
	tests := []struct {
		name     string
		versions []string
		expected string
	}{
		{
			name:     "single version",
			versions: []string{"0.1.0"},
			expected: "0.1.0",
		},
		{
			name:     "latest semantic version",
			versions: []string{"0.9.0", "0.10.0", "0.2.0"},
			expected: "0.10.0",
		},
		{
			name:     "latest v-prefixed semantic version",
			versions: []string{"v1.0.0", "v1.2.0", "v1.1.9"},
			expected: "v1.2.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sourcePath := t.TempDir()
			modelName := "exoscale-icons"
			for _, version := range tt.versions {
				if err := os.MkdirAll(filepath.Join(sourcePath, modelName, version), 0o755); err != nil {
					t.Fatal(err)
				}
			}

			actual, err := resolveModelOCIVersion(sourcePath, modelName, meshkitRegistryUtils.ModelCSV{Model: modelName})
			if err != nil {
				t.Fatal(err)
			}
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestResolveModelOCISourceDirWithRawModelName(t *testing.T) {
	sourcePath := t.TempDir()
	rawModelName := "confidential containers"
	if err := os.MkdirAll(filepath.Join(sourcePath, rawModelName, "0.1.0"), 0o755); err != nil {
		t.Fatal(err)
	}

	modelDir, version, err := resolveModelOCISource(sourcePath, "confidential-containers", meshkitRegistryUtils.ModelCSV{Model: rawModelName})
	if err != nil {
		t.Fatal(err)
	}

	if modelDir != filepath.Join(sourcePath, rawModelName) {
		t.Fatalf("expected model dir %q, got %q", filepath.Join(sourcePath, rawModelName), modelDir)
	}
	if version != "0.1.0" {
		t.Fatalf("expected version %q, got %q", "0.1.0", version)
	}
}

func TestResolveModelOCIVersionMissingDirectory(t *testing.T) {
	_, err := resolveModelOCIVersion(t.TempDir(), "missing-model", meshkitRegistryUtils.ModelCSV{Model: "missing-model"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
