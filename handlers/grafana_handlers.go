package handlers

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"

	"github.com/layer5io/meshery/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&helpers.GrafanaClient{})
}

// GrafanaConfigHandler is used for persisting or removing Grafana configuration
func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
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
		grafanaURL := req.FormValue("grafanaURL")
		grafanaAPIKey := req.FormValue("grafanaAPIKey")

		sessObj.Grafana = &models.Grafana{
			GrafanaURL:    grafanaURL,
			GrafanaAPIKey: grafanaAPIKey,
		}

		g, err := helpers.NewGrafanaClient(grafanaURL, grafanaAPIKey, true)
		if err != nil {
			http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
			return
		}
		defer g.Close()
		logrus.Debugf("connection to grafana @ %s succeeded", grafanaURL)
	} else if req.Method == http.MethodDelete {
		sessObj.Grafana = nil
	}
	err = h.config.SessionPersister.Write(user.UserID, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// GrafanaBoardsHandler is used for fetching Grafana boards and panels
func (h *Handler) GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet && req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if req.Method == http.MethodPost {
		h.SaveSelectedGrafanaBoardsHandler(w, req, session, user)
		return
	}

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Grafana == nil || sessObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	grafanaClient, err := helpers.NewGrafanaClient(sessObj.Grafana.GrafanaURL, sessObj.Grafana.GrafanaAPIKey, true)
	if err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}
	defer grafanaClient.Close()

	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	boards, err := grafanaClient.GetGrafanaBoards(dashboardSearch)
	if err != nil {
		msg := "unable to get grafana boards"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		logrus.Errorf("error marshalling boards: %v", err)
		http.Error(w, "unable to marshal boards payload", http.StatusInternalServerError)
		return
	}
}

// GrafanaQueryHandler is used for handling Grafana queries
func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reqQuery := req.URL.Query()

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Grafana == nil || sessObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	grafanaClient, err := helpers.NewGrafanaClientWithHTTPClient(sessObj.Grafana.GrafanaURL, sessObj.Grafana.GrafanaAPIKey, &http.Client{
		Timeout: time.Second,
	}, true)
	if err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}
	defer grafanaClient.Close()

	data, err := grafanaClient.GrafanaQuery(req.Context(), &reqQuery)
	if err != nil {
		msg := "unable to query grafana"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// GrafanaQueryRangeHandler is used for handling Grafana Range queries
func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reqQuery := req.URL.Query()

	sessObj, err := h.config.SessionPersister.Read(user.UserID)
	if err != nil {
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}

	if sessObj.Grafana == nil || sessObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	grafanaClient, err := helpers.NewGrafanaClientWithHTTPClient(sessObj.Grafana.GrafanaURL, sessObj.Grafana.GrafanaAPIKey, &http.Client{
		Timeout: time.Second,
	}, true)
	if err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}
	defer grafanaClient.Close()

	data, err := grafanaClient.GrafanaQueryRange(req.Context(), &reqQuery)
	if err != nil {
		msg := "unable to query grafana"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// SaveSelectedGrafanaBoardsHandler is used to persist board and panel selection
func (h *Handler) SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, user *models.User) {
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

	if sessObj.Grafana == nil || sessObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	// if sessObj.Grafana.GrafanaBoards == nil {
	// 	sessObj.Grafana.GrafanaBoards = []*models.SelectedGrafanaConfig{}
	// }

	defer func() {
		req.Body.Close()
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
		sessObj.Grafana.GrafanaBoards = boards
	} else {
		sessObj.Grafana.GrafanaBoards = nil
	}
	err = h.config.SessionPersister.Write(user.UserID, sessObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
