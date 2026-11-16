// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/viper"
)

func init() {
	gob.Register([]*models.Adapter{})
}

// AvailableAdaptersHandler returns the list of available adapters
func (h *Handler) AvailableAdaptersHandler(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(models.ListAvailableAdapters)
	if err != nil {
		obj := "data"
		// Encode streams to w; failure here means the body has already
		// started — log only, don't double-write.
		if h.log != nil {
			if isClientDisconnect(err) {
				h.log.Debug(models.ErrMarshal(err, obj))
			} else {
				h.log.Error(models.ErrMarshal(err, obj))
			}
		}
		return
	}
}

// AdaptersHandler is used to fetch all the adapters
func (h *Handler) AdaptersHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// if adapter found in query user is trying to ping an adapter
	adapterLoc := req.URL.Query().Get("adapter")
	if adapterLoc != "" {
		if h.log != nil {
			h.log.Debug("adapter pinging")
		}
		h.AdapterPingHandler(w, req, prefObj, user, provider)
		return
	}

	var adapters []models.Adapter
	if h.config != nil && h.config.AdapterTracker != nil {
		adapters = h.config.AdapterTracker.GetAdapters(req.Context())
	}
	if adapters == nil {
		adapters = []models.Adapter{}
	}

	err := json.NewEncoder(w).Encode(adapters)
	if err != nil {
		obj := "data"
		// Streaming encode failure — log only.
		if h.log != nil {
			if isClientDisconnect(err) {
				h.log.Debug(models.ErrMarshal(err, obj))
			} else {
				h.log.Error(models.ErrMarshal(err, obj))
			}
		}
		return
	}
}

// AdapterPingHandler is used to ping a given adapter
func (h *Handler) AdapterPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, _ models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	var meshAdapters []*models.Adapter
	if prefObj != nil && prefObj.MeshAdapters != nil {
		meshAdapters = prefObj.MeshAdapters
	} else {
		meshAdapters = []*models.Adapter{}
	}

	// adapterLoc := req.PostFormValue("adapter")
	adapterLoc := req.URL.Query().Get("adapter")
	if h.log != nil {
		h.log.Debug("Adapter url to ping: ", adapterLoc)
	}
	//logrus.Debug("Adapter url to ping: ", adapterLoc)

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		if h.log != nil {
			h.log.Error(ErrValidAdapter)
		}
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), meshAdapters[aID].Location)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrMeshClient)
		}
		writeMeshkitError(w, ErrMeshClient, http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = mClient.Close()
	}()
	_, err = mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrMeshClient)
		}
		writeMeshkitError(w, ErrMeshClient, http.StatusInternalServerError)
		return
	}

	writeJSONEmptyObject(w, http.StatusOK)
}

// MeshAdapterConfigHandler is used to persist adapter config
func (h *Handler) MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
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
	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}
	var err error

	switch req.Method {
	case http.MethodPost:
		meshLocationURL := req.FormValue("meshLocationURL")

		if h.log != nil {
			h.log.Debug("meshLocationURL: ", meshLocationURL)
		}
		if strings.TrimSpace(meshLocationURL) == "" {
			if h.log != nil {
				h.log.Error(ErrAddAdapter)
			}
			writeMeshkitError(w, ErrAddAdapter, http.StatusBadRequest)
			return
		}
		meshAdapters, err = h.addAdapter(req.Context(), meshAdapters, prefObj, meshLocationURL, provider)
		if err != nil {
			if h.log != nil {
				h.log.Error(ErrRetrieveData(err))
			}
			writeMeshkitError(w, ErrRetrieveData(err), http.StatusInternalServerError)
			return // error is handled appropriately in the relevant method
		}
	case http.MethodDelete:
		meshAdapters, err = h.deleteAdapter(meshAdapters, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}
	default:
		// Correct status is 405 (Method Not Allowed), not 404 — the
		// endpoint exists but the verb is wrong. Fixing alongside JSON.
		writeMeshkitError(w, ErrMethodNotAllowed(req.Method), http.StatusMethodNotAllowed)
		return
	}

	prefObj.MeshAdapters = meshAdapters
	err = provider.RecordPreferences(req, user.ID.String(), prefObj)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRecordPreferences(err))
		}
		writeMeshkitError(w, ErrRecordPreferences(err), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(meshAdapters)
	if err != nil {
		obj := "data"
		// Streaming encode failure — log only.
		if h.log != nil {
			if isClientDisconnect(err) {
				h.log.Debug(models.ErrMarshal(err, obj))
			} else {
				h.log.Error(models.ErrMarshal(err, obj))
			}
		}
		return
	}
}

