package models

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	"github.com/gosimple/slug"
	"github.com/sirupsen/logrus"

	"fmt"
	"time"

	"github.com/grafana-tools/sdk"
)

// GrafanaClient represents a client to Grafana in Meshery
type GrafanaClient struct {
	httpClient *http.Client

	promMode bool
}

// NewGrafanaClient returns a new GrafanaClient
func NewGrafanaClient() *GrafanaClient {
	return NewGrafanaClientWithHTTPClient(&http.Client{
		Timeout: 25 * time.Second,
	})
}

// NewGrafanaClientWithHTTPClient returns a new GrafanaClient with the given HTTP Client
func NewGrafanaClientWithHTTPClient(client *http.Client) *GrafanaClient {
	return &GrafanaClient{
		httpClient: client,
	}
}

// NewGrafanaClientForPrometheusWithHTTPClient returns a limited GrafanaClient for use with Prometheus with the given HTTP client
func NewGrafanaClientForPrometheusWithHTTPClient(client *http.Client) *GrafanaClient {
	// if strings.HasSuffix(promURL, "/") {
	// 	promURL = strings.Trim(promURL, "/")
	// }
	g := &GrafanaClient{
		promMode:   true,
		httpClient: client,
	}
	return g
}

// Validate - helps validate grafana connection
func (g *GrafanaClient) Validate(ctx context.Context, BaseURL, APIKey string) error {
	if strings.HasSuffix(BaseURL, "/") {
		BaseURL = strings.Trim(BaseURL, "/")
	}
	c := sdk.NewClient(BaseURL, APIKey, g.httpClient)
	if _, err := c.GetActualOrg(ctx); err != nil {
		err = errors.Wrapf(err, "connection to grafana failed")
		logrus.Error(err)
		return err
	}
	return nil
}
func (g *GrafanaClient) makeRequest(ctx context.Context, queryURL, APIKey string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, queryURL, nil)
	if err != nil {
		return nil, err
	}
	if !g.promMode {
		req.Header.Set("Authorization", APIKey)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "autograf")
	// c := &http.Client{}
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
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

// GetGrafanaBoards retrieves all the Grafana boards matching a search
func (g *GrafanaClient) GetGrafanaBoards(ctx context.Context, BaseURL, APIKey, dashboardSearch string) ([]*GrafanaBoard, error) {
	if strings.HasSuffix(BaseURL, "/") {
		BaseURL = strings.Trim(BaseURL, "/")
	}
	c := sdk.NewClient(BaseURL, APIKey, g.httpClient)

	boardLinks, err := c.SearchDashboards(ctx, dashboardSearch, false)
	if err != nil {
		logrus.Error(errors.Wrapf(err, "error getting boards from grafana"))
		return nil, errors.New("unable to fetch boards from grafana")
	}
	boards := []*GrafanaBoard{}
	for _, link := range boardLinks {
		if link.Type != "dash-db" {
			continue
		}
		// TODO Need to do the unitest for Grafana helper
		board, _, err := c.GetDashboardByUID(ctx, link.UID)
		if err != nil {
			err1 := fmt.Errorf("error getting board from grafana for URI - %s", link.URI)
			logrus.Error(errors.Wrapf(err, err1.Error()))
			return nil, err1
		}
		grafBoard, err := g.ProcessBoard(ctx, c, &board, &link)
		if err != nil {
			return nil, err
		}
		boards = append(boards, grafBoard)
	}
	return boards, nil
}

