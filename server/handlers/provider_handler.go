package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/meshery/meshery/server/models"
)

// ProviderHandler - handles the choice of provider
func (h *Handler) ProviderHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	for _, p := range h.config.Providers {
		if provider == p.Name() {
			http.SetCookie(w, &http.Cookie{
				Name:     h.config.ProviderCookieName,
				Value:    p.Name(),
				Path:     "/",
				HttpOnly: true,
			})
			redirectURL := "/user/login?" + r.URL.RawQuery
			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
	}
}

// ProvidersHandler returns a list of providers
func (h *Handler) ProvidersHandler(w http.ResponseWriter, _ *http.Request) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 	return
	// }

	providers := map[string]models.ProviderProperties{}
	for _, p := range h.config.Providers {
		providers[p.Name()] = (p.GetProviderProperties())
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
	if h.Provider != "" {
		if h.config.Providers[h.Provider] != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     h.config.ProviderCookieName,
				Value:    h.Provider,
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
