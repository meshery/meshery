package models

import (
	"encoding/json"
	"testing"
)

// TestMeshmodelsAPIResponse_MarshalEmitsBothKeyFlavors locks in the
// deprecation-window contract: MarshalJSON emits BOTH the canonical
// camelCase keys (`pageSize`, `totalCount`) AND the legacy snake_case
// keys (`page_size`, `total_count`) so external consumers on either
// vocabulary keep working while they migrate.
func TestMeshmodelsAPIResponse_MarshalEmitsBothKeyFlavors(t *testing.T) {
	p := MeshmodelsAPIResponse{Page: 1, PageSize: 25, TotalCount: 7}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(b, &top); err != nil {
		t.Fatalf("re-decode: %v (out=%s)", err, string(b))
	}
	for _, want := range []string{"pageSize", "page_size", "totalCount", "total_count", "page", "models"} {
		if _, ok := top[want]; !ok {
			t.Errorf("marshal output missing top-level key %q; got keys %v", want, topLevelKeysMeshmodels(top))
		}
	}
}

// TestMeshmodelsAPIResponse_UnmarshalAcceptsEitherKeyFlavor locks in the
// accept-either-on-input rule: clients sending canonical keys work,
// clients sending legacy keys work, and when both are present canonical
// wins.
func TestMeshmodelsAPIResponse_UnmarshalAcceptsEitherKeyFlavor(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		wantPageSize int
		wantTotal    int64
	}{
		{"canonical only", `{"page":1,"pageSize":25,"totalCount":7}`, 25, 7},
		{"legacy only", `{"page":1,"page_size":25,"total_count":7}`, 25, 7},
		{"canonical wins over legacy when both present", `{"page":1,"pageSize":25,"totalCount":7,"page_size":99,"total_count":999}`, 25, 7},
		{"both absent zeroes the fields", `{"page":1}`, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got MeshmodelsAPIResponse
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.PageSize != tc.wantPageSize {
				t.Errorf("PageSize = %d, want %d", got.PageSize, tc.wantPageSize)
			}
			if got.TotalCount != tc.wantTotal {
				t.Errorf("TotalCount = %d, want %d", got.TotalCount, tc.wantTotal)
			}
		})
	}
}

// TestMeshmodelsAPIResponse_UnmarshalResetsFieldsOnReuse locks in
// stdlib json.Unmarshal semantics: reusing a receiver across decodes,
// if the second input omits PageSize and TotalCount entirely, both
// fields must reset to zero rather than carry stale values forward.
func TestMeshmodelsAPIResponse_UnmarshalResetsFieldsOnReuse(t *testing.T) {
	var p MeshmodelsAPIResponse
	if err := json.Unmarshal([]byte(`{"page":1,"pageSize":50,"totalCount":100}`), &p); err != nil {
		t.Fatalf("prime: %v", err)
	}
	if p.PageSize != 50 || p.TotalCount != 100 {
		t.Fatalf("prime decode wrong: PageSize=%d TotalCount=%d", p.PageSize, p.TotalCount)
	}
	if err := json.Unmarshal([]byte(`{"page":2}`), &p); err != nil {
		t.Fatalf("reuse: %v", err)
	}
	if p.PageSize != 0 {
		t.Errorf("PageSize leaked stale value %d across reuse", p.PageSize)
	}
	if p.TotalCount != 0 {
		t.Errorf("TotalCount leaked stale value %d across reuse", p.TotalCount)
	}
	if p.Page != 2 {
		t.Errorf("Page = %d, want 2", p.Page)
	}
}

// TestMeshmodelsDuplicateAPIResponse_Marshal locks the same dual-emit
// contract on the duplicates variant envelope — they share the wire
// rules and must not diverge in future edits.
func TestMeshmodelsDuplicateAPIResponse_MarshalEmitsBothKeyFlavors(t *testing.T) {
	p := MeshmodelsDuplicateAPIResponse{Page: 1, PageSize: 25, TotalCount: 7}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(b, &top); err != nil {
		t.Fatalf("re-decode: %v", err)
	}
	for _, want := range []string{"pageSize", "page_size", "totalCount", "total_count"} {
		if _, ok := top[want]; !ok {
			t.Errorf("duplicate-variant marshal missing %q; got %v", want, topLevelKeysMeshmodels(top))
		}
	}
}

func topLevelKeysMeshmodels(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
