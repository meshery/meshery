package models

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/meshery/meshkit/logger"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

func newTestRemoteProvider(t *testing.T, remoteProviderURL string) *RemoteProvider {
	t.Helper()

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	return &RemoteProvider{
		RemoteProviderURL: remoteProviderURL,
		TokenStore:        make(map[string]string),
		Log:               log,
	}
}

func encodeTestToken(t *testing.T, token oauth2.Token) string {
	t.Helper()

	data, err := json.Marshal(token)
	if err != nil {
		t.Fatalf("failed to marshal test token: %v", err)
	}

	return base64.RawStdEncoding.EncodeToString(data)
}

func setViperKeyForRemoteAuthTest(t *testing.T, key string, value any) {
	t.Helper()

	hadPrevious := viper.IsSet(key)
	previous := viper.Get(key)
	viper.Set(key, value)

	t.Cleanup(func() {
		if hadPrevious {
			viper.Set(key, previous)
			return
		}
		viper.Set(key, nil)
	})
}

func newRemoteLoginRequest(t *testing.T, rawURL string) *http.Request {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, rawURL, nil)
	ctx := context.WithValue(req.Context(), MesheryServerCallbackURL, "http://localhost:9081/api/user/token")
	return req.WithContext(ctx)
}

func redirectQuery(t *testing.T, location string) url.Values {
	t.Helper()

	redirectURL, err := url.Parse(location)
	if err != nil {
		t.Fatalf("parse redirect location %q: %v", location, err)
	}

	return redirectURL.Query()
}

func TestRemoteProviderInitiateLogin_OmitsPlaceholderVersions(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://localhost:9876")
	provider.RefCookieName = "localhost:9876_ref"
	provider.LoginCookieDuration = time.Hour
	provider.ProviderVersion = "Not Set"

	setViperKeyForRemoteAuthTest(t, "BUILD", "Not Set")

	req := newRemoteLoginRequest(t, "http://localhost:9081/user/login")
	rec := httptest.NewRecorder()

	provider.InitiateLogin(rec, req, false)

	resp := rec.Result()
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected %d, got %d", http.StatusFound, resp.StatusCode)
	}

	query := redirectQuery(t, resp.Header.Get("Location"))
	if got := query.Get("meshery_version"); got != "" {
		t.Fatalf("expected meshery_version to be omitted, got %q", got)
	}
	if got := query.Get("provider_version"); got != "" {
		t.Fatalf("expected provider_version to be omitted, got %q", got)
	}
	if got := query.Get("source"); got == "" {
		t.Fatal("expected source to be present")
	}
}

