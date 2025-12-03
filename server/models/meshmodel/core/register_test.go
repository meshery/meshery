package core

import (
	"encoding/json"
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
		wantError       bool // New field to check for errors
	}{
		{
			name:            "Valid Pod Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Pod"},
			expectedCount:   1,
			expectedKind:    "pod", // Fixed: Expect lowercase
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:            "Valid Service Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Service"},
			expectedCount:   1,
			expectedKind:    "service", // Fixed: Expect lowercase
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:            "Case Sensitivity Check",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"pod"},
			expectedCount:   1,
			expectedKind:    "pod",
			expectedVersion: "v1",
			expectSchema:    true,
			wantError:       false,
		},
		{
			name:          "Multiple Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{"Pod", "Service"},
			expectedCount: 2,
			expectSchema:  true, // We will validate this manually in the loop
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
			expectedKind:    "deployment", // Fixed: Expect lowercase
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
			wantError:     true, // We now expect an error returned
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
			wantError:     true, // We now expect an error returned
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Updated to handle 2 return values
			results, err := getCRDsFromManifest(tt.manifest, tt.apiResources)

			// Check if error matches expectation
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

			// Enhanced check for Multiple Resources
			if tt.name == "Multiple Resources" && len(results) == 2 {
				kinds := make(map[string]bool)
				for _, r := range results {
					kinds[r.kind] = true
					// Reviewer requested schema validation
					if r.schema == "" {
						t.Errorf("expected schema for resource %s, got empty", r.kind)
					}
				}
				// Reviewer requested lowercase checks
				if !kinds["pod"] || !kinds["service"] {
					t.Error("expected to find both pod and service in results")
				}
			}
		})
	}
}
