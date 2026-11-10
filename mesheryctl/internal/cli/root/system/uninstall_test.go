package system

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mesheryerrors "github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	// Docker cleanup already forces container removal, so --force only changes
	// Kubernetes teardown. The help output must say so.
	assert.Contains(t, output, "has no effect for Docker installations")
}

// TestUninstallDockerErrorRemediation verifies that a Docker Compose removal
// failure points at Docker remediation and never at the Kubernetes-only
// --force flag.
func TestUninstallDockerErrorRemediation(t *testing.T) {
	err := utils.ErrDockerComposeRemove(errors.New("docker compose remove failed"))

	var meshkitErr *mesheryerrors.Error
	require.ErrorAs(t, err, &meshkitErr)
	assert.Equal(t, utils.ErrDockerComposeRemoveCode, meshkitErr.Code)
	assert.Equal(t, "Failed to delete Meshery containers", meshkitErr.ShortDescription[0])

	remediation := strings.Join(meshkitErr.SuggestedRemediation, " ")
	assert.Contains(t, remediation, "Docker")
	assert.NotContains(t, remediation, "--force")
}

// TestUninstallKubernetesErrorRemediation verifies that Kubernetes teardown
// failures keep the --force remediation, because --force is the flag that
// changes the Kubernetes uninstall path.
func TestUninstallKubernetesErrorRemediation(t *testing.T) {
	err := ErrUninstallMeshery(errors.New("helm uninstall failed"))

	var meshkitErr *mesheryerrors.Error
	require.ErrorAs(t, err, &meshkitErr)
	assert.Equal(t, ErrUninstallMesheryCode, meshkitErr.Code)

	remediation := strings.Join(meshkitErr.SuggestedRemediation, " ")
	assert.Contains(t, remediation, "mesheryctl system uninstall --force")
	assert.Contains(t, remediation, "cluster")
	assert.NotContains(t, remediation, "Docker")
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
