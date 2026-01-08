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
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
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
		h.log.Debug("adapter pinging")
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
	// logrus.Debug("Adapter url to ping: ", adapterLoc)

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

// For Configuration category, it further sub-categorizes by prefix (Add-ons, Policies, etc.)
func categorizeAdapterOps(ops []*meshes.SupportedOperation) map[string]interface{} {
	categoryNames := map[meshes.OpCategory]string{
		meshes.OpCategory_INSTALL:            "Service Mesh",
		meshes.OpCategory_SAMPLE_APPLICATION: "Sample Applications",
		meshes.OpCategory_CONFIGURE:          "Configuration",
		meshes.OpCategory_CUSTOM:             "Custom Operations",
	}

	result := make(map[string]interface{})
	configSubcategories := make(map[string][]string)

	for _, op := range ops {
		if op.Category == meshes.OpCategory_VALIDATE {
			continue
		}

		catName, ok := categoryNames[op.Category]
		if !ok {
			catName = "Custom Operations"
		}

		if op.Category == meshes.OpCategory_CONFIGURE {
			value := op.Value
			switch {
			case strings.HasPrefix(value, "Add-on:"):
				name := strings.TrimSpace(strings.TrimPrefix(value, "Add-on:"))
				configSubcategories["Add-ons"] = append(configSubcategories["Add-ons"], name)
			case strings.HasPrefix(value, "Policy:"):
				name := strings.TrimSpace(strings.TrimPrefix(value, "Policy:"))
				configSubcategories["Policies"] = append(configSubcategories["Policies"], name)
			default:
				configSubcategories["Other"] = append(configSubcategories["Other"], value)
			}
		} else {
			// For non-Configuration categories, store as simple string array
			if existing, ok := result[catName].([]string); ok {
				result[catName] = append(existing, op.Value)
			} else {
				result[catName] = []string{op.Value}
			}
		}
	}

	// Add Configuration sub-categories if any exist
	if len(configSubcategories) > 0 {
		result["Configuration"] = configSubcategories
	}

	return result
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
	var targetAdapter *models.Adapter
	var event *events.Event

	userID := uuid.FromStringOrNil(user.ID)

	// Build event object
	eventBuilder := events.NewEvent().
		FromUser(userID).
		FromSystem(*h.SystemID).
		WithCategory("adapter")

	switch req.Method {
	case http.MethodPost:
		meshLocationURL := req.FormValue("meshLocationURL")

		h.log.Debug("meshLocationURL: ", meshLocationURL)
		if strings.TrimSpace(meshLocationURL) == "" {
			h.log.Error(ErrAddAdapter)
			http.Error(w, ErrAddAdapter.Error(), http.StatusBadRequest)
			return
		}

		meshAdapters, targetAdapter, err = h.addAdapter(req.Context(), meshAdapters, prefObj, meshLocationURL, provider)
		if err != nil {
			err = ErrRetrieveData(err)

			h.log.Error(err)

			errDescription := errors.GetSDescription(err)
			errRemedy := errors.GetRemedy(err)
			errProbableCause := errors.GetCause(err)

			metadata := map[string]interface{}{
				"Details":               errDescription,
				"Suggested_Remediation": errRemedy,
				"Probable Cause":        errProbableCause,
				"location":              meshLocationURL,
			}

			event = eventBuilder.
				WithSeverity(events.Error).
				WithDescription("Error adding Adapter").
				WithAction(models.Register).
				WithMetadata(metadata).
				Build()

			_err := provider.PersistEvent(*event, nil)
			if _err != nil {
				h.log.Debug(fmt.Sprintf("Failed to persist event: %v", _err))
			}
			go h.config.EventBroadcaster.Publish(userID, event)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return // error is handled appropriately in the relevant method
		}

		if targetAdapter == nil {
			description := fmt.Sprintf("An adapter is already configured at %s", meshLocationURL)
			metadata := map[string]interface{}{
				"location": meshLocationURL,
			}
			event = eventBuilder.
				WithSeverity(events.Informational).
				WithDescription(description).
				WithAction(models.Register).
				WithMetadata(metadata).
				Build()

		} else {

			description := fmt.Sprintf("%s Adapter (%s) configured", targetAdapter.Name, targetAdapter.Version)
			metadata := map[string]interface{}{
				"adapter_name": targetAdapter.Name,
				"version":      targetAdapter.Version,
				"location":     targetAdapter.Location,
				"manifest":     categorizeAdapterOps(targetAdapter.Ops),
			}

			event = eventBuilder.
				WithSeverity(events.Success).
				WithDescription(description).
				WithAction(models.Register).
				WithMetadata(metadata).
				Build()
		}

	case http.MethodDelete:
		meshAdapters, targetAdapter, err = h.deleteAdapter(meshAdapters, w, req)
		if err != nil {
			return // error is handled appropriately in the relevant method
		}

		description := fmt.Sprintf("Removed %s Adapter", targetAdapter.Name)
		metadata := map[string]interface{}{
			"location": targetAdapter.Location,
		}

		event = eventBuilder.
			WithSeverity(events.Success).
			WithDescription(description).
			WithAction(models.Unregister).
			WithMetadata(metadata).
			Build()

	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	err = provider.PersistEvent(*event, nil)
	if err != nil {
		h.log.Debug(fmt.Sprintf("Failed to persist event: %v", err))
	}
	go h.config.EventBroadcaster.Publish(userID, event)

	prefObj.MeshAdapters = meshAdapters
	err = provider.RecordPreferences(req, user.UserId, prefObj)
	if err != nil {
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

func (h *Handler) addAdapter(ctx context.Context, meshAdapters []*models.Adapter, _ *models.Preference, meshLocationURL string, _ models.Provider) ([]*models.Adapter, *models.Adapter, error) {
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
		return meshAdapters, nil, nil
	}
	mClient, err := meshes.CreateClient(ctx, meshLocationURL)
	if err != nil {
		h.log.Error(ErrMeshClient)
		// http.Error(w, ErrMeshClient.Error(), http.StatusInternalServerError)
		return meshAdapters, nil, ErrMeshClient
	}
	h.log.Debug("created client for adapter: ", meshLocationURL)
	defer func() {
		_ = mClient.Close()
	}()
	respOps, err := mClient.MClient.SupportedOperations(ctx, &meshes.SupportedOperationsRequest{})
	if err != nil {
		h.log.Error(ErrRetrieveMeshData(err))
		// http.Error(w, ErrRetrieveMeshData(err).Error(), http.StatusInternalServerError)
		return meshAdapters, nil, err
	}
	h.log.Debug("retrieved supported ops for adapter: ", meshLocationURL)
	meshInfo, err := mClient.MClient.ComponentInfo(ctx, &meshes.ComponentInfoRequest{})
	if err != nil {
		h.log.Error(ErrRetrieveMeshData(err))
		// http.Error(w, ErrRetrieveMeshData(err).Error(), http.StatusInternalServerError)
		return meshAdapters, nil, err
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
	return meshAdapters, result, nil
}

func (h *Handler) deleteAdapter(meshAdapters []*models.Adapter, w http.ResponseWriter, req *http.Request) ([]*models.Adapter, *models.Adapter, error) {
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
		return meshAdapters, nil, ErrValidAdapter
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
	return newMeshAdapters, meshAdapters[aID], nil
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
		Username:    user.UserId,
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
