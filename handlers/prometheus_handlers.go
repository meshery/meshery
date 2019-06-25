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

	promURL := req.FormValue("prometheusURL")
	if _, err = helpers.NewPrometheusClient(req.Context(), promURL, true); err != nil {
		logrus.Errorf("unable to connect to prometheus: %v", err)
		http.Error(w, "unable to connect to prometheus", http.StatusInternalServerError)
		return
	}
	session.Values["promURL"] = promURL
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
	promURL, _ := session.Values["promURL"].(string)

	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), promURL, false)
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
	reqQuery := req.URL.Query()
	promURL, _ := session.Values["promURL"].(string)

	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), promURL, false)
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
	reqQuery := req.URL.Query()
	promURL, _ := session.Values["promURL"].(string)
	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), promURL, false)
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
	promURL, _ := session.Values["promURL"].(string)
	if promURL == "" {
		w.Write([]byte("{}"))
		return
	}
	prometheusClient, err := helpers.NewPrometheusClient(req.Context(), promURL, true)
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