func (h *Handler) addAdapter(ctx context.Context, meshAdapters []*models.Adapter, _ *models.Preference, meshLocationURL string, _ models.Provider) ([]*models.Adapter, error) {
	alreadyConfigured := false
	for _, adapter := range meshAdapters {
		if adapter.Location == meshLocationURL {
			alreadyConfigured = true
			break
		}
	}

	if alreadyConfigured {
		if h.log != nil {
			h.log.Debug("Adapter already configured...")
		}
		return meshAdapters, nil
	}
	mClient, err := meshes.CreateClient(ctx, meshLocationURL)
	if err != nil {
		// Caller decides how loudly to surface this. SessionSyncHandler probes
		// tracked adapters on every page load and these probes routinely fail
		// for stale or never-deployed adapters (e.g., the Meshery Playground),
		// which otherwise floods server logs with ERROR entries.
		//
		// Log both the raw transport `err` and the MeshKit error code so
		// operators can distinguish DNS resolution from TLS handshake from
		// dial timeout when troubleshooting, while still being able to grep
		// the stable error code for cross-referencing with docs or support
		// tickets. The wrapped ErrMeshClient is still what we return — the
		// public error shape stays unchanged.
		if h.log != nil {
			h.log.Debugf("adapter unreachable at %s (%s): %v", meshLocationURL, ErrMeshClientCode, err)
		}
		return meshAdapters, ErrMeshClient
	}
	if h.log != nil {
		h.log.Debug("created client for adapter: ", meshLocationURL)
	}
	defer func() {
		_ = mClient.Close()
	}()
	respOps, err := mClient.MClient.SupportedOperations(ctx, &meshes.SupportedOperationsRequest{})
	if err != nil {
		if h.log != nil {
			h.log.Debugf("adapter %s did not return supported operations: %v", meshLocationURL, err)
		}
		return meshAdapters, err
	}
	if h.log != nil {
		h.log.Debug("retrieved supported ops for adapter: ", meshLocationURL)
	}
	meshInfo, err := mClient.MClient.ComponentInfo(ctx, &meshes.ComponentInfoRequest{})
	if err != nil {
		if h.log != nil {
			h.log.Debugf("adapter %s did not return component info: %v", meshLocationURL, err)
		}
		return meshAdapters, err
	}
	if h.log != nil {
		h.log.Debug("retrieved name for adapter: ", meshLocationURL)
	}
	result := &models.Adapter{
		Location:     meshLocationURL,
		Name:         meshInfo.Name,
		Version:      meshInfo.Version,
		GitCommitSHA: meshInfo.GitSha,
		Ops:          respOps.GetOps(),
	}

	if h.config != nil && h.config.AdapterTracker != nil {
		h.config.AdapterTracker.AddAdapter(ctx, *result)
	}
	meshAdapters = append(meshAdapters, result)
	return meshAdapters, nil
}

func (h *Handler) deleteAdapter(meshAdapters []*models.Adapter, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, error) {
	adapterLoc := req.URL.Query().Get("adapter")
	if h.log != nil {
		h.log.Debug("URL of adapter to be removed: ", adapterLoc)
	}

	adaptersLen := len(meshAdapters)

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		if h.log != nil {
			h.log.Error(ErrValidAdapter)
		}
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return meshAdapters, ErrValidAdapter
	}

	if h.config != nil && h.config.AdapterTracker != nil {
		h.config.AdapterTracker.RemoveAdapter(req.Context(), *meshAdapters[aID])
	}

	newMeshAdapters := []*models.Adapter{}
	switch aID {
	case 0:
		newMeshAdapters = append(newMeshAdapters, meshAdapters[1:]...)

	case adaptersLen - 1:
		newMeshAdapters = append(newMeshAdapters, meshAdapters[:adaptersLen-1]...)

	default:
		newMeshAdapters = append(newMeshAdapters, meshAdapters[0:aID]...)
		newMeshAdapters = append(newMeshAdapters, meshAdapters[aID+1:]...)
	}
	if h.log != nil {
		b, _ := json.Marshal(meshAdapters)
		h.log.Debug("Old adapters: ", b)
		b, _ = json.Marshal(newMeshAdapters)
		h.log.Debug("New adapters: ", b)
	}
	return newMeshAdapters, nil
}

