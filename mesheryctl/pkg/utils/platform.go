package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/constants"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gopkg.in/yaml.v2"

	v1core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	meshkitutils "github.com/layer5io/meshkit/utils"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
)

var (
	// ManifestsFolder is where the Kubernetes manifests are stored
	ManifestsFolder = "manifests"
	ReleaseTag      string
)

// GetManifestTreeURL returns the manifest tree url based on version
func GetManifestTreeURL(version string) (string, error) {
	url := "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/" + version + "?recursive=1"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer SafeClose(resp.Body)

	body, err := io.ReadAll(resp.Body)
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

	body, err := io.ReadAll(resp.Body)
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
			if err := meshkitutils.DownloadFile(filepath, manifestFile); err != nil {
				return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s", filepath, manifestFile)))
			}
		}
	}
	return nil
}

// DownloadOperatorManifest downloads the operator manifest files
func DownloadOperatorManifest() error {
	operatorFilepath := filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperator)
	err := meshkitutils.DownloadFile(operatorFilepath, OperatorURL)
	if err != nil {
		return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s operator file", operatorFilepath, MesheryOperator)))
	}

	brokerFilepath := filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperatorBroker)
	err = meshkitutils.DownloadFile(brokerFilepath, BrokerURL)
	if err != nil {
		return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s operator file", brokerFilepath, MesheryOperatorBroker)))
	}

	meshsyncFilepath := filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperatorMeshsync)
	err = meshkitutils.DownloadFile(meshsyncFilepath, MeshsyncURL)
	if err != nil {
		return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to download %s file from %s operator file", meshsyncFilepath, MesheryOperatorMeshsync)))
	}

	return nil
}

// returns the Channel and Version given a context
func GetChannelAndVersion(currCtx *(config.Context)) (string, string, error) {
	var version, channel string
	var err error

	version = currCtx.GetVersion()
	channel = currCtx.GetChannel()
	if version == "latest" {
		if channel == "edge" {
			version = "master"
		} else {
			version, err = GetLatestStableReleaseTag()
			if err != nil {
				return "", "", err
			}
		}
	}

	return channel, version, nil
}

func GetDeploymentVersion(filePath string) (string, error) {
	// setting up config type to yaml files
	ViperCompose.SetConfigType("yaml")

	// setting up config file
	ViperCompose.SetConfigFile(filePath)
	err := ViperCompose.ReadInConfig()
	if err != nil {
		return "", fmt.Errorf("unable to read config %s | %s", MesheryDeployment, err)
	}

	compose := K8sCompose{}
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// unmarshal the file into structs
	err = yaml.Unmarshal(yamlFile, &compose)
	if err != nil {
		return "", fmt.Errorf("unable to unmarshal config %s | %s", MesheryDeployment, err)
	}

	image := compose.Spec.Template.Spec.Containers[0].Image
	spliter := strings.Split(image, ":")
	version := strings.Split(spliter[1], "-")[1]

	return version, nil
}

// CanUseCachedOperatorManifests returns an error if it is not possible to use cached operator manifests
func CanUseCachedOperatorManifests(currCtx *(config.Context)) error {
	if _, err := os.Stat(filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperator)); os.IsNotExist(err) {
		return errors.New("operator manifest file does not exist")
	}

	if _, err := os.Stat(filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperatorBroker)); os.IsNotExist(err) {
		return errors.New("broker manifest file does not exist")
	}

	if _, err := os.Stat(filepath.Join(MesheryFolder, ManifestsFolder, MesheryOperatorMeshsync)); os.IsNotExist(err) {
		return errors.New("meshsync manifest file does not exist")
	}

	return nil
}

// CanUseCachedManifests returns an error if it is not possible to use cached manifests
func CanUseCachedManifests(currCtx *(config.Context)) error {
	// checks if meshery folder are present
	if _, err := os.Stat(MesheryFolder); os.IsNotExist(err) {
		return errors.New("Manifests folder does not exist")
	}

	// check if meshery deployment file is present
	deploymentsPath := filepath.Join(MesheryFolder, ManifestsFolder, MesheryDeployment)
	if _, err := os.Stat(deploymentsPath); os.IsNotExist(err) {
		return errors.New("Deployments file does not exist")
	}

	// compare versions in currCtx and meshery-deployment.yaml
	deploymentVersion, err := GetDeploymentVersion(deploymentsPath)
	if err != nil {
		return errors.Wrap(err, "could not get deployment file version")
	}
	var currVersion string
	if currCtx.GetVersion() != "latest" {
		currVersion = currCtx.GetVersion()
		if currVersion != deploymentVersion {
			return errors.New("deployment version mismatch")
		}
	}

	switch currCtx.GetPlatform() {
	case "kubernetes":
		// check if component manifests are present
		for _, component := range currCtx.GetComponents() {
			serviceFile := filepath.Join(MesheryFolder, ManifestsFolder, component+"-service.yaml")
			if _, err := os.Stat(serviceFile); os.IsNotExist(err) {
				return errors.New("service file does not exist")
			}

			componentDeploymentFile := filepath.Join(MesheryFolder, ManifestsFolder, component+"-deployment.yaml")
			if _, err := os.Stat(componentDeploymentFile); os.IsNotExist(err) {
				return errors.New("component deployment file does not exist")
			}
		}
	}

	return nil
}

