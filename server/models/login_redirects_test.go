package models

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolvePostLoginRedirect(t *testing.T) {
	t.Parallel()

	const fallback = "/"

	tests := []struct {
		name     string
		rawRef   string
		expected string
	}{
		{
			name:     "empty ref falls back",
			rawRef:   "",
			expected: fallback,
		},
		{
			name:     "encoded in app path is decoded",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/extension/meshmap")),
			expected: "/extension/meshmap",
		},
		{
			name:     "plain in app path is preserved",
			rawRef:   "/extension/meshmap",
			expected: "/extension/meshmap",
		},
		{
			name:     "encoded absolute url falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/phish")),
			expected: fallback,
		},
		{
			name:     "plain absolute url falls back",
			rawRef:   "https://evil.example/phish",
			expected: fallback,
		},
		{
			name:     "invalid base64 falls back",
			rawRef:   "not-base64",
			expected: fallback,
		},
		// Regression coverage: /user/login and /api/user/token are auth
		// initiation paths. Redirecting to them after a successful token
		// exchange re-enters the OAuth dance and caused Kanvas to hang on
		// the loading splash indefinitely (meshery-server-1345 followed by
		// a second InitiateLogin in the same second).
		{
			name:     "plain /user/login ref falls back",
			rawRef:   "/user/login",
			expected: fallback,
		},
		{
			name:     "/user/login with query falls back",
			rawRef:   "/user/login?provider=Layer5",
			expected: fallback,
		},
		{
			name:     "encoded /user/login ref falls back",
			rawRef:   base64.RawURLEncoding.EncodeToString([]byte("/user/login?provider=Layer5")),
			expected: fallback,
		},
		{
			name:     "plain /api/user/token ref falls back",
			rawRef:   "/api/user/token",
			expected: fallback,
		},
		{
			name:     "/provider ref falls back",
			rawRef:   "/provider?ref=xyz",
			expected: fallback,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := resolvePostLoginRedirect(tc.rawRef, fallback)
			if actual != tc.expected {
				t.Fatalf("expected redirect %q, got %q", tc.expected, actual)
			}
		})
	}
}

func TestSelectPostLoginRefValue(t *testing.T) {
	t.Parallel()

	const cookieName = "playground.meshery.io_ref"
	const cookieValue = "L2V4dGVuc2lvbi9tZXNobWFw" // base64 of /extension/meshmap
	const queryValue = "L2Rhc2hib2FyZA"             // base64 of /dashboard

	tests := []struct {
		name     string
		cookie   *http.Cookie
		query    string
		expected string
	}{
		{
			name:     "cookie value is used when present",
			cookie:   &http.Cookie{Name: cookieName, Value: cookieValue},
			expected: cookieValue,
		},
		// Regression: the cookie is the SOLE source of truth. A ?ref= the
		// remote provider echoes back must never override (or fill in for)
		// the cookie — that's how the playground.meshery.io 404 escaped in
		// the first place. resolvePostLoginRedirect's "/" fallback handles
		// the missing-cookie case without us re-trusting provider state.
		{
			name:     "ignores ?ref= query param even when cookie missing",
			query:    "?ref=" + queryValue,
			expected: "",
		},
		{
			name:     "ignores ?ref= query param when cookie is empty",
			cookie:   &http.Cookie{Name: cookieName, Value: ""},
			query:    "?ref=" + queryValue,
			expected: "",
		},
		{
			name:     "returns empty when cookie is missing",
			expected: "",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/api/user/token"+tc.query, nil)
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}
			actual := selectPostLoginRefValue(req, cookieName)
			if actual != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, actual)
			}
		})
	}
}
