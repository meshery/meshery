package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"k8s.io/apimachinery/pkg/version"
)

var minAPIVersion = [3]int{1, 12, 0}
var minKubectlVersion = minAPIVersion

//GetK8sVersionInfo returns version.Info for the Kubernetes cluster.
func GetK8sVersionInfo() (*version.Info, error) {
	// create an kubernetes client
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return nil, err
	}
	return client.KubeClient.Discovery().ServerVersion()
}

func CheckK8sVersion(versionInfo *version.Info) error {
	apiVersion, err := getK8sVersion(versionInfo.String())
	if err != nil {
		return err
	}

	if !isCompatibleVersion(minAPIVersion, apiVersion) {
		return fmt.Errorf("kubernetes is on version [%d.%d.%d], but version [%d.%d.%d] or more recent is required",
			apiVersion[0], apiVersion[1], apiVersion[2],
			minAPIVersion[0], minAPIVersion[1], minAPIVersion[2])
	}

	return nil
}

func getK8sVersion(versionString string) ([3]int, error) {
	var version [3]int
	var revisionSeparator = regexp.MustCompile("[^0-9.]")

	justTheVersionString := strings.TrimPrefix(versionString, "v")
	justTheMajorMinorRevisionNumbers := revisionSeparator.Split(justTheVersionString, -1)[0]
	split := strings.Split(justTheMajorMinorRevisionNumbers, ".")

	if len(split) < 3 {
		return version, fmt.Errorf("unknown version string format [%s]", versionString)
	}

	for i, segment := range split {
		v, err := strconv.Atoi(strings.TrimSpace(segment))
		if err != nil {
			return version, fmt.Errorf("unknown version string format [%s]", versionString)
		}
		version[i] = v
	}

	return version, nil
}

// CheckKubectlVersion validates whether the installed kubectl version is
// running a minimum kubectl version.
func CheckKubectlVersion() error {
	cmd := exec.Command("kubectl", "version", "--client", "--short")
	bytes, err := cmd.Output()
	if err != nil {
		return err
	}

	clientVersion := fmt.Sprintf("%s\n", bytes)
	kubectlVersion, err := parseKubectlShortVersion(clientVersion)
	if err != nil {
		return err
	}

	if !isCompatibleVersion(minKubectlVersion, kubectlVersion) {
		return fmt.Errorf("kubectl is on version [%d.%d.%d], but version [%d.%d.%d] or more recent is required",
			kubectlVersion[0], kubectlVersion[1], kubectlVersion[2],
			minKubectlVersion[0], minKubectlVersion[1], minKubectlVersion[2])
	}

	return nil
}

func isCompatibleVersion(minimalRequirementVersion [3]int, actualVersion [3]int) bool {
	if minimalRequirementVersion[0] < actualVersion[0] {
		return true
	}

	if (minimalRequirementVersion[0] == actualVersion[0]) && minimalRequirementVersion[1] < actualVersion[1] {
		return true
	}

	if (minimalRequirementVersion[0] == actualVersion[0]) && (minimalRequirementVersion[1] == actualVersion[1]) && (minimalRequirementVersion[2] <= actualVersion[2]) {
		return true
	}

	return false
}

var semVer = regexp.MustCompile("(v[0-9]+.[0-9]+.[0-9]+)")

func parseKubectlShortVersion(version string) ([3]int, error) {
	versionString := semVer.FindString(version)
	return getK8sVersion(versionString)
}
