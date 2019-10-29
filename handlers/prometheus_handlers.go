package handlers

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"

	"github.com/layer5io/meshery/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&helpers.PrometheusClient{})
}

// PrometheusConfigHandler is used for persisting prometheus configuration
func (h *Handler) PrometheusConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodPost && req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if req.Method == http.MethodPost {
		promURL := req.FormValue("prometheusURL")
		if _, err = helpers.NewPrometheusClient(req.Context(), promURL, true); err != nil {
			logrus.Errorf("unable to connect to prometheus: %v", err)
			http.Error(w, "unable to connect to prometheus", http.StatusInternalServerError)
			return
		}
		sessObj.Prometheus = &models.Prometheus{
			PrometheusURL: promURL,
		}
		logrus.Debugf("Prometheus URL %s successfully saved", promURL)
	} else if req.Method == http.MethodDelete {
		sessObj.Prometheus = nil
	}

	err = h.config.SessionPersister.Write(user.UserID, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// GrafanaBoardImportForPrometheusHandler accepts a Grafana board json, parses it and returns the list of panels
func (h *Handler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		http.Error(w, "Prometheus URL is not configured", http.StatusBadRequest)
		return
	}

	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), sessObj.Prometheus.PrometheusURL, false)
	if err != nil {
		msg := "unable to initiate a client to connect to prometheus"
		err = errors.Wrap(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer prometheusClient.Close()

	defer func() {
		_ = req.Body.Close()
	}()
	boardData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := "unable to read the board payload"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	board, err := prometheusClient.ImportGrafanaBoard(req.Context(), boardData)
	if err != nil {
		msg := "unable to import the boards"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(board)
	if err != nil {
		logrus.Errorf("error marshalling board: %v", err)
		http.Error(w, "unable to marshal the board instance", http.StatusInternalServerError)
		return
	}
}

// PrometheusQueryHandler handles prometheus queries
func (h *Handler) PrometheusQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		http.Error(w, "Prometheus URL is not configured", http.StatusBadRequest)
		return
	}

	reqQuery := req.URL.Query()

	prometheusClient, err := helpers.NewPrometheusClientWithHTTPClient(req.Context(), sessObj.Prometheus.PrometheusURL, &http.Client{
		Timeout: time.Second,
	}, false)
	if err != nil {
		msg := "unable to create a client to talk to prometheus"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer prometheusClient.Close()

	data, err := prometheusClient.Query(req.Context(), &reqQuery)
	if err != nil {
		msg := "connection to prometheus failed"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// PrometheusQueryRangeHandler handles prometheus range queries
func (h *Handler) PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		http.Error(w, "Prometheus URL is not configured", http.StatusBadRequest)
		return
	}

	reqQuery := req.URL.Query()

	prometheusClient, err := helpers.NewPrometheusClientWithHTTPClient(req.Context(), sessObj.Prometheus.PrometheusURL, &http.Client{
		Timeout: time.Second,
	}, false)
	if err != nil {
		msg := "unable to create a client to talk to prometheus"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer prometheusClient.Close()

	testUUID := reqQuery.Get("uuid")
	if testUUID != "" {
		q := reqQuery.Get("query")
		h.config.QueryTracker.AddOrFlagQuery(req.Context(), testUUID, q, false)
	}

	data, err := prometheusClient.QueryRange(req.Context(), &reqQuery)
	if err != nil {
		msg := "connection to prometheus failed"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// PrometheusStaticBoardHandler returns the static board
func (h *Handler) PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		_, _ = w.Write([]byte("{}"))
		return
	}
	prometheusClient, err := helpers.NewPrometheusClientWithHTTPClient(req.Context(), sessObj.Prometheus.PrometheusURL, &http.Client{
		Timeout: time.Second,
	}, true)
	if err != nil {
		msg := "unable to create a client to talk to prometheus"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer prometheusClient.Close()

	result := map[string]*models.GrafanaBoard{}
	resultLock := &sync.Mutex{}
	resultWG := &sync.WaitGroup{}

	boardFunc := map[string]func(context.Context) (*models.GrafanaBoard, error){
		"cluster": prometheusClient.GetClusterStaticBoard,
		"node":    prometheusClient.GetNodesStaticBoard,
	}

	for key, bfunc := range boardFunc {
		resultWG.Add(1)
		go func(k string, bfun func(context.Context) (*models.GrafanaBoard, error)) {
			defer resultWG.Done()

			board, err := bfun(req.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			resultLock.Lock()
			defer resultLock.Unlock()
			result[k] = board
		}(key, bfunc)
	}
	resultWG.Wait()

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logrus.Errorf("error marshalling board: %v", err)
		http.Error(w, "unable to marshal board instance", http.StatusInternalServerError)
		return
	}
}

// SaveSelectedPrometheusBoardsHandler persists selected board and panels
func (h *Handler) SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		http.Error(w, "Prometheus URL is not configured", http.StatusBadRequest)
		return
	}

	// if sessObj.Prometheus.SelectedPrometheusBoardsConfigs == nil {
	// 	sessObj.Prometheus.SelectedPrometheusBoardsConfigs = []*models.GrafanaBoard{}
	// }

	defer func() {
		_ = req.Body.Close()
	}()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := "unable to read the request body"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	boards := []*models.SelectedGrafanaConfig{}
	err = json.Unmarshal(body, &boards)
	if err != nil {
		msg := "unable to parse the request body"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if len(boards) > 0 {
		sessObj.Prometheus.SelectedPrometheusBoardsConfigs = boards
	} else {
		sessObj.Prometheus.SelectedPrometheusBoardsConfigs = nil
	}
	err = h.config.SessionPersister.Write(user.UserID, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
