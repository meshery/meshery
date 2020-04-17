package models

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/grafana-tools/sdk"
	"github.com/pkg/errors"
	promAPI "github.com/prometheus/client_golang/api"
	promQAPI "github.com/prometheus/client_golang/api/prometheus/v1"
	promModel "github.com/prometheus/common/model"
	"github.com/sirupsen/logrus"
)

// PrometheusClient represents a prometheus client in Meshery
type PrometheusClient struct {
	grafanaClient *GrafanaClient
	//lint:ignore U1000 PromURL is not useless field over here but the rule will not consider function arguments as a valid option.
	promURL string
}

// NewPrometheusClient returns a PrometheusClient
func NewPrometheusClient() *PrometheusClient {
	return NewPrometheusClientWithHTTPClient(&http.Client{})
}

// NewPrometheusClientWithHTTPClient returns a PrometheusClient with a given http.Client
func NewPrometheusClientWithHTTPClient(client *http.Client) *PrometheusClient {
	return &PrometheusClient{
		grafanaClient: NewGrafanaClientForPrometheusWithHTTPClient(client),
	}
}

// Validate - helps validate the connection
func (p *PrometheusClient) Validate(ctx context.Context, promURL string) error {
	_, err := p.grafanaClient.makeRequest(ctx, promURL+"/api/v1/status/config", "")
	if err != nil {
		return err
	}
	return nil
}

// ImportGrafanaBoard takes raw Grafana board json and returns GrafanaBoard pointer for use in Meshery
func (p *PrometheusClient) ImportGrafanaBoard(ctx context.Context, boardData []byte) (*GrafanaBoard, error) {
	board := &sdk.Board{}
	if err := json.Unmarshal(boardData, board); err != nil {
		msg := errors.New("unable to parse grafana board data")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
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

// QueryRange queries prometheus using the GrafanaClient
func (p *PrometheusClient) QueryRange(ctx context.Context, promURL string, queryData *url.Values) ([]byte, error) {
	return p.grafanaClient.GrafanaQueryRange(ctx, promURL, "", queryData)
}

// GetClusterStaticBoard retrieves the cluster static board config
func (p *PrometheusClient) GetClusterStaticBoard(ctx context.Context, promURL string) (*GrafanaBoard, error) {
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
		err = errors.Wrapf(err, "unable to get all the nodes")
		logrus.Error(err)
		return nil, err
	}
	logrus.Debugf("Instances: %v, length: %d", instances, len(instances))
	tpl := template.Must(ttt.Parse(staticBoardNodes))
	if err := tpl.Execute(&buf, map[string]interface{}{
		"instances":  instances,
		"indexCheck": len(instances) - 1,
	}); err != nil {
		err = errors.Wrapf(err, "unable to get the static board")
		logrus.Error(err)
		return nil, err
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
		err = errors.Wrapf(err, "unable to get the label set series")
		logrus.Error(err)
		return nil, err
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
		err := errors.Wrapf(err, "error fetching data for query: %s, with start: %v, end: %v, step: %v", query, startTime, endTime, step)
		logrus.Error(err)
		return nil, err
	}
	return result, nil
}

// ComputeStep computes the step size for a window
func (p *PrometheusClient) ComputeStep(ctx context.Context, start, end time.Time) time.Duration {
	step := 5 * time.Second
	diff := end.Sub(start)
	// all calc. here are approx.
	if diff <= 10*time.Minute { // 10 mins
		step = 5 * time.Second
	} else if diff <= 30*time.Minute { // 30 mins
		step = 10 * time.Second
	} else if diff > 30*time.Minute && diff <= time.Hour { // 60 mins/1hr
		step = 20 * time.Second
	} else if diff > 1*time.Hour && diff <= 3*time.Hour { // 3 time.Hour
		step = 1 * time.Minute
	} else if diff > 3*time.Hour && diff <= 6*time.Hour { // 6 time.Hour
		step = 2 * time.Minute
	} else if diff > 6*time.Hour && diff <= 1*24*time.Hour { // 24 time.Hour/1 day
		step = 8 * time.Minute
	} else if diff > 1*24*time.Hour && diff <= 2*24*time.Hour { // 2 24*time.Hour
		step = 16 * time.Minute
	} else if diff > 2*24*time.Hour && diff <= 4*24*time.Hour { // 4 24*time.Hour
		step = 32 * time.Minute
	} else if diff > 4*24*time.Hour && diff <= 7*24*time.Hour { // 7 24*time.Hour
		step = 56 * time.Minute
	} else if diff > 7*24*time.Hour && diff <= 15*24*time.Hour { // 15 24*time.Hour
		step = 2 * time.Hour
	} else if diff > 15*24*time.Hour && diff <= 1*30*24*time.Hour { // 30 24*time.Hour/1 month
		step = 4 * time.Hour
	} else if diff > 1*30*24*time.Hour && diff <= 3*30*24*time.Hour { // 3 months
		step = 12 * time.Hour
	} else if diff > 3*30*24*time.Hour && diff <= 6*30*24*time.Hour { // 6 months
		step = 1 * 24 * time.Hour
	} else if diff > 6*30*24*time.Hour && diff <= 1*12*30*24*time.Hour { // 1 year/12 months
		step = 2 * 24 * time.Hour
	} else if diff > 1*12*30*24*time.Hour && diff <= 2*12*30*24*time.Hour { // 2 years
		step = 4 * 24 * time.Hour
	} else if diff > 2*12*30*24*time.Hour && diff <= 5*12*30*24*time.Hour { // 5 years
		step = 10 * 24 * time.Hour
	} else {
		step = 30 * 24 * time.Hour
	}
	return step
}
