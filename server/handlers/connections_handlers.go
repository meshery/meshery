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
)

// swagger:route POST /api/integrations/connections PostConnection idPostConnection
// Handle POST request for creating a new connection
//
// Creates a new connection
// responses:
// 201: noContentWrapper
func (h *Handler) SaveConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	connection := models.ConnectionPayload{}
	err = json.Unmarshal(bd, &connection)
	obj := "connection"

	if err != nil {
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(w, ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	err = provider.SaveConnection(req, &connection, "", false)
	if err != nil {
		h.log.Error(ErrFailToSave(err, obj))
		http.Error(w, ErrFailToSave(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("connection saved successfully")
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
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
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
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}


// swagger:route GET /api/integrations/connections/status GetConnectionsStatus idGetConnectionsStatus
// Handle GET request for getting all connections status
//
// responses:
// 200: ConnectionsStatusPage
func (h *Handler) GetConnectionsStatus(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	connectionsStatusPage, err := provider.GetConnectionsStatus(req, user.ID)
	obj := "connections status"

	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(connectionsStatusPage); err != nil {
		h.log.Error(ErrEncoding(err, obj))
		http.Error(w, ErrEncoding(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}


// swagger:route PUT /api/integrations/connections/{connectionKind} PutConnection idPutConnection
// Handle PUT request for updating an existing connection
//
// Updates existing connection
// responses:
// 200: noContentWrapper
func (h *Handler) UpdateConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	connection := &models.Connection{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(w, ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	_, err = provider.UpdateConnection(req, connection)
	if err != nil {
		h.log.Error(ErrQueryGet(obj))
		http.Error(w, ErrQueryGet(obj).Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("connection updated successfully")
	w.WriteHeader(http.StatusOK)
}

// swagger:route PUT /api/integrations/connections/{connectionId} PutConnectionById idPutConnectionById
// Handle PUT request for updating an existing connection by connection ID
//
// Updates existing connection using ID
// responses:
// 200: mesheryConnectionResponseWrapper
func (h *Handler) UpdateConnectionById(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	connection := &models.ConnectionPayload{}
	err = json.Unmarshal(bd, connection)
	obj := "connection"
	if err != nil {
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(w, ErrUnmarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	_, err = provider.UpdateConnectionById(req, connection, mux.Vars(req)["connectionId"])
	if err != nil {
		h.log.Error(ErrFailToSave(err, obj))
		http.Error(w, ErrFailToSave(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("connection updated successfully")
	w.WriteHeader(http.StatusOK)
}


// swagger:route DELETE /api/integrations/connections/{connectionId} DeleteConnection idDeleteConnection
// Handle DELETE request for deleting an existing connection by connection ID
//
// Deletes existing connection
// responses:
// 200: noContentWrapper
func (h *Handler) DeleteConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionId"])
	_, err := provider.DeleteConnection(req, connectionID)
	if err != nil {
		obj := "connection"
		h.log.Error(ErrFailToDelete(err, obj))
		http.Error(w, ErrFailToDelete(err, obj).Error(), http.StatusInternalServerError)
		return
	}

	h.log.Info("connection deleted successfully")
	w.WriteHeader(http.StatusOK)
}
