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

// TestFindAdapter covers resolving a registered mesh adapter by the host
// portion of its Location, including the no-port case from
// https://github.com/meshery/meshery/issues/21630 where matching used to
// leak adapter.Location into the caller's meshName instead of adapter.Name.
func TestFindAdapter(t *testing.T) {
	adapters := []*models.Adapter{
		{Name: "ISTIO", Location: "meshery-istio:10000"},
		{Name: "NSM", Location: "meshery-nsm"}, // registered without a port
	}

	tests := []struct {
		name       string
		adapterURL string
		wantFound  bool
		wantName   string
		wantLoc    string
	}{
		{
			name:       "adapter registered with a port",
			adapterURL: "meshery-istio",
			wantFound:  true,
			wantName:   "ISTIO",
			wantLoc:    "meshery-istio:10000",
		},
		{
			name:       "adapter registered without a port",
			adapterURL: "meshery-nsm",
			wantFound:  true,
			wantName:   "NSM",
			wantLoc:    "meshery-nsm",
		},
		{
			name:       "no matching adapter",
			adapterURL: "meshery-linkerd",
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := findAdapter(adapters, tt.adapterURL)

			if ok != tt.wantFound {
				t.Fatalf("findAdapter() found = %v, want %v", ok, tt.wantFound)
			}
			if !tt.wantFound {
				if got != nil {
					t.Fatalf("findAdapter() adapter = %#v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("findAdapter() returned nil adapter")
			}
			if got.Name != tt.wantName {
				t.Errorf("adapter Name = %q, want %q", got.Name, tt.wantName)
			}
			if got.Location != tt.wantLoc {
				t.Errorf("adapter Location = %q, want %q", got.Location, tt.wantLoc)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Reset()
			httpmock.RegisterResponder(http.MethodGet, testContext.BaseURL+"/api/system/sync",
				httpmock.NewStringResponder(200, utils.NewGoldenFile(t, "sync.adapters.golden", fixturesDir).Load()))

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
