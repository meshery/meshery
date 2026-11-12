package grafana

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func TestAuthHeader(t *testing.T) {
	cases := []struct {
		name   string
		secret string
		want   string
	}{
		{"empty", "", ""},
		{"bearer token", "glsa_abc123", "Bearer glsa_abc123"},
		{"basic auth", "admin:password", "Basic YWRtaW46cGFzc3dvcmQ="},
		{"trimmed", "  token  ", "Bearer token"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := authHeader(tc.secret); got != tc.want {
				t.Fatalf("authHeader(%q) = %q, want %q", tc.secret, got, tc.want)
			}
		})
	}
}

func TestSearchBoards(t *testing.T) {
	var gotAuth, gotPath, gotQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		gotQuery = r.URL.Query().Get("query")
		_, _ = w.Write([]byte(`[
			{"uid":"abc","title":"Node Exporter","url":"/d/abc/node","type":"dash-db","tags":["prom"]},
			{"uid":"fold","title":"A Folder","type":"dash-folder"}
		]`))
	}))
	defer srv.Close()

	c := New(srv.URL, "glsa_token", nil)
	boards, err := c.SearchBoards(context.Background(), "node")
	if err != nil {
		t.Fatalf("SearchBoards: %v", err)
	}
	if gotAuth != "Bearer glsa_token" {
		t.Errorf("auth header = %q", gotAuth)
	}
	if gotPath != "/api/search" {
		t.Errorf("path = %q", gotPath)
	}
	if gotQuery != "node" {
		t.Errorf("query = %q", gotQuery)
	}
	if len(boards) != 1 {
		t.Fatalf("expected folders filtered out, got %d boards", len(boards))
	}
	if boards[0].UID != "abc" || boards[0].Title != "Node Exporter" {
		t.Errorf("unexpected board: %+v", boards[0])
	}
}

func TestGetBoardNormalizesPanelsAndDatasources(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"dashboard":{
			"uid":"abc","title":"Node","tags":["prom"],
			"templating":{"list":[{"name":"node","label":"Node","type":"query","query":{"query":"up"}}]},
			"panels":[
				{"id":1,"type":"row","title":"sep"},
				{"id":2,"type":"timeseries","title":"CPU","gridPos":{"x":0,"y":0,"w":12,"h":8},
					"datasource":{"type":"prometheus","uid":"ds1"},
					"fieldConfig":{"defaults":{"unit":"percent"}},
					"targets":[
						{"refId":"A","expr":"rate(cpu[5m])","legendFormat":"{{cpu}}"},
						{"refId":"B","expr":"","legendFormat":"empty"},
						{"refId":"C","expr":"mem","datasource":"ds2"}
					]}
			]}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	board, err := c.GetBoard(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetBoard: %v", err)
	}
	if len(board.Panels) != 1 {
		t.Fatalf("expected row panel skipped, got %d panels", len(board.Panels))
	}
	p := board.Panels[0]
	if p.Title != "CPU" || p.Unit != "percent" {
		t.Errorf("unexpected panel: %+v", p)
	}
	if len(p.Targets) != 2 {
		t.Fatalf("expected empty-expr target skipped, got %d targets", len(p.Targets))
	}
	if p.Targets[0].DatasourceUID != "ds1" { // inherits panel datasource
		t.Errorf("target A datasource = %q, want ds1", p.Targets[0].DatasourceUID)
	}
	if p.Targets[1].DatasourceUID != "ds2" { // legacy string datasource
		t.Errorf("target C datasource = %q, want ds2", p.Targets[1].DatasourceUID)
	}
	if len(board.TemplateVars) != 1 || board.TemplateVars[0].Query != "up" {
		t.Errorf("unexpected template vars: %+v", board.TemplateVars)
	}
}

func TestDatasourceUID(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"missing", ``, ""},
		{"null", `null`, ""},
		{"legacy string", `"ds-uid"`, "ds-uid"},
		{"object form", `{"type":"prometheus","uid":"ds1"}`, "ds1"},
		{"object without uid", `{"type":"prometheus"}`, ""},
		{"template ref", `"${datasource}"`, "${datasource}"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := datasourceUID(json.RawMessage(tc.raw)); got != tc.want {
				t.Errorf("datasourceUID(%s) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

func TestCurrentValues(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want []string
	}{
		{"missing", ``, nil},
		{"empty string", `""`, nil},
		{"single string", `"node-1"`, []string{"node-1"}},
		{"multi value", `["node-1","node-2"]`, []string{"node-1", "node-2"}},
		{"empty array", `[]`, []string{}},
		{"unexpected number", `5`, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := currentValues(json.RawMessage(tc.raw)); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("currentValues(%s) = %#v, want %#v", tc.raw, got, tc.want)
			}
		})
	}
}

