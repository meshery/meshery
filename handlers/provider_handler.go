package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	models "github.com/layer5io/meshery/models"
)

// ProviderHandler - handles the choice of provider
func (h *Handler) ProviderHandler(w http.ResponseWriter, r *http.Request) {
	provider := r.URL.Query().Get("provider")
	for _, p := range h.config.Providers {
		if provider == p.Name() {
			http.SetCookie(w, &http.Cookie{
				Name:     h.config.ProviderCookieName,
				Value:    p.Name(),
				Expires:  time.Now().Add(h.config.ProviderCookieDuration),
				Path:     "/",
				HttpOnly: true,
			})
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	}
}

// ProvidersHandler returns a list of providers
func (h *Handler) ProvidersHandler(w http.ResponseWriter, r *http.Request) {
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
		h.log.Error(ErrMarshal(err, obj))
		http.Error(w, ErrMarshal(err, obj).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bd)
}

// ProviderUIHandler - serves providers UI
func (h *Handler) ProviderUIHandler(w http.ResponseWriter, r *http.Request) {
	ServeUI(w, r, "/provider", "../provider-ui/out/")
}

// ProviderCapabilityHandler returns the capabilities.json for the provider
func (h *Handler) ProviderCapabilityHandler(
	w http.ResponseWriter,
	r *http.Request,
	prefObj *models.Preference,
	user *models.User,
	provider models.Provider,
) {
	provider.GetProviderCapabilities(w, r)
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
			http.Error(w, ErrFailToLoadExtensions(err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		_, _ = w.Write([]byte("{}"))
	} else {
		ServeReactComponentFromPackage(w, r, uiReqBasePath, provider)
	}
}
