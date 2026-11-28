package handlers

import (
	"net/http"
	"time"

	"encoding/json"

	"github.com/meshery/meshery/server/models"
)

type SessionSyncData struct {
	*models.Preference `json:",inline"`
	K8sConfigs         []SessionSyncDataK8sConfig `json:"k8sConfig,omitempty"`
}

type SessionSyncDataK8sConfig struct {
	ContextID         string     `json:"id,omitempty"`
	ContextName       string     `json:"name,omitempty"`
	ClusterConfigured bool       `json:"clusterConfigured,omitempty"`
	ConfiguredServer  string     `json:"server,omitempty"`
	ClusterID         string     `json:"clusterId,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

// SessionSyncHandler is used to send session data to the UI for initial sync
func (h *Handler) SessionSyncHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}

	// To get fresh copy of User
	_, _ = provider.GetUserDetails(req)

	if prefObj == nil {
		var err error
		prefObj, err = provider.ReadFromPersister(user.ID.String())
		if err != nil || prefObj == nil {
			if h.log != nil {
				h.log.Error(ErrReadSessionPersistor)
			}
			writeMeshkitError(w, ErrReadSessionPersistor, http.StatusInternalServerError)
			return
		}
	}

	if h.config != nil && h.config.AdapterTracker != nil {
		meshAdapters := []*models.Adapter{}
		adapters := h.config.AdapterTracker.GetAdapters(req.Context())
		for _, adapter := range adapters {
			meshAdapters, _ = h.addAdapter(req.Context(), meshAdapters, prefObj, adapter.Location, provider)
		}
		if h.log != nil {
			h.log.Debug("final list of active adapters: ", meshAdapters)
		}
		prefObj.MeshAdapters = meshAdapters
		err := provider.RecordPreferences(req, user.ID.String(), prefObj)
		if err != nil { // ignoring errors in this context
			if h.log != nil {
				h.log.Error(ErrSaveSession(err))
			}
		}
	}
	s := []SessionSyncDataK8sConfig{}
	k8scontexts, ok := req.Context().Value(models.AllKubeClusterKey).([]*models.K8sContext)
	if ok {
		for _, k8scontext := range k8scontexts {
			if k8scontext == nil {
				continue
			}
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
				CreatedAt:         k8scontext.CreatedAt,
				UpdatedAt:         k8scontext.UpdatedAt,
			})
		}
	}
	data := SessionSyncData{
		Preference: prefObj,
		K8sConfigs: s,
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		obj := "user config data"
		if h.log != nil {
			h.log.Error(models.ErrMarshal(err, obj))
		}
		writeMeshkitError(w, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
}
