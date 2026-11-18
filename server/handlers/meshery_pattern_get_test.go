package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// TestGetMesheryPatternHandler_LocalProvider drives GET /api/pattern/{id}
// through the Local provider and a real database. A design that does not
// exist must be reported as 404 with an error body, never as a design; a
// malformed ID stays a 400 and a database failure is not a client error.
func TestGetMesheryPatternHandler_LocalProvider(t *testing.T) {
	tests := []struct {
		name       string
		pathID     string // overrides the generated ID when set
		setup      func(t *testing.T, f patternPOSTFixture, id uuid.UUID)
		wantStatus int
	}{
		{
			name:       "stored design is returned",
			setup:      func(t *testing.T, f patternPOSTFixture, id uuid.UUID) { f.seed(t, id) },
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing design is not found",
			setup:      func(*testing.T, patternPOSTFixture, uuid.UUID) {},
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "malformed ID is a bad request",
			pathID:     "not-a-uuid",
			setup:      func(*testing.T, patternPOSTFixture, uuid.UUID) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "database failure is a server error",
			setup: func(t *testing.T, f patternPOSTFixture, _ uuid.UUID) {
				if err := f.db.Migrator().DropTable(&models.MesheryPattern{}); err != nil {
					t.Fatalf("drop table: %v", err)
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			f := newPatternPOSTFixture(t)
			id := uuid.Must(uuid.NewV4())
			tc.setup(t, f, id)

			pathID := id.String()
			if tc.pathID != "" {
				pathID = tc.pathID
			}
			req := httptest.NewRequest(http.MethodGet, "/api/pattern/"+pathID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": pathID})
			rec := httptest.NewRecorder()

			f.h.GetMesheryPatternHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, f.provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}

			var decoded struct {
				ID   *uuid.UUID `json:"id"`
				Code string     `json:"code"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
				t.Fatalf("body did not parse as JSON: %v (body=%q)", err, rec.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				if decoded.ID == nil || *decoded.ID != id {
					t.Errorf("design id = %v, want %s (body=%q)", decoded.ID, id, rec.Body.String())
				}
				return
			}
			if decoded.Code == "" {
				t.Errorf("expected a non-empty error code (body=%q)", rec.Body.String())
			}
			if decoded.ID != nil {
				t.Errorf("error response carries a design id %s (body=%q)", decoded.ID, rec.Body.String())
			}
		})
	}
}
