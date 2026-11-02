// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetManifestTreeURL(t *testing.T) {
	t.Run("returns error on non-200 HTTP status code (e.g. 403 Rate Limit)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		origURL := githubBaseURL
		githubBaseURL = server.URL
		defer func() { githubBaseURL = origURL }()

		_, err := GetManifestTreeURL("v0.6.0")
		expectedURL := githubBaseURL + "/repos/" + "meshery" + "/" + "meshery" + "/git/trees/" + "v0.6.0" + "?recursive=1"
		expectedErr := ErrGitHubAPIResponse(http.StatusForbidden, expectedURL, `{"message": "API rate limit exceeded"}`)
		AssertMeshkitErrorsEqual(t, err, expectedErr)
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

		origURL := githubBaseURL
		githubBaseURL = server.URL
		defer func() { githubBaseURL = origURL }()

		treeURL, err := GetManifestTreeURL("v0.6.0")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
		expected := "https://api.github.com/repos/meshery/meshery/git/trees/k8s-tree-sha"
		if treeURL != expected {
			t.Errorf("expected treeURL %q, got %q", expected, treeURL)
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
		expectedErr := ErrGitHubAPIResponse(http.StatusForbidden, server.URL, `{"message": "API rate limit exceeded"}`)
		AssertMeshkitErrorsEqual(t, err, expectedErr)
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
			t.Fatalf("expected no error, got: %v", err)
		}
		if len(manifests) != 1 {
			t.Fatalf("expected 1 manifest, got %d", len(manifests))
		}
		expectedPath := "install/deployment_yamls/k8s/meshery-deployment.yaml"
		if manifests[0].Path != expectedPath {
			t.Errorf("expected manifest path %q, got %q", expectedPath, manifests[0].Path)
		}
	})
}

func TestGetCleanPodName(t *testing.T) {
	tests := []struct {
		name     string
		podName  string
		expected string
	}{
		{
			name:     "standard deployment pod name",
			podName:  "meshery-broker-7d9f8c9b47-x8vqp",
			expected: "meshery-broker",
		},
		{
			name:     "two segments, alphabetic second segment joins both",
			podName:  "brokerpod-abc",
			expected: "brokerpod-abc",
		},
		{
			name:     "statefulset pod name, numeric ordinal falls back to first segment",
			podName:  "broker-0",
			expected: "broker",
		},
		{
			name:     "single word pod name with no hyphen",
			podName:  "meshery",
			expected: "meshery",
		},
		{
			name:     "empty pod name",
			podName:  "",
			expected: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Names with fewer than two hyphen-separated segments used to panic with
			// "slice bounds out of range"; GetCleanPodName must now return cleanly.
			got := GetCleanPodName(tc.podName)
			if got != tc.expected {
				t.Errorf("GetCleanPodName(%q) = %q, want %q", tc.podName, got, tc.expected)
			}
		})
	}
}

type readErrorBody struct {
	closed bool
}

func (b *readErrorBody) Read([]byte) (int, error) {
	return 0, errors.New("read failed")
}

func (b *readErrorBody) Close() error {
	b.closed = true
	return nil
}

type readErrorTransport struct {
	body *readErrorBody
}

func (t *readErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.body = &readErrorBody{}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       t.body,
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestResolveDockerComposeBinaryURL(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	t.Run("falls back to default version and closes response body on read error", func(t *testing.T) {
		origURL := dockerComposeWebURL
		dockerComposeWebURL = "http://read-error.test"
		defer func() { dockerComposeWebURL = origURL }()

		transport := &readErrorTransport{}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		got := resolveDockerComposeBinaryURL()
		expected := dockerComposeBinaryURL + defaultDockerComposeVersion
		if got != expected {
			t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
		}
		if transport.body == nil || !transport.body.closed {
			t.Errorf("expected response body to be closed on read-error path")
		}
	})

	t.Run("success returns release binary url and closes response body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"tag_name": "v2.24.5"}`))
		}))
		defer server.Close()

		origURL := dockerComposeWebURL
		dockerComposeWebURL = server.URL
		defer func() { dockerComposeWebURL = origURL }()

		transport := &trackingTransport{base: http.DefaultTransport}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		got := resolveDockerComposeBinaryURL()
		expected := dockerComposeBinaryURL + "v2.24.5/docker-compose"
		if got != expected {
			t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
		}

		if transport.body == nil || !transport.body.closed {
			t.Errorf("expected response body to be closed on success path")
		}
	})

	for _, tc := range []struct {
		name string
		body string
	}{
		{"missing tag_name", `{}`},
		{"empty tag_name", `{"tag_name": ""}`},
	} {
		t.Run("falls back to default version on "+tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer server.Close()

			origURL := dockerComposeWebURL
			dockerComposeWebURL = server.URL
			defer func() { dockerComposeWebURL = origURL }()

			transport := &trackingTransport{base: http.DefaultTransport}
			previousTransport := http.DefaultTransport
			http.DefaultTransport = transport
			t.Cleanup(func() { http.DefaultTransport = previousTransport })

			got := resolveDockerComposeBinaryURL()
			expected := dockerComposeBinaryURL + defaultDockerComposeVersion
			if got != expected {
				t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
			}

			if transport.body == nil || !transport.body.closed {
				t.Errorf("expected response body to be closed on %s path", tc.name)
			}
		})
	}

	t.Run("falls back to default version and closes response body on non-200 status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"message": "API rate limit exceeded"}`))
		}))
		defer server.Close()

		origURL := dockerComposeWebURL
		dockerComposeWebURL = server.URL
		defer func() { dockerComposeWebURL = origURL }()

		transport := &trackingTransport{base: http.DefaultTransport}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		got := resolveDockerComposeBinaryURL()
		expected := dockerComposeBinaryURL + defaultDockerComposeVersion
		if got != expected {
			t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
		}

		if transport.body == nil || !transport.body.closed {
			t.Errorf("expected response body to be closed on non-200 status code path")
		}
	})

	t.Run("falls back to default version and closes response body on malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"invalid": "json"`))
		}))
		defer server.Close()

		origURL := dockerComposeWebURL
		dockerComposeWebURL = server.URL
		defer func() { dockerComposeWebURL = origURL }()

		transport := &trackingTransport{base: http.DefaultTransport}
		previousTransport := http.DefaultTransport
		http.DefaultTransport = transport
		t.Cleanup(func() { http.DefaultTransport = previousTransport })

		got := resolveDockerComposeBinaryURL()
		expected := dockerComposeBinaryURL + defaultDockerComposeVersion
		if got != expected {
			t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
		}

		if transport.body == nil || !transport.body.closed {
			t.Errorf("expected response body to be closed on malformed json path")
		}
	})

	t.Run("falls back to default version on unreachable endpoint", func(t *testing.T) {
		origURL := dockerComposeWebURL
		dockerComposeWebURL = "http://127.0.0.1:0/unreachable"
		defer func() { dockerComposeWebURL = origURL }()

		got := resolveDockerComposeBinaryURL()
		expected := dockerComposeBinaryURL + defaultDockerComposeVersion
		if got != expected {
			t.Errorf("resolveDockerComposeBinaryURL() = %q, want %q", got, expected)
		}
	})
}
