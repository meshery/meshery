package root

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

func TestPrepareConfig_ReadOnlyBehavior(t *testing.T) {
	tmp := t.TempDir()

	tests := []struct {
		name        string
		setup       func(t *testing.T)
		expectError bool
	}{
		{
			name: "given missing config file when prepareConfig then error returned",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-missing")
				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
			},
			expectError: true,
		},
		{
			name: "given empty config file when prepareConfig then error returned",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-empty")
				_ = os.MkdirAll(utils.MesheryFolder, 0o755)

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				_ = os.WriteFile(utils.DefaultConfigPath, []byte(""), 0o644)
			},
			expectError: false,
		},
		{
			name: "given permission denied config when prepareConfig then error returned",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-permission")
				_ = os.MkdirAll(utils.MesheryFolder, 0o755)

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				_ = os.WriteFile(utils.DefaultConfigPath, []byte("abc"), 0o644)
				_ = os.Chmod(utils.DefaultConfigPath, 0o000)

				t.Cleanup(func() {
					_ = os.Chmod(utils.DefaultConfigPath, 0o644)
				})
			},
			expectError: true,
		},
		{
			name: "given existing config when prepareConfig then loads successfully without mutation",
			setup: func(t *testing.T) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-existing")
				_ = os.MkdirAll(utils.MesheryFolder, 0o755)

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				_ = os.WriteFile(utils.DefaultConfigPath, []byte("test"), 0o644)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()

			origMesheryFolder := utils.MesheryFolder
			origDefaultConfigPath := utils.DefaultConfigPath
			origCfgFile := cfgFile

			t.Cleanup(func() {
				utils.MesheryFolder = origMesheryFolder
				utils.DefaultConfigPath = origDefaultConfigPath
				cfgFile = origCfgFile
			})

			tt.setup(t)
			cfgFile = utils.DefaultConfigPath

			err := prepareConfig()

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
