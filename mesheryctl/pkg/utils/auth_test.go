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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
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

// roundTripFunc adapts a custom function to the http.RoundTripper interface for testing.
type roundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip executes the custom transport function f on the incoming HTTP request.
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// redirectTransport redirects HTTP requests to a target test server URL.
type redirectTransport struct {
	target *url.URL
	base   http.RoundTripper
}

// RoundTrip rewrites the scheme and host of req to target and forwards it to the base transport.
func (t *redirectTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqCopy := req.Clone(req.Context())
	reqCopy.URL.Scheme = t.target.Scheme
	reqCopy.URL.Host = t.target.Host
	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(reqCopy)
}

// createTestTokenFile creates a temporary token file containing contents and returns its file path.
func createTestTokenFile(t *testing.T, contents string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "token.json")
	if err := os.WriteFile(path, []byte(contents), 0600); err != nil {
		t.Fatalf("failed to create test token file: %v", err)
	}
	return path
}

// TestUpdateAuthDetails tests UpdateAuthDetails for network errors, successful token updates, HTML responses, and non-existent token files.
func TestUpdateAuthDetails(t *testing.T) {
	viper.Set("current-context", "local")
	viper.Set("contexts.local.endpoint", "http://localhost:9081")
	t.Cleanup(func() {
		viper.Reset()
	})

	validTokenContent := `{"token":"test-token","meshery-provider":"Local"}`

	t.Run("nil response and error", func(t *testing.T) {
		t.Parallel()
		tokenPath := createTestTokenFile(t, validTokenContent)

		expectedErr := errors.New("network failure")
		client := &http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, expectedErr
			}),
		}

		err := updateAuthDetails(client, tokenPath)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "error dispatching the request") {
			t.Errorf("expected error message to contain 'error dispatching the request', got: %v", err)
		}
		if !strings.Contains(err.Error(), "network failure") {
			t.Errorf("expected underlying error 'network failure', got: %v", err)
		}
	})

	t.Run("successful response", func(t *testing.T) {
		t.Parallel()
		tokenPath := createTestTokenFile(t, validTokenContent)

		refreshedTokenJSON := `{"token":"refreshed-token","meshery-provider":"Local"}`
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = fmt.Fprint(w, refreshedTokenJSON)
		}))
		defer srv.Close()

		srvURL, err := url.Parse(srv.URL)
		if err != nil {
			t.Fatalf("failed to parse test server URL: %v", err)
		}

		client := &http.Client{
			Transport: &redirectTransport{
				target: srvURL,
				base:   srv.Client().Transport,
			},
		}

		if err := updateAuthDetails(client, tokenPath); err != nil {
			t.Fatalf("unexpected error updating auth details: %v", err)
		}

		updatedContent, err := os.ReadFile(tokenPath)
		if err != nil {
			t.Fatalf("failed to read updated token file: %v", err)
		}

		if string(updatedContent) != refreshedTokenJSON {
			t.Errorf("expected updated token content %q, got %q", refreshedTokenJSON, string(updatedContent))
		}
	})

	t.Run("HTML unexpected response", func(t *testing.T) {
		t.Parallel()
		tokenPath := createTestTokenFile(t, validTokenContent)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = fmt.Fprint(w, "<html><body>Login Page</body></html>")
		}))
		defer srv.Close()

		srvURL, err := url.Parse(srv.URL)
		if err != nil {
			t.Fatalf("failed to parse test server URL: %v", err)
		}

		client := &http.Client{
			Transport: &redirectTransport{
				target: srvURL,
				base:   srv.Client().Transport,
			},
		}

		err = updateAuthDetails(client, tokenPath)
		if err == nil {
			t.Fatal("expected error for HTML response, got nil")
		}
		if !strings.Contains(err.Error(), "invalid body") {
			t.Errorf("expected error 'invalid body', got: %v", err)
		}
	})

	t.Run("non-existent token file", func(t *testing.T) {
		t.Parallel()
		nonExistentPath := filepath.Join(t.TempDir(), "nonexistent_token.json")

		err := UpdateAuthDetails(nonExistentPath)
		if err == nil {
			t.Fatal("expected error for non-existent token file, got nil")
		}
	})
}
