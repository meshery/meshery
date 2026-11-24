package handlers

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/meshery/meshery/server/models"
)

// TestProvidersStreamHandler_RefreshOutlivesTheSubscriberThatStartedIt opens
// two provider streams. The first one starts the availability refresh and
// disconnects while the remote is still being probed; the second one
// connected during the probe, so it relies on that same refresh. The remote
// then answers. The second stream must see the remote online, and the tracker
// must not have cached it as offline.
func TestProvidersStreamHandler_RefreshOutlivesTheSubscriberThatStartedIt(t *testing.T) {
	probed := make(chan struct{})
	release := make(chan struct{})
	var probedOnce, releaseOnce sync.Once
	releaseRemote := func() { releaseOnce.Do(func() { close(release) }) }
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		probedOnce.Do(func() { close(probed) })
		select {
		case <-release:
		case <-r.Context().Done():
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"providerName":"cloud","providerType":"remote","capabilities":[]}`))
	}))
	t.Cleanup(remote.Close)
	t.Cleanup(releaseRemote)

	log := newTestLogger(t)
	cloud := &models.RemoteProvider{RemoteProviderURL: remote.URL, Log: log}
	cloud.Initialize()
	tracker := models.NewProviderTracker(map[string]models.Provider{"cloud": cloud}, log)
	h := &Handler{config: &models.HandlerConfig{ProviderTracker: tracker}, log: log}

	firstDone := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ProvidersStreamHandler(w, r)
		if r.URL.Query().Get("client") == "first" {
			close(firstDone)
		}
	}))
	t.Cleanup(srv.Close)

	// The first subscriber starts the refresh.
	firstCtx, cancelFirst := context.WithCancel(context.Background())
	defer cancelFirst()
	first := openProviderStream(t, firstCtx, srv.URL+"?client=first")
	defer func() { _ = first.Body.Close() }()
	select {
	case <-probed:
	case <-time.After(testTimeout):
		t.Fatal("the first subscriber did not start a refresh")
	}

	// The second subscriber connects while the probe is in flight. Its own
	// refresh is coalesced away, so it waits on the first one's.
	second := openProviderStream(t, context.Background(), srv.URL+"?client=second")
	defer func() { _ = second.Body.Close() }()
	events := readProviderEvents(t, second)
	if evt := nextProviderEvent(t, events); evt.Status != models.ProviderStatusChecking {
		t.Fatalf("second subscriber's snapshot: status %q, want %q", evt.Status, models.ProviderStatusChecking)
	}

	// The first subscriber goes away mid-probe, then the remote answers.
	cancelFirst()
	select {
	case <-firstDone:
	case <-time.After(testTimeout):
		t.Fatal("the first subscriber's handler did not return after it disconnected")
	}
	releaseRemote()

	evt := nextProviderEvent(t, events)
	if evt.Status != models.ProviderStatusOnline {
		t.Fatalf("second subscriber got status %q (error %q) for the remote, want %q", evt.Status, evt.Error, models.ProviderStatusOnline)
	}
	if got := tracker.Snapshot()["cloud"]; got.Status != models.ProviderStatusOnline {
		t.Fatalf("tracker cached status %q (error %q) for the remote, want %q", got.Status, got.Error, models.ProviderStatusOnline)
	}
}

func openProviderStream(t *testing.T, ctx context.Context, url string) *http.Response {
	t.Helper()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("building stream request: %v", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("opening provider stream: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("provider stream: status %d", resp.StatusCode)
	}
	return resp
}

// readProviderEvents decodes the stream's "data:" lines in the background.
func readProviderEvents(t *testing.T, resp *http.Response) <-chan models.ProviderStatusEvent {
	t.Helper()
	events := make(chan models.ProviderStatusEvent, 16)
	go func() {
		defer close(events)
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			data, ok := strings.CutPrefix(scanner.Text(), "data: ")
			if !ok {
				continue
			}
			var evt models.ProviderStatusEvent
			if err := json.Unmarshal([]byte(data), &evt); err != nil {
				t.Errorf("decoding provider event %q: %v", data, err)
				return
			}
			events <- evt
		}
	}()
	return events
}

func nextProviderEvent(t *testing.T, events <-chan models.ProviderStatusEvent) models.ProviderStatusEvent {
	t.Helper()
	select {
	case evt, ok := <-events:
		if !ok {
			t.Fatal("provider stream closed")
		}
		return evt
	case <-time.After(testTimeout):
		t.Fatal("no provider event")
	}
	return models.ProviderStatusEvent{}
}
