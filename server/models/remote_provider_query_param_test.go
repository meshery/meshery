package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
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


// TestRemoteProvider_FetchCapabilities_EscapesQueryParams pins URL query parameter
// encoding for remote provider capability requests. When OS or PLAYGROUND environment
// variables contain reserved URL characters (such as '&', '=', or '#'), they must be
// safely percent-encoded so that parameter keys and values are preserved over the wire
// rather than being split into separate parameters or truncated as URL fragments.
func TestRemoteProvider_FetchCapabilities_EscapesQueryParams(t *testing.T) {
	const (
		testOS         = "linux&arm=64"
		testPlayground = "preview#blue"
		testVersion    = "v1.2.3"
	)

	origOS := viper.GetString("OS")
	origPlayground := viper.GetString("PLAYGROUND")
	origBuild := viper.GetString("BUILD")
	t.Cleanup(func() {
		viper.Set("OS", origOS)
		viper.Set("PLAYGROUND", origPlayground)
		viper.Set("BUILD", origBuild)
	})

	viper.Set("OS", testOS)
	viper.Set("PLAYGROUND", testPlayground)
	viper.Set("BUILD", testVersion)

	tests := []struct {
		name      string
		versioned bool
		wantPath  string
	}{
		{
			name:      "unversioned capabilities request encodes reserved characters",
			versioned: false,
			wantPath:  "/capabilities",
		},
		{
			name:      "versioned capabilities request encodes reserved characters",
			versioned: true,
			wantPath:  "/" + testVersion + "/capabilities",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var gotQuery url.Values
			var gotPath string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				gotQuery = r.URL.Query()
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			defer server.Close()

			provider := newTestRemoteProvider(t, server.URL)
			if _, err := provider.fetchCapabilities(context.Background(), "test-token", tc.versioned); err != nil {
				t.Fatalf("fetchCapabilities: %v", err)
			}

			if gotPath != tc.wantPath {
				t.Errorf("path = %q, want %q", gotPath, tc.wantPath)
			}
			if gotOS := gotQuery.Get("os"); gotOS != testOS {
				t.Errorf("os query param = %q, want %q (full query: %v)", gotOS, testOS, gotQuery)
			}
			if gotPG := gotQuery.Get("playground"); gotPG != testPlayground {
				t.Errorf("playground query param = %q, want %q (full query: %v)", gotPG, testPlayground, gotQuery)
			}
			if _, stray := gotQuery["arm"]; stray {
				t.Errorf("unexpected stray query param 'arm' found due to unencoded '&': %v", gotQuery)
			}
		})
	}
}
