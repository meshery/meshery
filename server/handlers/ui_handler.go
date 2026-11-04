package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var dynamicUIEndpoints = map[string]string{
	"/extension": "/extension/[...component].html",
}

// immutableAssetPrefixes are the URL prefixes (evaluated against the
// base-path-stripped request URL) whose responses Next.js fingerprints with a
// content hash. Because the hash is part of the URL, the bytes at a given URL
// never change, so they are safe to cache forever and self-bust across releases.
var immutableAssetPrefixes = []string{
	"/_next/static/chunks/",
	"/_next/static/css/",
	"/_next/static/media/",
}

// uiCacheHeaders computes the release-scoped caching headers for a UI response.
//
// It implements a two-tier scheme that lets a CDN or caching reverse proxy in
// front of Meshery cache UI content for the lifetime of a release and auto-bust
// it on the next deploy WITHOUT any CDN purge API call:
//
//   - Tier A - immutable versioned assets: requests under
//     /_next/static/{chunks,css,media}/ carry a content hash in the URL, so they
//     are returned with "public, max-age=31536000, immutable". A new release ships
//     new hashes (new URLs), so old entries simply fall out of use - no purge.
//
//   - Tier B - version-validated HTML: HTML documents live at stable URLs
//     (/ -> index.html, <route>.html, /extension/[...component].html), so they
//     must never be cached blindly. They are returned with "public, no-cache"
//     plus a strong ETag set to the build/release version. "no-cache" means the
//     edge/browser may store the response but must revalidate on every use; while
//     the release is unchanged the origin answers If-None-Match with 304, and a
//     new release changes BUILD -> changes the ETag -> the next revalidation
//     returns 200 with fresh HTML. Browser caches cannot be purged, so HTML is
//     never given a positive max-age.
//
// reqURL is the request URL AFTER reqBasePath has been stripped; filePath is the
// resolved on-disk relative path (used only to detect HTML via its suffix).
// It returns the Cache-Control value and the ETag value (already wrapped in the
// HTTP strong-validator double-quote form, or "" when no validator applies).
func uiCacheHeaders(reqURL, filePath string) (cacheControl, etag string) {
	for _, prefix := range immutableAssetPrefixes {
		if strings.HasPrefix(reqURL, prefix) {
			// Tier A.
			return "public, max-age=31536000, immutable", ""
		}
	}

	if strings.HasSuffix(filePath, ".html") {
		// Tier B.
		cacheControl = "public, no-cache"
		if v := viper.GetString("BUILD"); v != "" && v != "Not set" {
			etag = `"` + v + `"`
		}
		return cacheControl, etag
	}

	return "", ""
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

	// Apply release-scoped caching headers BEFORE http.ServeFile, but ONLY when
	// the target is an existing regular file. If it is missing, http.ServeFile
	// responds 404; setting long-lived/immutable headers first would let a CDN or
	// browser cache that 404 - e.g. a content-hashed asset requested mid-deploy
	// against a pod that doesn't have it yet would be pinned as "immutable" for a
	// year, breaking the URL until a manual purge: the exact failure this scheme
	// exists to avoid. Skipping the headers leaves the 404 uncacheable.
	//
	// Ordering is load-bearing for the success path: http.ServeFile ->
	// http.ServeContent reads the ETag header to emit conditional 304 responses
	// and then writes the body, so the headers must already be present on
	// w.Header() at call time. See uiCacheHeaders for the two-tier (immutable
	// assets vs. version-validated HTML) cache-busting scheme.
	if stat, err := os.Stat(finalPath); err == nil && !stat.IsDir() {
		if cacheControl, etag := uiCacheHeaders(reqURL, filePath.String()); cacheControl != "" {
			w.Header().Set("Cache-Control", cacheControl)
			if etag != "" {
				w.Header().Set("ETag", etag)
			}
		}
	}

	http.ServeFile(w, r, finalPath)

}
