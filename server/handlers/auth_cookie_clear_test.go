package handlers

import (
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

func TestEffectiveTLDPlusOne_PSLParity(t *testing.T) {
	cases := []struct {
		name string
		host string
		want string
	}{
		{
			name: "single dot eTLD",
			host: "foo.example.com",
			want: "example.com",
		},
		{
			name: "deeper subdomain on single eTLD",
			host: "a.b.c.example.com",
			want: "example.com",
		},
		// The PSL case that a naive strings.SplitN(host, ".", 2)[1]
		// gets wrong: "example.co.uk" would yield "co.uk", which is a
		// public suffix — and any cookie scoped there would be
		// rejected by the browser as a public-suffix domain.
		{
			name: "multi-label public suffix .co.uk",
			host: "foo.example.co.uk",
			want: "example.co.uk",
		},
		{
			name: "multi-label public suffix .com.au",
			host: "deep.sub.example.com.au",
			want: "example.com.au",
		},
		// IP literals must NOT be treated as cross-host scopes. Returning
		// an empty string here is the signal to callers to skip the
		// Domain-scoped clearance.
		{
			name: "IPv4 literal returns empty",
			host: "127.0.0.1",
			want: "",
		},
		{
			name: "IPv6 literal returns empty",
			host: "::1",
			want: "",
		},
		// "localhost" is a single label — no cross-subdomain cookie is
		// even possible, and Domain=localhost has surprising semantics.
		{
			name: "localhost single label returns empty",
			host: "localhost",
			want: "",
		},
		// A host that IS itself a public suffix has no eTLD+1; reject it.
		{
			name: "bare public suffix returns empty",
			host: "co.uk",
			want: "",
		},
		{
			name: "empty input returns empty",
			host: "",
			want: "",
		},
		{
			name: "whitespace-only input returns empty",
			host: "   ",
			want: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := effectiveTLDPlusOne(tc.host)
			if got != tc.want {
				t.Fatalf("effectiveTLDPlusOne(%q) = %q, want %q", tc.host, got, tc.want)
			}
		})
	}
}

