package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
)

// TestUpdateConnectionById_NonKubernetesKindRejected verifies that when a
// MeshSync deployment mode is specified in the metadata, the handler rejects
// non-kubernetes connection kinds with 400 Bad Request.
func TestUpdateConnectionById_NonKubernetesKindRejected(t *testing.T) {
	connectionID := uuid.Must(uuid.NewV4())
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	payload := connections.ConnectionPayload{
		Kind: "grafana",
		MetaData: map[string]interface{}{
			"meshsync_deployment_mode": "operator",
		},
	}
	body, _ := json.Marshal(payload)

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPut,
		"/api/integrations/connections/"+connectionID.String(),
		bytes.NewBuffer(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	req = mux.SetURLVars(req, map[string]string{"connectionId": connectionID.String()})
	rec := httptest.NewRecorder()

	h.UpdateConnectionById(rec, req, nil, user, p)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for non-kubernetes kind with meshsync mode, got %d (body=%q)",
			rec.Code, rec.Body.String())
	}
}
