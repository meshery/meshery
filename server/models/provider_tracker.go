package models

import (
	"context"
	"sync"
	"time"

	"github.com/meshery/meshkit/logger"
)

// ProviderTracker owns the authoritative availability state for every
// registered Meshery provider (local + remotes) and broadcasts each
// transition to subscribers (the /api/providers/stream SSE handler).
//
// It exists to decouple three things that previously collapsed into a single
// blocking call chain (handler -> RemoteProvider.GetProviderProperties ->
// loadCapabilities):
//
//  1. /api/providers can now respond from cached state, so a single
//     unreachable remote can no longer freeze the provider chooser.
//  2. Availability probes run concurrently across providers; a single 3s
//     timeout caps the total wait regardless of how many remotes are
//     configured.
//  3. Each remote's status reaches the UI independently as its probe
//     resolves, instead of arriving as one all-or-nothing payload after
//     the slowest remote completes.
type ProviderTracker struct {
	log logger.Handler

	mu        sync.RWMutex
	providers map[string]Provider
	statuses  map[string]ProviderStatusEvent
	subs      map[chan ProviderStatusEvent]struct{}

	// refreshMu serializes VerifyAll runs so a burst of SSE subscriptions
	// cannot stampede remotes with overlapping capability probes. A caller
	// that arrives while a refresh is in flight returns immediately; the
	// in-flight refresh's publishes will still reach all subscribers.
	refreshMu sync.Mutex
}

// NewProviderTracker builds a tracker over the given registration map. Local
// providers start ProviderStatusOnline (no network dependency to verify);
// remote providers start ProviderStatusChecking so the UI renders an entry
// while the boot-time probe is still in flight.
//
// The tracker holds a reference to the same map main.go populates; it does
// not copy it. Callers must finish populating the map before constructing
// the tracker.
func NewProviderTracker(providers map[string]Provider, log logger.Handler) *ProviderTracker {
	t := &ProviderTracker{
		log:       log,
		providers: providers,
		statuses:  make(map[string]ProviderStatusEvent, len(providers)),
		subs:      map[chan ProviderStatusEvent]struct{}{},
	}
	for key, p := range providers {
		status := ProviderStatusOnline
		if p.GetProviderType() == RemoteProviderType {
			status = ProviderStatusChecking
		}
		t.statuses[key] = ProviderStatusEvent{
			Key:        key,
			Status:     status,
			Properties: p.GetProviderProperties(),
		}
	}
	return t
}

// Snapshot returns a copy of the current per-provider statuses keyed by
// registration key (the same key callers must put in the meshery-provider
// cookie to route to that provider). Used by /api/providers when an SSE
// client is not available and by /api/providers/stream to seed a new
// subscriber before live updates begin.
func (t *ProviderTracker) Snapshot() map[string]ProviderStatusEvent {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[string]ProviderStatusEvent, len(t.statuses))
	for k, v := range t.statuses {
		out[k] = v
	}
	return out
}

// Subscribe registers a new listener for status transitions. The returned
// channel receives the current snapshot first (one event per provider)
// followed by every future Publish. The returned unsubscribe must be called
// when the listener is done; calling it twice is safe.
//
// Sends to a subscriber that is not draining its channel are dropped
// (logged at debug) rather than blocking the publisher. A slow SSE client
// will simply miss intermediate updates and pick up state on its next
// reconnection.
func (t *ProviderTracker) Subscribe(ctx context.Context) (<-chan ProviderStatusEvent, func()) {
	ch := make(chan ProviderStatusEvent, 32)

	t.mu.Lock()
	t.subs[ch] = struct{}{}
	snapshot := make([]ProviderStatusEvent, 0, len(t.statuses))
	for _, evt := range t.statuses {
		snapshot = append(snapshot, evt)
	}
	t.mu.Unlock()

	go func() {
		for _, evt := range snapshot {
			select {
			case ch <- evt:
			case <-ctx.Done():
				return
			}
		}
	}()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			t.mu.Lock()
			defer t.mu.Unlock()
			if _, ok := t.subs[ch]; ok {
				delete(t.subs, ch)
				close(ch)
			}
		})
	}
	return ch, unsubscribe
}

