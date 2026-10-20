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
)

// defaultModelsLocation resolves the default relative path to the repo-root
// "models" directory based on the current working directory.
//
// The registry commands document being executed from the root of a
// meshery/meshery fork, but they are also run from subdirectories (e.g.
// mesheryctl/ in CI). A single hardcoded relative default cannot satisfy both:
// "./models" resolves outside the repo when run from a subdirectory, while
// "../models" resolves outside the repo when run from the root. Probe for the
// "models" directory and return a relative path that resolves to it from the
// current directory.
//
// This is only used when the user has not explicitly supplied --output/--input;
// an explicit value is always honored as-is.
func defaultModelsLocation() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "models"
	}
	// Run from the repository root (the documented prerequisite).
	if isDirectory(filepath.Join(cwd, "models")) {
		return "models"
	}
	// Run from an immediate subdirectory such as mesheryctl/.
	if isDirectory(filepath.Join(cwd, "..", "models")) {
		return filepath.Join("..", "models")
	}
	// Fall back to the repo-root assumption; generate will create the
	// directory if it does not yet exist.
	return "models"
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// isMesheryRepoRoot checks whether dir contains the marker directories of the Meshery repository root.
// A valid Meshery repository root contains a "models" directory and at least one of "server" or "mesheryctl".
func isMesheryRepoRoot(dir string) bool {
	return isDirectory(filepath.Join(dir, "models")) &&
		(isDirectory(filepath.Join(dir, "server")) || isDirectory(filepath.Join(dir, "mesheryctl")))
}

// findMesheryRepoRoot walks up the directory tree from startDir looking for the
// root of a Meshery repository fork. Returns the root directory path and true if found,
// or an empty string and false if not found.
func findMesheryRepoRoot(startDir string) (string, bool) {
	dir := filepath.Clean(startDir)
	for {
		if isMesheryRepoRoot(dir) {
			return dir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", false
}

// verifyPrerequisiteRepo checks whether the command is being executed from within
// a Meshery repository fork (at the root of the fork or an immediate subdirectory).
func verifyPrerequisiteRepo() error {
	cwd, err := os.Getwd()
	if err != nil {
		return ErrPrereqRepository(err)
	}

	if _, ok := findMesheryRepoRoot(cwd); ok {
		return nil
	}

	return ErrPrereqRepository(fmt.Errorf("current working directory %q is not within a meshery/meshery repository fork", cwd))
}
