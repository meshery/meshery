package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/workspaces WorkspacesAPI idGetWorkspaces
// Handles GET for all Workspaces
//
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 20
//
// ```?search={name}``` If search is non empty then a greedy search is performed
//
// ```?orgID={orgid}``` orgID is used to retrieve workspaces belonging to a particular org *required*
//
// ```?filter={condition}```
// responses:
// 	200: workspacesResponseWrapper

func (h *Handler) GetWorkspacesHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	resp, err := provider.GetWorkspaces(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"), q.Get("orgID"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/workspaces/{id} WorkspacesAPI idGetWorkspacesByIdHandler
// Handle GET for Workspace info by ID
//
// ```?orgID={orgid}``` orgID is used to retrieve workspaces belonging to a particular org
//
// Returns Workspace info
// responses:
//   200: workspaceResponseWrapper

func (h *Handler) GetWorkspaceByIdHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(r)["id"]
	q := r.URL.Query()
	resp, err := provider.GetWorkspaceByID(r, workspaceID, q.Get("orgID"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route POST /api/workspaces PostWorkspace idSaveWorkspace
// Handle POST request for creating a new workspace
//
// Creates a new workspace
// responses:
// 201: noContentWrapper
func (h *Handler) SaveWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	workspace := models.WorkspacePayload{}
	err = json.Unmarshal(bd, &workspace)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	err = provider.SaveWorkspace(req, &workspace, "", false)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	description := fmt.Sprintf("Workspace %s created.", workspace.Name)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
}

// swagger:route DELETE /api/workspaces/{id} WorkspaceAPI idDeleteWorkspaceHandler
// Handle DELETE for Workspace based on ID
//
// Deletes a workspace
// responses:
// 201: noContentWrapper

func (h *Handler) DeleteWorkspaceHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	workspaceID := mux.Vars(r)["id"]
	resp, err := provider.DeleteWorkspace(r, workspaceID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route PUT /api/workspaces/{id} PostWorkspace idUpdateWorkspaceHandler
// Handle PUT request for updating a workspace
//
// Updates a workspace
// responses:
//
//	200: workspaceResponseWrapper
func (h *Handler) UpdateWorkspaceHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	workspaceID := mux.Vars(req)["id"]
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	workspace := models.WorkspacePayload{}
	err = json.Unmarshal(bd, &workspace)
	obj := "workspace"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	resp, err := provider.UpdateWorkspace(req, &workspace, workspaceID)
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
	description := fmt.Sprintf("Workspace %s updated.", workspace.Name)
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
