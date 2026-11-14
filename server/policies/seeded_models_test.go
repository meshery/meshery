package policies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	meshkitutils "github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// seededModelDefDir is one model definition directory the server seeds, with
// the model name it registers under.
type seededModelDefDir struct {
	name string
	dir  string
}

// seededModelDefDirs returns the definition directories inside each model's
// newest version directory, picking the newest version with the same meshkit
// ordering that server/models.GetModelDirectoryPaths uses for seeding. Seeding
// takes the most recently modified definition directory of that version; every
// one is returned here so the result does not depend on checkout timestamps,
// and they must all register the same model name, so that whichever one is
// seeded registers the model the tests check against.
func seededModelDefDirs(t *testing.T, root string) []seededModelDefDir {
	t.Helper()
	modelsDir := filepath.Join(root, "models")
	entries, err := os.ReadDir(modelsDir)
	if err != nil {
		t.Fatalf("reading %s: %v", modelsDir, err)
	}
	seeded := []seededModelDefDir{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		versionDirs, err := meshkitutils.GetAllVersionDirsSortedDesc(filepath.Join(modelsDir, entry.Name()))
		if err != nil || len(versionDirs) == 0 {
			continue
		}
		defDirs, err := os.ReadDir(versionDirs[0])
		if err != nil {
			t.Fatalf("reading %s: %v", versionDirs[0], err)
		}
		versionModel := ""
		for _, defDir := range defDirs {
			if !defDir.IsDir() {
				continue
			}
			dir := filepath.Join(versionDirs[0], defDir.Name())
			raw, err := os.ReadFile(filepath.Join(dir, "model.json"))
			if err != nil {
				t.Fatalf("reading model definition in %s: %v", dir, err)
			}
			var model struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(raw, &model); err != nil {
				t.Fatalf("parsing model definition in %s: %v", dir, err)
			}
			if versionModel != "" && model.Name != versionModel {
				t.Fatalf("definition directories in %s register different models (%q and %q), but seeding loads only one of them", versionDirs[0], versionModel, model.Name)
			}
			versionModel = model.Name
			seeded = append(seeded, seededModelDefDir{name: model.Name, dir: dir})
		}
	}
	return seeded
}

// allSelectorItems flattens a definition's allow and deny selector items, from
// and to. The engine matches both against components (feasibility.go).
func allSelectorItems(definition *relationship.RelationshipDefinition) []relationship.SelectorItem {
	if definition.Selectors == nil {
		return nil
	}
	items := []relationship.SelectorItem{}
	for _, selector := range *definition.Selectors {
		items = append(items, selector.Allow.From...)
		items = append(items, selector.Allow.To...)
		if selector.Deny != nil {
			items = append(items, selector.Deny.From...)
			items = append(items, selector.Deny.To...)
		}
	}
	return items
}

// summarizeCounts renders a name->count map as "a (2), b (1)", sorted by name,
// and returns the total count.
func summarizeCounts(counts map[string]int) (int, string) {
	names := make([]string, 0, len(counts))
	total := 0
	for name, count := range counts {
		names = append(names, name)
		total += count
	}
	sort.Strings(names)
	parts := make([]string, 0, len(names))
	for _, name := range names {
		parts = append(parts, fmt.Sprintf("%s (%d)", name, counts[name]))
	}
	return total, strings.Join(parts, ", ")
}