func TestTemplateQueryString(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{"missing", ``, ""},
		{"plain string", `"up"`, "up"},
		{"object form", `{"query":"up","refId":"A"}`, "up"},
		{"object without query", `{"refId":"A"}`, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := templateQueryString(json.RawMessage(tc.raw)); got != tc.want {
				t.Errorf("templateQueryString(%s) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

// TestNormalizeBoardCollapsedRow asserts panels nested inside a collapsed row
// are recovered rather than dropped with the row separator.
func TestNormalizeBoardCollapsedRow(t *testing.T) {
	const body = `{"dashboard":{
		"uid":"abc","title":"Nested",
		"panels":[
			{"id":1,"type":"timeseries","title":"Top","gridPos":{"x":0,"y":0,"w":12,"h":8},
				"datasource":{"uid":"ds1"},"targets":[{"refId":"A","expr":"up"}]},
			{"id":2,"type":"row","title":"Collapsed","collapsed":true,"panels":[
				{"id":3,"type":"stat","title":"Hidden","gridPos":{"x":0,"y":9,"w":6,"h":4},
					"datasource":{"uid":"ds1"},"targets":[{"refId":"B","expr":"mem"}]}
			]}
		]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	board, err := New(srv.URL, "", nil).GetBoard(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetBoard: %v", err)
	}
	if len(board.Panels) != 2 {
		t.Fatalf("expected row dropped and its child recovered (2 panels), got %d: %+v", len(board.Panels), board.Panels)
	}
	titles := []string{board.Panels[0].Title, board.Panels[1].Title}
	if !reflect.DeepEqual(titles, []string{"Top", "Hidden"}) {
		t.Errorf("unexpected panel titles: %v", titles)
	}
}

// TestNormalizeBoardTemplateVars exercises multi-value, include-all, and
// single-value template variable normalization end to end.
func TestNormalizeBoardTemplateVars(t *testing.T) {
	const body = `{"dashboard":{
		"uid":"abc","title":"Vars","panels":[],
		"templating":{"list":[
			{"name":"node","label":"Node","type":"query","query":"up","multi":true,"includeAll":true,
				"current":{"value":["n1","n2"]}},
			{"name":"job","type":"custom","query":{"query":"a,b"},"current":{"value":"a"}}
		]}}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()

	board, err := New(srv.URL, "", nil).GetBoard(context.Background(), "abc")
	if err != nil {
		t.Fatalf("GetBoard: %v", err)
	}
	if len(board.TemplateVars) != 2 {
		t.Fatalf("expected 2 template vars, got %d", len(board.TemplateVars))
	}
	node := board.TemplateVars[0]
	if !node.Multi || !node.IncludeAll || node.Query != "up" {
		t.Errorf("unexpected node var: %+v", node)
	}
	if !reflect.DeepEqual(node.Current, []string{"n1", "n2"}) {
		t.Errorf("node.Current = %v, want [n1 n2]", node.Current)
	}
	job := board.TemplateVars[1]
	if job.Query != "a,b" || !reflect.DeepEqual(job.Current, []string{"a"}) {
		t.Errorf("unexpected job var: %+v", job)
	}
}

func TestQueryRangeProxiesVerbatim(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	var gotPath string
	var gotParams url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotParams = r.URL.Query()
		_, _ = w.Write([]byte(promResp))
	}))
	defer srv.Close()

	c := New(srv.URL, "admin:pass", nil)
	params := url.Values{}
	params.Set("query", "up")
	params.Set("start", "1")
	params.Set("end", "2")
	params.Set("step", "15")
	body, err := c.QueryRange(context.Background(), "ds1", params)
	if err != nil {
		t.Fatalf("QueryRange: %v", err)
	}
	if string(body) != promResp {
		t.Errorf("body not proxied verbatim: %s", body)
	}
	if gotPath != "/api/datasources/proxy/uid/ds1/api/v1/query_range" {
		t.Errorf("path = %q", gotPath)
	}
	if gotParams.Get("query") != "up" || gotParams.Get("step") != "15" {
		t.Errorf("params not forwarded: %v", gotParams)
	}
}

// TestQueryRangeResolvesDatasourceName covers a provisioned-dashboard panel that
// references its datasource by name: the first proxy-by-uid attempt 404s, the
// client resolves the name to a real uid via /api/datasources, and retries.
func TestQueryRangeResolvesDatasourceName(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	var proxiedUIDs []string
	dsListCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/datasources":
			dsListCalls++
			_, _ = w.Write([]byte(`[{"uid":"abc123","name":"dev-cortex","type":"prometheus"}]`))
		case strings.HasPrefix(r.URL.Path, "/api/datasources/proxy/uid/"):
			uid := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/datasources/proxy/uid/"), "/api/v1/query_range")
			proxiedUIDs = append(proxiedUIDs, uid)
			if uid == "abc123" {
				_, _ = w.Write([]byte(promResp))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"Unable to find datasource"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	body, err := c.QueryRange(context.Background(), "dev-cortex", url.Values{})
	if err != nil {
		t.Fatalf("QueryRange: %v", err)
	}
	if string(body) != promResp {
		t.Errorf("body = %s, want proxied prom response", body)
	}
	if dsListCalls != 1 {
		t.Errorf("expected 1 datasource-list lookup, got %d", dsListCalls)
	}
	want := []string{"dev-cortex", "abc123"} // failed name attempt, then resolved uid
	if !reflect.DeepEqual(proxiedUIDs, want) {
		t.Errorf("proxied uids = %v, want %v", proxiedUIDs, want)
	}
}

// TestQueryRangeValidUIDNoResolution ensures a working uid is proxied directly
// without an extra datasource-list lookup (no overhead on the happy path).
func TestQueryRangeValidUIDNoResolution(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	dsListCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/datasources" {
			dsListCalls++
		}
		_, _ = w.Write([]byte(promResp))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	if _, err := c.QueryRange(context.Background(), "good-uid", url.Values{}); err != nil {
		t.Fatalf("QueryRange: %v", err)
	}
	if dsListCalls != 0 {
		t.Errorf("expected no datasource-list lookup on success, got %d", dsListCalls)
	}
}

// TestQueryRangeUnresolvableSurfacesOriginalError verifies that when the ref
// can't be resolved, the original 404 is returned rather than masked.
func TestQueryRangeUnresolvableSurfacesOriginalError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/datasources" {
			_, _ = w.Write([]byte(`[{"uid":"abc123","name":"dev-cortex","type":"prometheus"}]`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"Unable to find datasource"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	_, err := c.QueryRange(context.Background(), "ghost", url.Values{})
	if err == nil {
		t.Fatal("expected error for unresolvable datasource, got nil")
	}
	var dsErr *DatasourceNotFoundError
	if !errors.As(err, &dsErr) {
		t.Fatalf("expected *DatasourceNotFoundError, got %T: %v", err, err)
	}
	if dsErr.Ref != "ghost" {
		t.Errorf("DatasourceNotFoundError.Ref = %q, want %q", dsErr.Ref, "ghost")
	}
	// The error should self-diagnose by listing what Grafana actually has.
	if !reflect.DeepEqual(dsErr.Available, []string{"dev-cortex (uid=abc123)"}) {
		t.Errorf("Available = %v, want [dev-cortex (uid=abc123)]", dsErr.Available)
	}
}

// TestQueryRangeBatch exercises the fan-out path end to end: a query with a valid
// uid succeeds directly, a query whose ds is a NAME 404s then resolves via
// /api/datasources and succeeds, and a query with an unresolvable ds yields a
// per-result error while the others still succeed. It also asserts the datasource
// list is fetched at most once for the whole batch.
func TestQueryRangeBatch(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	var mu sync.Mutex
	dsListCalls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/datasources":
			mu.Lock()
			dsListCalls++
			mu.Unlock()
			_, _ = w.Write([]byte(`[{"uid":"abc123","name":"dev-cortex","type":"prometheus"}]`))
		case strings.HasPrefix(r.URL.Path, "/api/datasources/proxy/uid/"):
			uid := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/datasources/proxy/uid/"), "/api/v1/query_range")
			if uid == "good-uid" || uid == "abc123" {
				_, _ = w.Write([]byte(promResp))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"Unable to find datasource"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	queries := []BatchQuery{
		{ID: "1:A", DS: "good-uid", Query: "up"},      // valid uid, direct success
		{ID: "2:B", DS: "dev-cortex", Query: "mem"},   // name -> resolves to abc123
		{ID: "3:C", DS: "ghost", Query: "cpu"},        // unresolvable -> per-result error
	}
	results := c.QueryRangeBatch(context.Background(), "1", "2", "15", queries)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	byID := map[string]BatchResult{}
	for _, r := range results {
		byID[r.ID] = r
	}

	if r := byID["1:A"]; r.Err != nil || string(r.Body) != promResp {
		t.Errorf("1:A = %+v, want success with prom body", r)
	}
	if r := byID["2:B"]; r.Err != nil || string(r.Body) != promResp {
		t.Errorf("2:B = %+v, want resolved success with prom body", r)
	}
	r := byID["3:C"]
	if r.Err == nil {
		t.Fatalf("3:C expected a per-result error, got nil")
	}
	var dsErr *DatasourceNotFoundError
	if !errors.As(r.Err, &dsErr) || dsErr.Ref != "ghost" {
		t.Errorf("3:C error = %v, want DatasourceNotFoundError{Ref: ghost}", r.Err)
	}

	if dsListCalls != 1 {
		t.Errorf("expected 1 datasource-list lookup for the whole batch, got %d", dsListCalls)
	}
}

func TestQueryRangeBatchEmpty(t *testing.T) {
	c := New("http://example.invalid", "", nil)
	if got := c.QueryRangeBatch(context.Background(), "1", "2", "15", nil); len(got) != 0 {
		t.Errorf("empty batch = %v, want no results", got)
	}
}

func TestDoSurfacesNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"message":"invalid api key"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "bad", nil)
	if _, err := c.Health(context.Background()); err == nil {
		t.Fatal("expected error on 401, got nil")
	}
}

func TestHealthReturnsVersion(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"commit":"abc","database":"ok","version":"10.4.2"}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "glsa_token", nil)
	info, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if info.Version != "10.4.2" || info.Database != "ok" {
		t.Errorf("unexpected health info: %+v", info)
	}
}
