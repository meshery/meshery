package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
)

// A rejected page or pagesize is a client error. Uses the real local provider,
// which parses both values before it touches the database.
func TestGetMesheryFiltersHandler_InvalidPagination_ReturnsBadRequest(t *testing.T) {
	for _, query := range []string{"page=-1", "pagesize=abc"} {
		t.Run(query, func(t *testing.T) {
			provider := &models.DefaultLocalProvider{}
			h := &Handler{log: newTestLogger(t)}

			req := httptest.NewRequest(http.MethodGet, "/api/filter?"+query, nil)
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.GetMesheryFiltersHandler(rec, req, nil, nil, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status %d, got %d. body: %s", http.StatusBadRequest, rec.Code, rec.Body.String())
			}
		})
	}
}
