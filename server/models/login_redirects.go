package models

import (
	"net/http"
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

// selectPostLoginRefValue returns the raw (encoded or plaintext) value to
// feed into resolvePostLoginRedirect when the auth flow returns to
// TokenHandler.
// Meshery is the source of truth for its own post-login destination: the value
// is captured into a cookie at InitiateLogin time and read back here. The
// ?ref= query param is a fallback for callers that never went through
// InitiateLogin (mesheryctl, direct extension callbacks) and for older
// provider deployments that still echo a ref back to us. We deliberately do
// not try to merge the two — when the cookie is present it wins outright,
// since stale provider-side state (e.g. a synthesized ref baked into Hydra
// state during a custom-domain bounce) was the bug this routing change was
// introduced to fix.
func selectPostLoginRefValue(r *http.Request, cookieName string) string {
	if ck, err := r.Cookie(cookieName); err == nil  {
		return ck.Value
	}
	return ""
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
