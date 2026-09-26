package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

// TestHandlers_TokenlessContext_ReturnErrorNotPanic is the regression guard
// for issue #22095. Several handlers read the auth token from the request
// context with an unchecked type assertion:
//
//	tokenString := req.Context().Value(models.TokenCtxKey).(string)
//
// When the token is absent from the context (an unauthenticated request, an
// unexpected middleware ordering, or a direct handler unit test) that
// assertion panics with:
//
//	panic: interface conversion: interface {} is nil, not string
//
// crashing the serving goroutine. The fix switches every such site to the
// comma-ok idiom and returns a structured ErrFetchToken (HTTP 500) instead.
//
// Each case invokes a handler with a request whose context carries NO token.
// Pre-fix the assertion panics and the deferred recover turns that into a
// test failure; post-fix the handler returns 500 without panicking. The eight
// handlers pinned here are the six named in #22095 plus the two
// GetMesheryFilters variants, which carried the identical anti-pattern.
func TestHandlers_TokenlessContext_ReturnErrorNotPanic(t *testing.T) {
	h := &Handler{log: newTestLogger(t)}

	// A syntactically valid result id so GetResultHandler clears its
	// missing-id and invalid-UUID guards (both 400 branches) and actually
	// reaches the token assertion. Every other handler reads the token
	// before any request-specific validation.
	validResultID := uuid.Must(uuid.NewV4()).String()

	cases := []struct {
		name   string
		vars   map[string]string
		invoke func(w http.ResponseWriter, r *http.Request)
	}{
		{"FetchResultsHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.FetchResultsHandler(w, r, nil, nil, nil)
		}},
		{"FetchAllResultsHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.FetchAllResultsHandler(w, r, nil, nil, nil)
		}},
		{"GetResultHandler", map[string]string{"id": validResultID}, func(w http.ResponseWriter, r *http.Request) {
			h.GetResultHandler(w, r, nil, nil, nil)
		}},
		{"GetPerformanceProfilesHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.GetPerformanceProfilesHandler(w, r, nil, nil, nil)
		}},
		{"GetMesheryPatternsHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.GetMesheryPatternsHandler(w, r, nil, nil, nil)
		}},
		{"GetCatalogMesheryPatternsHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.GetCatalogMesheryPatternsHandler(w, r, nil, nil, nil)
		}},
		{"GetMesheryFiltersHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.GetMesheryFiltersHandler(w, r, nil, nil, nil)
		}},
		{"GetCatalogMesheryFiltersHandler", nil, func(w http.ResponseWriter, r *http.Request) {
			h.GetCatalogMesheryFiltersHandler(w, r, nil, nil, nil)
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.vars != nil {
				req = mux.SetURLVars(req, tc.vars)
			}
			// Intentionally do NOT inject models.TokenCtxKey into the
			// context: this is the exact condition that used to panic.
			rec := httptest.NewRecorder()

			defer func() {
				if p := recover(); p != nil {
					t.Fatalf("%s panicked on token-less context (issue #22095 regression): %v", tc.name, p)
				}
			}()

			tc.invoke(rec, req)

			if rec.Code != http.StatusInternalServerError {
				t.Errorf("%s: status = %d, want %d (ErrFetchToken)", tc.name, rec.Code, http.StatusInternalServerError)
			}
		})
	}
}
