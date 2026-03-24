package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestDefaultAPIContentTypeMiddlewareSetsJSONForAliasedRoutes(t *testing.T) {
	r := mux.NewRouter()
	r.Use(defaultAPIContentTypeMiddleware)

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	})

	r.Handle("/api/test", handler).Methods(http.MethodGet)
	r.Handle("/api/test/alias", handler).Methods(http.MethodGet)

	for _, path := range []string{"/api/test", "/api/test/alias"} {
		t.Run(path, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			r.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
			}

			contentType := recorder.Header().Get("Content-Type")
			if !strings.HasPrefix(contentType, defaultAPIContentType) {
				t.Fatalf("content-type = %q, want prefix %q", contentType, defaultAPIContentType)
			}
		})
	}
}

func TestDefaultAPIContentTypeMiddlewarePreservesExplicitContentType(t *testing.T) {
	r := mux.NewRouter()
	r.Use(defaultAPIContentTypeMiddleware)
	r.HandleFunc("/api/file", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("plain text"))
	}).Methods(http.MethodGet)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/file", nil)

	r.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	contentType := recorder.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/plain") {
		t.Fatalf("content-type = %q, want prefix %q", contentType, "text/plain")
	}
}

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
