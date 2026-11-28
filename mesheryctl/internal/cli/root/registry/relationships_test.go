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
	"os"
	"path/filepath"
	"testing"
)

func TestCarryForwardRelationshipsToNewModelVersion(t *testing.T) {
	modelsPath := t.TempDir()
	modelPath := filepath.Join(modelsPath, "sample-model")
	previousRelationships := filepath.Join(modelPath, "v1.2.0", defVersion, "relationships")
	if err := os.MkdirAll(previousRelationships, 0755); err != nil {
		t.Fatal(err)
	}

	previousDefinition := []byte(`{"name":"sample-edge","model":{"name":"sample-model","model":{"version":"v1.2.0"}}}`)
	previousPath := filepath.Join(previousRelationships, "sample-edge.json")
	if err := os.WriteFile(previousPath, previousDefinition, 0644); err != nil {
		t.Fatal(err)
	}

	before, err := snapshotModelVersions(modelsPath)
	if err != nil {
		t.Fatalf("snapshot model versions: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(modelPath, "v1.3.0", defVersion, "components"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := carryForwardRelationships(modelsPath, before); err != nil {
		t.Fatalf("carry forward relationships: %v", err)
	}

	carriedPath := filepath.Join(modelPath, "v1.3.0", defVersion, "relationships", "sample-edge.json")
	carriedDefinition, err := os.ReadFile(carriedPath)
	if err != nil {
		t.Fatalf("read carried relationship: %v", err)
	}
	var relationship struct {
		Name  string `json:"name"`
		Model struct {
			Model struct {
				Version string `json:"version"`
			} `json:"model"`
		} `json:"model"`
	}
	if err := json.Unmarshal(carriedDefinition, &relationship); err != nil {
		t.Fatalf("decode carried relationship: %v", err)
	}
	if relationship.Name != "sample-edge" {
		t.Errorf("carried relationship name = %q, want sample-edge", relationship.Name)
	}
	if relationship.Model.Model.Version != "v1.3.0" {
		t.Errorf("embedded model version = %q, want v1.3.0", relationship.Model.Model.Version)
	}

	unchangedDefinition, err := os.ReadFile(previousPath)
	if err != nil {
		t.Fatalf("read previous relationship: %v", err)
	}
	if string(unchangedDefinition) != string(previousDefinition) {
		t.Errorf("previous relationship was modified: got %s, want %s", unchangedDefinition, previousDefinition)
	}
}

func TestCarryForwardRelationshipsMergesMissingDefinitions(t *testing.T) {
	modelsPath := t.TempDir()
	modelPath := filepath.Join(modelsPath, "sample-model")
	previousRelationships := filepath.Join(modelPath, "v1.2.0", defVersion, "relationships")
	if err := os.MkdirAll(previousRelationships, 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"inherited-edge.json", "generated-edge.json"} {
		definition := []byte(`{"model":{"model":{"version":"v1.2.0"}}}`)
		if err := os.WriteFile(filepath.Join(previousRelationships, name), definition, 0644); err != nil {
			t.Fatal(err)
		}
	}

	before, err := snapshotModelVersions(modelsPath)
	if err != nil {
		t.Fatalf("snapshot model versions: %v", err)
	}
	newRelationships := filepath.Join(modelPath, "v1.3.0", defVersion, "relationships")
	if err := os.MkdirAll(newRelationships, 0755); err != nil {
		t.Fatal(err)
	}
	generatedDefinition := []byte(`{"model":{"model":{"version":"v1.3.0"}},"source":"generated"}`)
	if err := os.WriteFile(filepath.Join(newRelationships, "generated-edge.json"), generatedDefinition, 0644); err != nil {
		t.Fatal(err)
	}
	newOnlyDefinition := []byte(`{"model":{"model":{"version":"v1.3.0"}},"source":"generated"}`)
	if err := os.WriteFile(filepath.Join(newRelationships, "new-edge.json"), newOnlyDefinition, 0644); err != nil {
		t.Fatal(err)
	}

	if err := carryForwardRelationships(modelsPath, before); err != nil {
		t.Fatalf("carry forward relationships: %v", err)
	}

	inheritedDefinition, err := os.ReadFile(filepath.Join(newRelationships, "inherited-edge.json"))
	if err != nil {
		t.Fatalf("read inherited relationship: %v", err)
	}
	var inherited struct {
		Model struct {
			Model struct {
				Version string `json:"version"`
			} `json:"model"`
		} `json:"model"`
	}
	if err := json.Unmarshal(inheritedDefinition, &inherited); err != nil {
		t.Fatalf("decode inherited relationship: %v", err)
	}
	if inherited.Model.Model.Version != "v1.3.0" {
		t.Errorf("inherited relationship version = %q, want v1.3.0", inherited.Model.Model.Version)
	}

	got, err := os.ReadFile(filepath.Join(newRelationships, "generated-edge.json"))
	if err != nil {
		t.Fatalf("read generated relationship: %v", err)
	}
	if string(got) != string(generatedDefinition) {
		t.Errorf("generated relationship changed: got %s, want %s", got, generatedDefinition)
	}
	got, err = os.ReadFile(filepath.Join(newRelationships, "new-edge.json"))
	if err != nil {
		t.Fatalf("read new generated relationship: %v", err)
	}
	if string(got) != string(newOnlyDefinition) {
		t.Errorf("new generated relationship changed: got %s, want %s", got, newOnlyDefinition)
	}
}

func TestCarryForwardRelationshipsSkipsRemovedComponentKinds(t *testing.T) {
	modelsPath := t.TempDir()
	modelPath := filepath.Join(modelsPath, "sample-model")
	previousRelationships := filepath.Join(modelPath, "v1.2.0", defVersion, "relationships")
	if err := os.MkdirAll(previousRelationships, 0755); err != nil {
		t.Fatal(err)
	}
	staleDefinition := []byte(`{"model":{"model":{"version":"v1.2.0"}},"selectors":[{"allow":{"from":[{"kind":"ExistingKind","model":{"name":"sample-model"}}],"to":[{"kind":"RemovedKind","model":{"name":"sample-model"}}]}}]}`)
	if err := os.WriteFile(filepath.Join(previousRelationships, "stale-edge.json"), staleDefinition, 0644); err != nil {
		t.Fatal(err)
	}
	externalDefinition := []byte(`{"model":{"model":{"version":"v1.2.0"}},"selectors":[{"allow":{"from":[{"kind":"ExistingKind","model":{"name":"sample-model"}}],"to":[{"kind":"ExternalKind","model":{"name":"another-model"}}]}}]}`)
	if err := os.WriteFile(filepath.Join(previousRelationships, "external-edge.json"), externalDefinition, 0644); err != nil {
		t.Fatal(err)
	}
	unlistedDefinition := []byte(`{"model":{"model":{"version":"v1.2.0"}},"selectors":[{"allow":{"from":[{"kind":"ExistingKind","model":{"name":"sample-model"}}],"to":[{"kind":"UnlistedKind","model":{"name":"sample-model"}}]}}]}`)
	if err := os.WriteFile(filepath.Join(previousRelationships, "unlisted-edge.json"), unlistedDefinition, 0644); err != nil {
		t.Fatal(err)
	}
	previousComponents := filepath.Join(modelPath, "v1.2.0", defVersion, "components")
	if err := os.MkdirAll(previousComponents, 0755); err != nil {
		t.Fatal(err)
	}
	for _, kind := range []string{"ExistingKind", "RemovedKind"} {
		if err := os.WriteFile(filepath.Join(previousComponents, kind+".json"), []byte(`{}`), 0644); err != nil {
			t.Fatal(err)
		}
	}

	before, err := snapshotModelVersions(modelsPath)
	if err != nil {
		t.Fatalf("snapshot model versions: %v", err)
	}
	components := filepath.Join(modelPath, "v1.3.0", defVersion, "components")
	if err := os.MkdirAll(components, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(components, "ExistingKind.json"), []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := carryForwardRelationships(modelsPath, before); err != nil {
		t.Fatalf("carry forward relationships: %v", err)
	}
	newRelationships := filepath.Join(modelPath, "v1.3.0", defVersion, "relationships")
	if _, err := os.Stat(filepath.Join(newRelationships, "stale-edge.json")); !os.IsNotExist(err) {
		t.Errorf("relationship referencing a removed component kind was copied; stat error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(newRelationships, "external-edge.json")); err != nil {
		t.Errorf("relationship referencing another model's kind was not copied: %v", err)
	}
	if _, err := os.Stat(filepath.Join(newRelationships, "unlisted-edge.json")); err != nil {
		t.Errorf("relationship referencing a kind absent from both versions was not copied: %v", err)
	}
}

func TestCarryForwardRelationshipsOnlyTouchesNewVersions(t *testing.T) {
	modelsPath := t.TempDir()
	modelPath := filepath.Join(modelsPath, "sample-model")
	relationships := filepath.Join(modelPath, "v1.2.0", defVersion, "relationships")
	if err := os.MkdirAll(relationships, 0755); err != nil {
		t.Fatal(err)
	}
	definition := []byte(`{"model":{"model":{"version":"v1.2.0"}}}`)
	if err := os.WriteFile(filepath.Join(relationships, "edge.json"), definition, 0644); err != nil {
		t.Fatal(err)
	}
	existingComponents := filepath.Join(modelPath, "v1.3.0", defVersion, "components")
	if err := os.MkdirAll(existingComponents, 0755); err != nil {
		t.Fatal(err)
	}

	before, err := snapshotModelVersions(modelsPath)
	if err != nil {
		t.Fatalf("snapshot model versions: %v", err)
	}
	if err := carryForwardRelationships(modelsPath, before); err != nil {
		t.Fatalf("carry forward relationships: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(relationships, "edge.json"))
	if err != nil {
		t.Fatalf("read existing relationship: %v", err)
	}
	if string(got) != string(definition) {
		t.Errorf("existing relationship changed: got %s, want %s", got, definition)
	}
	if _, err := os.Stat(filepath.Join(modelPath, "v1.3.0", defVersion, "relationships", "edge.json")); !os.IsNotExist(err) {
		t.Errorf("relationship was copied into a version that existed before generation; stat error = %v", err)
	}
}

func TestRelationshipKindsMissingFromModelChecksMatchSelectors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "policy-edge.json")
	definition := []byte(`{"selectors":[{"allow":{"from":[{"kind":"Pod","match":{"from":[{"kind":"self"}],"to":[{"kind":"NetworkPolicy"}]}}],"to":[{"kind":"Pod"}]}}]}`)
	if err := os.WriteFile(path, definition, 0644); err != nil {
		t.Fatal(err)
	}

	previousKinds := map[string]struct{}{"Pod": {}, "NetworkPolicy": {}}
	componentKinds := map[string]struct{}{"Pod": {}}
	missing, err := relationshipKindsMissingFromModel(path, "sample-model", previousKinds, componentKinds)
	if err != nil {
		t.Fatalf("check relationship kinds: %v", err)
	}
	if len(missing) != 1 || missing[0] != "NetworkPolicy" {
		t.Errorf("missing kinds = %v, want [NetworkPolicy]", missing)
	}
}

func TestCopyRelationshipsRemovesPartialWrites(t *testing.T) {
	source := filepath.Join(t.TempDir(), "relationships")
	if err := os.MkdirAll(source, 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a-edge.json", "b-edge.json"} {
		if err := os.WriteFile(filepath.Join(source, name), []byte(`{"model":{"model":{"version":"v1.2.0"}}}`), 0644); err != nil {
			t.Fatal(err)
		}
	}
	files, err := relationshipFiles(source)
	if err != nil {
		t.Fatal(err)
	}

	destination := filepath.Join(t.TempDir(), "relationships")
	// A directory with the second file's name makes its create fail after
	// the first file has been written.
	if err := os.MkdirAll(filepath.Join(destination, "b-edge.json"), 0755); err != nil {
		t.Fatal(err)
	}

	if err := copyRelationships(files, source, destination, "v1.3.0"); err == nil {
		t.Fatal("copy relationships succeeded, want an error")
	}
	if _, err := os.Stat(filepath.Join(destination, "a-edge.json")); !os.IsNotExist(err) {
		t.Errorf("partially copied relationship was left behind; stat error = %v", err)
	}
}
