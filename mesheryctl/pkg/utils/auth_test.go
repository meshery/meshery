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

package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
)

// testcases for auth.go
func TestAuth(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "A simple server only for testing")
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	// testcases for GetTokenLocation(token config.Token) (string, error)
	t.Run("GetTokenLocation", func(t *testing.T) {
		token := config.Token{
			Name:     "test",
			Location: "test",
		}
		_, err := GetTokenLocation(token)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("MakeRequest", func(t *testing.T) {
		_, err := MakeRequest(req)
		if err != nil {
			t.Fatal(err)
		}
	})
	//@Aisuko Need a token file to do other testings
}

func TestProviderUnmarshalJSON(t *testing.T) {
	t.Run("Given canonical camelCase provider fields, When unmarshaled, Then it populates Provider correctly", func(t *testing.T) {
		payload := []byte(`{"providerUrl":"https://cloud.meshery.io","providerName":"Meshery"}`)
		provider := Provider{}

		if err := json.Unmarshal(payload, &provider); err != nil {
			t.Fatalf("failed to unmarshal provider payload: %v", err)
		}

		if provider.ProviderName != "Meshery" {
			t.Fatalf("expected provider name Meshery, got %q", provider.ProviderName)
		}
		if provider.ProviderURL != "https://cloud.meshery.io" {
			t.Fatalf("expected provider URL https://cloud.meshery.io, got %q", provider.ProviderURL)
		}
	})

	t.Run("Given legacy snake_case provider fields, When unmarshaled, Then it populates Provider correctly", func(t *testing.T) {
		payload := []byte(`{"provider_url":"https://cloud.meshery.io","provider_name":"Meshery"}`)
		provider := Provider{}

		if err := json.Unmarshal(payload, &provider); err != nil {
			t.Fatalf("failed to unmarshal provider payload: %v", err)
		}

		if provider.ProviderName != "Meshery" {
			t.Fatalf("expected provider name Meshery, got %q", provider.ProviderName)
		}
		if provider.ProviderURL != "https://cloud.meshery.io" {
			t.Fatalf("expected provider URL https://cloud.meshery.io, got %q", provider.ProviderURL)
		}
	})
}

func TestMakeRequest_Status302AndHTMLResponseHandling(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	t.Run("MakeRequest returns ErrInvalidToken on 302 redirect and closes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusFound)
		}))
		defer server.Close()

		transport := &trackingTransport{base: http.DefaultTransport}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("unexpected error creating request: %v", err)
		}

		resp, err := MakeRequest(req)
		if resp != nil {
			t.Errorf("expected nil response on 302, got %v", resp)
		}
		if err == nil {
			t.Fatal("expected error on 302, got nil")
		}
		expectedErr := ErrInvalidToken()
		if errors.GetCode(err) != errors.GetCode(expectedErr) {
			t.Errorf("expected error code %s, got %s", errors.GetCode(expectedErr), errors.GetCode(err))
		}
		if transport.body == nil || !transport.body.closed {
			t.Fatal("expected response body to be closed on 302 redirect")
		}
	})

	t.Run("MakeRequest returns ErrUnauthenticated on HTML Content-Type and closes body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, "<html><body>Login</body></html>")
		}))
		defer server.Close()

		transport := &trackingTransport{base: http.DefaultTransport}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("unexpected error creating request: %v", err)
		}

		resp, err := MakeRequest(req)
		if resp != nil {
			t.Errorf("expected nil response on HTML content, got %v", resp)
		}
		if err == nil {
			t.Fatal("expected error on HTML response, got nil")
		}
		expectedErr := ErrUnauthenticated()
		if errors.GetCode(err) != errors.GetCode(expectedErr) {
			t.Errorf("expected error code %s, got %s", errors.GetCode(expectedErr), errors.GetCode(err))
		}
		if transport.body == nil || !transport.body.closed {
			t.Fatal("expected response body to be closed on HTML content-type response")
		}
	})

	t.Run("MakeRequest returns ErrMesheryServerInternalError on HTML 500 response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "<html><body>500 Internal Server Error</body></html>")
		}))
		defer server.Close()

		req, err := http.NewRequest(http.MethodGet, server.URL, nil)
		if err != nil {
			t.Fatalf("unexpected error creating request: %v", err)
		}

		resp, err := MakeRequest(req)
		if resp != nil {
			t.Errorf("expected nil response on 500, got %v", resp)
		}
		if err == nil {
			t.Fatal("expected error on 500, got nil")
		}
		if errors.GetCode(err) != ErrMesheryServerInternalErrorCode {
			t.Errorf("expected error code %s, got %s", ErrMesheryServerInternalErrorCode, errors.GetCode(err))
		}
	})
}

