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
	log "github.com/sirupsen/logrus"

	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
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

// FetchManifests is a wrapper function that identifies the required manifest files as downloads them
func FetchManifests(version string) ([]Manifest, error) {
	log.Debug("fetching required Kubernetes manifest files...")
	// get correct minfestsURL based on version
	manifestsURL, err := GetManifestTreeURL(version)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make GET request")
	}
	// pick all the manifest files stored in minfestsURL
	manifests, err := ListManifests(manifestsURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make GET request")
	}

	log.Debug("deleting ~/.meshery/manifests folder...")
	// delete manifests folder if it already exists
	if err := os.RemoveAll(ManifestsFolder); err != nil {
		return nil, err
	}
	log.Info("creating ~/.meshery/manifests folder...")
	// create a manifests folder under ~/.meshery to store the manifest files
	if err := os.MkdirAll(filepath.Join(MesheryFolder, ManifestsFolder), os.ModePerm); err != nil {
		return nil, errors.Wrapf(err, SystemError(fmt.Sprintf("failed to make %s directory", ManifestsFolder)))
	}
	log.Debug("created manifests folder...")

	gitHubFolder := "https://github.com/layer5io/meshery/tree/" + version + "/install/deployment_yamls/k8s"
	log.Info("downloading manifest files from ", gitHubFolder)

	// download all the manifest files to the ~/.meshery/manifests folder
	rawManifestsURL := "https://raw.githubusercontent.com/layer5io/meshery/" + version + "/install/deployment_yamls/k8s/"
	err = DownloadManifests(manifests, rawManifestsURL)

	if err != nil {
		return nil, errors.Wrap(err, "failed to download manifests")
	}

	return manifests, nil
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

// ApplyManifest is a wrapper function for client.ApplyManifest
func ApplyManifest(manifest []byte, client *meshkitkube.Client, update bool, delete bool) error {
	// ApplyManifest applies/updates/deletes the given manifest file to/from the Kubernetes cluster
	err := client.ApplyManifest(manifest, meshkitkube.ApplyOptions{
		Namespace: MesheryNamespace,
		Update:    update,
		Delete:    delete,
	})

	if err != nil {
		return errors.Wrap(err, "failed to apply manifests")
	}
	return nil
}

// ApplyManifestFiles applies/updates/deletes all the required manifests into the Kubernetes cluster
func ApplyManifestFiles(manifestArr []Manifest, requestedAdapters []string, client *meshkitkube.Client, update bool, delete bool) error {
	// path to the manifest files ~/.meshery/manifests
	manifestFiles := filepath.Join(MesheryFolder, ManifestsFolder)

	// read the manifest files as strings
	// other than the adapters, meshery-deployment.yaml, meshery-service.yaml and service-account.yaml should be applied
	MesheryDeploymentManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, MesheryDeployment))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	mesheryServiceManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, MesheryService))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}
	serviceAccountManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, ServiceAccount))
	if err != nil {
		return errors.Wrap(err, "failed to read manifest files")
	}

	// apply/update/delete the manifest files
	if err = ApplyManifest([]byte(MesheryDeploymentManifest), client, update, delete); err != nil {
		return err
	}
	if err = ApplyManifest([]byte(mesheryServiceManifest), client, update, delete); err != nil {
		return err
	}
	if err = ApplyManifest([]byte(serviceAccountManifest), client, update, delete); err != nil {
		return err
	}

	// loop through the required adapters as specified in the config.yaml file and apply/update/delete each
	for _, adapter := range requestedAdapters {
		// for each adapter, there is a meshery-adapterName-deployment.yaml and meshery-adapterName-service.yaml
		// manifest file. See- https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		adapterFile := filepath.Join(manifestFiles, adapter)
		adapterDeployment := adapterFile + "-deployment.yaml"
		adapterService := adapterFile + "-service.yaml"

		if !IsAdapterValid(manifestArr, adapter+"-deployment.yaml") {
			return fmt.Errorf("invalid adapter %s specified. Please check %s/config.yaml file", adapter, MesheryFolder)
		}

		// read manifest files as strings and apply/update/delete
		manifestDepl, err := meshkitutils.ReadLocalFile(adapterDeployment)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}
		manifestService, err := meshkitutils.ReadLocalFile(adapterService)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}

		if err = ApplyManifest([]byte(manifestDepl), client, update, delete); err != nil {
			return err
		}
		if err = ApplyManifest([]byte(manifestService), client, update, delete); err != nil {
			return err
		}
	}
	log.Debug("applied manifests to the Kubernetes cluster.")

	return nil
}
