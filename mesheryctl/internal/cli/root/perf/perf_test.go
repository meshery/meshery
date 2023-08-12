package perf

import (
	"strings"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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
	}
}
