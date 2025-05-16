package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/grafana-tools/sdk"
	"github.com/layer5io/meshkit/logger"
	promAPI "github.com/prometheus/client_golang/api"
	promQAPI "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
)

// PrometheusClient represents a prometheus client in Meshery
type PrometheusClient struct {
	grafanaClient *GrafanaClient
	logger        *logger.Handler
	//lint:ignore U1000 PromURL is not useless field over here but the rule will not consider function arguments as a valid option.
	promURL string
}

// NewPrometheusClient returns a PrometheusClient
func NewPrometheusClient(log *logger.Handler) *PrometheusClient {
	return NewPrometheusClientWithHTTPClient(&http.Client{}, log)
}

// NewPrometheusClientWithHTTPClient returns a PrometheusClient with a given http.Client
func NewPrometheusClientWithHTTPClient(client *http.Client, log *logger.Handler) *PrometheusClient {
	return &PrometheusClient{
		grafanaClient: NewGrafanaClientForPrometheusWithHTTPClient(client, log),
	}
}

// Validate - helps validate the connection
func (p *PrometheusClient) Validate(ctx context.Context, promURL, apiKeyOrBasicAuth string) error {
	_, err := p.grafanaClient.makeRequest(ctx, promURL+"/api/v1/status/config", apiKeyOrBasicAuth)
	if err != nil {
		return err
	}
	return nil
}

// ImportGrafanaBoard takes raw Grafana board json and returns GrafanaBoard pointer for use in Meshery
func (p *PrometheusClient) ImportGrafanaBoard(ctx context.Context, boardData []byte) (*GrafanaBoard, error) {
	board := &sdk.Board{}
	if err := json.Unmarshal(boardData, board); err != nil {
		return nil, ErrUnmarshal(err, "Grafana Board")
	}
	return p.grafanaClient.ProcessBoard(ctx, nil, board, &sdk.FoundBoard{
		Title: board.Title,
		URI:   board.Slug,
	})
}

// Query queries prometheus using the GrafanaClient
func (p *PrometheusClient) Query(ctx context.Context, promURL string, queryData *url.Values) ([]byte, error) {
	return p.grafanaClient.GrafanaQuery(ctx, promURL, "", queryData)
}

// QueryRange parses the given params and performs Prometheus range queries
func (p *PrometheusClient) QueryRange(ctx context.Context, promURL string, queryData *url.Values) ([]byte, error) {
	if queryData == nil {
		return nil, ErrNilQuery
	}
	reqURL := fmt.Sprintf("%s/api/v1/query_range", promURL)
	newURL, _ := url.Parse(reqURL)
	q := url.Values{}
	q.Set("query", queryData.Get("query"))
	q.Set("start", queryData.Get("start"))
	q.Set("end", queryData.Get("end"))
	q.Set("step", queryData.Get("step"))
	newURL.RawQuery = q.Encode()
	queryURL := newURL.String()
	data, err := p.grafanaClient.makeRequest(ctx, queryURL, "")
	if err != nil {
		return nil, ErrGrafanaData(err, queryURL)
	}
	return data, nil
}

// GetClusterStaticBoard retrieves the cluster static board config
func (p *PrometheusClient) GetClusterStaticBoard(ctx context.Context, _ string) (*GrafanaBoard, error) {
	return p.ImportGrafanaBoard(ctx, []byte(staticBoardCluster))
}

// Close - closes idle connections
func (p *PrometheusClient) Close() {
	p.grafanaClient.Close()
}

// GetNodesStaticBoard retrieves the per node static board config
func (p *PrometheusClient) GetNodesStaticBoard(ctx context.Context, promURL string) (*GrafanaBoard, error) {
	var buf bytes.Buffer
	ttt := template.New("staticBoard").Delims("[[", "]]")
	instances, err := p.getAllNodes(ctx, promURL)
	if err != nil {
		return nil, ErrPrometheusGetNodes(err)
	}
	(*p.logger).Debug(fmt.Sprintf("Instances: %v, length: %d", instances, len(instances)))
	tpl := template.Must(ttt.Parse(staticBoardNodes))
	if err := tpl.Execute(&buf, map[string]interface{}{
		"instances":  instances,
		"indexCheck": len(instances) - 1,
	}); err != nil {
		return nil, ErrPrometheusStaticBoard(err)
	}
	// logrus.Debugf("Board json: %s", buf.String())
	return p.ImportGrafanaBoard(ctx, buf.Bytes())
}

