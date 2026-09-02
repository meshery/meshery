package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newPatternUUIDTestHandler(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}
}

func TestGetMesheryPatternHandler_InvalidIDReturns400(t *testing.T) {
	cases := []struct {
		name string
		id   string
	}{
		{"empty string", ""},
		{"undefined", "undefined"},
		{"nil UUID", uuid.Nil.String()},
		{"garbage", "not-a-uuid"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newPatternUUIDTestHandler(t)
			req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+tc.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			rec := httptest.NewRecorder()

			h.GetMesheryPatternHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("id=%q: status = %d, want 400; body: %s", tc.id, rec.Code, rec.Body.String())
			}
		})
	}
}

func TestCloneMesheryPatternHandler_InvalidIDReturns400(t *testing.T) {
	cases := []struct {
		name string
		id   string
	}{
		{"empty string", ""},
		{"undefined", "undefined"},
		{"nil UUID", uuid.Nil.String()},
		{"garbage", "not-a-uuid"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newPatternUUIDTestHandler(t)
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/clone/"+tc.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			rec := httptest.NewRecorder()

			h.CloneMesheryPatternHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, nil)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("id=%q: status = %d, want 400; body: %s", tc.id, rec.Code, rec.Body.String())
			}
		})
	}
}
