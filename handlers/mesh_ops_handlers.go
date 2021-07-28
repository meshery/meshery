//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

func init() {
	gob.Register([]*models.Adapter{})
}

// GetAllAdaptersHandler is used to fetch all the adapters
func (h *Handler) GetAllAdaptersHandler(w http.ResponseWriter, req *http.Request, provider models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err := json.NewEncoder(w).Encode(h.config.AdapterTracker.GetAdapters(req.Context()))
	if err != nil {
		obj := "data"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// MeshAdapterConfigHandler is used to persist adapter config
func (h *Handler) MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}
	var err error
	switch req.Method {
	case http.MethodPost:
		meshLocationURL := req.FormValue("meshLocationURL")

		h.log.Debug("meshLocationURL: ", meshLocationURL)
		if strings.TrimSpace(meshLocationURL) == "" {
			h.log.Error(ErrAddAdapter)
			http.Error(w, ErrAddAdapter.Error(), http.StatusBadRequest)
			return
		}

		if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
			h.log.Error(ErrInvalidK8SConfig)
			http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
			return
		}

		meshAdapters, err = h.addAdapter(req.Context(), meshAdapters, prefObj, meshLocationURL, provider)
		if err != nil {
			// h.log.Error(ErrRetrieveData(err))
			http.Error(w, ErrRetrieveData(err).Error(), http.StatusInternalServerError)
			return // error is handled appropriately in the relevant method
		}
	case http.MethodDelete:
		meshAdapters, err = h.deleteAdapter(meshAdapters, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	prefObj.MeshAdapters = meshAdapters
	err = provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(meshAdapters)
	if err != nil {
		obj := "data"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) addAdapter(ctx context.Context, meshAdapters []*models.Adapter, prefObj *models.Preference, meshLocationURL string, provider models.Provider) ([]*models.Adapter, error) {
	alreadyConfigured := false
	for _, adapter := range meshAdapters {
		if adapter.Location == meshLocationURL {
			// err := errors.New("Adapter with the given meshLocationURL already exists.")
			// h.log.Error(err)
			// http.Error(w, err.Error(), http.StatusForbidden)
			// return nil, err
			alreadyConfigured = true
			break
		}
	}

	if alreadyConfigured {
		h.log.Debug("Adapter already configured...")
		return meshAdapters, nil
	}
	if prefObj.K8SConfig == nil {
		h.log.Error(ErrInvalidK8SConfig)
		return nil, ErrInvalidK8SConfig
	}

	kubeclient, err := mesherykube.New(prefObj.K8SConfig.Config)
	if err != nil {
		h.log.Error(ErrNilClient)
		return nil, ErrNilClient
	}
	*h.config.KubeClient = *kubeclient
	provider.SetKubeClient(h.config.KubeClient)

	mClient, err := meshes.CreateClient(ctx, prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, meshLocationURL)
	if err != nil || prefObj.K8SConfig == nil {
		h.log.Error(ErrMeshClient)
		// http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
		return meshAdapters, ErrMeshClient
	}
	h.log.Debug("created client for adapter: ", meshLocationURL)
	defer func() {
		_ = mClient.Close()
	}()
	respOps, err := mClient.MClient.SupportedOperations(ctx, &meshes.SupportedOperationsRequest{})
	if err != nil {
		h.log.Error(ErrRetrieveMeshData(err))
		// http.Error(w, ErrRetrieveMeshData(err).Error(), http.StatusInternalServerError)
		return meshAdapters, err
	}
	h.log.Debug("retrieved supported ops for adapter: ", meshLocationURL)
	meshInfo, err := mClient.MClient.ComponentInfo(ctx, &meshes.ComponentInfoRequest{})
	if err != nil {
		h.log.Error(ErrRetrieveMeshData(err))
		// http.Error(w, ErrRetrieveMeshData(err).Error(), http.StatusInternalServerError)
		return meshAdapters, err
	}
	h.log.Debug("retrieved name for adapter: ", meshLocationURL)
	result := &models.Adapter{
		Location:     meshLocationURL,
		Name:         meshInfo.Name,
		Version:      meshInfo.Version,
		GitCommitSHA: meshInfo.GitSha,
		Ops:          respOps.GetOps(),
	}

	h.config.AdapterTracker.AddAdapter(ctx, *result)
	meshAdapters = append(meshAdapters, result)
	return meshAdapters, nil
}

func (h *Handler) deleteAdapter(meshAdapters []*models.Adapter, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, error) {
	adapterLoc := req.URL.Query().Get("adapter")
	h.log.Debug("URL of adapter to be removed: ", adapterLoc)

	adaptersLen := len(meshAdapters)

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		h.log.Error(ErrValidAdapter)
		http.Error(w, ErrValidAdapter.Error(), http.StatusBadRequest)
		return meshAdapters, ErrValidAdapter
	}

	newMeshAdapters := []*models.Adapter{}
	if aID == 0 {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[1:]...)
	} else if aID == adaptersLen-1 {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[:adaptersLen-1]...)
	} else {
		newMeshAdapters = append(newMeshAdapters, meshAdapters[0:aID]...)
		newMeshAdapters = append(newMeshAdapters, meshAdapters[aID+1:]...)
	}
	b, _ := json.Marshal(meshAdapters)
	h.log.Debug("Old adapters: ", b)
	b, _ = json.Marshal(newMeshAdapters)
	h.log.Debug("New adapters: ", b)
	return newMeshAdapters, nil
}

