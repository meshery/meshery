package models

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
)

// TestGetMesheryPatternResources_SendsCanonicalOamTypeParam pins the
// pattern-resource filter to the canonical camelCase `oamType` the remote
// provider reads (the identifier-naming contract: URL query params are
// camelCase).
//
// This is an outbound-only contract - nothing in Meshery Server observes the
// spelling - so a snake_case `oam_type` regression is invisible locally: the
// provider ignores the unknown key and answers 200 with an unfiltered page.
// Design-engine resource lookups then match the wrong resource, or none.
func TestGetMesheryPatternResources_SendsCanonicalOamTypeParam(t *testing.T) {
	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"resources":[]}`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistMesheryPatternResources, Endpoint: "/pattern-resources"}}

	if _, err := provider.GetMesheryPatternResources(
		"token", "0", "10", "", "", "demo", "default", "Deployment", "workload",
	); err != nil {
		t.Fatalf("GetMesheryPatternResources: %v", err)
	}

	if got := gotQuery.Get("oamType"); got != "workload" {
		t.Errorf(`oamType = %q, want "workload" (query was %v)`, got, gotQuery)
	}
	if _, legacy := gotQuery["oam_type"]; legacy {
		t.Errorf("request still carries the legacy snake_case `oam_type`: %v", gotQuery)
	}
}

// TestRemoteProviderDeleteUserCredential_SendsCanonicalCredentialIdParam pins
// the delete param to the canonical camelCase `credentialId`, which is both what
// the schemas `deleteUserCredential` operation declares and what Meshery
// Server's own handler now reads.
//
// Also outbound-only: a snake_case regression here leaves the provider with no
// id to delete, and the caller cannot tell.
func TestRemoteProviderDeleteUserCredential_SendsCanonicalCredentialIdParam(t *testing.T) {
	credentialID := core.Uuid(uuid.Must(uuid.FromString("55555555-5555-5555-5555-555555555555")))

	var gotQuery url.Values

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistCredentials, Endpoint: "/credentials"}}

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "token"})

	if _, err := provider.DeleteUserCredential(req, credentialID); err != nil {
		t.Fatalf("DeleteUserCredential: %v", err)
	}

	if got := gotQuery.Get("credentialId"); got != credentialID.String() {
		t.Errorf("credentialId = %q, want %q (query was %v)", got, credentialID, gotQuery)
	}
	if _, legacy := gotQuery["credential_id"]; legacy {
		t.Errorf("request still carries the legacy snake_case `credential_id`: %v", gotQuery)
	}
}
