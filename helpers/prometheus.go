package helpers

import (
	"context"
	"encoding/json"
	"net/url"

	// promAPI "github.com/prometheus/client_golang/api"
	// promQAPI "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/grafana-tools/sdk"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type PrometheusClient struct {
	grafanaClient *GrafanaClient
}

func NewPrometheusClient(ctx context.Context, promURL string, validate bool) (*PrometheusClient, error) {
	// client, err := promAPI.NewClient(promAPI.Config{Address: promURL})
	// if err != nil {
	// 	msg := errors.New("unable to connect to prometheus")
	// 	logrus.Error(errors.Wrap(err, msg.Error()))
	// 	return nil, msg
	// }
	// queryAPI := promQAPI.NewAPI(client)
	// return &PrometheusClient{
	// 	client:      client,
	// 	queryClient: queryAPI,
	// }, nil
	p := &PrometheusClient{
		grafanaClient: NewGrafanaClientForPrometheus(promURL),
	}
	if validate {
		_, err := p.grafanaClient.makeRequest(ctx, promURL+"/api/v1/status/config")
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *PrometheusClient) ImportGrafanaBoard(ctx context.Context, boardData []byte) (*models.GrafanaBoard, error) {
	board := &sdk.Board{}
	err := json.Unmarshal(boardData, board)
	if err != nil {
		msg := errors.New("unable to parse grafana board data")
		logrus.Error(errors.Wrap(err, msg.Error()))
		return nil, msg
	}
	return p.grafanaClient.ProcessBoard(board, &sdk.FoundBoard{
		Title: board.Title,
		URI:   board.Slug,
	})
}

func (p *PrometheusClient) Query(ctx context.Context, queryData *url.Values) ([]byte, error) {
	return p.grafanaClient.GrafanaQuery(ctx, queryData)
}

func (p *PrometheusClient) QueryRange(ctx context.Context, queryData *url.Values) ([]byte, error) {
	return p.grafanaClient.GrafanaQueryRange(ctx, queryData)
}
