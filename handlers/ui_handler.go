package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ServeUI - helps serve static files for both meshery ui and provider ui
func ServeUI(w http.ResponseWriter, r *http.Request, reqBasePath, baseFolderPath string) {
	// if r.Method != http.MethodGet {
	// 	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	// 	return
	// }
	reqURL := r.URL.Path
	reqURL = strings.Replace(reqURL, reqBasePath, "", 1)

	var filePath strings.Builder

	filePath.WriteString(reqURL)
	if reqURL == "/" || reqURL == "" {
		filePath.WriteString("index.html")
	} else if filepath.Ext(reqURL) == "" && strings.HasPrefix(reqURL, "/extension") {
		filePath.Reset()
		filePath.WriteString("/extension/")
		filePath.WriteString("[component]")
		filePath.WriteString(".html")

		logrus.Debug("Generated path: ", filePath.String())
	} else if filepath.Ext(reqURL) == "" {
		filePath.WriteString(".html")
	}

	finalPath := filepath.Join(baseFolderPath, filePath.String())
	http.ServeFile(w, r, finalPath)
}

// else if strings.HasPrefix(reqURL, "/extension") {
// 		filePath.WriteString("/extension/")
// 		filePath.WriteString("[component]")
// 		filePath.WriteString(".html")
// 		logrus.Debug("Generated path: ", filePath.String())
// 	}