func TestRemoteProviderInitiateLogin_PropagatesKnownVersions(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://localhost:9876")
	provider.RefCookieName = "localhost:9876_ref"
	provider.LoginCookieDuration = time.Hour
	provider.ProviderVersion = "v0.8.0"

	setViperKeyForRemoteAuthTest(t, "BUILD", "v0.8.0")

	req := newRemoteLoginRequest(t, "http://localhost:9081/user/login")
	rec := httptest.NewRecorder()

	provider.InitiateLogin(rec, req, false)

	resp := rec.Result()
	t.Cleanup(func() {
		_ = resp.Body.Close()
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected %d, got %d", http.StatusFound, resp.StatusCode)
	}

	query := redirectQuery(t, resp.Header.Get("Location"))
	if got := query.Get("meshery_version"); got != "v0.8.0" {
		t.Fatalf("expected meshery_version %q, got %q", "v0.8.0", got)
	}
	if got := query.Get("provider_version"); got != "v0.8.0" {
		t.Fatalf("expected provider_version %q, got %q", "v0.8.0", got)
	}
}

func TestRemoteProviderDoRequest_DoesNotRefreshAnonymousToken(t *testing.T) {
	var resourceCalls int32
	var refreshCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/resource":
			atomic.AddInt32(&resourceCalls, 1)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
		case "/refresh":
			atomic.AddInt32(&refreshCalls, 1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"token":"should-not-be-used"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/resource", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := provider.DoRequest(req, GlobalTokenForAnonymousResults)
	if err != nil {
		t.Fatalf("expected anonymous token request to return original 401 response, got error: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, resp.StatusCode)
	}

	if got := atomic.LoadInt32(&resourceCalls); got != 1 {
		t.Fatalf("expected one request attempt, got %d", got)
	}

	if got := atomic.LoadInt32(&refreshCalls); got != 0 {
		t.Fatalf("expected anonymous token request not to hit /refresh, got %d calls", got)
	}
}

func TestRemoteProviderDoRequest_RejectsNonOKRefreshResponses(t *testing.T) {
	var resourceCalls int32
	var refreshCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/resource":
			atomic.AddInt32(&resourceCalls, 1)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
		case "/refresh":
			atomic.AddInt32(&refreshCalls, 1)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"bad refresh"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token := encodeTestToken(t, oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
	})

	req, err := http.NewRequest(http.MethodPost, server.URL+"/resource", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := provider.DoRequest(req, token)
	if resp != nil {
		t.Fatal("expected DoRequest to fail before retrying the protected request")
	}
	if err == nil {
		t.Fatal("expected refresh failure to return an error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "status code 400") {
		t.Fatalf("expected error to mention refresh status code, got: %v", err)
	}

	if got := atomic.LoadInt32(&resourceCalls); got != 1 {
		t.Fatalf("expected one protected request before refresh failure, got %d", got)
	}

	if got := atomic.LoadInt32(&refreshCalls); got != 1 {
		t.Fatalf("expected exactly one refresh attempt, got %d", got)
	}
}

func TestRemoteProviderDoRequest_RejectsRefreshResponsesWithoutToken(t *testing.T) {
	var resourceCalls int32
	var refreshCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/resource":
			atomic.AddInt32(&resourceCalls, 1)
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte("unauthorized"))
		case "/refresh":
			atomic.AddInt32(&refreshCalls, 1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"error":"missing token field"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token := encodeTestToken(t, oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
	})

	req, err := http.NewRequest(http.MethodPost, server.URL+"/resource", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := provider.DoRequest(req, token)
	if resp != nil {
		t.Fatal("expected DoRequest to stop when refresh response omits token")
	}
	if err == nil {
		t.Fatal("expected missing token in refresh response to return an error")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "missing \"token\"") {
		t.Fatalf("expected error to mention missing token field, got: %v", err)
	}

	if got := atomic.LoadInt32(&resourceCalls); got != 1 {
		t.Fatalf("expected one protected request before refresh validation failed, got %d", got)
	}

	if got := atomic.LoadInt32(&refreshCalls); got != 1 {
		t.Fatalf("expected exactly one refresh attempt, got %d", got)
	}
}

func TestRemoteProviderDoRequest_RetriesWithFreshBodyAfterRefresh(t *testing.T) {
	var resourceCalls int32
	var refreshCalls int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/resource":
			call := atomic.AddInt32(&resourceCalls, 1)
			if call == 1 {
				// Simulate auth middleware rejecting the request before the body is read.
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte("unauthorized"))
				return
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("failed to read retried request body: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if string(body) != `{"kind":"event"}` {
				t.Errorf("expected retried request body to be preserved, got %q", string(body))
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("ok"))
		case "/refresh":
			atomic.AddInt32(&refreshCalls, 1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"token":"refreshed-token"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token := encodeTestToken(t, oauth2.Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		TokenType:    "Bearer",
	})

	req, err := http.NewRequest(http.MethodPost, server.URL+"/resource", bytes.NewBufferString(`{"kind":"event"}`))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	resp, err := provider.DoRequest(req, token)
	if err != nil {
		t.Fatalf("expected retried request to succeed after refresh, got error: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d after retry, got %d", http.StatusOK, resp.StatusCode)
	}

	if got := atomic.LoadInt32(&resourceCalls); got != 2 {
		t.Fatalf("expected the protected endpoint to be called twice, got %d", got)
	}

	if got := atomic.LoadInt32(&refreshCalls); got != 1 {
		t.Fatalf("expected exactly one refresh request, got %d", got)
	}
}

// TestRemoteProviderDoRequest_XAPIKeyAnonymousOnly guards the anonymous
// passphrase contract: X-API-Key MUST only be set when the outbound request
// is being made with the global anonymous token, never with a real user JWT.
// Setting X-API-Key with a user JWT pollutes the cloud's static-token
// fallback path and broke anonymous flows. The test also
// covers the defensive Del path: a caller-supplied X-API-Key on a
// non-anonymous outbound request must be stripped before reaching cloud.
func TestRemoteProviderDoRequest_XAPIKeyAnonymousOnly(t *testing.T) {
	var (
		mu      sync.Mutex
		headers = map[string]string{}
	)
	record := func(label, value string) {
		mu.Lock()
		defer mu.Unlock()
		headers[label] = value
	}
	read := func(label string) string {
		mu.Lock()
		defer mu.Unlock()
		return headers[label]
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/anon":
			record("anon", r.Header.Get("X-API-Key"))
			w.WriteHeader(http.StatusOK)
		case "/user":
			record("user", r.Header.Get("X-API-Key"))
			w.WriteHeader(http.StatusOK)
		case "/user-with-stale-api-key":
			record("stale", r.Header.Get("X-API-Key"))
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)

	assertOK := func(resp *http.Response, label string) {
		t.Helper()
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				t.Logf("%s: error closing response body: %v", label, cerr)
			}
		}()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Fatalf("%s: status %d, body: %s", label, resp.StatusCode, strings.TrimSpace(string(body)))
		}
	}

	anonReq, err := http.NewRequest(http.MethodGet, server.URL+"/anon", nil)
	if err != nil {
		t.Fatalf("failed to build anonymous request: %v", err)
	}
	resp, err := provider.DoRequest(anonReq, GlobalTokenForAnonymousResults)
	if err != nil {
		t.Fatalf("anon request failed: %v", err)
	}
	assertOK(resp, "anon request")

	userToken := encodeTestToken(t, oauth2.Token{AccessToken: "real-user-jwt", TokenType: "Bearer"})
	userReq, err := http.NewRequest(http.MethodGet, server.URL+"/user", nil)
	if err != nil {
		t.Fatalf("failed to build user request: %v", err)
	}
	resp, err = provider.DoRequest(userReq, userToken)
	if err != nil {
		t.Fatalf("user request failed: %v", err)
	}
	assertOK(resp, "user request")

	// Inbound proxy header that should be stripped, not forwarded.
	staleReq, err := http.NewRequest(http.MethodGet, server.URL+"/user-with-stale-api-key", nil)
	if err != nil {
		t.Fatalf("failed to build stale request: %v", err)
	}
	staleReq.Header.Set("X-API-Key", "leaked-from-inbound-request")
	resp, err = provider.DoRequest(staleReq, userToken)
	if err != nil {
		t.Fatalf("stale request failed: %v", err)
	}
	assertOK(resp, "stale request")

	if got := read("anon"); got != GlobalTokenForAnonymousResults {
		t.Fatalf("anonymous request must carry X-API-Key=%q, got %q", GlobalTokenForAnonymousResults, got)
	}
	if got := read("user"); got != "" {
		t.Fatalf("authenticated request must NOT set X-API-Key, got %q", got)
	}
	if got := read("stale"); got != "" {
		t.Fatalf("authenticated request with inbound X-API-Key must strip it, got %q", got)
	}
}

func refCookie(rec *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// The ref cookie is the post-login destination Meshery owns on its own domain,
// and TokenHandler prefers it over the ?ref= a provider sends back. Only the
// Cloud-bound exit returns through TokenHandler to redeem and clear it, so only
// that exit writes one - carrying this attempt's ref, or deleted when there is
// none. The anonymous exit resolves its destination from the request itself and
// leaves nothing behind. A cookie that outlives its attempt shadows the ?ref= a
// later Sign In carries for the whole LoginCookieDuration, whether it was
// written empty or written with an abandoned attempt's destination, and that is
// how mode=design was dropped from the return address.
func TestRemoteProviderInitiateLogin_RefCookieCarriesOnlyThisAttemptsDestination(t *testing.T) {
	const refCookieName = "cloud.layer5.io_ref"
	refValue := base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap?mode=design"))
	staleValue := base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap?mode=abandoned"))

	// l.ProviderProperties is replaced at runtime by TokenHandler, so one
	// deployment does cross between these two shapes from one login to the next
	// - which is why the anonymous exit has to delete a cookie it never writes.
	shapes := []struct {
		name        string
		newProvider func(*testing.T) *RemoteProvider
		keepsRef    bool
	}{
		{
			name: "cloud-bound login",
			newProvider: func(t *testing.T) *RemoteProvider {
				return newTestRemoteProvider(t, "http://localhost:9876")
			},
			keepsRef: true,
		},
		{
			name:        "anonymous session",
			newProvider: newAnonymousSessionProvider,
		},
	}

	tests := []struct {
		name        string
		requestURL  string
		staleCookie string
		expected    string
	}{
		{
			name:       "no ref query deletes the ref cookie",
			requestURL: "http://localhost:9081/user/login",
		},
		{
			name:       "empty ref query deletes the ref cookie",
			requestURL: "http://localhost:9081/user/login?ref=",
		},
		{
			name:        "a login with no ref deletes a stale ref cookie",
			requestURL:  "http://localhost:9081/user/login",
			staleCookie: staleValue,
		},
		{
			name:       "ref query is captured in the cookie",
			requestURL: "http://localhost:9081/user/login?ref=" + refValue,
			expected:   refValue,
		},
		{
			name:        "ref query overwrites a stale ref cookie",
			requestURL:  "http://localhost:9081/user/login?ref=" + refValue,
			staleCookie: staleValue,
			expected:    refValue,
		},
	}

	for _, shape := range shapes {
		shape := shape
		t.Run(shape.name, func(t *testing.T) {
			for _, tc := range tests {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					provider := shape.newProvider(t)
					provider.RefCookieName = refCookieName
					provider.LoginCookieDuration = time.Hour

					req := newRemoteLoginRequest(t, tc.requestURL)
					req = req.WithContext(context.WithValue(req.Context(), MesheryServerURL, "http://localhost:9081"))
					if tc.staleCookie != "" {
						req.AddCookie(&http.Cookie{Name: refCookieName, Value: tc.staleCookie})
					}

					rec := httptest.NewRecorder()
					provider.InitiateLogin(rec, req, false)

					if rec.Code != http.StatusFound {
						t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
					}

					want := tc.expected
					if !shape.keepsRef {
						want = ""
					}

					ck := refCookie(rec, refCookieName)
					if ck == nil {
						t.Fatal("expected a ref cookie header, got none")
					}
					if want == "" {
						if ck.Value != "" || ck.MaxAge >= 0 {
							t.Fatalf("expected the ref cookie to be deleted, got value %q with MaxAge %d", ck.Value, ck.MaxAge)
						}
						return
					}
					if ck.Value != want {
						t.Fatalf("ref cookie = %q, want %q", ck.Value, want)
					}
					if ck.MaxAge < 0 {
						t.Fatalf("ref cookie carrying %q was emitted as a deletion", ck.Value)
					}
				})
			}
		})
	}
}

