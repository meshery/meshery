package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// SessionSyncHandler is used to send session data to the UI for initial sync
func (h *Handler) SessionSyncHandler(w http.ResponseWriter, req *http.Request, _ *sessions.Session, sessObj *models.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := h.checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req, user, sessObj)
	// if err != nil {
	// // We can ignore the errors here. They are logged in the other method
	// }

	// meshAdapters := sessObj.MeshAdapters
	// if meshAdapters == nil {
	// meshAdapters = []*models.Adapter{}
	// }

	// this is just called for getting a fresh copy of preferences
	h.config.Provider.GetUserDetails(req)

	meshAdapters := []*models.Adapter{}

	adapters := h.config.AdapterTracker.GetAdapters(req.Context())
	for _, adapterURL := range adapters {
		meshAdapters, _ = h.addAdapter(req.Context(), meshAdapters, sessObj, adapterURL)
	}
	logrus.Debugf("final list of active adapters: %+v", meshAdapters)
	sessObj.MeshAdapters = meshAdapters
	err = h.config.Provider.RecordPreferences(req, user.UserID, sessObj)
	if err != nil { // ignoring errors in this context
		logrus.Errorf("unable to save session: %v", err)
		// http.Error(w, "unable to save session", http.StatusInternalServerError)
		// return
	}

	if sessObj.K8SConfig != nil {
		if sessObj.K8SConfig.ServerVersion == "" {
			// fetching server version, if it has not already been
			sessObj.K8SConfig.ServerVersion, _ = helpers.FetchKubernetesVersion(sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName)
			// if err != nil {
			// 	http.Error(w, "unable to ping the kubernetes server", http.StatusInternalServerError)
			// 	return
			// }
		}

		if len(sessObj.K8SConfig.Nodes) == 0 {
			// fetching nodes, if it has not already been
			sessObj.K8SConfig.Nodes, _ = helpers.FetchKubernetesNodes(sessObj.K8SConfig.Config, sessObj.K8SConfig.ContextName)
			// if err != nil {
			// 	http.Error(w, "unable to fetch nodes metadata from the kubernetes server", http.StatusInternalServerError)
			// 	return
			// }
		}

		// clearing out the config just for displaying purposes
		if len(sessObj.K8SConfig.Config) > 0 {
			sessObj.K8SConfig.Config = nil
		}
	}

	err = json.NewEncoder(w).Encode(sessObj)
	if err != nil {
		logrus.Errorf("error marshalling user config data: %v", err)
		http.Error(w, "unable to process the request", http.StatusInternalServerError)
		return
	}
}
