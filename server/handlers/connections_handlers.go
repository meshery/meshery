package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/helpers"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
)

type connectionStatusPayload map[uuid.UUID]connections.ConnectionStatus

func (h *Handler) ProcessConnectionRegistration(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if req.Method == http.MethodDelete {
		h.handleProcessTermination(w, req)
		return
	}

	connectionRegisterPayload := connections.ConnectionPayload{}
	userUUID := uuid.FromStringOrNil(user.ID)
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
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userUUID, event)
		}

		event, err := inst.SendEvent(req.Context(), machines.EventType(connectionRegisterPayload.Status), connectionRegisterPayload)
		if err != nil {
			h.log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			_ = provider.PersistEvent(event)
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
	userID := uuid.FromStringOrNil(user.ID)
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

	_, err = provider.SaveConnection(&connection, "", false)
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Error creating connection %s", connection.Name)).WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s created.", connection.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(event)
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
// responses:
// 200: mesheryConnectionsResponseWrapper
func (h *Handler) GetConnections(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSizeStr := q.Get("pagesize")
	filter := q.Get("filter")

	var pageSize int
	if pageSizeStr == "all" {
		pageSize = 100
	} else {
		pageSize, _ = strconv.Atoi(pageSizeStr)
	}

	if pageSize > 50 {
		pageSize = 50
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
	}{}

	status := q.Get("status")
	kind := q.Get("kind")
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

	connectionsPage, err := provider.GetConnections(req, user.ID, page, pageSize, search, order, filter, queryParam.Status, queryParam.Kind)
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

// swagger:route GET /api/integrations/connections/{connectionKind} GetConnectionsByKind idGetConnectionsByKind
// Handle GET request for getting all connections for a given kind.
//
// ```?order={field}``` orders on the passed field
//
// ```?search={}``` If search is non empty then a greedy search is performed
//
// ```?page={page-number}``` Default page number is 0
//
// ```?pagesize={pagesize}``` Default pagesize is 10
// responses:
// 200: ConnectionPage
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

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s", page+1, pageSize, search, order))

	connectionsPage, err := provider.GetConnectionsByKind(req, user.ID, page, pageSize, search, order, connectionKind)
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

// swagger:route GET /api/integrations/connections/status GetConnectionsStatus idGetConnectionsStatus
// Handle GET request for getting all connections status
//
// Get all connections status
// responses:
// 200: mesheryConnectionsStatusPage
func (h *Handler) GetConnectionsStatus(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionsStatusPage, err := provider.GetConnectionsStatus(req, user.ID)
	obj := "connections status"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsStatusPage); err != nil {
		h.log.Error(models.ErrEncoding(err, obj))
		http.Error(w, models.ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateConnectionStatus(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionStatusPayload := &connectionStatusPayload{}
	defer func() {
		_ = req.Body.Close()
	}()

	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().FromSystem(*h.SystemID).FromUser(userID).WithCategory("connection").WithAction("update").ActedUpon(userID)

	err := json.NewDecoder(req.Body).Decode(connectionStatusPayload)
	if err != nil {
		errUnmarshal := models.ErrUnmarshal(err, "connection status payload")
		eventBuilder.WithSeverity(events.Error).WithDescription("Unable to update connection status.").
			WithMetadata(map[string]interface{}{
				"error": errUnmarshal,
			})
		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
		return
	}

	connKind := mux.Vars(req)["connectionKind"]
	if connKind == "kubernetes" {
		smInstanceTracker := h.ConnectionToStateMachineInstanceTracker
		token, _ := req.Context().Value(models.TokenCtxKey).(string)
		for id, status := range *connectionStatusPayload {
			eventBuilder.ActedUpon(id)
			k8scontext, err := provider.GetK8sContext(token, id.String())

			if err != nil {
				eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", id)).WithMetadata(map[string]interface{}{
					"error": err,
				})
				_event := eventBuilder.Build()
				_ = provider.PersistEvent(_event)
				go h.config.EventBroadcaster.Publish(userID, _event)
				continue
			}

			event := eventBuilder.WithSeverity(events.Informational).
				WithDescription(fmt.Sprintf("Processing status update to \"%s\" for connection %s", status, k8scontext.Name)).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userID, event)

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
				req.Context(),
				id,
				userID,
				smInstanceTracker,
				h.log,
				provider,
				machines.InitialState,
				"kubernetes",
				kubernetes.AssignInitialCtx,
			)

			if err != nil {
				event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", id)).WithMetadata(map[string]interface{}{
					"error": err,
				}).Build()
				_ = provider.PersistEvent(event)
				go h.config.EventBroadcaster.Publish(userID, event)
				continue
			}

			go func(inst *machines.StateMachine, status connections.ConnectionStatus) {
				event, err = inst.SendEvent(req.Context(), machines.EventType(helpers.StatusToEvent(status)), nil)
				if err != nil {
					h.log.Error(err)
					_ = provider.PersistEvent(event)
					h.config.EventBroadcaster.Publish(userID, event)
					return
				}

				if status == connections.DELETED {
					smInstanceTracker.Remove(inst.ID)
				}

				_ = provider.PersistEvent(event)
				h.config.EventBroadcaster.Publish(userID, event)
			}(inst, status)
		}
	} else {
		token, _ := req.Context().Value(models.TokenCtxKey).(string)
		for id, status := range *connectionStatusPayload {
			connection, statusCode, err := provider.UpdateConnectionStatusByID(token, id, status)

			if err != nil {
				event := events.NewEvent().WithCategory("connection").WithAction("update").WithSeverity(events.Error).ActedUpon(id).FromUser(userID).FromSystem(*h.SystemID).WithDescription(fmt.Sprintf("Unable to update connection status to %s", status)).WithMetadata(map[string]interface{}{"error": err}).Build()
				_ = provider.PersistEvent(event)
				h.config.EventBroadcaster.Publish(userID, event)
				h.log.Error(err)
				continue
			}
			eb := events.NewEvent()
			eb.WithDescription(fmt.Sprintf("Connection \"%s\" status updated to %s", connection.Name, connection.Status)).WithAction("update")
			if status == connections.DELETED {
				eb.WithDescription(fmt.Sprintf("Connection \"%s\" deleted", connection.Name)).WithAction("delete")
			}
			event := eb.WithCategory("connection").WithSeverity(events.Success).FromUser(userID).FromSystem(*h.SystemID).ActedUpon(id).Build()
			_ = provider.PersistEvent(event)
			h.config.EventBroadcaster.Publish(userID, event)

			h.log.Debug("connection", connection, statusCode)
		}

	}
	w.WriteHeader(http.StatusAccepted)
}

