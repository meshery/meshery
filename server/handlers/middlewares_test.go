package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
)

// TestResolveProviderName covers the precedence rules for picking the
// provider name from a request, plus the enforced-default fallback that
// breaks the /user/login ⇄ /provider redirect loop on enforced-provider
// deployments where the cookie isn't propagating back (SameSite/popup/CDN).
func TestResolveProviderName(t *testing.T) {
	const cookieName = "meshery-provider"

	mkReq := func(setup func(*http.Request)) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		if setup != nil {
			setup(req)
		}
		return req
	}

	tests := []struct {
		name             string
		setup            func(*http.Request)
		enforcedProvider string
		want             string
	}{
		{
			name:             "no signals, no enforced default → empty",
			setup:            nil,
			enforcedProvider: "",
			want:             "",
		},
		{
			name:             "no signals, enforced default returned",
			setup:            nil,
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over cookie",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "Meshery"})
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over header when cookie empty",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: ""})
				r.Header.Set(cookieName, "Meshery")
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over header and query",
			setup: func(r *http.Request) {
				r.Header.Set(cookieName, "Meshery")
				r.URL.RawQuery = "provider=None"
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "enforced default wins over query (None alias does not select Local)",
			setup: func(r *http.Request) {
				r.URL.RawQuery = "provider=None"
			},
			enforcedProvider: "Enforced",
			want:             "Enforced",
		},
		{
			name: "cookie with whitespace value is honored verbatim",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "  Meshery  "})
			},
			enforcedProvider: "",
			want:             "  Meshery  ",
		},
		// Backward-compat guard for the "None" -> "Local" rename. Stale
		// browser cookies, hand-written scripts, and older mesheryctl all
		// keep sending "None" — the request edge must accept it.
		{
			name: "legacy None cookie normalizes to Local",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "None"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name:             "legacy None enforced default normalizes to Local",
			setup:            nil,
			enforcedProvider: "None",
			want:             "Local",
		},
		{
			name: "alias matching is case-insensitive",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "none"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name: "lowercase canonical name normalizes to canonical Local",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "local"})
			},
			enforcedProvider: "",
			want:             "Local",
		},
		{
			name:             "uppercase canonical name normalizes to canonical Local",
			setup:            nil,
			enforcedProvider: "LOCAL",
			want:             "Local",
		},
		{
			name: "non-local provider name is returned verbatim (remote names not normalized)",
			setup: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: cookieName, Value: "Meshery"})
			},
			enforcedProvider: "",
			want:             "Meshery",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, _ := resolveProviderName(mkReq(tc.setup), cookieName, tc.enforcedProvider)
			if got != tc.want {
				t.Fatalf("resolveProviderName: want %q, got %q", tc.want, got)
			}
		})
	}
}

// TestResolveProviderName_ReturnsLegacyAliasFlag locks in the second return
// value used by ProviderMiddleware to drive the one-shot deprecation warning.
// The boolean must be true for any casing of the legacy alias and false for
// canonical names (including the canonical local provider name).
func TestResolveProviderName_ReturnsLegacyAliasFlag(t *testing.T) {
	const cookieName = "meshery-provider"
	mkReq := func(cookieValue string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
		if cookieValue != "" {
			req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
		}
		return req
	}

	cases := []struct {
		name     string
		cookie   string
		wantFlag bool
	}{
		{"canonical Local does not flag legacy", "Local", false},
		{"lowercase local does not flag legacy", "local", false},
		{"legacy None flags legacy", "None", true},
		{"lowercase none flags legacy", "none", true},
		{"uppercase NONE flags legacy", "NONE", true},
		{"remote Meshery does not flag legacy", "Meshery", false},
		{"empty does not flag legacy", "", false},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			_, got := resolveProviderName(mkReq(c.cookie), cookieName, "")
			if got != c.wantFlag {
				t.Errorf("usedLegacyAlias = %v, want %v", got, c.wantFlag)
			}
		})
	}
}

