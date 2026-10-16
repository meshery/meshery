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
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
)

// trackedBody records whether the consumer closed it. An unclosed response
// body holds its connection out of the pool for the life of the process,
// which is why these paths are worth pinning rather than reading.
type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error {
	b.closed = true
	return nil
}

// stubTransport answers every request with a canned response, so the body it
// hands back is one the test holds a reference to.
type stubTransport struct {
	status  int
	header  http.Header
	payload string
	body    *trackedBody
}

func (t *stubTransport) RoundTrip(*http.Request) (*http.Response, error) {
	t.body = &trackedBody{Reader: bytes.NewBufferString(t.payload)}
	header := t.header
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: t.status,
		Header:     header,
		Body:       t.body,
	}, nil
}

// installTransport swaps the default transport, which both http.Get and a
// zero-valued http.Client use, and puts the original back afterwards.
func installTransport(t *testing.T, rt http.RoundTripper) {
	t.Helper()
	original := http.DefaultTransport
	http.DefaultTransport = rt
	t.Cleanup(func() { http.DefaultTransport = original })
}

// testConfig returns a config with a valid current context. An empty one is
// not usable here: GetBaseMesheryURL calls Log.Fatal when the context is
// invalid, which would take the whole test binary down rather than fail.
func testConfig() *config.MesheryCtlConfig {
	return &config.MesheryCtlConfig{
		CurrentContext: "local",
		Contexts: map[string]config.Context{
			"local": {Endpoint: "http://example.invalid"},
		},
	}
}

func TestGetProviderInfoClosesBody(t *testing.T) {
	stub := &stubTransport{status: http.StatusOK, payload: `{}`}
	installTransport(t, stub)

	if _, err := GetProviderInfo(testConfig()); err != nil {
		t.Fatalf("GetProviderInfo: %v", err)
	}

	if !stub.body.closed {
		t.Error("response body was not closed; the connection cannot return to the pool")
	}
}

func TestGetProviderInfoClosesBodyOnDecodeError(t *testing.T) {
	// Malformed JSON takes the early return, which is the easier path to
	// abandon the body on.
	stub := &stubTransport{status: http.StatusOK, payload: `{`}
	installTransport(t, stub)

	if _, err := GetProviderInfo(testConfig()); err == nil {
		t.Fatal("expected a decode error for malformed JSON")
	}

	if !stub.body.closed {
		t.Error("response body was not closed on the decode-error path")
	}
}

func TestMakeRequestClosesBodyOnEarlyReturns(t *testing.T) {
	htmlHeader := http.Header{}
	htmlHeader.Set("Content-Type", "text/html; charset=utf-8")

	cases := []struct {
		name   string
		status int
		header http.Header
	}{
		// Both return before reading the body, so neither reaches the
		// deferred closes further down MakeRequest.
		{"302 invalid token", http.StatusFound, nil},
		{"HTML unauthenticated", http.StatusOK, htmlHeader},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			stub := &stubTransport{status: tc.status, header: tc.header, payload: "<html></html>"}
			installTransport(t, stub)

			req, err := http.NewRequest(http.MethodGet, "http://example.invalid/api", nil)
			if err != nil {
				t.Fatalf("NewRequest: %v", err)
			}

			if _, err := MakeRequest(req); err == nil {
				t.Fatal("expected an error from this status")
			}

			if !stub.body.closed {
				t.Error("response body was not closed before returning")
			}
		})
	}
}
