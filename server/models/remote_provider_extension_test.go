package models

import (
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshkit/errors"
)

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

func TestExtensionProxy_ValidPathPreservation(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://remote.example")

	validPaths := []struct {
		name string
		path string
	}{
		{
			name: "exact extensions root",
			path: "/api/extensions",
		},
		{
			name: "extensions trailing slash",
			path: "/api/extensions/",
		},
		{
			name: "valid subpath",
			path: "/api/extensions/users",
		},
		{
			name: "nested extensions substring",
			path: "/api/extensions/deep/api/extensions/sub",
		},
	}

	for _, tc := range validPaths {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)

			// We expect the call to pass path validation.
			// It will fail at GetToken (because no token cookie is provided in this test request),
			// which proves path validation succeeded without ErrInvalidExtensionProxyPathCode.
			_, err := provider.ExtensionProxy(req)
			if err == nil {
				t.Fatalf("expected error at GetToken stage, got nil")
			}

			code := errors.GetCode(err)
			if code == ErrInvalidExtensionProxyPathCode {
				t.Fatalf("path %q was unexpectedly rejected as invalid extension proxy path: %v", tc.path, err)
			}
			if code != ErrGetTokenCode {
				t.Fatalf("expected error code %q (GetToken stage), got %q (err: %v)", ErrGetTokenCode, code, err)
			}
		})
	}
}
