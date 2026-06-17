package grafana

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/meshery/meshkit/logger"
)

// defaultTimeout bounds every outbound call to Grafana so a slow or unreachable
// instance can never hang a Meshery request indefinitely.
const defaultTimeout = 30 * time.Second

// Client is a minimal, dependency-light Grafana HTTP API client. A Client is
// scoped to a single Grafana instance + credential and is safe for concurrent
// use.
type Client struct {
	baseURL    string
	authHeader string
	httpClient *http.Client
	log        logger.Handler
}

// New builds a Client for the given Grafana base URL and secret.
//
// The secret is interpreted the same way the connection wizard captures it:
//   - "username:password" (contains a colon) -> HTTP Basic auth
//   - any other non-empty value             -> Bearer token (API key / service-account token)
//   - empty                                  -> no auth (anonymous Grafana)
func New(baseURL, secret string, log logger.Handler) *Client {
	return &Client{
		baseURL:    strings.TrimRight(strings.TrimSpace(baseURL), "/"),
		authHeader: authHeader(secret),
		httpClient: &http.Client{Timeout: defaultTimeout},
		log:        log,
	}
}

func authHeader(secret string) string {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return ""
	}
	if strings.Contains(secret, ":") {
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(secret))
	}
	return "Bearer " + secret
}

// BaseURL returns the normalized Grafana base URL.
func (c *Client) BaseURL() string { return c.baseURL }

// do issues a request against the Grafana API and returns the response body.
// Non-2xx responses are surfaced as errors carrying the status and a trimmed body.
func (c *Client) do(ctx context.Context, method, path string, query url.Values) ([]byte, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("grafana: empty base URL")
	}
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("grafana: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("grafana: request to %s failed: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20)) // cap at 16MiB
	if err != nil {
		return nil, fmt.Errorf("grafana: read response from %s: %w", path, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("grafana: %s %s returned %d: %s", method, path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}

// Health verifies the instance is reachable and returns its reported version
// and database status. A nil error means the instance answered successfully.
func (c *Client) Health(ctx context.Context) (*HealthInfo, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/health", nil)
	if err != nil {
		return nil, err
	}
	var info HealthInfo
	if err := json.Unmarshal(body, &info); err != nil {
		// Reachable, but the body wasn't the expected shape — report reachable
		// with an unknown version rather than failing the ping.
		return &HealthInfo{}, nil
	}
	return &info, nil
}

// SearchBoards lists dashboards, optionally filtered by a free-text query.
func (c *Client) SearchBoards(ctx context.Context, query string) ([]BoardSummary, error) {
	q := url.Values{}
	q.Set("type", "dash-db")
	q.Set("limit", "5000")
	if strings.TrimSpace(query) != "" {
		q.Set("query", query)
	}
	body, err := c.do(ctx, http.MethodGet, "/api/search", q)
	if err != nil {
		return nil, err
	}
	var hits []rawSearchHit
	if err := json.Unmarshal(body, &hits); err != nil {
		return nil, fmt.Errorf("grafana: decode search response: %w", err)
	}
	boards := make([]BoardSummary, 0, len(hits))
	for _, h := range hits {
		if h.Type != "" && h.Type != "dash-db" {
			continue
		}
		boards = append(boards, BoardSummary{
			UID:   h.UID,
			Title: h.Title,
			URL:   h.URL,
			Tags:  h.Tags,
		})
	}
	return boards, nil
}

// GetBoard fetches a single dashboard by UID and normalizes it to a Board.
func (c *Client) GetBoard(ctx context.Context, uid string) (*Board, error) {
	if strings.TrimSpace(uid) == "" {
		return nil, fmt.Errorf("grafana: empty board uid")
	}
	body, err := c.do(ctx, http.MethodGet, "/api/dashboards/uid/"+url.PathEscape(uid), nil)
	if err != nil {
		return nil, err
	}
	var raw rawDashboardResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("grafana: decode dashboard response: %w", err)
	}
	return normalizeBoard(raw.Dashboard), nil
}

// ListDatasources returns the datasources configured in the Grafana instance.
func (c *Client) ListDatasources(ctx context.Context) ([]Datasource, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/datasources", nil)
	if err != nil {
		return nil, err
	}
	var ds []Datasource
	if err := json.Unmarshal(body, &ds); err != nil {
		return nil, fmt.Errorf("grafana: decode datasources response: %w", err)
	}
	return ds, nil
}

// QueryRange proxies a Prometheus-style range query through Grafana's datasource
// proxy and returns the raw Prometheus response body verbatim, so the UI can
// render it without the server reshaping time-series data.
//
// dsUID is the target datasource UID; params carries the Prometheus query_range
// parameters (query, start, end, step).
func (c *Client) QueryRange(ctx context.Context, dsUID string, params url.Values) ([]byte, error) {
	if strings.TrimSpace(dsUID) == "" {
		return nil, fmt.Errorf("grafana: empty datasource uid")
	}
	path := "/api/datasources/proxy/uid/" + url.PathEscape(dsUID) + "/api/v1/query_range"
	return c.do(ctx, http.MethodGet, path, params)
}

func normalizeBoard(d rawDashboard) *Board {
	board := &Board{
		UID:    d.UID,
		Title:  d.Title,
		Tags:   d.Tags,
		Panels: make([]Panel, 0, len(d.Panels)),
	}
	for _, rp := range d.Panels {
		if rp.Type == "row" { // row separators carry no data
			continue
		}
		p := Panel{
			ID:      rp.ID,
			Title:   rp.Title,
			Type:    rp.Type,
			GridPos: rp.GridPos,
			Unit:    rp.FieldConfig.Defaults.Unit,
			Targets: make([]Target, 0, len(rp.Targets)),
		}
		panelDS := datasourceUID(rp.Datasource)
		for _, rt := range rp.Targets {
			if strings.TrimSpace(rt.Expr) == "" {
				continue
			}
			dsUID := datasourceUID(rt.Datasource)
			if dsUID == "" {
				dsUID = panelDS
			}
			p.Targets = append(p.Targets, Target{
				RefID:         rt.RefID,
				DatasourceUID: dsUID,
				Expr:          rt.Expr,
				LegendFormat:  rt.LegendFormat,
			})
		}
		board.Panels = append(board.Panels, p)
	}
	for _, rv := range d.Templating.List {
		board.TemplateVars = append(board.TemplateVars, TemplateVar{
			Name:       rv.Name,
			Label:      rv.Label,
			Type:       rv.Type,
			Query:      templateQueryString(rv.Query),
			Multi:      rv.Multi,
			IncludeAll: rv.IncludeAll,
			Current:    currentValues(rv.Current.Value),
		})
	}
	return board
}

// currentValues flattens a template variable's "current.value", which may be a
// single string or an array of strings, into a string slice.
func currentValues(raw json.RawMessage) []string {
	if len(raw) == 0 {
		return nil
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		if s == "" {
			return nil
		}
		return []string{s}
	}
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		return arr
	}
	return nil
}

// datasourceUID extracts a datasource UID from Grafana's polymorphic
// "datasource" field, which may be a bare string (legacy) or an object
// {"type":..,"uid":..} (current).
func datasourceUID(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var obj struct {
		UID string `json:"uid"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil {
		return obj.UID
	}
	return ""
}

// templateQueryString flattens a template variable's "query", which may be a
// string or an object, into a display string.
func templateQueryString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return s
	}
	var obj struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(raw, &obj); err == nil {
		return obj.Query
	}
	return ""
}
