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
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestInvokeCompUpdateReturnsComponentsSheetError(t *testing.T) {
	utils.SetupMeshkitLoggerTesting(t, false)

	origComponentCSVFilePath, origLogFile := componentCSVFilePath, logFile
	t.Cleanup(func() {
		componentCSVFilePath, logFile = origComponentCSVFilePath, origLogFile
	})

	var err error
	logFile, err = os.Create(filepath.Join(t.TempDir(), "registry-update"))
	require.NoError(t, err)

	// A components CSV that cannot be read stands in for a Components
	// sheet that fails to download or parse.
	componentCSVFilePath = filepath.Join(t.TempDir(), "missing-components.csv")

	err = InvokeCompUpdate()
	require.Error(t, err)
	require.Contains(t, err.Error(), "error updating registry")
}
