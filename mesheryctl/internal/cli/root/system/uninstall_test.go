package system

import (
	"slices"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// TestUninstallCmdRegistration verifies that the uninstall command is wired
// into `mesheryctl system` and exposes the flags that drive the teardown
// behavior.
func TestUninstallCmdRegistration(t *testing.T) {
	assert.Contains(t, availableSubcommands, uninstallCmd)
	assert.Equal(t, "uninstall", uninstallCmd.Use)

	// --force drives the ForceCleanupCluster path and the tolerance of the
	// shared deleteCR/deleteCRD helpers toward missing resources.
	assert.NotNil(t, uninstallCmd.Flags().Lookup("force"))

	// --keep-namespace skips the Meshery namespace deletion.
	assert.NotNil(t, uninstallCmd.Flags().Lookup("keep-namespace"))
}

// TestUninstallCmdHelp verifies the help output of the uninstall command. Help
// short-circuits before the preflight health checks, so this test runs without
// a Docker daemon or a Kubernetes cluster.
func TestUninstallCmdHelp(t *testing.T) {
	buf := setupSystemOutCmdTest(t)
	defer buf.Reset()

	SystemCmd.SetArgs([]string{"uninstall", "--help"})
	err = SystemCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Uninstall Meshery")
	assert.Contains(t, output, "--keep-namespace")
	assert.Contains(t, output, "--force")
}

// TestSystemUninstallErrorReference verifies the "uninstall" error formatting
// points at the uninstall command reference page.
func TestSystemUninstallErrorReference(t *testing.T) {
	msg := utils.SystemLifeCycleError("this command takes no arguments.\n", "uninstall")

	assert.Contains(t, msg, "this command takes no arguments")
	assert.Contains(t, msg, "https://docs.meshery.io/reference/references/mesheryctl/system/uninstall")
}

// TestUninstallCmdInSubcommandList guards the subcommand registration against
// accidental removal that would silently change `mesheryctl system --help`.
func TestUninstallCmdInSubcommandList(t *testing.T) {
	names := make([]string, 0, len(availableSubcommands))
	for _, cmd := range availableSubcommands {
		names = append(names, cmd.Use)
	}

	assert.True(t, slices.Contains(names, "delete"))
	assert.True(t, slices.Contains(names, "uninstall"))
}