// FetchManifests is a wrapper function that identifies the required manifest files as downloads them
func FetchManifests(currCtx *(config.Context)) ([]Manifest, error) {
	_, version, err := GetChannelAndVersion(currCtx)
	if err != nil {
		return []Manifest{}, err
	}

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

	err = CreateManifestsFolder()

	if err != nil {
		return nil, err
	}

	gitHubFolder := "https://github.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/tree/" + version + "/install/deployment_yamls/k8s"
	log.Info("downloading manifest files from ", gitHubFolder)

	// download all the manifest files to the ~/.meshery/manifests folder
	rawManifestsURL := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/" + version + "/install/deployment_yamls/k8s/"
	err = DownloadManifests(manifests, rawManifestsURL)

	if err != nil {
		return nil, errors.Wrap(err, "failed to download manifests")
	}

	return manifests, nil
}

// GetLatestStableReleaseTag fetches and returns the latest release tag from GitHub
func GetLatestStableReleaseTag() (string, error) {
	url := "https://github.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to get latest stable release tag")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	re := regexp.MustCompile("/releases/tag/(.*?)\"")
	releases := re.FindAllString(string(body), -1)
	latest := strings.ReplaceAll(releases[0], "/releases/tag/", "")
	latest = strings.ReplaceAll(latest, "\"", "")
	return latest, nil
}

// IsAdapterValid checks if the component mentioned by the user is a valid component
func IsAdapterValid(manifestArr []Manifest, componentManifest string) bool {
	for _, v := range manifestArr {
		if v.Path == componentManifest {
			return true
		}
	}

	return false
}

