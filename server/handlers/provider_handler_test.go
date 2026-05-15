package handlers

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshsync/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestHandler returns a minimal Handler wired with the given providers and provider override.
func newTestHandler(t *testing.T, providers map[string]models.Provider, providerOverride string) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return &Handler{
		config: &models.HandlerConfig{
			Providers:          providers,
			ProviderCookieName: "meshery-provider",
		},
		Provider: providerOverride,
		log:      log,
	}
}

func TestProviderHandler_NoneRedirectsThroughLogin(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/api/provider?provider=None", nil)
	rec := httptest.NewRecorder()

	h.ProviderHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/user/login?provider=None" {
		t.Errorf("expected redirect to /user/login?provider=None, got %s", loc)
	}

	// Verify provider cookie is set
	var found bool
	for _, c := range resp.Cookies() {
		if c.Name == "meshery-provider" && c.Value == "None" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected meshery-provider cookie to be set to None")
	}
}

func TestProviderUIHandler_NoneEnvVarRedirectsThroughLogin(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "None")

	req := httptest.NewRequest(http.MethodGet, "/provider", nil)
	rec := httptest.NewRecorder()

	h.ProviderUIHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/user/login" {
		t.Errorf("expected redirect to /user/login, got %s", loc)
	}

	// Regression: previous version of this test did not assert on the cookie
	// value, which let the empty-cookie loop trigger ship. The
	// auto-select branch must always set a non-empty cookie naming a
	// registered provider.
	assertProviderCookieIsRegistered(t, resp, "meshery-provider", providers)
}

// assertProviderCookieIsRegistered fails the test if any meshery-provider
// cookie in the response has an empty value or names a provider that is
// not in the registered providers map. Empty / unregistered cookies are
// the trigger for the /user/login ⇄ /provider redirect loop documented in
// ProviderUIHandler.
func assertProviderCookieIsRegistered(t *testing.T, resp *http.Response, cookieName string, providers map[string]models.Provider) {
	t.Helper()
	for _, c := range resp.Cookies() {
		if c.Name != cookieName {
			continue
		}
		// LogoutHandler-style cookie deletions (MaxAge<0) are intentional
		// empties — those are not the bug we are guarding against.
		if c.MaxAge < 0 {
			continue
		}
		if c.Value == "" {
			t.Fatalf("ProviderUIHandler emitted empty %s cookie — this is a redirect-loop trigger", cookieName)
		}
		if _, ok := providers[c.Value]; !ok {
			t.Fatalf("ProviderUIHandler emitted %s=%q which is not a registered provider — this would loop via ProviderMiddleware nil-provider path", cookieName, c.Value)
		}
	}
}

// assertNoProviderCookieSet fails the test if any non-deletion meshery-provider
// cookie was written. Used for the degraded-fallthrough paths where we MUST
// reach ServeUI without setting the cookie.
func assertNoProviderCookieSet(t *testing.T, resp *http.Response, cookieName string) {
	t.Helper()
	for _, c := range resp.Cookies() {
		if c.Name == cookieName && c.MaxAge >= 0 {
			t.Fatalf("ProviderUIHandler set %s=%q on a fallthrough path — must serve provider chooser instead of writing a cookie that would close the redirect loop", cookieName, c.Value)
		}
	}
}

