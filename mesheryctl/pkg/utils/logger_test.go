package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupMeshkitLogger(t *testing.T) {
	t.Run("Given_SetupMeshkitLogger_is_called_When_debug_is_enabled_Then_it_should_return_a_valid_handler", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger-debug"
		
		handler := SetupMeshkitLogger(name, true, &buf)
		
		assert.NotNil(t, handler, "Handler should not be nil")
	})

	t.Run("Given_SetupMeshkitLogger_is_called_When_debug_is_disabled_Then_it_should_return_a_valid_handler", func(t *testing.T) {
		var buf bytes.Buffer
		name := "test-logger-no-debug"
		
		handler := SetupMeshkitLogger(name, false, &buf)
		
		assert.NotNil(t, handler, "Handler should not be nil")
	})
}