// DownloadDockerComposeFile fetches docker-compose.yaml based on passed context if it does not exists.
// Use force to override download anyway
func DownloadDockerComposeFile(ctx *config.Context, force bool) error {
	if _, err := os.Stat(DockerComposeFile); os.IsNotExist(err) || force {
		fileURL := ""

		if ctx.Channel == "edge" {
			fileURL = "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/master/install/docker/docker-compose.yaml"
		} else if ctx.Channel == "stable" {
			if ctx.Version == "latest" {
				ReleaseTag, err = GetLatestStableReleaseTag()
				if err != nil {
					return errors.Wrapf(err, "failed to fetch latest stable release tag")
				}
			} else { // else we get version tag from the config file
				mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
				if err != nil {
					return errors.Wrap(err, "error processing meshconfig")
				}

				currCtx, err := mctlCfg.GetCurrentContext()
				if err != nil {
					return err
				}
				ReleaseTag = currCtx.GetVersion()
			}

			fileURL = "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/" + ReleaseTag + "/install/docker/docker-compose.yaml"
		} else {
			return errors.Errorf("unknown channel %s", ctx.Channel)
		}

		if err := meshkitutils.DownloadFile(DockerComposeFile, fileURL); err != nil {
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
	// other than the components, meshery-deployment.yaml, meshery-service.yaml and service-account.yaml should be applied
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

	// Transform Manifests for custom configurations
	MesheryDeploymentManifestByt, err := TransformYAML([]byte(MesheryDeploymentManifest), func(i interface{}) (interface{}, error) {
		envVarI, ok := i.([]interface{})
		if !ok {
			return i, fmt.Errorf("unexpected data type")
		}

		return append(envVarI, map[string]interface{}{
			"name":  "MESHERY_SERVER_CALLBACK_URL",
			"value": viper.GetString("MESHERY_SERVER_CALLBACK_URL"),
		}), nil
	}, "spec", "template", "spec", "containers", "0", "env")
	if err != nil {
		log.Error(err)
		return errors.Wrap(err, "failed to transform manifest")
	}

	// apply/update/delete the manifest files
	if err = ApplyManifest(MesheryDeploymentManifestByt, client, update, delete); err != nil {
		return err
	}
	if err = ApplyManifest([]byte(mesheryServiceManifest), client, update, delete); err != nil {
		return err
	}
	if err = ApplyManifest([]byte(serviceAccountManifest), client, update, delete); err != nil {
		return err
	}

	// loop through the required components as specified in the config.yaml file and apply/update/delete each
	for _, component := range requestedAdapters {
		// for each component, there is a meshery-componentName-deployment.yaml and meshery-componentName-service.yaml
		// manifest file. See- https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
		componentFile := filepath.Join(manifestFiles, component)
		componentDeployment := componentFile + "-deployment.yaml"
		componentService := componentFile + "-service.yaml"

		if !IsAdapterValid(manifestArr, component+"-deployment.yaml") {
			return fmt.Errorf("invalid component %s specified. Please check %s/config.yaml file", component, MesheryFolder)
		}

		// read manifest files as strings and apply/update/delete
		manifestDepl, err := meshkitutils.ReadLocalFile(componentDeployment)
		if err != nil {
			return errors.Wrap(err, "failed to read manifest files")
		}
		manifestService, err := meshkitutils.ReadLocalFile(componentService)
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

// ApplyOperatorManifest applies/updates/deletes the operator manifest
func ApplyOperatorManifest(client *meshkitkube.Client, update bool, delete bool) error {
	// path to the manifest files ~/.meshery/manifests
	manifestFiles := filepath.Join(MesheryFolder, ManifestsFolder)

	//applying meshery operator file
	MesheryOperatorManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, MesheryOperator))

	if err != nil {
		return errors.Wrap(err, "failed to read operator manifest files")
	}

	if err = ApplyManifest([]byte(MesheryOperatorManifest), client, update, delete); err != nil {
		return err
	}

	//condition to check for system stop
	if !delete {
		MesheryBrokerManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, MesheryOperatorBroker))

		if err != nil {
			return errors.Wrap(err, "failed to read broker manifest files")
		}

		if err = ApplyManifest([]byte(MesheryBrokerManifest), client, update, delete); err != nil {
			return err
		}

		MesheryMeshsyncManifest, err := meshkitutils.ReadLocalFile(filepath.Join(manifestFiles, MesheryOperatorMeshsync))

		if err != nil {
			return errors.Wrap(err, "failed to read meshsync manifest files")
		}

		if err = ApplyManifest([]byte(MesheryMeshsyncManifest), client, update, delete); err != nil {
			return err
		}
	}

	log.Debug("applied operator manifest.")

	return nil
}

// ChangeManifestVersion changes the tag of the images in the manifest according to the pinned version
func ChangeManifestVersion(channel, version, filePath string) error {
	// setting up config type to yaml files
	ViperCompose.SetConfigType("yaml")

	// setting up config file
	ViperCompose.SetConfigFile(filePath)
	err := ViperCompose.ReadInConfig()
	if err != nil {
		return fmt.Errorf("unable to read config %s | %s", filePath, err)
	}

	compose := K8sCompose{}
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// unmarshal the file into structs
	err = yaml.Unmarshal(yamlFile, &compose)
	if err != nil {
		return fmt.Errorf("unable to unmarshal config %s | %s", filePath, err)
	}

	// for edge channel only the latest tag exist in Docker Hub
	if channel == "edge" {
		version = "latest"
	}

	image := compose.Spec.Template.Spec.Containers[0].Image
	spliter := strings.Split(image, ":")
	compose.Spec.Template.Spec.Containers[0].Image = fmt.Sprintf("%s:%s-%s", spliter[0], channel, version)

	log.Debug(image, " changed to ", compose.Spec.Template.Spec.Containers[0].Image)

	ViperCompose.Set("apiVersion", compose.APIVersion)
	ViperCompose.Set("kind", compose.Kind)
	ViperCompose.Set("metadata", compose.Metadata)
	ViperCompose.Set("spec", compose.Spec)
	ViperCompose.Set("status", compose.Status)

	// Marshal the structs
	newConfig, err := yaml.Marshal(compose)
	if err != nil {
		return fmt.Errorf("unable to marshal config %s | %s", filePath, err)
	}
	err = os.WriteFile(filePath, newConfig, 0644)
	if err != nil {
		return fmt.Errorf("unable to update config %s | %s", filePath, err)
	}

	return nil
}

