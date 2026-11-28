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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

// modelVersions records the version directories present before generation so
// the sync can distinguish new versions from existing versions it refreshed.
type modelVersions map[string]map[string]struct{}

func snapshotModelVersions(modelsPath string) (modelVersions, error) {
	models := make(modelVersions)
	entries, err := os.ReadDir(modelsPath)
	if errors.Is(err, os.ErrNotExist) {
		return models, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read models directory %q: %w", modelsPath, err)
	}

	for _, modelEntry := range entries {
		if !modelEntry.IsDir() {
			continue
		}
		modelPath := filepath.Join(modelsPath, modelEntry.Name())
		versionEntries, err := os.ReadDir(modelPath)
		if err != nil {
			return nil, fmt.Errorf("read model directory %q: %w", modelPath, err)
		}

		for _, versionEntry := range versionEntries {
			if !versionEntry.IsDir() {
				continue
			}
			if _, err := semver.NewVersion(versionEntry.Name()); err != nil {
				continue
			}
			if models[modelEntry.Name()] == nil {
				models[modelEntry.Name()] = make(map[string]struct{})
			}
			models[modelEntry.Name()][versionEntry.Name()] = struct{}{}
		}
	}

	return models, nil
}

// carryForwardRelationships copies missing relationships from the closest
// older model version into each version directory created by this generation
// run. Explicit relationships generated for the new version take precedence,
// and inherited relationships whose component kinds were removed are skipped.
func carryForwardRelationships(modelsPath string, before modelVersions) error {
	after, err := snapshotModelVersions(modelsPath)
	if err != nil {
		return err
	}

	var carryErrors []error
	modelNames := make([]string, 0, len(after))
	for modelName := range after {
		modelNames = append(modelNames, modelName)
	}
	sort.Strings(modelNames)

	for _, modelName := range modelNames {
		versions := after[modelName]
		newVersions := make([]string, 0)
		for version := range versions {
			if _, existed := before[modelName][version]; !existed {
				newVersions = append(newVersions, version)
			}
		}
		sort.Slice(newVersions, func(i, j int) bool {
			left, _ := semver.NewVersion(newVersions[i])
			right, _ := semver.NewVersion(newVersions[j])
			return left.LessThan(right)
		})

		modelPath := filepath.Join(modelsPath, modelName)
		// A version that failed or is incomplete is not a source for newer
		// versions; they fall back to the next older one and are reported as
		// incomplete too, so removing them all and rerunning retries the chain.
		failed := make(map[string]struct{})
		for _, newVersion := range newVersions {
			sources := make(map[string]struct{}, len(versions))
			for version := range versions {
				if _, isFailed := failed[version]; !isFailed {
					sources[version] = struct{}{}
				}
			}
			previousVersion := latestPreviousVersion(newVersion, sources)
			skipped := failedVersionsBetween(failed, previousVersion, newVersion)

			var versionErrors []error
			if previousVersion != "" {
				if err := carryVersionRelationships(modelPath, modelName, previousVersion, newVersion); err != nil {
					versionErrors = append(versionErrors, err)
				}
			}
			if len(skipped) > 0 {
				versionErrors = append(versionErrors, fmt.Errorf("%s %s did not inherit from %s, which failed", modelName, newVersion, strings.Join(skipped, ", ")))
			}
			if len(versionErrors) > 0 {
				failed[newVersion] = struct{}{}
				carryErrors = append(carryErrors, fmt.Errorf("carry relationships to %s %s (remove %q and rerun to retry): %w", modelName, newVersion, filepath.Join(modelPath, newVersion), errors.Join(versionErrors...)))
			}
		}
	}

	return errors.Join(carryErrors...)
}

