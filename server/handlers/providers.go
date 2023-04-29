package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
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
			cookie := &http.Cookie{
				Name:     h.config.ProviderCookieName,
				Value:    p.Name(),
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)

			redirectURL := "/user/login"
			if provider == "None" {
				redirectURL = "/"
			}

			http.Redirect(w, r, redirectURL, http.StatusFound)
			return
		}
	}

	// If the provider was not found, return a 404 error.
	http.NotFound(w, r)
}

// swagger:route GET /api/providers ProvidersAPI idGetProvidersList
// Handle GET request for list of providers
//
// Returns the available list of providers
// responses:
// 	200: listProvidersRespWrapper

// ProvidersHandler returns a list of providers
func (h *Handler) ProvidersHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests.
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get a map of providers and their properties.
	providers := make(map[string]models.ProviderProperties)
	for _, p := range h.config.Providers {
		providers[p.Name()] = p.GetProviderProperties()
	}

	// Marshal the providers map to JSON.
	b, err := json.Marshal(providers)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		h.log.Errorf("failed to marshal JSON: %v", err)
		return
	}

	// Write the JSON response.
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		h.log.Errorf("failed to write response: %v", err)
		return
	}
}

// swagger:route GET /provider ProvidersAPI idProvider
// Handle GET request to provider UI
//
// Servers providers UI
// responses:
// 	200:

// ProviderStaticHandler - servers Next.js UI
func (h *Handler) ProviderStaticHandler(w http.ResponseWriter, r *http.Request) {
	if h.config.PlaygroundBuild || h.Provider == "Meshery" { //Always use Remote provider for Playground build or when Provider is enforced
		http.SetCookie(w, &http.Cookie{
			Name:     h.config.ProviderCookieName,
			Value:    "Meshery",
			Path:     "/",
			HttpOnly: true,
		})
		redirectURL := "/user/login"
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
	}

	// Serve the application using the Next.js development server in live hot reload mode
	proxyURL, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	proxy.ServeHTTP(w, r)
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
	_ *models.User,
	provider models.Provider,
) {
	provider.GetProviderCapabilities(w, r)
}

const (
	uiReqBasePath     = "/api/provider/extension"
	serverReqBasePath = "/api/provider/extension/server/"
	loadReqBasePath   = "/api/provider/extension/"
)

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
	switch {
	case strings.HasPrefix(r.URL.Path, serverReqBasePath):
		h.ExtensionsEndpointHandler(w, r, prefObj, user, provider)
	case r.URL.Path == loadReqBasePath:
		err := h.LoadExtensionFromPackage(w, r, provider)
		if err != nil {
			// failed to load extensions from package
			//h.log.Error(ErrFailToLoadExtensions(err))
			http.Error(w, ErrFailToLoadExtensions(err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte("{}"))
	default:
		ServeReactComponentFromPackage(w, r, uiReqBasePath, provider)
	}
}
