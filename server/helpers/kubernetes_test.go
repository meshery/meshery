package helpers

import (
	"os"
	"testing"

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
	// Provide an intentionally invalid kubeconfig to ensure the function returns
	// a wrapped, actionable error instead of requiring a live cluster.
	_, err = FetchKubernetesVersion([]byte{0, 0}, "", log)
	if err == nil {
		t.Error("expected FetchKubernetesVersion() to fail with invalid kubeconfig")
	}
}
