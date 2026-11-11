package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	hostconn "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/stretchr/testify/require"
)

// seedConnectionDefinition registers one connection definition through the
// registry (the same path RegisterConnectionDefinition uses) and returns it.
func seedConnectionDefinition(t *testing.T, rm *registry.RegistryManager) connectionv1beta3.ConnectionDefinition {
	t.Helper()

	registrant := hostconn.Connection{
		Name:    "test-registrant",
		Kind:    "kubernetes",
		Type:    "platform",
		SubType: "orchestration",
		Status:  hostconn.ConnectionStatusConnected,
	}

	def := connectionv1beta3.ConnectionDefinition{
		Name:           "Kubernetes",
		Kind:           "kubernetes",
		ConnectionType: "platform",
		SubType:        "orchestration",
		Status:         connectionv1beta3.ConnectionStatusDiscovered,
		ModelReference: &model.ModelReference{
			Name:    "kubernetes",
			Version: "v1.25.0",
		},
	}
	id, err := def.GenerateID()
	require.NoError(t, err)
	def.ID = id

	_, _, err = rm.RegisterEntity(registry.RegistrantHostToV1beta3(registrant), &def)
	require.NoError(t, err, "RegisterEntity failed")
	return def
}

func callUpdate(t *testing.T, h *Handler, id string, body connectionv1beta3.ConnectionDefinition) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/registry/connections/"+id, bytes.NewReader(payload))
	req = mux.SetURLVars(req, map[string]string{"connectionDefinitionId": id})
	rr := httptest.NewRecorder()

	h.UpdateConnectionDefinition(rr, req, nil, nil, nil)
	return rr
}

// TestUpdateConnectionDefinitionUnknownIDDoesNotCreate guards the first half of
// issue #23070: gorm's Save() upserts, so a PUT to an id that matches no row
// used to silently insert a new definition. UpdateConnectionDefinition must
// instead return 404 and leave the table untouched.
func TestUpdateConnectionDefinitionUnknownIDDoesNotCreate(t *testing.T) {
	rm, db := newTestRegistryManager(t)
	h := &Handler{dbHandler: db, log: newTestLogger(t)}

	unknownID := uuid.Must(uuid.NewV4()).String()

	rr := callUpdate(t, h, unknownID, connectionv1beta3.ConnectionDefinition{
		Name: "ghost",
		Kind: "kubernetes",
	})

	require.Equal(t, http.StatusNotFound, rr.Code, "unknown id must 404, not upsert")

	var count int64
	require.NoError(t, db.Model(&connectionv1beta3.ConnectionDefinition{}).
		Where("id = ?", unknownID).Count(&count).Error)
	require.EqualValues(t, 0, count, "update of an unknown id must not create a row")

	_ = rm
}

// TestUpdateConnectionDefinitionPreservesCreatedAt guards the second half of
// issue #23070: the request body does not carry created_at, so a plain Save()
// overwrote the stored timestamp with the zero value. The update must preserve
// created_at while still applying the other fields.
func TestUpdateConnectionDefinitionPreservesCreatedAt(t *testing.T) {
	rm, db := newTestRegistryManager(t)
	h := &Handler{dbHandler: db, log: newTestLogger(t)}

	def := seedConnectionDefinition(t, rm)

	var before connectionv1beta3.ConnectionDefinition
	require.NoError(t, db.First(&before, "id = ?", def.ID).Error)
	require.False(t, before.CreatedAt.IsZero(), "seeded row should have a created_at")

	// Body mutates Name and carries no createdAt (the zero value).
	update := def
	update.Name = "Kubernetes (renamed)"
	update.CreatedAt = connectionv1beta3.ConnectionDefinition{}.CreatedAt // zero

	rr := callUpdate(t, h, def.ID.String(), update)
	require.Equal(t, http.StatusOK, rr.Code)

	var after connectionv1beta3.ConnectionDefinition
	require.NoError(t, db.First(&after, "id = ?", def.ID).Error)

	require.Equal(t, "Kubernetes (renamed)", after.Name, "update should apply the new name")
	require.False(t, after.CreatedAt.IsZero(), "created_at must not be reset to zero")
	require.True(t, after.CreatedAt.Equal(before.CreatedAt),
		"created_at must be preserved: before=%s after=%s", before.CreatedAt, after.CreatedAt)
}
