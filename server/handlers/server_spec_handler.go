package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Masterminds/semver/v3"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/utils"
	system "github.com/meshery/schemas/models/v1beta1/system"
	"github.com/spf13/viper"
)

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, _ *http.Request) {
	build := viper.GetString("BUILD")
	commitSHA := viper.GetString("COMMITSHA")
	releaseChannel := viper.GetString("RELEASE_CHANNEL")

	version := &system.SystemVersion{}
	if build != "" {
		version.Build = &build
	}
	if commitSHA != "" {
		version.Commitsha = &commitSHA
	}
	if releaseChannel != "" {
		version.ReleaseChannel = &releaseChannel
	}

	// compare the server build with the target build
	isOutdated, latestVersion, err := CheckLatestVersion(build)
	if err != nil {
		h.log.Error(err)
	} else {
		// Add "Latest" and "Outdated" fields to the response
		if latestVersion != "" {
			version.Latest = &latestVersion
		}
		version.Outdated = isOutdated
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(version)
	if err != nil {
		h.log.Error(models.ErrEncoding(err, "server-version"))
		writeMeshkitError(w, models.ErrEncoding(err, "server-version"), http.StatusNotFound)
	}
}

// CheckLatestVersion takes in the current server version compares it with the target
// and returns the (isOutdated, latestVersion, error)
func CheckLatestVersion(serverVersion string) (*bool, string, error) {
	// Inform user of the latest release version
	latestVersions, err := utils.GetLatestReleaseTagsSorted(constants.GetMesheryGitHubOrg(), constants.GetMesheryGitHubRepo())
	isOutdated := false
	if err != nil {
		return nil, "", ErrGetLatestVersion(err)
	}
	if len(latestVersions) == 0 {
		return &isOutdated, "", fmt.Errorf("no versions found")
	}
	latestVersion := latestVersions[len(latestVersions)-1]
	// Compare current running Meshery server version to the latest available Meshery release on GitHub.
	isOutdated = isCoreVersionOutdated(serverVersion, latestVersion)
	return &isOutdated, latestVersion, nil
}

// isCoreVersionOutdated reports whether serverVersion's core vX.Y.Z is
// strictly older than latestVersion's. Prerelease or revision suffixes
// (for example the "-1" on a kanvas-channel build "v1.0.82-1") and build
// metadata are ignored, so a server at v1.0.82-1 is not outdated relative
// to GitHub release v1.0.82. A newer or equal core is not outdated. If
// either string is not a semver version, the result falls back to exact
// string inequality, which is the previous comparison.
func isCoreVersionOutdated(serverVersion, latestVersion string) bool {
	server, err := semver.NewVersion(serverVersion)
	if err != nil {
		return latestVersion != serverVersion
	}
	latest, err := semver.NewVersion(latestVersion)
	if err != nil {
		return latestVersion != serverVersion
	}
	serverCore := semver.New(server.Major(), server.Minor(), server.Patch(), "", "")
	latestCore := semver.New(latest.Major(), latest.Minor(), latest.Patch(), "", "")
	return serverCore.LessThan(latestCore)
}
