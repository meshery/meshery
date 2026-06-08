package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	models "github.com/meshery/meshery/server/models"
)

// ProviderHandler - handles the choice of provider
func (h *Handler) ProviderHandler(w http.ResponseWriter, r *http.Request) {
	providerKey, ok := models.ResolveProviderKey(r.URL.Query().Get("provider"), h.config.Providers)
	if !ok {
		// An unknown or unregistered provider (e.g. a stale UI sending a
		// name that no longer maps to a registered provider) would
		// otherwise fall through to a silent 200 with an empty body,
		// leaving the browser with no actionable signal. Send the user
		// back to the provider chooser so they can pick a valid one.
		http.Redirect(w, r, "/provider", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     h.config.ProviderCookieName,
		Value:    providerKey,
		Path:     "/",
		HttpOnly: true,
	})
	redirectURL := "/user/login?" + r.URL.RawQuery
	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// ProvidersHandler returns a list of providers
//
// The response is sourced from ProviderTracker's snapshot (the same state
// /api/providers/stream broadcasts) so this handler never blocks on a
// remote probe. Existing clients keep working: each map value is the same
// ProviderProperties shape they expect. Clients that want live updates as
// remotes settle should subscribe to /api/providers/stream instead of
// polling this endpoint.
func (h *Handler) ProvidersHandler(w http.ResponseWriter, _ *http.Request) {
	// Key the response by the registration map key - that is the value a
	// client must put in the meshery-provider cookie to route to this
	// provider. Using p.Name() here would collapse entries whenever two
	// remote URLs report the same providerName from /capabilities (their
	// shared name overwrites in the output even though the registration
	// map disambiguates them by URL host).
	providers := map[string]models.ProviderProperties{}
	if h.config.ProviderTracker != nil {
		for key, evt := range h.config.ProviderTracker.Snapshot() {
			providers[key] = evt.Properties
		}
	} else {
		// Fallback for tests / older wiring that did not build a
		// tracker. Reads cached state from each provider directly;
		// since GetProviderProperties no longer blocks on HTTP, this
		// is still bounded.
		for key, p := range h.config.Providers {
			providers[key] = p.GetProviderProperties()
		}
	}
	bd, err := json.Marshal(providers)
	if err != nil {
		obj := "provider"
		h.log.Error(models.ErrMarshal(err, obj))
		writeMeshkitError(w, models.ErrMarshal(err, obj), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bd)
}

// ProvidersStreamHandler streams provider availability over Server-Sent
// Events. Each event is a JSON ProviderStatusEvent: the registration key,
// status (checking|online|offline), current ProviderProperties, and an
// optional error message. The subscriber first receives one event per
// registered provider (current state), then continues to receive updates
// as the tracker re-probes remotes.
//
// On every subscription we also kick off a tracker.VerifyAll in the
// background, so opening the chooser refreshes status without the UI
// having to poll. Overlapping VerifyAll calls are collapsed by the
// tracker's refresh mutex - the SSE flood from many tabs cannot stampede
// the configured remotes.
func (h *Handler) ProvidersStreamHandler(w http.ResponseWriter, r *http.Request) {
	if h.config.ProviderTracker == nil {
		writeJSONError(w, "provider tracker not initialized", http.StatusInternalServerError)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeJSONError(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// X-Accel-Buffering disables buffering at any nginx hop in front of
	// Meshery so the events reach the browser immediately.
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	ctx := r.Context()
	events, unsubscribe := h.config.ProviderTracker.Subscribe(ctx)
	defer unsubscribe()

	// Refresh status on subscribe so a freshly-opened chooser shows live
	// state. The call is non-blocking from the SSE writer's perspective:
	// VerifyAll publishes to all subscribers as it runs, and overlapping
	// runs are coalesced by the tracker's refresh mutex.
	go h.config.ProviderTracker.VerifyAll(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case evt, ok := <-events:
			if !ok {
				return
			}
			data, err := json.Marshal(evt)
			if err != nil {
				h.log.Error(models.ErrMarshal(err, "provider status event"))
				continue
			}
			if _, err := fmt.Fprintf(w, "event: provider\ndata: %s\n\n", data); err != nil {
				// Client disconnected mid-write; ctx cancellation
				// will close the channel and we'll exit cleanly.
				return
			}
			flusher.Flush()
		}
	}
}

// ProviderUIHandler serves the provider-selection UI, or — when the deployment
// has a single provider it can auto-select — sets the meshery-provider cookie
// and bounces the request to /user/login so the user never sees the chooser.
//
// Auto-select is gated on h.Provider being non-empty AND registered in
// h.config.Providers. We must never set an empty meshery-provider cookie or
// one that names an unregistered provider: ProviderMiddleware ignores empty
// cookie values, and an unregistered name resolves to a nil Provider in the
// request context. Either way, AuthMiddleware (and the inline /user/login
// handler) will fail the type assertion and 302 the request right back to
// /provider — producing an infinite /user/login ⇄ /provider redirect loop.
//
// We hit exactly that loop with meshery extension images that are built
// with PLAYGROUND=true baked in, and their deployment was missing
// the PROVIDER env var. ProviderUIHandler entered the auto-select branch
// (because PlaygroundBuild was true), wrote `meshery-provider=` (empty
// because h.Provider was unset), and looped. PR #19006 already closed the
// matching path inside ProviderMiddleware (cookie-roundtrip fragility on
// SameSite/CDN), but it relied on h.Provider being set — so it didn't help
// the deployment-misconfig case at all. This handler is the actual source
// of the empty cookie, and the fix belongs here.
//
// In any "can't safely auto-select" branch we fall through to ServeUI so
// the operator at least gets the provider-chooser page (degraded UX, but
// reachable) and a clear log line pointing at the deployment misconfig.
func (h *Handler) ProviderUIHandler(w http.ResponseWriter, r *http.Request) {
	enforcedProvider := models.NormalizeProviderName(h.Provider)
	if enforcedProvider != "" {
		if providerKey, ok := models.ResolveProviderKey(enforcedProvider, h.config.Providers); ok {
			http.SetCookie(w, &http.Cookie{
				Name:     h.config.ProviderCookieName,
				Value:    providerKey,
				Path:     "/",
				HttpOnly: true,
			})
			// Propagate existing request parameters, if present.
			redirectURL := "/user/login"
			if r.URL.RawQuery != "" {
				redirectURL += "?" + r.URL.RawQuery
			}
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
		h.log.Errorf(
			"enforced provider %q is not registered in h.config.Providers; serving the provider-selection UI instead of looping. Register %q in PROVIDERS or unset PROVIDER on this deployment.",
			h.Provider, h.Provider,
		)
	} else if h.config.PlaygroundBuild {
		h.log.Errorf(
			"PLAYGROUND=true but PROVIDER env var is unset on this deployment; serving the provider-selection UI instead of looping between /user/login and /provider. Set PROVIDER (e.g. PROVIDER=Meshery) to restore one-click playground access.",
		)
	}
	h.ServeUI(w, r, "/provider", "../../provider-ui/out/")
}

// ProviderCapabilityHandler returns the capabilities.json for the provider
func (h *Handler) ProviderCapabilityHandler(
	w http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// change it to use fethc from the meshery server cache
	providerCapabilities, err := provider.ReadCapabilitiesForUser(user.ID.String())
	if err != nil {
		h.log.Debugf("User capabilities not found in server store for user_id: %s, trying to fetch capabilities from the remote provider", user.ID.String())
		provider.GetProviderCapabilities(w, r, user.ID.String())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(providerCapabilities)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "provider capabilities"))
		writeMeshkitError(w, models.ErrMarshal(err, "provider capabilities"), http.StatusInternalServerError)
		return
	}
}

// ProviderComponentsHandler handlers the requests to serve react
// components from the provider package
func (h *Handler) ProviderComponentsHandler(
	w http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	uiReqBasePath := "/api/provider/extension"
	serverReqBasePath := "/api/provider/extension/server/"
	loadReqBasePath := "/api/provider/extension/"

	if strings.HasPrefix(r.URL.Path, serverReqBasePath) {
		h.ExtensionsEndpointHandler(w, r, prefObj, user, provider)
	} else if r.URL.Path == loadReqBasePath {
		err := h.LoadExtensionFromPackage(w, r, provider)
		if err != nil {
			// failed to load extensions from package
			h.log.Error(ErrFailToLoadExtensions(err))
			writeMeshkitError(w, ErrFailToLoadExtensions(err), http.StatusInternalServerError)
			return
		}
		writeJSONEmptyObject(w, http.StatusOK)
	} else {
		ServeReactComponentFromPackage(w, r, uiReqBasePath, provider)
	}
}
