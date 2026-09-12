package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TestRegisterRegistryRoute locks the capabilities-registry alias contract:
// every subpath registered through registerRegistryRoute must resolve at both
// the canonical /api/registry path and the deprecated /api/meshmodels alias,
// for exactly the registered methods, and both must hit the same handler.
func TestRegisterRegistryRoute(t *testing.T) {
	gMux := mux.NewRouter()

	var hits []string
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits = append(hits, r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	registerRegistryRoute(gMux, "/models", handler, "GET")
	registerRegistryRoute(gMux, "/models/{model}", handler, "GET")
	registerRegistryRoute(gMux, "/{entityType}/status", handler, "POST")

	cases := []struct {
		method string
		path   string
		want   int
	}{
		{"GET", "/api/registry/models", http.StatusNoContent},
		{"GET", "/api/meshmodels/models", http.StatusNoContent},
		{"GET", "/api/registry/models/istio", http.StatusNoContent},
		{"GET", "/api/meshmodels/models/istio", http.StatusNoContent},
		{"POST", "/api/registry/models/status", http.StatusNoContent},
		{"POST", "/api/meshmodels/models/status", http.StatusNoContent},
		// Methods not registered for a subpath must not resolve on either form.
		{"POST", "/api/registry/models", http.StatusMethodNotAllowed},
		{"POST", "/api/meshmodels/models", http.StatusMethodNotAllowed},
		// Unregistered prefixes must not resolve.
		{"GET", "/api/meshmodel/models", http.StatusNotFound},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		gMux.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Errorf("%s %s: got status %d, want %d", tc.method, tc.path, rec.Code, tc.want)
		}
	}

	// Six successful requests above must all have reached the shared handler.
	if len(hits) != 6 {
		t.Errorf("handler hit %d times, want 6 (canonical and alias must share one handler)", len(hits))
	}
}
