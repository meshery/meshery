package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/meshery/schemas/models/core"
	systemv1beta1 "github.com/meshery/schemas/models/v1beta1/system"

	"github.com/meshery/meshery/server/models"
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
	err := provider.RecordPreferences(req, user.UserId, prefObj)
	if err != nil { // ignoring errors in this context
		h.log.Error(ErrSaveSession(err))
	}
	k8sConfigs := []systemv1beta1.SystemSessionSyncK8sContext{}
	k8scontexts, ok := req.Context().Value(models.AllKubeClusterKey).([]*models.K8sContext)
	if ok {
		for _, k8scontext := range k8scontexts {
			if k8scontext == nil {
				continue
			}
			clusterConfigured := true
			contextID := k8scontext.ID
			contextName := k8scontext.Name
			server := k8scontext.Server
			context := systemv1beta1.SystemSessionSyncK8sContext{
				ID:                &contextID,
				Name:              &contextName,
				ClusterConfigured: &clusterConfigured,
				ClusterId:         k8scontext.KubernetesServerID,
				Server:            &server,
			}
			if k8scontext.CreatedAt != nil {
				context.CreatedAt = core.CreatedAt(*k8scontext.CreatedAt)
			}
			if k8scontext.UpdatedAt != nil {
				context.UpdatedAt = core.UpdatedAt(*k8scontext.UpdatedAt)
			}
			k8sConfigs = append(k8sConfigs, context)
		}
	}
	data := systemv1beta1.SystemSessionSync{
		K8sConfig: &k8sConfigs,
	}
	prefData, err := json.Marshal(prefObj)
	if err != nil {
		obj := "user config data"
		h.log.Error(models.ErrMarshal(err, obj))
		writeMeshkitError(w, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
if err := json.Unmarshal(prefData, &data.AdditionalProperties); err != nil {
	obj := "user config data"
	h.log.Error(models.ErrUnmarshal(err, obj))
	writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusInternalServerError)
	return
}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		obj := "user config data"
		h.log.Error(models.ErrMarshal(err, obj))
		writeMeshkitError(w, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
}
