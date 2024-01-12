// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"

	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&models.GrafanaClient{})
}

// swagger:route GET /api/telemetry/metrics/grafana/config GrafanaAPI idGetGrafanaConfig
// Handle GET request for Grafana configuration
//
// Used for fetching Grafana configuration
// responses:
// 	200: grafanaConfigResponseWrapper

// swagger:route POST /api/telemetry/metrics/grafana/config GrafanaAPI idPostGrafanaConfig
// Handle POST request for Grafana configuration
//
// Used for persisting Grafana configuration
// responses:
// 	200:

// swagger:route DELETE /api/telemetry/metrics/grafana/config GrafanaAPI idDeleteGrafanaConfig
// Handle DELETE request for Grafana configuration
//
// Used for Delete Grafana configuration
// responses:
// 	200:

// GrafanaConfigHandler is used for fetching or persisting or removing Grafana configuration
func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	sysID := h.SystemID
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	if req.Method == http.MethodGet {
		req = mux.SetURLVars(req, map[string]string{"connectionKind": "grafana"})
		h.GetConnectionsByKind(w, req, prefObj, user, p)
		return
	}

	if req.Method == http.MethodPost {
		grafanaURL := req.FormValue("grafanaURL")
		grafanaAPIKey := req.FormValue("grafanaAPIKey")
		credName := req.FormValue("grafanaCredentialName")
		connName := req.FormValue("grafanaConnectioncName")
		u, err := url.Parse(grafanaURL)
		if err != nil {
			return
		}
		if strings.Contains(grafanaURL, u.RequestURI()) {
			grafanaURL = strings.TrimSuffix(grafanaURL, u.RequestURI())
		}

		grafanaConn := map[string]interface{}{
			"url": grafanaURL,
		}
		grafanaCred := map[string]interface{}{
			"secret": grafanaAPIKey,
		}

		if err := h.config.GrafanaClient.Validate(req.Context(), grafanaURL, grafanaAPIKey); err != nil {
			h.log.Error(models.ErrGrafanaScan(err))
			http.Error(w, models.ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
			return
		}

		userUUID := uuid.FromStringOrNil(user.ID)
		credential, err := p.SaveUserCredential(token, &models.Credential{
			UserID: &userUUID,
			Type:   "grafana",
			Secret: grafanaCred,
			Name:   credName,
		})
		if err != nil {
			_err := models.ErrPersistCredential(err)
			event := eventBuilder.WithDescription(fmt.Sprintf("Unable to persist credential information for the connection %s", credName)).
				WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": _err}).Build()
			_ = p.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userUUID, event)
			http.Error(w, _err.Error(), http.StatusInternalServerError)
			return
		}
		connection, err := p.SaveConnection(&models.ConnectionPayload{
			Kind:             "grafana",
			Type:             "observability",
			SubType:          "monitoring",
			Status:           connections.CONNECTED,
			MetaData:         grafanaConn,
			CredentialSecret: grafanaCred,
			Name:             connName,
			CredentialID:     &credential.ID,
		}, token, false)

		if err != nil {
			_err := models.ErrPersistConnection(err)
			event := eventBuilder.WithDescription(fmt.Sprintf("Unable to perisit the \"%s\" connection details", connName)).WithMetadata(map[string]interface{}{"error": _err}).Build()
			_ = p.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userUUID, event)
			http.Error(w, _err.Error(), http.StatusInternalServerError)
			return
		}
		event := eventBuilder.WithDescription(fmt.Sprintf("Connection %s with grafana created at %s", connName, grafanaURL)).WithSeverity(events.Success).ActedUpon(connection.ID).Build()
		_ = p.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userUUID, event)

		logrus.Debugf("connection to grafana @ %s succeeded", grafanaURL)

		_ = json.NewEncoder(w).Encode(connection)
	} else if req.Method == http.MethodDelete {
		http.Error(w, "API is deprecated, please use connections API", http.StatusGone)
		return
	}

}

// swagger:route GET /api/telemetry/metrics/grafana/ping/{connectionID} GrafanaAPI idGetGrafanaPing
// Handle GET request for Grafana ping
//
// Used to initiate a Grafana ping
// responses:
// 	200:

