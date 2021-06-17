package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
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
	logrus.Debugf("final list of active adapters: %+v", meshAdapters)
	prefObj.MeshAdapters = meshAdapters
	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil { // ignoring errors in this context
		logrus.Errorf("unable to save session: %v", err)
	}

	if prefObj.K8SConfig != nil && h.kubeclient != nil {
		if prefObj.K8SConfig.ServerVersion == "" {
			// fetching server version, if it has not already been
			version, err := h.kubeclient.KubeClient.ServerVersion()
			if err != nil {
				logrus.Errorf("unable to ping the Kubernetes server")
			}
			prefObj.K8SConfig.ServerVersion = version.String()
		}

		//if len(prefObj.K8SConfig.Nodes) == 0 {
		//	// fetching nodes, if it has not already been
		//	prefObj.K8SConfig.Nodes, _ = helpers.FetchKubernetesNodes(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
		//}

		//// clearing out the config just for displaying purposes
		//if len(prefObj.K8SConfig.Config) > 0 {
		//	prefObj.K8SConfig.Config = nil
		//}
	} else {
		err = h.checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req, user, prefObj, provider)
		if err != nil {
			logrus.Errorf("unable to initialize kubernetes config")
		}
	}

	err = json.NewEncoder(w).Encode(prefObj)
	if err != nil {
		logrus.Errorf("error marshaling user config data: %v", err)
		http.Error(w, "unable to process the request", http.StatusInternalServerError)
		return
	}
}
