package handlers

// When the registry datastore fails, the read handlers must surface a 500 so a
// client can tell a real failure apart from a genuinely empty result, instead
// of the earlier behaviour of encoding an empty 200. The datastore is forced to
// fail by closing the in-memory database before the request runs, which makes
// every subsequent query return an error.

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestRegistryReadHandlers_DatastoreErrorReturns500(t *testing.T) {
	rm, db := newTestRegistryManager(t)
	require.NoError(t, db.DBClose())
	h := &Handler{registryManager: rm, log: newTestLogger(t)}

	cases := []struct {
		name string
		url  string
		vars map[string]string
		fn   http.HandlerFunc
	}{
		{"models", "/api/registry/models", nil, h.GetMeshmodelModels},
		{"components", "/api/registry/components", nil, h.GetAllMeshmodelComponents},
		{"components by model", "/api/registry/models/kubernetes/components", map[string]string{"model": "kubernetes"}, h.GetMeshmodelComponentByModel},
		{"relationships", "/api/registry/relationships", nil, h.GetAllMeshmodelRelationships},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			if tc.vars != nil {
				req = mux.SetURLVars(req, tc.vars)
			}
			rec := httptest.NewRecorder()

			tc.fn(rec, req)

			require.Equal(t, http.StatusInternalServerError, rec.Code,
				"a datastore failure must be reported as 500, not masked as an empty 200; body: %s", rec.Body.String())
		})
	}
}