// MeshOpsHandler is used to send operations to the adapters
func (h *Handler) MeshOpsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
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

	if provider.GetProviderType() == models.RemoteProviderType {
		token, err := provider.GetProviderToken(req)

		if err == nil {
			viper.SetDefault("opt-token", token)
		}
	}

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
	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}

	adapterLoc := req.PostFormValue("adapter")
	if h.log != nil {
		h.log.Debug("Adapter URL to execute operations on: ", adapterLoc)
	}

	aID := -1
	for i, ad := range meshAdapters {
		if adapterLoc == ad.Location {
			aID = i
		}
	}
	if aID < 0 {
		if h.log != nil {
			h.log.Error(ErrValidAdapter)
		}
		writeMeshkitError(w, ErrValidAdapter, http.StatusBadRequest)
		return
	}

	opName := req.PostFormValue("query")
	customBody := req.PostFormValue("customBody")
	namespace := req.PostFormValue("namespace")
	deleteOp := req.PostFormValue("deleteOp")
	version := req.PostFormValue("version")
	if namespace == "" {
		namespace = "default"
	}
	mk8sContexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok || len(mk8sContexts) == 0 {
		if h.log != nil {
			h.log.Error(ErrInvalidK8SConfigNil)
		}
		writeMeshkitError(w, ErrInvalidK8SConfigNil, http.StatusBadRequest)
		return
	}
	var configs []string
	for _, c := range mk8sContexts {
		// Generate Kube Handler
		kc, err := c.GenerateKubeConfig()
		if err != nil {
			if h.log != nil {
				h.log.Error(ErrInvalidKubeConfig(err, c.Name))
			}
			writeMeshkitError(w, ErrInvalidKubeConfig(err, c.Name), http.StatusBadRequest)
			return
		}
		configs = append(configs, string(kc))
	}
	mClient, err := meshes.CreateClient(req.Context(), meshAdapters[aID].Location)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrMeshClient)
		}
		writeMeshkitError(w, ErrMeshClient, http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = mClient.Close()
	}()
	operationID, err := uuid.NewV4()

	if err != nil {
		if h.log != nil {
			h.log.Error(ErrGenerateUUID(err))
		}
		writeMeshkitError(w, ErrGenerateUUID(err), http.StatusInternalServerError)
		return
	}
	// Record the owner first: the adapter may send the operation's events
	// before ApplyOperation returns.
	if h.adapterOperations != nil {
		h.adapterOperations.record(operationID.String(), user.ID.String())
	}
	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OperationId: operationID.String(),
		OpName:      opName,
		Username:    user.ID.String(),
		Namespace:   namespace,
		CustomBody:  customBody,
		DeleteOp:    (deleteOp != ""),
		KubeConfigs: configs,
		Version:     version,
	})
	if err != nil {
		if h.adapterOperations != nil {
			h.adapterOperations.forget(operationID.String())
		}
		if h.log != nil {
			h.log.Error(ErrApplyChange(err))
		}
		writeMeshkitError(w, ErrApplyChange(err), http.StatusInternalServerError)
		return
	}

	// Return the operation ID, so a client can pick its own operation's events
	// out of the event stream.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(adapterOperationResponse{OperationID: operationID.String()}); err != nil {
		if h.log != nil {
			h.log.Error(ErrEncodeResponse(err))
		}
	}
}

// adapterOperationResponse is the response to POST /api/system/adapter/operation.
type adapterOperationResponse struct {
	OperationID string `json:"operationId"`
}
