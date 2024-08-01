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
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
)

// testcases for auth.go
func TestAuth(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A simple server only for testing")
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
	defer resp.Body.Close()

	t.Run("NewRequest", func(t *testing.T) {
		tests := []struct {
			name        string
			method      string
			url         string
			body        io.Reader
			tokenPath   string
			setup       func()
			expectedErr bool
		}{
			{
				name:      "Valid token path",
				method:    "GET",
				url:       server.URL,
				body:      nil,
				tokenPath: "valid-token-path",
				setup: func() {
					TokenFlag = "valid-token-path"
					err := os.WriteFile("valid-token-path", []byte(`{"token": "valid-token", "meshery-provider": "valid-provider"}`), 0644)
					if err != nil {
						t.Fatalf("failed to write token file: %v", err)
					}
				},
				expectedErr: false,
			},
			{
				name:      "Invalid token path",
				method:    "GET",
				url:       server.URL,
				body:      nil,
				tokenPath: "",
				setup: func() {
					TokenFlag = ""
				},
				expectedErr: true,
			},
			{
				name:      "Token file does not exist",
				method:    "GET",
				url:       server.URL,
				body:      nil,
				tokenPath: "non-existent-token-path",
				setup: func() {
					TokenFlag = "non-existent-token-path"
				},
				expectedErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tt.setup()
				req, err := NewRequest(tt.method, tt.url, tt.body)
				if (err != nil) != tt.expectedErr {
					t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
				}
				if err == nil {
					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						t.Fatalf("Failed to make request: %v", err)
					}
					defer resp.Body.Close()
					defer os.Remove(tt.tokenPath)
					if resp.StatusCode != http.StatusOK {
						t.Errorf("Expected status code 200, got: %v", resp.StatusCode)
					}
				}
			})
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

func TestGetTokenLocation(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user's home directory: %v", err)
	}

	tests := []struct {
		name     string
		token    config.Token
		expected string
	}{
		{
			name: "Token location includes home directory",
			token: config.Token{
				Name:     "test",
				Location: filepath.Join(home, "test-token"),
			},
			expected: filepath.Join(home, "test-token"),
		},
		{
			name: "Token location does not include home directory",
			token: config.Token{
				Name:     "test",
				Location: "test-token",
			},
			expected: filepath.Join(MesheryFolder, "test-token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, err := GetTokenLocation(tt.token)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if location != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, location)
			}
		})
	}
}

func TestReadToken(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		tokenFile    string
		tokenContent string
		expected     map[string]string
		expectedErr  bool
	}{
		{
			name:         "Valid token file",
			tokenFile:    "valid-token.json",
			tokenContent: `{"token": "valid-token", "meshery-provider": "valid-provider"}`,
			expected:     map[string]string{"token": "valid-token", "meshery-provider": "valid-provider"},
			expectedErr:  false,
		},
		{
			name:         "Invalid token file",
			tokenFile:    "invalid-token.json",
			tokenContent: `{"token": "valid-token", "meshery-provider": "valid-provider"`,
			expected:     nil,
			expectedErr:  true,
		},
		{
			name:         "Non-existent token file",
			tokenFile:    "non-existent-token.json",
			tokenContent: "",
			expected:     nil,
			expectedErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenFilePath := filepath.Join(tempDir, tt.tokenFile)
			if tt.tokenContent != "" {
				err := os.WriteFile(tokenFilePath, []byte(tt.tokenContent), 0644)
				if err != nil {
					t.Fatalf("Failed to write token file: %v", err)
				}
			}

			token, err := ReadToken(tokenFilePath)
			if (err != nil) != tt.expectedErr {
				t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
			}
			if !tt.expectedErr && !reflect.DeepEqual(token, tt.expected) {
				t.Errorf("Expected token: %v, got: %v", tt.expected, token)
			}
		})
	}
}

func TestChooseDirectProvider(t *testing.T) {
	providers := map[string]Provider{
		"local": {
			ProviderName: "Local Provider",
			ProviderURL:  "http://localhost:9081",
		},
		"remote": {
			ProviderName: "Remote Provider",
			ProviderURL:  "http://remote.meshery.io",
		},
	}

	testCases := []struct {
		name          string
		option        string
		expectedName  string
		expectedError bool
	}{
		{
			name:          "Valid provider name",
			option:        "Local Provider",
			expectedName:  "Local Provider",
			expectedError: false,
		},
		{
			name:          "Invalid provider name",
			option:        "Invalid Provider",
			expectedName:  "",
			expectedError: true,
		},
		{
			name:          "Case-insensitive match",
			option:        "remote provider",
			expectedName:  "Remote Provider",
			expectedError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := chooseDirectProvider(providers, tc.option)
			if tc.expectedError {
				if err == nil {
					t.Fatalf("Expected an error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}
				if provider.ProviderName != tc.expectedName {
					t.Fatalf("Expected provider name %s, got %s", tc.expectedName, provider.ProviderName)
				}
			}
		})
	}
}

func TestCreateProviderURI(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		host     string
		port     int
		want     string
		wantErr  bool
	}{
		{
			name: "Valid URI",
			provider: Provider{
				ProviderURL: "http://example.com",
			},
			host:    "localhost",
			port:    8080,
			want:    "http://example.com/login?provider_version=v0.3.14&source=bG9jYWxob3N0OjgwODA", // Base64 encoded "localhost:8080"
			wantErr: false,
		},
		{
			name: "Invalid Provider URL",
			provider: Provider{
				ProviderURL: "http://[::1]:namedport",
			},
			host:    "localhost",
			port:    8080,
			want:    "",
			wantErr: true,
		},
		{
			name: "Empty Provider URL",
			provider: Provider{
				ProviderURL: "",
			},
			host:    "localhost",
			port:    8080,
			want:    "/login?provider_version=v0.3.14&source=bG9jYWxob3N0OjgwODA", // Base64 encoded "localhost:8080"
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createProviderURI(tt.provider, tt.host, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("createProviderURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createProviderURI() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTokenObjFromMesheryServer(t *testing.T) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(tokenName)
		if err != nil || tokenCookie.Value != "valid-token" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		providerCookie, err := r.Cookie("meshery-provider")
		if err != nil || providerCookie.Value != "valid-provider" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"token": "mock-token"}`)); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	mctlCfg := &config.MesheryCtlConfig{
		CurrentContext: "default",
		Contexts: map[string]config.Context{
			"default": {
				Endpoint: server.URL,
			},
		},
	}

	t.Run("Valid Token and Provider", func(t *testing.T) {
		data, err := getTokenObjFromMesheryServer(mctlCfg, "valid-provider", "valid-token")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		expected := `{"token": "mock-token"}`
		if string(data) != expected {
			t.Fatalf("Expected %s, got %s", expected, string(data))
		}
	})
}

func TestIsServerRunning(t *testing.T) {

	t.Run("Server is running", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Server is running")
		}))
		defer server.Close()

		serverAddr := strings.TrimPrefix(server.URL, "http://")

		err := IsServerRunning(serverAddr)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
	})

	t.Run("Server is not running", func(t *testing.T) {
		invalidAddr := "localhost:12345"

		err := IsServerRunning(invalidAddr)
		if err == nil {
			t.Fatalf("Expected an error, got nil")
		}
	})
}
