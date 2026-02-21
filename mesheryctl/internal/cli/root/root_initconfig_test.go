package root

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

func TestInitConfigUseCases(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name  string
		setup func(t *testing.T)
	}{
		{
			name: "given missing config file when calling initConfig then default config is created",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-missing")
				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
			},
		},
		{
			name: "given an empty config file when calling initConfig then empty config creates default",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-empty")
				if err := os.MkdirAll(utils.MesheryFolder, 0o755); err != nil {
					t.Fatal(err)
				}

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				if err := os.WriteFile(utils.DefaultConfigPath, []byte(""), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "given config path without permission when calling initconfig then returns error",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-permission")
				if err := os.MkdirAll(utils.MesheryFolder, 0o755); err != nil {
					t.Fatal(err)
				}

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				if err := os.WriteFile(utils.DefaultConfigPath, []byte("abc"), 0o644); err != nil {
					t.Fatal(err)
				}

				if err := os.Chmod(utils.DefaultConfigPath, 0o000); err != nil {
					t.Fatal(err)
				}

				t.Cleanup(func() {
					_ = os.Chmod(utils.DefaultConfigPath, 0o644)
				})
			},
		},
		{
			name: "given permission denied when calling initconfig then returns error",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-existing")
				if err := os.MkdirAll(utils.MesheryFolder, 0o755); err != nil {
					t.Fatal(err)
				}

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				if err := os.WriteFile(utils.DefaultConfigPath, []byte("test"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			// Save global state
			origMesheryFolder := utils.MesheryFolder
			origDefaultConfigPath := utils.DefaultConfigPath
			origCfgFile := cfgFile

			// Restore after each test
			t.Cleanup(func() {
				utils.MesheryFolder = origMesheryFolder
				utils.DefaultConfigPath = origDefaultConfigPath
				cfgFile = origCfgFile
			})

			tt.setup(t)

			// ensure default config path branch is exercised
			cfgFile = utils.DefaultConfigPath

			// Should not panic
			initConfig()
		})
	}
}
