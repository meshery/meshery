package root

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func TestInitConfigUseCases(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, tmp string)
		assert    func(t *testing.T)
		skipOnWin bool
	}{
		{
			name: "given missing config file when calling initConfig then default config is created",
			setup: func(t *testing.T, tmp string) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-missing")
				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
			},
			assert: func(t *testing.T) {
				// initConfig should create the missing config file with default content
				info, err := os.Stat(utils.DefaultConfigPath)
				if os.IsNotExist(err) {
					t.Errorf("expected config file to be created at %s, but it does not exist", utils.DefaultConfigPath)
					return
				}
				if err != nil {
					t.Errorf("unexpected error checking config file: %v", err)
					return
				}
				if info.Size() == 0 {
					t.Errorf("expected config file at %s to be non-empty after default creation, but it is empty", utils.DefaultConfigPath)
				}
			},
		},
		{
			name: "given an empty config file when calling initConfig then empty config creates default",
			setup: func(t *testing.T, tmp string) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-empty")
				if err := os.MkdirAll(utils.MesheryFolder, 0o755); err != nil {
					t.Fatal(err)
				}

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				if err := os.WriteFile(utils.DefaultConfigPath, []byte(""), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			assert: func(t *testing.T) {
				// initConfig should repopulate the empty config file with default content
				info, err := os.Stat(utils.DefaultConfigPath)
				if os.IsNotExist(err) {
					t.Errorf("expected config file to exist at %s after repopulation, but it does not", utils.DefaultConfigPath)
					return
				}
				if err != nil {
					t.Errorf("unexpected error checking config file: %v", err)
					return
				}
				if info.Size() == 0 {
					t.Errorf("expected config file at %s to be non-empty after repopulation, but it is still empty", utils.DefaultConfigPath)
				}
			},
		},
		{
			name: "given config path without permission when calling initconfig then returns error",
			setup: func(t *testing.T, tmp string) {
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
			assert: func(t *testing.T) {
				// initConfig should return early on permission error without removing the file
				if _, err := os.Stat(utils.DefaultConfigPath); os.IsNotExist(err) {
					t.Errorf("expected config file to still exist at %s after permission-denied early return, but it was removed", utils.DefaultConfigPath)
				}
			},
			skipOnWin: true, // chmod 0o000 has no effect on Windows
		},
		{
			name: "given permission denied when calling initconfig then returns error",
			setup: func(t *testing.T, tmp string) {
				utils.MesheryFolder = filepath.Join(tmp, "meshery-existing")
				if err := os.MkdirAll(utils.MesheryFolder, 0o755); err != nil {
					t.Fatal(err)
				}

				utils.DefaultConfigPath = filepath.Join(utils.MesheryFolder, "config.yaml")
				if err := os.WriteFile(utils.DefaultConfigPath, []byte("test"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			assert: func(t *testing.T) {
				// initConfig should leave a pre-existing non-empty config file in place
				if _, err := os.Stat(utils.DefaultConfigPath); os.IsNotExist(err) {
					t.Errorf("expected pre-existing config file to remain at %s after initConfig, but it was removed", utils.DefaultConfigPath)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipOnWin && isWindows() {
				t.Skip("skipping permission test on Windows: chmod 0o000 has no effect")
			}

			// Each subtest owns its TempDir so viper.Reset() (registered below,
			// LIFO) releases file handles before this TempDir is removed.
			tmp := t.TempDir()

			// Save global state
			origMesheryFolder := utils.MesheryFolder
			origDefaultConfigPath := utils.DefaultConfigPath
			origCfgFile := cfgFile
			origLog := utils.Log

			t.Cleanup(func() {
				utils.MesheryFolder = origMesheryFolder
				utils.DefaultConfigPath = origDefaultConfigPath
				cfgFile = origCfgFile
				utils.Log = origLog
			})

			utils.SetupMeshkitLoggerTesting(t, false)

			tt.setup(t, tmp)

			// ensure default config path branch is exercised
			cfgFile = utils.DefaultConfigPath

			initConfig()
			// Reset viper immediately to release any file handles it holds
			// (required on Windows where open handles block TempDir cleanup).
			viper.Reset()

			tt.assert(t)
		})
	}
}
