package modeloci

import "testing"

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
	actual := CompileFolderName("/tmp/models", "exoscale-icons", "0.1.0")
	expected := "/tmp/models/exoscale-icons/0.1.0"
	if actual != expected {
		t.Fatalf("expected %q, got %q", expected, actual)
	}
}
