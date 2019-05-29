package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gosimple/slug"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"

	"github.com/grafana-tools/sdk"
)

func (h *Handler) grafanaRequest(ctx context.Context, queryURL, key string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, queryURL, nil)
	req.Header.Set("Authorization", key)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "autograf")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		// return nil, fmt.Errorf("%s", data)
		logrus.Errorf("unable to get data from URL: %s due to status code: %d", queryURL, resp.StatusCode)
		return nil, fmt.Errorf("unable to fetch data from url: %s", queryURL)
	}
	return data, nil
}

func (h *Handler) GrafanaConfigHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	grafanaURL := req.FormValue("grafanaURL")
	grafanaAPIKey := req.FormValue("grafanaAPIKey")
	c := sdk.NewClient(grafanaURL, grafanaAPIKey, &http.Client{})
	org, err := c.GetActualOrg()
	if err != nil {
		http.Error(w, "connection to grafana failed", http.StatusForbidden)
		return
	}
	session.Values["grafanaURL"] = grafanaURL
	session.Values["grafanaAPIKey"] = grafanaAPIKey
	session.Values["grafanaOrgID"] = org.ID
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
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)
	grafanaOrg, _ := session.Values["grafanaOrgID"].(uint)

	dashboardSearch := req.URL.Query().Get("dashboardSearch")
	c := sdk.NewClient(grafanaURL, grafanaAPIKey, &http.Client{})
	boardLinks, err := c.SearchDashboards(dashboardSearch, false)
	if err != nil {
		logrus.Errorf("error getting boards from grafana: %v", err)
		http.Error(w, fmt.Sprintf("unable to get boards from grafana: %v", err), http.StatusInternalServerError)
		return
	}
	boards := []*models.GrafanaBoard{}
	for _, link := range boardLinks {
		if link.Type != "dash-db" {
			continue
		}
		board, _, err := c.GetDashboard(link.URI)
		if err != nil {
			logrus.Errorf("error getting board from grafana for URI - %s: %v", link.URI, err)
			http.Error(w, fmt.Sprintf("unable to get board from grafana for URI - %s: %v", link.URI, err), http.StatusInternalServerError)
			return
		}
		grafBoard := &models.GrafanaBoard{
			URI:          link.URI,
			Title:        link.Title,
			UID:          board.UID,
			Slug:         slug.Make(board.Title),
			TemplateVars: []*models.GrafanaTemplateVars{},
			// Panels:       []*models.GrafanaPanel{},
			Panels: []*sdk.Panel{},
			OrgID:  grafanaOrg,
		}
		if len(board.Templating.List) > 0 {
			for _, tmpVar := range board.Templating.List {
				ds, err := c.GetDatasourceByName(*tmpVar.Datasource)
				if err != nil {
					logrus.Errorf("error getting board datasource with name - %s: %v", *tmpVar.Datasource, err)
					http.Error(w, fmt.Sprintf("unable to get datasource with name - %s: %v", *tmpVar.Datasource, err), http.StatusInternalServerError)
					return
				}
				grafBoard.TemplateVars = append(grafBoard.TemplateVars, &models.GrafanaTemplateVars{
					Name:  tmpVar.Name,
					Query: tmpVar.Query,
					Datasource: &models.GrafanaDataSource{
						ID:   ds.ID,
						Name: ds.Name,
					},
				})
			}
		}
		for _, panel := range board.Panels {
			if panel.Type != "text" && panel.Type != "table" { // turning off text and table panels for now
				// grafPanel := &models.GrafanaPanel{
				// 	ID:    panel.ID,
				// 	PType: panel.Type,
				// 	Title: panel.Title,
				// }
				// grafBoard.Panels = append(grafBoard.Panels, grafPanel)
				grafBoard.Panels = append(grafBoard.Panels, panel)
			}
		}
		boards = append(boards, grafBoard)
	}
	err = json.NewEncoder(w).Encode(boards)
	if err != nil {
		logrus.Errorf("error marshalling boards: %v", err)
		http.Error(w, fmt.Sprintf("unable to marshal boards: %v", err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GrafanaQueryHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)

	query := strings.TrimSpace(req.URL.Query().Get("query"))
	dsID := req.URL.Query().Get("dsid")
	// c := sdk.NewClient(grafanaURL, grafanaAPIKey, &http.Client{})
	if strings.HasSuffix(grafanaURL, "/") {
		grafanaURL = strings.Trim(grafanaURL, "/")
	}
	var queryURL string
	switch {
	case strings.HasPrefix(query, "label_values("):
		val := strings.Replace(query, "label_values(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		queryURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/label/%s/values", grafanaURL, dsID, val)
	case strings.HasPrefix(query, "query_result("):
		val := strings.Replace(query, "query_result(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		for key := range req.URL.Query() {
			if key != "query" && key != "dsid" {
				kVal := req.URL.Query().Get(key)
				val = strings.Replace(val, "$"+key, kVal, -1)
			}
		}
		newURL, _ := url.Parse(fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query", grafanaURL, dsID))
		q := url.Values{}
		q.Set("query", val)
		newURL.RawQuery = q.Encode()
		queryURL = newURL.String()
	}
	logrus.Infof("derived query url: %s", queryURL)

	data, err := h.grafanaRequest(req.Context(), queryURL, grafanaAPIKey)
	if err != nil {
		logrus.Errorf("error getting data from grafana: %v", err)
		http.Error(w, fmt.Sprintf("error getting data from grafana: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
