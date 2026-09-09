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
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
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

func TestChooseDirectProvider(t *testing.T) {
	tests := []struct {
		name           string
		provs          map[string]Provider
		option         string
		expectedProv   Provider
		expectError    bool
		expectedErrSub string
	}{
		{
			name: "single provider with invalid requested provider (regression test)",
			provs: map[string]Provider{
				"None": {
					ProviderName: "None",
					ProviderURL:  "",
				},
			},
			option:         "invalid-provider",
			expectedProv:   Provider{},
			expectError:    true,
			expectedErrSub: "the specified provider 'invalid-provider' is not available. Please try giving correct provider name",
		},
		{
			name: "single provider with valid requested provider",
			provs: map[string]Provider{
				"None": {
					ProviderName: "None",
					ProviderURL:  "",
				},
			},
			option: "None",
			expectedProv: Provider{
				ProviderName: "None",
				ProviderURL:  "",
			},
			expectError: false,
		},
		{
			name: "single provider with valid requested provider case-insensitive",
			provs: map[string]Provider{
				"None": {
					ProviderName: "None",
					ProviderURL:  "",
				},
			},
			option: "none",
			expectedProv: Provider{
				ProviderName: "None",
				ProviderURL:  "",
			},
			expectError: false,
		},
		{
			name: "multiple providers with invalid requested provider",
			provs: map[string]Provider{
				"None": {
					ProviderName: "None",
					ProviderURL:  "",
				},
				"Meshery": {
					ProviderName: "Meshery",
					ProviderURL:  "https://cloud.meshery.io",
				},
			},
			option:         "invalid-provider",
			expectedProv:   Provider{},
			expectError:    true,
			expectedErrSub: "the specified provider 'invalid-provider' is not available. Please try giving correct provider name",
		},
		{
			name: "multiple providers with valid requested provider",
			provs: map[string]Provider{
				"None": {
					ProviderName: "None",
					ProviderURL:  "",
				},
				"Meshery": {
					ProviderName: "Meshery",
					ProviderURL:  "https://cloud.meshery.io",
				},
			},
			option: "Meshery",
			expectedProv: Provider{
				ProviderName: "Meshery",
				ProviderURL:  "https://cloud.meshery.io",
			},
			expectError: false,
		},
		{
			name:           "empty providers with invalid requested provider",
			provs:          map[string]Provider{},
			option:         "Meshery",
			expectedProv:   Provider{},
			expectError:    true,
			expectedErrSub: "the specified provider 'Meshery' is not available. Please try giving correct provider name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProv, err := chooseDirectProvider(tt.provs, tt.option)

			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if tt.expectedErrSub != "" && err.Error() != tt.expectedErrSub {
					t.Fatalf("expected error %q, got %q", tt.expectedErrSub, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if gotProv != tt.expectedProv {
				t.Fatalf("expected provider %+v, got %+v", tt.expectedProv, gotProv)
			}
		})
	}
}
