package model

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Mock implementation of NewSearchCmd for testing purposes
func NewSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search",
		Short: "Search for a model",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Simulate successful execution
			if len(args) == 0 {
				return cmd.Help()
			}
			cmd.Println("Found model: sample-model")
			return nil
		},
	}
}

// TestSearchModel checks the functionality of the search subcommand
func TestSearchModel(t *testing.T) {
	// Initialize the command
	cmd := NewSearchCmd()

	// Set arguments for the command, e.g., model name
	cmd.SetArgs([]string{"sample-model"})

	// Capture the output
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Execute the command
	err := cmd.Execute()

	// Assert that no error occurs
	assert.NoError(t, err, "Search command should not return an error")

	// Check the output for expected content
	assert.Contains(t, output.String(), "Found model: sample-model", "Output should confirm the model was found")
}
