package helpers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/gosimple/slug"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"

	"fmt"
	"time"

	"github.com/grafana-tools/sdk"
)

// GrafanaClient represents a client to Grafana in Meshery
type GrafanaClient struct {
	BaseURL string
	APIKey  string
	OrgID   uint
	c       *sdk.Client

	promMode bool
	promURL  string
}

// NewGrafanaClient returns a new GrafanaClient
func NewGrafanaClient(BaseURL, APIKey string, validateConfig bool) (*GrafanaClient, error) {
	if strings.HasSuffix(BaseURL, "/") {
		BaseURL = strings.Trim(BaseURL, "/")
	}
	g := &GrafanaClient{
		BaseURL: BaseURL,
		APIKey:  APIKey,
	}
	if validateConfig {
		var err error
		g.c = sdk.NewClient(g.BaseURL, g.APIKey, &http.Client{
Timeout: 25 * time.Second,})
		if g.OrgID, err = g.GrafanaConfigValidator(); err != nil {
			return nil, err
		}
	}
	return g, nil
}

// NewGrafanaClientForPrometheus returns a limited GrafanaClient for use with Prometheus
func NewGrafanaClientForPrometheus(promURL string) *GrafanaClient {
	if strings.HasSuffix(promURL, "/") {
		promURL = strings.Trim(promURL, "/")
	}
	g := &GrafanaClient{
		promURL:  promURL,
		promMode: true,
	}
	return g
}