func TestProviderMiddlewareBlocksCrossOriginLocalMutations(t *testing.T) {
	localProvider := &models.DefaultLocalProvider{
		ProviderProperties: models.ProviderProperties{
			ProviderName: models.LocalProviderName,
			ProviderType: models.LocalProviderType,
		},
	}
	h := &Handler{
		config: &models.HandlerConfig{
			ProviderCookieName: "meshery-provider",
			Providers:          map[string]models.Provider{models.LocalProviderName: localProvider},
		},
	}

	tests := []struct {
		name       string
		method     string
		target     string
		origin     string
		referer    string
		fetchSite  string
		headers    http.Header
		uiOrigins  string
		wantStatus int
	}{
		{
			name:       "cross-site form submission with Local query selector is rejected",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "http://attacker.example",
			fetchSite:  "cross-site",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "mismatched Origin is rejected without Fetch Metadata",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "http://attacker.example",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "mismatched Referer is rejected when Origin is absent",
			method:     http.MethodDelete,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			referer:    "http://attacker.example/submit",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "opaque Origin is rejected",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "null",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "same-origin mutation is allowed",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "http://localhost:9081",
			fetchSite:  "same-origin",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "unconfigured localhost UI port cannot mutate localhost API",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "http://localhost:3000",
			fetchSite:  "same-site",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "unconfigured named host is not trusted as the server origin",
			method:     http.MethodPost,
			target:     "http://attacker.example:9081/api/system/database/reset?provider=Local",
			origin:     "http://attacker.example:9081",
			fetchSite:  "same-origin",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "non-loopback origin must match the full host and port",
			method:     http.MethodPost,
			target:     "http://meshery.example.com:9081/api/system/database/reset?provider=Local",
			origin:     "http://meshery.example.com:3000",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "TLS reverse proxy origin is allowed when explicitly configured",
			method:     http.MethodPost,
			target:     "http://meshery.example.com:9081/api/system/database/reset?provider=Local",
			origin:     "https://ui.meshery.example.com",
			uiOrigins:  "https://ui.meshery.example.com",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "spoofed X-Forwarded-Proto cannot change the target origin",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "https://localhost:9081",
			fetchSite:  "same-origin",
			headers:    http.Header{"X-Forwarded-Proto": []string{"https"}},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "cross-site GET to logout is rejected",
			method:     http.MethodGet,
			target:     "http://localhost:9081/user/logout?provider=Local",
			origin:     "http://attacker.example",
			fetchSite:  "cross-site",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "configured localhost UI origin can mutate localhost API",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "http://localhost:3000",
			fetchSite:  "same-site",
			uiOrigins:  "http://localhost:3000",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "configured UI origin is allowed with cross-site metadata",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			origin:     "https://ui.example.com",
			fetchSite:  "cross-site",
			uiOrigins:  "https://ui.example.com",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "wildcard UI origin does not allow a different host",
			method:     http.MethodPost,
			target:     "http://meshery.example.com:9081/api/system/database/reset?provider=Local",
			origin:     "http://ui.example.com",
			uiOrigins:  "http://*.example.com",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "same-origin GET remains allowed",
			method:     http.MethodGet,
			target:     "http://localhost:9081/user/logout?provider=Local",
			origin:     "http://localhost:9081",
			fetchSite:  "same-origin",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "headerless GET login remains allowed",
			method:     http.MethodGet,
			target:     "http://localhost:9081/user/login?provider=Local",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "headerless GET API remains allowed",
			method:     http.MethodGet,
			target:     "http://localhost:9081/api/system/database?provider=Local",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "state-changing session sync GET requires origin metadata",
			method:     http.MethodGet,
			target:     "http://localhost:9081/api/system/sync?provider=Local",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "state-changing performance run GET rejects a foreign origin",
			method:     http.MethodGet,
			target:     "http://localhost:9081/api/user/performance/profiles/123/run?provider=Local",
			origin:     "http://attacker.example",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "profile execution GET accepts same-origin UI request",
			method:     http.MethodGet,
			target:     "http://localhost:9081/api/perf/profile?provider=Local",
			origin:     "http://localhost:9081",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "request without Origin or Referer is rejected",
			method:     http.MethodGet,
			target:     "http://localhost:9081/user/logout?provider=Local",
			fetchSite:  "same-site",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "request without browser origin metadata is rejected",
			method:     http.MethodGet,
			target:     "http://localhost:9081/user/logout?provider=Local",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "headerless POST is rejected",
			method:     http.MethodPost,
			target:     "http://localhost:9081/api/system/database/reset?provider=Local",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h.config.AllowedUIOrigins = tc.uiOrigins
			req := httptest.NewRequest(tc.method, tc.target, nil)
			if tc.origin != "" {
				req.Header.Set("Origin", tc.origin)
			}
			if tc.referer != "" {
				req.Header.Set("Referer", tc.referer)
			}
			if tc.fetchSite != "" {
				req.Header.Set("Sec-Fetch-Site", tc.fetchSite)
			}
			for name, values := range tc.headers {
				for _, value := range values {
					req.Header.Add(name, value)
				}
			}

			nextCalled := false
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusNoContent)
			})
			response := httptest.NewRecorder()
			h.ProviderMiddleware(next).ServeHTTP(response, req)

			if response.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body = %q", response.Code, tc.wantStatus, response.Body.String())
			}
			if tc.wantStatus == http.StatusForbidden && response.Header().Get("Content-Type") != "application/json; charset=utf-8" {
				t.Errorf("Content-Type = %q, want application/json; charset=utf-8", response.Header().Get("Content-Type"))
			}
			if got := nextCalled; got != (tc.wantStatus == http.StatusNoContent) {
				t.Errorf("nextCalled = %v, want %v", got, tc.wantStatus == http.StatusNoContent)
			}
		})
	}
}

