package handlers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/models"

	"github.com/layer5io/meshery/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&helpers.PrometheusClient{})
}

func (h *Handler) PrometheusConfigHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodDelete {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
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
		logrus.Debugf("Prometheus URL %s succeefully saved", promURL)
	} else if req.Method == http.MethodDelete {
		sessObj.Prometheus = nil
	}

	err = h.config.SessionPersister.Write(user.UserId, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("{}"))
}

func (h *Handler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	boardData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := "unable to get the board payload"
		logrus.Error(errors.Wrap(err, msg))
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	board, err := prometheusClient.ImportGrafanaBoard(req.Context(), boardData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(board)
	if err != nil {
		logrus.Errorf("error marshalling board: %v", err)
		http.Error(w, fmt.Sprintf("unable to marshal board: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) PrometheusQueryHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
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

	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), sessObj.Prometheus.PrometheusURL, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := prometheusClient.Query(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *Handler) PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
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

	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), sessObj.Prometheus.PrometheusURL, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	testUUID := reqQuery.Get("uuid")
	if testUUID != "" {
		q := reqQuery.Get("query")
		h.config.QueryTracker.AddOrFlagQuery(req.Context(), testUUID, q, false)
	}

	data, err := prometheusClient.QueryRange(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *Handler) PrometheusStaticBoardHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Prometheus == nil || sessObj.Prometheus.PrometheusURL == "" {
		w.Write([]byte("{}"))
		return
	}
	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), sessObj.Prometheus.PrometheusURL, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	board, err := prometheusClient.GetStaticBoard(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(board)
	if err != nil {
		logrus.Errorf("error marshalling board: %v", err)
		http.Error(w, fmt.Sprintf("unable to marshal board: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SaveSelectedPrometheusBoardsHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	h.config.SessionPersister.Lock(user.UserId)
	defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
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

	if sessObj.Prometheus.SelectedPrometheusBoardsConfigs == nil {
		sessObj.Prometheus.SelectedPrometheusBoardsConfigs = []*models.GrafanaBoard{}
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		msg := "unable to read the request body"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	boards := []*models.GrafanaBoard{}
	err = json.Unmarshal(body, boards)
	if err != nil {
		msg := "unable to parse the request body"
		err = errors.Wrapf(err, msg)
		logrus.Error(err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if len(boards) > 0 {
		sessObj.Prometheus.SelectedPrometheusBoardsConfigs = boards
	} else {
		sessObj.Prometheus.SelectedPrometheusBoardsConfigs = nil
	}
	err = h.config.SessionPersister.Write(user.UserId, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{}"))
}
