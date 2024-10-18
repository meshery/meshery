package utils

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetupMeshkitLogger(t *testing.T) {
	// Setup viper configuration
	viper.SetConfigType("yaml")
	err := viper.ReadConfig(bytes.NewBuffer([]byte("LOG_LEVEL: 4")))
	if err != nil {
		t.Fatalf("failed to read viper configuration: %v", err)
	}

	t.Run("DebugLevelFalse", func(t *testing.T) {
		var buf bytes.Buffer
		logHandler := SetupMeshkitLogger("testLogger", false, &buf)

		assert.NotNil(t, logHandler, "Logger should not be nil")
		assert.Equal(t, log.InfoLevel, log.GetLevel(), "Log level should be set from viper configuration")
	})

	t.Run("LoggerOutput", func(t *testing.T) {
		var buf bytes.Buffer
		logHandler := SetupMeshkitLogger("testLogger", true, &buf)

		assert.NotNil(t, logHandler, "Logger should not be nil")

		logHandler.Info("Test log message")
		assert.Contains(t, buf.String(), "Test log message", "Logger output should contain the test message")
	})
}