func generateTestRSAKeyAndJWK(t *testing.T, kid string) (*rsa.PrivateKey, map[string]string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	eBytes := big.NewInt(int64(priv.E)).Bytes()
	jwk := map[string]string{
		"kid": kid,
		"kty": "RSA",
		"alg": "RS256",
		"use": "sig",
		"n":   base64.RawURLEncoding.EncodeToString(priv.N.Bytes()),
		"e":   base64.RawURLEncoding.EncodeToString(eBytes),
	}
	return priv, jwk
}

func signTestJWT(t *testing.T, priv *rsa.PrivateKey, kid string, exp time.Time) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": "user-123",
		"iat": time.Now().Add(-time.Minute).Unix(),
	}
	if !exp.IsZero() {
		claims["exp"] = exp.Unix()
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid
	signed, err := tok.SignedString(priv)
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	return encodeTestToken(t, oauth2.Token{
		AccessToken:  signed,
		RefreshToken: "refresh-123",
		TokenType:    "Bearer",
	})
}

func TestRemoteProviderHTTPClient_SharedPoolAndScopedTimeouts(t *testing.T) {
	origTimeout := remoteAuthRequestTimeout
	remoteAuthRequestTimeout = 50 * time.Millisecond
	t.Cleanup(func() { remoteAuthRequestTimeout = origTimeout })

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys", "/introspect", "/profile":
			time.Sleep(150 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		case "/download":
			time.Sleep(120 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("large-transfer-ok"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: UsersProfile, Endpoint: "/profile"},
	}

	c1 := provider.getHTTPClient()
	c2 := provider.getHTTPClient()
	if c1 == nil {
		t.Fatal("expected non-nil shared HTTP client")
	}
	if c1 != c2 {
		t.Fatal("expected getHTTPClient to return the same shared *http.Client instance")
	}
	if c1.Timeout != 0 {
		t.Fatalf("expected shared http.Client.Timeout to be 0 so transfers are not cut off, got %v", c1.Timeout)
	}

	if err := provider.UpdateJWKs(); err == nil || !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected UpdateJWKs to honor auth timeout (%v), got %v", context.DeadlineExceeded, err)
	}

	token := encodeTestToken(t, oauth2.Token{AccessToken: "tok", TokenType: "Bearer"})
	if err := provider.introspectToken(token); err == nil || !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected introspectToken to honor auth timeout (%v), got %v", context.DeadlineExceeded, err)
	}

	profileReq := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	profileReq.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
	if _, err := provider.GetUserDetails(profileReq); err == nil || !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected GetUserDetails to honor auth timeout (%v), got %v", context.DeadlineExceeded, err)
	}

	downloadReq, err := http.NewRequest(http.MethodGet, server.URL+"/download", nil)
	if err != nil {
		t.Fatalf("build download request: %v", err)
	}
	resp, err := provider.DoRequest(downloadReq, token)
	if err != nil {
		t.Fatalf("expected DoRequest /download not to be cut off by auth timeout, got: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read download body: %v", err)
	}
	if string(body) != "large-transfer-ok" {
		t.Fatalf("expected full download body, got %q", string(body))
	}
}

