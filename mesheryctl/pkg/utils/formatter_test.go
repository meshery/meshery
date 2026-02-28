package utils

import (
	"bytes"
	"testing"

	log "github.com/sirupsen/logrus"
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
// is correctly initialized with provided options and output writer.
func TestSetupMeshkitLogger(t *testing.T) {
	// We use a buffer to capture output instead of printing to console during tests
	var buf bytes.Buffer
	
	t.Run("Initialize logger successfully", func(t *testing.T) {
		name := "test-logger"
		debugLevel := true
		
		handler := SetupMeshkitLogger(name, debugLevel, &buf)
		
		// Assert that the handler was created and is not nil
		assert.NotNil(t, handler, "Handler should not be nil")
	})
}