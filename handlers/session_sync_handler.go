package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
)

type SessionSyncData struct {
	*models.Preference `json:",inline"`
	K8sConfig          SessionSyncDataK8sConfig `json:"k8sConfig,omitempty"`
}

type SessionSyncDataK8sConfig struct {
	K8sFile           string `json:"k8sfile,omitempty"`
	ContextName       string `json:"contextName,omitempty"`
	ClusterConfigured bool   `json:"clusterConfigured,omitempty"`
	ConfiguredServer  string `json:"configuredServer,omitempty"`
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

		// Get the k8sconfig
		k8scontexts, ok := req.Context().Value(models.KubeContextKey).([]models.K8sContext)
		if !ok || len(k8scontexts) == 0 {
			h.log.Error(ErrInvalidK8SConfig)
			http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
			return
		}
		var mClient *meshes.MeshClient
		for _, mk8scontext := range k8scontexts {
			k8sconfig, err := mk8scontext.GenerateKubeConfig()
			if err != nil {
				continue
				// h.log.Error(ErrInvalidK8SConfig)
				// return nil, ErrInvalidK8SConfig
			}

			mClient, err = meshes.CreateClient(req.Context(), k8sconfig, mk8scontext.Name, adapter.Location)
			if err != nil {
				continue
				// http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
				// return
			}
		}
		if mClient == nil {
			http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
			return
		}
		meshAdapters, _ = h.addAdapter(req.Context(), meshAdapters, prefObj, adapter.Location, provider, mClient)
	}
	h.log.Debug("final list of active adapters: ", meshAdapters)
	prefObj.MeshAdapters = meshAdapters
	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil { // ignoring errors in this context
		h.log.Error(ErrSaveSession(err))
	}

	// Get the kubernetes context
	mk8scontext, ok := req.Context().Value(models.KubeContextKey).(*models.K8sContext)
	if !ok || mk8scontext == nil {
		h.log.Error(ErrInvalidK8SConfig)
		http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
		return
	}

	// Get the k8sconfig
	k8sconfig, ok := req.Context().Value(models.KubeConfigKey).([]byte)
	if !ok || k8sconfig == nil {
		h.log.Error(ErrInvalidK8SConfig)
		http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
		return
	}

	data := SessionSyncData{
		Preference: prefObj,
		K8sConfig: SessionSyncDataK8sConfig{
			ContextName:       mk8scontext.Name,
			ClusterConfigured: true,
			ConfiguredServer:  mk8scontext.Server,
		},
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		obj := "user config data"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}
