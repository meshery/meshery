package core

import (
	"strings"
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
		{
			name:         "schemas field is not an object",
			manifest:     `{"components": {"schemas": "invalid-string"}}`,
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

func TestGetCRDsFromManifest_AbsentSchemas(t *testing.T) {
	manifest := `{"openapi": "3.0.0", "info": {"title": "test without schemas"}}`
	res, err := getCRDsFromManifest(manifest, []string{"Pod"})
	if err != nil {
		t.Fatalf("unexpected error for manifest without schemas: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected 0 crds for manifest without schemas, got %d", len(res))
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

func TestExtractCRDsFromManifests_AllSuccess(t *testing.T) {
	validPodManifest := `{
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

	contents := []apiResult{
		{
			path:    "/openapi/v3/api/v1",
			content: []byte(validPodManifest),
		},
	}

	crds, err := extractCRDsFromManifests(contents, []string{"Pod"})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(crds) != 1 {
		t.Fatalf("expected 1 CRD, got %d", len(crds))
	}
	if crds[0].kind != "Pod" {
		t.Errorf("expected kind 'Pod', got '%s'", crds[0].kind)
	}
}

func TestExtractCRDsFromManifests_PartialFailure(t *testing.T) {
	validPodManifest := `{
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

	contents := []apiResult{
		{
			path:    "/openapi/v3/api/v1",
			content: []byte(validPodManifest),
		},
		{
			path:    "/openapi/v3/apis/faulty.k8s.io/v1",
			content: []byte("{invalid-json"),
		},
	}

	crds, err := extractCRDsFromManifests(contents, []string{"Pod"})
	if err == nil {
		t.Fatal("expected error due to malformed manifest, got nil")
	}
	// Verify valid CRD was extracted despite the error
	if len(crds) != 1 {
		t.Fatalf("expected 1 CRD extracted from valid manifest, got %d", len(crds))
	}
	if crds[0].kind != "Pod" {
		t.Errorf("expected kind 'Pod', got '%s'", crds[0].kind)
	}
	// Verify API path retention in the error message
	if !strings.Contains(err.Error(), "/openapi/v3/apis/faulty.k8s.io/v1") {
		t.Errorf("expected error message to retain API path '/openapi/v3/apis/faulty.k8s.io/v1', got: %v", err)
	}
}

func TestExtractCRDsFromManifests_AllFailure(t *testing.T) {
	contents := []apiResult{
		{
			path:    "/openapi/v3/apis/bad1.k8s.io/v1",
			content: []byte("{bad-json-1"),
		},
		{
			path:    "/openapi/v3/apis/bad2.k8s.io/v1",
			content: []byte("{bad-json-2"),
		},
	}

	crds, err := extractCRDsFromManifests(contents, []string{"Pod"})
	if err == nil {
		t.Fatal("expected error when all manifests fail, got nil")
	}
	if len(crds) != 0 {
		t.Fatalf("expected 0 CRDs, got %d", len(crds))
	}
	if !strings.Contains(err.Error(), "/openapi/v3/apis/bad1.k8s.io/v1") {
		t.Errorf("expected error to contain bad1 path, got: %v", err)
	}
	if !strings.Contains(err.Error(), "/openapi/v3/apis/bad2.k8s.io/v1") {
		t.Errorf("expected error to contain bad2 path, got: %v", err)
	}
}
