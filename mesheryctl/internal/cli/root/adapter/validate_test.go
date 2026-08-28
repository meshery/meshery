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
	"testing"

	"github.com/meshery/meshery/server/models"
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
