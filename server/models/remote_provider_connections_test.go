package models

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/errors"
)

// newDeleteConnectionRequest builds the inbound request DeleteConnection reads
// its provider token off. RemoteProvider.GetToken takes the token from the
// session cookie, so a request without one never reaches the outbound call.
func newDeleteConnectionRequest(t *testing.T) *http.Request {
	t.Helper()

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/connections/x", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})
	return req
}

// TestRemoteProviderDeleteConnectionPropagatesStatus pins the contract the
// connection handler depends on: whatever status the remote provider returned
// for a DELETE must survive on the error, so the handler can answer with it
// rather than inventing one.
//
// The case that mattered is 404. `mesheryctl connection delete <unknown-uuid>`
// is expected to warn "No connection with ID ..." and exit 0, which it decides
// from an HTTP 404. Before ErrDelete carried the provider status, a remote 404
// reached the handler as an opaque error, was wrapped in ErrFailToSave and
// answered 500 "Failed to Save: .connection" (meshery-server-1051) with a
// "Provider Database could be down or not reachable" remediation - for an id
// that was simply never issued.
func TestRemoteProviderDeleteConnectionPropagatesStatus(t *testing.T) {
	cases := []struct {
		name       string
		remoteCode int
		wantStatus int
	}{
		{name: "not found", remoteCode: http.StatusNotFound, wantStatus: http.StatusNotFound},
		{name: "forbidden", remoteCode: http.StatusForbidden, wantStatus: http.StatusForbidden},
		{name: "provider error", remoteCode: http.StatusInternalServerError, wantStatus: http.StatusInternalServerError},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE, got %s", r.Method)
				}
				w.WriteHeader(tc.remoteCode)
				_, _ = w.Write([]byte(`{"error":"no such connection"}`))
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = Capabilities{
				{Feature: PersistConnection, Endpoint: "/api/integrations/connections"},
			}

			conn, err := provider.DeleteConnection(newDeleteConnectionRequest(t), uuid.Must(uuid.NewV4()))
			if err == nil {
				t.Fatalf("expected an error for remote status %d, got connection %+v", tc.remoteCode, conn)
			}
			if conn != nil {
				t.Errorf("expected no connection alongside the error, got %+v", conn)
			}

			got, ok := httputil.ProviderStatusCode(err)
			if !ok {
				t.Fatalf("error did not carry the provider status: %v", err)
			}
			if got != tc.wantStatus {
				t.Errorf("provider status = %d, want %d (err: %v)", got, tc.wantStatus, err)
			}
		})
	}
}

// TestRemoteProviderDeleteConnectionUnreachableProvider covers the error path
// that used to panic: DoRequest returns (nil, err) when the provider cannot be
// reached, and every caller then read resp.StatusCode off that nil response.
// The failure has to come back as an error, not a nil-pointer dereference in
// the request goroutine.
func TestRemoteProviderDeleteConnectionUnreachableProvider(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	unreachableURL := server.URL
	server.Close() // nothing is listening on this port any more

	provider := newTestRemoteProvider(t, unreachableURL)
	provider.Capabilities = Capabilities{
		{Feature: PersistConnection, Endpoint: "/api/integrations/connections"},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("DeleteConnection panicked on an unreachable provider: %v", r)
		}
	}()

	conn, err := provider.DeleteConnection(newDeleteConnectionRequest(t), uuid.Must(uuid.NewV4()))
	if err == nil {
		t.Fatalf("expected an error for an unreachable provider, got connection %+v", conn)
	}
	if errors.GetCode(err) != ErrDeleteCode {
		t.Errorf("error code = %q, want %q (err: %v)", errors.GetCode(err), ErrDeleteCode, err)
	}
	// The request never produced an HTTP status, so none must be claimed - a
	// fabricated status is what makes an unreachable provider look like a
	// deliberate refusal.
	if status, ok := httputil.ProviderStatusCode(err); ok {
		t.Errorf("expected no provider status for an unreachable provider, got %d", status)
	}
}
