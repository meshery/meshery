package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta1/workspace"
)

// workspacePayloadWire is a handler-local dual-accept wrapper around
// workspace.WorkspacePayload. The schemas-generated struct tags
// OrganizationID as json:"organization_id" (required by the current
// published v1beta1 contract), but the canonical wire contract and every
// in-repo consumer now emit the camelCase `organizationId`. Go's
// encoding/json case-insensitive tag fallback does NOT match across an
// underscore boundary, so a struct tagged `organization_id` silently drops
// a JSON key of `organizationId`. This wrapper intercepts both spellings
// during the Phase 2 deprecation window. Canonical wins when both are
// present. Retire once schemas v1beta2 flips the tag to `organizationId`
// and this handler consumes the new version.
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
		http.Error(w, "failed to get token", http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
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
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	wire := workspacePayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	workspace := wire.WorkspacePayload
	bf, err := provider.SaveWorkspace(req, &workspace, "", false)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
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
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	wire := workspaceUpdatePayloadWire{}
	err = json.Unmarshal(bd, &wire)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	workspacePayload := wire.WorkspaceUpdatePayload
	resp, err := provider.UpdateWorkspace(req, &workspacePayload, workspaceID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, "Failed to marshal response to JSON", http.StatusInternalServerError)
		return
	}
	description := fmt.Sprintf("Workspace %s updated.", resp.Name)
	h.log.Info(description)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(respJSON)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetEnvironmentsOfWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	q := req.URL.Query()
	resp, err := provider.GetEnvironmentsOfWorkspace(req, workspaceID, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
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
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(resp); err != nil {
		h.log.Error(err)
	}
}
