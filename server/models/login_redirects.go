package models

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/meshery/meshery/server/core"
)

func resolvePostLoginRedirect(rawRef, fallback, host string) string {
	if rawRef == "" {
		return fallback
	}

	if decoded, ok := decodePostLoginRef(rawRef); ok {
		if target, ok := safePostLoginTarget(decoded, host); ok {
			return target
		}
	}

	if target, ok := safePostLoginTarget(rawRef, host); ok {
		return target
	}

	return fallback
}

// decodePostLoginRef accepts the server's raw-url encoding and the standard
// base64 produced by btoa in the Sign In link. A failed decode is not fatal:
// callers also try the raw value as a plaintext path.
func decodePostLoginRef(rawRef string) (string, bool) {
	if decoded, err := core.DecodeRefURL(rawRef); err == nil {
		return decoded, true
	}
	decoded, err := base64.StdEncoding.DecodeString(rawRef)
	if err != nil {
		return "", false
	}
	return string(decoded), true
}

// selectPostLoginRefValue returns the raw (encoded or plaintext) value to
// feed into resolvePostLoginRedirect when the auth flow returns to
// TokenHandler.
// Meshery is the source of truth for its own post-login destination: the value
// is captured into a cookie at InitiateLogin time and read back here. A cookie
// that carries a destination wins outright over a ?ref= the provider echoes
// back, since stale provider-side state (e.g. a synthesized ref baked into
// Hydra state during a custom-domain bounce) was the bug this routing change
// was introduced to fix. A cookie whose value is empty carries no destination
// and is therefore absent: ?ref= is used instead, which is how a Sign In taken
// from a page reached during an anonymous session keeps its query.
func selectPostLoginRefValue(r *http.Request, cookieName string) string {
	if ck, err := r.Cookie(cookieName); err == nil && ck.Value != "" {
		return ck.Value
	}
	return r.URL.Query().Get("ref")
}

// postLoginHost is the host an absolute ref is compared against. The
// configured public server URL wins over the inbound Host, which a proxy or
// client can rewrite.
func postLoginHost(r *http.Request) string {
	if raw, ok := r.Context().Value(MesheryServerURL).(string); ok {
		if configured, err := url.Parse(strings.TrimSpace(raw)); err == nil && configured.Hostname() != "" {
			return configured.Hostname()
		}
	}
	return (&url.URL{Host: r.Host}).Hostname()
}

// authInitiationPaths are routes whose job is to *start* authentication.
// Post-login redirects must never land on one of these, otherwise the browser
// immediately re-enters the OAuth dance and the original target is lost. The
// intermittent not-loading behavior was reproduced as exactly this:
// TokenHandler succeeded and then redirected to /user/login?provider=Meshery,
// which restarted InitiateLogin mid-mount.
// "/login" is the remote provider's own login page, not a Meshery route at all:
// a custom-domain bounce synthesizes it into the ref it echoes back, and honoring
// it served the catch-all handler as a 404 at playground.meshery.io/login.
var authInitiationPaths = []string{
	"/login",
	"/user/login",
	"/auth/login",
	"/api/user/token",
	"/provider",
}

// escapesOrigin reports whether a redirect target can resolve to an authority
// other than this server's. A leading "//" is protocol-relative, so it is
// disqualifying anywhere in the value. A backslash only delimits an authority
// where a browser reads one, which is the path - ending at the first "?" or
// "#". Past that point it is a query or fragment code point that browsers pass
// through without percent-encoding, so rejecting a whole ref over one would
// silently drop the user's page.
// This has to be asked of the value actually handed to http.Redirect, not only
// of the ref as received: reducing a same-host absolute ref to
// parsed.RequestURI() turns "https://kanvas.new//evil.example" into the
// protocol-relative "//evil.example", which http.Redirect writes out verbatim
// because the target it parses carries a Host and so skips its own
// normalization.
func escapesOrigin(target string) bool {
	if strings.HasPrefix(target, "//") {
		return true
	}
	authority := target
	if i := strings.IndexAny(authority, "?#"); i != -1 {
		authority = authority[:i]
	}
	return strings.Contains(authority, `\`)
}

// safePostLoginTarget validates a ref and returns the in-app path to redirect
// to. Anything that can leave the origin is rejected, both as the ref was
// received and as it was reduced - see escapesOrigin. The auth-path denylist
// runs over two spellings of the destination, because http.Redirect normalizes
// the Location it writes: the decoded parsed path, which catches "%2e%2e"
// traversal, and the effective Location, which is everything before the first
// "?" run through path.Clean exactly as the stdlib does it - it does not split
// on "#", so a fragment falls inside that when there is no query. Without the
// second check, "/extension/meshmap#/../../user/login" passes as its own path
// and is then cleaned to /user/login on the way out.
// Relative refs are kept as-is (path, query, and hash). An absolute ref on
// Meshery's own host is reduced to its path and query so an older client that
// sent window.location.href still lands on the design page. The scheme is not
// part of that comparison: the callback URL carries a hard-coded http:// on
// every deployment that leaves MESHERY_SERVER_CALLBACK_URL unset, so an https
// page would otherwise fail to match its own host. Every ref on another host
// is rejected.
func safePostLoginTarget(rawURL, host string) (string, bool) {
	if rawURL == "" || escapesOrigin(rawURL) {
		return "", false
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", false
	}

	target := rawURL
	if parsed.Scheme != "" || parsed.Host != "" {
		if !sameHost(parsed, host) {
			return "", false
		}
		target = parsed.RequestURI()
	} else if !strings.HasPrefix(rawURL, "/") {
		return "", false
	}

	if escapesOrigin(target) {
		return "", false
	}

	effective := target
	if i := strings.Index(effective, "?"); i != -1 {
		effective = effective[:i]
	}

	if !isAllowedAppPath(path.Clean(parsed.Path)) || !isAllowedAppPath(path.Clean(effective)) {
		return "", false
	}

	return target, true
}

func isAllowedAppPath(appPath string) bool {
	for _, p := range authInitiationPaths {
		if appPath == p || strings.HasPrefix(appPath, p+"/") {
			return false
		}
	}
	return true
}

func sameHost(ref *url.URL, host string) bool {
	return host != "" && ref.Hostname() != "" && strings.EqualFold(ref.Hostname(), host)
}
