package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/models/meshmodel/entity"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newConnectionDefTestDB(t *testing.T) *database.Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		if strings.Contains(err.Error(), "requires cgo") {
			t.Skip("skipping database test: sqlite requires cgo")
		}
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	dbHandler := &database.Handler{DB: db, Mutex: &sync.Mutex{}}
	if err := dbHandler.AutoMigrate(
		&connectionv1beta3.ConnectionDefinition{},
		&registry.Registry{},
	); err != nil {
		t.Fatalf("failed to auto migrate connection definition and registry: %v", err)
	}
	return dbHandler
}

func TestUpdateConnectionDefinition(t *testing.T) {
	db := newConnectionDefTestDB(t)
	h := &Handler{config: &models.HandlerConfig{}, log: newTestLogger(t), dbHandler: db}
	id := uuid.Must(uuid.NewV4())
	put := func(id uuid.UUID, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPut, "/api/registry/connections/"+id.String(), strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"connectionDefinitionId": id.String()})
		rw := httptest.NewRecorder()
		h.UpdateConnectionDefinition(rw, req, nil, nil, nil)
		return rw
	}

	rw := put(id, `{"name":"ghost"}`)
	require.Equal(t, http.StatusNotFound, rw.Code, rw.Body.String())
	var count int64
	require.NoError(t, db.Model(&connectionv1beta3.ConnectionDefinition{}).Count(&count).Error)
	assert.Zero(t, count)

	original := connectionv1beta3.ConnectionDefinition{ID: id, Name: "original"}
	require.NoError(t, db.Create(&original).Error)
	require.NoError(t, db.First(&original, "id = ?", id).Error)
	for _, body := range []string{`{"name":"updated"}`, `{"name":"updated","createdAt":"2000-01-01T00:00:00Z"}`} {
		rw = put(id, body)
		require.Equal(t, http.StatusOK, rw.Code, rw.Body.String())
		var stored, response connectionv1beta3.ConnectionDefinition
		require.NoError(t, db.First(&stored, "id = ?", id).Error)
		require.NoError(t, json.Unmarshal(rw.Body.Bytes(), &response))
		assert.Equal(t, "updated", stored.Name)
		assert.True(t, original.CreatedAt.Equal(stored.CreatedAt))
		assert.True(t, original.CreatedAt.Equal(response.CreatedAt))
	}
}

// TestDeleteConnectionDefinition_InvalidUUID verifies that DeleteConnectionDefinition returns
// HTTP 400 Bad Request when an invalid UUID string is provided in the route parameters.
func TestDeleteConnectionDefinition_InvalidUUID(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{},
		log:    newTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/registry/connections/invalid-uuid", nil)
	req = mux.SetURLVars(req, map[string]string{
		"connectionDefinitionId": "invalid-uuid",
	})
	rw := httptest.NewRecorder()

	h.DeleteConnectionDefinition(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// TestDeleteConnectionDefinition_NotFound verifies that DeleteConnectionDefinition returns
// HTTP 404 Not Found when attempting to delete a connection definition that does not exist in the database.
func TestDeleteConnectionDefinition_NotFound(t *testing.T) {
	dbHandler := newConnectionDefTestDB(t)
	h := &Handler{
		config:    &models.HandlerConfig{},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
	}

	missingID := uuid.Must(uuid.NewV4())
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/registry/connections/%s", missingID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"connectionDefinitionId": missingID.String(),
	})
	rw := httptest.NewRecorder()

	h.DeleteConnectionDefinition(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "deleting non-existent connection definition must return HTTP 404")
	assert.Contains(t, rw.Body.String(), fmt.Sprintf("connection definition with id %s not found", missingID))
}

// TestDeleteConnectionDefinition_Success verifies that DeleteConnectionDefinition returns
// HTTP 204 No Content and successfully removes the definition and its registry entry from the database.
func TestDeleteConnectionDefinition_Success(t *testing.T) {
	dbHandler := newConnectionDefTestDB(t)
	h := &Handler{
		config:    &models.HandlerConfig{},
		log:       newTestLogger(t),
		dbHandler: dbHandler,
	}

	connDefID := uuid.Must(uuid.NewV4())
	connRecord := connectionv1beta3.ConnectionDefinition{
		ID:   connDefID,
		Name: "test-connection-def",
	}
	require.NoError(t, dbHandler.Create(&connRecord).Error)

	registryRecord := registry.Registry{
		Entity: connDefID,
		Type:   entity.ConnectionDefinition,
	}
	require.NoError(t, dbHandler.Create(&registryRecord).Error)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/registry/connections/%s", connDefID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{
		"connectionDefinitionId": connDefID.String(),
	})
	rw := httptest.NewRecorder()

	h.DeleteConnectionDefinition(rw, req, nil, nil, nil)
	resp := rw.Result()
	defer func() {
		_ = resp.Body.Close()
	}()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode, "deleting existing connection definition must return HTTP 204")

	var connCount int64
	dbHandler.Model(&connectionv1beta3.ConnectionDefinition{}).Where("id = ?", connDefID).Count(&connCount)
	assert.Equal(t, int64(0), connCount, "connection definition must be deleted")

	var regCount int64
	dbHandler.Model(&registry.Registry{}).Where("entity = ?", connDefID).Count(&regCount)
	assert.Equal(t, int64(0), regCount, "registry entry must be deleted")
}
