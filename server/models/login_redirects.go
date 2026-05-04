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

// authInitiationPaths are server routes whose job is to *start* authentication.
// Post-login redirects must never land on one of these, otherwise the browser
// immediately re-enters the OAuth dance and the original target is lost. The
// intermittent Kanvas-never-loads behavior was reproduced as exactly this:
// TokenHandler succeeded and then redirected to /user/login?provider=Layer5,
// which restarted InitiateLogin mid-mount.
var authInitiationPaths = []string{
	"/user/login",
	"/auth/login",
	"/api/user/token",
	"/provider",
}

// isSafeRedirect validates that a decoded ref URL is a relative in-app path
// to prevent open redirects. It rejects absolute URLs (with scheme/host),
// protocol-relative URLs (starting with //), and auth-initiation paths that
// would cause a post-login redirect loop.
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

	if !strings.HasPrefix(rawURL, "/") {
		return false
	}

	for _, p := range authInitiationPaths {
		if parsed.Path == p || strings.HasPrefix(parsed.Path, p+"/") {
			return false
		}
	}

	return true
}
