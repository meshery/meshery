package environments

import (
	"testing"
	// "github.com/layer5io/meshkit/utils/kubernetes/describe"
	"github.com/spf13/cobra"
	// "github.com/stretchr/testify/assert"
	// "github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestCreateEnvironmentCmdValidation(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("orgId", "", "")
	cmd.Flags().String("name", "", "")
	cmd.Flags().String("description", "", "")

	tests := []struct {
		name          string
		args          []string
		expectedError bool
	}{
		{
			name:          "All flags present",
			args:          []string{"--orgId", "org123", "--name", "TestEnv", "--description", "Test Description"},
			expectedError: false,
		},
		{
			name:          "Missing orgId",
			args:          []string{"--name", "TestEnv", "--description", "Test Description"},
			expectedError: true,
		},
		{
			name:          "Missing name",
			args:          []string{"--orgId", "org123", "--description", "Test Description"},
			expectedError: true,
		},
		{
			name:          "Missing description",
			args:          []string{"--orgId", "org123", "--name", "TestEnv"},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd.Flags().Set("orgId", "")
			cmd.Flags().Set("name", "")
			cmd.Flags().Set("description", "")

			cmd.SilenceUsage = true
			cmd.SilenceErrors = true
			// cmd.SetArgs(tt.args)
			cmd.ParseFlags(tt.args)

			err := createEnvironmentCmd.Args(cmd, []string{})
			if (err != nil) != tt.expectedError {
				t.Errorf("createEnvironmentCmd.Args() error = %v, wantErr %v", err, tt.expectedError)
			}

		})
	}
}
