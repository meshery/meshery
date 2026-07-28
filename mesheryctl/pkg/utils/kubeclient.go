package utils

import (
	"fmt"
	"os"

	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
)

// NewKubeClient returns a Kubernetes client for mesheryctl.
// Prefer ConfigPath (~/.meshery/kubeconfig.yaml) when present; otherwise use ambient kubeconfig via meshkit.
func NewKubeClient() (*meshkitkube.Client, error) {
	if ConfigPath == "" {
		return meshkitkube.New([]byte(""))
	}

	data, err := os.ReadFile(ConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return meshkitkube.New([]byte(""))
		}
		return nil, ErrFileRead(fmt.Errorf("unable to read Meshery kubeconfig at %s: %w", ConfigPath, err))
	}

	if len(data) == 0 {
		return nil, ErrInvalidFile(fmt.Errorf("Meshery kubeconfig at %s is empty; re-run `mesheryctl connection create`, or remove the file to use ambient kubeconfig", ConfigPath))
	}

	client, err := meshkitkube.New(data)
	if err != nil {
		return nil, ErrInvalidFile(fmt.Errorf("invalid Meshery kubeconfig at %s: %w", ConfigPath, err))
	}

	if Log != nil {
		Log.Debugf("Using Meshery-managed kubeconfig: %s", ConfigPath)
	}
	return client, nil
}