// CreateManifestsFolder creates a new folder (.meshery/manifests)
func CreateManifestsFolder() error {
	log.Debug("deleting " + ManifestsFolder + " folder...")
	// delete manifests folder if it already exists
	if err := os.RemoveAll(ManifestsFolder); err != nil {
		return err
	}
	log.Debug("creating " + ManifestsFolder + "folder...")
	// create a manifests folder under ~/.meshery to store the manifest files
	if err := os.MkdirAll(filepath.Join(MesheryFolder, ManifestsFolder), os.ModePerm); err != nil {
		return errors.Wrapf(err, SystemError(fmt.Sprintf("failed to make %s directory", ManifestsFolder)))
	}
	log.Debug("created manifests folder...")

	return nil
}

// GetPods lists all the available pods in the MesheryNamespace
func GetPodList(client *meshkitkube.Client, namespace string) (*v1core.PodList, error) {
	// Create a pod interface for the given namespace
	podInterface := client.KubeClient.CoreV1().Pods(namespace)

	// List the pods in the given namespace
	podList, err := podInterface.List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	return podList, nil
}

// GetRequiredPods checks if the pods specified by the user is valid returns a list of the required pods
func GetRequiredPods(specifiedPods []string, availablePods []v1core.Pod) ([]string, error) {
	var requiredPods []string
	var availablePodsName []string
	for _, pod := range availablePods {
		availablePodsName = append(availablePodsName, pod.GetName())
	}
	for _, sp := range specifiedPods {
		if index := StringContainedInSlice(sp, availablePodsName); index != -1 {
			requiredPods = append(requiredPods, availablePodsName[index])
		} else {
			return nil, fmt.Errorf("invalid pod \"%s\" specified. Run mesheryctl `system status` to view the available pods", sp)
		}
	}
	return requiredPods, nil
}

// GetCleanPodName cleans the pod names in the MesheryNamespace to make it more readable
func GetCleanPodName(name string) string {
	// The pod names are of the form meshery-<component name>-dasd67qwe-jka244asd where the last characters are generated by kubernetes
	// Only the first two splits contain useful information
	split := strings.Split(name, "-")[:2]

	// Checks if the string in the split contains any useful info
	var IsRequired = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

	if IsRequired(split[1]) {
		return strings.Join(split, "-")
	}

	return split[0]
}

func Startdockerdaemon(subcommand string) error {
	userResponse := false
	// read user input on whether to start Docker daemon or not.
	if SilentFlag {
		userResponse = true
	} else {
		userResponse = AskForConfirmation("Start Docker now")
	}
	if !userResponse {
		return errors.Errorf("Please start Docker, then run the command `mesheryctl system %s`", subcommand)
	}

	log.Info("Attempting to start Docker...")
	// once user gaves permission, start docker daemon on linux/macOS
	if runtime.GOOS == "linux" {
		if err := exec.Command("sudo", "service", "docker", "start").Run(); err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
	} else {
		// Assuming we are on macOS, try to start Docker from default path
		cmd := exec.Command("/Applications/Docker.app/Contents/MacOS/Docker")
		err := cmd.Start()
		if err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
		// wait for few seconds for docker to start
		err = exec.Command("sleep", "30").Run()
		if err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
		// check whether docker started successfully or not, throw an error message otherwise
		if err := exec.Command("docker", "ps").Run(); err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
	}
	log.Info("Prerequisite Docker started.")
	return nil
}

func InstallprereqDocker() error {
	log.Info("Attempting Docker-Compose installation...")
	ostype, osarch, err := prereq()
	if err != nil {
		return errors.Wrap(err, "failed to get prerequisites")
	}

	osdetails := strings.TrimRight(string(ostype), "\r\n") + "-" + strings.TrimRight(string(osarch), "\r\n")

	dockerComposeBinaryURL := dockerComposeBinaryURL
	//checks for the latest docker-compose
	resp, err := http.Get(dockerComposeWebURL)
	if err != nil {
		dockerComposeBinaryURL = dockerComposeBinaryURL + defaultDockerComposeVersion
	} else {
		var dat map[string]interface{}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}
		if err := json.Unmarshal(body, &dat); err != nil {
			return errors.Wrap(err, "failed to unmarshal json into object")
		}
		num := dat["tag_name"]
		dockerComposeBinaryURL = fmt.Sprintf(dockerComposeBinaryURL+"%v/docker-compose", num)
	}
	dockerComposeBinaryURL = dockerComposeBinaryURL + "-" + osdetails
	if err := meshkitutils.DownloadFile(dockerComposeBinary, dockerComposeBinaryURL); err != nil {
		return errors.Wrapf(err, "failed to download %s from %s", dockerComposeBinary, dockerComposeBinaryURL)
	}
	if err := exec.Command("chmod", "+x", dockerComposeBinary).Run(); err != nil {
		return errors.Wrap(err, "failed to execute command")
	}
	log.Info("Prerequisite Docker Compose is installed.")
	return nil
}

