package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/models/events"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	schemasConnection "github.com/meshery/schemas/models/v1beta1/connection"
)

func (h *Handler) ProcessConnectionRegistration(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodDelete {
		h.handleProcessTermination(w, req)
		return
	}

	connectionRegisterPayload := connections.ConnectionPayload{}
	userUUID := user.ID
	err := json.NewDecoder(req.Body).Decode(&connectionRegisterPayload)
	if err != nil {
		http.Error(w, models.ErrUnmarshal(err, "connection registration payload").Error(), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*h.SystemID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	if string(connectionRegisterPayload.Status) == string(machines.Init) {
		h.handleRegistrationInitEvent(w, req, &connectionRegisterPayload)
	} else {
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker

		machineCtx := make(map[string]string, 0)
		inst, err := helpers.InitializeMachineWithContext(
			machineCtx,
			req.Context(),
			connectionRegisterPayload.ID,
			userUUID,
			smInstanceTracker,
			h.log,
			provider,
			machines.DISCOVERED,
			strings.ToLower(connectionRegisterPayload.Kind),
			nil,
		)
		if err != nil {
			event := eventBuilder.WithSeverity(events.Error).WithDescription("Unable to perisit the \"%s\" connection details").WithMetadata(map[string]interface{}{
				"error": err,
			}).Build()
			_ = provider.PersistEvent(*event, nil)
			go h.config.EventBroadcaster.Publish(userUUID, event)
		}

		event, err := inst.SendEvent(req.Context(), machines.EventType(connectionRegisterPayload.Status), connectionRegisterPayload)
		if err != nil {
			h.log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			_ = provider.PersistEvent(*event, nil)
			go h.config.EventBroadcaster.Publish(userUUID, event)
		}
	}
}

func (h *Handler) handleProcessTermination(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]string, 0)
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		_err := models.ErrUnmarshal(err, "request body")
		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}
	smInstancetracker := h.ConnectionToStateMachineInstanceTracker

	id, ok := body["id"]
	if ok {
		smInstancetracker.Remove(uuid.FromStringOrNil(id))
	}
}

func (h *Handler) handleRegistrationInitEvent(w http.ResponseWriter, req *http.Request, payload *connections.ConnectionPayload) {
	compFilter := &regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sConnection", payload.Kind),
		Limit: 1,
	}
	schema := make(map[string]interface{}, 1)
	connectionComponent, _, _, _ := h.registryManager.GetEntities(compFilter)
	if len(connectionComponent) == 0 {
		http.Error(w, "Unable to register resource as connection. No matching connection definition found in the registry", http.StatusInternalServerError)
		return
	}

	schema["connection"] = connectionComponent[0]
	credential, _, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sCredential", payload.Kind),
		Limit: 1,
	})

	if len(credential) > 0 {
		schema["credential"] = credential[0]
	}
	// id act as a connection registration process tracker.
	// The clients should always include this "id" in the subsequent API calls until the process is completed or terminated.
	id, _ := uuid.NewV4()
	schema["id"] = id

	err := json.NewEncoder(w).Encode(&schema)
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

// swagger:route POST /api/integrations/connections PostConnection idPostConnection
// Handle POST request for creating a new connection
//
// Creates a new connection
// responses:
// 201: noContentWrapper
func (h *Handler) SaveConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	userID := user.ID
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	connection := connections.ConnectionPayload{}
	err = json.Unmarshal(bd, &connection)
	obj := "connection"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("create")

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	if token == "" {
		if ck, err := req.Cookie(models.TokenCookieName); err == nil {
			token = ck.Value
		}
	}

	_, err = provider.SaveConnection(&connection, token, false)
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Error creating connection %s", connection.Name)).WithMetadata(metadata).Build()
		_ = provider.PersistEvent(*event, nil)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s created.", connection.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(*event, nil)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /api/integrations/connections GetConnections idGetConnections
