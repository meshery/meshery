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
	"github.com/meshery/meshkit/errors"
)

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
		resp, err := MakeRequest(req)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close()
	})
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

func TestMakeRequestSuccess(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		body     string
		ctype    string
	}{
		{
			name:   "200 OK",
			status: http.StatusOK,
			body:   `{"status":"ok"}`,
			ctype:  "application/json",
		},
		{
			name:   "201 Created",
			status: http.StatusCreated,
			body:   `{"id":"abc"}`,
			ctype:  "application/json",
		},
		{
			name:   "204 No Content",
			status: http.StatusNoContent,
			body:   "",
			ctype:  "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.ctype)
				w.WriteHeader(tt.status)
				if tt.body != "" {
					_, _ = fmt.Fprint(w, tt.body)
				}
			}
			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := MakeRequest(req)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			_ = resp.Body.Close()
		})
	}
}

func TestMakeRequestErrors(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		ctype      string
		wantErr error
	}{
		{
			name:       "302 redirect returns invalid token",
			status:     http.StatusFound,
			body:       "",
			ctype:      "application/json",
			wantErr:    ErrInvalidToken(),
		},
		{
			name:       "HTML response returns unauthenticated",
			status:     http.StatusOK,
			body:       "<html>login</html>",
			ctype:      "text/html",
			wantErr:    ErrUnauthenticated(),
		},
		{
			name:    "404 returns not found",
			status:  http.StatusNotFound,
			body:    `{"error":"not found"}`,
			ctype:   "application/json",
			wantErr: ErrNotFound(fmt.Errorf(`{"error":"not found"}`)),
		},
		{
			name:    "500 returns internal server error",
			status:  http.StatusInternalServerError,
			body:    `{"error":"server error"}`,
			ctype:   "application/json",
			wantErr: ErrMesheryServerInternalError(fmt.Errorf(`{"error":"server error"}`)),
		},
		{
			name:    "418 teapot returns fail request error",
			status:  http.StatusTeapot,
			body:    `{"error":"teapot"}`,
			ctype:   "application/json",
			wantErr: ErrFailReqStatus(http.StatusTeapot, `{"error":"teapot"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.ctype)
				w.WriteHeader(tt.status)
				if tt.body != "" {
					_, _ = fmt.Fprint(w, tt.body)
				}
			}
			srv := httptest.NewServer(http.HandlerFunc(handler))
			defer srv.Close()

			req, err := http.NewRequest("GET", srv.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = MakeRequest(req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if gotCode, wantCode := errors.GetCode(err), errors.GetCode(tt.wantErr); gotCode != wantCode {
				t.Fatalf("expected error code %s, got %s", wantCode, gotCode)
			}
		})
	}
}
