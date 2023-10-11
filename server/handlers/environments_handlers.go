package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/integrations/environments EnvironmentsAPI idGetEnvironments
// Handles GET for all Environments
//
// # Environments can be further filtered through query parameters
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 20
//
// ```?search={environments_name}``` If search is non empty then a greedy search is performed
//
// ```?filter={condition}```
// responses:
// 	200: environmentsResponseWrapper

func (h *Handler) GetEnvironments(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	resp, err := provider.GetEnvironments(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("filter"))
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route GET /api/integrations/environments/{id} EnvironmentAPI idGetEnvironmentByIDHandler
// Handle GET for Environment info by ID
//
// Returns Environment info
// responses:
//   200: environmentResponseWrapper

func (h *Handler) GetEnvironmentByIDHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["id"]
	resp, err := provider.GetEnvironmentByID(r, environmentID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route POST /api/integrations/environments PostEnvironment idSaveEnvironment
// Handle POST request for creating a new environment
//
// Creates a new environment
// responses:
// 201: noContentWrapper
func (h *Handler) SaveEnvironment(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	environment := models.EnvironmentPayload{}
	err = json.Unmarshal(bd, &environment)
	obj := "environment"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	err = provider.SaveEnvironment(req, &environment, "", false)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	description := fmt.Sprintf("Environment %s created.", environment.Name)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
}

// swagger:route DELETE /api/integrations/environments/{id} EnvironmentAPI idDeleteEnvironmentHandler
// Handle DELETE for Environment based on ID
//
// Returns Environment info
// responses:
// 201: noContentWrapper

func (h *Handler) DeleteEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["id"]
	resp, err := provider.DeleteEnvironment(r, environmentID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route PUT /api/integrations/environments/{id} PostEnvironment idUpdateEnvironmentHandler
// Handle PUT request for updating a environment
//
// Updates a environment
// responses:
//
//	200: environmentResponseWrapper
func (h *Handler) UpdateEnvironmentHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	environmentID := mux.Vars(req)["id"]
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	environment := models.EnvironmentPayload{}
	err = json.Unmarshal(bd, &environment)
	obj := "environment"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	resp, err := provider.UpdateEnvironment(req, &environment, environmentID)
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
	description := fmt.Sprintf("Environment %s updated.", environment.Name)
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

// swagger:route POST /api/integrations/environments/{environmentID}/connections/{connectionID} EnvironmentAPI idAddConnectionToEnvironmentHandler
// Handle POST to add Connection in Environment by ID
//
// Returns Environment connection mapping info
// responses:
//   200: environmentResponseWrapper

func (h *Handler) AddConnectionToEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["environmentID"]
	connectionID := mux.Vars(r)["connectionID"]
	resp, err := provider.AddConnectionToEnvironment(r, environmentID, connectionID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}

// swagger:route DELETE /api/integrations/environments/{environmentID}/connections/{connectionID} EnvironmentAPI idRemoveConnectionFromEnvironmentHandler
// Handle DELETE to remove Connection from Environment by ID
//
// Removes connection from environment
// responses:
// 201: noContentWrapper

func (h *Handler) RemoveConnectionFromEnvironmentHandler(w http.ResponseWriter, r *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	environmentID := mux.Vars(r)["environmentID"]
	connectionID := mux.Vars(r)["connectionID"]
	resp, err := provider.RemoveConnectionFromEnvironment(r, environmentID, connectionID)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(resp))
}
