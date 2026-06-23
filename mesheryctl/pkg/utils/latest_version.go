package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const latestVersionURL = "https://docs.meshery.io/project/releases/latest"
const latestVersionHTTPTimeout = 5 * time.Second

func GetLatestVersionForMesheryctl() (string, error) {
	client := &http.Client{Timeout: latestVersionHTTPTimeout}

	return getLatestVersionForMesheryctl(latestVersionURL, client)
}

func getLatestVersionForMesheryctl(url string, client *http.Client) (string, error) {
	if client == nil {
		client = &http.Client{Timeout: latestVersionHTTPTimeout}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return "", ErrCreatingRequest(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		var nerr net.Error
		if (errors.As(err, &nerr) && nerr.Timeout()) || errors.Is(err, context.DeadlineExceeded) {
			return "", nil
		}
		return "", ErrFailRequest(err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", ErrFailReqStatus(resp.StatusCode, "unexpected response when checking mesheryctl version")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		var nerr net.Error
		if (errors.As(err, &nerr) && nerr.Timeout()) || errors.Is(err, context.DeadlineExceeded) {
			return "", nil
		}
		return "", ErrReadResponseBody(err)
	}

	return strings.TrimSpace(string(data)), nil
}

func CheckMesheryctlClientVersion(build string) {
	Log.Info("Checking for latest version of mesheryctl...\n")

	latestVersion, err := GetLatestVersionForMesheryctl()
	if err != nil {

		Log.Warn(fmt.Errorf("unable to check for latest version of mesheryctl. %s", err))
		latestVersion = build
	}
	if latestVersion == "" {

		Log.Debug("latest version lookup returned empty; assuming current build is latest")
		latestVersion = build
	}

	if latestVersion != build {
		Log.Infof("A new release of mesheryctl is available: %s → %s", build, latestVersion)
		Log.Info("https://docs.meshery.io/project/releases/latest")
		Log.Info("Check https://docs.meshery.io/installation/upgrades#upgrading-meshery-cli for instructions on how to update mesheryctl\n")
	} else {
		Log.Info(latestVersion, " is the latest release.")
	}
}
