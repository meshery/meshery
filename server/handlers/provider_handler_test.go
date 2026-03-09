package handlers

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/user/login" {
		t.Errorf("expected redirect to /user/login, got %s", loc)
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
	defer resp.Body.Close()

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
	defer resp.Body.Close()

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
			defer resp.Body.Close()

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
	defer resp.Body.Close()

	// Should be a no-op pass-through (200 with no body), not a redirect
	if resp.StatusCode == http.StatusFound {
		t.Fatalf("expected no redirect from middleware path, got 302 to %s", resp.Header.Get("Location"))
	}
}
