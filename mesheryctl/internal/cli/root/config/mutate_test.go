package config_test

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

func TestNeedsMutation(t *testing.T) {
	tests := []struct {
		name  string
		setup func(string, *testing.T)
		want  bool
	}{
		{
			name:  "Given missing config file, When NeedsMutation is called, Then it returns true",
			setup: func(_ string, _ *testing.T) {},
			want:  true,
		},
		{
			name: "Given empty config file, When NeedsMutation is called, Then it returns true",
			setup: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			want: true,
		},
		{
			name: "Given non-empty config file, When NeedsMutation is called, Then it returns false",
			setup: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")

			tt.setup(configPath, t)

			got, err := config.NeedsMutation(configPath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("NeedsMutation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNeedsMutation_ErrorCase(t *testing.T) {
	t.Run("Given invalid config path, When NeedsMutation is called, Then it returns an error", func(t *testing.T) {
		invalidPath := string([]byte{0})

		_, err := config.NeedsMutation(invalidPath)
		if err == nil {
			t.Fatal("expected error for invalid path, got nil")
		}
	})
}

func TestInitDefaultConfig(t *testing.T) {
	t.Run("Given no existing config, When InitDefaultConfig is called, Then config file is created", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")
		mesheryFolder := filepath.Join(tmpDir, ".meshery")

		minimalMeshConfig := `contexts: {}
current-context: ""
tokens: []
`

		createConfig := func() error {
			return os.WriteFile(configPath, []byte(minimalMeshConfig), 0o644)
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

		data, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("config not found: %v", err)
		}

		if len(data) == 0 {
			t.Fatal("expected config to be created")
		}
	})
}

func TestConfigMutation(t *testing.T) {
	minimalMeshConfig := `contexts: {}
current-context: ""
tokens: []
`

	tests := []struct {
		name  string
		given func(string, *testing.T)
		want  bool
	}{
		{
			name:  "Given missing config file, When mutation flow runs, Then config is created",
			given: func(_ string, _ *testing.T) {},
			want:  true,
		},
		{
			name: "Given empty config file, When mutation flow runs, Then config is created",
			given: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte(""), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			want: true,
		},
		{
			name: "Given existing config file, When mutation flow runs, Then config is NOT modified",
			given: func(path string, t *testing.T) {
				if err := os.WriteFile(path, []byte("already-exists"), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			viper.Reset()
			t.Cleanup(viper.Reset)

			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			mesheryFolder := filepath.Join(tmpDir, ".meshery")

			tt.given(configPath, t)

			createConfig := func() error {
				return os.WriteFile(configPath, []byte(minimalMeshConfig), 0o644)
			}

			got, err := config.NeedsMutation(configPath)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tt.want {
				t.Fatalf("NeedsMutation() = %v, want %v", got, tt.want)
			}

			if got {
				if err := config.InitDefaultConfig(
					configPath,
					mesheryFolder,
					utils.TemplateToken,
					utils.TemplateContext,
					createConfig,
				); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			data, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("config not found: %v", err)
			}

			if tt.want {
				if len(data) == 0 {
					t.Fatal("expected config to be created")
				}
			} else {
				if string(data) != "already-exists" {
					t.Fatal("existing config was modified")
				}
			}
		})
	}
}
