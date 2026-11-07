package handlers

import (
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
