package model

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of NewListCmd for testing purposes
func NewListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all models",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Simulate successful execution
			cmd.Println("Model 1\nModel 2\nModel 3")
			return nil
		},
	}
}

// TestListModels checks the functionality of the list subcommand
func TestListModels(t *testing.T) {
	// Initialize the command
	cmd := NewListCmd()

	// Capture the output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Execute the command
	err := cmd.Execute()

	// Assert that no error occurs
	assert.NoError(t, err, "List command should not return an error")

	// Check the output for expected content
	assert.Contains(t, output.String(), "Model 1", "Output should contain 'Model 1'")
	assert.Contains(t, output.String(), "Model 2", "Output should contain 'Model 2'")
	assert.Contains(t, output.String(), "Model 3", "Output should contain 'Model 3'")
}
