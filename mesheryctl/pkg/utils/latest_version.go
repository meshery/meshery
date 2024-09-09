package utils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GetLatestVersionForMesheryctl() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://docs.meshery.io/project/releases/latest", nil)

	if err != nil {
		return "", err
	}

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func CheckMesheryctlClientVersion(build string) {
	Log.Info("Checking for latest version of mesheryctl...\n")

	// Inform user of the latest release version
	latestVersion, err := GetLatestVersionForMesheryctl()
	if err != nil {
		Log.Warn(fmt.Errorf("unable to check for latest version of mesheryctl. %s", err))
		return
	}
	if latestVersion == "" {
		Log.Warn(fmt.Errorf("unable to check for latest version of mesheryctl. %s", fmt.Errorf("no version found")))
		return
	}
	// If user is running an outdated release, let them know.
	if latestVersion != build {
		Log.Infof("A new release of mesheryctl is available: %s → %s", build, latestVersion)
		Log.Info("https://docs.meshery.io/project/releases/latest")
		Log.Info("Check https://docs.meshery.io/installation/upgrades#upgrading-meshery-cli for instructions on how to update mesheryctl\n")
	} else { // If user is running the latest release, let them know.
		Log.Info(latestVersion, " is the latest release.")
	}
}
