package handlers

import (
	"encoding/json"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
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

// -----------------------------------------------------------------------------
// MesheryPatternPOSTRequestBody wrapper-key dual-accept tests
// -----------------------------------------------------------------------------

// buildPatternPostBody returns a JSON body that wraps a minimal
// DesignPostPayload (carrying a name) under the provided wrapper key.
// Tests use the inner design name to assert which spelling "won" when
// multiple spellings are present.
func buildPatternPostBody(key, designName string) string {
	inner := map[string]any{
		"id":   "00000000-0000-0000-0000-000000000000",
		"name": designName,
	}
	body := map[string]any{
		"name": "envelope",
		key:    inner,
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// buildPatternPostBodyMulti returns a JSON body with multiple wrapper
// spellings present, each carrying a distinct inner name so precedence
// tests can assert which spelling won.
func buildPatternPostBodyMulti(entries map[string]string) string {
	body := map[string]any{"name": "e"}
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

// TestMesheryPatternPOSTRequestBody_UnmarshalAcceptsBothSpellings locks
// in the deprecation-window contract for POST /api/pattern: the
// handler accepts both the canonical `patternData` (camelCase) and the
// legacy `pattern_data` (snake_case) wrapper keys.
func TestMesheryPatternPOSTRequestBody_UnmarshalAcceptsBothSpellings(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantName string
	}{
		{
			name:     "canonical patternData only",
			body:     buildPatternPostBody("patternData", "from-patternData"),
			wantName: "from-patternData",
		},
		{
			name:     "legacy pattern_data only",
			body:     buildPatternPostBody("pattern_data", "from-pattern_data"),
			wantName: "from-pattern_data",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got MesheryPatternPOSTRequestBody
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.PatternData == nil {
				t.Fatalf("PatternData is nil, want populated")
			}
			if got.PatternData.Name != tc.wantName {
				t.Errorf("PatternData.Name = %q, want %q", got.PatternData.Name, tc.wantName)
			}
		})
	}
}

// TestMesheryPatternPOSTRequestBody_UnmarshalPrecedenceCanonicalWins
// locks the canonical-wins rule: when both `patternData` and
// `pattern_data` are present in one payload, the canonical camelCase
// spelling wins. Migrating clients temporarily emit both.
func TestMesheryPatternPOSTRequestBody_UnmarshalPrecedenceCanonicalWins(t *testing.T) {
	body := buildPatternPostBodyMulti(map[string]string{
		"patternData":  "canonical",
		"pattern_data": "legacy",
	})
	var got MesheryPatternPOSTRequestBody
	if err := json.Unmarshal([]byte(body), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.PatternData == nil {
		t.Fatalf("PatternData is nil, want canonical")
	}
	if got.PatternData.Name != "canonical" {
		t.Errorf("PatternData.Name = %q, want %q (precedence violated)", got.PatternData.Name, "canonical")
	}
}

// TestMesheryPatternPOSTRequestBody_UnmarshalAbsentZeroes verifies that
// a payload omitting both wrapper spellings leaves PatternData nil
// without returning an error, matching stdlib json.Unmarshal behavior.
func TestMesheryPatternPOSTRequestBody_UnmarshalAbsentZeroes(t *testing.T) {
	var got MesheryPatternPOSTRequestBody
	if err := json.Unmarshal([]byte(`{"name":"envelope-only"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != "envelope-only" {
		t.Errorf("Name = %q, want %q", got.Name, "envelope-only")
	}
	if got.PatternData != nil {
		t.Errorf("PatternData = %+v, want nil", got.PatternData)
	}
}

// TestMesheryPatternPOSTRequestBody_UnmarshalResetsOnReuse locks in
// stdlib json.Unmarshal semantics across decodes: when a single
// MesheryPatternPOSTRequestBody is reused, PatternData must reset to
// nil if the next payload omits both spellings — otherwise stale
// design data leaks across requests.
func TestMesheryPatternPOSTRequestBody_UnmarshalResetsOnReuse(t *testing.T) {
	var p MesheryPatternPOSTRequestBody
	first := buildPatternPostBody("patternData", "first-design")
	if err := json.Unmarshal([]byte(first), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.PatternData == nil || p.PatternData.Name != "first-design" {
		t.Fatalf("prime decode wrong: %+v", p.PatternData)
	}
	if err := json.Unmarshal([]byte(`{"name":"second-envelope"}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.PatternData != nil {
		t.Errorf("PatternData leaked stale value %+v across reuse", p.PatternData)
	}
	if p.Name != "second-envelope" {
		t.Errorf("Name = %q, want %q", p.Name, "second-envelope")
	}
}

// -----------------------------------------------------------------------------
// MesheryPatternUPDATERequestBody wrapper-key dual-accept tests
// -----------------------------------------------------------------------------

// buildPatternUpdateBody returns a JSON body for PUT /api/pattern.
// The MesheryPattern wrapped under patternKey carries a name; the
// cytoscape payload is wrapped under cytoscapeKey as a plain string.
// Empty key strings skip emitting that wrapper.
func buildPatternUpdateBody(patternKey, patternName, cytoscapeKey, cytoscapeVal string) string {
	body := map[string]any{
		"name": "envelope",
		"url":  "",
	}
	if patternKey != "" {
		body[patternKey] = map[string]any{
			"id":   "00000000-0000-0000-0000-000000000000",
			"name": patternName,
		}
	}
	if cytoscapeKey != "" {
		body[cytoscapeKey] = cytoscapeVal
	}
	out, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	return string(out)
}

// TestMesheryPatternUPDATERequestBody_UnmarshalPatternDataAcceptsBothSpellings
// locks in the dual-accept contract on PatternData for PUT /api/pattern.
func TestMesheryPatternUPDATERequestBody_UnmarshalPatternDataAcceptsBothSpellings(t *testing.T) {
	cases := []struct {
		name     string
		body     string
		wantName string
	}{
		{
			name:     "canonical patternData only",
			body:     buildPatternUpdateBody("patternData", "from-patternData", "", ""),
			wantName: "from-patternData",
		},
		{
			name:     "legacy pattern_data only",
			body:     buildPatternUpdateBody("pattern_data", "from-pattern_data", "", ""),
			wantName: "from-pattern_data",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got MesheryPatternUPDATERequestBody
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.PatternData == nil {
				t.Fatalf("PatternData is nil, want populated")
			}
			if got.PatternData.Name != tc.wantName {
				t.Errorf("PatternData.Name = %q, want %q", got.PatternData.Name, tc.wantName)
			}
		})
	}
}

// TestMesheryPatternUPDATERequestBody_UnmarshalCytoscapeAcceptsBothSpellings
// locks in the dual-accept contract on CytoscapeJSON for PUT /api/pattern.
func TestMesheryPatternUPDATERequestBody_UnmarshalCytoscapeAcceptsBothSpellings(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
	}{
		{
			name: "canonical cytoscapeJSON only",
			body: buildPatternUpdateBody("", "", "cytoscapeJSON", "graph-canonical"),
			want: "graph-canonical",
		},
		{
			name: "legacy cytoscape_json only",
			body: buildPatternUpdateBody("", "", "cytoscape_json", "graph-legacy"),
			want: "graph-legacy",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got MesheryPatternUPDATERequestBody
			if err := json.Unmarshal([]byte(tc.body), &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got.CytoscapeJSON != tc.want {
				t.Errorf("CytoscapeJSON = %q, want %q", got.CytoscapeJSON, tc.want)
			}
		})
	}
}

// TestMesheryPatternUPDATERequestBody_UnmarshalPrecedenceCanonicalWins
// locks the canonical-wins rule: when both spellings are present on
// either field, the camelCase spelling wins.
func TestMesheryPatternUPDATERequestBody_UnmarshalPrecedenceCanonicalWins(t *testing.T) {
	raw, err := json.Marshal(map[string]any{
		"name":           "envelope",
		"patternData":    map[string]any{"id": "00000000-0000-0000-0000-000000000000", "name": "canonical"},
		"pattern_data":   map[string]any{"id": "00000000-0000-0000-0000-000000000000", "name": "legacy"},
		"cytoscapeJSON":  "canonical-graph",
		"cytoscape_json": "legacy-graph",
	})
	if err != nil {
		t.Fatalf("marshal input: %v", err)
	}
	var got MesheryPatternUPDATERequestBody
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.PatternData == nil || got.PatternData.Name != "canonical" {
		t.Errorf("PatternData = %+v, want name=%q", got.PatternData, "canonical")
	}
	if got.CytoscapeJSON != "canonical-graph" {
		t.Errorf("CytoscapeJSON = %q, want %q", got.CytoscapeJSON, "canonical-graph")
	}
}

// TestMesheryPatternUPDATERequestBody_UnmarshalAbsentZeroes verifies that
// a payload containing neither PatternData spelling nor either
// CytoscapeJSON spelling leaves both fields at zero values without
// returning an error.
func TestMesheryPatternUPDATERequestBody_UnmarshalAbsentZeroes(t *testing.T) {
	var got MesheryPatternUPDATERequestBody
	if err := json.Unmarshal([]byte(`{"name":"envelope-only"}`), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Name != "envelope-only" {
		t.Errorf("Name = %q, want %q", got.Name, "envelope-only")
	}
	if got.PatternData != nil {
		t.Errorf("PatternData = %+v, want nil", got.PatternData)
	}
	if got.CytoscapeJSON != "" {
		t.Errorf("CytoscapeJSON = %q, want empty", got.CytoscapeJSON)
	}
}

// TestMesheryPatternUPDATERequestBody_UnmarshalResetsOnReuse locks in
// reuse semantics: a reused MesheryPatternUPDATERequestBody must reset
// PatternData to nil and CytoscapeJSON to "" when the next payload
// omits both spellings.
func TestMesheryPatternUPDATERequestBody_UnmarshalResetsOnReuse(t *testing.T) {
	var p MesheryPatternUPDATERequestBody
	first := buildPatternUpdateBody("patternData", "first-design", "cytoscapeJSON", "first-graph")
	if err := json.Unmarshal([]byte(first), &p); err != nil {
		t.Fatalf("first unmarshal: %v", err)
	}
	if p.PatternData == nil || p.PatternData.Name != "first-design" {
		t.Fatalf("prime decode PatternData wrong: %+v", p.PatternData)
	}
	if p.CytoscapeJSON != "first-graph" {
		t.Fatalf("prime decode CytoscapeJSON = %q, want %q", p.CytoscapeJSON, "first-graph")
	}
	if err := json.Unmarshal([]byte(`{"name":"second-envelope"}`), &p); err != nil {
		t.Fatalf("second unmarshal: %v", err)
	}
	if p.PatternData != nil {
		t.Errorf("PatternData leaked stale value %+v across reuse", p.PatternData)
	}
	if p.CytoscapeJSON != "" {
		t.Errorf("CytoscapeJSON leaked stale value %q across reuse", p.CytoscapeJSON)
	}
	if p.Name != "second-envelope" {
		t.Errorf("Name = %q, want %q", p.Name, "second-envelope")
	}
}

// Compile-time type reference: pulls models into the test graph so the
// MesheryPattern field type resolves without introducing import churn
// if the test file grows additional model-backed cases later.
var _ = models.MesheryPattern{}

// TestConvertFileToDesign_UnsupportedExtensionWrappedAsErrConvertToDesign
// guards the contract that VerifyAndConvertToDesign promotes raw
// ConvertFileToDesign failures into MeshKit-wrapped errors. Without a
// wrap, regressions back to `return err` would silently strip MeshKit
// metadata from the JSON envelope returned to the client. This test
// asserts that:
//  1. ConvertFileToDesign produces an error for an unsupported file
//     extension (the cheapest deterministic failure path that does not
//     require a registry, temp-FS, or fixture data).
//  2. Wrapping that error with ErrConvertToDesign rewrites the MeshKit
//     code to ErrConvertToDesignCode, which is what
//     VerifyAndConvertToDesign at line 450 of meshery_pattern_handler.go
//     returns to the handler.
func TestConvertFileToDesign_UnsupportedExtensionWrappedAsErrConvertToDesign(t *testing.T) {
	log := newTestLogger(t)

	// A bogus extension is rejected by SanitizeFile before any registry
	// or pattern-engine work, so passing a nil registry is safe.
	_, _, err := ConvertFileToDesign(FileToImport{
		Data:     []byte("not a real design"),
		FileName: "fixture.unsupported",
	}, nil, log)
	if err == nil {
		t.Fatal("ConvertFileToDesign accepted an unsupported extension; expected a sanitization failure")
	}

	wrapped := ErrConvertToDesign(err)
	if got, want := errors.GetCode(wrapped), ErrConvertToDesignCode; got != want {
		t.Fatalf("ErrConvertToDesign produced code %q, want %q", got, want)
	}
}

// TestErrEncodePattern_PreservesMeshKitCode mirrors the guard above for
// the encoding.Marshal failure path inside VerifyAndConvertToDesign at
// line 455 of meshery_pattern_handler.go: a regression back to
// `return err` would surface a raw stdlib error with no MeshKit code,
// and the JSON envelope sent to clients would lose its
// code/severity/remediation metadata. The Marshal failure is hard to
// reach deterministically without depending on internal pattern shape,
// so this test asserts the simpler invariant: ErrEncodePattern
// rewrites whatever it wraps to ErrEncodePatternCode.
func TestErrEncodePattern_PreservesMeshKitCode(t *testing.T) {
	wrapped := ErrEncodePattern(json.Unmarshal([]byte("not-json"), &struct{}{}))
	if got, want := errors.GetCode(wrapped), ErrEncodePatternCode; got != want {
		t.Fatalf("ErrEncodePattern produced code %q, want %q", got, want)
	}
}
