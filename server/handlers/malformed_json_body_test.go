package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// A request body that is not valid JSON is a client error and must be answered
// with 400, not 500.
func TestMalformedJSONBodyIsBadRequest(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	h := newTestHandler(t, map[string]models.Provider{local.Name(): local}, "")
	user := &models.User{}
	const malformed = `{"name":`

	cases := []struct {
		name string
		req  func() *http.Request
		call func(http.ResponseWriter, *http.Request)
	}{
		{
			name: "POST /api/workspaces",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/workspaces", strings.NewReader(malformed))
			},
			call: func(w http.ResponseWriter, r *http.Request) { h.SaveWorkspaceHandler(w, r, nil, user, local) },
		},
		{
			name: "PUT /api/workspaces/{id}",
			req: func() *http.Request {
				id := uuid.Must(uuid.NewV4()).String()
				r := httptest.NewRequest(http.MethodPut, "/api/workspaces/"+id, strings.NewReader(malformed))
				return mux.SetURLVars(r, map[string]string{"id": id})
			},
			call: func(w http.ResponseWriter, r *http.Request) { h.UpdateWorkspaceHandler(w, r, nil, user, local) },
		},
		{
			name: "POST /api/user/prefs",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodPost, "/api/user/prefs", strings.NewReader(malformed))
			},
			call: func(w http.ResponseWriter, r *http.Request) {
				h.UserPrefsHandler(w, r, &models.Preference{}, user, local)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.call(rec, tc.req())
			if rec.Code != http.StatusBadRequest {
				t.Errorf("%s with body %q: status = %d, want %d", tc.name, malformed, rec.Code, http.StatusBadRequest)
			}
		})
	}
}
