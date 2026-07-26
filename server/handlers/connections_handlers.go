package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/meshery/schemas/models/core"
	connectionv1beta3 "github.com/meshery/schemas/models/v1beta3/connection"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/helpers"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	regv1beta1 "github.com/meshery/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/meshkit/utils"
)

// ProcessConnectionRegistration drives the connection registration state
// machine (POST /api/integrations/connections/register). The wire contract is
// the schemas-defined ConnectionRegistrationEvent: an `initialize` event
// returns the registration bootstrap, every other event advances the tracked
// registration and responds with an empty body.
func (h *Handler) ProcessConnectionRegistration(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodDelete {
		// Deprecated: cancelling via DELETE on this route with the tracker id in
		// the body predates the schemas contract (v1.3.32). New clients cancel via
		// DELETE /api/integrations/connections/register/{registrationId}
		// (CancelConnectionRegister). Retire this branch once no supported client
		// sends the legacy shape.
		h.handleProcessTermination(w, req)
		return
	}

	registrationEvent := connectionv1beta3.ConnectionRegistrationEvent{}
	userUUID := user.ID
	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}
	err = json.NewDecoder(req.Body).Decode(&registrationEvent)
	if err != nil {
		writeMeshkitError(w, models.ErrUnmarshal(err, "connection registration payload"), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*h.SystemID).FromOwner(userUUID).WithDescription("Failed to interact with the connection.")

	if registrationEvent.Status == connectionv1beta3.ConnectionRegistrationEventStatusInitialize {
		h.handleRegistrationInitEvent(w, &registrationEvent)
	} else {
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
		connectionRegisterPayload := registrationEventToConnectionPayload(&registrationEvent)

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
			wrappedErr := ErrInitializeMachine(err)
			event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to persist the \"%s\" connection details", connectionRegisterPayload.Kind)).WithMetadata(map[string]interface{}{
				"error": wrappedErr,
			}).Build()
			if event != nil {
				_ = provider.PersistEvent(*event, token)
				go h.config.EventBroadcaster.Publish(userUUID, event)
			}
			h.log.Error(wrappedErr)
			writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
			return
		}

		event, err := inst.SendEvent(req.Context(), machines.EventType(registrationEvent.Status), connectionRegisterPayload)
		if err != nil {
			wrappedErr := ErrSendMachineEvent(err)
			h.log.Error(wrappedErr)
			writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
			if event != nil {
				_ = provider.PersistEvent(*event, token)
				go h.config.EventBroadcaster.Publish(userUUID, event)
			}
			return
		}
	}
}

// registrationEventToConnectionPayload converts the schemas wire event into the
// internal payload the connection state machines consume
// (machines/*/register.go casts the event data to connections.ConnectionPayload).
// The two shapes share the same wire fields; only the Go types differ.
func registrationEventToConnectionPayload(event *connectionv1beta3.ConnectionRegistrationEvent) connections.ConnectionPayload {
	payload := connections.ConnectionPayload{
		Kind:                       event.Kind,
		SubType:                    event.SubType,
		Type:                       event.Type,
		MetaData:                   event.Metadata,
		Status:                     connections.ConnectionStatus(event.Status),
		CredentialSecret:           event.CredentialSecret,
		Name:                       event.Name,
		CredentialID:               event.CredentialID,
		Model:                      event.Model,
		SkipCredentialVerification: event.SkipCredentialVerification,
	}
	if event.ID != nil {
		payload.ID = *event.ID
	}
	return payload
}

// handleProcessTermination is the deprecated body-based cancel path; see the
// deprecation note in ProcessConnectionRegistration.
func (h *Handler) handleProcessTermination(w http.ResponseWriter, req *http.Request) {
	body := make(map[string]string, 0)
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		_err := models.ErrUnmarshal(err, "request body")
		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusBadRequest)
		return
	}
	smInstancetracker := h.ConnectionToStateMachineInstanceTracker

	id, ok := body["id"]
	if ok {
		smInstancetracker.Remove(uuid.FromStringOrNil(id))
	}
}

