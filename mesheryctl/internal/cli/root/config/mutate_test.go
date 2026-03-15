package config_test

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

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
			name: "Given missing config file, when mutation is needed, then default config is created",
			given: func(_ string, _ *testing.T) {
				// file intentionally missing
			},
			want: true,
		},
		{
			name: "Given empty config file, when mutation is needed, then default config is created",
			given: func(configPath string, t *testing.T) {
				if err := os.WriteFile(configPath, []byte(""), 0o644); err != nil {
					t.Fatalf("failed to write file: %v", err)
				}
			},
			want: true,
		},
		{
			name: "Given existing config file, when mutation is not needed, then config is NOT modified",
			given: func(configPath string, t *testing.T) {
				if err := os.WriteFile(configPath, []byte("already-exists"), 0o644); err != nil {
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

			// GIVEN
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			mesheryFolder := filepath.Join(tmpDir, ".meshery")

			tt.given(configPath, t)

			createConfig := func() error {
				return os.WriteFile(configPath, []byte(minimalMeshConfig), 0o644)
			}

			// WHEN
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

			// THEN
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
