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
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// verifyRepoRoot checks that the current working directory is the root of a
// meshery/meshery repository fork.
func verifyRepoRoot() error {
	cwd, err := os.Getwd()
	if err != nil {
		return ErrPrerequisiteFailed(err, "Ensure current working directory is accessible and retry")
	}
	return verifyRepoRootDir(cwd)
}

// verifyRepoRootDir validates that dir is the root of a meshery/meshery repository fork.
func verifyRepoRootDir(dir string) error {
	cleanDir, err := filepath.EvalSymlinks(filepath.Clean(dir))
	if err != nil {
		cleanDir = filepath.Clean(dir)
	}

	gitRoot, err := findGitRoot(cleanDir)
	if err != nil {
		return ErrPrerequisiteFailed(
			fmt.Errorf("not inside a git repository: %s", dir),
			"Execute this command from the root of a meshery/meshery repository fork",
		)
	}

	cleanGitRoot, err := filepath.EvalSymlinks(filepath.Clean(gitRoot))
	if err != nil {
		cleanGitRoot = filepath.Clean(gitRoot)
	}

	if !strings.EqualFold(cleanDir, cleanGitRoot) {
		relPath, relErr := filepath.Rel(cleanGitRoot, cleanDir)
		if relErr != nil {
			relPath = cleanDir
		}
		return ErrPrerequisiteFailed(
			fmt.Errorf("command executed from subdirectory '%s'", relPath),
			fmt.Sprintf("Navigate to the repository root directory at '%s' before executing this command", cleanGitRoot),
		)
	}

	if !isMesheryRepo(cleanGitRoot) {
		return ErrPrerequisiteFailed(
			fmt.Errorf("current git repository at '%s' is not a meshery/meshery repository", cleanGitRoot),
			"Ensure this command is executed from a clone or fork of github.com/meshery/meshery",
		)
	}

	return nil
}

// findGitRoot traverses upwards starting from startDir to locate the nearest directory
// containing a .git entry (file or directory).
func findGitRoot(startDir string) (string, error) {
	curr := filepath.Clean(startDir)
	for {
		gitPath := filepath.Join(curr, ".git")
		if _, err := os.Stat(gitPath); err == nil {
			return curr, nil
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return "", fmt.Errorf("no .git repository found in %s or any parent directory", startDir)
}

// isMesheryRepo checks if the given directory contains characteristic markers
// of a meshery/meshery repository fork.
func isMesheryRepo(dir string) bool {
	if !isDirectory(filepath.Join(dir, "mesheryctl")) || !isDirectory(filepath.Join(dir, "server")) {
		return false
	}

	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "module github.com/meshery/meshery")
}
