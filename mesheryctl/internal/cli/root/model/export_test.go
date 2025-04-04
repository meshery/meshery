package model

import (
	"os"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/model/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestModelExportHelpOutput(t *testing.T) {
	modelCmd := &cobra.Command{Use: "model"}
	modelCmd.AddCommand(ExportCmd())

	modelCmd.SetArgs([]string{"export", "--help"})

	output, err := utils.ExecuteCommand(modelCmd)
	assert.NoError(t, err)

	expected, err := os.ReadFile("testdata/export.model.help.golden")
	assert.NoError(t, err)

	assert.Equal(t, string(expected), output)
}
