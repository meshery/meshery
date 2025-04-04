package utils

import (
	"bytes"

	"github.com/spf13/cobra"
)

func ExecuteCommand(cmd *cobra.Command) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	err := cmd.Execute()
	return buf.String(), err
}