func TestRemoteProviderVerifyToken_CachesParsedRSAKeyAndDeduplicatesJWKS(t *testing.T) {
	const kid1 = "kid-1"
	const kid2 = "kid-2"
	priv1, jwk1 := generateTestRSAKeyAndJWK(t, kid1)
	priv2, jwk2 := generateTestRSAKeyAndJWK(t, kid2)

	var keysCalls atomic.Int32
	var activeKid atomic.Value
	activeKid.Store(kid1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/keys" {
			http.NotFound(w, r)
			return
		}
		keysCalls.Add(1)
		time.Sleep(15 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		if activeKid.Load().(string) == kid1 {
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk1}})
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk2}})
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token1 := signTestJWT(t, priv1, kid1, time.Now().Add(5*time.Minute))

	const workers = 16
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := provider.VerifyToken(token1); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent VerifyToken failed: %v", err)
	}

	if got := keysCalls.Load(); got != 1 {
		t.Fatalf("expected /keys to be fetched once under concurrent cold miss, got %d", got)
	}

	provider.keysMu.RLock()
	cached1 := provider.parsedKeys[kid1]
	provider.keysMu.RUnlock()
	if cached1 == nil {
		t.Fatal("expected parsed RSA public key to be cached for kid-1")
	}

	// Verify another call reuses the exact same cached *rsa.PublicKey pointer.
	if _, err := provider.VerifyToken(token1); err != nil {
		t.Fatalf("second VerifyToken failed: %v", err)
	}
	provider.keysMu.RLock()
	cached1Again := provider.parsedKeys[kid1]
	provider.keysMu.RUnlock()
	if cached1 != cached1Again {
		t.Fatal("expected VerifyToken to reuse cached *rsa.PublicKey pointer")
	}

	// Rotate key ID on the remote JWKS endpoint and verify refresh works.
	activeKid.Store(kid2)
	token2 := signTestJWT(t, priv2, kid2, time.Now().Add(5*time.Minute))
	if _, err := provider.VerifyToken(token2); err != nil {
		t.Fatalf("VerifyToken after key rotation failed: %v", err)
	}
	if got := keysCalls.Load(); got != 2 {
		t.Fatalf("expected /keys to be fetched again for new kid, got %d", got)
	}
}

