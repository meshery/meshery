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
	nestedSubDir := filepath.Join(subDir, "internal", "cli")
	if err := os.MkdirAll(nestedSubDir, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/meshery/meshery\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
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
			name: "from nested subdirectory",
			cwd:  nestedSubDir,
			want: filepath.Join("..", "..", "..", "models"),
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

func TestVerifyPrerequisiteRepo(t *testing.T) {
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "models"), 0o755); err != nil {
		t.Fatalf("mkdir models: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	subDir := filepath.Join(root, "mesheryctl")
	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("mkdir mesheryctl: %v", err)
	}
	nestedSubDir := filepath.Join(subDir, "internal", "cli")
	if err := os.MkdirAll(nestedSubDir, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module github.com/meshery/meshery\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	// Unrelated repository with matching marker directories but a different module name.
	unrelatedRoot := t.TempDir()
	unrelatedRoot, err = filepath.EvalSymlinks(unrelatedRoot)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(unrelatedRoot, "models"), 0o755); err != nil {
		t.Fatalf("mkdir models: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(unrelatedRoot, "server"), 0o755); err != nil {
		t.Fatalf("mkdir server: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(unrelatedRoot, "mesheryctl"), 0o755); err != nil {
		t.Fatalf("mkdir mesheryctl: %v", err)
	}
	if err := os.WriteFile(filepath.Join(unrelatedRoot, "go.mod"), []byte("module github.com/unrelated/project\n\ngo 1.26\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	// Repository with marker directories but missing go.mod.
	missingGoModRoot := t.TempDir()
	missingGoModRoot, err = filepath.EvalSymlinks(missingGoModRoot)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(missingGoModRoot, "models"), 0o755); err != nil {
		t.Fatalf("mkdir models: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(missingGoModRoot, "mesheryctl"), 0o755); err != nil {
		t.Fatalf("mkdir mesheryctl: %v", err)
	}

	emptyDir := t.TempDir()

	tests := []struct {
		name    string
		cwd     string
		wantErr bool
	}{
		{
			name:    "repo root succeeds",
			cwd:     root,
			wantErr: false,
		},
		{
			name:    "immediate subdirectory succeeds",
			cwd:     subDir,
			wantErr: false,
		},
		{
			name:    "nested subdirectory within repo succeeds",
			cwd:     nestedSubDir,
			wantErr: false,
		},
		{
			name:    "unrelated repository with matching folder markers but non-meshery module fails",
			cwd:     unrelatedRoot,
			wantErr: true,
		},
		{
			name:    "unrelated repository with matching folder markers but missing go.mod fails",
			cwd:     missingGoModRoot,
			wantErr: true,
		},
		{
			name:    "empty directory outside repo fails",
			cwd:     emptyDir,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chdir(t, tt.cwd)
			err := verifyPrerequisiteRepo()
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyPrerequisiteRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

