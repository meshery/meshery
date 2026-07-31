package modeloci

import (
	"path/filepath"
	"testing"
)

func TestCompileImageName(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		extension string
		expected  string
	}{
		{
			name:      "exoscale-icons",
			extension: "tar",
			expected:  "exoscale-icons.tar",
		},
		{
			name:      "exoscale-icons",
			version:   "0.1.0",
			extension: "tar",
			expected:  "exoscale-icons-0-1-0.tar",
		},
		{
			name:      "exoscale-icons",
			version:   "v1.0.0",
			extension: "tar",
			expected:  "exoscale-icons-v1-0-0.tar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			actual := CompileImageName(tt.name, tt.version, tt.extension)
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestCompileFolderName(t *testing.T) {
	base := filepath.Join(string(filepath.Separator), "tmp", "models")
	actual := CompileFolderName(base, "exoscale-icons", "0.1.0")
	expected := filepath.Join(base, "exoscale-icons", "0.1.0")
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}

func TestValidateArtifactIdentifier(t *testing.T) {
	tests := []struct {
		name        string
		modelName   string
		version     string
		expectError bool
	}{
		{name: "valid name and version", modelName: "exoscale-icons", version: "0.1.0"},
		{name: "valid name without version", modelName: "exoscale-icons"},
		{name: "empty model name", modelName: "", expectError: true},
		{name: "dot model name", modelName: ".", expectError: true},
		{name: "parent model name", modelName: "..", expectError: true},
		{name: "slash in model name", modelName: "exoscale/icons", expectError: true},
		{name: "backslash in model name", modelName: `exoscale\icons`, expectError: true},
		{name: "dot version", modelName: "exoscale-icons", version: ".", expectError: true},
		{name: "parent version", modelName: "exoscale-icons", version: "..", expectError: true},
		{name: "slash in version", modelName: "exoscale-icons", version: "0.1/0", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateArtifactIdentifier(tt.modelName, tt.version)
			if tt.expectError && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}
