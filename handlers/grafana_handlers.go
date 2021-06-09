//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/layer5io/meshery/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&models.GrafanaClient{})
}

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

// GrafanaConfigHandler is used for persisting or removing Grafana configuration
func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if req.Method == http.MethodPost {
		grafanaURL := req.FormValue("grafanaURL")
		grafanaAPIKey := req.FormValue("grafanaAPIKey")

		u, err := url.Parse(grafanaURL)
		if err != nil {
			return
		}
		if strings.Contains(grafanaURL, u.RequestURI()) {
			grafanaURL = strings.TrimSuffix(grafanaURL, u.RequestURI())
		}

		prefObj.Grafana = &models.Grafana{
			GrafanaURL:    grafanaURL,
			GrafanaAPIKey: grafanaAPIKey,
		}

		if err := h.config.GrafanaClient.Validate(req.Context(), grafanaURL, grafanaAPIKey); err != nil {
			http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
			return
		}
		logrus.Debugf("connection to grafana @ %s succeeded", grafanaURL)
	} else if req.Method == http.MethodDelete {
		prefObj.Grafana = nil
	}
	err := p.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// swagger:route GET /api/telemetry/metrics/grafana/ping GrafanaAPI idGetGrafanaPing
// Handle GET request for Grafana ping
//
// Used to initiate a Grafana ping
// responses:
// 	200:

// GrafanaPingHandler - used to initiate a Grafana ping
func (h *Handler) GrafanaPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		logrus.Error("No valid kubernetes config found.")
		http.Error(w, `No valid kubernetes config found.`, http.StatusBadRequest)
		return
	}

	if err := h.config.GrafanaClient.Validate(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey); err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

// swagger:route GET /api/telemetry/metrics/grafana/boards GrafanaAPI idGetGrafanaBoards
// Handle GET request for Grafana boards
//
// Used for fetching Grafana boards and panels
// responses:
// 	200: grafanaBoardsResponseWrapper

// GrafanaBoardsHandler is used for fetching Grafana boards and panels
func (h *Handler) GrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet && req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	if req.Method == http.MethodPost {
		h.SaveSelectedGrafanaBoardsHandler(w, req, prefObj, user, p)
		return
	}

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	if err := h.config.GrafanaClient.Validate(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey); err != nil {
		http.Error(w, "connection to grafana failed", http.StatusInternalServerError)
		return
	}

	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	boards, err := h.config.GrafanaClient.GetGrafanaBoards(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, dashboardSearch)
	if err != nil {
		msg := "unable to get grafana boards"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		logrus.Errorf("error marshaling boards: %v", err)
		http.Error(w, "unable to marshal boards payload", http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/telemetry/metrics/grafana/query GrafanaAPI idGetGrafanaQuery
// Handle GET request for Grafana queries
//
// Used for handling Grafana queries
// responses:
// 	200:

// GrafanaQueryHandler is used for handling Grafana queries
func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	data, err := h.config.GrafanaClientForQuery.GrafanaQuery(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, &reqQuery)
	if err != nil {
		msg := "unable to query grafana"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// GrafanaQueryRangeHandler is used for handling Grafana Range queries
func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	data, err := h.config.GrafanaClientForQuery.GrafanaQueryRange(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, &reqQuery)
	if err != nil {
		msg := "unable to query grafana"
		logrus.Error(errors.Wrapf(err, msg))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// swagger:route POST /api/telemetry/metrics/grafana/boards GrafanaAPI idPostGrafanaBoards
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

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		http.Error(w, "Grafana URL is not configured", http.StatusBadRequest)
		return
	}

	// if prefObj.Grafana.GrafanaBoards == nil {
	// 	prefObj.Grafana.GrafanaBoards = []*models.SelectedGrafanaConfig{}
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
		prefObj.Grafana.GrafanaBoards = boards
	} else {
		prefObj.Grafana.GrafanaBoards = nil
	}
	err = p.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		logrus.Errorf("unable to save user config data: %v", err)
		http.Error(w, "unable to save user config data", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}
