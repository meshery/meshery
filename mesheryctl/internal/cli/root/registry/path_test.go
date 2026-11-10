package registry

import (
	"os"
	"path/filepath"
	"testing"

	meshkiterrors "github.com/meshery/meshkit/errors"
)

func TestDefaultModelsLocation(t *testing.T) {
	// Lay out a fake repo root containing a models directory, a server
	// directory with its own models directory, and a mesheryctl subdirectory,
	// mirroring how the command is actually run.
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
	serverDir := filepath.Join(root, "server")
	if err := os.MkdirAll(filepath.Join(serverDir, "models"), 0o755); err != nil {
		t.Fatalf("mkdir server/models: %v", err)
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
			name: "from subdirectory having its own models directory",
			cwd:  serverDir,
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
			t.Chdir(tt.cwd)
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

	// Repository with marker directories but malformed go.mod.
	malformedGoModRoot := t.TempDir()
	malformedGoModRoot, err = filepath.EvalSymlinks(malformedGoModRoot)
	if err != nil {
		t.Fatalf("evalsymlinks: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(malformedGoModRoot, "models"), 0o755); err != nil {
		t.Fatalf("mkdir models: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(malformedGoModRoot, "mesheryctl"), 0o755); err != nil {
		t.Fatalf("mkdir mesheryctl: %v", err)
	}
	if err := os.WriteFile(filepath.Join(malformedGoModRoot, "go.mod"), []byte("invalid go.mod syntax !!!"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
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
			name:    "repository with matching folder markers but malformed go.mod fails",
			cwd:     malformedGoModRoot,
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
			t.Chdir(tt.cwd)
			err := verifyPrerequisiteRepo()
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyPrerequisiteRepo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil {
				if code := meshkiterrors.GetCode(err); code != ErrPrereqRepositoryCode {
					t.Errorf("verifyPrerequisiteRepo() errorCode = %q, want %q", code, ErrPrereqRepositoryCode)
				}
			}
		})
	}
}

