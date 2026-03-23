package models

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/meshery/meshkit/logger"
	"github.com/spf13/viper"
)

func newRemoteProviderTestLogger(t *testing.T) logger.Handler {
	t.Helper()

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	return log
}

func newSignedJWT(t *testing.T, key *rsa.PrivateKey, kid string) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "user-1",
	})
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("failed to sign JWT: %v", err)
	}

	return tokenString
}

func publicJWK(key *rsa.PrivateKey, kid string) map[string]string {
	return map[string]string{
		"kid": kid,
		"n":   base64.RawURLEncoding.EncodeToString(key.PublicKey.N.Bytes()),
		"e":   "AQAB",
	}
}

func writeCapabilitiesFixture(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "capabilities.json")
	payload := ProviderProperties{
		ProviderType: RemoteProviderType,
		ProviderName: "Meshery Cloud",
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal capabilities fixture: %v", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("failed to write capabilities fixture: %v", err)
	}

	return path
}

func TestRemoteProviderGetSessionUsesStoredRefreshedToken(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	const (
		kid        = "test-key"
		staleToken = "stale-token"
	)

	introspectCalled := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/introspect" {
			http.NotFound(w, r)
			return
		}
		introspectCalled = true
		w.WriteHeader(http.StatusUnauthorized)
	}))
	t.Cleanup(server.Close)

	provider := &RemoteProvider{
		RemoteProviderURL: server.URL,
		TokenStore: map[string]string{
			staleToken: newSignedJWT(t, key, kid),
		},
		Keys: []map[string]string{
			publicJWK(key, kid),
		},
		Log: newRemoteProviderTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/system/events", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: staleToken})

	if err := provider.GetSession(req); err != nil {
		t.Fatalf("expected session validation to succeed using refreshed token, got: %v", err)
	}

	if introspectCalled {
		t.Fatal("expected GetSession to validate the refreshed token from store without introspecting the stale token")
	}
}

func TestRemoteProviderTokenHandlerSetsLatestTokenCookie(t *testing.T) {
	capabilitiesPath := writeCapabilitiesFixture(t)
	viper.Set(PROVIDER_CAPABILITIES_FILEPATH_ENV, capabilitiesPath)
	t.Cleanup(func() {
		viper.Set(PROVIDER_CAPABILITIES_FILEPATH_ENV, "")
	})

	const (
		staleToken   = "stale-token"
		refreshedJWT = "refreshed-token"
		sessionValue = "provider-session"
	)

	provider := &RemoteProvider{
		RemoteProviderURL: "http://provider.example",
		TokenStore: map[string]string{
			staleToken: refreshedJWT,
		},
		CookieDuration: time.Hour,
		Log:            newRemoteProviderTestLogger(t),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user/token?token="+staleToken+"&session_cookie="+sessionValue, nil)
	req = req.WithContext(context.WithValue(req.Context(), MesheryServerURL, "http://meshery.example"))
	rec := httptest.NewRecorder()

	provider.TokenHandler(rec, req, false)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected redirect after token handling, got %d", resp.StatusCode)
	}

	var tokenCookie *http.Cookie
	var sessionCookie *http.Cookie
	for _, cookie := range resp.Cookies() {
		switch cookie.Name {
		case TokenCookieName:
			tokenCookie = cookie
		case ProviderSessionCookieName:
			sessionCookie = cookie
		}
	}

	if tokenCookie == nil {
		t.Fatal("expected token cookie to be set")
	}
	if tokenCookie.Value != refreshedJWT {
		t.Fatalf("expected token cookie %q, got %q", refreshedJWT, tokenCookie.Value)
	}
	if sessionCookie == nil {
		t.Fatal("expected provider session cookie to be set")
	}
	if sessionCookie.Value != sessionValue {
		t.Fatalf("expected provider session cookie %q, got %q", sessionValue, sessionCookie.Value)
	}
}
