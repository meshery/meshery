package model

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
)

// captureStdout runs f and returns everything it wrote to os.Stdout. The
// display helpers print with fmt.Println, so this is how a test asserts that a
// malformed relationship renders nothing at all rather than a header or a
// %!s(...) row.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	f()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return buf.String()
}

// buildEntityTypeLine receives names/entityTypes as []interface{} decoded from a
// registry API JSON response. Its sibling builders (buildDescription,
// buildDescriptionList) guard every element with the comma-ok form, but
// buildEntityTypeLine asserted .(string) directly, so a non-string element
// panicked the mesheryctl model import error display. This test asserts it
// degrades gracefully instead of panicking.
func TestBuildEntityTypeLineHandlesNonStringEntries(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name        string
		names       []interface{}
		entityTypes []interface{}
		modelName   string
	}{
		{"non-string name with model filter", []interface{}{123}, []interface{}{"unknown"}, "mymodel"},
		{"non-string entityType", []interface{}{"mymodel"}, []interface{}{42}, "mymodel"},
		{"mixed non-string entries", []interface{}{true, "mymodel"}, []interface{}{nil, "unknown"}, "mymodel"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Must not panic on malformed (non-string) interface elements.
			_ = buildEntityTypeLine(tt.names, tt.entityTypes, "desc", "", "", tt.modelName)
		})
	}
}

// displaySuccessfulRelationships walks SuccessfulRelationships ([]map[string]interface{})
// straight from the registry API JSON. It asserted .(string)/.(map)/.([]interface{})
// directly and indexed from[0]/to[0] without a length check, so a malformed or
// sparse relationship entry panicked mesheryctl. This asserts it skips bad entries
// instead of crashing, matching displaySuccessfulComponents in the same file.
func TestDisplaySuccessfulRelationshipsHandlesMalformed(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name string
		rel  map[string]interface{}
	}{
		{
			name: "empty from/to lists (index out of range)",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{}, "to": []interface{}{},
					}},
				},
			},
		},
		{
			name: "non-string kind with a valid selector",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": 123, "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{map[string]interface{}{"kind": "X"}},
						"to":   []interface{}{map[string]interface{}{"kind": "Y"}},
					}},
				},
			},
		},
		{
			name: "selector is not a map",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{"not-a-map"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{tt.rel},
				},
			}
			// Must not panic on malformed relationship data.
			displaySuccessfulRelationships(resp, "mymodel")
		})
	}
}

// A relationship whose display fields (Kind/Subtype/RelationshipType) or endpoint
// kinds are not strings is skipped entirely: it must not render a RELATIONSHIP
// header or a From/To row, and must never emit a %!s(...) verb. Only well-formed
// relationships produce output.
func TestDisplaySuccessfulRelationshipsSkipsMalformedOutput(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	malformed := []struct {
		name string
		rel  map[string]interface{}
	}{
		{
			name: "non-string endpoint kind",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{map[string]interface{}{"kind": 7}},
						"to":   []interface{}{map[string]interface{}{"kind": "Y"}},
					}},
				},
			},
		},
		{
			name: "non-string RelationshipType",
			rel: map[string]interface{}{
				"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": 9,
				"Selectors": []interface{}{
					map[string]interface{}{"allow": map[string]interface{}{
						"from": []interface{}{map[string]interface{}{"kind": "X"}},
						"to":   []interface{}{map[string]interface{}{"kind": "Y"}},
					}},
				},
			},
		},
	}
	for _, tt := range malformed {
		t.Run(tt.name, func(t *testing.T) {
			resp := &models.RegistryAPIResponse{
				EntityTypeSummary: models.EntityTypeSummary{
					SuccessfulRelationships: []map[string]interface{}{tt.rel},
				},
			}
			out := captureStdout(t, func() {
				displaySuccessfulRelationships(resp, "mymodel")
			})
			if out != "" {
				t.Errorf("expected no output for malformed relationship, got:\n%s", out)
			}
		})
	}

	// Positive control: a well-formed relationship still renders.
	t.Run("well-formed relationship renders", func(t *testing.T) {
		resp := &models.RegistryAPIResponse{
			EntityTypeSummary: models.EntityTypeSummary{
				SuccessfulRelationships: []map[string]interface{}{{
					"Model": "mymodel", "Kind": "a", "Subtype": "b", "RelationshipType": "c",
					"Selectors": []interface{}{
						map[string]interface{}{"allow": map[string]interface{}{
							"from": []interface{}{map[string]interface{}{"kind": "X"}},
							"to":   []interface{}{map[string]interface{}{"kind": "Y"}},
						}},
					},
				}},
			},
		}
		out := captureStdout(t, func() {
			displaySuccessfulRelationships(resp, "mymodel")
		})
		// Assert the actual rendered content, not just that something printed:
		// displaySuccessfulRelationships emits a leading `fmt.Println("")`, so a
		// non-empty check alone would pass even if the table never rendered.
		// The table headers render uppercase ("FROM"/"TO"); X and Y are the
		// endpoint kinds from the selector.
		for _, want := range []string{"FROM", "TO", "X", "Y"} {
			if !strings.Contains(out, want) {
				t.Errorf("expected rendered output to contain %q, got:\n%s", want, out)
			}
		}
	})
}
