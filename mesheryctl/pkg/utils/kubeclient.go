package utils

import (
	"fmt"
	"os"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
)

// NewKubeClient returns a Kubernetes client for mesheryctl operations.
//
// Resolution order (Meshery-managed cluster access is authoritative when present):
//  1. If ConfigPath (~/.meshery/kubeconfig.yaml) exists and is non-empty, use it.
//     If the file exists but cannot be loaded into a client, return an error —
//     never silent-fallthrough to ambient kubeconfig (that deploys to the wrong cluster).
//  2. If ConfigPath is missing, use ambient resolution via meshkit
//     (KUBECONFIG env, then ~/.kube/config).
//
// This is the single entry point for mesheryctl Kubernetes client construction.
// Lifecycle commands (start/stop/status/…) and helpers must call NewKubeClient
// instead of meshkitkube.New([]byte("")) so connection create / system config
// and system start agree on the target cluster.
func NewKubeClient() (*meshkitkube.Client, error) {
	if ConfigPath == "" {
		return meshkitkube.New([]byte(""))
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return meshkitkube.New([]byte(""))
		}
		return nil, fmt.Errorf("unable to read Meshery kubeconfig at %s: %w", ConfigPath, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("Meshery kubeconfig at %s is empty; re-run `mesheryctl connection create` or `mesheryctl system config`, or remove the file to use ambient kubeconfig", ConfigPath)
	}

	client, err := meshkitkube.New(data)
	if err != nil {
		return nil, fmt.Errorf("invalid Meshery kubeconfig at %s: %w; fix or remove the file, or re-run `mesheryctl connection create` / `mesheryctl system config`", ConfigPath, err)
	}

	if Log != nil {
		Log.Debugf("Using Meshery-managed kubeconfig: %s", ConfigPath)
	}
	return client, nil
}
