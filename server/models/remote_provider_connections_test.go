package models

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/errors"
)

const (
	remoteShapeConnectionID = "11111111-1111-1111-1111-111111111111"
	remoteShapeOtherID      = "22222222-2222-2222-2222-222222222222"
)

func remoteConnectionJSON(id, name string) string {
	return fmt.Sprintf(`{"id":%q,"name":%q,"kind":"kubernetes","type":"platform","subType":"orchestrator","status":"discovered"}`, id, name)
}

func persistConnectionCapability() Capabilities {
	return Capabilities{
		{Feature: PersistConnection, Endpoint: "/api/integrations/connections"},
	}
}

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

func TestRemoteProviderSaveConnectionResponseShapes(t *testing.T) {
	prod := remoteConnectionJSON(remoteShapeConnectionID, "prod")
	other := remoteConnectionJSON(remoteShapeOtherID, "other")
	page := fmt.Sprintf(`{"page":0,"page_size":2,"total_count":2,"connections":[%s,%s]}`, prod, other)

	cases := []struct {
		name     string
		status   int
		body     string
		wantErr  bool
		wantID   string
		wantName string
	}{
		{name: "201 page", status: http.StatusCreated, body: page, wantID: remoteShapeConnectionID, wantName: "prod"},
		{name: "201 bare connection", status: http.StatusCreated, body: prod, wantID: remoteShapeConnectionID, wantName: "prod"},
		{name: "200 bare connection", status: http.StatusOK, body: prod, wantID: remoteShapeConnectionID, wantName: "prod"},
		{name: "201 empty page", status: http.StatusCreated, body: `{"page":0,"page_size":0,"total_count":0,"connections":[]}`, wantErr: true},
		{name: "400", status: http.StatusBadRequest, body: `{"error":"bad"}`, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("method = %s, want POST", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/api/integrations/connections") {
					t.Errorf("path = %s, want the persist-connection endpoint", r.URL.Path)
				}
				w.WriteHeader(tc.status)
				_, _ = io.WriteString(w, tc.body)
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = persistConnectionCapability()

			conn, err := provider.SaveConnection(&connections.ConnectionPayload{
				Name:   "prod",
				Kind:   "kubernetes",
				Type:   "platform",
				Status: connections.DISCOVERED,
			}, "test-token", true)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got connection %+v", conn)
				}
				if conn != nil {
					t.Errorf("connection = %+v, want nil alongside the error", conn)
				}
				if errors.GetCode(err) != ErrPostCode {
					t.Errorf("error code = %q, want %q (err: %v)", errors.GetCode(err), ErrPostCode, err)
				}
				// WithProviderStatus only records 4xx and 5xx. A 2xx that
				// still failed to yield a connection stays an ErrPost without
				// a provider status tag.
				if tc.status >= http.StatusBadRequest {
					got, ok := httputil.ProviderStatusCode(err)
					if !ok || got != tc.status {
						t.Errorf("provider status = %d (ok=%v), want %d", got, ok, tc.status)
					}
				}
				return
			}
			if err != nil {
				t.Fatalf("SaveConnection: %v", err)
			}
			if conn == nil {
				t.Fatal("connection is nil")
			}
			if conn.ID.String() != tc.wantID {
				t.Errorf("id = %s, want %s", conn.ID, tc.wantID)
			}
			if conn.Name != tc.wantName {
				t.Errorf("name = %q, want %q", conn.Name, tc.wantName)
			}
		})
	}
}

func TestRemoteProviderDeleteConnectionResponseShapes(t *testing.T) {
	bodyID := remoteShapeConnectionID
	bare := remoteConnectionJSON(bodyID, "prod")
	page := fmt.Sprintf(`{"page":0,"page_size":1,"total_count":1,"connections":[%s]}`, bare)

	cases := []struct {
		name     string
		status   int
		body     string
		wantErr  bool
		wantBody bool
		wantName string
	}{
		{name: "200 connection", status: http.StatusOK, body: bare, wantBody: true, wantName: "prod"},
		{name: "200 page", status: http.StatusOK, body: page, wantBody: true, wantName: "prod"},
		{name: "204 empty", status: http.StatusNoContent, wantName: ""},
		{name: "200 empty body", status: http.StatusOK, wantName: ""},
		{name: "404", status: http.StatusNotFound, body: `{"error":"no such connection"}`, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			deletedID := uuid.Must(uuid.NewV4())
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("method = %s, want DELETE", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/"+deletedID.String()) {
					t.Errorf("path = %s, want id %s", r.URL.Path, deletedID)
				}
				w.WriteHeader(tc.status)
				if tc.body != "" {
					_, _ = io.WriteString(w, tc.body)
				}
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = persistConnectionCapability()

			conn, err := provider.DeleteConnection(newDeleteConnectionRequest(t), deletedID)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got connection %+v", conn)
				}
				if conn != nil {
					t.Errorf("connection = %+v, want nil alongside the error", conn)
				}
				got, ok := httputil.ProviderStatusCode(err)
				if !ok || got != tc.status {
					t.Errorf("provider status = %d (ok=%v), want %d", got, ok, tc.status)
				}
				return
			}
			if err != nil {
				t.Fatalf("DeleteConnection: %v", err)
			}
			if conn == nil {
				t.Fatal("connection is nil")
			}
			wantID := deletedID.String()
			if tc.wantBody {
				wantID = bodyID
			}
			if conn.ID.String() != wantID {
				t.Errorf("id = %s, want %s", conn.ID, wantID)
			}
			if conn.Name != tc.wantName {
				t.Errorf("name = %q, want %q", conn.Name, tc.wantName)
			}
		})
	}
}

func TestRemoteProviderDeleteMesheryConnectionAcceptsNoContent(t *testing.T) {
	setViperKeyForRemoteAuthTest(t, "INSTANCE_ID", "instance-1")

	cases := []int{http.StatusOK, http.StatusNoContent}
	for _, code := range cases {
		t.Run(http.StatusText(code), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("method = %s, want DELETE", r.Method)
				}
				if !strings.HasSuffix(r.URL.Path, "/api/integrations/connections/meshery/instance-1") {
					t.Errorf("path = %s", r.URL.Path)
				}
				w.WriteHeader(code)
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			provider.Capabilities = persistConnectionCapability()
			if err := provider.DeleteMesheryConnection(); err != nil {
				t.Fatalf("DeleteMesheryConnection(%d): %v", code, err)
			}
		})
	}
}
