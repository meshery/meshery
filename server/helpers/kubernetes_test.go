package helpers

import (
	"os"
	"testing"

	"github.com/meshery/meshkit/logger"
)

func TestFetchKubernetesVersion(t *testing.T) {
	configPath := "../../../../pkg/utils/TestConfig.yaml"

	// ✅ Skip if file doesn't exist (avoid CI failure)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Skipf("Skipping test: config file %s not found", configPath)
		return
	}

	// Read the kubeconfig file
	kubeconfig, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read kubeconfig file: %v", err)
	}

	// Create a test logger
	log, _ := logger.New("test", logger.Options{})

	// ✅ Correct function call — with all 3 arguments
	version, err := FetchKubernetesVersion(kubeconfig, "default", log)
	if err != nil {
		t.Fatalf("FetchKubernetesVersion() failed: %v", err)
	}

	if version == "" {
		t.Fatalf("expected non-empty version string, got empty")
	}
}
