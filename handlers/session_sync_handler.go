package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/models"
)

type SessionSyncData struct {
	*models.Preference `json:",inline"`
	K8sConfigs         []SessionSyncDataK8sConfig `json:"k8sConfig,omitempty"`
}

type SessionSyncDataK8sConfig struct {
	ContextID         string `json:"contextID,omitempty"`
	ContextName       string `json:"contextName,omitempty"`
	ClusterConfigured bool   `json:"clusterConfigured,omitempty"`
	ConfiguredServer  string `json:"configuredServer,omitempty"`
	ClusterID         string `json:"clusterID,omitempty"`
}

// swagger:route GET /api/system/sync SystemAPI idSystemSync
// Handle GET request for config sync
//
// Used to send session data to the UI for initial sync
// responses:
// 	200: userLoadTestPrefsRespWrapper

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
	s := []SessionSyncDataK8sConfig{}
	k8scontexts, ok := req.Context().Value(models.AllKubeClusterKey).([]models.K8sContext)
	if ok {
		for _, k8scontext := range k8scontexts {
			var cid string
			if k8scontext.KubernetesServerID != nil {
				cid = k8scontext.KubernetesServerID.String()
			}
			s = append(s, SessionSyncDataK8sConfig{
				ContextID:         k8scontext.ID,
				ContextName:       k8scontext.Name,
				ClusterConfigured: true,
				ClusterID:         cid,
				ConfiguredServer:  k8scontext.Server,
			})
		}
	}
	data := SessionSyncData{
		Preference: prefObj,
		K8sConfigs: s,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		obj := "user config data"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
