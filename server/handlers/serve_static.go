package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshkit/logger"
)

var dynamicUIEndpointMappings = map[string]string{
	"/extension":   "/extension/[...component].html",
	"/management/": "/management/[mesh].html",
}

func getDynamicUIEndpoint(reqURL string) string {
	for prefix, mapping := range dynamicUIEndpointMappings {
		if strings.Contains(reqURL, prefix) {
			return strings.Replace(reqURL, prefix, mapping, 1)
		}
	}
	return ""
}

func isDynamicUIEndpoint(reqURL string) bool {
	for endpoint := range dynamicUIEndpointMappings {
		if strings.Contains(reqURL, endpoint) {
			return true
		}
	}
	return false
}

// ServerUI - helps serve static files for Meshery UI
func ServeUI(w http.ResponseWriter, r *http.Request, reqBasePath, baseFolderPath string, log logger.Handler) {
	reqURL := r.URL.Path
	reqURL = strings.Replace(reqURL, reqBasePath, "", 1)

	var filePath strings.Builder

	filePath.WriteString(reqURL)
	if reqURL == "/" || reqURL == "" {
		filePath.WriteString("index.html")
	} else if isDynamicUIEndpoint(reqURL) {
		filePath.Reset()
		filePath.WriteString(getDynamicUIEndpoint(reqURL))

		log.Debug("Generated path: ", filePath.String())
	} else if filepath.Ext(reqURL) == "" {
		filePath.WriteString(".html")
	}

	finalPath := filepath.Join(baseFolderPath, filePath.String())
	http.ServeFile(w, r, finalPath)
}