func TestRemoteProviderGetSession_SingleflightAndTTLCacheAndLogoutEviction(t *testing.T) {
	const kid = "kid-session"
	priv, jwk := generateTestRSAKeyAndJWK(t, kid)

	var introspectCalls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk}})
		case "/introspect":
			introspectCalls.Add(1)
			time.Sleep(25 * time.Millisecond)
			w.WriteHeader(http.StatusOK)
		case "/logout":
			w.WriteHeader(http.StatusOK)
		case "/revoke":
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token := signTestJWT(t, priv, kid, time.Now().Add(5*time.Minute))

	makeReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)
		req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
		req.AddCookie(&http.Cookie{Name: ProviderSessionCookieName, Value: "sess-1"})
		return req
	}

	// 12 concurrent GetSession calls on cold cache must collapse to 1 /introspect POST.
	const workers = 12
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := provider.GetSession(makeReq()); err != nil {
				errs <- err
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Fatalf("concurrent GetSession failed: %v", err)
	}
	if got := introspectCalls.Load(); got != 1 {
		t.Fatalf("expected 1 /introspect call across %d concurrent requests, got %d", workers, got)
	}

	// Subsequent call within TTL must hit the cache (0 additional /introspect calls).
	if err := provider.GetSession(makeReq()); err != nil {
		t.Fatalf("cached GetSession failed: %v", err)
	}
	if got := introspectCalls.Load(); got != 1 {
		t.Fatalf("expected cached GetSession not to call /introspect again, got %d", got)
	}

	// Logout must evict the cached introspection entry so the next check hits /introspect again.
	rec := httptest.NewRecorder()
	if err := provider.Logout(rec, makeReq()); err != nil {
		t.Fatalf("Logout failed: %v", err)
	}
	if err := provider.GetSession(makeReq()); err != nil {
		t.Fatalf("GetSession after Logout cache eviction failed: %v", err)
	}
	if got := introspectCalls.Load(); got != 2 {
		t.Fatalf("expected GetSession after Logout to re-introspect, got %d calls", got)
	}
}