// CancelConnectionRegister discards the in-progress registration state machine
// tracked by the {registrationId} path parameter
// (DELETE /api/integrations/connections/register/{registrationId}). Idempotent:
// unknown ids are ignored. Nothing is persisted for the abandoned process.
func (h *Handler) CancelConnectionRegister(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	registrationID, err := uuid.FromString(mux.Vars(req)["registrationId"])
	if err != nil {
		invalidIDErr := models.ErrInvalidUUID(err)
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}

	h.ConnectionToStateMachineInstanceTracker.Remove(registrationID)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) handleRegistrationInitEvent(w http.ResponseWriter, event *connectionv1beta3.ConnectionRegistrationEvent) {
	compFilter := &regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sConnection", event.Kind),
		Limit: 1,
	}
	connectionComponent, _, _, _ := h.registryManager.GetEntities(compFilter)
	if len(connectionComponent) == 0 {
		writeMeshkitError(w, ErrUnknownConnectionKind(event.Kind), http.StatusBadRequest)
		return
	}

	connectionDefinition, err := utils.MarshalAndUnmarshal[interface{}, core.Map](connectionComponent[0])
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
		writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
		return
	}

	// The bootstrap id acts as the connection registration process tracker.
	// Clients echo it on every subsequent event until the process completes or
	// is cancelled.
	bootstrap := connectionv1beta3.ConnectionRegistrationBootstrap{
		Connection: connectionDefinition,
		ID:         uuid.Must(uuid.NewV4()),
	}

	credential, _, _, _ := h.registryManager.GetEntities(&regv1beta1.ComponentFilter{
		Name:  fmt.Sprintf("%sCredential", event.Kind),
		Limit: 1,
	})
	if len(credential) > 0 {
		credentialDefinition, err := utils.MarshalAndUnmarshal[interface{}, core.Map](credential[0])
		if err != nil {
			h.log.Error(ErrWriteResponse(err))
			writeMeshkitError(w, ErrWriteResponse(err), http.StatusInternalServerError)
			return
		}
		bootstrap.Credential = credentialDefinition
	}

	err = json.NewEncoder(w).Encode(&bootstrap)
	if err != nil {
		h.log.Error(ErrWriteResponse(err))
	}
}

