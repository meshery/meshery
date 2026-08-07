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

package registry

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestDeleteCmdArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no args provided",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid UUID format",
			args:    []string{"invalid-uuid-string"},
			wantErr: true,
		},
		{
			name:    "valid connection UUID",
			args:    []string{"50bef83c-dad7-9977-952c-099321286a6a"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deleteCmd.Args(deleteCmd, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteCmd.Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteCmdRunE_MockServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message":"Deleted 9 models for registrant \"meshery\"","count":9,"connectionName":"meshery"}`))
	}))
	defer ts.Close()

	utils.MesheryEndpoint = ts.URL

	err := deleteCmd.RunE(deleteCmd, []string{"50bef83c-dad7-9977-952c-099321286a6a"})
	if err != nil {
		t.Errorf("deleteCmd.RunE() error = %v, want nil", err)
	}
}
