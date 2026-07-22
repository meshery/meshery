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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	purgeRetain  int
	purgeExclude string
	purgeDryRun  bool
)

// alwaysExcludedModels are never purged, regardless of --exclude.
var alwaysExcludedModels = []string{"kubernetes", "meshery-core"}

// modelVersionDirs is the purge decision for a single model: which version
// directories are kept and which are slated for removal, both sorted oldest
// to newest.
type modelVersionDirs struct {
	retain []string
	remove []string
}

// purgePlan is the outcome of scanning ./models: what was found, what will
// be skipped, and what will be removed. Building a plan never touches the
// filesystem beyond reading directories.
type purgePlan struct {
	modelsDir       string
	scannedModels   []string
	excludedModels  []string
	noVersionModels []string
	models          map[string]modelVersionDirs
}

func (p *purgePlan) totalToRemove() int {
	n := 0
	for _, m := range p.models {
		n += len(m.remove)
	}
	return n
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Prune old model version directories",
	Long: `Prune older versions of Meshery Models under ./models, retaining only the most recent --retain versions of each model.

Prerequisite: Execute this command from the root of a meshery/meshery repo fork; it operates on the "./models" directory relative to the current working directory.

The "kubernetes" and "meshery-core" models are always excluded from purging, regardless of flags. Use --exclude to skip additional models.
Find more information at: https://docs.meshery.io/reference/mesheryctl/registry/purge`,
	Example: `
// Retain only the latest version of every model (default).
mesheryctl registry purge

// Retain the 3 most recent versions of every model.
mesheryctl registry purge --retain 3

// Preview what would be removed without deleting anything.
mesheryctl registry purge --dry-run

// Additionally skip specific models.
mesheryctl registry purge --exclude aws-ec2-controller,cilium

// Skip the confirmation prompt.
mesheryctl registry purge --retain 2 -y
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if purgeRetain < 1 {
			return ErrPurgeInvalidRetain(purgeRetain)
		}

		cwd, err := os.Getwd()
		if err != nil {
			return ErrPurgeReadModelsDir(err, "models")
		}
		modelsDir := filepath.Join(cwd, "models")

		// Scoped deliberately to modelsDir itself: "no ./models here" is a
		// benign no-op, but a not-exist error surfacing from anywhere deeper in
		// the scan (a model subdirectory removed mid-scan, say) is a real
		// failure and must not be reported as a missing ./models.
		if _, err := os.Stat(modelsDir); errors.Is(err, os.ErrNotExist) {
			utils.Log.Warn(fmt.Errorf("./models directory not found at %s; nothing to purge. Run this command from the root of a meshery/meshery fork", modelsDir))
			return nil
		}

		plan, err := buildPurgePlan(modelsDir, purgeRetain, buildExcludeSet(purgeExclude))
		if err != nil {
			return ErrPurgeReadModelsDir(err, modelsDir)
		}
		utils.Log.Info(fmt.Sprintf("Found ./models directory at %s", modelsDir))

		printPurgePlan(plan)

		if plan.totalToRemove() == 0 {
			utils.Log.Info("Nothing to purge.")
			return nil
		}

		if purgeDryRun {
			utils.Log.Info("Dry run: no directories were removed.")
			return nil
		}

		userResponse := utils.SilentFlag
		if !userResponse {
			userResponse = utils.AskForConfirmation(fmt.Sprintf("This will permanently delete %s under %s. Are you sure you want to continue", pluralize(plan.totalToRemove(), "model version directory", "model version directories"), modelsDir))
		}
		if !userResponse {
			utils.Log.Info("Purge aborted.")
			return nil
		}

		removed, failed := applyPurgePlan(modelsDir, plan)
		utils.Log.Info(fmt.Sprintf("Purge complete: removed %s, %s.", pluralize(removed, "model version directory", "model version directories"), pluralize(len(failed), "failure", "failures")))
		// Surface the per-path structured errors (and their error codes) rather
		// than collapsing them into a generic count.
		return errors.Join(failed...)
	},
}

func init() {
	purgeCmd.Flags().IntVar(&purgeRetain, "retain", 1, "number of most-recent versions to retain per model")
	purgeCmd.Flags().StringVar(&purgeExclude, "exclude", "", "comma-delimited list of additional model names to exclude from purging")
	purgeCmd.Flags().BoolVar(&purgeDryRun, "dry-run", false, "print what would be removed without deleting anything")
	purgeCmd.Flags().BoolVarP(&utils.SilentFlag, "yes", "y", false, "(optional) assume yes for user interactive prompts.")
}

// pluralize renders a count with the noun form that matches it, so
// user-facing messages read naturally for both one and many.
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

// buildExcludeSet merges the always-excluded models with the comma-delimited
// --exclude flag value.
func buildExcludeSet(extra string) map[string]bool {
	set := make(map[string]bool, len(alwaysExcludedModels))
	for _, name := range alwaysExcludedModels {
		set[name] = true
	}
	for _, name := range strings.Split(extra, ",") {
		name = strings.TrimSpace(name)
		if name != "" {
			set[name] = true
		}
	}
	return set
}

// buildPurgePlan scans modelsDir and decides, per model, which version
// directories to retain and which to remove. It performs no mutation.
func buildPurgePlan(modelsDir string, retain int, exclude map[string]bool) (*purgePlan, error) {
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		return nil, err
	}

	plan := &purgePlan{
		modelsDir: modelsDir,
		models:    map[string]modelVersionDirs{},
	}

	for _, entry := range entries {
		// entry.IsDir() is false for symlinks (even ones pointing at a
		// directory) as well as regular files, so this also satisfies
		// "never follow symlinks" and "non-version entries are left alone".
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		plan.scannedModels = append(plan.scannedModels, name)

		if exclude[name] {
			plan.excludedModels = append(plan.excludedModels, name)
			continue
		}

		retainDirs, removeDirs, err := planModelVersions(filepath.Join(modelsDir, name), retain)
		if err != nil {
			return nil, fmt.Errorf("model %s: %w", name, err)
		}

		if len(retainDirs)+len(removeDirs) == 0 {
			plan.noVersionModels = append(plan.noVersionModels, name)
			continue
		}

		plan.models[name] = modelVersionDirs{retain: retainDirs, remove: removeDirs}
	}

	return plan, nil
}

// namedVersion pairs a model version directory's name with its parsed
// semver, so directories can be sorted correctly while keeping the original
// (possibly "v"-prefixed) directory name for filesystem operations.
type namedVersion struct {
	name    string
	version *semver.Version
}

// planModelVersions decides which version directories under modelPath to
// retain (the retain most recent, by real semver ordering) and which to
// remove. Directory entries that are not directories, or whose name does not
// parse as a semver version, are left out of both lists entirely - they are
// never counted toward the retain budget and never removed.
func planModelVersions(modelPath string, retain int) (retainNames, removeNames []string, err error) {
	entries, err := os.ReadDir(modelPath)
	if err != nil {
		return nil, nil, err
	}

	var versions []namedVersion
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		v, err := semver.NewVersion(entry.Name())
		if err != nil {
			continue
		}
		versions = append(versions, namedVersion{name: entry.Name(), version: v})
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[i].version.LessThan(versions[j].version)
	})

	if len(versions) <= retain {
		for _, v := range versions {
			retainNames = append(retainNames, v.name)
		}
		return retainNames, nil, nil
	}

	cut := len(versions) - retain
	for _, v := range versions[:cut] {
		removeNames = append(removeNames, v.name)
	}
	for _, v := range versions[cut:] {
		retainNames = append(retainNames, v.name)
	}
	return retainNames, removeNames, nil
}

// printPurgePlan reports what was found and what will happen, before any
// confirmation prompt or deletion.
func printPurgePlan(plan *purgePlan) {
	utils.Log.Info(fmt.Sprintf("Scanned %s under %s", pluralize(len(plan.scannedModels), "model directory", "model directories"), plan.modelsDir))

	if len(plan.excludedModels) > 0 {
		excluded := append([]string(nil), plan.excludedModels...)
		sort.Strings(excluded)
		utils.Log.Info("Excluded models (skipped entirely): ", strings.Join(excluded, ", "))
	}

	if len(plan.noVersionModels) > 0 {
		noVersion := append([]string(nil), plan.noVersionModels...)
		sort.Strings(noVersion)
		utils.Log.Warn(fmt.Errorf("models with no parseable version directories (left untouched): %s", strings.Join(noVersion, ", ")))
	}

	names := make([]string, 0, len(plan.models))
	for name := range plan.models {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		mv := plan.models[name]
		if len(mv.remove) == 0 {
			utils.Log.Info(fmt.Sprintf("  %s: retaining %s (nothing to remove)", name, strings.Join(mv.retain, ", ")))
			continue
		}
		utils.Log.Info(fmt.Sprintf("  %s: retaining %s; removing %s", name, strings.Join(mv.retain, ", "), strings.Join(mv.remove, ", ")))
	}
}

// applyPurgePlan deletes the version directories marked for removal. It
// re-verifies each target path resolves within modelsDir before deleting it,
// as defense in depth against ever deleting above ./models. Each failure is
// both logged (utils.Log.Error is what renders the MeshKit code, cause and
// remediation) and returned, so the caller can propagate it.
func applyPurgePlan(modelsDir string, plan *purgePlan) (removed int, failed []error) {
	names := make([]string, 0, len(plan.models))
	for name := range plan.models {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		for _, version := range plan.models[name].remove {
			path := filepath.Join(modelsDir, name, version)
			if !isWithinDir(modelsDir, path) {
				unsafe := ErrPurgeUnsafePath(path, modelsDir)
				utils.Log.Error(unsafe)
				failed = append(failed, unsafe)
				continue
			}
			if err := os.RemoveAll(path); err != nil {
				wrapped := ErrPurgeRemove(err, path)
				utils.Log.Error(wrapped)
				failed = append(failed, wrapped)
				continue
			}
			utils.Log.Info("Removed ", filepath.Join(name, version))
			removed++
		}
	}
	return removed, failed
}

// isWithinDir reports whether target is a strict descendant of base.
func isWithinDir(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != "." && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
