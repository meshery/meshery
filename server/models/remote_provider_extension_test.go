package models

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshkit/errors"
)

// TestExtensionProxy_InvalidPaths verifies that malformed paths or SSRF authority injection
// attempts are rejected with ErrInvalidExtensionProxyPathCode without crashing the server.
func TestExtensionProxy_InvalidPaths(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://remote.example")

	invalidPaths := []struct {
		name string
		path string
	}{
		{
			name: "missing prefix entirely",
			path: "/invalid/path",
		},
		{
			name: "empty path",
			path: "",
		},
		{
			name: "partial prefix match without boundary",
			path: "/api/extensionsfoo",
		},
		{
			name: "ssrf attempt with at-sign host",
			path: "/api/extensions@evil.com",
		},
		{
			name: "ssrf attempt with at-sign path",
			path: "/api/extensions@evil.com/subpath",
		},
		{
			name: "prefix in middle of path",
			path: "/other/api/extensions/users",
		},
	}

	for _, tc := range invalidPaths {
		t.Run(tc.name, func(t *testing.T) {
			target := tc.path
			if target == "" {
				target = "/"
			}
			req := httptest.NewRequest("GET", target, nil)
			if tc.path == "" {
				req.URL.Path = ""
			}

			resp, err := provider.ExtensionProxy(req)
			if resp != nil {
				t.Fatalf("expected nil response, got %+v", resp)
			}
			if err == nil {
				t.Fatalf("expected error for path %q, got nil", tc.path)
			}

			code := errors.GetCode(err)
			if code != ErrInvalidExtensionProxyPathCode {
				t.Fatalf("expected error code %q, got %q (err: %v)", ErrInvalidExtensionProxyPathCode, code, err)
			}
		})
	}
}

// TestExtensionProxy_ValidPathPreservation asserts that valid extension proxy paths are
// correctly forwarded to the remote provider without truncation or corruption.
func TestExtensionProxy_ValidPathPreservation(t *testing.T) {
	var receivedPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer ts.Close()

	provider := newTestRemoteProvider(t, ts.URL)

	validPaths := []struct {
		name         string
		path         string
		expectedPath string
	}{
		{
			name:         "exact extensions root",
			path:         "/api/extensions",
			expectedPath: "/",
		},
		{
			name:         "extensions trailing slash",
			path:         "/api/extensions/",
			expectedPath: "/",
		},
		{
			name:         "valid subpath",
			path:         "/api/extensions/users",
			expectedPath: "/users",
		},
		{
			name:         "nested extensions substring",
			path:         "/api/extensions/deep/api/extensions/sub",
			expectedPath: "/deep/api/extensions/sub",
		},
	}

	for _, tc := range validPaths {
		t.Run(tc.name, func(t *testing.T) {
			receivedPath = ""
			req := httptest.NewRequest("GET", tc.path, nil)
			req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "test-token"})

			resp, err := provider.ExtensionProxy(req)
			if err != nil {
				t.Fatalf("unexpected error for path %q: %v", tc.path, err)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200 OK, got %+v", resp)
			}

			if receivedPath != tc.expectedPath {
				t.Fatalf("expected upstream forwarded path %q, got %q", tc.expectedPath, receivedPath)
			}
		})
	}
}
