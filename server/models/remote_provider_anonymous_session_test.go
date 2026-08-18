package models

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newAnonymousFlowProvider returns a RemoteProvider whose anonymous-user mint
// is served by handler, with the PersistAnonymousUser capability wired to it.
func newAnonymousFlowProvider(t *testing.T, handler http.HandlerFunc) (*RemoteProvider, func()) {
	t.Helper()

	server := httptest.NewServer(handler)

	provider := newTestRemoteProvider(t, server.URL)
	provider.ProviderURL = server.URL
	provider.ProviderProperties = ProviderProperties{
		ProviderType: RemoteProviderType,
		ProviderURL:  server.URL,
		ProviderName: "test-remote",
		Capabilities: Capabilities{{Feature: PersistAnonymousUser, Endpoint: "/anonymous"}},
	}
	// A working capabilities store, so removing the nil-id refusal lets the
	// flow run to completion and fail on the assertions below rather than on a
	// nil persister - the guard is what has to be under test, not the fixture.
	provider.UserCapabilitiesPersister = &UserCapabilitiesPersister{DB: newMigratedDB(t)}

	return provider, server.Close
}

func newAnonymousFlowRequest() *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	return req.WithContext(context.WithValue(req.Context(), MesheryServerURL, "http://localhost:9081"))
}

func jwtCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == TokenCookieName {
			return c
		}
	}
	return nil
}

// TestInterceptAnonymousSession_RefusesReplyWithoutUserID pins the refusal added
// alongside the AnonymousFlowResponse repoint.
//
// The reply is the schemas v1beta2 user.AnonymousFlowResponse construct, keyed
// on `userId`. When a local copy of the struct read `owner` instead, every
// anonymous sign-in decoded a nil id and the session's capabilities were written
// under the zero UUID - one shared row for every anonymous user, with no error
// anywhere. A reply that genuinely carries no id must therefore fail the sign-in
// rather than proceed, and it must fail BEFORE the JWT cookie is set so a
// refused session leaves no half-authenticated browser behind.
func TestInterceptAnonymousSession_RefusesReplyWithoutUserID(t *testing.T) {
	for _, tc := range []struct {
		name  string
		reply string
	}{
		{"no userId at all", `{"accessToken":"eyJhbGciOiJIUzI1NiJ9.e30.sig"}`},
		{"explicit nil uuid", `{"accessToken":"eyJhbGciOiJIUzI1NiJ9.e30.sig","userId":"00000000-0000-0000-0000-000000000000"}`},
		{"legacy owner key", `{"accessToken":"eyJhbGciOiJIUzI1NiJ9.e30.sig","owner":"0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			provider, closeServer := newAnonymousFlowProvider(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.reply))
			})
			defer closeServer()

			rec := httptest.NewRecorder()
			provider.InterceptLoginAndInitiateAnonymousUserSession(newAnonymousFlowRequest(), rec)

			if rec.Code != http.StatusFound {
				t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
			}
			if got := rec.Header().Get("Location"); got != "/error" {
				t.Errorf("redirected to %q, want %q", got, "/error")
			}
			// Presence, not value: a refused session must emit no Set-Cookie for
			// the token at all. Allowing an empty-valued cookie through would let
			// the handler start writing one again without failing this test.
			if ck := jwtCookie(rec); ck != nil {
				t.Errorf("a JWT cookie was set for a refused session: %q", ck.Value)
			}
		})
	}
}
