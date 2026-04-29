package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/meshery/meshery/server/models"
	workspace "github.com/meshery/schemas/models/v1beta3/workspace"
)

// workspacePayloadWire is a handler-local dual-accept wrapper around
// workspace.WorkspacePayload (now v1beta3, canonical-camelCase). The
// canonical wire form emits `organizationId`; the legacy `organization_id`
// spelling is still accepted for the Phase 5 deprecation window so any
// unmigrated client (mesheryctl, older Kanvas releases) keeps working.
// Go's encoding/json case-insensitive tag fallback does NOT match across
// an underscore boundary, so the legacy spelling cannot simply piggy-back
// on the canonical tag. Canonical wins when both are present. Retire once
// every known consumer is on the canonical spelling.
type workspacePayloadWire struct {
	workspace.WorkspacePayload
}

func (p *workspacePayloadWire) UnmarshalJSON(data []byte) error {
	type alias workspace.WorkspacePayload
	aux := struct {
		*alias
		OrganizationIDCamel *openapi_types.UUID `json:"organizationId,omitempty"`
		OrganizationIDSnake *openapi_types.UUID `json:"organization_id,omitempty"`
	}{alias: (*alias)(&p.WorkspacePayload)}

	// Zero OrganizationID so a reused receiver does not carry stale data
	// when the next payload omits both spellings.
	p.OrganizationID = openapi_types.UUID{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Canonical wins when both are supplied.
	switch {
	case aux.OrganizationIDCamel != nil:
		p.OrganizationID = *aux.OrganizationIDCamel
	case aux.OrganizationIDSnake != nil:
		p.OrganizationID = *aux.OrganizationIDSnake
	}
	return nil
}

// workspaceUpdatePayloadWire mirrors workspacePayloadWire for the PUT path.
type workspaceUpdatePayloadWire struct {
	workspace.WorkspaceUpdatePayload
}

func (p *workspaceUpdatePayloadWire) UnmarshalJSON(data []byte) error {
	type alias workspace.WorkspaceUpdatePayload
	aux := struct {
		*alias
		OrganizationIDCamel *openapi_types.UUID `json:"organizationId,omitempty"`
		OrganizationIDSnake *openapi_types.UUID `json:"organization_id,omitempty"`
	}{alias: (*alias)(&p.WorkspaceUpdatePayload)}

	p.OrganizationID = openapi_types.UUID{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	switch {
	case aux.OrganizationIDCamel != nil:
		p.OrganizationID = *aux.OrganizationIDCamel
	case aux.OrganizationIDSnake != nil:
		p.OrganizationID = *aux.OrganizationIDSnake
	}
	return nil
}

func (h *Handler) GetWorkspacesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		writeMeshkitError(w, ErrFetchToken(fmt.Errorf("token not found in request context")), http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	// Canonical form is `orgId`; `orgID` is dual-accepted during the Phase 2
	// deprecation window so mesheryctl and any other legacy client keeps
	// working. Retire the fallback once Phase 3 consumer migration completes.
	orgID := q.Get("orgId")
	if orgID == "" {
		orgID = q.Get("orgID")
	}
	if orgID == "" {
		missingInput := models.ErrWorkspaceMissingInput()
		h.log.Error(missingInput)
		writeJSONError(w, missingInput.Error(), http.StatusBadRequest)
		return
	}
	resp, err := provider.GetWorkspaces(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), orgID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetWorkspaceByIdHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(r)["id"]
	q := r.URL.Query()
	// Canonical form is `orgId`; `orgID` is dual-accepted during the Phase 2
	// deprecation window. Retire once Phase 3 consumer migration completes.
	orgID := q.Get("orgId")
	if orgID == "" {
		orgID = q.Get("orgID")
	}
	if orgID == "" {
		missingInput := models.ErrWorkspaceMissingInput()
		h.log.Error(missingInput)
		writeJSONError(w, missingInput.Error(), http.StatusBadRequest)
		return
	}
	resp, err := provider.GetWorkspaceByID(r, workspaceID, orgID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) SaveWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	wire := workspacePayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	workspace := wire.WorkspacePayload
	bf, err := provider.SaveWorkspace(req, &workspace, "", false)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	description := fmt.Sprintf("Workspace %s created.", workspace.Name)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(bf); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) DeleteWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(r)["id"]
	resp, err := provider.DeleteWorkspace(r, workspaceID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) UpdateWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	wire := workspaceUpdatePayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
		return
	}

	workspacePayload := wire.WorkspaceUpdatePayload
	resp, err := provider.UpdateWorkspace(req, &workspacePayload, workspaceID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound)
		return
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, models.ErrMarshal(err, "workspace update response"), http.StatusInternalServerError)
		return
	}
	description := fmt.Sprintf("Workspace %s updated.", resp.Name)
	h.log.Info(description)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respJSON)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		// Headers already committed; log only. Writing another body would corrupt the stream.
		return
	}
}

func (h *Handler) GetEnvironmentsOfWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	q := req.URL.Query()
	resp, err := provider.GetEnvironmentsOfWorkspace(req, workspaceID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetDesignsOfWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	q := req.URL.Query()
	resp, err := provider.GetDesignsOfWorkspace(req, workspaceID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), q["visibility"])
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) AddEnvironmentToWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	environmentID := mux.Vars(req)["environmentID"]
	resp, err := provider.AddEnvironmentToWorkspace(req, workspaceID, environmentID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) RemoveEnvironmentFromWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	environmentID := mux.Vars(req)["environmentID"]
	resp, err := provider.RemoveEnvironmentFromWorkspace(req, workspaceID, environmentID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) AddDesignToWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	designID := mux.Vars(req)["designID"]
	resp, err := provider.AddDesignToWorkspace(req, workspaceID, designID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) RemoveDesignFromWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	designID := mux.Vars(req)["designID"]
	resp, err := provider.RemoveDesignFromWorkspace(req, workspaceID, designID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetViewsOfWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	q := req.URL.Query()
	resp, err := provider.GetViewsOfWorkspace(req, workspaceID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) AddViewToWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	viewID := mux.Vars(req)["viewID"]
	resp, err := provider.AddViewToWorkspace(req, workspaceID, viewID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) RemoveViewFromWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	viewID := mux.Vars(req)["viewID"]
	resp, err := provider.RemoveViewFromWorkspace(req, workspaceID, viewID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) GetTeamsOfWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	q := req.URL.Query()
	resp, err := provider.GetTeamsOfWorkspace(req, workspaceID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) AddTeamToWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	teamID := mux.Vars(req)["teamID"]
	resp, err := provider.AddTeamToWorkspace(req, workspaceID, teamID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}

func (h *Handler) RemoveTeamFromWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	teamID := mux.Vars(req)["teamID"]
	resp, err := provider.RemoveTeamFromWorkspace(req, workspaceID, teamID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		writeMeshkitError(w, ErrGetResult(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}
