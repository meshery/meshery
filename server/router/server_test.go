package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestClose(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run Close() skipping")
	//err := r.Close()
	//if err != nil {
	//	t.Errorf("Close() failed with error: %s", err)
	//}
}

func TestDeprecatedMeshmodelHandler_ForwardsAndSetsHeaders(t *testing.T) {
	gMux := mux.NewRouter()
	gMux.HandleFunc("/api/registry/models", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok:" + r.URL.RawQuery))
	}).Methods("GET")

	gMux.PathPrefix("/api/meshmodels").Handler(deprecatedMeshmodelHandler(gMux))
	gMux.PathPrefix("/api/meshmodel").Handler(deprecatedMeshmodelHandler(gMux))

	req := httptest.NewRequest("GET", "/api/meshmodels/models?page=1&pagesize=5", nil)
	rec := httptest.NewRecorder()

	gMux.ServeHTTP(rec, req)

	resp := rec.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Deprecation"); got != "true" {
		t.Errorf("Deprecation header = %q, want %q", got, "true")
	}
	if got := resp.Header.Get("Sunset"); got != "Tue, 01 Dec 2026 00:00:00 GMT" {
		t.Errorf("Sunset header = %q, want %q", got, "Tue, 01 Dec 2026 00:00:00 GMT")
	}
	wantLink := `</api/registry/models?page=1&pagesize=5>; rel="successor-version"`
	if got := resp.Header.Get("Link"); got != wantLink {
		t.Errorf("Link header = %q, want %q", got, wantLink)
	}
	if got := rec.Body.String(); got != "ok:page=1&pagesize=5" {
		t.Errorf("body = %q, want re-dispatched body with query intact", got)
	}
}
