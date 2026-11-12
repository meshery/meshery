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

package design

import (
	"io"
	"os"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiplePatternsConfirmationShowsEachDesign verifies that when several
// designs share a name, `design deploy` lists every design's own file so the
// user can tell them apart before picking an index.
func TestMultiplePatternsConfirmationShowsEachDesign(t *testing.T) {
	stdinR, stdinW, err := os.Pipe()
	require.NoError(t, err)
	stdoutR, stdoutW, err := os.Pipe()
	require.NoError(t, err)

	origStdin, origStdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = stdinR, stdoutW
	t.Cleanup(func() { os.Stdin, os.Stdout = origStdin, origStdout })

	_, err = stdinW.WriteString("1\n")
	require.NoError(t, err)
	require.NoError(t, stdinW.Close())

	first, second := uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4())
	designs := []models.MesheryPattern{
		{ID: &first, Name: "nginx", PatternFile: "name: nginx\nversion: 0.0.1\n"},
		{ID: &second, Name: "nginx", PatternFile: "name: nginx\nversion: 0.0.2\n"},
	}

	index := multiplepatternsConfirmation(designs)

	require.NoError(t, stdoutW.Close())
	out, err := io.ReadAll(stdoutR)
	require.NoError(t, err)

	assert.Equal(t, 1, index)
	assert.Contains(t, string(out), "version: 0.0.1")
	assert.Contains(t, string(out), "version: 0.0.2")
}
