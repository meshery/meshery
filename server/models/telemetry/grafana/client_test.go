package grafana

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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
