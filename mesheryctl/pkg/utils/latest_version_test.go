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
	"bytes"
	"io"
	"net/http"
	"sync/atomic"
	"testing"
)

type mockVersionTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockVersionTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestGetLatestVersionForMesheryctl(t *testing.T) {
	origTransport := http.DefaultTransport
	t.Cleanup(func() {
		http.DefaultTransport = origTransport
	})

	t.Run("successfully retrieves latest version and closes response body", func(t *testing.T) {
		var bodyClosed atomic.Bool

		http.DefaultTransport = &mockVersionTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body: &trackingBody{
						ReadCloser: io.NopCloser(bytes.NewBufferString("v0.7.20\n")),
						closed:     &bodyClosed,
					},
					Header: make(http.Header),
				}, nil
			},
		}

		version, err := GetLatestVersionForMesheryctl()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if version != "v0.7.20" {
			t.Fatalf("expected version 'v0.7.20', got: %q", version)
		}
		if !bodyClosed.Load() {
			t.Fatal("expected response body to be closed on success, but it was not")
		}
	})

	t.Run("returns error on non-200 status code and closes response body", func(t *testing.T) {
		var bodyClosed atomic.Bool

		http.DefaultTransport = &mockVersionTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadGateway,
					Body: &trackingBody{
						ReadCloser: io.NopCloser(bytes.NewBufferString("<html>502 Bad Gateway</html>")),
						closed:     &bodyClosed,
					},
					Header: make(http.Header),
				}, nil
			},
		}

		version, err := GetLatestVersionForMesheryctl()
		if err == nil {
			t.Fatalf("expected error for 502 status, got version: %q", version)
		}
		if !bodyClosed.Load() {
			t.Fatal("expected response body to be closed on non-200 status, but it was not")
		}
	})
}
