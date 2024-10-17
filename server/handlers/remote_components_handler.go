package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/server/models"
)

// ServeReactComponentFromPackage serves static assets from the package
// located at provider.PackageLocation()
func ServeReactComponentFromPackage(
	w http.ResponseWriter,
	r *http.Request,
	reqBasePath string,
	provider models.Provider,
) {
	escPath := r.URL.EscapedPath()
	fmt.Println("[package location] ", provider.PackageLocation())
	assetPath := strings.Replace(escPath, reqBasePath, provider.PackageLocation(), 1)

	http.ServeFile(w, r, assetPath)
}
