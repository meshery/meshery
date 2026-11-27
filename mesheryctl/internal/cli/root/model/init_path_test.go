package model

import (
	"os"
	"path/filepath"
	"testing"

	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

// newModelInitCommandForTest builds a throwaway `mesheryctl model` command with
// the init subcommand attached, so flag state does not leak between tests.
func newModelInitCommandForTest() *cobra.Command {
	cmd := &cobra.Command{
		Use:     ModelCmd.Use,
		Short:   ModelCmd.Short,
		Long:    ModelCmd.Long,
		Example: ModelCmd.Example,
		Args:    ModelCmd.Args,
		RunE:    ModelCmd.RunE,
	}
	ModelCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		cmd.Flags().AddFlag(flag)
	})
	cmd.AddCommand(initModelCmd)
	return cmd
}

// TestModelInitRejectsNonDirectoryPath is a regression test for a nil pointer
// dereference: `os.Stat` on a path below a regular file returns ENOTDIR, which
// os.IsNotExist does not classify as "not exist". The command then used the
// nil FileInfo and panicked instead of reporting the invalid --path.
func TestModelInitRejectsNonDirectoryPath(t *testing.T) {
	utils.SetupContextEnv(t)
	mesheryctlflags.InitValidators(initModelCmd)

	dir := t.TempDir()
	regularFile := filepath.Join(dir, "not-a-directory")
	require.NoError(t, os.WriteFile(regularFile, []byte("regular file"), 0o644))

	tests := map[string]string{
		"--path is a regular file":                regularFile,
		"--path is nested below a regular file":   filepath.Join(regularFile, "nested"),
		"--path is nested deeper below a regular": filepath.Join(regularFile, "nested", "deeper"),
	}

	for name, path := range tests {
		t.Run(name, func(t *testing.T) {
			cmd := newModelInitCommandForTest()
			defer utils.ResetCommandFlags(cmd, t)
			defer utils.ResetCommandFlags(initModelCmd, t)
			cmd.SetOut(utils.SetupMeshkitLoggerTesting(t, false))
			cmd.SetArgs([]string{"init", "test-case-not-a-directory", "--path", path})

			var err error
			require.NotPanics(t, func() {
				err = cmd.Execute()
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), "cannot use")
		})
	}
}
