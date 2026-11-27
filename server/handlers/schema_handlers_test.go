package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func serveResourceSchema(t *testing.T, name string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/schema/resource/x", nil)
	req = mux.SetURLVars(req, map[string]string{"resourceName": name})
	rec := httptest.NewRecorder()
	(&Handler{}).HandleResourceSchemas(rec, req)
	return rec
}

func TestHandleResourceSchemasKnownResource(t *testing.T) {
	for _, name := range []string{"design", "workspace", "environment", "helmRepo"} {
		t.Run(name, func(t *testing.T) {
			rec := serveResourceSchema(t, name)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
			}
			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", ct)
			}
			var body map[string]map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("invalid json: %v", err)
			}
			if len(body["rjsfSchema"]) == 0 {
				t.Fatalf("rjsfSchema is empty: %s", rec.Body.String())
			}
			if _, ok := body["uiSchema"]; !ok {
				t.Fatalf("uiSchema missing: %s", rec.Body.String())
			}
		})
	}
}

func TestHandleResourceSchemasUnknownResource(t *testing.T) {
	for _, name := range []string{"", "doesNotExist", "../go.mod", "configuration/designImport.json", "DESIGN"} {
		t.Run(name, func(t *testing.T) {
			rec := serveResourceSchema(t, name)
			if rec.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusNotFound, rec.Body.String())
			}
		})
	}
}

type failingWriter struct {
	header       http.Header
	writeHeaders int
}

func (w *failingWriter) Header() http.Header { return w.header }
func (w *failingWriter) WriteHeader(int)     { w.writeHeaders++ }
func (w *failingWriter) Write([]byte) (int, error) {
	if w.writeHeaders == 0 {
		w.writeHeaders++ // net/http sends an implicit 200 on the first Write
	}
	return 0, errors.New("connection reset")
}

func TestHandleResourceSchemasWriteFailureDoesNotRewriteHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/schema/resource/design", nil)
	req = mux.SetURLVars(req, map[string]string{"resourceName": "design"})
	w := &failingWriter{header: http.Header{}}
	(&Handler{}).HandleResourceSchemas(w, req)
	if w.writeHeaders != 1 {
		t.Fatalf("WriteHeader called %d times, want 1", w.writeHeaders)
	}
}
