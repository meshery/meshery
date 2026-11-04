package utils

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
)

const (
	// chartIndexTTL is how long a fetched repository index is reused. Charts
	// are published at release cadence, so minutes of staleness are harmless,
	// while re-downloading a multi-megabyte index.yaml on every cluster
	// connection is not.
	chartIndexTTL = 15 * time.Minute

	// chartIndexTimeout bounds the index fetch. Without it an unresponsive
	// repository stalls every caller indefinitely.
	chartIndexTimeout = 30 * time.Second
)

// chartIndexCache memoizes Helm repository indexes, keyed by repository URL.
//
// Only successful fetches are cached: a transient repository outage is retried
// on the next call rather than being pinned in for the whole TTL. The fetch
// itself runs without the lock held and is single-flighted per repository, so
// an unreachable repository costs all concurrent callers one chartIndexTimeout
// between them rather than one each, serially.
//
// Expiry is stale-while-revalidate, not a miss. This index is read on the
// cluster-connect request path, and index.yaml is several megabytes: treating
// an expired entry as absent puts a full chartIndexTimeout in front of a
// connect every time the TTL happens to lapse against a slow or unreachable
// repository. Published charts are append-mostly, so the worst a stale
// catalogue costs is not yet knowing about the newest chart - a case the
// version resolver already handles - which is strictly better than stalling the
// request. Only a cold cache, with no catalogue to serve at all, blocks.
type chartIndexCache struct {
	mu      sync.Mutex
	entries map[string]chartIndexCacheEntry

	// inflight tracks the fetch currently in progress for a repository, so
	// callers that miss the cache at the same time share its result instead of
	// each issuing their own request.
	inflight map[string]*chartIndexFetch

	// now and fetch are injectable so tests can drive expiry and serve a fixed
	// index without touching the network.
	now   func() time.Time
	fetch func(repo string) (map[string][]string, error)
}

type chartIndexCacheEntry struct {
	versions  map[string][]string
	fetchedAt time.Time
}

// chartIndexFetch is one in-progress fetch. done is closed once versions and
// err are final, and neither field is read before that.
type chartIndexFetch struct {
	done     chan struct{}
	versions map[string][]string
	err      error
}

func newChartIndexCache() *chartIndexCache {
	return &chartIndexCache{
		entries:  map[string]chartIndexCacheEntry{},
		inflight: map[string]*chartIndexFetch{},
		now:      time.Now,
		fetch:    fetchChartIndexVersions,
	}
}

var defaultChartIndexCache = newChartIndexCache()

// PublishedChartVersions returns every version of chart published in the Helm
// repository at repo, in the order that repository's index.yaml lists them
// (newest first by convention). The slice is a fresh copy the caller owns and
// may sort or otherwise modify; the cached catalogue behind it is shared across
// goroutines and is never handed out.
//
// A reachable repository that publishes no such chart yields an empty slice and
// no error: "the repository has nothing under that name" is a different
// condition from "the repository could not be read", and callers act on them
// differently.
func PublishedChartVersions(repo, chart string) ([]string, error) {
	return defaultChartIndexCache.list(repo, chart)
}

func (c *chartIndexCache) list(repo, chart string) ([]string, error) {
	versions, err := c.index(repo)
	if err != nil {
		return nil, err
	}
	return slices.Clone(versions[chart]), nil
}

// index returns repo's chart catalogue. A fresh entry is served as-is, an
// expired one is served while a refresh runs behind the caller, and only a
// caller that finds no entry at all waits for a fetch. The lock is held around
// the cache and the in-flight bookkeeping, never across the fetch.
func (c *chartIndexCache) index(repo string) (map[string][]string, error) {
	c.mu.Lock()
	entry, cached := c.entries[repo]
	if cached && c.now().Sub(entry.fetchedAt) <= chartIndexTTL {
		c.mu.Unlock()
		return entry.versions, nil
	}
	pending := c.startFetchLocked(repo)
	c.mu.Unlock()

	if cached {
		return entry.versions, nil
	}
	<-pending.done
	return pending.versions, pending.err
}

// startFetchLocked returns the fetch in progress for repo, starting one when
// none is. c.mu must be held.
//
// The fetch runs on its own goroutine whether or not anyone waits for it, which
// is what lets an expired entry be refreshed without a caller paying for it,
// and keeps every caller - waiting or not - behind the same single flight.
func (c *chartIndexCache) startFetchLocked(repo string) *chartIndexFetch {
	if pending, ok := c.inflight[repo]; ok {
		return pending
	}
	pending := &chartIndexFetch{done: make(chan struct{})}
	if c.inflight == nil {
		c.inflight = map[string]*chartIndexFetch{}
	}
	c.inflight[repo] = pending
	go func() {
		_, _ = c.runFetch(repo, pending)
	}()
	return pending
}

// runFetch performs the one fetch that waiters on pending are blocked behind,
// and publishes its outcome.
//
// Releasing the in-flight slot and closing pending.done are deferred because
// they must happen on every exit path. A panic that skipped them would wedge
// this repository permanently: the slot would still be occupied, so every later
// caller would block forever on a done channel nobody can close, leaking a
// goroutine each until the server is restarted. It would also take the process
// with it: this runs on its own goroutine, where an unrecovered panic is fatal
// with no handler above it to absorb it. A panic is therefore converted into a
// structured failure, which the normal error path already handles (installation
// is withheld, the cause is recorded for diagnostics, and the next call retries
// because failures are not cached).
func (c *chartIndexCache) runFetch(repo string, pending *chartIndexFetch) (versions map[string][]string, err error) {
	defer func() {
		if r := recover(); r != nil {
			versions = nil
			err = ErrHelmChartIndex(chartIndexURL(repo), fmt.Sprintf("reading the repository index panicked: %v", r))
		}

		c.mu.Lock()
		pending.versions, pending.err = versions, err
		delete(c.inflight, repo)
		if err == nil {
			if c.entries == nil {
				c.entries = map[string]chartIndexCacheEntry{}
			}
			c.entries[repo] = chartIndexCacheEntry{versions: versions, fetchedAt: c.now()}
		}
		c.mu.Unlock()
		close(pending.done)
	}()

	return c.fetch(repo)
}

// chartIndexURL is where a Helm repository serves its index.
func chartIndexURL(repo string) string {
	return strings.TrimSuffix(repo, "/") + "/index.yaml"
}

// fetchChartIndexVersions downloads repo's index.yaml and reduces it to
// chart name -> published versions. It decodes into MeshKit's HelmIndex rather
// than a local copy of the same shape, so the index contract keeps one owner.
func fetchChartIndexVersions(repo string) (map[string][]string, error) {
	url := chartIndexURL(repo)

	// The repository URL is server configuration, never user input. #nosec G107
	client := &http.Client{Timeout: chartIndexTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, ErrHelmChartIndex(url, err.Error())
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, ErrHelmChartIndex(url, "the repository responded with HTTP "+strconv.Itoa(resp.StatusCode))
	}

	var index mesherykube.HelmIndex
	if err := yaml.NewDecoder(resp.Body).Decode(&index); err != nil {
		return nil, ErrHelmChartIndex(url, err.Error())
	}

	versions := make(map[string][]string, len(index.Entries))
	for name, entries := range index.Entries {
		for _, entry := range entries {
			if entry.Version != "" {
				versions[name] = append(versions[name], entry.Version)
			}
		}
	}
	return versions, nil
}
