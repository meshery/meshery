package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func ServeUI(w http.ResponseWriter, r *http.Request, reqBasePath, baseFolderPath string) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	reqURL := r.URL.RequestURI()
	reqURL = strings.Replace(reqURL, reqBasePath, "", 1)
	logrus.Debugf("reqURL: %s", reqURL)
	var filePath strings.Builder
	filePath.WriteString(reqURL)
	if reqURL == "/" || reqURL == "" {
		filePath.WriteString("index.html")
	} else if filepath.Ext(r.URL.RequestURI()) == "" {
		filePath.WriteString(".html")
	}
	finalPath := filepath.Join(baseFolderPath, filePath.String())
	logrus.Debugf("file path derived: %s", finalPath)
	http.ServeFile(w, r, finalPath)
}
