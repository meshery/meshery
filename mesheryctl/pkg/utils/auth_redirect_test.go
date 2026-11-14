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
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshkit/errors"
)

// itemHits counts the requests each handler of newItemServer received.
type itemHits struct {
	get, delete atomic.Int32
}

// newItemServer serves /api/items/{id} for GET and DELETE through a real
// gorilla/mux router, plus routes that redirect to it with a given status.
// Its GET route answers 200 with JSON, the way a list or show endpoint would,
// so a DELETE or POST that gets re-sent as a GET there looks like a success.
func newItemServer(t *testing.T) (*httptest.Server, *itemHits) {
	t.Helper()
	hits := &itemHits{}

	r := mux.NewRouter()
	r.HandleFunc("/api/items/{id}", func(w http.ResponseWriter, _ *http.Request) {
		hits.get.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"id":"1"}`)
	}).Methods(http.MethodGet)
	r.HandleFunc("/api/items/{id}", func(w http.ResponseWriter, _ *http.Request) {
		hits.delete.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodDelete)
	for _, status := range []int{http.StatusFound, http.StatusSeeOther, http.StatusTemporaryRedirect} {
		r.HandleFunc(fmt.Sprintf("/redirect/%d", status), func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, "/api/items/1", status)
		})
	}
	r.HandleFunc("/loop", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/loop", http.StatusTemporaryRedirect)
	})

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)
	return server, hits
}

func TestMakeRequestRefusesMethodChangingRedirect(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name     string
		method   string
		path     string
		body     string
		wantCode string
		// wantInMessage is checked only for ErrRedirectChangesMethod.
		wantInMessage []string
	}{
		{
			// gorilla/mux answers a path it has to clean with 301 to the cleaned path.
			name:          "DELETE answered with 301 by gorilla/mux path cleaning",
			method:        http.MethodDelete,
			path:          "//api/items/1",
			wantCode:      ErrRedirectChangesMethodCode,
			wantInMessage: []string{"DELETE", "//api/items/1", "301 Moved Permanently", "/api/items/1"},
		},
		{
			name:          "POST answered with 303 See Other",
			method:        http.MethodPost,
			path:          "/redirect/303",
			body:          `{"name":"item"}`,
			wantCode:      ErrRedirectChangesMethodCode,
			wantInMessage: []string{"POST", "/redirect/303", "303 See Other", "/api/items/1"},
		},
		{
			// Meshery Server's provider and session middleware answer a request it
			// cannot authenticate with 302 to /provider; MakeRequest maps a 302 to
			// ErrInvalidToken, which a DELETE now reaches instead of following it.
			name:     "DELETE answered with 302 Found reports an invalid token",
			method:   http.MethodDelete,
			path:     "/redirect/302",
			wantCode: ErrInvalidTokenCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, hits := newItemServer(t)

			transport := &trackingTransport{base: http.DefaultTransport}
			previousTransport := http.DefaultTransport
			http.DefaultTransport = transport
			t.Cleanup(func() { http.DefaultTransport = previousTransport })

			req, err := http.NewRequest(tt.method, server.URL+tt.path, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			resp, err := MakeRequest(req)
			if resp != nil {
				SafeClose(resp.Body)
				t.Errorf("got response %d, want nil", resp.StatusCode)
			}
			if got := hits.get.Load(); got != 0 {
				t.Errorf("the %s was re-sent as a GET and ran the GET handler %d time(s)", tt.method, got)
			}
			if got := hits.delete.Load(); got != 0 {
				t.Errorf("DELETE handler ran %d time(s), want 0", got)
			}
			if got := errors.GetCode(err); got != tt.wantCode {
				t.Fatalf("error code = %q (err: %v), want %q", got, err, tt.wantCode)
			}
			for _, want := range tt.wantInMessage {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q does not mention %q", err.Error(), want)
				}
			}
			if transport.body == nil || !transport.body.closed {
				t.Error("redirect response body was not closed")
			}
		})
	}
}

func TestMakeRequestFollowsRedirectsThatKeepTheMethod(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantGet    int32
		wantDelete int32
	}{
		{
			name:       "GET answered with 301 by gorilla/mux path cleaning",
			method:     http.MethodGet,
			path:       "//api/items/1",
			wantStatus: http.StatusOK,
			wantGet:    1,
		},
		{
			name:       "DELETE answered with 307 Temporary Redirect",
			method:     http.MethodDelete,
			path:       "/redirect/307",
			wantStatus: http.StatusNoContent,
			wantDelete: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, hits := newItemServer(t)

			req, err := http.NewRequest(tt.method, server.URL+tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := MakeRequest(req)
			if err != nil {
				t.Fatalf("MakeRequest: %v", err)
			}
			SafeClose(resp.Body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.wantStatus)
			}
			if got := hits.get.Load(); got != tt.wantGet {
				t.Errorf("GET handler ran %d time(s), want %d", got, tt.wantGet)
			}
			if got := hits.delete.Load(); got != tt.wantDelete {
				t.Errorf("DELETE handler ran %d time(s), want %d", got, tt.wantDelete)
			}
		})
	}
}

// MakeRequest replaces net/http's default redirect policy, so it must keep
// that policy's limit rather than follow a redirect loop forever.
func TestMakeRequestStopsRedirectLoops(t *testing.T) {
	_ = SetupMeshkitLoggerTesting(t, false)
	server, _ := newItemServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/loop", nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := MakeRequest(req)
	if resp != nil {
		SafeClose(resp.Body)
		t.Errorf("got response %d, want nil", resp.StatusCode)
	}
	if got := errors.GetCode(err); got != ErrFailRequestCode {
		t.Fatalf("error code = %q (err: %v), want %q", got, err, ErrFailRequestCode)
	}
	if want := fmt.Sprintf("stopped after %d redirects", maxRedirects); !strings.Contains(err.Error(), want) {
		t.Errorf("error %q does not mention %q", err.Error(), want)
	}
}
