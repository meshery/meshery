package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
)

func init() {
	gob.Register(&models.PrometheusClient{})
}

// ScanPromGrafanaHandler - fetches  Prometheus and Grafana
func (h *Handler) ScanPromGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		err := ErrInvalidK8SConfig
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	availablePromGrafana, err := helpers.ScanPromGrafana(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
	if err != nil {
		h.log.Error(err)
		http.Error(w, "unable to scan Kubernetes", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(availablePromGrafana); err != nil {
		obj := "payloads"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// ScanPrometheusHandler - fetches  Prometheus
func (h *Handler) ScanPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		err := ErrInvalidK8SConfig
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	availablePrometheus, err := helpers.ScanPrometheus(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(availablePrometheus); err != nil {
		obj := "payloads"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/telemetry/metrics/grafana/scan GrafanaAPI idGetGrafana
// Handle GET request for Grafana
//
// Fetches and returns Grafana
// responses:
// 	200: v1ServicesMapResponseWrapper

// ScanGrafanaHandler - fetches  Grafana
func (h *Handler) ScanGrafanaHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		err := ErrInvalidK8SConfig
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	availableGrafana, err := helpers.ScanGrafana(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
	if err != nil {
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(availableGrafana); err != nil {
		obj := "payloads"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

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
		err := json.NewEncoder(w).Encode(prefObj.Prometheus)
		if err != nil {
			obj := "Prometheus config"
			h.log.Error(ErrMarshal(err, obj))
			http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if req.Method == http.MethodPost {
		promURL := req.FormValue("prometheusURL")
		if err := h.config.PrometheusClient.Validate(req.Context(), promURL); err != nil {
			h.log.Error(ErrPrometheusScan(err))
			http.Error(w, ErrPrometheusScan(err).Error(), http.StatusInternalServerError)
			return
		}

		u, err := url.Parse(promURL)
		if err != nil {
			return
		}
		if strings.Contains(promURL, u.RequestURI()) {
			promURL = strings.TrimSuffix(promURL, u.RequestURI())
		}

		prefObj.Prometheus = &models.Prometheus{
			PrometheusURL: promURL,
		}
		h.log.Debug("Prometheus URL %s successfully saved", promURL)
	} else if req.Method == http.MethodDelete {
		prefObj.Prometheus = nil
	}

	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// swagger:route GET /api/telemetry/metrics/ping PrometheusAPI idGetPrometheusPing
// Handle GET request for Prometheus Ping
//
// Used to ping prometheus
// responses:
// 	200:

// PrometheusPingHandler - fetches server version to simulate ping
func (h *Handler) PrometheusPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Prometheus == nil || prefObj.Prometheus.PrometheusURL == "" {
		h.log.Error(ErrPrometheusConfig)
		http.Error(w, ErrPrometheusConfig.Error(), http.StatusBadRequest)
		return
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		err := ErrInvalidK8SConfig
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.config.PrometheusClient.Validate(req.Context(), prefObj.Prometheus.PrometheusURL); err != nil {
		h.log.Error(ErrPrometheusScan(err))
		http.Error(w, ErrPrometheusScan(err).Error(), http.StatusInternalServerError)
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
func (h *Handler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Prometheus == nil || prefObj.Prometheus.PrometheusURL == "" {
		h.log.Error(ErrPrometheusConfig)
		http.Error(w, ErrPrometheusConfig.Error(), http.StatusBadRequest)
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	boardData, err := ioutil.ReadAll(req.Body)
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
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/telemetry/metrics/query PrometheusAPI idGetPrometheusQuery
// Handle GET request for Prometheus Query
//
// Used to prometheus queries
// responses:
// 	200:

// PrometheusQueryHandler handles prometheus queries
func (h *Handler) PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Prometheus == nil || prefObj.Prometheus.PrometheusURL == "" {
		h.log.Error(ErrPrometheusConfig)
		http.Error(w, ErrPrometheusConfig.Error(), http.StatusBadRequest)
		return
	}

	reqQuery := req.URL.Query()

	data, err := h.config.PrometheusClientForQuery.Query(req.Context(), prefObj.Prometheus.PrometheusURL, &reqQuery)
	if err != nil {
		h.log.Error(ErrPrometheusQuery(err))
		http.Error(w, ErrPrometheusQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// PrometheusQueryRangeHandler handles prometheus range queries
func (h *Handler) PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()

	testUUID := reqQuery.Get("uuid")
	if testUUID != "" {
		q := reqQuery.Get("query")
		h.config.QueryTracker.AddOrFlagQuery(req.Context(), testUUID, q, false)
	}

	data, err := h.config.PrometheusClientForQuery.QueryRange(req.Context(), reqQuery.Get("url"), &reqQuery)
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
func (h *Handler) PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Prometheus == nil || prefObj.Prometheus.PrometheusURL == "" {
		_, _ = w.Write([]byte("{}"))
		return
	}

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

			board, err := bfun(req.Context(), prefObj.Prometheus.PrometheusURL)
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

	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		obj := "board instance"
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route POST /api/telemetry/metrics/boards PrometheusAPI idPostPrometheusBoard
// Handle POST request for Prometheus board
//
// Used to persist selected board and panels
// responses:
// 	200:

// SaveSelectedPrometheusBoardsHandler persists selected board and panels
func (h *Handler) SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Prometheus == nil || prefObj.Prometheus.PrometheusURL == "" {
		h.log.Error(ErrPrometheusConfig)
		http.Error(w, ErrPrometheusConfig.Error(), http.StatusBadRequest)
		return
	}

	// if prefObj.Prometheus.SelectedPrometheusBoardsConfigs == nil {
	// 	prefObj.Prometheus.SelectedPrometheusBoardsConfigs = []*models.GrafanaBoard{}
	// }

	defer func() {
		_ = req.Body.Close()
	}()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	boards := []*models.SelectedGrafanaConfig{}
	err = json.Unmarshal(body, &boards)
	if err != nil {
		obj := "request body"
		h.log.Error(ErrUnmarshal(err, obj))
		http.Error(w, ErrUnmarshal(err, obj).Error(), http.StatusBadRequest)
		return
	}
	if len(boards) > 0 {
		prefObj.Prometheus.SelectedPrometheusBoardsConfigs = boards
	} else {
		prefObj.Prometheus.SelectedPrometheusBoardsConfigs = nil
	}
	err = provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
