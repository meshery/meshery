package models

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/core"
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

func TestRemoteProviderGetCredentialByID_UsesResponseCache(t *testing.T) {
	credentialUUID := uuid.FromStringOrNil("1f8b84e2-5af1-4c8c-a1d4-6f91f8c0c2b1")
	credentialID := core.UUIDOrUUIDNil(&credentialUUID)

	var credentialCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/integrations/credentials/") {
			atomic.AddInt32(&credentialCalls, 1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"` + credentialID.String() + `","secret":{"secret":"abc"}}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: PersistCredentials, Endpoint: "/api/integrations/credentials"},
	}

	cred1, statusCode1, err := provider.GetCredentialByID("token-a", credentialID)
	if err != nil {
		t.Fatalf("unexpected error on first fetch: %v", err)
	}
	if statusCode1 != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode1)
	}
	if cred1 == nil {
		t.Fatal("expected non-nil credential on first fetch")
	}

	cred2, statusCode2, err := provider.GetCredentialByID("token-a", credentialID)
	if err != nil {
		t.Fatalf("unexpected error on second fetch: %v", err)
	}
	if statusCode2 != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode2)
	}
	if cred2 == nil {
		t.Fatal("expected non-nil credential on second fetch")
	}

	if got := atomic.LoadInt32(&credentialCalls); got != 1 {
		t.Fatalf("expected one remote credential request, got %d", got)
	}
}

func TestRemoteProviderLoadAllK8sContext_UsesResponseCache(t *testing.T) {
	instanceUUID := uuid.FromStringOrNil("90d7a88c-68a7-4f58-a67d-e9e3f8b4f3f0")
	instanceID := core.UUIDOrUUIDNil(&instanceUUID)
	viper.Set("INSTANCE_ID", &instanceID)

	var k8sContextCalls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/integrations/connections") {
			atomic.AddInt32(&k8sContextCalls, 1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"page":0,"page_size":25,"total_count":0,"connections":[]}`))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: PersistConnection, Endpoint: "/api/integrations/connections"},
	}

	contexts1, err := provider.LoadAllK8sContext("token-b")
	if err != nil {
		t.Fatalf("unexpected error on first contexts fetch: %v", err)
	}
	if len(contexts1) != 0 {
		t.Fatalf("expected 0 contexts on first fetch, got %d", len(contexts1))
	}

	contexts2, err := provider.LoadAllK8sContext("token-b")
	if err != nil {
		t.Fatalf("unexpected error on second contexts fetch: %v", err)
	}
	if len(contexts2) != 0 {
		t.Fatalf("expected 0 contexts on second fetch, got %d", len(contexts2))
	}

	if got := atomic.LoadInt32(&k8sContextCalls); got != 1 {
		t.Fatalf("expected one remote k8s contexts request, got %d", got)
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