// ProcessBoard accepts raw Grafana board and returns a processed GrafanaBoard to be used in Meshery
func (g *GrafanaClient) ProcessBoard(ctx context.Context, c *sdk.Client, board *sdk.Board, link *sdk.FoundBoard) (*GrafanaBoard, error) {
	var orgID uint
	if !g.promMode {
		org, err := c.GetActualOrg(ctx)
		if err != nil {
			err = errors.Wrapf(err, "connection to grafana failed")
			logrus.Error(err)
			return nil, err
		}
		orgID = org.ID
	}
	grafBoard := &GrafanaBoard{
		URI:          link.URI,
		Title:        link.Title,
		UID:          board.UID,
		Slug:         slug.Make(board.Title),
		TemplateVars: []*GrafanaTemplateVars{},
		Panels:       []*sdk.Panel{},
		OrgID:        orgID,
	}
	var err error
	tmpDsName := map[string]string{}
	if len(board.Templating.List) > 0 {
		for _, tmpVar := range board.Templating.List {
			var ds sdk.Datasource
			var dsName string
			if tmpVar.Type == "datasource" {
				dsName = strings.Title(strings.ToLower(tmpVar.Query)) // datasource name can be found in the query field
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
			if c != nil {
				ds, err = c.GetDatasourceByName(ctx, dsName)
				if err != nil {
					msg := fmt.Errorf("error getting board datasource with name - %s", dsName)
					logrus.Error(errors.Wrapf(err, msg.Error()))
					return nil, msg
				}
			} else {
				ds.Name = dsName
			}

			tvVal := tmpVar.Current.Text
			grafBoard.TemplateVars = append(grafBoard.TemplateVars, &GrafanaTemplateVars{
				Name:  tmpVar.Name,
				Query: tmpVar.Query,
				Datasource: &GrafanaDataSource{
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
					}
				}
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
					_ = p4.UnmarshalJSON(p3)
					logrus.Debugf("board: %d, Row panel id: %d", board.ID, p4.ID)
					grafBoard.Panels = append(grafBoard.Panels, p4)
				}
			}
		}
	}
	return grafBoard, nil
}

// GrafanaQuery parses the provided query data and queries Grafana and streams response
func (g *GrafanaClient) GrafanaQuery(ctx context.Context, BaseURL, APIKey string, queryData *url.Values) ([]byte, error) {
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
				reqURL = fmt.Sprintf("%s/api/v1/series", BaseURL)
			} else {
				reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/series", BaseURL, dsID)
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
				queryURL = fmt.Sprintf("%s/api/v1/label/%s/values", BaseURL, val)
			} else {
				queryURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/label/%s/values", BaseURL, dsID, val)
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
			reqURL = fmt.Sprintf("%s/api/v1/query", BaseURL)
		} else {
			reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query", BaseURL, dsID)
		}
		newURL, _ := url.Parse(reqURL)
		q := url.Values{}
		q.Set("query", val)
		newURL.RawQuery = q.Encode()
		queryURL = newURL.String()
	default:
		return json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   []string{query},
		})
	}
	logrus.Debugf("derived query url: %s", queryURL)

	data, err := g.makeRequest(ctx, queryURL, APIKey)
	if err != nil {
		msg := errors.New("error getting data from grafana")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
	}
	return data, nil
}

// GrafanaQueryRange parses the given params and performs Grafana range queries
func (g *GrafanaClient) GrafanaQueryRange(ctx context.Context, BaseURL, APIKey string, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		err := errors.New("query data is empty")
		logrus.Error(err)
		return nil, err
	}
	ds := queryData.Get("ds")
	var reqURL string
	if g.promMode {
		reqURL = fmt.Sprintf("%s/api/v1/query_range", BaseURL)
	} else {
		reqURL = fmt.Sprintf("%s/api/datasources/proxy/%s/api/v1/query_range", BaseURL, ds)
	}

	newURL, _ := url.Parse(reqURL)
	q := url.Values{}
	q.Set("query", queryData.Get("query"))
	q.Set("start", queryData.Get("start"))
	q.Set("end", queryData.Get("end"))
	q.Set("step", queryData.Get("step"))
	newURL.RawQuery = q.Encode()
	queryURL := newURL.String()
	data, err := g.makeRequest(ctx, queryURL, APIKey)
	if err != nil {
		msg := errors.New("error getting data from grafana")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
	}
	return data, nil
}

// Close - closes idle connections
func (g *GrafanaClient) Close() {
	g.httpClient = nil
}
