package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

var dynamicUIEndpoints = map[string]string{
	"/extension": "/extension/[...component].html",
}

func isDynamicUIEndpoint(reqURL string) bool {
	if filepath.Ext(reqURL) != "" {
		return false
	}

	for extension := range dynamicUIEndpoints {
		if strings.HasPrefix(reqURL, extension) {
			return true
		}
	}

	return false
}

func getDynamicUIEndpoint(reqURL string) string {
	for prefix, mapping := range dynamicUIEndpoints {
		if strings.HasPrefix(reqURL, prefix) {
			return mapping
		}
	}

	return ""
}

// ServeUI - helps serve static files for both meshery ui and provider ui
func (h *Handler) ServeUI(w http.ResponseWriter, r *http.Request, reqBasePath, baseFolderPath string) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 	return
	// }
	reqURL := r.URL.Path
	reqURL = strings.Replace(reqURL, reqBasePath, "", 1)
	var filePath strings.Builder

	provider, ok := h.config.Providers[h.Provider]
	if ok && provider != nil {
		redirect, hasRedirect := provider.GetProviderProperties().Redirects[reqURL]

		if hasRedirect {
			http.Redirect(w, r, redirect, http.StatusPermanentRedirect)
			return
		}
	}

	filePath.WriteString(reqURL)
	if reqURL == "/" || reqURL == "" {
		filePath.WriteString("index.html")
	} else if isDynamicUIEndpoint(reqURL) {
		fmt.Println("serving dynamic ui endpoint: ", r.URL.Path, reqURL)
		filePath.Reset()
		filePath.WriteString(getDynamicUIEndpoint(reqURL))

		fmt.Println("Generated path: ", filePath.String())
	} else if filepath.Ext(reqURL) == "" {
		filePath.WriteString(".html")
	}
	finalPath := filepath.Join(baseFolderPath, filePath.String())
	http.ServeFile(w, r, finalPath)

}
