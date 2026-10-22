package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestOrganizationsPage_MarshalEmitsBothKeyFlavors locks in the
// deprecation-window contract: MarshalJSON emits both the canonical
// camelCase (`totalCount`, `pageSize`) AND the legacy snake_case
// (`total_count`, `page_size`) spellings so external consumers on
// either vocabulary keep working while they migrate.
func TestOrganizationsPage_MarshalEmitsBothKeyFlavors(t *testing.T) {
	p := OrganizationsPage{
		TotalCount: 7,
		Page:       2,
		PageSize:   25,
	}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	out := string(b)
	for _, want := range []string{`"totalCount":7`, `"total_count":7`, `"pageSize":25`, `"page_size":25`, `"page":2`} {
		if !strings.Contains(out, want) {
			t.Errorf("marshal output missing %q; got %s", want, out)
		}
	}
}

// TestOrganizationsPage_UnmarshalAcceptsEitherKeyFlavor verifies that
// either wire form round-trips back into the struct. Canonical wins if
// both are present.
func TestOrganizationsPage_UnmarshalAcceptsEitherKeyFlavor(t *testing.T) {
	cases := []struct {
		name        string
		body        string
		wantTotal   int
		wantPageSz  uint64
	}{
		{"canonical only", `{"totalCount":5,"page":1,"pageSize":10}`, 5, 10},
		{"legacy only", `{"total_count":5,"page":1,"page_size":10}`, 5, 10},
		{"canonical wins when both present", `{"totalCount":9,"total_count":1,"page":1,"pageSize":50,"page_size":5}`, 9, 50},
		{"both flavors absent zeroes the fields", `{"page":1}`, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got OrganizationsPage
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.TotalCount != tc.wantTotal {
				t.Errorf("TotalCount = %d, want %d", got.TotalCount, tc.wantTotal)
			}
			if got.PageSize != tc.wantPageSz {
				t.Errorf("PageSize = %d, want %d", got.PageSize, tc.wantPageSz)
			}
		})
	}
}

// TestOrganizationsPage_UnmarshalResetsFieldsOnReuse locks in std-library
// json.Unmarshal semantics: when the destination struct is reused across
// decodes and the second input omits TotalCount / PageSize entirely, the
// destination's fields must reset to zero rather than carry stale values
// forward from the prior decode.
func TestOrganizationsPage_UnmarshalResetsFieldsOnReuse(t *testing.T) {
	var p OrganizationsPage
	if err := json.Unmarshal([]byte(`{"totalCount":99,"pageSize":77,"page":1}`), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.TotalCount != 99 || p.PageSize != 77 {
		t.Fatalf("prime decode wrong: got TotalCount=%d PageSize=%d", p.TotalCount, p.PageSize)
	}
	if err := json.Unmarshal([]byte(`{"page":2}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.TotalCount != 0 {
		t.Errorf("TotalCount leaked stale value %d across reuse", p.TotalCount)
	}
	if p.PageSize != 0 {
		t.Errorf("PageSize leaked stale value %d across reuse", p.PageSize)
	}
	if p.Page != 2 {
		t.Errorf("Page = %d, want 2", p.Page)
	}
}
