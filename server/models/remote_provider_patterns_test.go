package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
)

// TestSaveMesheryPattern_SendsPatternDataWrapperKey verifies that
// SaveMesheryPattern serializes the upsert body with the "patternData"
// wrapper key that Meshery Cloud's MesheryPatternRequestBody.PatternData
// JSON tag requires. A regression to "patternFile" or "pattern_data" at
// the wrapper would cause Cloud to decode PatternData as nil and return
// 400 "invalid request", which the server translates to the user-visible
// "Bad request. The design might be corrupt" on the Kanvas save path.
func TestSaveMesheryPattern_SendsPatternDataWrapperKey(t *testing.T) {
	var capturedBody []byte

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/patterns" {
			http.NotFound(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		capturedBody = body

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistMesheryPatterns, Endpoint: "/patterns"}}

	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate pattern id: %v", err)
	}
	pattern := &MesheryPattern{
		ID:          &id,
		Name:        "Untitled Design",
		PatternFile: "name: Untitled Design\n",
	}

	if _, err := provider.SaveMesheryPattern("token", pattern); err != nil {
		t.Fatalf("expected save to succeed, got error: %v", err)
	}

	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(capturedBody, &envelope); err != nil {
		t.Fatalf("expected JSON envelope, got %v\nbody: %s", err, string(capturedBody))
	}

	if _, ok := envelope["patternData"]; !ok {
		t.Errorf("expected request body to carry wrapper key %q (Cloud's MesheryPatternRequestBody.PatternData JSON tag), got keys: %v\nbody: %s", "patternData", keysOf(envelope), string(capturedBody))
	}
	if _, ok := envelope["patternFile"]; ok {
		t.Errorf("unexpected wrapper key %q in request body — that's the legacy drift Cloud no longer accepts\nbody: %s", "patternFile", string(capturedBody))
	}
	if _, ok := envelope["pattern_data"]; ok {
		t.Errorf("unexpected snake_case wrapper key %q in request body — Cloud's tag is camelCase\nbody: %s", "pattern_data", string(capturedBody))
	}

	// Also verify that the inner pattern still carries its own "patternFile"
	// content field (distinct from the wrapper): this is the persisted
	// design yaml. Cloud's models.MesheryPattern.PatternFile JSON tag is
	// "patternFile", so the inner field tag must stay camelCase too.
	inner, ok := envelope["patternData"]
	if !ok {
		return // first assertion already failed
	}
	var innerFields map[string]json.RawMessage
	if err := json.Unmarshal(inner, &innerFields); err != nil {
		t.Fatalf("expected inner patternData object, got %v\nraw: %s", err, string(inner))
	}
	if _, ok := innerFields["patternFile"]; !ok {
		t.Errorf("expected inner %q camelCase tag on the pattern content field, got keys: %v", "patternFile", keysOf(innerFields))
	}
	if _, ok := innerFields["pattern_file"]; ok {
		t.Errorf("unexpected snake_case inner key %q — that's the drift; Cloud's tag is camelCase", "pattern_file")
	}
}

func keysOf(m map[string]json.RawMessage) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
