package utils

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper" // Added this
	"github.com/stretchr/testify/assert"
)

// TestTerminalFormatter_Format verifies that the custom logrus formatter
// correctly appends a newline character to log messages.
func TestTerminalFormatter_Format(t *testing.T) {
	t.Parallel()

	f := &TerminalFormatter{}

	cases := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "Simple message",
			message:  "Hello Meshery",
			expected: "Hello Meshery\n",
		},
		{
			name:     "Empty message",
			message:  "",
			expected: "\n",
		},
		{
			name:     "Message with special characters",
			message:  "Status: 200 OK! #uuid-123",
			expected: "Status: 200 OK! #uuid-123\n",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			entry := &log.Entry{Message: tc.message}
			got, err := f.Format(entry)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(got))
		})
	}
}

// TestSetupMeshkitLogger ensures that the Meshkit logger handler
// is correctly initialized and messages are correctly written to the output.
func TestSetupMeshkitLogger(t *testing.T) {
	t.Run("Initialize logger with debug enabled", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger"
		
		// Force Viper to a log level that allows Info/Debug
		viper.Set("LOG_LEVEL", int(log.DebugLevel)) 
		
		handler := SetupMeshkitLogger(name, true, &buf)
		assert.NotNil(t, handler)

		handler.Info("info message")

		output := buf.String()
		assert.Contains(t, output, "info message", "Buffer should contain the logged message")
	})

	t.Run("Initialize logger with debug disabled", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger"
		
		viper.Set("LOG_LEVEL", int(log.InfoLevel))

		handler := SetupMeshkitLogger(name, false, &buf)
		assert.NotNil(t, handler)

		handler.Info("info message")
		output := buf.String()
		assert.Contains(t, output, "info message")
	})
}