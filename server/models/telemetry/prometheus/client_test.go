package prometheus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
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

func TestHealthReturnsVersion(t *testing.T) {
	var gotAuth, gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"status":"success","data":{"version":"2.51.0","revision":"abc"}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "glsa_token", nil)
	info, err := c.Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if gotAuth != "Bearer glsa_token" {
		t.Errorf("auth header = %q", gotAuth)
	}
	if gotPath != "/api/v1/status/buildinfo" {
		t.Errorf("path = %q", gotPath)
	}
	if info.Version != "2.51.0" {
		t.Errorf("version = %q, want 2.51.0", info.Version)
	}
}

// TestHealthUnparseableStillReachable ensures a reachable instance whose body
// isn't the expected shape still counts as reachable (versionless).
func TestHealthUnparseableStillReachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	info, err := New(srv.URL, "", nil).Health(context.Background())
	if err != nil {
		t.Fatalf("Health: %v", err)
	}
	if info.Version != "" {
		t.Errorf("version = %q, want empty", info.Version)
	}
}

func TestHealthSurfacesNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"status":"error"}`))
	}))
	defer srv.Close()

	_, err := New(srv.URL, "bad", nil).Health(context.Background())
	if err == nil {
		t.Fatal("expected error on 401, got nil")
	}
	if StatusCode(err) != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", StatusCode(err))
	}
}

// TestMetricNamesFiltersAndLimits exercises the case-insensitive substring filter
// and the limit cap.
func TestMetricNamesFiltersAndLimits(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"status":"success","data":["up","go_goroutines","go_memstats_alloc","node_cpu_seconds_total"]}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)

	all, err := c.MetricNames(context.Background(), "", 0)
	if err != nil {
		t.Fatalf("MetricNames: %v", err)
	}
	if gotPath != "/api/v1/label/__name__/values" {
		t.Errorf("path = %q", gotPath)
	}
	if len(all) != 4 {
		t.Fatalf("unfiltered = %v, want 4", all)
	}

	// Case-insensitive substring filter on "GO".
	filtered, err := c.MetricNames(context.Background(), "GO", 0)
	if err != nil {
		t.Fatalf("MetricNames: %v", err)
	}
	if !reflect.DeepEqual(filtered, []string{"go_goroutines", "go_memstats_alloc"}) {
		t.Errorf("filtered = %v, want [go_goroutines go_memstats_alloc]", filtered)
	}

	// Limit caps the result.
	limited, err := c.MetricNames(context.Background(), "go", 1)
	if err != nil {
		t.Fatalf("MetricNames: %v", err)
	}
	if len(limited) != 1 || limited[0] != "go_goroutines" {
		t.Errorf("limited = %v, want [go_goroutines]", limited)
	}
}

func TestLabelNamesDecodesData(t *testing.T) {
	var gotMatch string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMatch = r.URL.Query().Get("match[]")
		_, _ = w.Write([]byte(`{"status":"success","data":["__name__","instance","job"]}`))
	}))
	defer srv.Close()

	names, err := New(srv.URL, "", nil).LabelNames(context.Background(), "up")
	if err != nil {
		t.Fatalf("LabelNames: %v", err)
	}
	if gotMatch != "up" {
		t.Errorf("match[] = %q, want up", gotMatch)
	}
	if !reflect.DeepEqual(names, []string{"__name__", "instance", "job"}) {
		t.Errorf("names = %v", names)
	}
}

func TestLabelValuesDecodesData(t *testing.T) {
	var gotPath, gotMatch string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMatch = r.URL.Query().Get("match[]")
		_, _ = w.Write([]byte(`{"status":"success","data":["node","prometheus"]}`))
	}))
	defer srv.Close()

	values, err := New(srv.URL, "", nil).LabelValues(context.Background(), "job", "up")
	if err != nil {
		t.Fatalf("LabelValues: %v", err)
	}
	if gotPath != "/api/v1/label/job/values" {
		t.Errorf("path = %q", gotPath)
	}
	if gotMatch != "up" {
		t.Errorf("match[] = %q, want up", gotMatch)
	}
	if !reflect.DeepEqual(values, []string{"node", "prometheus"}) {
		t.Errorf("values = %v", values)
	}
}

func TestQueryRangePassthrough(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	var gotPath string
	var gotParams map[string]string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		q := r.URL.Query()
		gotParams = map[string]string{
			"query": q.Get("query"),
			"start": q.Get("start"),
			"end":   q.Get("end"),
			"step":  q.Get("step"),
		}
		_, _ = w.Write([]byte(promResp))
	}))
	defer srv.Close()

	body, err := New(srv.URL, "admin:pass", nil).QueryRange(context.Background(), "up", "1", "2", "15")
	if err != nil {
		t.Fatalf("QueryRange: %v", err)
	}
	if string(body) != promResp {
		t.Errorf("body not proxied verbatim: %s", body)
	}
	if gotPath != "/api/v1/query_range" {
		t.Errorf("path = %q", gotPath)
	}
	want := map[string]string{"query": "up", "start": "1", "end": "2", "step": "15"}
	if !reflect.DeepEqual(gotParams, want) {
		t.Errorf("params = %v, want %v", gotParams, want)
	}
}

func TestQueryPassthrough(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"vector","result":[]}}`
	var gotPath, gotTime string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotTime = r.URL.Query().Get("time")
		_, _ = w.Write([]byte(promResp))
	}))
	defer srv.Close()

	body, err := New(srv.URL, "", nil).Query(context.Background(), "up", "12345")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if string(body) != promResp {
		t.Errorf("body = %s", body)
	}
	if gotPath != "/api/v1/query" {
		t.Errorf("path = %q", gotPath)
	}
	if gotTime != "12345" {
		t.Errorf("time = %q, want 12345", gotTime)
	}
}