// swagger:route PUT /api/integrations/connections/{connectionKind} PutConnection idPutConnection
// Handle PUT request for updating an existing connection
//
// Updates existing connection
// responses:
// 200: noContentWrapper
func (h *Handler) UpdateConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	userID := uuid.FromStringOrNil(user.ID)

	connection := &connections.Connection{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	connectionID := connection.ID
	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("update")

	updatedConnection, err := provider.UpdateConnection(req, connection)
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error updating connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s updated.", updatedConnection.Name)

	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info(description)
	w.WriteHeader(http.StatusOK)
}

// swagger:route PUT /api/integrations/connections/{connectionId} PutConnectionById idPutConnectionById
// Handle PUT request for updating an existing connection by connection ID
//
// Updates existing connection using ID
// responses:
// 200: mesheryConnectionResponseWrapper
func (h *Handler) UpdateConnectionById(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := uuid.FromStringOrNil(user.ID)

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

	updatedConnection, err := provider.UpdateConnectionById(req, connection, mux.Vars(req)["connectionId"])
	if err != nil {
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error updating connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s updated.", updatedConnection.Name)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()

	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
	h.log.Info(description)
	w.WriteHeader(http.StatusOK)
}

// swagger:route DELETE /api/integrations/connections/{connectionId} DeleteConnection idDeleteConnection
// Handle DELETE request for deleting an existing connection by connection ID
//
// Deletes existing connection
// responses:
// 200: noContentWrapper
func (h *Handler) DeleteConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	userID := uuid.FromStringOrNil(user.ID)
	eventBuilder := events.NewEvent().ActedUpon(connectionID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("delete")

	deletedConnection, err := provider.DeleteConnection(req, connectionID)
	if err != nil {
		obj := "connection"
		_err := ErrFailToSave(err, obj)
		metadata := map[string]interface{}{
			"error": _err,
		}
		event := eventBuilder.WithSeverity(events.Error).WithDescription("Error deleting connection").WithMetadata(metadata).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)

		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

	description := fmt.Sprintf("Connection %s deleted.", deletedConnection.Name)
	event := eventBuilder.WithSeverity(events.Informational).WithDescription(description).Build()

	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)

	h.log.Info("connection deleted.")
	w.WriteHeader(http.StatusOK)
}
