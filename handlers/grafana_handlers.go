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

	"github.com/sirupsen/logrus"
)

func init() {
	gob.Register(&models.GrafanaClient{})
}

// GrafanaConfigHandler is used for persisting or removing Grafana configuration
func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if req.Method == http.MethodGet {
		err := json.NewEncoder(w).Encode(prefObj.Grafana)
		if err != nil {
			h.log.Error(ErrMarshal(err))
			http.Error(w, ErrMarshal(err).Error(), http.StatusInternalServerError)
			return
		}

		return
	}

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
			h.log.Error(ErrGrafanaScan(err))
			http.Error(w, ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
			return
		}
		logrus.Debugf("connection to grafana @ %s succeeded", grafanaURL)
	} else if req.Method == http.MethodDelete {
		prefObj.Grafana = nil
	}
	err := p.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// GrafanaPingHandler - used to initiate a Grafana ping
func (h *Handler) GrafanaPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		h.log.Error(ErrGrafanaConfig)
		http.Error(w, ErrGrafanaConfig.Error(), http.StatusBadRequest)
		return
	}

	if prefObj.K8SConfig == nil || !prefObj.K8SConfig.InClusterConfig && (prefObj.K8SConfig.Config == nil || len(prefObj.K8SConfig.Config) == 0) {
		h.log.Error(ErrInvalidK8SConfig)
		http.Error(w, ErrInvalidK8SConfig.Error(), http.StatusBadRequest)
		return
	}

	if err := h.config.GrafanaClient.Validate(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey); err != nil {
		h.log.Error(ErrGrafanaScan(err))
		http.Error(w, ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte("{}"))
}

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
		h.log.Error(ErrGrafanaConfig)
		http.Error(w, ErrGrafanaConfig.Error(), http.StatusBadRequest)
		return
	}

	if err := h.config.GrafanaClient.Validate(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey); err != nil {
		h.log.Error(ErrGrafanaScan(err))
		http.Error(w, ErrGrafanaScan(err).Error(), http.StatusInternalServerError)
		return
	}

	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	boards, err := h.config.GrafanaClient.GetGrafanaBoards(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, dashboardSearch)
	if err != nil {
		h.log.Error(ErrGrafanaBoards(err))
		http.Error(w, ErrGrafanaBoards(err).Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		h.log.Error(ErrMarshal(err))
		http.Error(w, ErrMarshal(err).Error(), http.StatusInternalServerError)
		return
	}
}

// GrafanaQueryHandler is used for handling Grafana queries
func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		err := ErrGrafanaConfig
		h.log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := h.config.GrafanaClientForQuery.GrafanaQuery(req.Context(), prefObj.Grafana.GrafanaURL, prefObj.Grafana.GrafanaAPIKey, &reqQuery)
	if err != nil {
		h.log.Error(ErrGrafanaQuery(err))
		http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// GrafanaQueryRangeHandler is used for handling Grafana Range queries
func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	reqQuery := req.URL.Query()

	data, err := h.config.GrafanaClientForQuery.GrafanaQueryRange(req.Context(), reqQuery.Get("url"), reqQuery.Get("api-key"), &reqQuery)
	if err != nil {
		h.log.Error(ErrGrafanaQuery(err))
		http.Error(w, ErrGrafanaQuery(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

// SaveSelectedGrafanaBoardsHandler is used to persist board and panel selection
func (h *Handler) SaveSelectedGrafanaBoardsHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if prefObj.Grafana == nil || prefObj.Grafana.GrafanaURL == "" {
		h.log.Error(ErrGrafanaConfig)
		http.Error(w, ErrGrafanaConfig.Error(), http.StatusBadRequest)
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
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}
	boards := []*models.SelectedGrafanaConfig{}
	err = json.Unmarshal(body, &boards)
	if err != nil {
		h.log.Error(ErrUnmarshal(err))
		http.Error(w, ErrUnmarshal(err).Error(), http.StatusBadRequest)
		return
	}
	if len(boards) > 0 {
		prefObj.Grafana.GrafanaBoards = boards
	} else {
		prefObj.Grafana.GrafanaBoards = nil
	}
	err = p.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		h.log.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}
	h.log.Info("Board selection updated")
	_, _ = w.Write([]byte("{}"))
}
