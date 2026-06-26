package prometheus

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/meshery/meshkit/logger"
)

// defaultTimeout bounds every outbound call to Prometheus so a slow or
// unreachable instance can never hang a Meshery request indefinitely.
const defaultTimeout = 30 * time.Second

// Client is a minimal, dependency-light Prometheus HTTP API client. A Client is
// scoped to a single Prometheus instance + credential and is safe for concurrent
// use.
type Client struct {
	baseURL    string
	authHeader string
	httpClient *http.Client
	log        logger.Handler
}

// New builds a Client for the given Prometheus base URL and secret.
//
// The secret is interpreted the same way the connection wizard captures it:
//   - "username:password" (contains a colon) -> HTTP Basic auth
//   - any other non-empty value             -> Bearer token (API key / service-account token)
//   - empty                                  -> no auth (anonymous Prometheus)
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

// BaseURL returns the normalized Prometheus base URL.
func (c *Client) BaseURL() string { return c.baseURL }

// apiError is returned by do for non-2xx Prometheus responses. It carries the
// status code so callers can classify failures (auth vs not-found vs ...).
type apiError struct {
	StatusCode int
	Method     string
	Path       string
	Body       string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("prometheus: %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Body)
}

// StatusCode returns the HTTP status carried by a Prometheus API error, or 0 when
// err is not one. It lets callers classify failures without depending on the
// unexported error type.
func StatusCode(err error) int {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return 0
}

// do issues a request against the Prometheus API and returns the response body.
// Non-2xx responses are surfaced as *apiError carrying the status and a trimmed body.
func (c *Client) do(ctx context.Context, method, path string, query url.Values) ([]byte, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("prometheus: empty base URL")
	}
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("prometheus: build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("prometheus: request to %s failed: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20)) // cap at 16MiB
	if err != nil {
		return nil, fmt.Errorf("prometheus: read response from %s: %w", path, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &apiError{
			StatusCode: resp.StatusCode,
			Method:     method,
			Path:       path,
			Body:       strings.TrimSpace(string(body)),
		}
	}
	return body, nil
}

// Health verifies the instance is reachable and returns its reported version. A
// nil error means the instance answered successfully.
func (c *Client) Health(ctx context.Context) (*HealthInfo, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/status/buildinfo", nil)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Data struct {
			Version string `json:"version"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		// Reachable, but the body wasn't the expected shape — report reachable
		// with an unknown version rather than failing the ping.
		return &HealthInfo{}, nil
	}
	return &HealthInfo{Version: raw.Data.Version}, nil
}

// QueryRange runs a Prometheus range query and returns the raw response body
// verbatim, so the UI can render it without the server reshaping time-series data.
func (c *Client) QueryRange(ctx context.Context, query, start, end, step string) ([]byte, error) {
	params := url.Values{
		"query": {query},
		"start": {start},
		"end":   {end},
		"step":  {step},
	}
	return c.do(ctx, http.MethodGet, "/api/v1/query_range", params)
}

// Query runs an instant Prometheus query and returns the raw response body
// verbatim. ts is the optional evaluation timestamp (empty means "now").
func (c *Client) Query(ctx context.Context, query, ts string) ([]byte, error) {
	params := url.Values{"query": {query}}
	if strings.TrimSpace(ts) != "" {
		params.Set("time", ts)
	}
	return c.do(ctx, http.MethodGet, "/api/v1/query", params)
}

// batchConcurrency bounds how many query_range calls a single batch fans out to
// Prometheus at once, so one dashboard with many panels can't open an unbounded
// number of simultaneous connections to the upstream instance.
const batchConcurrency = 8

// QueryRangeBatch runs many Prometheus range queries against this instance
// concurrently and returns one result per input query, in input order.
//
// A single failing query never fails the whole batch: its result simply carries
// an Err while the others succeed.
//
// start, end and step are the shared Prometheus query_range window for the batch.
func (c *Client) QueryRangeBatch(ctx context.Context, start, end, step string, queries []BatchQuery) []BatchResult {
	results := make([]BatchResult, len(queries))
	if len(queries) == 0 {
		return results
	}

	sem := make(chan struct{}, batchConcurrency)
	var wg sync.WaitGroup
	for i, q := range queries {
		wg.Add(1)
		go func(i int, q BatchQuery) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results[i] = BatchResult{ID: q.ID, Err: ctx.Err()}
				return
			}
			defer func() { <-sem }()
			body, err := c.QueryRange(ctx, q.Query, start, end, step)
			results[i] = BatchResult{ID: q.ID, Body: body, Err: err}
		}(i, q)
	}
	wg.Wait()
	return results
}

// labelValuesResponse is the shape of Prometheus' label / series-name responses:
// a status plus a flat string array.
type labelValuesResponse struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// MetricNames returns the instance's metric names (the values of the __name__
// label), optionally filtered to those containing search (case-insensitive) and
// capped to limit (when limit > 0).
func (c *Client) MetricNames(ctx context.Context, search string, limit int) ([]string, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/label/__name__/values", nil)
	if err != nil {
		return nil, err
	}
	var resp labelValuesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("prometheus: decode metric names: %w", err)
	}
	search = strings.ToLower(strings.TrimSpace(search))
	out := make([]string, 0, len(resp.Data))
	for _, name := range resp.Data {
		if search != "" && !strings.Contains(strings.ToLower(name), search) {
			continue
		}
		out = append(out, name)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

// LabelNames returns the instance's label names, optionally constrained to a
// series selector (match[]) when match is non-empty.
func (c *Client) LabelNames(ctx context.Context, match string) ([]string, error) {
	q := url.Values{}
	if strings.TrimSpace(match) != "" {
		q.Set("match[]", match)
	}
	body, err := c.do(ctx, http.MethodGet, "/api/v1/labels", q)
	if err != nil {
		return nil, err
	}
	var resp labelValuesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("prometheus: decode label names: %w", err)
	}
	return resp.Data, nil
}

// LabelValues returns the values of a single label, optionally constrained to a
// series selector (match[]) when match is non-empty.
func (c *Client) LabelValues(ctx context.Context, label, match string) ([]string, error) {
	if strings.TrimSpace(label) == "" {
		return nil, fmt.Errorf("prometheus: empty label name")
	}
	q := url.Values{}
	if strings.TrimSpace(match) != "" {
		q.Set("match[]", match)
	}
	body, err := c.do(ctx, http.MethodGet, "/api/v1/label/"+url.PathEscape(label)+"/values", q)
	if err != nil {
		return nil, err
	}
	var resp labelValuesResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("prometheus: decode label values: %w", err)
	}
	return resp.Data, nil
}

// Metadata returns Prometheus' metric metadata response for metric (or all
// metrics when metric is empty), raw and unmodified for the UI to render.
func (c *Client) Metadata(ctx context.Context, metric string) ([]byte, error) {
	q := url.Values{}
	if strings.TrimSpace(metric) != "" {
		q.Set("metric", metric)
	}
	return c.do(ctx, http.MethodGet, "/api/v1/metadata", q)
}
