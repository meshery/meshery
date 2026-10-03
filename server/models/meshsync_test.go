package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// The entries inside `kinds` once serialized under their Go field names — the
// anonymous struct they lived in carried no json tags, so a client reading
// `kind` got an empty string and no error (#21715). Pin the wire spelling so a
// dropped tag fails here instead of silently emptying a consumer.
func TestMeshSyncResourcesSummaryKindsAreCamelCase(t *testing.T) {
	payload, err := json.Marshal(MeshSyncResourcesSummaryAPIResponse{
		Kinds: []MeshSyncResourceKindSummary{{
			Kind:  "Pod",
			Model: "kubernetes",
			Count: 5,
		}},
		Namespaces: []string{"default"},
	})
	if err != nil {
		t.Fatalf("marshalling summary response: %v", err)
	}

	got := string(payload)

	for _, want := range []string{`"kinds"`, `"kind":"Pod"`, `"model":"kubernetes"`, `"count":5`, `"namespaces"`, `"labels"`} {
		if !strings.Contains(got, want) {
			t.Errorf("summary response missing %s\ngot: %s", want, got)
		}
	}

	for _, unwanted := range []string{`"Kind"`, `"Model"`, `"Count"`} {
		if strings.Contains(got, unwanted) {
			t.Errorf("summary response still emits PascalCase %s\ngot: %s", unwanted, got)
		}
	}
}

// Each summary label is exactly {key, value} (MeshSyncLabel in meshery/schemas
// v1beta1/meshsync). `value` must be emitted even when empty: node-role labels
// such as node-role.kubernetes.io/control-plane carry an empty value, and the
// meshsync KubernetesKeyValue these used to be encoded as dropped it.
func TestMeshSyncResourcesSummaryLabelsAreKeyValuePairs(t *testing.T) {
	payload, err := json.Marshal(MeshSyncResourcesSummaryAPIResponse{
		Labels: []MeshSyncLabel{{Key: "node-role.kubernetes.io/control-plane", Value: ""}},
	})
	if err != nil {
		t.Fatalf("marshalling summary response: %v", err)
	}

	got := string(payload)
	want := `"labels":[{"key":"node-role.kubernetes.io/control-plane","value":""}]`
	if !strings.Contains(got, want) {
		t.Errorf("summary labels are not {key, value} pairs\nwant substring: %s\ngot: %s", want, got)
	}
}
