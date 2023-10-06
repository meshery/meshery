package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
)

// swagger:route GET /api/system/kubernetes/contexts GetAllContexts idGetAllContexts
// Handle GET request for all kubernetes contexts.
//
// # Contexts can be further filtered through query parameter
//
// ```?order={field}``` orders on the passed field
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// ```?search={contextname}``` If search is non empty then a greedy search is performed
// responses:
//
//	200: systemK8sContextsResponseWrapper
func (h *Handler) GetAllContexts(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("retrieve")

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()
	// Don't fetch credentials as UI has no use case.
	vals, err := provider.GetK8sContexts(token, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), false)
	if err != nil {
		eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while retrieving Kubernetes contexts.").
			WithMetadata(map[string]interface{}{
				"error":                 err,
				"probable_cause":        "You might have provided incorrect query params or there may be an issue with connectivity between Meshery and kubernetes.",
				"suggested_remediation": "Please ensure that you've provided valid query params (page, pagesize, search, order) and Meshery can establish connectivity with kubernetes.",
			})
		h.broadcastEvent(eventBuilder, provider, userID)
		http.Error(w, "failed to get contexts", http.StatusInternalServerError)
		return
	}
	var mesheryK8sContextPage models.MesheryK8sContextPage
	err = json.Unmarshal(vals, &mesheryK8sContextPage)
	if err != nil {
		obj := "k8s context"
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	contextNames := []string{}
	for _, context := range mesheryK8sContextPage.Contexts {
		contextNames = append(contextNames, context.Name)
	}
	if len(contextNames) > 0 {
		eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Successfully retrieved %d kubernetes contexts.", len(contextNames))).
			WithMetadata(map[string]interface{}{
				"context_name": strings.Join(contextNames, ","),
			})
		h.broadcastEvent(eventBuilder, provider, userID)
	}

	if err := json.NewEncoder(w).Encode(mesheryK8sContextPage); err != nil {
		http.Error(w, "failed to encode contexts", http.StatusInternalServerError)
		return
	}
}

// not being used....
func (h *Handler) GetContext(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	// if req.URL.Query().Get("current") != "" {
	// 	context, ok := req.Context().Value(models.KubeContextKey).(*models.K8sContext)
	// 	if !ok || context == nil {
	// 		http.Error(w, "failed to get context", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if err := json.NewEncoder(w).Encode(context); err != nil {
	// 		http.Error(w, "failed to encode context", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	return
	// }

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("retrieve")

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}
	id := mux.Vars(req)["id"]
	connectionID := uuid.FromStringOrNil(id)
	eventBuilder.ActedUpon(connectionID)
	h.log.Info("this is being used\n\n\n")
	val, err := provider.GetK8sContext(token, id)
	if err != nil {
		eventBuilder.WithSeverity(events.Error).WithDescription("Encountered an issue while retrieving Kubernetes context.").
			WithMetadata(map[string]interface{}{
				"error":                 err,
				"probable_cause":        "You might have passed incorrect query params or there may be an issue with the connectivity between meshery and kubernetes.",
				"suggested_remediation": "Please ensure that you've provided valid query params (page, pagesize, search, order), and Meshery can establish connectivity with kubernetes.",
			})
		h.broadcastEvent(eventBuilder, provider, userID)
		http.Error(w, "failed to get context", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(val); err != nil {
		http.Error(w, "failed to encode context", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteContext(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	userID := uuid.FromStringOrNil(user.ID)
	contextID := mux.Vars(req)["id"]

	eventBuilder := events.NewEvent().ActedUpon(uuid.FromStringOrNil(contextID)).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("delete")

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	obj := "connection"

	// id is the connection_id of the specific cluster in connections table
	deletedContext, err := provider.DeleteK8sContext(token, contextID)
	if err != nil {
		_err := ErrFailToDelete(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error deleting Kubernetes context").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Kubernetes context %s deleted.", deletedContext.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.config.K8scontextChannel.PublishContext()
	go models.FlushMeshSyncData(req.Context(), deletedContext, provider, h.config.EventBroadcaster, user.ID, h.SystemID)
}

// func (h *Handler) GetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	token, ok := req.Context().Value(models.TokenCtxKey).(string)
// 	if !ok {
// 		http.Error(w, "failed to get token", http.StatusInternalServerError)
// 		return
// 	}

// 	val, err := h.GetCurrentContext(token, provider)
// 	if err != nil {
// 		http.Error(w, "failed to get current context", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := json.NewEncoder(w).Encode(val); err != nil {
// 		http.Error(w, "failed to encode context", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) SetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	token, ok := req.Context().Value(models.TokenCtxKey).(string)
// 	if !ok {
// 		http.Error(w, "failed to get token", http.StatusInternalServerError)
// 		return
// 	}

// 	_, err := provider.SetCurrentContext(token, mux.Vars(req)["id"])
// 	if err != nil {
// 		http.Error(w, "failed to set current context", http.StatusInternalServerError)
// 		return
// 	}
// }