func (p *PrometheusClient) getAllNodes(ctx context.Context, promURL string) ([]string, error) {
	// api/datasources/proxy/1/api/v1/series?match[]=node_boot_time_seconds%7Bcluster%3D%22%22%2C%20job%3D%22node-exporter%22%7D&start=1568392571&end=1568396171
	c, _ := promAPI.NewClient(promAPI.Config{
		Address: promURL,
	})
	qc := promQAPI.NewAPI(c)
	labelSet, _, err := qc.Series(ctx, []string{`node_boot_time_seconds{cluster="", job="node-exporter"}`}, time.Now().Add(-5*time.Minute), time.Now())
	if err != nil {
		return nil, ErrPrometheusLabelSeries(err)
	}
	result := []string{}
	for _, l := range labelSet {
		inst := l["instance"]
		ins := string(inst)
		if ins != "" {
			result = append(result, ins)
		}
	}
	return result, nil
}

// QueryRangeUsingClient performs a range query within a window
func (p *PrometheusClient) QueryRangeUsingClient(ctx context.Context, promURL, query string, startTime, endTime time.Time, step time.Duration) (promModel.Value, error) {
	c, _ := promAPI.NewClient(promAPI.Config{
		Address: promURL,
	})
	qc := promQAPI.NewAPI(c)
	result, _, err := qc.QueryRange(ctx, query, promQAPI.Range{
		Start: startTime,
		End:   endTime,
		Step:  step,
	})
	if err != nil {
		return nil, ErrPrometheusQueryRange(err, query, startTime, endTime, step)
	}
	return result, nil
}

// ComputeStep computes the step size for a window
func (p *PrometheusClient) ComputeStep(_ context.Context, start, end time.Time) time.Duration {
	diff := end.Sub(start)
	// all calc. here are approx.
	if diff <= 10*time.Minute { // 10 mins
		return 5 * time.Second
	} else if diff <= 30*time.Minute { // 30 mins
		return 10 * time.Second
	} else if diff > 30*time.Minute && diff <= time.Hour { // 60 mins/1hr
		return 20 * time.Second
	} else if diff > 1*time.Hour && diff <= 3*time.Hour { // 3 time.Hour
		return 1 * time.Minute
	} else if diff > 3*time.Hour && diff <= 6*time.Hour { // 6 time.Hour
		return 2 * time.Minute
	} else if diff > 6*time.Hour && diff <= 1*24*time.Hour { // 24 time.Hour/1 day
		return 8 * time.Minute
	} else if diff > 1*24*time.Hour && diff <= 2*24*time.Hour { // 2 24*time.Hour
		return 16 * time.Minute
	} else if diff > 2*24*time.Hour && diff <= 4*24*time.Hour { // 4 24*time.Hour
		return 32 * time.Minute
	} else if diff > 4*24*time.Hour && diff <= 7*24*time.Hour { // 7 24*time.Hour
		return 56 * time.Minute
	} else if diff > 7*24*time.Hour && diff <= 15*24*time.Hour { // 15 24*time.Hour
		return 2 * time.Hour
	} else if diff > 15*24*time.Hour && diff <= 1*30*24*time.Hour { // 30 24*time.Hour/1 month
		return 4 * time.Hour
	} else if diff > 1*30*24*time.Hour && diff <= 3*30*24*time.Hour { // 3 months
		return 12 * time.Hour
	} else if diff > 3*30*24*time.Hour && diff <= 6*30*24*time.Hour { // 6 months
		return 1 * 24 * time.Hour
	} else if diff > 6*30*24*time.Hour && diff <= 1*12*30*24*time.Hour { // 1 year/12 months
		return 2 * 24 * time.Hour
	} else if diff > 1*12*30*24*time.Hour && diff <= 2*12*30*24*time.Hour { // 2 years
		return 4 * 24 * time.Hour
	} else if diff > 2*12*30*24*time.Hour && diff <= 5*12*30*24*time.Hour { // 5 years
		return 10 * 24 * time.Hour
	} else {
		return 30 * 24 * time.Hour
	}
}
