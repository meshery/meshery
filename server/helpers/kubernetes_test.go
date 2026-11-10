package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func TestFetchKubernetesVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tests")
	}
	logLevel := viper.GetInt("LOG_LEVEL")
	if viper.GetBool("DEBUG") {
		logLevel = int(logrus.DebugLevel)
	}
	// Initialize Logger instance
	log, err := logger.New("test", logger.Options{
		Format:   logger.SyslogLogFormat,
		LogLevel: logLevel,
	})
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// SetupContextEnv uses a relative path that assumes mesheryctl/ depth.
	// Override with the correct path for server/helpers/.
	path, err := os.Getwd()
	if err != nil {
		t.Fatal("unable to locate working directory")
	}
	configPath := filepath.Join(path, "..", "..", "mesheryctl", "pkg", "utils", "TestConfig.yaml")
	viper.Reset()
	viper.SetConfigFile(configPath)
	utils.DefaultConfigPath = configPath
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("unable to read configuration from %v: %v", configPath, err)
	}
	utils.StartMockery(t)

	// FetchKubernetesVersion requires a real K8s cluster.
	// Skip when no kubeconfig is available.
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("skipping: unable to determine home directory")
	}
	kubeconfig, err := os.ReadFile(filepath.Join(home, ".kube", "config"))
	if err != nil {
		t.Skip("skipping: no kubeconfig available at ~/.kube/config")
	}

	_, err = FetchKubernetesVersion(kubeconfig, "", log)
	if err != nil {
		t.Errorf("FetchKubernetesVersion() failed: %v", err)
	}
}
