package utils

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestSetupMeshkitLogger ensures that the Meshkit logger handler
// is correctly initialized and messages are correctly filtered based on level.
func TestSetupMeshkitLogger(t *testing.T) {
	viper.Reset()

	t.Run("Given SetupMeshkitLogger is called with debug enabled, then it should capture debug logs", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger-debug"
		
		handler := SetupMeshkitLogger(name, true, &buf)
		assert.NotNil(t, handler)

		handler.Info("info message")
		handler.Debug("debug message")

		output := buf.String()
		assert.Contains(t, output, "info message")
		assert.Contains(t, output, "debug message")
	})

	t.Run("Given SetupMeshkitLogger is called with debug disabled, then it should hide debug logs", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger-no-debug"
		
		viper.Set("LOG_LEVEL", int(log.InfoLevel))

		handler := SetupMeshkitLogger(name, false, &buf)
		assert.NotNil(t, handler)

		handler.Info("info message")
		handler.Debug("hidden debug message") 

		output := buf.String()
		assert.Contains(t, output, "info message")
		assert.NotContains(t, output, "hidden debug message")
	})
}