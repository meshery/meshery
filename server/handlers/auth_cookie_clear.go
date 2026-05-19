// Package handlers : CDA cookie-lifecycle helpers.
//
// This file owns the cross-scope auth-cookie clearance used by LogoutHandler.
// It pairs with the symmetric clearance on the Meshery Cloud side: when a user
// logs out from a Meshery Server (e.g. playground.meshery.io) that bridges
// auth through Cloud's per-org Custom Domain Authentication (CDA), cookies
// could have been set on multiple host scopes:
//
//   - The Meshery Server's own host (no Domain attribute).
//   - The Meshery Cloud canonical scope (.layer5.io).
//   - A custom-domain BYOC scope (e.g. .meshery.io, or an off-eTLD apex like
//     example.com).
//
// A one-scope logout leaves the other scopes' cookies in place, and the next
// auth-required request can pick them back up and re-authenticate the
// browser against a provider the user thought they'd left. We mirror what
// meshery-cloud does on its side and iterate every host scope the server
// knows about, emitting a Set-Cookie clearance for each.
//
// The eTLD+1 derivation uses golang.org/x/net/publicsuffix (the PSL) so that
// hosts like foo.example.co.uk resolve to example.co.uk rather than co.uk.
// A naive strings.SplitN("example.co.uk", ".", 2)[1] returns "co.uk" — which
// is a public suffix, so any cookie scoped to it would be rejected by the
// browser as a "public suffix domain" anyway. The PSL path is the only
// correct way to do this in Go.
package handlers

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// cookieScope represents one host scope at which a Set-Cookie clearance
// should be emitted. An empty Domain means "no Domain attribute" — the
// clearance applies only to the request's exact host.
type cookieScope struct {
	// Domain is the value placed in the Set-Cookie Domain attribute. An
	// empty string means "omit the Domain attribute entirely", which
	// targets the current host only.
	Domain string
}

// hostFromRequest returns the bare hostname (no port) of the incoming
// request. It prefers r.Host (which Go populates from the Host header) and
// strips an optional port suffix. Falls back to an empty string when the
// host is not parseable; callers must treat the empty result as "no host
// scope known" and skip Domain attachment.
func hostFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	host := r.Host
	if host == "" {
		return ""
	}
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

// hostFromURL returns the bare hostname for a configured provider URL.
// Returns empty string if the URL cannot be parsed.
func hostFromURL(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return parsed.Hostname()
}

// effectiveTLDPlusOne returns the PSL-derived eTLD+1 for host (e.g.
// "foo.bar.example.co.uk" -> "example.co.uk"). It returns the empty string
// for IP literals, hosts that ARE a public suffix, single-label hosts like
// "localhost", and any input the PSL rejects. Callers must treat the empty
// result as "no cross-host clearance possible at this scope" and skip
// emitting a Domain-scoped Set-Cookie.
//
// We never want to set Domain=co.uk or Domain=localhost on a cookie: the
// browser will silently drop the attribute (public-suffix rule) or — worse
// for development — pin the cookie to a single label that overlaps with
// every dev project on the same machine.
func effectiveTLDPlusOne(host string) string {
	host = strings.TrimSpace(host)
	if host == "" {
		return ""
	}
	// Public-suffix lookup is undefined for IP literals; bail out early.
	if ip := net.ParseIP(host); ip != nil {
		return ""
	}
	// Reject hosts without at least one dot — a single label is either
	// "localhost" or a development hostname that cannot host a
	// cross-subdomain cookie regardless of PSL behavior.
	if !strings.Contains(host, ".") {
		return ""
	}
	etld1, err := publicsuffix.EffectiveTLDPlusOne(host)
	if err != nil {
		return ""
	}
	return etld1
}

// authCookieScopes returns the ordered, deduplicated list of host scopes
// at which auth cookies should be cleared on logout.
//
// The order is: (1) the current request's host (no Domain attribute), then
// (2) the current host's PSL-derived eTLD+1, then (3) for each configured
// remote provider URL the same pair (bare host + eTLD+1). Duplicates are
// dropped to avoid emitting redundant Set-Cookie headers.
//
// Empty scope.Domain values are preserved in the list — they signal "emit
// the no-Domain clearance once". Callers must handle that case explicitly.
func authCookieScopes(requestHost string, providerURLs []string) []cookieScope {
	scopes := []cookieScope{}
	seen := map[string]struct{}{}

	addDomain := func(domain string) {
		domain = strings.TrimSpace(strings.ToLower(domain))
		if domain == "" {
			return
		}
		if _, ok := seen[domain]; ok {
			return
		}
		seen[domain] = struct{}{}
		scopes = append(scopes, cookieScope{Domain: domain})
	}

	// (1) Always emit the no-Domain clearance first. This targets the
	// current host exactly and is what the legacy single-scope logout
	// already did — we never want to regress that case.
	scopes = append(scopes, cookieScope{Domain: ""})

	// (2) Current host's eTLD+1.
	if requestHost != "" {
		addDomain(effectiveTLDPlusOne(requestHost))
	}

	// (3) For each configured provider URL, add both the bare host
	// (as a Domain-scoped clearance, which targets that exact host
	// from the browser's perspective) and its eTLD+1. The bare-host
	// entry covers the case where the provider sets a cookie without
	// a Domain attribute on its own host; the eTLD+1 entry covers
	// cookies set with Domain=.example.com on a sibling subdomain.
	for _, raw := range providerURLs {
		host := hostFromURL(raw)
		if host == "" {
			continue
		}
		// Avoid re-adding the current host as a Domain-scoped
		// clearance — the no-Domain entry above already covers it,
		// and emitting Domain=<current-host> on the same response
		// is redundant.
		if !strings.EqualFold(host, requestHost) {
			addDomain(host)
		}
		addDomain(effectiveTLDPlusOne(host))
	}

	return scopes
}

// clearAuthCookies emits a Set-Cookie header for every (cookieName, scope)
// combination, with MaxAge=-1 to request immediate browser-side deletion.
// Callers should pass the scopes returned by authCookieScopes.
func clearAuthCookies(w http.ResponseWriter, cookieNames []string, scopes []cookieScope) {
	for _, name := range cookieNames {
		if strings.TrimSpace(name) == "" {
			continue
		}
		for _, s := range scopes {
			http.SetCookie(w, &http.Cookie{
				Name:     name,
				Value:    "",
				Path:     "/",
				Domain:   s.Domain,
				HttpOnly: true,
				MaxAge:   -1,
			})
		}
	}
}