func newTestMctCfg(endpoint string) *config.MesheryCtlConfig {
	return &config.MesheryCtlConfig{
		CurrentContext: "test",
		Contexts:       map[string]config.Context{"test": {Endpoint: endpoint}},
	}
}

func TestGetProviderInfo_BodyClosed(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	t.Run("closes body on 200 OK success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `{"Local":{"provider_name":"Local","provider_type":"local"}}`)
		}))
		defer server.Close()

		transport := installTrackingTransport(t)
		mctCfg := newTestMctCfg(server.URL)

		providers, err := GetProviderInfo(mctCfg)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		if _, ok := providers["Local"]; !ok {
			t.Fatalf("expected Local provider, got: %v", providers)
		}
		assertBodyClosed(t, transport, "expected response body to be closed on success")
	})

	t.Run("closes body on decode error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, `invalid-json`)
		}))
		defer server.Close()

		transport := installTrackingTransport(t)
		mctCfg := newTestMctCfg(server.URL)

		_, err := GetProviderInfo(mctCfg)
		if err == nil {
			t.Fatal("expected error decoding invalid JSON, got nil")
		}
		assertBodyClosed(t, transport, "expected response body to be closed on decode error")
	})
}

func TestChooseDirectProvider(t *testing.T) {
	// A server pinned to one provider, reported under a key that differs from
	// its display name - the shape playground.meshery.io returns.
	single := map[string]Provider{
		"cloud.meshery.io": {ProviderName: "Meshery", ProviderURL: "https://cloud.meshery.io"},
	}

	t.Run("matches display name case-insensitively", func(t *testing.T) {
		got, key, err := chooseDirectProvider(single, "meshery")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "cloud.meshery.io" {
			t.Fatalf("got key %q, want %q", key, "cloud.meshery.io")
		}
		if got.ProviderURL != "https://cloud.meshery.io" {
			t.Fatalf("got provider %+v", got)
		}
	})

	t.Run("matches the registry key", func(t *testing.T) {
		got, key, err := chooseDirectProvider(single, "cloud.meshery.io")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "cloud.meshery.io" {
			t.Fatalf("got key %q, want %q", key, "cloud.meshery.io")
		}
		if got.ProviderName != "Meshery" {
			t.Fatalf("got provider %+v", got)
		}
	})

	// An exact key match must win over a display-name match, deterministically:
	// with key "cloud-a" named "Meshery" beside key "cloud-b" named "cloud-a",
	// --provider cloud-a selects cloud-a, not whichever entry a map traversal
	// happens to visit first.
	t.Run("exact key beats another entry's display name", func(t *testing.T) {
		provs := map[string]Provider{
			"cloud-a": {ProviderName: "Meshery", ProviderURL: "https://cloud-a.example"},
			"cloud-b": {ProviderName: "cloud-a", ProviderURL: "https://cloud-b.example"},
		}
		for i := 0; i < 20; i++ {
			got, key, err := chooseDirectProvider(provs, "cloud-a")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if key != "cloud-a" || got.ProviderURL != "https://cloud-a.example" {
				t.Fatalf("attempt %d: got key %q provider %+v", i, key, got)
			}
		}
	})

	// Regression: two providers sharing a display name used to resolve to
	// whichever entry a map traversal visited first. --provider with that
	// display name must now fail and ask for a registration key.
	t.Run("ambiguous display name returns an error naming the keys", func(t *testing.T) {
		provs := map[string]Provider{
			"cloud-a": {ProviderName: "Acme", ProviderURL: "https://cloud-a.example"},
			"cloud-b": {ProviderName: "Acme", ProviderURL: "https://cloud-b.example"},
		}
		for i := 0; i < 20; i++ {
			got, key, err := chooseDirectProvider(provs, "Acme")
			if err == nil {
				t.Fatalf("attempt %d: expected an error, got key %q provider %+v", i, key, got)
			}
			if code := errors.GetCode(err); code != ErrProviderAmbiguousCode {
				t.Fatalf("attempt %d: got error code %q, want %q", i, code, ErrProviderAmbiguousCode)
			}
			if key != "" {
				t.Fatalf("attempt %d: expected empty key, got %q", i, key)
			}
			if got != (Provider{}) {
				t.Fatalf("attempt %d: expected empty provider, got %+v", i, got)
			}
		}
	})

	// A display name nobody else shares still resolves.
	t.Run("unique display name still resolves", func(t *testing.T) {
		provs := map[string]Provider{
			"cloud-a": {ProviderName: "Acme", ProviderURL: "https://cloud-a.example"},
			"cloud-b": {ProviderName: "Other", ProviderURL: "https://cloud-b.example"},
		}
		got, key, err := chooseDirectProvider(provs, "acme")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if key != "cloud-a" || got.ProviderURL != "https://cloud-a.example" {
			t.Fatalf("got key %q provider %+v", key, got)
		}
	})

	// Regression: an unmatched name used to return provArray[1], panicking
	// with "index out of range [1] with length 1" on a single-provider server.
	for name, provs := range map[string]map[string]Provider{
		"single provider": single,
		"no providers":    {},
	} {
		t.Run("unknown provider on "+name+" returns a structured error", func(t *testing.T) {
			got, key, err := chooseDirectProvider(provs, "does-not-exist")
			if err == nil {
				t.Fatal("expected an error")
			}
			if code := errors.GetCode(err); code != ErrProviderNotFoundCode {
				t.Fatalf("got error code %q, want %q", code, ErrProviderNotFoundCode)
			}
			if key != "" {
				t.Fatalf("expected empty key, got %q", key)
			}
			if got != (Provider{}) {
				t.Fatalf("expected zero Provider, got %+v", got)
			}
		})
	}
}

