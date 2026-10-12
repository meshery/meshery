package design

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_getUniqueFilename(t *testing.T) {
	dir := t.TempDir()

	seed := func(t *testing.T, names ...string) {
		t.Helper()
		for _, name := range names {
			if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), 0o600); err != nil {
				t.Fatalf("failed to seed %q: %v", name, err)
			}
		}
	}

	t.Run("returns the original path when no collision exists", func(t *testing.T) {
		want := filepath.Join(dir, "report.yaml")
		if got := getUniqueFilename(want); got != want {
			t.Errorf("getUniqueFilename() = %q, want %q", got, want)
		}
	})

	t.Run("appends counter before the extension on collision", func(t *testing.T) {
		seed(t, "collision.yaml")

		want := filepath.Join(dir, "collision(1).yaml")
		if got := getUniqueFilename(filepath.Join(dir, "collision.yaml")); got != want {
			t.Errorf("getUniqueFilename() = %q, want %q", got, want)
		}
	})

	t.Run("increments the counter until a free name is found", func(t *testing.T) {
		seed(t, "multi.yaml", "multi(1).yaml", "multi(2).yaml")

		want := filepath.Join(dir, "multi(3).yaml")
		if got := getUniqueFilename(filepath.Join(dir, "multi.yaml")); got != want {
			t.Errorf("getUniqueFilename() = %q, want %q", got, want)
		}
	})
}
