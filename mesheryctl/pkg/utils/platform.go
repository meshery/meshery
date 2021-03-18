package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
)

var (
	// ManifestsFolder is where the Kubernetes manifests are stored
	ManifestsFolder = "manifests"
)

// GetManifestTreeURL returns the manifest tree url based on version
func GetManifestTreeURL(version string) (string, error) {
	url := "https://api.github.com/repos/layer5io/meshery/git/trees/" + version + "?recursive=1"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer SafeClose(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}

	var manLis ManifestList

	err = json.Unmarshal([]byte(body), &manLis)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	for i := range manLis.Tree {
		if manLis.Tree[i].Path == "install/deployment_yamls/k8s" {
			return manLis.Tree[i].URL, nil
		}
	}
	return "", errors.New("could not find path: install/deployment_yamls/k8s in the manifest tree")
}

// ListManifests lists the manifest files stored in GitHub
func ListManifests(url string) ([]Manifest, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer SafeClose(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var manLis ManifestList

	err = json.Unmarshal([]byte(body), &manLis)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	return manLis.Tree, nil
}

// GetManifestURL returns the URLs for the manifest files
func GetManifestURL(manifest Manifest, rawManifestsURL string) string {
	var manifestURL string

	if manifest.Typ == "blob" {
		manifestURL = rawManifestsURL + manifest.Path
		return manifestURL
	}
	return ""
}

// DownloadManifests downloads all the Kubernetes manifest files
func DownloadManifests(manifestArr []Manifest, rawManifestsURL string) error {
	for _, manifest := range manifestArr {
		if manifestFile := GetManifestURL(manifest, rawManifestsURL); manifestFile != "" {
			// download the manifest files to ~/.meshery/manifests folder
			filepath := filepath.Join(MesheryFolder, ManifestsFolder, manifest.Path)
			if err := DownloadFile(filepath, manifestFile); err != nil {
				return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s", filepath, manifestFile)))
			}
		}
	}
	return nil
}

// GetLatestStableReleaseTag fetches and returns the latest release tag from GitHub
func GetLatestStableReleaseTag() (string, error) {
	url := "https://api.github.com/repos/layer5io/meshery/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer SafeClose(resp.Body)

	var dat map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal json into object")
	}

	return dat["tag_name"].(string), nil
}

// IsAdapterValid checks if the adapter mentioned by the user is a valid adapter
func IsAdapterValid(manifestArr []Manifest, adapterManifest string) bool {
	for _, v := range manifestArr {
		if v.Path == adapterManifest {
			return true
		}
	}

	return false
}

// DownloadDockerComposeFile fetches docker-compose.yaml based on passed context if it does not exists.
// Use force to override download anyway
func DownloadDockerComposeFile(ctx config.Context, force bool) error {
	if _, err := os.Stat(DockerComposeFile); os.IsNotExist(err) || force {
		fileURL := ""

		if ctx.Channel == "edge" {
			fileURL = "https://raw.githubusercontent.com/layer5io/meshery/master/docker-compose.yaml"
		} else if ctx.Channel == "stable" {
			if ctx.Version == "latest" {
				ctx.Version, err = GetLatestStableReleaseTag()
				if err != nil {
					return errors.Wrapf(err, fmt.Sprintf("failed to fetch latest stable release tag"))
				}
			}
			fileURL = "https://raw.githubusercontent.com/layer5io/meshery/" + ctx.Version + "/docker-compose.yaml"
		} else {
			return errors.Errorf("unknown channel %s", ctx.Channel)
		}

		if err := DownloadFile(DockerComposeFile, fileURL); err != nil {
			return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s", DockerComposeFile, fileURL)))
		}
	}
	return nil
}
