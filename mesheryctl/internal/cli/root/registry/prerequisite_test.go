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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupMockMesheryRepo creates a temporary directory structured as a meshery/meshery repo.
func setupMockMesheryRepo(t *testing.T, isWorktree bool) string {
	t.Helper()
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	require.NoError(t, err)

	// Create .git entry (directory or file)
	gitPath := filepath.Join(root, ".git")
	if isWorktree {
		err = os.WriteFile(gitPath, []byte("gitdir: /path/to/worktrees/mock\n"), 0o644)
		require.NoError(t, err)
	} else {
		err = os.MkdirAll(gitPath, 0o755)
		require.NoError(t, err)
	}

	// Create mesheryctl and server directories
	require.NoError(t, os.MkdirAll(filepath.Join(root, "mesheryctl"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "server"), 0o755))

	// Create go.mod
	goModContent := "module github.com/meshery/meshery\n\ngo 1.23\n"
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte(goModContent), 0o644))

	return root
}

func TestVerifyRepoRoot_ValidRoot(t *testing.T) {
	root := setupMockMesheryRepo(t, false)
	err := verifyRepoRootDir(root)
	assert.NoError(t, err)
}

func TestVerifyRepoRoot_WorktreeFile(t *testing.T) {
	root := setupMockMesheryRepo(t, true)
	err := verifyRepoRootDir(root)
	assert.NoError(t, err)
}

func TestVerifyRepoRoot_Subdirectory(t *testing.T) {
	root := setupMockMesheryRepo(t, false)
	subDir := filepath.Join(root, "mesheryctl")

	err := verifyRepoRootDir(subDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ErrPrerequisiteFailedCode)
	assert.Contains(t, err.Error(), "command executed from subdirectory")
	assert.Contains(t, err.Error(), "mesheryctl")
}

func TestVerifyRepoRoot_NestedSubdirectory(t *testing.T) {
	root := setupMockMesheryRepo(t, false)
	nestedDir := filepath.Join(root, "mesheryctl", "internal", "cli")
	require.NoError(t, os.MkdirAll(nestedDir, 0o755))

	err := verifyRepoRootDir(nestedDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ErrPrerequisiteFailedCode)
	assert.Contains(t, err.Error(), "command executed from subdirectory")
}

func TestVerifyRepoRoot_NonGitDirectory(t *testing.T) {
	emptyDir := t.TempDir()
	emptyDir, err := filepath.EvalSymlinks(emptyDir)
	require.NoError(t, err)

	err = verifyRepoRootDir(emptyDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ErrPrerequisiteFailedCode)
	assert.Contains(t, err.Error(), "not inside a git repository")
}

func TestVerifyRepoRoot_ForeignGitRepository(t *testing.T) {
	foreignRoot := t.TempDir()
	foreignRoot, err := filepath.EvalSymlinks(foreignRoot)
	require.NoError(t, err)

	// Create .git but without meshery markers
	require.NoError(t, os.MkdirAll(filepath.Join(foreignRoot, ".git"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(foreignRoot, "go.mod"), []byte("module github.com/other/project\n"), 0o644))

	err = verifyRepoRootDir(foreignRoot)
	require.Error(t, err)
	assert.Contains(t, err.Error(), ErrPrerequisiteFailedCode)
	assert.Contains(t, err.Error(), "not a meshery/meshery repository")
}

func TestVerifyRepoRoot_CurrentWorkingDir(t *testing.T) {
	root := setupMockMesheryRepo(t, false)
	chdir(t, root)

	err := verifyRepoRoot()
	assert.NoError(t, err)

	// Change to subdirectory
	subDir := filepath.Join(root, "mesheryctl")
	chdir(t, subDir)

	err = verifyRepoRoot()
	require.Error(t, err)
	assert.Contains(t, err.Error(), ErrPrerequisiteFailedCode)
}
