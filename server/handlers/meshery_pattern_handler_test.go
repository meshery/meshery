package handlers

import (
	"encoding/json"
	"testing"
)

// buildDesignPostBody returns a JSON body that wraps a minimal PatternFile
// under the provided key. The PatternFile carries a name so tests can
// assert which spelling "won" when multiple spellings are present. Values
// are json-escaped via json.Marshal so names containing quotes, newlines,
// or backslashes do not produce invalid JSON.
func buildDesignPostBody(key, designName string) string {
	patternFile := map[string]any{
		"id":            "00000000-0000-0000-0000-000000000000",
		"name":          designName,
		"schemaVersion": "designs.meshery.io/v1beta1",
		"version":       "0.0.1",
		"components":    []any{},
		"relationships": []any{},
	}
	body := map[string]any{
		"name": "envelope",
		key:    patternFile,
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// buildDesignPostBodyMulti returns a JSON body with multiple design-file
// spellings present, each carrying a distinct name so precedence tests
// can assert which spelling won.
func buildDesignPostBodyMulti(entries map[string]string) string {
	body := map[string]any{"name": "e"}
	for key, designName := range entries {
		body[key] = map[string]any{
			"id":            "00000000-0000-0000-0000-000000000000",
			"name":          designName,
			"schemaVersion": "designs.meshery.io/v1beta1",
			"version":       "0.0.1",
			"components":    []any{},
			"relationships": []any{},
		}
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// TestDesignPostPayload_UnmarshalAcceptsAllDesignFileKeyFlavors locks in
// the deprecation-window contract for POST /api/pattern: the handler
// accepts the canonical `designFile`, the alternate `patternFile`, and
// the legacy snake_case spellings `design_file` and `pattern_file`.
// Canonical camelCase wins over legacy; `designFile` wins over
// `patternFile` when both are present.
func TestDesignPostPayload_UnmarshalAcceptsAllDesignFileKeyFlavors(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantName string
	}{
		{
			name:     "canonical designFile only",
			body:     buildDesignPostBody("designFile", "from-designFile"),
			wantName: "from-designFile",
		},
		{
			name:     "alternate patternFile only",
			body:     buildDesignPostBody("patternFile", "from-patternFile"),
			wantName: "from-patternFile",
		},
		{
			name:     "legacy design_file only",
			body:     buildDesignPostBody("design_file", "from-design_file"),
			wantName: "from-design_file",
		},
		{
			name:     "legacy pattern_file only",
			body:     buildDesignPostBody("pattern_file", "from-pattern_file"),
			wantName: "from-pattern_file",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got DesignPostPayload
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.DesignFile.Name != tc.wantName {
				t.Errorf("DesignFile.Name = %q, want %q", got.DesignFile.Name, tc.wantName)
			}
		})
	}
}

// TestDesignPostPayload_UnmarshalPrecedenceCanonicalWins locks the
// canonical-wins-on-conflict rule: when multiple spellings of the
// design-file field are present in one payload, `designFile` wins over
// every other spelling, and `patternFile` wins over the snake_case
// legacies. This is important because clients migrating incrementally
// may temporarily emit both spellings.
func TestDesignPostPayload_UnmarshalPrecedenceCanonicalWins(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantName string
	}{
		{
			name:     "designFile beats patternFile",
			body:     buildDesignPostBodyMulti(map[string]string{"designFile": "canonical", "patternFile": "alternate"}),
			wantName: "canonical",
		},
		{
			name:     "designFile beats design_file",
			body:     buildDesignPostBodyMulti(map[string]string{"designFile": "canonical", "design_file": "legacy-snake"}),
			wantName: "canonical",
		},
		{
			name:     "patternFile beats design_file",
			body:     buildDesignPostBodyMulti(map[string]string{"patternFile": "alternate", "design_file": "legacy-snake"}),
			wantName: "alternate",
		},
		{
			name:     "design_file beats pattern_file",
			body:     buildDesignPostBodyMulti(map[string]string{"design_file": "design-legacy", "pattern_file": "pattern-legacy"}),
			wantName: "design-legacy",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got DesignPostPayload
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.DesignFile.Name != tc.wantName {
				t.Errorf("DesignFile.Name = %q, want %q (precedence violated)", got.DesignFile.Name, tc.wantName)
			}
		})
	}
}

// TestDesignPostPayload_UnmarshalAllSpellingsAbsentZeroes verifies that
// a payload containing none of the four design-file spellings leaves
// DesignFile at its zero value and does NOT return an error. This
// matches stdlib json.Unmarshal behavior for omitted fields.
func TestDesignPostPayload_UnmarshalAllSpellingsAbsentZeroes(t *testing.T) {
	var got DesignPostPayload
	if err := json.Unmarshal([]byte(`{"name":"envelope-only"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != "envelope-only" {
		t.Errorf("Name = %q, want %q", got.Name, "envelope-only")
	}
	if got.DesignFile.Name != "" {
		t.Errorf("DesignFile.Name = %q, want empty when no design-file spelling present", got.DesignFile.Name)
	}
}

// TestDesignPostPayload_UnmarshalResetsDesignFileOnReuse locks in
// stdlib json.Unmarshal semantics: when a single DesignPostPayload is
// reused across decodes, the DesignFile field must reset to its zero
// value before the second payload is applied. Otherwise a prior
// decode's components/relationships could leak into the next request.
func TestDesignPostPayload_UnmarshalResetsDesignFileOnReuse(t *testing.T) {
	var p DesignPostPayload
	first := buildDesignPostBody("designFile", "first-design")
	if err := json.Unmarshal([]byte(first), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.DesignFile.Name != "first-design" {
		t.Fatalf("prime decode wrong: DesignFile.Name = %q", p.DesignFile.Name)
	}
	// Second payload has no design-file spelling at all; DesignFile
	// must reset to zero rather than carry "first-design" forward.
	if err := json.Unmarshal([]byte(`{"name":"second-envelope"}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.DesignFile.Name != "" {
		t.Errorf("DesignFile.Name leaked stale value %q across reuse", p.DesignFile.Name)
	}
	if p.Name != "second-envelope" {
		t.Errorf("Name = %q, want %q", p.Name, "second-envelope")
	}
}

// TestDesignPostPayload_MarshalEmitsBothDesignFileKeyFlavors locks in
// the deprecation-window contract: MarshalJSON emits both the
// canonical `designFile` key AND the legacy `design_file` key so any
// external consumer still reading either spelling keeps working while
// they migrate. Once all known consumers have migrated, MarshalJSON
// and this test can be dropped.
func TestDesignPostPayload_MarshalEmitsBothDesignFileKeyFlavors(t *testing.T) {
	var p DesignPostPayload
	body := buildDesignPostBody("designFile", "round-trip")
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatalf("seed unmarshal: %v", err)
	}
	out, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	// Assert against top-level JSON keys rather than substring-matching
	// the serialized output: the latter can false-pass if the same
	// substrings ever appear inside nested objects or escaped strings.
	var top map[string]json.RawMessage
	if err := json.Unmarshal(out, &top); err != nil {
		t.Fatalf("could not re-decode marshal output to top-level keys: %v (out=%s)", err, string(out))
	}
	for _, wantKey := range []string{"designFile", "design_file"} {
		if _, ok := top[wantKey]; !ok {
			t.Errorf("marshal output missing top-level key %q; got keys %v", wantKey, topLevelKeys(top))
		}
	}
}

func topLevelKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
