package handlers

import (
	"encoding/json"
	"net/http"
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
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	providers := map[string]models.ProviderProperties{}
	for _, p := range h.config.Providers {
		providers[p.Name()] = (p.GetProviderProperties())
	}
	bd, err := json.Marshal(providers)
	if err != nil {
		http.Error(w, "unable to marshal the providers", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(bd)
}

// ProviderUIHandler - serves providers UI
func (h *Handler) ProviderUIHandler(w http.ResponseWriter, r *http.Request) {
	ServeUI(w, r, "/provider", "../provider-ui/out/")
}
