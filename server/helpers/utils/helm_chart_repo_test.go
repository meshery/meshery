package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

const testIndexYAML = `apiVersion: v1
entries:
  meshery-operator:
    - apiVersion: v2
      appVersion: v1.0.64
      name: meshery-operator
      version: v1.0.64
    - apiVersion: v2
      appVersion: v1.0.63
      name: meshery-operator
      version: v1.0.63
  meshery:
    - apiVersion: v2
      name: meshery
      version: v1.0.64
`

func TestFetchChartIndexVersionsReadsEveryChart(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/index.yaml" {
			t.Errorf("requested %q, want /index.yaml", r.URL.Path)
		}
		_, _ = w.Write([]byte(testIndexYAML))
	}))
	defer srv.Close()

	versions, err := fetchChartIndexVersions(srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := versions["meshery-operator"]; len(got) != 2 || got[0] != "v1.0.64" || got[1] != "v1.0.63" {
		t.Fatalf("meshery-operator versions = %v, want the index order [v1.0.64 v1.0.63]", got)
	}
	if got := versions["meshery"]; len(got) != 1 {
		t.Fatalf("meshery versions = %v, want one entry", got)
	}
}

// TestFetchChartIndexVersionsToleratesATrailingSlash guards the URL join: a
// repository configured with a trailing slash must not produce
// ".../charts//index.yaml".
func TestFetchChartIndexVersionsToleratesATrailingSlash(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/index.yaml" {
			t.Errorf("requested %q, want /index.yaml", r.URL.Path)
		}
		_, _ = w.Write([]byte(testIndexYAML))
	}))
	defer srv.Close()

	if _, err := fetchChartIndexVersions(srv.URL + "/"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFetchChartIndexVersionsReportsStructuredFailures(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{
			name:    "an error status",
			handler: func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNotFound) },
		},
		{
			name:    "a body that is not an index",
			handler: func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("\t- not: [yaml")) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(tt.handler)
			defer srv.Close()

			_, err := fetchChartIndexVersions(srv.URL)
			if err == nil {
				t.Fatal("expected a failure")
			}
			if code := meshkiterrors.GetCode(err); code != ErrHelmChartIndexCode {
				t.Fatalf("error code = %q, want %q (from %v)", code, ErrHelmChartIndexCode, err)
			}
		})
	}
}

// TestChartIndexCacheServesFromCacheUntilTheTTLExpires keeps the index off the
// wire on the hot path: it is fetched on every cluster connection, and the
// Meshery index is several megabytes.
func TestChartIndexCacheServesFromCacheUntilTheTTLExpires(t *testing.T) {
	now := time.Unix(0, 0)
	fetches := 0
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     func() time.Time { return now },
		fetch: func(string) (map[string][]string, error) {
			fetches++
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	for range 3 {
		if _, err := cache.list("repo", "meshery-operator"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if fetches != 1 {
		t.Fatalf("fetched %d times within the TTL, want 1", fetches)
	}

	now = now.Add(chartIndexTTL + time.Second)
	if _, err := cache.list("repo", "meshery-operator"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fetches != 2 {
		t.Fatalf("fetched %d times after the TTL expired, want 2", fetches)
	}
}

// TestChartIndexCacheDoesNotCacheFailures asserts a transient repository outage
// is retried on the next connection instead of being pinned in for a full TTL,
// which would extend a momentary blip into fifteen minutes of failed operator
// deployments.
func TestChartIndexCacheDoesNotCacheFailures(t *testing.T) {
	now := time.Unix(0, 0)
	calls := 0
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     func() time.Time { return now },
		fetch: func(string) (map[string][]string, error) {
			calls++
			if calls == 1 {
				return nil, errors.New("transient outage")
			}
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	if _, err := cache.list("repo", "meshery-operator"); err == nil {
		t.Fatal("expected the first call to report the outage")
	}
	versions, err := cache.list("repo", "meshery-operator")
	if err != nil {
		t.Fatalf("the retry must not be served from a cached failure: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("versions = %v, want the successfully fetched catalogue", versions)
	}
}

// TestChartIndexCacheDistinguishesEmptyFromUnreadable: a repository that is
// reachable but publishes no such chart is not an error. Callers decide what an
// empty catalogue means; only an unreadable index is a failure.
func TestChartIndexCacheDistinguishesEmptyFromUnreadable(t *testing.T) {
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch: func(string) (map[string][]string, error) {
			return map[string][]string{"meshery": {"v1.0.64"}}, nil
		},
	}

	versions, err := cache.list("repo", "meshery-operator")
	if err != nil {
		t.Fatalf("an absent chart is not a read failure: %v", err)
	}
	if len(versions) != 0 {
		t.Fatalf("versions = %v, want none", versions)
	}
}
