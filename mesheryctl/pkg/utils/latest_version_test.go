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
	"strings"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

type versionTrackingBody struct {
	io.ReadCloser
	readCalls int
	closed    bool
}

func (b *versionTrackingBody) Read(p []byte) (n int, err error) {
	b.readCalls++
	return b.ReadCloser.Read(p)
}

func (b *versionTrackingBody) Close() error {
	b.closed = true
	return b.ReadCloser.Close()
}

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
		body := &versionTrackingBody{
			ReadCloser: io.NopCloser(bytes.NewBufferString("v0.7.20\n")),
		}

		http.DefaultTransport = &mockVersionTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       body,
					Header:     make(http.Header),
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
		if body.readCalls == 0 {
			t.Fatal("expected response body to be read on success, but it was not")
		}
		if !body.closed {
			t.Fatal("expected response body to be closed on success, but it was not")
		}
	})

	t.Run("returns error on non-200 status code and closes response body without reading", func(t *testing.T) {
		body := &versionTrackingBody{
			ReadCloser: io.NopCloser(bytes.NewBufferString("<html>502 Bad Gateway</html>")),
		}

		http.DefaultTransport = &mockVersionTransport{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusBadGateway,
					Body:       body,
					Header:     make(http.Header),
				}, nil
			},
		}

		version, err := GetLatestVersionForMesheryctl()
		if err == nil {
			t.Fatalf("expected error for 502 status, got version: %q", version)
		}
		if meshkiterrors.GetCode(err) != ErrGetChannelVersionCode {
			t.Fatalf("expected error code %s, got: %s", ErrGetChannelVersionCode, meshkiterrors.GetCode(err))
		}
		if !strings.Contains(err.Error(), "502") || !strings.Contains(err.Error(), "docs.meshery.io") {
			t.Fatalf("expected error identifying docs.meshery.io and 502, got: %v", err)
		}
		if body.readCalls != 0 {
			t.Fatalf("expected 0 reads on 502 status code, got %d reads", body.readCalls)
		}
		if !body.closed {
			t.Fatal("expected response body to be closed on non-200 status, but it was not")
		}
	})
}
