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
		// logrus.Debugf("board: %+#v", board)
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
		tmpDsName := map[string]string{}
		if len(board.Templating.List) > 0 {
			for _, tmpVar := range board.Templating.List {
				// logrus.Debugf("tmpvar: %+#v", tmpVar)
				var ds sdk.Datasource
				var dsName string
				if tmpVar.Type == "datasource" {
					dsName = tmpVar.Query // datasource name can be found in the query field
					tmpDsName[tmpVar.Name] = dsName
				} else if tmpVar.Type == "query" && tmpVar.Datasource != nil {
					if !strings.HasPrefix(*tmpVar.Datasource, "$") {
						dsName = *tmpVar.Datasource
					} else {
						dsName = tmpDsName[strings.Replace(*tmpVar.Datasource, "$", "", 1)]
					}
				} else {
					logrus.Errorf("unable to get datasource name for tmpvar: %+#v", tmpVar)
					logrus.Errorf("error getting board datasource with name - %s: %v", *tmpVar.Datasource, err)
					http.Error(w, fmt.Sprintf("unable to get datasource with name - %s: %v", *tmpVar.Datasource, err), http.StatusInternalServerError)
					return
				}
				ds, err = c.GetDatasourceByName(dsName)
				if err != nil {
					logrus.Errorf("error getting board datasource with name - %s: %v", dsName, err)
					http.Error(w, fmt.Sprintf("unable to get datasource with name - %s: %v", dsName, err), http.StatusInternalServerError)
					return
				}

				tvVal := tmpVar.Current.Text
				// if tmpVar.Current. {
				// 	tvVal = tmpVar.Current.Text
				// }
				grafBoard.TemplateVars = append(grafBoard.TemplateVars, &models.GrafanaTemplateVars{
					Name:  tmpVar.Name,
					Query: tmpVar.Query,
					Datasource: &models.GrafanaDataSource{
						ID:   ds.ID,
						Name: ds.Name,
					},
					Hide:  tmpVar.Hide,
					Value: tvVal,
				})
			}
		}
		if len(board.Panels) > 0 {
			for _, p1 := range board.Panels {
				if p1.OfType != sdk.TextType && p1.OfType != sdk.TableType && p1.Type != "row" { // turning off text and table panels for now
					if p1.Datasource != nil {
						if strings.HasPrefix(*p1.Datasource, "$") {
							*p1.Datasource = tmpDsName[strings.Replace(*p1.Datasource, "$", "", 1)]
							// logrus.Debugf("updated panel datasource: %s", *panel.Datasource)
						}
					}
					// grafPanel := &models.GrafanaPanel{
					// 	ID:    panel.ID,
					// 	PType: panel.Type,
					// 	Title: panel.Title,
					// }
					// grafBoard.Panels = append(grafBoard.Panels, grafPanel)
					logrus.Debugf("board: %d, panel id: %d", board.ID, p1.ID)
					grafBoard.Panels = append(grafBoard.Panels, p1)
				}
			}
		} else if len(board.Rows) > 0 {
			for _, r1 := range board.Rows {
				for _, p2 := range r1.Panels {
					if p2.OfType != sdk.TextType && p2.OfType != sdk.TableType && p2.Type != "row" { // turning off text and table panels for now
						if strings.HasPrefix(*p2.Datasource, "$") {
							*p2.Datasource = tmpDsName[strings.Replace(*p2.Datasource, "$", "", 1)]
							// logrus.Debugf("updated panel datasource: %s", *panel.Datasource)
						}
						p3, _ := p2.MarshalJSON()
						p4 := &sdk.Panel{}
						p4.UnmarshalJSON(p3)
						logrus.Debugf("board: %d, Row panel id: %d", board.ID, p4.ID)
						grafBoard.Panels = append(grafBoard.Panels, p4)
					}
				}
			}
		}
		for _, p := range grafBoard.Panels {
			logrus.Debugf("board: %s, panel id: %d", grafBoard.URI, p.ID)
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
	reqQuery := req.URL.Query()
	grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)

	query := strings.TrimSpace(reqQuery.Get("query"))
	dsID := reqQuery.Get("dsid")
	// c := sdk.NewClient(grafanaURL, grafanaAPIKey, &http.Client{})
	if strings.HasSuffix(grafanaURL, "/") {
		grafanaURL = strings.Trim(grafanaURL, "/")
	}

	var queryURL string
	switch {
	case strings.HasPrefix(query, "label_values("):
		val := strings.Replace(query, "label_values(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		if strings.Contains(val, ",") {
			start := reqQuery.Get("start")
			end := reqQuery.Get("end")
			comInd := strings.LastIndex(val, ", ")
			if comInd > -1 {
				val = val[:comInd]
			}
			for key := range reqQuery {
				if key != "query" && key != "dsid" && key != "start" && key != "end" {
					kVal := reqQuery.Get(key)
					val = strings.Replace(val, "$"+key, kVal, -1)
				}
			}
			queryURLInst, _ := url.Parse(fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/series", grafanaURL, dsID))
			qParams := queryURLInst.Query()
			qParams.Set("match[]", val)
			if start != "" && end != "" {
				qParams.Set("start", start)
				qParams.Set("end", end)
			}
			queryURLInst.RawQuery = qParams.Encode()
			queryURL = queryURLInst.String()
		} else {
			queryURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/label/%s/values", grafanaURL, dsID, val)
		}
	case strings.HasPrefix(query, "query_result("):
		val := strings.Replace(query, "query_result(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		for key := range reqQuery {
			if key != "query" && key != "dsid" {
				kVal := reqQuery.Get(key)
				val = strings.Replace(val, "$"+key, kVal, -1)
			}
		}
		newURL, _ := url.Parse(fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query", grafanaURL, dsID))
		q := url.Values{}
		q.Set("query", val)
		newURL.RawQuery = q.Encode()
		queryURL = newURL.String()
	default:
		// {"status":"success","data":["istio-pilot.istio-system.svc.cluster.local","istio-telemetry.istio-system.svc.cluster.local"]}
		data, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   []string{query},
		})
		w.Write(data)
		return
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

func (h *Handler) GrafanaQueryRangeHandler(w http.ResponseWriter, req *http.Request) {
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}
	reqQuery := req.URL.Query()
	// grafanaURL, _ := session.Values["grafanaURL"].(string)
	grafanaAPIKey, _ := session.Values["grafanaAPIKey"].(string)

	queryURL := reqQuery.Get("query")
	gReq, err := http.NewRequest(http.MethodGet, queryURL, nil)
	if err != nil {
		logrus.Errorf("error creating a request: %v", err)
		http.Error(w, fmt.Sprintf("error creating a request: %v", err), http.StatusInternalServerError)
		return
	}
	if grafanaAPIKey != "" {
		gReq.Header.Set("Authorization", "Bearer "+grafanaAPIKey)
	}
	client := &http.Client{}
	resp, err := client.Do(gReq)
	if err != nil {
		logrus.Errorf("error calling requested url: %s - %v", queryURL, err)
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("error parsing response body: %v", err)
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