// MeshOpsHandler is used to send operations to the adapters
func (h *Handler) MeshOpsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if provider.GetProviderType() == models.RemoteProviderType {
		token, err := provider.GetProviderToken(req)

		if err == nil {
			viper.SetDefault("opt-token", token)
		}
	}

	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	adapterLoc := req.PostFormValue("adapter")
	h.log.Debug("Adapter URL to execute operations on: ", adapterLoc)

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		h.log.Error(ErrValidAdapter)
		http.Error(w, ErrValidAdapter.Error(), http.StatusBadRequest)
		return
	}

	opName := req.PostFormValue("query")
	customBody := req.PostFormValue("customBody")
	namespace := req.PostFormValue("namespace")
	delete := req.PostFormValue("deleteOp")
	if namespace == "" {
		namespace = "default"
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		h.log.Error(ErrInvalidK8SConfig)
		http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, meshAdapters[aID].Location)
	if err != nil {
		//h.log.Error(ErrMeshClient)
		http.Error(w, ErrMeshClient.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = mClient.Close()
	}()

	operationID, err := uuid.NewV4()

	if err != nil {
		h.log.Error(ErrOperationID(err))
		http.Error(w, ErrOperationID(err).Error(), http.StatusInternalServerError)
		return
	}

	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OperationId: operationID.String(),
		OpName:      opName,
		Username:    user.UserID,
		Namespace:   namespace,
		CustomBody:  customBody,
		DeleteOp:    (delete != ""),
	})
	if err != nil {
		h.log.Error(ErrApplyChange(err))
		http.Error(w, ErrApplyChange(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// AdapterPingHandler is used to ping a given adapter
func (h *Handler) AdapterPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	// adapterLoc := req.PostFormValue("adapter")
	adapterLoc := req.URL.Query().Get("adapter")
	h.log.Debug("Adapter url to ping: ", adapterLoc)

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		h.log.Error(ErrValidAdapter)
		http.Error(w, ErrValidAdapter.Error(), http.StatusBadRequest)
		return
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		h.log.Error(ErrInvalidK8SConfig)
		http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName, meshAdapters[aID].Location)
	if err != nil {
		h.log.Error(ErrMeshClient)
		http.Error(w, ErrMeshClient.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = mClient.Close()
	}()

	_, err = mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
	if err != nil {
		h.log.Error(ErrMeshClient)
		http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
