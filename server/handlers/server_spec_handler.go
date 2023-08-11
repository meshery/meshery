package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/pkg/constants"
	"github.com/layer5io/meshkit/utils"
	"github.com/spf13/viper"
)

// Version defines the Json payload structure for version api\
type Version struct {
	Build          string `json:"build,omitempty"`
	Latest         string `json:"latest,omitempty"`
	Outdated       *bool  `json:"outdated,omitempty"`
	CommitSHA      string `json:"commitsha,omitempty"`
	ReleaseChannel string `json:"release_channel,omitempty"`
}

// swagger:route GET /api/system/version SystemAPI idGetSystemVersion
// Handle GET request for system/server version
//
// Returns the running Meshery version
// responses:
// 	200: mesheryVersionRespWrapper

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, _ *http.Request) {
	// Default values incase any errors
	version := &Version{
		Build:          viper.GetString("BUILD"),
		CommitSHA:      viper.GetString("COMMITSHA"),
		ReleaseChannel: viper.GetString("RELEASE_CHANNEL"),
	}

	// if r.Method != http.MethodGet {
	//      w.WriteHeader(http.StatusNotFound)
	//      return
	// }

	// compare the server build with the target build
	isOutdated, latestVersion, err := CheckLatestVersion(version.Build)
	if err != nil {
		h.log.Error(err)
	} else {
		// Add "Latest" and "Outdated" fields to the response
		version.Latest = latestVersion
		version.Outdated = isOutdated
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(version)
	if err != nil {
		h.log.Error(ErrEncoding(err, "server-version"))
		http.Error(w, ErrEncoding(err, "server-version").Error(), http.StatusNotFound)
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
	if latestVersion != serverVersion {
		isOutdated = true
		return &isOutdated, latestVersion, nil
	}

	return &isOutdated, latestVersion, nil
}
