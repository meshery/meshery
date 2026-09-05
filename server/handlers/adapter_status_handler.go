package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/meshery/meshery/server/models"
)

// changeAdapterStatusRequest is the REST request body for deploying or
// undeploying an adapter.
type changeAdapterStatusRequest struct {
	AdapterName  string `json:"adapterName"`
	TargetPort   string `json:"targetPort,omitempty"`
	TargetStatus string `json:"targetStatus"` // "enabled" (deploy) or "disabled" (undeploy)
}

// changeAdapterStatusResponse is the REST response body.
type changeAdapterStatusResponse struct {
	Status string `json:"status"`
}

// ChangeAdapterStatusHandler deploys or undeploys an adapter over REST.
// POST /api/system/adapter/status (auth required via ProviderMiddleware/AuthMiddleware).
// Request body: {adapterName, targetPort?, targetStatus: "enabled"|"disabled"}.
// Returns 400 for invalid or missing input, 401 if unauthenticated, 500 on
// deploy/undeploy failure. This mirrors the behavior of the GraphQL
// changeAdapterStatus mutation, including validation, default port
// resolution, and proper error propagation (see #19221).
func (h *Handler) ChangeAdapterStatusHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	if req.Method != http.MethodPost {
		writeMeshkitError(w, ErrMethodNotAllowed(req.Method), http.StatusMethodNotAllowed)
		return
	}

	var body changeAdapterStatusRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		h.log.Error(ErrRetrieveData(err))
		writeMeshkitError(w, ErrRetrieveData(err), http.StatusBadRequest)
		return
	}

	// adapterName is always required
	if body.AdapterName == "" {
		h.log.Error(ErrValidAdapter)
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return
	}

	targetPort := body.TargetPort

	// in case of empty target, prefer the default ports
	if targetPort == "" {
		h.log.Debug("target port is not provided, looking for default ports")
		selectedAdapter := getAdapterInformationByNameREST(body.AdapterName)
		if selectedAdapter == nil {
			h.log.Error(ErrValidAdapter)
			writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
			return
		}
		targetPort = selectedAdapter.Location
	}

	// reject malformed or out-of-range target ports before they reach the adapter tracker
	if port, err := strconv.Atoi(targetPort); err != nil || port < 1 || port > 65535 {
		h.log.Error(ErrValidAdapter)
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return
	}

	// reject any targetStatus value that isn't exactly "enabled" or "disabled"
	var deploy bool
	switch body.TargetStatus {
	case "enabled":
		deploy = true
	case "disabled":
		deploy = false
	default:
		h.log.Error(ErrValidAdapter)
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return
	}

	var operation string
	var err error
	adapter := models.Adapter{Name: body.AdapterName, Location: body.AdapterName + ":" + targetPort}
	if deploy {
		operation = "Deploy"
		h.log.Info("Deploying Adapter")
		err = h.config.AdapterTracker.DeployAdapter(req.Context(), adapter)
	} else {
		operation = "Undeploy"
		h.log.Info("Undeploying Adapter")
		err = h.config.AdapterTracker.UndeployAdapter(req.Context(), adapter)
	}

	if err != nil {
		h.log.Info("Failed to " + operation + " adapter")
		h.log.Error(err)
		writeMeshkitError(w, ErrRetrieveData(err), http.StatusInternalServerError)
		return
	}

	h.log.Info(operation + "ed adapter")

	err = json.NewEncoder(w).Encode(changeAdapterStatusResponse{Status: "processing"})
	if err != nil {
		obj := "data"
		if isClientDisconnect(err) {
			h.log.Debug(models.ErrMarshal(err, obj))
		} else {
			h.log.Error(models.ErrMarshal(err, obj))
		}
		return
	}
}

// getAdapterInformationByNameREST mirrors resolver.getAdapterInformationByName
// so the REST handler can resolve default ports without importing the
// resolver package.
func getAdapterInformationByNameREST(adapterName string) *models.Adapter {
	var adapter *models.Adapter
	for _, v := range models.ListAvailableAdapters {
		if adapterName == v.Name {
			adapter = &v
		}
	}
	return adapter
}
