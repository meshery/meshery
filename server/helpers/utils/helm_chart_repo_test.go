package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
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
	var mu sync.Mutex
	now := time.Unix(0, 0)
	var fetches atomic.Int32
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now: func() time.Time {
			mu.Lock()
			defer mu.Unlock()
			return now
		},
		fetch: func(string) (map[string][]string, error) {
			fetches.Add(1)
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	for range 3 {
		if _, err := cache.list("repo", "meshery-operator"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if got := fetches.Load(); got != 1 {
		t.Fatalf("fetched %d times within the TTL, want 1", got)
	}

	mu.Lock()
	now = now.Add(chartIndexTTL + time.Second)
	mu.Unlock()
	if _, err := cache.list("repo", "meshery-operator"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	waitFor(t, "the expired entry to be refreshed", func() bool { return fetches.Load() == 2 })
}

// TestChartIndexCacheServesStaleWhileRefreshing is the reason expiry is not a
// miss. This index is read inside the cluster-connect request handler; blocking
// that request on a multi-megabyte fetch bounded only by chartIndexTimeout
// turned a slow or unreachable repository into a 30-second connect (and every
// concurrent connect joined the same wait) every time the TTL happened to
// lapse. Published charts are append-mostly, so a stale catalogue costs at most
// not knowing about the newest chart, which the resolver already handles.
func TestChartIndexCacheServesStaleWhileRefreshing(t *testing.T) {
	var mu sync.Mutex
	now := time.Unix(0, 0)
	release := make(chan struct{})
	var fetches atomic.Int32

	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now: func() time.Time {
			mu.Lock()
			defer mu.Unlock()
			return now
		},
		fetch: func(string) (map[string][]string, error) {
			if fetches.Add(1) > 1 {
				<-release
				return map[string][]string{"meshery-operator": {"v1.0.65", "v1.0.64"}}, nil
			}
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	// Warm the cache, then expire it. The refresh that the next call kicks off
	// is blocked for the rest of the test.
	if _, err := cache.list("repo", "meshery-operator"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mu.Lock()
	now = now.Add(chartIndexTTL + time.Second)
	mu.Unlock()

	served := make(chan []string, 1)
	go func() {
		versions, err := cache.list("repo", "meshery-operator")
		if err != nil {
			t.Errorf("serving a stale catalogue must not fail: %v", err)
		}
		served <- versions
	}()

	select {
	case versions := <-served:
		if len(versions) != 1 || versions[0] != "v1.0.64" {
			t.Fatalf("versions = %v, want the stale catalogue served while the refresh runs", versions)
		}
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("an expired entry blocked the caller on a refetch instead of being served stale")
	}

	// The refresh really was started, and its result replaces the stale entry.
	waitFor(t, "the background refresh to start", func() bool { return fetches.Load() == 2 })
	close(release)
	waitFor(t, "the refreshed catalogue to be installed", func() bool {
		versions, err := cache.list("repo", "meshery-operator")
		return err == nil && len(versions) == 2
	})
}

// TestChartIndexCacheBlocksOnlyOnAColdCache: with no catalogue to serve there
// is nothing to be stale about, so a first caller does wait - and gets the
// failure rather than an empty slice that would read as "no chart published".
func TestChartIndexCacheBlocksOnlyOnAColdCache(t *testing.T) {
	outage := errors.New("repository unreachable")
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch:   func(string) (map[string][]string, error) { return nil, outage },
	}

	versions, err := cache.list("repo", "meshery-operator")
	if !errors.Is(err, outage) {
		t.Fatalf("err = %v, want the fetch failure reported to the first caller", err)
	}
	if versions != nil {
		t.Fatalf("versions = %v, want none: an unreadable index is not an empty catalogue", versions)
	}
}

// TestChartIndexCacheServesStaleThroughAnOutage: failures are not cached, so a
// refresh that keeps failing must neither poison the stale entry nor start
// blocking callers again.
func TestChartIndexCacheServesStaleThroughAnOutage(t *testing.T) {
	var mu sync.Mutex
	now := time.Unix(0, 0)
	var fetches atomic.Int32

	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now: func() time.Time {
			mu.Lock()
			defer mu.Unlock()
			return now
		},
		fetch: func(string) (map[string][]string, error) {
			if fetches.Add(1) > 1 {
				return nil, errors.New("repository unreachable")
			}
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	if _, err := cache.list("repo", "meshery-operator"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mu.Lock()
	now = now.Add(chartIndexTTL + time.Second)
	mu.Unlock()

	for range 3 {
		versions, err := cache.list("repo", "meshery-operator")
		if err != nil {
			t.Fatalf("a failing refresh must not fail the caller that could be served stale: %v", err)
		}
		if len(versions) != 1 || versions[0] != "v1.0.64" {
			t.Fatalf("versions = %v, want the stale catalogue", versions)
		}
		waitFor(t, "the failing refresh to finish", func() bool {
			cache.mu.Lock()
			defer cache.mu.Unlock()
			return len(cache.inflight) == 0
		})
	}
	if got := fetches.Load(); got < 2 {
		t.Fatalf("fetched %d times, want the refresh retried rather than a failure cached", got)
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

// waitFor spins until cond holds, failing the test if it has not within a
// generous deadline. Used instead of a fixed sleep so the tests below do not
// encode a timing guess.
func waitFor(t *testing.T, what string, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(time.Millisecond)
	}
}

// TestChartIndexCacheDoesNotHoldTheLockAcrossAFetch is the reason the fetch
// moved out from under the mutex. The fetch is an HTTP GET bounded only by
// chartIndexTimeout, so holding the lock across it stalled every other caller -
// including callers for an entirely different repository - for the full
// timeout, one after another. On an egress-blocked server that turned several
// simultaneous cluster connections into minutes of blocked request handlers.
func TestChartIndexCacheDoesNotHoldTheLockAcrossAFetch(t *testing.T) {
	inFlight := make(chan struct{})
	release := make(chan struct{})

	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch: func(repo string) (map[string][]string, error) {
			if repo == "slow" {
				close(inFlight)
				<-release
			}
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	slowDone := make(chan struct{})
	go func() {
		defer close(slowDone)
		if _, err := cache.list("slow", "meshery-operator"); err != nil {
			t.Errorf("slow repository: unexpected error: %v", err)
		}
	}()
	<-inFlight

	fastDone := make(chan struct{})
	go func() {
		defer close(fastDone)
		if _, err := cache.list("fast", "meshery-operator"); err != nil {
			t.Errorf("fast repository: unexpected error: %v", err)
		}
	}()

	select {
	case <-fastDone:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("a lookup for another repository blocked behind an in-flight fetch: the cache lock is held across the fetch")
	}

	close(release)
	<-slowDone
}

// TestChartIndexCacheSingleFlightsConcurrentMisses asserts concurrent callers
// that miss the cache share one fetch instead of each issuing their own.
//
// The fetch here fails, which is the case that matters: failures are
// deliberately not cached, so without single-flight every caller pays its own
// full chartIndexTimeout - N cluster connections against an unreachable
// repository cost N timeouts rather than one.
func TestChartIndexCacheSingleFlightsConcurrentMisses(t *testing.T) {
	const callers = 8

	release := make(chan struct{})
	var fetches atomic.Int32
	outage := errors.New("repository unreachable")

	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch: func(string) (map[string][]string, error) {
			fetches.Add(1)
			<-release
			return nil, outage
		},
	}

	var wg sync.WaitGroup
	errs := make([]error, callers)
	for i := range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, errs[i] = cache.list("repo", "meshery-operator")
		}()
	}

	// The first caller is now inside the blocked fetch; the rest are joining it.
	waitFor(t, "the first fetch to start", func() bool { return fetches.Load() >= 1 })
	time.Sleep(200 * time.Millisecond)
	close(release)
	wg.Wait()

	// Serializing on the mutex instead of sharing the result costs exactly one
	// fetch per caller.
	if got := fetches.Load(); got >= callers {
		t.Fatalf("fetched %d times for %d concurrent callers: each caller paid its own timeout", got, callers)
	}
	for i := range callers {
		if !errors.Is(errs[i], outage) {
			t.Fatalf("caller %d: err = %v, want the shared fetch failure", i, errs[i])
		}
	}
}

// TestChartIndexCacheSurvivesAPanickingFetch: the in-flight slot must be
// released on every exit path. Leaving it behind after a panic wedges the
// repository permanently - the done channel nobody can close blocks every
// later caller forever, leaking a goroutine each - and net/http recovers the
// panic in the handler goroutine, so the process stays up to keep hanging.
func TestChartIndexCacheSurvivesAPanickingFetch(t *testing.T) {
	calls := 0
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch: func(string) (map[string][]string, error) {
			calls++
			if calls == 1 {
				panic("index decode blew up")
			}
			return map[string][]string{"meshery-operator": {"v1.0.64"}}, nil
		},
	}

	// The panicking flight reports a structured failure rather than escaping as
	// a panic or handing back a nil catalogue with a nil error.
	versions, err := cache.list("repo", "meshery-operator")
	if err == nil {
		t.Fatalf("a panicking fetch must be reported as a failure, got versions %v", versions)
	}
	if code := meshkiterrors.GetCode(err); code != ErrHelmChartIndexCode {
		t.Fatalf("error code = %q, want %q (from %v)", code, ErrHelmChartIndexCode, err)
	}

	// And the cache is still usable: no slot left occupied, so this neither
	// blocks nor is served from a cached failure.
	done := make(chan struct{})
	go func() {
		defer close(done)
		versions, err := cache.list("repo", "meshery-operator")
		if err != nil {
			t.Errorf("the cache is wedged after a panicking fetch: %v", err)
			return
		}
		if len(versions) != 1 || versions[0] != "v1.0.64" {
			t.Errorf("versions = %v, want the successfully fetched catalogue", versions)
		}
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("a lookup after a panicking fetch blocked forever: the in-flight slot was never released")
	}
}

// TestChartIndexCacheReturnsACopyOfTheCatalogue: the cached catalogue is shared
// across goroutines for a full TTL. Handing out its own slice lets any caller
// that sorts or appends corrupt it for everyone else, with no lock held.
func TestChartIndexCacheReturnsACopyOfTheCatalogue(t *testing.T) {
	cache := &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch: func(string) (map[string][]string, error) {
			return map[string][]string{"meshery-operator": {"v1.0.64", "v1.0.63"}}, nil
		},
	}

	first, err := cache.list("repo", "meshery-operator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	first[0] = "tampered"

	second, err := cache.list("repo", "meshery-operator")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if second[0] != "v1.0.64" {
		t.Fatalf("versions = %v, want the cached catalogue to be unaffected by a caller mutating its result", second)
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
