package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWriteJSONError_ShapeIsParseableJSON guards the response shape of the
// validation-failure path on /api/workspaces and /api/environments. The
// symptom this prevents is RTK Query's default baseQuery (which dispatches
// on Content-Type) throwing `SyntaxError: Unexpected token 'W', "WorkspaceI"...`
// when the server emitted a plain-text 400 body like
// "WorkspaceID or OrgID cannot be empty". The contract: status code is
// honored, Content-Type is application/json, and the body JSON-parses to
// {"error": "<message>"}.
func TestWriteJSONError_ShapeIsParseableJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSONError(rec, "WorkspaceID or OrgID cannot be empty", http.StatusBadRequest)

	resp := rec.Result()
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type %q, got %q — a non-JSON Content-Type is what broke RTK Query", "application/json; charset=utf-8", ct)
	}

	if nosniff := resp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got %q", nosniff)
	}

	var decoded map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("expected body to parse as JSON, got %v", err)
	}

	if got := decoded["error"]; got != "WorkspaceID or OrgID cannot be empty" {
		t.Errorf("expected error field %q, got %q", "WorkspaceID or OrgID cannot be empty", got)
	}
}

// TestWriteJSONError_DoesNotStartWithBareWord pins the regression-of-interest:
// a plain-text body beginning with "W" (as http.Error would emit for the
// "WorkspaceID or OrgID cannot be empty" message) is exactly what crashed
// RTK Query's JSON parser. A JSON-encoded body must start with '{'.
func TestWriteJSONError_DoesNotStartWithBareWord(t *testing.T) {
	rec := httptest.NewRecorder()
	writeJSONError(rec, "WorkspaceID or OrgID cannot be empty", http.StatusBadRequest)

	body := rec.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected a non-empty body")
	}
	if body[0] != '{' {
		t.Errorf("expected body to start with '{' (JSON object), got %q — this is the hazard RTK Query trips on", string(body[:min(20, len(body))]))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
