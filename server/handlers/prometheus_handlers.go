package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/connections"
	"github.com/layer5io/meshkit/models/events"
)

func init() {
	gob.Register(&models.PrometheusClient{})
}

// swagger:route GET /api/system/meshsync/grafana SystemAPI idMeshSyncGrafana
// Handle GET request for mesh-sync grafana
//
// Fetches Prometheus and Grafana
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanPromGrafanaHandler - fetches  Prometheus and Grafana
// func (h *Handler) ScanPromGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanPromGrafana(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/system/meshsync/prometheus SystemAPI idMeshSyncPrometheus
// Handle GET request for fetching prometheus
//
// Fetches Prometheus
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanPrometheusHandler - fetches  Prometheus
// func (h *Handler) ScanPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanPrometheus(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/telemetry/metrics/grafana/scan GrafanaAPI idGetGrafana
// Handle GET request for Grafana
//
// Fetches and returns Grafana
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanGrafanaHandler - fetches  Grafana
// func (h *Handler) ScanGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	errs := []string{}
// 	var wg sync.WaitGroup
// 	customK8scontexts, ok := req.Context().Value(models.KubeClustersKey).([]models.K8sContext)
// 	if ok && len(customK8scontexts) > 0 {
// 		for _, mk8scontext := range customK8scontexts {
// 			wg.Add(1)
// 			go func(mk8scontext models.K8sContext) {
// 				defer wg.Done()
// 				k8sconfig, err := mk8scontext.GenerateKubeConfig()
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				availablePromGrafana, err := helpers.ScanGrafana(k8sconfig, mk8scontext.Name)
// 				if err != nil {
// 					errs = append(errs, err.Error())
// 					h.log.Error(err)
// 					return
// 				}
// 				if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
// 					obj := "payloads"
// 					h.log.Error(ErrMarshal(err, obj))
// 					errs = append(errs, ErrMarshal(err, obj).Error())
// 					return
// 				}
// 			}(mk8scontext)
// 		}
// 	}
// 	if len(errs) != 0 {
// 		http.Error(w, mergeMsgs(errs), http.StatusInternalServerError)
// 	}
// 	wg.Wait()
// }

// swagger:route GET /api/telemetry/metrics/config PrometheusAPI idGetPrometheusConfig
// Handle GET for Prometheus configuration
//
// Used for fetching Prometheus configuration
// responses:
//  200: prometheusConfigResponseWrapper

// swagger:route POST /api/telemetry/metrics/config PrometheusAPI idPostPrometheusConfig
// Handle POST for Prometheus configuration
//
// Used for persisting Prometheus configuration
// responses:
//  200:

// swagger:route DELETE /api/telemetry/metrics/config PrometheusAPI idDeletePrometheusConfig
// Handle DELETE for Prometheus configuration
//
// Used for deleting Prometheus configuration
// responses:
//  200:

