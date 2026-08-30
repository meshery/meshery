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
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
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

func TestUpdateAuthDetails(t *testing.T) {
	configPath := CopyMeshconfigFixture(t, "TestConfig.yaml")
	viper.Reset()
	viper.SetConfigFile(configPath)
	DefaultConfigPath = configPath
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("unable to read configuration from %v: %v", viper.ConfigFileUsed(), err)
	}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}
	baseURL := mctlCfg.GetBaseMesheryURL()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	tempDir := t.TempDir()
	tokenFile := filepath.Join(tempDir, "token.json")
	tokenContent := `{"token": "test-token", "meshery-provider": "Local"}`
	if err := os.WriteFile(tokenFile, []byte(tokenContent), 0600); err != nil {
		t.Fatalf("failed to create dummy token file: %v", err)
	}

	t.Run("Given client.Do returns nil response and error, When UpdateAuthDetails is called, Then it returns error without panic", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder("GET", baseURL+"/api/user/token",
			httpmock.NewErrorResponder(fmt.Errorf("connection refused")))

		err := UpdateAuthDetails(tokenFile)
		if err == nil {
			t.Fatalf("expected error from UpdateAuthDetails, got nil")
		}
		if !strings.Contains(err.Error(), "error dispatching there request") {
			t.Fatalf("expected error containing 'error dispatching there request', got %q", err.Error())
		}
	})

	t.Run("Given client.Do returns success response, When UpdateAuthDetails is called, Then token is updated successfully", func(t *testing.T) {
		httpmock.Reset()
		updatedToken := `{"token": "refreshed-token", "meshery-provider": "Local"}`
		httpmock.RegisterResponder("GET", baseURL+"/api/user/token",
			httpmock.NewStringResponder(200, updatedToken))

		err := UpdateAuthDetails(tokenFile)
		if err != nil {
			t.Fatalf("expected no error from UpdateAuthDetails, got %v", err)
		}

		data, err := os.ReadFile(tokenFile)
		if err != nil {
			t.Fatalf("failed to read updated token file: %v", err)
		}
		if string(data) != updatedToken {
			t.Fatalf("expected updated token %q, got %q", updatedToken, string(data))
		}
	})

	t.Run("Given client.Do returns HTML response, When UpdateAuthDetails is called, Then it returns invalid body error", func(t *testing.T) {
		httpmock.Reset()
		resp := httpmock.NewStringResponse(200, "<html><body>Login</body></html>")
		resp.Header.Set("Content-Type", "text/html")
		httpmock.RegisterResponder("GET", baseURL+"/api/user/token",
			httpmock.ResponderFromResponse(resp))

		err := UpdateAuthDetails(tokenFile)
		if err == nil {
			t.Fatalf("expected error from UpdateAuthDetails for HTML body, got nil")
		}
		if err.Error() != "invalid body" {
			t.Fatalf("expected 'invalid body' error, got %q", err.Error())
		}
	})

	t.Run("Given non-existent token file, When UpdateAuthDetails is called, Then it returns error reading token", func(t *testing.T) {
		httpmock.Reset()
		err := UpdateAuthDetails(filepath.Join(tempDir, "non_existent.json"))
		if err == nil {
			t.Fatalf("expected error for non-existent token file, got nil")
		}
		if !strings.Contains(err.Error(), "could not read token") {
			t.Fatalf("expected error containing 'could not read token', got %q", err.Error())
		}
	})
}
