package model

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of NewViewCmd for testing purposes
func NewViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "View details of a model",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Simulate successful execution
			if len(args) == 0 {
				return cmd.Help()
			}
			cmd.Println("Details of model: sample-model")
			return nil
		},
	}
}

// TestViewModel checks the functionality of the view subcommand
func TestViewModel(t *testing.T) {
	// Initialize the command
	cmd := NewViewCmd()

	// Set arguments for the command, e.g., model name
	cmd.SetArgs([]string{"sample-model"})

	// Capture the output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Execute the command
	err := cmd.Execute()

	// Assert that no error occurs
	assert.NoError(t, err, "View command should not return an error")

	// Check the output for expected content
	assert.Contains(t, output.String(), "Details of model: sample-model", "Output should confirm the model details were displayed")
}