func (h *Handler) SaveConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	userID := user.ID
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	connection := connections.ConnectionPayload{}
	err = json.Unmarshal(bd, &connection)
	obj := "connection"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("create")

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
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s created.", connection.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info(description)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetConnections(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	// Canonical camelCase `pageSize` (matches the schemas wire contract); fall
	// back to legacy lowercase `pagesize` for older clients.
	pageSizeStr := q.Get("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = q.Get("pagesize")
	}
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
		writeMeshkitError(w, ErrGetConnections(err), http.StatusInternalServerError)
		return
	}

	// Filters are passed as repeated query params — the standard convention
	// across the API (e.g. ?kind=kubernetes&kind=meshery&status=connected). No
	// per-param JSON decoding.
	queryParam := struct {
		Status []string
		Kind   []string
		Type   []string
	}{
		Status: q["status"],
		Kind:   q["kind"],
		Type:   q["type"],
	}

	connectionsPage, err := provider.GetConnections(req, user.ID.String(), page, pageSize, search, order, filter, queryParam.Status, queryParam.Kind, queryParam.Type, name)
	obj := "connections"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		writeMeshkitError(w, ErrQueryGet(obj), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
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
		// Provider implementations return a mix of bare errors and
		// MeshKit-wrapped ones depending on whether the failure was
		// inside DoRequest, in unmarshal, or in the local DAO. Wrap
		// uniformly so the JSON envelope always carries MeshKit
		// metadata.
		wrappedErr := ErrGetConnections(err)
		h.log.Error(wrappedErr)
		writeMeshkitError(w, wrappedErr, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetConnectionByID(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	if connectionID == uuid.Nil {
		invalidIDErr := models.ErrInvalidUUID(fmt.Errorf("invalid connection ID"))
		h.log.Error(invalidIDErr)
		writeMeshkitError(w, invalidIDErr, http.StatusBadRequest)
		return
	}

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	obj := "connection"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		writeMeshkitError(w, ErrQueryGet(obj), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(connection); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		writeMeshkitError(w, models.ErrEncoding(err, obj), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateConnectionById(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := user.ID

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	connection := &connections.ConnectionPayload{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		writeMeshkitError(w, models.ErrUnmarshal(err, obj), http.StatusBadRequest)
		return
	}

	// MeshSync deployment-mode changes are handled by the dedicated
	// POST /api/integrations/connections/{connectionId}/actions endpoint
	// (PerformConnectionAction), which owns the metadata merge and cluster-side
	// redeploy. PUT here only updates connection fields; it no longer sniffs the
	// metadata for a mode change.

	token, err := provider.GetProviderToken(req)
	if err != nil {
		event := eventBuilder.WithSeverity(events.Critical).WithMetadata(map[string]interface{}{
			"error": ErrRetrieveUserToken(err),
		}).WithDescription("No auth token provided in the request.").Build()

		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}

	// A PUT here is a partial update (the UI's connect action sends only
	// {status}); the provider's UpdateConnectionById backfills any omitted field
	// from the persisted row so a partial payload never clobbers columns like
	// kind/name/type/metadata.
	updatedConnection, err := provider.UpdateConnectionById(token, connection, mux.Vars(req)["connectionId"])
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error updating connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusInternalServerError)
		return
	}

	// TODO enhance event with information about meshsync deployment mode change
	description := fmt.Sprintf("Connection %s updated.", updatedConnection.Name)
	eventBuilder = eventBuilder.WithDescription(description)

	if connection.Status != "" {
		event, _ := h.NotifySmOfConnectionStatusChange(req.Context(), userID, provider, token, connection)
		_ = provider.PersistEvent(event, token)
	}

	event := eventBuilder.WithSeverity(events.Informational).Build()
	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)
	h.log.Info(description)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) NotifySmOfConnectionStatusChange(ctx context.Context, userID core.Uuid, provider models.Provider, token string, connection *connections.ConnectionPayload) (events.Event, error) {
	connectionID := connection.ID

	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

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
			ctx,
			connectionID,
			userID,
			smInstanceTracker,
			h.log,
			provider,
			machines.InitialState,
			"kubernetes",
			kubernetes.AssignInitialCtx,
		)

		// A connection being deleted must not leave a tracker entry behind,
		// whichever way its machine failed. InitializeMachineWithContext caches
		// the instance *before* surfacing a Start error, so both the error return
		// below and the no-context return after it would otherwise strand the
		// entry - the first for a fresh failure, the second for every later cache
		// hit. Nothing can drive that machine afterwards, and the delete paths
		// that normally Remove it (the goroutine below, DeleteContext) both do so
		// only after a SendEvent that cannot succeed without a Context.
		usable := err == nil && helpers.HasMachineContext(inst)
		if !usable && connection.Status == connections.DELETED {
			smInstanceTracker.Remove(connectionID)
		}

		if err != nil {
			eventBuilder = eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", connectionID)).WithMetadata(map[string]interface{}{
				"error": err,
			})
			return *eventBuilder.Build(), err
		}
		// A machine whose Context was never assigned cannot service the event:
		// SendEvent would only fail on ErrAssertMachineCtx and publish an error
		// the user can do nothing about. Same shape as the DeleteContext guard in
		// contexts_handler.go; see helpers.HasMachineContext.
		if !usable {
			h.log.Debug(fmt.Sprintf("machine instance for connection %s has no context assigned, skipping the %q event", connectionID, connection.Status))
			return *eventBuilder.Build(), nil
		}

		// detach from the http request lifecycle so that the goroutine isn't cancelled when
		// the handler returns, while preserving context values (e.g. TokenCtxKey) that downstream calls depend on.
		detachedCtx := context.WithoutCancel(ctx)
		go func(inst *machines.StateMachine, status connections.ConnectionStatus) {
			event, err := inst.SendEvent(detachedCtx, machines.EventType(helpers.StatusToEvent(status)), nil)
			if err != nil {
				h.log.Error(err)
				_ = provider.PersistEvent(*event, token)
				h.config.EventBroadcaster.Publish(userID, event)
				return
			}

			if status == connections.DELETED {
				smInstanceTracker.Remove(inst.ID)
			}

			_ = provider.PersistEvent(*event, token)
			h.config.EventBroadcaster.Publish(userID, event)
		}(inst, connection.Status)
	}

	return *eventBuilder.Build(), nil
}

