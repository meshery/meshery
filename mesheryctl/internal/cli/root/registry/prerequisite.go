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

	gitRoot, err := findMesheryGitRoot(cleanDir)
	if err != nil {
		return ErrPrerequisiteFailed(
			err,
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

// findMesheryGitRoot locates the root of the enclosing git repository. If the nearest
// git root is a nested repo (e.g. submodule), it continues upwards to find an enclosing
// repository that contains Meshery markers.
func findMesheryGitRoot(startDir string) (string, error) {
	curr := filepath.Clean(startDir)
	var firstGitRoot string

	for {
		gitPath := filepath.Join(curr, ".git")
		if isValidGitEntry(gitPath) {
			if firstGitRoot == "" {
				firstGitRoot = curr
			}
			if isMesheryRepo(curr) {
				return curr, nil
			}
		}

		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}

	if firstGitRoot != "" {
		return firstGitRoot, nil
	}
	return "", fmt.Errorf("not inside a git repository: %s", startDir)
}

// isValidGitEntry verifies that gitPath is either a valid .git directory
// or a .git worktree file pointing to an existing directory.
func isValidGitEntry(gitPath string) bool {
	info, err := os.Stat(gitPath)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	_, err = resolveGitFileTarget(gitPath)
	return err == nil
}

// resolveGitFileTarget reads a .git file (worktree/submodule) and verifies that the
// target gitdir exists and is a directory.
func resolveGitFileTarget(gitFilePath string) (string, error) {
	data, err := os.ReadFile(gitFilePath)
	if err != nil {
		return "", err
	}
	content := strings.TrimSpace(string(data))
	if !strings.HasPrefix(content, "gitdir:") {
		return "", fmt.Errorf("invalid gitfile format in %s", gitFilePath)
	}
	target := strings.TrimSpace(strings.TrimPrefix(content, "gitdir:"))
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(gitFilePath), target)
	}
	target = filepath.Clean(target)
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("gitfile target does not exist or is not a directory: %s", target)
	}
	return target, nil
}

// isMesheryRepo checks if the given directory contains characteristic markers
// of a meshery/meshery repository fork, including exact module line and remote URL.
func isMesheryRepo(dir string) bool {
	if !isDirectory(filepath.Join(dir, "mesheryctl")) || !isDirectory(filepath.Join(dir, "server")) {
		return false
	}

	goModPath := filepath.Join(dir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return false
	}

	// Compare parsed module directive exactly
	hasExactModule := false
	for _, line := range strings.Split(string(content), "\n") {
		trimmed := strings.TrimSpace(line)
		fields := strings.Fields(trimmed)
		if len(fields) >= 2 && fields[0] == "module" && fields[1] == "github.com/meshery/meshery" {
			hasExactModule = true
			break
		}
	}
	if !hasExactModule {
		return false
	}

	// Validate git remote points to a meshery/meshery fork/clone
	return hasMesheryRemote(dir)
}

// hasMesheryRemote checks if any configured git remote (e.g. origin, upstream)
// references a meshery/meshery repository or fork.
func hasMesheryRemote(repoRoot string) bool {
	configPaths := getGitConfigFiles(repoRoot)
	for _, cfgPath := range configPaths {
		content, err := os.ReadFile(cfgPath)
		if err != nil {
			continue
		}
		if containsMesheryRemoteURL(string(content)) {
			return true
		}
	}
	return false
}

// getGitConfigFiles discovers candidate git config files for a repo or worktree.
func getGitConfigFiles(repoRoot string) []string {
	var paths []string
	gitEntry := filepath.Join(repoRoot, ".git")
	info, err := os.Stat(gitEntry)
	if err != nil {
		return paths
	}

	var gitDir string
	if info.IsDir() {
		gitDir = gitEntry
	} else {
		target, err := resolveGitFileTarget(gitEntry)
		if err == nil {
			gitDir = target
		}
	}

	if gitDir != "" {
		paths = append(paths, filepath.Join(gitDir, "config"))
		commondirFile := filepath.Join(gitDir, "commondir")
		if cdata, err := os.ReadFile(commondirFile); err == nil {
			cdir := strings.TrimSpace(string(cdata))
			if !filepath.IsAbs(cdir) {
				cdir = filepath.Join(gitDir, cdir)
			}
			paths = append(paths, filepath.Join(cdir, "config"))
		}
	}
	return paths
}

// isMesheryRemoteURL checks if a raw Git remote URL identifies a meshery repository clone or fork,
// comparing the complete repository identity and rejecting names like meshery-other.
func isMesheryRemoteURL(rawURL string) bool {
	u := strings.TrimSpace(rawURL)
	u = strings.Trim(u, "\"")
	u = strings.TrimSuffix(u, "/")
	u = strings.TrimSuffix(u, ".git")
	idx := strings.LastIndexAny(u, "/:")
	if idx == -1 {
		return false
	}
	repoName := u[idx+1:]
	return strings.EqualFold(repoName, "meshery")
}

// cleanGitConfigValue decodes Git configuration value quotes and strips trailing inline comments.
func cleanGitConfigValue(val string) string {
	trimmed := strings.TrimSpace(val)
	if strings.HasPrefix(trimmed, "\"") {
		idx := -1
		escaped := false
		for i := 1; i < len(trimmed); i++ {
			if escaped {
				escaped = false
				continue
			}
			if trimmed[i] == '\\' {
				escaped = true
				continue
			}
			if trimmed[i] == '"' {
				idx = i
				break
			}
		}
		if idx != -1 {
			trimmed = trimmed[1:idx]
		} else {
			trimmed = strings.Trim(trimmed, "\"")
		}
	} else {
		if idx := strings.IndexAny(trimmed, "#;"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
	}
	return strings.TrimSpace(trimmed)
}

// containsMesheryRemoteURL checks if git config text defines a remote URL for meshery.
// Section and key names are recognized case-insensitively per Git configuration spec,
// while preserving subsection names and URL values.
func containsMesheryRemoteURL(configText string) bool {
	lines := strings.Split(configText, "\n")
	inRemoteSection := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, ";") {
			continue
		}
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			sectionContent := strings.TrimSpace(trimmed[1 : len(trimmed)-1])
			fields := strings.Fields(sectionContent)
			if (len(fields) >= 2 && strings.EqualFold(fields[0], "remote")) ||
				strings.HasPrefix(strings.ToLower(sectionContent), "remote.") {
				inRemoteSection = true
			} else {
				inRemoteSection = false
			}
			continue
		}
		if inRemoteSection {
			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				key := strings.ToLower(strings.TrimSpace(parts[0]))
				if key == "url" {
					urlVal := cleanGitConfigValue(parts[1])
					if isMesheryRemoteURL(urlVal) {
						return true
					}
				}
			}
		}
	}
	return false
}
