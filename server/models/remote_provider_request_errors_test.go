package models

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// invalidRemoteProviderURL does not parse as a URL: the IPv6 host is never
// closed with "]". Building the outbound request from it must fail.
const invalidRemoteProviderURL = "http://[::1"

// expectRequestBuildError calls save and fails the test when it panics or
// returns no error. While the errors from url.Parse and http.NewRequest were
// discarded, a remote provider URL that did not parse left the parsed URL nil,
// and calling its String method panicked before any request was sent.
func expectRequestBuildError(t *testing.T, name string, save func() ([]byte, error)) {
	t.Helper()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("%s panicked on a remote provider URL that does not parse: %v", name, r)
		}
	}()

	body, err := save()
	if err == nil {
		t.Fatalf("%s returned no error for a remote provider URL that does not parse", name)
	}
	if body != nil {
		t.Errorf("%s returned a body alongside the error: %q", name, body)
	}
}

func TestRemotePatternFile_InvalidRemoteProviderURL(t *testing.T) {
	provider := newTestRemoteProvider(t, invalidRemoteProviderURL)
	provider.Capabilities = Capabilities{{Feature: PersistMesheryPatterns, Endpoint: "/api/content/patterns"}}
	req := httptest.NewRequest(http.MethodPost, "/api/pattern", nil)

	expectRequestBuildError(t, "RemotePatternFile", func() ([]byte, error) {
		return provider.RemotePatternFile(req, "https://example.com/design.yaml", "", true)
	})
}

func TestSaveMesheryFilter_InvalidRemoteProviderURL(t *testing.T) {
	provider := newTestRemoteProvider(t, invalidRemoteProviderURL)
	provider.Capabilities = Capabilities{{Feature: PersistMesheryFilters, Endpoint: "/api/content/filters"}}

	expectRequestBuildError(t, "SaveMesheryFilter", func() ([]byte, error) {
		return provider.SaveMesheryFilter("test-token", &MesheryFilter{})
	})
}

func TestRemoteFilterFile_InvalidRemoteProviderURL(t *testing.T) {
	provider := newTestRemoteProvider(t, invalidRemoteProviderURL)
	provider.Capabilities = Capabilities{{Feature: PersistMesheryFilters, Endpoint: "/api/content/filters"}}
	req := httptest.NewRequest(http.MethodPost, "/api/filter", nil)

	expectRequestBuildError(t, "RemoteFilterFile", func() ([]byte, error) {
		return provider.RemoteFilterFile(req, "https://example.com/filter.wasm", "", true, "")
	})
}