// TestProviderUIHandler_LoopGuards covers the (PlaygroundBuild, h.Provider,
// providerRegistered) combinations that decide whether ProviderUIHandler may
// auto-select and redirect, or must degrade to the provider-selection UI.
//
// Meshery extension image outages are exactly the (PlaygroundBuild=true, h.Provider="",
// any registry) row: the handler wrote an empty meshery-provider cookie,
// which ProviderMiddleware rightly ignored, which made AuthMiddleware bounce
// the request back to /provider — infinite loop. Guarding that row, plus the
// adjacent "h.Provider names something not in the registry" row, is the
// whole point of this matrix.
func TestProviderUIHandler_LoopGuards(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	registered := map[string]models.Provider{local.Name(): local} // local.Name() == "None"

	tests := []struct {
		name             string
		playgroundBuild  bool
		providerOverride string
		providers        map[string]models.Provider
		wantRedirect     bool   // true means 302 to /user/login + cookie set
		wantCookieValue  string // empty means assert no cookie set
	}{
		{
			name:             "production playground regression: PLAYGROUND=true, PROVIDER unset → must NOT loop",
			playgroundBuild:  true,
			providerOverride: "",
			providers:        registered,
			wantRedirect:     false,
			wantCookieValue:  "",
		},
		{
			name:             "PLAYGROUND=false, PROVIDER unset → standard multi-provider chooser, no cookie",
			playgroundBuild:  false,
			providerOverride: "",
			providers:        registered,
			wantRedirect:     false,
			wantCookieValue:  "",
		},
		{
			name:             "PLAYGROUND=true, PROVIDER set + registered → auto-select happy path",
			playgroundBuild:  true,
			providerOverride: local.Name(),
			providers:        registered,
			wantRedirect:     true,
			wantCookieValue:  local.Name(),
		},
		{
			name:             "PLAYGROUND=false, PROVIDER set + registered → enforced-provider auto-select",
			playgroundBuild:  false,
			providerOverride: local.Name(),
			providers:        registered,
			wantRedirect:     true,
			wantCookieValue:  local.Name(),
		},
		{
			name:             "PROVIDER names an unregistered provider → must NOT set bogus cookie that would loop",
			playgroundBuild:  false,
			providerOverride: "Meshery",
			providers:        registered, // does not contain "Meshery"
			wantRedirect:     false,
			wantCookieValue:  "",
		},
		{
			name:             "PLAYGROUND=true + PROVIDER unregistered → must degrade, not loop",
			playgroundBuild:  true,
			providerOverride: "Meshery",
			providers:        registered, // does not contain "Meshery"
			wantRedirect:     false,
			wantCookieValue:  "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h := newTestHandler(t, tc.providers, tc.providerOverride)
			h.config.PlaygroundBuild = tc.playgroundBuild

			req := httptest.NewRequest(http.MethodGet, "/provider", nil)
			rec := httptest.NewRecorder()

			h.ProviderUIHandler(rec, req)

			resp := rec.Result()
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			if tc.wantRedirect {
				if resp.StatusCode != http.StatusFound {
					t.Fatalf("want 302 Found, got %d", resp.StatusCode)
				}
				if loc := resp.Header.Get("Location"); loc != "/user/login" {
					t.Errorf("want redirect to /user/login, got %q", loc)
				}
				assertProviderCookieIsRegistered(t, resp, "meshery-provider", tc.providers)

				var got string
				for _, c := range resp.Cookies() {
					if c.Name == "meshery-provider" {
						got = c.Value
					}
				}
				if got != tc.wantCookieValue {
					t.Errorf("want cookie value %q, got %q", tc.wantCookieValue, got)
				}
				return
			}

			// Fallthrough paths: must NOT 302, must NOT set the cookie. The
			// status code itself depends on whether ServeUI finds the static
			// file under ../../provider-ui/out/ — in unit tests we don't
			// require a built provider-ui, so we accept any non-302 response.
			if resp.StatusCode == http.StatusFound {
				t.Fatalf("ProviderUIHandler must not redirect on fallthrough (got 302 to %q) — this is the loop trigger", resp.Header.Get("Location"))
			}
			assertNoProviderCookieSet(t, resp, "meshery-provider")
		})
	}
}

// TestProviderUIHandler_PreservesQueryStringOnAutoSelect makes sure deep-link
// query params (e.g. ?ref=<base64>) survive the auto-select redirect, which
// is the only path that takes the user to /user/login through this handler.
func TestProviderUIHandler_PreservesQueryStringOnAutoSelect(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, local.Name())

	req := httptest.NewRequest(http.MethodGet, "/provider?ref=L2V4dGVuc2lvbi9tZXNobWFw", nil)
	rec := httptest.NewRecorder()
	h.ProviderUIHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("want 302 Found, got %d", resp.StatusCode)
	}
	wantLoc := "/user/login?ref=L2V4dGVuc2lvbi9tZXNobWFw"
	if loc := resp.Header.Get("Location"); loc != wantLoc {
		t.Errorf("want %q, got %q", wantLoc, loc)
	}
}

