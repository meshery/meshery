package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/models"
)

// SessionSyncHandler is used to send session data to the UI for initial sync
func (h *Handler) SessionSyncHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	// To get fresh copy of User
	_, _ = provider.GetUserDetails(req)

	meshAdapters := []*models.Adapter{}

	adapters := h.config.AdapterTracker.GetAdapters(req.Context())
	for _, adapter := range adapters {
		meshAdapters, _ = h.addAdapter(req.Context(), meshAdapters, prefObj, adapter.Location, provider)
	}
	h.log.Debug("final list of active adapters: ", meshAdapters)
	prefObj.MeshAdapters = meshAdapters
	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil { // ignoring errors in this context
		h.log.Error(ErrSaveSession(err))
	}

	if prefObj.K8SConfig != nil && h.config.KubeClient != nil {
		if prefObj.K8SConfig.ServerVersion == "" {
			// fetching server version, if it has not already been
			version, err := h.config.KubeClient.KubeClient.ServerVersion()
			if err != nil {
				h.log.Error(ErrFetchKubernetes(err))
			}
			prefObj.K8SConfig.ServerVersion = version.String()
		}

		//if len(prefObj.K8SConfig.Nodes) == 0 {
		//	// fetching nodes, if it has not already been
		//	prefObj.K8SConfig.Nodes, _ = helpers.FetchKubernetesNodes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
		//}

		// clearing out the config just for displaying purposes
		if len(prefObj.K8SConfig.Config) > 0 {
			prefObj.K8SConfig.Config = nil
		}
	} else {
		err = h.checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req, user, prefObj, provider)
		if err != nil {
			h.log.Error(ErrFetchKubernetes(err))
		}
	}

	err = json.NewEncoder(w).Encode(prefObj)
	if err != nil {
		obj := "user config data"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
