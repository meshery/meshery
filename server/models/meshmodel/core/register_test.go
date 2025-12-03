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
	}{
		{
			name:            "Valid Pod Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Pod"},
			expectedCount:   1,
			expectedKind:    "Pod",
			expectedVersion: "v1",
			expectSchema:    true,
		},
		{
			name:            "Valid Service Parsing",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Service"},
			expectedCount:   1,
			expectedKind:    "Service",
			expectedVersion: "v1",
			expectSchema:    true,
		},
		{
			name:            "Case Sensitivity Check",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"pod"}, // Lowercase input
			expectedCount:   1,
			expectedKind:    "pod",
			expectedVersion: "v1",
			expectSchema:    true,
		},
		{
			name:          "Multiple Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{"Pod", "Service"},
			expectedCount: 2,
			expectSchema:  false,
		},
		{
			name:          "Filter Unwanted Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{"ConfigMap"},
			expectedCount: 0,
			expectSchema:  false,
		},
		{
			name:            "Deployment with Group",
			manifest:        mockOpenAPIJSON,
			apiResources:    []string{"Deployment"},
			expectedCount:   1,
			expectedKind:    "Deployment",
			expectedVersion: "apps/v1",
			expectSchema:    true,
		},
		{
			name:          "Malformed JSON",
			manifest:      malformedJSON,
			apiResources:  []string{"Pod"},
			expectedCount: 0,
			expectSchema:  false,
		},
		{
			name:          "Empty API Resources",
			manifest:      mockOpenAPIJSON,
			apiResources:  []string{},
			expectedCount: 0,
			expectSchema:  false,
		},
		{
			name:          "Empty Manifest",
			manifest:      "",
			apiResources:  []string{"Pod"},
			expectedCount: 0,
			expectSchema:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := getCRDsFromManifest(tt.manifest, tt.apiResources)

			if len(results) != tt.expectedCount {
				t.Errorf("expected %d results, got %d", tt.expectedCount, len(results))
				return
			}

			// If we expect a specific single result, verify its details
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

				if tt.expectSchema && result.name == "" {
					t.Error("expected non-empty name, got empty string")
				}

				// Verify schema is valid JSON
				if tt.expectSchema {
					var schemaCheck map[string]interface{}
					if err := json.Unmarshal([]byte(result.schema), &schemaCheck); err != nil {
						t.Errorf("schema is not valid JSON: %v", err)
					}
				}
			}

			// Special check for Multiple Resources case
			if tt.name == "Multiple Resources" && len(results) == 2 {
				kinds := make(map[string]bool)
				for _, r := range results {
					kinds[r.kind] = true
				}
				if !kinds["Pod"] || !kinds["Service"] {
					t.Error("expected to find both Pod and Service in results")
				}
			}
		})
	}
}
