package models

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestResolvePostLoginRedirect(t *testing.T) {
	t.Parallel()

	const fallback = "/"

	const host = "kanvas.new"

	tests := []struct {
		name     string
		rawRef   string
		host     string
		expected string
	}{
		{
			name:     "empty ref falls back",
			rawRef:   "",
			expected: fallback,
		},
		{
			name:     "encoded in app path is decoded",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap")),
			expected: "/extension/meshmap",
		},
		{
			name:     "plain in app path is preserved",
			rawRef:   "/extension/meshmap",
			expected: "/extension/meshmap",
		},
		{
			name:     "relative ref keeps search and hash",
			rawRef:   "/extension/meshmap?mode=design#canvas",
			expected: "/extension/meshmap?mode=design#canvas",
		},
		{
			name:     "encoded absolute url falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/phish")),
			host:     host,
			expected: fallback,
		},
		{
			name:     "plain absolute url falls back",
			rawRef:   "https://evil.example/phish",
			host:     host,
			expected: fallback,
		},
		{
			name:     "same-origin absolute ref reduces to path and query",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://kanvas.new/extension/meshmap?mode=design")),
			host:     host,
			expected: "/extension/meshmap?mode=design",
		},
		{
			name:     "standard base64 same-origin absolute ref reduces to path and query",
			rawRef:   base64.StdEncoding.EncodeToString([]byte("https://kanvas.new/extension/meshmap?mode=design#canvas")),
			host:     host,
			expected: "/extension/meshmap?mode=design",
		},
		{
			name:     "cross-origin absolute ref is rejected",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/extension/meshmap?mode=design")),
			host:     host,
			expected: fallback,
		},
		// A deployment that leaves MESHERY_SERVER_CALLBACK_URL unset derives its
		// own origin as http://<host> even when it is served over TLS, so the
		// scheme of an absolute ref carries no information about whether it is
		// ours. Comparing it dropped mode=design from a genuinely same-host ref.
		{
			name:     "same-host absolute ref with a different scheme is accepted",
			rawRef:   "http://kanvas.new/extension/meshmap?mode=design",
			host:     host,
			expected: "/extension/meshmap?mode=design",
		},
		{
			name:     "same-host absolute ref on a different port is accepted",
			rawRef:   "https://kanvas.new:8443/extension/meshmap?mode=design",
			host:     host,
			expected: "/extension/meshmap?mode=design",
		},
		{
			name:     "host suffix of the server host is rejected",
			rawRef:   "https://kanvas.new.evil.example/extension/meshmap?mode=design",
			host:     host,
			expected: fallback,
		},
		{
			name:     "absolute ref is rejected when the server host is unknown",
			rawRef:   "https://kanvas.new/extension/meshmap?mode=design",
			expected: fallback,
		},
		{
			name:     "same-origin absolute auth path falls back",
			rawRef:   "https://kanvas.new/user/login?provider=Meshery",
			host:     host,
			expected: fallback,
		},
		{
			name:     "invalid base64 falls back",
			rawRef:   "not-base64",
			expected: fallback,
		},
		// A browser resolves a backslash in the relative-slash state as an
		// authority delimiter, so Location: /\evil.example leaves the origin.
		// DefaultLocalProvider.InitiateLogin reaches this with no authentication
		// at all, via GET /user/login?ref=<b64>.
		{
			name:     "backslash ref falls back",
			rawRef:   `/\evil.example`,
			expected: fallback,
		},
		{
			name:     "encoded backslash ref falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte(`/\evil.example`)),
			expected: fallback,
		},
		{
			name:     "backslash after a leading slash pair falls back",
			rawRef:   `/\/evil.example`,
			expected: fallback,
		},
		{
			name:     "backslash in a same-host absolute ref falls back",
			rawRef:   `https://kanvas.new/\evil.example`,
			host:     host,
			expected: fallback,
		},
		// http.Redirect runs path.Clean over the Location it writes, so a ref
		// that only reaches an auth-initiation path after normalization still
		// restarts the OAuth dance. The denylist has to see the cleaned path.
		{
			name:     "dot-dot traversal onto an auth path falls back",
			rawRef:   "/../user/login",
			expected: fallback,
		},
		{
			name:     "traversal from an app path onto an auth path falls back",
			rawRef:   "/extension/meshmap/../../user/login",
			expected: fallback,
		},
		{
			name:     "encoded traversal onto an auth path falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/extension/../api/user/token")),
			expected: fallback,
		},
		{
			name:     "traversal that stays in the app is preserved",
			rawRef:   "/extension/meshmap/../meshmap?mode=design",
			expected: "/extension/meshmap/../meshmap?mode=design",
		},
		// A custom-domain Cloud bounce synthesizes its own /login into the ref it
		// echoes back. Meshery serves no /login route, so honoring it landed
		// playground.meshery.io on the catch-all handler as a 404.
		{
			name:     "provider echoed /login falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/login")),
			expected: fallback,
		},
		{
			name:     "/login with query falls back",
			rawRef:   "/login?provider=Meshery",
			expected: fallback,
		},
		// Browsers do not percent-encode a backslash in a query or fragment, so
		// window.location can carry a literal one into the ref. Only the path
		// can be read as an authority, so rejecting the whole ref over one in
		// the query or fragment would silently drop the user's own page.
		{
			name:     "backslash in the query is preserved",
			rawRef:   `/extension/meshmap?design=a\b`,
			expected: `/extension/meshmap?design=a\b`,
		},
		{
			name:     "backslash in the fragment is preserved",
			rawRef:   `/extension/meshmap#c\d`,
			expected: `/extension/meshmap#c\d`,
		},
		{
			name:     "backslash in the query of a same-host absolute ref is preserved",
			rawRef:   `https://kanvas.new/extension/meshmap?design=a\b`,
			host:     host,
			expected: `/extension/meshmap?design=a\b`,
		},
		{
			name:     "backslash in the path still falls back when a query follows",
			rawRef:   `/\evil.example?x=1`,
			expected: fallback,
		},
		{
			name:     "backslash in the path still falls back when a fragment follows",
			rawRef:   `/\evil.example#x`,
			expected: fallback,
		},
		// Reducing a same-host absolute ref to its path and query can itself
		// produce a protocol-relative target. http.Redirect parses that back,
		// finds a Host, and so skips the normalization that would have
		// collapsed it - writing "//evil.example/..." out verbatim.
		{
			name:     "same-host absolute ref with a double-slash path falls back",
			rawRef:   "https://kanvas.new//evil.example/pretend-login?x=1",
			host:     host,
			expected: fallback,
		},
		{
			name:     "encoded same-host absolute ref with a double-slash path falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://kanvas.new//evil.example")),
			host:     host,
			expected: fallback,
		},
		{
			name:     "same-host absolute ref with a triple-slash path falls back",
			rawRef:   "https://kanvas.new///evil.example",
			host:     host,
			expected: fallback,
		},
		// http.Redirect splits the Location at the first "?" and runs path.Clean
		// over everything before it, the fragment included. A ref with a
		// fragment and no query therefore has its ".." applied to the path
		// after the denylist has already passed the parsed path.
		{
			name:     "fragment traversal onto an auth path falls back",
			rawRef:   "/extension/meshmap#/../../user/login",
			expected: fallback,
		},
		{
			name:     "fragment traversal onto the token endpoint falls back",
			rawRef:   "/extension/meshmap#/../../api/user/token",
			expected: fallback,
		},
		{
			name:     "encoded fragment traversal onto an auth path falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap#/../../login")),
			expected: fallback,
		},
		{
			name:     "plain fragment is preserved",
			rawRef:   "/extension/meshmap#canvas",
			expected: "/extension/meshmap#canvas",
		},
		{
			name:     "same-host absolute /login falls back",
			rawRef:   "https://kanvas.new/login",
			host:     host,
			expected: fallback,
		},
		// Regression coverage: /user/login and /api/user/token are auth
		// initiation paths. Redirecting to them after a successful token
		// exchange re-enters the OAuth dance and caused to hang on
		// the loading splash indefinitely (meshery-server-1345 followed by
		// a second InitiateLogin in the same second).
		{
			name:     "plain /user/login ref falls back",
			rawRef:   "/user/login",
			expected: fallback,
		},
		{
			name:     "/user/login with query falls back",
			rawRef:   "/user/login?provider=Meshery",
			expected: fallback,
		},
		{
			name:     "encoded /user/login ref falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/user/login?provider=Meshery")),
			expected: fallback,
		},
		{
			name:     "plain /api/user/token ref falls back",
			rawRef:   "/api/user/token",
			expected: fallback,
		},
		{
			name:     "/provider ref falls back",
			rawRef:   "/provider?ref=xyz",
			expected: fallback,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := resolvePostLoginRedirect(tc.rawRef, fallback, tc.host)
			if actual != tc.expected {
				t.Fatalf("expected redirect %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestSelectPostLoginRefValue(t *testing.T) {
	t.Parallel()

	const cookieName = "playground.meshery.io_ref"
	const cookieValue = "L2V4dGVuc2lvbi9tZXNobWFw" // base64 of /extension/meshmap
	const queryValue = "L2Rhc2hib2FyZA"            // base64 of /dashboard

	tests := []struct {
		name     string
		cookie   *http.Cookie
		query    string
		expected string
	}{
		{
			name:     "cookie value is used when present",
			cookie:   &http.Cookie{Name: cookieName, Value: cookieValue},
			expected: cookieValue,
		},
		// The cookie wins when it carries a destination, including when a
		// provider echoes a different ?ref=. That echo is what landed
		// playground.meshery.io on a 404. When the cookie was never set (Sign In
		// goes straight to the provider), ?ref= is the fallback the comment on
		// selectPostLoginRefValue describes.
		{
			name:     "cookie wins over ?ref= query param",
			cookie:   &http.Cookie{Name: cookieName, Value: cookieValue},
			query:    "?ref=" + queryValue,
			expected: cookieValue,
		},
		{
			name:     "uses ?ref= query param when cookie is missing",
			query:    "?ref=" + queryValue,
			expected: queryValue,
		},
		// An empty-valued cookie still parses as present. Treating it as a
		// winning destination suppressed the ?ref= that a Sign In taken after an
		// anonymous bootstrap carries, and dropped mode=design from the return.
		{
			name:     "uses ?ref= query param when cookie value is empty",
			cookie:   &http.Cookie{Name: cookieName, Value: ""},
			query:    "?ref=" + queryValue,
			expected: queryValue,
		},
		{
			name:     "returns empty when cookie is missing",
			expected: "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/user/token"+tc.query, nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			actual := selectPostLoginRefValue(req, cookieName)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestPostLoginHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		mesheryServeURL any
		requestHost     string
		expected        string
	}{
		{
			name:            "configured server url wins over inbound host",
			mesheryServeURL: "http://kanvas.new",
			requestHost:     "proxy.internal:8080",
			expected:        "kanvas.new",
		},
		{
			name:            "configured server url port is dropped",
			mesheryServeURL: "http://localhost:9081",
			requestHost:     "localhost:9081",
			expected:        "localhost",
		},
		{
			name:        "falls back to the inbound host without its port",
			requestHost: "kanvas.new:9081",
			expected:    "kanvas.new",
		},
		{
			name:            "blank configured server url falls back to the inbound host",
			mesheryServeURL: "   ",
			requestHost:     "kanvas.new",
			expected:        "kanvas.new",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/user/token", nil)
			req.Host = tc.requestHost
			if tc.mesheryServeURL != nil {
				req = req.WithContext(context.WithValue(req.Context(), MesheryServerURL, tc.mesheryServeURL))
			}
			if actual := postLoginHost(req); actual != tc.expected {
				t.Fatalf("expected host %q, got %q", tc.expected, actual)
			}
		})
	}
}

// Regression for the kanvas.new report: an anonymous bootstrap leaves an
// empty-valued ref cookie behind, the Sign In that follows returns to
// /api/user/token carrying the absolute page URL in ?ref=, and the deployment
// leaves MESHERY_SERVER_CALLBACK_URL unset so its own origin is http:// while
// the page is https://. Either the shadowing cookie or the scheme comparison
// alone was enough to drop mode=design from the return address.
func TestPostLoginRedirect_KeepsQueryAfterAnonymousBootstrap(t *testing.T) {
	t.Parallel()

	const cookieName = "kanvas.new_ref"
	ref := base64.StdEncoding.EncodeToString([]byte("https://kanvas.new/extension/meshmap?mode=design"))

	req := httptest.NewRequest(http.MethodGet, "/api/user/token?token=abc&ref="+url.QueryEscape(ref), nil)
	req.Host = "kanvas.new"
	req.AddCookie(&http.Cookie{Name: cookieName, Value: ""})
	req = req.WithContext(context.WithValue(req.Context(), MesheryServerURL, "http://kanvas.new"))

	actual := resolvePostLoginRedirect(selectPostLoginRefValue(req, cookieName), "/", postLoginHost(req))
	if actual != "/extension/meshmap?mode=design" {
		t.Fatalf("expected redirect %q, got %q", "/extension/meshmap?mode=design", actual)
	}
}

// The denylist exists to stop a post-login redirect from re-entering the OAuth
// dance, and what decides that is the Location the browser receives - not the
// ref as written. http.Redirect normalizes the Location, so the guarantee is
// only worth as much as it is against the real normalizer. Every hostile ref
// here reached an auth-initiation path through that normalization.
func TestPostLoginRedirect_LocationNeverReachesAnAuthPath(t *testing.T) {
	t.Parallel()

	const host = "kanvas.new"

	hostile := []string{
		"/extension/meshmap#/../../user/login",
		"/extension/meshmap#/../../api/user/token",
		"/extension/meshmap#/../../login",
		"/extension/meshmap#/../../auth/login",
		"/extension/meshmap#/../../provider",
		"/../user/login",
		"/extension/meshmap/../../user/login",
		base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap#/../../user/login")),
		base64.StdEncoding.EncodeToString([]byte("/extension/meshmap#/../../api/user/token")),
	}

	for _, rawRef := range hostile {
		rawRef := rawRef
		t.Run(rawRef, func(t *testing.T) {
			t.Parallel()

			target := resolvePostLoginRedirect(rawRef, "/", host)

			rec := httptest.NewRecorder()
			http.Redirect(rec, httptest.NewRequest(http.MethodGet, "/api/user/token", nil), target, http.StatusFound)

			location := rec.Header().Get("Location")
			landed, err := url.Parse(location)
			if err != nil {
				t.Fatalf("parse Location %q: %v", location, err)
			}
			if !isAllowedAppPath(landed.Path) {
				t.Fatalf("ref %q resolved to target %q, which http.Redirect emitted as Location %q - an auth-initiation path", rawRef, target, location)
			}
		})
	}
}

// The ordinary destinations have to survive the same normalization, or the
// guard above would be satisfied by rejecting everything.
func TestPostLoginRedirect_LocationKeepsOrdinaryDestinations(t *testing.T) {
	t.Parallel()

	const host = "kanvas.new"

	tests := map[string]string{
		"/extension/meshmap?mode=design":                   "/extension/meshmap?mode=design",
		"/extension/meshmap?mode=design#canvas":            "/extension/meshmap?mode=design#canvas",
		"/extension/meshmap#canvas":                        "/extension/meshmap#canvas",
		"https://kanvas.new/extension/meshmap?mode=design": "/extension/meshmap?mode=design",
		`/extension/meshmap?design=a\b`:                    `/extension/meshmap?design=a\b`,
		`/extension/meshmap#c\d`:                           `/extension/meshmap#c\d`,
	}

	for rawRef, want := range tests {
		rawRef, want := rawRef, want
		t.Run(rawRef, func(t *testing.T) {
			t.Parallel()

			target := resolvePostLoginRedirect(rawRef, "/", host)

			rec := httptest.NewRecorder()
			http.Redirect(rec, httptest.NewRequest(http.MethodGet, "/api/user/token", nil), target, http.StatusFound)

			if got := rec.Header().Get("Location"); got != want {
				t.Fatalf("ref %q emitted Location %q, want %q", rawRef, got, want)
			}
		})
	}
}

// "Keep every cross-origin ref rejected" is a property of the Location the
// browser receives, so it is asserted against the real normalizer. A Location
// that parses with a Host, or that carries a backslash, resolves to another
// authority - which is a phishing bounce off a just-authenticated session.
func TestPostLoginRedirect_LocationNeverLeavesTheOrigin(t *testing.T) {
	t.Parallel()

	const host = "kanvas.new"

	hostile := []string{
		"https://kanvas.new//evil.example/pretend-login?x=1",
		"https://kanvas.new//evil.example",
		"https://kanvas.new///evil.example",
		"http://kanvas.new//evil.example/pretend-login",
		"//evil.example/pretend-login",
		"///evil.example",
		`/\evil.example`,
		`/\/evil.example`,
		"https://evil.example/pretend-login",
		base64.RawURLEncoding.EncodeToString([]byte("https://kanvas.new//evil.example/pretend-login")),
		base64.StdEncoding.EncodeToString([]byte("https://kanvas.new//evil.example/pretend-login")),
	}

	for _, rawRef := range hostile {
		rawRef := rawRef
		t.Run(rawRef, func(t *testing.T) {
			t.Parallel()

			target := resolvePostLoginRedirect(rawRef, "/", host)

			rec := httptest.NewRecorder()
			http.Redirect(rec, httptest.NewRequest(http.MethodGet, "/api/user/token", nil), target, http.StatusFound)

			location := rec.Header().Get("Location")
			landed, err := url.Parse(location)
			if err != nil {
				t.Fatalf("parse Location %q: %v", location, err)
			}
			if landed.Host != "" || landed.Scheme != "" {
				t.Fatalf("ref %q resolved to target %q, which http.Redirect emitted as Location %q - authority %q", rawRef, target, location, landed.Host)
			}
			if strings.Contains(location, `\`) {
				t.Fatalf("ref %q resolved to target %q, which http.Redirect emitted as Location %q - a backslash a browser reads as an authority delimiter", rawRef, target, location)
			}
		})
	}
}
