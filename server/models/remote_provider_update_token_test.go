package models

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/errors"
)

// newRequestWithoutToken builds a request with no provider session cookie.
// GetToken reads the token off that cookie, so its absence is what fails.
func newRequestWithoutToken(t *testing.T, target string) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodPut, target, nil)
}

// countingProvider stands up a remote provider that records whether it was
// reached at all.
func countingProvider(t *testing.T, hits *int32) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	return server
}

// TestRemoteProviderUpdateConnectionRequiresToken and its credential twin below
// pin one invariant on the two update paths: a request carrying no session
// cookie must not reach the remote provider at all.
//
// Both paths used to discard GetToken's error and guard on a stale one, so an
// empty token was PUT to the provider anyway - and against a provider answering
// 200, the caller got a zero-valued object back with no error at all. That is
// why the assertion is the hit counter and not just the returned error.
func TestRemoteProviderUpdateConnectionRequiresToken(t *testing.T) {
	var hits int32
	server := countingProvider(t, &hits)

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: PersistConnection, Endpoint: "/api/integrations/connections"},
	}

	conn, err := provider.UpdateConnection(
		newRequestWithoutToken(t, "/api/integrations/connections"),
		&connections.Connection{Kind: "kubernetes"},
	)
	if err == nil {
		t.Fatalf("expected an error when the request carries no token, got connection %+v", conn)
	}
	if conn != nil {
		t.Errorf("expected no connection alongside the error, got %+v", conn)
	}
	if code := errors.GetCode(err); code != ErrGetTokenCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrGetTokenCode, err)
	}
	if got := atomic.LoadInt32(&hits); got != 0 {
		t.Errorf("remote provider was called %d time(s) with an empty token; a missing token must short-circuit the request", got)
	}
}

// TestRemoteProviderUpdateUserCredentialRequiresToken covers the same swallowed
// GetToken error on the credential path, where the body being PUT is a
// credential and an unauthenticated write is one the provider must never make.
func TestRemoteProviderUpdateUserCredentialRequiresToken(t *testing.T) {
	var hits int32
	server := countingProvider(t, &hits)

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{
		{Feature: PersistCredentials, Endpoint: "/api/integrations/credentials"},
	}

	cred, err := provider.UpdateUserCredential(
		newRequestWithoutToken(t, "/api/integrations/credentials"),
		&Credential{},
	)
	if err == nil {
		t.Fatalf("expected an error when the request carries no token, got credential %+v", cred)
	}
	if cred != nil {
		t.Errorf("expected no credential alongside the error, got %+v", cred)
	}
	if code := errors.GetCode(err); code != ErrGetTokenCode {
		t.Errorf("error code = %q, want %q (err: %v)", code, ErrGetTokenCode, err)
	}
	if got := atomic.LoadInt32(&hits); got != 0 {
		t.Errorf("remote provider was called %d time(s) with an empty token; a missing token must short-circuit the request", got)
	}
}
