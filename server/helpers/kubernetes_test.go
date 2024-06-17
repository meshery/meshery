package helpers

import (
	"os"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshkit/logger"
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
	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testByte := []byte{0, 0}
	var testContext = "nil"

	_, err = FetchKubernetesVersion(testByte, testContext, log)
	if err != nil {
		t.Error("FetchKubernetesVersion() failed")
	}
}
