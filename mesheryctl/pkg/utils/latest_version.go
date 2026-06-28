package utils

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	latestVersionURL            = "https://docs.meshery.io/project/releases/latest"
	latestVersionRequestTimeout = 5 * time.Second
)

func GetLatestVersionForMesheryctl() (string, error) {
	req, err := http.NewRequest(http.MethodGet, latestVersionURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "mesheryctl/"+runtime.Version())

	client := &http.Client{
		Timeout: latestVersionRequestTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("request to %s failed with status %s", latestVersionURL, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func CheckMesheryctlClientVersion(build string) {
	Log.Info("Checking for latest version of mesheryctl...\n")

	latestVersion, err := GetLatestVersionForMesheryctl()
	if err != nil || latestVersion == "" {
		if err == nil {
			err = fmt.Errorf("no version found")
		}
		Log.Warn(fmt.Errorf("unable to check for latest version of mesheryctl. %s", err))
		return
	}
	if latestVersion != build {
		Log.Infof("A new release of mesheryctl is available: %s \u2192 %s", build, latestVersion)
		Log.Info(latestVersionURL)
		Log.Info("Check https://docs.meshery.io/installation/upgrades#upgrading-meshery-cli for instructions on how to update mesheryctl\n")
	} else {
		Log.Infof("%s is the latest release.", latestVersion)
	}
}