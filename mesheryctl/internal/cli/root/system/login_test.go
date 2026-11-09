// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
)

const existingAuth = `{"meshery-provider":"Meshery","token":"EXISTING-VALID-TOKEN"}`

// setupLoginEnv points the current context at a fake Meshery Server whose
// /api/token endpoint is served by tokenHandler, and the context's token at a
// private file. It returns the token file path, the log buffer and a counter
// of requests that reached the redirect target.
func setupLoginEnv(t *testing.T, tokenHandler http.HandlerFunc, authContent string, authMode os.FileMode) (string, *bytes.Buffer, *atomic.Int32) {
	t.Helper()
	setupContextTestEnv(t)
	logs := utils.SetupMeshkitLoggerTesting(t, false)

	redirectHits := &atomic.Int32{}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/system/version", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"build":"test"}`)
	})
	mux.HandleFunc("/api/providers", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"Local":{"providerName":"Local"}}`)
	})
	mux.HandleFunc("/provider", func(w http.ResponseWriter, _ *http.Request) {
		redirectHits.Add(1)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, "<html><body>Choose a provider</body></html>")
	})
	mux.HandleFunc("/api/token", tokenHandler)
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	authPath := filepath.Join(t.TempDir(), "auth.json")
	if authContent != "" {
		if err := os.WriteFile(authPath, []byte(authContent), authMode); err != nil {
			t.Fatalf("unable to write auth fixture: %v", err)
		}
	}
	viper.Set("contexts.local.endpoint", srv.URL)
	viper.Set("tokens", []map[string]string{{"name": "default", "location": authPath}})

	return authPath, logs, redirectHits
}

// respond returns a handler that writes a fixed status, content type and body.
func respond(status int, contentType, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = fmt.Fprint(w, body)
	}
}

// redirectTo returns a handler that redirects to the HTML provider page with code.
func redirectTo(code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/provider?ref=x", code)
	}
}

// runLogin executes `system login -p Local` through the root system command.
func runLogin(t *testing.T) error {
	t.Helper()
	SystemCmd.SetArgs([]string{"login", "-p", "Local"})
	return SystemCmd.Execute()
}

// TestLoginPersistsValidTokenResponse checks that a valid token response,
// including the Local provider's empty token, is written and reported.
func TestLoginPersistsValidTokenResponse(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "Local provider with an empty token",
			body: `{"meshery-provider":"Local","token":""}` + "\n",
			want: `{"meshery-provider":"Local","token":""}`,
		},
		{
			name: "remote provider with a token and extra fields",
			body: `{"meshery-provider":"Meshery","token":"abc.def.ghi","extra":1}`,
			want: `{"meshery-provider":"Meshery","token":"abc.def.ghi"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authPath, logs, _ := setupLoginEnv(t, respond(http.StatusOK, "application/json", tt.body), existingAuth, 0o600)

			if err := runLogin(t); err != nil {
				t.Fatalf("expected login to succeed, got %v", err)
			}

			got, err := os.ReadFile(authPath)
			if err != nil {
				t.Fatalf("unable to read auth file: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("auth file = %q, want %q", got, tt.want)
			}
			if !strings.Contains(logs.String(), "authenticated") {
				t.Errorf("expected success message after persisting the token, got %q", logs.String())
			}
		})
	}
}