func TestRemoteProviderGetSession_CapsCacheTTLAtJWTExpiry(t *testing.T) {
	if ttl := authCacheTTLFromExpiry(time.Now().Add(2 * time.Second)); ttl <= 0 || ttl > 2*time.Second {
		t.Fatalf("expected TTL capped at <=2s, got %v", ttl)
	}
	if ttl := authCacheTTLFromExpiry(time.Now().Add(-1 * time.Second)); ttl != 0 {
		t.Fatalf("expected expired token TTL to be 0, got %v", ttl)
	}
	if ttl := authCacheTTLFromExpiry(time.Time{}); ttl != defaultAuthCacheTTL {
		t.Fatalf("expected zero-expiry token to use default TTL %v, got %v", defaultAuthCacheTTL, ttl)
	}
}

func TestRemoteProviderGetUserDetails_SingleflightAndTTLCache(t *testing.T) {
	const kid = "kid-profile"
	priv, _ := generateTestRSAKeyAndJWK(t, kid)
	userID := uuid.Must(uuid.NewV4())

	var profileCalls atomic.Int32
	var failProfile atomic.Bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/profile" {
			http.NotFound(w, r)
			return
		}
		profileCalls.Add(1)
		if failProfile.Load() {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("internal error"))
			return
		}
		time.Sleep(25 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(UserPref{
			User: User{
				ID:        userID,
				FirstName: "Tony",
				LastName:  "Stark",
				Email:     "tony@stark.io",
			},
		})
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: UsersProfile, Endpoint: "/profile"},
	}

	token := signTestJWT(t, priv, kid, time.Now().Add(5*time.Minute))
	makeReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
		return req
	}

	// Failing responses must not be cached.
	failProfile.Store(true)
	if _, err := provider.GetUserDetails(makeReq()); err == nil {
		t.Fatal("expected error when remote /profile returns 500")
	}
	failProfile.Store(false)

	// 12 concurrent calls after recovery must collapse into 1 HTTP GET /profile.
	const workers = 12
	var wg sync.WaitGroup
	users := make(chan *User, workers)
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			u, err := provider.GetUserDetails(makeReq())
			if err != nil {
				errs <- err
				return
			}
			users <- u
		}()
	}
	wg.Wait()
	close(errs)
	close(users)
	for err := range errs {
		t.Fatalf("concurrent GetUserDetails failed: %v", err)
	}
	if got := profileCalls.Load(); got != 2 { // 1 failed + 1 singleflight success
		t.Fatalf("expected 2 total /profile calls (1 failed + 1 singleflight), got %d", got)
	}

	var first *User
	for u := range users {
		if u.ID != userID || u.FirstName != "Tony" {
			t.Fatalf("unexpected user returned: %+v", u)
		}
		if first == nil {
			first = u
		}
	}
	// Mutating returned user must not corrupt the cached copy.
	first.FirstName = "Mutated"
	cachedUser, err := provider.GetUserDetails(makeReq())
	if err != nil {
		t.Fatalf("cached GetUserDetails failed: %v", err)
	}
	if cachedUser.FirstName != "Tony" {
		t.Fatalf("expected cached user FirstName to remain Tony, got %q", cachedUser.FirstName)
	}
	if got := profileCalls.Load(); got != 2 {
		t.Fatalf("expected cached GetUserDetails not to hit /profile again, got %d calls", got)
	}
}

