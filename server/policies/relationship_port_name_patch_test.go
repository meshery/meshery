package policies

import (
	"path/filepath"
	"strconv"
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

// countTargetPortToPortNamePatches counts mutatorRef/mutatedRef pairs that copy
// a ports[*].targetPort value into a ports[*].name. It pairs refs the way the
// engine does: for every From/To item combination in a selector,
// resolveMutatorMutatedRefs picks the mutator and mutated side, and
// mutatorRef[i] patches mutatedRef[i] up to the shorter list
// (identifyRelationshipsBasedOnMatchingMutatorAndMutatedFields, matchingMutators).
func countTargetPortToPortNamePatches(definition *relationship.RelationshipDefinition) int {
	if definition.Selectors == nil {
		return 0
	}
	count := 0
	for _, selector := range *definition.Selectors {
		for _, from := range selector.Allow.From {
			for _, to := range selector.Allow.To {
				if from.RelationshipDefinitionSelectorsPatch == nil || to.RelationshipDefinitionSelectorsPatch == nil {
					continue
				}
				mutators, mutated, _, _ := resolveMutatorMutatedRefs(from.RelationshipDefinitionSelectorsPatch, to.RelationshipDefinitionSelectorsPatch, nil, nil)
				for i := 0; i < len(mutators) && i < len(mutated); i++ {
					if isPortField(mutators[i], "targetPort") && isPortField(mutated[i], "name") {
						count++
					}
				}
			}
		}
	}
	return count
}

// isPortField reports whether path ends in ports, <numeric index>, field, so
// only array entries such as ports/0/name match.
func isPortField(path []string, field string) bool {
	n := len(path)
	if n < 3 {
		return false
	}
	_, err := strconv.Atoi(path[n-2])
	return err == nil && path[n-1] == field && path[n-3] == "ports"
}

func TestCountTargetPortToPortNamePatches(t *testing.T) {
	targetPort := []string{"configuration", "spec", "ports", "0", "targetPort"}
	containerPort := []string{"configuration", "spec", "template", "spec", "containers", "0", "ports", "0", "containerPort"}
	portName := []string{"configuration", "spec", "template", "spec", "containers", "0", "ports", "0", "name"}
	selectorRef := []string{"configuration", "spec", "selector"}
	matchLabels := []string{"configuration", "spec", "selector", "matchLabels"}

	mutator := func(refs ...[]string) relationship.SelectorItem {
		ref := relationship.MutatorRef(refs)
		return relationship.SelectorItem{RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{MutatorRef: &ref}}
	}
	mutated := func(refs ...[]string) relationship.SelectorItem {
		ref := relationship.MutatedRef(refs)
		return relationship.SelectorItem{RelationshipDefinitionSelectorsPatch: &relationship.RelationshipDefinitionSelectorsPatch{MutatedRef: &ref}}
	}
	definition := func(from, to []relationship.SelectorItem) *relationship.RelationshipDefinition {
		selectors := relationship.SelectorSet{{Allow: relationship.Selector{From: from, To: to}}}
		return &relationship.RelationshipDefinition{Selectors: &selectors}
	}

	tests := []struct {
		name       string
		definition *relationship.RelationshipDefinition
		want       int
	}{
		{
			name:       "nil selectors",
			definition: &relationship.RelationshipDefinition{},
			want:       0,
		},
		{
			name:       "targetPort patched into containerPort",
			definition: definition([]relationship.SelectorItem{mutator(selectorRef, targetPort)}, []relationship.SelectorItem{mutated(matchLabels, containerPort)}),
			want:       0,
		},
		{
			name:       "targetPort patched into port name",
			definition: definition([]relationship.SelectorItem{mutator(selectorRef, targetPort)}, []relationship.SelectorItem{mutated(matchLabels, portName)}),
			want:       1,
		},
		{
			name:       "mutator on the To side and mutated on the From side",
			definition: definition([]relationship.SelectorItem{mutated(portName)}, []relationship.SelectorItem{mutator(targetPort)}),
			want:       1,
		},
		{
			name:       "uneven ref lists pair only up to the shorter list",
			definition: definition([]relationship.SelectorItem{mutator(selectorRef, targetPort, targetPort)}, []relationship.SelectorItem{mutated(matchLabels, portName)}),
			want:       1,
		},
		{
			name:       "mutatedRef longer than mutatorRef",
			definition: definition([]relationship.SelectorItem{mutator(targetPort)}, []relationship.SelectorItem{mutated(containerPort, portName)}),
			want:       0,
		},
		{
			// Flattening both From items would pair selectorRef with portName and
			// miss the targetPort patch; the engine pairs each From/To item on its own.
			name: "several From items are paired with each To item separately",
			definition: definition(
				[]relationship.SelectorItem{mutator(selectorRef), mutator(targetPort)},
				[]relationship.SelectorItem{mutated(portName)},
			),
			want: 1,
		},
		{
			name: "path segment before name is not a numeric index",
			definition: definition(
				[]relationship.SelectorItem{mutator(targetPort)},
				[]relationship.SelectorItem{mutated([]string{"configuration", "spec", "ports", "http", "name"})},
			),
			want: 0,
		},
		{
			name:       "selector item without a patch",
			definition: definition([]relationship.SelectorItem{{}}, []relationship.SelectorItem{mutated(portName)}),
			want:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countTargetPortToPortNamePatches(tt.definition); got != tt.want {
				t.Errorf("countTargetPortToPortNamePatches() = %d, want %d", got, tt.want)
			}
		})
	}
}
