package handlers

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/helpers"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&helpers.PrometheusClient{})
}

func (h *Handler) PrometheusConfigHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	promURL := req.FormValue("promURL")
	promClient := helpers.NewPrometheusClient(promURL)
	session.Values["promURL"] = promURL
	session.Values["prometheusClient"] = promClient
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}
	logrus.Debugf("Prometheus URL %s succeefully saved", promURL)

	w.Write([]byte("{}"))
}

func (h *Handler) GrafanaBoardImportForPrometheusHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	promURL, _ := session.Values["promURL"].(string)
	prometheusClient, _ := session.Values["prometheusClient"].(*helpers.PrometheusClient)

	if prometheusClient == nil {
		prometheusClient = helpers.NewPrometheusClient(promURL)
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
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	reqQuery := req.URL.Query()
	promURL, _ := session.Values["promURL"].(string)
	prometheusClient, _ := session.Values["prometheusClient"].(*helpers.PrometheusClient)

	if prometheusClient == nil {
		prometheusClient = helpers.NewPrometheusClient(promURL)
	}
	data, err := prometheusClient.Query(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (h *Handler) PrometheusQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	reqQuery := req.URL.Query()
	promURL, _ := session.Values["promURL"].(string)
	prometheusClient, _ := session.Values["prometheusClient"].(*helpers.PrometheusClient)

	if prometheusClient == nil {
		prometheusClient = helpers.NewPrometheusClient(promURL)
	}
	data, err := prometheusClient.Query(req.Context(), &reqQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