// Publish records the new event in the status map and broadcasts it to every
// subscriber. Sends are non-blocking; a subscriber whose channel is full is
// skipped for this event so a single stalled SSE writer cannot stall the
// tracker for everyone else.
func (t *ProviderTracker) Publish(evt ProviderStatusEvent) {
	t.mu.Lock()
	t.statuses[evt.Key] = evt
	subs := make([]chan ProviderStatusEvent, 0, len(t.subs))
	for ch := range t.subs {
		subs = append(subs, ch)
	}
	t.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- evt:
		default:
			if t.log != nil {
				t.log.Debugf("provider-tracker: subscriber channel full; dropping event for %q", evt.Key)
			}
		}
	}
}

// VerifyAll fans the availability probe out across every registered remote
// provider, publishing each result independently as its probe completes.
// Local providers are reported online with no network call. The function
// returns once every probe has settled; pass a context to bound the total
// runtime (e.g. an http.Request context for an SSE-driven refresh).
//
// Concurrency: each remote is probed in its own goroutine, so total wall
// time is bounded by the slowest single probe rather than by their sum.
// refreshMu collapses overlapping VerifyAll calls: a refresh that arrives
// while another is in flight is a no-op (the in-flight refresh's publishes
// reach all current subscribers already).
func (t *ProviderTracker) VerifyAll(ctx context.Context) {
	if !t.refreshMu.TryLock() {
		return
	}
	defer t.refreshMu.Unlock()

	t.mu.RLock()
	type item struct {
		key string
		p   Provider
	}
	work := make([]item, 0, len(t.providers))
	for k, p := range t.providers {
		work = append(work, item{k, p})
	}
	t.mu.RUnlock()

	// Mark remotes that are about to be re-probed as "checking" so the UI
	// can show a refresh-in-progress state instead of stale "online"/"offline"
	// while the probe is running.
	for _, w := range work {
		if w.p.GetProviderType() == RemoteProviderType {
			t.Publish(ProviderStatusEvent{
				Key:        w.key,
				Status:     ProviderStatusChecking,
				Properties: w.p.GetProviderProperties(),
			})
		}
	}

	var wg sync.WaitGroup
	for _, w := range work {
		wg.Add(1)
		go func(key string, p Provider) {
			defer wg.Done()
			t.Publish(verifyOne(ctx, key, p))
		}(w.key, w.p)
	}
	wg.Wait()
}

// verifyOne performs the availability probe for a single provider and
// returns the resulting status event. Local providers are reported online
// without a network call; remote providers fall back to their cached
// properties (URL, fallback host name) when the probe fails so the chooser
// can still render a recognizable label in the offline section.
func verifyOne(ctx context.Context, key string, p Provider) ProviderStatusEvent {
	if p.GetProviderType() == LocalProviderType {
		return ProviderStatusEvent{
			Key:        key,
			Status:     ProviderStatusOnline,
			Properties: p.GetProviderProperties(),
		}
	}
	rp, ok := p.(*RemoteProvider)
	if !ok {
		return ProviderStatusEvent{
			Key:        key,
			Status:     ProviderStatusOnline,
			Properties: p.GetProviderProperties(),
		}
	}
	probeCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	props, err := rp.VerifyAvailability(probeCtx)
	if err != nil {
		return ProviderStatusEvent{
			Key:        key,
			Status:     ProviderStatusOffline,
			Properties: rp.GetProviderProperties(),
			Error:      err.Error(),
		}
	}
	return ProviderStatusEvent{
		Key:        key,
		Status:     ProviderStatusOnline,
		Properties: props,
	}
}