func (g *GrafanaClient) makeRequest(ctx context.Context, queryURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, queryURL, nil)
	if !g.promMode {
		req.Header.Set("Authorization", g.APIKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "autograf")
	c := &http.Client{}
	resp, err := c.Do(req)
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

// GrafanaConfigValidator validates connection to Grafana
func (g *GrafanaClient) GrafanaConfigValidator() (uint, error) {
	if g.c == nil {
		g.c = sdk.NewClient(g.BaseURL, g.APIKey, &http.Client{})
	}
	org, err := g.c.GetActualOrg()
	if err != nil {
		err = errors.Wrapf(err, "connection to grafana failed")
		logrus.Error(err)
		return 0, errors.New("connection to grafana failed")
	}
	// logrus.Infof("connection to grafana succeeded @ %s", g.BaseURL)
	return org.ID, nil
}

// GetGrafanaBoards retrieves all the Grafana boards matching a search
func (g *GrafanaClient) GetGrafanaBoards(dashboardSearch string) ([]*models.GrafanaBoard, error) {
	if g.c == nil {
		g.c = sdk.NewClient(g.BaseURL, g.APIKey, &http.Client{})
	}
	boardLinks, err := g.c.SearchDashboards(dashboardSearch, false)
	if err != nil {
		logrus.Error(errors.Wrapf(err, "error getting boards from grafana"))
		return nil, errors.New("unable to fetch boards from grafana")
	}
	boards := []*models.GrafanaBoard{}
	for _, link := range boardLinks {
		if link.Type != "dash-db" {
			continue
		}
		board, _, err := g.c.GetDashboard(link.URI)
		if err != nil {
			err1 := fmt.Errorf("error getting board from grafana for URI - %s", link.URI)
			logrus.Error(errors.Wrapf(err, err1.Error()))
			return nil, err1
		}
		grafBoard, err := g.ProcessBoard(&board, &link)
		boards = append(boards, grafBoard)
	}
	return boards, nil
}

// ProcessBoard accepts raw Grafana board and returns a processed GrafanaBoard to be used in Meshery
func (g *GrafanaClient) ProcessBoard(board *sdk.Board, link *sdk.FoundBoard) (*models.GrafanaBoard, error) {
	grafBoard := &models.GrafanaBoard{
		URI:          link.URI,
		Title:        link.Title,
		UID:          board.UID,
		Slug:         slug.Make(board.Title),
		TemplateVars: []*models.GrafanaTemplateVars{},
		// Panels:       []*models.GrafanaPanel{},
		Panels: []*sdk.Panel{},
		OrgID:  g.OrgID,
	}
	var err error

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
				err := fmt.Errorf("unable to get datasource name for tmpvar: %+#v", tmpVar)
				logrus.Error(err)
				return nil, err
			}
			if g.c != nil {
				ds, err = g.c.GetDatasourceByName(dsName)
				if err != nil {
					msg := fmt.Errorf("error getting board datasource with name - %s", dsName)
					logrus.Error(errors.Wrapf(err, msg.Error()))
					return nil, msg
				}
			} else {
				ds.Name = dsName
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
				// logrus.Debugf("board: %d, panel id: %d", board.ID, p1.ID)
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
	return grafBoard, nil
}

// GrafanaQuery parses the provided query data and queries Grafana and streams response
func (g *GrafanaClient) GrafanaQuery(ctx context.Context, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		err := errors.New("query data is empty")
		logrus.Error(err)
		return nil, err
	}
	query := strings.TrimSpace(queryData.Get("query"))
	dsID := queryData.Get("dsid")
	var queryURL string
	switch {
	case strings.HasPrefix(query, "label_values("):
		val := strings.Replace(query, "label_values(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		if strings.Contains(val, ",") {
			start := queryData.Get("start")
			end := queryData.Get("end")
			comInd := strings.LastIndex(val, ", ")
			if comInd > -1 {
				val = val[:comInd]
			}
			for key := range *queryData {
				if key != "query" && key != "dsid" && key != "start" && key != "end" {
					val1 := queryData.Get(key)
					val = strings.Replace(val, "$"+key, val1, -1)
				}
			}
			var reqURL string
			if g.promMode {
				reqURL = fmt.Sprintf("%s/api/v1/series", g.promURL)
			} else {
				reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/series", g.BaseURL, dsID)
			}
			queryURLInst, _ := url.Parse(reqURL)
			qParams := queryURLInst.Query()
			qParams.Set("match[]", val)
			if start != "" && end != "" {
				qParams.Set("start", start)
				qParams.Set("end", end)
			}
			queryURLInst.RawQuery = qParams.Encode()
			queryURL = queryURLInst.String()
		} else {
			if g.promMode {
				queryURL = fmt.Sprintf("%s/api/v1/label/%s/values", g.promURL, val)
			} else {
				queryURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/label/%s/values", g.BaseURL, dsID, val)
			}
		}
	case strings.HasPrefix(query, "query_result("):
		val := strings.Replace(query, "query_result(", "", 1)
		val = strings.TrimSpace(strings.TrimSuffix(val, ")"))
		for key := range *queryData {
			if key != "query" && key != "dsid" {
				val1 := queryData.Get(key)
				val = strings.Replace(val, "$"+key, val1, -1)
			}
		}
		var reqURL string
		if g.promMode {
			// reqURL = fmt.Sprintf("%s/api/v1/query", g.promURL, dsID)
			reqURL = fmt.Sprintf("%s/api/v1/query", g.promURL)
		} else {
			reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query", g.BaseURL, dsID)
		}
		newURL, _ := url.Parse(reqURL)
		q := url.Values{}
		q.Set("query", val)
		newURL.RawQuery = q.Encode()
		queryURL = newURL.String()
	default:
		// {"status":"success","data":["istio-pilot.istio-system.svc.cluster.local","istio-telemetry.istio-system.svc.cluster.local"]}
		return json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   []string{query},
		})
	}
	logrus.Debugf("derived query url: %s", queryURL)

	data, err := g.makeRequest(ctx, queryURL)
	if err != nil {
		msg := errors.New("error getting data from grafana")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
	}
	return data, nil
}

// GrafanaQueryRange parses the given params and performs Grafana range queries
func (g *GrafanaClient) GrafanaQueryRange(ctx context.Context, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		err := errors.New("query data is empty")
		logrus.Error(err)
		return nil, err
	}
	ds := queryData.Get("ds")
	var reqURL string
	if g.promMode {
		reqURL = fmt.Sprintf("%s/api/v1/query_range", g.promURL)
	} else {
		reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query_range", g.BaseURL, ds)
	}

	newURL, _ := url.Parse(reqURL)
	q := url.Values{}
	q.Set("query", queryData.Get("query"))
	q.Set("start", queryData.Get("start"))
	q.Set("end", queryData.Get("end"))
	q.Set("step", queryData.Get("step"))
	newURL.RawQuery = q.Encode()
	queryURL := newURL.String()
	// logrus.Debugf("Query range url: %s", queryURL)
	data, err := g.makeRequest(ctx, queryURL)
	if err != nil {
		msg := errors.New("error getting data from grafana")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
	}
	return data, nil
}