func TestRemoteProviderTokenExpiry_DecodesBase64URLWithDashAndUnderscore(t *testing.T) {
	const kid = "kid-b64url"
	priv, jwk := generateTestRSAKeyAndJWK(t, kid)

	expTime := time.Now().Add(3 * time.Second).Truncate(time.Second)
	claims := jwt.MapClaims{
		"sub": "user-b64url",
		"n1":  "~~~",
		"n2":  "~~~",
		"n3":  "~~~",
		"q1":  "???",
		"q2":  "???",
		"q3":  "???", // guarantees both '-' and '_' in Base64URL-encoded claims segment
		"exp": expTime.Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = kid
	signed, err := tok.SignedString(priv)
	if err != nil {
		t.Fatalf("sign jwt: %v", err)
	}
	parts := strings.Split(signed, ".")
	if len(parts) != 3 || !strings.Contains(parts[1], "-") || !strings.Contains(parts[1], "_") {
		t.Fatalf("test precondition failed: expected Base64URL claims segment with both '-' and '_', got %q", parts[1])
	}

	wrappedToken := encodeTestToken(t, oauth2.Token{
		AccessToken:  signed,
		RefreshToken: "refresh-b64url",
		TokenType:    "Bearer",
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk}})
		case "/introspect":
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)

	gotExp := provider.tokenExpiry(wrappedToken)
	if !gotExp.Equal(expTime) {
		t.Fatalf("tokenExpiry() = %v, want %v (RawStdEncoding would fail on '-'/'_')", gotExp, expTime)
	}

	ttl := provider.authCacheTTLForToken(wrappedToken)
	if ttl <= 0 || ttl > 3*time.Second {
		t.Fatalf("expected TTL capped at <= 3s by JWT exp, got %v", ttl)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: wrappedToken})
	if err := provider.GetSession(req); err != nil {
		t.Fatalf("GetSession failed for Base64URL '-'/'_' payload: %v", err)
	}

	cacheKey := tokenCacheKey(wrappedToken)
	provider.authCacheMu.RLock()
	cachedUntil := provider.introspectCache[cacheKey]
	provider.authCacheMu.RUnlock()
	if cachedUntil.IsZero() || cachedUntil.After(expTime.Add(time.Second)) {
		t.Fatalf("expected introspectCache expiry <= %v, got %v", expTime, cachedUntil)
	}
}

func TestRemoteProviderGetParsedKey_SameKidRotationDoesNotPersistStaleKey(t *testing.T) {
	const sharedKid = "kid-shared"
	priv1, jwk1 := generateTestRSAKeyAndJWK(t, sharedKid)
	priv2, jwk2 := generateTestRSAKeyAndJWK(t, sharedKid)

	provider := newTestRemoteProvider(t, "http://localhost:9876")
	provider.Keys = []map[string]string{jwk2}

	// Simulate an in-flight goroutine attempting to store the old parsed key
	// after UpdateJWKs rotated key material under the same kid.
	provider.storeParsedKeyIfCurrent(sharedKid, jwk1, &priv1.PublicKey)

	provider.keysMu.RLock()
	stale := provider.parsedKeys[sharedKid]
	provider.keysMu.RUnlock()
	if stale != nil {
		t.Fatal("expected stale parsed key from rotated same-kid JWK not to be cached")
	}

	got, err := provider.getParsedKey(sharedKid)
	if err != nil {
		t.Fatalf("getParsedKey after rotation failed: %v", err)
	}
	if got.N.Cmp(priv2.N) != 0 {
		t.Fatal("expected getParsedKey to parse and cache the current key material for sharedKid")
	}
}

func TestRemoteProviderGetSession_MidFlightInvalidationSkipsCacheWrite(t *testing.T) {
	const kid = "kid-inflight-introspect"
	priv, jwk := generateTestRSAKeyAndJWK(t, kid)

	var introspectCalls atomic.Int32
	inFlightStarted := make(chan struct{})
	releaseIntrospect := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/keys":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{jwk}})
		case "/introspect":
			if introspectCalls.Add(1) == 1 {
				close(inFlightStarted)
				<-releaseIntrospect
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	token := signTestJWT(t, priv, kid, time.Now().Add(5*time.Minute))

	makeReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
		return req
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- provider.GetSession(makeReq())
	}()

	<-inFlightStarted
	provider.invalidateTokenCache(token)
	close(releaseIntrospect)

	if err := <-errCh; err != nil {
		t.Fatalf("in-flight GetSession failed: %v", err)
	}

	if err := provider.GetSession(makeReq()); err != nil {
		t.Fatalf("post-invalidation GetSession failed: %v", err)
	}
	if got := introspectCalls.Load(); got != 2 {
		t.Fatalf("expected post-invalidation GetSession to re-introspect (2 calls), got %d", got)
	}
}

