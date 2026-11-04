package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/registry"
	hostconn "github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteConnectionDefinition_NotFound(t *testing.T) {
	rm, db := newTestRegistryManager(t)
	log, err := logger.New("connection_definition_handler_test", logger.Options{})
	require.NoError(t, err)

	h := &Handler{
		registryManager: rm,
		dbHandler:       db,
		log:             log,
	}

	nonExistentID := "00000000-0000-0000-0000-000000000000"
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/registry/connections/%s", nonExistentID), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionDefinitionId": nonExistentID})
	rw := httptest.NewRecorder()

	h.DeleteConnectionDefinition(rw, req, nil, nil, nil)

	resp := rw.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteConnectionDefinition_Success(t *testing.T) {
	rm, db := newTestRegistryManager(t)

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

	log, err := logger.New("connection_definition_handler_test", logger.Options{})
	require.NoError(t, err)

	h := &Handler{
		registryManager: rm,
		dbHandler:       db,
		log:             log,
	}

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/registry/connections/%s", def.ID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"connectionDefinitionId": def.ID.String()})
	rw := httptest.NewRecorder()

	h.DeleteConnectionDefinition(rw, req, nil, nil, nil)

	resp := rw.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
