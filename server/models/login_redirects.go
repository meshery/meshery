package models

import (
	"net/url"
	"strings"

	"github.com/meshery/meshery/server/core"
)

func resolvePostLoginRedirect(rawRef, fallback string) string {
	if rawRef == "" {
		return fallback
	}

	if decoded, err := core.DecodeRefURL(rawRef); err == nil && isSafeRedirect(decoded) {
		return decoded
	}

	if isSafeRedirect(rawRef) {
		return rawRef
	}

	return fallback
}

// isSafeRedirect validates that a decoded ref URL is a relative in-app path
// to prevent open redirects. It rejects absolute URLs (with scheme/host) and
// protocol-relative URLs (starting with //).
func isSafeRedirect(rawURL string) bool {
	if rawURL == "" || strings.HasPrefix(rawURL, "//") {
		return false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	if parsed.Scheme != "" || parsed.Host != "" {
		return false
	}

	return strings.HasPrefix(rawURL, "/")
}