// TestQueryRangeBatch exercises the fan-out path end to end: valid queries return
// their bodies, an erroring query yields a per-result error while the others
// succeed, and every input query produces exactly one result (concurrency drops
// nothing).
func TestQueryRangeBatch(t *testing.T) {
	const promResp = `{"status":"success","data":{"resultType":"matrix","result":[]}}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The "bad" query is rejected by Prometheus with a 400.
		if r.URL.Query().Get("query") == "bad" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"status":"error","errorType":"bad_data","error":"parse error"}`))
			return
		}
		_, _ = w.Write([]byte(promResp))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	queries := []BatchQuery{
		{ID: "1", Query: "up"},
		{ID: "2", Query: "bad"},
		{ID: "3", Query: "rate(node_cpu_seconds_total[5m])"},
	}
	results := c.QueryRangeBatch(context.Background(), "1", "2", "15", queries)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	byID := map[string]BatchResult{}
	for _, r := range results {
		byID[r.ID] = r
	}
	if r := byID["1"]; r.Err != nil || string(r.Body) != promResp {
		t.Errorf("1 = %+v, want success with prom body", r)
	}
	if r := byID["3"]; r.Err != nil || string(r.Body) != promResp {
		t.Errorf("3 = %+v, want success with prom body", r)
	}
	if r := byID["2"]; r.Err == nil {
		t.Errorf("2 expected a per-result error, got nil (body=%s)", r.Body)
	}
}

func TestQueryRangeBatchEmpty(t *testing.T) {
	c := New("http://example.invalid", "", nil)
	if got := c.QueryRangeBatch(context.Background(), "1", "2", "15", nil); len(got) != 0 {
		t.Errorf("empty batch = %v, want no results", got)
	}
}

// TestQueryRangeBatchConcurrencyKeepsAllResults fans out more queries than the
// concurrency cap and asserts every query produces exactly one ordered result.
func TestQueryRangeBatchConcurrencyKeepsAllResults(t *testing.T) {
	var mu sync.Mutex
	seen := map[string]bool{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen[r.URL.Query().Get("query")] = true
		mu.Unlock()
		_, _ = w.Write([]byte(`{"status":"success","data":{"resultType":"matrix","result":[]}}`))
	}))
	defer srv.Close()

	c := New(srv.URL, "", nil)
	const n = 25 // > batchConcurrency (8)
	queries := make([]BatchQuery, n)
	for i := 0; i < n; i++ {
		id := "q" + itoa(i)
		queries[i] = BatchQuery{ID: id, Query: id}
	}
	results := c.QueryRangeBatch(context.Background(), "1", "2", "15", queries)
	if len(results) != n {
		t.Fatalf("expected %d results, got %d", n, len(results))
	}
	gotIDs := make([]string, len(results))
	for i, r := range results {
		if r.Err != nil {
			t.Errorf("result %d errored: %v", i, r.Err)
		}
		gotIDs[i] = r.ID
	}
	wantIDs := make([]string, n)
	for i := range queries {
		wantIDs[i] = queries[i].ID
	}
	if !reflect.DeepEqual(gotIDs, wantIDs) {
		t.Errorf("result order mismatch: got %v, want %v", gotIDs, wantIDs)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(seen) != n {
		uniq := make([]string, 0, len(seen))
		for k := range seen {
			uniq = append(uniq, k)
		}
		sort.Strings(uniq)
		t.Errorf("expected %d distinct queries reaching server, got %d: %v", n, len(seen), uniq)
	}
}

func TestMetadataPassthrough(t *testing.T) {
	const resp = `{"status":"success","data":{"up":[{"type":"gauge","help":"up","unit":""}]}}`
	var gotMetric string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMetric = r.URL.Query().Get("metric")
		_, _ = w.Write([]byte(resp))
	}))
	defer srv.Close()

	body, err := New(srv.URL, "", nil).Metadata(context.Background(), "up")
	if err != nil {
		t.Fatalf("Metadata: %v", err)
	}
	if string(body) != resp {
		t.Errorf("body = %s", body)
	}
	if gotMetric != "up" {
		t.Errorf("metric = %q, want up", gotMetric)
	}
}

// itoa is a tiny dependency-free int->string helper for test query IDs.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}
