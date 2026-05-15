package models

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
// fallback path and broke anonymous flows in kanvas.new. The test also
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
