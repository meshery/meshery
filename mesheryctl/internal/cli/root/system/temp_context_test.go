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

package system

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogoutHonorsContextFlag verifies that `system logout -c <context>` clears
// the token of the named context and leaves the current context's token alone.
func TestLogoutHonorsContextFlag(t *testing.T) {
	setupContextTestEnv(t)
	utils.SetupMeshkitLoggerTesting(t, false)
	t.Cleanup(BreakupFunc)

	content := []byte(`{"auths":{"meshery":{"token":"secret"}}}`)
	dir := t.TempDir()
	localToken := filepath.Join(dir, "local.json")
	gkeToken := filepath.Join(dir, "gke.json")
	require.NoError(t, os.WriteFile(localToken, content, 0o600))
	require.NoError(t, os.WriteFile(gkeToken, content, 0o600))
	viper.Set("contexts.gke.token", "gke")
	viper.Set("tokens", []map[string]string{
		{"name": "default", "location": localToken},
		{"name": "gke", "location": gkeToken},
	})

	SystemCmd.SetArgs([]string{"logout", "-c", "gke"})
	require.NoError(t, SystemCmd.Execute())

	got, err := os.ReadFile(gkeToken)
	require.NoError(t, err)
	assert.Empty(t, got, "token of the context passed with -c should be cleared")

	got, err = os.ReadFile(localToken)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(content, got), "token of the current context should be left alone")
}

// TestContextFlagRejectsUnknownContext verifies that dashboard, login and
// logout use the context passed with -c, so an unknown context is reported
// instead of the command silently running against the current context.
func TestContextFlagRejectsUnknownContext(t *testing.T) {
	for _, command := range []string{"dashboard", "login", "logout"} {
		t.Run(command, func(t *testing.T) {
			setupContextTestEnv(t)
			utils.SetupMeshkitLoggerTesting(t, false)
			t.Cleanup(BreakupFunc)

			tokenPath := filepath.Join(t.TempDir(), "token.json")
			content := []byte(`{"auths":{"meshery":{"token":"secret"}}}`)
			require.NoError(t, os.WriteFile(tokenPath, content, 0o600))
			viper.Set("tokens", []map[string]string{{"name": "default", "location": tokenPath}})

			SystemCmd.SetArgs([]string{command, "-c", "no-such-context"})
			err := SystemCmd.Execute()

			require.Error(t, err)
			assert.Equal(t, ErrGetCurrentContextCode, errors.GetCode(err))

			got, readErr := os.ReadFile(tokenPath)
			require.NoError(t, readErr)
			assert.True(t, bytes.Equal(content, got), "the current context's token should be left alone")
		})
	}
}