// Sets the path to user's kubeconfig file into global variables
func SetKubeConfig() {
	// Define the path where the kubeconfig.yaml will be written to
	usr, err := user.Current()
	if err != nil {
		ConfigPath = filepath.Join(".meshery", KubeConfigYaml)
		KubeConfig = filepath.Join(".kube", "config")
	} else {
		ConfigPath = filepath.Join(usr.HomeDir, ".meshery", KubeConfigYaml)
		KubeConfig = filepath.Join(usr.HomeDir, ".kube", "config")
	}
}

func ForceCleanupCluster() error {
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return errors.Wrap(err, "failed to create new client")
	}

	// deletePolicy to delete related resources
	deletePolicy := metav1.DeletePropagationForeground

	// Deleting resources via kube-client
	// Deployments
	// Services
	// StatefulSets
	// ClusterRoleBindings
	// ClusterRoles
	// ServiceAccount
	// Secrets
	// ConfigMaps
	// Roles
	// RoleBindings

	deploymentInterface := client.KubeClient.AppsV1().Deployments(MesheryNamespace)
	deploymentList, err := deploymentInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, deployment := range deploymentList.Items {
		if err := deploymentInterface.Delete(context.TODO(), deployment.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	serviceInterface := client.KubeClient.CoreV1().Services(MesheryNamespace)
	serviceList, err := serviceInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, service := range serviceList.Items {
		if err := serviceInterface.Delete(context.TODO(), service.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	statefulSetInterface := client.KubeClient.AppsV1().StatefulSets(MesheryNamespace)
	statefulSetList, err := statefulSetInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, statefulSet := range statefulSetList.Items {
		if err := statefulSetInterface.Delete(context.TODO(), statefulSet.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	clusterRoleBindingInterface := client.KubeClient.RbacV1().ClusterRoleBindings()
	clusterRoleBindingList, err := clusterRoleBindingInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindingList.Items {
		// clusterrolebindings aren't namespaced so gotta
		// check if "meshery" is present in the name or not
		if strings.Contains(string(clusterRoleBinding.GetName()), "meshery") {
			if err := clusterRoleBindingInterface.Delete(context.TODO(), clusterRoleBinding.GetName(), metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Debug(err)
			}
		}
	}

	clusterRoleInterface := client.KubeClient.RbacV1().ClusterRoles()
	clusterRoleList, err := clusterRoleInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRole := range clusterRoleList.Items {
		// clusterroles aren't namespaced so gotta
		// check if "meshery" is present in the name or not
		if strings.Contains(string(clusterRole.GetName()), "meshery") {
			if err := clusterRoleInterface.Delete(context.TODO(), clusterRole.GetName(), metav1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				log.Debug(err)
			}
		}
	}

	serviceAccountInterface := client.KubeClient.CoreV1().ServiceAccounts(MesheryNamespace)
	serviceAccountList, err := serviceAccountInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, serviceAccount := range serviceAccountList.Items {
		if err := serviceAccountInterface.Delete(context.TODO(), serviceAccount.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	secretsInterface := client.KubeClient.CoreV1().Secrets(MesheryNamespace)
	secretsList, err := secretsInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, secrets := range secretsList.Items {
		if err := secretsInterface.Delete(context.TODO(), secrets.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	configMapsInterface := client.KubeClient.CoreV1().ConfigMaps(MesheryNamespace)
	configMapsList, err := configMapsInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, configMap := range configMapsList.Items {
		if err := configMapsInterface.Delete(context.TODO(), configMap.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	rolesInterface := client.KubeClient.RbacV1().Roles(MesheryNamespace)
	rolesList, err := rolesInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, role := range rolesList.Items {
		if err := rolesInterface.Delete(context.TODO(), role.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	roleBindingsInterface := client.KubeClient.RbacV1().RoleBindings(MesheryNamespace)
	roleBindingsList, err := roleBindingsInterface.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, roleBinding := range roleBindingsList.Items {
		if err := roleBindingsInterface.Delete(context.TODO(), roleBinding.GetName(), metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			log.Debug(err)
		}
	}

	return nil
}
