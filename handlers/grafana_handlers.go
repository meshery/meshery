package handlers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/helpers"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&helpers.GrafanaClient{})
}

func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request) {
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

	grafanaURL := req.FormValue("grafanaURL")
	grafanaAPIKey := req.FormValue("grafanaAPIKey")

	grafanaClient, err := helpers.NewGrafanaClient(grafanaURL, grafanaAPIKey, true)
	if err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}
	session.Values["grafanaURL"] = grafanaURL
	session.Values["grafanaAPIKey"] = grafanaAPIKey
	session.Values["grafanaClient"] = grafanaClient
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}
	logrus.Debugf("connection to grafana @ %s succeeded", grafanaURL)
	// TODO: save the grafana configs with user data
	w.Write([]byte("{}"))
}

func (h *Handler) GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request) {
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
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)
	// grafanaOrg, _ := session.Values["grafanaOrgID"].(uint)
	grafanaClient, _ := session.Values["grafanaClient"].(*helpers.GrafanaClient)

	if grafanaClient == nil {
		grafanaClient, err = helpers.NewGrafanaClient(grafanaURL, grafanaAPIKey, true)
		if err != nil {
			http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
			return
		}
	}
	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	boards, err := grafanaClient.GetGrafanaBoards(dashboardSearch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		logrus.Errorf("error marshalling boards: %v", err)
		http.Error(w, fmt.Sprintf("unable to marshal boards: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request) {
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
	reqQuery := req.URL.Query()
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)

	grafanaClient, _ := session.Values["grafanaClient"].(*helpers.GrafanaClient)

	if grafanaClient == nil {
		grafanaClient, err = helpers.NewGrafanaClient(grafanaURL, grafanaAPIKey, true)
		if err != nil {
			http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
			return
		}
	}
	data, err := grafanaClient.GrafanaQuery(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
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
	reqQuery := req.URL.Query()
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)
	grafanaClient, _ := session.Values["grafanaClient"].(*helpers.GrafanaClient)

	if grafanaClient == nil {
		grafanaClient, err = helpers.NewGrafanaClient(grafanaURL, grafanaAPIKey, true)
		if err != nil {
			http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
			return
		}
	}
	data, err := grafanaClient.GrafanaQueryRange(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
