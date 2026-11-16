package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func newWorkspaceNilGuardsFixture(t *testing.T) *Handler {
	t.Helper()
	return &Handler{
		log: newTestLogger(t),
	}
}

func TestGetWorkspacesHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)

	ctx := context.WithValue(context.Background(), models.TokenCtxKey, "fake-token")
	req := httptest.NewRequest(http.MethodGet, "/api/workspaces?orgId=11111111-1111-1111-1111-111111111111", nil).WithContext(ctx)
	rec := httptest.NewRecorder()

	h.GetWorkspacesHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetWorkspaceByIdHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+workspaceID+"?orgId=11111111-1111-1111-1111-111111111111", nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.GetWorkspaceByIdHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestSaveWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)

	body := `{"name":"test-workspace","organizationId":"11111111-1111-1111-1111-111111111111"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workspaces", strings.NewReader(body))
	rec := httptest.NewRecorder()

	h.SaveWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDeleteWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/workspaces/"+workspaceID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.DeleteWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUpdateWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	body := `{"name":"updated-workspace"}`
	req := httptest.NewRequest(http.MethodPut, "/api/workspaces/"+workspaceID, strings.NewReader(body))
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.UpdateWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetEnvironmentsOfWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+workspaceID+"/environments", nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.GetEnvironmentsOfWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetDesignsOfWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+workspaceID+"/designs", nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.GetDesignsOfWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestAddEnvironmentToWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	envID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodPost, "/api/workspaces/"+workspaceID+"/environments/"+envID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "environmentID": envID})
	rec := httptest.NewRecorder()

	h.AddEnvironmentToWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestRemoveEnvironmentFromWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	envID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/workspaces/"+workspaceID+"/environments/"+envID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "environmentID": envID})
	rec := httptest.NewRecorder()

	h.RemoveEnvironmentFromWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestAddDesignToWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	designID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodPost, "/api/workspaces/"+workspaceID+"/designs/"+designID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "designID": designID})
	rec := httptest.NewRecorder()

	h.AddDesignToWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestRemoveDesignFromWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	designID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/workspaces/"+workspaceID+"/designs/"+designID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "designID": designID})
	rec := httptest.NewRecorder()

	h.RemoveDesignFromWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetViewsOfWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+workspaceID+"/views", nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.GetViewsOfWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestAddViewToWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	viewID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodPost, "/api/workspaces/"+workspaceID+"/views/"+viewID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "viewID": viewID})
	rec := httptest.NewRecorder()

	h.AddViewToWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestRemoveViewFromWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	viewID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/workspaces/"+workspaceID+"/views/"+viewID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "viewID": viewID})
	rec := httptest.NewRecorder()

	h.RemoveViewFromWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetTeamsOfWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodGet, "/api/workspaces/"+workspaceID+"/teams", nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID})
	rec := httptest.NewRecorder()

	h.GetTeamsOfWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestAddTeamToWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodPost, "/api/workspaces/"+workspaceID+"/teams/"+teamID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "teamID": teamID})
	rec := httptest.NewRecorder()

	h.AddTeamToWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestRemoveTeamFromWorkspaceHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newWorkspaceNilGuardsFixture(t)
	workspaceID := uuid.Must(uuid.NewV4()).String()
	teamID := uuid.Must(uuid.NewV4()).String()

	req := httptest.NewRequest(http.MethodDelete, "/api/workspaces/"+workspaceID+"/teams/"+teamID, nil)
	req = mux.SetURLVars(req, map[string]string{"id": workspaceID, "teamID": teamID})
	rec := httptest.NewRecorder()

	h.RemoveTeamFromWorkspaceHandler(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}