func TestRemoteProviderGetUserDetails_MidFlightInvalidationSkipsCacheWrite(t *testing.T) {
	const kid = "kid-inflight-profile"
	priv, _ := generateTestRSAKeyAndJWK(t, kid)
	userID := uuid.Must(uuid.NewV4())

	var profileCalls atomic.Int32
	inFlightStarted := make(chan struct{})
	releaseProfile := make(chan struct{})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/profile" {
			http.NotFound(w, r)
			return
		}
		if profileCalls.Add(1) == 1 {
			close(inFlightStarted)
			<-releaseProfile
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(UserPref{
			User: User{
				ID:        userID,
				FirstName: "Steve",
				LastName:  "Rogers",
				Email:     "cap@avengers.io",
			},
		})
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: UsersProfile, Endpoint: "/profile"},
	}

	token := signTestJWT(t, priv, kid, time.Now().Add(5*time.Minute))
	makeReq := func() *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/api/user", nil)
		req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: token})
		return req
	}

	errCh := make(chan error, 1)
	go func() {
		_, err := provider.GetUserDetails(makeReq())
		errCh <- err
	}()

	<-inFlightStarted
	provider.invalidateTokenCache(token)
	close(releaseProfile)

	if err := <-errCh; err != nil {
		t.Fatalf("in-flight GetUserDetails failed: %v", err)
	}

	if _, err := provider.GetUserDetails(makeReq()); err != nil {
		t.Fatalf("post-invalidation GetUserDetails failed: %v", err)
	}
	if got := profileCalls.Load(); got != 2 {
		t.Fatalf("expected post-invalidation GetUserDetails to re-fetch /profile (2 calls), got %d", got)
	}
}

func TestRemoteProviderAuthCache_UpdateExistingKeyDoesNotEvictAtCapacity(t *testing.T) {
	provider := &RemoteProvider{
		introspectCache:  make(map[string]time.Time),
		userProfileCache: make(map[string]cachedUserEntry),
	}

	future := time.Now().Add(10 * time.Minute)
	for i := 0; i < maxAuthCacheEntries; i++ {
		key := fmt.Sprintf("key-%d", i)
		provider.introspectCache[key] = future
		provider.userProfileCache[key] = cachedUserEntry{
			user:      User{FirstName: fmt.Sprintf("User-%d", i)},
			expiresAt: future,
		}
	}

	existingKey := "key-0"
	provider.setIntrospectCacheIfEpoch(existingKey, 5*time.Minute, provider.authCacheEpoch)
	if len(provider.introspectCache) != maxAuthCacheEntries {
		t.Fatalf("expected introspectCache size to stay at %d after updating existing key, got %d", maxAuthCacheEntries, len(provider.introspectCache))
	}
	if _, exists := provider.introspectCache[existingKey]; !exists {
		t.Fatalf("expected existing key %q to remain in introspectCache", existingKey)
	}
	if _, exists := provider.introspectCache["key-1"]; !exists {
		t.Fatalf("expected unrelated key to not be evicted when updating existing key")
	}

	provider.setCachedUserProfileIfEpoch(existingKey, &User{FirstName: "Updated-0"}, 5*time.Minute, provider.authCacheEpoch)
	if len(provider.userProfileCache) != maxAuthCacheEntries {
		t.Fatalf("expected userProfileCache size to stay at %d after updating existing key, got %d", maxAuthCacheEntries, len(provider.userProfileCache))
	}
	if _, exists := provider.userProfileCache[existingKey]; !exists {
		t.Fatalf("expected existing key %q to remain in userProfileCache", existingKey)
	}
	if _, exists := provider.userProfileCache["key-1"]; !exists {
		t.Fatalf("expected unrelated key to not be evicted when updating existing key")
	}

	newKey := "key-new"
	provider.setIntrospectCacheIfEpoch(newKey, 5*time.Minute, provider.authCacheEpoch)
	if len(provider.introspectCache) != maxAuthCacheEntries {
		t.Fatalf("expected introspectCache size to remain capped at %d after inserting new key, got %d", maxAuthCacheEntries, len(provider.introspectCache))
	}
	if _, exists := provider.introspectCache[newKey]; !exists {
		t.Fatalf("expected new key %q to be present in introspectCache", newKey)
	}

	provider.setCachedUserProfileIfEpoch(newKey, &User{FirstName: "New"}, 5*time.Minute, provider.authCacheEpoch)
	if len(provider.userProfileCache) != maxAuthCacheEntries {
		t.Fatalf("expected userProfileCache size to remain capped at %d after inserting new key, got %d", maxAuthCacheEntries, len(provider.userProfileCache))
	}
	if _, exists := provider.userProfileCache[newKey]; !exists {
		t.Fatalf("expected new key %q to be present in userProfileCache", newKey)
	}
}
