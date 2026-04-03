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

	t.Run("Given valid setup, When InitDefaultConfig is called, Then config is created", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		mesheryFolder := filepath.Join(tmpDir, ".meshery")

		createConfig := func() error {
			return os.WriteFile(configPath, []byte("data"), 0o644)
		}

		err := config.InitDefaultConfig(
			configPath,
			mesheryFolder,
			utils.TemplateToken,
			utils.TemplateContext,
			createConfig,
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Given createConfig fails, When InitDefaultConfig is called, Then error is returned", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		mesheryFolder := filepath.Join(tmpDir, ".meshery")

		createConfig := func() error {
			return errors.New("create failed")
		}

		err := config.InitDefaultConfig(
			configPath,
			mesheryFolder,
			utils.TemplateToken,
			utils.TemplateContext,
			createConfig,
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("Given invalid directory, When MkdirAll fails, Then error is returned", func(t *testing.T) {
		configPath := "/invalid/path/config.yaml"
		mesheryFolder := "/invalid/path/.meshery"

		createConfig := func() error { return nil }

		err := config.InitDefaultConfig(
			configPath,
			mesheryFolder,
			utils.TemplateToken,
			utils.TemplateContext,
			createConfig,
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