func TestAuthCookieScopes(t *testing.T) {
	cases := []struct {
		name         string
		requestHost  string
		providerURLs []string
		// wantDomains is the expected SORTED set of non-empty Domain
		// values. We always expect at least one entry with empty
		// Domain ("no Domain attribute"), and we check that separately.
		wantDomains []string
	}{
		{
			name:        "playground -> cloud (canonical scope)",
			requestHost: "playground.meshery.io",
			providerURLs: []string{
				"https://cloud.layer5.io",
			},
			wantDomains: []string{
				"cloud.layer5.io", // bare host
				"layer5.io",       // eTLD+1 of cloud.layer5.io
				"meshery.io",      // eTLD+1 of playground.meshery.io
			},
		},
		{
			name:        "self-hosted localhost with cloud provider — no eTLD+1 for current host",
			requestHost: "localhost",
			providerURLs: []string{
				"https://cloud.layer5.io",
			},
			wantDomains: []string{
				"cloud.layer5.io",
				"layer5.io",
			},
		},
		{
			name:        "BYOC custom-domain provider",
			requestHost: "playground.meshery.io",
			providerURLs: []string{
				"https://cloud.example.co.uk",
			},
			wantDomains: []string{
				"cloud.example.co.uk", // bare host
				"example.co.uk",       // PSL-correct eTLD+1
				"meshery.io",          // request host eTLD+1
			},
		},
		{
			name:        "multiple providers deduped",
			requestHost: "playground.meshery.io",
			providerURLs: []string{
				"https://cloud.layer5.io",
				"https://cloud.layer5.io", // exact dup
				"https://other.layer5.io", // shares eTLD+1
			},
			wantDomains: []string{
				"cloud.layer5.io",
				"layer5.io",
				"meshery.io",
				"other.layer5.io",
			},
		},
		{
			name:         "no providers, current host eTLD+1 still emitted",
			requestHost:  "app.example.com",
			providerURLs: nil,
			wantDomains:  []string{"example.com"},
		},
		{
			name:        "current host equals provider host — no redundant Domain entry for current host",
			requestHost: "cloud.layer5.io",
			providerURLs: []string{
				"https://cloud.layer5.io",
			},
			wantDomains: []string{
				// "cloud.layer5.io" is NOT here: it's the request
				// host and the no-Domain entry already covers it.
				"layer5.io",
			},
		},
		{
			name:         "ip-literal request host, no providers",
			requestHost:  "127.0.0.1",
			providerURLs: nil,
			wantDomains:  []string{},
		},
		{
			name:        "malformed provider URL skipped",
			requestHost: "playground.meshery.io",
			providerURLs: []string{
				"::::not-a-url",
				"",
				"   ",
				"https://cloud.layer5.io",
			},
			wantDomains: []string{
				"cloud.layer5.io",
				"layer5.io",
				"meshery.io",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			scopes := authCookieScopes(tc.requestHost, tc.providerURLs)

			// 1. The first scope must always be the no-Domain
			// clearance — that's the back-compat guarantee.
			if len(scopes) == 0 || scopes[0].Domain != "" {
				t.Fatalf("expected first scope to be no-Domain, got %+v", scopes)
			}

			// 2. Collect non-empty Domains and compare as a set.
			var gotDomains []string
			seen := map[string]struct{}{}
			for _, s := range scopes[1:] {
				if s.Domain == "" {
					t.Errorf("unexpected duplicate no-Domain entry in scopes: %+v", scopes)
					continue
				}
				if _, dup := seen[s.Domain]; dup {
					t.Errorf("duplicate Domain %q in scopes: %+v", s.Domain, scopes)
				}
				seen[s.Domain] = struct{}{}
				gotDomains = append(gotDomains, s.Domain)
			}
			sort.Strings(gotDomains)
			want := append([]string(nil), tc.wantDomains...)
			sort.Strings(want)
			if !equalStringSlices(gotDomains, want) {
				t.Fatalf("Domain set mismatch:\ngot:  %v\nwant: %v", gotDomains, want)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestClearAuthCookies_EmitsAllScopes(t *testing.T) {
	w := httptest.NewRecorder()
	cookieNames := []string{"token", "session_cookie", "meshery-provider"}
	scopes := []cookieScope{
		{Domain: ""},             // no Domain
		{Domain: "layer5.io"},    // PSL-correct eTLD+1
		{Domain: "meshery.io"},   // BYOC eTLD+1
		{Domain: "cloud.layer5.io"},
	}

	clearAuthCookies(w, cookieNames, scopes)

	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })
	setCookies := resp.Header.Values("Set-Cookie")
	if len(setCookies) != len(cookieNames)*len(scopes) {
		t.Fatalf("expected %d Set-Cookie headers, got %d: %v", len(cookieNames)*len(scopes), len(setCookies), setCookies)
	}

	// For each (cookieName, scope) pair, assert exactly one matching
	// Set-Cookie header that's a clearance (Max-Age=0 from MaxAge=-1) and
	// that the Domain attribute is set correctly (or absent).
	for _, name := range cookieNames {
		for _, s := range scopes {
			matches := 0
			for _, c := range setCookies {
				if !strings.HasPrefix(c, name+"=") {
					continue
				}
				domainAttr := extractCookieAttr(c, "Domain")
				if !strings.EqualFold(domainAttr, s.Domain) {
					continue
				}
				if !strings.Contains(c, "Max-Age=0") {
					t.Errorf("Set-Cookie for %q@%q should be a clearance (Max-Age=0), got %q", name, s.Domain, c)
				}
				matches++
			}
			if matches != 1 {
				t.Errorf("expected exactly one Set-Cookie for (%q, Domain=%q), got %d (headers=%v)", name, s.Domain, matches, setCookies)
			}
		}
	}
}

// extractCookieAttr returns the value of the given attribute in a
// Set-Cookie header, or empty string if absent. Case-insensitive on the
// attribute name, as is RFC 6265 standard.
func extractCookieAttr(setCookie, attr string) string {
	parts := strings.Split(setCookie, ";")
	for _, p := range parts[1:] {
		p = strings.TrimSpace(p)
		eq := strings.IndexByte(p, '=')
		if eq < 0 {
			continue
		}
		if strings.EqualFold(p[:eq], attr) {
			return p[eq+1:]
		}
	}
	return ""
}

// logoutSpyProvider is a minimal Provider implementation used to exercise
// LogoutHandler without standing up a real persister or remote-provider
// chain. It embeds DefaultLocalProvider so we inherit all method
// implementations we don't need to override.
type logoutSpyProvider struct {
	*models.DefaultLocalProvider
	providerURL string
	name        string
}

func (p *logoutSpyProvider) Name() string             { return p.name }
func (p *logoutSpyProvider) GetProviderURL() string   { return p.providerURL }
func (p *logoutSpyProvider) Logout(_ http.ResponseWriter, _ *http.Request) error { return nil }

// DeleteCapabilitiesForUser short-circuits the embedded persister so the
// test doesn't need to stand up a real DB. LogoutHandler ignores the
// return value anyway, so a no-op is fine.
func (p *logoutSpyProvider) DeleteCapabilitiesForUser(_ string) error { return nil }

// TestLogoutHandler_ClearsAcrossAllProviderScopes asserts the headline CDA
// cookie-lifecycle invariant: when a user logs out, the response carries
// Set-Cookie clearances on every host scope a CDA flow might have set
// cookies on — current host, current host eTLD+1 (PSL-correct), and
// each configured remote provider's host + eTLD+1.
func TestLogoutHandler_ClearsAcrossAllProviderScopes(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	cloudProv := &logoutSpyProvider{DefaultLocalProvider: local, providerURL: "https://cloud.layer5.io", name: "Cloud"}
	byocProv := &logoutSpyProvider{DefaultLocalProvider: local, providerURL: "https://cloud.example.co.uk", name: "BYOC"}

	providers := map[string]models.Provider{
		"Cloud": cloudProv,
		"BYOC":  byocProv,
	}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/api/user/logout", nil)
	req.Host = "playground.meshery.io"
	rec := httptest.NewRecorder()

	user := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}
	h.LogoutHandler(rec, req, user, cloudProv)

	resp := rec.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected logout redirect (302), got %d", resp.StatusCode)
	}

	setCookies := resp.Header.Values("Set-Cookie")
	if len(setCookies) == 0 {
		t.Fatal("expected logout to emit Set-Cookie clearances, got none")
	}

	// The three cookie names we expect on every scope.
	expectedNames := []string{"meshery-provider", "token", "session_cookie"}

	// The expected NON-EMPTY Domain scopes — a clearance with no Domain
	// attribute is also expected but tested separately.
	expectedDomains := []string{
		"cloud.layer5.io",     // bare provider host
		"layer5.io",           // cloud eTLD+1
		"cloud.example.co.uk", // BYOC bare provider host
		"example.co.uk",       // PSL-correct BYOC eTLD+1 (NOT co.uk!)
		"meshery.io",          // request host eTLD+1
	}

	for _, name := range expectedNames {
		// Bare/no-Domain clearance.
		if !findClearanceCookie(setCookies, name, "") {
			t.Errorf("missing no-Domain clearance for cookie %q (headers=%v)", name, setCookies)
		}
		for _, dom := range expectedDomains {
			if !findClearanceCookie(setCookies, name, dom) {
				t.Errorf("missing clearance for cookie %q on Domain=%q (headers=%v)", name, dom, setCookies)
			}
		}
	}

	// Specifically assert that NO clearance was emitted for Domain=co.uk —
	// that's the Gap A regression we're guarding against.
	for _, c := range setCookies {
		if strings.EqualFold(strings.TrimSpace(extractCookieAttr(c, "Domain")), "co.uk") {
			t.Fatalf("PSL parity violated: Set-Cookie emitted at Domain=co.uk: %q", c)
		}
	}
}

// findClearanceCookie reports whether setCookies contains a Set-Cookie
// header with the given name, Domain (empty string = no Domain attribute),
// and a clearance Max-Age of 0.
func findClearanceCookie(setCookies []string, name, domain string) bool {
	for _, c := range setCookies {
		if !strings.HasPrefix(c, name+"=") {
			continue
		}
		got := extractCookieAttr(c, "Domain")
		if !strings.EqualFold(strings.TrimSpace(got), strings.TrimSpace(domain)) {
			continue
		}
		if !strings.Contains(c, "Max-Age=0") {
			continue
		}
		return true
	}
	return false
}

func TestHostFromRequest(t *testing.T) {
	cases := []struct {
		name     string
		hostHdr  string
		want     string
	}{
		{name: "bare host", hostHdr: "playground.meshery.io", want: "playground.meshery.io"},
		{name: "host:port", hostHdr: "localhost:9081", want: "localhost"},
		{name: "ipv4:port", hostHdr: "127.0.0.1:9081", want: "127.0.0.1"},
		{name: "ipv6:port", hostHdr: "[::1]:9081", want: "::1"},
		{name: "empty host header", hostHdr: "", want: ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Host = tc.hostHdr
			got := hostFromRequest(req)
			if got != tc.want {
				t.Fatalf("hostFromRequest(host=%q) = %q, want %q", tc.hostHdr, got, tc.want)
			}
		})
	}
}
