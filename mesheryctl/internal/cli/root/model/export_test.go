package model

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/spf13/cobra"
)

// TestExportModel checks the functionality of the export subcommand
func TestExportModel(t *testing.T) {
	// Initialize the command (you might need to modify this to fit your structure)
	cmd := NewExportCmd()

	// Set arguments for the command, e.g., model name
	cmd.SetArgs([]string{"sample-model"})

	// Execute the command
	err := cmd.Execute()

	// Assert that no error occurs
	assert.NoError(t, err, "Export command should not return an error")
}
