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
// compose file is fetched for the version of the context passed in, such as
// the one chosen with --context, and not for the version of the current
// context in the meshconfig. An empty version means the latest release.
func TestDownloadDockerComposeFileUsesGivenContextVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantTag string
	}{
		{name: "pinned version", version: "v0.8.0", wantTag: "v0.8.0"},
		{name: "empty version means latest", version: "", wantTag: "v0.8.5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			repo := constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo()
			httpmock.RegisterResponder("GET", "https://github.com/"+repo+"/releases",
				httpmock.NewStringResponder(200, `<a href="/`+repo+`/releases/tag/v0.8.5">v0.8.5</a>`))
			url := "https://raw.githubusercontent.com/" + repo + "/" + tt.wantTag + "/install/docker/docker-compose.yaml"
			httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, "services: {}\n"))

			ctx := &config.Context{Channel: "stable", Version: tt.version, Platform: "docker"}
			if err := DownloadDockerComposeFile(ctx, true); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ReleaseTag != tt.wantTag {
				t.Errorf("expected release tag %s, got %s", tt.wantTag, ReleaseTag)
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
		})
	}
}