func (h *Handler) DeleteConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := user.ID
	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		writeMeshkitError(w, ErrRetrieveUserToken(err), http.StatusInternalServerError)
		return
	}
	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromOwner(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("delete")

	deletedConnection, err := provider.DeleteConnection(req, connectionID)
	if err != nil {
		obj := "connection"
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error deleting connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(*event, token)
		go h.config.EventBroadcaster.Publish(userID, event)

		if errors.GetCode(err) == models.ErrResultNotFoundCode {
			h.log.Warnf("No connection with ID %q found to delete", connectionID)
			writeMeshkitError(w, _err, http.StatusNotFound)
			return
		}
		h.log.Error(_err)
		writeMeshkitError(w, _err, http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s deleted.", deletedConnection.Name)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()

	_ = provider.PersistEvent(*event, token)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info("connection deleted.")
	w.WriteHeader(http.StatusOK)
}

// handleMeshSyncDeploymentModeChange retrieves existing connection, compares meshsync deployment modes
// between existing and new connections, and performs necessary actions when they differ
// Returns: oldMode, newMode, changed, error
func (h *Handler) handleMeshSyncDeploymentModeChange(
	ctx context.Context,
	connectionID core.Uuid,
	newConnection *connections.ConnectionPayload,
	token string,
	userID core.Uuid,
	provider models.Provider,
) (connections.MeshsyncDeploymentMode, connections.MeshsyncDeploymentMode, bool, error) {
	if newConnection == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("new connection is nil, cannot compare meshsync deployment modes")
	}

	if h.SystemID == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("system ID is not configured in handler")
	}
	// TODO is h.SystemID a correct instance id here?
	mesheryInstanceID := *h.SystemID

	// Retrieve existing connection for mode comparison
	existingConnection, statusCode, err := provider.GetConnectionByID(token, connectionID)
	if err != nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("failed to retrieve existing connection (status %d): %w", statusCode, err)
	}

	if existingConnection == nil {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("existing connection is nil, cannot compare meshsync deployment modes")
	}

	if existingConnection.Kind != "kubernetes" {
		return connections.MeshsyncDeploymentModeUndefined, connections.MeshsyncDeploymentModeUndefined, false, fmt.Errorf("connection is not of kind kubernetes")
	}

	existingMeshSyncMode := connections.MeshsyncDeploymentModeFromMetadata(existingConnection.Metadata)
	newMeshSyncMode := connections.MeshsyncDeploymentModeFromMetadata(newConnection.MetaData)

	// draw back to default mode
	if newMeshSyncMode == connections.MeshsyncDeploymentModeUndefined {
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
				AddMeshsyncDataHandlers(ctx, machineCtx.K8sContext, userID, mesheryInstanceID, provider)
		}

	}

	return existingMeshSyncMode, newMeshSyncMode, meshSyncModeChanged, nil
}
