package core

import (
	"testing"
)

// Mock data representing a simplified K8s OpenAPI v3 response
const mockOpenAPIJSON = `
{
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
          "metadata": { "type": "object" },
          "spec": {
            "type": "object",
            "properties": {
              "containers": { "type": "array", "items": { "type": "object" } },
              "nodeName": { "type": "string" }
            }
          },
          "status": { "type": "object" }
        }
      },
      "io.k8s.api.core.v1.Service": {
        "type": "object",
        "x-kubernetes-group-version-kind": [
          {
            "group": "",
            "kind": "Service",
            "version": "v1"
          }
        ],
        "properties": {
          "apiVersion": { "type": "string" },
          "kind": { "type": "string" },
          "spec": {
            "type": "object",
            "properties": { "ports": { "type": "array" } }
          }
        }
      },
      "io.k8s.api.apps.v1.Deployment": {
        "type": "object",
        "x-kubernetes-group-version-kind": [
          {
            "group": "apps",
            "kind": "Deployment",
            "version": "v1"
          }
        ],
        "properties": {
          "spec": { "type": "object" }
        }
      }
    }
  }
}
`

const malformedJSON = `{"components": {"schemas": {`

func TestGetCRDsFromManifest(t *testing.T) {
	tests := []struct {
		name            string
		manifest        string
		apiResources    []string
		expectedCount   int
		expectedKind    string
		expectedVersion string
		expectSchema    bool
		wantError       bool
	}{
		{
			name:            "Valid Pod Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Pod"},
			expectedCount:   1,
			expectedKind:    "Pod", // Fixed: Expect CamelCase (K8s convention)
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:            "Valid Service Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Service"},
			expectedCount:   1,
			expectedKind:    "Service", // Fixed: Expect CamelCase
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:            "Case Sensitivity Check",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"pod"},
			expectedCount:   1,
			expectedKind:    "pod", // Expects input casing
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:          "Multiple Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{"Pod", "Service"},
			expectedCount: 2,
			expectSchema:  true,
			wantError:     false,
		},
		{
			name:          "Filter Unwanted Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{"ConfigMap"},
			expectedCount: 0,
			wantError:     false,
		},
		{
			name:            "Deployment with Group",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Deployment"},
			expectedCount:   1,
			expectedKind:    "Deployment", // Fixed: Expect CamelCase
			expectedVersion: "apps/v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:          "Malformed JSON",
			manifest:      malformedJSON,
			apiResources:  []string{"Pod"},
			expectedCount: 0,
			expectSchema:  false,
			wantError:     true,
		},
		{
			name:          "Empty API Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{},
			expectedCount: 0,
			expectSchema:  false,
			wantError:     false,
		},
		{
			name:          "Empty Manifest",
			manifest:      "",
			apiResources:  []string{"Pod"},
			expectedCount: 0,
			expectSchema:  false,
			wantError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := getCRDsFromManifest(tt.manifest, tt.apiResources)

			if (err != nil) != tt.wantError {
				t.Errorf("getCRDsFromManifest() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
				return
			}

			if tt.expectedCount == 1 {
				result := results[0]

				if result.kind != tt.expectedKind {
					t.Errorf("expected kind %s, got %s", tt.expectedKind, result.kind)
				}

				if result.apiVersion != tt.expectedVersion {
					t.Errorf("expected version %s, got %s", tt.expectedVersion, result.apiVersion)
				}

				if tt.expectSchema && result.schema == "" {
					t.Error("expected non-empty schema, got empty string")
				}
			}

			if tt.name == "Multiple Resources" && len(results) == 2 {
				kinds := make(map[string]bool)
				for _, r := range results {
					kinds[r.kind] = true
					if r.schema == "" {
						t.Errorf("expected schema for resource %s, got empty", r.kind)
					}
				}
				// Updated to check for CamelCase
				if !kinds["Pod"] || !kinds["Service"] {
					t.Error("expected to find both Pod and Service in results")
				}
			}
		})
	}
}
