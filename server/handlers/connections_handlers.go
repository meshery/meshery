package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
)

type connectionStatusPayload map[uuid.UUID]connections.ConnectionStatus

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

	connection := models.ConnectionPayload{}
	err = json.Unmarshal(bd, &connection)
	obj := "connection"

	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	eventBuilder := events.NewEvent().ActedUpon(userID).FromUser(userID).FromSystem(*h.SystemID).WithCategory("connection").WithAction("create")

	_, err = provider.SaveConnection(req, &connection, "", false)
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
// responses:
// 200: ConnectionPage
func (h *Handler) GetConnections(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(q.Get("pagesize"))

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

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s", page+1, pageSize, search, order))

	connectionsPage, err := provider.GetConnections(req, user.ID, page, pageSize, search, order)
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

// swagger:route GET /api/integrations/connections/{connectionKind}/transitions GetAvailableTransitionsByKind idGetConnectionsStatus
// Handle GET request for getting all possible connection transitions
//
// Get all possible state transitions for a particular connection kind.
// responses:
// 200: mesheryConnectionsStatusPage
func (h *Handler) GetPossibleTransitionsByKind(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionKind := mux.Vars(req)["connectionKind"]
	transitions := connections.PossibleTransitionnsMap[connectionKind]

	err := json.NewEncoder(w).Encode(transitions)
	if err != nil {
		http.Error(w, models.ErrMarshal(err, "connection transitions").Error(), http.StatusInternalServerError)
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
	var statusCode int 
	for id, status := range *connectionStatusPayload {
		eventBuilder.ActedUpon(id)
		var updatedConnection *connections.Connection
		var err error
		updatedConnection, statusCode, err = provider.UpdateConnectionStatusByID(req, id, status)
		fmt.Println("tetete_", updatedConnection)
		if err != nil {
			eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Failed to update connection status for %s", id)).WithMetadata(map[string]interface{}{
				"error": err,
			})
			// Get the connection name as id doesn't convey much to the user.
		} else {
			eventBuilder.WithSeverity(events.Success).WithDescription(fmt.Sprintf("Connection \"%s\" status updated to %s for %s", updatedConnection.Name, status, id))
			go h.config.K8scontextChannel.PublishContext()
		}

		event := eventBuilder.Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userID, event)
	}
	w.WriteHeader(statusCode)
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

	connection := &models.ConnectionPayload{}
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

	h.log.Info("connection deleted successfully")
	w.WriteHeader(http.StatusOK)
}
