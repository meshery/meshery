// # Copyright Meshery Authors
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestPublishCmdPreRunE_ArgValidation verifies that the publish command's
// argument validation accepts exactly 4 or 5 positional arguments and rejects
// everything else. This prevents the regression where the official help-text
// example (4 args for the "meshery" system) was rejected by a len(args)!=5
// guard, and ensures that remote-provider / website callers still supply the
// required imgs-output-path.
func TestPublishCmdPreRunE_ArgValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			// Matches the official "Publish To Meshery" example in --help (4 args).
			name:    "4 args accepted (meshery system without imgs-output-path)",
			args:    []string{"meshery", "GoogleCredential", "GoogleSheetID", "path/to/models"},
			wantErr: false,
		},
		{
			// Full 5-arg form used by remote-provider and website targets.
			name:    "5 args accepted (all targets with imgs-output-path)",
			args:    []string{"remote-provider", "GoogleCredential", "GoogleSheetID", "path/to/models", "path/to/imgs"},
			wantErr: false,
		},
		{
			// Too few arguments – missing at least models-output-path.
			name:    "3 args rejected",
			args:    []string{"meshery", "GoogleCredential", "GoogleSheetID"},
			wantErr: true,
		},
		{
			// Too many arguments – no valid target needs more than 5 positional args.
			name:    "6 args rejected",
			args:    []string{"meshery", "GoogleCredential", "GoogleSheetID", "path/to/models", "path/to/imgs", "extra"},
			wantErr: true,
		},
		{
			// No arguments at all.
			name:    "0 args rejected",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := publishCmd.PreRunE(publishCmd, tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