func TestAuthMiddleWare(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode.")
	}

	t.Log("Need to run AuthMiddleware() skipping")
	//_, err := AuthMiddleware(nil, nil)
	//if err != nil {
	//	t.Errorf("AuthMiddleWare() failed with error: %s", err)
	//}
}

// TestWriteMeshkitError_ErrTransientProvider is a focused test for the
// transient-provider response shape emitted by SessionInjectorMiddleware when
// Meshery Cloud is temporarily unreachable. A full middleware-level test would
// require stubbing the 100+ method models.Provider interface, which the
// package doesn't yet have. Instead, we test the migration end-to-end at the
// response-helper level: the same two lines the middleware now executes.
//
// If a provider mocking pattern lands in this package later, promote this to
// a middleware-dispatch test by wiring a stub provider whose GetUserDetails
// returns "Could not reach remote provider".
func TestWriteMeshkitError_ErrTransientProvider(t *testing.T) {
	rec := httptest.NewRecorder()

	writeMeshkitError(
		rec,
		ErrTransientProvider(errors.New("Could not reach remote provider: dial tcp: i/o timeout")),
		http.StatusServiceUnavailable,
	)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}
	if nosniff := resp.Header.Get("X-Content-Type-Options"); nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got %q", nosniff)
	}

	var decoded struct {
		Error                string   `json:"error"`
		Code                 string   `json:"code"`
		Severity             string   `json:"severity"`
		ProbableCause        []string `json:"probableCause"`
		SuggestedRemediation []string `json:"suggestedRemediation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		t.Fatalf("body did not parse as JSON: %v", err)
	}
	if decoded.Code != ErrTransientProviderCode {
		t.Errorf("expected code %q, got %q", ErrTransientProviderCode, decoded.Code)
	}
	if decoded.Error == "" {
		t.Errorf("expected non-empty error message")
	}
	if len(decoded.SuggestedRemediation) == 0 {
		t.Errorf("expected suggestedRemediation to be populated for transient-provider errors")
	}
}

// TestIsTransientProviderError exercises the detector used by
// SessionInjectorMiddleware to decide between a 401 auth failure and a 503
// transient-provider failure. Regression guard in case the error-message
// substrings the detector relies on drift.
func TestIsTransientProviderError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil is not transient", nil, false},
		{"generic auth failure is not transient", fmt.Errorf("token expired"), false},
		{"Could not reach remote provider is transient", fmt.Errorf("Could not reach remote provider: dial tcp: i/o timeout"), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isTransientProviderError(tc.err); got != tc.want {
				t.Errorf("isTransientProviderError(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
