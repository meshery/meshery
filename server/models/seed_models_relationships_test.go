package models

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSeededModelsRetainRelationships guards issue #23071. Meshery seeds only a
// model's newest version directory (GetModelDirectoryPaths uses
// sortedVersionDirs[0]), so a relationship that exists only in an older version
// is never registered. cert-manager and dapr both regressed this way when the
// model generator added newer version directories without carrying their
// relationships forward.
//
// The test reads the directories GetModelDirectoryPaths would seed and asserts
// that the newest version of each affected model still carries its
// relationships. It fails on the pre-fix tree (the newest version directories
// have no relationships/ directory) and passes once the relationships are
// present in the seeded version.
func TestSeededModelsRetainRelationships(t *testing.T) {
	seededDirs, err := GetModelDirectoryPaths(ModelsPath)
	if err != nil {
		t.Fatalf("GetModelDirectoryPaths(%q): %v", ModelsPath, err)
	}

	// Models known to carry relationships; each must have at least this many
	// relationship files in its seeded (newest) version directory.
	wantRelationships := map[string]int{
		"cert-manager": 3,
		"dapr":         1,
	}

	for model, want := range wantRelationships {
		needle := "/" + model + "/"
		var defDir string
		for _, dir := range seededDirs {
			if strings.Contains(filepath.ToSlash(dir), needle) {
				defDir = dir
				break
			}
		}
		if defDir == "" {
			t.Errorf("model %q has no seeded directory", model)
			continue
		}

		relDir := filepath.Join(defDir, "relationships")
		entries, err := os.ReadDir(relDir)
		if err != nil {
			t.Errorf("model %q: seeded directory %q has no relationships: %v", model, defDir, err)
			continue
		}

		got := 0
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
				got++
			}
		}
		if got < want {
			t.Errorf("model %q: seeded directory %q has %d relationship files, want at least %d", model, defDir, got, want)
		}
	}
}