// PrometheusConfigHandler is used for persisting prometheus configuration
func (h *Handler) PrometheusConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if req.Method == http.MethodGet {
		req = mux.SetURLVars(req, map[string]string{"connectionKind": "prometheus"})
		h.GetConnectionsByKind(w, req, prefObj, user, provider)
		return
	}

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	sysID := h.SystemID
	userUUID := uuid.FromStringOrNil(user.ID)

	eventBuilder := events.NewEvent().ActedUpon(userUUID).WithCategory("connection").WithAction("update").FromSystem(*sysID).FromUser(userUUID).WithDescription("Failed to interact with the connection.")

	if req.Method == http.MethodPost {
		promURL := req.FormValue("prometheusURL")
		promKey := req.FormValue("prometheusKey")
		connName := req.FormValue("prometheusConnectionName")
		credName := req.FormValue("prometheusCredentialName")

		u, err := url.Parse(promURL)
		if err != nil {
			return
		}
		if strings.Contains(promURL, u.RequestURI()) {
			promURL = strings.TrimSuffix(promURL, u.RequestURI())
		}

		promConn := map[string]interface{}{
			"url": promURL,
		}

		promCred := map[string]interface{}{
			"secret": promKey,
		}

		if err := h.config.PrometheusClient.Validate(req.Context(), promURL, promKey); err != nil {
			h.log.Error(models.ErrPrometheusScan(err))
			http.Error(w, models.ErrPrometheusScan(err).Error(), http.StatusInternalServerError)
			return
		}

		userUUID := uuid.FromStringOrNil(user.ID)
		credential, err := provider.SaveUserCredential(token, &models.Credential{
			UserID: &userUUID,
			Type:   "prometheus",
			Secret: promCred,
			Name:   credName,
		})
		if err != nil {
			_err := models.ErrPersistCredential(err)
			event := eventBuilder.WithDescription(fmt.Sprintf("Unable to persist credential information for the connection %s", credName)).
				WithSeverity(events.Error).WithMetadata(map[string]interface{}{"error": _err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userUUID, event)
			http.Error(w, _err.Error(), http.StatusInternalServerError)
			return
		}
		connection, err := provider.SaveConnection(&connections.ConnectionPayload{
			Kind:             "prometheus",
			Type:             "observability",
			SubType:          "monitoring",
			Status:           connections.CONNECTED,
			MetaData:         promConn,
			CredentialSecret: promCred,
			Name:             connName,
			CredentialID:     &credential.ID,
		}, token, false)

		if err != nil {
			_err := models.ErrPersistConnection(err)
			event := eventBuilder.WithDescription(fmt.Sprintf("Unable to perisit the \"%s\" connection details", connName)).WithMetadata(map[string]interface{}{"error": _err}).Build()
			_ = provider.PersistEvent(event)
			go h.config.EventBroadcaster.Publish(userUUID, event)
			http.Error(w, _err.Error(), http.StatusInternalServerError)
			return
		}
		event := eventBuilder.WithDescription(fmt.Sprintf("Connection %s with Prometheus created at %s", connName, promURL)).WithSeverity(events.Success).ActedUpon(connection.ID).Build()
		_ = provider.PersistEvent(event)
		go h.config.EventBroadcaster.Publish(userUUID, event)

		h.log.Debug("Prometheus URL %s saved", promURL)
	} else if req.Method == http.MethodDelete {
		http.Error(w, "API is deprecated, please use connections API", http.StatusGone)
		return
	}

	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// swagger:route GET /api/telemetry/metrics/ping/{connectionID} PrometheusAPI idGetPrometheusPing
// Handle GET request for Prometheus Ping
//
// Used to ping prometheus
// responses:
// 	200:

// PrometheusPingHandler - fetches server version to simulate ping
func (h *Handler) PrometheusPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := p.GetConnectionByIDAndKind(token, connectionID, "prometheus")
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

	if err := h.config.PrometheusClient.Validate(req.Context(), url, apiKeyOrBasicAuth); err != nil {
		h.log.Error(models.ErrPrometheusScan(err))
		http.Error(w, models.ErrPrometheusScan(err).Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// swagger:route POST /api/telemetry/metrics/board_import PrometheusAPI idPostPrometheusBoardImport
// Handle POST request for Prometheus board import
//
// Used for importing Grafana board for Prometheus
// responses:
// 	200: prometheusBoardImportRespWrapper

// GrafanaBoardImportForPrometheusHandler accepts a Grafana board json, parses it and returns the list of panels
func (h *Handler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, _ models.Provider) {
	defer func() {
		_ = req.Body.Close()
	}()

	boardData, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	board, err := h.config.PrometheusClient.ImportGrafanaBoard(req.Context(), boardData)
	if err != nil {
		h.log.Error(ErrPrometheusBoards(err))
		http.Error(w, ErrPrometheusBoards(err).Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(board)
	if err != nil {
		obj := "board instance"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/telemetry/metrics/query/{connectionID}  PrometheusAPI idGetPrometheusQuery
// Handle GET request for Prometheus Query
//
// Used to prometheus queries
// responses:
// 	200:

// PrometheusQueryHandler handles prometheus queries
func (h *Handler) PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, p models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := p.GetConnectionByIDAndKind(token, connectionID, "prometheus")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	url, _ := connection.Metadata["url"].(string)

	reqQuery := req.URL.Query()

	data, err := h.config.PrometheusClientForQuery.Query(req.Context(), url, &reqQuery)
	if err != nil {
		h.log.Error(ErrPrometheusQuery(err))
		http.Error(w, ErrPrometheusQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// PrometheusQueryRangeHandler handles prometheus range queries
func (h *Handler) PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := provider.GetConnectionByIDAndKind(token, connectionID, "prometheus")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	reqQuery := req.URL.Query()
	url, _ := connection.Metadata["url"].(string)

	testUUID := reqQuery.Get("uuid")
	if testUUID != "" {
		q := reqQuery.Get("query")
		h.config.QueryTracker.AddOrFlagQuery(req.Context(), testUUID, q, false)
	}

	data, err := h.config.PrometheusClientForQuery.QueryRange(req.Context(), url, &reqQuery)
	if err != nil {
		h.log.Error(ErrPrometheusQuery(err))
		http.Error(w, ErrPrometheusQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// swagger:route GET /api/telemetry/metrics/static-board PrometheusAPI idGetPrometheusStaticBoard
// Handle GET request for Prometheus static board
//
// Used to fetch the static board
// responses:
// 	200: prometheusStaticBoardRespWrapper

// PrometheusStaticBoardHandler returns the static board
func (h *Handler) PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, _ *models.User, provider models.Provider) {
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	connectionID := uuid.FromStringOrNil(mux.Vars(req)["connectionID"])

	connection, statusCode, err := provider.GetConnectionByIDAndKind(token, connectionID, "prometheus")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}
	url, _ := connection.Metadata["url"].(string)

	result := map[string]*models.GrafanaBoard{}
	resultLock := &sync.Mutex{}
	resultWG := &sync.WaitGroup{}

	boardFunc := map[string]func(context.Context, string) (*models.GrafanaBoard, error){
		"cluster": h.config.PrometheusClient.GetClusterStaticBoard,
		"node":    h.config.PrometheusClient.GetNodesStaticBoard,
	}

	for key, bfunc := range boardFunc {
		resultWG.Add(1)
		go func(k string, bfun func(context.Context, string) (*models.GrafanaBoard, error)) {
			defer resultWG.Done()

			board, err := bfun(req.Context(), url)
			if err != nil {
				// error is already logged
				return
			}
			resultLock.Lock()
			defer resultLock.Unlock()
			result[k] = board
		}(key, bfunc)
	}
	resultWG.Wait()

	if len(result) != len(boardFunc) {
		h.log.Error(ErrStaticBoards)
		http.Error(w, ErrStaticBoards.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		obj := "board instance"
		h.log.Error(models.ErrMarshal(err, obj))
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /api/telemetry/metrics/boards/{connectionID} PrometheusAPI idPostPrometheusBoard
// Handle POST request for Prometheus board
//
// Used to persist selected board and panels
// responses:
// 	200:

// SaveSelectedPrometheusBoardsHandler persists selected board and panels
func (h *Handler) SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

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
	connection, statusCode, err := provider.GetConnectionByIDAndKind(token, connectionID, "prometheus")
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	if len(boards) > 0 {
		connection.Metadata["prometheus_boards"] = boards
	} else {
		delete(connection.Metadata, "prometheus_boards")
	}
	updatedConnection, err := provider.UpdateConnection(req, connection)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Debug("Board selection updated", updatedConnection.Metadata)
	_, _ = w.Write([]byte("{}"))
}
