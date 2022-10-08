package utils

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/version"
)

var minAPIVersion = [3]int{1, 12, 0}
var minKubectlVersion = minAPIVersion

// GetK8sVersionInfo returns version.Info for the Kubernetes cluster.
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

// IsMesheryRunning checks if the meshery server containers are up and running
func IsMesheryRunning(currPlatform string) (bool, error) {
	// Get viper instance used for context to extract the endpoint from config file
	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())

	currCtx, _ := mctlCfg.GetCurrentContext()

	urlEndpoint := currCtx.GetEndpoint()

	urlTest := urlEndpoint + "/api/system/version"

	// Checking if Meshery is running with the URL obtained
	resp, _ := http.Get(urlTest)

	if resp != nil && resp.StatusCode == 200 {
		return true, nil
	}

	//If not, use the platforms to check if Meshery is running or not
	switch currPlatform {
	case "docker":
		{
			op, err := exec.Command("docker-compose", "-f", DockerComposeFile, "ps").Output()
			if err != nil {
				return false, errors.Wrap(err, " required dependency, docker-compose, is not present or docker is not available. Please run `mesheryctl system check --preflight` to verify system readiness")
			}
			return strings.Contains(string(op), "meshery"), nil
		}
	case "kubernetes":
		{
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return false, errors.Wrap(err, "failed to create new client")
			}

			//podInterface := client.KubeClient.CoreV1().Pods(MesheryNamespace)
			deploymentInterface := client.KubeClient.AppsV1().Deployments(MesheryNamespace)
			//podList, err := podInterface.List(context.TODO(), v1.ListOptions{})
			deploymentList, err := deploymentInterface.List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return false, err
			}
			for _, deployment := range deploymentList.Items {
				if deployment.GetName() == "meshery" {
					return true, nil
				}
			}

			return false, err
		}
	}

	return false, nil
}

// AreMesheryComponentsRunning checks if the meshery containers are up and running
func AreMesheryComponentsRunning(currPlatform string) (bool, error) {
	//If not, use the platforms to check if Meshery is running or not
	switch currPlatform {
	case "docker":
		{
			op, err := exec.Command("docker-compose", "-f", DockerComposeFile, "ps").Output()
			if err != nil {
				return false, errors.Wrap(err, " required dependency, docker-compose, is not present or docker is not available. Please run `mesheryctl system check --preflight` to verify system readiness")
			}
			return strings.Contains(string(op), "meshery"), nil
		}
	case "kubernetes":
		{
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return false, errors.Wrap(err, "failed to create new client")
			}

			deploymentInterface := client.KubeClient.AppsV1().Deployments(MesheryNamespace)
			deploymentList, err := deploymentInterface.List(context.TODO(), metav1.ListOptions{})

			if err != nil {
				return false, err
			}
			for _, deployment := range deploymentList.Items {
				if strings.Contains(string(deployment.GetName()), "meshery") {
					return true, nil
				}
			}

			return false, err
		}
	}

	return false, nil
}

// AreAllPodsRunning checks if all the deployment pods under kubernetes are running
func AreAllPodsRunning() (bool, error) {
	// create an kubernetes client
	client, err := meshkitkube.New([]byte(""))

	if err != nil {
		return false, err
	}

	// List the pods in the MesheryNamespace
	podList, err := GetPodList(client, MesheryNamespace)

	if err != nil {
		return false, err
	}

	// List all the pods similar to kubectl get pods -n MesheryNamespace
	for _, pod := range podList.Items {
		// Check the status of each of the pods
		if pod.Status.Phase != "Running" {
			return false, nil
		}
	}
	return true, nil
}

// CheckMesheryNsDelete waits for Meshery namespace to be deleted, returns (done, error)
func CheckMesheryNsDelete() (bool, error) {
	client, err := meshkitkube.New([]byte(""))
	if err != nil {
		return false, err
	}

	if err := WaitForNamespaceDeleted(client, MesheryNamespace, 300); err != nil {
		return false, err
	}

	return true, nil
}

// return a condition function that indicates whether the given pod is
// currently running
func isPodRunning(c *meshkitkube.Client, podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.KubeClient.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		switch pod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, fmt.Errorf("%s failed to start or never reached Running state", podName)
		}
		return false, nil
	}
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func pollForPodRunning(c *meshkitkube.Client, namespace, podName string, timeout time.Duration) error {
	return wait.PollImmediate(time.Second, timeout, isPodRunning(c, podName, namespace))
}

// Wait up to timeout seconds for pod in 'namespace' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func WaitForPodRunning(c *meshkitkube.Client, desiredPod, namespace string, timeout int) error {
	podList, err := GetPodList(c, namespace)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		return fmt.Errorf("no pods in %s", namespace)
	}
	var desiredPodName string
	for _, pod := range podList.Items {
		if GetCleanPodName(pod.Name) == desiredPod {
			desiredPodName = pod.Name
			break
		}
	}

	if desiredPodName == "" {
		return fmt.Errorf("`%s` pod not found", desiredPod)
	}

	return pollForPodRunning(c, namespace, desiredPodName, time.Duration(timeout)*time.Second)
}

// Returns condition function to indicate that the `namespace` does not exist anymore.
func isNamespaceDeleted(c *meshkitkube.Client, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		namespaces, err := c.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		// Check if namespace exists in namespaces list
		for _, ns := range namespaces.Items {
			if ns.Name == namespace {
				return false, nil
			}
		}

		return true, nil
	}
}

// Poll up to timeout seconds every 5 seconds until the namespace no more exists.
func pollForNamespaceDeleted(c *meshkitkube.Client, namespace string, timeout time.Duration) error {
	return wait.Poll(5*time.Second, timeout, isNamespaceDeleted(c, namespace))
}

// Wait up to timeout seconds for `namespace` to be deleted.
func WaitForNamespaceDeleted(c *meshkitkube.Client, namespace string, timeout int) error {
	return pollForNamespaceDeleted(c, namespace, time.Duration(timeout)*time.Second)
}
