package perf

import (
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPerfCmd(t *testing.T) {
	utils.SetupContextEnv(t)

	// test scenrios for fetching data
	tests := []struct {
		Name        string
		Args        []string
		ExpectError bool
	}{
		{"Invalid command", []string{"invalid"}, true},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			PerfCmd.SetArgs(tt.Args)
			err := PerfCmd.Execute()
			if err != nil {
				if tt.ExpectError && strings.Contains(err.Error(), "invalid") {
					return
				}
				t.Error("Failed")
			}
			t.Error("Failed")
		})
		t.Log("Perf tests Passed")
	}
}

func TestPerfCmdSubcommands(t *testing.T) {
	// Test that all expected subcommands are available
	expectedSubcommands := []string{"apply", "profile", "result", "start", "stop"}
	
	for _, subCmd := range expectedSubcommands {
		found := false
		for _, cmd := range availableSubcommands {
			if cmd.Use == subCmd {
				found = true
				break
			}
		}
		assert.True(t, found, "Subcommand %s should be available", subCmd)
	}
}

func TestPerfCmdStructure(t *testing.T) {
	// Test command structure
	assert.Equal(t, "perf", PerfCmd.Use)
	assert.Equal(t, "Performance Management & Benchmarking", PerfCmd.Short)
	assert.NotEmpty(t, PerfCmd.Long)
	assert.NotEmpty(t, PerfCmd.Example)
	assert.NotNil(t, PerfCmd.RunE)
	
	// Test that examples include new commands
	assert.Contains(t, PerfCmd.Example, "mesheryctl perf start")
	assert.Contains(t, PerfCmd.Example, "mesheryctl perf stop")
}
