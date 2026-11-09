package stages

import (
	"encoding/json"
	"strings"
	"testing"

	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

// Test_matchPattern verifies that matchPattern correctly detects valid reference syntax and extracts the query token.
func Test_matchPattern(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		wantKey   string
		wantMatch bool
	}{
		{
			name:      "Exact pattern matches",
			str:       "$(#ref.services.test)",
			wantKey:   "services.test",
			wantMatch: true,
		},
		{
			name:      "Pattern without closing parenthesis fails",
			str:       "$(#ref.services.test",
			wantKey:   "",
			wantMatch: false,
		},
		{
			name:      "String without pattern fails",
			str:       "plain-unreferenced-string",
			wantKey:   "",
			wantMatch: false,
		},
		{
			name:      "Pattern embedded with surrounding characters extracts key",
			str:       "prefix_$(#ref.components.0.name)_suffix",
			wantKey:   "components.0.name",
			wantMatch: true,
		},
		{
			name:      "Multiple patterns extracts first key",
			str:       "$(#ref.first.key):$(#ref.second.key)",
			wantKey:   "first.key",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKey, gotMatch := matchPattern(tt.str)
			if gotKey != tt.wantKey {
				t.Errorf("matchPattern() gotKey = %v, want %v", gotKey, tt.wantKey)
			}
			if gotMatch != tt.wantMatch {
				t.Errorf("matchPattern() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

// Test_fillMapString tests composite URLs, multi-reference strings, repeated references, scalar formatting, error handling, and non-reprocessing.
func Test_fillMapString(t *testing.T) {
	flattened := map[string]interface{}{
		"services.api.name":     "api-gateway",
		"services.api.host":     "10.0.0.1",
		"services.api.port":       float64(8080),
		"services.api.replicas":   3,
		"services.api.enabled":    true,
		"services.api.rawPattern": "$(#ref.services.other.ignored)",
	}

	tests := []struct {
		name        string
		input       string
		wantResult  string
		wantOk      bool
		wantErrSub  string
	}{
		{
			name:       "No reference present",
			input:      "regular string without references",
			wantResult: "",
			wantOk:     false,
		},
		{
			name:       "Exact single reference",
			input:      "$(#ref.services.api.name)",
			wantResult: "api-gateway",
			wantOk:     true,
		},
		{
			name:       "Composite string preserves prefix and suffix",
			input:      "http://$(#ref.services.api.name):8080/healthz",
			wantResult: "http://api-gateway:8080/healthz",
			wantOk:     true,
		},
		{
			name:       "Multiple references on single line",
			input:      "$(#ref.services.api.host):$(#ref.services.api.port)",
			wantResult: "10.0.0.1:8080",
			wantOk:     true,
		},
		{
			name:       "Repeated same reference",
			input:      "/api/$(#ref.services.api.name)/v1/$(#ref.services.api.name)",
			wantResult: "/api/api-gateway/v1/api-gateway",
			wantOk:     true,
		},
		{
			name:       "Numeric and boolean scalars formatted cleanly",
			input:      "replicas=$(#ref.services.api.replicas)&enabled=$(#ref.services.api.enabled)",
			wantResult: "replicas=3&enabled=true",
			wantOk:     true,
		},
		{
			name:       "Resolved value containing reference pattern is not reprocessed",
			input:      "$(#ref.services.api.rawPattern)",
			wantResult: "$(#ref.services.other.ignored)",
			wantOk:     true,
		},
		{
			name:       "Invalid reference returns error",
			input:      "http://$(#ref.services.unknown.host):8080",
			wantErrSub: "invalid reference query: services.unknown.host",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotOk, err := fillMapString(tt.input, flattened)
			if tt.wantErrSub != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErrSub) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErrSub, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotOk != tt.wantOk {
				t.Errorf("fillMapString() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotResult != tt.wantResult {
				t.Errorf("fillMapString() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

// Test_fillMap verifies that fillMap recursively traverses nested maps and multi-dimensional slices.
func Test_fillMap(t *testing.T) {
	flattened := map[string]interface{}{
		"host": "localhost",
		"port": float64(9090),
		"db":   "mydb",
	}

	input := map[string]interface{}{
		"url": "http://$(#ref.host):$(#ref.port)/query",
		"nested": map[string]interface{}{
			"conn": "$(#ref.db)",
		},
		"list": []interface{}{
			"direct-$(#ref.host)",
			map[string]interface{}{
				"port": "$(#ref.port)",
			},
			[]interface{}{
				"nested-slice-$(#ref.db)",
			},
		},
	}

	filled, err := fillMap(input, flattened)
	if err != nil {
		t.Fatalf("fillMap failed: %v", err)
	}

	if filled["url"] != "http://localhost:9090/query" {
		t.Errorf("expected composite url, got %v", filled["url"])
	}

	nested, ok := filled["nested"].(map[string]interface{})
	if !ok || nested["conn"] != "mydb" {
		t.Errorf("expected nested conn 'mydb', got %v", nested)
	}

	list, ok := filled["list"].([]interface{})
	if !ok || len(list) != 3 {
		t.Fatalf("expected list of length 3, got %v", list)
	}

	if list[0] != "direct-localhost" {
		t.Errorf("expected list[0] 'direct-localhost', got %v", list[0])
	}

	innerMap, ok := list[1].(map[string]interface{})
	if !ok || innerMap["port"] != "9090" {
		t.Errorf("expected innerMap port '9090', got %v", innerMap)
	}

	innerSlice, ok := list[2].([]interface{})
	if !ok || len(innerSlice) != 1 || innerSlice[0] != "nested-slice-mydb" {
		t.Errorf("expected innerSlice[0] 'nested-slice-mydb', got %v", innerSlice)
	}
}

// TestFiller_StageEndToEnd verifies full end-to-end execution of the Filler stage on an unmarshaled design pattern.
func TestFiller_StageEndToEnd(t *testing.T) {
	rawDesign := `{
		"id": "11111111-2222-3333-4444-555555555555",
		"name": "filler-test-design",
		"schemaVersion": "designs.meshery.io/v1beta1",
		"version": "1.0.0",
		"components": [
			{
				"id": "11111111-1111-1111-1111-111111111111",
				"displayName": "database",
				"component": {
					"kind": "StatefulSet",
					"version": "v1"
				},
				"model": {
					"name": "kubernetes",
					"model": {
						"version": "v1.25.0"
					}
				},
				"configuration": {
					"metadata": {
						"namespace": "production"
					},
					"spec": {
						"serviceName": "db-service",
						"port": 5432
					}
				},
				"metadata": {}
			},
			{
				"id": "22222222-2222-2222-2222-222222222222",
				"displayName": "backend",
				"component": {
					"kind": "Deployment",
					"version": "$(#ref.components.0.component.version)"
				},
				"model": {
					"name": "kubernetes",
					"model": {
						"version": "$(#ref.components.0.model.model.version)"
					}
				},
				"configuration": {
					"metadata": {
						"namespace": "$(#ref.components.0.configuration.metadata.namespace)"
					},
					"spec": {
						"databaseUrl": "postgres://$(#ref.components.0.configuration.spec.serviceName):$(#ref.components.0.configuration.spec.port)/prod_db",
						"endpoints": [
							"http://$(#ref.components.0.configuration.spec.serviceName):$(#ref.components.0.configuration.spec.port)/status"
						]
					}
				},
				"metadata": {
					"dependsOn": [
						"$(#ref.components.0.displayName)"
					]
				}
			}
		]
	}`

	var p pattern.PatternFile
	if err := json.Unmarshal([]byte(rawDesign), &p); err != nil {
		t.Fatalf("failed to unmarshal test design: %v", err)
	}

	data := &Data{
		Pattern: &p,
	}

	var stageErr error
	Filler(true)(data, nil, func(d *Data, err error) {
		stageErr = err
	})

	if stageErr != nil {
		t.Fatalf("Filler stage returned error: %v", stageErr)
	}

	c2 := p.Components[1]

	// Check version reference resolved
	if c2.Component.Version != "v1" {
		t.Errorf("expected component version 'v1', got %q", c2.Component.Version)
	}
	if c2.Model.Model.Version != "v1.25.0" {
		t.Errorf("expected model version 'v1.25.0', got %q", c2.Model.Model.Version)
	}

	// Check namespace reference resolved
	meta, ok := c2.Configuration["metadata"].(map[string]interface{})
	if !ok || meta["namespace"] != "production" {
		t.Errorf("expected namespace 'production', got %v", meta["namespace"])
	}

	// Check composite databaseUrl preserved prefixes, suffixes, and multiple references
	spec, ok := c2.Configuration["spec"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected spec map in configuration, got %v", c2.Configuration["spec"])
	}

	expectedDBURL := "postgres://db-service:5432/prod_db"
	if spec["databaseUrl"] != expectedDBURL {
		t.Errorf("expected databaseUrl %q, got %v", expectedDBURL, spec["databaseUrl"])
	}

	endpoints, ok := spec["endpoints"].([]interface{})
	if !ok || len(endpoints) != 1 {
		t.Fatalf("expected endpoints slice with 1 element, got %v", spec["endpoints"])
	}
	expectedEndpoint := "http://db-service:5432/status"
	if endpoints[0] != expectedEndpoint {
		t.Errorf("expected endpoints[0] %q, got %v", expectedEndpoint, endpoints[0])
	}

	// Check dependsOn resolved
	dependsOn, ok := c2.Metadata.AdditionalProperties["dependsOn"].([]string)
	if !ok || len(dependsOn) != 1 || dependsOn[0] != "database" {
		t.Errorf("expected dependsOn ['database'], got %v", c2.Metadata.AdditionalProperties["dependsOn"])
	}
}