func TestInitiateLogin_PersistsRegistrationKey(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	// Two providers share the display name "Meshery". The server answers
	// /api/token with that display name, as RemoteProvider.ExtractToken
	// does (meshery-provider is l.Name()), and login persists the response
	// unchanged for later requests. Login must persist the registration
	// key instead, or later requests would send the shared display name
	// and route nondeterministically. The providers carry no URL so the
	// test exercises the token exchange without the browser flow.
	var tokenCookie string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/providers":
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"cloud-a":{"providerName":"Meshery"},"cloud-b":{"providerName":"Meshery"}}`)
		case "/api/token":
			for _, c := range r.Cookies() {
				if c.Name == "meshery-provider" {
					tokenCookie = c.Value
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, `{"meshery-provider":"Meshery","meshery-token":"tok"}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	mctCfg := newTestMctCfg(server.URL)
	data, err := InitiateLogin(mctCfg, "cloud-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenCookie != "cloud-b" {
		t.Fatalf("/api/token got meshery-provider cookie %q, want %q", tokenCookie, "cloud-b")
	}

	var persisted map[string]string
	if err := json.Unmarshal(data, &persisted); err != nil {
		t.Fatalf("login returned invalid JSON: %v", err)
	}
	if persisted["meshery-provider"] != "cloud-b" {
		t.Fatalf("persisted meshery-provider %q, want %q", persisted["meshery-provider"], "cloud-b")
	}
	if persisted["meshery-token"] != "tok" {
		t.Fatalf("persisted meshery-token %q, want %q", persisted["meshery-token"], "tok")
	}

	// A later authenticated request built from the persisted token must
	// carry the registration key, not the shared display name.
	tokenPath := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(tokenPath, data, 0o600); err != nil {
		t.Fatalf("unexpected error writing token file: %v", err)
	}
	req, err := http.NewRequest(http.MethodGet, server.URL+"/api/user", nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}
	if err := AddAuthDetails(req, tokenPath); err != nil {
		t.Fatalf("unexpected error adding auth details: %v", err)
	}
	var providerCookie string
	for _, c := range req.Cookies() {
		if c.Name == "meshery-provider" {
			providerCookie = c.Value
		}
	}
	if providerCookie != "cloud-b" {
		t.Fatalf("authenticated request got meshery-provider cookie %q, want %q", providerCookie, "cloud-b")
	}
}

func TestWriteTokenFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	// A token file left behind by an older mesheryctl is world-readable.
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := WriteTokenFile(path, []byte(`{"token":"t"}`)); err != nil {
		t.Fatalf("WriteTokenFile: %v", err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	// Windows only tracks the read-only attribute, so a mode check is
	// meaningless there; the content check below still applies.
	if runtime.GOOS != "windows" {
		if mode := info.Mode().Perm(); mode != 0o600 {
			t.Fatalf("token file mode = %o, want 600", mode)
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"token":"t"}` {
		t.Fatalf("token file content = %q", data)
	}
}

// A context passed with --config that names a token its file does not define
// must get that token added to the same file, and later requests must keep
// using that file's server rather than the default meshconfig's.
func TestGetCurrentAuthTokenAddsMissingTokenToActiveConfig(t *testing.T) {
	mesheryDir := IsolateMesheryHome(t)

	defaultConfig := filepath.Join(mesheryDir, "config.yaml")
	defaultContents := []byte(`contexts:
  local:
    endpoint: http://127.0.0.1:19999
    token: default
    platform: docker
current-context: local
tokens:
- name: default
  location: auth.json
`)
	if err := os.WriteFile(defaultConfig, defaultContents, 0o600); err != nil {
		t.Fatal(err)
	}

	activeConfig := filepath.Join(t.TempDir(), "custom.yaml")
	if err := os.WriteFile(activeConfig, []byte(`contexts:
  staging:
    endpoint: http://127.0.0.1:19181
    token: missingtok
    platform: docker
current-context: staging
tokens:
- name: default
  location: auth.json
`), 0o600); err != nil {
		t.Fatal(err)
	}

	origDefaultConfigPath, origCfgFile, origAuthConfigFile := DefaultConfigPath, CfgFile, AuthConfigFile
	t.Cleanup(func() {
		DefaultConfigPath, CfgFile, AuthConfigFile = origDefaultConfigPath, origCfgFile, origAuthConfigFile
		viper.Reset()
	})
	DefaultConfigPath = defaultConfig
	CfgFile = activeConfig
	AuthConfigFile = filepath.Join(mesheryDir, "auth.json")

	viper.Reset()
	viper.SetConfigFile(activeConfig)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	location, err := GetCurrentAuthToken()
	if err != nil {
		t.Fatalf("GetCurrentAuthToken: %v", err)
	}
	if location != AuthConfigFile {
		t.Errorf("token location = %q, want %q", location, AuthConfigFile)
	}

	if got := viper.ConfigFileUsed(); got != activeConfig {
		t.Errorf("viper now reads %q, want the --config file %q: later requests would go to another server", got, activeConfig)
	}

	gotDefault, err := os.ReadFile(defaultConfig)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotDefault) != string(defaultContents) {
		t.Errorf("default meshconfig was rewritten:\n%s", gotDefault)
	}

	v := viper.New()
	v.SetConfigFile(activeConfig)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	active, err := config.GetMesheryCtl(v)
	if err != nil {
		t.Fatal(err)
	}
	if token, err := active.GetTokenForContext("staging"); err != nil || token.Location != AuthConfigFile {
		t.Errorf("--config file token for context staging = %+v, %v; want location %q", token, err, AuthConfigFile)
	}

	// The entry now exists, so a second command, which starts by reading the
	// --config file afresh, must reuse it rather than fail on adding a duplicate.
	viper.Reset()
	viper.SetConfigFile(activeConfig)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	if _, err := GetCurrentAuthToken(); err != nil {
		t.Errorf("second GetCurrentAuthToken: %v", err)
	}
}
