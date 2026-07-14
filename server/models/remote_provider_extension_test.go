package models

import (
	"net/http"
	"net/http/httptest"
	"testing"

	mkerrors "github.com/meshery/meshkit/errors"
)

func TestExtensionProxy_PanicFix(t *testing.T) {
	l := newTestRemoteProvider(t, "http://remote.example")

	// This path does not contain /api/extensions, which previously caused an index out of bounds panic
	req := httptest.NewRequest(http.MethodGet, "http://meshery.local/invalid/path", nil)

	// The function should now return an error instead of panicking
	_, err := l.ExtensionProxy(req)
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrInvalidExtensionProxyPathCode {
		t.Fatalf("expected error code %s, got %s", ErrInvalidExtensionProxyPathCode, code)
	}
}

func TestExtensionProxy_NoDelimiterBoundary(t *testing.T) {
	l := newTestRemoteProvider(t, "http://remote.example")

	// Path contains "extensions" but not as a proper "/api/extensions" prefix boundary.
	// "/api/extensionsfoo" should NOT be treated as a valid extension proxy path.
	req := httptest.NewRequest(http.MethodGet, "http://meshery.local/api/extensionsfoo", nil)

	_, err := l.ExtensionProxy(req)
	if err == nil {
		t.Fatal("expected error for path without proper /api/extensions boundary, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrInvalidExtensionProxyPathCode {
		t.Fatalf("expected error code %s, got %s", ErrInvalidExtensionProxyPathCode, code)
	}
}

func TestExtensionProxy_EmptyRemainder(t *testing.T) {
	l := newTestRemoteProvider(t, "http://remote.example")

	// Path is exactly "/api/extensions" with no trailing path. split[1] will be "".
	// The function should proceed (empty path appended to remote provider URL).
	req := httptest.NewRequest(http.MethodGet, "http://meshery.local/api/extensions", nil)

	// Request has a valid /api/extensions path, so validation should pass.
	// The call should fail at token retrieval because the request has no token cookie.
	_, err := l.ExtensionProxy(req)
	if err == nil {
		t.Fatal("expected error due to missing token cookie, got nil")
	}
	if code := mkerrors.GetCode(err); code != ErrGetTokenCode {
		t.Fatalf("expected error code %s, got %s", ErrGetTokenCode, code)
	}
}
