package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-latest"
)

// Version defines the Json payload structure for version api\
type Version struct {
	Build     string `json:"build,omitempty"`
	Latest    string `json:"latest,omitempty"`
	Outdated  *bool  `json:"outdated,omitempty"`
	CommitSHA string `json:"commitsha,omitempty"`
}

// ServerVersionHandler handles the version api request for the server
func (h *Handler) ServerVersionHandler(w http.ResponseWriter, r *http.Request) {
	// Default values incase any errors
	version := &Version{
		Build:     viper.GetString("BUILD"),
		CommitSHA: viper.GetString("COMMITSHA"),
	}

	// if r.Method != http.MethodGet {
	//      w.WriteHeader(http.StatusNotFound)
	//      return
	// }

	// compare the server build with the target build
	res, err := CheckLatestVersion(version.Build)
	if err != nil {
		logrus.Errorln(err)
	} else {
		// Add "Latest" and "Outdated" fields to the response
		version.Latest = res.Current
		version.Outdated = &res.Outdated
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(version)
	if err != nil {
		logrus.Errorf("unable to send data: %v", err)
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
		return nil, errors.Wrap(err, "failed to compare latest and current version of Meshery")
	}
	// If user is running an outdated release, let them know.
	if res.Outdated {
		logrus.Info("\n  ", serverVersion, " is not the latest Meshery release. Update to v", res.Current, ". Run `mesheryctl system update`")
		promptLabel := fmt.Sprintf("Would you like to upgrade to v%s now [y/n]?", res.Current)

		prompt := promptui.Select{
			Label: promptLabel,
			Items: []string{"y", "n"},
		}
		_, result, err := prompt.Run()
		if err != nil {
			logrus.Error("Prompt failed %w\n", err)
		}
		result = strings.TrimSpace(result)
		result = strings.ToLower(result)
		if result != "n" && result != "no" {
			err = utils.UpdateMesheryContainers()
			if err != nil {
				logrus.Error("Unable to update meshery: %w", err)
			}
		}
	}

	// If user is running the latest release, let them know.
	if res.Latest {
		logrus.Info("\n  ", serverVersion, " is the latest Meshery release.")
	}

	// Add "v" to the "Current" property of the CheckResponse
	res.Current = fmt.Sprintf("v%s", res.Current)

	return res, nil
}
