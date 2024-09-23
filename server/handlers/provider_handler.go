package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/layer5io/meshery/server/models"
)

// swagger:route GET /api/provider ProvidersAPI idChoiceProvider
// Handle GET request for the choice of provider
//
// Update the choice of provider in system
// responses:
// 	200:

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
			if provider == "None" {
				redirectURL = "/"
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
	}
}

// swagger:route GET /api/providers ProvidersAPI idGetProvidersList
// Handle GET request for list of providers
//
// Returns the available list of providers
// responses:
// 	200: listProvidersRespWrapper

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
		http.Error(w, models.ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bd)
}

// swagger:route GET /provider ProvidersAPI idProvider
// Handle GET request to provider UI
//
// Servers providers UI
// responses:
// 	200:

// ProviderUIHandler - serves providers UI
func (h *Handler) ProviderUIHandler(w http.ResponseWriter, r *http.Request) {
	if h.config.PlaygroundBuild || h.Provider != "" { //Always use Remote provider for Playground build or when Provider is enforced
		http.SetCookie(w, &http.Cookie{
			Name:     h.config.ProviderCookieName,
			Value:    h.Provider,
			Path:     "/",
			HttpOnly: true,
		})
		// Propagate existing request parameters, if present.
		redirectURL := "/user/login?" + r.URL.RawQuery
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}
	ServeUI(w, r, "/provider", "../../provider-ui/out/")
}

// swagger:route GET /api/provider/capabilities ProvidersAPI idGetProviderCapabilities
// Handle GET requests for Provider
//
// Returns the capabilities.json for the provider
// responses:
// 	200:

// ProviderCapabilityHandler returns the capabilities.json for the provider
func (h *Handler) ProviderCapabilityHandler(
	w http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	// change it to use fethc from the meshery server cache
	providerCapabilities, err := provider.ReadCapabilitiesForUser(user.ID)
	if err != nil {
		h.log.Debugf("User capabilities not found in server store for user_id: %s, trying to fetch capabilities from the remote provider", user.ID)
		provider.GetProviderCapabilities(w, r, user.ID)
		return
	}

	err = json.NewEncoder(w).Encode(providerCapabilities)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "provider capabilities"))
		http.Error(w, models.ErrMarshal(err, "provider capabilities").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/provider/extension ProvidersAPI idReactComponents
// Handle GET request for React Components
//
// handles the requests to serve react components from the provider package
// responses:
// 	200:

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
			http.Error(w, ErrFailToLoadExtensions(err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte("{}"))
	} else {
		ServeReactComponentFromPackage(w, r, uiReqBasePath, provider)
	}
}
