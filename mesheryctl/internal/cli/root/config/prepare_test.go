package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/spf13/viper"
)

func TestPrepareConfig_ReadOnlyBehavior(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name        string
		setup       func(t *testing.T) string
		expectError bool
	}{
		{
			name: "missing config file is allowed",
			setup: func(t *testing.T) string {
				return filepath.Join(tmp, "missing", "config.yaml")
			},
			expectError: false,
		},
		{
			name: "empty config file is allowed",
			setup: func(t *testing.T) string {
				dir := filepath.Join(tmp, "empty")
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}

				cfg := filepath.Join(dir, "config.yaml")
				if err := os.WriteFile(cfg, []byte(""), 0o644); err != nil {
					t.Fatal(err)
				}
				return cfg
			},
			expectError: false,
		},
		{
			name: "permission denied config returns error",
			setup: func(t *testing.T) string {
				dir := filepath.Join(tmp, "permission")
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}

				cfg := filepath.Join(dir, "config.yaml")
				if err := os.WriteFile(cfg, []byte("abc"), 0o644); err != nil {
					t.Fatal(err)
				}
				if err := os.Chmod(cfg, 0o000); err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() {
					_ = os.Chmod(cfg, 0o644)
				})

				return cfg
			},
			expectError: true,
		},
		{
			name: "existing valid config loads successfully",
			setup: func(t *testing.T) string {
				dir := filepath.Join(tmp, "existing")
				if err := os.MkdirAll(dir, 0o755); err != nil {
					t.Fatal(err)
				}

				cfg := filepath.Join(dir, "config.yaml")

				content := "current-context: local\n" +
					"contexts:\n" +
					"  local:\n" +
					"    endpoint: http://localhost:9081\n" +
					"    platform: docker\n"
				if err := os.WriteFile(cfg, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
				return cfg
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			cfgPath := tt.setup(t)
			err := config.PrepareConfig(cfgPath)

			if tt.expectError && err == nil {
				t.Fatalf("expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
