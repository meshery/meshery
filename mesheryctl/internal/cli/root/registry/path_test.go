package registry

import (
	"os"
	"path/filepath"
	"testing"
)

// chdir changes the working directory for the duration of a test and restores
// it afterwards.
func chdir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir %s: %v", dir, err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
}

func TestDefaultModelsLocation(t *testing.T) {
	// Lay out a fake repo root containing a models directory and a
	// mesheryctl subdirectory, mirroring how the command is actually run.
	root := t.TempDir()
	// On macOS t.TempDir() lives under /private/var symlinks; resolve so the
	// directory probes below compare against a stable path.
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "models"), 0o755); err != nil {
		t.Fatalf("mkdir models: %v", err)
	}
	subDir := filepath.Join(root, "mesheryctl")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("mkdir mesheryctl: %v", err)
	}
	emptyDir := t.TempDir() // no models directory anywhere

	tests := []struct {
		name string
		cwd  string
		want string
	}{
		{
			name: "from repo root",
			cwd:  root,
			want: "models",
		},
		{
			name: "from subdirectory",
			cwd:  subDir,
			want: filepath.Join("..", "models"),
		},
		{
			name: "no models directory falls back to repo-root assumption",
			cwd:  emptyDir,
			want: "models",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chdir(t, tt.cwd)
			if got := defaultModelsLocation(); got != tt.want {
				t.Errorf("defaultModelsLocation() = %q, want %q", got, tt.want)
			}
		})
	}
}
