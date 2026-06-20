package grafana

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

// apiError is returned by do for non-2xx Grafana responses. It carries the
// status code so callers can react to specific conditions (e.g. a 404 when a
// datasource reference turns out to be a name rather than a uid).
type apiError struct {
	StatusCode int
	Method     string
	Path       string
	Body       string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("grafana: %s %s returned %d: %s", e.Method, e.Path, e.StatusCode, e.Body)
}

// DatasourceNotFoundError indicates a panel referenced a datasource the Grafana
// instance could not resolve — typically a name or stale uid carried by a
// provisioned dashboard, or a datasource the credential is not allowed to see.
// It is distinct from a generic API error so callers can explain it precisely.
//
// Available lists the datasources Grafana actually returned ("name (uid=...)"),
// so the user can immediately see whether Ref is missing entirely (dangling
// reference) or the list came back empty (credential can't read datasources).
type DatasourceNotFoundError struct {
	Ref       string
	Available []string
}

func (e *DatasourceNotFoundError) Error() string {
	if len(e.Available) == 0 {
		return fmt.Sprintf("grafana: datasource %q was not found, and no datasources were listable — the credential may lack permission to read datasources", e.Ref)
	}
	return fmt.Sprintf("grafana: datasource %q was not found; available datasources: %s", e.Ref, strings.Join(e.Available, ", "))
}

// StatusCode returns the HTTP status carried by a Grafana API error, or 0 when
// err is not one. It lets callers classify failures (auth vs not-found vs ...)
// without depending on the unexported error type.
func StatusCode(err error) int {
	var apiErr *apiError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode
	}
	return 0
}

// do issues a request against the Grafana API and returns the response body.
// Non-2xx responses are surfaced as *apiError carrying the status and a trimmed body.
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
		return nil, &apiError{
			StatusCode: resp.StatusCode,
			Method:     method,
			Path:       path,
			Body:       strings.TrimSpace(string(body)),
		}
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
func (c *Client) QueryRange(ctx context.Context, dsRef string, params url.Values) ([]byte, error) {
	if strings.TrimSpace(dsRef) == "" {
		return nil, fmt.Errorf("grafana: empty datasource uid")
	}
	body, err := c.queryRangeByUID(ctx, dsRef, params)
	if err == nil {
		return body, nil
	}
	// A 404 means Grafana has no datasource with that uid. Provisioned dashboards
	// frequently reference datasources by name (or a stale uid) instead, so try to
	// resolve the reference to a real uid and retry once before giving up.
	var apiErr *apiError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		// dsRef may be a name (or stale uid) rather than a real uid. Resolve it
		// against the live datasource list and retry once before giving up.
		sources, _ := c.ListDatasources(ctx)
		if uid := datasourceUIDByName(sources, dsRef); uid != "" && uid != dsRef {
			if body, retryErr := c.queryRangeByUID(ctx, uid, params); retryErr == nil {
				return body, nil
			}
		}
		// Still unresolved. Surface a typed error naming the reference and what
		// datasources are actually available, so the user can see precisely why.
		return nil, &DatasourceNotFoundError{Ref: dsRef, Available: datasourceIdentifiers(sources)}
	}
	return nil, err
}

func (c *Client) queryRangeByUID(ctx context.Context, uid string, params url.Values) ([]byte, error) {
	path := "/api/datasources/proxy/uid/" + url.PathEscape(uid) + "/api/v1/query_range"
	return c.do(ctx, http.MethodGet, path, params)
}

// batchConcurrency bounds how many query_range calls a single batch fans out to
// Grafana at once, so one board with many panels can't open an unbounded number
// of simultaneous connections to the upstream instance.
const batchConcurrency = 8

// BatchQuery is a single resolved query within a batch. ID is an opaque,
// caller-chosen key (e.g. "<panelID>:<refID>") echoed back on the matching
// result; DS is the datasource reference (uid or name); Query is the resolved
// PromQL expression.
type BatchQuery struct {
	ID    string
	DS    string
	Query string
}

// BatchResult carries the outcome of one BatchQuery. ID matches the input query.
// Exactly one of Body (the raw Prometheus query_range JSON, unmodified) or Err
// (a per-query failure) is set.
type BatchResult struct {
	ID   string
	Body json.RawMessage
	Err  error
}

