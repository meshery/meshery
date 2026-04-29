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
		redirectURL := "/user/login"
		if r.URL.RawQuery != "" {
			redirectURL += "?" + r.URL.RawQuery
		}
		http.Redirect(w, r, redirectURL, http.StatusFound)
		return
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
