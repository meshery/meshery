package policies

import (
	"path/filepath"
	"testing"

	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// knownLegacyPortNamePatches lists seeded relationship files that predate this
// test and still patch a container port name from a Service targetPort
// (meshery/meshery#23971). The selector comes from the Integrations sheet, so
// an entry can be removed only after the sheet is corrected and the nightly
// model sync regenerates the file.
var knownLegacyPortNamePatches = map[string]bool{
	"kubernetes/edge-non-binding-network-duixv.json": true,
}

// A Service's spec.ports[*].targetPort is an IntOrString and is usually a
// number such as 80, while a container's ports[*].name must be an IANA service
// name (a string). A selector that copies targetPort into a port name produces
// a design that fails Validate and Dry Run with
// `Field ".spec.template.spec.containers.0.ports.0.name" must be string`
// (meshery/meshery#23971). targetPort already maps to containerPort through a
// separate selector, so no seeded relationship should write it into a name.
func TestRelationshipsDoNotPatchPortNameFromTargetPort(t *testing.T) {
	root := repoRoot(t)

	violations := map[string]int{}
	legacy := map[string]int{}
	for _, s := range seededModelDefDirs(t, root) {
		files, err := filepath.Glob(filepath.Join(s.dir, "relationships", "*.json"))
		if err != nil {
			t.Fatalf("globbing relationships in %s: %v", s.dir, err)
		}
		for _, file := range files {
			definition, err := loadRelationshipDefinition(file)
			if err != nil {
				t.Fatal(err)
			}
			count := countTargetPortToPortNamePatches(definition)
			if count == 0 {
				continue
			}
			key := s.name + "/" + filepath.Base(file)
			if knownLegacyPortNamePatches[key] {
				legacy[key] += count
			} else {
				violations[key] += count
			}
		}
	}

	if len(legacy) > 0 {
		total, summary := summarizeCounts(legacy)
		t.Logf("%d legacy selectors patch a port name from targetPort (meshery/meshery#23971): %s", total, summary)
	}
	for key := range knownLegacyPortNamePatches {
		if legacy[key] == 0 {
			t.Logf("%s no longer patches a port name from targetPort; remove it from knownLegacyPortNamePatches", key)
		}
	}
	if len(violations) > 0 {
		total, summary := summarizeCounts(violations)
		t.Errorf("%d selectors patch a container port name from a Service targetPort, which is usually a number (meshery/meshery#23971): %s", total, summary)
	}
}

// countTargetPortToPortNamePatches counts mutatorRef/mutatedRef pairs, within
// one selector, that copy a ports[*].targetPort value into a ports[*].name.
// The two lists are paired by position: mutatorRef[i] patches mutatedRef[i].
func countTargetPortToPortNamePatches(definition *relationship.RelationshipDefinition) int {
	if definition.Selectors == nil {
		return 0
	}
	count := 0
	for _, selector := range *definition.Selectors {
		items := append(append([]relationship.SelectorItem{}, selector.Allow.From...), selector.Allow.To...)
		var mutators, mutated [][]string
		for _, item := range items {
			patch := item.RelationshipDefinitionSelectorsPatch
			if patch == nil {
				continue
			}
			if patch.MutatorRef != nil {
				mutators = append(mutators, *patch.MutatorRef...)
			}
			if patch.MutatedRef != nil {
				mutated = append(mutated, *patch.MutatedRef...)
			}
		}
		for i := 0; i < len(mutators) && i < len(mutated); i++ {
			if isPortField(mutators[i], "targetPort") && isPortField(mutated[i], "name") {
				count++
			}
		}
	}
	return count
}

// isPortField reports whether path ends in ports, <index>, field.
func isPortField(path []string, field string) bool {
	n := len(path)
	return n >= 3 && path[n-1] == field && path[n-3] == "ports"
}