// QueryRangeBatch runs many Prometheus range queries against this Grafana
// instance concurrently and returns one result per input query, in input order.
//
// A single failing query never fails the whole batch: its result simply carries
// an Err while the others succeed. The datasource list is fetched at most once
// for the entire batch (best-effort; an error is treated as an empty list) and
// shared across every query's name-resolution fallback, so a board with many
// panels referencing datasources by name doesn't trigger a list lookup per query.
//
// start, end and step are the shared Prometheus query_range window for the batch.
func (c *Client) QueryRangeBatch(ctx context.Context, start, end, step string, queries []BatchQuery) []BatchResult {
	results := make([]BatchResult, len(queries))
	if len(queries) == 0 {
		return results
	}

	// Resolve the datasource list once for the whole batch (best-effort). A
	// failure here is non-fatal: name resolution simply has nothing to match
	// against, and individual queries fall back to a DatasourceNotFoundError.
	sources, _ := c.ListDatasources(ctx)

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
			params := url.Values{
				"query": {q.Query},
				"start": {start},
				"end":   {end},
				"step":  {step},
			}
			body, err := c.queryRangeWith(ctx, q.DS, params, sources)
			results[i] = BatchResult{ID: q.ID, Body: body, Err: err}
		}(i, q)
	}
	wg.Wait()
	return results
}

// queryRangeWith runs a single range query against dsRef, reusing the already
// resolved datasource list (sources) instead of fetching it per query. It mirrors
// QueryRange's 404 -> resolve-name -> retry-once dance, surfacing a typed
// DatasourceNotFoundError when the reference still can't be resolved.
func (c *Client) queryRangeWith(ctx context.Context, dsRef string, params url.Values, sources []Datasource) (json.RawMessage, error) {
	if strings.TrimSpace(dsRef) == "" {
		return nil, fmt.Errorf("grafana: empty datasource uid")
	}
	body, err := c.queryRangeByUID(ctx, dsRef, params)
	if err == nil {
		return body, nil
	}
	var apiErr *apiError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
		if uid := datasourceUIDByName(sources, dsRef); uid != "" && uid != dsRef {
			if body, retryErr := c.queryRangeByUID(ctx, uid, params); retryErr == nil {
				return body, nil
			}
		}
		return nil, &DatasourceNotFoundError{Ref: dsRef, Available: datasourceIdentifiers(sources)}
	}
	return nil, err
}

// datasourceUIDByName maps a datasource reference that is actually a name (the
// common provisioned-dashboard case) to its uid. A uid match takes precedence,
// then a name match. Returns "" when nothing matches.
func datasourceUIDByName(sources []Datasource, ref string) string {
	for _, ds := range sources {
		if ds.UID == ref { // already a valid uid
			return ds.UID
		}
	}
	for _, ds := range sources {
		if ds.Name == ref {
			return ds.UID
		}
	}
	return ""
}

// datasourceIdentifiers renders each datasource as "name (uid=...)" for surfacing
// in a not-found error so the user can see what is actually available.
func datasourceIdentifiers(sources []Datasource) []string {
	out := make([]string, 0, len(sources))
	for _, ds := range sources {
		out = append(out, fmt.Sprintf("%s (uid=%s)", ds.Name, ds.UID))
	}
	return out
}

func normalizeBoard(d rawDashboard) *Board {
	board := &Board{
		UID:    d.UID,
		Title:  d.Title,
		Tags:   d.Tags,
		Panels: make([]Panel, 0, len(d.Panels)),
	}
	for _, rp := range flattenPanels(d.Panels) {
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

// flattenPanels drops row separators while recovering panels nested inside a
// collapsed row. Grafana stores a collapsed row's children in the row's own
// "panels" array, whereas an expanded row's children are top-level siblings;
// without descending, panels under a collapsed row would be silently lost.
func flattenPanels(panels []rawPanel) []rawPanel {
	out := make([]rawPanel, 0, len(panels))
	for _, p := range panels {
		if p.Type == "row" {
			out = append(out, flattenPanels(p.Panels)...)
			continue
		}
		out = append(out, p)
	}
	return out
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