// TestLoginRejectsInvalidTokenResponse checks that every unsuccessful or
// malformed /api/token response fails the login, leaves the existing auth
// file byte-identical, does not report success and does not follow redirects.
func TestLoginRejectsInvalidTokenResponse(t *testing.T) {
	tests := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"401 JSON", respond(http.StatusUnauthorized, "application/json", `{"error":"unauthorized"}`)},
		{"403 text", respond(http.StatusForbidden, "text/plain", "forbidden")},
		{"500 JSON", respond(http.StatusInternalServerError, "application/json", `{"message":"boom"}`)},
		{"200 HTML", respond(http.StatusOK, "text/html; charset=utf-8", "<html><body>Sign in</body></html>")},
		{"200 malformed JSON", respond(http.StatusOK, "application/json", `{"token": "abc"`)},
		{"200 empty body", respond(http.StatusOK, "", "")},
		{"200 JSON without either field", respond(http.StatusOK, "application/json", `{"status":"error"}`)},
		{"200 JSON without provider", respond(http.StatusOK, "application/json", `{"token":"abc"}`)},
		{"200 JSON with empty provider", respond(http.StatusOK, "application/json", `{"meshery-provider":"","token":"abc"}`)},
		{"200 JSON without token", respond(http.StatusOK, "application/json", `{"meshery-provider":"Local"}`)},
		{"200 JSON with null token", respond(http.StatusOK, "application/json", `{"meshery-provider":"Local","token":null}`)},
		{"200 JSON with non-string provider", respond(http.StatusOK, "application/json", `{"meshery-provider":7,"token":""}`)},
		{"200 JSON with non-string token", respond(http.StatusOK, "application/json", `{"meshery-provider":"Local","token":["x"]}`)},
		{"200 JSON array", respond(http.StatusOK, "application/json", `[]`)},
		{"connection closed", func(w http.ResponseWriter, _ *http.Request) {
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				_ = conn.Close()
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authPath, logs, redirectHits := setupLoginEnv(t, tt.handler, existingAuth, 0o600)

			assertLoginRejected(t, runLogin(t), authPath, logs, redirectHits)
		})
	}
}

// TestLoginRejectsRedirectAsInvalidToken checks that every 3xx from /api/token
// is reported as an invalid token rather than followed or echoed back.
func TestLoginRejectsRedirectAsInvalidToken(t *testing.T) {
	codes := []int{
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect,
	}

	for _, code := range codes {
		t.Run(strconv.Itoa(code), func(t *testing.T) {
			authPath, logs, redirectHits := setupLoginEnv(t, redirectTo(code), existingAuth, 0o600)

			err := runLogin(t)
			assertLoginRejected(t, err, authPath, logs, redirectHits)
			if got := meshkiterrors.GetCode(err); got != utils.ErrInvalidTokenCode {
				t.Errorf("error code = %q, want %q (%v)", got, utils.ErrInvalidTokenCode, err)
			}
		})
	}
}

// assertLoginRejected checks that a failed login returned an error, left the
// existing auth file byte-identical, did not report success and did not
// follow a redirect.
func assertLoginRejected(t *testing.T, loginErr error, authPath string, logs *bytes.Buffer, redirectHits *atomic.Int32) {
	t.Helper()
	if loginErr == nil {
		t.Fatal("expected login to fail, got nil")
	}

	got, err := os.ReadFile(authPath)
	if err != nil {
		t.Fatalf("unable to read auth file: %v", err)
	}
	if string(got) != existingAuth {
		t.Errorf("existing credentials were replaced: got %q, want %q", got, existingAuth)
	}
	if strings.Contains(logs.String(), "authenticated") {
		t.Errorf("failed login reported success: %q", logs.String())
	}
	if n := redirectHits.Load(); n != 0 {
		t.Errorf("token request followed a redirect %d time(s)", n)
	}
}

// TestLoginDoesNotReportSuccessWhenTokenCannotBeWritten checks that a valid
// response whose token cannot be persisted fails without claiming success.
func TestLoginDoesNotReportSuccessWhenTokenCannotBeWritten(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("read-only file permissions are not enforced in this environment")
	}

	valid := respond(http.StatusOK, "application/json", `{"meshery-provider":"Local","token":""}`)
	authPath, logs, _ := setupLoginEnv(t, valid, existingAuth, 0o400)

	if err := runLogin(t); err == nil {
		t.Fatal("expected login to fail when the auth file cannot be written, got nil")
	}

	got, err := os.ReadFile(authPath)
	if err != nil {
		t.Fatalf("unable to read auth file: %v", err)
	}
	if string(got) != existingAuth {
		t.Errorf("auth file = %q, want %q", got, existingAuth)
	}
	if strings.Contains(logs.String(), "authenticated") {
		t.Errorf("login reported success although the token was not written: %q", logs.String())
	}
}
