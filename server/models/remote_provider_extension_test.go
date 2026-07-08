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
