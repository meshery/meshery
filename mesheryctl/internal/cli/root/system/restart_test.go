package system

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestRestart(t *testing.T) {
	// skipping this integration test with --short flag
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// setup current context
	//utils.SetupContextEnv(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	// update all location
	utils.SetFileLocationTesting(t, currDir)
	utils.SetupCustomContextEnv(t, currDir+"/testdata/restart/TestRestart.yaml")

	tests := []struct {
		Name        string
		Action      string
		Args        []string
		ExpectError bool
	}{
		{
			Name:        "Restart Meshery ",
			Action:      "restart",
			Args:        []string{"restart"},
			ExpectError: false,
		},
		{
			Name:        "Restart Meshery and skiping update ",
			Action:      "restart",
			Args:        []string{"restart", "--skip-update"},
			ExpectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			SystemCmd.SetArgs(tt.Args)

			t.Logf("%sing meshery", tt.Action)
			err := SystemCmd.Execute()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("Meshery %sed", tt.Action)

		})
	}
}
