package design

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/stretchr/testify/assert"
)

// captureActionOutput returns what print writes to stdout.
func captureActionOutput(t *testing.T, print func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	print()

	_ = w.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatal(err)
	}
	_ = r.Close()
	return buf.String()
}

// TestPrintDeleteActionsShowRemovedNames verifies that design evaluate names
// the removed component or relationship, which the server sends as a
// definition object, and falls back to the ID when there is no definition.
func TestPrintDeleteActionsShowRemovedNames(t *testing.T) {
	tests := []struct {
		name   string
		print  func(pattern.Action)
		action pattern.Action
		want   string
	}{
		{
			name:  "component with definition",
			print: printDeleteComponentAction,
			action: pattern.Action{Op: "delete_component", Value: map[string]interface{}{
				"id": "c1",
				"component": map[string]interface{}{
					"displayName": "my-nginx",
					"component":   map[string]interface{}{"kind": "Deployment"},
				},
			}},
			want: "  - Deployment \"my-nginx\" (id: c1)\n",
		},
		{
			name:  "component without definition",
			print: printDeleteComponentAction,
			action: pattern.Action{Op: "delete_component", Value: map[string]interface{}{
				"id": "c1",
			}},
			want: "  - component (id: c1)\n",
		},
		{
			name:  "relationship with definition",
			print: printDeleteRelationshipAction,
			action: pattern.Action{Op: "delete_relationship", Value: map[string]interface{}{
				"id": "r1",
				"relationship": map[string]interface{}{
					"kind":    "edge",
					"subType": "network",
					"selectors": []interface{}{
						map[string]interface{}{
							"allow": map[string]interface{}{
								"from": []interface{}{map[string]interface{}{"kind": "Service"}},
								"to":   []interface{}{map[string]interface{}{"kind": "Deployment"}},
							},
						},
					},
				},
			}},
			want: "  - edge-network relationship from Service to Deployment (id: r1)\n",
		},
		{
			name:  "relationship without definition",
			print: printDeleteRelationshipAction,
			action: pattern.Action{Op: "delete_relationship", Value: map[string]interface{}{
				"id": "r1",
			}},
			want: "  - relationship (id: r1)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := captureActionOutput(t, func() { tt.print(tt.action) })
			assert.Equal(t, tt.want, got)
		})
	}
}
