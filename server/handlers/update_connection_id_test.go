package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// putConnection sends PUT /api/integrations/connections/{connectionId} with a
// body that carries no id, as the UI's partial updates do.
func putConnection(h *Handler, provider models.Provider, id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+id, strings.NewReader(`{"name":"renamed","kind":"prometheus"}`))
	req = mux.SetURLVars(req, map[string]string{"connectionId": id})
	rec := httptest.NewRecorder()
	h.UpdateConnectionById(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)
	return rec
}

func countConnections(t *testing.T, provider *models.DefaultLocalProvider) int64 {
	t.Helper()
	var n int64
	require.NoError(t, provider.ConnectionPersister.DB.Model(&connections.Connection{}).Count(&n).Error)
	return n
}

// TestUpdateConnectionByIDMalformedIDIsBadRequest checks that an ID that is not
// a UUID is a client error and does not reach the database.
func TestUpdateConnectionByIDMalformedIDIsBadRequest(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)

	rec := putConnection(h, provider, "not-a-uuid")

	assert.Equal(t, http.StatusBadRequest, rec.Code, rec.Body.String())
	assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
}

// TestUpdateConnectionByIDUnknownIDIsNotFound checks that updating the nil ID,
// or an ID that was never issued, answers 404 instead of creating a connection
// under it.
func TestUpdateConnectionByIDUnknownIDIsNotFound(t *testing.T) {
	for _, c := range []struct{ name, id string }{
		{"nil", uuid.Nil.String()},
		{"never-issued", uuid.Must(uuid.NewV4()).String()},
	} {
		t.Run(c.name, func(t *testing.T) {
			h, provider := newDeleteConnectionFixture(t)

			rec := putConnection(h, provider, c.id)

			assert.Equal(t, http.StatusNotFound, rec.Code, rec.Body.String())
			assert.Zero(t, countConnections(t, provider), "an update must not create a connection")
		})
	}
}

// TestUpdateConnectionByIDExistingIsUpdated guards the fix's blast radius: the
// lookup added before the save must not stop a stored connection from updating.
func TestUpdateConnectionByIDExistingIsUpdated(t *testing.T) {
	h, provider := newDeleteConnectionFixture(t)
	saved, err := provider.ConnectionPersister.SaveConnection(&connections.Connection{
		Name: "original", Kind: "prometheus", ConnectionType: "observability",
	})
	require.NoError(t, err)

	rec := putConnection(h, provider, saved.ID.String())

	assert.Equal(t, http.StatusOK, rec.Code, rec.Body.String())
	got, err := provider.ConnectionPersister.GetConnection(saved.ID, "")
	require.NoError(t, err)
	assert.Equal(t, "renamed", got.Name)
	assert.Equal(t, int64(1), countConnections(t, provider), "an update must not create a connection")
}
