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

	"golang.org/x/mod/modfile"
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
	if repoRoot, ok := findMesheryRepoRoot(cwd); ok {
		if rel, err := filepath.Rel(cwd, filepath.Join(repoRoot, "models")); err == nil {
			return rel
		}
	}
	return "models"
}

func isDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// hasMesheryModuleDeclaration parses the go.mod file at path with the Go module
// parser and returns true only if parsing succeeds and the declared module path
// is "github.com/meshery/meshery".
func hasMesheryModuleDeclaration(goModPath string) bool {
	data, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}
	f, err := modfile.Parse(goModPath, data, nil)
	if err != nil || f.Module == nil {
		return false
	}
	return f.Module.Mod.Path == "github.com/meshery/meshery"
}

// isMesheryRepoRoot checks whether dir contains the marker directories and root go.mod of a Meshery repository.
// A valid Meshery repository root contains a "models" directory, at least one of "server" or "mesheryctl",
// and a go.mod file declaring "module github.com/meshery/meshery".
func isMesheryRepoRoot(dir string) bool {
	return isDirectory(filepath.Join(dir, "models")) &&
		(isDirectory(filepath.Join(dir, "server")) || isDirectory(filepath.Join(dir, "mesheryctl"))) &&
		hasMesheryModuleDeclaration(filepath.Join(dir, "go.mod"))
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
// a Meshery repository fork (at the root of the fork or any subdirectory).
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
