package utils

import (
	"fmt"
	"testing"

	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorDescriptionsSeparateCause(t *testing.T) {
	cause := fmt.Errorf("connection refused")

	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrLoadConfig", ErrLoadConfig(cause), "Error processing config: connection refused"},
		{"ErrParseGithubFile", ErrParseGithubFile(cause, "https://github.com/meshery/meshery"), "Failed to parse github file: connection refused"},
		{"ErrRequestResponse", ErrRequestResponse(cause), "Unable to create a response from request: connection refused"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, errors.GetLDescription(tt.err))
		})
	}
}
