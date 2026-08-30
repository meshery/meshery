// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetManifestTreeURL(t *testing.T) {
	t.Run("returns error on non-200 HTTP status code (e.g. 403 Rate Limit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		origURL := gitHubBaseURL
		gitHubBaseURL = server.URL
		defer func() { gitHubBaseURL = origURL }()

		_, err := GetManifestTreeURL("v0.6.0")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "403") {
			t.Errorf("expected error message to contain '403', got: %v", err)
		}
		if !strings.Contains(err.Error(), "API rate limit exceeded") {
			t.Errorf("expected error to contain rate limit message, got: %v", err)
		}
		if strings.Contains(err.Error(), "could not find path") {
			t.Errorf("expected error not to be misleading 'could not find path', got: %v", err)
		}
	})

	t.Run("happy path with valid ManifestList", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			respBody := `{
				"sha": "12345",
				"tree": [
					{
						"path": "install/deployment_yamls/k8s",
						"url": "https://api.github.com/repos/meshery/meshery/git/trees/k8s-tree-sha"
					}
				]
			}`
			_, _ = w.Write([]byte(respBody))
		}))
		defer server.Close()

		origURL := gitHubBaseURL
		gitHubBaseURL = server.URL
		defer func() { gitHubBaseURL = origURL }()

		treeURL, err := GetManifestTreeURL("v0.6.0")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "https://api.github.com/repos/meshery/meshery/git/trees/k8s-tree-sha"
		if treeURL != expected {
			t.Errorf("expected treeURL %s, got %s", expected, treeURL)
		}
	})
}

func TestListManifests(t *testing.T) {
	t.Run("returns error on non-200 HTTP status code (e.g. 403 Rate Limit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		_, err := ListManifests(server.URL)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "403") {
			t.Errorf("expected error message to contain '403', got: %v", err)
		}
		if !strings.Contains(err.Error(), "API rate limit exceeded") {
			t.Errorf("expected error to contain rate limit message, got: %v", err)
		}
	})

	t.Run("happy path with valid ManifestList", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			respBody := `{
				"sha": "12345",
				"tree": [
					{
						"path": "install/deployment_yamls/k8s/meshery-deployment.yaml",
						"type": "blob"
					}
				]
			}`
			_, _ = w.Write([]byte(respBody))
		}))
		defer server.Close()

		manifests, err := ListManifests(server.URL)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(manifests) != 1 {
			t.Fatalf("expected 1 manifest, got %d", len(manifests))
		}
		if manifests[0].Path != "install/deployment_yamls/k8s/meshery-deployment.yaml" {
			t.Errorf("unexpected manifest path: %s", manifests[0].Path)
		}
	})
}

