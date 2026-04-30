package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerWrappers_ForwardToHttputil is a smoke test for the four thin
// wrappers in utils.go that forward to server/models/httputil. The full
// behavioral surface (headers, MeshKit serialization, empty-object body)
// is covered by server/models/httputil/httputil_test.go; this test only
// guards the wrapper layer itself against accidental signature drift or
// a missed rename during a future refactor.
func TestHandlerWrappers_ForwardToHttputil(t *testing.T) {
	t.Run("writeJSONError", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONError(rec, "some error", http.StatusBadRequest)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected 400, got %d", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("expected JSON Content-Type, got %q", ct)
		}

		var decoded map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("body did not parse as JSON: %v", err)
		}
		if decoded["error"] != "some error" {
			t.Errorf("expected wrapper to pass message through, got %q", decoded["error"])
		}
	})

	t.Run("writeMeshkitError", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeMeshkitError(rec, nil, http.StatusInternalServerError)

		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("expected 500, got %d", rec.Code)
		}
		if ct := rec.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Errorf("expected JSON Content-Type, got %q", ct)
		}
	})

	t.Run("writeJSONMessage", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONMessage(rec, map[string]string{"message": "ok"}, http.StatusOK)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}

		var decoded map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&decoded); err != nil {
			t.Fatalf("body did not parse as JSON: %v", err)
		}
		if decoded["message"] != "ok" {
			t.Errorf("expected wrapper to pass payload through, got %q", decoded["message"])
		}
	})

	t.Run("writeJSONEmptyObject", func(t *testing.T) {
		rec := httptest.NewRecorder()
		writeJSONEmptyObject(rec, http.StatusOK)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if body := rec.Body.String(); body != "{}" {
			t.Errorf("expected body %q, got %q", "{}", body)
		}
	})
}
