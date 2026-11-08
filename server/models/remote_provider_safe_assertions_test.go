package models

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRemoteProvider_InitiateLogin_MissingCallbackURL asserts that calling InitiateLogin
// when MesheryServerCallbackURL is missing from the request context does not trigger
// an unchecked type assertion panic, but safely logs the error and redirects to /error.
func TestRemoteProvider_InitiateLogin_MissingCallbackURL(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://localhost:9876")
	provider.RefCookieName = "test_ref"
	provider.LoginCookieDuration = time.Hour

	req := httptest.NewRequest(http.MethodGet, "http://localhost:9081/user/login", nil)
	rec := httptest.NewRecorder()

	provider.InitiateLogin(rec, req, false)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "/error" {
		t.Errorf("expected redirect to /error, got %s", loc)
	}
}

// TestRemoteProvider_GetProviderCapabilities_MissingToken asserts that calling
// GetProviderCapabilities without TokenCtxKey in context does not panic, but returns 401 Unauthorized.
func TestRemoteProvider_GetProviderCapabilities_MissingToken(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://localhost:9876")

	req := httptest.NewRequest(http.MethodGet, "http://localhost:9081/api/provider/capabilities", nil)
	rec := httptest.NewRecorder()

	provider.GetProviderCapabilities(rec, req, "test-user-id")

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized on missing token, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

// TestRemoteProvider_TokenHandler_MissingServerURL asserts that TokenHandler does not panic
// when MesheryServerURL is absent from the request context, and that spawned goroutines
// are safely protected.
func TestRemoteProvider_TokenHandler_MissingServerURL(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ProviderProperties{
			ProviderType: RemoteProviderType,
			ProviderName: "test-remote",
			Capabilities: Capabilities{},
		})
	}))
	defer mockServer.Close()

	provider := newTestRemoteProvider(t, mockServer.URL)
	provider.RefCookieName = "test_ref"

	req := httptest.NewRequest(http.MethodGet, "http://localhost:9081/api/user/token?meshery-provider-token=test-token", nil)
	rec := httptest.NewRecorder()

	provider.TokenHandler(rec, req, false)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status 302 redirect, got %d", rec.Code)
	}

	// Give background goroutines time to execute and ensure no panic
	time.Sleep(50 * time.Millisecond)
}

// TestRemoteProvider_InterceptAnonymousSession_MissingServerURL asserts that
// InterceptLoginAndInitiateAnonymousUserSession safely handles missing MesheryServerURL
// in request context without panicking.
func TestRemoteProvider_InterceptAnonymousSession_MissingServerURL(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://localhost:9876")
	provider.ProviderProperties = ProviderProperties{
		Capabilities: Capabilities{
			{
				Feature:  PersistAnonymousUser,
				Endpoint: "/api/user/anonymous",
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "http://localhost:9081/user/login", nil)
	rec := httptest.NewRecorder()

	provider.InterceptLoginAndInitiateAnonymousUserSession(req, rec)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected status 302 redirect, got %d", rec.Code)
	}
}