// Handle GET request for getting all connections
//
// ```?order={field}``` orders on the passed field
//
// ```?search={}``` If search is non empty then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10
//
// ```?filter={filter}``` Filter connections with type or sub_type, eg /api/integrations/connections?filter=type%20platform or /api/integrations/connections?filter=sub_type%20management
//
// ```?status={status}``` Status takes array as param to filter connections based on status, eg /api/integrations/connections?status=["connected", "deleted"]
//
// ```?kind={kind}``` Kind takes array as param to filter connections based on kind, eg /api/integrations/connections?kind=["meshery", "kubernetes"]
//
// ```?type={type}``` Type takes array as param to filter connections based on type, eg /api/integrations/connections?type=["platform", "observability"]
//
// ```?name={name}``` Name filters connections by name (partial match), eg /api/integrations/connections?name=my-cluster
// responses:
// 200: mesheryConnectionsResponseWrapper
func (h *Handler) GetConnections(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSizeStr := q.Get("pagesize")
	filter := q.Get("filter")
	name := q.Get("name")

	var pageSize int
	if pageSizeStr == "all" {
		pageSize = 100
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if pageSize > 100 {
		pageSize = 100
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "updated_at desc"
	}

	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrGetConnections(err))
		http.Error(w, ErrGetConnections(err).Error(), http.StatusInternalServerError)
		return
	}

	queryParam := struct {
		Status []string `json:"status"`
		Kind   []string `json:"kind"`
		Type   []string `json:"type"`
	}{}

	status := q.Get("status")
	kind := q.Get("kind")
	connType := q.Get("type")
	if status != "" {
		err := json.Unmarshal([]byte(status), &queryParam.Status)
		if err != nil {
			h.log.Error(ErrGetConnections(err))
			http.Error(w, ErrGetConnections(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	if kind != "" {
		err := json.Unmarshal([]byte(kind), &queryParam.Kind)
		if err != nil {
			h.log.Error(ErrGetConnections(err))
			http.Error(w, ErrGetConnections(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	if connType != "" {
		err := json.Unmarshal([]byte(connType), &queryParam.Type)
		if err != nil {
			h.log.Error(ErrGetConnections(err))
			http.Error(w, ErrGetConnections(err).Error(), http.StatusInternalServerError)
			return
		}
	}

	connectionsPage, err := provider.GetConnections(req, user.ID.String(), page, pageSize, search, order, filter, queryParam.Status, queryParam.Kind, queryParam.Type, name)
	obj := "connections"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		http.Error(w, models.ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// GetConnectionsByKind is an internal handler that fetches connections filtered by kind.
// Note: This handler is used internally by other handlers (e.g., GrafanaConfigHandler, PrometheusConfigHandler)
// and is not exposed as an HTTP route since the route was deprecated in schemas v0.8.115.
func (h *Handler) GetConnectionsByKind(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	connectionKind := mux.Vars(req)["connectionKind"]
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(q.Get("pagesize"))

	if pageSize > 25 {
		pageSize = 25
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "updated_at desc"
	}

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s, kind: %s", page+1, pageSize, search, order, connectionKind))

	// Use GetConnections with kind filter
	connectionsPage, err := provider.GetConnections(req, user.ID.String(), page, pageSize, search, order, "", nil, []string{connectionKind}, nil, "")
	obj := "connections"

	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		http.Error(w, models.ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/integrations/connections/{connectionId} GetConnectionById idGetConnectionById
// Handle GET request for getting a single connection by its ID
//
// Fetches a single connection by its ID
// responses:
// 200: mesheryConnectionResponseWrapper
func (h *Handler) GetConnectionByID(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	if connectionID == uuid.Nil {
		h.log.Error(ErrGetConnections(fmt.Errorf("invalid connection ID")))
		http.Error(w, "Invalid connection ID", http.StatusBadRequest)
		return
	}

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	obj := "connection"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(connection); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		http.Error(w, models.ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route PUT /api/integrations/connections/{connectionId} PutConnectionById idPutConnectionById
// Handle PUT request for updating an existing connection by connection ID
//
// Updates existing connection using ID
// responses:
// 200: mesheryConnectionResponseWrapper
func (h *Handler) UpdateConnectionById(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := user.ID

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	connection := &connections.ConnectionPayload{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	// Check if meshsync deployment mode is specified in the payload metadata.
	// Only handle mode changes if a mode is explicitly provided (not undefined).
	// In fact this method is used (for now) only for perform meshsync deployment mode change.
	// If mode change fails return error.
	// TODO: also check that kind = "kubernetes" (when client starts to send full connection object)
	if schemasConnection.MeshsyncDeploymentModeFromMetadata(connection.MetaData) != schemasConnection.MeshsyncDeploymentModeUndefined {
		// Handle meshsync deployment mode changes before connection update
		token, _ := req.Context().Value(models.TokenCtxKey).(string)
		oldMode, newMode, modeChanged, err := h.handleMeshSyncDeploymentModeChange(
			req.Context(),
			connectionID,
			connection,
			token,
			userID,
			provider,
		)
		if err != nil {
			meshSyncErr := fmt.Errorf("error handling meshsync deployment mode change: %w", err)
			metadata := map[string]any{
				"error":         meshSyncErr,
				"connection_id": connectionID.String(),
			}
			event := eventBuilder.WithSeverity(events.Error).WithDescription("Failed to handle meshsync deployment mode change").WithMetadata(metadata).Build()
			_ = provider.PersistEvent(*event, nil)
			go h.config.EventBroadcaster.Publish(userID, event)

			h.log.Error(meshSyncErr)
			http.Error(w, meshSyncErr.Error(), http.StatusInternalServerError)
			return
		}

		// Log and emit event if mode actually changed
		if modeChanged {
			description := fmt.Sprintf("MeshSync deployment mode changed from '%s' to '%s' for connection %s", oldMode, newMode, connectionID)
			metadata := map[string]any{
				"meshsync_deployment_mode_old": oldMode,
				"meshsync_deployment_mode_new": newMode,
				"connection_id":                connectionID.String(),
			}
			event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).WithMetadata(metadata).Build()
			_ = provider.PersistEvent(*event, nil)
			go h.config.EventBroadcaster.Publish(userID, event)

			h.log.Info(description)
		}
	}

	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, ErrGetResult(err).Error(), http.StatusInternalServerError)
		return
	}
	updatedConnection, err := provider.UpdateConnectionById(token, connection, mux.Vars(req)["connectionId"])
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error updating connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(*event, nil)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO enhance event with information about meshsync deployment mode change
	description := fmt.Sprintf("Connection %s updated.", updatedConnection.Name)
	eventBuilder = eventBuilder.WithDescription(description)

	if connection.Status != "" {
		event, _ := h.NotifySmOfConnectionStatusChange(req.Context(), userID, provider, token, connection)
		_ = provider.PersistEvent(event, nil)
	}

	event := eventBuilder.WithSeverity(events.Informational).Build()
	_ = provider.PersistEvent(*event, nil)
	go h.config.EventBroadcaster.Publish(userID, event)
	h.log.Info(description)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NotifySmOfConnectionStatusChange(context context.Context, userID uuid.UUID, provider models.Provider, token string, connection *connections.ConnectionPayload) (events.Event, error) {
	connectionID := connection.ID

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	if connection.Status != "" {
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
		// token, _ := req.Context().Value(models.TokenCtxKey).(string)
		k8scontext, err := provider.GetK8sContext(token, connectionID.String())

		if err != nil {
			eventBuilder = eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", connectionID)).WithMetadata(map[string]interface{}{
				"error": err,
			})

			return *eventBuilder.Build(), err
		}

		eventBuilder = eventBuilder.WithSeverity(events.Informational).
			WithDescription(fmt.Sprintf("Processing status update to \"%s\" for connection %s", connection.Status, k8scontext.Name)).
			WithMetadata(map[string]interface{}{
				"connectionName": k8scontext.Name,
			})

		machineCtx := &kubernetes.MachineCtx{
			K8sContext:         k8scontext,
			MesheryCtrlsHelper: h.MesheryCtrlsHelper,
			K8sCompRegHelper:   h.K8sCompRegHelper,
			OperatorTracker:    h.config.OperatorTracker,
			K8scontextChannel:  h.config.K8scontextChannel,
			EventBroadcaster:   h.config.EventBroadcaster,
			RegistryManager:    h.registryManager,
		}

		inst, err := helpers.InitializeMachineWithContext(
			machineCtx,
			context,
			connectionID,
			userID,
			smInstanceTracker,
			h.log,
			provider,
			machines.InitialState,
			"kubernetes",
			kubernetes.AssignInitialCtx,
		)

		if err != nil {
			eventBuilder = eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", connectionID)).WithMetadata(map[string]interface{}{
				"error": err,
			})
			return *eventBuilder.Build(), err
		}

		go func(inst *machines.StateMachine, status connections.ConnectionStatus) {
			event, err := inst.SendEvent(context, machines.EventType(helpers.StatusToEvent(status)), nil)
			if err != nil {
				h.log.Error(err)
				_ = provider.PersistEvent(*event, nil)
				h.config.EventBroadcaster.Publish(userID, event)
				return
			}

			if status == connections.DELETED {
				smInstanceTracker.Remove(inst.ID)
			}

			_ = provider.PersistEvent(*event, nil)
			h.config.EventBroadcaster.Publish(userID, event)
		}(inst, connection.Status)
	}

	return *eventBuilder.Build(), nil
}

// swagger:route DELETE /api/integrations/connections/{connectionId} DeleteConnection idDeleteConnection
// Handle DELETE request for deleting an existing connection by connection ID
//
// Deletes existing connection
// responses:
// 200: noContentWrapper
func (h *Handler) DeleteConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := user.ID
	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("delete")

	deletedConnection, err := provider.DeleteConnection(req, connectionID)
	if err != nil {
		obj := "connection"
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error deleting connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(*event, nil)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s deleted.", deletedConnection.Name)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()

	_ = provider.PersistEvent(*event, nil)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info("connection deleted.")
	w.WriteHeader(http.StatusOK)
}

// handleMeshSyncDeploymentModeChange retrieves existing connection, compares meshsync deployment modes
// between existing and new connections, and performs necessary actions when they differ
// Returns: oldMode, newMode, changed, error
func (h *Handler) handleMeshSyncDeploymentModeChange(
	ctx context.Context,
	connectionID uuid.UUID,
	newConnection *connections.ConnectionPayload,
	token string,
	userID uuid.UUID,
	provider models.Provider,
) (schemasConnection.MeshsyncDeploymentMode, schemasConnection.MeshsyncDeploymentMode, bool, error) {
	if newConnection == nil {
		return schemasConnection.MeshsyncDeploymentModeUndefined, schemasConnection.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("new connection is nil, cannot compare meshsync deployment modes")
	}

	if h.SystemID == nil {
		return schemasConnection.MeshsyncDeploymentModeUndefined, schemasConnection.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("system ID is not configured in handler")
	}
	// TODO is h.SystemID a correct instance id here?
	mesheryInstanceID := *h.SystemID

	// Retrieve existing connection for mode comparison
	existingConnection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		return schemasConnection.MeshsyncDeploymentModeUndefined, schemasConnection.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("failed to retrieve existing connection (status %d): %w", statusCode, err)
	}

	if existingConnection == nil {
		return schemasConnection.MeshsyncDeploymentModeUndefined, schemasConnection.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("existing connection is nil, cannot compare meshsync deployment modes")
	}

	if existingConnection.Kind != "kubernetes" {
		return schemasConnection.MeshsyncDeploymentModeUndefined, schemasConnection.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("connection is not of kind kubernetes")
	}

	existingMeshSyncMode := schemasConnection.MeshsyncDeploymentModeFromMetadata(existingConnection.Metadata)
	newMeshSyncMode := schemasConnection.MeshsyncDeploymentModeFromMetadata(newConnection.MetaData)

	// draw back to default mode
	if newMeshSyncMode == schemasConnection.MeshsyncDeploymentModeUndefined {
		newMeshSyncMode = h.MeshsyncDefaultDeploymentMode
	}

	meshSyncModeChanged := existingMeshSyncMode != newMeshSyncMode
	if meshSyncModeChanged {
		instanceTracker := h.ConnectionToStateMachineInstanceTracker
		if instanceTracker == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("instance tracker is nil in handler instance")
		}

		machine, ok := instanceTracker.Get(connectionID)
		if !ok || machine == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("instance tracker does not contain machine for connection %s", connectionID)
		}

		machineCtx, err := kubernetes.GetMachineCtx(machine.Context, nil)
		if err != nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("failed to get machine context for connection %s: %w", connectionID, err)
		}

		if machineCtx == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("machine context is nil for connection %s", connectionID)
		}

		ctrlHelper := machineCtx.MesheryCtrlsHelper
		if ctrlHelper == nil {
			return existingMeshSyncMode, newMeshSyncMode, false, fmt.Errorf("machine context does not contain reference to MesheryCtrlsHelper for connection %s", connectionID)
		}

		// disconnect
		{
			contextID := machineCtx.K8sContext.ID
			ctrlHelper.
				UpdateOperatorsStatusMap(machineCtx.OperatorTracker).
				UndeployDeployedOperators(machineCtx.OperatorTracker).
				RemoveCtxControllerHandler(ctx, contextID)
			ctrlHelper.RemoveMeshSyncDataHandler(ctx, contextID)
		}
		// connect
		{
			ctrlHelper.
				AddCtxControllerHandlers(machineCtx.K8sContext).
				SetMeshsyncDeploymentMode(newMeshSyncMode).
				UpdateOperatorsStatusMap(machineCtx.OperatorTracker).
				DeployUndeployedOperators(machineCtx.OperatorTracker).
				AddMeshsynDataHandlers(ctx, machineCtx.K8sContext, userID, mesheryInstanceID, provider)
		}

	}

	return existingMeshSyncMode, newMeshSyncMode, meshSyncModeChanged, nil
}