// GrafanaPingHandler - used to initiate a Grafana ping
func (h *Handler) GrafanaPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := p.GetConnectionByID(token, connectionID, "grafana")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	url, _ := connection.Metadata["url"].(string)
	cred, statusCode, err := p.GetCredentialByID(token, connection.CredentialID)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	apiKeyOrBasicAuth, _ := cred.Secret["secret"].(string)
	if err := h.config.GrafanaClient.Validate(req.Context(), url, apiKeyOrBasicAuth); err != nil {
		h.log.Error(models.ErrGrafanaScan(err))
		http.Error(w, models.ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// swagger:route GET /api/telemetry/metrics/grafana/boards/{connectionID} GrafanaAPI idGetGrafanaBoards
// Handle GET request for Grafana boards
//
// Used for fetching Grafana boards and panels
// responses:
// 	200: grafanaBoardsResponseWrapper

// GrafanaBoardsHandler is used for fetching Grafana boards and panels
func (h *Handler) GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	if req.Method == http.MethodPost {
		h.SaveSelectedGrafanaBoardsHandler(w, req, prefObj, user, p)
		return
	}
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])
	connection, statusCode, err := p.GetConnectionByID(token, connectionID, "grafana")
	fmt.Println("CONNECTION ID : ", connectionID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), statusCode)
		return
	}

	url, _ := connection.Metadata["url"].(string)
	cred, statusCode, err := p.GetCredentialByID(token, connection.CredentialID)
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), statusCode)
		return
	}
	apiKeyOrBasicAuth, _ := cred.Secret["secret"].(string)
	fmt.Println(apiKeyOrBasicAuth, "GRAFANA KEY", cred.Secret)
	if err := h.config.GrafanaClient.Validate(req.Context(), url, apiKeyOrBasicAuth); err != nil {
		h.log.Error(models.ErrGrafanaScan(err))
		http.Error(w, models.ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	boards, err := h.config.GrafanaClient.GetGrafanaBoards(req.Context(), url, apiKeyOrBasicAuth, dashboardSearch)
	if err != nil {
		h.log.Error(models.ErrGrafanaBoards(err))
		http.Error(w, models.ErrGrafanaBoards(err).Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		obj := "boards payload"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/telemetry/metrics/grafana/query/{connectionID} GrafanaAPI idGetGrafanaQuery
// Handle GET request for Grafana queries
//
// Used for handling Grafana queries
// responses:
// 	200:

// GrafanaQueryHandler is used for handling Grafana queries
func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, p models.Provider) {

	reqQuery := req.URL.Query()
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])
	connection, statusCode, err := p.GetConnectionByID(token, connectionID, "grafana")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	url, _ := connection.Metadata["url"].(string)
	cred, statusCode, err := p.GetCredentialByID(token, connection.CredentialID)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	apiKeyOrBasicAuth, _ := cred.Secret["secret"].(string)
	if err := h.config.GrafanaClient.Validate(req.Context(), url, apiKeyOrBasicAuth); err != nil {
		h.log.Error(models.ErrGrafanaScan(err))
		http.Error(w, models.ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	data, err := h.config.GrafanaClientForQuery.GrafanaQuery(req.Context(), url, apiKeyOrBasicAuth, &reqQuery)
	if err != nil {
		h.log.Error(ErrGrafanaQuery(err))
		http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// GrafanaQueryRangeHandler is used for handling Grafana Range queries
func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	reqQuery := req.URL.Query()

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])
	connection, statusCode, err := provider.GetConnectionByID(token, connectionID, "grafana")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	url, _ := connection.Metadata["url"].(string)
	cred, statusCode, err := provider.GetCredentialByID(token, connection.CredentialID)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	apiKeyOrBasicAuth, _ := cred.Secret["secret"].(string)
	if err := h.config.GrafanaClient.Validate(req.Context(), url, apiKeyOrBasicAuth); err != nil {
		h.log.Error(models.ErrGrafanaScan(err))
		http.Error(w, models.ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	data, err := h.config.GrafanaClientForQuery.GrafanaQueryRange(req.Context(), url, apiKeyOrBasicAuth, &reqQuery)
	if err != nil {
		h.log.Error(ErrGrafanaQuery(err))
		http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// swagger:route POST /api/telemetry/metrics/grafana/boards/{connectionID} GrafanaAPI idPostGrafanaBoards
// Handle POST request for Grafana boards
//
// Used for persist Grafana boards and panel selections
// responses:
// 	200:

// SaveSelectedGrafanaBoardsHandler is used to persist board and panel selection
func (h *Handler) SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	boards := []*models.SelectedGrafanaConfig{}
	err = json.Unmarshal(body, &boards)
	if err != nil {
		obj := "request body"
		h.log.Error(models.ErrUnmarshal(err, obj))
		http.Error(w, models.ErrUnmarshal(err, obj).Error(), http.StatusBadRequest)
		return
	}

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])
	connection, statusCode, err := p.GetConnectionByID(token, connectionID, "grafana")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	if len(boards) > 0 {
		connection.Metadata["grafana_boards"] = boards
	} else {
		delete(connection.Metadata, "grafana_boards")
	}

	updatedConnection, err := p.UpdateConnection(req, connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.log.Debug("Board selection updated", updatedConnection.Metadata)
	_, _ = w.Write([]byte("{}"))
}
