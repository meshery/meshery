package config_test

import (
	"os"
	"path/filepath"
	"testing"

	config "github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

const minimalMeshConfig = "contexts: {}\ncurrent-context: \"\"\ntokens: []\n"

func TestMutateConfigIfNeeded(t *testing.T) {

	tests := []struct {
		name        string
		given       func(string)
		thenCreated bool
	}{
		{
			name: "Given missing config file, when mutate is called, then default config is created",
			given: func(_ string) {
				// file intentionally missing
			},
			thenCreated: true,
		},
		{
			name: "Given empty config file, when mutate is called, then default config is created",
			given: func(configPath string) {
				_ = os.WriteFile(configPath, []byte(""), 0o644)
			},
			thenCreated: true,
		},
		{
			name: "Given existing config file, when mutate is called, then config is NOT modified",
			given: func(configPath string) {
				_ = os.WriteFile(configPath, []byte("already-exists"), 0o644)
			},
			thenCreated: false,
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

			tt.given(configPath)

			createConfig := func() error {
				return os.WriteFile(configPath, []byte(minimalMeshConfig), 0o644)
			}

			// WHEN
			err := config.MutateConfigIfNeeded(
				configPath,
				mesheryFolder,
				utils.TemplateToken,
				utils.TemplateContext,
				createConfig,
			)

			// THEN
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			data, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("config not found: %v", err)
			}

			if tt.thenCreated {
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
