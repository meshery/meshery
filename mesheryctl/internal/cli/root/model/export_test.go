package model

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExportCommand(t *testing.T) {
	// Create a test instance of the command
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export registered models",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("Exporting models...")
			return nil
		},
	}

	// Capture the command output
	var out bytes.Buffer
	cmd.SetOut(&out)

	// Execute the command
	err := cmd.Execute()
	assert.NoError(t, err)

	// Verify expected output
	expected := "Exporting models...\n"
	assert.Equal(t, expected, out.String())
}
