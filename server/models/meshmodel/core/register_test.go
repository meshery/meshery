package core

import (
	"encoding/json"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func TestGetResolvedManifest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		// checkPath is a dot-separated path into the parsed output to verify
		// that $ref was resolved. Empty means just check valid JSON output.
		checkPath string
		wantType  string
	}{
		{
			name: "No refs passes through unchanged",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {},
				"components": {
					"schemas": {
						"Foo": {
							"type": "object",
							"properties": {
								"name": {"type": "string"}
							}
						}
					}
				}
			}`,
			checkPath: "components.schemas.Foo.properties.name",
			wantType:  "string",
		},
		{
			name: "Property $ref is inlined",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {},
				"components": {
					"schemas": {
						"Address": {
							"type": "object",
							"properties": {
								"street": {"type": "string"}
							}
						},
						"Person": {
							"type": "object",
							"properties": {
								"name": {"type": "string"},
								"address": {"$ref": "#/components/schemas/Address"}
							}
						}
					}
				}
			}`,
			checkPath: "components.schemas.Person.properties.address",
			wantType:  "object",
		},
		{
			name: "Nested $ref chain is fully inlined",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {},
				"components": {
					"schemas": {
						"Zip": {"type": "string"},
						"Address": {
							"type": "object",
							"properties": {
								"zip": {"$ref": "#/components/schemas/Zip"}
							}
						},
						"Person": {
							"type": "object",
							"properties": {
								"address": {"$ref": "#/components/schemas/Address"}
							}
						}
					}
				}
			}`,
			checkPath: "components.schemas.Person.properties.address.properties.zip",
			wantType:  "string",
		},
		{
			name: "Array items $ref is inlined",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {},
				"components": {
					"schemas": {
						"Tag": {
							"type": "object",
							"properties": {"label": {"type": "string"}}
						},
						"TagList": {
							"type": "array",
							"items": {"$ref": "#/components/schemas/Tag"}
						}
					}
				}
			}`,
			checkPath: "components.schemas.TagList.items",
			wantType:  "object",
		},
		{
			name: "additionalProperties $ref is inlined",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {},
				"components": {
					"schemas": {
						"Value": {"type": "string"},
						"Map": {
							"type": "object",
							"additionalProperties": {"$ref": "#/components/schemas/Value"}
						}
					}
				}
			}`,
			checkPath: "components.schemas.Map.additionalProperties",
			wantType:  "string",
		},
		{
			name: "No components does not error",
			input: `{
				"openapi": "3.0.0",
				"info": {"title": "test", "version": "1.0"},
				"paths": {}
			}`,
		},
		{
			name:    "Invalid JSON returns error",
			input:   "not json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := getResolvedManifest(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var parsed map[string]any
			if err := json.Unmarshal([]byte(out), &parsed); err != nil {
				t.Fatalf("output is not valid JSON: %v", err)
			}

			if tt.checkPath == "" {
				return
			}

			node := navigatePath(t, parsed, tt.checkPath)
			m, ok := node.(map[string]any)
			if !ok {
				t.Fatalf("expected map at path %s, got %T", tt.checkPath, node)
			}
			if _, hasRef := m["$ref"]; hasRef {
				t.Errorf("$ref still present at path %s", tt.checkPath)
			}
			if tt.wantType != "" {
				if got := m["type"]; got != tt.wantType {
					t.Errorf("type at %s = %v, want %s", tt.checkPath, got, tt.wantType)
				}
			}
		})
	}
}

func TestGetResolvedManifest_AllOf(t *testing.T) {
	input := `{
		"openapi": "3.0.0",
		"info": {"title": "test", "version": "1.0"},
		"paths": {},
		"components": {
			"schemas": {
				"Base": {
					"type": "object",
					"properties": {"id": {"type": "integer"}}
				},
				"Extended": {
					"allOf": [
						{"$ref": "#/components/schemas/Base"},
						{"type": "object", "properties": {"name": {"type": "string"}}}
					]
				}
			}
		}
	}`

	out, err := getResolvedManifest(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	json.Unmarshal([]byte(out), &parsed)

	extended := navigatePath(t, parsed, "components.schemas.Extended").(map[string]any)
	allOf, ok := extended["allOf"].([]any)
	if !ok {
		t.Fatal("expected allOf to be an array")
	}

	for i, item := range allOf {
		m := item.(map[string]any)
		if _, hasRef := m["$ref"]; hasRef {
			t.Errorf("allOf[%d] still has $ref", i)
		}
	}

	base := allOf[0].(map[string]any)
	if base["type"] != "object" {
		t.Errorf("allOf[0] type = %v, want object", base["type"])
	}
}

func TestClearSchemaRefs(t *testing.T) {
	tests := []struct {
		name string
		sr   *openapi3.SchemaRef
	}{
		{
			name: "Nil SchemaRef does not panic",
			sr:   nil,
		},
		{
			name: "Nil Value clears Ref",
			sr:   &openapi3.SchemaRef{Ref: "#/components/schemas/Foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visited := make(map[*openapi3.SchemaRef]bool)
			clearSchemaRefs(tt.sr, visited)
			if tt.sr != nil && tt.sr.Ref != "" {
				t.Errorf("Ref = %q, want empty", tt.sr.Ref)
			}
		})
	}
}

func TestClearSchemaRefs_Circular(t *testing.T) {
	// Build a circular reference: A -> B -> A
	a := &openapi3.SchemaRef{Ref: "#/components/schemas/A", Value: &openapi3.Schema{}}
	b := &openapi3.SchemaRef{Ref: "#/components/schemas/B", Value: &openapi3.Schema{}}
	a.Value.Properties = openapi3.Schemas{"b": b}
	b.Value.Properties = openapi3.Schemas{"a": a}

	visited := make(map[*openapi3.SchemaRef]bool)
	clearSchemaRefs(a, visited) // must not hang or panic

	if a.Ref != "" {
		t.Errorf("a.Ref = %q, want empty", a.Ref)
	}
	if b.Ref != "" {
		t.Errorf("b.Ref = %q, want empty", b.Ref)
	}
}

// navigatePath walks a dot-separated path through nested maps.
func navigatePath(t *testing.T, data map[string]any, path string) any {
	t.Helper()
	keys := splitDot(path)
	var current any = data
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			t.Fatalf("expected map at key %q in path %s, got %T", key, path, current)
		}
		current, ok = m[key]
		if !ok {
			t.Fatalf("key %q not found in path %s", key, path)
		}
	}
	return current
}

func splitDot(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	parts = append(parts, s[start:])
	return parts
}
