package utils

import (
	"net/http"
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
// on the next call rather than being pinned in for the whole TTL.
type chartIndexCache struct {
	mu      sync.Mutex
	entries map[string]chartIndexCacheEntry

	// now and fetch are injectable so tests can drive expiry and serve a fixed
	// index without touching the network.
	now   func() time.Time
	fetch func(repo string) (map[string][]string, error)
}

type chartIndexCacheEntry struct {
	versions  map[string][]string
	fetchedAt time.Time
}

func newChartIndexCache() *chartIndexCache {
	return &chartIndexCache{
		entries: map[string]chartIndexCacheEntry{},
		now:     time.Now,
		fetch:   fetchChartIndexVersions,
	}
}

var defaultChartIndexCache = newChartIndexCache()

// PublishedChartVersions returns every version of chart published in the Helm
// repository at repo, in the order that repository's index.yaml lists them
// (newest first by convention; callers that depend on ordering must sort).
//
// A reachable repository that publishes no such chart yields an empty slice and
// no error: "the repository has nothing under that name" is a different
// condition from "the repository could not be read", and callers act on them
// differently.
func PublishedChartVersions(repo, chart string) ([]string, error) {
	return defaultChartIndexCache.list(repo, chart)
}

func (c *chartIndexCache) list(repo, chart string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[repo]
	if !ok || c.now().Sub(entry.fetchedAt) > chartIndexTTL {
		versions, err := c.fetch(repo)
		if err != nil {
			return nil, err
		}
		entry = chartIndexCacheEntry{versions: versions, fetchedAt: c.now()}
		c.entries[repo] = entry
	}
	return entry.versions[chart], nil
}

// fetchChartIndexVersions downloads repo's index.yaml and reduces it to
// chart name -> published versions. It decodes into MeshKit's HelmIndex rather
// than a local copy of the same shape, so the index contract keeps one owner.
func fetchChartIndexVersions(repo string) (map[string][]string, error) {
	url := strings.TrimSuffix(repo, "/") + "/index.yaml"

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
