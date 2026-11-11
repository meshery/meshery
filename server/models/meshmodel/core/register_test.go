package core

import (
	"testing"
)

func TestGetResolvedManifest_InvalidOrEmpty(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
	}{
		{
			name:     "empty string",
			manifest: "",
		},
		{
			name:     "malformed json",
			manifest: "not-json",
		},
		{
			name:     "unclosed braces",
			manifest: `{"components":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolved, err := getResolvedManifest(tt.manifest)
			if err == nil {
				t.Errorf("expected error for %s, got nil (resolved: %s)", tt.name, resolved)
			}
		})
	}
}

func TestGetCRDsFromManifest_InvalidOrEmpty(t *testing.T) {
	tests := []struct {
		name         string
		manifest     string
		apiResources []string
	}{
		{
			name:         "empty string",
			manifest:     "",
			apiResources: []string{"Pod"},
		},
		{
			name:         "malformed json",
			manifest:     "{broken",
			apiResources: []string{"Pod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := getCRDsFromManifest(tt.manifest, tt.apiResources)
			if len(res) != 0 {
				t.Errorf("expected 0 crds for %s, got %d", tt.name, len(res))
			}
		})
	}
}

func TestGetResolvedManifest_ValidManifest(t *testing.T) {
	manifest := `{
		"components": {
			"schemas": {
				"io.k8s.api.core.v1.Pod": {
					"type": "object",
					"properties": {
						"apiVersion": { "type": "string" },
						"kind": { "type": "string" }
					}
				}
			}
		}
	}`

	resolved, err := getResolvedManifest(manifest)
	if err != nil {
		t.Fatalf("unexpected error for valid manifest: %v", err)
	}
	if resolved == "" {
		t.Error("expected non-empty resolved manifest")
	}
}
