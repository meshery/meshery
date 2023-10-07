// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	gob.Register([]*models.Adapter{})
}

// swagger:route GET /api/system/adapters/available SystemAPI idGetAvailableAdapters
// Handle GET request for available adapters
//
// Fetches and returns all the adapters available for deployment
// Responses:
//  200: systemAdaptersRespWrapper

// AdaptersHandler is used to fetch all the adapters
func (h *Handler) AvailableAdaptersHandler(w http.ResponseWriter, _ *http.Request) {
	err := json.NewEncoder(w).Encode(models.ListAvailableAdapters)
	if err != nil {
		obj := "data"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/adapters SystemAPI idGetSystemAdapters
// Handle GET request for adapters
//
// Fetches and returns all the adapters and ping adapters
// Responses:
//  200: systemAdaptersRespWrapper

// AdaptersHandler is used to fetch all the adapters
func (h *Handler) AdaptersHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// if adapter found in query user is trying to ping an adapter
	adapterLoc := req.URL.Query().Get("adapter")
	if adapterLoc != "" {
		logrus.Debug("adapter pinging")
		h.AdapterPingHandler(w, req, prefObj, user, provider)
		return
	}

	err := json.NewEncoder(w).Encode(h.config.AdapterTracker.GetAdapters(req.Context()))
	if err != nil {
		obj := "data"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// AdapterPingHandler is used to ping a given adapter
func (h *Handler) AdapterPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, _ models.Provider) {
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
	logrus.Debug("Adapter url to ping: ", adapterLoc)

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

	mClient, err := meshes.CreateClient(req.Context(), meshAdapters[aID].Location)
	if err != nil {
		http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
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

// swagger:route POST /api/system/adapter/manage SystemAPI idPostAdapterConfig
// Handle POST requests to persist adapter config
//
// Used to persist adapter config
// responses:
// 	200: mesheryAdaptersRespWrapper

// swagger:route DELETE /api/system/adapter/manage SystemAPI idDeleteAdapterConfig
// Handle DELETE requests to delete adapter config
//
// Used to delete adapter configuration
// responses:
// 	200:

// MeshAdapterConfigHandler is used to persist adapter config
func (h *Handler) MeshAdapterConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	meshAdapters := prefObj.MeshAdapters
	if meshAdapters == nil {
		meshAdapters = []*models.Adapter{}
	}
	var err error
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("mesh_adapter")
	switch req.Method {
	case http.MethodPost:
		eventBuilder.WithAction("add")
		meshLocationURL := req.FormValue("meshLocationURL")

		h.log.Debug("meshLocationURL: ", meshLocationURL)
		if strings.TrimSpace(meshLocationURL) == "" {
			eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while parsing meshLocationURL").
				WithMetadata(map[string]interface{}{
					"error":                 "meshLocationURL is empty",
					"probable_cause":        "The request body is missing value for `meshLocationURL` field.",
					"suggested_remediation": "Please ensure that you're providing a valid adapter location URL in the request body",
				})
			h.broadcastEvent(eventBuilder, provider, userID)
			h.log.Error(ErrAddAdapter)
			http.Error(w, ErrAddAdapter.Error(), http.StatusBadRequest)
			return
		}
		meshAdapters, err = h.addAdapter(req.Context(), meshAdapters, prefObj, meshLocationURL, provider)
		if err != nil {
			eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while retrieving adapter information.").
				WithMetadata(map[string]interface{}{
					"error":                 "Meshery encounters issues when attempting to add and connect with adapter.",
					"probable_cause":        "The requested adapter service might not be up and running.",
					"suggested_remediation": "Please ensure that adapter service is up and running and adapter has been deployed successfully.",
				})
			h.broadcastEvent(eventBuilder, provider, userID)
			// h.log.Error(ErrRetrieveData(err))
			http.Error(w, ErrRetrieveData(err).Error(), http.StatusInternalServerError)
			return // error is handled appropriately in the relevant method
		}
		// if h.addAdapter succeeds the last adapter in meshAdapters slice is the recently added adapter.
		adapter := meshAdapters[len(meshAdapters)-1]
		eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Successfully connected to `%s` adapter", adapter.Name)).
			WithMetadata(map[string]interface{}{
				"name":     adapter.Name,
				"location": adapter.Location,
				"version":  adapter.Version,
			})
		h.broadcastEvent(eventBuilder, provider, userID)
	case http.MethodDelete:
		eventBuilder.WithAction("delete")
		meshAdapters, err = h.deleteAdapter(meshAdapters, w, req)
		if err != nil {
			eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while deleting adapter.").
				WithMetadata(map[string]interface{}{
					"error":                 "Meshery encounters issues when attempting to delete adapter.",
					"probable_cause":        "The requested adapter might not have been deployed and added in preferences.",
					"suggested_remediation": "Please ensure that adapter had been deployed successfully.",
				})
			h.broadcastEvent(eventBuilder, provider, userID)
			return // error is handled appropriately in the relevant method
		}
		adapterLocation := req.URL.Query().Get("adapter")
		var adapter *models.Adapter
		for _, ad := range meshAdapters {
			if ad.Location == adapterLocation {
				adapter = ad
			}
		}
		eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Successfully deleted adapter `%s`", adapter.Name)).
			WithMetadata(map[string]interface{}{
				"name":     adapter.Name,
				"location": adapter.Location,
				"version":  adapter.Version,
			})
		h.broadcastEvent(eventBuilder, provider, userID)
	default:
		eventBuilder.WithAction("manage").WithSeverity(events.Error).WithDescription(fmt.Sprintf("The method `%s` is not allowed while managing adapter", req.Method)).
			WithMetadata(map[string]interface{}{
				"error":                 "Unsupported http verb",
				"probable_cause":        "You've called api endpoint with wronge method.",
				"suggested_remediation": "Please ensure that you're using the allowed HTTP methods (POST, DELETE) for this operation.",
			})
		h.broadcastEvent(eventBuilder, provider, userID)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	prefObj.MeshAdapters = meshAdapters
	err = provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while persisting user config data.").
			WithMetadata(map[string]interface{}{
				"error":                 "Unable to persist user config data",
				"probable_cause":        "User token might be invalid or db migh have been corrupted.",
				"suggested_remediation": "Logging out and then logging in again may resolve the issue.",
			})
		h.broadcastEvent(eventBuilder, provider, userID)
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(meshAdapters)
	if err != nil {
		obj := "data"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) addAdapter(ctx context.Context, meshAdapters []*models.Adapter, _ *models.Preference, meshLocationURL string, _ models.Provider) ([]*models.Adapter, error) {
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
	mClient, err := meshes.CreateClient(ctx, meshLocationURL)
	if err != nil {
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

// swagger:route POST /api/system/adapter/operation SystemAPI idPostAdapterOperation
// Handle POST requests for Adapter Operations
//
// Used to send operations to the adapters
// responses:
// 	200:

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
	deleteOp := req.PostFormValue("deleteOp")
	version := req.PostFormValue("version")
	if namespace == "" {
		namespace = "default"
	}
	mk8sContexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
	if !ok || len(mk8sContexts) == 0 {
		h.log.Error(ErrInvalidK8SConfigNil)
		http.Error(w, ErrInvalidK8SConfigNil.Error(), http.StatusBadRequest)
		return
	}
	var configs []string
	for _, c := range mk8sContexts {
		// Generate Kube Handler
		kc, err := c.GenerateKubeConfig()
		if err != nil {
			return
		}
		configs = append(configs, string(kc))
	}
	mClient, err := meshes.CreateClient(req.Context(), meshAdapters[aID].Location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = mClient.Close()
	}()
	operationID, err := uuid.NewV4()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OperationId: operationID.String(),
		OpName:      opName,
		Username:    user.UserID,
		Namespace:   namespace,
		CustomBody:  customBody,
		DeleteOp:    (deleteOp != ""),
		KubeConfigs: configs,
		Version:     version,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
