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

package adapter

import (
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
)

// TestFindAdapters covers resolving a registered mesh adapter from --adapter by
// exact location or by host, including the no-port case from
// https://github.com/meshery/meshery/issues/21630 and the host collision, exact
// host:port and IPv6 cases raised in review on
// https://github.com/meshery/meshery/pull/21702.
func TestFindAdapters(t *testing.T) {
	adapters := []*models.Adapter{
		{Name: "ISTIO", Location: "meshery-istio:10000"},
		{Name: "NSM", Location: "meshery-nsm"},
		{Name: "LOCAL_ISTIO", Location: "localhost:10000"},
		{Name: "LOCAL_LINKERD", Location: "localhost:10001"},
		{Name: "IPV6_CONSUL", Location: "[::1]:10002"},
	}

	tests := []struct {
		name       string
		adapterURL string
		want       []string
	}{
		{
			name:       "host of an adapter registered with a port",
			adapterURL: "meshery-istio",
			want:       []string{"ISTIO"},
		},
		{
			name:       "adapter registered without a port",
			adapterURL: "meshery-nsm",
			want:       []string{"NSM"},
		},
		{
			name:       "exact host:port location",
			adapterURL: "meshery-istio:10000",
			want:       []string{"ISTIO"},
		},
		{
			name:       "exact location disambiguates a shared host",
			adapterURL: "localhost:10001",
			want:       []string{"LOCAL_LINKERD"},
		},
		{
			name:       "host shared by several adapters returns all of them",
			adapterURL: "localhost",
			want:       []string{"LOCAL_ISTIO", "LOCAL_LINKERD"},
		},
		{
			name:       "exact IPv6 location",
			adapterURL: "[::1]:10002",
			want:       []string{"IPV6_CONSUL"},
		},
		{
			name:       "IPv6 host without a port",
			adapterURL: "::1",
			want:       []string{"IPV6_CONSUL"},
		},
		{
			name:       "bracketed IPv6 host without a port",
			adapterURL: "[::1]",
			want:       []string{"IPV6_CONSUL"},
		},
		{
			name:       "known host with an unregistered port",
			adapterURL: "localhost:19999",
			want:       nil,
		},
		{
			name:       "no matching adapter",
			adapterURL: "meshery-linkerd",
			want:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findAdapters(adapters, tt.adapterURL)

			names := make([]string, 0, len(got))
			for _, adapter := range got {
				names = append(names, adapter.Name)
			}

			if len(names) != len(tt.want) {
				t.Fatalf("findAdapters() = %v, want %v", names, tt.want)
			}
			for i := range names {
				if names[i] != tt.want[i] {
					t.Fatalf("findAdapters() = %v, want %v", names, tt.want)
				}
			}
		})
	}
}

// TestValidateAdapterSelection pins that `--adapter` selects the adapter the
// operation is sent to, per https://github.com/meshery/meshery/issues/21679
// where an unrecognized name was silently replaced by whichever adapter
// happened to be connected and the command still exited 0.
func TestValidateAdapterSelection(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	t.Cleanup(func() { utils.StopMockery(t) })

	testContext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	fixturesDir := filepath.Join(filepath.Dir(filename), "fixtures")

	tests := []struct {
		name string
		args []string
		// fixture is the session data the server responds with, defaulting to
		// adapters on distinct hosts
		fixture string
		// wantErrCode, when set, is the meshkit code the command must fail with
		wantErrCode string
		// wantMesh and wantAdapter are what a successful run must resolve to
		wantMesh    string
		wantAdapter string
	}{
		{
			name:        "unrecognized adapter is rejected",
			args:        []string{"validate", "--adapter", "meshery-doesnotexist", "--spec", "smi"},
			wantErrCode: ErrAdapterNotFoundCode,
		},
		{
			name:        "mesh name contradicting the adapter is rejected",
			args:        []string{"validate", "linkerd", "--adapter", "meshery-istio", "--spec", "smi"},
			wantErrCode: ErrAdapterMeshMismatchCode,
		},
		{
			name:        "connected adapter resolves to its own mesh",
			args:        []string{"validate", "--adapter", "meshery-istio", "--spec", "smi"},
			wantMesh:    "ISTIO",
			wantAdapter: "meshery-istio:10000",
		},
		{
			name:        "mesh name agreeing with the adapter is accepted",
			args:        []string{"validate", "istio", "--adapter", "meshery-istio", "--spec", "smi"},
			wantMesh:    "ISTIO",
			wantAdapter: "meshery-istio:10000",
		},
		{
			name:        "exact host:port location is accepted",
			args:        []string{"validate", "--adapter", "meshery-istio:10000", "--spec", "smi"},
			wantMesh:    "ISTIO",
			wantAdapter: "meshery-istio:10000",
		},
		{
			name:        "host shared by several adapters is rejected",
			args:        []string{"validate", "--adapter", "localhost", "--spec", "smi"},
			fixture:     "sync.localhost.golden",
			wantErrCode: ErrAmbiguousAdapterCode,
		},
		{
			name:        "exact location selects one of several adapters on a host",
			args:        []string{"validate", "--adapter", "localhost:10001", "--spec", "smi"},
			fixture:     "sync.localhost.golden",
			wantMesh:    "LINKERD",
			wantAdapter: "localhost:10001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixture := tt.fixture
			if fixture == "" {
				fixture = "sync.adapters.golden"
			}

			httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/system/sync",
				httpmock.NewStringResponder(200, utils.NewGoldenFile(t, fixture, fixturesDir).Load()))

			// capture what actually goes over the wire as the adapter to operate on
			var sentAdapter string
			httpmock.RegisterResponder(http.MethodPost, testContext.BaseURL+"/api/system/adapter/operation",
				func(req *http.Request) (*http.Response, error) {
					if err := req.ParseForm(); err != nil {
						return nil, err
					}
					sentAdapter = req.PostFormValue("adapter")
					return httpmock.NewStringResponse(200, ""), nil
				})

			// AdapterCmd and the vars its flags bind to are package-level
			// singletons; meshName is not flag-bound, so it would otherwise
			// survive from one invocation to the next
			meshName = ""
			adapterURL = ""
			watch = false
			utils.TokenFlag = filepath.Join(fixturesDir, "token.golden")
			buff := utils.SetupMeshkitLoggerTesting(t, false)

			AdapterCmd.SetArgs(tt.args)
			AdapterCmd.SetOut(buff)
			err := AdapterCmd.Execute()

			if tt.wantErrCode != "" {
				if err == nil {
					t.Fatalf("expected %s, got no error (adapter sent: %q)", tt.wantErrCode, sentAdapter)
				}
				if got := errors.GetCode(err); got != tt.wantErrCode {
					t.Fatalf("error code = %q, want %q (err: %v)", got, tt.wantErrCode, err)
				}
				if sentAdapter != "" {
					t.Errorf("rejected run still sent an operation for adapter %q", sentAdapter)
				}
				return
			}

			if err != nil {
				t.Fatal(err)
			}
			if meshName != tt.wantMesh {
				t.Errorf("meshName = %q, want %q", meshName, tt.wantMesh)
			}
			if sentAdapter != tt.wantAdapter {
				t.Errorf("adapter sent = %q, want %q", sentAdapter, tt.wantAdapter)
			}
		})
	}
}
