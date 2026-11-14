package policies

import (
	"path/filepath"
	"testing"
)

// A selector's model.name is compared with a component's model name through
// matchName (isRelationshipFeasible). A name that is not a registered model
// can never match, and a name that matches only because matchName is an
// unanchored regex stops matching as soon as matchName is anchored
// (meshery/meshery#22147). Every selector in the seeded relationships must name
// a registered model exactly (meshery/meshery#23092).
func TestRelationshipSelectorsNameRegisteredModels(t *testing.T) {
	root := repoRoot(t)
	seeded := seededModelDefDirs(t, root)

	registered := map[string]bool{}
	for _, s := range seeded {
		registered[s.name] = true
	}

	missing := map[string]int{}
	substringOnly := map[string]int{}
	for _, s := range seeded {
		files, err := filepath.Glob(filepath.Join(s.dir, "relationships", "*.json"))
		if err != nil {
			t.Fatalf("globbing relationships in %s: %v", s.dir, err)
		}
		for _, file := range files {
			definition, err := loadRelationshipDefinition(file)
			if err != nil {
				t.Fatal(err)
			}
			for _, item := range allSelectorItems(definition) {
				if item.Model == nil {
					continue
				}
				name := item.Model.Name
				if name == "" || name == "*" || registered[name] {
					continue
				}
				matched := false
				for model := range registered {
					if matchName(model, name) {
						matched = true
						break
					}
				}
				if matched {
					substringOnly[name]++
				} else {
					missing[name]++
				}
			}
		}
	}

	if len(missing) > 0 {
		total, summary := summarizeCounts(missing)
		t.Errorf("%d selectors name a model that does not exist: %s", total, summary)
	}
	if len(substringOnly) > 0 {
		total, summary := summarizeCounts(substringOnly)
		t.Errorf("%d selectors match a model only by unanchored substring: %s", total, summary)
	}
}
