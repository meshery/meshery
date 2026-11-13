package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/spf13/viper"
)

// TestDownloadDockerComposeFileUsesGivenContextVersion verifies that the
// compose file is fetched for the version pinned in the context passed in,
// such as the one chosen with --context, and not for the version of the
// current context in the meshconfig.
func TestDownloadDockerComposeFileUsesGivenContextVersion(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("current-context", "current")
	viper.Set("contexts", map[string]interface{}{
		"current": map[string]interface{}{
			"channel":  "stable",
			"version":  "v0.0.1",
			"platform": "docker",
		},
	})

	composeFile := DockerComposeFile
	DockerComposeFile = filepath.Join(t.TempDir(), "docker-compose.yaml")
	t.Cleanup(func() { DockerComposeFile = composeFile })

	releaseTag := ReleaseTag
	t.Cleanup(func() { ReleaseTag = releaseTag })

	StartMockery(t)
	t.Cleanup(func() { StopMockery(t) })
	url := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/v0.8.0/install/docker/docker-compose.yaml"
	httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "services: {}\n"))

	ctx := &config.Context{Channel: "stable", Version: "v0.8.0", Platform: "docker"}
	if err := DownloadDockerComposeFile(ctx, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ReleaseTag != "v0.8.0" {
		t.Errorf("expected release tag v0.8.0, got %s", ReleaseTag)
	}
	if got := httpmock.GetCallCountInfo()["GET "+url]; got != 1 {
		t.Errorf("expected one download from %s, got %d", url, got)
	}
	data, err := os.ReadFile(DockerComposeFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "services: {}\n" {
		t.Errorf("unexpected compose file contents %q", data)
	}
}