// modelComponentKinds returns the kinds declared by the component files in
// directory. Filenames don't always match the kind ("digital ocean.json"
// declares "digital-ocean"), so the filename is only used when a file
// declares no kind.
// carryVersionRelationships copies the relationships of previousVersion that
// newVersion does not already have and whose component kinds it still has.
func carryVersionRelationships(modelPath, modelName, previousVersion, newVersion string) error {
	source := filepath.Join(modelPath, previousVersion, defVersion, "relationships")
	destination := filepath.Join(modelPath, newVersion, defVersion, "relationships")
	files, err := relationshipFiles(source)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read relationships for %s %s: %w", modelName, previousVersion, err)
	}
	if len(files) == 0 {
		return nil
	}

	destinationFiles, err := relationshipFiles(destination)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read relationships for new %s %s: %w", modelName, newVersion, err)
	}
	existingFiles := make(map[string]struct{}, len(destinationFiles))
	for _, file := range destinationFiles {
		existingFiles[file.Name()] = struct{}{}
	}

	previousKinds, err := modelComponentKinds(filepath.Join(modelPath, previousVersion, defVersion, "components"))
	if err != nil {
		return fmt.Errorf("read components for %s %s: %w", modelName, previousVersion, err)
	}
	componentKinds, err := modelComponentKinds(filepath.Join(modelPath, newVersion, defVersion, "components"))
	if err != nil {
		return fmt.Errorf("read components for new %s %s: %w", modelName, newVersion, err)
	}

	var checkErrors []error
	missingFiles := make([]os.DirEntry, 0, len(files))
	for _, file := range files {
		if _, exists := existingFiles[file.Name()]; exists {
			continue
		}
		staleKinds, err := relationshipKindsMissingFromModel(filepath.Join(source, file.Name()), modelName, previousKinds, componentKinds)
		if err != nil {
			checkErrors = append(checkErrors, fmt.Errorf("check relationship %s: %w", file.Name(), err))
			continue
		}
		if len(staleKinds) > 0 {
			if utils.Log != nil {
				utils.Log.Warnf("Skipping inherited relationship %s for %s %s: referenced component kinds are missing: %s", file.Name(), modelName, newVersion, strings.Join(staleKinds, ", "))
			}
			continue
		}
		missingFiles = append(missingFiles, file)
	}
	if len(missingFiles) == 0 {
		return errors.Join(checkErrors...)
	}
	return errors.Join(append(checkErrors, copyRelationships(missingFiles, source, destination, newVersion))...)
}

// failedVersionsBetween lists the failed versions newer than source and older
// than version, which version would otherwise have inherited from.
func failedVersionsBetween(failed map[string]struct{}, source, version string) []string {
	upper, err := semver.NewVersion(version)
	if err != nil {
		return nil
	}
	var lower *semver.Version
	if source != "" {
		if lower, err = semver.NewVersion(source); err != nil {
			return nil
		}
	}

	skipped := make([]string, 0)
	for candidate := range failed {
		parsed, err := semver.NewVersion(candidate)
		if err != nil || !parsed.LessThan(upper) || (lower != nil && !parsed.GreaterThan(lower)) {
			continue
		}
		skipped = append(skipped, candidate)
	}
	sort.Strings(skipped)
	return skipped
}

func modelComponentKinds(directory string) (map[string]struct{}, error) {
	files, err := relationshipFiles(directory)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]struct{}{}, nil
	}
	if err != nil {
		return nil, err
	}

	kinds := make(map[string]struct{}, len(files))
	for _, file := range files {
		path := filepath.Join(directory, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		var definition struct {
			Component struct {
				Kind string `json:"kind"`
			} `json:"component"`
		}
		if err := json.Unmarshal(data, &definition); err != nil {
			return nil, fmt.Errorf("decode %q: %w", path, err)
		}
		kind := definition.Component.Kind
		if kind == "" {
			kind = strings.TrimSuffix(file.Name(), ".json")
		}
		kinds[kind] = struct{}{}
	}
	return kinds, nil
}

func relationshipKindsMissingFromModel(path, modelName string, previousKinds, componentKinds map[string]struct{}) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var relationship struct {
		Selectors []struct {
			Allow struct {
				From []relationshipEndpoint `json:"from"`
				To   []relationshipEndpoint `json:"to"`
			} `json:"allow"`
		} `json:"selectors"`
	}
	if err := json.Unmarshal(data, &relationship); err != nil {
		return nil, fmt.Errorf("decode relationship: %w", err)
	}

	missing := make(map[string]struct{})
	checkKind := func(endpoint relationshipEndpoint) {
		// Kinds absent from the previous version, such as "self", are not
		// component kinds this model could have dropped.
		if endpoint.Kind == "" || (endpoint.Model.Name != "" && endpoint.Model.Name != modelName) {
			return
		}
		if _, existed := previousKinds[endpoint.Kind]; !existed {
			return
		}
		if _, exists := componentKinds[endpoint.Kind]; !exists {
			missing[endpoint.Kind] = struct{}{}
		}
	}
	for _, selector := range relationship.Selectors {
		for _, endpoint := range append(selector.Allow.From, selector.Allow.To...) {
			checkKind(endpoint)
			for _, matched := range append(endpoint.Match.From, endpoint.Match.To...) {
				checkKind(matched)
			}
		}
	}

	kinds := make([]string, 0, len(missing))
	for kind := range missing {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	return kinds, nil
}

type relationshipEndpoint struct {
	Kind  string `json:"kind"`
	Model struct {
		Name string `json:"name"`
	} `json:"model"`
	Match struct {
		From []relationshipEndpoint `json:"from"`
		To   []relationshipEndpoint `json:"to"`
	} `json:"match"`
}

