package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"
)

// Version defines the Json payload structure for version api\
type Version struct {
	Build          string `json:"build,omitempty"`
	Latest         string `json:"latest,omitempty"`
	Outdated       *bool  `json:"outdated,omitempty"`
	CommitSHA      string `json:"commitsha,omitempty"`
	ReleaseChannel string `json:"release_channel,omitempty"`
}

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, r *http.Request) {
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
	res, err := CheckLatestVersion(version.Build)
	if err != nil {
		h.log.Error(err)
	} else {
		// Add "Latest" and "Outdated" fields to the response
		version.Latest = res.Current
		version.Outdated = &res.Outdated
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(version)
	if err != nil {
		h.log.Error(ErrEncoding(err, "server-version"))
		http.Error(w, ErrEncoding(err, "server-version").Error(), http.StatusNotFound)
	}
}

// CheckLatestVersion takes in the current server version compares it with the target
// and returns the result (latest.CheckResponse)
func CheckLatestVersion(serverVersion string) (*latest.CheckResponse, error) {
	githubTag := &latest.GithubTag{
		Owner:      mesheryGitHubOrg,
		Repository: mesheryGitHubRepo,
	}

	// Compare current running Meshery server version to the latest available Meshery release on GitHub.
	res, err := latest.Check(githubTag, serverVersion)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compare latest and current version")
	}

	// Add "v" to the "Current" property of the CheckResponse
	res.Current = fmt.Sprintf("v%s", res.Current)

	return res, nil
}
