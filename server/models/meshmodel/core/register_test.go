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
			res, err := getCRDsFromManifest(tt.manifest, tt.apiResources)
			if err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}
			if len(res) != 0 {
				t.Errorf("expected 0 crds for %s, got %d", tt.name, len(res))
			}
		})
	}
}

func TestGetCRDsFromManifest_ValidManifest(t *testing.T) {
	manifest := `{
		"components": {
			"schemas": {
				"io.k8s.api.core.v1.Pod": {
					"type": "object",
					"x-kubernetes-group-version-kind": [
						{
							"group": "",
							"kind": "Pod",
							"version": "v1"
						}
					],
					"properties": {
						"apiVersion": { "type": "string" },
						"kind": { "type": "string" },
						"spec": { "type": "object" }
					}
				}
			}
		}
	}`

	crds, err := getCRDsFromManifest(manifest, []string{"Pod"})
	if err != nil {
		t.Fatalf("unexpected error for valid manifest: %v", err)
	}
	if len(crds) != 1 {
		t.Fatalf("expected 1 CRD, got %d", len(crds))
	}
	if crds[0].kind != "Pod" {
		t.Errorf("expected kind 'Pod', got '%s'", crds[0].kind)
	}
	if crds[0].apiVersion != "v1" {
		t.Errorf("expected apiVersion 'v1', got '%s'", crds[0].apiVersion)
	}
	if crds[0].schema == "" {
		t.Error("expected non-empty schema")
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