func latestPreviousVersion(version string, versions map[string]struct{}) string {
	current, err := semver.NewVersion(version)
	if err != nil {
		return ""
	}

	var previous string
	var previousSemver *semver.Version
	for candidate := range versions {
		parsed, err := semver.NewVersion(candidate)
		if err != nil || !parsed.LessThan(current) {
			continue
		}
		if previousSemver == nil || parsed.GreaterThan(previousSemver) {
			previous = candidate
			previousSemver = parsed
		}
	}

	return previous
}

// relationshipFiles lists the regular .json files in directory. Component
// directories use the same layout, so it serves both.
func relationshipFiles(directory string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	files := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("inspect %q: %w", filepath.Join(directory, entry.Name()), err)
		}
		if info.Mode().IsRegular() {
			files = append(files, entry)
		}
	}
	return files, nil
}

// copyRelationships copies the files that can be read and updated, and
// reports the ones that cannot. A write failure is all or nothing: it removes
// the files this call already wrote so a version is never left half carried.
func copyRelationships(files []os.DirEntry, source, destination, version string) error {
	type fileCopy struct {
		name string
		data []byte
		mode os.FileMode
	}

	var readErrors []error
	copies := make([]fileCopy, 0, len(files))
	for _, file := range files {
		path := filepath.Join(source, file.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			readErrors = append(readErrors, fmt.Errorf("read %q: %w", path, err))
			continue
		}
		data, err = setRelationshipModelVersion(data, version)
		if err != nil {
			readErrors = append(readErrors, fmt.Errorf("update embedded model version in %q: %w", path, err))
			continue
		}
		info, err := file.Info()
		if err != nil {
			readErrors = append(readErrors, fmt.Errorf("read file info for %q: %w", path, err))
			continue
		}
		copies = append(copies, fileCopy{name: file.Name(), data: data, mode: info.Mode().Perm()})
	}
	if len(copies) == 0 {
		return errors.Join(readErrors...)
	}

	_, statErr := os.Stat(destination)
	createdDestination := errors.Is(statErr, os.ErrNotExist)
	if err := os.MkdirAll(destination, 0755); err != nil {
		return errors.Join(append(readErrors, fmt.Errorf("create %q: %w", destination, err))...)
	}

	written := make([]string, 0, len(copies))
	for _, file := range copies {
		path := filepath.Join(destination, file.name)
		if err := writeNewFile(path, file.data, file.mode); err != nil {
			for _, path := range written {
				_ = os.Remove(path)
			}
			if createdDestination {
				_ = os.Remove(destination)
			}
			return errors.Join(append(readErrors, err)...)
		}
		written = append(written, path)
	}
	return errors.Join(readErrors...)
}

// writeNewFile creates path, failing if it already exists. A file it created
// but could not finish writing is removed.
func writeNewFile(path string, data []byte, mode os.FileMode) error {
	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	if _, err := output.Write(data); err != nil {
		_ = output.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write %q: %w", path, err)
	}
	if err := output.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("close %q: %w", path, err)
	}
	return nil
}

func setRelationshipModelVersion(data []byte, version string) ([]byte, error) {
	var relationship map[string]json.RawMessage
	if err := json.Unmarshal(data, &relationship); err != nil {
		return nil, fmt.Errorf("decode relationship: %w", err)
	}

	var model map[string]json.RawMessage
	if err := json.Unmarshal(relationship["model"], &model); err != nil || model == nil {
		return nil, errors.New("relationship has no model definition")
	}
	var modelVersion map[string]json.RawMessage
	if err := json.Unmarshal(model["model"], &modelVersion); err != nil || modelVersion == nil {
		return nil, errors.New("relationship has no embedded model version")
	}
	if _, ok := modelVersion["version"]; !ok {
		return nil, errors.New("relationship has no embedded model version")
	}

	encodedVersion, err := json.Marshal(version)
	if err != nil {
		return nil, fmt.Errorf("encode model version: %w", err)
	}
	modelVersion["version"] = encodedVersion
	modelBytes, err := json.Marshal(modelVersion)
	if err != nil {
		return nil, fmt.Errorf("encode embedded model: %w", err)
	}
	model["model"] = modelBytes
	updatedModel, err := json.Marshal(model)
	if err != nil {
		return nil, fmt.Errorf("encode relationship model: %w", err)
	}
	relationship["model"] = updatedModel
	relationshipBytes, err := json.MarshalIndent(relationship, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode relationship: %w", err)
	}
	return append(relationshipBytes, '\n'), nil
}
