package config_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	config "github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestNeedsMutation(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(string, *testing.T)
		path    func(string) string
		want    bool
		wantErr bool
	}{
		{
			name:    "Given missing config file, When NeedsMutation is called, Then it returns true",
			setup:   func(_ string, _ *testing.T) {},
			path:    func(p string) string { return p },
			want:    true,
			wantErr: false,
		},
		{
			name: "Given empty config file, When NeedsMutation is called, Then it returns true",
			setup: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			path:    func(p string) string { return p },
			want:    true,
			wantErr: false,
		},
		{
			name: "Given non-empty config file, When NeedsMutation is called, Then it returns false",
			setup: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			path:    func(p string) string { return p },
			want:    false,
			wantErr: false,
		},
		{
			name:    "Given invalid path, When NeedsMutation is called, Then it returns error",
			setup:   func(_ string, _ *testing.T) {},
			path:    func(_ string) string { return string([]byte{0}) },
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			tt.setup(configPath, t)
			configPath = tt.path(configPath)

			got, err := config.NeedsMutation(configPath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Fatalf("NeedsMutation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitDefaultConfig(t *testing.T) {
	tests := []struct {
	name         string
	setup        func(string, *testing.T)
	configPath   func(string) string
	mesheryFolder func(string) string
	createConfig func(string) func() error
	wantErr      bool
}{
		{
			name: "Given valid setup, When InitDefaultConfig is called, Then no error is returned",
			setup: func(_ string, _ *testing.T) {},
			configPath: func(p string) string { return p },
			createConfig: func(path string) func() error {
				return func() error {
					return os.WriteFile(path, []byte("contexts: {}\n"), 0o644)
				}
			},
			wantErr: false,
		},
		{
			name: "Given createConfigFile fails, When InitDefaultConfig is called, Then error is returned",
			setup: func(_ string, _ *testing.T) {},
			configPath: func(p string) string { return p },
			createConfig: func(_ string) func() error {
				return func() error {
					return errors.New("create failed")
				}
			},
			wantErr: true,
		},
		{
			name: "Given os.MkdirAll fails, When InitDefaultConfig is called, Then error is returned",
			setup: func(_ string, _ *testing.T) {},
			configPath: func(p string) string { return p },
			mesheryFolder: func(_ string) string {
				return string([]byte{0}) // portable invalid path
			},
			createConfig: func(_ string) func() error {
				return func() error { return nil }
			},
			wantErr: true,
		},
		{
			name: "Given AddTokenToConfig fails, When InitDefaultConfig is called, Then error is returned",
			setup: func(path string, t *testing.T) {
				// Create directory instead of file → guaranteed failure
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
			},
			configPath: func(p string) string { return p },
			createConfig: func(_ string) func() error {
				return func() error { return nil }
			},
			wantErr: true,
		},
		{
			name: "Given AddContextToConfig fails, When InitDefaultConfig is called, Then error is returned",
			setup: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte("valid: yaml"), 0o644); err != nil {
					t.Fatalf("write failed: %v", err)
				}
				
				dir := filepath.Dir(path)
				// Save original permissions
				info, err := os.Stat(dir)
				if err != nil {
					t.Fatalf("stat failed: %v", err)
				}
				originalPerm := info.Mode()

				// Make directory read-only
				if err := os.Chmod(dir, 0o500); err != nil {
					t.Fatalf("chmod failed: %v", err)
				}

				// Restore permissions after test
				t.Cleanup(func() {
					_ = os.Chmod(dir, originalPerm)
				})
			},
			configPath: func(p string) string { return p },
			createConfig: func(_ string) func() error {
				return func() error { return nil }
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			basePath := filepath.Join(tmpDir, "config.yaml")
			configPath := tt.configPath(basePath)
			mesheryFolder := filepath.Join(tmpDir, ".meshery")
			
			if tt.mesheryFolder != nil {
				mesheryFolder = tt.mesheryFolder(mesheryFolder)
			}

			tt.setup(configPath, t)

			err := config.InitDefaultConfig(
				configPath,
				mesheryFolder,
				utils.TemplateToken,
				utils.TemplateContext,
				tt.createConfig(configPath),
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tt.wantErr)
			}

			// Strong assertion for success case
			if !tt.wantErr {
				if _, err := os.Stat(configPath); err != nil {
					t.Fatalf("expected config file to exist: %v", err)
				}
			}
		})
	}
}