func TestLoginHandler_NoneProviderRedirectsHome(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None", nil)
	rec := httptest.NewRecorder()

	h.LoginHandler(rec, req, local, false)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/" {
		t.Errorf("expected redirect to /, got %s", loc)
	}
}

func TestLoginHandler_NoneProviderDeepLink(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	// ref is base64-raw-url-encoded "/extension/meshmap"
	req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None&ref=L2V4dGVuc2lvbi9tZXNobWFw", nil)
	rec := httptest.NewRecorder()

	h.LoginHandler(rec, req, local, false)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/extension/meshmap" {
		t.Errorf("expected redirect to /extension/meshmap, got %s", loc)
	}
}

func TestLoginHandler_NoneProviderRejectsOpenRedirect(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	cases := []struct {
		name string
		ref  string
	}{
		{"absolute URL", base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/phish"))},
		{"protocol-relative", base64.RawURLEncoding.EncodeToString([]byte("//evil.example/phish"))},
		{"javascript scheme", base64.RawURLEncoding.EncodeToString([]byte("javascript:alert(1)"))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None&ref="+tc.ref, nil)
			rec := httptest.NewRecorder()

			h.LoginHandler(rec, req, local, false)

			resp := rec.Result()
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			loc := resp.Header.Get("Location")
			if loc != "/" {
				t.Errorf("expected redirect to / for unsafe ref, got %s", loc)
			}
		})
	}
}

func TestLoginHandler_NoneProviderFromMiddleware(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/api/some/endpoint", nil)
	rec := httptest.NewRecorder()

	// fromMiddleWare=true: simulates AuthMiddleware intercepting an API request
	h.LoginHandler(rec, req, local, true)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	// Should be a no-op pass-through (200 with no body), not a redirect
	if resp.StatusCode == http.StatusFound {
		t.Fatalf("expected no redirect from middleware path, got 302 to %s", resp.Header.Get("Location"))
	}
}

func TestDeleteMeshSyncResource(t *testing.T) {
	handler := newTestHandler(t, map[string]models.Provider{}, "")

		tests := []struct {
			name                 string
			migrateResourceTable bool
			expectedStatus       int
			expectedContentType  string
			expectedBody         string
		}{
			{
				name: "given resource table migrated when DeleteMeshSyncResource then return status 200 and deleted true",
				migrateResourceTable: true,
				expectedStatus:       http.StatusOK,
				expectedContentType:  "application/json",
				expectedBody:         `"deleted":true`,
			},
			{
				name: "given resource table not migrated when DeleteMeshSyncResource then return status 500 and failed to delete",
				migrateResourceTable: false,
				expectedStatus:       http.StatusInternalServerError,
				expectedContentType:  "application/json",
				// Error response now carries the MeshKit envelope; match on the
				// error code emitted by ErrFailToDelete so the assertion is
				// stable across short-description wording changes.
				expectedBody:         "meshery-server-",
			},
		}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to create in-memory database: %v", err)
			}

			if tt.migrateResourceTable {
				if err := db.AutoMigrate(&model.KubernetesResource{}); err != nil {
					t.Fatalf("failed to migrate kubernetes resource table: %v", err)
				}
			}

			provider := &models.DefaultLocalProvider{}
			provider.Initialize()
			provider.GenericPersister = &database.Handler{DB: db}

			req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/resource-1", nil)
			req = mux.SetURLVars(req, map[string]string{"id": "resource-1"})
			rw := httptest.NewRecorder()

			handler.DeleteMeshSyncResource(rw, req, nil, nil, provider)

			if rw.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rw.Code)
			}

			if contentType := rw.Header().Get("Content-Type"); !strings.HasPrefix(contentType, tt.expectedContentType) {
				t.Fatalf("expected content type prefix %q, got %q", tt.expectedContentType, contentType)
			}

			// Success path decodes the deleted envelope; error path asserts
			// the MeshKit JSON error envelope contains the expected substring.
			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Deleted bool `json:"deleted"`
				}

				if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode json response: %v", err)
				}

				if !response.Deleted {
					t.Fatal("expected 'deleted' field to be true")
				}
			} else {
				if body := rw.Body.String(); !strings.Contains(body, tt.expectedBody) {
					t.Fatalf("expected body to contain %q, got %q", tt.expectedBody, body)
				}
			}
		})
	}
}
