package models

import (
	"encoding/json"
	"errors"
	"testing"
)

// buildFilterRequestBody returns a JSON body wrapping a minimal
// MesheryFilterPayload (carrying a name) under the provided wrapper
// key. The inner name lets tests assert which spelling "won" when
// multiple spellings are present.
func buildFilterRequestBody(key, filterName string) string {
	body := map[string]any{
		"url":  "",
		"save": false,
	}
	if key != "" {
		body[key] = map[string]any{
			"id":   "00000000-0000-0000-0000-000000000000",
			"name": filterName,
		}
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// buildFilterRequestBodyMulti returns a JSON body with multiple
// wrapper spellings present, each carrying a distinct inner name so
// precedence tests can assert which spelling won.
func buildFilterRequestBodyMulti(entries map[string]string) string {
	body := map[string]any{"url": ""}
	for key, innerName := range entries {
		body[key] = map[string]any{
			"id":   "00000000-0000-0000-0000-000000000000",
			"name": innerName,
		}
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// TestMesheryFilterRequestBody_UnmarshalAcceptsBothSpellings locks in
// the deprecation-window contract for POST /api/filter: the handler
// accepts both the canonical `filterData` (camelCase) and the legacy
// `filter_data` (snake_case) wrapper keys.
func TestMesheryFilterRequestBody_UnmarshalAcceptsBothSpellings(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantName string
	}{
		{
			name:     "canonical filterData only",
			body:     buildFilterRequestBody("filterData", "from-filterData"),
			wantName: "from-filterData",
		},
		{
			name:     "legacy filter_data only",
			body:     buildFilterRequestBody("filter_data", "from-filter_data"),
			wantName: "from-filter_data",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got MesheryFilterRequestBody
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.FilterData == nil {
				t.Fatalf("FilterData is nil, want populated")
			}
			if got.FilterData.Name != tc.wantName {
				t.Errorf("FilterData.Name = %q, want %q", got.FilterData.Name, tc.wantName)
			}
		})
	}
}

// TestMesheryFilterRequestBody_UnmarshalPrecedenceCanonicalWins locks
// the canonical-wins rule: when both spellings are present on
// FilterData, the camelCase spelling wins. Migrating clients may
// temporarily emit both.
func TestMesheryFilterRequestBody_UnmarshalPrecedenceCanonicalWins(t *testing.T) {
	body := buildFilterRequestBodyMulti(map[string]string{
		"filterData":  "canonical",
		"filter_data": "legacy",
	})
	var got MesheryFilterRequestBody
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.FilterData == nil {
		t.Fatalf("FilterData is nil, want canonical")
	}
	if got.FilterData.Name != "canonical" {
		t.Errorf("FilterData.Name = %q, want %q (precedence violated)", got.FilterData.Name, "canonical")
	}
}

// TestMesheryFilterRequestBody_UnmarshalAbsentZeroes verifies that a
// payload omitting both wrapper spellings leaves FilterData nil
// without returning an error, matching stdlib json.Unmarshal.
func TestMesheryFilterRequestBody_UnmarshalAbsentZeroes(t *testing.T) {
	var got MesheryFilterRequestBody
	if err := json.Unmarshal([]byte(`{"url":"https://example.com"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.URL != "https://example.com" {
		t.Errorf("URL = %q, want %q", got.URL, "https://example.com")
	}
	if got.FilterData != nil {
		t.Errorf("FilterData = %+v, want nil", got.FilterData)
	}
}

// TestMesheryFilterRequestBody_UnmarshalResetsOnReuse locks in stdlib
// json.Unmarshal semantics across decodes: a reused
// MesheryFilterRequestBody must reset FilterData to nil if the next
// payload omits both spellings — otherwise stale filter data leaks
// across requests.
func TestMesheryFilterRequestBody_UnmarshalResetsOnReuse(t *testing.T) {
	var p MesheryFilterRequestBody
	first := buildFilterRequestBody("filterData", "first-filter")
	if err := json.Unmarshal([]byte(first), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.FilterData == nil || p.FilterData.Name != "first-filter" {
		t.Fatalf("prime decode wrong: %+v", p.FilterData)
	}
	if err := json.Unmarshal([]byte(`{"url":"https://second"}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.FilterData != nil {
		t.Errorf("FilterData leaked stale value %+v across reuse", p.FilterData)
	}
	if p.URL != "https://second" {
		t.Errorf("URL = %q, want %q", p.URL, "https://second")
	}
}

// TestFilterNameForURLImport pins how a filter imported from a URL is named.
// A WebAssembly module carries no name, so it is named after the file in the
// URL path; any other body is a filter YAML file and keeps its own name field.
func TestFilterNameForURLImport(t *testing.T) {
	module := []byte("\x00asm\x01\x00\x00\x00")
	cases := []struct {
		name    string
		url     string
		body    []byte
		want    string
		wantErr error
	}{
		{name: "wasm named after the file", url: "https://example.com/filters/rate-limit.wasm", body: module, want: "rate-limit"},
		{name: "wasm extension in any case", url: "https://example.com/Header.WASM", body: module, want: "Header"},
		{name: "query string is ignored", url: "https://example.com/f/rate_limit.wasm?token=secret", body: module, want: "rate_limit"},
		{name: "escaped path is decoded", url: "https://example.com/my%20filter.wasm", body: module, want: "my filter"},
		{name: "file without the wasm extension", url: "https://example.com/download/filter", body: module, want: "filter"},
		{name: "no file in the path falls back to the host", url: "https://filters.example.com/", body: module, want: "filters.example.com"},
		{name: "directory path falls back to the host", url: "https://example.com/downloads/", body: module, want: "example.com"},
		{name: "yaml keeps its name field", url: "https://example.com/filter.yaml", body: []byte("name: from-yaml\n"), want: "from-yaml"},
		{name: "yaml without a name is rejected", url: "https://example.com/filter.yaml", body: []byte("kind: filter\n"), wantErr: ErrFilterFileName},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterNameForURLImport(tc.url, tc.body)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("name = %q, want %q", got, tc.want)
			}
		})
	}
}
